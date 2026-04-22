<script setup lang="ts">
import type {HTMLAttributes, Component} from "vue"
import Tooltip from "~/components/Tooltip.vue"
import {NuxtLink} from "#components"
interface Props {
  ariaExpanded?: boolean
  type?: "button" | "submit" | "reset"
  label?: string
  loading?: boolean
  disabled?: boolean
  icon?: Component
  iconClass?: HTMLAttributes["class"]
  iconPosition?: "left" | "right"
  iconTooltip?: string[]
  tooltipIndex?: number
  to?: string
}

const props = withDefaults(defineProps<Props>(), {
  ariaExpanded: false,
  type: "button",
  loading: false,
  disabled: false,
  iconPosition: "left",
  iconTooltip: () => []
})

const slots = useSlots()
const route = useRoute()

const isActive = computed(() => route.path === props.to)
const hasContent = computed(() => !!(props.label || slots.default))
const tooltip = computed(() => props.iconTooltip?.[props.tooltipIndex ?? +props.ariaExpanded])
const hasTooltip = computed(() => !!tooltip.value)

const isExternal = computed(() => /^(https?:|mailto:|tel:)/.test(props.to || ""))

const tag = computed(() => {
  if (!props.to) return "button"
  return isExternal.value ? "a" : NuxtLink
})

const isDisabled = computed(() => props.disabled || props.loading)

const buttonClass = computed(() => [
  "cursor-pointer rounded-md text-sm hover:bg-(--elevated)",
  hasContent.value ? "px-2 py-1" : "p-1",
  hasContent.value && props.icon && "inline-flex items-center justify-center gap-1",
  props.icon && props.iconPosition === "right" && "flex-row-reverse",
  isActive.value && "bg-(--footed)"
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
    :class="buttonClass"
    @click="preventIfDisabled"
  >
    <template v-if="loading">
      <slot name="loader">...</slot>
    </template>
    <template v-else>
      <template v-if="icon">
        <Tooltip v-if="hasTooltip" :text="tooltip">
          <component :is="icon" aria-hidden="true" :class="iconClass" />
        </Tooltip>
        <component v-else :is="icon" aria-hidden="true" :class="iconClass" />
      </template>
      <slot v-if="hasContent">{{ label }}</slot>
    </template>
  </component>
</template>
