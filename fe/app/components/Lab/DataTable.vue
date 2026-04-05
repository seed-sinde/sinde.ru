<template>
  <div class="relative overflow-auto" :class="[maxHeightClass, containerClass]">
    <table :class="tableClassList">
      <thead :class="theadClassList">
        <tr class="border-b border-zinc-800 text-zinc-400">
          <th class="w-7 py-2 pr-2 text-left font-medium">
            <slot name="loading" :loading="loading">
              <div class="flex h-4 items-center">
                <Icon
                  name="ic:round-autorenew"
                  class="h-3.5 w-3.5 text-zinc-500 transition-opacity"
                  :class="loading ? 'animate-spin opacity-100' : 'opacity-0'"
                  aria-hidden="true" />
                <span class="sr-only">{{ loading ? 'Загрузка' : 'Ожидание' }}</span>
              </div>
            </slot>
          </th>
          <th v-for="column in columns" :key="column.key" :class="headerClass(column)">
            <slot :name="`header-${column.key}`" :column="column">
              {{ column.label }}
            </slot>
          </th>
        </tr>
      </thead>
      <tbody :class="tbodyClass">
        <tr v-if="normalizedRows.length === 0" class="border-b border-zinc-900">
          <td :colspan="Math.max(columns.length + 1, 1)" class="py-3 pr-3 text-zinc-500">
            <slot name="empty">{{ emptyText }}</slot>
          </td>
        </tr>
        <tr v-for="(row, index) in normalizedRows" :key="resolveRowKey(row, index)" :class="rowClass">
          <td class="w-7 py-2 pr-2"></td>
          <td v-for="column in columns" :key="`${String(resolveRowKey(row, index))}:${column.key}`" :class="cellClass(column)">
            <slot :name="`cell-${column.key}`" :row="row" :column="column" :value="row?.[column.key]" :index="index">
              <slot name="cell" :row="row" :column="column" :value="row?.[column.key]" :index="index">
                {{ formatCellValue(row?.[column.key]) }}
              </slot>
            </slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      columns: LabDataTableColumn[]
      rows: any[]
      loading?: boolean
      emptyText?: string
      rowKey?: string | ((row: any, index: number) => PropertyKey)
      tableClass?: string
      theadClass?: string
      tbodyClass?: string
      rowClass?: string
      maxHeightClass?: string
      containerClass?: string
    }>(),
    {
      loading: false,
      emptyText: 'Ничего не найдено.',
      rowKey: 'id',
      tableClass: 'min-w-full text-xs',
      theadClass: 'sticky top-0 z-10 bg-zinc-950/95',
      tbodyClass: '',
      rowClass: 'border-b border-zinc-900 align-top',
      maxHeightClass: 'max-h-136',
      containerClass: ''
    }
  )
  const normalizedRows = computed(() => (Array.isArray(props.rows) ? props.rows : []))
  const tableClassList = computed(() => [props.tableClass])
  const theadClassList = computed(() => [props.theadClass])
  const resolveRowKey = (row: LabDataTableRow, index: number): PropertyKey => {
    if (typeof props.rowKey === 'function') return props.rowKey(row, index)
    const keyName = String(props.rowKey || 'id')
    const value = row?.[keyName]
    if (typeof value === 'string' || typeof value === 'number' || typeof value === 'symbol') return value
    if (value !== undefined && value !== null && String(value) !== '') return String(value)
    return `${index}`
  }
  const headerClass = (column: LabDataTableColumn) => [
    'py-2 pr-3 text-left font-medium',
    column.nowrap ? 'whitespace-nowrap' : '',
    column.widthClass || '',
    column.headerClass || ''
  ]
  const cellClass = (column: LabDataTableColumn) => [
    'py-2 pr-3',
    column.nowrap ? 'whitespace-nowrap' : '',
    column.widthClass || '',
    column.cellClass || ''
  ]
  const formatCellValue = (value: unknown) => {
    if (value === null || value === undefined) return '—'
    if (typeof value === 'string') return value || '—'
    return String(value)
  }
</script>
