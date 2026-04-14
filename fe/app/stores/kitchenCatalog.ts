export const useKitchenCatalogStore = defineStore('kitchenCatalog', () => {
  const categories = ref<KitchenCategory[]>([])
  const ingredients = ref<KitchenCatalogIngredient[]>([])
  const filterOptions = ref<KitchenFilterOption[]>([])
  const loaded = ref(false)
  const loading = ref(false)
  const error = ref<string | null>(null)
  let loadPromise: Promise<void> | null = null
  const categoryLabels = computed(() => categories.value.map((item) => item.label))
  const ingredientItems = computed(() =>
    ingredients.value.map((item) => ({
      id: `catalog:${item.ingredient_id}`,
      ingredient_id: item.ingredient_id,
      name: item.name,
      category: item.category_label,
      description: item.description,
      protein_g: item.protein_g ?? null,
      fat_g: item.fat_g ?? null,
      carbs_g: item.carbs_g ?? null,
      kcal: item.kcal ?? null
    }))
  )
  const filterOptionsByKind = computed<Record<string, KitchenFilterOption[]>>(() => {
    const grouped: Record<string, KitchenFilterOption[]> = {}
    for (const item of filterOptions.value) {
      const kind = String(item.kind || '').trim()
      if (!kind) continue
      if (!grouped[kind]) grouped[kind] = []
      grouped[kind].push(item)
    }
    return grouped
  })
  const getOptions = (kind: string) => filterOptionsByKind.value[String(kind || '').trim()] || []
  const labelFor = (kind: string, value?: string | null) => {
    const raw = String(value || '').trim()
    if (!raw) return ''
    const normalized = raw.toLowerCase()
    const match = getOptions(kind).find(
      (option) => option.code.toLowerCase() === normalized || option.label.toLowerCase() === normalized
    )
    return match?.label || raw
  }
  const normalizeOptionValue = (kind: string, value?: string | null) => {
    const raw = String(value || '').trim()
    if (!raw) return ''
    const normalized = raw.toLowerCase()
    const match = getOptions(kind).find(
      (option) => option.code.toLowerCase() === normalized || option.label.toLowerCase() === normalized
    )
    return match?.code || raw
  }
  const load = async (force = false) => {
    if (loaded.value && !force) return
    if (loadPromise) {
      await loadPromise
      return
    }
    loadPromise = (async () => {
      loading.value = true
      error.value = null
      try {
        const res = await getKitchenCatalog()
        categories.value = res.data.categories || []
        ingredients.value = res.data.ingredients || []
        filterOptions.value = res.data.filter_options || []
        loaded.value = true
      } catch (err: any) {
        error.value = err?.data?.message || err?.message || 'Не удалось загрузить каталог кухни.'
        if (!loaded.value) {
          categories.value = []
          ingredients.value = []
          filterOptions.value = []
        }
      } finally {
        loading.value = false
      }
    })()
    try {
      await loadPromise
    } finally {
      loadPromise = null
    }
  }
  const ensureLoaded = async () => {
    await load(false)
  }
  return {
    categories,
    ingredients,
    filterOptions,
    loaded,
    loading,
    error,
    categoryLabels,
    ingredientItems,
    filterOptionsByKind,
    getOptions,
    labelFor,
    normalizeOptionValue,
    load,
    ensureLoaded
  }
})
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useKitchenCatalogStore, import.meta.hot))
}
