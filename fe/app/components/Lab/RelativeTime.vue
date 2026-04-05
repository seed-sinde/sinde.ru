<script setup lang="ts">
const props = withDefaults(defineProps<{
  datetime: number | string | Date | null | undefined
  locale?: string
  compact?: boolean
  fallback?: string
  relative?: boolean
  showTime?: boolean
  titleMode?: 'datetime' | 'date' | 'none'
}>(), {
  locale: 'ru-RU',
  compact: false,
  fallback: '—',
  relative: true,
  showTime: true,
  titleMode: 'datetime',
})
const normalizeDate = (value: number | string | Date | null | undefined): Date | null => {
  if (!value && value !== 0) return null
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
const date = computed(() => normalizeDate(props.datetime))
const rootClass = computed(() => props.compact ? 'text-xs text-zinc-500' : 'text-sm text-zinc-400')
const nowTs = ref(Date.now())
const isLessThanMinute = computed(() => {
  if (!date.value) return false
  return Math.abs(nowTs.value - date.value.getTime()) < 60_000
})
const titleText = computed(() => {
  if (!date.value || props.titleMode === 'none') return undefined
  return new Intl.DateTimeFormat(props.locale, props.titleMode === 'date'
    ? { day: '2-digit', month: '2-digit', year: 'numeric' }
    : { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' },
  ).format(date.value)
})
const absoluteDateProps = computed(() => {
  const options = {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  } as const
  if (!props.showTime) return options
  return {
    ...options,
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  } as const
})
let nowTimer: ReturnType<typeof setInterval> | null = null
onMounted(() => {
  nowTs.value = Date.now()
  nowTimer = setInterval(() => {
    nowTs.value = Date.now()
  }, 30_000)
})
onBeforeUnmount(() => {
  if (nowTimer) clearInterval(nowTimer)
})
</script>
<template>
  <span v-if="date" :class="rootClass" :title="titleText">
    <template v-if="relative">
      <template v-if="isLessThanMinute">меньше минуты назад</template>
      <NuxtTime v-else :datetime="date" :locale="locale" relative />
    </template>
    <NuxtTime
      v-else
      :datetime="date"
      :locale="locale"
      v-bind="absoluteDateProps" />
  </span>
  <span v-else :class="rootClass">{{ fallback }}</span>
</template>
