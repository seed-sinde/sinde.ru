<script setup lang="ts">
const route = useRoute()
const { localeCode } = useInterfacePreferences()
const props = withDefaults(
  defineProps<{
    currentTraitUuid?: string | null
  }>(),
  {
    currentTraitUuid: ''
  }
)
const copy = computed(() => TRAITS_WORKSPACE_COPY[localeCode.value] || TRAITS_WORKSPACE_COPY.ru)
const activeTab = ref<'traits' | 'saved'>(
  normalizeTabRouteValue(route.query.tab, ['traits', 'saved'], 'traits') as 'traits' | 'saved'
)
const activatedTabs = reactive(new Set<'traits' | 'saved'>([activeTab.value]))
const tabItems = computed(() => [
  { value: 'traits', label: copy.value.tabs.traits },
  { value: 'saved', label: copy.value.tabs.saved }
])
watch(
  activeTab,
  (next) => {
    activatedTabs.add(next)
  },
  { immediate: true }
)
</script>
<template>
  <LabNavTabs v-model="activeTab" :items="tabItems" route-query-key="tab" route-default-value="traits">
    <template #panel-traits>
      <LazyTraitsManager v-if="activatedTabs.has('traits')" />
    </template>
    <template #panel-saved>
      <LazyTraitsLibraryPanel v-if="activatedTabs.has('saved')" :current-trait-uuid="props.currentTraitUuid" />
    </template>
  </LabNavTabs>
</template>
