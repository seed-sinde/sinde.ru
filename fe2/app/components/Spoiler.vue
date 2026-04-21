<template>
  <section class="leading-none">
    <Button
      :aria-expanded="isOpen"
      icon="ic:baseline-arrow-right"
      :icon-tooltip="['Показать', 'Скрыть']"
      :icon-class="['text-xl transition-transform origin-center inline-block', {'rotate-90': isOpen}]"
      @click="toggle"
    >
      {{ label }}
    </Button>
    <div v-show="isOpen">
      <slot :expanded="isOpen" />
    </div>
  </section>
</template>
<script setup lang="ts">
interface Props {
  label?: string
  modelValue?: boolean
  defaultExpanded?: boolean
}
const props = withDefaults(defineProps<Props>(), {
  label: "Спойлер",
  defaultExpanded: false
})

const emit = defineEmits<{
  "update:modelValue": [value: boolean]
  toggle: [value: boolean]
}>()

const internalOpen = ref(props.modelValue ?? props.defaultExpanded)
watch(
  () => props.modelValue,
  v => {
    if (typeof v === "boolean") internalOpen.value = v
  },
  {immediate: true}
)

const isOpen = computed(() => internalOpen.value)

const setOpen = (value: boolean) => {
  internalOpen.value = value
  emit("update:modelValue", value)
  emit("toggle", value)
}
const toggle = () => setOpen(!isOpen.value)
</script>
