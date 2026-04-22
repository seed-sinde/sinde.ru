<script setup lang="ts">
interface Props {
  text?: string
}
const props = withDefaults(defineProps<Props>(), {text: ""})
</script>

<template>
  <p v-if="text" role="alert" class="text-xs text-red-500">
    {{ text }}
  </p>
</template>
