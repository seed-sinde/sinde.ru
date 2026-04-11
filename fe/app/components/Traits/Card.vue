<template>
  <div
    v-if="trait"
    role="checkbox"
    :aria-checked="selected ? 'true' : 'false'"
    tabindex="0"
    class="group cursor-pointer select-none border p-3 transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-(--lab-focus-ring)"
    :class="
      selected
        ? 'border-[color-mix(in_srgb,var(--lab-accent)_52%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-accent)_10%,var(--lab-bg-surface))]'
        : 'border-[color-mix(in_srgb,var(--lab-border)_82%,transparent)] bg-[color-mix(in_srgb,var(--lab-bg-surface)_88%,transparent)] hover:border-[color-mix(in_srgb,var(--lab-border-strong)_88%,transparent)] hover:bg-[color-mix(in_srgb,var(--lab-bg-surface-hover)_78%,transparent)]'
    "
    @click="toggleSelection"
    @keydown.enter.prevent="toggleSelection"
    @keydown.space.prevent="toggleSelection">
    <div class="flex items-start gap-2.5">
      <LabBaseCheckbox
        :name="trait.t_uuid"
        :model-value="Boolean(selected)"
        bare
        checkbox-class="mt-0.5 shrink-0"
        @click.stop="toggleSelection"
        aria-hidden="true" />
      <div class="min-w-0 flex-1">
        <div class="flex items-start justify-between gap-2">
          <div class="min-w-0">
            <div class="flex min-w-0 flex-wrap items-center gap-x-2 gap-y-1">
              <span class="truncate text-sm font-semibold">
                {{ trait.t_key }}
              </span>
              <span class="lab-text-soft text-xs">—</span>
              <span class="lab-text-secondary inline-flex min-w-0 items-center gap-1.5 text-sm">
                <span
                  v-if="colorInfo"
                  class="h-2 w-2 shrink-0 rounded-full border"
                  :style="{ backgroundColor: colorInfo.css }"
                  aria-hidden="true"></span>
                <span class="min-w-0 wrap-break-word leading-5">
                  {{ colorInfo ? colorInfo.text : displayValue }}
                </span>
              </span>
              <span v-if="numberUnit" class="lab-text-muted shrink-0 text-xs font-mono">
                {{ numberUnit }}
              </span>
            </div>
            <div class="lab-text-muted mt-1 flex flex-wrap items-center gap-x-2 gap-y-1 text-xs leading-4">
              <span class="uppercase tracking-wide">{{ dataTypeLabel }}</span>
              <span v-if="dataCategory" class="lab-text-soft">•</span>
              <span v-if="dataCategory">{{ dataCategory }}</span>
              <span v-if="colorModeLabel" class="lab-text-soft">•</span>
              <span v-if="colorModeLabel">{{ colorModeLabel }}</span>
            </div>
          </div>
          <NuxtLink
            :to="detailLinkTo"
            class="shrink-0 text-xs uppercase tracking-[0.14em] text-(--lab-text-muted) transition hover:text-(--lab-accent) focus-visible:outline-none focus-visible:text-(--lab-accent)"
            :aria-label="copy.card.openAria.replace('{key}', trait.t_key)"
            :title="copy.card.openTitle"
            @click.stop
            @keydown.stop>
            {{ copy.card.openLabel }}
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
  const { localeCode } = useInterfacePreferences()
  const { trait, selected } = defineProps<{
    trait: Trait | null
    selected?: boolean
  }>()
  const emit = defineEmits<{ (e: 'update:selected', value: boolean): void }>()
  const copy = computed(() => TRAITS_WORKSPACE_COPY[localeCode.value] || TRAITS_WORKSPACE_COPY.ru)
  const toggleSelection = () => emit('update:selected', !selected)
  const store = useTraitsStore()
  const route = useRoute()
  const detailLinkTo = computed(() => {
    const uuid = String(trait?.t_uuid || '').trim()
    if (!uuid) return '/traits'
    const rawSetUuid = route.params.uuid
    const setUuid = typeof rawSetUuid === 'string' ? rawSetUuid : Array.isArray(rawSetUuid) ? rawSetUuid[0] : ''
    const isSetUuidPage = /^\/traits\/[0-9a-f-]+$/i.test(route.path) && !route.path.startsWith('/traits/trait/')
    return isSetUuidPage && setUuid
      ? { path: `/traits/trait/${uuid}`, query: { set: setUuid } }
      : { path: `/traits/trait/${uuid}` }
  })
  const metaEntry = computed(() => (trait?.t_key_id ? store.keyMetaById[trait.t_key_id] : null))
  const metaRaw = computed<Record<string, any>>(() => {
    const raw = metaEntry.value?.meta
    return raw && typeof raw === 'object' ? raw : {}
  })
  const metaType = computed(() => String(metaRaw.value.dataType || metaRaw.value.type || '').trim())
  const dataTypeLabel = computed(() => getDataTypeLabel(metaType.value, copy.value.detail.metaMissing))
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
    return formatTraitValueForDisplay(String(trait?.t_value || ''), metaType.value)
  })
  const colorInfo = computed(() => {
    if (metaType.value !== 'color' || !trait?.t_value) return null
    return resolveColorPreviewInfo(trait.t_value, metaRaw.value)
  })
</script>
