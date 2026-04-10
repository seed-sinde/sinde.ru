<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      items: AstronomyCardStripItem[]
      selectedId?: string | null
      compact?: boolean
      largeSymbol?: boolean
      wide?: boolean
    }>(),
    {
      selectedId: null,
      compact: false,
      largeSymbol: false,
      wide: false
    }
  )
  const emit = defineEmits<{
    (e: 'select', id: string): void
  }>()
  const colorHexToRgba = (hexColor: string, alpha: number) => {
    const normalized = String(hexColor || '').replace('#', '')
    if (!/^[0-9a-f]{6}$/i.test(normalized)) {
      return `rgba(82, 82, 91, ${alpha})`
    }
    const red = Number.parseInt(normalized.slice(0, 2), 16)
    const green = Number.parseInt(normalized.slice(2, 4), 16)
    const blue = Number.parseInt(normalized.slice(4, 6), 16)
    return `rgba(${red}, ${green}, ${blue}, ${alpha})`
  }
  const itemStyle = (color: string) => ({
    borderColor: colorHexToRgba(color, 0.4),
    backgroundImage: `linear-gradient(160deg, ${colorHexToRgba(color, 0.24)} 0%, rgba(24, 24, 27, 0.96) 76%)`
  })
</script>
<template>
  <div class="flex gap-3 overflow-x-auto overflow-y-hidden" :class="compact ? 'pb-0' : 'pb-1'">
    <button
      v-for="item in props.items"
      :key="item.id"
      type="button"
      class="flex min-h-0 flex-col items-stretch justify-start gap-2 border px-3 py-3 text-left text-(--lab-text-primary) transition"
      :class="[
        item.id === props.selectedId ? 'astro-body-card-active' : '',
        compact ? 'basis-33 p-2.5' : 'basis-40',
        wide ? 'basis-48' : '',
        'border-(--lab-border)'
      ]"
      :style="itemStyle(item.color)"
      @click="emit('select', item.id)">
      <span class="flex min-w-0 items-center gap-2">
        <span class="shrink-0 text-base font-semibold leading-none" :class="props.largeSymbol ? 'text-2xl' : ''">
          {{ item.symbol }}
        </span>
        <span class="min-w-0 flex-1 text-xs leading-5 wrap-break-word">{{ item.name }}</span>
      </span>
      <span class="min-w-0 text-[0.6875rem] leading-[1.35] text-(--lab-text-muted) wrap-break-word">{{ item.meta }}</span>
    </button>
  </div>
</template>
