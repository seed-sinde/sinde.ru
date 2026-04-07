<script setup lang="ts">
  const title = 'Кухня'
  const runtimeConfig = useRuntimeConfig()
  const { effectiveTheme } = useInterfacePreferences()
  usePageSeo({
    title,
    description: 'Каталог рецептов, ингредиентов, поиск и публикация рецептов.'
  })
  const activeKitchenTab = ref<KitchenMainTab>('recipes')
  const kitchenBaseTabItems: LabTabItem[] = [
    { value: 'recipes', label: 'Рецепты' },
    { value: 'ingredients', label: 'Ингредиенты' },
    { value: 'my-recipes', label: 'Сохранённые рецепты' }
  ]
  const RECIPE_TIME_FILTER_MAX = 360
  const RECIPE_AUTO_SEARCH_DELAY_MS = 300
  const searchUiHydrated = ref(false)
  const queuedAutoSearchAfterLoading = ref(false)
  const recipeToolbarActionLoading = ref<RecipeToolbarAction | null>(null)
  const recipes = ref<KitchenRecipe[]>([])
  const publicRecipesTotal = ref<number | null>(null)
  const recipesLoading = ref(false)
  const recipesError = ref<string | null>(null)
  const resultMode = ref<'latest' | 'search'>('latest')
  const isCompactRecipeSearch = ref(false)
  const recipeFavoriteMap = reactive<Record<string, boolean>>({})
  const recipeFavoritePendingMap = reactive<Record<string, boolean>>({})
  const deletePendingRecipeId = ref<string | null>(null)
  const myRecipes = ref<KitchenRecipe[]>([])
  const myRecipesLoading = ref(false)
  const myRecipesError = ref<string | null>(null)
  const userAttentionLoading = ref(false)
  const userAttentionReadPending = ref(false)
  const accountIngredientsLoading = ref(false)
  const accountIngredientsError = ref<string | null>(null)
  const accountStateInitialized = ref(false)
  const customIngredients = ref<KitchenUserIngredient[]>([])
  const favoriteIncludeIngredients = ref<KitchenFavoriteIngredient[]>([])
  const favoriteExcludeIngredients = ref<KitchenFavoriteIngredient[]>([])
  const customIngredientPending = ref(false)
  const deletingCustomIngredientId = ref<string | null>(null)
  const deletingFavoriteKey = ref<string | null>(null)
  const customIngredientForm = reactive({
    name: '',
    category: 'другое'
  })
  const {
    user,
    isAuthenticated,
    ensureLoaded,
    loaded: authLoaded,
    sharedUserSummary,
    loadSharedUserSummary,
    markUserSummaryRead
  } = useAuth()
  const isCurrentUserAdmin = computed(() => Boolean(user.value?.roles?.includes('admin')))
  const catalogStore = useKitchenCatalogStore()
  await callOnce('kitchen-catalog', async () => {
    await catalogStore.ensureLoaded()
  })
  const userAttentionSummary = computed(() => sharedUserSummary.value)
  const difficultyOptions = computed(() => catalogStore.getOptions('difficulty'))
  const defaultRecipeDifficulty = computed(() => difficultyOptions.value[0]?.code || 'easy')
  const kitchenImageMaxBytes = Number(runtimeConfig.public.mediaImageMaxBytes || 8388608)
  const KITCHEN_IMAGE_MAX_MB = Math.max(1, Math.ceil(kitchenImageMaxBytes / (1024 * 1024)))
  const KITCHEN_IMAGE_FORMATS_LABEL = 'jpg, jpeg, png, webp, gif, avif'
  const KITCHEN_STEP_IMAGE_RECOMMENDED_MIN_WIDTH = 1200
  const KITCHEN_STEP_IMAGE_RECOMMENDED_MIN_HEIGHT = 900
  const catalogIngredients = computed<KitchenCatalogItem[]>(() => {
    const merged: KitchenCatalogItem[] = catalogStore.ingredientItems.map(item => ({ ...item }))
    const seen = new Set(merged.map(item => `${item.category}:${item.name.toLowerCase()}`))
    for (const item of customIngredients.value) {
      const key = `${item.category}:${item.name.toLowerCase()}`
      if (seen.has(key)) continue
      seen.add(key)
      const customItem: KitchenCatalogItem = {
        id: `custom:${item.ingredient_id}`,
        name: item.name,
        category: item.category,
        description: '',
        protein_g: null,
        fat_g: null,
        carbs_g: null,
        kcal: null
      }
      merged.push(customItem)
    }
    return merged
  })
  const {
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
    ingredientCategoryOptions,
    filteredCatalog,
    catalogSortItems,
    catalogSortLabel,
    catalogFrequencyHint,
    visibleCatalogGroups,
    ingredientSuggestions,
    favoriteIncludeIngredientNames,
    favoriteExcludeIngredientNames,
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
    displaySelectedIngredients: resolvedDisplaySelectedIngredients,
    displayExcludedIngredients: resolvedDisplayExcludedIngredients,
    normalizeTag,
    getIngredientSuggestions,
    addSearchIngredientFromInput,
    addExcludedIngredientFromInput,
    removeEffectiveSelectedIngredient,
    removeEffectiveExcludedIngredient,
    resetRecipeSearchFilters
  } = useKitchenCatalogSearch({
    catalogItems: catalogIngredients,
    categoryLabels: computed(() => catalogStore.categoryLabels),
    favoriteIncludeIngredients,
    favoriteExcludeIngredients
  })
  const customIngredientCategoryOptions = computed<SelectOptionInput[]>(() =>
    catalogStore.categoryLabels.map(category => ({
      value: category,
      label: category
    }))
  )
  const visibleRecipes = computed(() => {
    if (!showFavoriteRecipesOnly.value) return recipes.value
    return recipes.value.filter(recipe => isRecipeFavorite(recipe.id))
  })
  const favoriteIncludeIngredientIds = computed(
    () => new Set(favoriteIncludeIngredients.value.map(item => item.ingredient_id))
  )
  const favoriteExcludeIngredientIds = computed(
    () => new Set(favoriteExcludeIngredients.value.map(item => item.ingredient_id))
  )
  const showCatalogFavoriteActions = computed(() => searchUiHydrated.value && isAuthenticated.value)
  const canManageRecipe = (recipe: KitchenRecipe) =>
    Boolean(
      user.value?.user_id &&
      ((recipe.owner_user_id && user.value.user_id === recipe.owner_user_id) || user.value.roles?.includes('admin'))
    )
  const isFavoriteIngredient = (ingredientId?: number) =>
    typeof ingredientId === 'number' && ingredientId > 0 && favoriteIncludeIngredientIds.value.has(ingredientId)
  const isExcludeFavoriteIngredient = (ingredientId?: number) =>
    typeof ingredientId === 'number' && ingredientId > 0 && favoriteExcludeIngredientIds.value.has(ingredientId)
  const upsertFavoriteIngredient = (item: KitchenFavoriteIngredient) => {
    const listType = normalizeKitchenFavoriteListType(item?.list_type)
    if (listType === 'exclude') {
      favoriteExcludeIngredients.value = upsertKitchenFavoriteIngredient(
        favoriteExcludeIngredients.value,
        item,
        'exclude'
      )
      return
    }
    favoriteIncludeIngredients.value = upsertKitchenFavoriteIngredient(
      favoriteIncludeIngredients.value,
      item,
      'include'
    )
  }
  const removeFavoriteIngredientById = (ingredientId: number, listType: 'include' | 'exclude') => {
    const items = listType === 'exclude' ? favoriteExcludeIngredients : favoriteIncludeIngredients
    items.value = removeKitchenFavoriteIngredientById(items.value, ingredientId)
  }
  const loadMyRecipes = async () => {
    if (!isAuthenticated.value) {
      myRecipes.value = []
      myRecipesError.value = null
      return
    }
    myRecipesLoading.value = true
    myRecipesError.value = null
    try {
      const res = await getMyKitchenRecipes()
      myRecipes.value = res?.data?.items || []
    } catch (err: any) {
      myRecipesError.value = err?.data?.message || err?.message || 'Не удалось загрузить ваши рецепты.'
      myRecipes.value = []
    } finally {
      myRecipesLoading.value = false
    }
  }
  const loadUserAttentionSummary = async () => {
    if (!isAuthenticated.value) {
      return
    }
    if (userAttentionLoading.value) return
    userAttentionLoading.value = true
    try {
      await loadSharedUserSummary()
    } catch {
      // Keep the last shared snapshot if the refresh failed.
    } finally {
      userAttentionLoading.value = false
    }
  }
  const markUserAttentionRead = async () => {
    if (userAttentionReadPending.value) return
    userAttentionReadPending.value = true
    try {
      await markUserSummaryRead()
      await loadUserAttentionSummary()
    } catch (err: any) {
      myRecipesError.value = err?.data?.message || err?.message || 'Не удалось отметить уведомления как просмотренные.'
    } finally {
      userAttentionReadPending.value = false
    }
  }
  const loadKitchenAccountIngredients = async () => {
    if (!isAuthenticated.value) {
      customIngredients.value = []
      favoriteIncludeIngredients.value = []
      favoriteExcludeIngredients.value = []
      accountIngredientsError.value = null
      return
    }
    accountIngredientsLoading.value = true
    accountIngredientsError.value = null
    try {
      const res = await getKitchenAccountIngredients()
      customIngredients.value = res?.data?.custom || []
      const includeRaw = res?.data?.favorites_include || res?.data?.favorites || []
      const excludeRaw = res?.data?.favorites_exclude || []
      favoriteIncludeIngredients.value = includeRaw.map(item => normalizeKitchenFavoriteIngredient(item, 'include'))
      favoriteExcludeIngredients.value = excludeRaw.map(item => normalizeKitchenFavoriteIngredient(item, 'exclude'))
    } catch (err: any) {
      accountIngredientsError.value =
        err?.data?.message || err?.message || 'Не удалось загрузить персональные ингредиенты.'
    } finally {
      accountIngredientsLoading.value = false
    }
  }
  const route = useRoute()
  const router = useRouter()
  const loginPath = computed(() => buildLoginPath(route.fullPath))
  const isKitchenEditRoute = computed(() => route.path.startsWith('/kitchen/edit/'))
  const routeEditSlug = computed(() => {
    if (!isKitchenEditRoute.value) return ''
    const byParams = String(route.params.slug || '').trim()
    if (byParams) return byParams
    return extractKitchenEditSlug(route.path)
  })
  const routeEditRecipeId = computed(() => extractKitchenRecipeId(routeEditSlug.value))
  const requestedKitchenTab = computed(() => {
    const raw = Array.isArray(route.query.tab) ? route.query.tab[0] : route.query.tab
    return String(raw || '').trim()
  })
  const loadingRouteRecipeId = ref<string | null>(null)
  const {
    data: routeEditRecipeData,
    error: routeEditRecipeError,
    pending: routeEditRecipePending
  } = await useAsyncData(
    () => `kitchen-edit-recipe:${isKitchenEditRoute.value ? routeEditRecipeId.value || 'invalid' : 'none'}`,
    async () => {
      if (!isKitchenEditRoute.value) return null
      if (!routeEditRecipeId.value) return null
      const res = await getKitchenManageRecipeById(routeEditRecipeId.value)
      return res?.data || null
    },
    { watch: [isKitchenEditRoute, routeEditRecipeId] }
  )
  const kitchenTabItems = computed<LabTabItem[]>(() => {
    const items: LabTabItem[] = [...kitchenBaseTabItems]
    if (!isKitchenEditRoute.value && (activeKitchenTab.value === 'create' || requestedKitchenTab.value === 'create')) {
      items.push({ value: 'create', label: 'Новый рецепт' })
    }
    if (isKitchenEditRoute.value || activeKitchenTab.value === 'edit') {
      items.push({ value: 'edit', label: 'Редактирование' })
    }
    return items
  })
  const kitchenTabRouteMap = computed<Record<string, string | undefined>>(() => ({
    edit: isKitchenEditRoute.value ? route.fullPath : undefined
  }))
  const goToKitchenTab = async (tab: Exclude<KitchenMainTab, 'edit'>) => {
    const nextQuery: Record<string, any> = { ...route.query }
    if (tab === 'recipes') {
      delete nextQuery.tab
    } else {
      nextQuery.tab = tab
    }
    await router.replace({ path: '/kitchen', query: nextQuery })
  }
  const openCreateRecipeTab = async () => {
    activeKitchenTab.value = 'create'
    await goToKitchenTab('create')
  }
  const openMyRecipesTab = async () => {
    activeKitchenTab.value = 'my-recipes'
    await goToKitchenTab('my-recipes')
  }
  const recipeEditor = useKitchenRecipeEditor({
    activeKitchenTab,
    isAuthenticated,
    isCurrentUserAdmin,
    defaultRecipeDifficulty,
    isKitchenEditRoute,
    getIngredientSuggestions,
    normalizeMealTypeInput,
    normalizeDietTypeInput,
    normalizeCookingMethodInput,
    onOpenRecipeEdit: async recipe => {
      await router.push(kitchenRecipeEditLink(recipe))
    },
    onResetEditRoute: async () => {
      await goToKitchenTab('my-recipes')
    },
    onAfterSubmit: async () => {
      await Promise.all([loadLatestRecipes(), loadMyRecipes(), loadUserAttentionSummary()])
    }
  })
  const { form, createError, createSuccess, editRecipeId, resetRecipeForm, openRecipeEditing, applyRecipeToForm } =
    recipeEditor
  const recipeImages = useKitchenRecipeImages({
    activeRecipeUploadId: recipeEditor.activeRecipeUploadId,
    formSteps: recipeEditor.formSteps,
    stepDraft: recipeEditor.stepDraft,
    coverImageDraft: recipeEditor.coverImageDraft
  })
  const recipeStepsDnD = useKitchenRecipeStepsDnD({
    formSteps: recipeEditor.formSteps,
    moveFormStep: recipeEditor.moveFormStep,
    buildStepImageUrl: recipeImages.buildStepImageUrl
  })
  recipeEditor.setTransientResetHandler(recipeImages.resetRecipeImageUi)
  watch(
    () =>
      [
        isKitchenEditRoute.value,
        routeEditRecipeId.value,
        isAuthenticated.value,
        routeEditRecipePending.value,
        String(routeEditRecipeData.value?.id || ''),
        Boolean(routeEditRecipeError.value)
      ] as const,
    () => {
      void syncEditRouteRecipe()
    },
    { immediate: true }
  )
  watch(isKitchenEditRoute, (isEdit, wasEdit) => {
    if (isEdit || !wasEdit) return
    if (!editRecipeId.value) return
    resetRecipeForm()
  })
  const showRecipeManageSection = computed(
    () =>
      activeKitchenTab.value === 'my-recipes' ||
      activeKitchenTab.value === 'create' ||
      activeKitchenTab.value === 'edit'
  )
  const showRecipeForm = computed(() => activeKitchenTab.value === 'create' || activeKitchenTab.value === 'edit')
  const showMyRecipesList = computed(() => activeKitchenTab.value === 'my-recipes')
  const visibleFavoriteStatusRecipes = computed(() => {
    const items = [...recipes.value]
    if (showMyRecipesList.value) {
      items.push(...myRecipes.value)
    }
    return items
  })
  const visibleFavoriteStatusRecipeIds = computed(() => {
    const ids = new Set<string>()
    for (const item of visibleFavoriteStatusRecipes.value) {
      const id = String(item?.id || '').trim()
      if (!id) continue
      ids.add(id)
    }
    return Array.from(ids)
  })
  const kitchenBreadcrumbItems = computed<BreadcrumbItem[]>(() => {
    const root =
      activeKitchenTab.value === 'recipes' ? [{ label: 'Кухня', current: true }] : [{ label: 'Кухня', to: '/kitchen' }]
    const recipeTitle = form.title.trim() || 'Рецепт'
    if (activeKitchenTab.value === 'ingredients') {
      return [...root, { label: 'Ингредиенты', current: true, kind: 'tab' }]
    }
    if (activeKitchenTab.value === 'my-recipes') {
      return [...root, { label: 'Сохранённые рецепты', current: true, kind: 'tab' }]
    }
    if (activeKitchenTab.value === 'create') {
      return [...root, { label: 'Новый рецепт', current: true, kind: 'tab' }]
    }
    if (activeKitchenTab.value === 'edit') {
      return [...root, { label: 'Редактирование рецепта', kind: 'tab' }, { label: recipeTitle, current: true }]
    }
    return root
  })
  const activeSearchSummary = computed(() => {
    const count = visibleRecipes.value.length
    const favoriteOnlyText = showFavoriteRecipesOnly.value ? ' Включен фильтр: только избранные.' : ''
    if (resultMode.value === 'latest')
      return `Показаны последние публичные рецепты. Найдено рецептов: ${count}.${favoriteOnlyText}`
    const resolved = resolvedSearchFilters.value
    const modeText = requireAllIngredients.value ? 'все выбранные ингредиенты' : 'любой из выбранных ингредиентов'
    const excludedText = resolved.excluded.length > 0 ? ` Исключено ингредиентов: ${resolved.excluded.length}.` : ''
    const favoriteScope: string[] = []
    if (includeFavoriteIngredientsInSearch.value) favoriteScope.push('любимые добавлены')
    if (excludeFavoriteIngredientsInSearch.value) favoriteScope.push('исключаемые добавлены')
    const favoriteText = favoriteScope.length ? ` ${favoriteScope.join(', ')}.` : ''
    return `Режим поиска: ${modeText}.${excludedText}${favoriteText}${favoriteOnlyText} Найдено рецептов: ${count}.`
  })
  const recipeAdvancedSearchContainerClass = computed(() =>
    recipeAdvancedSearchOpen.value
      ? 'space-y-0 fixed inset-0 z-60 bg-zinc-950 md:static md:z-auto md:bg-transparent'
      : 'space-y-0'
  )
  const recipeAdvancedSearchContentClass = computed(() =>
    recipeAdvancedSearchOpen.value
      ? 'flex h-full min-h-0 flex-col gap-4 overflow-y-auto px-3 pb-4 pt-4 md:block md:space-y-3 md:px-0 md:pb-2 md:pt-3'
      : 'space-y-3'
  )
  const recipeAdvancedSearchFoundLabel = computed(() => {
    if (recipesLoading.value) return 'Ищем рецепты...'
    return `Найдено рецептов: ${visibleRecipes.value.length}`
  })
  const recipeAdvancedSearchCloseLabel = computed(() => {
    if (recipesLoading.value) return 'Вернуться к результатам'
    return `Показать рецепты: ${visibleRecipes.value.length}`
  })
  const closeRecipeAdvancedSearch = () => {
    recipeAdvancedSearchOpen.value = false
  }
  let recipeAdvancedSearchMediaQuery: MediaQueryList | null = null
  const syncRecipeAdvancedSearchViewport = (event?: MediaQueryList | MediaQueryListEvent) => {
    if (event) {
      isCompactRecipeSearch.value = event.matches
      return
    }
    isCompactRecipeSearch.value = Boolean(recipeAdvancedSearchMediaQuery?.matches)
  }
  onMounted(() => {
    if (!import.meta.client) return
    recipeAdvancedSearchMediaQuery = window.matchMedia('(max-width: 767px)')
    syncRecipeAdvancedSearchViewport()
    recipeAdvancedSearchMediaQuery.addEventListener('change', syncRecipeAdvancedSearchViewport)
  })
  onBeforeUnmount(() => {
    recipeAdvancedSearchMediaQuery?.removeEventListener('change', syncRecipeAdvancedSearchViewport)
  })
  const clearRecipeFavoriteState = () => {
    for (const key of Object.keys(recipeFavoriteMap)) {
      delete recipeFavoriteMap[key]
    }
    for (const key of Object.keys(recipeFavoritePendingMap)) {
      delete recipeFavoritePendingMap[key]
    }
  }
  const isRecipeFavorite = (recipeId?: string) => {
    const id = String(recipeId || '').trim()
    if (!id) return false
    return Boolean(recipeFavoriteMap[id])
  }
  const isRecipeFavoritePending = (recipeId?: string) => {
    const id = String(recipeId || '').trim()
    if (!id) return false
    return Boolean(recipeFavoritePendingMap[id])
  }
  const syncVisibleRecipeFavoriteStatuses = async () => {
    if (!import.meta.client) return
    if (!isAuthenticated.value) {
      clearRecipeFavoriteState()
      return
    }
    const ids = visibleFavoriteStatusRecipeIds.value
    const visibleIdSet = new Set(ids)
    for (const key of Object.keys(recipeFavoriteMap)) {
      if (!visibleIdSet.has(key) && !recipeFavoritePendingMap[key]) {
        delete recipeFavoriteMap[key]
      }
    }
    if (!ids.length) return
    await Promise.all(
      ids.map(async id => {
        try {
          const res = await getKitchenRecipeFavoriteStatus(id)
          recipeFavoriteMap[id] = Boolean(res?.data?.favorited)
        } catch {
          recipeFavoriteMap[id] = false
        }
      })
    )
  }
  const toggleRecipeFavorite = async (recipe: KitchenRecipe) => {
    const id = String(recipe?.id || '').trim()
    if (!id) return
    if (!isAuthenticated.value) {
      await router.push(loginPath.value)
      return
    }
    if (recipeFavoritePendingMap[id]) return
    recipeFavoritePendingMap[id] = true
    recipesError.value = null
    try {
      if (isRecipeFavorite(id)) {
        await unfavoriteKitchenRecipe(id)
        recipeFavoriteMap[id] = false
      } else {
        await favoriteKitchenRecipe(id)
        recipeFavoriteMap[id] = true
      }
    } catch (err: any) {
      recipesError.value = err?.data?.message || err?.message || 'Не удалось обновить избранное для рецепта.'
    } finally {
      delete recipeFavoritePendingMap[id]
    }
  }
  const toggleFavoriteRecipesOnlyFilter = async () => {
    if (!isAuthenticated.value) {
      await router.push(loginPath.value)
      return
    }
    showFavoriteRecipesOnly.value = !showFavoriteRecipesOnly.value
  }
  const toggleFavoriteRecipesOnlyFilterFromToolbar = async () => {
    if (recipeToolbarActionLoading.value) return
    recipeToolbarActionLoading.value = 'favorite'
    try {
      await toggleFavoriteRecipesOnlyFilter()
    } finally {
      if (recipeToolbarActionLoading.value === 'favorite') {
        recipeToolbarActionLoading.value = null
      }
    }
  }
  const ingredientCategoryByName = computed(() => {
    const map: Record<string, string> = {}
    for (const item of catalogIngredients.value) {
      const key = normalizeTag(item.name)
      if (!map[key]) {
        map[key] = item.category
      }
    }
    return map
  })
  const ingredientCategory = (name: string) => ingredientCategoryByName.value[normalizeTag(name)] || 'другое'
  const categoryTagStyle = (category: string) => {
    const color = kitchenCategoryColor(category)
    const isLight = effectiveTheme.value === 'light'
    return {
      borderColor: kitchenColorWithAlpha(color, isLight ? 0.42 : 0.58),
      backgroundColor: kitchenColorWithAlpha(color, isLight ? 0.18 : 0.24),
      color: kitchenColorWithAlpha(color, isLight ? 0.82 : 0.98)
    }
  }
  const selectedTagStyle = (name: string) => categoryTagStyle(ingredientCategory(name))
  const customIngredientChipItems = computed(() =>
    toKitchenChipItems(customIngredients.value, 'custom:').map(item => ({
      ...item,
      disabled: deletingCustomIngredientId.value === String((item.payload as KitchenUserIngredient).ingredient_id || '')
    }))
  )
  const favoriteIncludeChipItems = computed(() =>
    toKitchenChipItems(favoriteIncludeIngredients.value, 'favorite:include:').map(item => ({
      ...item,
      disabled:
        deletingFavoriteKey.value ===
        `include:${String((item.payload as KitchenFavoriteIngredient).ingredient_id || '')}`
    }))
  )
  const favoriteExcludeChipItems = computed(() =>
    toKitchenChipItems(favoriteExcludeIngredients.value, 'favorite:exclude:').map(item => ({
      ...item,
      disabled:
        deletingFavoriteKey.value ===
        `exclude:${String((item.payload as KitchenFavoriteIngredient).ingredient_id || '')}`
    }))
  )
  const selectedIngredientChipItems = computed(() =>
    toKitchenTagChipItems(displaySelectedIngredients.value, 'selected:')
  )
  const excludedIngredientChipItems = computed(() =>
    toKitchenTagChipItems(displayExcludedIngredients.value, 'excluded:')
  )
  /** Removes a custom ingredient from a chip list event payload. */
  const removeCustomIngredientChip = (item: { payload?: unknown }) => {
    if (!item.payload) return
    void removeCustomIngredient(item.payload as KitchenUserIngredient)
  }
  /** Removes a favorite include ingredient from a chip list event payload. */
  const removeFavoriteIncludeChip = (item: { payload?: unknown }) => {
    if (!item.payload) return
    void removeFavoriteIngredientChip(item.payload as KitchenFavoriteIngredient, 'include')
  }
  /** Removes a favorite exclude ingredient from a chip list event payload. */
  const removeFavoriteExcludeChip = (item: { payload?: unknown }) => {
    if (!item.payload) return
    void removeFavoriteIngredientChip(item.payload as KitchenFavoriteIngredient, 'exclude')
  }
  const displaySelectedIngredients = computed(() =>
    searchUiHydrated.value ? resolvedDisplaySelectedIngredients.value : selectedIngredients.value
  )
  const displayExcludedIngredients = computed(() =>
    searchUiHydrated.value ? resolvedDisplayExcludedIngredients.value : excludedIngredients.value
  )
  const catalogCardVisualClass = () =>
    effectiveTheme.value === 'light'
      ? 'bg-white/80 hover:bg-white border-zinc-300'
      : 'bg-zinc-900/70 hover:bg-zinc-900 border-zinc-700'
  const catalogShellClass = computed(() => (effectiveTheme.value === 'light' ? 'bg-zinc-50 p-2' : 'bg-zinc-900 p-2'))
  const catalogGroupDividerClass = computed(() =>
    effectiveTheme.value === 'light'
      ? 'inline-flex w-full items-center gap-2 text-[10px] uppercase tracking-[0.08em] text-zinc-500'
      : 'inline-flex w-full items-center gap-2 text-[10px] uppercase tracking-[0.08em] text-zinc-500'
  )
  const catalogGroupDividerLineClass = computed(() =>
    effectiveTheme.value === 'light' ? 'h-px flex-1 rounded-full bg-zinc-300' : 'h-px flex-1 rounded-full bg-zinc-700'
  )
  const catalogCardMetaClass = computed(() =>
    effectiveTheme.value === 'light'
      ? 'text-[10px] leading-tight text-zinc-600'
      : 'text-[10px] leading-tight text-zinc-500'
  )
  const catalogFavoriteButtonClass = (active: boolean) =>
    active
      ? effectiveTheme.value === 'light'
        ? 'inline-flex h-6 w-6 items-center justify-center rounded-full border border-amber-500/70 bg-amber-100 text-amber-700 opacity-100 pointer-events-auto transition-colors duration-150'
        : 'inline-flex h-6 w-6 items-center justify-center rounded-full border border-amber-300/80 bg-amber-400/20 text-amber-300 opacity-100 pointer-events-auto transition-colors duration-150'
      : effectiveTheme.value === 'light'
        ? 'inline-flex h-6 w-6 items-center justify-center rounded-full border border-zinc-300 bg-white text-zinc-700 opacity-100 pointer-events-auto transition-colors duration-150 hover:border-amber-400 hover:bg-amber-50 hover:text-amber-700 lg:opacity-0 lg:pointer-events-none lg:group-hover:opacity-100 lg:group-hover:pointer-events-auto lg:focus-visible:opacity-100 lg:focus-visible:pointer-events-auto'
        : 'inline-flex h-6 w-6 items-center justify-center rounded-full border border-zinc-700 bg-zinc-950/85 text-zinc-300 opacity-100 pointer-events-auto transition-colors duration-150 hover:border-amber-400/70 hover:bg-amber-400/10 hover:text-amber-300 lg:opacity-0 lg:pointer-events-none lg:group-hover:opacity-100 lg:group-hover:pointer-events-auto lg:focus-visible:opacity-100 lg:focus-visible:pointer-events-auto'
  const catalogExcludeButtonClass = (active: boolean) =>
    active
      ? effectiveTheme.value === 'light'
        ? 'inline-flex h-6 w-6 items-center justify-center rounded-full border border-rose-500/70 bg-rose-100 text-rose-700 opacity-100 pointer-events-auto transition-colors duration-150'
        : 'inline-flex h-6 w-6 items-center justify-center rounded-full border border-rose-300/80 bg-rose-400/20 text-rose-200 opacity-100 pointer-events-auto transition-colors duration-150'
      : effectiveTheme.value === 'light'
        ? 'inline-flex h-6 w-6 items-center justify-center rounded-full border border-zinc-300 bg-white text-zinc-700 opacity-100 pointer-events-auto transition-colors duration-150 hover:border-rose-400 hover:bg-rose-50 hover:text-rose-700 lg:opacity-0 lg:pointer-events-none lg:group-hover:opacity-100 lg:group-hover:pointer-events-auto lg:focus-visible:opacity-100 lg:focus-visible:pointer-events-auto'
        : 'inline-flex h-6 w-6 items-center justify-center rounded-full border border-zinc-700 bg-zinc-950/85 text-zinc-300 opacity-100 pointer-events-auto transition-colors duration-150 hover:border-rose-400/70 hover:bg-rose-400/10 hover:text-rose-200 lg:opacity-0 lg:pointer-events-none lg:group-hover:opacity-100 lg:group-hover:pointer-events-auto lg:focus-visible:opacity-100 lg:focus-visible:pointer-events-auto'
  const recipeDifficultyScaleOptions = computed(() =>
    difficultyOptions.value.map(item => ({
      value: item.code,
      label: item.label
    }))
  )
  const cookingMethodOptions = computed(() => catalogStore.getOptions('cooking_method'))
  const mealTypeOptions = computed(() => catalogStore.getOptions('meal_type'))
  const dietTypeOptions = computed(() => catalogStore.getOptions('diet_type'))
  const mealTypeSelectOptions = computed<SelectOptionInput[]>(() => [
    { value: 'all', label: 'любой' },
    ...mealTypeOptions.value.map(item => ({
      value: item.code,
      label: item.label
    }))
  ])
  const dietTypeSelectOptions = computed<SelectOptionInput[]>(() => [
    { value: 'all', label: 'любой' },
    ...dietTypeOptions.value.map(item => ({
      value: item.code,
      label: item.label
    }))
  ])
  const cuisineOptions = computed(() => catalogStore.getOptions('cuisine'))
  const unitOptions = computed(() => catalogStore.getOptions('unit'))
  const KITCHEN_MEAL_TYPE_SUGGESTIONS = computed(() => mealTypeOptions.value.map(item => item.label))
  const KITCHEN_DIET_TYPE_SUGGESTIONS = computed(() => dietTypeOptions.value.map(item => item.label))
  const KITCHEN_COOKING_METHOD_SUGGESTIONS = computed(() => cookingMethodOptions.value.map(item => item.label))
  const KITCHEN_UNIT_SUGGESTIONS = computed(() => unitOptions.value.map(item => item.label))
  const KITCHEN_NATIONALITY_SUGGESTIONS = computed(() => cuisineOptions.value.map(item => item.label))
  function difficultyLabel(value: string) {
    return catalogStore.labelFor('difficulty', value) || value
  }
  function mealTypeLabel(value: string) {
    const raw = String(value || '').trim()
    if (!raw) return 'другое'
    return catalogStore.labelFor('meal_type', raw) || raw
  }
  function dietTypeLabel(value?: string) {
    const raw = String(value || '').trim()
    if (!raw) return 'без ограничений'
    return catalogStore.labelFor('diet_type', raw) || raw
  }
  function cookingMethodLabel(value?: string) {
    const raw = String(value || '').trim()
    if (!raw) return 'не указан'
    return catalogStore.labelFor('cooking_method', raw) || raw
  }
  const recipeAttentionUpdatedAt = (recipe: KitchenRecipe) =>
    Date.parse(String(recipe.moderated_at || recipe.updated_at || recipe.created_at || ''))
  const userAttentionRecipes = computed(() => {
    const sinceRaw = String(userAttentionSummary.value?.notifications_since_at || '').trim()
    const sinceTs = sinceRaw ? Date.parse(sinceRaw) : Number.NaN
    if (!Number.isFinite(sinceTs)) return [] as KitchenRecipe[]
    return myRecipes.value
      .filter(recipe => {
        if (recipe.moderation_status !== 'approved' && recipe.moderation_status !== 'rejected') return false
        const eventTs = recipeAttentionUpdatedAt(recipe)
        return Number.isFinite(eventTs) && eventTs > sinceTs
      })
      .sort((left, right) => recipeAttentionUpdatedAt(right) - recipeAttentionUpdatedAt(left))
  })
  const userAttentionSummaryText = computed(() => {
    const approved = Number(userAttentionSummary.value?.new_approved_recipes_since_last_login || 0)
    const rejected = Number(userAttentionSummary.value?.new_rejected_recipes_since_last_login || 0)
    const parts: string[] = []
    if (approved > 0) parts.push(`одобрено: ${approved}`)
    if (rejected > 0) parts.push(`отклонено: ${rejected}`)
    return parts.length ? `Новые события по вашим рецептам: ${parts.join(', ')}.` : ''
  })
  function normalizeMealTypeInput(value: string) {
    return catalogStore.normalizeOptionValue('meal_type', value)
  }
  function normalizeDietTypeInput(value: string) {
    return catalogStore.normalizeOptionValue('diet_type', value)
  }
  function normalizeCookingMethodInput(value: string) {
    return catalogStore.normalizeOptionValue('cooking_method', value)
  }
  const recipeCoverSrc = (recipe: KitchenRecipe) => recipeImages.buildStepImageUrl(recipe.cover_image_key)
  const kitchenRecipeLink = (recipe: KitchenRecipe) => `/kitchen/${kitchenRecipeSlug(recipe)}`
  const kitchenRecipeEditLink = (recipe: KitchenRecipe) => `/kitchen/edit/${kitchenRecipeSlug(recipe)}`
  const stepsCountLabel = (count: number) => kitchenStepsCountLabel(count)
  const recipeMetaItems = (recipe: KitchenRecipe, options: { includeSteps?: boolean } = {}) =>
    buildKitchenRecipeMetaItems(
      recipe,
      {
        mealTypeLabel,
        dietTypeLabel,
        cookingMethodLabel,
        difficultyLabel,
        stepsCountLabel
      },
      options
    )
  const scheduleRecipeAutoSearch = debounce(() => {
    if (!import.meta.client) return
    if (activeKitchenTab.value !== 'recipes') return
    if (recipesLoading.value) {
      queuedAutoSearchAfterLoading.value = true
      return
    }
    void runRecipeSearch()
  }, RECIPE_AUTO_SEARCH_DELAY_MS)
  const applyLatestRecipesSnapshot = (payload?: { items?: KitchenRecipe[]; total_public?: number | null } | null) => {
    recipes.value = payload?.items || []
    publicRecipesTotal.value = Number(payload?.total_public ?? 0)
    resultMode.value = 'latest'
  }
  const latestRecipesInitialized = ref(false)
  const { data: initialLatestRecipesData, error: initialLatestRecipesError } = await useAsyncData(
    'kitchen-latest-recipes-initial',
    async () => {
      const res = await getKitchenLatestRecipes(24)
      return {
        items: res?.data?.items || [],
        total_public: Number(res?.data?.total_public ?? 0)
      }
    },
    {
      default: () => ({ items: [], total_public: 0 })
    }
  )
  if (initialLatestRecipesError.value) {
    const initialLatestError = initialLatestRecipesError.value as any
    recipesError.value =
      initialLatestError?.data?.message || initialLatestError?.message || 'Не удалось загрузить последние рецепты.'
    recipes.value = []
    publicRecipesTotal.value = null
    resultMode.value = 'latest'
  } else {
    applyLatestRecipesSnapshot(initialLatestRecipesData.value)
    latestRecipesInitialized.value = true
  }
  const loadLatestRecipes = async () => {
    recipesLoading.value = true
    recipesError.value = null
    resultMode.value = 'latest'
    try {
      const res = await getKitchenLatestRecipes(24)
      applyLatestRecipesSnapshot(res?.data)
      latestRecipesInitialized.value = true
    } catch (err: any) {
      recipesError.value = err?.data?.message || err?.message || 'Не удалось загрузить последние рецепты.'
      recipes.value = []
      publicRecipesTotal.value = null
      latestRecipesInitialized.value = false
    } finally {
      recipesLoading.value = false
    }
  }
  const runRecipeSearch = async () => {
    recipeServingsRangeError.value = ''
    recipeKcalRangeError.value = ''
    if (isServingsRangeInvalid.value) {
      recipeServingsRangeError.value = 'Значение «от» не может быть больше значения «до».'
      return
    }
    if (isKcalRangeInvalid.value) {
      recipeKcalRangeError.value = 'Значение «от» не может быть больше значения «до».'
      return
    }
    const query = recipeQuery.value.trim()
    const resolvedFilters = resolvedSearchFilters.value
    const ingredients = resolvedFilters.selected
    const excluded = resolvedFilters.excluded
    const parsedMaxTotalMinutes = maxTotalMinutesFilter.value
    const { min: servingsMin, max: servingsMax } = servingsRange.value
    const { min: kcalMin, max: kcalMax } = kcalRange.value
    const hasSearchScope = Boolean(query) || ingredients.length > 0 || excluded.length > 0 || hasMetaRecipeFilters.value
    if (!hasSearchScope) {
      if (latestRecipesInitialized.value && resultMode.value === 'latest' && !recipesError.value) {
        return
      }
      await loadLatestRecipes()
      return
    }
    recipesLoading.value = true
    recipesError.value = null
    resultMode.value = 'search'
    try {
      const payload: KitchenSearchInput = {
        query,
        ingredients,
        excludedIngredients: excluded,
        requireAllIngredients: requireAllIngredients.value,
        limit: 50,
        offset: 0
      }
      if (recipeDifficultyFilter.value !== 'all') payload.difficulty = recipeDifficultyFilter.value
      if (recipeMealTypeFilter.value !== 'all') payload.mealType = recipeMealTypeFilter.value
      if (recipeDietTypeFilter.value !== 'all') payload.dietType = recipeDietTypeFilter.value
      if (recipeCookingMethodFilter.value.trim()) {
        payload.cookingMethod = normalizeCookingMethodInput(recipeCookingMethodFilter.value)
      }
      if (recipeCuisineFilter.value.trim()) payload.cuisine = recipeCuisineFilter.value.trim()
      payload.kcalMin = kcalMin
      payload.kcalMax = kcalMax
      payload.maxTotalMinutes = parsedMaxTotalMinutes
      payload.servingsMin = servingsMin
      payload.servingsMax = servingsMax
      const res = await searchKitchenRecipes(payload)
      recipes.value = res?.data?.items || []
    } catch (err: any) {
      recipesError.value = err?.data?.message || err?.message || 'Не удалось выполнить поиск рецептов.'
      recipes.value = []
    } finally {
      recipesLoading.value = false
    }
  }
  const resetRecipeSearch = async () => {
    resetRecipeSearchFilters()
    await loadLatestRecipes()
  }
  const runRecipeSearchFromToolbar = async () => {
    if (recipesLoading.value || recipeToolbarActionLoading.value) return
    recipeToolbarActionLoading.value = 'search'
    try {
      await runRecipeSearch()
    } finally {
      if (recipeToolbarActionLoading.value === 'search') {
        recipeToolbarActionLoading.value = null
      }
    }
  }
  const resetRecipeSearchFromToolbar = async () => {
    if (recipesLoading.value || recipeToolbarActionLoading.value) return
    recipeToolbarActionLoading.value = 'reset'
    try {
      await resetRecipeSearch()
    } finally {
      if (recipeToolbarActionLoading.value === 'reset') {
        recipeToolbarActionLoading.value = null
      }
    }
  }
  watch(
    () => [
      recipeQuery.value,
      selectedIngredients.value.join('|'),
      excludedIngredients.value.join('|'),
      requireAllIngredients.value,
      includeFavoriteIngredientsInSearch.value,
      excludeFavoriteIngredientsInSearch.value,
      recipeDifficultyFilter.value,
      recipeMealTypeFilter.value,
      recipeDietTypeFilter.value,
      recipeCookingMethodFilter.value,
      recipeCuisineFilter.value,
      recipeKcalMin.value,
      recipeKcalMax.value,
      recipeMaxTotalMinutes.value,
      recipeServingsMin.value,
      recipeServingsMax.value
    ],
    () => {
      if (!isServingsRangeInvalid.value) recipeServingsRangeError.value = ''
      if (!isKcalRangeInvalid.value) recipeKcalRangeError.value = ''
      scheduleRecipeAutoSearch()
    }
  )
  watch(recipesLoading, loading => {
    if (loading) return
    if (!queuedAutoSearchAfterLoading.value) return
    queuedAutoSearchAfterLoading.value = false
    scheduleRecipeAutoSearch()
  })
  watch(activeKitchenTab, tab => {
    if (tab === 'recipes') return
    queuedAutoSearchAfterLoading.value = false
  })
  watch(showMyRecipesList, async next => {
    if (!next || !isAuthenticated.value) return
    await Promise.all([loadMyRecipes(), loadUserAttentionSummary()])
  })
  watch(
    () => [isAuthenticated.value, showMyRecipesList.value, visibleFavoriteStatusRecipeIds.value.join('|')] as const,
    () => {
      void syncVisibleRecipeFavoriteStatuses()
    },
    { immediate: true }
  )
  const openNewRecipeForm = async () => {
    createError.value = null
    createSuccess.value = null
    resetRecipeForm({ skipEditRouteRedirect: true })
    await openCreateRecipeTab()
  }
  const cancelRecipeEditing = async () => {
    createError.value = null
    createSuccess.value = null
    resetRecipeForm({ skipEditRouteRedirect: true })
    await openMyRecipesTab()
  }
  async function syncEditRouteRecipe() {
    if (!isKitchenEditRoute.value) return
    activeKitchenTab.value = 'edit'
    const id = routeEditRecipeId.value
    if (!id) {
      createError.value = 'Не удалось определить рецепт для редактирования.'
      return
    }
    if (editRecipeId.value === id) return
    if (loadingRouteRecipeId.value === id) return
    if (!authLoaded.value) {
      await ensureLoaded()
    }
    if (!isAuthenticated.value) {
      createError.value = 'Для редактирования рецепта нужен активный аккаунт.'
      return
    }
    const prefetchedRecipe = routeEditRecipeData.value
    if (prefetchedRecipe && prefetchedRecipe.id === id) {
      if (!canManageRecipe(prefetchedRecipe)) {
        createError.value = 'У вас нет доступа к редактированию этого рецепта.'
        return
      }
      applyRecipeToForm(prefetchedRecipe, {
        mealTypeLabel,
        cookingMethodLabel,
        dietTypeLabel
      })
      return
    }
    if (routeEditRecipeError.value) {
      createError.value = 'Не удалось загрузить рецепт для редактирования.'
      return
    }
    const inMemoryRecipe =
      myRecipes.value.find(item => item.id === id) || recipes.value.find(item => item.id === id) || null
    if (inMemoryRecipe) {
      if (!canManageRecipe(inMemoryRecipe)) {
        createError.value = 'У вас нет доступа к редактированию этого рецепта.'
        return
      }
      applyRecipeToForm(inMemoryRecipe, {
        mealTypeLabel,
        cookingMethodLabel,
        dietTypeLabel
      })
      return
    }
    loadingRouteRecipeId.value = id
    createError.value = null
    try {
      const res = await getKitchenManageRecipeById(id)
      const recipe = res?.data || null
      if (!recipe) {
        throw new Error('Рецепт не найден.')
      }
      if (!canManageRecipe(recipe)) {
        createError.value = 'У вас нет доступа к редактированию этого рецепта.'
        return
      }
      applyRecipeToForm(recipe, {
        mealTypeLabel,
        cookingMethodLabel,
        dietTypeLabel
      })
    } catch (err: any) {
      createError.value = err?.data?.message || err?.message || 'Не удалось загрузить рецепт для редактирования.'
    } finally {
      loadingRouteRecipeId.value = null
    }
  }
  const findExistingIngredient = (name: string) => {
    const normalized = normalizeTag(name)
    if (!normalized) return null
    const customMatch = customIngredients.value.find(item => normalizeTag(item.name) === normalized)
    if (customMatch) {
      return { name: customMatch.name, category: customMatch.category }
    }
    const catalogMatch = catalogStore.ingredientItems.find(item => normalizeTag(item.name) === normalized)
    if (catalogMatch) {
      return { name: catalogMatch.name, category: catalogMatch.category }
    }
    return null
  }
  const deleteRecipe = async (recipe: KitchenRecipe) => {
    if (!canManageRecipe(recipe) || deletePendingRecipeId.value) return
    deletePendingRecipeId.value = recipe.id
    createError.value = null
    createSuccess.value = null
    try {
      await deleteKitchenRecipe(recipe.id)
      if (editRecipeId.value === recipe.id) {
        resetRecipeForm()
      }
      if (isKitchenEditRoute.value && routeEditRecipeId.value === recipe.id) {
        await goToKitchenTab('my-recipes')
      }
      createSuccess.value = 'Рецепт удалён.'
      await Promise.all([loadLatestRecipes(), loadMyRecipes(), loadUserAttentionSummary()])
    } catch (err: any) {
      createError.value = err?.data?.message || err?.message || 'Не удалось удалить рецепт.'
    } finally {
      deletePendingRecipeId.value = null
    }
  }
  const addCustomIngredient = async () => {
    if (!isAuthenticated.value) {
      accountIngredientsError.value = 'Сначала войдите в аккаунт.'
      return
    }
    const name = customIngredientForm.name.trim()
    if (!name) {
      accountIngredientsError.value = 'Укажите название ингредиента.'
      return
    }
    const existingIngredient = findExistingIngredient(name)
    if (existingIngredient) {
      accountIngredientsError.value = `Такой ингредиент уже есть: «${existingIngredient.name}», категория «${existingIngredient.category}».`
      return
    }
    customIngredientPending.value = true
    accountIngredientsError.value = null
    try {
      await createKitchenCustomIngredient({
        name,
        category: customIngredientForm.category.trim() || 'другое'
      })
      customIngredientForm.name = ''
      customIngredientForm.category = catalogStore.categoryLabels.includes('другое')
        ? 'другое'
        : catalogStore.categoryLabels[0] || 'другое'
      await loadKitchenAccountIngredients()
    } catch (err: any) {
      accountIngredientsError.value = err?.data?.message || err?.message || 'Не удалось сохранить ингредиент.'
    } finally {
      customIngredientPending.value = false
    }
  }
  const removeCustomIngredient = async (item: KitchenUserIngredient) => {
    if (!isAuthenticated.value || deletingCustomIngredientId.value) return
    deletingCustomIngredientId.value = item.ingredient_id
    accountIngredientsError.value = null
    try {
      await deleteKitchenCustomIngredient(item.ingredient_id)
      await loadKitchenAccountIngredients()
    } catch (err: any) {
      accountIngredientsError.value =
        err?.data?.message || err?.message || 'Не удалось удалить персональный ингредиент.'
    } finally {
      deletingCustomIngredientId.value = null
    }
  }
  const toggleFavoriteIngredient = async (
    item: { ingredient_id?: number; name: string; category: string },
    listType: 'include' | 'exclude' = 'include'
  ) => {
    if (!isAuthenticated.value) {
      accountIngredientsError.value = 'Сначала войдите в аккаунт.'
      return
    }
    if (typeof item.ingredient_id !== 'number' || item.ingredient_id <= 0) {
      accountIngredientsError.value = 'Для этого ингредиента недоступно действие избранного.'
      return
    }
    accountIngredientsError.value = null
    try {
      const exists =
        listType === 'exclude'
          ? isExcludeFavoriteIngredient(item.ingredient_id)
          : isFavoriteIngredient(item.ingredient_id)
      const oppositeListType: 'include' | 'exclude' = listType === 'exclude' ? 'include' : 'exclude'
      const existsInOppositeList =
        oppositeListType === 'exclude'
          ? isExcludeFavoriteIngredient(item.ingredient_id)
          : isFavoriteIngredient(item.ingredient_id)
      if (exists) {
        await unfavoriteKitchenIngredient(item.ingredient_id, listType)
        removeFavoriteIngredientById(item.ingredient_id, listType)
      } else {
        if (existsInOppositeList) {
          await unfavoriteKitchenIngredient(item.ingredient_id, oppositeListType)
          removeFavoriteIngredientById(item.ingredient_id, oppositeListType)
        }
        const res = await favoriteKitchenIngredient({ ingredient_id: item.ingredient_id, list_type: listType })
        if (res?.data) {
          upsertFavoriteIngredient(normalizeKitchenFavoriteIngredient(res.data, listType))
        }
      }
    } catch (err: any) {
      accountIngredientsError.value = err?.data?.message || err?.message || 'Не удалось обновить избранное.'
    }
  }
  const removeFavoriteIngredientChip = async (item: KitchenFavoriteIngredient, listType: 'include' | 'exclude') => {
    if (!isAuthenticated.value || deletingFavoriteKey.value) return
    deletingFavoriteKey.value = `${listType}:${item.ingredient_id}`
    accountIngredientsError.value = null
    try {
      await unfavoriteKitchenIngredient(item.ingredient_id, listType)
      removeFavoriteIngredientById(item.ingredient_id, listType)
    } catch (err: any) {
      accountIngredientsError.value =
        err?.data?.message || err?.message || 'Не удалось удалить ингредиент из избранного.'
    } finally {
      deletingFavoriteKey.value = null
    }
  }
  onMounted(async () => {
    searchUiHydrated.value = true
    await ensureLoaded()
    const startupTasks: Promise<unknown>[] = [loadKitchenAccountIngredients()]
    if (showMyRecipesList.value && isAuthenticated.value) {
      startupTasks.push(loadMyRecipes(), loadUserAttentionSummary())
    } else {
      myRecipes.value = []
      myRecipesError.value = null
    }
    await Promise.all(startupTasks)
    accountStateInitialized.value = true
  })
  watch(isAuthenticated, async () => {
    if (!accountStateInitialized.value) return
    if (!isAuthenticated.value) {
      showFavoriteRecipesOnly.value = false
      myRecipes.value = []
      myRecipesError.value = null
    }
    const refreshTasks: Promise<unknown>[] = [loadKitchenAccountIngredients()]
    if (isAuthenticated.value && showMyRecipesList.value) {
      refreshTasks.push(loadMyRecipes(), loadUserAttentionSummary())
    }
    await Promise.all(refreshTasks)
  })
  onBeforeUnmount(() => {
    queuedAutoSearchAfterLoading.value = false
  })
</script>
<template>
  <LabBaseSection variant="plain">
    <LabNavHeader :title :breadcrumb-items="kitchenBreadcrumbItems" />
    <LabNavTabs
      v-model="activeKitchenTab"
      :items="kitchenTabItems"
      :no-select="true"
      :render-panels="false"
      route-query-key="tab"
      route-default-value="recipes"
      route-path="/kitchen"
      :route-active-value="isKitchenEditRoute ? 'edit' : null"
      :route-to-map="kitchenTabRouteMap"
      container-class="w-full" />
    <section v-show="activeKitchenTab === 'ingredients'" class="p-3 space-y-4 sm:p-5">
      <LabNotify
        :text="catalogStore.error ? `Справочник кухни загружен не полностью: ${catalogStore.error}` : ''"
        tone="warning"
        size="xs" />
      <div class="space-y-2 bg-zinc-900 p-2">
        <div class="flex flex-wrap items-end gap-2">
          <LabField
            label="Название ингредиента"
            for-id="ingredient-catalog-query"
            field-class="min-w-56 max-w-120 flex-1">
            <LabBaseInput
              id="ingredient-catalog-query"
              v-model="ingredientCatalogQuery"
              name="ingredient_catalog_query"
              type="text"
              placeholder="например, огурец" />
          </LabField>
          <LabField label="Категория" for-id="ingredient-category-filter">
            <LabBaseSelect
              id="ingredient-category-filter"
              v-model="categoryFilter"
              name="ingredient_category_filter"
              :options="ingredientCategoryOptions" />
          </LabField>
          <div class="flex flex-col gap-1">
            <span class="text-[11px] text-zinc-500">Сортировка: {{ catalogSortLabel }}</span>
            <LabNavTabs
              v-model="catalogSortMode"
              :items="catalogSortItems"
              :no-select="true"
              :render-panels="false"
              route-query-key="sort"
              route-default-value="alpha"
              container-class="space-y-0"
              list-class="flex flex-wrap gap-2" />
          </div>
          <KitchenGroupByCategoryToggle v-model="groupCatalogByCategory" mode="button" />
        </div>
        <ClientOnly>
          <span v-if="catalogFrequencyHint" class="text-[11px] text-zinc-500">
            {{ catalogFrequencyHint }}
          </span>
        </ClientOnly>
      </div>
      <div :class="catalogShellClass">
        <p class="text-xs leading-none text-zinc-500">Найдено: {{ filteredCatalog.length }}</p>
        <div class="max-h-96 overflow-y-auto pt-2 pr-1">
          <div class="space-y-3">
            <div v-for="group in visibleCatalogGroups" :key="group.category || 'all'" class="space-y-2">
              <div v-if="groupCatalogByCategory" :class="catalogGroupDividerClass">
                <span :class="catalogGroupDividerLineClass"></span>
                <span class="whitespace-nowrap">{{ group.category }}</span>
                <span :class="catalogGroupDividerLineClass"></span>
              </div>
              <div class="grid grid-cols-2 gap-0.5 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
                <div
                  v-for="item in group.items"
                  :key="item.id"
                  class="group relative flex aspect-4/3 w-full overflow-hidden border p-0 text-left transition-colors duration-200 hover:ring-1 hover:ring-zinc-500/70"
                  :class="catalogCardVisualClass()">
                  <div class="absolute inset-0 flex items-center justify-center bg-zinc-950/70">
                    <div class="flex flex-col items-center gap-1 text-zinc-500">
                      <Icon name="ic:round-image" class="h-5 w-5" />
                      <span class="text-[10px] leading-none">изображение позже</span>
                    </div>
                  </div>
                  <div class="absolute inset-x-0 top-0 z-10 flex items-start justify-between gap-2 p-2">
                    <h3
                      class="max-w-[calc(100%-4.5rem)] bg-zinc-950/70 px-2 py-1 text-[11px] font-medium leading-tight text-zinc-100 wrap-break-word">
                      {{ item.name }}
                    </h3>
                    <div
                      v-if="showCatalogFavoriteActions"
                      class="pointer-events-auto inline-flex shrink-0 items-center gap-1">
                      <LabBaseButton
                        :button-class="catalogFavoriteButtonClass(isFavoriteIngredient(item.ingredient_id))"
                        :aria-label="
                          isFavoriteIngredient(item.ingredient_id)
                            ? `Убрать ${item.name} из любимых`
                            : `Добавить ${item.name} в любимые`
                        "
                        :icon="isFavoriteIngredient(item.ingredient_id) ? 'ic:round-star' : 'ic:round-star-border'"
                        :icon-class="isFavoriteIngredient(item.ingredient_id) ? 'text-amber-300' : ''"
                        size="xs"
                        icon-only
                        @click.stop="toggleFavoriteIngredient(item, 'include')" />
                      <LabBaseButton
                        :button-class="catalogExcludeButtonClass(isExcludeFavoriteIngredient(item.ingredient_id))"
                        :aria-label="
                          isExcludeFavoriteIngredient(item.ingredient_id)
                            ? `Убрать ${item.name} из исключаемых`
                            : `Добавить ${item.name} в исключаемые`
                        "
                        :icon="
                          isExcludeFavoriteIngredient(item.ingredient_id)
                            ? 'ic:round-block'
                            : 'ic:round-do-not-disturb-on'
                        "
                        icon-only
                        size="xs"
                        @click.stop="toggleFavoriteIngredient(item, 'exclude')" />
                    </div>
                  </div>
                  <div
                    class="pointer-events-none absolute inset-x-0 bottom-0 hidden translate-y-2 bg-zinc-950/78 p-2 opacity-0 transition duration-200 group-hover:translate-y-0 group-hover:opacity-100 md:block">
                    <div class="flex flex-col gap-1 text-left">
                      <p v-if="!groupCatalogByCategory" class="text-[10px] leading-tight text-zinc-400 wrap-break-word">
                        {{ item.category }}
                      </p>
                      <p v-if="item.description" :class="['line-clamp-2', catalogCardMetaClass]">
                        {{ item.description }}
                      </p>
                      <p
                        v-if="
                          item.kcal !== null || item.protein_g !== null || item.fat_g !== null || item.carbs_g !== null
                        "
                        :class="catalogCardMetaClass">
                        {{ item.kcal ?? '—' }} ккал · Б {{ item.protein_g ?? '—' }} · Ж {{ item.fat_g ?? '—' }} · У
                        {{ item.carbs_g ?? '—' }}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <p v-if="!authLoaded" class="text-xs text-zinc-500">Проверяем сессию...</p>
      <div v-else-if="isAuthenticated" class="space-y-3 bg-zinc-900 p-3">
        <div class="flex flex-wrap items-center justify-between gap-2">
          <div>
            <p class="text-sm text-zinc-200">Персональные ингредиенты и избранное</p>
          </div>
          <LabLoader
            v-if="accountIngredientsLoading"
            variant="inline"
            label="Загрузка персональных ингредиентов"
            class="lab-text-muted" />
        </div>
        <div class="grid gap-3 md:grid-cols-[minmax(0,1.2fr)_minmax(0,1fr)]">
          <div class="space-y-2">
            <p class="text-xs uppercase tracking-[0.06em] text-zinc-500">Добавить свой ингредиент</p>
            <div class="flex flex-wrap gap-2">
              <LabBaseInput
                id="custom-ingredient-name"
                v-model="customIngredientForm.name"
                name="custom_ingredient_name"
                type="text"
                placeholder="например, соус кимчи"
                class="min-w-56 flex-1" />
              <LabBaseSelect
                id="custom-ingredient-category"
                v-model="customIngredientForm.category"
                name="custom_ingredient_category"
                :options="customIngredientCategoryOptions" />
              <LabBaseButton
                variant="primary"
                size="md"
                :disabled="customIngredientPending"
                @click="addCustomIngredient">
                Добавить
              </LabBaseButton>
            </div>
            <div class="flex flex-wrap gap-2">
              <KitchenRemovableChipList
                :items="customIngredientChipItems"
                empty-text="Своих ингредиентов пока нет."
                :item-style="item => selectedTagStyle(item.label)"
                :item-title="item => `Удалить персональный ингредиент «${item.label}»`"
                :item-aria-label="item => `Удалить персональный ингредиент ${item.label}`"
                @remove="removeCustomIngredientChip" />
            </div>
          </div>
          <div class="space-y-3">
            <div class="space-y-2">
              <p class="text-xs uppercase tracking-[0.06em] text-zinc-500">Любимые ингредиенты (для включения)</p>
              <KitchenRemovableChipList
                :items="favoriteIncludeChipItems"
                empty-text="Любимые пока пусто."
                :item-style="item => selectedTagStyle(item.label)"
                :item-title="item => `Убрать «${item.label}» из любимых`"
                :item-aria-label="item => `Убрать ${item.label} из любимых`"
                @remove="removeFavoriteIncludeChip" />
            </div>
            <div class="space-y-2">
              <p class="text-xs uppercase tracking-[0.06em] text-zinc-500">Исключаемые ингредиенты</p>
              <KitchenRemovableChipList
                :items="favoriteExcludeChipItems"
                empty-text="Список исключений пока пуст."
                :item-style="item => selectedTagStyle(item.label)"
                :item-title="item => `Убрать «${item.label}» из исключаемых`"
                :item-aria-label="item => `Убрать ${item.label} из исключаемых`"
                @remove="removeFavoriteExcludeChip" />
            </div>
          </div>
        </div>
        <LabNotify :text="accountIngredientsError" tone="error" size="xs" />
      </div>
      <p v-else class="text-xs text-zinc-500">Войдите в аккаунт, чтобы добавлять свои ингредиенты и вести избранное.</p>
    </section>
    <section v-show="activeKitchenTab === 'recipes'" class="p-3 space-y-4 sm:p-5">
      <div
        class="sticky top-0 z-20 -mx-3 border-y border-zinc-800 px-3 py-2 backdrop-blur sm:-mx-5 sm:px-5 md:static md:mx-0 md:border-y-0 md:bg-transparent md:px-0 md:py-0 md:backdrop-blur-none">
        <div class="flex items-end gap-2">
          <LabField label="Поиск по названию рецепта" for-id="recipe-query" field-class="min-w-0 max-w-[34rem] flex-1">
            <LabBaseInput
              id="recipe-query"
              v-model="recipeQuery"
              name="recipe_query"
              type="text"
              placeholder="например, салат" />
          </LabField>
          <div class="flex shrink-0 items-center gap-2">
            <LabBaseButton
              :button-class="[
                'inline-flex h-11 w-11 items-center justify-center border transition',
                showFavoriteRecipesOnly
                  ? 'border-rose-400/90 bg-rose-600 text-white hover:bg-rose-500'
                  : 'border-zinc-700 bg-zinc-900 text-rose-300 hover:border-rose-500/70 hover:bg-rose-500/10'
              ]"
              :title="showFavoriteRecipesOnly ? 'Показать все рецепты' : 'Показать только избранные рецепты'"
              :aria-label="showFavoriteRecipesOnly ? 'Показать все рецепты' : 'Показать только избранные рецепты'"
              :aria-pressed="showFavoriteRecipesOnly ? 'true' : 'false'"
              :disabled="recipeToolbarActionLoading === 'favorite'"
              :icon="
                recipeToolbarActionLoading === 'favorite'
                  ? ''
                  : showFavoriteRecipesOnly
                    ? 'ic:round-favorite'
                    : 'ic:round-favorite-border'
              "
              icon-class="h-5 w-5 text-xl"
              icon-only
              @click="toggleFavoriteRecipesOnlyFilterFromToolbar">
              <LabLoader
                v-if="recipeToolbarActionLoading === 'favorite'"
                variant="icon"
                label="Обновление избранного" />
            </LabBaseButton>
            <LabBaseButton
              :button-class="[
                'inline-flex h-11 w-11 items-center justify-center border transition disabled:cursor-not-allowed disabled:opacity-60',
                recipesLoading
                  ? 'border-emerald-500/50 bg-emerald-500/10 text-emerald-300'
                  : 'border-zinc-700 bg-zinc-900 text-emerald-300 hover:border-emerald-400/70 hover:bg-emerald-500/10'
              ]"
              title="Найти рецепты"
              aria-label="Найти рецепты"
              :disabled="recipesLoading"
              icon-only
              @click="runRecipeSearchFromToolbar">
              <span class="inline-flex h-6 w-6 items-center justify-center">
                <LabLoader
                  v-if="recipeToolbarActionLoading === 'search'"
                  variant="icon"
                  icon-name="ic:round-autorenew"
                  label="Выполняется поиск" />
                <Icon v-else name="ic:round-search" class="h-6 w-6 text-xl" />
              </span>
            </LabBaseButton>
            <LabBaseButton
              :button-class="[
                'inline-flex h-11 w-11 items-center justify-center border transition disabled:cursor-not-allowed disabled:opacity-60',
                recipesLoading
                  ? 'border-zinc-600 bg-zinc-900 text-zinc-400'
                  : 'border-zinc-700 bg-zinc-900 text-zinc-200 hover:border-zinc-500 hover:bg-zinc-800'
              ]"
              title="Сбросить фильтры"
              aria-label="Сбросить фильтры"
              :disabled="recipesLoading"
              icon-only
              @click="resetRecipeSearchFromToolbar">
              <span class="inline-flex h-6 w-6 items-center justify-center">
                <LabLoader
                  v-if="recipeToolbarActionLoading === 'reset'"
                  variant="icon"
                  icon-name="ic:round-autorenew"
                  label="Сброс фильтров" />
                <Icon v-else name="ic:round-restart-alt" class="h-6 w-6 text-xl" />
              </span>
            </LabBaseButton>
            <LabBaseButton
              :button-class="[
                'inline-flex h-11 w-11 items-center justify-center border transition',
                recipeAdvancedSearchOpen
                  ? 'border-cyan-400/80 bg-cyan-500/20 text-cyan-100 hover:bg-cyan-500/30'
                  : 'border-zinc-700 bg-zinc-900 text-cyan-300 hover:border-cyan-400/70 hover:bg-cyan-500/10'
              ]"
              :title="recipeAdvancedSearchOpen ? 'Скрыть расширенный поиск' : 'Показать расширенный поиск'"
              :aria-label="recipeAdvancedSearchOpen ? 'Скрыть расширенный поиск' : 'Показать расширенный поиск'"
              :aria-pressed="recipeAdvancedSearchOpen ? 'true' : 'false'"
              icon="ic:round-manage-search"
              icon-class="h-5 w-5 text-xl"
              icon-only
              @click="recipeAdvancedSearchOpen = !recipeAdvancedSearchOpen" />
          </div>
        </div>
        <Teleport to="body" :disabled="!isCompactRecipeSearch">
          <LabSpoiler
            v-model="recipeAdvancedSearchOpen"
            title="Расширенный поиск"
            show-label="Показать"
            hide-label="Свернуть"
            :inline-toggle="true"
            :container-class="recipeAdvancedSearchContainerClass"
            header-class="hidden"
            title-class="text-xs font-semibold uppercase tracking-[0.08em] text-zinc-300"
            toggle-button-class="inline-flex h-8 items-center gap-1.5 rounded-md border border-zinc-600/90 bg-zinc-800/70 px-2.5 text-xs text-zinc-100 transition hover:border-zinc-400 hover:bg-zinc-700"
            :content-class="recipeAdvancedSearchContentClass">
            <div class="flex items-start justify-between gap-3 border-b border-zinc-800 pb-3 md:hidden">
              <div class="min-w-0 space-y-1">
                <p class="text-base font-semibold text-zinc-100">Расширенный поиск</p>
                <p class="text-xs text-zinc-500">{{ recipeAdvancedSearchFoundLabel }}</p>
              </div>
              <LabBaseButton variant="secondary" size="md" icon="ic:round-close" @click="closeRecipeAdvancedSearch">
                Закрыть
              </LabBaseButton>
            </div>
            <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
              <LabField label="Тип блюда" for-id="recipe-meal-type-filter" field-class="min-w-0">
                <LabBaseSelect
                  id="recipe-meal-type-filter"
                  v-model="recipeMealTypeFilter"
                  name="recipe_meal_type_filter"
                  :options="mealTypeSelectOptions" />
              </LabField>
              <LabField label="Тип питания" for-id="recipe-diet-type-filter" field-class="min-w-0">
                <LabBaseSelect
                  id="recipe-diet-type-filter"
                  v-model="recipeDietTypeFilter"
                  name="recipe_diet_type_filter"
                  :options="dietTypeSelectOptions" />
              </LabField>
              <LabField label="Способ приготовления" for-id="recipe-cooking-method-filter" field-class="min-w-0">
                <LabBaseInput
                  id="recipe-cooking-method-filter"
                  v-model="recipeCookingMethodFilter"
                  name="recipe_cooking_method_filter"
                  type="text"
                  list="kitchen-cooking-method-list"
                  placeholder="например, на гриле" />
                <datalist id="kitchen-cooking-method-list">
                  <option
                    v-for="item in KITCHEN_COOKING_METHOD_SUGGESTIONS"
                    :key="`cooking-method:${item}`"
                    :value="item"></option>
                </datalist>
              </LabField>
              <LabField label="Время приготовления" for-id="recipe-max-total-minutes" field-class="min-w-0">
                <div class="relative">
                  <span
                    class="pointer-events-none absolute left-2 top-1/2 -translate-y-1/2 text-[11px] font-medium text-zinc-500">
                    до
                  </span>
                  <LabBaseInput
                    id="recipe-max-total-minutes"
                    v-model="recipeMaxTotalMinutes"
                    name="recipe_max_total_minutes"
                    type="number"
                    min="1"
                    :max="RECIPE_TIME_FILTER_MAX"
                    step="1"
                    placeholder="∞"
                    class="w-full"
                    input-class="pl-7 pr-12" />
                  <span
                    class="pointer-events-none absolute right-2 top-1/2 -translate-y-1/2 text-[11px] font-medium text-zinc-500">
                    мин.
                  </span>
                </div>
              </LabField>
              <LabField label="Порции" field-class="min-w-0">
                <div class="flex flex-wrap items-center gap-2">
                  <div class="relative min-w-0 flex-1">
                    <span
                      class="pointer-events-none absolute left-2 top-1/2 -translate-y-1/2 text-[11px] font-medium text-zinc-500">
                      от
                    </span>
                    <LabBaseInput
                      id="recipe-servings-min"
                      v-model="recipeServingsMin"
                      name="recipe_servings_min"
                      type="number"
                      min="1"
                      inputmode="numeric"
                      placeholder="1"
                      class="w-full min-w-0"
                      :input-class="['pl-7 pr-2', isServingsRangeInvalid ? 'lab-control-invalid' : '']" />
                  </div>
                  <div class="relative min-w-0 flex-1">
                    <span
                      class="pointer-events-none absolute left-2 top-1/2 -translate-y-1/2 text-[11px] font-medium text-zinc-500">
                      до
                    </span>
                    <LabBaseInput
                      id="recipe-servings-max"
                      v-model="recipeServingsMax"
                      name="recipe_servings_max"
                      type="number"
                      min="1"
                      inputmode="numeric"
                      placeholder="8"
                      class="w-full min-w-0"
                      :input-class="['pl-7 pr-2', isServingsRangeInvalid ? 'lab-control-invalid' : '']" />
                  </div>
                </div>
                <LabNotify :text="recipeServingsRangeError" tone="error" size="xs" class="mt-2" />
              </LabField>
              <LabField label="Калорийность (ккал)" field-class="min-w-0">
                <div class="flex flex-wrap items-center gap-2">
                  <div class="relative min-w-0 flex-1">
                    <span
                      class="pointer-events-none absolute left-2 top-1/2 -translate-y-1/2 text-[11px] font-medium text-zinc-500">
                      от
                    </span>
                    <LabBaseInput
                      id="recipe-kcal-min"
                      v-model="recipeKcalMin"
                      name="recipe_kcal_min"
                      type="number"
                      min="1"
                      inputmode="numeric"
                      placeholder="200"
                      class="w-full min-w-0"
                      :input-class="['pl-7 pr-2', isKcalRangeInvalid ? 'lab-control-invalid' : '']" />
                  </div>
                  <div class="relative min-w-0 flex-1">
                    <span
                      class="pointer-events-none absolute left-2 top-1/2 -translate-y-1/2 text-[11px] font-medium text-zinc-500">
                      до
                    </span>
                    <LabBaseInput
                      id="recipe-kcal-max"
                      v-model="recipeKcalMax"
                      name="recipe_kcal_max"
                      type="number"
                      min="1"
                      inputmode="numeric"
                      placeholder="700"
                      class="w-full min-w-0"
                      :input-class="['pl-7 pr-2', isKcalRangeInvalid ? 'lab-control-invalid' : '']" />
                  </div>
                </div>
                <LabNotify :text="recipeKcalRangeError" tone="error" size="xs" class="mt-2" />
              </LabField>
              <LabField label="Национальная кухня" for-id="recipe-cuisine-filter" field-class="min-w-0">
                <LabBaseInput
                  id="recipe-cuisine-filter"
                  v-model="recipeCuisineFilter"
                  name="recipe_cuisine_filter"
                  type="text"
                  list="kitchen-search-nationality-list"
                  placeholder="например, итальянская" />
                <datalist id="kitchen-search-nationality-list">
                  <option
                    v-for="item in KITCHEN_NATIONALITY_SUGGESTIONS"
                    :key="`search-nationality:${item}`"
                    :value="item"></option>
                </datalist>
              </LabField>
              <LabBaseScale
                v-model="recipeDifficultyFilter"
                id="recipe-difficulty-filter"
                name="recipe_difficulty_filter"
                label="Сложность"
                :options="recipeDifficultyScaleOptions"
                field-class="min-w-0" />
            </div>
            <div class="space-y-2">
              <LabField label="Ингредиенты для фильтра" for-id="search-ingredient-input">
                <div class="flex flex-col gap-2 sm:flex-row sm:items-end">
                  <LabBaseInput
                    id="search-ingredient-input"
                    v-model="searchIngredientInput"
                    name="search_ingredient_input"
                    type="text"
                    placeholder="например, томат"
                    class="min-w-0 flex-1"
                    @keydown.enter.prevent="addSearchIngredientFromInput" />
                  <div class="flex shrink-0 items-center gap-2">
                    <LabBaseButton variant="primary" size="xl" @click="addSearchIngredientFromInput">
                      Включить
                    </LabBaseButton>
                    <LabBaseButton variant="danger" size="xl" @click="addExcludedIngredientFromInput">
                      Исключить
                    </LabBaseButton>
                  </div>
                </div>
              </LabField>
              <div v-if="ingredientSuggestions.length" class="rounded-md border border-zinc-700 bg-zinc-900">
                <p class="px-2 py-1 text-[11px] text-zinc-500">Найдено: {{ ingredientSuggestions.length }}</p>
                <div class="max-h-28 overflow-y-auto px-2 pb-2">
                  <div class="flex flex-wrap gap-2">
                    <LabBaseButton
                      v-for="item in ingredientSuggestions"
                      :key="item.name"
                      rounded-xl
                      button-class="transition hover:brightness-110"
                      :style="categoryTagStyle(item.category)"
                      size="xs"
                      @click="searchIngredientInput = item.name">
                      {{ item.name }}
                    </LabBaseButton>
                  </div>
                </div>
              </div>
              <div class="grid gap-2 bg-zinc-900 md:grid-cols-2 md:gap-3">
                <KitchenIngredientScopePanel
                  title="Включить"
                  container-class="min-h-16 overflow-visible md:max-h-36 md:overflow-y-auto"
                  :items="selectedIngredientChipItems"
                  empty-text="Пока ничего не выбрано."
                  :show-favorite-empty-text="
                    includeFavoriteIngredientsInSearch && !favoriteIncludeIngredientNames.length
                  "
                  favorite-empty-text="В списке любимых пока нет ингредиентов."
                  toggle-id="include-favorites-in-search"
                  :toggle-value="includeFavoriteIngredientsInSearch"
                  toggle-name="include_favorites_in_search"
                  toggle-label="Добавить любимые"
                  toggle-tone="cyan"
                  :toggle-visual-state="includeFavoriteToggleVisualState"
                  :list-container-class="
                    displaySelectedIngredients.length
                      ? requireAllIngredients
                        ? 'rounded-[13px] border border-zinc-700 bg-zinc-800/50 p-2'
                        : 'rounded-[13px] border border-transparent bg-transparent p-2'
                      : ''
                  "
                  :item-style="item => selectedTagStyle(item.label)"
                  :item-title="item => `Нажмите, чтобы убрать «${item.label}» из фильтра`"
                  :item-aria-label="item => `Убрать ${item.label} из фильтра`"
                  @update:toggle-value="includeFavoriteIngredientsInSearch = $event"
                  @remove="removeEffectiveSelectedIngredient($event.label)">
                  <template #header-extra>
                    <div class="inline-flex items-center gap-2">
                      <LabBaseSwitch
                        id="require-all-ingredients"
                        v-model="requireAllIngredients"
                        name="require_all_ingredients"
                        label="AND"
                        tone="emerald" />
                      <LabHelpTooltip
                        text="AND-режим: рецепт попадет в выдачу, только если содержит все ингредиенты из блока «Включить»."
                        side="right" />
                    </div>
                  </template>
                </KitchenIngredientScopePanel>
                <KitchenIngredientScopePanel
                  title="Исключить"
                  title-class="text-rose-300"
                  container-class="min-h-16 overflow-visible border-t border-zinc-700 pt-2 md:max-h-36 md:overflow-y-auto md:border-l md:border-t-0 md:pl-3 md:pt-2"
                  :items="excludedIngredientChipItems"
                  empty-text="Ингредиенты не исключены."
                  :show-favorite-empty-text="
                    excludeFavoriteIngredientsInSearch && !favoriteExcludeIngredientNames.length
                  "
                  favorite-empty-text="В списке исключаемых пока нет ингредиентов."
                  toggle-id="exclude-favorites-in-search"
                  :toggle-value="excludeFavoriteIngredientsInSearch"
                  toggle-name="exclude_favorites_in_search"
                  :toggle-label="excludeFavoriteSwitchLabel"
                  toggle-tone="rose"
                  :toggle-visual-state="excludeFavoriteToggleVisualState"
                  :item-style="item => selectedTagStyle(item.label)"
                  :item-title="item => `Нажмите, чтобы убрать «${item.label}» из исключений`"
                  :item-aria-label="item => `Убрать ${item.label} из исключений`"
                  @update:toggle-value="excludeFavoriteIngredientsInSearch = $event"
                  @remove="removeEffectiveExcludedIngredient($event.label)" />
              </div>
            </div>
            <div class="sticky bottom-0 mt-auto border-t border-zinc-800 bg-zinc-950 pt-3 md:hidden">
              <LabBaseButton variant="primary" size="xl" block @click="closeRecipeAdvancedSearch">
                {{ recipeAdvancedSearchCloseLabel }}
              </LabBaseButton>
            </div>
          </LabSpoiler>
        </Teleport>
      </div>
      <p class="text-xs text-zinc-500">{{ activeSearchSummary }}</p>
      <LabNotify v-if="recipesError" :text="recipesError" tone="error" />
      <p v-else-if="recipesLoading" class="text-sm text-zinc-400">Загрузка рецептов...</p>
      <p v-else-if="visibleRecipes.length === 0" class="text-sm text-zinc-500">
        Рецепты не найдены по текущим фильтрам.
      </p>
      <div v-else class="grid gap-3 md:grid-cols-2">
        <KitchenRecipeCard
          v-for="recipe in visibleRecipes"
          :key="recipe.id"
          :recipe="recipe"
          :to="kitchenRecipeLink(recipe)"
          :cover-src="recipeCoverSrc(recipe)"
          :meta-items="recipeMetaItems(recipe, { includeSteps: false })"
          :can-manage="canManageRecipe(recipe)"
          :show-favorite="true"
          :favorite="isRecipeFavorite(recipe.id)"
          :favorite-pending="isRecipeFavoritePending(recipe.id)"
          :delete-disabled="deletePendingRecipeId === recipe.id"
          @edit="openRecipeEditing(recipe)"
          @delete="deleteRecipe(recipe)"
          @toggle-favorite="toggleRecipeFavorite(recipe)" />
      </div>
    </section>
    <section v-show="showRecipeManageSection" class="bg-zinc-900 p-3 space-y-4 sm:p-5">
      <div class="flex flex-wrap items-center gap-2">
        <LabBaseButton
          v-if="showMyRecipesList && isAuthenticated"
          variant="primary"
          size="md"
          button-class="self-start"
          @click="openNewRecipeForm">
          Новый рецепт
        </LabBaseButton>
      </div>
      <div v-if="!authLoaded" class="rounded-lg border border-zinc-700 bg-zinc-900 p-4 text-sm text-zinc-400">
        Проверяем сессию...
      </div>
      <AuthFeatureGateNotice
        v-else-if="!isAuthenticated"
        message="Войдите в аккаунт, чтобы создавать, редактировать и сохранять свои рецепты." />
      <template v-else>
        <div v-if="showMyRecipesList" class="space-y-3 bg-zinc-900 p-3">
          <div
            v-if="userAttentionSummary?.has_unread && userAttentionRecipes.length"
            class="space-y-2 rounded-xl border border-amber-500/30 bg-amber-500/10 p-3">
            <div class="flex flex-wrap items-center justify-between gap-2">
              <p class="text-sm font-medium text-amber-100">{{ userAttentionSummaryText }}</p>
              <LabBaseButton
                variant="ghost"
                size="xs"
                button-class="px-0 text-amber-200 hover:text-amber-100"
                :disabled="userAttentionReadPending"
                @click="markUserAttentionRead">
                Отметить просмотренным
              </LabBaseButton>
            </div>
            <div class="space-y-1 text-xs text-amber-50/90">
              <p v-for="recipe in userAttentionRecipes.slice(0, 4)" :key="`attention:${recipe.id}`">
                {{ recipe.title }}: {{ kitchenRecipeModerationLabel(recipe.moderation_status) }}
                <span v-if="recipe.moderation_status === 'rejected' && recipe.moderation_note">
                  . Причина: {{ recipe.moderation_note }}
                </span>
              </p>
            </div>
          </div>
          <LabNotify v-if="myRecipesError" :text="myRecipesError" tone="error" size="xs" />
          <p v-else-if="myRecipesLoading" class="text-xs text-zinc-500">Загрузка ваших рецептов...</p>
          <div v-else-if="myRecipes.length" class="grid gap-2 lg:grid-cols-2">
            <KitchenRecipeCard
              v-for="recipe in myRecipes"
              :key="`my:${recipe.id}`"
              :recipe="recipe"
              :to="kitchenRecipeLink(recipe)"
              :cover-src="recipeCoverSrc(recipe)"
              :meta-items="recipeMetaItems(recipe, { includeSteps: false })"
              :can-manage="canManageRecipe(recipe)"
              :show-favorite="true"
              :favorite="isRecipeFavorite(recipe.id)"
              :favorite-pending="isRecipeFavoritePending(recipe.id)"
              :show-moderation="true"
              :moderation-text="kitchenRecipeModerationLabel(recipe.moderation_status)"
              :moderation-note="
                recipe.moderation_status === 'rejected' ? String(recipe.moderation_note || '').trim() : ''
              "
              :delete-disabled="deletePendingRecipeId === recipe.id"
              @edit="openRecipeEditing(recipe)"
              @delete="deleteRecipe(recipe)"
              @toggle-favorite="toggleRecipeFavorite(recipe)" />
          </div>
          <p v-else class="text-xs text-zinc-500">У вас пока нет сохранённых рецептов.</p>
        </div>
        <KitchenRecipeEditor
          v-if="showRecipeForm"
          :editor="recipeEditor"
          :images="recipeImages"
          :steps-dn-d="recipeStepsDnD"
          rounded-xl
          :is-current-user-admin="isCurrentUserAdmin"
          :ingredient-category-by-name="ingredientCategoryByName"
          :category-tag-style="categoryTagStyle"
          :recipe-difficulty-scale-options="recipeDifficultyScaleOptions"
          :meal-type-suggestions="KITCHEN_MEAL_TYPE_SUGGESTIONS"
          :diet-type-suggestions="KITCHEN_DIET_TYPE_SUGGESTIONS"
          :cooking-method-suggestions="KITCHEN_COOKING_METHOD_SUGGESTIONS"
          :nationality-suggestions="KITCHEN_NATIONALITY_SUGGESTIONS"
          :unit-suggestions="KITCHEN_UNIT_SUGGESTIONS"
          :image-max-mb="KITCHEN_IMAGE_MAX_MB"
          :image-formats-label="KITCHEN_IMAGE_FORMATS_LABEL"
          :recommended-step-image-min-width="KITCHEN_STEP_IMAGE_RECOMMENDED_MIN_WIDTH"
          :recommended-step-image-min-height="KITCHEN_STEP_IMAGE_RECOMMENDED_MIN_HEIGHT"
          @cancel="cancelRecipeEditing" />
      </template>
    </section>
  </LabBaseSection>
</template>
