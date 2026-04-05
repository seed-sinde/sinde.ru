<template>
  <div :class="['flex flex-col', gapClass, fieldClass]">
    <component
      :is="labelTag"
      v-if="hasLabel"
      :for="forId || undefined"
      :class="['lab-field-label', labelClass]">
      <slot name="label">{{ label }}</slot>
    </component>
    <slot></slot>
    <LabHint v-if="resolvedHint && !resolvedError" :text="resolvedHint" :hint-class="hintClass" />
    <LabErrorMessage v-if="resolvedError" :text="resolvedError" :error-class="errorClass" />
  </div>
</template>
<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      label?: string
      forId?: string
      hint?: string
      error?: string | string[] | null
      fieldClass?: string
      labelClass?: string
      hintClass?: string
      errorClass?: string
      gapClass?: string
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
      gapClass: 'gap-1'
    }
  )
  const hasLabel = computed(() => Boolean(props.label || useSlots().label))
  const labelTag = computed(() => (props.forId ? 'label' : 'div'))
  const resolvedHint = computed(() => String(props.hint || '').trim())
  const resolvedError = computed(() => {
    if (Array.isArray(props.error)) {
      return props.error.map(item => String(item || '').trim()).find(Boolean) || ''
    }
    return String(props.error || '').trim()
  })
</script>
