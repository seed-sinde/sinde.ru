<script setup lang="ts">
const props = defineProps<{
  label: string
  value: string | number | null | undefined
  containerClass?: string
}>()

const displayValue = computed(() => {
  if (props.value === null || props.value === undefined || props.value === '') return '—'
  return props.value
})
</script>

<template>
  <div :class="props.containerClass">
    <p class="text-zinc-500">{{ props.label }}</p>
    <p class="text-zinc-100">{{ displayValue }}</p>
  </div>
</template>
