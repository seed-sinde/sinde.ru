<template>
  <form
    ref="formRef"
    class="mx-auto flex w-full max-w-2xl flex-col gap-4 border-[color-mix(in_srgb,var(--lab-border)_82%,transparent)] bg-[color-mix(in_srgb,var(--lab-bg-surface)_86%,transparent)] p-4 sm:p-5"
    @submit.prevent="onSubmit">
    <div class="space-y-3">
      <LabField :label="copy.addForm.keyName" for-id="trait-key" class="w-full min-w-0">
        <LabBaseInput
          id="trait-key"
          ref="keyInput"
          v-model="key"
          name="TraitKey"
          :placeholder="copy.addForm.keyPlaceholder"
          autocomplete="off"
          spellcheck="false"
          inputmode="text"
          class="w-full min-w-0" />
      </LabField>
      <div class="w-full min-w-0">
        <LabField v-if="showValueLabel" :label="copy.addForm.value" :for-id="valueLabelFor" class="w-full min-w-0">
          <component
            :is="valueComponent"
            :id="valueInputId"
            v-model="valueModel"
            :meta="metaForInput"
            class="block w-full min-w-0" />
        </LabField>
        <component
          v-else
          :is="valueComponent"
          :id="valueInputId"
          v-model="valueModel"
          :meta="metaForInput"
          class="block w-full min-w-0" />
      </div>
    </div>
    <div class="space-y-3">
      <div class="min-w-0">
        <TraitsFormKeyMeta v-model="meta" :syn="key" />
      </div>
      <div class="flex sm:justify-end">
        <LabBaseButton
          type="submit"
          :label="copy.addForm.add"
          :disabled="!key || !isValueFilled"
          class="w-full sm:w-auto" />
      </div>
    </div>
  </form>
</template>
<script setup lang="ts">
  const { localeCode } = useInterfacePreferences()
  import TraitsInputString from './InputString.vue'
  import TraitsInputNumber from './InputNumber.vue'
  import TraitsInputBoolean from './InputBoolean.vue'
  import TraitsInputDatetime from './InputDatetime.vue'
  import TraitsInputDateRange from './InputDateRange.vue'
  import TraitsInputInterval from './InputInterval.vue'
  import TraitsInputSchedule from './InputSchedule.vue'
  import TraitsInputGeoPoint from './InputGeoPoint.vue'
  import TraitsInputEnum from './InputEnum.vue'
  import TraitsInputValidity from './InputValidity.vue'
  import TraitsInputColor from './InputColor.vue'
  const copy = computed(() => TRAITS_WORKSPACE_COPY[localeCode.value] || TRAITS_WORKSPACE_COPY.ru)
  const emit = defineEmits<{ (e: 'add', trait: TraitInput): void }>()
  const key = ref('')
  const meta = ref<KeyMeta>(defaultKeyMeta('string'))
  const value = ref<string | number | boolean | Record<string, any>>('')
  const enumOptions = ref<string[]>([])
  const keyInput = ref<{ focus: () => void } | null>(null)
  const formRef = ref<HTMLFormElement | null>(null)
  const valueModel = computed<any>({
    get: () => value.value,
    set: next => {
      value.value = next as any
    }
  })
  const valueComponent = computed(() => {
    switch (meta.value.dataType) {
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
  const valueComponentHasLabel = computed(() => {
    return ['enum', 'boolean', 'datetime-range', 'interval', 'schedule', 'geo-point', 'validity', 'color'].includes(
      meta.value.dataType as any
    )
  })
  const showValueLabel = computed(() => !valueComponentHasLabel.value)
  const valueLabelFor = computed(() => {
    if (meta.value.dataType === 'datetime-range') return `${valueInputId.value}-start`
    if (meta.value.dataType === 'interval') return `${valueInputId.value}-start`
    if (meta.value.dataType === 'schedule') return `${valueInputId.value}-from-day`
    return valueInputId.value
  })
  const metaForInput = computed(() => {
    if (meta.value.dataType !== 'enum') return meta.value
    const options = Array.from(new Set([...((meta.value as any)?.options ?? []), ...enumOptions.value]))
    return {
      ...(meta.value as any),
      options
    }
  })
  const valueInputId = computed(() => `TraitValue-${meta.value.dataType || 'string'}`)
  const isValueFilled = computed(() => {
    return isTraitFormValueFilled(meta.value, value.value)
  })
  const loadEnumOptions = async () => {
    if (meta.value.dataType !== 'enum' || !key.value.trim()) {
      enumOptions.value = []
      return
    }
    try {
      const res = await getEnumOptions(key.value.trim())
      enumOptions.value = res?.data?.options ?? []
    } catch (error) {
      console.warn('enum options load failed', error)
      enumOptions.value = []
    }
  }
  const captureFocusSnapshot = (): TraitAddFormFocusSnapshot | null => {
    if (!import.meta.client || !formRef.value) return null
    const active = document.activeElement as HTMLElement | null
    if (!active || !formRef.value.contains(active)) return null
    const inputLike = active as HTMLInputElement | HTMLTextAreaElement
    const id = String(active.id || '')
    return {
      id,
      name: String(active.getAttribute('name') || ''),
      wasValueField: id.startsWith('TraitValue-'),
      selectionStart: typeof inputLike.selectionStart === 'number' ? inputLike.selectionStart : null,
      selectionEnd: typeof inputLike.selectionEnd === 'number' ? inputLike.selectionEnd : null
    }
  }
  const restoreFocusSnapshot = (snapshot: TraitAddFormFocusSnapshot | null) => {
    if (!import.meta.client || !snapshot || !formRef.value) return
    let target: HTMLElement | null = null
    if (snapshot.id) {
      const byId = document.getElementById(snapshot.id)
      if (byId && formRef.value.contains(byId)) target = byId as HTMLElement
    }
    if (!target && snapshot.wasValueField) {
      const fallbackId = valueLabelFor.value
      const byValueId = document.getElementById(fallbackId)
      if (byValueId && formRef.value.contains(byValueId)) target = byValueId as HTMLElement
    }
    if (!target && snapshot.name) {
      const byName = formRef.value.querySelector(`[name="${snapshot.name}"]`)
      if (byName) target = byName as HTMLElement
    }
    if (!target) return
    target.focus()
    if (
      (target instanceof HTMLInputElement || target instanceof HTMLTextAreaElement) &&
      snapshot.selectionStart !== null &&
      snapshot.selectionEnd !== null
    ) {
      try {
        const length = target.value.length
        const start = Math.min(snapshot.selectionStart, length)
        const end = Math.min(snapshot.selectionEnd, length)
        target.setSelectionRange(start, end)
      } catch {
        // no-op: some input types do not support setSelectionRange
      }
    }
  }
  watch([() => meta.value.dataType, key], async ([nextDataType], [prevDataType]) => {
    const focusSnapshot = nextDataType !== prevDataType ? captureFocusSnapshot() : null
    loadEnumOptions()
    value.value = defaultTraitFormValue(meta.value)
    if (focusSnapshot) {
      await nextTick()
      restoreFocusSnapshot(focusSnapshot)
    }
  })
  const onSubmit = async () => {
    if (!key.value || !isValueFilled.value) return
    emit('add', {
      t_key: key.value,
      t_value: serializeTraitFormValue(meta.value, value.value),
      meta: meta.value
    })
    value.value = defaultTraitFormValue(meta.value)
    await nextTick()
    keyInput.value?.focus()
  }
  onMounted(() => {
    keyInput.value?.focus()
  })
</script>
