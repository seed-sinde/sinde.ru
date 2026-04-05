<template>
  <div
    v-if="trait"
    role="checkbox"
    :aria-checked="selected ? 'true' : 'false'"
    tabindex="0"
    class="traits-card group cursor-pointer select-none p-3 focus-visible:outline-none"
    :class="selected ? 'traits-card-selected' : 'traits-card-idle'"
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
              <span class="lab-text-primary truncate text-sm font-semibold">
                {{ trait.t_key }}
              </span>
              <span class="lab-text-soft text-xs">—</span>
              <span class="lab-text-secondary inline-flex min-w-0 items-center gap-1.5 text-sm">
                <span
                  v-if="colorInfo"
                  class="h-2 w-2 shrink-0 rounded-full border border-zinc-500/70"
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
            class="traits-card-link inline-flex h-7 w-7 shrink-0 items-center justify-center focus-visible:outline-none"
            :aria-label="`Открыть особенность ${trait.t_key}`"
            title="Открыть отдельную особенность"
            @click.stop
            @keydown.stop>
            <Icon name="ic:round-open-in-new" class="h-4 w-4" />
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
  const { trait, selected } = defineProps<{
    trait: Trait | null
    selected?: boolean
  }>()
  const emit = defineEmits<{ (e: 'update:selected', value: boolean): void }>()
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
  const dataTypeLabel = computed(() => getDataTypeLabel(metaType.value, 'тип'))
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
<style scoped>
  .traits-card {
    border: 1px solid color-mix(in srgb, var(--lab-border) 82%, transparent);
    background: color-mix(in srgb, var(--lab-bg-surface) 88%, transparent);
    transition:
      border-color 0.15s ease,
      background-color 0.15s ease,
      color 0.15s ease;
  }
  .traits-card-idle:hover {
    border-color: color-mix(in srgb, var(--lab-border-strong) 88%, transparent);
    background: color-mix(in srgb, var(--lab-bg-surface-hover) 78%, transparent);
  }
  .traits-card-selected {
    border-color: color-mix(in srgb, var(--lab-accent) 52%, var(--lab-border));
    background: color-mix(in srgb, var(--lab-accent) 10%, var(--lab-bg-surface));
  }
  .traits-card:focus-visible {
    box-shadow: 0 0 0 2px color-mix(in srgb, var(--lab-focus-ring) 100%, transparent);
  }
  .traits-card-link {
    color: var(--lab-text-muted);
    transition: color 0.15s ease;
  }
  .traits-card-link:hover,
  .traits-card-link:focus-visible {
    color: var(--lab-accent);
  }
</style>
