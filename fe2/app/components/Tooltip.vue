<script setup lang="ts">
interface Props {
  text?: string
  side?: Side
}

const props = withDefaults(defineProps<Props>(), {
  text: '',
  side: 'bottom'
})

const triggerRef = ref<HTMLElement | null>(null)
const panelRef = ref<HTMLElement | null>(null)
const open = ref(false)

const { panelStyle, updatePosition } =
  useFloatingPanelPosition({ triggerRef, panelRef, side: () => props.side })

const show = () => {
  if (!props.text) return
  open.value = true
  nextTick(updatePosition)
}

const hide = () => {
  open.value = false
}
</script>

<template>
  <span ref="triggerRef" class="inline-flex cursor-pointer" @mouseenter="show" @mouseleave="hide">
    <slot />
  </span>
  <div
    v-if="open && text"
    ref="panelRef"
    :style="panelStyle"
    class="pointer-events-none z-60 w-max max-w-[min(16rem,calc(100vw-24px))] rounded-lg bg-zinc-900 p-1 px-2 text-xs font-bold wrap-break-word whitespace-normal text-white"
  >
    {{ text }}
  </div>
</template>
