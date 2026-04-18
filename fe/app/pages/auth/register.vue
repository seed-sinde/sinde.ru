<script setup lang="ts">
definePageMeta({
  middleware: 'guest-only'
})
const title = 'Регистрация'
usePageSeo({
  title,
  description: 'Регистрация нового аккаунта sinde с подтверждением email.'
})
const router = useRouter()
const runtimeConfig = useRuntimeConfig()
const { localeTag } = useInterfacePreferences()
const { locale, key, load, t } = useI18nSection('auth')
await useAsyncData(key.value, load, { watch: [locale] })
const { register } = useAuth()
const passwordMinLength = Math.max(1, Number(runtimeConfig.public.authPasswordMinLength || 12))
const form = reactive({
  email: '',
  password: '',
  display_name: '',
  locale: localeTag.value,
  timezone: Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC'
})
const pending = ref(false)
const errorText = ref('')
const successText = ref('')
const hasInvalidEmail = computed(() => Boolean(String(form.email || '').trim()) && !isValidEmail(form.email))
const passwordPolicyRules = computed(() => getPasswordPolicyRules(form.password, form.email, passwordMinLength))
const passwordPolicyItems = computed(() =>
  passwordPolicyRules.value.map(item => ({
    ...item,
    label:
      item.key === 'length'
        ? t('register.password_length', { min: passwordMinLength })
        : item.key === 'email'
          ? t('register.password_email')
          : t('register.password_common')
  }))
)
const canSubmit = computed(() => isValidEmail(form.email) && passwordPolicyRules.value.every(item => item.passed))
const submit = async () => {
  if (!isValidEmail(form.email)) {
    errorText.value = 'Укажите корректный email.'
    successText.value = ''
    return
  }
  if (!canSubmit.value) {
    errorText.value = 'Проверь требования к паролю перед регистрацией.'
    successText.value = ''
    return
  }
  pending.value = true
  errorText.value = ''
  successText.value = ''
  try {
    const res = await register(form)
    successText.value = `Аккаунт создан. Проверь email и подтверди адрес в течение ${res.data.verification_ttl}.`
    setTimeout(() => router.push('/auth/login'), 1200)
  } catch (err: any) {
    errorText.value = err?.data?.message || err?.message || 'Регистрация не удалась.'
  } finally {
    pending.value = false
  }
}
watch(localeTag, nextLocale => {
  form.locale = nextLocale
})
</script>
<template>
  <div>
    <LabNavHeader :title />
    <section class="w-full space-y-4 p-4 sm:max-w-sm">
      <div class="space-y-2">
        <h1 class="text-2xl font-semibold text-(--lab-text-primary)">{{ t('register.title') }}</h1>
        <p class="text-sm text-(--lab-text-muted)">{{ t('register.description') }}</p>
      </div>
      <form class="space-y-2" @submit.prevent="submit">
        <LabField :label="t('register.name')" for-id="register-name">
          <LabBaseInput
            id="register-name"
            v-model="form.display_name"
            name="display_name"
            type="text"
            input-class="w-full"
          />
        </LabField>
        <LabField :label="t('register.email')" for-id="register-email">
          <LabBaseInput
            id="register-email"
            v-model="form.email"
            name="email"
            type="email"
            autocomplete="email"
            :invalid="hasInvalidEmail"
            input-class="w-full"
          />
        </LabField>
        <LabField :label="t('register.password')" for-id="register-password">
          <LabBaseInput
            id="register-password"
            v-model="form.password"
            name="password"
            type="password"
            autocomplete="new-password"
            input-class="w-full"
            :placeholder="`Минимум ${passwordMinLength} символов`"
          />
          <div class="p-3">
            <p class="text-xs text-(--lab-text-muted)">{{ t('register.password_rules') }}</p>
            <ul class="mt-2 space-y-1">
              <li
                v-for="item in passwordPolicyItems"
                :key="item.label"
                class="flex items-start gap-2 text-xs"
                :class="
                  item.state === 'passed'
                    ? 'lab-text-success'
                    : item.state === 'failed'
                      ? 'lab-text-danger'
                      : 'lab-text-muted'
                "
              >
                <Icon
                  :name="
                    item.state === 'passed'
                      ? 'ic:round-check-circle'
                      : item.state === 'failed'
                        ? 'ic:round-cancel'
                        : 'ic:round-radio-button-unchecked'
                  "
                  class="mt-0.5 h-3.5 w-3.5 shrink-0"
                />
                <span>{{ item.label }}</span>
              </li>
            </ul>
          </div>
        </LabField>
        <LabBaseButton
          type="submit"
          variant="primary"
          size="xl"
          :disabled="pending || !canSubmit"
          button-class="text-sm font-medium"
        >
          {{ t('register.submit') }}
        </LabBaseButton>
      </form>
      <LabNotify :text="errorText" tone="error" />
      <LabNotify :text="successText" tone="success" />
      <NuxtLink to="/auth/login" class="text-sm text-(--lab-accent) transition hover:text-(--lab-accent-hover)">
        {{ t('register.has_account') }}
      </NuxtLink>
    </section>
  </div>
</template>
