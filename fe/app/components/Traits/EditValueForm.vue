<template>
  <section
    v-if="trait"
    class="w-full p-3 sm:p-4"
    :class="
      inline
        ? ''
        : 'mx-auto max-w-3xl border-[color-mix(in_srgb,var(--lab-border)_82%,transparent)] bg-[color-mix(in_srgb,var(--lab-bg-surface)_84%,transparent)]'
    "
  >
    <form class="space-y-4" @submit.prevent="onSubmit">
      <component :is="valueComponent" :id="valueInputId" ref="valueRootRef" v-model="valueModel" :meta="resolvedMeta" />
      <div class="flex flex-col gap-2 sm:flex-row sm:justify-end">
        <LabBaseButton
          variant="ghost"
          :label="t('edit_form.cancel')"
          :disabled="pending"
          @click="emit('cancel')"
        />
        <LabBaseButton
          variant="primary"
          type="submit"
          :label="t('edit_form.save')"
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
import TraitsInputInterval from './Form/InputInterval.vue'
import TraitsInputSchedule from './Form/InputSchedule.vue'
import TraitsInputValidity from './Form/InputValidity.vue'
import TraitsInputColor from './Form/InputColor.vue'
import TraitsInputList from './Form/InputList.vue'
import TraitsInputRange from './Form/InputRange.vue'
import TraitsInputGeo from './Form/InputGeo.vue'
import TraitsInputSurface from './Form/InputSurface.vue'
import { defaultKeyMeta, normalizeKeyMeta } from '../../utils/traitMeta'
import { isTraitFormValueFilled, parseTraitStoredValue, serializeTraitFormValue } from '../../utils/traitValueCodec'
const { locale, key, load, t } = useI18nSection('traits')
await useAsyncData(key.value, load, { watch: [locale] })
const props = withDefaults(
  defineProps<{
    trait: Trait | null
    meta?: KeyMeta | null
    pending?: boolean
    inline?: boolean
  }>(),
  {
    trait: null,
    meta: null,
    pending: false,
    inline: false
  }
)
const emit = defineEmits<{
  (e: 'save', payload: { traitUuid: string; t_key: string; t_value: string }): void
  (e: 'cancel'): void
}>()
const value = ref<TraitDynamicValueModel>('')
const valueRootRef = ref<ComponentPublicInstance<{ focus?: () => void }> | HTMLElement | null>(null)
const resolvedMeta = computed<KeyMeta>(() => {
  return normalizeKeyMeta(props.meta || defaultKeyMeta('string'))
})
const valueModel = computed<TraitDynamicValueModel>({
  get: () => value.value,
  set: next => {
    value.value = next
  }
})
const valueInputId = computed(() => `TraitEditValue-${resolvedMeta.value.dataType || 'string'}`)
const valueComponent = computed(() => {
  switch (resolveDataType(resolvedMeta.value.dataType)) {
    case 'number':
      return TraitsInputNumber
    case 'boolean':
      return TraitsInputBoolean
    case 'datetime':
      return TraitsInputDatetime
    case 'range':
      return TraitsInputRange
    case 'interval':
      return TraitsInputInterval
    case 'schedule':
      return TraitsInputSchedule
    case 'geo':
      return TraitsInputGeo
    case 'list':
      return TraitsInputList
    case 'validity':
      return TraitsInputValidity
    case 'color':
      return TraitsInputColor
    case 'surface':
      return TraitsInputSurface
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
  const candidateIds = [valueInputId.value].filter(Boolean)
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
