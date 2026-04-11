<template>
  <component
    v-if="resolvedText"
    :is="as"
    :class="resolvedClassWithVisibility"
    :aria-live="props.tone === 'error' ? 'assertive' : 'polite'">
    {{ resolvedText }}
  </component>
</template>
<script setup lang="ts">
  import { computed, onBeforeUnmount, ref, watch } from 'vue'
  const props = withDefaults(
    defineProps<{
      text?: string | null
      tone?: NotifyTone
      size?: NotifySize
      as?: NotifyTag
      className?: string
      temporary?: boolean
      durationMs?: number
    }>(),
    {
      text: '',
      tone: 'info',
      size: 'sm',
      as: 'p',
      className: '',
      durationMs: 5000
    }
  )
  const toneClass: Record<NotifyTone, string> = {
    error:
      'lab-text-danger bg-[color-mix(in_srgb,var(--lab-danger)_12%,transparent)] ring-1 ring-inset ring-[color-mix(in_srgb,var(--lab-danger)_22%,transparent)]',
    success:
      'lab-text-success bg-[color-mix(in_srgb,var(--lab-success)_12%,transparent)] ring-1 ring-inset ring-[color-mix(in_srgb,var(--lab-success)_22%,transparent)]',
    info:
      'lab-text-secondary bg-[color-mix(in_srgb,var(--lab-text-secondary)_10%,transparent)] ring-1 ring-inset ring-[color-mix(in_srgb,var(--lab-text-secondary)_18%,transparent)]',
    warning:
      'lab-text-warning bg-[color-mix(in_srgb,var(--lab-warning)_12%,transparent)] ring-1 ring-inset ring-[color-mix(in_srgb,var(--lab-warning)_22%,transparent)]'
  }
  const sizeClass: Record<NotifySize, string> = {
    xs: 'min-h-4 px-2.5 py-1.5 text-xs leading-4',
    sm: 'min-h-5 px-3 py-2 text-sm leading-5',
    base: 'min-h-6 px-3.5 py-2.5 text-base leading-6'
  }
  const resolvedText = computed(() => String(props.text || '').trim())
  const visible = ref(false)
  let hideTimer: ReturnType<typeof setTimeout> | null = null
  const shouldAutoHide = computed(() => {
    if (typeof props.temporary === 'boolean') return props.temporary
    if (props.tone === 'error') return false
    return props.tone === 'success'
  })
  const clearHideTimer = () => {
    if (!hideTimer) return
    clearTimeout(hideTimer)
    hideTimer = null
  }
  watch(
    () => [resolvedText.value, shouldAutoHide.value, props.durationMs] as const,
    ([text, autoHide, duration]) => {
      clearHideTimer()
      if (!text) {
        visible.value = false
        return
      }
      visible.value = true
      if (!autoHide) return
      hideTimer = setTimeout(
        () => {
          visible.value = false
          hideTimer = null
        },
        Math.max(0, Number(duration) || 0)
      )
    },
    { immediate: true }
  )
  onBeforeUnmount(() => {
    clearHideTimer()
  })
  const resolvedClassWithVisibility = computed(() => [
    sizeClass[props.size],
    toneClass[props.tone],
    props.className,
    'block max-w-full wrap-break-word transition-opacity duration-200',
    visible.value ? 'visible opacity-100' : 'invisible opacity-0'
  ])
</script>
