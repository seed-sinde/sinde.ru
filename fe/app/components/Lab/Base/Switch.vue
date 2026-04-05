<template>
  <label
    :for="switchId"
    :class="[
      'relative inline-flex select-none items-center gap-2',
      disabled ? 'cursor-default opacity-70' : 'cursor-pointer',
      containerClass
    ]">
    <input
      :id="switchId"
      ref="inputRef"
      :name="name"
      type="checkbox"
      class="peer sr-only"
      :checked="modelValue"
      :disabled="disabled"
      @change="onChange" />
    <span
      aria-hidden="true"
      :class="[
        'relative inline-flex h-5 w-9 shrink-0 rounded-full border transition-colors',
        resolvedTrackClass,
        trackClass
      ]"
      :style="resolvedTrackStyle">
      <span
        :class="[
          'absolute left-0.5 top-0.5 h-3.5 w-3.5 rounded-full border transition-transform transition-colors',
          resolvedThumbClass,
          thumbClass
        ]"
        :style="resolvedThumbStyle"></span>
    </span>
    <span v-if="label" :class="['whitespace-nowrap text-[11px] text-zinc-300', labelClass]">
      {{ label }}
    </span>
    <span
      v-if="showStateLabel"
      :class="['whitespace-nowrap text-[11px] uppercase tracking-[0.06em] text-zinc-500', stateClass]">
      {{ resolvedStateLabel }}
    </span>
  </label>
</template>
<script setup lang="ts">
  import { useId } from 'vue'
  const props = withDefaults(
    defineProps<{
      modelValue: boolean
      id?: string
      name?: string
      label?: string
      trueLabel?: string
      falseLabel?: string
      partialLabel?: string
      disabled?: boolean
      tone?: ToggleTone
      visualState?: ToggleVisualState
      containerClass?: string
      labelClass?: string
      stateClass?: string
      trackClass?: string
      thumbClass?: string
    }>(),
    {
      id: '',
      name: '',
      label: '',
      trueLabel: '',
      falseLabel: '',
      partialLabel: '',
      disabled: false,
      tone: 'cyan',
      visualState: undefined,
      containerClass: '',
      labelClass: '',
      stateClass: '',
      trackClass: '',
      thumbClass: ''
    }
  )
  const emit = defineEmits<{ (e: 'update:modelValue', value: boolean): void }>()
  const uid = useId()
  const inputRef = ref<HTMLInputElement | null>(null)
  const switchId = computed(() => props.id || `LabToggle-${uid}`)
  const resolvedVisualState = computed<ToggleVisualState>(() => props.visualState || (props.modelValue ? 'on' : 'off'))
  const showStateLabel = computed(() => Boolean(props.trueLabel || props.falseLabel || props.partialLabel))
  const resolvedStateLabel = computed(() => {
    if (resolvedVisualState.value === 'partial') return props.partialLabel || props.falseLabel || props.trueLabel
    return props.modelValue ? props.trueLabel : props.falseLabel
  })
  const toneTrackClass = {
    cyan: {
      on: 'border-cyan-500/70 bg-cyan-600/30'
    },
    emerald: {
      on: 'border-emerald-500/70 bg-emerald-600/30'
    },
    rose: {
      on: 'border-rose-500/70 bg-rose-600/30'
    },
    amber: {
      on: 'border-amber-500/70 bg-amber-600/30'
    }
  } satisfies Record<ToggleTone, Record<'on', string>>
  const toneThumbClass = {
    cyan: {
      on: 'translate-x-4 border-cyan-50 bg-cyan-100'
    },
    emerald: {
      on: 'translate-x-4 border-emerald-50 bg-emerald-100'
    },
    rose: {
      on: 'translate-x-4 border-rose-50 bg-rose-100'
    },
    amber: {
      on: 'translate-x-4 border-amber-50 bg-amber-100'
    }
  } satisfies Record<ToggleTone, Record<'on', string>>
  const partialTrackStyle = {
    borderColor: 'color-mix(in srgb, var(--lab-warning) 78%, var(--lab-border))',
    background: 'color-mix(in srgb, var(--lab-warning) 24%, var(--lab-bg-control))'
  }
  const partialThumbStyle = {
    borderColor: 'color-mix(in srgb, var(--lab-warning) 70%, var(--lab-text-primary))',
    background: 'color-mix(in srgb, var(--lab-warning) 56%, var(--lab-bg-surface))'
  }
  const resolvedTrackClass = computed(() => {
    if (resolvedVisualState.value === 'off') return 'border-zinc-600 bg-zinc-800'
    if (resolvedVisualState.value === 'partial') return ''
    return toneTrackClass[props.tone][resolvedVisualState.value]
  })
  const resolvedThumbClass = computed(() => {
    if (resolvedVisualState.value === 'off') return 'translate-x-0 border-zinc-400/40 bg-zinc-300'
    if (resolvedVisualState.value === 'partial') return 'translate-x-2'
    return toneThumbClass[props.tone][resolvedVisualState.value]
  })
  const resolvedTrackStyle = computed(() => (resolvedVisualState.value === 'partial' ? partialTrackStyle : undefined))
  const resolvedThumbStyle = computed(() => (resolvedVisualState.value === 'partial' ? partialThumbStyle : undefined))
  watchEffect(() => {
    if (!inputRef.value) return
    inputRef.value.indeterminate = resolvedVisualState.value === 'partial'
  })
  const onChange = (event: Event) => {
    const target = event.target as HTMLInputElement | null
    if (!target || props.disabled) return
    emit('update:modelValue', target.checked)
  }
</script>
