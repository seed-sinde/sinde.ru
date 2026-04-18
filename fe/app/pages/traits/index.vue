<script setup lang="ts">
useTraitsStore().clear()
useTraitsMobileHeader()
const { locale, key, load, t } = useI18nSection('traits')
await useAsyncData(key.value, load, { watch: [locale] })
const route = useRoute()
const { error: pasteError, pasteUuidAndNavigate } = usePasteUuid()
const activeWorkspaceTab = computed<'traits' | 'saved'>(
  () => normalizeTabRouteValue(route.query.tab, ['traits', 'saved'], 'traits') as 'traits' | 'saved'
)
const pageTitle = computed(() => t('page.title'))
const pageDescription = computed(() => t('page.description.root'))
const mobileActions = computed<MobileHeaderAction[]>(() => [{ kind: 'traits-paste-uuid', mode: 'set' }])
const breadcrumbItems = computed<BreadcrumbItem[]>(() =>
  activeWorkspaceTab.value === 'saved'
    ? [
        { label: t('page.title'), to: '/traits' },
        { label: t('saved.breadcrumb'), current: true, kind: 'tab' }
      ]
    : [{ label: t('page.title'), to: '/traits', current: true }]
)
const pasteUuid = async () => await pasteUuidAndNavigate(uuid => `/traits/${uuid}`)
usePageSeo({ title: pageTitle, description: pageDescription })
</script>
<template>
  <div>
    <LabNavHeader :title="pageTitle" :breadcrumb-items="breadcrumbItems" :mobile-actions="mobileActions">
      <template #actions="{ compact }">
        <TraitsUuidButton action="paste" :compact="compact" @click="pasteUuid" />
      </template>
    </LabNavHeader>
    <LabNotify :text="pasteError" tone="error" :temporary="true" class-name="px-4" />
    <TraitsWorkspaceTabs />
  </div>
</template>
