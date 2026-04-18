<template>
  <div
    v-if="trait"
    role="checkbox"
    :aria-checked="selected ? 'true' : 'false'"
    tabindex="0"
    class="lab-focus group cursor-pointer px-2 py-1.5 select-none"
    :class="[selected && 'lab-focus-manual', stateToneClass]"
    @click="toggleSelection"
    @dblclick.stop="onRequestEdit"
    @keydown.enter.prevent="toggleSelection"
    @keydown.space.prevent="toggleSelection"
    @pointerdown="onPointerDown"
    @pointerup="onPointerUp"
    @pointerleave="onPointerUp"
    @pointercancel="onPointerUp"
  >
    <div class="flex items-start gap-2.5">
      <span
        aria-hidden="true"
        class="mt-0.5 inline-flex h-4 w-4 shrink-0 items-center justify-center border transition"
        :class="
          selected
            ? 'border-(--lab-accent) text-(--lab-accent)'
            : 'border-(--lab-border) text-transparent group-hover:border-(--lab-border-strong)'
        "
      >
        <Icon name="ic:round-check" class="text-sm" />
      </span>
      <div class="min-w-0 flex-1">
        <div class="flex items-start justify-between gap-2">
          <div class="min-w-0">
            <div class="flex min-w-0 flex-wrap items-center gap-x-2 gap-y-1">
              <span class="truncate text-sm font-semibold">
                {{ trait.t_key }}
              </span>
              <span class="text-xs text-(--lab-text-soft)">—</span>
              <span class="inline-flex min-w-0 items-center gap-1.5 text-sm text-(--lab-text-secondary)">
                <span
                  v-if="colorInfo"
                  class="h-2 w-2 shrink-0 rounded-full border"
                  :style="{ backgroundColor: colorInfo.css }"
                  aria-hidden="true"
                />
                <span class="min-w-0 leading-5 wrap-break-word">
                  {{ colorInfo ? colorInfo.text : displayValue }}
                </span>
              </span>
              <span v-if="numberUnit" class="shrink-0 font-mono text-xs text-(--lab-text-muted)">
                {{ numberUnit }}
              </span>
            </div>
            <div class="mt-1 flex flex-wrap items-center gap-x-2 gap-y-1 text-xs leading-4 text-(--lab-text-muted)">
              <span class="tracking-wide uppercase">{{ dataTypeLabel }}</span>
              <span v-if="dataCategory" class="text-(--lab-text-soft)">•</span>
              <span v-if="dataCategory">{{ dataCategory }}</span>
              <span v-if="colorModeLabel" class="text-(--lab-text-soft)">•</span>
              <span v-if="colorModeLabel">{{ colorModeLabel }}</span>
            </div>
          </div>
          <NuxtLink
            :to="detailLinkTo"
            class="lab-focus shrink-0 text-[11px] tracking-[0.14em] text-(--lab-text-muted) uppercase"
            :aria-label="t('card.open.aria', { key: trait.t_key })"
            :title="t('card.open.title')"
            @click.stop
            @keydown.stop
          >
            {{ t('card.open.label') }}
          </NuxtLink>
        </div>
      </div>
    </div>
    <div v-if="editing" class="mt-2" @click.stop @keydown.stop>
      <TraitsEditValueForm
        :trait="trait"
        :meta="editingMeta"
        :pending="editPending"
        inline
        @save="emit('save-inline', $event)"
        @cancel="emit('cancel-inline')"
      />
    </div>
  </div>
</template>
<script setup lang="ts">
const { locale, key, load, t } = useI18nSection('traits')
await useAsyncData(key.value, load, { watch: [locale] })
const props = withDefaults(
  defineProps<{
    trait: Trait | null
    selected?: boolean
    editing?: boolean
    editingMeta?: KeyMeta | null
    editPending?: boolean
    stateTone?: 'added' | 'updated' | 'removing' | ''
  }>(),
  {
    selected: false,
    editing: false,
    editingMeta: null,
    editPending: false,
    stateTone: ''
  }
)
const { trait, selected, editing, editingMeta, editPending, stateTone } = toRefs(props)
const emit = defineEmits<{
  (e: 'update:selected', value: boolean): void
  (e: 'request-edit' | 'cancel-inline'): void
  (e: 'save-inline', payload: { traitUuid: string; t_key: string; t_value: string }): void
}>()
const stateToneClass = computed(() => {
  if (stateTone.value === 'removing') {
    return 'bg-[color-mix(in_srgb,#ef4444_20%,transparent)]'
  }
  if (stateTone.value === 'updated') {
    return 'bg-[color-mix(in_srgb,#f59e0b_20%,transparent)]'
  }
  if (stateTone.value === 'added') {
    return 'bg-[color-mix(in_srgb,#3b82f6_20%,transparent)]'
  }
  return ''
})
const LONG_PRESS_MS = 520
const longPressTimer = ref<ReturnType<typeof setTimeout> | null>(null)
const longPressTriggered = ref(false)
const clearLongPress = () => {
  if (longPressTimer.value) {
    clearTimeout(longPressTimer.value)
    longPressTimer.value = null
  }
}
const onRequestEdit = () => {
  emit('request-edit')
}
const onPointerDown = (event: PointerEvent) => {
  if (event.pointerType !== 'touch') return
  clearLongPress()
  longPressTriggered.value = false
  longPressTimer.value = setTimeout(() => {
    longPressTriggered.value = true
    onRequestEdit()
  }, LONG_PRESS_MS)
}
const onPointerUp = () => {
  clearLongPress()
}
const toggleSelection = () => {
  if (longPressTriggered.value) {
    longPressTriggered.value = false
    return
  }
  emit('update:selected', !selected.value)
}
const store = useTraitsStore()
const route = useRoute()
const detailLinkTo = computed(() => {
  const uuid = String(trait.value?.t_uuid || '').trim()
  if (!uuid) return '/traits'
  const rawSetUuid = route.params.uuid
  const setUuid = typeof rawSetUuid === 'string' ? rawSetUuid : Array.isArray(rawSetUuid) ? rawSetUuid[0] : ''
  const isSetUuidPage = /^\/traits\/[0-9a-f-]+$/i.test(route.path) && !route.path.startsWith('/traits/trait/')
  return isSetUuidPage && setUuid
    ? { path: `/traits/trait/${uuid}`, query: { set: setUuid } }
    : { path: `/traits/trait/${uuid}` }
})
const metaEntry = computed(() => (trait.value?.t_key_id ? store.keyMetaById[trait.value.t_key_id] : null))
const metaRaw = computed<Record<string, any>>(() => {
  const raw = metaEntry.value?.meta
  return raw && typeof raw === 'object' ? raw : {}
})
const metaType = computed(() => String(metaRaw.value.dataType || metaRaw.value.type || '').trim())
const dataTypeLabel = computed(() => getDataTypeLabel(metaType.value, t('detail.meta_missing')))
const dataCategory = computed(() => {
  return getDataCategoryLabel(metaRaw.value)
})
const colorModeLabel = computed(() => {
  if (metaType.value !== 'color') return ''
  return getColorModeLabel(metaRaw.value)
})
const numberUnit = computed(() => {
  if (metaType.value !== 'number') return ''
  return getNumberUnitLabel(String(metaRaw.value.unit || ''))
})
const displayValue = computed(() => {
  return formatTraitValueForDisplay(String(trait.value?.t_value || ''), metaType.value)
})
const colorInfo = computed(() => {
  if (metaType.value !== 'color' || !trait.value?.t_value) return null
  return resolveColorPreviewInfo(trait.value.t_value, metaRaw.value)
})
onBeforeUnmount(() => {
  clearLongPress()
})
</script>
