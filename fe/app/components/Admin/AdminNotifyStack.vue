<script setup lang="ts">
type NotifyTone = 'error' | 'success' | 'warning' | 'info'

type NotifyItem = {
  text?: string | null
  tone?: NotifyTone
}

const props = defineProps<{
  items: NotifyItem[]
}>()

const visibleItems = computed(() =>
  props.items.flatMap(item => {
    if (typeof item.text !== 'string') return []
    const text = item.text.trim()
    if (!text) return []

    return [
      {
        text,
        ...(item.tone ? { tone: item.tone } : {})
      }
    ]
  })
)
</script>

<template>
  <div v-if="visibleItems.length" class="space-y-1">
    <LabNotify v-for="(item, i) in visibleItems" :key="i" v-bind="item" size="xs" />
  </div>
</template>
