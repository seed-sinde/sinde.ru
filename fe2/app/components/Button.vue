<script setup lang="ts">
interface Props {
  type?: "button" | "submit" | "reset";
  label?: string;
  loading?: boolean;
  disabled?: boolean;
  variant?: string;
  size?: string;
  icon?: string;
  iconPosition?: "left" | "right";
}
const props = withDefaults(defineProps<Props>(), {
  type: "button",
  label: "",
  loading: false,
  disabled: false,
  variant: "primary",
  size: "sm",
  icon: "",
  iconPosition: "left",
});
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    class="inline-flex items-center gap-2"
    :class="icon && iconPosition === 'right' ? 'flex-row-reverse' : ''"
  >
    <template v-if="loading">
      <slot name="loader">Loading...</slot>
    </template>

    <template v-else>
      <Icon v-if="icon" :name="icon" />

      <span v-if="label || $slots.default">
        <slot>{{ label }}</slot>
      </span>
    </template>
  </button>
</template>
