<template>
  <div class="space-y-1">
    <div class="grid gap-2">
      <LabField label="Начало" :for-id="startId">
        <LabBaseInput
          v-if="rangeType === 'number'"
          :id="startId"
          v-model="model.start"
          type="number"
          inputmode="decimal"
          :invalid="Boolean(model.start) && !startValid"
        />
        <TraitsFormDateTextInput
          v-else
          :id="startId"
          v-model="model.start"
          :placeholder="placeholder"
          :invalid="Boolean(model.start) && !startValid"
        />
      </LabField>
      <LabField label="Конец" :for-id="endId">
        <LabBaseInput
          v-if="rangeType === 'number'"
          :id="endId"
          v-model="model.end"
          type="number"
          inputmode="decimal"
          :invalid="Boolean(model.end) && !endValid"
        />
        <TraitsFormDateTextInput
          v-else
          :id="endId"
          v-model="model.end"
          :placeholder="placeholder"
          :invalid="Boolean(model.end) && !endValid"
        />
      </LabField>
    </div>
    <LabErrorMessage v-if="hasFormatError" :text="formatError" />
    <LabErrorMessage v-else-if="isRangeInvalid" text="Начало не может быть больше конца." />
  </div>
</template>

<script setup lang="ts">
import TraitsFormDateTextInput from "./DateTextInput.vue"

const props = withDefaults(
  defineProps<{
    id?: string
    meta?: Pick<KeyMeta, "rangeType">
  }>(),
  {
    id: ""
  }
)
const model = defineModel<TraitRangeModel>({required: true})
const startId = computed(() => (props.id ? `${props.id}-start` : "TraitRangeStart"))
const endId = computed(() => (props.id ? `${props.id}-end` : "TraitRangeEnd"))
const rangeType = computed(() => resolveRangeValueType(props.meta?.rangeType))
const placeholder = computed(() => {
  if (rangeType.value === "time") return datePlaceholder("time")
  if (rangeType.value === "date") return datePlaceholder("date")
  if (rangeType.value === "datetime") return datePlaceholder("datetime")
  return ""
})
const isValidValue = (value: string) => {
  if (!value) return false
  if (rangeType.value === "number") return !Number.isNaN(Number(value))
  return isValidDateText(value, rangeType.value === "datetime" ? "datetime" : rangeType.value)
}
const startValid = computed(() => !model.value.start || isValidValue(model.value.start))
const endValid = computed(() => !model.value.end || isValidValue(model.value.end))
const isRangeInvalid = computed(() => {
  const start = String(model.value.start || "").trim()
  const end = String(model.value.end || "").trim()
  if (!start || !end || !startValid.value || !endValid.value) return false
  if (rangeType.value === "number") return Number(start) > Number(end)
  const mode = rangeType.value === "datetime" ? "datetime" : rangeType.value
  const cmp = compareDateText(start, end, mode)
  return cmp !== null && cmp > 0
})
const hasFormatError = computed(() => (Boolean(model.value.start) && !startValid.value) || (Boolean(model.value.end) && !endValid.value))
const formatError = computed(() =>
  rangeType.value === "number" ? "Введите число." : `Неверный формат: ${placeholder.value}`
)
</script>
