<template>
  <input
    :id="id"
    ref="inputRef"
    :name="name"
    :type="type"
    :value="modelValue ?? ''"
    :disabled="disabled"
    :placeholder="placeholder"
    :aria-label="ariaLabel"
    :aria-invalid="invalid || undefined"
    class="bg-(--elevated) p-2 px-3 focus:ring focus:ring-(--accent) focus:outline-none"
    @input="onInput"
    @blur="emit('blur', $event)"
    @focus="emit('focus', $event)"
  />
</template>

<script setup lang="ts">
interface Props {
  id?: string
  type?: HTMLInputElement["type"]
  name?: string
  invalid?: boolean
  disabled?: boolean
  ariaLabel?: string
  modelValue?: string | number | null
  placeholder?: string
}
const injectedId = inject<string | undefined>("field-id", undefined)
const props = withDefaults(defineProps<Props>(), {
  type: "text",
  disabled: false,
  invalid: false
})
const id = props.id || injectedId || `input-${useId()}`
const emit = defineEmits<{
  (e: "update:modelValue", value: string): void
  (e: "input", event: Event): void
  (e: "blur" | "focus", event: FocusEvent): void
}>()
const inputRef = ref<HTMLInputElement>()
const onInput = (e: Event) => {
  const v = (e.target as HTMLInputElement).value
  emit("update:modelValue", v)
  emit("input", e)
}
defineExpose({
  focus: () => inputRef.value?.focus()
})
</script>
