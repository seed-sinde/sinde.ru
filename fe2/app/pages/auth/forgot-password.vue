<script setup lang="ts">
type ApiError = Error & {
  data?: {
    message?: string
  }
}

const errorMessage = (err: unknown, fallback: string) =>
  (err as ApiError)?.data?.message || (err as Error)?.message || fallback
const isValidEmail = (value: string) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value.trim())

useSeoMeta({
  title: "Восстановление доступа",
  description: "Запрос ссылки на сброс пароля."
})

const {forgotPassword} = useAuth()
const email = ref("")
const pending = ref(false)
const errorText = ref("")
const successText = ref("")

const submit = async () => {
  if (!isValidEmail(email.value)) {
    errorText.value = "Укажите корректный email."
    successText.value = ""
    return
  }

  pending.value = true
  errorText.value = ""
  successText.value = ""

  try {
    await forgotPassword(email.value)
    successText.value = "Если аккаунт существует, инструкция отправлена на email."
  } catch (err: unknown) {
    errorText.value = errorMessage(err, "Не удалось отправить письмо.")
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <main class="p-4">
    <h1 class="text-2xl font-semibold">Восстановление доступа</h1>

    <form class="space-y-2" @submit.prevent="submit">
      <UiField
        label="Email"
        for-id="forgot-password-email"
        hint="Введите email, чтобы получить ссылку на сброс пароля."
      >
        <UiInput
          v-model="email"
          placeholder="seed@sinde.ru"
          type="email"
          id="forgot-password-email"
          autocomplete="username"
        />
      </UiField>
      <UiButton label="Отправить" :loading="pending" :disabled="pending" type="submit" />
    </form>

    <MessageError :text="errorText" />
    <p v-if="successText" class="text-sm text-(--accent)">{{ successText }}</p>
  </main>
</template>
