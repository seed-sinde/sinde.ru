<template>
  <p :class="['lab-text-soft text-xs', hintClass]">
    <slot>{{ text }}</slot>
  </p>
</template>
<script setup lang="ts">
withDefaults(defineProps<{ text?: string; hintClass?: string }>(), {
  text: '',
  hintClass: ''
})
</script>
