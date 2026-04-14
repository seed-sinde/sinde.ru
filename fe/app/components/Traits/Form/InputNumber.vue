<template>
  <div class="flex items-stretch">
    <LabBaseInput
      :id="inputId"
      v-model="valueText"
      name="TraitInputNumber"
      type="number"
      inputmode="decimal"
      step="any"
      :aria-label="ariaLabel"
    />
    <div v-if="unitLabel" class="relative flex items-center justify-center px-3 py-2.5 font-mono text-xs text-zinc-300">
      <span class="-t-2 absolute text-xs text-zinc-400" />
      {{ unitLabel }}
    </div>
  </div>
</template>
<script setup lang="ts">
const props = withDefaults(defineProps<{ id?: string; label?: string; meta?: { unit?: string } }>(), {
  label: 'Число'
})
const model = defineModel<number | string>({ required: true })
const inputId = computed(() => props.id || 'TraitInputNumber')
const ariaLabel = computed(() => props.label || 'Число')
const unitLabel = computed(() => props.meta?.unit || '')
const valueText = computed<string>({
  get: () => (model.value === null || model.value === undefined ? '' : String(model.value)),
  set: (next) => {
    model.value = next
  }
})
</script>
