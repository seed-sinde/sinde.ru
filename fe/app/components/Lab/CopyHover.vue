<template>
  <div class="inline-block" @mouseenter="onEnter" @mouseleave="onLeave" @click="onClick">
    <slot :display="display" :label-style="labelStyle" />
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    idle: string
    hover?: string
    done?: string
    doneTimeoutMs?: number
    hasIcon?: boolean
  }>(),
  {
    hover: 'Копировать',
    done: 'Скопировано!',
    doneTimeoutMs: 1200,
    hasIcon: true
  }
)

const emit = defineEmits<{ (e: 'click'): void }>()
const hovered = ref(false)
const copied = ref(false)
let t: number | undefined

const maxChars = Math.max(props.idle.length, props.hover.length, props.done.length)
const iconReserve = props.hasIcon ? 2 : 0 // reserve space for icon in ch
const paddingReserve = '0rem' // keep padding on button; exclude from text width

const labelStyle = computed(() => ({
  width: `calc(${maxChars + iconReserve}ch + ${paddingReserve})`,
  fontFamily: 'monospace',
  display: 'inline-block', // to make width work
  textAlign: 'left'
}))

const display = computed(() => (copied.value ? props.done : hovered.value ? props.hover : props.idle))
const onEnter = () => {
  hovered.value = true
}
const onLeave = () => {
  hovered.value = false
}
const onClick = () => {
  emit('click')
  copied.value = true
  clearTimeout(t)
  t = window.setTimeout(() => {
    copied.value = false
  }, props.doneTimeoutMs)
}
onBeforeUnmount(() => clearTimeout(t))
</script>
