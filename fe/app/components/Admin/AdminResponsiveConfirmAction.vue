<script setup lang="ts">
const props = defineProps<{
  desktopLabel: string
  mobileAriaLabel?: string
  confirmLabel?: string
  mobileConfirmLabel?: string
  mobileConfirmAriaLabel?: string
  tooltip?: string
  icon?: string
  buttonClass?: string
  idleClass?: string
  confirmClass?: string
  progressClass?: string
  tooltipClass?: string
  disabled?: boolean
  type?: 'button' | 'submit' | 'reset'
  title?: string
  variant?: LabButtonVariant
}>()

defineEmits<{
  (e: 'confirm'): void
}>()

const desktopProps = computed(() => ({
  label: props.desktopLabel,
  confirmLabel: props.confirmLabel ?? 'Подтвердить',
  ...(props.tooltip ? { tooltip: props.tooltip } : {}),
  ...(props.icon ? { icon: props.icon } : {}),
  ...(props.buttonClass ? { buttonClass: props.buttonClass } : {}),
  ...(props.idleClass ? { idleClass: props.idleClass } : {}),
  ...(props.confirmClass ? { confirmClass: props.confirmClass } : {}),
  ...(props.progressClass ? { progressClass: props.progressClass } : {}),
  ...(props.tooltipClass ? { tooltipClass: props.tooltipClass } : {}),
  ...(props.disabled !== undefined ? { disabled: props.disabled } : {}),
  ...(props.type ? { type: props.type } : {}),
  ...(props.title ? { title: props.title } : {}),
  ...(props.variant ? { variant: props.variant } : {})
}))

const mobileProps = computed(() => ({
  iconOnly: true,
  confirmLabel: props.mobileConfirmLabel ?? 'Ок',
  ariaLabel: props.mobileAriaLabel ?? props.desktopLabel,
  confirmAriaLabel: props.mobileConfirmAriaLabel ?? props.confirmLabel ?? 'Подтвердить',
  ...(props.tooltip ? { tooltip: props.tooltip } : {}),
  ...(props.icon ? { icon: props.icon } : {}),
  ...(props.buttonClass ? { buttonClass: props.buttonClass } : {}),
  ...(props.idleClass ? { idleClass: props.idleClass } : {}),
  ...(props.confirmClass ? { confirmClass: props.confirmClass } : {}),
  ...(props.progressClass ? { progressClass: props.progressClass } : {}),
  ...(props.tooltipClass ? { tooltipClass: props.tooltipClass } : {}),
  ...(props.disabled !== undefined ? { disabled: props.disabled } : {}),
  ...(props.type ? { type: props.type } : {}),
  ...(props.title ? { title: props.title } : {}),
  ...(props.variant ? { variant: props.variant } : {})
}))
</script>

<template>
  <div>
    <span class="max-sm:hidden">
      <LabConfirmActionButton v-bind="desktopProps" @confirm="$emit('confirm')" />
    </span>

    <span class="sm:hidden">
      <LabConfirmActionButton v-bind="mobileProps" @confirm="$emit('confirm')" />
    </span>
  </div>
</template>
