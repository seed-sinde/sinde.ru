<template>
  <div
    :class="[rootClass, !iconOnly && 'animate-pulse']"
    :aria-label="label"
    aria-live="polite"
    aria-busy="true"
    role="status">
    <template v-if="iconOnly">
      <Icon :name="iconName" class="h-6 w-6 shrink-0 animate-spin" />
    </template>
    <template v-else-if="inline">
      <div
        class="h-4 w-4 shrink-0 rounded-full bg-[color-mix(in_srgb,var(--lab-text-muted)_24%,var(--lab-bg-surface))]"></div>
      <div class="h-4 w-32 bg-[color-mix(in_srgb,var(--lab-text-muted)_24%,var(--lab-bg-surface))]"></div>
    </template>
    <template v-else-if="variant === 'list'">
      <div
        v-for="item in normalizedCount"
        :key="`list-${item}`"
        class="flex items-start gap-3 border-[color-mix(in_srgb,var(--lab-border)_72%,transparent)] bg-[color-mix(in_srgb,var(--lab-bg-surface)_92%,transparent)] p-3">
        <div
          v-if="showAvatar"
          class="h-10 w-10 shrink-0 bg-[color-mix(in_srgb,var(--lab-text-muted)_24%,var(--lab-bg-surface))]" />
        <div class="min-w-0 flex-1 space-y-2">
          <div class="h-4 w-2/5 bg-[color-mix(in_srgb,var(--lab-text-muted)_24%,var(--lab-bg-surface))]" />
          <div class="h-3 w-full bg-[color-mix(in_srgb,var(--lab-text-muted)_20%,var(--lab-bg-surface))]" />
          <div class="h-3 w-4/5 bg-[color-mix(in_srgb,var(--lab-text-muted)_20%,var(--lab-bg-surface))]" />
        </div>
        <div v-if="showMeta" class="hidden w-16 shrink-0 sm:block">
          <div class="ml-auto h-3 w-12 bg-[color-mix(in_srgb,var(--lab-text-muted)_20%,var(--lab-bg-surface))]" />
        </div>
      </div>
    </template>
    <template v-else>
      <div :class="gridClass">
        <div
          v-for="item in normalizedCount"
          :key="`card-${item}`"
          class="border-[color-mix(in_srgb,var(--lab-border)_72%,transparent)] bg-[color-mix(in_srgb,var(--lab-bg-surface)_92%,transparent)] p-4">
          <div
            v-if="showMedia"
            class="mb-4 h-32 w-full bg-[color-mix(in_srgb,var(--lab-text-muted)_24%,var(--lab-bg-surface))]" />
          <div class="space-y-3">
            <div class="h-5 w-3/5 bg-[color-mix(in_srgb,var(--lab-text-muted)_24%,var(--lab-bg-surface))]" />
            <div class="h-3 w-full bg-[color-mix(in_srgb,var(--lab-text-muted)_20%,var(--lab-bg-surface))]" />
            <div class="h-3 w-5/6 bg-[color-mix(in_srgb,var(--lab-text-muted)_20%,var(--lab-bg-surface))]" />
            <div class="h-3 w-2/3 bg-[color-mix(in_srgb,var(--lab-text-muted)_20%,var(--lab-bg-surface))]" />
          </div>
          <div v-if="showFooter" class="mt-4 flex items-center justify-between gap-3">
            <div class="h-8 w-20 bg-[color-mix(in_srgb,var(--lab-text-muted)_24%,var(--lab-bg-surface))]" />
            <div class="h-4 w-12 bg-[color-mix(in_srgb,var(--lab-text-muted)_20%,var(--lab-bg-surface))]" />
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
