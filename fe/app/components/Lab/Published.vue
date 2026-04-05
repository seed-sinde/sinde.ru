<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      unix: number | string | Date
      label?: string
      locale?: string
      showTime?: boolean
      compact?: boolean
      relative?: boolean
    }>(),
    {
      label: 'опубликовано',
      locale: 'ru-RU',
      showTime: true,
      compact: false,
      relative: false
    }
  )
  const normalizeDate = (value: number | string | Date): Date | null => {
    if (value instanceof Date) {
      return Number.isNaN(value.getTime()) ? null : value
    }
    if (typeof value === 'number') {
      if (!Number.isFinite(value)) return null
      const ms = value > 1e11 ? value : value * 1000
      const result = new Date(ms)
      return Number.isNaN(result.getTime()) ? null : result
    }
    const normalized = value.trim()
    if (!normalized) return null
    if (/^\d+$/.test(normalized)) {
      const asNumber = Number(normalized)
      if (!Number.isFinite(asNumber)) return null
      const ms = asNumber > 1e11 ? asNumber : asNumber * 1000
      const result = new Date(ms)
      return Number.isNaN(result.getTime()) ? null : result
    }
    const result = new Date(normalized)
    return Number.isNaN(result.getTime()) ? null : result
  }
  const date = computed(() => normalizeDate(props.unix))
  const dateIso = computed(() => (date.value ? date.value.toISOString() : ''))
  const rootClass = computed(() => {
    if (props.compact) {
      return 'lab-text-soft inline-flex max-w-full items-center gap-1 text-xs leading-4'
    }
    return 'lab-text-soft inline-flex max-w-full flex-wrap items-center gap-x-1.5 gap-y-0.5 text-xs leading-5 sm:text-sm'
  })
  const dateClass = computed(() => {
    if (props.compact) return 'lab-text-muted tabular-nums'
    return 'lab-text-muted tabular-nums'
  })
</script>
<template>
  <div :class="rootClass">
    <span class="inline-flex items-center gap-1 whitespace-nowrap">{{ label }}:</span>
    <time v-if="date" :datetime="dateIso" :class="dateClass" :title="dateIso">
      <LabRelativeTime
        v-if="relative"
        :datetime="date"
        :locale="locale"
        :compact="compact"
        :title-mode="showTime ? 'datetime' : 'date'" />
      <NuxtTime
        v-else
        :datetime="date"
        :locale="locale"
        day="2-digit"
        month="2-digit"
        year="numeric"
        :hour="showTime ? '2-digit' : undefined"
        :minute="showTime ? '2-digit' : undefined" />
    </time>
    <span v-else class="lab-text-soft">дата не указана</span>
  </div>
</template>
