<template>
  <input
    id="TraitsColorPicker"
    v-model="model"
    type="color"
    class="h-11 w-16 cursor-pointer rounded-[14px] border border-zinc-700/80 bg-zinc-950/80 p-1"
  >
</template>
<script setup lang="ts">
const model = defineModel<string>({ required: true })
</script>
