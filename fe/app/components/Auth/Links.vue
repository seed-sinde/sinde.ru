<script setup lang="ts">
  const { t } = useInterfacePreferences()
  const props = withDefaults(
    defineProps<{
      variant?: 'buttons' | 'inline'
      containerClass?: string
      loginLabel?: string
      registerLabel?: string
    }>(),
    {
      variant: 'buttons',
      containerClass: '',
      loginLabel: '',
      registerLabel: ''
    }
  )
  const resolvedContainerClass = computed(() => {
    if (props.containerClass) return props.containerClass
    return props.variant === 'buttons' ? 'flex flex-wrap gap-2' : 'flex gap-3 text-xs'
  })
  const resolvedLoginClass = computed(() =>
    props.variant === 'buttons'
      ? 'inline-flex h-10 items-center justify-center bg-emerald-700 px-4 text-sm font-medium text-white transition hover:bg-emerald-600'
      : 'text-emerald-300 hover:underline'
  )
  const resolvedRegisterClass = computed(() =>
    props.variant === 'buttons'
      ? 'inline-flex h-10 items-center justify-center border border-zinc-700 bg-zinc-900 px-4 text-sm font-medium text-zinc-100 transition hover:bg-zinc-800'
      : 'text-cyan-300 hover:underline'
  )
  const route = useRoute()
  const loginTo = computed(() => buildLoginPath(route.fullPath))
  const resolvedLoginLabel = computed(() => props.loginLabel || t('auth.login.submit'))
  const resolvedRegisterLabel = computed(() => props.registerLabel || t('auth.login.register'))
</script>
<template>
  <div :class="resolvedContainerClass">
    <NuxtLink :to="loginTo" :class="resolvedLoginClass">{{ resolvedLoginLabel }}</NuxtLink>
    <NuxtLink to="/auth/register" :class="resolvedRegisterClass">{{ resolvedRegisterLabel }}</NuxtLink>
  </div>
</template>
