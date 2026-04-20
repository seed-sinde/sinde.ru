<template>
  <div class="space-y-1">
    <div class="grid grid-cols-1 gap-2 sm:grid-cols-2">
      <LabBaseField label="С дня" :for-id="fromDayId">
        <LabBaseSelect :id="fromDayId" v-model="fromDay" :aria-label="fromDayLabel" :options="weekDayOptions" />
      </LabBaseField>
      <LabBaseField label="По день" :for-id="toDayId">
        <LabBaseSelect :id="toDayId" v-model="toDay" :aria-label="toDayLabel" :options="weekDayOptions" />
      </LabBaseField>
      <LabBaseField label="С времени" :for-id="fromTimeId">
        <TraitsFormDateTextInput
          :id="fromTimeId"
          v-model="fromTime"
          mode="time"
          :placeholder="timePlaceholder"
          :aria-label="fromTimeLabel"
          :invalid="Boolean(fromTime) && !fromTimeValid"
        />
      </LabBaseField>
      <LabBaseField label="По время" :for-id="toTimeId">
        <TraitsFormDateTextInput
          :id="toTimeId"
          v-model="toTime"
          mode="time"
          :placeholder="timePlaceholder"
          :aria-label="toTimeLabel"
          :invalid="Boolean(toTime) && !toTimeValid"
        />
      </LabBaseField>
    </div>
    <LabErrorMessage
      v-if="(fromTime && !fromTimeValid) || (toTime && !toTimeValid)"
      :text="`Неверный формат: ${timePlaceholder}`"
    />
    <LabErrorMessage
      v-else-if="isTimeRangeInvalid"
      text='В пределах одного дня время "С" не может быть позже времени "По".'
    />
  </div>
</template>
<script setup lang="ts">
import TraitsFormDateTextInput from './DateTextInput.vue'
const model = defineModel<TraitScheduleModel>({ required: true })
const props = withDefaults(
  defineProps<{
    id?: string
    labelFromDay?: string
    labelToDay?: string
    labelFromTime?: string
    labelToTime?: string
    meta?: unknown
  }>(),
  {}
)
const weekDays = [
  { value: '1', label: 'Понедельник' },
  { value: '2', label: 'Вторник' },
  { value: '3', label: 'Среда' },
  { value: '4', label: 'Четверг' },
  { value: '5', label: 'Пятница' },
  { value: '6', label: 'Суббота' },
  { value: '7', label: 'Воскресенье' }
] as const
const weekDayOptions = weekDays.map(day => ({ value: day.value, label: day.label }))
const fromDayId = computed(() => (props.id ? `${props.id}-from-day` : 'TraitScheduleFromDay'))
const toDayId = computed(() => (props.id ? `${props.id}-to-day` : 'TraitScheduleToDay'))
const fromTimeId = computed(() => (props.id ? `${props.id}-from-time` : 'TraitScheduleFromTime'))
const toTimeId = computed(() => (props.id ? `${props.id}-to-time` : 'TraitScheduleToTime'))
const fromDayLabel = computed(() => props.labelFromDay || 'С дня')
const toDayLabel = computed(() => props.labelToDay || 'По день')
const fromTimeLabel = computed(() => props.labelFromTime || 'С времени')
const toTimeLabel = computed(() => props.labelToTime || 'По время')
const timePlaceholder = computed(() => datePlaceholder('time'))
const fromDay = ref(String(model.value?.fromDay ?? '1'))
const toDay = ref(String(model.value?.toDay ?? '5'))
const fromTime = ref(String(model.value?.fromTime ?? ''))
const toTime = ref(String(model.value?.toTime ?? ''))
const fromTimeValid = computed(() => !fromTime.value || isValidDateText(fromTime.value, 'time'))
const toTimeValid = computed(() => !toTime.value || isValidDateText(toTime.value, 'time'))
const isTimeRangeInvalid = computed(() => {
  if (!fromTime.value || !toTime.value || !fromTimeValid.value || !toTimeValid.value) return false
  if (fromDay.value !== toDay.value) return false
  const cmp = compareDateText(fromTime.value, toTime.value, 'time')
  return cmp !== null && cmp > 0
})
watch([fromDay, toDay, fromTime, toTime], ([nextFromDay, nextToDay, nextFromTime, nextToTime]) => {
  const current = model.value ?? ({} as TraitScheduleModel)
  if (
    String(current.fromDay ?? '') === nextFromDay &&
    String(current.toDay ?? '') === nextToDay &&
    String(current.fromTime ?? '') === nextFromTime &&
    String(current.toTime ?? '') === nextToTime
  )
    return
  model.value = {
    ...current,
    fromDay: nextFromDay,
    toDay: nextToDay,
    fromTime: nextFromTime,
    toTime: nextToTime
  }
})
watch(
  (): [string, string, string, string] => [
    String(model.value?.fromDay ?? '1'),
    String(model.value?.toDay ?? '5'),
    String(model.value?.fromTime ?? ''),
    String(model.value?.toTime ?? '')
  ],
  ([nextFromDay, nextToDay, nextFromTime, nextToTime]) => {
    if (fromDay.value !== nextFromDay) fromDay.value = nextFromDay
    if (toDay.value !== nextToDay) toDay.value = nextToDay
    if (fromTime.value !== nextFromTime) fromTime.value = nextFromTime
    if (toTime.value !== nextToTime) toTime.value = nextToTime
  }
)
</script>
