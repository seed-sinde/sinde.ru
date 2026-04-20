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
    :aria-invalid="isInvalid ? 'true' : undefined"
    class="lab-control lab-focus"
    :class="{ 'lab-control-invalid': isInvalid }"
    @input="onInput"
    @blur="event => emit('blur', event)"
    @focus="event => emit('focus', event)"
  />
</template>

<script setup lang="ts">
import { isInvalidEmailInput } from '~/utils/email'

const props = withDefaults(
  defineProps<{
    modelValue: string | number | null | undefined
    id?: string | undefined
    name?: string | undefined
    type?: string | undefined
    placeholder?: string | undefined
    disabled?: boolean | undefined
    ariaLabel?: string | undefined
    invalid?: boolean | undefined
  }>(),
  {
    modelValue: '',
    id: '',
    name: '',
    type: 'text',
    placeholder: '',
    disabled: false,
    ariaLabel: '',
    invalid: false
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'input', event: Event): void
  (e: 'blur', event: FocusEvent): void
  (e: 'focus', event: FocusEvent): void
}>()

const inputRef = ref<HTMLInputElement | null>(null)

const isInvalid = computed(() =>
  props.invalid || (props.type === 'email' && isInvalidEmailInput(String(props.modelValue ?? '')))
)

const onInput = (event: Event) => {
  const input = event.target as HTMLInputElement
  emit('update:modelValue', input.value)
  emit('input', event)
}

defineExpose({
  focus: () => inputRef.value?.focus()
})
</script>
