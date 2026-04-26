<script setup lang="ts">
interface Props {
  text?: string
  side?: Side
}
const props = withDefaults(defineProps<Props>(), {
  text: "",
  side: "bottom"
})
const triggerRef = ref<HTMLElement | null>(null)
const panelRef = ref<HTMLElement | null>(null)
const isVisible = ref(false)
const {actualSide, panelStyle, updatePosition} = useFloatingPanelPosition({
  triggerRef,
  panelRef,
  side: () => props.side
})
const {arrowClass, arrowFullStyle, arrowColorStyle, updateArrow} = useFloatingArrow({
  triggerRef,
  panelRef,
  actualSide
})
const updateAll = () => {
  updatePosition()
  updateArrow()
}
watch(isVisible, async v => {
  if (!v) return
  await nextTick()
  requestAnimationFrame(updateAll)
})
watch(actualSide, () => {
  requestAnimationFrame(updateAll)
})
</script>
<template>
  <span
    ref="triggerRef"
    class="group relative inline-block cursor-help"
    @mouseenter="isVisible = true"
    @mouseleave="isVisible = false"
  >
    <slot />
    <div
      ref="panelRef"
      :style="panelStyle"
      class="pointer-events-none absolute z-60 w-max max-w-[min(26rem,calc(100vw-24px))] rounded-xl bg-(--tooltip-bg) p-1 px-2 text-xs font-bold wrap-break-word whitespace-normal text-white opacity-0 transition-opacity group-hover:opacity-100"
    >
      {{ text }}
      <div :class="arrowClass" :style="[arrowFullStyle, arrowColorStyle]" />
    </div>
  </span>
</template>
