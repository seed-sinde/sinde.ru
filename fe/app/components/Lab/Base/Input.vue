<template>
  <input
    :id="id"
    ref="inputRef"
    :name="name"
    :type="type"
    :value="resolvedValue"
    :disabled="disabled"
    :placeholder="placeholder"
    :aria-label="ariaLabel"
    :aria-invalid="resolvedInvalid ? 'true' : undefined"
    :class="inputClassList"
    v-bind="inputAttrs"
    @input="onInput"
    @blur="(event) => emit('blur', event)"
    @focus="(event) => emit('focus', event)"
  >
</template>
<script setup lang="ts">
import { twMerge } from 'tailwind-merge'
import { isInvalidEmailInput } from '~/utils/email'
defineOptions({ inheritAttrs: false })
const normalizeClassValue = (value: unknown): string => {
  if (!value) return ''
  if (typeof value === 'string') return value.trim()
  if (Array.isArray(value)) {
    return value
      .map((item) => normalizeClassValue(item))
      .filter(Boolean)
      .join(' ')
      .trim()
  }
  if (typeof value === 'object') {
    return Object.entries(value as Record<string, unknown>)
      .filter(([, enabled]) => Boolean(enabled))
      .map(([className]) => className)
      .join(' ')
      .trim()
  }
  return ''
}
const props = withDefaults(
  defineProps<{
    modelValue?: string | number | null | undefined
    id?: string
    name?: string
    type?: string
    placeholder?: string
    disabled?: boolean
    ariaLabel?: string
    invalid?: boolean
    inputClass?: ClassValue
  }>(),
  {
    modelValue: '',
    id: '',
    name: '',
    type: 'text',
    placeholder: '',
    disabled: false,
    ariaLabel: '',
    invalid: false,
    inputClass: ''
  }
)
const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'input', event: Event): void
  (e: 'blur' | 'focus', event: FocusEvent): void
}>()
const attrs = useAttrs()
const inputRef = ref<HTMLInputElement | null>(null)
const inputAttrs = computed(() => {
  const out: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(attrs)) {
    if (key === 'class') continue
    out[key] = value
  }
  return out
})
const externalClass = computed(() => normalizeClassValue(attrs.class))
const hasEmailValidationError = computed(() => props.type === 'email' && isInvalidEmailInput(resolvedValue.value))
const resolvedInvalid = computed(() => props.invalid || hasEmailValidationError.value)
const inputClassList = computed(() => {
  const stateClass = resolvedInvalid.value ? 'lab-control-invalid' : ''
  const overrideClass = normalizeClassValue(props.inputClass)
  return twMerge('lab-control lab-focus', stateClass, externalClass.value, overrideClass)
})
const resolvedValue = computed(() => {
  if (props.modelValue === null || props.modelValue === undefined) return ''
  return String(props.modelValue)
})
const onInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.value)
  emit('input', event)
}
const focus = () => {
  inputRef.value?.focus()
}
defineExpose({ focus })
</script>
