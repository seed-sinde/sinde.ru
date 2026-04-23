<script setup lang="ts">
interface Props {
  text?: string
}
const props = withDefaults(defineProps<Props>(), {text: ""})
</script>

<template>
  <p v-if="props.text" class="text-xs opacity-70">
    {{ props.text }}
  </p>
</template>
