<script setup lang="ts">
const { locale, key, load, t } = useI18nSection('auth')
await useAsyncData(key.value, load, { watch: [locale] })
const route = useRoute()
const router = useRouter()
const {
  user,
  ensureLoaded,
  loadMe,
  requestEmailChange,
  changePassword,
  setupTwoFactor,
  enableTwoFactor,
  disableTwoFactor
} = useAuth()

const disableCodeTabItems: LabTabItem[] = [
  { value: 'totp', label: t('account.twofa.code_label') },
  { value: 'backup', label: t('login.mfa_backup') }
]

const disableMethodTab = computed<'totp' | 'backup'>({
  get: () => normalizeTabRouteValue(route.query.disable2fa, ['totp', 'backup'], 'totp') as 'totp' | 'backup',
  set: value => {
    const target = buildTabRouteLocation(route, value, {
      queryKey: 'disable2fa',
      defaultValue: 'totp'
    })
    if (!target) return
    void router.replace(target)
  }
})
const showDisable2faForm = ref(false)

const backupCodes = ref<string[]>([])

const enableCode = ref('')
const disablePassword = ref('')
const disableTotpCode = ref('')
const disableBackupCode = ref('')
const setupSecret = ref('')
const setupOtpAuthUrl = ref('')
const setupQrDataUrl = ref('')

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

const twofaStatusLabel = computed(() =>
  user.value?.is_two_factor_enabled ? t('account.twofa.enabled') : t('account.twofa.disabled')
)
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
const refreshState = async () => {
  try {
    await loadMe()
    if (user.value) {
      emailChangeForm.email = user.value.email || ''
    }
  } catch (err: any) {
    emailChangeError.value = err?.data?.message || err?.message || 'Не удалось загрузить аккаунт.'
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
    emailChangeInfo.value = t('account.email.request_sent', { email: nextEmail })
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
    twofaInfo.value = t('account.twofa.setup_created')
  } catch (err) {
    twofaError.value = extractApiErrorMessage(err, 'Не удалось начать настройку 2FA.')
  }
}
const activate2fa = async () => {
  twofaError.value = ''
  clearSuccessNotice('twofa', twofaInfo)
  enableCode.value = enableCode.value.replace(/\D+/g, '').slice(0, 6)
  if (enableCode.value.length !== 6) {
    twofaError.value = t('account.twofa.code_hint')
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
</script>

<template>
  <section v-if="user" class="space-y-4 text-(--lab-text-primary)">
    <article class="space-y-3">
      <h2 class="text-xl sm:text-2xl">{{ t('account.email.title') }}</h2>
      <p class="text-sm text-(--lab-text-muted)">
        {{ t('account.email.description', { email: user.email }) }}
      </p>
      <div class="grid gap-3 sm:max-w-fit">
        <LabBaseField :label="t('account.email.new_label')" for-id="account-next-email">
          <LabBaseInput
            id="account-next-email"
            v-model="emailChangeForm.email"
            name="email"
            type="email"
            autocomplete="email"
            class="w-full"
            :placeholder="t('account.email.new_placeholder')"
            @keydown.enter.prevent="submitEmailChange"
          />
        </LabBaseField>
      </div>
      <div class="flex flex-wrap items-center">
        <LabBaseButton
          :label="t('account.email.submit')"
          variant="primary"
          size="lg"
          class="text-xs"
          :disabled="emailChangePending"
          @click="submitEmailChange"
        />
      </div>
      <LabNotify :text="emailChangeError" tone="error" size="xs" />
      <LabNotify :text="emailChangeInfo" tone="success" size="xs" />
    </article>
    <article class="space-y-3">
      <h2 class="flex flex-wrap items-center gap-3 text-xl sm:text-2xl">
        {{ t('account.password.title') }}
      </h2>
      <span class="text-sm text-(--lab-text-muted)">{{ t('account.password.description') }}</span>
      <div class="grid gap-3 sm:max-w-fit lg:grid-rows-2">
        <LabBaseField :label="t('account.password.current_label')" for-id="account-current-password">
          <LabBaseInput
            id="account-current-password"
            v-model="passwordForm.current"
            name="current_password"
            type="password"
            autocomplete="current-password"
            :class="['w-full', passwordCurrentInputError ? 'text-(--lab-danger)' : '']"
            :placeholder="t('account.password.current_placeholder')"
          />
        </LabBaseField>
        <LabBaseField :label="t('account.password.next_label')" for-id="account-next-password">
          <LabBaseInput
            id="account-next-password"
            v-model="passwordForm.next"
            name="new_password"
            type="password"
            autocomplete="new-password"
            class="w-full"
            :placeholder="t('account.password.next_placeholder')"
          />
        </LabBaseField>
      </div>
      <div class="flex flex-wrap items-center">
        <LabBaseButton
          :label="t('account.password.submit')"
          variant="primary"
          size="lg"
          class="text-xs"
          @click="submitPasswordChange"
        />
      </div>
      <LabNotify :text="passwordError" tone="error" size="xs" />
      <LabNotify :text="passwordInfo" tone="success" size="xs" />
    </article>
    <article class="space-y-3">
      <div class="space-y-3">
        <h2 class="flex flex-wrap items-center gap-3 text-xl sm:text-2xl">2FA</h2>
        <p class="text-sm text-(--lab-text-muted)">
          {{ twofaStatusLabel }}
        </p>

        <LabBaseButton
          :variant="user.is_two_factor_enabled ? 'danger' : 'success'"
          size="lg"
          class="text-xs"
          @click="user.is_two_factor_enabled ? toggleDisable2faForm() : begin2faSetup()"
        >
          {{
            user.is_two_factor_enabled
              ? showDisable2faForm
                ? t('account.twofa.hide_form')
                : t('account.twofa.disable')
              : t('account.twofa.get_secret')
          }}
        </LabBaseButton>
      </div>
      <LabNotify :text="twofaError" tone="error" />
      <LabNotify :text="twofaInfo" tone="success" :temporary="false" />
      <div
        v-if="!user.is_two_factor_enabled && setupSecret"
        class="grid items-start gap-3 xl:grid-cols-[auto_minmax(0,1fr)]"
      >
        <div v-if="setupQrDataUrl">
          <img :src="setupQrDataUrl" alt="QR-код для настройки 2FA" class="h-48 w-48 object-contain" loading="lazy" />
        </div>
        <div class="space-y-3">
          <p v-if="setupOtpAuthUrl" class="text-xs leading-5 wrap-break-word text-(--lab-text-muted)">
            {{ t('account.twofa.qr_fallback').split('{link}')[0] }}
            <a :href="setupOtpAuthUrl" target="_blank" rel="noopener noreferrer" class="lab-focus text-(--lab-info)">
              {{ t('account.twofa.qr_fallback_link') }}
            </a>
            {{ t('account.twofa.qr_fallback').split('{link}')[1] }}
          </p>
          <LabCopyBlock
            :label="t('account.twofa.copy_key_label')"
            :value="setupSecret"
            variant="dark-cyan"
            class="w-full"
            :title-idle="t('account.twofa.copy_key_idle')"
            :title-success="t('account.twofa.copy_key_success')"
            :title-error="t('account.twofa.copy_key_error')"
            :show-state-tooltip="true"
          />
          <AuthCodeInput
            id="account-enable-2fa-code"
            ref="enableCodeInputRef"
            :model-value="enableCode"
            name="one_time_code_enable"
            :label="t('account.twofa.code_label')"
            :hint="t('account.twofa.code_hint')"
            :invalid="Boolean(twofaError)"
            :valid="Boolean(twofaInfo) && !Boolean(twofaError)"
            @update:model-value="onEnableCodeInput"
            @complete="activate2fa"
          />
        </div>
      </div>
      <div v-if="user.is_two_factor_enabled && showDisable2faForm" class="space-y-3">
        <input
          type="text"
          name="account_disable_2fa_username"
          autocomplete="username"
          tabindex="-1"
          aria-hidden="true"
          class="pointer-events-none absolute -left-96 h-px w-px opacity-0"
        />
        <LabBaseField label="Текущий пароль" for-id="account-disable-2fa-password" field-class="min-w-0">
          <LabBaseInput
            id="account-disable-2fa-password"
            v-model="disablePassword"
            name="disable_2fa_password"
            type="password"
            autocomplete="current-password"
            class="w-full"
            placeholder="Текущий пароль"
          />
        </LabBaseField>
        <LabNavTabs v-model="disableMethodTab" :items="disableCodeTabItems" route-query-key="disable2fa">
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
              @complete="deactivate2fa"
            />
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
              @complete="deactivate2fa(`${$event.slice(0, 4)}-${$event.slice(4, 8)}`)"
            />
          </template>
        </LabNavTabs>
      </div>
      <div v-if="backupCodes.length" class="space-y-3 text-(--lab-info)">
        <div>
          <p class="text-sm font-medium">Резервные коды</p>
          <p class="text-xs leading-5">Сохрани все 10 кодов. Каждый код сработает только один раз.</p>
        </div>
        <div class="grid grid-cols-2 gap-2 sm:grid-cols-3">
          <div
            v-for="code in backupCodes"
            :key="code"
            class="inline-flex min-h-11 items-center justify-center text-center font-mono text-sm tracking-widest"
          >
            {{ code }}
          </div>
        </div>
        <p class="text-xs leading-5">Храни их вне браузера и не передавай через мессенджеры.</p>
      </div>
    </article>
  </section>
</template>
