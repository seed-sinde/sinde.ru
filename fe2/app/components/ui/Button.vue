<script setup lang="ts">
import {NuxtLink} from "#components"
type Tag = "button" | "a" | typeof NuxtLink
interface Props {
  type?: "button" | "submit" | "reset"
  loading?: boolean
  disabled?: boolean
  to?: string
  variant?: "default" | "ghost"
}
defineOptions({inheritAttrs: true})
const props = withDefaults(defineProps<Props>(), {
  type: "button",
  loading: false,
  disabled: false,
  variant: "default"
})
const slots = useSlots()
const attrs = useAttrs()
const route = useRoute()
const isHttp = computed(() => /^https?:/i.test(props.to || ""))
const isExternal = computed(() => /^(https?:\/\/|mailto:|tel:)/i.test(props.to || ""))
const isActive = computed(() => {
  if (!props.to) return false
  return route.path === props.to
})
const hasContent = computed(() => !!slots.default?.().length)
const tag = computed<Tag>(() => (!props.to ? "button" : isExternal.value ? "a" : NuxtLink))
const isLink = computed(() => tag.value !== "button")
const isDisabled = computed(() => props.disabled || props.loading)
const isDisabledLink = computed(() => isDisabled.value && isLink.value)
const tabIndex = computed(() =>
  isDisabledLink.value ? -1 : (attrs.tabindex as string | number | undefined)
)
const variantClass = computed(() => {
  if (props.variant === "ghost") {
    return isActive.value
      ? [
          "bg-transparent text-(--accent)",
          "hover:bg-[color-mix(in_oklch,var(--accent)_15%,transparent)]",
          "active:bg-[color-mix(in_oklch,var(--accent)_22%,transparent)]"
        ]
      : [
          "bg-transparent text-(--text)",
          "hover:bg-[color-mix(in_oklch,var(--elevated),var(--text)_8%)]",
          "active:bg-[color-mix(in_oklch,var(--elevated),var(--text)_13%)]"
        ]
  }

  // default
  return isActive.value
    ? [
        "bg-[color-mix(in_oklch,var(--accent)_18%,var(--bg))] text-(--accent)",
        "hover:bg-[color-mix(in_oklch,var(--accent)_25%,var(--bg))]",
        "active:bg-[color-mix(in_oklch,var(--accent)_32%,var(--bg))]"
      ]
    : [
        "bg-(--elevated) text-(--text)",
        "hover:bg-[color-mix(in_oklch,var(--elevated),var(--text)_8%)]",
        "active:bg-[color-mix(in_oklch,var(--elevated),var(--text)_13%)]"
      ]
})
const buttonClass = computed(() => [
  "ui-focus flex items-center justify-center",
  "leading-none text-sm font-medium select-none",
  hasContent.value ? "rounded-xl px-2 py-1.5 gap-1.5" : "rounded-full p-1",
  variantClass.value,
  isDisabled.value && ["opacity-50", "hover:bg-(--elevated) active:bg-(--elevated)"],
  props.loading ? "cursor-wait" : isDisabled.value && "cursor-not-allowed"
])
const preventIfDisabled = (e: Event) => {
  if (isDisabledLink.value) {
    e.preventDefault()
    e.stopPropagation()
  }
}
</script>

<template>
  <component
    :is="tag"
    v-bind="$attrs"
    :to="tag === NuxtLink ? props.to : undefined"
    :href="tag === 'a' ? props.to : undefined"
    :type="!isLink ? props.type : undefined"
    :disabled="!isLink ? isDisabled : undefined"
    :aria-disabled="isDisabledLink ? 'true' : undefined"
    :aria-busy="props.loading || undefined"
    :tabindex="tabIndex"
    :target="isHttp ? '_blank' : undefined"
    :rel="isHttp ? 'noopener noreferrer' : undefined"
    :class="buttonClass"
    @click="preventIfDisabled"
  >
    <UiLoader v-show="props.loading" />
    <span class="flex items-center gap-1.5 text-nowrap">
      <slot name="left" />
      <slot />
      <slot name="right" />
    </span>
  </component>
</template>
