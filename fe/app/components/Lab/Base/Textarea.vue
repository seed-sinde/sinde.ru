<template>
  <textarea
    :id="id"
    :name="name"
    :rows="rows"
    :disabled="disabled"
    :placeholder="placeholder"
    :aria-label="ariaLabel"
    :value="resolvedValue"
    :class="textareaClassList"
    v-bind="textareaAttrs"
    @input="onInput"></textarea>
</template>
<script setup lang="ts">
  import { twMerge } from 'tailwind-merge'
  defineOptions({ inheritAttrs: false })
  const props = withDefaults(
    defineProps<{
      modelValue?: string | null | undefined
      id?: string
      name?: string
      rows?: number | string
      disabled?: boolean
      placeholder?: string
      ariaLabel?: string
      invalid?: boolean
      textareaClass?: string
    }>(),
    {
      modelValue: '',
      id: '',
      name: '',
      rows: 3,
      disabled: false,
      placeholder: '',
      ariaLabel: '',
      invalid: false,
      textareaClass: ''
    }
  )
  const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void
    (e: 'input', event: Event): void
  }>()
  const attrs = useAttrs()
  const textareaAttrs = computed(() => {
    const out: Record<string, unknown> = {}
    for (const [key, value] of Object.entries(attrs)) {
      if (key === 'class') continue
      out[key] = value
    }
    return out
  })
  const externalClass = computed(() => String(attrs.class || '').trim())
  const textareaClassList = computed(() => {
    const stateClass = props.invalid ? 'lab-control-invalid' : ''
    return twMerge('lab-control min-h-24 resize-y', stateClass, props.textareaClass, externalClass.value)
  })
  const resolvedValue = computed(() => String(props.modelValue || ''))
  const onInput = (event: Event) => {
    const target = event.target as HTMLTextAreaElement
    emit('update:modelValue', target.value)
    emit('input', event)
  }
</script>
