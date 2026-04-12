<script setup lang="ts">
  const title = 'Минералы'
  usePageSeo({
    title,
    description: 'Каталог минералов с поиском, фильтрами по химическим элементам.'
  })
  const DEFAULT_LIMIT = 30
  const LIMIT_OPTIONS = [30, 60, 100]
  const SEARCH_DEBOUNCE_MS = 350
  const route = useRoute()
  const router = useRouter()
  const { t } = useInterfacePreferences()
  const { data: periodicTableElementsData, error: chemistryElementsError } = await useChemistryElements()
  if (chemistryElementsError.value) {
    throw createError({ statusCode: 500, statusMessage: 'Не удалось загрузить элементы' })
  }
  const isMediumViewport = ref(false)
  const isExtraLargeViewport = ref(false)
  const isElementsFilterOpen = ref(false)
  const preferencesStore = useUiPreferencesStore()
  preferencesStore.restorePersisted()
  const state = reactive(toRefs(preferencesStore.mineralsFilters))
  const periodicTableElements = computed(() => periodicTableElementsData.value || [])
  const crystalSystemOptions: Array<{ value: MineralCrystalSystem; label: string }> = [
    { value: 'cubic', label: 'minerals.crystal_system.cubic' },
    { value: 'hexagonal', label: 'minerals.crystal_system.hexagonal' },
    { value: 'monoclinic', label: 'minerals.crystal_system.monoclinic' },
    { value: 'orthorhombic', label: 'minerals.crystal_system.orthorhombic' },
    { value: 'tetragonal', label: 'minerals.crystal_system.tetragonal' },
    { value: 'triclinic', label: 'minerals.crystal_system.triclinic' },
    { value: 'unknown', label: 'minerals.crystal_system.unknown' }
  ]
  const mineralCrystalSystems = new Set<MineralCrystalSystem>(crystalSystemOptions.map(item => item.value))
  const hiddenMineralElementNumbers = computed(() =>
    !hasChemistryAvailability.value ?
      periodicTableElements.value
        .filter(element => {
          if (element.number === 90 || element.number === 92) return false
          if (element.category === 'noble gas' || element.category === 'unknown, predicted to be noble gas') return true
          return element.ypos === 10
        })
        .map(element => element.number)
    : periodicTableElements.value
        .filter(element => {
          const isSelected = Boolean(selectedBucketByElement.value[element.symbol])
          if (isSelected) return false
          if (element.number === 90 || element.number === 92) return false
          if (element.category === 'noble gas' || element.category === 'unknown, predicted to be noble gas') return true
          if (element.ypos === 10) return true
          return Number(chemistryAvailability.value[element.symbol] || 0) <= 0
        })
        .map(element => element.number)
  )
  const elementOrder = computed(
    () => new Map(periodicTableElements.value.map((element, index) => [element.symbol, index]))
  )
  const allElementSymbols = computed(() =>
    periodicTableElements.value
      .slice()
      .sort((left, right) => left.number - right.number)
      .map(element => element.symbol)
  )
  const compareElementOrder = (left: string, right: string) => {
    const leftOrder = elementOrder.value.get(left) ?? Number.MAX_SAFE_INTEGER
    const rightOrder = elementOrder.value.get(right) ?? Number.MAX_SAFE_INTEGER
    return leftOrder - rightOrder
  }
  const normalizeElementSymbol = (value: unknown) => normalizeMineralElementSymbol(value)
  const sortElements = (values: string[]) => sortMineralElementSymbols(values, compareElementOrder)
  const activeChemistryBucket = ref<ChemistryBucket>('all')
  const applyingRouteState = ref(false)
  let mediumViewportQuery: MediaQueryList | null = null
  let extraLargeViewportQuery: MediaQueryList | null = null
  const syncViewportState = () => {
    isMediumViewport.value = Boolean(mediumViewportQuery?.matches)
    isExtraLargeViewport.value = Boolean(extraLargeViewportQuery?.matches)
  }
  const restoreRouteFromStore = ref(
    import.meta.client &&
      !hasMineralsRouteQueryValues(route.query as Record<string, unknown>) &&
      preferencesStore.hasActiveMineralsFilters()
  )
  const applyRouteState = (query: Record<string, unknown>) => {
    const next = readMineralsRouteState(query, {
      defaultLimit: DEFAULT_LIMIT,
      allowedCrystalSystems: mineralCrystalSystems,
      compareElementOrder
    })
    applyingRouteState.value = true
    state.q = next.q
    state.sort = next.sort
    state.limit = next.limit
    state.offset = next.offset
    state.imageFilter = next.imageFilter
    state.crystalSystems = next.crystalSystems
    state.crystalSystemMode = next.crystalSystemMode
    state.chemistryAll = next.chemistryAll
    state.chemistryAny = next.chemistryAny
    state.chemistryNone = next.chemistryNone
    queueMicrotask(() => {
      applyingRouteState.value = false
    })
  }
  watch(
    () => route.query,
    query => {
      if (restoreRouteFromStore.value && !hasMineralsRouteQueryValues(query as Record<string, unknown>)) return
      applyRouteState(query as Record<string, unknown>)
    },
    { immediate: true, deep: true }
  )
  const routeMatchesState = () => {
    const current = readMineralsRouteState(route.query as Record<string, unknown>, {
      defaultLimit: DEFAULT_LIMIT,
      allowedCrystalSystems: mineralCrystalSystems,
      compareElementOrder
    })
    return (
      current.q === state.q.trim() &&
      current.sort === state.sort &&
      current.limit === state.limit &&
      current.offset === state.offset &&
      current.imageFilter === state.imageFilter &&
      JSON.stringify(current.crystalSystems) === JSON.stringify(state.crystalSystems) &&
      current.crystalSystemMode === state.crystalSystemMode &&
      JSON.stringify(current.chemistryAll) === JSON.stringify(sortElements(state.chemistryAll)) &&
      JSON.stringify(current.chemistryAny) === JSON.stringify(sortElements(state.chemistryAny)) &&
      JSON.stringify(current.chemistryNone) === JSON.stringify(sortElements(state.chemistryNone))
    )
  }
  const buildRouteQuery = () =>
    buildMineralsRouteQuery(
      {
        q: state.q,
        sort: state.sort,
        limit: state.limit,
        offset: state.offset,
        imageFilter: state.imageFilter,
        crystalSystems: state.crystalSystems,
        crystalSystemMode: state.crystalSystemMode,
        chemistryAll: state.chemistryAll,
        chemistryAny: state.chemistryAny,
        chemistryNone: state.chemistryNone
      },
      {
        defaultLimit: DEFAULT_LIMIT,
        sortElements
      }
    )
  const syncRouteFromState = async () => {
    if (applyingRouteState.value || routeMatchesState()) return
    await router.replace({
      path: '/edu/chemistry/minerals',
      query: buildRouteQuery()
    })
  }
  const scheduleSearchSync = debounce(() => {
    void syncRouteFromState()
  }, SEARCH_DEBOUNCE_MS)
  const setPageOffset = async (offset: number) => {
    state.offset = Math.max(0, offset)
    await syncRouteFromState()
  }
  const onSearchInput = (value: string) => {
    state.q = value
    state.offset = 0
    scheduleSearchSync()
  }
  const onSortChange = async (value: string) => {
    state.sort = value === 'name_desc' ? 'name_desc' : 'name_asc'
    state.offset = 0
    await syncRouteFromState()
  }
  const onLimitChange = async (value: string) => {
    const nextLimit = Number.parseInt(String(value || ''), 10)
    state.limit = Number.isFinite(nextLimit) && nextLimit > 0 ? Math.min(nextLimit, 100) : DEFAULT_LIMIT
    state.offset = 0
    await syncRouteFromState()
  }
  const onImageFilterChange = async (value: MineralImageFilter) => {
    state.imageFilter = value
    state.offset = 0
    await syncRouteFromState()
  }
  const removeFromAllBuckets = (symbol: string) => {
    const normalizedSymbol = normalizeElementSymbol(symbol)
    state.chemistryAll = state.chemistryAll.filter((item: string) => normalizeElementSymbol(item) !== normalizedSymbol)
    state.chemistryAny = state.chemistryAny.filter((item: string) => normalizeElementSymbol(item) !== normalizedSymbol)
    state.chemistryNone = state.chemistryNone.filter(
      (item: string) => normalizeElementSymbol(item) !== normalizedSymbol
    )
  }
  const selectedBucketByElement = computed<Record<string, ChemistryBucket | undefined>>(() => {
    const out: Record<string, ChemistryBucket | undefined> = {}
    for (const symbol of state.chemistryAll.map(normalizeElementSymbol)) out[symbol] = 'all'
    for (const symbol of state.chemistryAny.map(normalizeElementSymbol)) out[symbol] = 'any'
    for (const symbol of state.chemistryNone.map(normalizeElementSymbol)) out[symbol] = 'none'
    return out
  })
  const addToBucket = (bucket: ChemistryBucket, symbol: string) => {
    const normalizedSymbol = normalizeElementSymbol(symbol)
    if (!normalizedSymbol) return
    removeFromAllBuckets(normalizedSymbol)
    if (bucket === 'all') state.chemistryAll = sortElements([...state.chemistryAll, normalizedSymbol])
    if (bucket === 'any') state.chemistryAny = sortElements([...state.chemistryAny, normalizedSymbol])
    if (bucket === 'none') state.chemistryNone = sortElements([...state.chemistryNone, normalizedSymbol])
  }
  const toggleElementSelection = async (symbol: string) => {
    const normalizedSymbol = normalizeElementSymbol(symbol)
    if (!normalizedSymbol) return
    const currentBucket = selectedBucketByElement.value[normalizedSymbol]
    if (currentBucket === activeChemistryBucket.value) {
      removeFromAllBuckets(normalizedSymbol)
    } else {
      addToBucket(activeChemistryBucket.value, normalizedSymbol)
    }
    state.offset = 0
    await syncRouteFromState()
  }
  const clearChemistry = async () => {
    state.chemistryAll = []
    state.chemistryAny = []
    state.chemistryNone = []
    state.offset = 0
    await syncRouteFromState()
  }
  const excludeAllNonSelected = async () => {
    const selected = new Set([...state.chemistryAll, ...state.chemistryAny, ...state.chemistryNone])
    state.chemistryNone = allElementSymbols.value.filter(symbol => !selected.has(symbol))
    state.offset = 0
    await syncRouteFromState()
  }
  const clearAllFilters = async () => {
    state.q = ''
    state.sort = 'name_asc'
    state.limit = DEFAULT_LIMIT
    state.offset = 0
    state.imageFilter = 'any'
    state.crystalSystems = []
    state.crystalSystemMode = 'any'
    state.chemistryAll = []
    state.chemistryAny = []
    state.chemistryNone = []
    await syncRouteFromState()
  }
  const requestParams = computed(() =>
    readMineralsRouteState(route.query as Record<string, unknown>, {
      defaultLimit: DEFAULT_LIMIT,
      allowedCrystalSystems: mineralCrystalSystems,
      compareElementOrder
    })
  )
  const visibleLimitOptions = computed(() =>
    Array.from(new Set([...LIMIT_OPTIONS, state.limit]))
      .filter(value => value > 0 && value <= 100)
      .sort((left, right) => left - right)
  )
  const sortOptions: SelectOptionInput[] = [
    { value: 'name_asc', label: 'Имя А-Я' },
    { value: 'name_desc', label: 'Имя Я-А' }
  ]
  const imageFilterOptions: Array<{ value: MineralImageFilter; label: string }> = [
    { value: 'without', label: 'без' },
    { value: 'any', label: 'без/с' },
    { value: 'with', label: 'с' }
  ]
  const limitOptions = computed<SelectOptionInput[]>(() =>
    visibleLimitOptions.value.map(limitValue => ({
      value: String(limitValue),
      label: String(limitValue)
    }))
  )
  const { data, pending, error, refresh } = await useAsyncData(
    () => `minerals:list:${JSON.stringify(requestParams.value)}`,
    async () => {
      const response = await getMineralsList(requestParams.value)
      return response.data
    },
    {
      watch: [() => route.fullPath],
      default: () => ({
        items: [],
        meta: {
          limit: DEFAULT_LIMIT,
          offset: 0,
          total: 0
        },
        facets: {
          chemistryElementsAvailable: {}
        }
      })
    }
  )
  const minerals = computed(() => data.value?.items || [])
  const meta = computed(() => data.value?.meta || { limit: DEFAULT_LIMIT, offset: 0, total: 0 })
  const chemistryAvailability = computed<Record<string, number>>(
    () => data.value?.facets?.chemistryElementsAvailable || {}
  )
  const hasChemistryAvailability = computed(() => Object.keys(chemistryAvailability.value).length > 0)
  const totalPages = computed(() => {
    if (!meta.value.total || !meta.value.limit) return 1
    return Math.max(1, Math.ceil(meta.value.total / meta.value.limit))
  })
  const currentPage = computed(() => Math.floor(meta.value.offset / meta.value.limit) + 1)
  const hasPreviousPage = computed(() => meta.value.offset > 0)
  const hasNextPage = computed(() => meta.value.offset + meta.value.limit < meta.value.total)
  const errorMessage = computed(() => {
    const value = error.value as any
    return value?.data?.message || value?.message || 'Не удалось загрузить минералы.'
  })
  const mineralListItemName = (mineral: MineralListItem) => String(mineral.mineral_name || '').trim()
  const mineralListItemNameLatex = (mineral: MineralListItem) => mineralFormulaToLatex(mineral.mineral_name)
  const showMineralListItemNameAsFormula = (mineral: MineralListItem) => isMineralFormulaLike(mineral.mineral_name)
  const chemistryBucketOrder = ['all', 'any', 'none'] as const
  const chemistryBucketMeta: Record<ChemistryBucket, ChemistryBucketMeta> = {
    all: {
      title: 'ВСЕ',
      description: 'Минерал обязан содержать все выбранные элементы.',
      accentClass: 'text-(--lab-warning)',
      borderClass: 'border-l-(--lab-warning)',
      activeClass:
        'border-[color-mix(in_srgb,var(--lab-warning)_42%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-warning)_16%,var(--lab-bg-surface))] text-[color-mix(in_srgb,var(--lab-warning)_88%,black_18%)]',
      activeTextClass: 'text-[color-mix(in_srgb,var(--lab-warning)_88%,black_18%)]',
      dotClass: 'bg-(--lab-warning)'
    },
    any: {
      title: 'ЛЮБОЙ',
      description: 'Минерал должен содержать хотя бы один выбранный элемент.',
      accentClass: 'text-(--lab-accent)',
      borderClass: 'border-l-(--lab-accent)',
      activeClass:
        'border-[color-mix(in_srgb,var(--lab-accent)_42%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-accent)_16%,var(--lab-bg-surface))] text-[color-mix(in_srgb,var(--lab-accent)_84%,black_18%)]',
      activeTextClass: 'text-[color-mix(in_srgb,var(--lab-accent)_84%,black_18%)]',
      dotClass: 'bg-(--lab-accent)'
    },
    none: {
      title: 'ИСКЛЮЧЕНИЕ',
      description: 'Минерал не должен содержать ни одного выбранного элемента.',
      accentClass: 'text-(--lab-danger)',
      borderClass: 'border-l-(--lab-danger)',
      activeClass:
        'border-[color-mix(in_srgb,var(--lab-danger)_42%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-danger)_14%,var(--lab-bg-surface))] text-[color-mix(in_srgb,var(--lab-danger)_84%,black_22%)]',
      activeTextClass: 'text-[color-mix(in_srgb,var(--lab-danger)_84%,black_22%)]',
      dotClass: 'bg-(--lab-danger)'
    }
  }
  const chemistryBuckets = computed(() =>
    chemistryBucketOrder.map(key => ({
      key,
      ...chemistryBucketMeta[key],
      values:
        key === 'all' ? state.chemistryAll
        : key === 'any' ? state.chemistryAny
        : state.chemistryNone
    }))
  )
  const chemistryBucketClass = (bucket: ChemistryBucket) => {
    const baseClass = `lab-focus w-full border border-l-2 px-0 py-0 text-left transition ${chemistryBucketMeta[bucket].borderClass}`
    if (activeChemistryBucket.value === bucket) {
      return `${baseClass} ${chemistryBucketMeta[bucket].activeClass}`
    }
    if (bucket === 'all') {
      return `${baseClass} bg-(--lab-bg-surface) hover:border-[color-mix(in_srgb,var(--lab-warning)_38%,var(--lab-border))] hover:bg-[color-mix(in_srgb,var(--lab-warning)_16%,var(--lab-bg-surface))] hover:text-[color-mix(in_srgb,var(--lab-warning)_88%,black_18%)]`
    }
    if (bucket === 'any') {
      return `${baseClass} bg-(--lab-bg-surface) hover:border-[color-mix(in_srgb,var(--lab-accent)_38%,var(--lab-border))] hover:bg-[color-mix(in_srgb,var(--lab-accent)_16%,var(--lab-bg-surface))] hover:text-[color-mix(in_srgb,var(--lab-accent)_84%,black_18%)]`
    }
    return `${baseClass} bg-(--lab-bg-surface) hover:border-[color-mix(in_srgb,var(--lab-danger)_38%,var(--lab-border))] hover:bg-[color-mix(in_srgb,var(--lab-danger)_14%,var(--lab-bg-surface))] hover:text-[color-mix(in_srgb,var(--lab-danger)_84%,black_22%)]`
  }
  const chemistryBucketText = (values: string[]) => values.join(', ')
  const selectedElementNumbers = computed(() =>
    periodicTableElements.value
      .filter(element => Boolean(selectedBucketByElement.value[element.symbol]))
      .map(element => element.number)
  )
  const hiddenElementNumbers = computed(() =>
    periodicTableElements.value.filter(element => element.xpos === 18).map(element => element.number)
  )
  const selectionToneByNumber = computed(() =>
    periodicTableElements.value.reduce<Partial<Record<number, ChemistryBucket>>>((acc, element) => {
      const bucket = selectedBucketByElement.value[element.symbol]
      if (bucket) acc[element.number] = bucket
      return acc
    }, {})
  )
  const chemistryAvailabilityByNumber = computed(() =>
    periodicTableElements.value.reduce<Partial<Record<number, number>>>((acc, element) => {
      if (Object.prototype.hasOwnProperty.call(chemistryAvailability.value, element.symbol)) {
        acc[element.number] = Number(chemistryAvailability.value[element.symbol] || 0)
      }
      return acc
    }, {})
  )
  const useCompactMineralsTable = computed(() => !isMediumViewport.value)
  const onPeriodicElementClick = async (element: PeriodicTableElement) => {
    await toggleElementSelection(element.symbol)
  }
  const crystalSystemLabel = (value: MineralCrystalSystem) =>
    t(
      (crystalSystemOptions.find(item => item.value === value)?.label ||
        'minerals.crystal_system.unknown') as InterfaceMessageKey
    )
  const crystalSystemHelpText =
    'Кристаллическая система описывает симметрию решётки минерала. Кубическая: наиболее симметричная. Тетрагональная и гексагональная: оси с регулярной симметрией. Ромбическая, моноклинная и триклинная: менее симметричные. Неизвестно: система не указана или не определена.'
  const toggleCrystalSystem = async (value: MineralCrystalSystem) => {
    if (state.crystalSystems.includes(value)) {
      state.crystalSystems = state.crystalSystems.filter(item => item !== value)
    } else {
      state.crystalSystems = [...state.crystalSystems, value]
    }
    state.offset = 0
    await syncRouteFromState()
  }
  const onCrystalSystemModeChange = async (value: boolean) => {
    state.crystalSystemMode = value ? 'all' : 'any'
    state.offset = 0
    await syncRouteFromState()
  }

  onMounted(async () => {
    mediumViewportQuery = window.matchMedia('(min-width: 768px)')
    extraLargeViewportQuery = window.matchMedia('(min-width: 1280px)')
    syncViewportState()
    mediumViewportQuery.addEventListener('change', syncViewportState)
    extraLargeViewportQuery.addEventListener('change', syncViewportState)
    if (restoreRouteFromStore.value) {
      restoreRouteFromStore.value = false
      await router.replace({
        path: '/edu/chemistry/minerals',
        query: buildRouteQuery()
      })
    }
  })
  onBeforeUnmount(() => {
    mediumViewportQuery?.removeEventListener('change', syncViewportState)
    extraLargeViewportQuery?.removeEventListener('change', syncViewportState)
  })
</script>
<template>
  <div>
    <LabNavHeader
      :title
      :breadcrumb-items="[
        { label: 'Вики', to: '/edu' },
        { label: 'Химия', to: '/edu/chemistry' },
        { label: title, current: true }
      ]" />
    <section class="overflow-hidden">
      <div class="flex flex-wrap min-w-0 gap-3 border-b p-4">
        <LabField label="Поиск минерала" for-id="minerals-search">
          <LabBaseInput
            id="minerals-search"
            name="minerals_search"
            :model-value="state.q"
            placeholder="Название минерала"
            class="max-w-fit"
            @update:model-value="onSearchInput" />
        </LabField>
        <LabField label="Сортировка" for-id="minerals-sort">
          <LabBaseSelect
            id="minerals-sort"
            name="minerals_sort"
            :model-value="state.sort"
            :options="sortOptions"
            class="max-w-fit"
            @update:model-value="onSortChange" />
        </LabField>
        <LabField label="На страницу" for-id="minerals-limit">
          <LabBaseSelect
            id="minerals-limit"
            name="minerals_limit"
            :model-value="String(state.limit)"
            :options="limitOptions"
            class="max-w-fit"
            @update:model-value="onLimitChange" />
        </LabField>
        <div class="min-w-0 flex flex-wrap items-end justify-between gap-2 text-sm text-zinc-400 xl:col-span-12">
          <LabField label="Фото">
            <div class="flex min-w-0 items-center gap-1">
              <button
                v-for="option in imageFilterOptions"
                :key="option.value"
                type="button"
                class="lab-focus px-2 py-1 text-xs transition"
                :class="
                  state.imageFilter === option.value ?
                    'ring-1 ring-(--lab-accent) text-(--lab-text-primary)'
                  : 'text-(--lab-text-secondary) hover:ring-1 hover:ring-(--lab-border) hover:text-(--lab-text-primary)'
                "
                :aria-pressed="state.imageFilter === option.value"
                @click="onImageFilterChange(option.value)">
                {{ option.label }}
              </button>
            </div>
          </LabField>
          <LabBaseButton variant="secondary" size="sm" label="Сбросить все фильтры" @click="clearAllFilters" />
        </div>
        <LabSpoiler
          v-model="isElementsFilterOpen"
          label="Химический фильтр"
          container-class="min-w-0 max-w-6xl md:col-span-12">
          <div class="flex flex-wrap items-center gap-2">
            <p class="text-sm leading-6 text-zinc-400">Выберите режим и отмечайте элементы прямо в таблице.</p>
            <div class="flex flex-wrap items-center gap-2">
              <LabBaseButton variant="secondary" size="sm" label="Очистить химию" @click="clearChemistry" />
              <LabBaseButton variant="danger" size="sm" label="Исключить остальные" @click="excludeAllNonSelected" />
            </div>
          </div>
          <div class="grid grid-cols-1 gap-px overflow-hidden border sm:grid-cols-3">
            <button
              v-for="bucket in chemistryBuckets"
              :key="bucket.key"
              type="button"
              :aria-pressed="activeChemistryBucket === bucket.key"
              :class="chemistryBucketClass(bucket.key)"
              @click="activeChemistryBucket = bucket.key">
              <div class="flex items-center justify-between gap-2 px-2 py-1">
                <LabBaseTooltip
                  :text="bucket.description"
                  underline-trigger
                  :trigger-class="activeChemistryBucket === bucket.key ? bucket.accentClass : 'text-current'">
                  <template #trigger>
                    <span class="text-xs font-semibold tracking-[0.14em]">
                      {{ bucket.title }}
                    </span>
                  </template>
                </LabBaseTooltip>
                <span
                  class="uppercase tracking-[0.12em] opacity-65"
                  :class="activeChemistryBucket === bucket.key ? bucket.accentClass : ''">
                  {{ activeChemistryBucket === bucket.key ? 'выбран' : 'выбрать' }}
                </span>
              </div>
              <div
                class="min-h-4 px-2 pb-1 text-xs leading-4 wrap-break-word"
                :class="activeChemistryBucket === bucket.key ? bucket.activeTextClass : 'opacity-70'">
                {{ chemistryBucketText(bucket.values) || '\u00A0' }}
              </div>
            </button>
          </div>
          <div class="min-w-0 overflow-x-auto">
            <ChemistryPeriodicTableGrid
              layout="minerals"
              :elements="periodicTableElements"
              :hidden-element-numbers="hiddenElementNumbers"
              :highlighted-element-numbers="selectedElementNumbers"
              :selection-tone-by-number="selectionToneByNumber"
              :availability-by-number="chemistryAvailabilityByNumber"
              :disable-unavailable="hasChemistryAvailability"
              use-theme-card-colors
              compact-cell-height="1.7rem"
              compact-cell-height-wide="2.5rem"
              :preserve-compact-cell-width="useCompactMineralsTable"
              :compact="useCompactMineralsTable"
              interactive
              :show-series-labels="false"
              @element-click="onPeriodicElementClick" />
          </div>
          <div class="space-y-2">
            <div class="flex flex-wrap items-center gap-2">
              <div class="inline-flex min-w-0 shrink-0 items-center gap-2 whitespace-nowrap">
                <div class="text-sm text-zinc-100">Кристаллическая система</div>
                <LabHelpTooltip :text="crystalSystemHelpText" />
              </div>
              <LabBaseSwitch
                :model-value="state.crystalSystemMode === 'all'"
                label="совпадения"
                false-label="или"
                true-label="и"
                tone="amber"
                @update:model-value="onCrystalSystemModeChange" />
            </div>
            <div class="min-w-0 overflow-x-auto">
              <div class="flex w-max min-w-full flex-nowrap gap-2">
                <button
                  v-for="item in crystalSystemOptions"
                  :key="item.value"
                  type="button"
                  class="lab-focus shrink-0 border px-3 py-2 text-left text-sm transition"
                  :class="
                    state.crystalSystems.includes(item.value) ?
                      'border-zinc-500 bg-zinc-900 text-zinc-100 ring-1 ring-inset ring-zinc-500'
                    : 'border-zinc-800 bg-zinc-950/90 text-zinc-400 hover:border-zinc-700 hover:text-zinc-200'
                  "
                  @click="toggleCrystalSystem(item.value)">
                  {{ crystalSystemLabel(item.value) }}
                </button>
              </div>
            </div>
          </div>
        </LabSpoiler>
      </div>
      <div v-if="error" class="border-b bg-rose-950/20 p-5">
        <LabErrorMessage :text="errorMessage" error-class="text-sm" />
        <div class="mt-3">
          <LabBaseButton variant="secondary" size="sm" label="Повторить" @click="refresh" />
        </div>
      </div>
      <div v-else-if="pending && !minerals.length" class="flex min-h-72 items-center justify-center border-b p-6">
        <LabLoader size="md" variant="inline" label="Загружаем минералы..." />
      </div>
      <div v-else-if="!minerals.length" class="border-b p-6 text-sm leading-6 text-zinc-400">
        По текущим условиям ничего не найдено. Попробуйте ослабить поиск или очистить chemistry-фильтр.
      </div>
      <div v-else class="divide-y">
        <div class="text-xs text-zinc-100">Найдено: {{ meta.total }}</div>
        <div v-for="mineral in minerals" :key="mineral.database_id">
          <NuxtLink
            :to="`/edu/chemistry/minerals/${mineral.database_id}`"
            class="block px-4 py-2 transition hover:bg-white/10 hover:ring-1 hover:ring-inset hover:ring-zinc-600/70">
            <LabViewerLaTex
              v-if="showMineralListItemNameAsFormula(mineral)"
              :formula="mineralListItemNameLatex(mineral)"
              class="text-sm font-semibold text-zinc-100" />
            <div v-else class="text-sm font-semibold leading-5 text-zinc-100 wrap-break-word">
              {{ mineralListItemName(mineral) }}
            </div>
          </NuxtLink>
        </div>
      </div>
      <div class="flex flex-col gap-3 p-4 sm:flex-row sm:items-center sm:justify-between">
        <div class="text-sm text-zinc-400">
          Страница
          <span class="text-zinc-100">{{ currentPage }}</span>
          из
          <span class="text-zinc-100">{{ totalPages }}</span>
        </div>
        <div class="flex items-center gap-2">
          <LabBaseButton
            variant="secondary"
            size="sm"
            label="Назад"
            :disabled="!hasPreviousPage || pending"
            @click="setPageOffset(Math.max(0, meta.offset - meta.limit))" />
          <LabBaseButton
            variant="secondary"
            size="sm"
            label="Вперёд"
            :disabled="!hasNextPage || pending"
            @click="setPageOffset(meta.offset + meta.limit)" />
        </div>
      </div>
    </section>
  </div>
</template>
