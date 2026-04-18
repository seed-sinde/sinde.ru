<script setup lang="ts">
useTraitsMobileHeader()
const { locale, key, load, t } = useI18nSection('traits')
await useAsyncData(key.value, load, { watch: [locale] })
const pageTitle = computed(() => t('page.title'))
const pageDescription = computed(() => t('page.description.uuid'))
definePageMeta({
  robots: 'noindex, nofollow',
  key: route => route.path
})
usePageSeo({
  title: pageTitle,
  description: pageDescription
})
const route = useRoute()
const { user } = useAuth()
const { uuid, skipFetchUuid } = useTraitNavigation()
const { error: pasteError, pasteUuidAndNavigate } = usePasteUuid()
const pasteUuid = async () => await pasteUuidAndNavigate(uuid => `/traits/${uuid}`)
const activeWorkspaceTab = computed<'traits' | 'saved'>(
  () => normalizeTabRouteValue(route.query.tab, ['traits', 'saved'], 'traits') as 'traits' | 'saved'
)
const isPrimaryTraitOpen = computed(() => {
  const currentUuid = String(uuid.value || '').trim()
  const primaryUuid = String(user.value?.primary_trait_uuid || '').trim()
  return Boolean(currentUuid) && currentUuid === primaryUuid
})
const mobileActions = computed<MobileHeaderAction[]>(() => [
  ...(uuid.value ? [{ kind: 'traits-copy-uuid' as const, uuid: uuid.value }] : []),
  { kind: 'traits-paste-uuid', mode: 'set' }
])
const breadcrumbItems = computed<BreadcrumbItem[]>(() => {
  if (activeWorkspaceTab.value === 'saved') {
    return [
      { label: t('page.title'), to: '/traits' },
      { label: t('saved.breadcrumb'), current: true, kind: 'tab' }
    ]
  }
  if (isPrimaryTraitOpen.value) {
    return [
      { label: t('page.title'), to: '/traits' },
      { label: '', current: true, kind: 'tab', badge: t('primary_badge') }
    ]
  }
  return [
    {
      label: t('page.title'),
      ...(uuid.value ? { to: '/traits' } : { current: true })
    }
  ]
})
if (import.meta.client) {
  useTraitStream(uuid, skipFetchUuid)
}
</script>
<template>
  <div>
    <LabNavHeader :title="pageTitle" :breadcrumb-items="breadcrumbItems" :mobile-actions="mobileActions">
      <template #actions="{ compact }">
        <TraitsUuidButton v-if="uuid" action="copy" :uuid="uuid" :compact="compact" />
        <TraitsUuidButton action="paste" :compact="compact" @click="pasteUuid" />
      </template>
    </LabNavHeader>
    <LabNotify :text="pasteError" tone="error" class-name="px-4" />
    <TraitsWorkspaceTabs :current-trait-uuid="uuid || ''" />
  </div>
</template>
