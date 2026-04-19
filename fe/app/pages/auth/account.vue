<script setup lang="ts">
definePageMeta({
  path: '/auth/account/:account?',
  validate: route => {
    const account = Array.isArray(route.params.account) ? route.params.account[0] : route.params.account
    return !account || ['profile', 'balance', 'security', 'admin'].includes(String(account))
  }
})

const { readThemePreferenceFromSettings } = useInterfacePreferences()
const authI18n = useI18nSection('auth')
const paymentsI18n = useI18nSection('payments')
const uiI18n = useI18nSection('ui')
await Promise.all([
  useAsyncData(authI18n.key.value, authI18n.load, { watch: [authI18n.locale] }),
  useAsyncData(paymentsI18n.key.value, paymentsI18n.load, { watch: [paymentsI18n.locale] }),
  useAsyncData(uiI18n.key.value, uiI18n.load, { watch: [uiI18n.locale] })
])
const t = (messageKey: string, params?: Record<string, string | number>) => {
  if (messageKey.startsWith('payments.')) return paymentsI18n.t(messageKey.slice(9), params)
  if (messageKey.startsWith('settings.') || messageKey.startsWith('theme.')) return uiI18n.t(messageKey, params)
  return authI18n.t(messageKey.startsWith('auth.') ? messageKey.slice(5) : messageKey, params)
}
const title = computed(() => t('auth.account.title'))
const description = computed(() => t('auth.account.description'))
usePageSeo({
  title,
  description
})
const route = useRoute()
const accountTabPathMap: Record<'profile' | 'balance' | 'security' | 'admin', string> = {
  profile: '/auth/account/profile',
  balance: '/auth/account/balance',
  security: '/auth/account/security',
  admin: '/auth/admin/users'
}
const routeAccountParam = Array.isArray(route.params.account) ? route.params.account[0] : route.params.account
const routeAccountQuery = Array.isArray(route.query.account) ? route.query.account[0] : route.query.account
const routeAccountTarget = normalizeTabRouteValue(
  routeAccountParam || routeAccountQuery,
  ['profile', 'balance', 'security', 'admin'],
  'profile'
) as 'profile' | 'balance' | 'security' | 'admin'
if (routeAccountTarget === 'admin') {
  await navigateTo('/auth/admin/users', {
    redirectCode: 301,
    replace: true
  })
}
if (routeAccountQuery || routeAccountParam !== routeAccountTarget) {
  const { account: _account, ...query } = route.query
  await navigateTo(
    {
      path: accountTabPathMap[routeAccountTarget],
      query,
      hash: route.hash
    },
    {
      redirectCode: 301,
      replace: true
    }
  )
}
const router = useRouter()
const runtimeConfig = useRuntimeConfig()
const uiPreferences = useUiPreferencesStore()
const { formatAbsoluteDateTime } = useLocalizedDateTime()
const { saveInterfacePreferences: syncInterfacePreferences } = useInterfacePreferencesSync()
const {
  user,
  isAdmin,
  loaded,
  ensureLoaded,
  loadMe,
  logout,
  logoutAll,
  requestEmailChange,
  listSessions,
  listLoginAttempts,
  listSecurityEvents,
  revokeSession,
  updateProfile,
  changePassword,
  setupTwoFactor,
  enableTwoFactor,
  disableTwoFactor
} = useAuth()

const {
  access,
  accessLoading: _accessLoading,
  ensureAccessLoaded,
  history,
  historyLoading: _historyLoading,
  loadHistory,
  refundOrder,
  loadAccess
} = usePayments()

await ensureLoaded()
await ensureAccessLoaded()

const SUCCESS_NOTIFY_TIMEOUT_MS = 5000
const AVATAR_IMAGE_ACCEPT = 'image/jpeg,image/png,image/webp,image/gif,image/avif'
const avatarImageMaxBytes = Number(runtimeConfig.public.mediaImageMaxBytes || 8388608)
const AVATAR_IMAGE_MAX_MB = Math.max(1, Math.ceil(avatarImageMaxBytes / (1024 * 1024)))
const AVATAR_IMAGE_MAX_BYTES = avatarImageMaxBytes

const accountTabItems = computed<LabTabItem[]>(() => [
  { value: 'profile', label: t('auth.account.tab.profile') },
  { value: 'balance', label: t('auth.account.tab.balance') },
  { value: 'security', label: t('auth.account.tab.security') },
  ...(isAdmin.value ? [{ value: 'admin', label: t('auth.account.admin') }] : [])
])
const accountTabRouteTargetMap = computed<TabRouteTargetMap>(() => ({
  profile: '/auth/account/profile',
  balance: '/auth/account/balance',
  security: '/auth/account/security',
  admin: '/auth/admin/users'
}))
const accountTabValues = computed(() => accountTabItems.value.map(item => item.value))
const _activityTabItems = computed<LabTabItem[]>(() => [
  { value: 'sessions', label: t('auth.account.activity.sessions') },
  { value: 'attempts', label: t('auth.account.activity.attempts') },
  { value: 'events', label: t('auth.account.activity.events') }
])
const _disableCodeTabItems = computed<LabTabItem[]>(() => [
  { value: 'totp', label: t('auth.account.twofa.code_label') },
  { value: 'backup', label: t('auth.login.mfa_backup') }
])

const accountTab = computed<'profile' | 'balance' | 'security' | 'admin'>({
  get: () =>
    normalizeTabRouteValue(route.params.account || route.query.account, accountTabValues.value, 'profile') as
      | 'profile'
      | 'balance'
      | 'security'
      | 'admin',
  set: () => {}
})
const accountBreadcrumbItems = computed<BreadcrumbItem[]>(() => {
  const currentLabel = accountTabItems.value.find(item => item.value === accountTab.value)?.label || title.value
  return [
    { label: title.value, to: '/auth/account/profile' },
    { label: currentLabel, current: true, kind: 'tab' }
  ]
})
const _activityTab = ref<'sessions' | 'attempts' | 'events'>('sessions')
const disableMethodTab = ref<'totp' | 'backup'>('totp')

const initialLoading = ref(!loaded.value)
const editingDisplayName = ref(false)
const showLogoutActions = ref(false)
const showDisable2faForm = ref(false)
const avatarUploading = ref(false)
const refundPendingOrderId = ref('')
const displayNameKeyboardSubmitPending = ref(false)

const sessions = ref<AuthSessionView[]>([])
const loginAttempts = ref<AuthLoginAttemptView[]>([])
const securityEvents = ref<AuthSecurityEventView[]>([])
const backupCodes = ref<string[]>([])

const enableCode = ref('')
const disablePassword = ref('')
const disableTotpCode = ref('')
const disableBackupCode = ref('')
const setupSecret = ref('')
const setupOtpAuthUrl = ref('')
const setupQrDataUrl = ref('')

const activityLoading = ref(false)
const activityLoaded = ref(false)
const activityError = ref('')
const actionError = ref('')
const actionInfo = ref('')
const paymentHistoryError = ref('')
const paymentHistoryInfo = ref('')
const twofaError = ref('')
const twofaInfo = ref('')

const passwordError = ref('')
const passwordInfo = ref('')
const passwordCurrentInputError = ref(false)
const emailChangeError = ref('')
const emailChangeInfo = ref('')
const emailChangePending = ref(false)
const profileError = ref('')
const profileInfo = ref('')
const profileNoticeTone = ref<NotifyTone>('success')
const profileNoticeTemporary = ref(false)
const interfaceError = ref('')
const interfaceInfo = ref('')
const avatarError = ref('')
const avatarInfo = ref('')

const suppressInterfaceAutosave = ref(false)
const interfaceAutosavePending = ref(false)
const suppressProfileErrorReset = ref(false)
const profileInputError = ref(false)

const enableCodeInputRef = ref<{ focus: () => void } | null>(null)
const disableTotpInputRef = ref<{ focus: () => void } | null>(null)
const disableBackupCodeInputRef = ref<{ focus: () => void } | null>(null)
const displayNameInputRef = ref<HTMLInputElement | null>(null)
const profileNameEditorRef = ref<HTMLElement | null>(null)
const avatarPreviewScrollerRef = ref<HTMLElement | null>(null)

const successNoticeTimers = new Map<string, ReturnType<typeof setTimeout>>()

const passwordForm = reactive({
  current: '',
  next: ''
})
const emailChangeForm = reactive({
  email: ''
})
const profileForm = reactive({
  display_name: ''
})
const interfaceForm = reactive<{
  locale: InterfaceLocaleCode
  theme_preference: ThemePreference
}>({
  locale: 'ru',
  theme_preference: 'system'
})
const avatarCropDialog = reactive<{
  open: boolean
  file: File | null
}>({
  open: false,
  file: null
})
const avatarPreviewDialog = reactive({
  open: false,
  index: 0
})
const paymentAccessUntilText = computed(() => {
  if (!access.value?.access_until) return ''
  return formatAbsoluteDateTime(access.value.access_until)
})
const profileSubscriptionTooltipText = computed(() => {
  if (!access.value?.has_active_access || !paymentAccessUntilText.value) return ''
  return `Подписка активна до ${paymentAccessUntilText.value}`
})
const _balanceSubscriptionStatusLabel = computed(() =>
  access.value?.has_active_access ? `Подписка активна до ${paymentAccessUntilText.value || '—'}` : 'Подписка не активна'
)
const _paymentHistoryColumns = computed<LabDataTableColumn[]>(() => [
  { key: 'createdAt', label: 'Дата', nowrap: true },
  { key: 'plan', label: 'План' },
  { key: 'amount', label: 'Сумма', nowrap: true },
  { key: 'status', label: 'Статус', nowrap: true },
  { key: 'access', label: 'Доступ' },
  { key: 'refund', label: 'Возврат', nowrap: true }
])
const _paymentHistoryRows = computed(() =>
  history.value.map(item => ({
    id: item.order_id,
    orderId: item.order_id,
    createdAt: formatDateTime(item.created_at),
    plan: paymentPlanLabel(item.plan_code),
    amount: formatPaymentWholeRubles(item.amount),
    status: paymentStatusLabel(item.status),
    access: item.access_until
      ? `До ${formatDateTime(item.access_until)}`
      : item.access_from
        ? `С ${formatDateTime(item.access_from)}`
        : '—',
    refund: item.can_refund ? 'Доступен' : item.refunded_at ? 'Выполнен' : 'Недоступен',
    canRefund: item.can_refund
  }))
)
const _sessionActivityColumns = computed<LabDataTableColumn[]>(() => [
  { key: 'device', label: 'Устройство', cellClass: 'whitespace-normal wrap-break-word' },
  { key: 'status', label: '2FA', nowrap: true },
  { key: 'activity', label: 'Активность', cellClass: 'whitespace-normal wrap-break-word' },
  { key: 'action', label: 'Действие', nowrap: true }
])
const _sessionActivityRows = computed(() =>
  groupedSessions.value.map(item => ({
    id: item.key,
    device: item.deviceLabel,
    ip: item.ip,
    status: item.mfaVerified ? t('auth.account.activity.mfa_verified') : t('auth.account.activity.mfa_unverified'),
    count: item.count,
    lastSeenAt: item.lastSeenAt,
    revokableSessionIds: item.revokableSessionIds,
    hasCurrent: item.hasCurrent,
    action:
      item.revokableSessionIds.length === 0
        ? t('auth.account.activity.current')
        : item.hasCurrent
          ? t('auth.account.activity.revoke_and_logout')
          : t('auth.account.activity.revoke'),
    source: item
  }))
)
const _loginAttemptColumns = computed<LabDataTableColumn[]>(() => [
  { key: 'createdAt', label: 'Дата', nowrap: true },
  { key: 'outcome', label: 'Результат', nowrap: true },
  { key: 'ip', label: 'IP', nowrap: true },
  { key: 'risk', label: 'Риск', nowrap: true },
  { key: 'details', label: 'Детали', cellClass: 'whitespace-normal wrap-break-word' }
])
const _loginAttemptRows = computed(() =>
  loginAttempts.value.map(item => ({
    id: item.attempt_id,
    createdAt: formatDateTime(item.created_at),
    outcome: item.outcome || '—',
    ip: item.ip || '—',
    risk: String(item.risk_score ?? '—'),
    details: item.failure_reason || item.suspicious_reason || item.user_agent || '—',
    source: item
  }))
)
const _securityEventColumns = computed<LabDataTableColumn[]>(() => [
  { key: 'createdAt', label: 'Дата', nowrap: true },
  { key: 'event', label: 'Событие', cellClass: 'whitespace-normal wrap-break-word' },
  { key: 'ip', label: 'IP', nowrap: true },
  { key: 'payload', label: 'Payload', cellClass: 'whitespace-normal wrap-break-word' }
])
const _securityEventRows = computed(() =>
  securityEvents.value.map(item => ({
    id: item.event_id,
    createdAt: formatDateTime(item.created_at),
    event: `${item.event_type} · ${item.severity}`,
    ip: item.ip || '—',
    payload: JSON.stringify(item.payload || {}),
    source: item
  }))
)
const avatarGallery = computed(() => getAuthAvatarGallery(user.value))
const avatarGalleryItems = computed(() => avatarGallery.value.items)
const avatarPrimaryItemId = computed(() => avatarGallery.value.primaryId)
const avatarSecondaryGalleryItems = computed(() =>
  avatarGalleryItems.value.filter(item => item.id !== avatarPrimaryItemId.value)
)
const { edges: avatarPreviewScrollState, sync: syncAvatarPreviewScrollState } = useScrollableEdges(
  avatarPreviewScrollerRef,
  { axis: 'x' }
)
const avatarViewerItems = computed<ImageViewerItem[]>(() =>
  avatarGalleryItems.value.map((item, index) => ({
    src: buildMediaFileUrl(item.original_image_key || item.profile_image_key || item.icon_image_key),
    thumbnailSrc: buildMediaFileUrl(item.profile_image_key || item.icon_image_key || item.original_image_key),
    alt:
      avatarGalleryItems.value.length > 1
        ? `${t('auth.account.tab.profile')} ${index + 1}`
        : t('auth.account.tab.profile')
  }))
)
const activePreviewAvatarItem = computed(() => avatarGalleryItems.value[avatarPreviewDialog.index] || null)
const isPreviewAvatarPrimary = computed(() => activePreviewAvatarItem.value?.id === avatarPrimaryItemId.value)
const displayNameText = computed(() => profileForm.display_name.trim() || 'Имя профиля')
const _twofaStatusLabel = computed(() =>
  user.value?.is_two_factor_enabled ? t('auth.account.twofa.enabled') : t('auth.account.twofa.disabled')
)
const groupedSessions = computed<AuthSessionGroupView[]>(() => {
  const groups = new Map<
    string,
    {
      deviceLabel: string
      ip: string
      latestSession: AuthSessionView
      mfaVerified: boolean
      count: number
      revokableSessionIds: string[]
      currentSessionIds: string[]
      hasCurrent: boolean
    }
  >()
  for (const item of sessions.value) {
    if (item.revoked_at) continue
    const deviceLabel =
      String(item.device_label || '').trim() || String(item.user_agent || '').trim() || 'Неизвестное устройство'
    const ip = String(item.ip || '').trim() || '—'
    const key = `${deviceLabel}::${ip}`
    const current = groups.get(key)
    if (!current) {
      groups.set(key, {
        deviceLabel,
        ip,
        latestSession: item,
        mfaVerified: Boolean(item.mfa_verified),
        count: 1,
        revokableSessionIds: [item.session_id],
        currentSessionIds: item.is_current ? [item.session_id] : [],
        hasCurrent: Boolean(item.is_current)
      })
      continue
    }
    current.count += 1
    current.revokableSessionIds.push(item.session_id)
    if (item.is_current) {
      current.currentSessionIds.push(item.session_id)
      current.hasCurrent = true
    }
    const currentTs = new Date(current.latestSession.last_seen_at).getTime() || 0
    const nextTs = new Date(item.last_seen_at).getTime() || 0
    if (nextTs > currentTs) {
      current.latestSession = item
      current.mfaVerified = Boolean(item.mfa_verified)
    }
  }
  return Array.from(groups.entries())
    .map(([key, group]) => ({
      key,
      ip: group.ip,
      deviceLabel: group.deviceLabel,
      count: group.count,
      mfaVerified: group.mfaVerified,
      lastSeenAt: group.latestSession.last_seen_at,
      revokableSessionIds: group.revokableSessionIds,
      currentSessionIds: group.currentSessionIds,
      hasCurrent: group.hasCurrent
    }))
    .sort((a, b) => (new Date(b.lastSeenAt).getTime() || 0) - (new Date(a.lastSeenAt).getTime() || 0))
})

const formatDateTime = (value?: string | null) =>
  formatAbsoluteDateTime(value, { dateStyle: 'medium', timeStyle: 'short' })
const scheduleClientFrame = (cb: () => void) => {
  if (import.meta.client) {
    requestAnimationFrame(cb)
    return
  }
  cb()
}
const paymentStatusLabel = (status?: string | null) => {
  switch (String(status || '').trim()) {
    case 'success':
      return t('payments.status.success')
    case 'pending':
      return t('payments.status.pending')
    case 'failed':
      return t('payments.status.failed')
    case 'canceled':
      return t('payments.status.canceled')
    case 'refunded':
      return t('payments.status.refunded')
    default:
      return t('payments.status.unknown')
  }
}
const paymentPlanLabel = (planCode?: string | null) =>
  String(planCode || '').trim() === 'donation' ? t('payments.plan.donation') : t('payments.plan.pro')
const resolveUploadedImageKey = (res: ApiResponseWithData<MediaUploadResult> | MediaUploadResult) =>
  String('data' in res ? res.data.image_key : res.image_key).trim()
const findAvatarGalleryIndexById = (itemId?: string | null) => {
  const normalizedId = String(itemId || '').trim()
  if (!normalizedId) return 0
  const foundIndex = avatarGalleryItems.value.findIndex(item => item.id === normalizedId)
  return foundIndex >= 0 ? foundIndex : 0
}
const clearSuccessNotice = (key: string, target?: { value: string }) => {
  const existing = successNoticeTimers.get(key)
  if (existing) {
    clearTimeout(existing)
    successNoticeTimers.delete(key)
  }
  if (target) {
    target.value = ''
  }
}
const showSuccessNotice = (key: string, target: { value: string }, message: string) => {
  clearSuccessNotice(key)
  target.value = message
  const timer = setTimeout(() => {
    if (target.value === message) {
      target.value = ''
    }
    successNoticeTimers.delete(key)
  }, SUCCESS_NOTIFY_TIMEOUT_MS)
  successNoticeTimers.set(key, timer)
}
const showProfileNotice = (message: string, tone: NotifyTone, temporary = false) => {
  profileError.value = ''
  clearSuccessNotice('profile', profileInfo)
  profileNoticeTone.value = tone
  profileNoticeTemporary.value = temporary
  if (temporary) {
    showSuccessNotice('profile', profileInfo, message)
    return
  }
  profileInfo.value = message
}

const syncProfileForm = () => {
  if (!user.value) return
  suppressInterfaceAutosave.value = true
  emailChangeForm.email = user.value.email || ''
  profileForm.display_name = user.value.display_name || ''
  interfaceForm.locale = normalizeInterfaceLocale(user.value.locale || uiPreferences.interfaceLocale)
  interfaceForm.theme_preference = readThemePreferenceFromSettings(user.value.settings)
  nextTick(() => {
    suppressInterfaceAutosave.value = false
  })
}

const updateAvatarProfilePayload = async (nextAvatar: AuthAvatarPayload | null) => {
  if (!user.value) throw new Error('Требуется активная сессия.')
  await updateProfile({
    profile: {
      avatar: nextAvatar
    },
    ...(nextAvatar
      ? {}
      : {
          settings: {
            avatar: null
          }
        })
  })
  syncProfileForm()
}
const persistAvatarGallery = async (items: AuthAvatarGalleryItem[], primaryId: string, successMessage: string) => {
  avatarUploading.value = true
  avatarError.value = ''
  clearSuccessNotice('avatar', avatarInfo)
  try {
    const nextPayload = buildAuthAvatarPayload(items, primaryId, AVATAR_ICON_SIZE)
    await updateAvatarProfilePayload(nextPayload)
    showSuccessNotice('avatar', avatarInfo, successMessage)
  } catch (err) {
    avatarError.value = extractApiErrorMessage(err, 'Не удалось сохранить фотографии профиля.')
    throw err
  } finally {
    avatarUploading.value = false
  }
}
const closeAvatarCropDialog = (force = false) => {
  if (avatarUploading.value && !force) return
  avatarCropDialog.open = false
  avatarCropDialog.file = null
}
const openAvatarPreview = (index = findAvatarGalleryIndexById(avatarPrimaryItemId.value)) => {
  if (!avatarViewerItems.value.length) return
  avatarPreviewDialog.index = Math.min(Math.max(index, 0), avatarViewerItems.value.length - 1)
  avatarPreviewDialog.open = true
}
const makePreviewAvatarPrimary = async () => {
  const itemId = activePreviewAvatarItem.value?.id || ''
  if (!itemId) return
  await setPrimaryAvatar(itemId)
}
const removeAvatarItem = async (itemId: string) => {
  const nextItems = avatarGalleryItems.value.filter(item => item.id !== itemId)
  const nextPrimaryId = nextItems.some(item => item.id === avatarPrimaryItemId.value)
    ? avatarPrimaryItemId.value
    : nextItems[0]?.id || ''
  await persistAvatarGallery(
    nextItems,
    nextPrimaryId,
    nextItems.length ? 'Фотография профиля удалена.' : 'Фотографии профиля удалены.'
  )
  if (!nextItems.length) {
    onAvatarPreviewModelUpdate(false)
    return
  }
  avatarPreviewDialog.index = Math.min(avatarPreviewDialog.index, nextItems.length - 1)
}
const setPrimaryAvatar = async (itemId: string) => {
  if (!itemId || itemId === avatarPrimaryItemId.value) return
  await persistAvatarGallery(avatarGalleryItems.value, itemId, 'Основная фотография профиля обновлена.')
}
const refreshState = async () => {
  actionError.value = ''
  try {
    await loadMe()
    if (!user.value) return
    syncProfileForm()
  } catch (err) {
    actionError.value = extractApiErrorMessage(err, 'Не удалось загрузить аккаунт.')
  } finally {
    initialLoading.value = false
  }
}
const loadActivityState = async (force = false) => {
  if (activityLoading.value) return
  if (activityLoaded.value && !force) return
  activityLoading.value = true
  activityError.value = ''
  try {
    const [sessionRes, attemptsRes, eventsRes] = await Promise.all([
      listSessions(),
      listLoginAttempts(),
      listSecurityEvents()
    ])
    sessions.value = sessionRes.data.items || []
    loginAttempts.value = attemptsRes.data.items || []
    securityEvents.value = eventsRes.data.items || []
    activityLoaded.value = true
  } catch (err) {
    activityError.value = extractApiErrorMessage(err, 'Не удалось загрузить сессии и активность.')
  } finally {
    activityLoading.value = false
  }
}
const saveProfile = async () => {
  profileError.value = ''
  clearSuccessNotice('profile', profileInfo)
  profileNoticeTone.value = 'success'
  profileNoticeTemporary.value = true
  try {
    await updateProfile({
      display_name: profileForm.display_name
    })
    editingDisplayName.value = false
    profileInputError.value = false
    showProfileNotice(t('auth.account.display_name_updated'), 'success', true)
    syncProfileForm()
  } catch (err) {
    profileNoticeTemporary.value = false
    profileNoticeTone.value = 'error'
    profileInputError.value = true
    profileError.value = extractApiErrorMessage(err, 'Не удалось обновить профиль.')
  }
}
const submitDisplayNameFromKeyboard = (event: KeyboardEvent) => {
  event.preventDefault()
  event.stopPropagation()
  displayNameKeyboardSubmitPending.value = true
  const target = event.currentTarget as HTMLInputElement | null
  target?.blur()
  void nextTick(async () => {
    try {
      await saveProfile()
    } finally {
      displayNameKeyboardSubmitPending.value = false
    }
  })
}
const saveInterfacePreferences = async () => {
  if (interfaceAutosavePending.value) return
  interfaceError.value = ''
  clearSuccessNotice('interface', interfaceInfo)
  interfaceAutosavePending.value = true
  try {
    await syncInterfacePreferences({
      locale: interfaceForm.locale,
      theme: interfaceForm.theme_preference
    })
    showSuccessNotice('interface', interfaceInfo, t('settings.saved'))
    syncProfileForm()
  } catch (err) {
    interfaceError.value = extractApiErrorMessage(err, 'Не удалось сохранить настройки интерфейса.')
  } finally {
    interfaceAutosavePending.value = false
  }
}
const _submitPasswordChange = async () => {
  passwordError.value = ''
  clearSuccessNotice('password', passwordInfo)
  passwordCurrentInputError.value = false
  try {
    await changePassword(passwordForm.current, passwordForm.next)
    showSuccessNotice('password', passwordInfo, 'Пароль изменён. Выполнен принудительный выход из текущей сессии.')
    passwordForm.current = ''
    passwordForm.next = ''
    await router.push('/auth/login')
  } catch (err) {
    const message = extractApiErrorMessage(err, 'Не удалось изменить пароль.')
    passwordError.value = message
    if (message.toLowerCase().includes('current password') || message.toLowerCase().includes('текущ')) {
      passwordCurrentInputError.value = true
    }
  }
}
const _submitEmailChange = async () => {
  emailChangeError.value = ''
  emailChangeInfo.value = ''
  const nextEmail = emailChangeForm.email.trim()
  if (!nextEmail) {
    emailChangeError.value = 'Введите новый email.'
    return
  }
  emailChangePending.value = true
  try {
    await requestEmailChange(nextEmail)
    emailChangeInfo.value = t('auth.account.email.request_sent', { email: nextEmail })
  } catch (err) {
    emailChangeError.value = extractApiErrorMessage(err, 'Не удалось отправить письмо для подтверждения email.')
  } finally {
    emailChangePending.value = false
  }
}
const _begin2faSetup = async () => {
  if (user.value?.is_two_factor_enabled) return
  twofaError.value = ''
  clearSuccessNotice('twofa', twofaInfo)
  backupCodes.value = []
  try {
    const res = await setupTwoFactor()
    setupSecret.value = res.data.secret
    setupOtpAuthUrl.value = res.data.otpauth_url || ''
    setupQrDataUrl.value = res.data.qr_data_url || ''
    clearSuccessNotice('twofa', twofaInfo)
    twofaInfo.value = t('auth.account.twofa.setup_created')
  } catch (err) {
    twofaError.value = extractApiErrorMessage(err, 'Не удалось начать настройку 2FA.')
  }
}
const _activate2fa = async () => {
  twofaError.value = ''
  clearSuccessNotice('twofa', twofaInfo)
  enableCode.value = enableCode.value.replace(/\D+/g, '').slice(0, 6)
  if (enableCode.value.length !== 6) {
    twofaError.value = t('auth.account.twofa.code_hint')
    return
  }
  try {
    const res = await enableTwoFactor(enableCode.value)
    backupCodes.value = res.data.backup_codes || []
    showSuccessNotice('twofa', twofaInfo, '2FA включена. Backup codes показаны один раз.')
    setupSecret.value = ''
    setupOtpAuthUrl.value = ''
    setupQrDataUrl.value = ''
    enableCode.value = ''
    showDisable2faForm.value = false
    await refreshState()
  } catch (err) {
    twofaError.value = getTwoFactorErrorMessage(err, 'totp', 'Неверный код из приложения. Попробуйте ещё раз.')
  }
}
const focusDisableCodeInput = () => {
  if (disableMethodTab.value === 'backup') {
    disableBackupCodeInputRef.value?.focus()
    return
  }
  disableTotpInputRef.value?.focus()
}
const _deactivate2fa = async (rawCode?: string) => {
  twofaError.value = ''
  clearSuccessNotice('twofa', twofaInfo)
  const password = disablePassword.value.trim()
  if (!password) {
    twofaError.value = 'Введите текущий пароль.'
    return
  }
  const code =
    typeof rawCode === 'string'
      ? rawCode
      : disableMethodTab.value === 'backup'
        ? `${disableBackupCode.value.slice(0, 4)}-${disableBackupCode.value.slice(4, 8)}`
        : disableTotpCode.value
  const normalizedCode =
    disableMethodTab.value === 'backup'
      ? code.toUpperCase().replace(/[^A-Z0-9-]+/g, '')
      : code.replace(/\D+/g, '').slice(0, 6)
  if (disableMethodTab.value === 'backup') {
    const compact = normalizedCode.replace(/-/g, '')
    if (compact.length !== 8) {
      twofaError.value = 'Введите 8 символов кода сброса.'
      return
    }
  } else if (normalizedCode.length !== 6) {
    twofaError.value = 'Введите 6 цифр кода из приложения.'
    return
  }
  try {
    await disableTwoFactor(password, normalizedCode)
    showSuccessNotice('twofa', twofaInfo, '2FA отключена.')
    setupSecret.value = ''
    setupOtpAuthUrl.value = ''
    setupQrDataUrl.value = ''
    backupCodes.value = []
    disablePassword.value = ''
    disableTotpCode.value = ''
    disableBackupCode.value = ''
    showDisable2faForm.value = false
    await refreshState()
  } catch (err) {
    twofaError.value = getTwoFactorErrorMessage(
      err,
      disableMethodTab.value === 'backup' ? 'backup' : 'totp',
      disableMethodTab.value === 'backup'
        ? 'Неверный код сброса. Попробуйте ещё раз.'
        : 'Неверный код из приложения. Попробуйте ещё раз.'
    )
    if (disableMethodTab.value === 'backup') {
      disableBackupCode.value = ''
    } else {
      disableTotpCode.value = ''
    }
    await nextTick()
    focusDisableCodeInput()
  }
}
const _toggleDisable2faForm = () => {
  const next = !showDisable2faForm.value
  showDisable2faForm.value = next
  if (next) {
    backupCodes.value = []
    twofaError.value = ''
    disableMethodTab.value = 'totp'
    disableTotpCode.value = ''
    disableBackupCode.value = ''
    nextTick().then(() => {
      scheduleClientFrame(() => {
        focusDisableCodeInput()
      })
    })
  }
}
const leave = async (all: boolean) => {
  showLogoutActions.value = false
  if (all) {
    await logoutAll()
  } else {
    await logout()
  }
  await router.push('/auth/login')
}
const loadPaymentHistoryState = async (force = false) => {
  paymentHistoryError.value = ''
  try {
    await loadHistory(force)
  } catch (err) {
    paymentHistoryError.value = extractApiErrorMessage(err, 'Не удалось загрузить историю оплат.')
  }
}
const _submitRefund = async (orderID: string) => {
  if (!orderID || refundPendingOrderId.value) return
  paymentHistoryError.value = ''
  paymentHistoryInfo.value = ''
  refundPendingOrderId.value = orderID
  try {
    await refundOrder(orderID)
    await Promise.all([loadAccess(true), loadPaymentHistoryState(true)])
    paymentHistoryInfo.value = 'Возврат отправлен в платёжную систему.'
  } catch (err) {
    paymentHistoryError.value = extractApiErrorMessage(err, 'Не удалось выполнить возврат.')
  } finally {
    refundPendingOrderId.value = ''
  }
}
const _revokeSessionGroup = async (group: AuthSessionGroupView) => {
  if (!group.revokableSessionIds.length) return
  actionError.value = ''
  clearSuccessNotice('action', actionInfo)
  try {
    const currentSet = new Set(group.currentSessionIds)
    const revokeOrder = [
      ...group.revokableSessionIds.filter(sessionId => !currentSet.has(sessionId)),
      ...group.currentSessionIds
    ]
    for (const sessionId of revokeOrder) {
      await revokeSession(sessionId)
    }
    if (group.hasCurrent) {
      await logout()
      await router.push('/auth/login')
      return
    }
    showSuccessNotice('action', actionInfo, 'Сессии устройства отозваны.')
    await loadActivityState(true)
  } catch (err) {
    actionError.value = extractApiErrorMessage(err, 'Не удалось отозвать сессии устройства.')
  }
}
syncProfileForm()

const onAvatarFileChange = async (event: Event) => {
  const input = event.target as HTMLInputElement | null
  const file = input?.files?.[0] || null
  if (input) input.value = ''
  if (!file) return
  avatarError.value = ''
  clearSuccessNotice('avatar', avatarInfo)
  if (file.size > AVATAR_IMAGE_MAX_BYTES) {
    avatarError.value = `Файл слишком большой. Максимум ${AVATAR_IMAGE_MAX_MB} МБ.`
    return
  }
  avatarCropDialog.file = file
  avatarCropDialog.open = true
}
const onAvatarCropConfirm = async (file: File) => {
  avatarUploading.value = true
  avatarError.value = ''
  clearSuccessNotice('avatar', avatarInfo)
  try {
    const originalFile = avatarCropDialog.file
    if (!originalFile) {
      throw new Error('Не найден исходный файл аватарки.')
    }
    const icon64 = await resizeImageToSquareFile(file, AVATAR_ICON_SIZE, {
      suffix: `avatar-${AVATAR_ICON_SIZE}`
    })
    const [originalRes, profileRes, iconRes] = await Promise.all([
      uploadMediaFile(originalFile, { section: 'users', collection: 'avatars' }),
      uploadMediaFile(file, { section: 'users', collection: 'avatars' }),
      uploadMediaFile(icon64, { section: 'users', collection: 'avatars' })
    ])
    const originalKey = resolveUploadedImageKey(originalRes)
    const iconKey = resolveUploadedImageKey(iconRes)
    const profileKey = resolveUploadedImageKey(profileRes)
    if (!originalKey || !iconKey || !profileKey) {
      throw new Error('Сервер не вернул ключи для аватарки.')
    }
    const nextItem: AuthAvatarGalleryItem = {
      id: globalThis.crypto.randomUUID(),
      original_image_key: originalKey,
      icon_image_key: iconKey,
      profile_image_key: profileKey,
      icon_size: AVATAR_ICON_SIZE,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
    const nextItems = [nextItem, ...avatarGalleryItems.value]
    await updateAvatarProfilePayload(buildAuthAvatarPayload(nextItems, nextItem.id, AVATAR_ICON_SIZE))
    showSuccessNotice('avatar', avatarInfo, 'Фотография профиля добавлена.')
    avatarPreviewDialog.index = 0
    closeAvatarCropDialog(true)
  } catch (err) {
    avatarError.value = extractApiErrorMessage(err, 'Не удалось обработать и сохранить аватарку.')
  } finally {
    avatarUploading.value = false
  }
}
const onAvatarPreviewModelUpdate = (value: boolean) => {
  avatarPreviewDialog.open = value
}
const onAvatarPreviewIndexChange = (index: number) => {
  avatarPreviewDialog.index = index
}
const onWindowPointerDown = (event: PointerEvent) => {
  const target = event.target as Node | null
  if (!target) return
  if (editingDisplayName.value && !profileNameEditorRef.value?.contains(target)) {
    cancelProfileEdit()
  }
}
const onProfileEditorFocusout = (event: FocusEvent) => {
  if (!editingDisplayName.value) return
  if (displayNameKeyboardSubmitPending.value) return
  const nextTarget = event.relatedTarget as Node | null
  if (nextTarget && profileNameEditorRef.value?.contains(nextTarget)) return
  if (!import.meta.client) return
  queueMicrotask(() => {
    const activeElement = document.activeElement
    if (activeElement && profileNameEditorRef.value?.contains(activeElement)) return
    cancelProfileEdit()
  })
}
const _onEnableCodeInput = (value: string) => {
  enableCode.value = value
  if (twofaError.value) {
    twofaError.value = ''
  }
}
const _onDisableTotpCodeInput = (value: string) => {
  disableTotpCode.value = value
  if (twofaError.value) {
    twofaError.value = ''
  }
}
const _onDisableBackupCodeInput = (value: string) => {
  disableBackupCode.value = value
  if (twofaError.value) {
    twofaError.value = ''
  }
}

const cancelProfileEdit = () => {
  const actualDisplayName = user.value?.display_name || ''
  const hasChanges = profileForm.display_name !== actualDisplayName
  editingDisplayName.value = false
  suppressProfileErrorReset.value = true
  profileForm.display_name = actualDisplayName
  if (hasChanges) {
    profileInputError.value = false
    showProfileNotice(t('auth.account.display_name_cancelled'), 'info', true)
  }
  nextTick(() => {
    suppressProfileErrorReset.value = false
  })
}

onMounted(async () => {
  if (import.meta.client) {
    window.addEventListener('pointerdown', onWindowPointerDown)
  }
  initialLoading.value = false
})
onBeforeUnmount(() => {
  if (import.meta.client) {
    window.removeEventListener('pointerdown', onWindowPointerDown)
  }
  for (const timer of successNoticeTimers.values()) {
    clearTimeout(timer)
  }
  successNoticeTimers.clear()
})

watch(
  () => Boolean(user.value?.is_two_factor_enabled),
  enabled => {
    if (enabled) {
      setupSecret.value = ''
      setupOtpAuthUrl.value = ''
      setupQrDataUrl.value = ''
      enableCode.value = ''
      return
    }
    showDisable2faForm.value = false
    disablePassword.value = ''
    disableTotpCode.value = ''
    disableBackupCode.value = ''
  },
  { immediate: true }
)
watch(
  () => setupSecret.value,
  async value => {
    if (!value || user.value?.is_two_factor_enabled) return
    await nextTick()
    enableCodeInputRef.value?.focus()
  }
)
watch(disableMethodTab, async () => {
  if (!showDisable2faForm.value) return
  twofaError.value = ''
  await nextTick()
  scheduleClientFrame(() => {
    focusDisableCodeInput()
  })
})
watch(
  () => profileForm.display_name,
  () => {
    if (suppressProfileErrorReset.value) return
    if (profileInputError.value) {
      profileInputError.value = false
      profileError.value = ''
    }
  }
)
watch(
  () => passwordForm.current,
  () => {
    if (passwordCurrentInputError.value) {
      passwordCurrentInputError.value = false
    }
  }
)
watch(editingDisplayName, async value => {
  if (!value) return
  profileForm.display_name = user.value?.display_name || ''
  profileInputError.value = false
  profileError.value = ''
  clearSuccessNotice('profile', profileInfo)
  await nextTick()
  displayNameInputRef.value?.focus()
})
watch(
  () => [interfaceForm.locale, interfaceForm.theme_preference] as const,
  ([locale, theme], [prevLocale, prevTheme]) => {
    if (suppressInterfaceAutosave.value) return
    if (locale === prevLocale && theme === prevTheme) return
    void saveInterfacePreferences()
  }
)
watch(
  [accountTab, user],
  ([tab, currentUser]) => {
    if (!currentUser) return
    if (tab === 'security') {
      void loadActivityState()
    }
    if (tab === 'balance') {
      void loadPaymentHistoryState()
    }
  },
  { immediate: true }
)
watch(
  () => avatarGalleryItems.value.length,
  length => {
    if (!length) {
      avatarPreviewDialog.index = 0
      return
    }
    if (avatarPreviewDialog.index < length) return
    avatarPreviewDialog.index = length - 1
  }
)
watch(
  avatarSecondaryGalleryItems,
  () => {
    nextTick(syncAvatarPreviewScrollState)
  },
  { deep: true }
)
</script>
<template>
  <div>
    <LabNavHeader :title :breadcrumb-items="accountBreadcrumbItems" />
    <LabLoader v-if="initialLoading" variant="inline" label="Загрузка аккаунта…" />
    <AuthFeatureGateNotice
      v-else-if="!user"
      message="Войдите в аккаунт, чтобы управлять профилем, сессиями и параметрами безопасности."
    />
    <template v-else>
      <LabNavTabs
        v-model="accountTab"
        :items="accountTabItems"
        route-param-key="account"
        :route-target-map="accountTabRouteTargetMap"
      >
        <template #panel-profile>
          <section class="flex flex-wrap items-start gap-4">
            <div class="flex shrink-0 flex-col items-start gap-3">
              <div class="relative overflow-hidden">
                <LabAvatar version="profile" :user="user" clickable @click="openAvatarPreview()" />
                <label
                  class="absolute right-2 bottom-2"
                  :class="avatarUploading ? 'cursor-not-allowed' : 'cursor-pointer'"
                >
                  <LabBaseInput
                    id="account-avatar-file"
                    name="account_avatar_file"
                    type="file"
                    :accept="AVATAR_IMAGE_ACCEPT"
                    class="hidden"
                    :disabled="avatarUploading"
                    @change="onAvatarFileChange"
                  />
                  <LabBaseButton
                    tag="span"
                    variant="secondary"
                    size="lg"
                    button-class="rounded-full pointer-events-none"
                    icon="ic:round-photo-camera"
                    icon-only
                    title="Загрузить новую аватарку"
                    aria-label="Загрузить новую аватарку"
                    :disabled="avatarUploading"
                  />
                </label>
              </div>
              <div
                v-if="avatarSecondaryGalleryItems.length"
                class="relative w-36 [--avatar-preview-overlap:calc(var(--avatar-preview-size)/2)] [--avatar-preview-size:3rem] sm:w-40"
              >
                <div ref="avatarPreviewScrollerRef" class="lab-scroll-hidden overflow-x-auto">
                  <div class="flex min-w-max items-center pr-6 pb-2">
                    <button
                      v-for="(item, index) in avatarSecondaryGalleryItems"
                      :key="item.id"
                      type="button"
                      class="h-(--avatar-preview-size) w-(--avatar-preview-size) shrink-0 overflow-hidden rounded-full border-2 border-(--lab-bg-canvas) bg-(--lab-bg-canvas)"
                      :style="{
                        marginInlineStart: index === 0 ? '0' : 'calc(var(--avatar-preview-overlap) * -1)',
                        zIndex: avatarSecondaryGalleryItems.length - index
                      }"
                      aria-label="Открыть фотографию профиля"
                      @click="openAvatarPreview(findAvatarGalleryIndexById(item.id))"
                    >
                      <img
                        :src="
                          buildMediaFileUrl(item.profile_image_key || item.icon_image_key || item.original_image_key)
                        "
                        alt="Фотография профиля"
                        class="h-full w-full object-cover"
                      />
                    </button>
                  </div>
                </div>
                <div
                  class="lab-scroll-fade lab-scroll-fade-x-left bottom-2"
                  :class="{ 'lab-scroll-fade-visible': avatarPreviewScrollState.left }"
                  aria-hidden="true"
                />
                <div
                  class="lab-scroll-fade lab-scroll-fade-x-right bottom-2"
                  :class="{ 'lab-scroll-fade-visible': avatarPreviewScrollState.right }"
                  aria-hidden="true"
                />
              </div>
              <LabNotify :text="avatarError" tone="error" size="xs" />
              <LabNotify :text="avatarInfo" tone="success" size="xs" />
            </div>
            <div class="flex min-w-fit flex-1 flex-col gap-3">
              <article class="flex flex-col gap-2">
                <div class="flex min-w-fit flex-wrap items-start justify-between gap-3">
                  <div ref="profileNameEditorRef" class="flex min-w-fit flex-col gap-2">
                    <div
                      v-if="!editingDisplayName"
                      class="inline-grid grid-cols-[max-content_max-content] items-start gap-1.5"
                    >
                      <button
                        type="button"
                        class="lab-focus inline-flex items-center border-b border-transparent bg-transparent py-0 text-left text-xl font-semibold text-(--lab-text-primary) transition-colors"
                        @click="editingDisplayName = true"
                      >
                        <span>{{ displayNameText }}</span>
                      </button>
                      <LabBaseTooltip
                        v-if="profileSubscriptionTooltipText"
                        :text="profileSubscriptionTooltipText"
                        side="right"
                        align="left"
                        :offset="10"
                        :cross-axis-offset="0"
                      >
                        <template #trigger>
                          <LabBaseButton
                            icon="ic:round-auto-awesome"
                            icon-only
                            variant="ghost"
                            size="sm"
                            button-class="h-8 w-8 rounded-full border-transparent text-orange-300 hover:bg-(--lab-bg-surface-hover) focus:bg-(--lab-bg-surface-hover) focus-visible:bg-(--lab-bg-surface-hover)"
                            aria-label="Статус подписки"
                          />
                        </template>
                      </LabBaseTooltip>
                    </div>
                    <div v-else class="min-w-fit" @focusout="onProfileEditorFocusout">
                      <div class="inline-grid grid-cols-[max-content_2rem] items-start gap-1.5">
                        <input
                          id="account-profile-name"
                          ref="displayNameInputRef"
                          v-model="profileForm.display_name"
                          name="profile_name"
                          type="text"
                          autocomplete="off"
                          autocapitalize="words"
                          autocorrect="off"
                          spellcheck="false"
                          data-form-type="other"
                          data-lpignore="true"
                          data-1p-ignore="true"
                          :size="Math.max((profileForm.display_name || displayNameText).length, 1)"
                          class="lab-focus w-auto appearance-none border-b bg-transparent p-0 text-xl font-semibold text-(--lab-text-primary)"
                          :class="profileInputError ? 'text-(--lab-danger)' : ''"
                          placeholder="Имя профиля"
                          @keydown.enter="submitDisplayNameFromKeyboard"
                          @keydown.esc.prevent="cancelProfileEdit"
                        />
                        <LabBaseButton
                          icon="ic:round-check"
                          icon-only
                          variant="secondary"
                          size="sm"
                          button-class="text-(--lab-text-secondary) hover:text-(--lab-info) focus-visible:text-(--lab-info) h-8 w-8"
                          title="Сохранить никнейм"
                          aria-label="Сохранить никнейм"
                          @click="saveProfile"
                        />
                      </div>
                    </div>
                    <LabNotify
                      :text="profileError || profileInfo"
                      :tone="profileError ? 'error' : profileNoticeTone"
                      :temporary="profileError ? false : profileNoticeTemporary"
                      size="xs"
                    />
                  </div>
                </div>
                <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                  <span class="shrink-0 text-xs tracking-wide text-(--lab-text-muted) uppercase">
                    {{ t('auth.account.email_label') }}
                  </span>
                  <span class="min-w-56 flex-1 basis-56 text-sm wrap-break-word text-(--lab-text-primary)">
                    {{ user.email }}
                  </span>
                </div>
              </article>
              <article class="space-y-3">
                <LabNotify :text="interfaceError" tone="error" size="xs" />
                <LabNotify :text="interfaceInfo" tone="success" size="xs" />
              </article>
              <LabDropdown
                v-model="showLogoutActions"
                side="top"
                align="right"
                width-class="w-56"
                :match-trigger-width="false"
                :close-on-select="false"
              >
                <template #trigger="{ toggle }">
                  <LabBaseButton
                    variant="danger"
                    size="lg"
                    :label="t('auth.account.logout')"
                    button-class="text-xs"
                    @click="toggle"
                  />
                </template>
                <LabBaseButton
                  variant="plain"
                  size="lg"
                  block
                  :label="t('auth.account.logout.current')"
                  button-class="justify-start text-xs"
                  @click="leave(false)"
                />
                <LabBaseButton
                  variant="danger"
                  size="lg"
                  block
                  :label="t('auth.account.logout.all')"
                  button-class="justify-start text-xs"
                  @click="leave(true)"
                />
              </LabDropdown>
            </div>
          </section>
        </template>
        <template #panel-balance>
          <LazyAuthAccountBalancePanel v-if="accountTab === 'balance'" />
        </template>
        <template #panel-security>
          <LazyAuthAccountSecurityPanel v-if="accountTab === 'security'" />
        </template>
      </LabNavTabs>
      <LazyLabCropperModal
        v-if="avatarCropDialog.open && avatarCropDialog.file"
        :file="avatarCropDialog.file"
        title="Кадрирование аватарки"
        :loading="avatarUploading"
        aspect-preset="1:1"
        :aspect-locked="true"
        @cancel="closeAvatarCropDialog"
        @confirm="onAvatarCropConfirm"
      />
      <LazyLabViewerImage
        v-model="avatarPreviewDialog.open"
        :items="avatarViewerItems"
        :initial-index="avatarPreviewDialog.index"
        @active-index-change="onAvatarPreviewIndexChange"
        @update:model-value="onAvatarPreviewModelUpdate"
      >
        <template #toolbar-extra>
          <div v-if="activePreviewAvatarItem" class="flex items-center gap-2">
            <LabBaseButton
              :variant="isPreviewAvatarPrimary ? 'success' : 'primary'"
              size="sm"
              :disabled="avatarUploading"
              :aria-label="isPreviewAvatarPrimary ? 'Основная фотография профиля' : 'Сделать фотографию основной'"
              :label="isPreviewAvatarPrimary ? 'Основная фотография' : 'Сделать основной'"
              @click.stop="makePreviewAvatarPrimary"
            />
            <LabConfirmActionButton
              :disabled="avatarUploading"
              icon="ic:round-delete-outline"
              confirm-icon="ic:round-check"
              label="Удалить фотографию"
              title="Удалить текущую фотографию профиля"
              idle-class="lab-button lab-button-secondary"
              confirm-class="lab-button lab-button-danger text-white hover:text-white hover:bg-rose-500"
              progress-class="bg-rose-300/30"
              aria-label="Удалить текущую фотографию профиля"
              confirm-aria-label="Подтвердить удаление текущей фотографии профиля"
              tooltip="Удалить фото?"
              @confirm="removeAvatarItem(activePreviewAvatarItem.id)"
            />
          </div>
        </template>
      </LazyLabViewerImage>
    </template>
  </div>
</template>
