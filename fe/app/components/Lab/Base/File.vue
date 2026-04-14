<template>
  <label class="inline-block" :class="disabled ? 'cursor-not-allowed' : 'cursor-pointer'">
    <LabBaseInput
      :id="id"
      :name="name"
      type="file"
      :accept="accept"
      class="hidden"
      :disabled="disabled"
      @change="$emit('change', $event)"
    />
    <LabBaseButton tag="span" :variant="variant" :size="size" :disabled="disabled" :button-class="resolvedButtonClass">
      <slot>{{ label }}</slot>
    </LabBaseButton>
  </label>
</template>
<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    id?: string
    name?: string
    accept?: string
    label?: string
    disabled?: boolean
    buttonClass?: LabButtonClass
    variant?: LabButtonVariant
    size?: LabButtonSize
  }>(),
  {
    id: '',
    name: '',
    accept: '',
    label: 'Выбрать файл',
    disabled: false,
    buttonClass: '',
    variant: 'secondary',
    size: 'sm'
  }
)
defineEmits<{
  (e: 'change', event: Event): void
}>()
const resolvedButtonClass = computed(() => ['pointer-events-none', props.buttonClass].filter(Boolean).join(' '))
</script>
