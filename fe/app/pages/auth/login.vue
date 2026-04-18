<script setup lang="ts">
import {
  extractApiErrorMessage,
  getTwoFactorErrorMessage,
  isMfaTicketExpiredError,
  MFA_TICKET_EXPIRED_MESSAGE
} from '~/utils/authErrors'
import { isValidEmail, normalizeEmail } from '~/utils/email'
import { normalizeAuthNextPath } from '~/utils/authNavigation'
definePageMeta({
  middleware: 'guest-only'
})
const title = 'Вход'
usePageSeo({
  title,
  description: 'Вход в аккаунт.'
})
const router = useRouter()
const route = useRoute()
const { locale, key, load, t } = useI18nSection('auth')
await useAsyncData(key.value, load, { watch: [locale] })
const { login, completeMfa, mfaTicket, mfaExpiresAt, requestEmailVerification, resetMfaTicket } = useAuth()
resetMfaTicket()
const email = ref('')
const password = ref('')
const mfaCode = ref('')
const backupCode = ref('')
const mfaMethodTab = computed<'totp' | 'backup'>({
  get: () => normalizeTabRouteValue(route.query.mfa, ['totp', 'backup'], 'totp') as 'totp' | 'backup',
  set: value => {
    const target = buildTabRouteLocation(route, value, {
      queryKey: 'mfa',
      defaultValue: 'totp'
    })
    if (!target) return
    void router.replace(target)
  }
})
const emailInputRef = ref<{ focus: () => void } | null>(null)
const mfaCodeInputRef = ref<{ focus: () => void } | null>(null)
const backupCodeInputRef = ref<{ focus: () => void } | null>(null)
const lastSubmittedMfaCode = ref('')
const lastSubmittedBackupCode = ref('')
const pending = ref(false)
const errorText = ref('')
const resendPending = ref(false)
const resendFeedback = ref('')
const resendError = ref('')
const verificationEmail = ref('')
const resendCooldownUntil = ref(0)
const currentTimestamp = ref(Date.now())
const resendCooldownLeft = computed(() => {
  const diff = resendCooldownUntil.value - Date.now()
  return diff > 0 ? Math.ceil(diff / 1000) : 0
})
const canResendVerification = computed(() => {
  return !pending.value && !resendPending.value && resendCooldownLeft.value === 0 && Boolean(verificationEmail.value)
})
const mfaExpiresAtTimestamp = computed(() => {
  const parsed = Date.parse(String(mfaExpiresAt.value || ''))
  return Number.isFinite(parsed) ? parsed : 0
})
const mfaCountdownLeft = computed(() => {
  if (!mfaTicket.value || !mfaExpiresAtTimestamp.value) return 0
  const diff = mfaExpiresAtTimestamp.value - currentTimestamp.value
  return diff > 0 ? Math.ceil(diff / 1000) : 0
})
const mfaCountdownLabel = computed(() => {
  const minutes = Math.floor(mfaCountdownLeft.value / 60)
  const seconds = mfaCountdownLeft.value % 60
  return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
})
let resendCooldownTimer: ReturnType<typeof setInterval> | null = null
let mfaCountdownTimer: ReturnType<typeof setInterval> | null = null
const syncResendCooldown = () => {
  if (resendCooldownLeft.value > 0) return
  if (resendCooldownTimer) {
    clearInterval(resendCooldownTimer)
    resendCooldownTimer = null
  }
}
const startResendCooldown = (seconds = 60) => {
  if (!import.meta.client) return
  resendCooldownUntil.value = Date.now() + seconds * 1000
  if (resendCooldownTimer) return
  resendCooldownTimer = setInterval(syncResendCooldown, 1000)
}
const clearVerificationState = () => {
  if (verificationEmail.value && normalizeEmail(email.value) !== verificationEmail.value) {
    verificationEmail.value = ''
  }
  resendFeedback.value = ''
  resendError.value = ''
}
const stopMfaCountdown = () => {
  if (!mfaCountdownTimer) return
  clearInterval(mfaCountdownTimer)
  mfaCountdownTimer = null
}
const syncMfaCountdown = () => {
  currentTimestamp.value = Date.now()
  if (mfaCountdownLeft.value > 0) return
  stopMfaCountdown()
}
const startMfaCountdown = () => {
  if (!import.meta.client) return
  currentTimestamp.value = Date.now()
  if (mfaCountdownTimer) return
  mfaCountdownTimer = setInterval(syncMfaCountdown, 1000)
}
const clearMfaInputs = () => {
  mfaCode.value = ''
  backupCode.value = ''
  lastSubmittedMfaCode.value = ''
  lastSubmittedBackupCode.value = ''
  mfaMethodTab.value = 'totp'
}
const resetExpiredMfaFlow = async () => {
  pending.value = false
  clearMfaInputs()
  resetMfaTicket()
  errorText.value = MFA_TICKET_EXPIRED_MESSAGE
  await nextTick()
  emailInputRef.value?.focus()
}
const nextTarget = computed(() => {
  const next = typeof route.query.next === 'string' ? route.query.next : ''
  return normalizeAuthNextPath(next) || '/auth/account/profile'
})
const mfaTabItems = computed<LabTabItem[]>(() => [
  { value: 'totp', label: t('login.mfa_totp') },
  { value: 'backup', label: t('login.mfa_backup') }
])
const breadcrumbItems = computed<BreadcrumbItem[]>(() =>
  mfaTicket.value
    ? [
        { label: title, to: '/auth/login' },
        {
          label: mfaTabItems.value.find(item => item.value === mfaMethodTab.value)?.label || t('login.mfa_totp'),
          current: true,
          kind: 'tab'
        }
      ]
    : [{ label: title, current: true }]
)
const submitPassword = async () => {
  if (!isValidEmail(email.value)) {
    errorText.value = 'Укажите корректный email.'
    clearVerificationState()
    await nextTick()
    emailInputRef.value?.focus()
    return
  }
  pending.value = true
  errorText.value = ''
  clearVerificationState()
  try {
    const res = await login(email.value, password.value)
    if (res.data.mfa_required) return
    await router.push(nextTarget.value)
  } catch (err: any) {
    const message = extractApiErrorMessage(err, 'Не удалось войти.')
    errorText.value = message
    if (message === 'email is not verified') {
      verificationEmail.value = normalizeEmail(email.value)
      resendFeedback.value = ''
    } else {
      verificationEmail.value = ''
    }
  } finally {
    pending.value = false
  }
}
const focusMfaCodeInput = () => {
  mfaCodeInputRef.value?.focus()
}
const focusBackupCodeInput = () => {
  backupCodeInputRef.value?.focus()
}
const focusActiveMfaInput = () => {
  if (mfaMethodTab.value === 'backup') {
    focusBackupCodeInput()
    return
  }
  focusMfaCodeInput()
}
const onMfaCodeInput = (nextValue: string) => {
  mfaCode.value = nextValue
  errorText.value = ''
  if (nextValue !== lastSubmittedMfaCode.value) {
    lastSubmittedMfaCode.value = ''
  }
}
const onBackupCodeInput = (nextValue: string) => {
  backupCode.value = nextValue
  errorText.value = ''
  if (nextValue !== lastSubmittedBackupCode.value) {
    lastSubmittedBackupCode.value = ''
  }
}
watch(mfaTicket, async ticket => {
  clearMfaInputs()
  if (!ticket) {
    stopMfaCountdown()
    await nextTick()
    emailInputRef.value?.focus()
    return
  }
  startMfaCountdown()
  await nextTick()
  focusActiveMfaInput()
})
watch([mfaTicket, mfaCountdownLeft], async ([ticket, secondsLeft], previous) => {
  if (!ticket || secondsLeft > 0) return
  if (previous?.[0] === ticket && previous[1] === 0) return
  await resetExpiredMfaFlow()
})
watch(mfaMethodTab, async () => {
  errorText.value = ''
  await nextTick()
  if (import.meta.client) {
    requestAnimationFrame(() => {
      focusActiveMfaInput()
    })
    return
  }
  focusActiveMfaInput()
})
const submitMfa = async () => {
  if (mfaCountdownLeft.value === 0) {
    await resetExpiredMfaFlow()
    return
  }
  const normalizedCode = mfaCode.value.replace(/\D+/g, '').slice(0, 6)
  mfaCode.value = normalizedCode
  if (normalizedCode.length !== 6) {
    errorText.value = 'Введите 6 цифр кода 2FA.'
    return
  }
  pending.value = true
  errorText.value = ''
  try {
    await completeMfa(normalizedCode)
    await router.push(nextTarget.value)
  } catch (err: any) {
    if (isMfaTicketExpiredError(err)) {
      await resetExpiredMfaFlow()
      return
    }
    const message = String(err?.data?.message || err?.message || '')
    const chunkLoadLikeError = /dynamically imported module|Importing a module script failed|Failed to fetch/i.test(
      message
    )
    if (chunkLoadLikeError && import.meta.client) {
      window.location.assign(nextTarget.value)
      return
    }
    errorText.value = getTwoFactorErrorMessage(err, 'totp', 'Неверный код из приложения. Попробуйте ещё раз.')
    mfaCode.value = ''
    lastSubmittedMfaCode.value = ''
    await nextTick()
    focusMfaCodeInput()
  } finally {
    pending.value = false
  }
}
const submitBackupCode = async () => {
  if (mfaCountdownLeft.value === 0) {
    await resetExpiredMfaFlow()
    return
  }
  const normalizedCode = backupCode.value
    .toUpperCase()
    .replace(/[^A-Z0-9]+/g, '')
    .slice(0, 8)
  backupCode.value = normalizedCode
  if (normalizedCode.length !== 8) {
    errorText.value = 'Введите 8 символов кода сброса.'
    return
  }
  pending.value = true
  errorText.value = ''
  try {
    await completeMfa(`${normalizedCode.slice(0, 4)}-${normalizedCode.slice(4)}`)
    await router.push(nextTarget.value)
  } catch (err: any) {
    if (isMfaTicketExpiredError(err)) {
      await resetExpiredMfaFlow()
      return
    }
    const message = String(err?.data?.message || err?.message || '')
    const chunkLoadLikeError = /dynamically imported module|Importing a module script failed|Failed to fetch/i.test(
      message
    )
    if (chunkLoadLikeError && import.meta.client) {
      window.location.assign(nextTarget.value)
      return
    }
    errorText.value = getTwoFactorErrorMessage(err, 'backup', 'Неверный код сброса. Попробуйте ещё раз.')
    backupCode.value = ''
    lastSubmittedBackupCode.value = ''
    await nextTick()
    focusBackupCodeInput()
  } finally {
    pending.value = false
  }
}
const resendVerification = async () => {
  if (!canResendVerification.value) return
  resendPending.value = true
  resendError.value = ''
  resendFeedback.value = ''
  try {
    await requestEmailVerification(verificationEmail.value)
    resendFeedback.value = `Письмо с подтверждением повторно отправлено на ${verificationEmail.value}.`
    startResendCooldown(60)
  } catch (err: any) {
    const status = Number(err?.status || err?.statusCode || 0)
    if (status === 429) {
      resendError.value = 'Слишком много запросов. Подождите немного и попробуйте снова.'
      startResendCooldown(60)
    } else {
      resendError.value = extractApiErrorMessage(err, 'Не удалось отправить письмо повторно.')
    }
  } finally {
    resendPending.value = false
  }
}
onBeforeUnmount(() => {
  if (resendCooldownTimer) clearInterval(resendCooldownTimer)
  stopMfaCountdown()
})
onMounted(async () => {
  await nextTick()
  if (mfaTicket.value) return
  emailInputRef.value?.focus()
})
</script>
<template>
  <div>
    <LabNavHeader :title :breadcrumb-items="breadcrumbItems" />
    <section class="w-full space-y-4 p-4 sm:max-w-sm">
      <h1 class="text-2xl font-semibold text-(--lab-text-primary)">{{ t('login.title') }}</h1>
      <template v-if="!mfaTicket">
        <form class="space-y-4" @submit.prevent="submitPassword">
          <LabField :label="t('login.email')" for-id="auth-email">
            <LabBaseInput
              id="auth-email"
              ref="emailInputRef"
              v-model="email"
              name="email"
              type="email"
              autocomplete="email"
              autocapitalize="none"
              spellcheck="false"
              :invalid="Boolean(email.trim()) && !isValidEmail(email)"
              input-class="w-full"
              placeholder="you@example.com"
              @input="clearVerificationState"
            />
          </LabField>
          <LabField :label="t('login.password')" for-id="auth-password">
            <LabBaseInput
              id="auth-password"
              v-model="password"
              name="password"
              type="password"
              autocomplete="current-password"
              input-class="w-full"
              placeholder="Минимум 12 символов"
            />
          </LabField>
          <LabBaseButton
            type="submit"
            variant="primary"
            size="xl"
            :disabled="pending"
            :label="t('login.submit')"
          />
        </form>
        <div
          v-if="verificationEmail"
          class="space-y-4 border-[color-mix(in_srgb,var(--lab-accent)_34%,transparent)] bg-[color-mix(in_srgb,var(--lab-accent)_10%,var(--lab-bg-surface))] p-4 text-(--lab-text-primary)"
        >
          <div class="space-y-2">
            <p class="text-xs font-semibold tracking-[0.14em] text-(--lab-accent) uppercase">
              {{ t('verify.title') }}
            </p>
            <p class="text-sm">{{ t('login.verify_email_description') }}</p>
            <p class="text-2xl leading-tight font-semibold break-all md:text-3xl">{{ verificationEmail }}</p>
          </div>
          <LabBaseButton variant="secondary" size="xl" :disabled="!canResendVerification" @click="resendVerification">
            <span v-if="resendPending">{{ t('login.sending') }}</span>
            <span v-else-if="resendCooldownLeft > 0">
              {{ t('login.resend_wait', { seconds: resendCooldownLeft }) }}
            </span>
            <span v-else>{{ t('login.resend') }}</span>
          </LabBaseButton>
          <LabNotify :text="resendFeedback" tone="warning" />
          <LabNotify :text="resendError" tone="error" />
        </div>
      </template>
      <div v-else class="space-y-4">
        <div
          class="border-[color-mix(in_srgb,var(--lab-accent)_34%,transparent)] bg-[color-mix(in_srgb,var(--lab-accent)_8%,var(--lab-bg-surface))] p-3 text-sm text-(--lab-text-primary)"
        >
          <p>{{ t('login.mfa_description') }}</p>
        </div>
        <LabNavTabs
          v-model="mfaMethodTab"
          :items="mfaTabItems"
          route-query-key="mfa"
        >
          <template #panel-totp>
            <AuthCodeInput
              id="auth-mfa-code"
              ref="mfaCodeInputRef"
              :model-value="mfaCode"
              name="one_time_code"
              :label="t('login.mfa_totp')"
              hint="Проверка начнётся автоматически после ввода 6 цифр."
              :invalid="Boolean(errorText)"
              @update:model-value="onMfaCodeInput"
              @complete="submitMfa"
            />
          </template>
          <template #panel-backup>
            <AuthRecoveryCodeInput
              id="auth-backup-code"
              ref="backupCodeInputRef"
              :model-value="backupCode"
              name="backup_code"
              :label="t('login.mfa_backup')"
              hint="Введите 8 символов. Проверка начнётся автоматически."
              :invalid="Boolean(errorText)"
              @update:model-value="onBackupCodeInput"
              @complete="submitBackupCode"
            />
          </template>
        </LabNavTabs>
        <LabNotify
          :text="mfaTicket ? `до сброса авторизации: ${mfaCountdownLabel}` : ''"
          tone="warning"
          class-name="border-[color-mix(in_srgb,var(--lab-accent)_34%,transparent)] bg-[color-mix(in_srgb,var(--lab-accent)_10%,var(--lab-bg-surface))] px-3 py-2 font-mono text-(--lab-accent)"
        />
        <p v-if="pending" class="text-sm text-(--lab-text-muted)">Проверяем код…</p>
      </div>
      <LabNotify :text="errorText" tone="error" />
      <div class="flex flex-wrap gap-3 text-sm text-(--lab-text-muted)">
        <NuxtLink to="/auth/register" class="text-(--lab-accent) transition hover:text-(--lab-accent-hover)">
          {{ t('login.register') }}
        </NuxtLink>
        <NuxtLink to="/auth/forgot-password" class="text-(--lab-text-muted) transition hover:text-(--lab-text-primary)">
          {{ t('login.forgot') }}
        </NuxtLink>
      </div>
    </section>
  </div>
</template>
