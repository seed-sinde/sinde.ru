<script setup lang="ts">
import type { FloatingPanelSide } from '~/composables/useFloatingPanelPosition'

const props = withDefaults(
  defineProps<{
    text?: string
    triggerText?: string
    side?: FloatingPanelSide
    triggerClass?: string
    underlineTrigger?: boolean
  }>(),
  {
    text: '',
    triggerText: '',
    side: 'top',
    triggerClass: '',
    underlineTrigger: false
  }
)

const triggerRef = ref<HTMLElement | null>(null)
const panelRef = ref<HTMLElement | null>(null)
const open = ref(false)

const { panelStyle, resolvedSide, updatePosition, schedulePositionUpdate, resetPosition } = useFloatingPanelPosition({
  triggerRef,
  panelRef,
  side: () => props.side
})

const hasContent = computed(() => Boolean(props.text.trim()))
const hasTriggerText = computed(() => Boolean(props.triggerText.trim()))

const show = async () => {
  if (!hasContent.value) return
  open.value = true
  await updatePosition()
  schedulePositionUpdate()
}

const hide = () => {
  open.value = false
}

const onWindowChange = () => {
  if (open.value) schedulePositionUpdate()
}

watch(open, value => {
  if (!import.meta.client) return
  const method = value ? 'addEventListener' : 'removeEventListener'
  window[method]('resize', onWindowChange)
  window[method]('scroll', onWindowChange, true)
  value || resetPosition()
})

onBeforeUnmount(() => {
  if (!import.meta.client) return
  window.removeEventListener('resize', onWindowChange)
  window.removeEventListener('scroll', onWindowChange, true)
})
</script>

<template>
  <span
    ref="triggerRef"
    class="relative inline-flex align-top"
    @mouseenter="show"
    @mouseleave="hide"
    @focusin="show"
    @focusout="hide"
  >
    <span
      v-if="$slots.trigger || hasTriggerText"
      :class="['inline-flex items-center', underlineTrigger && 'border-b border-dotted', triggerClass]"
    >
      <slot name="trigger" :open="open" :show="show" :hide="hide" :side="resolvedSide">
        {{ triggerText }}
      </slot>
    </span>
  </span>

  <Teleport to="body">
    <transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="translate-y-1 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="translate-y-1 opacity-0"
    >
      <div
        v-if="open && hasContent"
        ref="panelRef"
        :style="panelStyle"
        class="lab-dropdown-panel pointer-events-none z-60 w-max max-w-[min(16rem,calc(100vw-24px))] px-2.5 py-1.5 text-xs leading-5 wrap-break-word whitespace-normal"
      >
        <slot>
          {{ text }}
        </slot>
      </div>
    </transition>
  </Teleport>
</template>
