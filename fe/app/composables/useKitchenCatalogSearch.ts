import { computed, ref, watch, type ComputedRef, type Ref } from 'vue'
type MaybeReadonlyRef<T> = Ref<T> | ComputedRef<T>
/**
 * Reusable kitchen search and catalog state for filtering, grouping and favorite-based ingredient scopes.
 */
export const useKitchenCatalogSearch = (options: {
  catalogItems: MaybeReadonlyRef<KitchenCatalogItem[]>
  categoryLabels: MaybeReadonlyRef<string[]>
  favoriteIncludeIngredients: Ref<KitchenFavoriteIngredient[]>
  favoriteExcludeIngredients: Ref<KitchenFavoriteIngredient[]>
}) => {
  const categoryFilter = ref('все')
  const ingredientCatalogQuery = ref('')
  const catalogSortMode = ref<'alpha' | 'alphaDesc' | 'freq'>('alpha')
  const groupCatalogByCategory = ref(false)
  const recipeQuery = ref('')
  const searchIngredientInput = ref('')
  const recipeAdvancedSearchOpen = ref(false)
  const showFavoriteRecipesOnly = ref(false)
  const selectedIngredients = ref<string[]>([])
  const excludedIngredients = ref<string[]>([])
  const dismissedFavoriteIncludedInSearch = ref<string[]>([])
  const dismissedFavoriteExcludedInSearch = ref<string[]>([])
  const requireAllIngredients = ref(false)
  const includeFavoriteIngredientsInSearch = ref(false)
  const excludeFavoriteIngredientsInSearch = ref(false)
  const recipeDifficultyFilter = ref('all')
  const recipeMealTypeFilter = ref('all')
  const recipeCuisineFilter = ref('')
  const recipeDietTypeFilter = ref('all')
  const recipeCookingMethodFilter = ref('')
  const recipeKcalMin = ref('')
  const recipeKcalMax = ref('')
  const recipeMaxTotalMinutes = ref('')
  const recipeServingsMin = ref('')
  const recipeServingsMax = ref('')
  const recipeServingsRangeError = ref('')
  const recipeKcalRangeError = ref('')
  const catalogSortItems: LabTabItem[] = [
    { value: 'alpha', label: 'А-Я' },
    { value: 'alphaDesc', label: 'Я-А' },
    { value: 'freq', label: 'Частота' }
  ]
  /** Normalizes a tag-like value for stable comparisons. */
  const normalizeTag = (value: string) =>
    String(value || '')
      .trim()
      .toLowerCase()
  /** Parses a positive integer filter value and returns null when the input is empty or invalid. */
  const normalizePositiveInt = (raw: string) => {
    const parsed = Number.parseInt(String(raw || '').trim(), 10)
    if (!Number.isFinite(parsed) || parsed <= 0) return null
    return parsed
  }
  /** Counts substring matches to support frequency-based catalog sorting. */
  const countOccurrences = (text: string, query: string) => {
    if (!query) return 0
    let count = 0
    let from = 0
    while (from < text.length) {
      const index = text.indexOf(query, from)
      if (index < 0) break
      count += 1
      from = index + query.length
    }
    return count
  }
  /** Applies alphabetical ordering for catalog items. */
  const compareByAlphabet = (a: KitchenCatalogItem, b: KitchenCatalogItem) => a.name.localeCompare(b.name, 'ru-RU')
  /** Applies query-sensitive ordering so the most relevant catalog items surface first. */
  const compareByFrequency = (a: KitchenCatalogItem, b: KitchenCatalogItem, query: string) => {
    const aName = a.name.toLowerCase()
    const bName = b.name.toLowerCase()
    const aCount = countOccurrences(aName, query)
    const bCount = countOccurrences(bName, query)
    if (aCount !== bCount) return bCount - aCount
    const aStarts = aName.startsWith(query)
    const bStarts = bName.startsWith(query)
    if (aStarts !== bStarts) return aStarts ? -1 : 1
    const aIndex = aName.indexOf(query)
    const bIndex = bName.indexOf(query)
    if (aIndex !== bIndex) return aIndex - bIndex
    if (a.name.length !== b.name.length) return a.name.length - b.name.length
    return compareByAlphabet(a, b)
  }
  /** Extracts unique favorite ingredient names while preserving display order. */
  const favoriteNamesFrom = (items: KitchenFavoriteIngredient[]) => {
    const seen = new Set<string>()
    const names: string[] = []
    for (const item of items) {
      const name = String(item?.name || '').trim()
      const key = normalizeTag(name)
      if (!key || seen.has(key)) continue
      seen.add(key)
      names.push(name)
    }
    return names
  }
  /** Counts active favorite names after dismissed entries are removed from the effective filter scope. */
  const countActiveFavoriteNames = (allNames: string[], dismissedNames: string[]) => {
    const dismissedSet = new Set(dismissedNames.map((item) => normalizeTag(item)))
    return allNames.filter((name) => !dismissedSet.has(normalizeTag(name))).length
  }
  /** Returns ingredient suggestions while excluding items already selected in the current filter state. */
  const getIngredientSuggestions = (query: string, selectedValues: string[]) => {
    const normalizedQuery = query.trim().toLowerCase()
    if (!normalizedQuery) return [] as Array<{ name: string; category: string }>
    const selected = new Set(selectedValues.map((value) => normalizeTag(value)))
    const seen = new Set<string>()
    const items = options.catalogItems.value
    const result: Array<{ name: string; category: string }> = []
    for (const item of items) {
      const normalizedName = normalizeTag(item.name)
      if (!normalizedName || seen.has(normalizedName)) continue
      seen.add(normalizedName)
      if (!normalizedName.includes(normalizedQuery)) continue
      if (selected.has(normalizedName)) continue
      result.push({ name: item.name, category: item.category })
    }
    return result
  }
  /** Adds an ingredient to the include-filter list and removes duplicates from exclusions. */
  const addSearchIngredient = (name: string) => {
    const trimmed = String(name || '').trim()
    if (!trimmed) return
    dismissedFavoriteIncludedInSearch.value = dismissedFavoriteIncludedInSearch.value.filter(
      (item) => normalizeTag(item) !== normalizeTag(trimmed)
    )
    if (!selectedIngredients.value.some((item) => normalizeTag(item) === normalizeTag(trimmed))) {
      selectedIngredients.value = [...selectedIngredients.value, trimmed]
    }
    excludedIngredients.value = excludedIngredients.value.filter((item) => normalizeTag(item) !== normalizeTag(trimmed))
    searchIngredientInput.value = ''
  }
  /** Removes a manually selected include-filter ingredient. */
  const removeSearchIngredient = (name: string) => {
    selectedIngredients.value = selectedIngredients.value.filter((item) => normalizeTag(item) !== normalizeTag(name))
  }
  /** Adds an ingredient to the exclusion list and removes duplicates from includes. */
  const addExcludedIngredient = (name: string) => {
    const trimmed = String(name || '').trim()
    if (!trimmed) return
    dismissedFavoriteExcludedInSearch.value = dismissedFavoriteExcludedInSearch.value.filter(
      (item) => normalizeTag(item) !== normalizeTag(trimmed)
    )
    if (!excludedIngredients.value.some((item) => normalizeTag(item) === normalizeTag(trimmed))) {
      excludedIngredients.value = [...excludedIngredients.value, trimmed]
    }
    selectedIngredients.value = selectedIngredients.value.filter((item) => normalizeTag(item) !== normalizeTag(trimmed))
    searchIngredientInput.value = ''
  }
  /** Removes a manually excluded ingredient. */
  const removeExcludedIngredient = (name: string) => {
    excludedIngredients.value = excludedIngredients.value.filter((item) => normalizeTag(item) !== normalizeTag(name))
  }
  /** Adds the current typed ingredient to includes. */
  const addSearchIngredientFromInput = () => addSearchIngredient(searchIngredientInput.value)
  /** Adds the current typed ingredient to exclusions. */
  const addExcludedIngredientFromInput = () => addExcludedIngredient(searchIngredientInput.value)
  /** Marks a favorite include ingredient as dismissed from the effective search scope. */
  const dismissFavoriteIncludedForSearch = (name: string) => {
    if (dismissedFavoriteIncludedInSearch.value.some((item) => normalizeTag(item) === normalizeTag(name))) return
    dismissedFavoriteIncludedInSearch.value = [...dismissedFavoriteIncludedInSearch.value, name]
  }
  /** Marks a favorite exclude ingredient as dismissed from the effective search scope. */
  const dismissFavoriteExcludedForSearch = (name: string) => {
    if (dismissedFavoriteExcludedInSearch.value.some((item) => normalizeTag(item) === normalizeTag(name))) return
    dismissedFavoriteExcludedInSearch.value = [...dismissedFavoriteExcludedInSearch.value, name]
  }
  /** Resets the full advanced recipe search state to its initial values. */
  const resetRecipeSearchFilters = () => {
    recipeQuery.value = ''
    searchIngredientInput.value = ''
    showFavoriteRecipesOnly.value = false
    selectedIngredients.value = []
    excludedIngredients.value = []
    dismissedFavoriteIncludedInSearch.value = []
    dismissedFavoriteExcludedInSearch.value = []
    requireAllIngredients.value = false
    includeFavoriteIngredientsInSearch.value = false
    excludeFavoriteIngredientsInSearch.value = false
    recipeDifficultyFilter.value = 'all'
    recipeMealTypeFilter.value = 'all'
    recipeCuisineFilter.value = ''
    recipeDietTypeFilter.value = 'all'
    recipeCookingMethodFilter.value = ''
    recipeKcalMin.value = ''
    recipeKcalMax.value = ''
    recipeMaxTotalMinutes.value = ''
    recipeServingsMin.value = ''
    recipeServingsMax.value = ''
    recipeServingsRangeError.value = ''
    recipeKcalRangeError.value = ''
  }
  const availableCategories = computed(() => ['все', ...options.categoryLabels.value])
  const ingredientCategoryOptions = computed(() =>
    availableCategories.value.map((category) => ({
      value: category,
      label: category
    }))
  )
  const filteredCatalog = computed(() => {
    const query = ingredientCatalogQuery.value.trim().toLowerCase()
    return options.catalogItems.value.filter((item) => {
      const byCategory = categoryFilter.value === 'все' || item.category === categoryFilter.value
      if (!byCategory) return false
      if (!query) return true
      return item.name.toLowerCase().includes(query)
    })
  })
  const catalogSortLabel = computed(() => {
    if (catalogSortMode.value === 'alphaDesc') return 'Я-А'
    if (catalogSortMode.value === 'freq') return 'Частота'
    return 'А-Я'
  })
  const catalogFrequencyHint = computed(() => {
    if (catalogSortMode.value !== 'freq') return ''
    const query = ingredientCatalogQuery.value.trim()
    if (!query) return 'Частотный режим: без запроса сортировка выполняется по алфавиту.'
    return `Частотный режим: выше показываются ингредиенты, где «${query}» встречается чаще и ближе к началу названия.`
  })
  const sortedCatalog = computed(() => {
    const query = ingredientCatalogQuery.value.trim().toLowerCase()
    const items = [...filteredCatalog.value]
    items.sort((a, b) => {
      if (catalogSortMode.value === 'alphaDesc') return compareByAlphabet(b, a)
      if (catalogSortMode.value === 'freq') return compareByFrequency(a, b, query)
      return compareByAlphabet(a, b)
    })
    return items
  })
  const visibleCatalogGroups = computed(() => {
    if (!groupCatalogByCategory.value) {
      return [{ category: '', items: sortedCatalog.value }]
    }
    const byCategory = new Map<string, KitchenCatalogItem[]>()
    for (const item of sortedCatalog.value) {
      const category = item.category || 'другое'
      const current = byCategory.get(category) || []
      current.push(item)
      byCategory.set(category, current)
    }
    const known = options.categoryLabels.value.filter((category) => byCategory.has(category))
    const unknown = Array.from(byCategory.keys())
      .filter((category) => !options.categoryLabels.value.includes(category))
      .sort((a, b) => a.localeCompare(b, 'ru-RU'))
    return [...known, ...unknown].map((category) => ({
      category,
      items: byCategory.get(category) || []
    }))
  })
  const ingredientSuggestions = computed(() =>
    getIngredientSuggestions(searchIngredientInput.value, [...selectedIngredients.value, ...excludedIngredients.value])
  )
  const favoriteIncludeIngredientNames = computed(() => favoriteNamesFrom(options.favoriteIncludeIngredients.value))
  const favoriteExcludeIngredientNames = computed(() => favoriteNamesFrom(options.favoriteExcludeIngredients.value))
  const isFavoriteIncludedDismissed = (name: string) =>
    dismissedFavoriteIncludedInSearch.value.some((item) => normalizeTag(item) === normalizeTag(name))
  const isFavoriteExcludedDismissed = (name: string) =>
    dismissedFavoriteExcludedInSearch.value.some((item) => normalizeTag(item) === normalizeTag(name))
  const isFavoriteIncludedName = (name: string) =>
    favoriteIncludeIngredientNames.value.some((item) => normalizeTag(item) === normalizeTag(name))
  const isFavoriteExcludedName = (name: string) =>
    favoriteExcludeIngredientNames.value.some((item) => normalizeTag(item) === normalizeTag(name))
  const activeFavoriteIncludeCount = computed(() =>
    countActiveFavoriteNames(favoriteIncludeIngredientNames.value, dismissedFavoriteIncludedInSearch.value)
  )
  const activeFavoriteExcludeCount = computed(() =>
    countActiveFavoriteNames(favoriteExcludeIngredientNames.value, dismissedFavoriteExcludedInSearch.value)
  )
  const includeFavoriteToggleVisualState = computed(() => {
    if (!includeFavoriteIngredientsInSearch.value) return 'off'
    const total = favoriteIncludeIngredientNames.value.length
    const active = activeFavoriteIncludeCount.value
    if (total <= 0 || active <= 0) return 'off'
    if (active < total) return 'partial'
    return 'on'
  })
  const excludeFavoriteToggleVisualState = computed(() => {
    if (!excludeFavoriteIngredientsInSearch.value) return 'off'
    const total = favoriteExcludeIngredientNames.value.length
    const active = activeFavoriteExcludeCount.value
    if (total <= 0 || active <= 0) return 'off'
    if (active < total) return 'partial'
    return 'on'
  })
  const excludeFavoriteSwitchLabel = computed(() =>
    excludeFavoriteToggleVisualState.value === 'off'
      ? 'Добавить из списка исключаемых'
      : 'Удалить все из списка исключаемых'
  )
  const servingsRange = computed(() => ({
    min: normalizePositiveInt(recipeServingsMin.value),
    max: normalizePositiveInt(recipeServingsMax.value)
  }))
  const isServingsRangeInvalid = computed(() => {
    const { min, max } = servingsRange.value
    return min !== null && max !== null && min > max
  })
  const maxTotalMinutesFilter = computed(() => normalizePositiveInt(recipeMaxTotalMinutes.value))
  const kcalRange = computed(() => ({
    min: normalizePositiveInt(recipeKcalMin.value),
    max: normalizePositiveInt(recipeKcalMax.value)
  }))
  const isKcalRangeInvalid = computed(() => {
    const { min, max } = kcalRange.value
    return min !== null && max !== null && min > max
  })
  const hasMetaRecipeFilters = computed(
    () =>
      recipeDifficultyFilter.value !== 'all' ||
      recipeMealTypeFilter.value !== 'all' ||
      recipeDietTypeFilter.value !== 'all' ||
      Boolean(recipeCookingMethodFilter.value.trim()) ||
      Boolean(recipeCuisineFilter.value.trim()) ||
      kcalRange.value.min !== null ||
      kcalRange.value.max !== null ||
      maxTotalMinutesFilter.value !== null ||
      servingsRange.value.min !== null ||
      servingsRange.value.max !== null
  )
  const resolveSearchIngredients = () => {
    const selected = [...selectedIngredients.value]
    const excluded = [...excludedIngredients.value]
    if (includeFavoriteIngredientsInSearch.value) {
      for (const name of favoriteIncludeIngredientNames.value) {
        if (isFavoriteIncludedDismissed(name)) continue
        if (!selected.some((item) => normalizeTag(item) === normalizeTag(name))) {
          selected.push(name)
        }
      }
    }
    if (excludeFavoriteIngredientsInSearch.value) {
      for (const name of favoriteExcludeIngredientNames.value) {
        if (isFavoriteExcludedDismissed(name)) continue
        if (!excluded.some((item) => normalizeTag(item) === normalizeTag(name))) {
          excluded.push(name)
        }
      }
    }
    const excludedNormalized = new Set(excluded.map((item) => normalizeTag(item)))
    return {
      selected: selected.filter((item) => !excludedNormalized.has(normalizeTag(item))),
      excluded
    }
  }
  const resolvedSearchFilters = computed(() => resolveSearchIngredients())
  const displaySelectedIngredients = computed(() => resolvedSearchFilters.value.selected)
  const displayExcludedIngredients = computed(() => resolvedSearchFilters.value.excluded)
  const isManuallySelectedIngredient = (name: string) =>
    selectedIngredients.value.some((item) => normalizeTag(item) === normalizeTag(name))
  const isManuallyExcludedIngredient = (name: string) =>
    excludedIngredients.value.some((item) => normalizeTag(item) === normalizeTag(name))
  const removeEffectiveSelectedIngredient = (name: string) => {
    if (isManuallySelectedIngredient(name)) removeSearchIngredient(name)
    if (includeFavoriteIngredientsInSearch.value && isFavoriteIncludedName(name)) {
      dismissFavoriteIncludedForSearch(name)
    }
  }
  const removeEffectiveExcludedIngredient = (name: string) => {
    if (isManuallyExcludedIngredient(name)) removeExcludedIngredient(name)
    if (excludeFavoriteIngredientsInSearch.value && isFavoriteExcludedName(name)) {
      dismissFavoriteExcludedForSearch(name)
    }
  }
  watch(includeFavoriteIngredientsInSearch, (enabled) => {
    if (enabled) return
    dismissedFavoriteIncludedInSearch.value = []
  })
  watch(excludeFavoriteIngredientsInSearch, (enabled) => {
    if (enabled) return
    dismissedFavoriteExcludedInSearch.value = []
  })
  watch(activeFavoriteIncludeCount, (count) => {
    if (!includeFavoriteIngredientsInSearch.value || count > 0) return
    includeFavoriteIngredientsInSearch.value = false
  })
  watch(activeFavoriteExcludeCount, (count) => {
    if (!excludeFavoriteIngredientsInSearch.value || count > 0) return
    excludeFavoriteIngredientsInSearch.value = false
  })
  return {
    categoryFilter,
    ingredientCatalogQuery,
    catalogSortMode,
    groupCatalogByCategory,
    recipeQuery,
    searchIngredientInput,
    recipeAdvancedSearchOpen,
    showFavoriteRecipesOnly,
    selectedIngredients,
    excludedIngredients,
    requireAllIngredients,
    includeFavoriteIngredientsInSearch,
    excludeFavoriteIngredientsInSearch,
    recipeDifficultyFilter,
    recipeMealTypeFilter,
    recipeCuisineFilter,
    recipeDietTypeFilter,
    recipeCookingMethodFilter,
    recipeKcalMin,
    recipeKcalMax,
    recipeMaxTotalMinutes,
    recipeServingsMin,
    recipeServingsMax,
    recipeServingsRangeError,
    recipeKcalRangeError,
    availableCategories,
    ingredientCategoryOptions,
    filteredCatalog,
    catalogSortItems,
    catalogSortLabel,
    catalogFrequencyHint,
    visibleCatalogGroups,
    ingredientSuggestions,
    favoriteIncludeIngredientNames,
    favoriteExcludeIngredientNames,
    activeFavoriteIncludeCount,
    activeFavoriteExcludeCount,
    includeFavoriteToggleVisualState,
    excludeFavoriteToggleVisualState,
    excludeFavoriteSwitchLabel,
    servingsRange,
    isServingsRangeInvalid,
    maxTotalMinutesFilter,
    kcalRange,
    isKcalRangeInvalid,
    hasMetaRecipeFilters,
    resolvedSearchFilters,
    displaySelectedIngredients,
    displayExcludedIngredients,
    normalizeTag,
    normalizePositiveInt,
    getIngredientSuggestions,
    addSearchIngredient,
    removeSearchIngredient,
    addSearchIngredientFromInput,
    addExcludedIngredient,
    removeExcludedIngredient,
    addExcludedIngredientFromInput,
    dismissFavoriteIncludedForSearch,
    dismissFavoriteExcludedForSearch,
    removeEffectiveSelectedIngredient,
    removeEffectiveExcludedIngredient,
    resetRecipeSearchFilters
  }
}
