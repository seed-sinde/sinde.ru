<template>
  <section>
    <button
      type="button"
      :aria-expanded="isOpen"
      class="flex flex-nowrap items-center gap-1 rounded-xl px-2 py-1 hover:bg-(--lab-bg-control-hover)"
      @click="toggle"
    >
      <span aria-hidden="true" :class="['transform text-xs transition-transform', isOpen && 'rotate-90']">►</span>
      <span>{{ label }}</span>
    </button>
    <div v-if="isOpen">
      <slot :expanded="isOpen" />
    </div>
  </section>
</template>
<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    label: string
    modelValue?: boolean
    defaultExpanded?: boolean
  }>(),
  {
    label: 'Спойлер',
    defaultExpanded: false
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  toggle: [value: boolean]
}>()

const internalOpen = ref(props.defaultExpanded)
const isControlled = computed(() => typeof props.modelValue === 'boolean')
const isOpen = computed(() => (isControlled.value ? props.modelValue === true : internalOpen.value))
const setOpen = (value: boolean) => {
  isControlled.value || (internalOpen.value = value)
  emit('update:modelValue', value)
  emit('toggle', value)
}
const toggle = () => setOpen(!isOpen.value)
</script>
