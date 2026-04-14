<template>
  <div class="space-y-1">
    <TraitsFormDateTextInput
      :id="inputId"
      v-model="model"
      name="TraitInputDatetime"
      :placeholder="placeholder"
      :aria-label="ariaLabel"
      :invalid="Boolean(model) && !isValid"
    />
    <LabErrorMessage v-if="model && !isValid" :text="`Неверный формат: ${placeholder}`" />
  </div>
</template>
<script setup lang="ts">
import TraitsFormDateTextInput from './DateTextInput.vue'
const model = defineModel<string>({ required: true })
const props = withDefaults(defineProps<{ min?: string; max?: string; id?: string; label?: string; meta?: unknown }>(), {
  label: 'Дата/время'
})
const placeholder = computed(() => datePlaceholder('datetime'))
const inputId = computed(() => props.id || 'TraitInputDatetime')
const ariaLabel = computed(() => props.label || 'Дата/время')
const isValid = computed(() => !model.value || isValidDateText(model.value, 'datetime'))
</script>
