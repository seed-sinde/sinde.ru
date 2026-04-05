<template>
  <LabBaseInput
    :id="inputId"
    name="TraitInputString"
    v-model="model"
    :placeholder="placeholderText"
    :aria-label="ariaLabel" />
</template>
<script setup lang="ts">
  const props = withDefaults(defineProps<{ label?: string; placeholder?: string; id?: string; meta?: unknown }>(), {
    label: 'Значение',
    placeholder: 'Значение',
    id: undefined,
    meta: undefined
  })
  const model = defineModel<string>({ required: true })
  const inputId = computed(() => props.id || 'TraitInputString')
  const placeholderText = computed(() => props.placeholder || '')
  const ariaLabel = computed(() => props.label || props.placeholder || 'Значение')
</script>
