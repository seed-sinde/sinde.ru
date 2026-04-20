<template>
  <div :class="['flex flex-col', gapClass, fieldClass]">
    <component
      v-if="hasLabel"
      :is="labelTag"
      :for="labelTag === 'label' ? forId : undefined"
      :class="['lab-field-label', labelClass]"
    >
      <slot name="label">
        <span>{{ label }}</span>
        <span v-if="required" class="text-rose-400">*</span>
      </slot>
    </component>

    <slot />

    <LabHint v-if="resolvedHint && !resolvedError" :text="resolvedHint" :hint-class="hintClass" />
    <LabErrorMessage v-if="resolvedError" :text="resolvedError" :error-class="errorClass" />
  </div>
</template>

<script setup lang="ts">
const slots = useSlots()

const props = withDefaults(
  defineProps<{
    label?: string | undefined
    forId?: string | undefined
    hint?: string | undefined
    error?: string | string[] | null | undefined
    fieldClass?: string | undefined
    labelClass?: string | undefined
    hintClass?: string | undefined
    errorClass?: string | undefined
    gapClass?: string | undefined
    required?: boolean | undefined
  }>(),
  {
    label: '',
    forId: '',
    hint: '',
    error: '',
    fieldClass: '',
    labelClass: '',
    hintClass: '',
    errorClass: '',
    gapClass: 'gap-1',
    required: false
  }
)

const hasLabel = computed(() => Boolean(props.label || slots.label))
const labelTag = computed<'label' | 'div'>(() => (props.forId ? 'label' : 'div'))
const resolvedHint = computed(() => (props.hint ?? '').trim())
const resolvedError = computed(() =>
  Array.isArray(props.error)
    ? props.error.map(item => String(item ?? '').trim()).find(Boolean) || ''
    : String(props.error ?? '').trim()
)
</script>
