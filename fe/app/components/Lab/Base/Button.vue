<template>
  <button
    :type="type"
    :disabled="disabled"
    :aria-busy="loading ? 'true' : undefined"
    :class="buttonClass"
    :style="[attrs.style, buttonStyle]"
    v-bind="passThroughAttrs">
    <slot name="leading">
      <Icon v-if="iconName" :name="iconName" :class="iconClass" :style="iconStyle" />
    </slot>

    <slot v-if="$slots.default" />
    <span v-else-if="labelText" :class="labelClass" :style="labelStyle">
      {{ labelText }}
    </span>
  </button>
</template>

<script setup lang="ts">
  defineOptions({ inheritAttrs: false })
  const props = withDefaults(
    defineProps<{
      type?: LabButtonType
      label?: string
      loading?: boolean
      loadingLabel?: string
      disabled?: boolean
      variant?: LabButtonVariant
      size?: LabButtonSize
      block?: boolean
      icon?: string
      iconOnly?: boolean
      iconSize?: LabButtonSize
      buttonClass?: LabButtonClass
      labelClass?: LabButtonClass
      iconClass?: LabButtonClass
      buttonStyle?: LabButtonStyle
      labelStyle?: LabButtonStyle
      iconStyle?: LabButtonStyle
    }>(),
    {
      type: 'button',
      label: '',
      loading: false,
      loadingLabel: '',
      disabled: false,
      variant: 'default',
      size: 'sm',
      block: false,
      icon: '',
      iconOnly: false,
      iconSize: 'md',
      buttonClass: '',
      labelClass: '',
      iconClass: '',
      buttonStyle: undefined,
      labelStyle: undefined,
      iconStyle: undefined
    }
  )
  const attrs = useAttrs()
  const labelText = computed(() => (props.loading && props.loadingLabel ? props.loadingLabel : props.label))
  const iconName = computed(() => (props.loading ? 'ic:round-autorenew' : props.icon))
  const passThroughAttrs = computed(() =>
    Object.fromEntries(Object.entries(attrs).filter(([key]) => key !== 'class' && key !== 'style'))
  )
  const buttonClass = computed(() => [
    'lab-button lab-focus min-w-0 select-none',
    props.block && 'w-full',
    props.iconOnly ? 'gap-0' : 'gap-2',
    props.iconOnly ? iconOnlySizeClassMap[props.size] : sizeClassMap[props.size],
    variantClassMap[props.variant],
    attrs.class,
    props.buttonClass
  ])

  const labelClass = computed(() => ['leading-none', props.labelClass])

  const iconClass = computed(() => [
    'shrink-0 leading-none text-current',
    iconSizeClassMap[props.iconSize],
    props.iconClass,
    props.loading && 'animate-spin'
  ])
</script>
