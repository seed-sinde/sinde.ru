<template>
  <section
    v-if="trait"
    class="mx-auto w-full max-w-3xl border-[color-mix(in_srgb,var(--lab-border)_82%,transparent)] bg-[color-mix(in_srgb,var(--lab-bg-surface)_84%,transparent)] p-3 sm:p-4"
  >
    <form class="space-y-4" @submit.prevent="onSubmit">
      <div class="space-y-1">
        <h4 class="text-sm font-medium sm:text-base">{{ copy.editForm.title }}</h4>
        <p class="lab-text-muted text-xs">{{ copy.editForm.description }}</p>
      </div>
      <div class="grid grid-cols-1 gap-3 lg:grid-cols-2">
        <LabField :label="copy.editForm.keyName" for-id="TraitEditKey" class="max-w-full">
          <LabBaseInput id="TraitEditKey" :model-value="trait.t_key" readonly disabled />
        </LabField>
        <LabField v-if="showValueLabel" :label="copy.editForm.value" :for-id="valueLabelFor" class="max-w-full">
          <component
            :is="valueComponent"
            :id="valueInputId"
            ref="valueRootRef"
            v-model="valueModel"
            :meta="resolvedMeta"
            class="block w-full"
          />
        </LabField>
        <component
          :is="valueComponent"
          v-else
          :id="valueInputId"
          ref="valueRootRef"
          v-model="valueModel"
          :meta="resolvedMeta"
          class="block w-full"
        />
      </div>
      <div class="flex flex-col gap-2 sm:flex-row sm:justify-end">
        <LabBaseButton variant="ghost" :label="copy.editForm.cancel" :disabled="pending" @click="emit('cancel')" />
        <LabBaseButton
          variant="primary"
          type="submit"
          :label="copy.editForm.save"
          :disabled="pending || !isValueFilled"
        />
      </div>
    </form>
  </section>
</template>
<script setup lang="ts">
import type { ComponentPublicInstance } from 'vue'
import TraitsInputString from './Form/InputString.vue'
import TraitsInputNumber from './Form/InputNumber.vue'
import TraitsInputBoolean from './Form/InputBoolean.vue'
import TraitsInputDatetime from './Form/InputDatetime.vue'
import TraitsInputDateRange from './Form/InputDateRange.vue'
import TraitsInputInterval from './Form/InputInterval.vue'
import TraitsInputSchedule from './Form/InputSchedule.vue'
import TraitsInputGeoPoint from './Form/InputGeoPoint.vue'
import TraitsInputEnum from './Form/InputEnum.vue'
import TraitsInputValidity from './Form/InputValidity.vue'
import TraitsInputColor from './Form/InputColor.vue'
import { defaultKeyMeta, normalizeKeyMeta } from '../../utils/traitMeta'
import { isTraitFormValueFilled, parseTraitStoredValue, serializeTraitFormValue } from '../../utils/traitValueCodec'
const { localeCode } = useInterfacePreferences()
const copy = computed(() => TRAITS_WORKSPACE_COPY[localeCode.value] || TRAITS_WORKSPACE_COPY.ru)
const props = withDefaults(
  defineProps<{
    trait: Trait | null
    meta?: KeyMeta | null
    pending?: boolean
  }>(),
  {
    trait: null,
    meta: null,
    pending: false
  }
)
const emit = defineEmits<{
  (e: 'save', payload: { traitUuid: string; t_key: string; t_value: string }): void
  (e: 'cancel'): void
}>()
const value = ref<string | number | boolean | Record<string, any>>('')
const valueRootRef = ref<ComponentPublicInstance<{ focus?: () => void }> | HTMLElement | null>(null)
const resolvedMeta = computed<KeyMeta>(() => {
  return normalizeKeyMeta(props.meta || defaultKeyMeta('string'))
})
const valueModel = computed<any>({
  get: () => value.value,
  set: (next) => {
    value.value = next as any
  }
})
const valueComponentHasLabel = computed(() => {
  return ['enum', 'boolean', 'datetime-range', 'interval', 'schedule', 'geo-point', 'validity', 'color'].includes(
    resolvedMeta.value.dataType as any
  )
})
const showValueLabel = computed(() => !valueComponentHasLabel.value)
const valueInputId = computed(() => `TraitEditValue-${resolvedMeta.value.dataType || 'string'}`)
const valueLabelFor = computed(() => {
  if (resolvedMeta.value.dataType === 'datetime-range') return `${valueInputId.value}-start`
  if (resolvedMeta.value.dataType === 'interval') return `${valueInputId.value}-start`
  if (resolvedMeta.value.dataType === 'schedule') return `${valueInputId.value}-from-day`
  if (resolvedMeta.value.dataType === 'boolean') return 'TraitInputBoolean'
  if (resolvedMeta.value.dataType === 'color') return 'color-hex'
  return valueInputId.value
})
const valueComponent = computed(() => {
  switch (resolvedMeta.value.dataType) {
    case 'number':
      return TraitsInputNumber
    case 'boolean':
      return TraitsInputBoolean
    case 'datetime':
      return TraitsInputDatetime
    case 'datetime-range':
      return TraitsInputDateRange
    case 'interval':
      return TraitsInputInterval
    case 'schedule':
      return TraitsInputSchedule
    case 'geo-point':
      return TraitsInputGeoPoint
    case 'enum':
      return TraitsInputEnum
    case 'validity':
      return TraitsInputValidity
    case 'color':
      return TraitsInputColor
    default:
      return TraitsInputString
  }
})
const isValueFilled = computed(() => {
  return isTraitFormValueFilled(resolvedMeta.value, value.value)
})
const hydrateValueFromTrait = () => {
  if (!props.trait) {
    value.value = ''
    return
  }
  value.value = parseTraitStoredValue(resolvedMeta.value, String(props.trait.t_value || ''))
}
watch(
  () => [
    String(props.trait?.t_uuid || ''),
    String(props.trait?.t_value || ''),
    String(resolvedMeta.value.dataType || ''),
    String(resolvedMeta.value.mode || ''),
    String(resolvedMeta.value.unit || '')
  ],
  hydrateValueFromTrait,
  { immediate: true }
)
const focusValueInput = async () => {
  if (!import.meta.client || !props.trait) return
  await nextTick()
  const componentRef = valueRootRef.value as ComponentPublicInstance<{ focus?: () => void }> | null
  if (componentRef && typeof componentRef.focus === 'function') {
    componentRef.focus()
    return
  }
  const candidateIds = [valueLabelFor.value, valueInputId.value].filter(Boolean)
  for (const id of candidateIds) {
    const root = document.getElementById(id)
    if (!root) continue
    if (typeof (root as HTMLElement).focus === 'function') {
      ;(root as HTMLElement).focus()
      return
    }
    const nested = root.querySelector<HTMLElement>('input, textarea, select, button, [tabindex]:not([tabindex="-1"])')
    if (nested && typeof nested.focus === 'function') {
      nested.focus()
      return
    }
  }
  const fallbackRoot = valueRootRef.value instanceof HTMLElement ? valueRootRef.value : null
  const fallbackTarget = fallbackRoot?.querySelector<HTMLElement>(
    'input, textarea, select, button, [tabindex]:not([tabindex="-1"])'
  )
  fallbackTarget?.focus()
}
watch(
  () => [String(props.trait?.t_uuid || ''), String(resolvedMeta.value.dataType || '')],
  () => {
    void focusValueInput()
  },
  { immediate: true }
)
const onSubmit = () => {
  if (!props.trait || !isValueFilled.value) return
  const serialized = serializeTraitFormValue(resolvedMeta.value, value.value).trim()
  if (!serialized) return
  emit('save', {
    traitUuid: props.trait.t_uuid,
    t_key: props.trait.t_key,
    t_value: serialized
  })
}
</script>
