<script setup lang="ts">
withDefaults(
  defineProps<{
    prevDisabled?: boolean
    nextDisabled?: boolean
  }>(),
  {
    prevDisabled: false,
    nextDisabled: false
  }
)

defineEmits<{
  (e: 'prev' | 'next'): void
}>()
</script>

<template>
  <div class="flex items-center gap-2">
    <LabBaseButton label="Назад" variant="secondary" size="xs" :disabled="prevDisabled" @click="$emit('prev')" />
    <LabBaseButton label="Вперёд" variant="secondary" size="xs" :disabled="nextDisabled" @click="$emit('next')" />
  </div>
</template>
