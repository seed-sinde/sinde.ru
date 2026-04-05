<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      currentTraitUuid?: string | null
    }>(),
    {
      currentTraitUuid: ''
    }
  )
  const activeTab = ref<'traits' | 'saved'>('traits')
  const tabItems = [
    { value: 'traits', label: 'Особенности' },
    { value: 'saved', label: 'Сохранённые наборы' }
  ]
</script>
<template>
  <LabNavTabs v-model="activeTab" :items="tabItems" route-query-key="tab" route-default-value="traits">
    <template #panel-traits>
      <TraitsManager />
    </template>
    <template #panel-saved>
      <TraitsLibraryPanel :current-trait-uuid="props.currentTraitUuid" />
    </template>
  </LabNavTabs>
</template>
