<template>
  <div class="max-w-sm space-y-2">
    <label v-if="label" :for="id" class="text-sm text-(--lab-text-secondary)">{{ label }}</label>
    <div class="relative" @click="focus">
      <input
        :id="id"
        ref="inputRef"
        :value="normalizedValue"
        :name="name"
        type="text"
        inputmode="numeric"
        pattern="[0-9]*"
        :autocomplete="autocomplete"
        autocapitalize="none"
        spellcheck="false"
        data-1p-ignore="true"
        data-lpignore="true"
        class="absolute inset-0 z-10 h-full w-full cursor-text opacity-0"
        :aria-label="ariaLabel || label || `Код на ${length} цифр`"
        @input="onInput"
        @focus="isFocused = true"
        @blur="isFocused = false"
      />
      <div class="grid gap-2" :style="{ gridTemplateColumns: `repeat(${length}, minmax(0, 1fr))` }">
        <div v-for="(digit, index) in digits" :key="index" :class="getDigitClass(index)">
          {{ digit }}
        </div>
      </div>
    </div>
    <p v-if="hint" class="text-xs text-(--lab-text-muted)">{{ hint }}</p>
  </div>
</template>
<script setup lang="ts">
const { rejectedTarget, markRejected, clearRejectFeedback } = useInputRejectFeedback<number>()
const props = withDefaults(
  defineProps<{
    modelValue?: string
    id?: string
    name?: string
    label?: string
    hint?: string
    ariaLabel?: string
    length?: number
    autocomplete?: string
    invalid?: boolean
    valid?: boolean
  }>(),
  {
    modelValue: '',
    id: '',
    name: '',
    label: '',
    hint: '',
    ariaLabel: '',
    length: 6,
    autocomplete: 'one-time-code',
    invalid: false,
    valid: false
  }
)
const emit = defineEmits<{
  (e: 'update:modelValue' | 'complete', value: string): void
}>()
const inputRef = ref<HTMLInputElement | null>(null)
const lastCompletedValue = ref('')
const isFocused = ref(false)
const normalizedValue = computed(() =>
  String(props.modelValue || '')
    .replace(/\D+/g, '')
    .slice(0, props.length)
)
const digits = computed(() => Array.from({ length: props.length }, (_, index) => normalizedValue.value[index] || ''))
const getDigitClass = (index: number) => {
  const isFilled = index < normalizedValue.value.length
  const isActive =
    index === Math.min(normalizedValue.value.length, props.length - 1) && normalizedValue.value.length < props.length
  const isSuccess = props.valid && !props.invalid
  const isError = props.invalid
  return [
    'border-(--lab-border) text-(--lab-text-primary) bg-(--lab-bg-surface) flex h-11 items-center justify-center border text-center text-xl font-semibold tabular-nums transition sm:h-13 sm:text-2xl',
    isFilled
      ? 'text-(--lab-text-primary) bg-[color-mix(in_srgb,var(--lab-info)_12%,var(--lab-bg-surface))]'
      : 'text-(--lab-text-muted)',
    isActive && isFocused.value ? 'ring-(--lab-accent) ring-2' : isActive ? 'ring-(--lab-border) ring-1' : '',
    isActive && isFocused.value
      ? 'after:bg-(--lab-accent) after:block after:h-7 after:w-[2px] after:animate-pulse sm:after:h-8'
      : '',
    isSuccess ? 'ring-(--lab-success) bg-[color-mix(in_srgb,var(--lab-success)_14%,var(--lab-bg-surface))] ring-2' : '',
    isError ? 'ring-(--lab-danger) bg-[color-mix(in_srgb,var(--lab-danger)_14%,var(--lab-bg-surface))] ring-2' : '',
    rejectedTarget.value === index
      ? 'ring-(--lab-danger) bg-[color-mix(in_srgb,var(--lab-danger)_14%,var(--lab-bg-surface))] ring-2'
      : ''
  ]
}
const onInput = (event: Event) => {
  const target = event.target as HTMLInputElement | null
  if (!target) return
  const rawValue = String(target.value || '')
  const previousLength = normalizedValue.value.length
  const nextValue = rawValue.replace(/\D+/g, '').slice(0, props.length)
  target.value = nextValue
  if (rawValue !== nextValue) {
    markRejected(Math.min(previousLength, props.length - 1))
  } else {
    clearRejectFeedback()
  }
  emit('update:modelValue', nextValue)
  if (nextValue.length === props.length && nextValue !== lastCompletedValue.value) {
    lastCompletedValue.value = nextValue
    emit('complete', nextValue)
    return
  }
  if (nextValue !== lastCompletedValue.value) {
    lastCompletedValue.value = ''
  }
}
watch(normalizedValue, nextValue => {
  if (nextValue !== lastCompletedValue.value && nextValue.length < props.length) {
    lastCompletedValue.value = ''
  }
  if (!nextValue) {
    clearRejectFeedback()
  }
})
const focus = () => {
  inputRef.value?.focus()
}
defineExpose({ focus })
</script>
