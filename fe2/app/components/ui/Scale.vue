<script setup lang="ts">
type ScaleOption = {
  value: string
  label?: string
  activeColor?: string
}
type NormalizedOption = {
  value: string
  label: string
  activeColor: string
}
interface Props {
  modelValue?: string | null
  options: ScaleOption[]
  id?: string
  name?: string
  label?: string
  disabled?: boolean
  clearValue?: string
}
const props = withDefaults(defineProps<Props>(), {
  modelValue: "",
  label: "Сложность",
  disabled: false,
  clearValue: ""
})
const emit = defineEmits<{
  "update:modelValue": [value: string]
  change: [value: string]
}>()
const generatedId = useId()
const id = computed(() => props.id || `scale-${generatedId}`)
const labelId = computed(() => `${id.value}-label`)
const value = computed(() => String(props.modelValue ?? "").trim())
const getColor = (i: number, total: number) => {
  const percent = i / Math.max(1, total - 1)
  return percent <= 0.33 ? "#16a34a" : percent <= 0.66 ? "#ca8a04" : "#dc2626"
}
const options = computed<NormalizedOption[]>(() =>
  props.options
    .map((option, i) => ({
      value: String(option.value).trim(),
      label: option.label?.trim() || `Уровень ${i + 1}`,
      activeColor: option.activeColor || getColor(i, props.options.length)
    }))
    .filter(option => option.value)
)
const selectedIndex = computed(() =>
  options.value.findIndex(option => option.value === value.value)
)
const gridStyle = computed(() => ({
  gridTemplateColumns: `repeat(${options.value.length}, minmax(0, 1fr))`
}))
const isSelected = (option: NormalizedOption) => option.value === value.value
const setValue = (option: NormalizedOption) => {
  if (props.disabled) return
  const next = isSelected(option) ? props.clearValue : option.value
  emit("update:modelValue", next)
  emit("change", next)
}
</script>
<template>
  <div class="inline-flex flex-col gap-1">
    <span :id="labelId" class="text-sm">
      <slot name="label">{{ label }}</slot>
    </span>
    <div
      :id="id"
      :aria-disabled="disabled ? 'true' : 'false'"
      :aria-labelledby="labelId"
      class="grid gap-1"
      :style="gridStyle"
      role="radiogroup"
    >
      <button
        v-for="(option, i) in options"
        :key="option.value"
        :aria-checked="isSelected(option) ? 'true' : 'false'"
        :aria-label="option.label"
        :disabled="disabled"
        :style="{
          backgroundColor:
            selectedIndex >= 0 && i <= selectedIndex ? option.activeColor : 'var(--elevated)'
        }"
        :title="option.label"
        :class="[
          'ui-focus h-5 min-w-0 rounded-lg',
          disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'
        ]"
        role="radio"
        type="button"
        @click="setValue(option)"
      >
        <span class="sr-only">{{ option.label }}</span>
      </button>
    </div>
    <div class="grid gap-1" :style="gridStyle">
      <span
        v-for="option in options"
        :key="`label-${option.value}`"
        :class="[
          'truncate text-center text-sm',
          isSelected(option) ? 'text-(--accent)' : 'opacity-70'
        ]"
        >{{ option.label }}</span
      >
    </div>
    <input v-if="name" type="hidden" :name="name" :value="value" />
  </div>
</template>
