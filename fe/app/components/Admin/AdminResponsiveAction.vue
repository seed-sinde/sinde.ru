<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    label: string
    ariaLabel?: string
    icon?: string
    variant?: LabButtonVariant
    size?: LabButtonSize
    buttonClass?: string
    disabled?: boolean
    title?: string
    type?: 'button' | 'submit' | 'reset'
  }>(),
  {
    variant: 'secondary',
    size: 'xs',
    disabled: false,
    type: 'button'
  }
)

defineEmits<{
  (e: 'click'): void
}>()

const desktopProps = computed(() => ({
  label: props.label,
  variant: props.variant,
  size: props.size,
  disabled: props.disabled,
  type: props.type,
  ...(props.icon ? { icon: props.icon } : {}),
  ...(props.buttonClass ? { buttonClass: props.buttonClass } : {}),
  ...(props.title ? { title: props.title } : {})
}))

const mobileProps = computed(() => ({
  iconOnly: true,
  ariaLabel: props.ariaLabel ?? props.label,
  variant: props.variant,
  size: props.size,
  disabled: props.disabled,
  type: props.type,
  ...(props.icon ? { icon: props.icon } : {}),
  ...(props.buttonClass ? { buttonClass: props.buttonClass } : {}),
  ...(props.title ? { title: props.title } : {})
}))
</script>

<template>
  <div>
    <span class="max-sm:hidden">
      <LabBaseButton v-bind="desktopProps" @click="$emit('click')" />
    </span>

    <span class="sm:hidden">
      <LabBaseButton v-bind="mobileProps" @click="$emit('click')" />
    </span>
  </div>
</template>
