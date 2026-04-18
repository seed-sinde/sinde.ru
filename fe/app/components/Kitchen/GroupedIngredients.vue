<template>
  <div class="space-y-4">
    <div
      v-if="groupByCategory && groupedItems.length"
      class="grid grid-cols-[repeat(auto-fit,minmax(18rem,1fr))] items-start gap-4"
    >
      <section v-for="group in groupedItems" :key="`group:${group.category}`" class="min-w-0">
        <div class="mb-2 flex items-center gap-2">
          <span class="h-px flex-1 bg-zinc-700/80" />
          <span class="text-xs tracking-[0.07em] whitespace-nowrap text-zinc-400 uppercase">
            {{ group.category }}
          </span>
          <span class="h-px flex-1 bg-zinc-700/80" />
        </div>
        <ul class="space-y-1.5">
          <li
            v-for="entry in group.items"
            :key="`ing:${entry.index}:${entry.item.name}`"
            class="grid w-full min-w-0 items-center gap-x-2 text-sm leading-tight whitespace-nowrap"
            :class="
              showActions
                ? 'grid-cols-[max-content_minmax(0,1.35fr)_minmax(0,0.75fr)_minmax(0,1fr)]'
                : 'grid-cols-[minmax(0,1.35fr)_minmax(0,0.75fr)_minmax(0,1fr)]'
            "
          >
            <div v-if="showActions" class="flex items-center gap-1">
              <LabBaseButton
                icon="ic:round-edit"
                icon-class="h-3.5 w-3.5"
                icon-only
                :title="`Редактировать ингредиент: ${ingredientLine(entry.item)}`"
                @click="emit('edit', entry.index)"
              />
              <LabBaseButton
                icon="ic:round-delete-outline"
                icon-class="h-3.5 w-3.5"
                icon-only
                :title="`Удалить ингредиент: ${ingredientLine(entry.item)}`"
                @click="emit('delete', entry.index)"
              />
            </div>
            <span class="min-w-0 truncate text-sm font-semibold text-zinc-100">{{ entry.item.name }}</span>
            <span class="min-w-0 truncate text-zinc-200 tabular-nums">{{ amountUnitText(entry.item) }}</span>
            <span class="min-w-0 truncate leading-tight text-zinc-400">{{ noteText(entry.item) }}</span>
          </li>
        </ul>
      </section>
    </div>
    <ul v-else-if="!groupByCategory && flatItems.length" class="w-full max-w-full space-y-1.5 pb-1">
      <li
        v-for="entry in flatItems"
        :key="`ing:flat:${entry.index}:${entry.item.name}`"
        class="grid w-full min-w-0 items-center gap-x-2 text-sm leading-tight whitespace-nowrap"
        :class="
          showActions
            ? 'grid-cols-[max-content_minmax(0,1.35fr)_minmax(0,0.75fr)_minmax(0,1fr)]'
            : 'grid-cols-[minmax(0,1.35fr)_minmax(0,0.75fr)_minmax(0,1fr)]'
        "
      >
        <div v-if="showActions" class="flex items-center gap-1">
          <LabBaseButton
            icon="ic:round-edit"
            icon-only
            variant="secondary"
            :title="`Редактировать ингредиент: ${ingredientLine(entry.item)}`"
            @click="emit('edit', entry.index)"
          />
          <LabBaseButton
            icon="ic:round-delete-outline"
            icon-only
            variant="danger"
            :title="`Удалить ингредиент: ${ingredientLine(entry.item)}`"
            @click="emit('delete', entry.index)"
          />
        </div>
        <span class="min-w-0 truncate text-sm font-semibold text-zinc-100">{{ entry.item.name }}</span>
        <span class="min-w-0 truncate text-zinc-200 tabular-nums">{{ amountUnitText(entry.item) }}</span>
        <span class="min-w-0 truncate leading-tight text-zinc-400">{{ noteText(entry.item) }}</span>
      </li>
    </ul>
    <p v-else class="text-xs text-zinc-500">{{ emptyText }}</p>
  </div>
</template>
<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    ingredients: KitchenIngredient[]
    categoryByName?: Record<string, string>
    groupByCategory?: boolean
    showActions?: boolean
    emptyText?: string
  }>(),
  {
    categoryByName: () => ({}),
    groupByCategory: true,
    showActions: false,
    emptyText: 'Еще не добавлены.'
  }
)
const emit = defineEmits<{
  edit: [index: number]
  delete: [index: number]
}>()
const normalizeTag = (value: string) =>
  String(value || '')
    .trim()
    .toLowerCase()
const categoryFor = (ingredient: KitchenIngredient) => {
  const byName = props.categoryByName?.[normalizeTag(ingredient.name)]
  return String(byName || 'другое').trim() || 'другое'
}
const groupedItems = computed(() => {
  const grouped = new Map<string, IndexedIngredient[]>()
  props.ingredients.forEach((item, index) => {
    const category = categoryFor(item)
    if (!grouped.has(category)) grouped.set(category, [])
    grouped.get(category)!.push({ index, item })
  })
  return Array.from(grouped.entries()).map(([category, items]) => ({ category, items }))
})
const flatItems = computed<IndexedIngredient[]>(() => props.ingredients.map((item, index) => ({ index, item })))
const amountUnitText = (ingredient: KitchenIngredient) => {
  const amount = String(ingredient.amount || '').trim()
  const unit = String(ingredient.unit || '').trim()
  if (!amount && !unit) return ''
  if (amount && unit) return `${amount} ${unit}`
  return amount || unit
}
const noteText = (ingredient: KitchenIngredient) => String(ingredient.note || '').trim()
const ingredientLine = (ingredient: KitchenIngredient) =>
  [ingredient.name, amountUnitText(ingredient), noteText(ingredient)]
    .map(value => String(value || '').trim())
    .filter(Boolean)
    .join(' ')
</script>
