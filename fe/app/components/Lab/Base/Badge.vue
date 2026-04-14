<template>
  <span :class="badgeClass">
    <slot />
  </span>
</template>
<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    variant?: 'default' | 'muted' | 'success' | 'warning' | 'danger' | 'info'
    size?: 'xs' | 'sm' | 'md'
    rounded?: boolean
  }>(),
  {
    variant: 'default',
    size: 'xs',
    rounded: true
  }
)
const baseClass = 'lab-badge'
const sizeClassMap = {
  xs: 'px-1.5 py-0.5 text-xs',
  sm: 'px-2 py-0.5 text-sm',
  md: 'px-2.5 py-1 text-sm'
}
const variantClassMap = {
  default: 'lab-badge-default',
  muted: 'lab-badge-muted',
  success: 'lab-badge-success',
  warning: 'lab-badge-warning',
  danger: 'lab-badge-danger',
  info: 'lab-badge-info'
}
const roundedClass = computed(() => (props.rounded ? 'rounded-full' : 'rounded'))
const badgeClass = computed(() => [
  baseClass,
  sizeClassMap[props.size],
  variantClassMap[props.variant],
  roundedClass.value
])
</script>
