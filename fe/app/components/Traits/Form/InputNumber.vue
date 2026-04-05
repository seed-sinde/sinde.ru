<template>
  <div class="flex items-stretch">
    <LabBaseInput
      :id="inputId"
      name="TraitInputNumber"
      type="number"
      v-model="valueText"
      inputmode="decimal"
      step="any"
      :aria-label="ariaLabel" />
    <div v-if="unitLabel" class="relative flex items-center justify-center px-3 py-2.5 text-xs font-mono text-zinc-300">
      <span class="absolute -t-2 text-xs text-zinc-400"></span>
      {{ unitLabel }}
    </div>
  </div>
</template>
<script setup lang="ts">
  const props = withDefaults(defineProps<{ id?: string; label?: string; meta?: { unit?: string } }>(), {
    id: undefined,
    label: 'Число',
    meta: undefined
  })
  const model = defineModel<number | string>({ required: true })
  const inputId = computed(() => props.id || 'TraitInputNumber')
  const ariaLabel = computed(() => props.label || 'Число')
  const unitLabel = computed(() => props.meta?.unit || '')
  const valueText = computed<string>({
    get: () => (model.value === null || model.value === undefined ? '' : String(model.value)),
    set: next => {
      model.value = next
    }
  })
</script>
