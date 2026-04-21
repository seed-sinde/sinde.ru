<template>
  <input
    ref="inputRef"
    :id="id"
    :name="name"
    :type="type"
    :value="modelValue ?? ''"
    :disabled="disabled"
    :placeholder="placeholder"
    :aria-label="ariaLabel"
    :aria-invalid="invalid || undefined"
    @input="onInput"
    @blur="emit('blur', $event)"
    @focus="emit('focus', $event)"
    class="rounded-md bg-(--elevated) p-2 px-3"
  />
</template>

<script setup lang="ts">
import {inject} from "vue"
const injectedId = inject<string | undefined>("field-id", undefined)
const props = withDefaults(
  defineProps<{
    modelValue?: string | number | null
    id?: string
    name?: string
    type?: string
    placeholder?: string
    disabled?: boolean
    ariaLabel?: string
    invalid?: boolean
  }>(),
  {
    type: "text"
  }
)
const id = props.id || injectedId || `input-${useId()}`
const emit = defineEmits<{
  (e: "update:modelValue", value: string): void
  (e: "input", event: Event): void
  (e: "blur", event: FocusEvent): void
  (e: "focus", event: FocusEvent): void
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
