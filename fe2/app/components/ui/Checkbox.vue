<script setup lang="ts">
export type CheckboxValue = boolean | "partial"
interface Props {
  modelValue?: CheckboxValue
  partial?: boolean
  id?: string
  name?: string
  label?: string
  disabled?: boolean
  invalid?: boolean
  ariaLabel?: string
}
defineOptions({inheritAttrs: false})
const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
  partial: false,
  disabled: false,
  invalid: false
})
const emit = defineEmits<{
  "update:modelValue": [value: CheckboxValue]
  change: [event: Event]
}>()
const attrs = useAttrs()
const inputRef = ref<HTMLInputElement | null>(null)
const generatedId = useId()
const injectedId = inject<string | undefined>("field-id", undefined)
const id = computed(() => props.id || injectedId || `checkbox-${generatedId}`)
const isPartial = computed(() => props.modelValue === "partial")
const isChecked = computed(() => props.modelValue === true)
const ariaChecked = computed(() => (isPartial.value ? "mixed" : isChecked.value))
const inputAttrs = computed(() => {
  const {class: _class, ...rest} = attrs
  return rest
})
const getNextValue = (): CheckboxValue =>
  props.partial
    ? props.modelValue === false
      ? true
      : props.modelValue === true
        ? "partial"
        : false
    : !isChecked.value
const onChange = (event: Event) => {
  emit("update:modelValue", getNextValue())
  emit("change", event)
}
watchEffect(() => {
  if (inputRef.value) inputRef.value.indeterminate = isPartial.value
})
</script>
<template>
  <label
    :for="id"
    :class="[
      'relative inline-flex w-fit items-center gap-2 select-none',
      disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'
    ]"
  >
    <input
      :id="id"
      ref="inputRef"
      v-bind="inputAttrs"
      type="checkbox"
      :name="name"
      :checked="isChecked"
      :disabled="disabled"
      :aria-label="ariaLabel"
      :aria-invalid="invalid || undefined"
      :aria-checked="ariaChecked"
      class="peer absolute size-4 opacity-0 outline-none"
      @change="onChange"
    />
    <span
      :class="[
        'relative flex size-4 shrink-0 items-center justify-center border-2 border-(--mantis)',
        'outline-2 outline-offset-2 outline-transparent transition-all duration-200 outline-dotted',
        'peer-focus-visible:outline-(--accent)',
        isChecked && 'bg-(--mantis)',
        $attrs.class
      ]"
      aria-hidden="true"
    >
      <span v-if="isPartial" class="h-0.5 w-2 rounded-xs bg-(--mantis)" />
    </span>
    <span v-if="label || $slots.default" class="text-(--text)">
      <slot>{{ label }}</slot>
    </span>
  </label>
</template>
