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
  <div class="astro-card-strip" :class="compact ? 'astro-card-strip-compact' : ''">
    <button
      v-for="item in props.items"
      :key="item.id"
      type="button"
      class="astro-body-card"
      :class="[
        item.id === props.selectedId ? 'astro-body-card-active' : '',
        compact ? 'astro-body-card-compact' : '',
        wide ? 'astro-body-card-wide' : ''
      ]"
      :style="itemStyle(item.color)"
      @click="emit('select', item.id)">
      <span class="astro-body-card-title">
        <span class="astro-body-card-symbol" :class="props.largeSymbol ? 'astro-body-card-symbol-lg' : ''">
          {{ item.symbol }}
        </span>
        <span class="astro-body-card-name">{{ item.name }}</span>
      </span>
      <span class="astro-body-card-meta">{{ item.meta }}</span>
    </button>
  </div>
</template>
<style scoped>
  .astro-card-strip {
    display: flex;
    gap: 0.75rem;
    overflow-x: auto;
    overflow-y: hidden;
    padding-bottom: 0.25rem;
  }
  .astro-card-strip-compact {
    padding-bottom: 0;
  }
  .astro-body-card {
    display: flex;
    flex: 0 0 10rem;
    flex-direction: column;
    align-items: stretch;
    justify-content: flex-start;
    gap: 0.5rem;
    min-height: 0;
    border: 1px solid rgb(63 63 70);
    padding: 0.75rem;
    text-align: left;
    color: rgb(244 244 245);
    transition: background-color 160ms ease;
  }
  .astro-body-card-compact {
    flex-basis: 8.25rem;
    padding: 0.625rem;
  }
  .astro-body-card-wide {
    flex-basis: 12rem;
  }
  .astro-body-card-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    min-width: 0;
  }
  .astro-body-card-symbol {
    flex-shrink: 0;
    font-size: 1rem;
    font-weight: 600;
    line-height: 1;
  }
  .astro-body-card-symbol-lg {
    font-size: 1.5rem;
  }
  .astro-body-card-name {
    min-width: 0;
    flex: 1 1 auto;
    font-size: 0.75rem;
    line-height: 1.25;
    white-space: normal;
    overflow-wrap: anywhere;
    word-break: break-word;
  }
  .astro-body-card-meta {
    min-width: 0;
    font-size: 0.6875rem;
    line-height: 1.35;
    color: rgb(212 212 216);
    white-space: normal;
    overflow-wrap: anywhere;
    word-break: break-word;
  }
</style>
