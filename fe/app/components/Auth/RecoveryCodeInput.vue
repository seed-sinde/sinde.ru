<template>
  <div class="max-w-sm space-y-2">
    <label v-if="label" :for="id" class="text-(--lab-text-secondary) text-sm">{{ label }}</label>
    <div class="relative" @click="focus">
      <input
        :id="id"
        ref="inputRef"
        :value="normalizedValue"
        :name="name"
        type="text"
        inputmode="text"
        :autocomplete="autocomplete"
        autocapitalize="characters"
        spellcheck="false"
        data-1p-ignore="true"
        data-lpignore="true"
        class="absolute inset-0 z-10 h-full w-full cursor-text opacity-0"
        :aria-label="ariaLabel || label || `Код сброса на ${length} символов`"
        @input="onInput"
        @focus="isFocused = true"
        @blur="isFocused = false" />
      <div :class="segmentsWrapClass">
        <div :class="segmentClass(leftValue, 'left')">
          {{ leftValue }}
        </div>
        <span class="text-(--lab-text-muted) text-lg font-semibold">-</span>
        <div :class="segmentClass(rightValue, 'right')">
          {{ rightValue }}
        </div>
      </div>
    </div>
    <p
      v-if="resolvedHint"
      class="text-xs"
      :class="rejectedHint ? 'text-(--lab-danger)' : 'text-(--lab-text-muted)'">
      {{ resolvedHint }}
    </p>
  </div>
</template>
<script setup lang="ts">
  const { rejectedTarget, rejectedHint, markRejected, clearRejectFeedback } = useInputRejectFeedback<'left' | 'right'>()
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
      length: 8,
      autocomplete: 'one-time-code',
      invalid: false,
      valid: false
    }
  )
  const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void
    (e: 'complete', value: string): void
  }>()
  const inputRef = ref<HTMLInputElement | null>(null)
  const lastCompletedValue = ref('')
  const isFocused = ref(false)
  const normalizedValue = computed(() =>
    String(props.modelValue || '')
      .toUpperCase()
      .replace(/[^A-Z0-9]+/g, '')
      .slice(0, props.length)
  )
  const resolvedHint = computed(() => rejectedHint.value || props.hint)
  const leftValue = computed(() => normalizedValue.value.slice(0, props.length / 2))
  const rightValue = computed(() => normalizedValue.value.slice(props.length / 2, props.length))
  const segmentsWrapClass = computed(() => [
    'flex items-center justify-center gap-3 p-1 transition-colors',
    props.invalid ?
      'ring-(--lab-danger) ring-1'
    : props.valid ?
      'ring-(--lab-success) ring-1'
    : ''
  ])
  const activeSegment = computed(() => {
    if (normalizedValue.value.length >= props.length / 2) return 'right'
    return 'left'
  })
  const segmentClass = (value: string, side: 'left' | 'right') => {
    const isSuccess = props.valid && !props.invalid
    const isError = props.invalid
    return [
      'border-(--lab-border) text-(--lab-text-primary) bg-(--lab-bg-surface) flex h-11 min-w-0 flex-1 items-center justify-center border px-3 font-mono text-lg font-semibold tracking-[0.22em] transition sm:h-13 sm:text-xl',
      value ? 'bg-[color-mix(in_srgb,var(--lab-info)_12%,var(--lab-bg-surface))]' : 'text-(--lab-text-muted)',
      isFocused.value && activeSegment.value === side ? 'ring-(--lab-accent) ring-2' : '',
      isFocused.value && activeSegment.value === side
        ? 'after:bg-(--lab-accent) after:block after:h-7 after:w-[2px] after:animate-pulse sm:after:h-8'
        : '',
      isSuccess ? 'ring-(--lab-success) bg-[color-mix(in_srgb,var(--lab-success)_14%,var(--lab-bg-surface))] ring-2' : '',
      isError ? 'ring-(--lab-danger) bg-[color-mix(in_srgb,var(--lab-danger)_14%,var(--lab-bg-surface))] ring-2' : '',
      rejectedTarget.value === side ?
        'ring-(--lab-danger) bg-[color-mix(in_srgb,var(--lab-danger)_14%,var(--lab-bg-surface))] ring-2'
      : ''
    ]
  }
  const onInput = (event: Event) => {
    const target = event.target as HTMLInputElement | null
    if (!target) return
    const rawValue = String(target.value || '').toUpperCase()
    const previousLength = normalizedValue.value.length
    const nextValue = rawValue.replace(/[^A-Z0-9]+/g, '').slice(0, props.length)
    target.value = nextValue
    if (rawValue !== nextValue) {
      markRejected(previousLength >= props.length / 2 ? 'right' : 'left', 'Только A-Z и 0-9 символы.')
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
