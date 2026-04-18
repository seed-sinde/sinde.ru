<script setup lang="ts">
const { locale, key, load, t } = useI18nSection('traits')
await useAsyncData(key.value, load, { watch: [locale] })
const route = useRoute()
const props = withDefaults(
  defineProps<{
    currentTraitUuid?: string | null
  }>(),
  {
    currentTraitUuid: ''
  }
)
const activeTab = computed<'traits' | 'saved'>({
  get: () => normalizeTabRouteValue(route.query.tab, ['traits', 'saved'], 'traits') as 'traits' | 'saved',
  set: () => {}
})
const activatedTabs = reactive(new Set<'traits' | 'saved'>([activeTab.value]))
const tabItems = computed(() => [
  { value: 'traits', label: t('tabs.traits') },
  { value: 'saved', label: t('tabs.saved') }
])
watch(
  activeTab,
  next => {
    activatedTabs.add(next)
  },
  { immediate: true }
)
</script>
<template>
  <LabNavTabs v-model="activeTab" :items="tabItems" route-query-key="tab">
    <template #panel-traits>
      <LazyTraitsManager v-if="activatedTabs.has('traits')" />
    </template>
    <template #panel-saved>
      <LazyTraitsLibraryPanel v-if="activatedTabs.has('saved')" :current-trait-uuid="props.currentTraitUuid" />
    </template>
  </LabNavTabs>
</template>
