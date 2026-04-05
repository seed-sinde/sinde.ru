<template>
  <div
    :class="[rootClass, !iconOnly && 'animate-pulse']"
    :aria-label="label"
    aria-live="polite"
    aria-busy="true"
    role="status">
    <template v-if="iconOnly">
      <Icon :name="iconName" class="lab-loader-icon-spin" />
    </template>
    <template v-else-if="inline">
      <div class="lab-loader-icon"></div>
      <div class="lab-loader-block h-4 w-32"></div>
    </template>
    <template v-else-if="variant === 'list'">
      <div v-for="item in normalizedCount" :key="`list-${item}`" class="lab-loader-surface flex items-start gap-3 p-3">
        <div v-if="showAvatar" class="lab-loader-block h-10 w-10 shrink-0" />
        <div class="min-w-0 flex-1 space-y-2">
          <div class="lab-loader-block h-4 w-2/5" />
          <div class="lab-loader-block-soft h-3 w-full" />
          <div class="lab-loader-block-soft h-3 w-4/5" />
        </div>
        <div v-if="showMeta" class="hidden w-16 shrink-0 sm:block">
          <div class="lab-loader-block-soft ml-auto h-3 w-12" />
        </div>
      </div>
    </template>
    <template v-else>
      <div :class="gridClass">
        <div v-for="item in normalizedCount" :key="`card-${item}`" class="lab-loader-surface p-4">
          <div v-if="showMedia" class="lab-loader-block mb-4 h-32 w-full" />
          <div class="space-y-3">
            <div class="lab-loader-block h-5 w-3/5" />
            <div class="lab-loader-block-soft h-3 w-full" />
            <div class="lab-loader-block-soft h-3 w-5/6" />
            <div class="lab-loader-block-soft h-3 w-2/3" />
          </div>
          <div v-if="showFooter" class="mt-4 flex items-center justify-between gap-3">
            <div class="lab-loader-block h-8 w-20" />
            <div class="lab-loader-block-soft h-4 w-12" />
          </div>
        </div>
      </div>
    </template>
    <span class="sr-only">{{ label }}</span>
  </div>
</template>
<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      variant?: 'list' | 'cards' | 'inline' | 'icon'
      count?: number
      columns?: 1 | 2 | 3 | 4
      label?: string
      showAvatar?: boolean
      showMeta?: boolean
      showMedia?: boolean
      showFooter?: boolean
      iconName?: string
    }>(),
    {
      variant: 'list',
      count: 3,
      columns: 3,
      label: 'Загрузка контента',
      showAvatar: true,
      showMeta: true,
      showMedia: true,
      showFooter: true,
      iconName: 'ic:round-autorenew'
    }
  )
  const inline = computed(() => props.variant === 'inline')
  const iconOnly = computed(() => props.variant === 'icon')
  const rootClass = computed(() => {
    if (iconOnly.value) return 'inline-flex h-6 w-6 items-center justify-center'
    if (inline.value) return 'flex items-center gap-3'
    return 'space-y-3'
  })
  const normalizedCount = computed(() => {
    if (!Number.isFinite(props.count)) return 3
    return Math.max(1, Math.min(12, Math.trunc(props.count)))
  })
  const gridClass = computed(() => {
    const columnClassMap = {
      1: 'grid-cols-1',
      2: 'grid-cols-1 md:grid-cols-2',
      3: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3',
      4: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-4'
    } as const
    return ['grid gap-4', columnClassMap[props.columns]]
  })
</script>
<style scoped>
  .lab-loader-surface {
    border: 1px solid color-mix(in srgb, var(--lab-border) 72%, transparent);
    background: color-mix(in srgb, var(--lab-bg-surface) 92%, transparent);
  }
  .lab-loader-block,
  .lab-loader-block-soft,
  .lab-loader-icon {
    background: color-mix(in srgb, var(--lab-text-muted) 24%, var(--lab-bg-surface));
  }
  .lab-loader-block-soft {
    background: color-mix(in srgb, var(--lab-text-soft) 20%, var(--lab-bg-surface));
  }
  .lab-loader-block,
  .lab-loader-block-soft {
    border-radius: 0;
  }
  .lab-loader-icon {
    height: 1rem;
    width: 1rem;
    flex-shrink: 0;
    border-radius: 9999px;
  }
  .lab-loader-icon-spin {
    height: 1.5rem;
    width: 1.5rem;
    flex-shrink: 0;
    animation: lab-loader-icon-spin 0.8s linear infinite;
  }
  @keyframes lab-loader-icon-spin {
    to {
      transform: rotate(360deg);
    }
  }
</style>
