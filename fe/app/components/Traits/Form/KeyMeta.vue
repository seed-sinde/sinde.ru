<template>
  <div class="w-full min-w-0 space-y-3">
    <LabField label="Тип данных" for-id="keyType" class="w-full min-w-0">
      <LabBaseSelect id="keyType" v-model="form.meta.dataType" :options="dataTypeOptions" />
    </LabField>
    <div
      v-if="form.meta.dataType === 'enum'"
      class="space-y-2 border-[color-mix(in_srgb,var(--lab-border)_78%,transparent)] bg-[color-mix(in_srgb,var(--lab-bg-surface-subtle)_72%,transparent)] p-3"
    >
      <LabField label="Тип опции" for-id="enumOptionType" class="w-full min-w-0">
        <LabBaseSelect id="enumOptionType" v-model="form.meta.optionType" :options="optionTypeOptions" />
      </LabField>
      <LabField label="Опции" for-id="enumOptions" class="w-full min-w-0">
        <LabBaseInput id="enumOptions" v-model="enumRaw" placeholder="красный, зелёный, синий" class="w-full min-w-0" />
      </LabField>
    </div>
    <div
      v-else-if="form.meta.dataType === 'color'"
      class="border-[color-mix(in_srgb,var(--lab-border)_78%,transparent)] bg-[color-mix(in_srgb,var(--lab-bg-surface-subtle)_72%,transparent)] p-3"
    >
      <LabField label="Режим цвета" for-id="colorMode" class="w-full min-w-0">
        <LabBaseSelect id="colorMode" v-model="form.meta.mode" :options="colorModeSelectOptions" />
      </LabField>
    </div>
    <div
      v-else-if="form.meta.dataType === 'number'"
      class="grid grid-cols-1 gap-2 border-[color-mix(in_srgb,var(--lab-border)_78%,transparent)] bg-[color-mix(in_srgb,var(--lab-bg-surface-subtle)_72%,transparent)] p-3 md:grid-cols-2"
    >
      <LabField label="Категория" for-id="numberUnitCategory" class="w-full min-w-0">
        <LabBaseSelect id="numberUnitCategory" v-model="form.meta.unitCategory" :options="unitCategoryOptions" />
      </LabField>
      <LabField label="Единица" for-id="numberUnit" class="w-full min-w-0">
        <LabBaseSelect id="numberUnit" v-model="form.meta.unit" :options="unitsByCategoryOptions" />
      </LabField>
    </div>
  </div>
</template>
<script setup lang="ts">
const props = defineProps<{
  syn: string
}>()
const metaModel = defineModel<KeyMeta>({ required: true })
const store = useTraitsStore()
const { dataTypes, createMeta } = useTraitTypes()
const form = reactive<{ meta: KeyMeta }>({
  meta: createMeta('string')
})
const enumRaw = ref('')
const metaDirty = ref(false)
let debounceTimer: ReturnType<typeof setTimeout> | undefined
let requestToken = 0
let suppressDirtyTracking = false
const dataTypeOptions = computed(() =>
  dataTypes.value.map((dt) => ({
    value: dt,
    label: DATA_TYPE_SELECT_LABELS[dt] || dt
  }))
)
const optionTypeOptions = computed(() =>
  dataTypes.value
    .filter((dt) => dt !== 'enum')
    .map((dt) => ({
      value: dt,
      label: DATA_TYPE_SELECT_LABELS[dt] || dt
    }))
)
const colorModeSelectOptions = computed(() =>
  COLOR_MODE_OPTIONS.map((opt) => ({
    value: opt.value,
    label: opt.label
  }))
)
const unitCategoryOptions = computed(() =>
  UNIT_CATEGORIES.map((cat) => ({
    value: cat.id,
    label: cat.label
  }))
)
const unitsByCategoryOptions = computed(() =>
  getUnitOptionsByCategory(form.meta.unitCategory).map((unit) => ({
    value: unit.value,
    label: `${unit.label} (${unit.value})`
  }))
)
const findCachedMeta = (syn: string): KeyMeta | null => {
  const targetSyn = String(syn || '')
    .trim()
    .toLowerCase()
  if (!targetSyn) return null
  const match = Object.values(store.keyMetaById || {}).find(
    (item) =>
      String(item?.syn || '')
        .trim()
        .toLowerCase() === targetSyn
  )
  return match?.meta ? (match.meta as KeyMeta) : null
}
const withProgrammaticMetaState = (apply: () => void) => {
  suppressDirtyTracking = true
  apply()
  metaDirty.value = false
  void nextTick(() => {
    suppressDirtyTracking = false
  })
}
const resetMetaState = () => {
  withProgrammaticMetaState(() => {
    form.meta = createMeta('string')
    enumRaw.value = ''
    metaModel.value = form.meta
  })
}
const emitMeta = () => {
  metaModel.value = normalizeKeyMeta(form.meta, enumRaw.value)
}
const applyLoadedMeta = (meta?: KeyMeta | null) => {
  const nextMeta = meta || createMeta('string')
  withProgrammaticMetaState(() => {
    form.meta = { ...createMeta(nextMeta.dataType ?? 'string'), ...nextMeta }
    enumRaw.value = Array.isArray(nextMeta.options) ? nextMeta.options.join(', ') : ''
    emitMeta()
  })
}
const loadMeta = async () => {
  const token = ++requestToken
  const syn = String(props.syn || '').trim()
  if (!syn) {
    if (!metaDirty.value) resetMetaState()
    return
  }
  try {
    const res = await getKeyMeta(syn)
    if (token !== requestToken) return
    if ((res?.data as { notFound?: boolean } | undefined)?.notFound) {
      if (!metaDirty.value) resetMetaState()
      return
    }
    const meta = (res?.data?.meta as KeyMeta | undefined) || null
    if (!meta) {
      if (!metaDirty.value) resetMetaState()
      return
    }
    if (metaDirty.value) return
    applyLoadedMeta(meta)
  } catch {
    if (token !== requestToken) return
    if (!metaDirty.value) resetMetaState()
  }
}
watch(
  () => props.syn,
  (nextValue) => {
    requestToken += 1
    if (debounceTimer) {
      clearTimeout(debounceTimer)
      debounceTimer = undefined
    }
    if (!String(nextValue || '').trim()) {
      resetMetaState()
      return
    }
    const cachedMeta = findCachedMeta(nextValue)
    if (cachedMeta) {
      applyLoadedMeta(cachedMeta)
      return
    }
    resetMetaState()
    debounceTimer = setTimeout(() => {
      void loadMeta()
    }, 450)
  },
  { immediate: true }
)
watch(
  () => [
    form.meta.dataType,
    form.meta.optionType,
    form.meta.mode,
    form.meta.unitCategory,
    form.meta.unit,
    enumRaw.value
  ],
  () => {
    emitMeta()
    if (!suppressDirtyTracking) {
      metaDirty.value = true
    }
  }
)
watch(
  () => form.meta.dataType,
  (nextType) => {
    if (nextType === 'enum') {
      form.meta.optionType = form.meta.optionType || 'string'
    } else {
      enumRaw.value = ''
      delete form.meta.optionType
    }
    if (nextType === 'color') {
      form.meta.mode = resolveColorMode(form.meta) || 'hex'
    } else {
      delete form.meta.mode
    }
    if (nextType === 'number') {
      const numberDefaults = createMeta('number')
      form.meta.unitCategory = form.meta.unitCategory || numberDefaults.unitCategory || 'unitless'
      form.meta.unit = form.meta.unit || numberDefaults.unit || ''
    } else {
      delete form.meta.unitCategory
      delete form.meta.unit
    }
  }
)
watch(
  () => form.meta.unitCategory,
  (nextCategory) => {
    if (form.meta.dataType !== 'number') return
    const nextUnits = getUnitOptionsByCategory(nextCategory)
    if (!nextUnits.length) {
      form.meta.unit = ''
      return
    }
    if (!nextUnits.some((unit) => unit.value === form.meta.unit)) {
      form.meta.unit = nextUnits[0]?.value || ''
    }
  }
)
onBeforeUnmount(() => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
    debounceTimer = undefined
  }
  requestToken += 1
})
</script>
