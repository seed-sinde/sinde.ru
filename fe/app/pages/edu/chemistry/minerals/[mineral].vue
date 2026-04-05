<script setup lang="ts">
  const route = useRoute()
  const requestEvent = useRequestEvent()
  const { t } = useInterfacePreferences()
  const databaseId = computed(() => String(route.params.mineral || '').trim())
  const notFound = ref(false)
  const formatText = (value: unknown) => String(value || '').trim()
  const formatList = (value?: string[] | null) =>
    (value || [])
      .map(item => String(item || '').trim())
      .filter(Boolean)
      .join(', ')
  const formatNumber = (value?: number | null) => {
    if (typeof value !== 'number' || !Number.isFinite(value)) return ''
    return new Intl.NumberFormat('ru-RU', { maximumFractionDigits: 4 }).format(value)
  }
  const compactItems = (items: Array<{ label: string; value: string }>) => items.filter(item => item.value)
  const toFormulaText = (value: unknown) => {
    const text = formatText(value)
    return {
      text,
      latex: mineralFormulaToLatex(text),
      isFormula: isMineralFormulaLike(text)
    }
  }
  const { data, pending, error, refresh } = await useAsyncData(
    () => `mineral:${databaseId.value || 'missing'}`,
    async () => {
      if (!databaseId.value) {
        notFound.value = true
        return null
      }
      try {
        notFound.value = false
        const response = await getMineralByDatabaseId(databaseId.value)
        return response.data
      } catch (err: any) {
        const status = Number(err?.statusCode || err?.status || err?.response?.status || 0)
        if (status === 404) {
          notFound.value = true
          return null
        }
        throw err
      }
    },
    {
      watch: [databaseId],
      default: () => null
    }
  )
  if (import.meta.server && notFound.value && requestEvent) {
    setResponseStatus(requestEvent, 404, 'Минерал не найден')
  }
  const isImageViewerOpen = ref(false)
  const activeImageIndex = ref(0)
  const mineral = computed(() => data.value)
  const mineralName = computed(() => formatText(mineral.value?.mineral_name))
  const mineralNamePlain = computed(() => formatText(mineral.value?.mineral_name_plain))
  const mineralSubtitle = computed(() =>
    mineralNamePlain.value && mineralNamePlain.value !== mineralName.value ? mineralNamePlain.value : ''
  )
  const mineralNameLatex = computed(() => mineralFormulaToLatex(mineralName.value))
  const showMineralNameAsFormula = computed(() => isMineralFormulaLike(mineralName.value))
  const mineralImages = computed(() => mineral.value?.images || [])
  const hasMineralImages = computed(() => mineralImages.value.length > 0)
  const mineralImageItems = computed(() =>
    mineralImages.value.map(item => ({
      key: `${item.file}:${item.rruff_id || ''}`,
      src: buildMediaFileUrl(`minerals/webp/${item.file}`),
      previewSrc: buildMediaFileUrl(`minerals/preview/${item.file}`),
      title: item.rruff_id || mineralName.value,
      alt: `${mineralName.value || 'Mineral'} ${item.rruff_id || ''}`.trim()
    }))
  )
  const mineralCoverImage = computed(
    () => mineralImageItems.value[activeImageIndex.value] || mineralImageItems.value[0] || null
  )
  const openMineralImageViewer = (index = 0) => {
    activeImageIndex.value = index
    isImageViewerOpen.value = true
  }
  const setMineralViewerIndex = (index: number) => {
    activeImageIndex.value = index
  }
  watch(
    mineralImageItems,
    items => {
      if (!items.length) {
        activeImageIndex.value = 0
        return
      }
      if (activeImageIndex.value >= items.length) {
        activeImageIndex.value = 0
      }
    },
    { immediate: true }
  )
  const chemistryElements = computed(() => mineral.value?.chemistry_elements || [])
  const valenceElements = computed(() => mineral.value?.valence_elements || [])
  const crystalSystems = computed(() => mineral.value?.crystal_systems || [])
  const crystalSystemLabel = (value: string) =>
    t(`minerals.crystal_system.${String(value || 'unknown').trim() || 'unknown'}` as InterfaceMessageKey)
  const spaceGroups = computed(() => mineral.value?.space_groups || [])
  const parageneticModes = computed(() => mineral.value?.paragenetic_modes || [])
  const goToMinerals = () => navigateTo('/edu/chemistry/minerals')
  const toElementChip = (symbol: string) => {
    const normalizedSymbol = String(symbol || '').trim()
    if (!normalizedSymbol) {
      return { symbol: '', to: null as string | null }
    }
    const element = resolvePeriodicTableElement(normalizedSymbol)
    return {
      symbol: normalizedSymbol,
      to: element ? getPeriodicTableElementRoute(element) : null
    }
  }
  const chemistryElementChips = computed(() => chemistryElements.value.map(toElementChip).filter(item => item.symbol))
  const valenceElementChips = computed(() =>
    valenceElements.value
      .map(symbol => {
        const chip = toElementChip(symbol)
        const formula = toFormulaText(chip.symbol)
        return {
          ...chip,
          latex: formula.latex,
          isFormula: formula.isFormula
        }
      })
      .filter(item => item.symbol)
  )
  const primaryChemistry = computed(() => formatText(mineral.value?.ima_chemistry))
  const primaryValenceChemistry = computed(() => formatText(mineral.value?.valence_chemistry))
  const referenceItems = computed(() =>
    compactItems([
      { label: 'IMA number', value: formatText(mineral.value?.ima_number) },
      { label: 'Database ID', value: formatText(mineral.value?.database_id) }
    ])
  )
  const chemistryItems = computed(() =>
    compactItems([
      { label: 'IMA chemistry', value: primaryChemistry.value },
      { label: 'Valence chemistry', value: primaryValenceChemistry.value }
    ])
  )
  const chemistryDisplayItems = computed(() =>
    chemistryItems.value.map(item => ({
      ...item,
      ...toFormulaText(item.value)
    }))
  )
  const classificationPrimaryItems = computed(() =>
    compactItems([
      { label: 'Crystal systems', value: crystalSystems.value.map(crystalSystemLabel).join(', ') },
      { label: 'Space groups', value: formatList(spaceGroups.value) },
      { label: 'IMA symbol', value: formatText(mineral.value?.ima_mineral_symbol) },
      { label: 'IMA status', value: formatText(mineral.value?.ima_status) }
    ])
  )
  const classificationSecondaryItems = computed(() =>
    compactItems([{ label: 'Structural group', value: formatText(mineral.value?.structural_groupname) }])
  )
  const originItems = computed(() =>
    compactItems([
      { label: 'Country', value: formatText(mineral.value?.country_of_type_locality) },
      { label: 'Published', value: formatText(mineral.value?.year_first_published) },
      { label: 'Oldest age Ma', value: formatNumber(mineral.value?.oldest_known_age_ma) }
    ])
  )
  const parageneticModeItems = computed(() => parageneticModes.value.map(item => formatText(item)).filter(Boolean))
  const notesItems = computed(() =>
    compactItems([{ label: 'Status notes', value: formatText(mineral.value?.status_notes) }])
  )
  const errorMessage = computed(() => {
    const value = error.value as any
    return value?.data?.message || value?.message || 'Не удалось загрузить минерал.'
  })
  const title = computed(() => mineralName.value || 'Минерал')
  const description = computed(() => {
    return (
      primaryChemistry.value ||
      formatText(mineral.value?.ima_status) ||
      'Подробная карточка минерала с химией, статусом и классификацией.'
    )
  })
  usePageSeo({
    title,
    description
  })
</script>
<template>
  <div class="space-y-4">
    <LabNavHeader
      :title
      :breadcrumb-items="[
        { label: 'Вики', to: '/edu' },
        { label: 'Химия', to: '/edu/chemistry' },
        { label: 'Минералы', to: '/edu/chemistry/minerals' },
        { label: title, current: true }
      ]" />
    <section class="px-3 sm:px-6">
      <div
        v-if="pending && !mineral && !notFound"
        class="flex min-h-72 items-center justify-center border border-zinc-800 bg-zinc-950/70">
        <div class="flex items-center gap-3 text-sm text-zinc-400">
          <LabLoader varian="list" label="Загружаем минерал..." />
        </div>
      </div>
      <div v-else-if="error" class="border border-rose-500/30 bg-rose-950/20 p-5">
        <LabErrorMessage :text="errorMessage" error-class="text-sm" />
        <div class="mt-3">
          <LabBaseButton variant="secondary" size="sm" label="Повторить" @click="refresh" />
        </div>
      </div>
      <div v-else-if="notFound" class="border border-zinc-800 bg-zinc-950/70 p-6">
        <div class="space-y-2">
          <div class="text-2xl font-semibold text-zinc-100">Минерал не найден</div>
          <p class="text-sm leading-6 text-zinc-400">
            Проверьте `database_id` в адресе или вернитесь к каталогу минералов.
          </p>
        </div>
        <div class="mt-4">
          <LabBaseButton variant="secondary" size="sm" label="К списку минералов" @click="goToMinerals" />
        </div>
      </div>
      <div v-else-if="mineral" class="border border-zinc-800 bg-zinc-950/70 p-4">
        <div :class="hasMineralImages ? 'grid gap-4 xl:grid-cols-[minmax(0,1fr)_minmax(0,2fr)]' : 'space-y-4'">
          <div class="space-y-4">
            <div class="space-y-3">
              <div class="space-y-2">
                <LabViewerLaTex
                  v-if="showMineralNameAsFormula"
                  :formula="mineralNameLatex"
                  class="text-2xl font-semibold text-zinc-100 sm:text-3xl" />
                <div v-else class="text-2xl font-semibold text-zinc-100 sm:text-3xl wrap-break-word">
                  {{ mineralName }}
                </div>
                <div v-if="mineralSubtitle" class="text-sm text-zinc-500 wrap-break-word">{{ mineralSubtitle }}</div>
              </div>
              <div class="flex flex-wrap items-center gap-2 text-sm text-zinc-400">
                <span
                  v-for="item in referenceItems"
                  :key="`reference-header:${item.label}`"
                  class="inline-flex items-baseline gap-1">
                  <span class="text-zinc-500">{{ item.label }}:</span>
                  <span>{{ item.value }}</span>
                </span>
              </div>
            </div>
            <div v-if="hasMineralImages" class="space-y-3">
              <div class="text-xs uppercase tracking-[0.08em] text-zinc-500">Изображения</div>
              <LabViewerPreviewButton
                :src="mineralCoverImage?.previewSrc || mineralCoverImage?.src || ''"
                :alt="mineralCoverImage?.alt || mineralName"
                :label="mineralCoverImage?.title || 'Открыть изображение'"
                button-class="w-full"
                frame-class="max-h-80"
                image-class="h-full w-full object-contain"
                @preview="openMineralImageViewer(activeImageIndex)" />
              <LabViewerImageThumbnails
                :items="
                  mineralImageItems.slice(0, 8).map(item => ({
                    key: item.key,
                    src: item.src,
                    thumbnailSrc: item.previewSrc,
                    alt: item.alt
                  }))
                "
                :active-index="activeImageIndex"
                @select="activeImageIndex = $event" />
              <div class="text-xs text-zinc-500">{{ mineralImageItems.length }} изображений</div>
            </div>
          </div>
          <div class="space-y-3">
            <div
              v-for="item in chemistryDisplayItems"
              :key="`chemistry:${item.label}`"
              class="min-w-0">
              <div class="text-[10px] uppercase tracking-[0.08em] text-zinc-500">{{ item.label }}</div>
              <LabViewerLaTex
                v-if="item.isFormula"
                :formula="item.latex"
                class="mt-1 text-base leading-6 text-zinc-100" />
              <div v-else class="mt-1 wrap-break-word text-base leading-6 text-zinc-100">{{ item.value }}</div>
            </div>
            <div v-if="chemistryElementChips.length" class="min-w-0">
              <div class="text-[10px] uppercase tracking-[0.08em] text-zinc-500">Chemistry elements</div>
              <div class="mt-2 flex flex-wrap gap-2">
                <template v-for="item in chemistryElementChips" :key="`chemistry:${item.symbol}`">
                  <NuxtLink
                    v-if="item.to"
                    :to="item.to"
                    class="inline-flex border border-zinc-700 px-2 py-1 text-sm text-zinc-100 hover:ring-1 hover:ring-zinc-500/70">
                    {{ item.symbol }}
                  </NuxtLink>
                  <span v-else class="inline-flex border border-zinc-700 px-2 py-1 text-sm text-zinc-100">
                    {{ item.symbol }}
                  </span>
                </template>
              </div>
            </div>
            <div v-if="valenceElementChips.length" class="min-w-0">
              <div class="text-[10px] uppercase tracking-[0.08em] text-zinc-500">Valence elements</div>
              <div class="mt-2 flex flex-wrap gap-2">
                <template v-for="item in valenceElementChips" :key="`valence:${item.symbol}`">
                  <NuxtLink
                    v-if="item.to"
                    :to="item.to"
                    class="inline-flex items-center border border-zinc-700 px-2 py-1 text-sm text-zinc-100 hover:ring-1 hover:ring-zinc-500/70">
                    <LabViewerLaTex v-if="item.isFormula" :formula="item.latex" class="text-sm text-zinc-100" />
                    <span v-else>{{ item.symbol }}</span>
                  </NuxtLink>
                  <span v-else class="inline-flex items-center border border-zinc-700 px-2 py-1 text-sm text-zinc-100">
                    <LabViewerLaTex v-if="item.isFormula" :formula="item.latex" class="text-sm text-zinc-100" />
                    <span v-else>{{ item.symbol }}</span>
                  </span>
                </template>
              </div>
            </div>
            <div class="text-xs uppercase tracking-[0.08em] text-zinc-500">Классификация</div>
            <div class="grid gap-3 grid-cols-1 min-[720px]:grid-cols-2 min-[1400px]:grid-cols-4">
              <div
                v-for="item in classificationPrimaryItems"
                :key="`classification-primary:${item.label}`"
                class="min-w-0">
                <div class="text-[10px] uppercase tracking-[0.08em] text-zinc-500">{{ item.label }}</div>
                <div class="mt-1 wrap-break-word text-sm leading-5 text-zinc-100">{{ item.value }}</div>
              </div>
            </div>
            <div
              v-for="item in classificationSecondaryItems"
              :key="`classification-secondary:${item.label}`"
              class="min-w-0">
              <div class="text-[10px] uppercase tracking-[0.08em] text-zinc-500">{{ item.label }}</div>
              <div class="mt-1 wrap-break-word text-sm leading-5 text-zinc-100">{{ item.value }}</div>
            </div>
            <div v-if="notesItems.length" class="space-y-3">
              <div class="text-xs uppercase tracking-[0.08em] text-zinc-500">Status notes</div>
              <div
                v-for="item in notesItems"
                :key="`notes:${item.label}`"
                class="min-w-0">
                <div class="mt-1 wrap-break-word text-sm leading-6 text-zinc-100">{{ item.value }}</div>
              </div>
            </div>
            <div
              v-for="item in originItems"
              :key="`origin:${item.label}`"
              class="min-w-0">
              <div class="text-[10px] uppercase tracking-[0.08em] text-zinc-500">{{ item.label }}</div>
              <div class="mt-1 wrap-break-word text-sm leading-5 text-zinc-100">{{ item.value }}</div>
            </div>
            <div v-if="parageneticModeItems.length" class="min-w-0">
              <div class="text-[10px] uppercase tracking-[0.08em] text-zinc-500">Paragenetic modes</div>
              <div class="mt-2 flex flex-wrap gap-2">
                <span
                  v-for="mode in parageneticModeItems"
                  :key="mode"
                  class="inline-flex border border-zinc-700 px-2 py-1 text-xs text-zinc-100">
                  {{ mode }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
    <LabViewerImage
      v-model="isImageViewerOpen"
      :items="
        mineralImageItems.map(item => ({
          src: item.src,
          title: item.title,
          alt: item.alt
        }))
      "
      :initial-index="activeImageIndex"
      :title="mineralName"
      @active-index-change="setMineralViewerIndex" />
  </div>
</template>
