<script setup lang="ts">
  const { effectiveTheme } = useInterfacePreferences()
  const props = withDefaults(
    defineProps<{
      elements: PeriodicTableElement[]
      dimmedElementNumbers?: number[]
      highlightedElementNumbers?: number[]
      selectionToneByNumber?: Partial<Record<number, 'all' | 'any' | 'none'>>
      availabilityByNumber?: Partial<Record<number, number>>
      disableUnavailable?: boolean
      useThemeCardColors?: boolean
      preserveCompactCellWidth?: boolean
      compactScrollBreakpoint?: 'md' | 'lg'
      compactCellHeight?: string
      compactCellHeightWide?: string
      compact?: boolean
      interactive?: boolean
      showGroupHeaders?: boolean
      showPeriodHeaders?: boolean
      showSeriesLabels?: boolean
      hiddenElementNumbers?: number[]
      layout?: 'default' | 'minerals'
    }>(),
    {
      dimmedElementNumbers: () => [],
      highlightedElementNumbers: () => [],
      selectionToneByNumber: () => ({}),
      availabilityByNumber: () => ({}),
      disableUnavailable: false,
      useThemeCardColors: false,
      preserveCompactCellWidth: false,
      compactScrollBreakpoint: 'lg',
      compactCellHeight: '',
      compactCellHeightWide: '',
      compact: false,
      interactive: false,
      showGroupHeaders: true,
      showPeriodHeaders: true,
      showSeriesLabels: true,
      hiddenElementNumbers: () => [],
      layout: 'default'
    }
  )
  const emit = defineEmits<{
    elementClick: [element: PeriodicTableElement, event: MouseEvent]
  }>()
  const periodicTableGroups = computed(() => Math.max(...props.elements.map(element => element.xpos), 18))
  const groupHeaders = computed(() => Array.from({ length: periodicTableGroups.value }, (_, index) => index + 1))
  const hasVisibleHeaders = computed(() => props.showGroupHeaders || props.showPeriodHeaders)
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
  const periodHeaders = computed(() =>
    visiblePeriodRows.value.map(row => {
      if (row === 9) return '6'
      if (row === 10) return '7'
      return String(row)
    })
  )
  const dimmedSet = computed(() => new Set(props.dimmedElementNumbers))
  const highlightedSet = computed(() => new Set(props.highlightedElementNumbers))
  const selectionToneMap = computed(() => props.selectionToneByNumber || {})
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
  const hiddenSet = computed(() => new Set(props.hiddenElementNumbers))
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
  const visibleGroupHeaders = computed(() =>
    isMineralsLayout.value ? groupHeaders.value.filter(group => group !== 18) : groupHeaders.value
  )
  const getDisplayRow = (row: number) => {
    const index = visiblePeriodRows.value.indexOf(row)
    return index >= 0 ? index + 1 : row
  }
  const getDisplayColumn = (column: number) => displayColumnBySourceColumn.value.get(column) ?? column
  const displayColumnBySourceColumn = computed(() => {
    const visibleGroups = isMineralsLayout.value ? groupHeaders.value.filter(group => group !== 18) : groupHeaders.value
    return new Map<number, number>(visibleGroups.map((group, index) => [group, index + 1]))
  })
  const tableGridStyle = computed(() => ({
    display: 'grid',
    gridTemplateColumns: `repeat(${visibleGroupHeaders.value.length}, minmax(0, 1fr))`,
    gridTemplateRows: `repeat(${visiblePeriodRows.value.length}, minmax(0, var(--periodic-cell-height)))`
  }))
  const frameStyle = computed(() => ({
    '--periodic-row-label-size': props.showPeriodHeaders ? undefined : '0px',
    '--periodic-column-label-size': props.showGroupHeaders ? undefined : '0px',
    '--periodic-compact-cell-height': props.compactCellHeight || undefined,
    '--periodic-compact-cell-height-wide': props.compactCellHeightWide || undefined,
    '--periodic-row-count': String(visiblePeriodRows.value.length)
  }))
  const layoutStyle = computed(() => ({
    gridTemplateColumns: hasVisibleHeaders.value
      ? `${props.showPeriodHeaders ? 'var(--periodic-row-label-size)' : '0px'} minmax(0, 1fr)`
      : 'minmax(0, 1fr)',
    gridTemplateRows: hasVisibleHeaders.value
      ? `${props.showGroupHeaders ? 'var(--periodic-column-label-size)' : '0px'} minmax(0, 1fr)`
      : 'minmax(0, 1fr)'
  }))
  const tableColumnHeaderStyle = computed(() => ({
    display: 'grid',
    gridTemplateColumns: `repeat(${visibleGroupHeaders.value.length}, minmax(0, 1fr))`
  }))
  const tableRowHeaderStyle = computed(() => ({
    display: 'grid',
    gridTemplateRows: `repeat(${visiblePeriodRows.value.length}, minmax(0, var(--periodic-cell-height)))`
  }))
  const elementSlotStyle = (xpos: number, ypos: number) => ({
    gridColumn: `${getDisplayColumn(xpos)} / span 1`,
    gridRow: `${getDisplayRow(ypos)} / span 1`
  })
  const columnGuideStyle = (column: number) => ({
    gridColumn: `${getDisplayColumn(column)} / span 1`,
    gridRow: `1 / span ${visiblePeriodRows.value.length}`
  })
  const rowGuideStyle = (row: number) => ({
    gridColumn: `1 / span ${visibleGroupHeaders.value.length}`,
    gridRow: `${getDisplayRow(row)} / span 1`
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
      return {
        borderColor: 'var(--lab-border)',
        background: 'var(--lab-bg-surface)'
      }
    }
    if (isElementDimmed(elementNumber)) {
      if (effectiveTheme.value === 'light') {
        return {
          borderColor: 'color-mix(in srgb, var(--lab-border) 88%, transparent)',
          background: `color-mix(in srgb, ${categoryColor} 52%, white)`
        }
      }
      return {
        borderColor: `color-mix(in srgb, ${categoryColor} 18%, var(--lab-border))`,
        background: `color-mix(in srgb, ${categoryColor} 8%, var(--lab-bg-surface))`
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
  const isElementDimmed = (elementNumber: number) =>
    dimmedSet.value.has(elementNumber) && !selectionToneMap.value[elementNumber]
  const isElementDisabled = (elementNumber: number) =>
    isElementUnavailable(elementNumber) && !selectionToneMap.value[elementNumber]
  const elementUnavailableClass = (elementNumber: number) => {
    if (!isElementUnavailable(elementNumber) || selectionToneMap.value[elementNumber]) return ''
    return 'text-zinc-500 ring-1 ring-inset ring-zinc-900'
  }
</script>
<template>
  <div
    class="periodic-table-scroll flex min-h-0 w-full max-w-full items-start justify-start"
    :class="[
      { 'periodic-table-scroll-compact': compact },
      preserveCompactCellWidth
        ? compactScrollBreakpoint === 'md'
          ? 'periodic-table-scroll-fixed-md'
          : 'periodic-table-scroll-fixed-lg'
        : ''
    ]">
    <div class="periodic-table-frame" :class="{ 'periodic-table-frame-compact': compact }" :style="frameStyle">
      <div class="periodic-table-layout" :style="layoutStyle">
        <div v-if="hasVisibleHeaders" aria-hidden="true"></div>
        <div v-if="showGroupHeaders" class="periodic-column-headers text-center" :style="tableColumnHeaderStyle">
          <div
            v-for="group in visibleGroupHeaders"
            :key="`group-header:${group}`"
            class="flex items-center justify-center text-[0.55rem] font-medium tracking-[0.12em] text-zinc-500 sm:text-[0.68rem]">
            {{ group }}
          </div>
        </div>
        <div v-if="showPeriodHeaders" class="periodic-row-headers text-right" :style="tableRowHeaderStyle">
          <div
            v-for="(period, index) in periodHeaders"
            :key="`period-header:${visiblePeriodRows[index]}`"
            class="flex items-center justify-end pr-1 text-[0.55rem] font-medium tracking-[0.12em] text-zinc-500 sm:text-[0.68rem]"
            aria-hidden="true">
            {{ period }}
          </div>
        </div>
        <div class="periodic-grid" :style="tableGridStyle">
          <div
            v-for="column in periodicTableGroups"
            :key="`column-guide:${column}`"
            class="periodic-column-guide border-(--lab-border) pointer-events-none z-0 border-l"
            :class="column === 1 ? 'border-l-transparent' : ''"
            :style="columnGuideStyle(column)"
            aria-hidden="true"></div>
          <div
            v-for="row in visiblePeriodRows"
            :key="`row-guide:${row}`"
            class="periodic-row-guide border-(--lab-border) pointer-events-none z-0 border-t"
            :class="row === 1 ? 'border-t-transparent' : ''"
            :style="rowGuideStyle(row)"
            aria-hidden="true"></div>
          <slot name="grid-overlay"></slot>
          <div
            v-if="showSeriesLabels"
            class="border-(--lab-border) bg-(--lab-bg-surface) text-(--lab-text-muted) relative z-10 flex items-center border border-dashed px-1 text-center text-[0.52rem] leading-tight sm:px-2 sm:text-[0.68rem]"
            :style="{ gridColumn: '1 / span 2', gridRow: `${getDisplayRow(9)} / span 1` }">
            Лантаноиды
          </div>
          <div
            v-if="showSeriesLabels"
            class="border-(--lab-border) bg-(--lab-bg-surface) text-(--lab-text-muted) relative z-10 flex items-center border border-dashed px-1 text-center text-[0.52rem] leading-tight sm:px-2 sm:text-[0.68rem]"
            :style="{ gridColumn: '1 / span 2', gridRow: `${getDisplayRow(10)} / span 1` }">
            Актиноиды
          </div>
          <div
            v-for="element in visibleElements"
            :key="element.number"
            :style="elementSlotStyle(element.xpos, element.ypos)"
            class="relative z-10 min-h-0">
            <button
              type="button"
              :style="elementCardStyle(element)"
              class="periodic-element-card relative flex h-full w-full min-w-0 flex-col overflow-hidden border px-1 py-1 text-left sm:px-2 sm:py-1.5"
              :disabled="isElementDisabled(element.number)"
              :class="[
                elementUnavailableClass(element.number),
                isElementDisabled(element.number)
                  ? 'opacity-30 saturate-[0.45]'
                  : isElementDimmed(element.number)
                    ? 'opacity-55 saturate-[0.7]'
                    : 'opacity-100',
                selectionToneMap[element.number]
                  ? elementSelectionToneClass(element.number)
                  : highlightedSet.has(element.number)
                    ? 'ring-1 ring-white/60'
                    : '',
                isElementDisabled(element.number)
                  ? 'cursor-not-allowed ring-1 ring-inset ring-zinc-800'
                  : interactive
                    ? 'cursor-pointer hover:ring-1 hover:ring-white/35'
                    : 'cursor-default'
              ]"
              @click="onElementClick(element, $event)">
              <div class="periodic-element-compact flex flex-1 items-center justify-center">
                <span
                  class="periodic-element-compact-symbol font-semibold leading-none"
                  :class="isElementDisabled(element.number) ? 'text-zinc-500' : 'text-(--lab-text-primary)'">
                  {{ element.displaySymbol }}
                </span>
              </div>
              <div class="periodic-element-full h-full min-w-0 flex-col">
                <div class="flex items-start justify-between gap-2">
                  <span
                    class="periodic-element-symbol font-semibold leading-none"
                    :class="isElementDisabled(element.number) ? 'text-zinc-500' : 'text-(--lab-text-primary)'">
                    {{ element.displaySymbol }}
                  </span>
                  <span
                    class="periodic-element-number pt-0.5 font-medium leading-none"
                    :class="isElementDisabled(element.number) ? 'text-zinc-600' : 'text-(--lab-text-muted)'">
                    {{ element.number }}
                  </span>
                </div>
                <div class="mt-1 min-w-0">
                  <span
                    class="periodic-element-name"
                    :class="isElementDisabled(element.number) ? 'text-zinc-500' : 'text-(--lab-text-primary)'">
                    {{ element.name }}
                  </span>
                  <span
                    class="periodic-element-russian-name mt-0.5"
                    :class="isElementDisabled(element.number) ? 'text-zinc-600' : 'text-(--lab-text-muted)'">
                    {{ element.russianName }}
                  </span>
                </div>
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<style scoped>
  .periodic-table-scroll {
    align-items: flex-start;
    overflow-x: auto;
    overflow-y: hidden;
  }
  .periodic-table-frame {
    --periodic-gap: 0.25rem;
    --periodic-highlight-reserve: 6px;
    --periodic-row-label-size: 1rem;
    --periodic-column-label-size: 1rem;
    --periodic-mobile-cell-width: 2.85rem;
    --periodic-mobile-cell-height: 2.2rem;
    --periodic-cell-width: var(--periodic-mobile-cell-width);
    --periodic-cell-height: var(--periodic-mobile-cell-height);
    --periodic-series-gap-height: 0rem;
    width: calc(
      var(--periodic-row-label-size) + var(--periodic-gap) + 18 * var(--periodic-cell-width) + 17 * var(--periodic-gap)
    );
    height: calc(
      var(--periodic-column-label-size) + var(--periodic-gap) + var(--periodic-row-count) *
        var(--periodic-cell-height) + var(--periodic-series-gap-height) + var(--periodic-row-count) *
        var(--periodic-gap)
    );
    padding: var(--periodic-highlight-reserve);
    box-sizing: border-box;
    flex: none;
  }
  .periodic-table-layout {
    display: grid;
    grid-template-columns: var(--periodic-row-label-size) minmax(0, 1fr);
    grid-template-rows: var(--periodic-column-label-size) minmax(0, 1fr);
    gap: var(--periodic-gap);
    width: 100%;
    height: 100%;
  }
  .periodic-column-headers,
  .periodic-row-headers,
  .periodic-grid {
    gap: var(--periodic-gap);
  }
  .periodic-row-headers,
  .periodic-grid {
    min-height: 0;
    height: 100%;
  }
  .periodic-grid {
    position: relative;
    min-width: 0;
  }
  .periodic-element-card {
    min-height: 0;
  }
  @media (max-width: 639px) {
    .periodic-column-guide,
    .periodic-row-guide {
      display: none;
    }
  }
  .periodic-element-compact {
    display: flex;
  }
  .periodic-element-full {
    display: none;
  }
  .periodic-element-compact-symbol {
    font-size: 1.05rem;
  }
  .periodic-element-symbol {
    font-size: clamp(0.95rem, 0.55vw + 0.65rem, 1.45rem);
  }
  .periodic-element-number {
    font-size: 0.68rem;
  }
  .periodic-element-name,
  .periodic-element-russian-name {
    display: block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .periodic-element-name {
    font-size: clamp(0.6rem, 0.18vw + 0.5rem, 0.82rem);
    line-height: 1.2;
    font-weight: 500;
  }
  .periodic-element-russian-name {
    font-size: clamp(0.56rem, 0.14vw + 0.48rem, 0.74rem);
    line-height: 1.2;
  }
  @media (min-width: 960px) {
    .periodic-table-frame {
      --periodic-gap: 0.35rem;
      --periodic-row-label-size: 1.5rem;
      --periodic-column-label-size: 1.5rem;
      --periodic-compact-cell-width: 2.1rem;
      --periodic-compact-cell-height: var(--periodic-compact-cell-height, 1.8rem);
      --periodic-cell-width: var(--periodic-compact-cell-width);
      --periodic-cell-height: var(--periodic-compact-cell-height);
      --periodic-series-gap-height: 0rem;
      width: 100%;
      min-width: 0;
      min-height: 0;
      height: calc(
        var(--periodic-column-label-size) + var(--periodic-gap) + var(--periodic-row-count) *
          var(--periodic-cell-height) + var(--periodic-series-gap-height) + var(--periodic-row-count) *
          var(--periodic-gap)
      );
      flex: none;
    }
    .periodic-table-frame-compact {
      --periodic-gap: 0.2rem;
      --periodic-row-label-size: 1rem;
      --periodic-column-label-size: 1rem;
      --periodic-compact-cell-width: 1.65rem;
      --periodic-compact-cell-height: var(--periodic-compact-cell-height, 1.45rem);
      --periodic-cell-width: var(--periodic-compact-cell-width);
      --periodic-cell-height: var(--periodic-compact-cell-height);
      --periodic-series-gap-height: 0rem;
    }
    .periodic-table-frame-compact .periodic-element-compact-symbol {
      font-size: 0.82rem;
    }
    .periodic-table-frame-compact .periodic-column-headers > div,
    .periodic-table-frame-compact .periodic-row-headers > div {
      font-size: 0.5rem;
    }
    .periodic-table-scroll-fixed-lg .periodic-table-frame {
      width: calc(
        var(--periodic-row-label-size) + var(--periodic-gap) + 18 * var(--periodic-cell-width) + 17 *
          var(--periodic-gap)
      );
      min-width: auto;
    }
  }
  @media (min-width: 768px) and (max-width: 959px) {
    .periodic-table-scroll-fixed-md .periodic-table-frame {
      width: calc(
        var(--periodic-row-label-size) + var(--periodic-gap) + 18 * var(--periodic-cell-width) + 17 *
          var(--periodic-gap)
      );
      min-width: auto;
    }
  }
  @media (min-width: 1280px) {
    .periodic-table-frame {
      --periodic-gap: 0.4rem;
      --periodic-compact-cell-width: 2.2rem;
    }
    .periodic-table-frame-compact {
      --periodic-gap: 0.24rem;
      --periodic-compact-cell-width: 1.75rem;
      --periodic-compact-cell-height: var(
        --periodic-compact-cell-height-wide,
        var(--periodic-compact-cell-height, 1.55rem)
      );
    }
  }
  @media (min-width: 1400px) {
    .periodic-table-scroll {
      align-items: flex-start;
    }
    .periodic-table-frame {
      --periodic-full-cell-height: 4.75rem;
      --periodic-cell-height: var(--periodic-full-cell-height);
      --periodic-series-gap-height: 0rem;
      height: calc(
        var(--periodic-column-label-size) + var(--periodic-gap) + var(--periodic-row-count) *
          var(--periodic-cell-height) + var(--periodic-series-gap-height) + var(--periodic-row-count) *
          var(--periodic-gap)
      );
      flex: none;
    }
    .periodic-element-compact {
      display: none;
    }
    .periodic-element-full {
      display: flex;
    }
    .periodic-table-frame-compact {
      --periodic-full-cell-height: var(
        --periodic-compact-cell-height-wide,
        var(--periodic-compact-cell-height, 1.55rem)
      );
      --periodic-cell-height: var(--periodic-full-cell-height);
      --periodic-series-gap-height: 0rem;
      height: calc(
        var(--periodic-column-label-size) + var(--periodic-gap) + var(--periodic-row-count) *
          var(--periodic-cell-height) + var(--periodic-series-gap-height) + var(--periodic-row-count) *
          var(--periodic-gap)
      );
    }
    .periodic-table-frame-compact .periodic-element-compact {
      display: flex;
    }
    .periodic-table-frame-compact .periodic-element-full {
      display: none;
    }
  }
</style>
