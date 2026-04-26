<script setup lang="ts">
interface Props {
  modelValue?: boolean
  id?: string
  name?: string
  label?: string
  disabled?: boolean
}
const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
  disabled: false
})
const emit = defineEmits<{
  "update:modelValue": [value: boolean]
  change: [event: Event]
}>()
const injectedId = inject<string | undefined>("field-id", undefined)
const id = computed(() => props.id || injectedId || `toggle-${useId()}`)
const onChange = (event: Event) => {
  emit("update:modelValue", (event.target as HTMLInputElement).checked)
  emit("change", event)
}
</script>
<template>
  <label
    :for="id"
    :class="[
      'relative inline-flex w-fit items-center gap-2 text-sm select-none',
      disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'
    ]"
  >
    <input
      :id="id"
      :name="name"
      type="checkbox"
      class="ui-focus peer absolute h-5 w-9 appearance-none rounded-full"
      :checked="modelValue"
      :disabled="disabled"
      @change="onChange"
    />
    <span
      aria-hidden="true"
      :class="[
        'relative h-4 w-9 shrink-0 rounded-full border-2 transition-colors',
        modelValue
          ? 'border-(--accent) bg-(--accent)'
          : 'border-(--border-color) bg-(--toggle-off-bg)'
      ]"
    >
      <span
        :class="[
          'absolute -top-0.5 -left-0.5 size-4 rounded-full border-2 transition-[background-color,border-color,transform]',
          modelValue
            ? 'translate-x-5 border-(--accent) bg-(--toggle-thumb)'
            : 'translate-x-0 border-(--border-color) bg-(--toggle-thumb)'
        ]"
      />
    </span>
    <span v-if="label || $slots.default" class="text-nowrap text-(--text)">
      <slot>{{ label }}</slot>
    </span>
  </label>
</template>
