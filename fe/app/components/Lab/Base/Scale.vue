<script setup lang="ts">
const model = defineModel<string | null>({ default: '' })
const props = withDefaults(
  defineProps<{
    options: DifficultyScaleOption[]
    label?: string
    id?: string
    name?: string
    disabled?: boolean
    clearValue?: string
    fieldClass?: string
    labelClass?: string
  }>(),
  {
    label: 'Сложность',
    id: '',
    name: '',
    disabled: false,
    clearValue: 'all',
    fieldClass: '',
    labelClass: ''
  }
)
const uid = useId()
const resolvedId = computed(() => props.id || `scale-${uid}`)
const labelId = computed(() => `${resolvedId.value}-label`)
const normalizedValue = computed(() => String(model.value || '').trim())
const fallbackScaleColor = (index: number, total: number) => {
  const maxIndex = Math.max(1, total - 1)
  const ratio = index / maxIndex
  if (ratio <= 0.33) return 'var(--lab-success)'
  if (ratio <= 0.66) return 'var(--lab-warning)'
  return 'var(--lab-danger)'
}
const normalizedOptions = computed(() =>
  props.options
    .map((option, index) => {
      const value = String(option.value || '').trim()
      if (!value) return null
      return {
        value,
        label: String(option.label || '').trim() || `Уровень ${index + 1}`,
        activeColor: String(option.activeColor || '').trim() || fallbackScaleColor(index, props.options.length)
      }
    })
    .filter((option): option is { value: string; label: string; activeColor: string } => Boolean(option))
)
const selectedIndex = computed(() =>
  normalizedOptions.value.findIndex(option => option.value === normalizedValue.value)
)
const isSelected = (value: string) => value === normalizedValue.value
const isActive = (index: number) => selectedIndex.value >= 0 && index <= selectedIndex.value
const segmentClass = (option: { value: string; label: string; activeColor: string }, index: number) => [
  'lab-focus h-8 min-w-0 border transition-colors',
  props.disabled ? 'cursor-not-allowed' : 'cursor-pointer',
  isActive(index)
    ? props.disabled
      ? 'border-[color-mix(in_srgb,var(--scale-segment-color)_28%,var(--lab-border))] bg-[color-mix(in_srgb,var(--scale-segment-color)_12%,var(--lab-bg-control))] hover:enabled:border-[color-mix(in_srgb,var(--scale-segment-color)_28%,var(--lab-border))] hover:enabled:bg-[color-mix(in_srgb,var(--scale-segment-color)_12%,var(--lab-bg-control))]'
      : 'border-[color-mix(in_srgb,var(--scale-segment-color)_55%,var(--lab-border))] bg-[color-mix(in_srgb,var(--scale-segment-color)_20%,var(--lab-bg-control))] hover:enabled:border-[color-mix(in_srgb,var(--scale-segment-color)_72%,var(--lab-border))] hover:enabled:bg-[color-mix(in_srgb,var(--scale-segment-color)_28%,var(--lab-bg-control))]'
    : props.disabled
      ? 'border-[color-mix(in_srgb,var(--lab-border)_85%,transparent)] bg-[color-mix(in_srgb,var(--lab-bg-control)_78%,transparent)]'
      : 'bg-(--lab-bg-control) hover:enabled:border-(--lab-border-strong) hover:enabled:bg-(--lab-bg-surface-hover)',
  isSelected(option.value)
    ? props.disabled
      ? 'border-[color-mix(in_srgb,var(--lab-text-primary)_50%,var(--lab-border))]'
      : 'border-(--lab-text-primary)'
    : ''
]
const selectValue = (index: number) => {
  if (props.disabled || index < 0 || index >= normalizedOptions.value.length) return
  const nextValue = normalizedOptions.value[index]?.value || ''
  if (!nextValue) return
  model.value = nextValue === normalizedValue.value ? props.clearValue : nextValue
}
const onKeydown = (event: KeyboardEvent) => {
  if (props.disabled || normalizedOptions.value.length === 0) return
  const lastIndex = normalizedOptions.value.length - 1
  const currentIndex = selectedIndex.value >= 0 ? selectedIndex.value : 0
  if (event.key === 'ArrowRight' || event.key === 'ArrowUp') {
    event.preventDefault()
    model.value = normalizedOptions.value[Math.min(lastIndex, currentIndex + 1)]?.value || normalizedValue.value
  }
  if (event.key === 'ArrowLeft' || event.key === 'ArrowDown') {
    event.preventDefault()
    model.value = normalizedOptions.value[Math.max(0, currentIndex - 1)]?.value || normalizedValue.value
  }
  if (event.key === 'Home') {
    event.preventDefault()
    model.value = normalizedOptions.value[0]?.value || normalizedValue.value
  }
  if (event.key === 'End') {
    event.preventDefault()
    model.value = normalizedOptions.value[lastIndex]?.value || normalizedValue.value
  }
}
</script>
<template>
  <LabBaseField :field-class="fieldClass" :label-class="labelClass">
    <template #label>
      <span :id="labelId">{{ label }}</span>
    </template>
    <div
      v-if="normalizedOptions.length"
      :id="resolvedId"
      :aria-disabled="disabled ? 'true' : 'false'"
      :aria-labelledby="labelId"
      class="grid gap-1"
      :style="{ gridTemplateColumns: `repeat(${normalizedOptions.length}, minmax(0, 1fr))` }"
      role="radiogroup"
      @keydown="onKeydown"
    >
      <button
        v-for="(option, index) in normalizedOptions"
        :key="option.value"
        :aria-checked="isSelected(option.value) ? 'true' : 'false'"
        :aria-label="option.label"
        :disabled="disabled"
        :style="{ '--scale-segment-color': option.activeColor }"
        :title="option.label"
        :class="segmentClass(option, index)"
        role="radio"
        type="button"
        @click="selectValue(index)"
      >
        <span class="sr-only">{{ option.label }}</span>
      </button>
    </div>
    <div
      v-if="normalizedOptions.length"
      class="grid gap-1"
      :class="disabled ? 'opacity-55' : ''"
      :style="{ gridTemplateColumns: `repeat(${normalizedOptions.length}, minmax(0, 1fr))` }"
    >
      <div
        v-for="option in normalizedOptions"
        :key="`label:${option.value}`"
        class="min-w-0 text-center text-xs leading-4"
        :class="isSelected(option.value) ? 'text-(--lab-text-primary)' : 'text-(--lab-text-soft)'"
      >
        <span class="block truncate">{{ option.label }}</span>
      </div>
    </div>
    <input v-if="name" type="hidden" :name="name" :value="normalizedValue" />
  </LabBaseField>
</template>
