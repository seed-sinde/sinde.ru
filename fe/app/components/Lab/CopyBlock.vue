<template>
  <div class="relative inline-flex max-w-full">
    <LabBaseButton
      :button-class="[
        'group inline-flex max-w-full items-center gap-2 px-2.5 py-1 text-left transition-colors duration-200',
        stateBlockClass,
        buttonClass
      ]"
      :aria-label="resolvedTitle"
      :icon="iconName"
      :icon-class="['order-last h-4 w-4 transition-colors', iconClass]"
      :label="resolvedTitle"
      @click="copyValue">
      <span class="text-[10px] uppercase tracking-[0.14em] transition-colors" :class="stateLabelClass">
        {{ label }}
      </span>
      <span class="min-w-0 break-all text-xs" :class="valueClass">{{ value }}</span>
    </LabBaseButton>
    <span
      v-if="showStateTooltip && copyState !== 'idle'"
      class="pointer-events-none absolute -top-1 right-2 -translate-y-full px-2 py-1 text-[11px] leading-none"
      :class="stateTooltipClass"
      role="status"
      aria-live="polite">
      {{ resolvedTitle }}
    </span>
  </div>
</template>
<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      label: string
      value: string
      variant?: CopyBlockVariant
      buttonClass?: string
      titleIdle?: string
      titleSuccess?: string
      titleError?: string
      monospace?: boolean
      showStateTooltip?: boolean
    }>(),
    {
      variant: 'default',
      buttonClass: '',
      titleIdle: 'Скопировать',
      titleSuccess: 'Скопировано',
      titleError: 'Ошибка копирования',
      monospace: true,
      showStateTooltip: false
    }
  )
  const emit = defineEmits<{
    copied: []
    error: []
  }>()
  const { copyFrom, copied } = useClipboard()
  const copyFailed = ref(false)
  let copyFailedTimer: ReturnType<typeof setTimeout> | null = null
  const clearCopyFailed = () => {
    copyFailed.value = false
    if (copyFailedTimer) {
      clearTimeout(copyFailedTimer)
      copyFailedTimer = null
    }
  }
  const markCopyFailed = () => {
    copyFailed.value = true
    if (copyFailedTimer) clearTimeout(copyFailedTimer)
    copyFailedTimer = setTimeout(() => {
      copyFailed.value = false
      copyFailedTimer = null
    }, 1400)
  }
  const copyState = computed<'idle' | 'success' | 'error'>(() => {
    if (copyFailed.value) return 'error'
    if (copied.value) return 'success'
    return 'idle'
  })
  const iconName = computed(() => {
    if (copyState.value === 'success') return 'ic:round-check'
    if (copyState.value === 'error') return 'ic:round-error'
    return 'ic:baseline-content-copy'
  })
  const resolvedTitle = computed(() => {
    if (copyState.value === 'success') return props.titleSuccess
    if (copyState.value === 'error') return props.titleError
    return props.titleIdle
  })
  const stateBlockClass = computed(() => {
    if (copyState.value === 'success') {
      return 'bg-[color-mix(in_srgb,var(--lab-success)_14%,var(--lab-bg-surface))] text-(--lab-text-primary) ring-(--lab-success) ring-1'
    }
    if (copyState.value === 'error') {
      return 'bg-[color-mix(in_srgb,var(--lab-danger)_14%,var(--lab-bg-surface))] text-(--lab-text-primary) ring-(--lab-danger) ring-1'
    }
    if (props.variant === 'dark-cyan') {
      return 'bg-[color-mix(in_srgb,var(--lab-info)_12%,var(--lab-bg-surface))] text-(--lab-text-primary) ring-(--lab-border) hover:ring-(--lab-info) ring-1'
    }
    return 'bg-(--lab-bg-surface) text-(--lab-text-primary) ring-(--lab-border) hover:ring-(--lab-text-muted) ring-1'
  })
  const stateLabelClass = computed(() => {
    if (copyState.value === 'success') return 'text-(--lab-success)'
    if (copyState.value === 'error') return 'text-(--lab-danger)'
    if (props.variant === 'dark-cyan') return 'text-(--lab-info) group-hover:text-(--lab-text-primary)'
    return 'text-(--lab-text-muted) group-hover:text-(--lab-text-primary)'
  })
  const iconClass = computed(() => {
    if (copyState.value === 'success') return 'animate-pulse text-(--lab-success)'
    if (copyState.value === 'error') return 'animate-pulse text-(--lab-danger)'
    if (props.variant === 'dark-cyan') return 'text-(--lab-info) group-hover:text-(--lab-text-primary)'
    return 'text-(--lab-text-muted) group-hover:text-(--lab-text-primary)'
  })
  const valueClass = computed(() => {
    if (props.monospace) return 'font-mono'
    return ''
  })
  const stateTooltipClass = computed(() => {
    if (copyState.value === 'success')
      return 'bg-[color-mix(in_srgb,var(--lab-success)_18%,var(--lab-bg-surface))] text-(--lab-text-primary) ring-(--lab-success) ring-1'
    if (copyState.value === 'error')
      return 'bg-[color-mix(in_srgb,var(--lab-danger)_18%,var(--lab-bg-surface))] text-(--lab-text-primary) ring-(--lab-danger) ring-1'
    return 'bg-(--lab-bg-surface) text-(--lab-text-primary) ring-(--lab-border) ring-1'
  })
  const copyValue = async () => {
    try {
      await copyFrom(props.value)
      clearCopyFailed()
      emit('copied')
    } catch {
      markCopyFailed()
      emit('error')
    }
  }
  onBeforeUnmount(() => {
    if (copyFailedTimer) {
      clearTimeout(copyFailedTimer)
    }
  })
</script>
