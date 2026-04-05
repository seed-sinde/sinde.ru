<script setup lang="ts">
  const model = defineModel<boolean>({ required: true })
  const props = withDefaults(
    defineProps<{
      mode?: 'button' | 'checkbox'
      label?: string
      buttonClass?: string
      textClass?: string
    }>(),
    {
      mode: 'checkbox',
      label: 'Группировка по категориям',
      buttonClass: '',
      textClass: ''
    }
  )
  const resolvedButtonClass = computed(() => [
    'inline-flex items-center gap-1 whitespace-nowrap',
    model.value ? 'kitchen-group-by-category-toggle-active' : '',
    props.buttonClass
  ])
  const resolvedTextClass = computed(() => props.textClass || 'text-xs lab-text-muted')
  const toggle = () => {
    model.value = !model.value
  }
</script>
<template>
  <LabBaseButton
    v-if="mode === 'button'"
    variant="secondary"
    size="md"
    :button-class="resolvedButtonClass"
    :icon="model ? 'ic:round-check-box' : 'ic:round-check-box-outline-blank'"
    icon-class="h-4 w-4"
    @click="toggle">
    {{ label }}
  </LabBaseButton>
  <label v-else class="inline-flex items-center gap-2" :class="resolvedTextClass">
    <LabBaseCheckbox v-model="model" bare />
    {{ label }}
  </label>
</template>
<style scoped>
  :deep(.kitchen-group-by-category-toggle-active) {
    border-color: color-mix(in srgb, var(--lab-success) 42%, var(--lab-border));
    background: color-mix(in srgb, var(--lab-success) 14%, var(--lab-bg-control));
    color: var(--lab-text-primary);
  }
  :deep(.kitchen-group-by-category-toggle-active:hover:enabled) {
    border-color: color-mix(in srgb, var(--lab-success) 56%, var(--lab-border-strong));
    background: color-mix(in srgb, var(--lab-success) 20%, var(--lab-bg-control-hover));
    color: var(--lab-text-primary);
  }
</style>
