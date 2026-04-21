<script setup lang="ts">
import {provide} from "vue"
const props = defineProps<{
  label?: string
  forId?: string
  hint?: string
  error?: string | string[] | null
  required?: boolean
}>()
const baseId = props.forId || `field-${useId()}`
provide("field-id", baseId)
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
      <Hint v-if="hint" :text="hint" class="text-xs text-(--text) opacity-70" />
    </div>
    <slot />
    <MessageError v-if="error" :text="Array.isArray(error) ? error[0] : error" />
  </div>
</template>
