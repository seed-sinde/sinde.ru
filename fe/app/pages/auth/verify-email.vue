<script setup lang="ts">
definePageMeta({
  title: 'подтверждение email',
  description: 'Завершение регистрации и активация аккаунта.'
})
const route = useRoute()
const router = useRouter()
const { t } = useInterfacePreferences()
const { verifyEmail, refresh, loadMe, canAttemptSessionRestore } = useAuth()
const token = computed(() => {
  const raw = route.query.token
  return typeof raw === 'string' ? raw : ''
})
const pending = ref(false)
const errorText = ref('')
const successText = ref('')
onMounted(async () => {
  if (!token.value) return
  pending.value = true
  try {
    const res = await verifyEmail(token.value)
    if (res.data.action === 'change_email') {
      if (canAttemptSessionRestore.value) {
        try {
          await refresh()
        } catch {
          await loadMe()
        }
      }
      successText.value = `Новый email подтверждён: ${res.data.email || 'адрес обновлён'}. Текущая сессия обновлена.`
      setTimeout(() => router.push('/auth/account?account=security'), 1200)
    } else {
      successText.value = 'Email подтверждён. Можно входить в систему.'
      setTimeout(() => router.push('/auth/login'), 1000)
    }
  } catch (err: any) {
    errorText.value = err?.data?.message || err?.message || 'Подтверждение не удалось.'
  } finally {
    pending.value = false
  }
})
</script>
<template>
  <div class="px-3 py-6 md:px-5">
    <section class="lab-surface max-w-xl space-y-4 border p-5">
      <h1 class="lab-text-primary text-2xl font-semibold">{{ t('auth.verify.title') }}</h1>
      <p v-if="pending" class="lab-text-secondary text-sm">{{ t('auth.verify.pending') }}</p>
      <LabNotify :text="errorText" tone="error" />
      <LabNotify :text="successText" tone="success" />
      <NuxtLink to="/auth/login" class="text-sm text-(--lab-accent) transition hover:text-(--lab-accent-hover)">
        {{ t('auth.verify.login_link') }}
      </NuxtLink>
    </section>
  </div>
</template>
