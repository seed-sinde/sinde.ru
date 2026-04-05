<template>
  <div class="flex flex-wrap gap-2" :class="listClass">
    <LabBaseButton
      v-for="item in items"
      :key="item.key"
      button-class="group relative inline-flex items-center justify-center rounded-full border px-3 py-1 text-xs transition hover:border-rose-500/50! hover:bg-rose-500/10! hover:text-rose-300!"
      icon="ic:round-delete-outline"
      icon-class="pointer-events-none absolute left-1/2 top-1/2 h-3.5 w-3.5 -translate-x-1/2 -translate-y-1/2 opacity-0 transition-opacity group-hover:opacity-100"
      :style="resolveItemStyle(item)"
      :disabled="Boolean(item.disabled)"
      :title="resolveItemTitle(item)"
      :aria-label="resolveItemAriaLabel(item)"
      size="xs"
      @click="$emit('remove', item)">
      <span class="group-hover:invisible">{{ item.label }}</span>
    </LabBaseButton>
    <span v-if="!items.length && emptyText" class="text-xs text-zinc-500">{{ emptyText }}</span>
  </div>
</template>
<script setup lang="ts">
  import type { StyleValue } from 'vue'
  export type KitchenRemovableChipItem = {
    key: string | number
    label: string
    disabled?: boolean
    payload?: unknown
  }
  const props = withDefaults(
    defineProps<{
      items?: KitchenRemovableChipItem[]
      emptyText?: string
      listClass?: string
      itemStyle?: StyleValue | ((item: KitchenRemovableChipItem) => StyleValue)
      itemTitle?: string | ((item: KitchenRemovableChipItem) => string)
      itemAriaLabel?: string | ((item: KitchenRemovableChipItem) => string)
    }>(),
    {
      items: () => [],
      emptyText: '',
      listClass: '',
      itemStyle: undefined,
      itemTitle: '',
      itemAriaLabel: ''
    }
  )
  defineEmits<{
    (e: 'remove', item: KitchenRemovableChipItem): void
  }>()
  /** Resolves an inline style or style factory for the chip. */
  const resolveItemStyle = (item: KitchenRemovableChipItem) =>
    typeof props.itemStyle === 'function' ? props.itemStyle(item) : props.itemStyle
  /** Resolves the title text shown on chip hover. */
  const resolveItemTitle = (item: KitchenRemovableChipItem) =>
    typeof props.itemTitle === 'function' ? props.itemTitle(item) : props.itemTitle
  /** Resolves the aria-label for the chip removal action. */
  const resolveItemAriaLabel = (item: KitchenRemovableChipItem) =>
    typeof props.itemAriaLabel === 'function' ? props.itemAriaLabel(item) : props.itemAriaLabel
</script>
