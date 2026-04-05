<template>
  <component :is="as" v-show="resolvedText" :class="resolvedClassWithVisibility">
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
    error: 'lab-text-danger',
    success: 'lab-text-success',
    info: 'lab-text-secondary',
    warning: 'lab-text-warning'
  }
  const sizeClass: Record<NotifySize, string> = {
    xs: 'min-h-4 text-xs',
    sm: 'min-h-5 text-sm',
    base: 'min-h-6 text-base'
  }
  const resolvedText = computed(() => String(props.text || '').trim())
  const visible = ref(false)
  let hideTimer: ReturnType<typeof setTimeout> | null = null
  const shouldAutoHide = computed(() => {
    if (props.tone === 'error') return false
    if (typeof props.temporary === 'boolean') return props.temporary
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
    'transition-opacity duration-200',
    visible.value ? 'visible opacity-100' : 'invisible opacity-0'
  ])
</script>
