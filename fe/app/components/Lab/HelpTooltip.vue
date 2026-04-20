<script setup lang="ts">
defineOptions({
  inheritAttrs: false
})
const attrs = useAttrs()
const props = withDefaults(
  defineProps<{
    text: string
    side?: 'left' | 'right'
    align?: 'left' | 'right'
    maxWidthClass?: string
    minWidthClass?: string
  }>(),
  {
    side: 'right',
    maxWidthClass: 'max-w-xs',
    minWidthClass: 'min-w-52'
  }
)
const open = ref(false)
const triggerRef = ref<HTMLElement | null>(null)
const tooltipPos = reactive({ left: 0, top: 0 })
const resolvedSide = computed(() => props.side || props.align || 'right')
const tooltipPositionClass = computed(() =>
  resolvedSide.value === 'left' ? '-translate-x-full origin-top-right' : 'origin-top-left'
)
const updateTooltipPosition = () => {
  if (!triggerRef.value) return
  const rect = triggerRef.value.getBoundingClientRect()
  tooltipPos.top = rect.bottom + 8
  tooltipPos.left = resolvedSide.value === 'left' ? rect.right : rect.left
}
const onViewportChange = () => {
  if (!open.value) return
  updateTooltipPosition()
}
const toggle = () => {
  open.value = !open.value
  if (open.value) updateTooltipPosition()
}
const close = () => {
  open.value = false
}
watch(open, isOpen => {
  if (!import.meta.client) return
  if (!isOpen) return
  updateTooltipPosition()
})
onMounted(() => {
  if (!import.meta.client) return
  window.addEventListener('resize', onViewportChange, { passive: true })
  window.addEventListener('scroll', onViewportChange, { passive: true, capture: true })
})
onBeforeUnmount(() => {
  if (!import.meta.client) return
  window.removeEventListener('resize', onViewportChange)
  window.removeEventListener('scroll', onViewportChange, true)
})
</script>
<template>
  <span ref="triggerRef" class="relative inline-flex items-center" v-bind="attrs" @mouseleave="close">
    <LabBaseButton
      variant="default"
      :aria-expanded="open ? 'true' : 'false'"
      aria-label="Показать подсказку"
      label="?"
      icon-only
      size="xs"
      @mouseenter="open = true"
      @focus="open = true"
      @blur="close"
      @click="toggle"
    />
  </span>
  <Teleport to="body">
    <span
      v-if="open"
      class="pointer-events-none fixed z-9999 rounded-md border px-2 py-1.5 text-xs leading-snug text-zinc-200 transition"
      :class="[tooltipPositionClass, maxWidthClass, minWidthClass, 'translate-y-0 opacity-100']"
      :style="{ left: `${tooltipPos.left}px`, top: `${tooltipPos.top}px` }"
    >
      {{ text }}
    </span>
  </Teleport>
</template>
