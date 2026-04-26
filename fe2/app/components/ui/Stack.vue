<script setup lang="ts">
import {useSlots, Comment, Text} from "vue"
const slots = useSlots()
</script>
<template>
  <main class="flex flex-col divide-y border-b">
    <div
      v-for="(child, i) in (slots.default?.() ?? []).filter(
        v => v.type !== Comment && v.type !== Text
      )"
      :key="child.key ?? i"
      class="p-4"
    >
      <component :is="child" />
    </div>
  </main>
</template>
