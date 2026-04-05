<script setup lang="ts">
  const title = 'Особенности'
  definePageMeta({
    robots: 'noindex, nofollow',
    key: route => route.path
  })
  usePageSeo({
    title,
    description: 'Раздел для просмотра и управления особенностями.'
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
  const breadcrumbItems = computed<BreadcrumbItem[]>(() => {
    if (activeWorkspaceTab.value === 'saved') {
      return [
        { label: title, to: '/traits' },
        { label: 'Сохранённые наборы', current: true, kind: 'tab' }
      ]
    }
    if (isPrimaryTraitOpen.value) {
      return [
        { label: title, to: '/traits' },
        { label: '', current: true, kind: 'tab', badge: 'Основной' }
      ]
    }
    return [
      {
        label: title,
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
    <LabNavHeader :title :breadcrumb-items="breadcrumbItems">
      <template #actions="{ compact }">
        <TraitsUuidButton v-if="uuid" action="copy" :uuid="uuid" :compact="compact" />
        <TraitsUuidButton action="paste" :compact="compact" @click="pasteUuid" />
      </template>
    </LabNavHeader>
    <LabNotify :text="pasteError" tone="error" class-name="px-4" />
    <TraitsWorkspaceTabs :current-trait-uuid="uuid || ''" />
  </div>
</template>
