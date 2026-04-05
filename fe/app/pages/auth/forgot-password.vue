<script setup lang="ts">
  import { isValidEmail } from '~/utils/email'
  definePageMeta({
    title: 'восстановление доступа',
    description: 'Запрос ссылки на сброс пароля.'
  })
  const { forgotPassword } = useAuth()
  const { t } = useInterfacePreferences()
  const email = ref('')
  const pending = ref(false)
  const errorText = ref('')
  const successText = ref('')
  const hasInvalidEmail = computed(() => Boolean(email.value.trim()) && !isValidEmail(email.value))
  const submit = async () => {
    if (!isValidEmail(email.value)) {
      errorText.value = 'Укажите корректный email.'
      successText.value = ''
      return
    }
    pending.value = true
    errorText.value = ''
    successText.value = ''
    try {
      await forgotPassword(email.value)
      successText.value = 'Если аккаунт существует, инструкция отправлена на email.'
    } catch (err: any) {
      errorText.value = err?.data?.message || err?.message || 'Не удалось отправить письмо.'
    } finally {
      pending.value = false
    }
  }
</script>
<template>
  <div class="px-3 py-6 md:px-5">
    <section class="max-w-xl space-y-5 border p-5 lab-surface">
      <h1 class="text-2xl font-semibold lab-text-primary">{{ t('auth.forgot.title') }}</h1>
      <p class="text-sm lab-text-muted">{{ t('auth.forgot.description') }}</p>
      <form class="space-y-4" @submit.prevent="submit">
        <LabField label="Email" forId="forgot-password-email">
          <LabBaseInput
            id="forgot-password-email"
            v-model="email"
            name="email"
            type="email"
            autocomplete="email"
            :invalid="hasInvalidEmail"
            input-class="w-full"
            placeholder="you@example.com" />
        </LabField>
        <LabBaseButton type="submit" variant="primary" size="xl" :disabled="pending">
          {{ t('auth.forgot.submit') }}
        </LabBaseButton>
      </form>
      <LabNotify :text="errorText" tone="error" />
      <LabNotify :text="successText" tone="success" />
    </section>
  </div>
</template>
