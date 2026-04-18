<script setup lang="ts">
useTraitsMobileHeader()
const { locale, key, load, t } = useI18nSection('traits')
await useAsyncData(key.value, load, { watch: [locale] })
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
    message: t('detail.not_found')
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
    } catch (error: any) {
      const status = Number(error?.statusCode ?? error?.status ?? error?.response?.status ?? 0)
      if (status === 404) {
        return fallback || null
      }
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
  const statusCode = Number((traitError.value as any)?.statusCode || (traitError.value as any)?.status || 500)
  throw createError({
    statusCode,
    message: statusCode === 404 ? t('detail.not_found') : t('detail.load_failed')
  })
}
const trait = computed<Trait | null>(() => traitData.value ?? null)
const traitKey = computed(() => String(trait.value?.t_key || '').trim())
const title = computed(
  () => `${traitKey.value || t('detail.title_fallback')} · ${t('detail.title_fallback')}`
)
const pageDescription = computed(() => t('detail.seo_description', { key: traitKey.value }))
const breadcrumbItems = computed<BreadcrumbItem[]>(() => [{ label: t('detail.breadcrumb'), to: '/traits' }])
const { data: metaEntryData } = await useAsyncData<TraitKey | null>(
  () => `trait-meta:${traitUuid.value}:${String(trait.value?.t_key_id || 0)}`,
  async () => {
    const keyId = Number(trait.value?.t_key_id || 0)
    if (keyId <= 0) return null
    try {
      const res = await getKeysMetaBulk([keyId])
      return res?.data?.items?.[0] || null
    } catch {
      return null
    }
  },
  {
    server: true,
    lazy: false,
    default: () => null,
    watch: [trait]
  }
)
const metaEntry = computed<TraitKey | null>(() => metaEntryData.value ?? null)
const metaRaw = computed<Record<string, any>>(() => {
  const raw = metaEntry.value?.meta
  return raw && typeof raw === 'object' ? raw : {}
})
const metaType = computed(() => String(metaRaw.value.dataType || metaRaw.value.type || '').trim())
const dataTypeLabel = computed(() => getDataTypeLabel(metaType.value, t('detail.meta_missing')))
const metaInfoEntries = computed(() => {
  const raw = metaRaw.value || {}
  const items: Array<{ key: string; label: string; value: string }> = []
  const keyId = trait.value?.t_key_id ?? metaEntry.value?.id
  if (typeof keyId === 'number' && Number.isFinite(keyId)) {
    items.push({
      key: 'keyId',
      label: t('detail.key.id'),
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
const mobileActions = computed<MobileHeaderAction[]>(() => [
  { kind: 'traits-copy-uuid', uuid: traitUuid.value },
  { kind: 'traits-paste-uuid', mode: 'trait' }
])
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
  description: pageDescription
})
</script>
<template>
  <div>
    <LabNavHeader
      :title
      :breadcrumb-items="breadcrumbItems"
      :mobile-actions="mobileActions"
    >
      <template #actions="{ compact }">
        <Icon
          v-if="setUuid && !compact"
          name="ic:round-chevron-right"
          class="h-4 w-4 shrink-0 text-zinc-600"
          aria-hidden="true"
        />
        <NuxtLink
          v-if="setUuid && !compact"
          :to="`/traits/${setUuid}`"
          class="min-w-0 truncate font-mono text-sm text-(--lab-text-muted) transition hover:text-(--lab-text-primary)"
          :title="t('detail.open_set_title', { uuid: setUuid })"
        >
          {{ formatShortUuid(setUuid, 5) || t('detail.copy_uuid_idle') }}
        </NuxtLink>
        <Icon
          v-if="setUuid && !compact"
          name="ic:round-chevron-right"
          class="h-4 w-4 shrink-0 text-zinc-600"
          aria-hidden="true"
        />
        <TraitsUuidButton action="copy" :uuid="traitUuid" :compact="compact" />
        <TraitsUuidButton action="paste" :compact="compact" @click="pasteUuid" />
      </template>
    </LabNavHeader>
    <LabNotify :text="pasteError" tone="error" class-name="px-4" />
    <section class="space-y-2 p-4">
      <div class="flex flex-wrap items-start gap-x-3 gap-y-2">
        <h2 class="text-lg font-semibold tracking-tight sm:text-2xl">{{ trait?.t_key }}</h2>
        <span class="text-lg text-(--lab-text-soft) sm:text-2xl">:</span>
        <div class="flex min-w-0 flex-1 items-center gap-2 text-base sm:text-xl">
          <span
            v-if="colorInfo"
            class="inline-block h-2.5 w-2.5 shrink-0 self-center rounded-full border"
            :style="{ backgroundColor: colorInfo.css }"
            aria-hidden="true"
          />
          <span class="min-w-0 leading-relaxed wrap-break-word whitespace-pre-wrap">
            {{ displayValueText }}
          </span>
        </div>
      </div>
      <p
        v-if="trait?.t_value && trait?.t_value !== displayValueText"
        class="mt-3 font-mono text-xs wrap-break-word whitespace-pre-wrap text-(--lab-text-muted)"
      >
        {{ trait?.t_value }}
      </p>
      <p class="w-full max-w-lg text-xs tracking-wide text-(--lab-text-muted) uppercase lg:w-80">
        {{ t('detail.key.params') }}
      </p>
      <div v-if="metaInfoEntries.length" class="flex w-full max-w-lg flex-wrap gap-2 text-xs lg:w-80">
        <span v-for="entry in metaInfoEntries" :key="entry.key" class="inline-flex items-center gap-1 px-2 py-1">
          <span class="text-(--lab-text-muted)">{{ entry.label }}:</span>
          <span class="text-(--lab-text-secondary)">{{ entry.value }}</span>
        </span>
      </div>
      <p v-else class="w-full max-w-lg text-sm sm:text-base lg:w-80">пусто</p>
    </section>
  </div>
</template>
