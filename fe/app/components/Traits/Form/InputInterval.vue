<template>
  <div class="space-y-1">
    <div class="flex flex-col gap-2">
      <LabBaseField label="Начало события" :for-id="startId">
        <TraitsFormDateTextInput
          :id="startId"
          v-model="start"
          :placeholder="dateTimePlaceholder"
          :aria-label="startLabel"
          :invalid="Boolean(start) && !startValid"
        />
      </LabBaseField>
      <LabBaseField label="Конец события" :for-id="endId">
        <TraitsFormDateTextInput
          :id="endId"
          v-model="end"
          :placeholder="dateTimePlaceholder"
          :aria-label="endLabel"
          :invalid="Boolean(end) && !endValid"
        />
      </LabBaseField>
      <LabBaseField label="Единица длительности" :for-id="unitId">
        <LabBaseSelect :id="unitId" v-model="unit" :aria-label="unitLabel" :options="intervalUnitOptions" />
      </LabBaseField>
    </div>
    <LabErrorMessage
      v-if="(start && !startValid) || (end && !endValid)"
      :text="`Неверный формат: ${dateTimePlaceholder}`"
    />
    <LabErrorMessage v-else-if="isRangeInvalid" text="Начало события не может быть позже конца." />
    <p v-else-if="durationPreview" class="text-xs text-zinc-400">
      {{ durationPreview }}
    </p>
  </div>
</template>
<script setup lang="ts">
import TraitsFormDateTextInput from './DateTextInput.vue'
const model = defineModel<TraitIntervalModel>({ required: true })
const props = withDefaults(
  defineProps<{
    id?: string
    labelStart?: string
    labelEnd?: string
    labelUnit?: string
    meta?: unknown
  }>(),
  {}
)
const intervalUnits: Array<{ value: TraitIntervalUnit; label: string; ms: number; short: string }> = [
  { value: 'seconds', label: 'Секунды', ms: 1000, short: 'сек' },
  { value: 'minutes', label: 'Минуты', ms: 60_000, short: 'мин' },
  { value: 'hours', label: 'Часы', ms: 3_600_000, short: 'ч' },
  { value: 'days', label: 'Дни', ms: 86_400_000, short: 'дн' },
  { value: 'years', label: 'Годы', ms: 31_536_000_000, short: 'г' }
]
const intervalUnitOptions = intervalUnits.map(unit => ({ value: unit.value, label: unit.label }))
const startId = computed(() => (props.id ? `${props.id}-start` : 'TraitIntervalStart'))
const endId = computed(() => (props.id ? `${props.id}-end` : 'TraitIntervalEnd'))
const unitId = computed(() => (props.id ? `${props.id}-unit` : 'TraitIntervalUnit'))
const startLabel = computed(() => props.labelStart || 'Начало события')
const endLabel = computed(() => props.labelEnd || 'Конец события')
const unitLabel = computed(() => props.labelUnit || 'Единица длительности')
const dateTimePlaceholder = computed(() => datePlaceholder('datetime'))
const start = ref(String(model.value?.start ?? ''))
const end = ref(String(model.value?.end ?? ''))
const unit = ref<TraitIntervalUnit>(
  intervalUnits.some(opt => opt.value === model.value?.unit) ? (model.value.unit as TraitIntervalUnit) : 'minutes'
)
const startValid = computed(() => !start.value || isValidDateText(start.value, 'datetime'))
const endValid = computed(() => !end.value || isValidDateText(end.value, 'datetime'))
const isRangeInvalid = computed(() => {
  if (!start.value || !end.value || !startValid.value || !endValid.value) return false
  const cmp = compareDateText(start.value, end.value, 'datetime')
  return cmp !== null && cmp > 0
})
const formatDurationValue = (value: number): string => {
  if (Number.isInteger(value)) return String(value)
  return value.toFixed(3).replace(/\.?0+$/, '')
}
const durationPreview = computed(() => {
  if (!start.value || !end.value || !startValid.value || !endValid.value || isRangeInvalid.value) return ''
  const startParsed = parseDateText(start.value, 'datetime')
  const endParsed = parseDateText(end.value, 'datetime')
  if (!startParsed || !endParsed) return ''
  const selectedUnit = intervalUnits.find(opt => opt.value === unit.value)
  if (!selectedUnit) return ''
  const durationRaw = (endParsed.ts - startParsed.ts) / selectedUnit.ms
  return `Интервал: ${formatDurationValue(durationRaw)} ${selectedUnit.short}`
})
watch([start, end, unit], ([nextStart, nextEnd, nextUnit]) => {
  const current = model.value ?? ({} as TraitIntervalModel)
  if (
    String(current.start ?? '') === nextStart &&
    String(current.end ?? '') === nextEnd &&
    String(current.unit ?? 'minutes') === nextUnit
  )
    return
  model.value = {
    ...current,
    start: nextStart,
    end: nextEnd,
    unit: nextUnit
  }
})
watch(
  (): [string, string, TraitIntervalUnit] => [
    String(model.value?.start ?? ''),
    String(model.value?.end ?? ''),
    intervalUnits.some(opt => opt.value === model.value?.unit) ? (model.value.unit as TraitIntervalUnit) : 'minutes'
  ],
  ([nextStart, nextEnd, nextUnit]) => {
    if (start.value !== nextStart) start.value = nextStart
    if (end.value !== nextEnd) end.value = nextEnd
    if (unit.value !== nextUnit) unit.value = nextUnit
  }
)
</script>
