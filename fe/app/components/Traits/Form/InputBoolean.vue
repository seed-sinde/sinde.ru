<template>
  <div class="flex min-h-11 items-center bg-transparent p-0">
    <LabBaseSwitch
      id="TraitInputBoolean"
      v-model="model"
      label="Булево"
      true-label="да"
      false-label="нет"
      tone="emerald"
    />
  </div>
</template>
<script setup lang="ts">
defineProps<{ meta?: unknown }>()
const model = defineModel<boolean>({ required: true })
</script>
