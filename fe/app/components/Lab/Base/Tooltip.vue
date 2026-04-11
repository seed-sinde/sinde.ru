<script setup lang="ts">
  import type { FloatingPanelAlign, FloatingPanelSide } from '~/composables/useFloatingPanelPosition'

  const props = withDefaults(
    defineProps<{
      text?: string
      align?: FloatingPanelAlign
      side?: FloatingPanelSide
      panelClass?: string
      offset?: number
      crossAxisOffset?: number
      viewportPadding?: number
    }>(),
    {
      text: '',
      align: 'left',
      side: 'top',
      panelClass: '',
      offset: 8,
      crossAxisOffset: 10,
      viewportPadding: 12
    }
  )

  const triggerRef = ref<HTMLElement | null>(null)
  const panelRef = ref<HTMLElement | null>(null)
  const open = ref(false)
  const { panelStyle, resolvedSide, resolvedAlign, updatePosition, schedulePositionUpdate, resetPosition } =
    useFloatingPanelPosition({
      triggerRef,
      panelRef,
      side: computed(() => props.side),
      align: computed(() => props.align),
      offset: computed(() => props.offset),
      crossAxisOffset: computed(() => props.crossAxisOffset),
      viewportPadding: computed(() => props.viewportPadding),
      matchTriggerWidth: false
    })

  const hasContent = computed(() => Boolean(String(props.text || '').trim()))
  const show = async () => {
    if (!hasContent.value) return
    open.value = true
    await updatePosition()
  }
  const hide = () => {
    open.value = false
  }
  const onWindowChange = () => {
    if (!open.value) return
    schedulePositionUpdate()
  }

  watch(
    () => [props.align, props.side, props.offset, props.crossAxisOffset, props.viewportPadding],
    () => {
      if (!open.value) return
      schedulePositionUpdate()
    }
  )

  watch(open, value => {
    if (!import.meta.client) return
    if (value) {
      window.addEventListener('resize', onWindowChange)
      window.addEventListener('scroll', onWindowChange, true)
      return
    }
    window.removeEventListener('resize', onWindowChange)
    window.removeEventListener('scroll', onWindowChange, true)
    resetPosition()
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
    @focusout="hide">
    <slot name="trigger" :open="open" :show="show" :hide="hide" :side="resolvedSide" :align="resolvedAlign" />
  </span>
  <Teleport to="body">
    <transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="translate-y-1 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="translate-y-1 opacity-0">
      <div
        v-if="open && hasContent"
        ref="panelRef"
        :style="panelStyle"
        :class="[
          'lab-dropdown-panel pointer-events-none z-60 max-w-64 rounded-2xl px-2.5 py-1.5 text-xs leading-5',
          panelClass
        ]">
        <slot>
          {{ text }}
        </slot>
      </div>
    </transition>
  </Teleport>
</template>
