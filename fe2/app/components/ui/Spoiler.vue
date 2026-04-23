<template>
  <section class="leading-none">
    <UiButton
      :label="label"
      :aria-expanded="isOpen"
      :icon="IcBaselineArrowRight"
      :icon-tooltip="['Показать', 'Скрыть']"
      :icon-class="['text-lg transition-transform', {'rotate-90': isOpen}]"
      @click="toggle"
    />
    <div v-show="isOpen">
      <slot :expanded="isOpen" />
    </div>
  </section>
</template>
<script setup lang="ts">
import IcBaselineArrowRight from "~icons/ic/baseline-arrow-right"
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
