<template>
  <div class="space-y-1">
    <div class="flex flex-col gap-2">
      <LabBaseField label="Начало" :for-id="startId">
        <TraitsFormDateTextInput
          :id="startId"
          v-model="start"
          :placeholder="placeholder"
          :aria-label="startLabel"
          :invalid="Boolean(start) && !startValid"
        />
      </LabBaseField>
      <LabBaseField label="Конец" :for-id="endId">
        <TraitsFormDateTextInput
          :id="endId"
          v-model="end"
          :placeholder="placeholder"
          :aria-label="endLabel"
          :invalid="Boolean(end) && !endValid"
        />
      </LabBaseField>
    </div>
    <LabErrorMessage v-if="(start && !startValid) || (end && !endValid)" :text="`Неверный формат: ${placeholder}`" />
    <LabErrorMessage v-else-if="isRangeInvalid" text="Начало не может быть позже конца." />
  </div>
</template>
<script setup lang="ts">
import TraitsFormDateTextInput from './DateTextInput.vue'
const model = defineModel<TraitDateRangeValue>({ required: true })
const props = withDefaults(
  defineProps<{
    min?: string
    max?: string
    id?: string
    labelStart?: string
    labelEnd?: string
    meta?: unknown
  }>(),
  {}
)
const startId = computed(() => (props.id ? `${props.id}-start` : 'TraitInputDatetimeStart'))
const endId = computed(() => (props.id ? `${props.id}-end` : 'TraitInputDatetimeEnd'))
const startLabel = computed(() => props.labelStart || 'Начало')
const endLabel = computed(() => props.labelEnd || 'Конец')
const {
  left: start,
  right: end,
  leftValid: startValid,
  rightValid: endValid,
  placeholder,
  isRangeInvalid
} = useDatePairModel(model as any, {
  leftKey: 'start',
  rightKey: 'end'
})
</script>
