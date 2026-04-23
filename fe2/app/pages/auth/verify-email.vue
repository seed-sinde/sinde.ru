<script setup lang="ts">
type ApiError = Error & {
  data?: {
    message?: string
  }
}

const errorMessage = (err: unknown, fallback: string) =>
  (err as ApiError)?.data?.message || (err as Error)?.message || fallback

useSeoMeta({
  title: "Подтверждение email",
  description: "Завершение регистрации и активация аккаунта."
})

const route = useRoute()
const router = useRouter()
const {verifyEmail, refresh, loadMe, canAttemptSessionRestore} = useAuth()
const token = computed(() => (typeof route.query.token === "string" ? route.query.token : ""))
const pending = ref(false)
const errorText = ref("")
const successText = ref("")

onMounted(async () => {
  if (!token.value) return

  pending.value = true

  try {
    const res = await verifyEmail(token.value)
    if (res.data.action === "change_email") {
      if (canAttemptSessionRestore.value) {
        try {
          await refresh()
        } catch {
          await loadMe()
        }
      }
      successText.value = `Новый email подтверждён: ${res.data.email || "адрес обновлён"}.`
      setTimeout(() => router.push("/auth/account/profile"), 1200)
      return
    }

    successText.value = "Email подтверждён. Можно входить в систему."
    setTimeout(() => router.push("/auth/login"), 1000)
  } catch (err: unknown) {
    errorText.value = errorMessage(err, "Подтверждение не удалось.")
  } finally {
    pending.value = false
  }
})
</script>

<template>
  <main class="mx-auto flex w-full max-w-md flex-col gap-5 p-4">
    <header class="space-y-1">
      <h1 class="text-2xl font-semibold">Подтверждение email</h1>
      <p v-if="pending" class="text-sm opacity-75">Проверяем токен подтверждения.</p>
      <p v-else-if="!token" class="text-sm opacity-75">Токен подтверждения не найден.</p>
    </header>

    <MessageError :text="errorText" />
    <p v-if="successText" class="text-sm text-(--accent)">{{ successText }}</p>
    <NuxtLink class="text-sm text-(--accent)" to="/auth/login">Войти</NuxtLink>
  </main>
</template>
