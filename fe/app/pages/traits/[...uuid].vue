<script setup lang="ts">
const { localeCode } = useInterfacePreferences()
const copy = computed(() => TRAITS_WORKSPACE_COPY[localeCode.value] || TRAITS_WORKSPACE_COPY.ru)
definePageMeta({
  robots: 'noindex, nofollow',
  key: (route) => route.path
})
usePageSeo({
  title: () => copy.value.pageTitle,
  description: () => copy.value.pageDescriptionUuid
})
const route = useRoute()
const { user } = useAuth()
const { uuid, skipFetchUuid } = useTraitNavigation()
const { error: pasteError, pasteUuidAndNavigate } = usePasteUuid()
const pasteUuid = async () => await pasteUuidAndNavigate((uuid) => `/traits/${uuid}`)
const activeWorkspaceTab = computed<'traits' | 'saved'>(
  () => normalizeTabRouteValue(route.query.tab, ['traits', 'saved'], 'traits') as 'traits' | 'saved'
)
const isPrimaryTraitOpen = computed(() => {
  const currentUuid = String(uuid.value || '').trim()
  const primaryUuid = String(user.value?.primary_trait_uuid || '').trim()
  return Boolean(currentUuid) && currentUuid === primaryUuid
})
const breadcrumbItems = computed<BreadcrumbItem[]>(() => {
  if (activeWorkspaceTab.value === 'saved') {
    return [
      { label: copy.value.pageTitle, to: '/traits' },
      { label: copy.value.savedBreadcrumb, current: true, kind: 'tab' }
    ]
  }
  if (isPrimaryTraitOpen.value) {
    return [
      { label: copy.value.pageTitle, to: '/traits' },
      { label: '', current: true, kind: 'tab', badge: copy.value.primaryBadge }
    ]
  }
  return [
    {
      label: copy.value.pageTitle,
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
    <LabNavHeader :title="copy.pageTitle" :breadcrumb-items="breadcrumbItems">
      <template #actions="{ compact }">
        <TraitsUuidButton v-if="uuid" action="copy" :uuid="uuid" :compact="compact" />
        <TraitsUuidButton action="paste" :compact="compact" @click="pasteUuid" />
      </template>
    </LabNavHeader>
    <LabNotify :text="pasteError" tone="error" class-name="px-4" />
    <TraitsWorkspaceTabs :current-trait-uuid="uuid || ''" />
  </div>
</template>
