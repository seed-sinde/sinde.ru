<script setup lang="ts">
definePageMeta({
  middleware: "guest-only"
})

type ApiError = Error & {
  data?: {
    message?: string
  }
}

const normalizeAuthNextPath = (value?: string | null) => {
  const raw = String(value || "").trim()
  if (!raw || !raw.startsWith("/") || raw.startsWith("//")) return ""
  return ["/auth/login", "/auth/register", "/auth/verify-email", "/auth/reset-password"].some(
    prefix => raw === prefix || raw.startsWith(`${prefix}/`)
  )
    ? ""
    : raw
}

const errorMessage = (err: unknown, fallback: string) =>
  (err as ApiError)?.data?.message || (err as Error)?.message || fallback

useSeoMeta({
  title: "Вход",
  description: "Вход в аккаунт."
})

const route = useRoute()
const router = useRouter()
const {login, completeMfa, mfaTicket} = useAuth()
const email = ref("")
const password = ref("")
const mfaCode = ref("")
const pending = ref(false)
const errorText = ref("")
const nextTarget = computed(
  () =>
    normalizeAuthNextPath(typeof route.query.next === "string" ? route.query.next : "") ||
    "/auth/account/profile"
)

const submitPassword = async () => {
  pending.value = true
  errorText.value = ""

  try {
    const res = await login(email.value, password.value)
    if (!res.data.mfa_required) await router.push(nextTarget.value)
  } catch (err: unknown) {
    errorText.value = errorMessage(err, "Не удалось войти.")
  } finally {
    pending.value = false
  }
}

const submitMfa = async () => {
  const code = mfaCode.value.replace(/\D+/g, "").slice(0, 6)
  mfaCode.value = code
  if (code.length !== 6) {
    errorText.value = "Введите 6 цифр кода 2FA."
    return
  }

  pending.value = true
  errorText.value = ""

  try {
    await completeMfa(code)
    await router.push(nextTarget.value)
  } catch (err: unknown) {
    errorText.value = errorMessage(err, "Неверный код 2FA.")
    mfaCode.value = ""
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <main class="p-4">
    <header class="space-y-1">
      <h1 class="text-2xl font-semibold">Вход</h1>
      <p class="text-sm opacity-75">Войдите по email и паролю.</p>
    </header>

    <form v-if="!mfaTicket" class="flex flex-col gap-3" @submit.prevent="submitPassword">
      <UiField label="Email" for-id="login-email">
        <UiInput id="login-email" v-model="email" type="email" autocomplete="username" />
      </UiField>
      <UiField label="Пароль" for-id="login-password">
        <UiInput
          id="login-password"
          v-model="password"
          type="password"
          autocomplete="current-password"
        />
      </UiField>
      <UiButton
        type="submit"
        label="Войти"
        :loading="pending"
        :disabled="pending || !email || !password"
      />
    </form>

    <form v-else class="flex flex-col gap-2" @submit.prevent="submitMfa">
      <UiField label="Код 2FA" for-id="login-mfa">
        <UiInput
          id="login-mfa"
          v-model="mfaCode"
          name="mfa_code"
          inputmode="numeric"
          autocomplete="one-time-code"
          class="w-full"
        />
      </UiField>
      <UiButton type="submit" label="Подтвердить" :loading="pending" :disabled="pending" />
    </form>

    <MessageError :text="errorText" />

    <nav class="flex gap-3 text-sm">
      <NuxtLink class="text-(--accent)" to="/auth/register">Регистрация</NuxtLink>
      <NuxtLink class="text-(--accent)" to="/auth/forgot-password">Забыли пароль?</NuxtLink>
    </nav>
  </main>
</template>
