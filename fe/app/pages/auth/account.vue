<script setup lang="ts">
  const { t, readThemePreferenceFromSettings } = useInterfacePreferences()
  const title = t('auth.account.title')
  usePageSeo({
    title,
    description: t('auth.account.description')
  })
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

  const { access, accessLoading, ensureAccessLoaded } = usePayments()

  await ensureLoaded()
  await ensureAccessLoaded()

  const SUCCESS_NOTIFY_TIMEOUT_MS = 5000
  const AVATAR_IMAGE_ACCEPT = 'image/jpeg,image/png,image/webp,image/gif,image/avif'
  const avatarImageMaxBytes = Number(runtimeConfig.public.mediaImageMaxBytes || 8388608)
  const AVATAR_IMAGE_MAX_MB = Math.max(1, Math.ceil(avatarImageMaxBytes / (1024 * 1024)))
  const AVATAR_IMAGE_MAX_BYTES = avatarImageMaxBytes

  const accountTabItems: LabTabItem[] = [
    { value: 'profile', label: t('auth.account.tab.profile') },
    { value: 'security', label: t('auth.account.tab.security') }
  ]
  const activityTabItems: LabTabItem[] = [
    { value: 'sessions', label: t('auth.account.activity.sessions') },
    { value: 'attempts', label: t('auth.account.activity.attempts') },
    { value: 'events', label: t('auth.account.activity.events') }
  ]
  const disableCodeTabItems: LabTabItem[] = [
    { value: 'totp', label: t('auth.account.twofa.code_label') },
    { value: 'backup', label: t('auth.login.mfa_backup') }
  ]

  const accountTab = ref<'profile' | 'security'>('profile')
  const activityTab = ref<'sessions' | 'attempts' | 'events'>('sessions')
  const disableMethodTab = ref<'totp' | 'backup'>('totp')

  const initialLoading = ref(!loaded.value)
  const editingDisplayName = ref(false)
  const showLogoutActions = ref(false)
  const showDisable2faForm = ref(false)
  const avatarUploading = ref(false)

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
  const paymentLatestOrderStatusLabel = computed(() => {
    const status = String(access.value?.latest_order?.status || '').trim()
    switch (status) {
      case 'success':
        return 'Оплачен'
      case 'pending':
        return 'Ожидает подтверждения'
      case 'failed':
        return 'Неуспешен'
      case 'canceled':
        return 'Отменён'
      case 'refunded':
        return 'Возвращён'
      default:
        return '—'
    }
  })
  const paymentLatestOrderAmountText = computed(() => {
    const amount = Number(access.value?.latest_order?.amount || 0)
    return new Intl.NumberFormat('ru-RU').format(Math.floor(amount / 100)) + ' ₽'
  })
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
  const twofaStatusLabel = computed(() =>
    user.value?.is_two_factor_enabled ? t('auth.account.twofa.enabled') : t('auth.account.twofa.disabled')
  )
  const groupedSessions = computed<AuthSessionGroupView[]>(() => {
    const groups = new Map<
      string,
      {
        browser: string
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
      const browser = browserLabelFromUserAgent(item.user_agent)
      const ip = String(item.ip || '').trim() || '—'
      const key = `${browser}::${ip}`
      const current = groups.get(key)
      if (!current) {
        groups.set(key, {
          browser,
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
      current.mfaVerified = current.mfaVerified || Boolean(item.mfa_verified)
      current.revokableSessionIds.push(item.session_id)
      if (item.is_current) {
        current.currentSessionIds.push(item.session_id)
        current.hasCurrent = true
      }
      const currentTs = new Date(current.latestSession.last_seen_at).getTime() || 0
      const nextTs = new Date(item.last_seen_at).getTime() || 0
      if (nextTs > currentTs) {
        current.latestSession = item
      }
    }
    return Array.from(groups.entries())
      .map(([key, group]) => ({
        key,
        browser: group.browser,
        ip: group.ip,
        deviceLabel: group.browser || group.latestSession.device_label || 'Неизвестное устройство',
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
  const browserLabelFromUserAgent = (userAgent?: string) => {
    const raw = String(userAgent || '').toLowerCase()
    if (!raw) return 'Браузер'
    if (raw.includes('yabrowser')) return 'Yandex Browser'
    if (raw.includes('edg/')) return 'Microsoft Edge'
    if (raw.includes('opr/') || raw.includes('opera')) return 'Opera'
    if (raw.includes('firefox/')) return 'Firefox'
    if (raw.includes('chrome/') && !raw.includes('chromium')) return 'Chrome'
    if (raw.includes('safari/') && !raw.includes('chrome/')) return 'Safari'
    if (raw.includes('curl/')) return 'cURL'
    return 'Браузер'
  }
  const resolveUploadedImageKey = (res: any) => String(res?.data?.image_key || res?.image_key || '').trim()
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

  const updateAvatarProfilePayload = async (nextAvatar: Record<string, any> | null) => {
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
    } catch (err: any) {
      avatarError.value = err?.data?.message || err?.message || 'Не удалось сохранить фотографии профиля.'
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
    } catch (err: any) {
      actionError.value = err?.data?.message || err?.message || 'Не удалось загрузить аккаунт.'
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
    } catch (err: any) {
      activityError.value = err?.data?.message || err?.message || 'Не удалось загрузить сессии и активность.'
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
    } catch (err: any) {
      profileNoticeTemporary.value = false
      profileNoticeTone.value = 'error'
      profileInputError.value = true
      profileError.value = err?.data?.message || err?.message || 'Не удалось обновить профиль.'
    }
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
    } catch (err: any) {
      interfaceError.value = err?.data?.message || err?.message || 'Не удалось сохранить настройки интерфейса.'
    } finally {
      interfaceAutosavePending.value = false
    }
  }
  const submitPasswordChange = async () => {
    passwordError.value = ''
    clearSuccessNotice('password', passwordInfo)
    passwordCurrentInputError.value = false
    try {
      await changePassword(passwordForm.current, passwordForm.next)
      showSuccessNotice('password', passwordInfo, 'Пароль изменён. Выполнен принудительный выход из текущей сессии.')
      passwordForm.current = ''
      passwordForm.next = ''
      await router.push('/auth/login')
    } catch (err: any) {
      const message = String(err?.data?.message || err?.message || 'Не удалось изменить пароль.')
      passwordError.value = message
      if (message.toLowerCase().includes('current password') || message.toLowerCase().includes('текущ')) {
        passwordCurrentInputError.value = true
      }
    }
  }
  const submitEmailChange = async () => {
    emailChangeError.value = ''
    clearSuccessNotice('email-change', emailChangeInfo)
    const nextEmail = emailChangeForm.email.trim()
    if (!nextEmail) {
      emailChangeError.value = 'Введите новый email.'
      return
    }
    emailChangePending.value = true
    try {
      await requestEmailChange(nextEmail)
      showSuccessNotice('email-change', emailChangeInfo, t('auth.account.email.request_sent', { email: nextEmail }))
    } catch (err: any) {
      emailChangeError.value =
        err?.data?.message || err?.message || 'Не удалось отправить письмо для подтверждения email.'
    } finally {
      emailChangePending.value = false
    }
  }
  const begin2faSetup = async () => {
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
    } catch (err: any) {
      twofaError.value = extractApiErrorMessage(err, 'Не удалось начать настройку 2FA.')
    }
  }
  const activate2fa = async () => {
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
    } catch (err: any) {
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
  const deactivate2fa = async (rawCode?: string) => {
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
    } catch (err: any) {
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
  const toggleDisable2faForm = () => {
    const next = !showDisable2faForm.value
    showDisable2faForm.value = next
    if (next) {
      backupCodes.value = []
      twofaError.value = ''
      disableMethodTab.value = 'totp'
      disableTotpCode.value = ''
      disableBackupCode.value = ''
      nextTick().then(() => {
        if (import.meta.client) {
          requestAnimationFrame(() => {
            focusDisableCodeInput()
          })
          return
        }
        focusDisableCodeInput()
      })
    }
  }
  const goToAdmin = async () => {
    await router.push('/auth/admin')
  }
  const leave = async (all: boolean) => {
    showLogoutActions.value = false
    all ? await logoutAll() : await logout()
    await router.push('/auth/login')
  }
  const revokeSessionGroup = async (group: AuthSessionGroupView) => {
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
    } catch (err: any) {
      actionError.value = err?.data?.message || err?.message || 'Не удалось отозвать сессии устройства.'
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
    } catch (err: any) {
      avatarError.value = err?.data?.message || err?.message || 'Не удалось обработать и сохранить аватарку.'
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
    const nextTarget = event.relatedTarget as Node | null
    if (nextTarget && profileNameEditorRef.value?.contains(nextTarget)) return
    if (!import.meta.client) return
    queueMicrotask(() => {
      const activeElement = document.activeElement
      if (activeElement && profileNameEditorRef.value?.contains(activeElement)) return
      cancelProfileEdit()
    })
  }
  const onEnableCodeInput = (value: string) => {
    enableCode.value = value
    if (twofaError.value) {
      twofaError.value = ''
    }
  }
  const onDisableTotpCodeInput = (value: string) => {
    disableTotpCode.value = value
    if (twofaError.value) {
      twofaError.value = ''
    }
  }
  const onDisableBackupCodeInput = (value: string) => {
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
    if (import.meta.client) {
      requestAnimationFrame(() => {
        focusDisableCodeInput()
      })
      return
    }
    focusDisableCodeInput()
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
      if (tab !== 'security' || !currentUser) return
      void loadActivityState()
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
    <LabNavHeader :title class="min-w-0 flex-1" />
    <section v-if="initialLoading" class="p-4">
      <LabLoader variant="inline" label="Загрузка аккаунта…" />
    </section>
    <section v-else-if="!user" class="p-4">
      <AuthFeatureGateNotice
        message="Войдите в аккаунт, чтобы управлять профилем, сессиями и параметрами безопасности." />
    </section>
    <template v-else>
      <LabNavTabs
        v-model="accountTab"
        :items="accountTabItems"
        :no-select="true"
        route-query-key="account"
        route-default-value="profile"
        panel-class="space-y-4 p-4">
        <template #panel-profile>
          <section class="text-(--lab-text-primary) space-y-4">
            <section class="flex flex-wrap items-start gap-4">
              <div class="flex shrink-0 flex-col items-start gap-3">
                <div class="relative overflow-hidden">
                  <LabAvatar version="profile" :user="user" clickable @click="openAvatarPreview()" />
                  <label
                    class="absolute bottom-2 right-2"
                    :class="avatarUploading ? 'cursor-not-allowed' : 'cursor-pointer'">
                    <LabBaseInput
                      id="account-avatar-file"
                      name="account_avatar_file"
                      type="file"
                      :accept="AVATAR_IMAGE_ACCEPT"
                      class="hidden"
                      :disabled="avatarUploading"
                      @change="onAvatarFileChange" />
                    <LabBaseButton
                      tag="span"
                      variant="secondary"
                      size="lg"
                      button-class="pointer-events-none"
                      icon="ic:round-photo-camera"
                      iconOnly
                      title="Загрузить новую аватарку"
                      aria-label="Загрузить новую аватарку"
                      :disabled="avatarUploading" />
                  </label>
                </div>
                <div
                  v-if="avatarSecondaryGalleryItems.length"
                  class="relative w-36 [--avatar-preview-size:3rem] [--avatar-preview-overlap:calc(var(--avatar-preview-size)/2)] sm:w-40">
                  <div ref="avatarPreviewScrollerRef" class="lab-scroll-hidden overflow-x-auto">
                    <div class="flex min-w-max items-center pb-2 pr-6">
                      <button
                        v-for="(item, index) in avatarSecondaryGalleryItems"
                        :key="item.id"
                        type="button"
                        class="border-(--lab-bg-canvas) h-(--avatar-preview-size) w-(--avatar-preview-size) shrink-0 overflow-hidden rounded-full border-2 bg-(--lab-bg-canvas)"
                        :style="{
                          marginInlineStart: index === 0 ? '0' : 'calc(var(--avatar-preview-overlap) * -1)',
                          zIndex: avatarSecondaryGalleryItems.length - index
                        }"
                        aria-label="Открыть фотографию профиля"
                        @click="openAvatarPreview(findAvatarGalleryIndexById(item.id))">
                        <img
                          :src="
                            buildMediaFileUrl(item.profile_image_key || item.icon_image_key || item.original_image_key)
                          "
                          alt="Фотография профиля"
                          class="h-full w-full object-cover" />
                      </button>
                    </div>
                  </div>
                  <div
                    class="lab-scroll-fade lab-scroll-fade-x-left bottom-2"
                    :class="{ 'lab-scroll-fade-visible': avatarPreviewScrollState.left }"
                    aria-hidden="true"></div>
                  <div
                    class="lab-scroll-fade lab-scroll-fade-x-right bottom-2"
                    :class="{ 'lab-scroll-fade-visible': avatarPreviewScrollState.right }"
                    aria-hidden="true"></div>
                </div>
                <LabNotify :text="avatarError" tone="error" size="xs" />
                <LabNotify :text="avatarInfo" tone="success" size="xs" />
              </div>
              <div class="min-w-fit flex flex-1 flex-col gap-3">
                <article class="flex flex-col gap-2">
                  <div class="flex min-w-fit flex-wrap items-start justify-between gap-3">
                    <div ref="profileNameEditorRef" class="min-w-fit flex flex-col gap-2">
                      <div
                        v-if="!editingDisplayName"
                        class="inline-grid grid-cols-[max-content_2rem] items-start gap-1.5">
                        <button
                          type="button"
                          class="text-(--lab-text-primary) border-b border-transparent hover:border-(--lab-border) focus-visible:ring-(--lab-accent) inline-flex items-center bg-transparent py-0 text-left text-xl font-semibold outline-none transition-colors ring-0"
                          @click="editingDisplayName = true">
                          <span>{{ displayNameText }}</span>
                        </button>
                        <span class="h-8 w-8" aria-hidden="true"></span>
                      </div>
                      <div v-else class="min-w-fit" @focusout="onProfileEditorFocusout">
                        <div class="inline-grid grid-cols-[max-content_2rem] items-start gap-1.5">
                          <input
                            ref="displayNameInputRef"
                            id="account-display-name"
                            v-model="profileForm.display_name"
                            name="display_name"
                            type="text"
                            autocomplete="new-password"
                            autocapitalize="words"
                            spellcheck="false"
                            data-form-type="other"
                            :size="Math.max((profileForm.display_name || displayNameText).length, 1)"
                            class="text-(--lab-text-primary) border-b focus-visible:ring-(--lab-accent) appearance-none w-auto bg-transparent p-0 text-xl font-semibold outline-none ring-0"
                            :class="profileInputError ? 'text-(--lab-danger)' : ''"
                            placeholder="Имя профиля"
                            @keydown.enter.prevent="saveProfile"
                            @keydown.esc.prevent="cancelProfileEdit" />
                          <LabBaseButton
                            icon="ic:round-check"
                            icon-only
                            variant="secondary"
                            size="sm"
                            button-class="text-(--lab-text-secondary) hover:text-(--lab-info) focus-visible:text-(--lab-info) h-8 w-8"
                            title="Сохранить никнейм"
                            aria-label="Сохранить никнейм"
                            @click="saveProfile" />
                        </div>
                      </div>
                      <LabNotify
                        :text="profileError || profileInfo"
                        :tone="profileError ? 'error' : profileNoticeTone"
                        :temporary="profileError ? false : profileNoticeTemporary"
                        size="xs" />
                    </div>
                    <LabBaseButton
                      v-if="isAdmin"
                      :label="t('auth.account.admin')"
                      variant="info"
                      size="lg"
                      button-class="text-xs"
                      @click="goToAdmin" />
                  </div>
                  <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                    <span class="lab-text-muted shrink-0 text-xs uppercase tracking-wide">
                      {{ t('auth.account.email_label') }}
                    </span>
                    <span class="text-(--lab-text-primary) min-w-56 flex-1 basis-56 wrap-break-word text-sm">
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
                  :close-on-select="false">
                  <template #trigger="{ toggle }">
                    <LabBaseButton
                      variant="danger"
                      size="lg"
                      :label="t('auth.account.logout')"
                      button-class="text-xs"
                      @click="toggle" />
                  </template>
                  <LabBaseButton
                    variant="plain"
                    size="lg"
                    block
                    :label="t('auth.account.logout.current')"
                    button-class="justify-start text-xs"
                    @click="leave(false)" />
                  <LabBaseButton
                    variant="danger"
                    size="lg"
                    block
                    :label="t('auth.account.logout.all')"
                    button-class="justify-start text-xs"
                    @click="leave(true)" />
                </LabDropdown>
              </div>
            </section>
          </section>
          <section class="space-y-4 border bg-(--lab-bg-elevated) p-5">
            <div class="space-y-1">
              <h2 class="text-base font-medium text-(--lab-text-primary)">Оплата и доступ</h2>
              <p class="text-sm text-(--lab-text-muted)">Состояние платного доступа и последний заказ.</p>
            </div>

            <div v-if="accessLoading" class="rounded-2xl border bg-(--lab-bg-soft) p-4 text-sm text-(--lab-text-muted)">
              Загрузка сведений о доступе…
            </div>

            <template v-else>
              <div
                v-if="access?.has_active_access"
                class="rounded-2xl border border-(--lab-success)/30 bg-(--lab-success)/8 p-4">
                <div class="text-sm font-medium text-(--lab-text-primary)">Доступ активен</div>
                <p class="mt-1 text-sm text-(--lab-text-muted)">До {{ paymentAccessUntilText || '—' }}.</p>
              </div>

              <div v-else class="rounded-2xl border border-(--lab-border) bg-(--lab-bg-soft) p-4">
                <div class="text-sm font-medium text-(--lab-text-primary)">Активного доступа нет</div>
                <p class="mt-1 text-sm text-(--lab-text-muted)">
                  Можно перейти к странице оплаты и открыть доступ на месяц.
                </p>
              </div>

              <dl v-if="access?.latest_order" class="grid gap-3 sm:grid-cols-2">
                <div class="space-y-1">
                  <dt class="text-xs uppercase tracking-wide text-(--lab-text-muted)">Последний статус</dt>
                  <dd class="text-sm text-(--lab-text-primary)">
                    {{ paymentLatestOrderStatusLabel }}
                  </dd>
                </div>

                <div class="space-y-1">
                  <dt class="text-xs uppercase tracking-wide text-(--lab-text-muted)">Последняя сумма</dt>
                  <dd class="text-sm text-(--lab-text-primary)">
                    {{ paymentLatestOrderAmountText }}
                  </dd>
                </div>
              </dl>

              <div class="pt-1">
                <NuxtLink
                  to="/payments"
                  class="inline-flex min-h-11 items-center justify-center rounded-2xl border border-(--lab-border-strong) px-4 text-sm font-medium text-(--lab-text-primary)">
                  Перейти к оплате
                </NuxtLink>
              </div>
            </template>
          </section>
        </template>
        <template #panel-security>
          <section class="text-(--lab-text-primary) space-y-4">
            <article class="space-y-3">
              <div>
                <h2 class="text-xl sm:text-2xl">{{ t('auth.account.email.title') }}</h2>
                <p class="text-(--lab-text-muted) text-sm">
                  {{ t('auth.account.email.description', { email: user.email }) }}
                </p>
              </div>
              <div class="grid gap-3 sm:max-w-fit">
                <LabField
                  :label="t('auth.account.email.new_label')"
                  forId="account-next-email"
                  label-class="lab-text-muted text-xs normal-case tracking-normal">
                  <LabBaseInput
                    id="account-next-email"
                    v-model="emailChangeForm.email"
                    name="email"
                    type="email"
                    autocomplete="email"
                    input-class="w-full"
                    :placeholder="t('auth.account.email.new_placeholder')"
                    @keydown.enter.prevent="submitEmailChange" />
                </LabField>
              </div>
              <div class="flex flex-wrap items-center">
                <LabBaseButton
                  :label="t('auth.account.email.submit')"
                  variant="primary"
                  size="lg"
                  button-class="text-xs"
                  :disabled="emailChangePending"
                  @click="submitEmailChange" />
              </div>
              <LabNotify :text="emailChangeError" tone="error" size="xs" />
              <LabNotify :text="emailChangeInfo" tone="success" size="xs" />
            </article>
            <article class="space-y-3">
              <div>
                <h2 class="text-xl sm:text-2xl flex items-center flex-wrap gap-3">
                  {{ t('auth.account.password.title') }}
                  <span class="text-(--lab-text-muted) text-sm">{{ t('auth.account.password.description') }}</span>
                </h2>
              </div>
              <div class="grid gap-3 lg:grid-rows-2 sm:max-w-fit">
                <LabField
                  :label="t('auth.account.password.current_label')"
                  forId="account-current-password"
                  label-class="lab-text-muted text-xs normal-case tracking-normal">
                  <LabBaseInput
                    id="account-current-password"
                    v-model="passwordForm.current"
                    name="current_password"
                    type="password"
                    autocomplete="current-password"
                    :input-class="['w-full', passwordCurrentInputError ? 'text-(--lab-danger)' : '']"
                    :placeholder="t('auth.account.password.current_placeholder')" />
                </LabField>
                <LabField
                  :label="t('auth.account.password.next_label')"
                  forId="account-next-password"
                  label-class="lab-text-muted text-xs normal-case tracking-normal">
                  <LabBaseInput
                    id="account-next-password"
                    v-model="passwordForm.next"
                    name="new_password"
                    type="password"
                    autocomplete="new-password"
                    input-class="w-full"
                    :placeholder="t('auth.account.password.next_placeholder')" />
                </LabField>
              </div>
              <div class="flex flex-wrap items-center">
                <LabBaseButton
                  :label="t('auth.account.password.submit')"
                  variant="primary"
                  size="lg"
                  button-class="text-xs"
                  @click="submitPasswordChange" />
              </div>
              <LabNotify :text="passwordError" tone="error" size="xs" />
              <LabNotify :text="passwordInfo" tone="success" size="xs" />
            </article>
            <article class="space-y-3">
              <div class="space-y-3">
                <h2 class="text-xl sm:text-2xl flex items-center flex-wrap gap-3">
                  2FA
                  <span class="text-(--lab-text-muted) text-sm">
                    {{ twofaStatusLabel }}
                  </span>
                </h2>

                <LabBaseButton
                  :variant="user.is_two_factor_enabled ? 'danger' : 'success'"
                  size="lg"
                  button-class="text-xs"
                  @click="user.is_two_factor_enabled ? toggleDisable2faForm() : begin2faSetup()">
                  {{
                    user.is_two_factor_enabled
                      ? showDisable2faForm
                        ? t('auth.account.twofa.hide_form')
                        : t('auth.account.twofa.disable')
                      : t('auth.account.twofa.get_secret')
                  }}
                </LabBaseButton>
              </div>
              <LabNotify :text="twofaError" tone="error" />
              <LabNotify :text="twofaInfo" tone="success" :temporary="false" />
              <div
                v-if="!user.is_two_factor_enabled && setupSecret"
                class="grid items-start gap-3 xl:grid-cols-[auto_minmax(0,1fr)]">
                <div v-if="setupQrDataUrl">
                  <img
                    :src="setupQrDataUrl"
                    alt="QR-код для настройки 2FA"
                    class="h-48 w-48 object-contain"
                    loading="lazy" />
                </div>
                <div class="space-y-3">
                  <p v-if="setupOtpAuthUrl" class="text-(--lab-text-muted) wrap-break-word text-xs leading-5">
                    {{ t('auth.account.twofa.qr_fallback').split('{link}')[0] }}
                    <a
                      :href="setupOtpAuthUrl"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="text-(--lab-info) ring-1 ring-transparent focus-visible:outline-none focus-visible:ring-2">
                      {{ t('auth.account.twofa.qr_fallback_link') }}
                    </a>
                    {{ t('auth.account.twofa.qr_fallback').split('{link}')[1] }}
                  </p>
                  <LabCopyBlock
                    :label="t('auth.account.twofa.copy_key_label')"
                    :value="setupSecret"
                    variant="dark-cyan"
                    button-class="w-full"
                    :title-idle="t('auth.account.twofa.copy_key_idle')"
                    :title-success="t('auth.account.twofa.copy_key_success')"
                    :title-error="t('auth.account.twofa.copy_key_error')"
                    :show-state-tooltip="true" />
                  <AuthCodeInput
                    id="account-enable-2fa-code"
                    ref="enableCodeInputRef"
                    :model-value="enableCode"
                    name="one_time_code_enable"
                    :label="t('auth.account.twofa.code_label')"
                    :hint="t('auth.account.twofa.code_hint')"
                    :invalid="Boolean(twofaError)"
                    :valid="Boolean(twofaInfo) && !Boolean(twofaError)"
                    @update:model-value="onEnableCodeInput"
                    @complete="activate2fa" />
                </div>
              </div>
              <div v-if="user.is_two_factor_enabled && showDisable2faForm" class="space-y-3">
                <input
                  type="text"
                  name="account_disable_2fa_username"
                  autocomplete="username"
                  tabindex="-1"
                  aria-hidden="true"
                  class="pointer-events-none absolute -left-96 h-px w-px opacity-0" />
                <LabField
                  label="Текущий пароль"
                  forId="account-disable-2fa-password"
                  field-class="min-w-0"
                  label-class="lab-text-muted text-xs normal-case tracking-normal">
                  <LabBaseInput
                    id="account-disable-2fa-password"
                    v-model="disablePassword"
                    name="disable_2fa_password"
                    type="password"
                    autocomplete="current-password"
                    input-class="w-full"
                    placeholder="Текущий пароль" />
                </LabField>
                <LabNavTabs
                  v-model="disableMethodTab"
                  :items="disableCodeTabItems"
                  route-query-key="disable2fa"
                  route-default-value="totp"
                  panel-class="space-y-3">
                  <template #panel-totp>
                    <AuthCodeInput
                      id="account-disable-2fa-totp"
                      ref="disableTotpInputRef"
                      :model-value="disableTotpCode"
                      name="disable_totp_code"
                      label="Код из приложения"
                      hint="Отключение начнётся автоматически после ввода 6 цифр."
                      :invalid="Boolean(twofaError)"
                      :valid="Boolean(twofaInfo) && !Boolean(twofaError)"
                      @update:model-value="onDisableTotpCodeInput"
                      @complete="deactivate2fa" />
                  </template>
                  <template #panel-backup>
                    <AuthRecoveryCodeInput
                      id="account-disable-2fa-backup"
                      ref="disableBackupCodeInputRef"
                      :model-value="disableBackupCode"
                      name="disable_backup_code"
                      label="Код сброса"
                      hint="Отключение начнётся автоматически после ввода 8 символов."
                      :invalid="Boolean(twofaError)"
                      :valid="Boolean(twofaInfo) && !Boolean(twofaError)"
                      @update:model-value="onDisableBackupCodeInput"
                      @complete="deactivate2fa(`${$event.slice(0, 4)}-${$event.slice(4, 8)}`)" />
                  </template>
                </LabNavTabs>
              </div>
              <div v-if="backupCodes.length" class="text-(--lab-info) space-y-3">
                <div>
                  <p class="text-sm font-medium">Резервные коды</p>
                  <p class="text-xs leading-5">Сохрани все 10 кодов. Каждый код сработает только один раз.</p>
                </div>
                <div class="grid grid-cols-2 gap-2 sm:grid-cols-3">
                  <div
                    v-for="code in backupCodes"
                    :key="code"
                    class="inline-flex min-h-11 items-center justify-center text-center font-mono text-sm tracking-widest">
                    {{ code }}
                  </div>
                </div>
                <p class="text-xs leading-5">Храни их вне браузера и не передавай через мессенджеры.</p>
              </div>
            </article>
            <article class="space-y-3">
              <div>
                <h2 class="text-xl sm:text-2xl">{{ t('auth.account.activity.title') }}</h2>
                <p class="text-(--lab-text-muted) text-sm">{{ t('auth.account.activity.description') }}</p>
              </div>
              <div class="space-y-2">
                <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                  <span class="lab-text-muted shrink-0 text-xs uppercase tracking-wide">
                    {{ t('auth.account.activity.last_login') }}
                  </span>
                  <span class="text-(--lab-text-primary) min-w-56 flex-1 basis-56 text-sm">
                    {{ formatDateTime(user.last_login_at) }}
                  </span>
                </div>
                <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                  <span class="lab-text-muted shrink-0 text-xs uppercase tracking-wide">
                    {{ t('auth.account.activity.created_at') }}
                  </span>
                  <span class="text-(--lab-text-primary) min-w-56 flex-1 basis-56 text-sm">
                    {{ formatDateTime(user.created_at) }}
                  </span>
                </div>
              </div>
              <LabNotify :text="actionError" tone="error" size="xs" />
              <LabNotify :text="actionInfo" tone="success" size="xs" />
              <p v-if="activityLoading" class="text-(--lab-text-muted) text-xs">Загрузка сессий и активности…</p>
              <LabNotify v-else-if="activityError" :text="activityError" tone="error" size="xs" />
              <LabNavTabs
                v-else
                v-model="activityTab"
                :items="activityTabItems"
                route-query-key="activity"
                route-default-value="sessions">
                <template #tab="{ item }">
                  <span>{{ item.label }}</span>
                  <span class="ml-1 opacity-80">
                    {{
                      item.value === 'sessions'
                        ? groupedSessions.length
                        : item.value === 'attempts'
                          ? loginAttempts.length
                          : securityEvents.length
                    }}
                  </span>
                </template>
                <template #panel-sessions>
                  <div class="max-h-96 overflow-y-auto space-y-3">
                    <article
                      v-for="item in groupedSessions"
                      :key="item.key"
                      class="flex flex-wrap items-center justify-between gap-3">
                      <div>
                        <p class="text-(--lab-text-primary) text-sm">{{ item.deviceLabel }}</p>
                        <p class="text-(--lab-text-muted) text-xs">
                          {{ item.ip }} ·
                          {{
                            item.mfaVerified
                              ? t('auth.account.activity.mfa_verified')
                              : t('auth.account.activity.mfa_unverified')
                          }}
                        </p>
                        <p class="text-(--lab-text-muted) text-xs">
                          {{ t('auth.account.activity.sessions_count', { count: item.count }) }} ·
                          {{ t('auth.account.activity.last_activity') }}
                          <LabRelativeTime :datetime="item.lastSeenAt" compact />
                        </p>
                      </div>
                      <LabBaseButton
                        variant="secondary"
                        size="lg"
                        button-class="text-xs"
                        :disabled="item.revokableSessionIds.length === 0"
                        @click="revokeSessionGroup(item)">
                        {{
                          item.revokableSessionIds.length === 0
                            ? t('auth.account.activity.current')
                            : item.hasCurrent
                              ? t('auth.account.activity.revoke_and_logout')
                              : t('auth.account.activity.revoke')
                        }}
                      </LabBaseButton>
                    </article>
                  </div>
                </template>
                <template #panel-attempts>
                  <div class="max-h-96 overflow-y-auto space-y-3">
                    <article v-for="item in loginAttempts" :key="item.attempt_id">
                      <div class="flex flex-wrap items-start justify-between gap-2">
                        <p class="text-(--lab-text-primary) text-sm">{{ item.outcome }}</p>
                        <p class="text-(--lab-text-muted) text-xs">
                          <LabRelativeTime :datetime="item.created_at" compact />
                        </p>
                      </div>
                      <p class="text-(--lab-text-muted) text-xs">
                        {{ item.ip || '—' }} · оценка риска {{ item.risk_score }}
                      </p>
                      <p class="text-(--lab-text-muted) text-xs">
                        {{ item.failure_reason || item.suspicious_reason || item.user_agent || '—' }}
                      </p>
                    </article>
                  </div>
                </template>
                <template #panel-events>
                  <div class="max-h-96 overflow-y-auto space-y-3">
                    <article v-for="item in securityEvents" :key="item.event_id">
                      <div class="flex flex-wrap items-start justify-between gap-2">
                        <p class="text-(--lab-text-primary) text-sm">
                          {{ item.event_type }}
                          <span class="text-(--lab-text-muted)">· {{ item.severity }}</span>
                        </p>
                        <p class="text-(--lab-text-muted) text-xs">
                          <LabRelativeTime :datetime="item.created_at" compact />
                        </p>
                      </div>
                      <p class="text-(--lab-text-muted) text-xs">{{ item.ip || '—' }}</p>
                      <p class="text-(--lab-text-muted) text-xs wrap-break-word">
                        {{ JSON.stringify(item.payload || {}) }}
                      </p>
                    </article>
                  </div>
                </template>
              </LabNavTabs>
            </article>
          </section>
        </template>
      </LabNavTabs>
      <LabCropperModal
        v-if="avatarCropDialog.open && avatarCropDialog.file"
        :file="avatarCropDialog.file"
        title="Кадрирование аватарки"
        :loading="avatarUploading"
        aspect-preset="1:1"
        :aspect-locked="true"
        @cancel="closeAvatarCropDialog"
        @confirm="onAvatarCropConfirm" />
      <LabViewerImage
        v-model="avatarPreviewDialog.open"
        :items="avatarViewerItems"
        :initial-index="avatarPreviewDialog.index"
        @active-index-change="onAvatarPreviewIndexChange"
        @update:model-value="onAvatarPreviewModelUpdate">
        <template #toolbar-extra>
          <div v-if="activePreviewAvatarItem" class="flex items-center gap-2">
            <LabBaseButton
              :variant="isPreviewAvatarPrimary ? 'success' : 'primary'"
              size="sm"
              :disabled="avatarUploading"
              :aria-label="isPreviewAvatarPrimary ? 'Основная фотография профиля' : 'Сделать фотографию основной'"
              :label="isPreviewAvatarPrimary ? 'Основная фотография' : 'Сделать основной'"
              @click.stop="makePreviewAvatarPrimary" />
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
              @confirm="removeAvatarItem(activePreviewAvatarItem.id)" />
          </div>
        </template>
      </LabViewerImage>
    </template>
  </div>
</template>
