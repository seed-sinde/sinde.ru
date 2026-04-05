<template>
  <LabField label="Значение" :for-id="inputId">
    <LabBaseInput
      :id="inputId"
      :list="listId"
      name="TraitSelectEnum"
      v-model="model"
      :placeholder="placeholder"
      :aria-label="ariaLabel" />
    <datalist :id="listId">
      <option v-for="opt in options" :key="opt" :value="opt"></option>
    </datalist>
  </LabField>
</template>
<script setup lang="ts">
  import { useId } from 'vue'
  const props = withDefaults(defineProps<{ meta: { options?: string[] }; id?: string; label?: string }>(), {
    id: undefined,
    label: 'Значение'
  })
  const model = defineModel<string>({ required: true })
  const options = computed(() => props.meta?.options ?? [])
  const uid = useId()
  const baseId = computed(() => props.id || `TraitEnum-${uid}`)
  const inputId = computed(() => baseId.value)
  const listId = computed(() => `${baseId.value}-list`)
  const placeholder = computed(() => {
    return options.value.length ? 'Выберите или введите свое значение' : 'Введите значение'
  })
  const ariaLabel = computed(() => props.label || 'Значение')
</script>
