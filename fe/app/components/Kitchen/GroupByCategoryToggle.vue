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
    model.value
      ? 'border-[color-mix(in_srgb,var(--lab-accent)_42%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-accent)_14%,var(--lab-bg-control))] text-(--lab-text-primary) hover:enabled:border-[color-mix(in_srgb,var(--lab-accent)_56%,var(--lab-border-strong))] hover:enabled:bg-[color-mix(in_srgb,var(--lab-accent)_20%,var(--lab-bg-control-hover))]'
      : '',
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
