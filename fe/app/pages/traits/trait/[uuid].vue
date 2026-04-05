<script setup lang="ts">
  const formatShortUuid = shortUuid
  const route = useRoute()
  const store = useTraitsStore()
  const { error: pasteError, pasteUuidAndNavigate } = usePasteUuid()
  const traitUuid = computed(() => {
    const raw = route.params.uuid
    const candidate = typeof raw === 'string' ? raw : Array.isArray(raw) ? raw[0] : ''
    return String(candidate || '').trim()
  })
  const setUuid = computed(() => {
    const raw = route.query.set
    const candidate = typeof raw === 'string' ? raw : Array.isArray(raw) ? raw[0] : ''
    const normalized = String(candidate || '').trim()
    return UUID_RE.test(normalized) ? normalized : ''
  })
  if (!UUID_RE.test(traitUuid.value)) {
    throw createError({
      statusCode: 404,
      message: 'Особенность не найдена'
    })
  }
  const pasteUuid = async () => await pasteUuidAndNavigate(uuid => `/traits/trait/${uuid}`)
  const resolveTraitFromStore = () => {
    const targetUuid = traitUuid.value
    if (!targetUuid) return null
    const match = store.traits.find(item => String(item?.t_uuid || '').trim() === targetUuid)
    return match ? { ...match } : null
  }
  const { data: traitData, error: traitError } = await useAsyncData<Trait | null>(
    () => `trait-detail:${traitUuid.value}`,
    async () => {
      const fallback = resolveTraitFromStore()
      try {
        const item = await fetchTraitByUuid(traitUuid.value)
        return item || fallback || null
      } catch (error) {
        if (fallback) return fallback
        throw error
      }
    },
    {
      watch: [traitUuid],
      server: true,
      default: () => resolveTraitFromStore()
    }
  )
  if (traitError.value && !traitData.value) {
    throw createError({
      statusCode: 500,
      message: 'Не удалось загрузить особенность'
    })
  }
  if (!traitData.value) {
    throw createError({
      statusCode: 404,
      message: 'Особенность не найдена'
    })
  }
  const trait = computed<Trait | null>(() => traitData.value ?? null)
  const traitKey = computed(() => String(trait.value?.t_key || '').trim())
  const title = computed(() => `${traitKey.value || 'особенность'} · особенность`)
  const { data: metaEntryData } = await useAsyncData<TraitKey | null>(
    () => `trait-meta:${traitUuid.value}:${traitKey.value}`,
    async () => {
      if (!traitKey.value) return null
      try {
        const res = await getKeyMeta(traitKey.value)
        return res?.data || null
      } catch {
        return null
      }
    },
    {
      server: true,
      lazy: false,
      default: () => null,
      watch: [traitKey]
    }
  )
  const metaEntry = computed<TraitKey | null>(() => metaEntryData.value ?? null)
  const metaRaw = computed<Record<string, any>>(() => {
    const raw = metaEntry.value?.meta
    return raw && typeof raw === 'object' ? raw : {}
  })
  const metaType = computed(() => String(metaRaw.value.dataType || metaRaw.value.type || '').trim())
  const dataTypeLabel = computed(() => getDataTypeLabel(metaType.value, 'не задан'))
  const metaInfoEntries = computed(() => {
    const raw = metaRaw.value || {}
    const items: Array<{ key: string; label: string; value: string }> = []
    const keyId = trait.value?.t_key_id ?? metaEntry.value?.id
    if (typeof keyId === 'number' && Number.isFinite(keyId)) {
      items.push({
        key: 'keyId',
        label: 'id ключа',
        value: String(keyId)
      })
    }
    for (const [key, value] of Object.entries(raw)) {
      if (key === 'format') continue
      const text = formatMetaValue(key, value, dataTypeLabel.value)
      if (!text) continue
      items.push({
        key,
        label: META_FIELD_LABELS[key] || key,
        value: text
      })
    }
    return items
  })
  const displayValue = computed(() => {
    return formatTraitValueForDisplay(String(trait.value?.t_value || ''), metaType.value)
  })
  const colorInfo = computed(() => {
    if (metaType.value !== 'color') return null
    return resolveColorPreviewInfo(String(trait.value?.t_value || ''), metaRaw.value)
  })
  const displayValueText = computed(() => colorInfo.value?.text || displayValue.value || '—')
  definePageMeta({
    robots: 'noindex, nofollow',
    key: route => route.path
  })
  usePageSeo({
    title,
    description: () => `Полная карточка особенности ${traitKey.value}`
  })
</script>
<template>
  <div v-if="trait" class="space-y-2">
    <LabNavHeader :title :breadcrumb-items="[{ label: 'Особенности', to: '/traits' }]">
      <template #actions="{ compact }">
        <Icon
          v-if="setUuid && !compact"
          name="ic:round-chevron-right"
          class="h-4 w-4 shrink-0 text-zinc-600"
          aria-hidden="true" />
        <NuxtLink
          v-if="setUuid && !compact"
          :to="`/traits/${setUuid}`"
          class="lab-text-muted min-w-0 truncate font-mono text-sm transition-colors hover:text-zinc-200"
          :title="`Открыть набор ${setUuid}`">
          {{ formatShortUuid(setUuid, 5) || '—' }}
        </NuxtLink>
        <Icon
          v-if="setUuid && !compact"
          name="ic:round-chevron-right"
          class="h-4 w-4 shrink-0 text-zinc-600"
          aria-hidden="true" />
        <TraitsUuidButton action="copy" :uuid="trait.t_uuid" :compact="compact" />
        <TraitsUuidButton action="paste" :compact="compact" @click="pasteUuid" />
      </template>
    </LabNavHeader>
    <LabNotify :text="pasteError" tone="error" class-name="px-4" />
    <section class="space-y-4 px-3 py-3 sm:px-4 sm:py-4 lg:px-5">
      <article class="trait-detail-card mx-auto max-w-5xl">
        <div class="border-b border-zinc-800/80 px-4 py-3 sm:px-5">
          <LabCopyBlock
            label="UUID"
            :value="trait.t_uuid"
            title-idle="Скопировать UUID"
            title-success="UUID скопирован"
            title-error="Ошибка копирования" />
        </div>
        <div class="border-b border-zinc-800/80 px-4 py-4 sm:px-5">
          <div
            class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between lg:gap-6">
            <div class="min-w-0 flex-1">
              <div class="flex flex-wrap items-start gap-x-3 gap-y-2">
                <h2 class="lab-text-primary text-lg font-semibold tracking-tight sm:text-2xl">{{ trait.t_key }}</h2>
                <span class="lab-text-soft text-lg sm:text-2xl">:</span>
                <div class="min-w-0 flex-1">
                  <div class="lab-text-primary flex items-center gap-2 text-base sm:text-xl">
                    <span
                      v-if="colorInfo"
                      class="inline-block h-2.5 w-2.5 shrink-0 self-center rounded-full border border-zinc-500/70"
                      :style="{ backgroundColor: colorInfo.css }"
                      aria-hidden="true" />
                    <span class="min-w-0 whitespace-pre-wrap wrap-break-word leading-relaxed">
                      {{ displayValueText }}
                    </span>
                  </div>
                </div>
              </div>
              <p
                v-if="trait.t_value && trait.t_value !== displayValueText"
                class="lab-text-muted mt-3 whitespace-pre-wrap wrap-break-word font-mono text-xs">
                {{ trait.t_value }}
              </p>
            </div>
            <div class="w-full max-w-lg space-y-2 lg:w-80 lg:flex-none">
              <p class="lab-text-muted text-xs uppercase tracking-wide">Параметры ключа</p>
              <div v-if="metaInfoEntries.length" class="flex flex-wrap gap-2 text-xs">
                <span
                  v-for="entry in metaInfoEntries"
                  :key="entry.key"
                  class="trait-meta-chip inline-flex items-center gap-1 px-2 py-1">
                  <span class="lab-text-muted">{{ entry.label }}:</span>
                  <span class="lab-text-secondary">{{ entry.value }}</span>
                </span>
              </div>
              <p v-else class="lab-text-primary text-sm sm:text-base">—</p>
            </div>
          </div>
        </div>
      </article>
    </section>
  </div>
</template>
<style scoped>
  .trait-detail-card {
    border: 1px solid color-mix(in srgb, var(--lab-border) 82%, transparent);
    background: color-mix(in srgb, var(--lab-bg-surface) 84%, transparent);
  }
  .trait-meta-chip {
    border: 1px solid color-mix(in srgb, var(--lab-border) 72%, transparent);
    background: color-mix(in srgb, var(--lab-bg-surface-subtle) 62%, transparent);
  }
</style>
