<script setup lang="ts">
defineProps<{text?: string}>()
</script>

<template>
  <p v-if="text" role="alert" class="text-xs text-red-500">
    {{ text }}
  </p>
</template>
