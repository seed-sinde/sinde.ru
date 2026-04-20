<template>
  <LabBaseField :label="label" :for-id="id" :class="disabled ? 'cursor-not-allowed' : 'cursor-pointer'">
    <input
      :id="id"
      ref="inputRef"
      :name="name"
      :accept="accept"
      :disabled="disabled"
      type="file"
      class="hidden"
      @change="$emit('change', $event)"
    />
    <LabBaseButton
      :label="label"
      :variant="variant"
      :size="size"
      :disabled="disabled"
      :icon="icon"
      :icon-only="iconOnly"
      :class="buttonClass"
      @click="openPicker"
    />
  </LabBaseField>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    id?: string
    name?: string
    accept?: string
    label?: string
    disabled?: boolean
    icon?: string
    iconOnly?: boolean
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
    icon: '',
    iconOnly: false,
    buttonClass: '',
    variant: 'secondary',
    size: 'sm'
  }
)

defineEmits<{
  (e: 'change', event: Event): void
}>()

const inputRef = ref<HTMLInputElement | null>(null)

const openPicker = () => !props.disabled && inputRef.value?.click()
</script>
