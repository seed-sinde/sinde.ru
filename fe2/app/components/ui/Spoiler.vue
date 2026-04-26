<template>
  <UiButton
    variant="ghost"
    :aria-expanded="isOpen"
    :aria-controls="contentId"
    @click="isOpen = !isOpen"
  >
    <template v-if="isIconPositionLeft" #left>
      <UiTooltip :text="isOpen ? 'Закрыть' : 'Открыть'">
        <IcBaselineArrowRight
          :class="['h-5 w-5 shrink-0 transition-transform', isOpen && 'rotate-90']"
        />
      </UiTooltip>
    </template>
    <template v-else #right>
      <UiTooltip :text="isOpen ? 'Закрыть' : 'Открыть'">
        <IcBaselineArrowRight
          :class="['h-5 w-5 shrink-0 transition-transform', isOpen && 'rotate-90']"
        />
      </UiTooltip>
    </template>
    {{ label }}
  </UiButton>
  <div v-show="isOpen" :id="contentId" class="mt-4 rounded-xl bg-(--elevated) p-4 leading-none">
    <slot />
  </div>
</template>
<script setup lang="ts">
import IcBaselineArrowRight from "~icons/ic/baseline-arrow-right"
interface Props {
  label?: string
  modelValue?: boolean
  defaultExpanded?: boolean
  iconPosition?: "left" | "right"
}
const props = withDefaults(defineProps<Props>(), {
  label: "Спойлер",
  defaultExpanded: false,
  iconPosition: "left"
})
const emit = defineEmits<{
  "update:modelValue": [value: boolean]
  toggle: [value: boolean]
}>()
const triggerLayout = useLayoutUpdateTrigger()
const contentId = `spoiler-${useId()}`
const isIconPositionLeft = computed(() => props.iconPosition === "left")
const isOpen = computed({
  get: () => props.modelValue ?? internal.value,
  set: v => {
    internal.value = v
    emit("update:modelValue", v)
    emit("toggle", v)
  }
})
const internal = ref(props.defaultExpanded)
watch(isOpen, async () => {
  await nextTick()
  triggerLayout()
})
</script>
