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
const {panelStyle, updatePosition} = useFloatingPanelPosition({triggerRef, panelRef, side: () => props.side})
onMounted(() => {
  nextTick(updatePosition)
})
</script>
<template>
  <span ref="triggerRef" class="group relative flex w-fit cursor-pointer">
    <slot />
    <div
      ref="panelRef"
      :style="panelStyle"
      class="pointer-events-none absolute z-60 w-max max-w-[min(16rem,calc(100vw-24px))] bg-zinc-900 p-1 px-2 text-xs font-bold wrap-break-word whitespace-normal text-white opacity-0 transition-opacity group-hover:opacity-100"
    >
      {{ text }}
    </div>
  </span>
</template>
