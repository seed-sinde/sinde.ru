<script setup lang="ts">
  const title = 'Особенности'
  usePageSeo({
    title,
    description: 'Раздел для просмотра и управления особенностями человека.'
  })
  const route = useRoute()
  const { error: pasteError, pasteUuidAndNavigate } = usePasteUuid()
  const activeWorkspaceTab = computed<'traits' | 'saved'>(
    () => normalizeTabRouteValue(route.query.tab, ['traits', 'saved'], 'traits') as 'traits' | 'saved'
  )
  const breadcrumbItems = computed<BreadcrumbItem[]>(() =>
    activeWorkspaceTab.value === 'saved'
      ? [
          { label: title, to: '/traits' },
          { label: 'Сохранённые наборы', current: true, kind: 'tab' }
        ]
      : [{ label: title, to: '/traits', current: true }]
  )
  const pasteUuid = async () => await pasteUuidAndNavigate(uuid => `/traits/${uuid}`)
</script>
<template>
  <div>
    <LabNavHeader :title :breadcrumb-items="breadcrumbItems">
      <template #actions="{ compact }">
        <TraitsUuidButton action="paste" :compact="compact" @click="pasteUuid" />
      </template>
    </LabNavHeader>
    <LabNotify :text="pasteError" tone="error" class-name="px-4" />
    <TraitsWorkspaceTabs />
  </div>
</template>
