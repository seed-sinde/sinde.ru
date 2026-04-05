<template>
  <fieldset :class="containerClass">
    <div class="mb-2 flex flex-wrap items-center gap-x-3 gap-y-1">
      <span class="px-1 text-[11px] uppercase tracking-[0.06em]" :class="titleClass">{{ title }}</span>
      <slot name="header-extra" />
      <LabBaseSwitch
        v-if="toggleLabel"
        :id="toggleId"
        :model-value="toggleValue"
        :name="toggleName"
        :label="toggleLabel"
        :tone="toggleTone"
        :visual-state="toggleVisualState"
        @update:model-value="$emit('update:toggleValue', Boolean($event))" />
    </div>
    <div v-if="items.length" :class="listContainerClass">
      <KitchenRemovableChipList
        :items="items"
        :item-style="itemStyle"
        :item-title="itemTitle"
        :item-aria-label="itemAriaLabel"
        @remove="$emit('remove', $event)" />
    </div>
    <p v-if="showFavoriteEmptyText" class="mt-2 text-xs text-zinc-500">{{ favoriteEmptyText }}</p>
    <p v-if="!items.length" class="text-xs text-zinc-500">{{ emptyText }}</p>
  </fieldset>
</template>
<script setup lang="ts">
  import type { StyleValue } from 'vue'
  type KitchenIngredientScopeChip = {
    key: string | number
    label: string
    disabled?: boolean
    payload?: unknown
  }
  const props = withDefaults(
    defineProps<{
      title: string
      titleClass?: string
      containerClass?: string
      listContainerClass?: string
      items?: KitchenIngredientScopeChip[]
      emptyText?: string
      favoriteEmptyText?: string
      showFavoriteEmptyText?: boolean
      toggleId?: string
      toggleName?: string
      toggleValue?: boolean
      toggleLabel?: string
      toggleTone?: ToggleTone
      toggleVisualState?: ToggleVisualState
      itemStyle?: StyleValue | ((item: KitchenIngredientScopeChip) => StyleValue)
      itemTitle?: string | ((item: KitchenIngredientScopeChip) => string)
      itemAriaLabel?: string | ((item: KitchenIngredientScopeChip) => string)
    }>(),
    {
      titleClass: 'text-zinc-400',
      containerClass: '',
      listContainerClass: '',
      items: () => [],
      emptyText: '',
      favoriteEmptyText: '',
      showFavoriteEmptyText: false,
      toggleId: undefined,
      toggleName: undefined,
      toggleValue: false,
      toggleLabel: '',
      toggleTone: 'cyan',
      toggleVisualState: undefined,
      itemStyle: undefined,
      itemTitle: '',
      itemAriaLabel: ''
    }
  )
  defineEmits<{
    (e: 'update:toggleValue', value: boolean): void
    (e: 'remove', item: KitchenIngredientScopeChip): void
  }>()
</script>
