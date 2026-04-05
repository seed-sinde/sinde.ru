<script setup lang="ts">
definePageMeta({
  title: 'подтверждение email',
  description: 'Завершение регистрации и активация аккаунта.',
})
const route = useRoute()
const router = useRouter()
const { t } = useInterfacePreferences()
const { verifyEmail } = useAuth()
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
    await verifyEmail(token.value)
    successText.value = 'Email подтверждён. Можно входить в систему.'
    setTimeout(() => router.push('/auth/login'), 1000)
  } catch (err: any) {
    errorText.value = err?.data?.message || err?.message || 'Подтверждение не удалось.'
  } finally {
    pending.value = false
  }
})
</script>
<template>
  <div class="px-3 py-6 md:px-5">
    <section class="max-w-xl space-y-4 border p-5 lab-surface">
      <h1 class="text-2xl font-semibold lab-text-primary">{{ t('auth.verify.title') }}</h1>
      <p v-if="pending" class="text-sm lab-text-secondary">{{ t('auth.verify.pending') }}</p>
      <LabNotify :text="errorText" tone="error" />
      <LabNotify :text="successText" tone="success" />
      <NuxtLink to="/auth/login" class="auth-verify-link text-sm">{{ t('auth.verify.login_link') }}</NuxtLink>
    </section>
  </div>
</template>
<style scoped>
  .auth-verify-link {
    color: var(--lab-accent);
    text-decoration: none;
    transition: color 0.2s ease;
  }
  .auth-verify-link:hover {
    color: var(--lab-accent-hover);
  }
</style>
