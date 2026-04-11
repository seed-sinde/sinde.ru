<script setup lang="ts">
  const { t } = useInterfacePreferences()
  const props = defineProps<{
    message: string
  }>()
  const route = useRoute()
  const loginTo = computed(() => buildLoginPath(route.fullPath))
  const registerTo = computed(() => {
    const next = normalizeAuthNextPath(route.fullPath)
    if (!next) return '/auth/register'
    return `/auth/register?next=${encodeURIComponent(next)}`
  })
  const goToLogin = () => navigateTo(loginTo.value)
  const goToRegister = () => navigateTo(registerTo.value)
</script>
<template>
  <div class="space-y-3 text-sm p-4">
    <p class="text-(--lab-text-muted)">{{ message }}</p>
    <div class="flex flex-wrap items-center gap-3">
      <LabBaseButton :label="t('auth.login.submit')" variant="primary" @click="goToLogin" />
      <LabBaseButton :label="t('auth.login.register')" variant="secondary" @click="goToRegister" />
    </div>
  </div>
</template>
