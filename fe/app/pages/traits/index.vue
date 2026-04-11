<script setup lang="ts">
  const { localeCode } = useInterfacePreferences()
  const copy = computed(() => TRAITS_WORKSPACE_COPY[localeCode.value] || TRAITS_WORKSPACE_COPY.ru)
  usePageSeo({
    title: () => copy.value.pageTitle,
    description: () => copy.value.pageDescriptionRoot
  })
  const route = useRoute()
  const { error: pasteError, pasteUuidAndNavigate } = usePasteUuid()
  const activeWorkspaceTab = computed<'traits' | 'saved'>(
    () => normalizeTabRouteValue(route.query.tab, ['traits', 'saved'], 'traits') as 'traits' | 'saved'
  )
  const breadcrumbItems = computed<BreadcrumbItem[]>(() =>
    activeWorkspaceTab.value === 'saved'
      ? [
          { label: copy.value.pageTitle, to: '/traits' },
          { label: copy.value.savedBreadcrumb, current: true, kind: 'tab' }
        ]
      : [{ label: copy.value.pageTitle, to: '/traits', current: true }]
  )
  const pasteUuid = async () => await pasteUuidAndNavigate(uuid => `/traits/${uuid}`)
</script>
<template>
  <div>
    <LabNavHeader :title="copy.pageTitle" :breadcrumb-items="breadcrumbItems">
      <template #actions="{ compact }">
        <TraitsUuidButton action="paste" :compact="compact" @click="pasteUuid" />
      </template>
    </LabNavHeader>
    <LabNotify :text="pasteError" tone="error" class-name="px-4" />
    <TraitsWorkspaceTabs />
  </div>
</template>
