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
</script>
<template>
  <div class="border border-zinc-800/80 bg-zinc-900/55 px-4 py-3 text-sm text-zinc-400">
    <p>{{ message }}</p>
    <div class="mt-2 flex flex-wrap items-center gap-3">
      <NuxtLink :to="loginTo" class="inline-flex text-sm transition">
        {{ t('auth.login.submit') }}
      </NuxtLink>
      <NuxtLink :to="registerTo" class="inline-flex text-sm transition">
        {{ t('auth.login.register') }}
      </NuxtLink>
    </div>
  </div>
</template>
