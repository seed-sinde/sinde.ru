<script setup lang="ts">
defineProps<{text?: string}>()
</script>

<template>
  <p v-if="text" class="text-xs opacity-70">
    {{ text }}
  </p>
</template>
