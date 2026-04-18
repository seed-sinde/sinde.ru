<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    datetime: number | string | Date | null | undefined
    locale?: string
    compact?: boolean
    fallback?: string
    relative?: boolean
    showTime?: boolean
    titleMode?: 'datetime' | 'date' | 'none'
  }>(),
  {
    locale: 'ru-RU',
    compact: false,
    fallback: '—',
    relative: true,
    showTime: true,
    titleMode: 'datetime'
  }
)
const { localeTag } = useInterfacePreferences()
const { locale: i18nLocale, key, load, t } = useI18nSection('ui')
await useAsyncData(key, load, { watch: [i18nLocale] })
const { normalizeLocalizedDate, formatRelativeTime } = useLocalizedDateTime()
const resolvedLocale = computed(() => props.locale || localeTag.value)
const date = computed(() => normalizeLocalizedDate(props.datetime))
const rootClass = computed(() => (props.compact ? 'text-xs text-zinc-500' : 'text-sm text-zinc-400'))
const hydrated = ref(false)
const nowTs = ref(0)
const titleText = computed(() => {
  if (!date.value || props.titleMode === 'none') return undefined
  return new Intl.DateTimeFormat(
    resolvedLocale.value,
    props.titleMode === 'date'
      ? { day: '2-digit', month: '2-digit', year: 'numeric' }
      : { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' }
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
  hydrated.value = true
  nowTs.value = Date.now()
  nowTimer = setInterval(() => {
    nowTs.value = Date.now()
  }, 30_000)
})
onBeforeUnmount(() => {
  if (nowTimer) clearInterval(nowTimer)
})
const relativeText = computed(() => {
  void nowTs.value
  return formatRelativeTime(date.value, props.fallback, 'auto')
})
</script>
<template>
  <span v-if="date" :class="rootClass" :title="titleText">
    <template v-if="relative && hydrated">
      {{ relativeText || t('time.less_than_minute_ago') }}
    </template>
    <NuxtTime v-else :datetime="date" :locale="resolvedLocale" v-bind="absoluteDateProps" />
  </span>
  <span v-else :class="rootClass">{{ fallback }}</span>
</template>
