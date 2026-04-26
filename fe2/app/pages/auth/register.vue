<script setup lang="ts">
definePageMeta({
  middleware: "guest-only"
})

type ApiError = Error & {
  data?: {
    message?: string
  }
}

const errorMessage = (err: unknown, fallback: string) =>
  (err as ApiError)?.data?.message || (err as Error)?.message || fallback
const isValidEmail = (value: string) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value.trim())

useSeoMeta({
  title: "Регистрация",
  description: "Регистрация нового аккаунта."
})

const router = useRouter()
const {register} = useAuth()
const form = reactive({
  email: "",
  password: "",
  display_name: "",
  locale: "ru",
  timezone: "UTC"
})
const pending = ref(false)
const errorText = ref("")
const successText = ref("")
const canSubmit = computed(() => isValidEmail(form.email) && form.password.length >= 12)

onMounted(() => {
  form.timezone = Intl.DateTimeFormat().resolvedOptions().timeZone || "UTC"
})

const submit = async () => {
  if (!canSubmit.value) {
    errorText.value = "Укажите корректный email и пароль от 12 символов."
    successText.value = ""
    return
  }

  pending.value = true
  errorText.value = ""
  successText.value = ""

  try {
    const res = await register(form)
    successText.value = `Аккаунт создан. Проверьте email в течение ${res.data.verification_ttl}.`
    setTimeout(() => router.push("/auth/login"), 1200)
  } catch (err: unknown) {
    errorText.value = errorMessage(err, "Регистрация не удалась.")
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <main class="p-4">
    <header class="space-y-1">
      <h1 class="text-2xl font-semibold">Регистрация</h1>
      <p class="text-sm opacity-75">Создайте аккаунт и подтвердите email.</p>
    </header>

    <form class="flex flex-col items-start gap-2" @submit.prevent="submit">
      <UiField label="Имя" for-id="register-name" hint="Отображаемое публично">
        <UiInput id="register-name" v-model="form.display_name" placeholder="Имя" />
      </UiField>
      <UiField label="Email" for-id="register-email" hint="Введите электронную почту">
        <UiInput
          id="register-email"
          v-model="form.email"
          placeholder="seed@sinde.ru"
          type="email"
          autocomplete="email"
        />
      </UiField>
      <UiField label="Пароль" hint="Минимум 12 символов" for-id="register-password">
        <UiInput
          id="register-password"
          v-model="form.password"
          type="password"
          autocomplete="new-password"
        />
      </UiField>
      <UiButton
        type="submit"
        label="Зарегистрироваться"
        :loading="pending"
        :disabled="pending || !canSubmit"
      />
    </form>

    <MessageError :text="errorText" />
    <p v-if="successText" class="text-sm text-(--accent)">{{ successText }}</p>
    <NuxtLink class="text-sm text-(--accent)" to="/auth/login">Уже есть аккаунт</NuxtLink>
  </main>
</template>
