<template>
  <LabBaseSwitch
    id="TraitInputBoolean"
    v-model="model"
    :label="displayLabel"
    :true-label="trueLabel"
    :false-label="falseLabel"
    tone="emerald"
  />
</template>
<script setup lang="ts">
const props = defineProps<{ meta?: Pick<KeyMeta, "booleanDisplay" | "booleanTrueLabel" | "booleanFalseLabel"> }>()
const model = defineModel<boolean>({ required: true })
const defaults = computed(() => getDefaultBooleanLabels(props.meta?.booleanDisplay))
const displayLabel = computed(() => defaults.value.label)
const trueLabel = computed(() => String(props.meta?.booleanTrueLabel || defaults.value.trueLabel).trim() || defaults.value.trueLabel)
const falseLabel = computed(() => String(props.meta?.booleanFalseLabel || defaults.value.falseLabel).trim() || defaults.value.falseLabel)
</script>
