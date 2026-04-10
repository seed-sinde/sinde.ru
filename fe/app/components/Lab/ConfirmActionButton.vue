<template>
  <div ref="rootRef" class="relative inline-flex">
    <LabBaseButton
      :type="type"
      :disabled="disabled"
      :icon-only="iconOnly"
      :aria-label="resolvedAriaLabel"
      :button-class="resolvedButtonClass"
      :title="title || undefined"
      @click="handleClick">
      <span
        v-if="armed && !iconOnly"
        class="pointer-events-none absolute inset-y-0 left-0 z-0"
        :class="resolvedProgressClass"
        :style="{ width: `${progress}%` }" />
      <svg
        v-if="armed && iconOnly"
        class="pointer-events-none absolute inset-0 z-0 h-full w-full -rotate-90 p-0.5"
        viewBox="0 0 40 40"
        aria-hidden="true">
        <circle cx="20" cy="20" r="18" fill="none" stroke="#ffffff" stroke-width="2" class="opacity-20" />
        <circle
          cx="20"
          cy="20"
          r="18"
          fill="none"
          stroke="#ffffff"
          stroke-width="2"
          stroke-linecap="round"
          pathLength="100"
          :stroke-dasharray="100"
          :stroke-dashoffset="100 - progress" />
      </svg>
      <span class="relative z-10 inline-flex items-center justify-center" :class="iconOnly ? 'h-full w-full' : 'gap-2'">
        <Icon
          v-if="resolvedIcon"
          :name="resolvedIcon"
          class="shrink-0 leading-none"
          :class="iconOnly ? 'text-lg' : 'text-base'" />
        <span v-if="!iconOnly">{{ armed ? confirmLabel : label }}</span>
      </span>
    </LabBaseButton>
  </div>
  <Teleport to="body">
    <div
      v-if="armed && tooltip"
      ref="tooltipRef"
      class="pointer-events-none fixed z-10000 max-w-xs whitespace-nowrap border px-2 py-1 text-sm"
      :class="resolvedTooltipClass"
      :style="resolvedTooltipStyle">
      <span class="inline-flex items-center gap-1">
        <Icon name="ic:round-warning-amber" class="text-(--lab-warning) shrink-0 text-sm" />
        {{ tooltip }}
      </span>
    </div>
  </Teleport>
</template>
<script setup lang="ts">
  defineOptions({ inheritAttrs: false })
  const props = withDefaults(
    defineProps<{
      label?: string
      confirmLabel?: string
      tooltip?: string
      timeoutMs?: number
      disabled?: boolean
      type?: 'button' | 'submit' | 'reset'
      class?: string
      buttonClass?: string
      idleClass?: string
      confirmClass?: string
      progressClass?: string
      tooltipClass?: string
      icon?: string
      confirmIcon?: string
      iconOnly?: boolean
      title?: string
      ariaLabel?: string
      confirmAriaLabel?: string
      variant?: LabButtonVariant
    }>(),
    {
      label: '',
      confirmLabel: 'Подтвердить',
      tooltip: '',
      timeoutMs: 5000,
      disabled: false,
      type: 'button',
      class: '',
      buttonClass: '',
      idleClass: '',
      confirmClass: '',
      progressClass: '',
      tooltipClass: '',
      icon: '',
      confirmIcon: 'ic:round-check',
      iconOnly: false,
      title: '',
      ariaLabel: '',
      confirmAriaLabel: '',
      variant: 'info'
    }
  )
  const emit = defineEmits<{
    confirm: []
  }>()
  const rootRef = ref<HTMLElement | null>(null)
  const tooltipRef = ref<HTMLElement | null>(null)
  const armed = ref(false)
  const progress = ref(0)
  const tooltipPosition = reactive({
    left: 0,
    top: 0
  })
  let timer: number | null = null
  const resolvedButtonClass = computed(() => [
    'relative isolate',
    props.iconOnly ? '' : 'overflow-hidden',
    armed.value
      ? [
          'border-(--lab-warning)',
          'bg-(--lab-warning)',
          'text-white',
          'hover:enabled:border-(--lab-warning-hover)',
          'hover:enabled:bg-(--lab-warning-hover)',
          'hover:enabled:text-white'
        ]
      : '',
    armed.value ? props.confirmClass : props.idleClass,
    props.class,
    props.buttonClass
  ])
  const resolvedTooltipClass = computed(() => [props.tooltipClass])
  const resolvedProgressClass = computed(() => [
    'bg-[color-mix(in_srgb,var(--lab-warning)_28%,transparent)]',
    props.progressClass
  ])
  const resolvedTooltipStyle = computed(() => ({
    left: `${tooltipPosition.left}px`,
    top: `${tooltipPosition.top}px`,
    borderColor: 'var(--lab-warning)',
    background: 'var(--lab-bg-overlay)',
    color: 'var(--lab-text-primary)'
  }))
  const resolvedIcon = computed(() => {
    if (armed.value && props.confirmIcon) return props.confirmIcon
    return props.icon || ''
  })
  const resolvedAriaLabel = computed(() => {
    if (armed.value) {
      return props.confirmAriaLabel || props.ariaLabel || props.confirmLabel
    }
    return props.ariaLabel || props.label || props.confirmLabel
  })
  const clearTimer = () => {
    if (timer !== null) {
      window.clearInterval(timer)
      timer = null
    }
  }
  const updateTooltipPosition = () => {
    if (!import.meta.client) return
    const root = rootRef.value
    const tooltip = tooltipRef.value
    if (!root || !tooltip) return
    const spacing = 8
    const viewportPadding = 8
    const rootRect = root.getBoundingClientRect()
    const tooltipRect = tooltip.getBoundingClientRect()
    const viewportWidth = window.innerWidth
    const viewportHeight = window.innerHeight
    const centeredLeft = rootRect.left + rootRect.width / 2 - tooltipRect.width / 2
    const minLeft = viewportPadding
    const maxLeft = Math.max(viewportPadding, viewportWidth - tooltipRect.width - viewportPadding)
    const left = Math.min(Math.max(centeredLeft, minLeft), maxLeft)
    const fitsAbove = rootRect.top >= tooltipRect.height + spacing + viewportPadding
    const top = fitsAbove
      ? rootRect.top - tooltipRect.height - spacing
      : Math.min(rootRect.bottom + spacing, viewportHeight - tooltipRect.height - viewportPadding)
    tooltipPosition.left = left
    tooltipPosition.top = top
  }
  const handleViewportChange = () => {
    if (!armed.value || !props.tooltip) return
    updateTooltipPosition()
  }
  const reset = () => {
    clearTimer()
    armed.value = false
    progress.value = 0
  }
  const start = () => {
    reset()
    armed.value = true
    progress.value = 100
    nextTick(() => {
      updateTooltipPosition()
    })
    const timeout = Math.max(250, props.timeoutMs)
    const startedAt = Date.now()
    timer = window.setInterval(() => {
      const elapsed = Date.now() - startedAt
      const remaining = Math.max(0, timeout - elapsed)
      progress.value = (remaining / timeout) * 100
      if (remaining <= 0) {
        reset()
      }
    }, 50)
  }
  const handleClick = () => {
    if (props.disabled) return
    if (!armed.value) {
      start()
      return
    }
    emit('confirm')
    reset()
  }
  onUnmounted(() => {
    clearTimer()
  })
  watch(
    () => armed.value,
    isArmed => {
      if (!import.meta.client || !isArmed || !props.tooltip) return
      nextTick(() => {
        updateTooltipPosition()
      })
    }
  )
  onMounted(() => {
    if (!import.meta.client) return
    window.addEventListener('resize', handleViewportChange, { passive: true })
    window.addEventListener('scroll', handleViewportChange, { passive: true, capture: true })
  })
  onBeforeUnmount(() => {
    if (!import.meta.client) return
    window.removeEventListener('resize', handleViewportChange)
    window.removeEventListener('scroll', handleViewportChange, true)
  })
  defineExpose({ reset })
</script>
