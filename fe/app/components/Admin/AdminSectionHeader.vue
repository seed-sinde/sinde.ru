<template>
  <div class="flex flex-wrap items-start justify-between space-y-1">
    <h2 class="text-base font-semibold text-zinc-100">{{ title }}</h2>
    <p class="text-xs text-zinc-500">{{ description }}</p>
  </div>
</template>
<script setup lang="ts">
withDefaults(
  defineProps<{
    title: string
    description: string
  }>(),
  {
    title: '',
    description: ''
  }
)
</script>
