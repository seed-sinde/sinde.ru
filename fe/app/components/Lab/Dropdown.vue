<script setup lang="ts">
import type { FloatingPanelSide } from '~/composables/useFloatingPanelPosition'
const props = withDefaults(
  defineProps<{
    modelValue?: boolean
    side?: FloatingPanelSide
    widthClass?: string
    panelClass?: string
    closeOnSelect?: boolean
    disabled?: boolean
    offset?: number
    crossAxisOffset?: number
    viewportPadding?: number
    matchTriggerWidth?: boolean
  }>(),
  {
    side: 'bottom',
    widthClass: 'w-56',
    panelClass: '',
    closeOnSelect: true,
    disabled: false,
    offset: 8,
    crossAxisOffset: 10,
    viewportPadding: 12,
    matchTriggerWidth: false
  }
)
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  open: []
  close: []
}>()
const triggerRef = ref<HTMLElement | null>(null)
const panelRef = ref<HTMLElement | null>(null)
const localOpen = ref(false)
const { panelStyle, resolvedSide,  updatePosition, schedulePositionUpdate, resetPosition } =
  useFloatingPanelPosition({
    triggerRef,
    panelRef,
    side: computed(() => props.side),
    offset: computed(() => props.offset),
    viewportPadding: computed(() => props.viewportPadding),
    matchTriggerWidth: computed(() => props.matchTriggerWidth)
  })
const isControlled = computed(() => props.modelValue !== undefined)
const isOpen = computed(() => (isControlled.value ? Boolean(props.modelValue) : localOpen.value))
const setOpen = (value: boolean) => {
  if (props.disabled) return
  if (!isControlled.value) {
    localOpen.value = value
  }
  emit('update:modelValue', value)
  if (value) emit('open')
  else emit('close')
}
const toggle = () => {
  setOpen(!isOpen.value)
}
const close = () => {
  if (!isOpen.value) return
  setOpen(false)
}
const onDocumentPointerDown = (event: Event) => {
  if (!isOpen.value) return
  const target = event.target as Node | null
  if (!target) return
  const triggerEl = triggerRef.value
  const panelEl = panelRef.value
  if (triggerEl?.contains(target) || panelEl?.contains(target)) {
    return
  }
  close()
}
const onDocumentKeydown = (event: KeyboardEvent) => {
  if (!isOpen.value) return
  if (event.key === 'Escape') {
    event.preventDefault()
    close()
    nextTick(() => {
      triggerRef.value?.focus()
    })
  }
}
const onWindowChange = () => {
  schedulePositionUpdate()
}
const handlePanelClick = () => {
  if (props.closeOnSelect) {
    close()
  }
}
watch(
  () => props.modelValue,
  value => {
    if (!isControlled.value) return
    localOpen.value = Boolean(value)
  },
  { immediate: true }
)
watch(isOpen, async value => {
  if (!import.meta.client) return
  if (value) {
    document.addEventListener('pointerdown', onDocumentPointerDown, true)
    document.addEventListener('keydown', onDocumentKeydown)
    window.addEventListener('resize', onWindowChange)
    window.addEventListener('scroll', onWindowChange, true)
    await updatePosition()
  } else {
    document.removeEventListener('pointerdown', onDocumentPointerDown, true)
    document.removeEventListener('keydown', onDocumentKeydown)
    window.removeEventListener('resize', onWindowChange)
    window.removeEventListener('scroll', onWindowChange, true)
    resetPosition()
  }
})
watch(
  () => [props.side, props.offset, props.crossAxisOffset, props.viewportPadding, props.matchTriggerWidth],
  () => {
    if (!isOpen.value) return
    schedulePositionUpdate()
  }
)
onBeforeUnmount(() => {
  if (!import.meta.client) return
  document.removeEventListener('pointerdown', onDocumentPointerDown, true)
  document.removeEventListener('keydown', onDocumentKeydown)
  window.removeEventListener('resize', onWindowChange)
  window.removeEventListener('scroll', onWindowChange, true)
})
defineExpose({
  open: () => setOpen(true),
  close,
  toggle,
  updatePosition
})
</script>
<template>
  <div class="inline-flex">
    <div ref="triggerRef" class="inline-flex" :aria-expanded="isOpen ? 'true' : 'false'" aria-haspopup="menu">
      <slot name="trigger" :open="isOpen" :toggle="toggle" :close="close" :disabled="disabled" :side="resolvedSide" />
    </div>
    <Teleport to="body">
      <transition
        enter-active-class="transition duration-150 ease-out"
        enter-from-class="translate-y-1 scale-98 opacity-0"
        enter-to-class="translate-y-0 scale-100 opacity-100"
        leave-active-class="transition duration-100 ease-in"
        leave-from-class="translate-y-0 scale-100 opacity-100"
        leave-to-class="translate-y-1 scale-98 opacity-0"
      >
        <div
          v-if="isOpen"
          ref="panelRef"
          role="menu"
          :style="panelStyle"
          :class="['lab-dropdown-panel z-50 rounded-2xl', widthClass, panelClass]"
          @click="handlePanelClick"
        >
          <slot :open="isOpen" :close="close" :toggle="toggle" :side="resolvedSide" />
        </div>
      </transition>
    </Teleport>
  </div>
</template>
