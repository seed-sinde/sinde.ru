<template>
  <button
    :type="type"
    :disabled="disabled"
    :aria-busy="loading ? 'true' : undefined"
    :class="buttonClassList"
    :style="buttonStyleList"
    v-bind="buttonAttrs">
    <slot name="leading">
      <Icon v-if="iconName" :key="iconName" :name="iconName" :class="iconClassList" :style="iconStyleList" />
    </slot>
    <slot v-if="hasDefaultSlot" />
    <span v-else-if="resolvedLabel" :class="labelClassList" :style="labelStyleList">
      {{ resolvedLabel }}
    </span>
  </button>
</template>
<script setup lang="ts">
  import { normalizeClass, normalizeStyle, type StyleValue } from 'vue'
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
      class?: LabButtonClass
      buttonClass?: LabButtonClass
      labelClass?: LabButtonClass
      iconClass?: LabButtonClass
      buttonStyle?: LabButtonStyle
      labelStyle?: LabButtonStyle
      iconStyle?: LabButtonStyle
    }>(),
    {
      label: '',
      icon: '',
      loading: false,
      loadingLabel: '',
      disabled: false,
      block: false,
      iconOnly: false,
      variant: 'default',
      size: 'sm',
      iconSize: 'md',
      type: 'button',
      class: '',
      buttonClass: '',
      labelClass: '',
      iconClass: '',
      buttonStyle: undefined,
      labelStyle: undefined,
      iconStyle: undefined
    }
  )
  const attrs = useAttrs()
  const slots = useSlots()
  const hasDefaultSlot = computed(() => Boolean(slots.default))
  const resolvedLabel = computed(() => (props.loading && props.loadingLabel ? props.loadingLabel : props.label))
  const iconName = computed(() => (props.loading ? 'ic:round-autorenew' : props.icon || ''))
  const buttonAttrs = computed(() => {
    const out: Record<string, unknown> = {}
    for (const [key, value] of Object.entries(attrs)) {
      if (key === 'class' || key === 'style') continue
      out[key] = value
    }
    return out
  })
  const externalClass = computed(() => normalizeClass([attrs.class, props.class, props.buttonClass]))
  const externalStyle = computed(() => normalizeStyle([attrs.style as StyleValue, props.buttonStyle]))
  const buttonClassList = computed(() => [
    'inline-flex shrink-0 items-center justify-center',
    props.block ? 'w-full' : '',
    props.iconOnly ? 'gap-0 rounded-full' : 'gap-2 rounded-none',
    'min-w-0 select-none',
    props.iconOnly ? iconOnlySizeClassMap[props.size] : sizeClassMap[props.size],
    variantClassMap[props.variant],
    externalClass.value
  ])
  const buttonStyleList = computed(() => externalStyle.value)
  const labelClassList = computed(() => ['leading-none', props.labelClass])
  const labelStyleList = computed(() => normalizeStyle(props.labelStyle))
  const iconClassList = computed(() => [
    'shrink-0 leading-none text-current',
    iconSizeClassMap[props.iconSize],
    props.iconClass,
    props.loading ? 'animate-spin' : ''
  ])
  const iconStyleList = computed(() => normalizeStyle(props.iconStyle))
</script>
