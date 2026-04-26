<script setup lang="ts">
type SelectValue = string | number | null | undefined
type SelectOption = {
  value: SelectValue
  label: string
  disabled?: boolean
}
interface Props {
  modelValue?: SelectValue
  options?: SelectOption[]
  id?: string
  name?: string
  disabled?: boolean
  invalid?: boolean
  placeholder?: string
  ariaLabel?: string
}
const props = withDefaults(defineProps<Props>(), {
  options: () => [],
  disabled: false,
  invalid: false,
  placeholder: "Выберите значение"
})
const injectedId = inject<string | undefined>("field-id", undefined)
const id = props.id || injectedId || `select-${useId()}`
const value = computed(() =>
  props.modelValue === null || props.modelValue === undefined ? "" : String(props.modelValue)
)
const normalizedOptions = computed(() =>
  props.options.map(option => ({
    ...option,
    value: option.value === null || option.value === undefined ? "" : String(option.value)
  }))
)
const emit = defineEmits<{
  "update:modelValue": [value: string]
  change: [event: Event]
}>()
const onChange = (e: Event) => {
  emit("update:modelValue", (e.target as HTMLSelectElement).value)
  emit("change", e)
}
</script>
<template>
  <select
    :id="id"
    :name="name"
    :value="value"
    :disabled="disabled"
    :aria-label="ariaLabel"
    :aria-invalid="invalid || undefined"
    class="ui-focus bg-(--elevated) p-2 px-3 text-sm disabled:cursor-not-allowed disabled:opacity-50"
    @change="onChange"
  >
    <option v-if="placeholder" value="" disabled>
      <slot name="placeholder">{{ placeholder }}</slot>
    </option>
    <option
      v-for="option in normalizedOptions"
      :key="option.value || option.label"
      :value="option.value"
      :disabled="option.disabled"
    >
      {{ option.label }}
    </option>
  </select>
</template>
