<script setup lang="ts">
interface Props {
  id?: string
  name?: string
  accept?: string
  label?: string
  buttonLabel?: string
  disabled?: boolean
  multiple?: boolean
}
const props = withDefaults(defineProps<Props>(), {
  label: "Файл",
  buttonLabel: "Выбрать файл",
  disabled: false,
  multiple: false
})
const injectedId = inject<string | undefined>("field-id", undefined)
const id = props.id || injectedId || `file-${useId()}`
const inputRef = ref<HTMLInputElement | null>(null)
const emit = defineEmits<{
  change: [event: Event, files: File[]]
}>()
const open = () => !props.disabled && inputRef.value?.click()
const onChange = (e: Event) =>
  emit("change", e, Array.from((e.target as HTMLInputElement).files ?? []))
</script>
<template>
  <UiField :label="label" :for-id="id">
    <input
      :id="id"
      ref="inputRef"
      :name="name"
      :accept="accept"
      :multiple="multiple"
      :disabled="disabled"
      type="file"
      class="sr-only"
      @change="onChange"
    />
    <UiButton type="button" :disabled="disabled" @click="open">
      <slot>{{ buttonLabel }}</slot>
    </UiButton>
  </UiField>
</template>
