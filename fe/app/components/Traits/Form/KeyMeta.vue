<template>
  <div class="w-full min-w-0 space-y-3">
    <LabField label="Тип данных" for-id="keyType" class="w-full min-w-0">
      <LabBaseSelect id="keyType" v-model="form.meta.dataType" :options="dataTypeOptions" />
    </LabField>
    <div v-if="form.meta.dataType === 'enum'" class="traits-meta-group space-y-2">
      <LabField label="Тип опции" for-id="enumOptionType" class="w-full min-w-0">
        <LabBaseSelect id="enumOptionType" v-model="form.meta.optionType" :options="optionTypeOptions" />
      </LabField>
      <LabField label="Опции" for-id="enumOptions" class="w-full min-w-0">
        <LabBaseInput id="enumOptions" v-model="enumRaw" placeholder="красный, зелёный, синий" class="w-full min-w-0" />
      </LabField>
    </div>
    <div v-else-if="form.meta.dataType === 'color'" class="traits-meta-group">
      <LabField label="Режим цвета" for-id="colorMode" class="w-full min-w-0">
        <LabBaseSelect id="colorMode" v-model="form.meta.mode" :options="colorModeSelectOptions" />
      </LabField>
    </div>
    <div
      v-else-if="form.meta.dataType === 'number'"
      class="traits-meta-group grid grid-cols-1 gap-2 md:grid-cols-2">
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
  const { dataTypes, createMeta } = useTraitTypes()
  const form = reactive<{ meta: KeyMeta }>({
    meta: createMeta('string')
  })
  const enumRaw = ref('')
  let debounceTimer: ReturnType<typeof setTimeout> | undefined
  let requestToken = 0
  const dataTypeOptions = computed(() =>
    dataTypes.value.map(dt => ({
      value: dt,
      label: DATA_TYPE_SELECT_LABELS[dt] || dt
    }))
  )
  const optionTypeOptions = computed(() =>
    dataTypes.value
      .filter(dt => dt !== 'enum')
      .map(dt => ({
        value: dt,
        label: DATA_TYPE_SELECT_LABELS[dt] || dt
      }))
  )
  const colorModeSelectOptions = computed(() =>
    COLOR_MODE_OPTIONS.map(opt => ({
      value: opt.value,
      label: opt.label
    }))
  )
  const unitCategoryOptions = computed(() =>
    UNIT_CATEGORIES.map(cat => ({
      value: cat.id,
      label: cat.label
    }))
  )
  const unitsByCategoryOptions = computed(() =>
    getUnitOptionsByCategory(form.meta.unitCategory).map(unit => ({
      value: unit.value,
      label: `${unit.label} (${unit.value})`
    }))
  )
  const resetMetaState = () => {
    form.meta = createMeta('string')
    enumRaw.value = ''
    metaModel.value = form.meta
  }
  const emitMeta = () => {
    metaModel.value = normalizeKeyMeta(form.meta, enumRaw.value)
  }
  const applyLoadedMeta = (meta?: KeyMeta | null) => {
    const nextMeta = meta || createMeta('string')
    form.meta = { ...createMeta(nextMeta.dataType ?? 'string'), ...nextMeta }
    enumRaw.value = Array.isArray(nextMeta.options) ? nextMeta.options.join(', ') : ''
    emitMeta()
  }
  const loadMeta = async () => {
    const token = ++requestToken
    const syn = String(props.syn || '').trim()
    if (!syn) {
      resetMetaState()
      return
    }
    try {
      const res = await getKeyMeta(syn)
      if (token !== requestToken) return
      const meta = (res?.data?.meta as KeyMeta | undefined) || null
      if (!meta) {
        resetMetaState()
        return
      }
      applyLoadedMeta(meta)
    } catch {
      if (token !== requestToken) return
      resetMetaState()
    }
  }
  watch(
    () => props.syn,
    nextValue => {
      if (debounceTimer) {
        clearTimeout(debounceTimer)
        debounceTimer = undefined
      }
      if (!String(nextValue || '').trim()) {
        requestToken += 1
        resetMetaState()
        return
      }
      debounceTimer = setTimeout(() => {
        loadMeta()
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
    emitMeta
  )
  watch(
    () => form.meta.dataType,
    nextType => {
      if (nextType !== 'enum') {
        enumRaw.value = ''
      }
      if (nextType !== 'color') {
        form.meta.mode = createMeta('color').mode
      }
      if (nextType !== 'number') {
        const numberDefaults = createMeta('number')
        form.meta.unitCategory = numberDefaults.unitCategory
        form.meta.unit = numberDefaults.unit
      }
    }
  )
  watch(
    () => form.meta.unitCategory,
    nextCategory => {
      if (form.meta.dataType !== 'number') return
      const nextUnits = getUnitOptionsByCategory(nextCategory)
      if (!nextUnits.length) {
        form.meta.unit = ''
        return
      }
      if (!nextUnits.some(unit => unit.value === form.meta.unit)) {
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
<style scoped>
  .traits-meta-group {
    border: 1px solid color-mix(in srgb, var(--lab-border) 78%, transparent);
    background: color-mix(in srgb, var(--lab-bg-surface-subtle) 72%, transparent);
    padding: 0.75rem;
  }
</style>
