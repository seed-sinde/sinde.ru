<template>
  <div class="relative min-w-0">
    <h1 v-if="title">{{ title }}</h1>
    <div ref="scrRef" class="lab-scroll-hidden max-h-136 overflow-auto">
      <table class="min-w-full text-xs">
        <thead class="sticky top-0 z-10">
          <tr class="border-b border-(--lab-border) bg-(--lab-bg-canvas)">
            <th class="w-7 py-2 pr-2 text-left font-medium">
              <slot name="loading" :loading="loading">
                <div class="flex h-4 items-center">
                  <Icon
                    name="ic:round-autorenew"
                    class="h-3.5 w-3.5 transition-opacity"
                    :class="loading ? 'animate-spin opacity-100' : 'opacity-0'"
                    aria-hidden="true"
                  />
                  <span class="sr-only">{{ loading ? 'Загрузка' : 'Ожидание' }}</span>
                </div>
              </slot>
            </th>
            <th v-for="col in cols" :key="col.key" :class="hdCls(col)">
              <slot :name="`header-${col.key}`" :column="col">
                {{ col.label }}
              </slot>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="!rowsSafe.length" class="border-b border-(--lab-border)">
            <td :colspan="colsLen + 1" class="py-3 pr-3">
              <slot name="empty">{{ emptyText }}</slot>
            </td>
          </tr>
          <tr v-for="(row,i) in rowsSafe" :key="rk(row, i)" class="border-b border-(--lab-border) align-top">
            <td class="w-7 py-2 pr-2" />
            <td v-for="col in cols" :key="`${String(rk(row, i))}:${col.key}`" :class="tdCls(col)">
              <slot :name="`cell-${col.key}`" :row="row" :column="col" :value="row?.[col.key]" :index="i">
                <slot name="cell" :row="row" :column="col" :value="row?.[col.key]" :index="i">
                  {{ fmt(row?.[col.key]) }}
                </slot>
              </slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div
      class="lab-scroll-fade lab-scroll-fade-x-left"
      :class="{ 'lab-scroll-fade-visible': edges.left }"
      aria-hidden="true"
    />
    <div
      class="lab-scroll-fade lab-scroll-fade-x-right"
      :class="{ 'lab-scroll-fade-visible': edges.right }"
      aria-hidden="true"
    />
  </div>
</template>

<script setup lang="ts">
type Row = LabDataTableRow
type Col = LabDataTableColumn

const props = withDefaults(defineProps<{
  columns: Col[]
  rows: Row[]
  loading?: boolean
  emptyText?: string
  rowKey?: string | ((row: Row, index: number) => PropertyKey)
  title?: string
  nowrap?: boolean
}>(), {
  loading: false,
  emptyText: 'Ничего не найдено.',
  rowKey: 'id',
  nowrap: false
})

const { columns: cols, loading, emptyText } = toRefs(props)
const colsLen = computed(() => cols.value.length)
const rowsSafe = computed<Row[]>(() => Array.isArray(props.rows) ? props.rows : [])
const scrRef = ref<HTMLElement | null>(null)
const { edges, sync } = useScrollableEdges(scrRef, { axis: 'x' })
const rk = (row: Row, i: number): PropertyKey => {
  if (typeof props.rowKey === 'function') return props.rowKey(row, i)
  const k = String(props.rowKey || 'id')
  const v = row?.[k]
  return typeof v === 'string' || typeof v === 'number' || typeof v === 'symbol' ? v : v !== undefined && v !== null && String(v) !== '' ? String(v) : i
}

const hdCls = (col: Col) => [
  'py-2 pr-3 text-left font-medium',
  col.widthClass || '',
  col.headerClass || '',
  props.nowrap || col.nowrap ? 'whitespace-nowrap' : ''
]

const tdCls = (col: Col) => [
  'py-2 pr-3',
  col.widthClass || '',
  col.cellClass || '',
  props.nowrap || col.nowrap ? 'whitespace-nowrap' : ''
]

const fmt = (v: unknown) => v === null || v === undefined ? '—' : typeof v === 'string' ? v || '—' : String(v)

watch(
  () => [rowsSafe.value.length, colsLen.value, loading.value] as const,
  () => void nextTick(sync),
  { immediate: true }
)
</script>
