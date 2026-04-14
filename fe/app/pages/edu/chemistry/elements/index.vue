<script setup lang="ts">
const title = 'Периодическая система химических элементов'
const { effectiveTheme } = useInterfacePreferences()
const { data: periodicTableElementsData, error: chemistryElementsError } = await useChemistryElements()
if (chemistryElementsError.value) {
  throw createError({
    statusCode: 500,
    statusMessage: 'Не удалось загрузить элементы'
  })
}
const pageHeaderRef = ref<HTMLElement | null>(null)
const pageHeaderHeight = ref(0)
const isDesktopViewport = ref(false)
const isMobileViewport = ref(false)
const searchQuery = ref('')
const selectedCategory = ref('')
const periodicTableElements = computed(() => periodicTableElementsData.value || [])
const periodicTableCategoryCounts = computed(() => buildPeriodicTableCategoryCounts(periodicTableElements.value))
let pageHeaderResizeObserver: ResizeObserver | null = null
let desktopMediaQuery: MediaQueryList | null = null
let mobileMediaQuery: MediaQueryList | null = null
let removeDesktopViewportListener: (() => void) | null = null
let removeMobileViewportListener: (() => void) | null = null
const tableShellStyle = computed(() => ({
  '--periodic-header-height': `${Math.ceil(pageHeaderHeight.value || 72)}px`
}))
const categoryOptions = computed(() => [
  { value: '', label: 'Все категории' },
  ...periodicTableCategoryCounts.value.map((entry) => ({
    value: entry.category,
    label: `${entry.label} · ${entry.count}`
  }))
])
const normalizedSearchQuery = computed(() => searchQuery.value.trim().toLowerCase())
const visibleElementIds = computed(() => {
  const query = normalizedSearchQuery.value
  if (!query) {
    return new Set(periodicTableElements.value.map((element) => element.number))
  }
  return new Set(
    periodicTableElements.value.filter((element) => element.searchText.includes(query)).map((element) => element.number)
  )
})
const selectedCategoryEntry = computed(() => {
  return periodicTableCategoryCounts.value.find((entry) => entry.category === selectedCategory.value) || null
})
const categoryPanelDescription = computed(() => {
  return (
    selectedCategoryEntry.value?.description ||
    'Выберите категорию, чтобы подсветить элементы в таблице и посмотреть краткое описание группы.'
  )
})
watch(isDesktopViewport, (isDesktop) => {
  if (isDesktop) {
    searchQuery.value = ''
  }
})
const measureHeader = () => {
  if (!pageHeaderRef.value) return
  pageHeaderHeight.value = pageHeaderRef.value.getBoundingClientRect().height || 72
}
onBeforeMount(() => {
  measureHeader()
})
const syncViewportFlags = () => {
  isDesktopViewport.value = Boolean(desktopMediaQuery?.matches)
  isMobileViewport.value = Boolean(mobileMediaQuery?.matches)
}
onMounted(() => {
  measureHeader()
  if (typeof window !== 'undefined') {
    desktopMediaQuery = window.matchMedia('(min-width: 1024px)')
    mobileMediaQuery = window.matchMedia('(max-width: 639px)')
    syncViewportFlags()
    const onDesktopChange = (event: MediaQueryListEvent) => {
      isDesktopViewport.value = event.matches
    }
    const onMobileChange = (event: MediaQueryListEvent) => {
      isMobileViewport.value = event.matches
    }
    desktopMediaQuery.addEventListener('change', onDesktopChange)
    mobileMediaQuery.addEventListener('change', onMobileChange)
    removeDesktopViewportListener = () => {
      desktopMediaQuery?.removeEventListener('change', onDesktopChange)
    }
    removeMobileViewportListener = () => {
      mobileMediaQuery?.removeEventListener('change', onMobileChange)
    }
  }
  if (typeof ResizeObserver !== 'undefined' && pageHeaderRef.value) {
    pageHeaderResizeObserver = new ResizeObserver(() => {
      measureHeader()
    })
    pageHeaderResizeObserver.observe(pageHeaderRef.value)
  }
  window.addEventListener('resize', measureHeader)
})
onBeforeUnmount(() => {
  pageHeaderResizeObserver?.disconnect()
  removeDesktopViewportListener?.()
  removeMobileViewportListener?.()
  window.removeEventListener('resize', measureHeader)
})
const categoryPanelToneClass = computed(() => {
  if (!selectedCategoryEntry.value) {
    return 'bg-(--lab-bg-surface-muted)'
  }
  return effectiveTheme.value === 'light' ? 'bg-(--lab-bg-surface)' : 'bg-(--lab-bg-surface-muted)'
})
const isElementDimmed = (number: number, category: string) => {
  if (!visibleElementIds.value.has(number)) return true
  if (!selectedCategory.value) return false
  return category !== selectedCategory.value
}
const dimmedElementNumbers = computed(() =>
  periodicTableElements.value
    .filter((element) => isElementDimmed(element.number, element.category))
    .map((element) => element.number)
)
const isElementHighlighted = (number: number, category: string) => {
  if (!visibleElementIds.value.has(number)) return false
  if (selectedCategory.value) return category === selectedCategory.value
  return Boolean(normalizedSearchQuery.value)
}
const highlightedElementNumbers = computed(() =>
  periodicTableElements.value
    .filter((element) => isElementHighlighted(element.number, element.category))
    .map((element) => element.number)
)
const onElementClick = (element: PeriodicTableElement) => {
  navigateTo(getPeriodicTableElementRoute(element))
}
usePageSeo({
  title,
  description: 'Периодическая система химических элементов упорядоченная по атомному номеру, периодам и группам.'
})
</script>
<template>
  <div>
    <div ref="pageHeaderRef">
      <LabNavHeader
        :title
        :breadcrumb-items="[
          { label: 'Вики', to: '/edu' },
          { label: 'Химия', to: '/edu/chemistry' },
          { label: title, current: true }
        ]"
      />
    </div>
    <section
      class="h-auto max-h-[calc(100dvh-var(--periodic-header-height,72px)-1rem)] min-h-0 overflow-auto sm:max-h-[calc(100dvh-var(--periodic-header-height,0px)-1.5rem)] sm:bg-(--lab-bg-surface-muted)"
      :style="tableShellStyle"
    >
      <div :class="isMobileViewport ? 'p-1' : 'p-2'">
        <div
          class="sticky top-0 z-20 border-b bg-(--lab-bg-overlay) p-2 lg:hidden"
          :class="isMobileViewport ? '-mx-1 px-1.5' : '-mx-2'"
        >
          <div class="flex flex-wrap items-center gap-2">
            <div class="min-w-0 flex-1 basis-full">
              <LabBaseInput
                id="elements-search-mobile"
                v-model="searchQuery"
                name="elementsSearchMobile"
                placeholder="Поиск"
                aria-label="Поиск элемента по периодической системе химических элементов"
                input-class="bg-(--lab-bg-control) text-(--lab-text-primary) placeholder:text-(--lab-text-soft) !min-h-10 px-3 py-2 text-xs"
              />
            </div>
            <div class="min-w-0 basis-full sm:basis-56">
              <LabBaseSelect
                v-model="selectedCategory"
                :options="categoryOptions"
                aria-label="Выбор категории элементов"
                select-class="bg-(--lab-bg-control) text-(--lab-text-primary) !min-h-10 w-full px-3 py-2 text-xs"
              />
            </div>
          </div>
        </div>
        <ChemistryPeriodicTableGrid
          :elements="periodicTableElements"
          :highlighted-element-numbers="highlightedElementNumbers"
          :dimmed-element-numbers="dimmedElementNumbers"
          use-theme-card-colors
          interactive
          @element-click="onElementClick"
        >
          <template #grid-overlay>
            <div
              class="relative z-10 col-[3/span_10] row-[1/span_3] hidden border p-2 lg:block"
              :class="categoryPanelToneClass"
            >
              <div class="flex items-start justify-between gap-2">
                <div class="min-w-44 shrink-0">
                  <LabBaseSelect
                    v-model="selectedCategory"
                    :options="categoryOptions"
                    aria-label="Выбор категории элементов"
                    select-class="bg-(--lab-bg-control) text-(--lab-text-primary) text-xs"
                  />
                </div>
                <div class="flex items-center gap-2 text-xs text-(--lab-text-secondary) xl:hidden">
                  <span
                    class="h-2.5 w-2.5 rounded-full"
                    :style="{
                      backgroundColor: selectedCategoryEntry?.color || 'var(--lab-text-muted)'
                    }"
                  />
                  <LabHelpTooltip class="" :text="categoryPanelDescription" />
                </div>
              </div>
              <p class="hidden px-2 text-sm leading-6 text-(--lab-text-secondary) xl:block">
                {{ categoryPanelDescription }}
              </p>
            </div>
          </template>
        </ChemistryPeriodicTableGrid>
      </div>
    </section>
  </div>
</template>
