<template>
  <div v-if="items.length > 1" class="flex gap-2 overflow-x-auto">
    <LabBaseButton
      v-for="(item, index) in items"
      :key="item.key"
      :button-class="[
        'p-0! relative h-14 w-14 shrink-0 overflow-hidden border transition ring-1 ring-inset',
        index === activeIndex
          ? 'border-zinc-100 ring-zinc-100 bg-zinc-100/10'
          : 'border-zinc-700 ring-zinc-800 hover:border-zinc-500 hover:ring-zinc-600'
      ]"
      :aria-label="`Показать изображение ${index + 1}`"
      @click="$emit('select', index)">
      <img
        :src="item.thumbnailSrc || item.src"
        :alt="item.alt"
        class="h-full w-full object-cover"
        loading="lazy"
        decoding="async" />
      <span class="absolute inset-x-0 bottom-0 bg-black/55 px-1 py-0.5 text-center text-xs font-medium text-white">
        {{ index + 1 }}
      </span>
    </LabBaseButton>
  </div>
</template>
<script setup lang="ts">
  defineProps<{
    items: Array<{
      key: string
      src: string
      thumbnailSrc?: string | null
      alt: string
    }>
    activeIndex: number
  }>()
  defineEmits<{
    (e: 'select', index: number): void
  }>()
</script>
