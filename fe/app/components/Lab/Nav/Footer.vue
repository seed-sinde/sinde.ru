<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      collapsed?: boolean
      showControls?: boolean
      showLinks?: boolean
    }>(),
    {
      collapsed: true,
      showControls: true,
      showLinks: true
    }
  )
  const { collapsed, showControls, showLinks } = toRefs(props)
  const { loaded, isAdmin } = useAuth()
  const { t } = useInterfacePreferences()
  const showAdminLinks = computed(() => loaded.value && isAdmin.value)
  const linkClass = 'text-(--lab-text-soft) hover:text-(--lab-text-primary) shrink-0 text-xs transition-colors'
</script>
<template>
  <div v-if="(showControls && !collapsed) || showLinks" class="space-y-2">
    <div v-if="showLinks" class="flex flex-wrap items-center gap-x-3 gap-y-1">
      <NuxtLink v-if="showAdminLinks" to="/auth/admin" :class="linkClass">
        {{ t('nav.admin') }}
      </NuxtLink>
      <NuxtLink to="/company" :class="linkClass">
        {{ t('nav.company') }}
      </NuxtLink>
      <NuxtLink to="/offer" :class="linkClass">
        {{ t('nav.offer') }}
      </NuxtLink>
      <NuxtLink to="/terms" :class="linkClass">
        {{ t('nav.code_of_conduct') }}
      </NuxtLink>
      <NuxtLink to="mailto:bug@sinde.ru" :class="linkClass">bug@sinde.ru</NuxtLink>
    </div>
  </div>
</template>
