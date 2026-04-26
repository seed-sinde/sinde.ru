<script setup lang="ts">
type ApiError = Error & {
  data?: {
    message?: string
  }
}

const errorMessage = (err: unknown, fallback: string) =>
  (err as ApiError)?.data?.message || (err as Error)?.message || fallback

useSeoMeta({
  title: "Сброс пароля",
  description: "Установка нового пароля по одноразовому токену."
})

const route = useRoute()
const router = useRouter()
const {resetPassword} = useAuth()
const token = computed(() => (typeof route.query.token === "string" ? route.query.token : ""))
const password = ref("")
const pending = ref(false)
const errorText = ref("")
const successText = ref("")

const submit = async () => {
  pending.value = true
  errorText.value = ""
  successText.value = ""

  try {
    await resetPassword(token.value, password.value)
    successText.value = "Пароль обновлён. Теперь можно войти."
    setTimeout(() => router.push("/auth/login"), 1000)
  } catch (err: unknown) {
    errorText.value = errorMessage(err, "Сброс пароля не удался.")
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <main class="mx-auto flex w-full max-w-md flex-col gap-5 p-4">
    <header class="space-y-1">
      <h1 class="text-2xl font-semibold">Сброс пароля</h1>
      <p class="text-sm opacity-75">Укажите новый пароль для аккаунта.</p>
    </header>

    <form class="flex flex-col items-start gap-3" @submit.prevent="submit">
      <UiField label="Новый пароль" hint="Минимум 12 символов" for-id="reset-password-new">
        <UiInput
          id="reset-password-new"
          v-model="password"
          type="password"
          autocomplete="new-password"
        />
      </UiField>
      <UiButton
        label="Обновить пароль"
        :loading="pending"
        :disabled="pending || !token || password.length < 12"
        type="submit"
      />
    </form>

    <MessageError :text="errorText" />
    <p v-if="successText" class="text-sm text-(--accent)">{{ successText }}</p>
  </main>
</template>
