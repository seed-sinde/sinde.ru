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
  <div v-if="(showControls && !collapsed) || showLinks" class="min-w-0 flex-1 basis-72">
    <div v-if="showLinks" class="flex flex-wrap items-center gap-x-3 gap-y-1">
      <NuxtLink v-if="showAdminLinks" to="/auth/admin" :class="linkClass">
        {{ t('nav.admin') }}
      </NuxtLink>
      <NuxtLink to="/docs/company" :class="linkClass">
        {{ t('nav.company') }}
      </NuxtLink>
      <NuxtLink to="/docs/offer" :class="linkClass">
        {{ t('nav.offer') }}
      </NuxtLink>
      <NuxtLink to="/docs/terms" :class="linkClass">
        {{ t('nav.terms_of_use') }}
      </NuxtLink>
      <NuxtLink to="/docs/privacy" :class="linkClass">
        {{ t('nav.privacy') }}
      </NuxtLink>
      <NuxtLink to="mailto:bug@sinde.ru" external :class="linkClass">bug@sinde.ru</NuxtLink>
    </div>
  </div>
</template>
