<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import Tooltip from '~/components/Tooltip.vue'
import { NuxtLink } from '#components'

interface Props {
  ariaExpanded?: boolean
  type?: "button" | "submit" | "reset"
  label?: string
  loading?: boolean
  disabled?: boolean
  icon?: string
  iconClass?: HTMLAttributes['class']
  iconPosition?: "left" | "right"
  iconTooltip?: string[]
  to?: string
}

const props = withDefaults(defineProps<Props>(), {
  ariaExpanded: false,
  type: "button",
  label: undefined,
  loading: false,
  disabled: false,
  icon: undefined,
  iconClass: undefined,
  iconPosition: "left",
  iconTooltip: () => [],
  to: undefined
})

const slots = useSlots()
const route = useRoute()

const isActive = computed(() =>route.path === props.to)
const hasContent = computed(() => !!(props.label || slots.default))

const tooltip = computed(() => props.iconTooltip?.[+props.ariaExpanded])
const hasTooltip = computed(() => !!tooltip.value)

const isExternal = computed(() =>
  /^(https?:|mailto:|tel:)/.test(props.to || "")
)

const tag = computed(() => {
  if (!props.to) return 'button'
  return isExternal.value ? 'a' : NuxtLink
})

const isDisabled = computed(() => props.disabled || props.loading)

const buttonClass = computed(() => [
  "inline-flex items-center justify-center gap-1 cursor-pointer rounded-md px-2 py-1 text-sm hover:bg-(--elevated)",
  props.icon && props.iconPosition === "right" && "flex-row-reverse",
  isActive.value && 'bg-(--footed)'
])

const preventIfDisabled = (e: Event) => {
  if (tag.value !== "button" && isDisabled.value) {
    e.preventDefault()
    e.stopPropagation()
  }
}
</script>

<template>
  <component
    :is="tag"
    :to="!isExternal ? to : undefined"
    :href="tag === 'a' ? to : undefined"
    :type="tag === 'button' ? type : undefined"
    :disabled="tag === 'button' && isDisabled"
    :aria-disabled="tag !== 'button' && isDisabled ? 'true' : undefined"
    :tabindex="tag !== 'button' && isDisabled ? -1 : undefined"
    :target="isExternal && to?.startsWith('http') ? '_blank' : undefined"
    :rel="isExternal && to?.startsWith('http') ? 'noopener noreferrer' : undefined"
    @click="preventIfDisabled"
    :class="buttonClass"
  >
    <template v-if="loading">
      <slot name="loader">...</slot>
    </template>

    <template v-else>
      <component :is="hasTooltip ? Tooltip : 'span'" :text="tooltip">
        <Icon v-if="icon" aria-hidden="true" :name="icon" :class="['align-middle', iconClass]" />
      </component>

      <span v-if="hasContent">
        <slot>{{ label }}</slot>
      </span>
    </template>
  </component>
</template>
