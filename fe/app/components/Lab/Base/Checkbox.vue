<template>
  <input
    v-if="bare"
    :id="id"
    :name="name"
    type="checkbox"
    :checked="modelValue"
    :disabled="disabled"
    :class="checkboxClassList"
    v-bind="checkboxAttrs"
    @change="onChange"
  >
  <label v-else :for="id" :class="['inline-flex items-center gap-2 text-sm lab-text-secondary', wrapperClass]">
    <input
      :id="id"
      :name="name"
      type="checkbox"
      :checked="modelValue"
      :disabled="disabled"
      :class="checkboxClassList"
      v-bind="checkboxAttrs"
      @change="onChange"
    >
    <span><slot>{{ label }}</slot></span>
  </label>
</template>
<script setup lang="ts">
defineOptions({ inheritAttrs: false })
const props = withDefaults(
  defineProps<{
    modelValue?: boolean
    id?: string
    name?: string
    label?: string
    bare?: boolean
    disabled?: boolean
    wrapperClass?: string
    checkboxClass?: string
  }>(),
  {
    modelValue: false,
    id: '',
    name: '',
    label: '',
    bare: false,
    disabled: false,
    wrapperClass: '',
    checkboxClass: ''
  }
)
const emit = defineEmits<{ (e: 'update:modelValue', value: boolean): void; (e: 'change', event: Event): void }>()
const attrs = useAttrs()
const checkboxAttrs = computed(() => {
  const out: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(attrs)) {
    if (key === 'class') continue
    out[key] = value
  }
  return out
})
const externalClass = computed(() => String(attrs.class || '').trim())
const checkboxClassList = computed(() => {
  return [
    'lab-checkbox h-4 w-4 rounded accent-amber-400',
    props.checkboxClass,
    externalClass.value
  ]
})
const onChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.checked)
  emit('change', event)
}
</script>
<style scoped>
  .lab-checkbox {
    border: 1px solid var(--lab-border);
    background: var(--lab-bg-control);
    color: var(--lab-accent);
  }
  .lab-checkbox:hover:enabled {
    border-color: var(--lab-border-strong);
    background: var(--lab-bg-control-hover);
  }
  .lab-checkbox:focus,
  .lab-checkbox:focus-visible {
    border-color: var(--lab-accent);
    outline: none;
  }
  .lab-checkbox:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }
</style>
