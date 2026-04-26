<script setup lang="ts">
type ApiError = Error & {
  data?: {
    message?: string
  }
}

const errorMessage = (err: unknown, fallback: string) =>
  (err as ApiError)?.data?.message || (err as Error)?.message || fallback
const buildLoginPath = (next: string) => `/auth/login?next=${encodeURIComponent(next)}`

definePageMeta({
  validate: route => {
    const account = Array.isArray(route.params.account)
      ? route.params.account[0]
      : route.params.account
    return account === "profile"
  }
})

useSeoMeta({
  title: "Аккаунт",
  description: "Профиль аккаунта."
})

const route = useRoute()
const router = useRouter()
const {user, ensureLoaded, updateProfile, logout} = useAuth()

await ensureLoaded()
if (!user.value) {
  await navigateTo(buildLoginPath(route.fullPath))
}

const form = reactive({
  display_name: user.value?.display_name || "",
  locale: user.value?.locale || "ru",
  timezone: user.value?.timezone || "UTC"
})
const pending = ref(false)
const logoutPending = ref(false)
const errorText = ref("")
const successText = ref("")

watch(
  user,
  next => {
    form.display_name = next?.display_name || ""
    form.locale = next?.locale || "ru"
    form.timezone = next?.timezone || "UTC"
  },
  {immediate: true}
)

const saveProfile = async () => {
  pending.value = true
  errorText.value = ""
  successText.value = ""

  try {
    await updateProfile(form)
    successText.value = "Профиль обновлён."
  } catch (err: unknown) {
    errorText.value = errorMessage(err, "Не удалось обновить профиль.")
  } finally {
    pending.value = false
  }
}

const submitLogout = async () => {
  logoutPending.value = true
  errorText.value = ""

  try {
    await logout()
    await router.push("/auth/login")
  } catch (err: unknown) {
    errorText.value = errorMessage(err, "Не удалось выйти.")
  } finally {
    logoutPending.value = false
  }
}
</script>

<template>
  <main class="mx-auto flex w-full max-w-2xl flex-col gap-5 p-4">
    <header class="space-y-1">
      <h1 class="text-2xl font-semibold">Аккаунт</h1>
      <p class="text-sm opacity-75">{{ user?.email }}</p>
    </header>

    <form class="flex flex-col gap-3" @submit.prevent="saveProfile">
      <UiField label="Имя" for-id="account-name">
        <UiInput id="account-name" v-model="form.display_name" />
      </UiField>
      <UiField label="Локаль" for-id="account-locale">
        <UiInput id="account-locale" v-model="form.locale" />
      </UiField>
      <UiField label="Часовой пояс" for-id="account-timezone">
        <UiInput id="account-timezone" v-model="form.timezone" />
      </UiField>
      <div class="flex flex-wrap gap-2">
        <UiButton type="submit" label="Сохранить" :loading="pending" :disabled="pending" />
        <UiButton
          type="button"
          label="Выйти"
          :loading="logoutPending"
          :disabled="logoutPending"
          @click="submitLogout"
        />
      </div>
    </form>

    <MessageError :text="errorText" />
    <p v-if="successText" class="text-sm text-(--accent)">{{ successText }}</p>
  </main>
</template>
