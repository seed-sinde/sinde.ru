<template>
  <p :class="['text-xs text-(--lab-text-soft)', hintClass]">
    <slot>{{ text }}</slot>
  </p>
</template>
<script setup lang="ts">
withDefaults(defineProps<{ text?: string; hintClass?: string }>(), {
  text: '',
  hintClass: ''
})
</script>
