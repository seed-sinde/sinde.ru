<template>
  {{ props.formula }}
</template>
<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    formula: string
  }>(),
  {
    formula: ''
  }
)
</script>
