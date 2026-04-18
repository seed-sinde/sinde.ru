<template>
  <div class="relative">
    <LabBaseInput
      :id="inputId"
      v-model="valueText"
      name="TraitInputNumber"
      type="number"
      inputmode="decimal"
      step="any"
      :aria-label="ariaLabel"
    />
    <span
      v-if="unitLabel"
      aria-hidden="true"
      class="pointer-events-none absolute inset-y-0 right-12 left-3.5 flex items-center overflow-hidden text-sm whitespace-nowrap"
    >
      <span v-if="valueText" class="invisible shrink-0">{{ valueText }}</span>
      <span :class="valueText ? 'pl-1' : ''" class="truncate text-(--lab-text-soft)">
        {{ unitLabel }}
      </span>
    </span>
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
  set: next => {
    model.value = next
  }
})
</script>
