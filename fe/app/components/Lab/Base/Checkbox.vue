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
  <label v-else :for="id" :class="['lab-text-secondary inline-flex items-center gap-2 text-sm', wrapperClass]">
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
    <span
      ><slot>{{ label }}</slot></span
    >
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
    'lab-focus h-4 w-4 border bg-(--lab-bg-control) text-(--lab-accent) accent-(--lab-accent) transition hover:enabled:border-(--lab-border-strong) hover:enabled:bg-(--lab-bg-control-hover) disabled:cursor-not-allowed disabled:opacity-60',
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
