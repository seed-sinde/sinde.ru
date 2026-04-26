<script setup lang="ts">
defineOptions({inheritAttrs: false})
type Item = {
  key?: string
  src: string
  thumbnailSrc?: string | null
  alt?: string | null
}
withDefaults(
  defineProps<{
    items: Item[]
    activeIndex?: number
    label?: string
  }>(),
  {
    activeIndex: 0,
    label: "Показать изображение"
  }
)
const emit = defineEmits<{select: [index: number]}>()
const scrollerRef = ref<HTMLElement | null>(null)
const {edges} = useScrollableEdges(scrollerRef, {axis: "x"})
</script>
<template>
  <div v-if="items.length > 1" class="relative min-w-0" v-bind="$attrs">
    <div ref="scrollerRef" class="ui-thumbnails-scroll min-w-0 overflow-x-auto pb-2">
      <div class="flex min-w-max gap-2">
        <slot :items="items" :active-index="activeIndex" :select="(i: number) => emit('select', i)">
          <button
            v-for="(item, i) in items"
            :key="item.key || `${item.src}:${i}`"
            type="button"
            :aria-label="`${label} ${i + 1}`"
            :aria-current="i === activeIndex ? 'true' : undefined"
            :data-thumb-index="i"
            :class="[
              'ui-focus relative h-16 w-16 shrink-0 overflow-visible transition sm:h-18 sm:w-18',
              i === activeIndex
                ? 'opacity-100 ring-2 ring-(--accent)'
                : 'opacity-70 hover:ring-2 hover:ring-(--accent)'
            ]"
            @click="emit('select', i)"
          >
            <span
              :class="[
                'absolute inset-x-0 top-0 h-1',
                i === activeIndex ? 'bg-(--accent)' : 'bg-transparent'
              ]"
              aria-hidden="true"
            />
            <img
              :src="item.thumbnailSrc || item.src"
              :alt="item.alt || ''"
              class="h-full w-full object-cover pt-1.5"
              loading="lazy"
              decoding="async"
            />
          </button>
        </slot>
      </div>
    </div>
    <span
      :class="[
        'pointer-events-none absolute top-0 bottom-2 left-0 w-6 bg-(--bg) opacity-0 transition-opacity',
        edges.left && 'opacity-90'
      ]"
      aria-hidden="true"
    />
    <span
      :class="[
        'pointer-events-none absolute top-0 right-0 bottom-2 w-6 bg-(--bg) opacity-0 transition-opacity',
        edges.right && 'opacity-90'
      ]"
      aria-hidden="true"
    />
  </div>
</template>
<style scoped>
.ui-thumbnails-scroll {
  scrollbar-width: none;
}
.ui-thumbnails-scroll::-webkit-scrollbar {
  display: none;
}
</style>
