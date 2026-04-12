<script setup lang="ts">
  const { t } = useInterfacePreferences()
  const router = useRouter()
  const { formatAbsoluteDateTime } = useLocalizedDateTime()
  const {
    user,
    ensureLoaded,
    loadMe,
    logout,
    requestEmailChange,
    listSessions,
    listLoginAttempts,
    listSecurityEvents,
    revokeSession,
    changePassword,
    setupTwoFactor,
    enableTwoFactor,
    disableTwoFactor
  } = useAuth()

  const activityTabItems: LabTabItem[] = [
    { value: 'sessions', label: t('auth.account.activity.sessions') },
    { value: 'attempts', label: t('auth.account.activity.attempts') },
    { value: 'events', label: t('auth.account.activity.events') }
  ]
  const disableCodeTabItems: LabTabItem[] = [
    { value: 'totp', label: t('auth.account.twofa.code_label') },
    { value: 'backup', label: t('auth.login.mfa_backup') }
  ]

  const activityTab = ref<'sessions' | 'attempts' | 'events'>('sessions')
  const disableMethodTab = ref<'totp' | 'backup'>('totp')
  const showDisable2faForm = ref(false)

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

  const enableCodeInputRef = ref<{ focus: () => void } | null>(null)
  const disableTotpInputRef = ref<{ focus: () => void } | null>(null)
  const disableBackupCodeInputRef = ref<{ focus: () => void } | null>(null)

  const successNoticeTimers = new Map<string, ReturnType<typeof setTimeout>>()

  const passwordForm = reactive({
    current: '',
    next: ''
  })
  const emailChangeForm = reactive({
    email: ''
  })

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
  const twofaStatusLabel = computed(() =>
    user.value?.is_two_factor_enabled ? t('auth.account.twofa.enabled') : t('auth.account.twofa.disabled')
  )
  const sessionActivityColumns = computed<LabDataTableColumn[]>(() => [
    { key: 'device', label: 'Устройство', cellClass: 'whitespace-normal wrap-break-word' },
    { key: 'status', label: '2FA', nowrap: true },
    { key: 'activity', label: 'Активность', cellClass: 'whitespace-normal wrap-break-word' },
    { key: 'action', label: 'Действие', nowrap: true }
  ])
  const sessionActivityRows = computed(() =>
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
        item.revokableSessionIds.length === 0 ? t('auth.account.activity.current')
        : item.hasCurrent ? t('auth.account.activity.revoke_and_logout')
        : t('auth.account.activity.revoke'),
      source: item
    }))
  )
  const loginAttemptColumns = computed<LabDataTableColumn[]>(() => [
    { key: 'createdAt', label: 'Дата', nowrap: true },
    { key: 'outcome', label: 'Результат', nowrap: true },
    { key: 'ip', label: 'IP', nowrap: true },
    { key: 'risk', label: 'Риск', nowrap: true },
    { key: 'details', label: 'Детали', cellClass: 'whitespace-normal wrap-break-word' }
  ])
  const loginAttemptRows = computed(() =>
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
  const securityEventColumns = computed<LabDataTableColumn[]>(() => [
    { key: 'createdAt', label: 'Дата', nowrap: true },
    { key: 'event', label: 'Событие', cellClass: 'whitespace-normal wrap-break-word' },
    { key: 'ip', label: 'IP', nowrap: true },
    { key: 'payload', label: 'Payload', cellClass: 'whitespace-normal wrap-break-word' }
  ])
  const securityEventRows = computed(() =>
    securityEvents.value.map(item => ({
      id: item.event_id,
      createdAt: formatDateTime(item.created_at),
      event: `${item.event_type} · ${item.severity}`,
      ip: item.ip || '—',
      payload: JSON.stringify(item.payload || {}),
      source: item
    }))
  )

  const formatDateTime = (value?: string | null) =>
    formatAbsoluteDateTime(value, { dateStyle: 'medium', timeStyle: 'short' })
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
    }, 5000)
    successNoticeTimers.set(key, timer)
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
  const refreshState = async () => {
    actionError.value = ''
    try {
      await loadMe()
      if (user.value) {
        emailChangeForm.email = user.value.email || ''
      }
    } catch (err: any) {
      actionError.value = err?.data?.message || err?.message || 'Не удалось загрузить аккаунт.'
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
    } catch (err) {
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
  const deactivate2fa = async (rawCode?: string) => {
    twofaError.value = ''
    clearSuccessNotice('twofa', twofaInfo)
    const password = disablePassword.value.trim()
    if (!password) {
      twofaError.value = 'Введите текущий пароль.'
      return
    }
    const code =
      typeof rawCode === 'string' ? rawCode
      : disableMethodTab.value === 'backup' ?
        `${disableBackupCode.value.slice(0, 4)}-${disableBackupCode.value.slice(4, 8)}`
      : disableTotpCode.value
    const normalizedCode =
      disableMethodTab.value === 'backup' ?
        code.toUpperCase().replace(/[^A-Z0-9-]+/g, '')
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
        disableMethodTab.value === 'backup' ?
          'Неверный код сброса. Попробуйте ещё раз.'
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
        focusDisableCodeInput()
      })
    }
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

  onBeforeUnmount(() => {
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
    focusDisableCodeInput()
  })
  watch(
    () => passwordForm.current,
    () => {
      if (passwordCurrentInputError.value) {
        passwordCurrentInputError.value = false
      }
    }
  )

  await ensureLoaded()
  if (user.value) {
    emailChangeForm.email = user.value.email || ''
  }
  await loadActivityState()
</script>

<template>
  <section v-if="user" class="text-(--lab-text-primary) space-y-4">
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
            user.is_two_factor_enabled ?
              showDisable2faForm ? t('auth.account.twofa.hide_form')
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
          <img :src="setupQrDataUrl" alt="QR-код для настройки 2FA" class="h-48 w-48 object-contain" loading="lazy" />
        </div>
        <div class="space-y-3">
          <p v-if="setupOtpAuthUrl" class="text-(--lab-text-muted) wrap-break-word text-xs leading-5">
            {{ t('auth.account.twofa.qr_fallback').split('{link}')[0] }}
            <a
              :href="setupOtpAuthUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="lab-focus text-(--lab-info)">
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
      <div class="flex flex-wrap items-baseline gap-x-6 gap-y-2">
        <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
          <span class="lab-text-muted shrink-0 text-xs uppercase tracking-wide">
            {{ t('auth.account.activity.last_login') }}
          </span>
          <span class="text-(--lab-text-primary) text-sm">
            {{ formatDateTime(user.last_login_at) }}
          </span>
        </div>
        <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
          <span class="lab-text-muted shrink-0 text-xs uppercase tracking-wide">
            {{ t('auth.account.activity.created_at') }}
          </span>
          <span class="text-(--lab-text-primary) text-sm">
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
              item.value === 'sessions' ? groupedSessions.length
              : item.value === 'attempts' ? loginAttempts.length
              : securityEvents.length
            }}
          </span>
        </template>
        <template #panel-sessions>
          <LabDataTable
            :columns="sessionActivityColumns"
            :rows="sessionActivityRows"
            :loading="activityLoading"
            max-height-class="max-h-96"
            empty-text="У вас пока нет активных сессий.">
            <template #cell-device="{ row }">
              <div class="space-y-1">
                <p class="text-sm">{{ row.device }}</p>
                <p class="text-(--lab-text-muted) text-xs">{{ row.ip }}</p>
              </div>
            </template>
            <template #cell-activity="{ row }">
              <div class="space-y-1 text-xs">
                <p>{{ t('auth.account.activity.sessions_count', { count: row.count }) }}</p>
                <p class="text-(--lab-text-muted)">
                  {{ t('auth.account.activity.last_activity') }}
                  <LabRelativeTime :datetime="row.lastSeenAt" compact />
                </p>
              </div>
            </template>
            <template #cell-action="{ row }">
              <LabBaseButton
                variant="secondary"
                size="lg"
                button-class="text-xs"
                :disabled="row.revokableSessionIds.length === 0"
                @click="revokeSessionGroup(row.source)">
                {{ row.action }}
              </LabBaseButton>
            </template>
          </LabDataTable>
        </template>
        <template #panel-attempts>
          <LabDataTable
            :columns="loginAttemptColumns"
            :rows="loginAttemptRows"
            :loading="activityLoading"
            max-height-class="max-h-96"
            empty-text="Попыток входа пока нет.">
            <template #cell-createdAt="{ row }">
              <div class="space-y-1 text-xs">
                <p>{{ row.createdAt }}</p>
                <p class="text-(--lab-text-muted)">
                  <LabRelativeTime :datetime="row.source.created_at" compact />
                </p>
              </div>
            </template>
            <template #cell-details="{ row }">
              <p class="text-xs wrap-break-word">{{ row.details }}</p>
            </template>
          </LabDataTable>
        </template>
        <template #panel-events>
          <LabDataTable
            :columns="securityEventColumns"
            :rows="securityEventRows"
            :loading="activityLoading"
            max-height-class="max-h-96"
            empty-text="Событий безопасности пока нет.">
            <template #cell-createdAt="{ row }">
              <div class="space-y-1 text-xs">
                <p>{{ row.createdAt }}</p>
                <p class="text-(--lab-text-muted)">
                  <LabRelativeTime :datetime="row.source.created_at" compact />
                </p>
              </div>
            </template>
            <template #cell-payload="{ row }">
              <p class="text-xs wrap-break-word">{{ row.payload }}</p>
            </template>
          </LabDataTable>
        </template>
      </LabNavTabs>
    </article>
  </section>
</template>
