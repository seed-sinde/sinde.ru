<template>
  <LabNotify :text="resolvedText" tone="error" size="xs" :class-name="[errorClass, 'text-xs'].join(' ')" />
</template>
<script setup lang="ts">
const props = withDefaults(defineProps<{ text?: string | string[] | null; errorClass?: string }>(), {
  text: '',
  errorClass: ''
})
const resolvedText = computed(() => {
  if (Array.isArray(props.text)) {
    return props.text.map((item) => String(item || '').trim()).find(Boolean) || ''
  }
  return String(props.text || '').trim()
})
</script>
