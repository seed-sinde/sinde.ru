<script setup lang="ts">
  definePageMeta({
    title: 'сброс пароля',
    description: 'Установка нового пароля по одноразовому токену.'
  })
  const route = useRoute()
  const router = useRouter()
  const { t } = useInterfacePreferences()
  const { resetPassword } = useAuth()
  const token = computed(() => {
    const raw = route.query.token
    return typeof raw === 'string' ? raw : ''
  })
  const password = ref('')
  const pending = ref(false)
  const errorText = ref('')
  const successText = ref('')
  const submit = async () => {
    pending.value = true
    errorText.value = ''
    successText.value = ''
    try {
      await resetPassword(token.value, password.value)
      successText.value = 'Пароль обновлён. Теперь можно войти.'
      setTimeout(() => router.push('/auth/login'), 1000)
    } catch (err: any) {
      errorText.value = err?.data?.message || err?.message || 'Сброс пароля не удался.'
    } finally {
      pending.value = false
    }
  }
</script>
<template>
  <div class="px-3 py-6 md:px-5">
    <section class="max-w-xl space-y-5 border p-5 lab-surface">
      <h1 class="text-2xl font-semibold lab-text-primary">{{ t('auth.reset.title') }}</h1>
      <p class="text-sm lab-text-muted">{{ t('auth.reset.description') }}</p>
      <form class="space-y-4" @submit.prevent="submit">
        <LabField label="Новый пароль" forId="reset-password-new">
          <LabBaseInput
            id="reset-password-new"
            v-model="password"
            name="new_password"
            type="password"
            autocomplete="new-password"
            input-class="w-full"
            placeholder="Минимум 12 символов" />
        </LabField>
        <LabBaseButton type="submit" variant="primary" size="xl" :disabled="pending || !token">
          {{ t('auth.reset.submit') }}
        </LabBaseButton>
      </form>
      <LabNotify :text="errorText" tone="error" />
      <LabNotify :text="successText" tone="success" />
    </section>
  </div>
</template>
