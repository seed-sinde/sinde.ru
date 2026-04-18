<template>
  <form ref="formRef" class="flex min-w-0 flex-col gap-3" @submit.prevent="onSubmit">
    <TraitsFormKeyMeta v-model="meta" :syn="key" />
    <LabField :label="t('add_form.key.name')" for-id="trait-key">
      <LabBaseInput
        id="trait-key"
        ref="keyInput"
        v-model="key"
        name="TraitKey"
        :placeholder="t('add_form.key.placeholder')"
        autocomplete="off"
        spellcheck="false"
        inputmode="text"
      />
    </LabField>
    <LabField v-if="showValueLabel" :label="t('add_form.value')" :for-id="valueLabelFor">
      <component :is="valueComponent" :id="valueInputId" v-model="valueModel" :meta="meta" />
    </LabField>
    <component :is="valueComponent" v-else :id="valueInputId" v-model="valueModel" :meta="meta" />
    <LabBaseButton type="submit" variant="primary" :label="t('add_form.add')" :disabled="!key.trim() || !isValueFilled" />
  </form>
</template>

<script setup lang="ts">
import TraitsInputString from "./InputString.vue"
import TraitsInputNumber from "./InputNumber.vue"
import TraitsInputBoolean from "./InputBoolean.vue"
import TraitsInputDatetime from "./InputDatetime.vue"
import TraitsInputInterval from "./InputInterval.vue"
import TraitsInputSchedule from "./InputSchedule.vue"
import TraitsInputValidity from "./InputValidity.vue"
import TraitsInputColor from "./InputColor.vue"
import TraitsInputList from "./InputList.vue"
import TraitsInputRange from "./InputRange.vue"
import TraitsInputGeo from "./InputGeo.vue"
import TraitsInputSurface from "./InputSurface.vue"

const {locale, key: i18nKey, load, t} = useI18nSection("traits")
await useAsyncData(i18nKey.value, load, {watch: [locale]})

const emit = defineEmits<{(e: "add", trait: TraitInput): void}>()

const key = ref("")
const meta = ref<KeyMeta>(defaultKeyMeta("string"))
const value = ref<TraitDynamicValueModel>(defaultTraitFormValue(meta.value))
const keyInput = ref<{focus: () => void} | null>(null)
const formRef = ref<HTMLFormElement | null>(null)
const valueModel = computed<TraitDynamicValueModel>({
  get: () => value.value,
  set: next => {
    value.value = next
  }
})
const valueComponent = computed(() => {
  switch (resolveDataType(meta.value.dataType)) {
    case "number":
      return TraitsInputNumber
    case "boolean":
      return TraitsInputBoolean
    case "datetime":
      return TraitsInputDatetime
    case "range":
      return TraitsInputRange
    case "interval":
      return TraitsInputInterval
    case "schedule":
      return TraitsInputSchedule
    case "geo":
      return TraitsInputGeo
    case "list":
      return TraitsInputList
    case "validity":
      return TraitsInputValidity
    case "color":
      return TraitsInputColor
    case "surface":
      return TraitsInputSurface
    default:
      return TraitsInputString
  }
})
const showValueLabel = computed(
  () => !["list", "boolean", "range", "interval", "schedule", "geo", "validity", "color", "surface"].includes(resolveDataType(meta.value.dataType))
)
const valueInputId = computed(() => `TraitValue-${resolveDataType(meta.value.dataType) || "string"}`)
const valueLabelFor = computed(() => {
  const dataType = resolveDataType(meta.value.dataType)
  if (dataType === "range" || dataType === "interval") return `${valueInputId.value}-start`
  if (dataType === "schedule") return `${valueInputId.value}-from-day`
  return valueInputId.value
})
const isValueFilled = computed(() => isTraitFormValueFilled(meta.value, value.value))
const captureFocusSnapshot = (): TraitAddFormFocusSnapshot | null => {
  if (!import.meta.client || !formRef.value) return null
  const active = document.activeElement as HTMLElement | null
  if (!active || !formRef.value.contains(active)) return null
  const inputLike = active as HTMLInputElement | HTMLTextAreaElement
  return {
    id: String(active.id || ""),
    name: String(active.getAttribute("name") || ""),
    wasValueField: String(active.id || "").startsWith("TraitValue-"),
    selectionStart: typeof inputLike.selectionStart === "number" ? inputLike.selectionStart : null,
    selectionEnd: typeof inputLike.selectionEnd === "number" ? inputLike.selectionEnd : null
  }
}
const restoreFocusSnapshot = (snapshot: TraitAddFormFocusSnapshot | null) => {
  if (!import.meta.client || !snapshot || !formRef.value) return
  const target =
    (snapshot.id ? (document.getElementById(snapshot.id) as HTMLElement | null) : null) ||
    (snapshot.wasValueField ? (document.getElementById(valueLabelFor.value) as HTMLElement | null) : null) ||
    (snapshot.name ? (formRef.value.querySelector(`[name="${snapshot.name}"]`) as HTMLElement | null) : null)
  if (!target) return
  target.focus()
  if (
    (target instanceof HTMLInputElement || target instanceof HTMLTextAreaElement) &&
    snapshot.selectionStart !== null &&
    snapshot.selectionEnd !== null
  ) {
    try {
      target.setSelectionRange(snapshot.selectionStart, snapshot.selectionEnd)
    } catch {
      // ignore unsupported input types
    }
  }
}
watch(
  () => meta.value.dataType,
  async (nextDataType, prevDataType) => {
    const focusSnapshot = nextDataType !== prevDataType ? captureFocusSnapshot() : null
    meta.value = normalizeKeyMeta(meta.value)
    value.value = defaultTraitFormValue(meta.value)
    if (!focusSnapshot) return
    await nextTick()
    restoreFocusSnapshot(focusSnapshot)
  }
)
const onSubmit = async () => {
  if (!key.value.trim() || !isValueFilled.value) return
  emit("add", {
    t_key: key.value.trim(),
    t_value: serializeTraitFormValue(meta.value, value.value),
    meta: normalizeKeyMeta(meta.value)
  })
  value.value = defaultTraitFormValue(meta.value)
  await nextTick()
  keyInput.value?.focus()
}
onMounted(() => {
  keyInput.value?.focus()
})
</script>
