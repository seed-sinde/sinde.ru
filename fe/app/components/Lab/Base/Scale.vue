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
  <LabField :field-class="fieldClass" :label-class="labelClass">
    <template #label>
      <span :id="labelId">{{ label }}</span>
    </template>
    <div
      v-if="normalizedOptions.length"
      :id="resolvedId"
      :aria-disabled="disabled ? 'true' : 'false'"
      :aria-labelledby="labelId"
      class="grid gap-1"
      :class="disabled ? 'scale-disabled' : ''"
      :style="{ gridTemplateColumns: `repeat(${normalizedOptions.length}, minmax(0, 1fr))` }"
      role="radiogroup"
      @keydown="onKeydown">
      <button
        v-for="(option, index) in normalizedOptions"
        :key="option.value"
        :aria-checked="isSelected(option.value) ? 'true' : 'false'"
        :aria-label="option.label"
        :disabled="disabled"
        :style="{ '--scale-segment-color': option.activeColor }"
        :title="option.label"
        class="h-8 min-w-0 rounded border outline-none transition-colors"
        :class="[
          isActive(index) ? 'scale-segment-active' : 'scale-segment-idle',
          isSelected(option.value) ? 'scale-segment-selected' : '',
          disabled ? 'scale-segment-disabled' : 'cursor-pointer'
        ]"
        role="radio"
        type="button"
        @click="selectValue(index)">
        <span class="sr-only">{{ option.label }}</span>
      </button>
    </div>
    <div
      v-if="normalizedOptions.length"
      class="grid gap-1"
      :class="disabled ? 'scale-labels-disabled' : ''"
      :style="{ gridTemplateColumns: `repeat(${normalizedOptions.length}, minmax(0, 1fr))` }">
      <div
        v-for="option in normalizedOptions"
        :key="`label:${option.value}`"
        class="min-w-0 text-center text-xs leading-4"
        :class="isSelected(option.value) ? 'text-(--lab-text-primary)' : 'text-(--lab-text-soft)'">
        <span class="block truncate">{{ option.label }}</span>
      </div>
    </div>
    <input v-if="name" type="hidden" :name="name" :value="normalizedValue" />
  </LabField>
</template>
<style scoped>
  .scale-segment-idle {
    background: var(--lab-bg-control);
  }
  .scale-segment-idle:hover:enabled {
    border-color: var(--lab-border-strong);
    background: var(--lab-bg-surface-hover);
  }
  .scale-segment-idle:focus-visible {
    border-color: var(--lab-accent);
  }
  .scale-segment-active {
    border-color: color-mix(in srgb, var(--scale-segment-color) 55%, var(--lab-border));
    background: color-mix(in srgb, var(--scale-segment-color) 20%, var(--lab-bg-control));
  }
  .scale-segment-active:hover:enabled {
    border-color: color-mix(in srgb, var(--scale-segment-color) 72%, var(--lab-border));
    background: color-mix(in srgb, var(--scale-segment-color) 28%, var(--lab-bg-control));
  }
  .scale-segment-active:focus-visible {
    border-color: var(--scale-segment-color);
  }
  .scale-segment-selected {
    border-color: var(--lab-text-primary);
  }
  .scale-segment-disabled {
    cursor: not-allowed;
  }
  .scale-disabled .scale-segment-idle {
    border-color: color-mix(in srgb, var(--lab-border) 85%, transparent);
    background: color-mix(in srgb, var(--lab-bg-control) 78%, transparent);
  }
  .scale-disabled .scale-segment-active {
    border-color: color-mix(in srgb, var(--scale-segment-color) 28%, var(--lab-border));
    background: color-mix(in srgb, var(--scale-segment-color) 12%, var(--lab-bg-control));
  }
  .scale-disabled .scale-segment-selected {
    border-color: color-mix(in srgb, var(--lab-text-primary) 50%, var(--lab-border));
  }
  .scale-labels-disabled {
    opacity: 0.55;
  }
</style>
