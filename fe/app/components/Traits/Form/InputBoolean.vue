<template>
  <div class="flex min-h-11 items-center bg-transparent p-0">
    <LabBaseSwitch
      v-model="model"
      id="TraitInputBoolean"
      label="Булево"
      true-label="да"
      false-label="нет"
      tone="emerald" />
  </div>
</template>
<script setup lang="ts">
  defineProps<{ meta?: unknown }>()
  const model = defineModel<boolean>({ required: true })
</script>
