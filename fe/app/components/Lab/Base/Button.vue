<template>
  <button
    :type="type"
    :disabled="isDisabled"
    :aria-busy="loading ? 'true' : undefined"
    :class="rootClass"
    v-bind="forwardedAttrs"
  >
    <Icon v-if="showLeftIcon" :name="resolvedIcon" :class="iconClass" />
    <slot v-if="$slots.default" />
    <span v-else-if="label" class="leading-none">
      {{ label }}
    </span>
    <Icon v-if="showRightIcon" :name="resolvedIcon" :class="iconClass" />
  </button>
</template>

<script setup lang="ts">
defineOptions({ inheritAttrs: false })

const props = withDefaults(
  defineProps<{
    type?: LabButtonType
    label?: string
    loading?: boolean
    disabled?: boolean
    variant?: LabButtonVariant
    size?: LabButtonSize
    block?: boolean
    icon?: string
    iconOnly?: boolean
    iconPosition?: 'left' | 'right'
  }>(),
  {
    type: 'button',
    label: '',
    loading: false,
    disabled: false,
    variant: 'default',
    size: 'sm',
    block: false,
    icon: '',
    iconOnly: false,
    iconPosition: 'left'
  }
)

const attrs = useAttrs()

const isDisabled = computed(() => props.disabled || props.loading)
const resolvedIcon = computed(() => props.loading ? 'ic:round-autorenew' : props.icon)
const hasIcon = computed(() => Boolean(resolvedIcon.value))
const showLeftIcon = computed(() => hasIcon.value && props.iconPosition === 'left')
const showRightIcon = computed(() => hasIcon.value && props.iconPosition === 'right')

const forwardedAttrs = computed(() =>
  Object.fromEntries(Object.entries(attrs).filter(([key]) => key !== 'class'))
)

const rootClass = computed(() => [
  'lab-button lab-focus min-w-0 select-none',
  props.block && 'w-full',
  props.iconOnly ? 'rounded-full gap-0' : 'gap-2',
  props.iconOnly ? iconOnlySizeClassMap[props.size] : sizeClassMap[props.size],
  variantClassMap[props.variant],
  attrs.class
])

const iconClass = computed(() => [
  'shrink-0 leading-none text-current',
  iconSizeClassMap[props.size],
  props.loading && 'animate-spin'
])
</script>
