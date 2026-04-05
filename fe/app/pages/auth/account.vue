<script setup lang="ts">
  const title = 'Аккаунт'
  usePageSeo({
    title,
    description: 'Профиль аккаунта, настройки, выход.'
  })
  const router = useRouter()
  const runtimeConfig = useRuntimeConfig()
  const uiPreferences = useUiPreferencesStore()
  const { t, readThemePreferenceFromSettings } = useInterfacePreferences()
  const { saveInterfacePreferences: syncInterfacePreferences } = useInterfacePreferencesSync()
  const {
    user,
    isAdmin,
    loaded,
    ensureLoaded,
    loadMe,
    logout,
    logoutAll,
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
  await ensureLoaded()
  const sessions = ref<AuthSessionView[]>([])
  const loginAttempts = ref<AuthLoginAttemptView[]>([])
  const securityEvents = ref<AuthSecurityEventView[]>([])
  const initialLoading = ref(!loaded.value)
  const activityTab = ref<'sessions' | 'attempts' | 'events'>('sessions')
  const securityTab = ref<'password' | 'twofa' | 'activity'>('password')
  const activityLoading = ref(false)
  const activityLoaded = ref(false)
  const activityError = ref('')
  const actionError = ref('')
  const actionInfo = ref('')
  const twofaError = ref('')
  const twofaInfo = ref('')
  const setupSecret = ref('')
  const setupOtpAuthUrl = ref('')
  const setupQrDataUrl = ref('')
  const enableCode = ref('')
  const enableCodeInputRef = ref<{ focus: () => void } | null>(null)
  const disablePassword = ref('')
  const disableMethodTab = ref<'totp' | 'backup'>('totp')
  const disableTotpCode = ref('')
  const disableBackupCode = ref('')
  const disableTotpInputRef = ref<{ focus: () => void } | null>(null)
  const disableBackupCodeInputRef = ref<{ focus: () => void } | null>(null)
  const showDisable2faForm = ref(false)
  const backupCodes = ref<string[]>([])
  const securityTabItems: LabTabItem[] = [
    { value: 'password', label: 'Изменить пароль' },
    { value: 'twofa', label: '2FA' },
    { value: 'activity', label: 'Активность' }
  ]
  const activityTabItems: LabTabItem[] = [
    { value: 'sessions', label: 'Сессии' },
    { value: 'attempts', label: 'Попытки входа' },
    { value: 'events', label: 'Журнал безопасности' }
  ]
  const disableCodeTabItems: LabTabItem[] = [
    { value: 'totp', label: 'Код из приложения' },
    { value: 'backup', label: 'Код сброса' }
  ]
  const passwordForm = reactive({
    current: '',
    next: ''
  })
  const passwordError = ref('')
  const passwordInfo = ref('')
  const passwordCurrentInputError = ref(false)
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
  const profileError = ref('')
  const profileInfo = ref('')
  const interfaceError = ref('')
  const interfaceInfo = ref('')
  const suppressInterfaceAutosave = ref(false)
  const interfaceAutosavePending = ref(false)
  const profileInputError = ref(false)
  const editingDisplayName = ref(false)
  const displayNameInputRef = ref<{ focus: () => void } | null>(null)
  const profileNameEditorRef = ref<HTMLElement | null>(null)
  const showLogoutActions = ref(false)
  const logoutMenuRef = ref<HTMLElement | null>(null)
  const avatarCropDialog = reactive<{
    open: boolean
    file: File | null
  }>({
    open: false,
    file: null
  })
  const avatarPreviewDialog = reactive({
    open: false,
    index: 0,
    title: 'Фотографии профиля'
  })
  const avatarUploading = ref(false)
  const avatarError = ref('')
  const avatarInfo = ref('')
  const SUCCESS_NOTIFY_TIMEOUT_MS = 5000
  const successNoticeTimers = new Map<string, ReturnType<typeof setTimeout>>()
  const avatarImageMaxBytes = Number(runtimeConfig.public.mediaImageMaxBytes || 8388608)
  const AVATAR_IMAGE_MAX_MB = Math.max(1, Math.ceil(avatarImageMaxBytes / (1024 * 1024)))
  const AVATAR_IMAGE_MAX_BYTES = avatarImageMaxBytes
  const AVATAR_IMAGE_ACCEPT = 'image/jpeg,image/png,image/webp,image/gif,image/avif'
  const avatarGallery = computed(() => getAuthAvatarGallery(user.value))
  const avatarGalleryItems = computed(() => avatarGallery.value.items)
  const avatarPrimaryItemId = computed(() => avatarGallery.value.primaryId)
  const avatar = computed(() => getAuthAvatarUrls(user.value))
  const avatarViewerItems = computed<ImageViewerItem[]>(() =>
    avatarGalleryItems.value.map((item, index) => ({
      src: buildMediaFileUrl(item.original_image_key || item.profile_image_key || item.icon_image_key),
      thumbnailSrc: buildMediaFileUrl(item.profile_image_key || item.icon_image_key || item.original_image_key),
      title: avatarGalleryItems.value.length > 1 ? `Фотография профиля ${index + 1}` : 'Фотография профиля',
      alt: avatarGalleryItems.value.length > 1 ? `Фотография профиля ${index + 1}` : 'Фотография профиля'
    }))
  )
  const activePreviewAvatarItem = computed(() => avatarGalleryItems.value[avatarPreviewDialog.index] || null)
  const isPreviewAvatarPrimary = computed(() => activePreviewAvatarItem.value?.id === avatarPrimaryItemId.value)
  const displayNameText = computed(() => profileForm.display_name.trim() || 'Имя профиля')
  const formatDateTime = (value?: string | null) => {
    const date = value ? new Date(value) : null
    if (!date || Number.isNaN(date.getTime())) return '—'
    return new Intl.DateTimeFormat('ru-RU', {
      dateStyle: 'medium',
      timeStyle: 'short'
    }).format(date)
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
  const resolveUploadedImageKey = (res: any) => {
    return String(res?.data?.image_key || res?.image_key || '').trim()
  }
  const closeAvatarCropDialog = (force = false) => {
    if (avatarUploading.value && !force) return
    avatarCropDialog.open = false
    avatarCropDialog.file = null
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
  const findAvatarGalleryIndexById = (itemId?: string | null) => {
    const normalizedId = String(itemId || '').trim()
    if (!normalizedId) return 0
    const foundIndex = avatarGalleryItems.value.findIndex(item => item.id === normalizedId)
    return foundIndex >= 0 ? foundIndex : 0
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
  const openAvatarPreview = (index = findAvatarGalleryIndexById(avatarPrimaryItemId.value)) => {
    if (!avatarViewerItems.value.length) return
    avatarPreviewDialog.index = Math.min(Math.max(index, 0), avatarViewerItems.value.length - 1)
    avatarPreviewDialog.open = true
    avatarPreviewDialog.title = 'Фотографии профиля'
  }
  const onAvatarPreviewModelUpdate = (value: boolean) => {
    avatarPreviewDialog.open = value
  }
  const onAvatarPreviewIndexChange = (index: number) => {
    avatarPreviewDialog.index = index
  }
  const makePreviewAvatarPrimary = async () => {
    const itemId = activePreviewAvatarItem.value?.id || ''
    if (!itemId) return
    await setPrimaryAvatar(itemId)
  }
  const syncProfileForm = () => {
    if (!user.value) return
    suppressInterfaceAutosave.value = true
    profileForm.display_name = user.value.display_name || ''
    interfaceForm.locale = normalizeInterfaceLocale(user.value.locale || uiPreferences.interfaceLocale)
    interfaceForm.theme_preference = readThemePreferenceFromSettings(user.value.settings)
    nextTick(() => {
      suppressInterfaceAutosave.value = false
    })
  }
  syncProfileForm()
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
  const onWindowPointerDown = (event: PointerEvent) => {
    const target = event.target as Node | null
    if (!target) return
    if (editingDisplayName.value && !profileNameEditorRef.value?.contains(target)) {
      editingDisplayName.value = false
    }
    if (showLogoutActions.value && !logoutMenuRef.value?.contains(target)) {
      showLogoutActions.value = false
    }
  }
  watch(securityTab, tab => {
    if (tab !== 'activity') return
    void loadActivityState()
  })
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
  const saveProfile = async () => {
    profileError.value = ''
    clearSuccessNotice('profile', profileInfo)
    try {
      await updateProfile({
        display_name: profileForm.display_name
      })
      editingDisplayName.value = false
      profileInputError.value = false
      showSuccessNotice('profile', profileInfo, 'Ник обновлён.')
      syncProfileForm()
    } catch (err: any) {
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
      showSuccessNotice(
        'twofa',
        twofaInfo,
        'Ключ для двухфакторной защиты создан. Добавь его в приложение и введи код подтверждения.'
      )
    } catch (err: any) {
      twofaError.value = extractApiErrorMessage(err, 'Не удалось начать настройку 2FA.')
    }
  }
  const activate2fa = async () => {
    twofaError.value = ''
    clearSuccessNotice('twofa', twofaInfo)
    enableCode.value = enableCode.value.replace(/\D+/g, '').slice(0, 6)
    if (enableCode.value.length !== 6) {
      twofaError.value = 'Введите 6 цифр кода подтверждения 2FA.'
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
  const onEnableCodeInput = (value: string) => {
    enableCode.value = value
    if (twofaError.value) {
      twofaError.value = ''
    }
  }
  const focusDisableCodeInput = () => {
    if (disableMethodTab.value === 'backup') {
      disableBackupCodeInputRef.value?.focus()
      return
    }
    disableTotpInputRef.value?.focus()
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
</script>
<template>
  <div>
    <LabNavHeader :title class="min-w-0 flex-1" />
    <LabNotify :text="actionError" tone="error" size="xs" />
    <LabNotify :text="actionInfo" tone="success" size="xs" />
    <section>
      <LabLoader v-if="initialLoading" variant="inline" label="Загрузка аккаунта…" />
      <AuthFeatureGateNotice
        v-else-if="!user"
        message="Войдите в аккаунт, чтобы управлять профилем, сессиями и параметрами безопасности." />
      <template v-else-if="user">
        <section class="text-(--lab-text-primary) relative overflow-hidden">
          <div class="flex min-w-0 flex-col">
            <div class="flex min-w-0 flex-wrap items-start justify-between">
              <div ref="profileNameEditorRef" class="relative h-12 min-w-0 flex-1 basis-72">
                <button
                  v-if="!editingDisplayName"
                  type="button"
                  class="text-(--lab-text-primary) focus-visible:ring-(--lab-accent) absolute inset-0 inline-flex h-full w-full items-center bg-transparent text-left text-xl font-semibold outline-none transition-colors ring-0"
                  @click="editingDisplayName = true">
                  <span class="truncate">{{ displayNameText }}</span>
                </button>
                <div v-else class="absolute inset-0 flex items-center">
                  <LabBaseInput
                    ref="displayNameInputRef"
                    id="account-display-name"
                    v-model="profileForm.display_name"
                    name="display_name"
                    type="text"
                    autocomplete="nickname"
                    autocapitalize="words"
                    spellcheck="false"
                    :input-class="[
                      'h-full min-h-0 w-full text-lg font-semibold leading-tight outline-none',
                      profileInputError ? 'text-(--lab-danger)' : ''
                    ]"
                    placeholder="Имя профиля"
                    @keydown.enter.prevent="saveProfile" />
                  <LabBaseButton
                    variant="secondary"
                    size="sm"
                    button-class="shrink-0"
                    icon="ic:round-save"
                    icon-only
                    icon-class="h-5 w-5"
                    @click="saveProfile"></LabBaseButton>
                </div>
              </div>
              <div class="flex shrink-0 flex-wrap items-center justify-end">
                <LabBaseButton
                  v-if="isAdmin"
                  label="Админка"
                  variant="info"
                  size="lg"
                  button-class="text-xs"
                  @click="goToAdmin" />
                <div ref="logoutMenuRef" class="relative shrink-0">
                  <LabBaseButton
                    variant="danger"
                    size="lg"
                    label="Выйти"
                    button-class="text-xs"
                    @click="showLogoutActions = !showLogoutActions" />
                  <div v-if="showLogoutActions" class="absolute right-0 top-full z-20 min-w-56">
                    <LabBaseButton
                      variant="plain"
                      size="lg"
                      block
                      label="На этом устройстве"
                      button-class="justify-start text-xs"
                      @click="leave(false)" />
                    <LabBaseButton
                      variant="danger"
                      size="lg"
                      block
                      label="На всех устройствах"
                      button-class="justify-start text-xs"
                      @click="leave(true)" />
                  </div>
                </div>
              </div>
            </div>
            <LabNotify :text="profileError || profileInfo" :tone="profileError ? 'error' : 'success'" size="xs" />
          </div>
          <div class="flex flex-col lg:flex-row lg:items-start">
            <div class="flex shrink-0 flex-col items-start">
              <div class="relative h-36 w-36 overflow-hidden sm:h-40 sm:w-40">
                <img
                  v-if="avatar.profileImageUrl"
                  :src="avatar.profileImageUrl"
                  alt="Аватар профиля"
                  class="h-full w-full cursor-zoom-in object-cover"
                  @click="openAvatarPreview()" />
                <div v-else class="text-(--lab-text-muted) flex h-full w-full items-center justify-center">
                  <Icon name="ic:round-account-circle" class="h-20 w-20 text-6xl" />
                </div>
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
              <div v-if="avatarGalleryItems.length > 1" class="w-full max-w-40 sm:max-w-44">
                <div class="flex items-center justify-between">
                  <p class="lab-text-muted text-xs uppercase tracking-wide">Фотографии профиля</p>
                </div>
                <div class="grid w-full grid-cols-3">
                  <div v-for="(item, index) in avatarGalleryItems" :key="item.id" class="relative overflow-hidden">
                    <button
                      type="button"
                      class="block aspect-square w-full overflow-hidden"
                      :aria-label="`Открыть фотографию профиля ${index + 1}`"
                      @click="openAvatarPreview(index)">
                      <img
                        :src="
                          buildMediaFileUrl(item.profile_image_key || item.icon_image_key || item.original_image_key)
                        "
                        :alt="`Фотография профиля ${index + 1}`"
                        class="h-full w-full object-cover" />
                    </button>
                  </div>
                </div>
              </div>
              <LabNotify :text="avatarError" tone="error" size="xs" />
              <LabNotify :text="avatarInfo" tone="success" size="xs" />
            </div>
            <div class="min-w-0 flex-1">
              <article>
                <div>
                  <div class="flex flex-wrap items-baseline">
                    <span class="lab-text-muted basis-full text-xs uppercase tracking-wide sm:basis-36">Email</span>
                    <span class="text-(--lab-text-primary) min-w-0 flex-1 basis-64 wrap-break-word text-sm">
                      {{ user.email }}
                    </span>
                  </div>
                  <div class="flex flex-wrap items-baseline">
                    <span class="lab-text-muted basis-full text-xs uppercase tracking-wide sm:basis-36">Статус</span>
                    <span class="text-(--lab-text-primary) min-w-0 flex-1 basis-64 text-sm">
                      {{ user.status }}
                      <span v-if="isAdmin" class="text-(--lab-text-muted)">· администратор</span>
                    </span>
                  </div>
                  <div class="flex flex-wrap items-baseline">
                    <span class="lab-text-muted basis-full text-xs uppercase tracking-wide sm:basis-36">Роли</span>
                    <span class="text-(--lab-text-primary) min-w-0 flex-1 basis-64 wrap-break-word text-sm">
                      {{ user.roles.length ? user.roles.join(', ') : '—' }}
                    </span>
                  </div>
                  <div class="flex flex-wrap items-baseline">
                    <span class="lab-text-muted basis-full text-xs uppercase tracking-wide sm:basis-36">
                      Последний вход
                    </span>
                    <span class="text-(--lab-text-primary) min-w-0 flex-1 basis-64 text-sm">
                      {{ formatDateTime(user.last_login_at) }}
                    </span>
                  </div>
                  <div class="flex flex-wrap items-baseline">
                    <span class="lab-text-muted basis-full text-xs uppercase tracking-wide sm:basis-36">
                      Регистрация
                    </span>
                    <span class="text-(--lab-text-primary) min-w-0 flex-1 basis-64 text-sm">
                      {{ formatDateTime(user.created_at) }}
                    </span>
                  </div>
                </div>
              </article>
              <article>
                <LabNotify :text="interfaceError" tone="error" size="xs" />
                <LabNotify :text="interfaceInfo" tone="success" size="xs" />
              </article>
            </div>
          </div>
        </section>
        <section class="text-(--lab-text-primary) relative overflow-hidden">
          <div>
            <h2 class="lab-text-primary text-xl sm:text-2xl">Безопасность</h2>
            <p class="text-(--lab-text-muted) text-sm">
              {{ user.is_two_factor_enabled ? '2FA включена' : '2FA выключена' }}
            </p>
          </div>
          <LabNavTabs
            v-model="securityTab"
            :items="securityTabItems"
            :no-select="true"
            list-class="flex flex-wrap"
            route-query-key="security"
            route-default-value="password"
            button-class="inline-flex min-h-10 items-center justify-center text-xs transition-colors"
            active-class="text-(--lab-accent)"
            inactive-class="text-(--lab-text-secondary) hover:text-(--lab-text-primary)"
            panel-class="">
            <template #panel-password>
              <section>
                <div class="grid lg:grid-cols-2">
                  <LabField
                    label="Текущий пароль"
                    forId="account-current-password"
                    label-class="lab-text-muted text-xs normal-case tracking-normal">
                    <LabBaseInput
                      id="account-current-password"
                      v-model="passwordForm.current"
                      name="current_password"
                      type="password"
                      autocomplete="current-password"
                      :input-class="['w-full', passwordCurrentInputError ? 'text-(--lab-danger)' : '']"
                      placeholder="Текущий пароль" />
                  </LabField>
                  <LabField
                    label="Новый пароль"
                    forId="account-next-password"
                    label-class="lab-text-muted text-xs normal-case tracking-normal">
                    <LabBaseInput
                      id="account-next-password"
                      v-model="passwordForm.next"
                      name="new_password"
                      type="password"
                      autocomplete="new-password"
                      input-class="w-full"
                      placeholder="Новый пароль" />
                  </LabField>
                </div>
                <div class="flex flex-wrap items-center">
                  <LabBaseButton
                    label="Изменить пароль"
                    variant="primary"
                    size="lg"
                    button-class="text-xs"
                    @click="submitPasswordChange" />
                </div>
                <div>
                  <LabNotify :text="passwordError" tone="error" size="xs" />
                  <LabNotify :text="passwordInfo" tone="success" size="xs" />
                </div>
              </section>
            </template>
            <template #panel-twofa>
              <section>
                <div class="flex flex-wrap items-center justify-between">
                  <p class="text-(--lab-text-primary) text-sm">
                    2FA {{ user.is_two_factor_enabled ? 'включена' : 'выключена' }}
                  </p>
                  <LabBaseButton
                    :variant="user.is_two_factor_enabled ? 'danger' : 'success'"
                    size="lg"
                    button-class="text-xs"
                    @click="user.is_two_factor_enabled ? toggleDisable2faForm() : begin2faSetup()">
                    {{
                      user.is_two_factor_enabled
                        ? showDisable2faForm
                          ? 'Скрыть форму'
                          : 'Отключить 2FA'
                        : 'Получить секрет'
                    }}
                  </LabBaseButton>
                </div>
                <div>
                  <LabNotify :text="twofaError" tone="error" />
                  <LabNotify :text="twofaInfo" tone="success" />
                </div>
                <div v-if="!user.is_two_factor_enabled && setupSecret" class="grid xl:grid-cols-2">
                  <div v-if="setupQrDataUrl">
                    <img
                      :src="setupQrDataUrl"
                      alt="QR-код для настройки 2FA"
                      class="h-48 w-48 object-contain"
                      loading="lazy" />
                  </div>
                  <div>
                    <p v-if="setupOtpAuthUrl" class="text-(--lab-text-muted) wrap-break-word text-xs leading-5">
                      Если QR не сканируется, используй ключ ниже или otpauth URL: {{ setupOtpAuthUrl }}
                    </p>
                    <LabCopyBlock
                      label="Ключ 2FA"
                      :value="setupSecret"
                      variant="dark-cyan"
                      button-class="w-full"
                      title-idle="Нажмите, чтобы скопировать ключ 2FA"
                      title-success="Ключ 2FA скопирован в буфер обмена."
                      title-error="Ошибка копирования ключа 2FA"
                      :show-state-tooltip="true" />
                    <AuthCodeInput
                      id="account-enable-2fa-code"
                      ref="enableCodeInputRef"
                      :model-value="enableCode"
                      name="one_time_code_enable"
                      label="Код из приложения"
                      hint="Проверка начнётся автоматически после ввода 6 цифр."
                      :invalid="Boolean(twofaError)"
                      :valid="Boolean(twofaInfo) && !Boolean(twofaError)"
                      @update:model-value="onEnableCodeInput"
                      @complete="activate2fa" />
                  </div>
                </div>
                <div v-if="user.is_two_factor_enabled && showDisable2faForm" class="">
                  <input
                    type="text"
                    name="account_disable_2fa_username"
                    autocomplete="username"
                    tabindex="-1"
                    aria-hidden="true"
                    class="pointer-events-none absolute -left-96 h-px w-px opacity-0" />
                  <div>
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
                      list-class="flex flex-wrap"
                      button-class="inline-flex min-h-10 items-center justify-center text-xs transition-colors"
                      active-class="text-(--lab-danger)"
                      inactive-class="text-(--lab-text-secondary) hover:text-(--lab-text-primary)"
                      panel-class="">
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
                </div>
                <div v-if="backupCodes.length" class="text-(--lab-info)">
                  <div>
                    <p class="text-sm font-medium">Резервные коды</p>
                    <p class="text-xs leading-5">Сохрани все 10 кодов. Каждый код сработает только один раз.</p>
                  </div>
                  <div class="grid grid-cols-2 sm:grid-cols-3">
                    <div
                      v-for="code in backupCodes"
                      :key="code"
                      class="inline-flex min-h-11 items-center justify-center text-center font-mono text-sm tracking-widest">
                      {{ code }}
                    </div>
                  </div>
                  <p class="text-xs leading-5">Храни их вне браузера и не передавай через мессенджеры.</p>
                </div>
              </section>
            </template>
            <template #panel-activity>
              <section>
                <p v-if="activityLoading" class="text-(--lab-text-muted) text-xs">Загрузка сессий и активности…</p>
                <template v-else-if="activityError">
                  <LabNotify :text="activityError" tone="error" size="xs" />
                </template>
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
                    <div class="max-h-96 overflow-y-auto">
                      <article
                        v-for="item in groupedSessions"
                        :key="item.key"
                        class="flex flex-wrap items-center justify-between">
                        <div>
                          <p class="text-(--lab-text-primary) text-sm">{{ item.deviceLabel }}</p>
                          <p class="text-(--lab-text-muted) text-xs">
                            {{ item.ip }} · 2FA {{ item.mfaVerified ? 'подтверждена' : 'не подтверждена' }}
                          </p>
                          <p class="text-(--lab-text-muted) text-xs">
                            Сессий: {{ item.count }} · Последняя активность:
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
                              ? 'Текущая'
                              : item.hasCurrent
                                ? 'Отозвать и выйти'
                                : 'Отозвать'
                          }}
                        </LabBaseButton>
                      </article>
                    </div>
                  </template>
                  <template #panel-attempts>
                    <div class="max-h-96 overflow-y-auto">
                      <article v-for="item in loginAttempts" :key="item.attempt_id" class="">
                        <div class="flex flex-wrap items-start justify-between">
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
                    <div class="max-h-96 overflow-y-auto">
                      <article v-for="item in securityEvents" :key="item.event_id" class="">
                        <div class="flex flex-wrap items-start justify-between">
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
              </section>
            </template>
          </LabNavTabs>
        </section>
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
          :title="avatarPreviewDialog.title"
          @active-index-change="onAvatarPreviewIndexChange"
          @update:model-value="onAvatarPreviewModelUpdate">
          <template #toolbar-extra>
            <div v-if="activePreviewAvatarItem" class="flex items-center">
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
      <template v-else>
        <p class="text-xs">Активная сессия не найдена.</p>
        <AuthLinks variant="inline" />
      </template>
    </section>
  </div>
</template>
