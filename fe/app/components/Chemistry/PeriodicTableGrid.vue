<script setup lang="ts">
  const { effectiveTheme } = useInterfacePreferences()
  const props = withDefaults(
    defineProps<{
      elements: PeriodicTableElement[]
      highlightedElementNumbers?: number[]
      selectionToneByNumber?: Partial<Record<number, 'all' | 'any' | 'none'>>
      availabilityByNumber?: Partial<Record<number, number>>
      hiddenElementNumbers?: number[]
      disableUnavailable?: boolean
      useThemeCardColors?: boolean
      preserveCompactCellWidth?: boolean
      compactScrollBreakpoint?: 'md' | 'lg'
      compactCellHeight?: string
      compactCellHeightWide?: string
      compact?: boolean
      interactive?: boolean
      showSeriesLabels?: boolean
      layout?: 'default' | 'minerals'
    }>(),
    {
      highlightedElementNumbers: () => [],
      selectionToneByNumber: () => ({}),
      availabilityByNumber: () => ({}),
      hiddenElementNumbers: () => [],
      disableUnavailable: false,
      useThemeCardColors: false,
      preserveCompactCellWidth: false,
      compactScrollBreakpoint: 'lg',
      compactCellHeight: '',
      compactCellHeightWide: '',
      compact: false,
      interactive: false,
      showSeriesLabels: true,
      layout: 'default'
    }
  )
  const emit = defineEmits<{
    elementClick: [element: PeriodicTableElement, event: MouseEvent]
  }>()
  const periodicTableGroups = computed(() => Math.max(...visibleElements.value.map(element => element.xpos), 1))
  const groupHeaders = computed(() => Array.from({ length: periodicTableGroups.value }, (_, index) => index + 1))
  const periodHeaders = computed(() =>
    visiblePeriodRows.value.map(row => {
      if (row === 9) return '6'
      if (row === 10) return '7'
      return String(row)
    })
  )
  const visiblePeriodRows = computed(() => {
    if (!isMineralsLayout.value) return [1, 2, 3, 4, 5, 6, 7, 9, 10]
    const rows = [1, 2, 3, 4, 5, 6]
    const hasRow7 = visibleElements.value.some(element => element.ypos === 7)
    const hasLanthanidesRow = visibleElements.value.some(element => element.ypos === 9)
    const hasActinidesRow = visibleElements.value.some(element => element.ypos === 10)
    if (hasRow7) rows.push(7)
    if (hasLanthanidesRow) rows.push(9)
    if (hasActinidesRow) rows.push(10)
    return rows
  })
  const highlightedSet = computed(() => new Set(props.highlightedElementNumbers))
  const selectionToneMap = computed(() => props.selectionToneByNumber || {})
  const hiddenSet = computed(() => new Set(props.hiddenElementNumbers))
  const availabilityMap = computed(() => props.availabilityByNumber || {})
  const unavailableSet = computed(() => {
    const out = new Set<number>()
    for (const [key, value] of Object.entries(availabilityMap.value)) {
      if (Number(value || 0) <= 0) {
        out.add(Number(key))
      }
    }
    return out
  })
  const hasAvailabilityState = computed(() => Object.keys(availabilityMap.value).length > 0)
  const visibleElements = computed(() =>
    props.elements
      .filter(element => !hiddenSet.value.has(element.number))
      .map(element => {
        if (!isMineralsLayout.value) return element
        if (element.number === 90) {
          return { ...element, xpos: 4, ypos: 10 }
        }
        if (element.number === 92) {
          return { ...element, xpos: 5, ypos: 10 }
        }
        return element
      })
  )
  const isMineralsLayout = computed(() => props.layout === 'minerals')
  const getDisplayRow = (row: number) => {
    const index = visiblePeriodRows.value.indexOf(row)
    return index >= 0 ? index + 1 : row
  }
  const getDisplayColumn = (column: number) => displayColumnBySourceColumn.value.get(column) ?? column
  const displayColumnBySourceColumn = computed(
    () => new Map<number, number>(groupHeaders.value.map((group, index) => [group, index + 1]))
  )
  const tableGridStyle = computed(() => ({
    display: 'grid',
    gridTemplateColumns: `repeat(${groupHeaders.value.length}, minmax(0, 1fr))`,
    gridTemplateRows: `repeat(${visiblePeriodRows.value.length}, minmax(0, var(--periodic-cell-height)))`
  }))
  const frameStyle = computed(() => ({
    '--periodic-gap': '0.25rem',
    '--periodic-highlight-reserve': '6px',
    '--periodic-mobile-cell-width': isMineralsLayout.value ? '2.15rem' : '2.85rem',
    '--periodic-mobile-cell-height': isMineralsLayout.value ? '1.85rem' : '2.2rem',
    '--periodic-compact-cell-height': props.compactCellHeight || undefined,
    '--periodic-compact-cell-height-wide': props.compactCellHeightWide || undefined,
    '--periodic-row-count': String(visiblePeriodRows.value.length),
    '--periodic-row-gap-count': String(Math.max(visiblePeriodRows.value.length - 1, 0)),
    '--periodic-column-count': String(groupHeaders.value.length),
    '--periodic-column-gap-count': String(Math.max(groupHeaders.value.length - 1, 0))
  }))
  const frameClass = computed(() => [
    'box-border shrink-0 p-(--periodic-highlight-reserve) [--periodic-cell-width:var(--periodic-mobile-cell-width)] [--periodic-cell-height:var(--periodic-mobile-cell-height)] [--periodic-series-gap-height:0rem] [--periodic-row-label-size:0px] [--periodic-column-label-size:0px] w-[calc(var(--periodic-column-count)*var(--periodic-cell-width)+var(--periodic-column-gap-count)*var(--periodic-gap))] h-[calc(var(--periodic-row-count)*var(--periodic-cell-height)+var(--periodic-series-gap-height)+var(--periodic-row-gap-count)*var(--periodic-gap))]',
    props.compact ?
      'sm:[--periodic-row-label-size:1rem] sm:[--periodic-column-label-size:1rem] sm:w-[calc(var(--periodic-row-label-size)+var(--periodic-gap)+var(--periodic-column-count)*var(--periodic-cell-width)+var(--periodic-column-gap-count)*var(--periodic-gap))] sm:h-[calc(var(--periodic-column-label-size)+var(--periodic-gap)+var(--periodic-row-count)*var(--periodic-cell-height)+var(--periodic-series-gap-height)+var(--periodic-row-gap-count)*var(--periodic-gap))] lg:[--periodic-gap:0.2rem] lg:[--periodic-compact-cell-width:1.65rem] lg:[--periodic-cell-width:var(--periodic-compact-cell-width)] lg:[--periodic-cell-height:var(--periodic-compact-cell-height,1.45rem)] xl:[--periodic-gap:0.24rem] xl:[--periodic-compact-cell-width:1.75rem] xl:[--periodic-cell-height:var(--periodic-compact-cell-height-wide,var(--periodic-compact-cell-height,1.55rem))] 2xl:[--periodic-cell-height:var(--periodic-compact-cell-height-wide,var(--periodic-compact-cell-height,1.55rem))]'
    : 'sm:[--periodic-row-label-size:1.5rem] sm:[--periodic-column-label-size:1.5rem] sm:w-[calc(var(--periodic-row-label-size)+var(--periodic-gap)+var(--periodic-column-count)*var(--periodic-cell-width)+var(--periodic-column-gap-count)*var(--periodic-gap))] sm:h-[calc(var(--periodic-column-label-size)+var(--periodic-gap)+var(--periodic-row-count)*var(--periodic-cell-height)+var(--periodic-series-gap-height)+var(--periodic-row-gap-count)*var(--periodic-gap))] lg:[--periodic-gap:0.35rem] lg:[--periodic-compact-cell-width:2.1rem] lg:[--periodic-cell-width:var(--periodic-compact-cell-width)] lg:[--periodic-cell-height:var(--periodic-compact-cell-height,1.8rem)] xl:[--periodic-gap:0.4rem] xl:[--periodic-compact-cell-width:2.2rem] 2xl:[--periodic-full-cell-height:4.75rem] 2xl:[--periodic-cell-height:var(--periodic-full-cell-height)]',
    props.preserveCompactCellWidth ?
      props.compactScrollBreakpoint === 'md' ?
        'md:w-[calc(var(--periodic-row-label-size)+var(--periodic-gap)+var(--periodic-column-count)*var(--periodic-cell-width)+var(--periodic-column-gap-count)*var(--periodic-gap))] lg:w-full'
      : 'lg:w-[calc(var(--periodic-row-label-size)+var(--periodic-gap)+var(--periodic-column-count)*var(--periodic-cell-width)+var(--periodic-column-gap-count)*var(--periodic-gap))]'
    : 'lg:min-h-0 lg:min-w-0 lg:w-full'
  ])
  const layoutStyle = computed(() => ({
    gridTemplateColumns: 'var(--periodic-row-label-size) minmax(0, 1fr)',
    gridTemplateRows: 'var(--periodic-column-label-size) minmax(0, 1fr)'
  }))
  const tableColumnHeaderStyle = computed(() => ({
    display: 'grid',
    gridTemplateColumns: `repeat(${groupHeaders.value.length}, minmax(0, 1fr))`
  }))
  const tableRowHeaderStyle = computed(() => ({
    display: 'grid',
    gridTemplateRows: `repeat(${visiblePeriodRows.value.length}, minmax(0, var(--periodic-cell-height)))`
  }))
  const elementSlotStyle = (xpos: number, ypos: number) => ({
    gridColumn: `${getDisplayColumn(xpos)} / span 1`,
    gridRow: `${getDisplayRow(ypos)} / span 1`
  })

  const colorHexToRgba = (hexColor: string, alpha: number) => {
    const normalized = String(hexColor || '').replace('#', '')
    if (!/^[0-9a-f]{6}$/i.test(normalized)) {
      return `rgba(82, 82, 91, ${alpha})`
    }
    const red = Number.parseInt(normalized.slice(0, 2), 16)
    const green = Number.parseInt(normalized.slice(2, 4), 16)
    const blue = Number.parseInt(normalized.slice(4, 6), 16)
    return `rgba(${red}, ${green}, ${blue}, ${alpha})`
  }
  const elementCardStyle = (element: PeriodicTableElement) => {
    const categoryColor = effectiveTheme.value === 'light' ? element.categoryLightColor : element.categoryColor
    const elementNumber = element.number
    if (isElementUnavailable(elementNumber) && !selectionToneMap.value[elementNumber]) {
      if (effectiveTheme.value === 'light') {
        return {
          borderColor: 'color-mix(in srgb, var(--lab-border) 92%, transparent)',
          background: 'color-mix(in srgb, var(--lab-bg-surface) 88%, rgb(228 228 231))'
        }
      }

      return {
        borderColor: 'color-mix(in srgb, var(--lab-border) 92%, transparent)',
        background: 'color-mix(in srgb, var(--lab-bg-surface) 90%, rgb(9 9 11))'
      }
    }
    if (props.useThemeCardColors) {
      if (effectiveTheme.value === 'light') {
        return {
          borderColor: `color-mix(in srgb, ${categoryColor} 72%, var(--lab-border))`,
          background: categoryColor
        }
      }
      return {
        borderColor: `color-mix(in srgb, ${categoryColor} 34%, var(--lab-border))`,
        background: `color-mix(in srgb, ${categoryColor} 18%, var(--lab-bg-surface))`
      }
    }
    return {
      borderColor: colorHexToRgba(categoryColor, 0.4),
      backgroundImage: `linear-gradient(160deg, ${colorHexToRgba(categoryColor, 0.26)} 0%, rgba(24, 24, 27, 0.96) 72%)`
    }
  }
  const onElementClick = (element: PeriodicTableElement, event: MouseEvent) => {
    if (isElementDisabled(element.number)) return
    emit('elementClick', element, event)
  }
  const elementSelectionToneClass = (elementNumber: number) => {
    const tone = selectionToneMap.value[elementNumber]
    if (tone === 'all') {
      return 'ring-2 ring-amber-300/80 bg-amber-200/18'
    }
    if (tone === 'any') {
      return 'ring-2 ring-cyan-300/80 bg-cyan-200/18'
    }
    if (tone === 'none') {
      return 'ring-2 ring-rose-300/80 bg-rose-200/18'
    }
    return ''
  }
  const isElementUnavailable = (elementNumber: number) => {
    if (!props.disableUnavailable || !hasAvailabilityState.value) return false
    return unavailableSet.value.has(elementNumber)
  }
  const isElementDisabled = (elementNumber: number) =>
    isElementUnavailable(elementNumber) && !selectionToneMap.value[elementNumber]
  const elementUnavailableClass = (elementNumber: number) => {
    if (!isElementUnavailable(elementNumber) || selectionToneMap.value[elementNumber]) return ''
    return 'text-(--lab-text-primary) ring-1 ring-inset'
  }
</script>
<template>
  <div class="flex min-h-0 w-full max-w-full items-start justify-start overflow-x-auto overflow-y-hidden">
    <div :class="frameClass" :style="frameStyle">
      <div class="grid h-full w-full gap-(--periodic-gap)" :style="layoutStyle">
        <div class="hidden sm:block" aria-hidden="true"></div>
        <div class="hidden gap-(--periodic-gap) text-center sm:grid" :style="tableColumnHeaderStyle">
          <div
            v-for="group in groupHeaders"
            :key="`group-header:${group}`"
            class="flex items-center justify-center font-medium tracking-[0.12em] text-(--lab-text-muted) text-[0.55rem] sm:text-[0.68rem]"
            :class="compact ? 'lg:text-[0.5rem]' : ''">
            {{ group }}
          </div>
        </div>
        <div class="hidden h-full min-h-0 gap-(--periodic-gap) text-right sm:grid" :style="tableRowHeaderStyle">
          <div
            v-for="(period, index) in periodHeaders"
            :key="`period-header:${visiblePeriodRows[index]}`"
            class="flex items-center justify-end pr-1 font-medium tracking-[0.12em] text-(--lab-text-muted) text-[0.55rem] sm:text-[0.68rem]"
            :class="compact ? 'lg:text-[0.5rem]' : ''"
            aria-hidden="true">
            {{ period }}
          </div>
        </div>
        <div class="relative h-full min-h-0 min-w-0 grid gap-(--periodic-gap)" :style="tableGridStyle">
          <slot name="grid-overlay"></slot>
          <div
            v-if="showSeriesLabels"
            class="relative z-10 flex items-center border border-dashed bg-(--lab-bg-surface) px-1 text-center text-[0.52rem] leading-tight text-(--lab-text-muted) sm:px-2 sm:text-[0.68rem]"
            :style="{ gridColumn: '1 / span 2', gridRow: `${getDisplayRow(9)} / span 1` }">
            Лантаноиды
          </div>
          <div
            v-if="showSeriesLabels"
            class="relative z-10 flex items-center border border-dashed border-(--lab-border) bg-(--lab-bg-surface) px-1 text-center text-[0.52rem] leading-tight text-(--lab-text-muted) sm:px-2 sm:text-[0.68rem]"
            :style="{ gridColumn: '1 / span 2', gridRow: `${getDisplayRow(10)} / span 1` }">
            Актиноиды
          </div>
          <button
            v-for="element in visibleElements"
            :key="element.number"
            type="button"
            :style="[elementSlotStyle(element.xpos, element.ypos), elementCardStyle(element)]"
            class="relative z-10 flex h-full w-full min-h-0 min-w-0 flex-col overflow-hidden border p-0.5 px-1 py-1 text-left lg:p-1 sm:px-2 sm:py-1.5"
            :disabled="isElementDisabled(element.number)"
            :class="[
              elementUnavailableClass(element.number),
              isElementDisabled(element.number) ?
                'cursor-not-allowed opacity-100 saturate-[0.8] ring-1 ring-inset'
              : 'opacity-100',
              selectionToneMap[element.number] ? elementSelectionToneClass(element.number)
              : highlightedSet.has(element.number) ? 'ring-1 ring-white/60'
              : '',
              !isElementDisabled(element.number) && interactive ?
                'cursor-pointer hover:ring-1 hover:ring-(--lab-accent)'
              : !isElementDisabled(element.number) ? 'cursor-default'
              : ''
            ]"
            @click="onElementClick(element, $event)">
            <span
              v-if="compact"
              class="flex flex-1 items-center justify-center font-semibold leading-none text-[1.05rem] lg:text-[0.82rem]"
              :class="isElementDisabled(element.number) ? 'text-(--lab-text-muted)' : 'text-(--lab-text-primary)'">
              {{ element.displaySymbol }}
            </span>

            <span
              v-else
              class="flex flex-1 items-center justify-center font-semibold leading-none text-[1.05rem] 2xl:hidden"
              :class="isElementDisabled(element.number) ? 'text-(--lab-text-muted)' : 'text-(--lab-text-primary)'">
              {{ element.displaySymbol }}
            </span>

            <div v-if="!compact" class="hidden h-full min-w-0 flex-col 2xl:flex">
              <div class="flex items-start justify-between gap-2">
                <span
                  class="font-semibold leading-none"
                  style="font-size: clamp(0.95rem, 0.55vw + 0.65rem, 1.45rem)"
                  :class="isElementDisabled(element.number) ? 'text-(--lab-text-muted)' : 'text-(--lab-text-primary)'">
                  {{ element.displaySymbol }}
                </span>
                <span class="pt-0.5 text-[0.68rem] font-medium leading-none text-(--lab-text-muted)">
                  {{ element.number }}
                </span>
              </div>

              <div class="mt-1 min-w-0">
                <span
                  class="block overflow-hidden text-ellipsis whitespace-nowrap font-medium leading-[1.2]"
                  style="font-size: clamp(0.6rem, 0.18vw + 0.5rem, 0.82rem)"
                  :class="isElementDisabled(element.number) ? 'text-(--lab-text-muted)' : 'text-(--lab-text-primary)'">
                  {{ element.name }}
                </span>
                <span
                  class="mt-0.5 block overflow-hidden text-ellipsis whitespace-nowrap leading-[1.2] text-(--lab-text-muted)"
                  style="font-size: clamp(0.56rem, 0.14vw + 0.48rem, 0.74rem)">
                  {{ element.russianName }}
                </span>
              </div>
            </div>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
