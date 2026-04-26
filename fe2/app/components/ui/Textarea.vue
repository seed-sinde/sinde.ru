<script setup lang="ts">
interface Props {
  modelValue?: string | null
  id?: string
  name?: string
  rows?: number | string
  disabled?: boolean
  placeholder?: string
  ariaLabel?: string
  invalid?: boolean
}
const props = withDefaults(defineProps<Props>(), {
  modelValue: "",
  rows: 3,
  disabled: false,
  invalid: false
})
const injectedId = inject<string | undefined>("field-id", undefined)
const id = props.id || injectedId || `textarea-${useId()}`
const emit = defineEmits<{
  "update:modelValue": [value: string]
  input: [event: Event]
}>()
const value = computed(() => String(props.modelValue ?? ""))
const onInput = (e: Event) => {
  emit("update:modelValue", (e.target as HTMLTextAreaElement).value)
  emit("input", e)
}
</script>
<template>
  <textarea
    :id="id"
    :name="name"
    :rows="rows"
    :disabled="disabled"
    :placeholder="placeholder"
    :aria-label="ariaLabel"
    :aria-invalid="invalid || undefined"
    :value="value"
    class="ui-focus min-h-24 resize-y bg-(--elevated) p-2 px-3 text-sm disabled:cursor-not-allowed disabled:opacity-50"
    @input="onInput"
  />
</template>
