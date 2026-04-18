<template>
  <div class="grid gap-3">
    <LabField label="Тип данных" for-id="trait-meta-type">
      <LabBaseSelect id="trait-meta-type" v-model="form.dataType" :options="dataTypeOptions" />
    </LabField>

    <div v-if="form.dataType === 'string'" class="grid gap-2 md:grid-cols-2">
      <LabField label="Минимальная длина" for-id="trait-meta-min-length">
        <LabBaseInput id="trait-meta-min-length" v-model="minLengthText" type="number" inputmode="numeric" />
      </LabField>
      <LabField label="Максимальная длина" for-id="trait-meta-max-length">
        <LabBaseInput id="trait-meta-max-length" v-model="maxLengthText" type="number" inputmode="numeric" />
      </LabField>
    </div>

    <div v-else-if="form.dataType === 'boolean'" class="grid gap-2 md:grid-cols-2">
      <LabField label="Способ отображения" for-id="trait-meta-boolean-display">
        <LabBaseSelect id="trait-meta-boolean-display" v-model="form.booleanDisplay" :options="booleanDisplayOptions" />
      </LabField>
      <div class="grid gap-2">
        <LabField label="Подпись true" for-id="trait-meta-boolean-true">
          <LabBaseInput id="trait-meta-boolean-true" v-model="form.booleanTrueLabel" />
        </LabField>
        <LabField label="Подпись false" for-id="trait-meta-boolean-false">
          <LabBaseInput id="trait-meta-boolean-false" v-model="form.booleanFalseLabel" />
        </LabField>
      </div>
    </div>

    <div v-else-if="form.dataType === 'number'" class="grid gap-2 md:grid-cols-2">
      <LabField label="Категория" for-id="trait-meta-unit-category">
        <LabBaseSelect id="trait-meta-unit-category" v-model="form.unitCategory" :options="unitCategoryOptions" />
      </LabField>
      <LabField label="Единица" for-id="trait-meta-unit">
        <LabBaseSelect id="trait-meta-unit" v-model="form.unit" :options="unitOptions" />
      </LabField>
    </div>

    <div v-else-if="form.dataType === 'list'" class="grid gap-2">
      <LabField label="Тип значений" for-id="trait-meta-list-value">
        <LabBaseSelect id="trait-meta-list-value" v-model="form.listValueType" :options="listValueTypeOptions" />
      </LabField>
      <div class="grid gap-2 md:grid-cols-2">
        <LabBaseSwitch v-model="listOrderedModel" label="Упорядоченный" />
        <LabBaseSwitch v-model="listUniqueModel" label="Уникальные значения" />
      </div>
      <LabField label="Минимум значений" for-id="trait-meta-list-min-items">
        <LabBaseInput id="trait-meta-list-min-items" v-model="minItemsText" type="number" inputmode="numeric" />
      </LabField>
      <LabField label="Предустановленные варианты" for-id="trait-meta-options" hint="Необязательно. По одному значению на строку.">
        <LabBaseTextarea id="trait-meta-options" v-model="optionsRaw" rows="3" />
      </LabField>
    </div>

    <div v-else-if="form.dataType === 'range'" class="grid gap-2 md:grid-cols-2">
      <LabField label="Тип диапазона" for-id="trait-meta-range-type">
        <LabBaseSelect id="trait-meta-range-type" v-model="form.rangeType" :options="rangeTypeOptions" />
      </LabField>
      <LabField label="Часовой пояс" for-id="trait-meta-range-timezone">
        <LabBaseInput id="trait-meta-range-timezone" v-model="form.timezone" placeholder="UTC" />
      </LabField>
    </div>

    <div v-else-if="form.dataType === 'datetime' || form.dataType === 'interval' || form.dataType === 'schedule' || form.dataType === 'validity'" class="grid gap-2 md:grid-cols-2">
      <LabField label="Часовой пояс" for-id="trait-meta-timezone">
        <LabBaseInput id="trait-meta-timezone" v-model="form.timezone" placeholder="UTC" />
      </LabField>
      <LabField v-if="form.dataType === 'interval'" label="Единица длительности" for-id="trait-meta-duration-unit">
        <LabBaseInput id="trait-meta-duration-unit" v-model="form.durationUnit" placeholder="minutes" />
      </LabField>
      <LabField v-if="form.dataType === 'validity'" label="Условное событие" for-id="trait-meta-event">
        <LabBaseInput id="trait-meta-event" v-model="form.eventLabel" placeholder="до события / после события" />
      </LabField>
    </div>

    <div v-else-if="form.dataType === 'geo'" class="grid gap-2 md:grid-cols-2">
      <LabField label="Подтип" for-id="trait-meta-geo-type">
        <LabBaseSelect id="trait-meta-geo-type" v-model="form.geoType" :options="geoTypeOptions" />
      </LabField>
      <LabField v-if="form.geoType === 'point'" label="Единица высоты" for-id="trait-meta-height-unit">
        <LabBaseInput id="trait-meta-height-unit" v-model="form.heightUnit" placeholder="m" />
      </LabField>
      <LabField v-if="form.geoType === 'zone'" label="Единица радиуса" for-id="trait-meta-radius-unit">
        <LabBaseInput id="trait-meta-radius-unit" v-model="form.radiusUnit" placeholder="m" />
      </LabField>
    </div>

    <div v-else-if="form.dataType === 'color'" class="grid gap-2 md:grid-cols-2">
      <LabField label="Режим" for-id="trait-meta-color-mode">
        <LabBaseSelect id="trait-meta-color-mode" v-model="form.mode" :options="colorModeOptions" />
      </LabField>
      <LabField label="Область применения" for-id="trait-meta-color-domain">
        <LabBaseInput id="trait-meta-color-domain" v-model="form.colorDomain" placeholder="дизайн, медицина" />
      </LabField>
      <LabField label="Стандарт" for-id="trait-meta-color-standard">
        <LabBaseInput id="trait-meta-color-standard" v-model="form.colorStandard" placeholder="RAL, Fitzpatrick, VITA" />
      </LabField>
      <LabField label="Условия измерения" for-id="trait-meta-color-conditions">
        <LabBaseInput id="trait-meta-color-conditions" v-model="form.measurementConditions" placeholder="свет, угол, наблюдение" />
      </LabField>
      <LabField label="Фактура поверхности" for-id="trait-meta-color-texture">
        <LabBaseInput id="trait-meta-color-texture" v-model="form.surfaceTexture" />
      </LabField>
      <LabField label="Прозрачность 0..1" for-id="trait-meta-color-opacity">
        <LabBaseInput id="trait-meta-color-opacity" v-model="opacityText" type="number" inputmode="decimal" step="0.01" />
      </LabField>
      <LabField label="Именованный цвет" for-id="trait-meta-color-name">
        <LabBaseInput id="trait-meta-color-name" v-model="form.namedColor" />
      </LabField>
      <LabField label="Палитра" for-id="trait-meta-color-palette">
        <LabBaseInput id="trait-meta-color-palette" v-model="form.palette" placeholder="OKLCH, HEX, CIE Lab" />
      </LabField>
    </div>

    <div v-else-if="form.dataType === 'surface'" class="grid gap-2 md:grid-cols-2">
      <LabField label="Категория блеска" for-id="trait-meta-gloss">
        <LabBaseInput id="trait-meta-gloss" v-model="form.glossCategory" placeholder="матовый, глянец" />
      </LabField>
      <LabField label="Тип рельефа" for-id="trait-meta-relief">
        <LabBaseInput id="trait-meta-relief" v-model="form.reliefType" placeholder="песок, кожа, сатин" />
      </LabField>
    </div>

    <LabField label="Проверка / валидация" for-id="trait-meta-validate">
      <LabBaseInput id="trait-meta-validate" v-model="form.validate" placeholder=">=2, 0..1, 3..n" />
    </LabField>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{syn: string}>()
void props
const metaModel = defineModel<KeyMeta>({required: true})
const form = reactive<KeyMeta>(normalizeKeyMeta(metaModel.value))
const syncing = ref(false)
const syncForm = (next: KeyMeta) => {
  for (const key of Object.keys(form) as Array<keyof KeyMeta>) {
    if (!(key in next)) delete form[key]
  }
  Object.assign(form, next)
}
const dataTypeOptions = computed(() => DATA_TYPES.map(value => ({value, label: DATA_TYPE_SELECT_LABELS[value] || value})))
const booleanDisplayOptions = computed(() =>
  BOOLEAN_DISPLAY_OPTIONS.map(option => ({value: option.value, label: option.label}))
)
const listValueTypeOptions = computed(() =>
  LIST_VALUE_TYPE_OPTIONS.map(option => ({value: option.value, label: option.label}))
)
const rangeTypeOptions = computed(() =>
  RANGE_VALUE_TYPE_OPTIONS.map(option => ({value: option.value, label: option.label}))
)
const geoTypeOptions = computed(() => GEO_TYPE_OPTIONS.map(option => ({value: option.value, label: option.label})))
const colorModeOptions = computed(() => COLOR_MODE_OPTIONS.map(option => ({value: option.value, label: option.label})))
const unitCategoryOptions = computed(() => UNIT_CATEGORIES.map(option => ({value: option.id, label: option.label})))
const unitOptions = computed(() =>
  getUnitOptionsByCategory(form.unitCategory).map(option => ({value: option.value, label: `${option.label} (${option.value})`}))
)
const clearNumericField = (key: "minLength" | "maxLength" | "minItems" | "opacity") => {
  delete form[key]
}
const bindNumberText = (key: keyof Pick<KeyMeta, "minLength" | "maxLength" | "minItems" | "opacity">) =>
  computed({
    get: () => (typeof form[key] === "number" ? String(form[key]) : ""),
    set: next => {
      const value = String(next || "").trim()
      if (!value) {
        clearNumericField(key)
        return
      }
      const parsed = Number(value)
      if (!Number.isNaN(parsed)) form[key] = parsed as never
    }
  })
const minLengthText = bindNumberText("minLength")
const maxLengthText = bindNumberText("maxLength")
const minItemsText = bindNumberText("minItems")
const opacityText = bindNumberText("opacity")
const listOrderedModel = computed({
  get: () => Boolean(form.listOrdered),
  set: next => {
    form.listOrdered = next
  }
})
const listUniqueModel = computed({
  get: () => Boolean(form.listUnique),
  set: next => {
    form.listUnique = next
  }
})
const optionsRaw = computed({
  get: () => (Array.isArray(form.options) ? form.options.join("\n") : ""),
  set: next => {
    form.options = String(next || "")
      .split("\n")
      .map(item => item.trim())
      .filter(Boolean)
  }
})
watch(
  () => metaModel.value,
  next => {
    if (syncing.value) return
    syncForm(normalizeKeyMeta(next))
  },
  {deep: true, immediate: true}
)
watch(
  () => form.dataType,
  next => {
    syncForm(normalizeKeyMeta({...form, dataType: next}))
  }
)
watch(
  form,
  next => {
    syncing.value = true
    metaModel.value = normalizeKeyMeta({...next})
    nextTick(() => {
      syncing.value = false
    })
  },
  {deep: true}
)
watch(
  () => form.booleanDisplay,
  next => {
    const labels = getDefaultBooleanLabels(next)
    if (!String(form.booleanTrueLabel || "").trim()) form.booleanTrueLabel = labels.trueLabel
    if (!String(form.booleanFalseLabel || "").trim()) form.booleanFalseLabel = labels.falseLabel
  },
  {immediate: true}
)
watch(
  () => form.unitCategory,
  next => {
    if (form.dataType !== "number") return
    const available = getUnitOptionsByCategory(next)
    if (available.some(option => option.value === form.unit)) return
    form.unit = available[0]?.value || ""
  }
)
</script>
