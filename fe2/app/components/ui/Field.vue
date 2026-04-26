<script setup lang="ts">
interface Props {
  label?: string
  forId?: string
  hint?: string
  error?: string | string[] | null
  required?: boolean
}
const props = withDefaults(defineProps<Props>(), {
  error: null,
  required: false
})
const baseId = props.forId || `field-${useId()}`
provide("field-id", baseId)
const error = computed(() => (Array.isArray(props.error) ? props.error[0] : props.error))
</script>

<template>
  <div class="inline-flex flex-col space-y-1">
    <div>
      <label v-if="label || $slots.label" :for="baseId" class="text-sm">
        <slot name="label">
          {{ label }}
          <span v-if="required" class="text-sm text-(--danger)">*</span>
        </slot>
      </label>
      <UiHint v-if="hint" :text="hint" class="text-xs text-(--text) opacity-70" />
    </div>
    <slot />
    <MessageError v-if="error" :text="error" />
  </div>
</template>
