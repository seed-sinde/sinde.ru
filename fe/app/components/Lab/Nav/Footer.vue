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
const { locale, key, load, t } = useI18nSection('nav')
await useAsyncData(key.value, load, { watch: [locale] })
const showAdminLinks = computed(() => loaded.value && isAdmin.value)
const linkClass = 'text-(--lab-text-soft) hover:text-(--lab-text-primary) shrink-0 text-xs transition-colors'
const links = computed(() =>[
  { to: '/docs/company', label: t('company') },
  { to: '/docs/offer', label: t('offer') },
  { to: '/docs/terms', label: t('terms_of_use') },
  { to: '/docs/privacy', label: t('privacy') },
  { to: 'mailto:bug@sinde.ru', label: 'bug@sinde.ru', external: true }
])
</script>
<template>
  <div v-if="(showControls && !collapsed) || showLinks" class="min-w-0 flex-1 basis-72">
    <div v-if="showLinks" class="space-y-2">
      <div v-if="showAdminLinks" class="flex w-fit flex-col gap-2 px-2">
        <NuxtLink to="/auth/admin/users" :class="linkClass">
          {{ t('admin') }}
        </NuxtLink>
        <AdminIPStat />
      </div>
      <div class="flex flex-wrap items-center gap-x-3 gap-y-1">
        <LabLocaleSwitcher />
        <LabThemeSwitcher />
        <NuxtLink v-for="l in links" :key="l.to" :to="l.to" :external="l.external ?? false" :class="linkClass">
          {{ l.label }}
        </NuxtLink>
      </div>
    </div>
  </div>
</template>
