import { computed, reactive, ref, watch, type ComputedRef, type Ref } from 'vue'
import { createKitchenRecipe, updateKitchenRecipe } from '~/composables/useKitchen'
import { buildMediaFileUrl } from '~/utils/mediaUrl'
type MaybeReadonlyRef<T> = Ref<T> | ComputedRef<T>
export type KitchenRecipeStepDraft = {
  text: string
  image_key: string
  image_url: string
}
export type KitchenRecipeImageDraft = {
  image_key: string
  image_url: string
}
export type UseKitchenRecipeEditorOptions = {
  activeKitchenTab: Ref<KitchenMainTab>
  isAuthenticated: MaybeReadonlyRef<boolean>
  isCurrentUserAdmin: MaybeReadonlyRef<boolean>
  defaultRecipeDifficulty: MaybeReadonlyRef<string>
  isKitchenEditRoute: MaybeReadonlyRef<boolean>
  getIngredientSuggestions: (query: string, selectedValues: string[]) => Array<{ name: string; category: string }>
  normalizeMealTypeInput: (value: string) => string
  normalizeDietTypeInput: (value: string) => string
  normalizeCookingMethodInput: (value: string) => string
  onOpenRecipeEdit: (recipe: KitchenRecipe) => Promise<void>
  onResetEditRoute: () => Promise<void>
  onAfterSubmit: () => Promise<void>
}
export const useKitchenRecipeEditor = (options: UseKitchenRecipeEditorOptions) => {
  const form = reactive<KitchenRecipeForm>({
    title: '',
    description: '',
    kcal: '',
    prep_minutes: '',
    cook_minutes: '',
    servings: '',
    difficulty: '',
    meal_type: '',
    cooking_method: '',
    cuisine: '',
    diet_type: '',
    is_public: true
  })
  const ingredientDraft = reactive<KitchenIngredient>({
    name: '',
    amount: '',
    unit: '',
    note: ''
  })
  const stepDraft = reactive<KitchenRecipeStepDraft>({
    text: '',
    image_key: '',
    image_url: ''
  })
  const coverImageDraft = reactive<KitchenRecipeImageDraft>({
    image_key: '',
    image_url: ''
  })
  const formIngredients = ref<KitchenIngredient[]>([])
  const formSteps = ref<KitchenFormStep[]>([])
  const formTagsText = ref('')
  const groupRecipeIngredientsByCategory = ref(false)
  const createPending = ref(false)
  const createError = ref<string | null>(null)
  const createSuccess = ref<string | null>(null)
  const editRecipeId = ref<string | null>(null)
  const createRecipeId = ref<string>(globalThis.crypto.randomUUID())
  const editingIngredientIndex = ref<number | null>(null)
  const hideFormIngredientSuggestions = ref(false)
  const ingredientActionInfo = ref<string | null>(null)
  const editingRecipeModerationStatus = ref<AdminModerationStatus | ''>('')
  const editingRecipeModerationNote = ref('')
  let ingredientActionNoticeTick = 0
  let resetTransientStateHandler: (() => void) | null = null
  const activeRecipeUploadId = computed(() => String(editRecipeId.value || '').trim())
  const isEditMode = computed(() => options.activeKitchenTab.value === 'edit')
  const hasRejectedModerationNote = computed(
    () => editingRecipeModerationStatus.value === 'rejected' && Boolean(editingRecipeModerationNote.value)
  )
  const formIngredientSuggestions = computed(() => {
    const selectedNames = formIngredients.value
      .filter((_, idx) => idx !== editingIngredientIndex.value)
      .map(item => item.name)
    return options.getIngredientSuggestions(ingredientDraft.name, selectedNames)
  })
  const normalizeIngredientAmount = (raw: string) => {
    const compact = String(raw || '')
      .trim()
      .replace(/\s+/g, '')
      .replace(/,/g, '.')
    if (!compact) return ''
    if (!/^\d+(?:\.\d+)?$/.test(compact)) return null
    return compact
  }
  const normalizeFormStep = (step: { text: string; image_key?: string | null }): KitchenFormStep => {
    const imageKey = String(step.image_key || '').trim()
    return imageKey ? { text: step.text, image_key: imageKey } : { text: step.text }
  }
  const sanitizeIngredientAmountDraft = () => {
    const raw = String(ingredientDraft.amount || '')
    let next = raw.replace(/,/g, '.').replace(/[^\d.]/g, '')
    const firstDotIndex = next.indexOf('.')
    if (firstDotIndex >= 0) {
      next = `${next.slice(0, firstDotIndex + 1)}${next.slice(firstDotIndex + 1).replace(/\./g, '')}`
    }
    ingredientDraft.amount = next
  }
  const showIngredientActionInfo = (text: string) => {
    ingredientActionNoticeTick += 1
    const currentTick = ingredientActionNoticeTick
    ingredientActionInfo.value = null
    setTimeout(() => {
      if (currentTick !== ingredientActionNoticeTick) return
      ingredientActionInfo.value = text
    }, 0)
  }
  const resetIngredientDraft = () => {
    ingredientDraft.name = ''
    ingredientDraft.amount = ''
    ingredientDraft.unit = ''
    ingredientDraft.note = ''
    hideFormIngredientSuggestions.value = false
  }
  const resetStepDraft = () => {
    stepDraft.text = ''
    stepDraft.image_key = ''
    stepDraft.image_url = ''
  }
  const resetCoverImageDraft = () => {
    coverImageDraft.image_key = ''
    coverImageDraft.image_url = ''
  }
  const runTransientReset = () => {
    resetTransientStateHandler?.()
  }
  const cancelIngredientEditing = () => {
    editingIngredientIndex.value = null
    resetIngredientDraft()
  }
  const selectFormIngredientSuggestion = (name: string) => {
    ingredientDraft.name = name
    hideFormIngredientSuggestions.value = true
  }
  const onIngredientDraftNameInput = () => {
    hideFormIngredientSuggestions.value = false
  }
  const addFormIngredient = () => {
    createError.value = null
    const name = ingredientDraft.name?.trim() || ''
    if (!name) return
    const normalizedAmount = normalizeIngredientAmount(ingredientDraft.amount || '')
    if (normalizedAmount === null) {
      createError.value = 'Количество ингредиента должно быть числом (например: 1, 0.5, 2.25).'
      return
    }
    const nextIngredient: KitchenIngredient = {
      name,
      amount: normalizedAmount,
      unit: ingredientDraft.unit?.trim() || '',
      note: ingredientDraft.note?.trim() || ''
    }
    const wasEditingIngredient = editingIngredientIndex.value !== null
    if (wasEditingIngredient && formIngredients.value[editingIngredientIndex.value!]) {
      const next = [...formIngredients.value]
      next[editingIngredientIndex.value!] = nextIngredient
      formIngredients.value = next
    } else {
      formIngredients.value = [...formIngredients.value, nextIngredient]
    }
    cancelIngredientEditing()
    if (isEditMode.value) {
      showIngredientActionInfo(wasEditingIngredient ? 'Ингредиент обновлён.' : 'Ингредиент добавлен.')
    }
  }
  const removeFormIngredient = (index: number) => {
    formIngredients.value = formIngredients.value.filter((_, idx) => idx !== index)
    if (editingIngredientIndex.value === null) return
    if (editingIngredientIndex.value === index) {
      cancelIngredientEditing()
      return
    }
    if (editingIngredientIndex.value > index) {
      editingIngredientIndex.value -= 1
    }
  }
  const startIngredientEditing = (index: number) => {
    const ingredient = formIngredients.value[index]
    if (!ingredient) return
    editingIngredientIndex.value = index
    ingredientDraft.name = String(ingredient.name || '')
    ingredientDraft.amount = String(ingredient.amount || '')
    ingredientDraft.unit = String(ingredient.unit || '')
    ingredientDraft.note = String(ingredient.note || '')
    hideFormIngredientSuggestions.value = true
  }
  const addFormStep = () => {
    const text = stepDraft.text.trim()
    if (!text) return
    const nextStep = normalizeFormStep({
      text,
      image_key: stepDraft.image_key
    })
    formSteps.value = [...formSteps.value, nextStep]
    resetStepDraft()
  }
  const removeFormStep = (index: number) => {
    formSteps.value = formSteps.value.filter((_, idx) => idx !== index)
  }
  const updateFormStepText = (index: number, value: string) => {
    if (!formSteps.value[index]) return
    const next = [...formSteps.value]
    next[index] = {
      ...next[index],
      text: value
    }
    formSteps.value = next
  }
  const moveFormStep = (fromIndex: number, toIndex: number) => {
    if (fromIndex === toIndex) return
    if (fromIndex < 0 || toIndex < 0) return
    if (fromIndex >= formSteps.value.length || toIndex >= formSteps.value.length) return
    const next = [...formSteps.value]
    const [moved] = next.splice(fromIndex, 1)
    if (!moved) return
    next.splice(toIndex, 0, moved)
    formSteps.value = next
  }
  const resetRecipeForm = (resetOptions: { skipEditRouteRedirect?: boolean } = {}) => {
    editRecipeId.value = null
    createRecipeId.value = globalThis.crypto.randomUUID()
    ingredientActionInfo.value = null
    editingRecipeModerationStatus.value = ''
    editingRecipeModerationNote.value = ''
    form.title = ''
    form.description = ''
    form.kcal = ''
    form.prep_minutes = ''
    form.cook_minutes = ''
    form.servings = ''
    form.difficulty = options.defaultRecipeDifficulty.value
    form.meal_type = ''
    form.cooking_method = ''
    form.cuisine = ''
    form.diet_type = ''
    form.is_public = true
    formIngredients.value = []
    formSteps.value = []
    formTagsText.value = ''
    cancelIngredientEditing()
    resetStepDraft()
    resetCoverImageDraft()
    runTransientReset()
    if (options.isKitchenEditRoute.value && !resetOptions.skipEditRouteRedirect) {
      void options.onResetEditRoute().catch(() => {
        // Ignore navigation duplication.
      })
    }
  }
  const applyRecipeToForm = (recipe: KitchenRecipe, labels: {
    mealTypeLabel: (value: string) => string
    cookingMethodLabel: (value?: string) => string
    dietTypeLabel: (value?: string) => string
  }) => {
    editRecipeId.value = recipe.id
    createRecipeId.value = recipe.id
    createError.value = null
    createSuccess.value = null
    form.title = recipe.title
    form.description = recipe.description || ''
    form.kcal = Number(recipe.kcal || 0)
    form.prep_minutes = Number(recipe.prep_minutes || 0)
    form.cook_minutes = Number(recipe.cook_minutes || 0)
    form.servings = Number(recipe.servings || 1)
    form.difficulty = recipe.difficulty || options.defaultRecipeDifficulty.value
    form.meal_type = labels.mealTypeLabel(recipe.meal_type)
    form.cooking_method = recipe.cooking_method ? labels.cookingMethodLabel(recipe.cooking_method) : ''
    form.cuisine = recipe.cuisine || ''
    form.diet_type = recipe.diet_type ? labels.dietTypeLabel(recipe.diet_type) : ''
    form.is_public = recipe.is_public
    coverImageDraft.image_key = String(recipe.cover_image_key || '')
    coverImageDraft.image_url = buildMediaFileUrl(recipe.cover_image_key)
    formIngredients.value = recipe.ingredients.map(item => ({ ...item }))
    cancelIngredientEditing()
    formSteps.value = recipe.steps.map(step => normalizeFormStep(step))
    formTagsText.value = recipe.tags.join(', ')
    resetStepDraft()
    editingRecipeModerationStatus.value = recipe.moderation_status
    editingRecipeModerationNote.value = String(recipe.moderation_note || '').trim()
    runTransientReset()
    options.activeKitchenTab.value = 'edit'
  }
  const submitRecipe = async () => {
    createError.value = null
    createSuccess.value = null
    if (!options.isAuthenticated.value) {
      createError.value = 'Для публикации рецепта нужен активный аккаунт.'
      return
    }
    if (!form.title.trim()) {
      createError.value = 'Укажите название рецепта.'
      return
    }
    if (!form.description.trim()) {
      createError.value = 'Добавьте описание рецепта.'
      return
    }
    if (formIngredients.value.length === 0) {
      createError.value = 'Добавьте хотя бы один ингредиент.'
      return
    }
    if (formSteps.value.length === 0) {
      createError.value = 'Добавьте хотя бы один шаг приготовления.'
      return
    }
    if (formIngredients.value.some(ingredient => normalizeIngredientAmount(ingredient.amount || '') === null)) {
      createError.value = 'У одного или нескольких ингредиентов неверное количество. Используйте только числа.'
      return
    }
    const tags = formTagsText.value
      .split(/[\n,;]+/)
      .map(item => item.trim())
      .filter(Boolean)
    createPending.value = true
    try {
      const recipeId = editRecipeId.value ? '' : createRecipeId.value
      const coverImageKey = String(coverImageDraft.image_key || '').trim()
      const cookingMethod = options.normalizeCookingMethodInput(form.cooking_method) || ''
      const dietType = options.normalizeDietTypeInput(form.diet_type) || ''
      const payload: KitchenRecipeCreateInput = {
        title: form.title.trim(),
        description: form.description.trim(),
        kcal: Number(form.kcal || 0),
        prep_minutes: Number(form.prep_minutes || 0),
        cook_minutes: Number(form.cook_minutes || 0),
        servings: Number(form.servings || 1),
        difficulty: form.difficulty || options.defaultRecipeDifficulty.value,
        meal_type: options.normalizeMealTypeInput(form.meal_type) || 'other',
        cuisine: form.cuisine.trim(),
        ingredients: formIngredients.value,
        steps: formSteps.value.map(step => normalizeFormStep(step)),
        tags,
        is_public: form.is_public
      }
      if (recipeId) {
        payload.recipe_id = recipeId
      }
      if (coverImageKey) {
        payload.cover_image_key = coverImageKey
      }
      if (cookingMethod) {
        payload.cooking_method = cookingMethod
      }
      if (dietType) {
        payload.diet_type = dietType
      }
      if (editRecipeId.value) {
        await updateKitchenRecipe(editRecipeId.value, payload)
        createSuccess.value = options.isCurrentUserAdmin.value
          ? 'Рецепт обновлён.'
          : 'Рецепт обновлён и отправлен на модерацию.'
      } else {
        await createKitchenRecipe(payload)
        createSuccess.value = options.isCurrentUserAdmin.value
          ? 'Рецепт опубликован.'
          : 'Рецепт отправлен на модерацию.'
        resetRecipeForm({ skipEditRouteRedirect: true })
        options.activeKitchenTab.value = 'create'
      }
      await options.onAfterSubmit()
    } catch (err: any) {
      createError.value = err?.data?.message || err?.message || 'Не удалось сохранить рецепт.'
    } finally {
      createPending.value = false
    }
  }
  const openRecipeEditing = async (recipe: KitchenRecipe) => {
    await options.onOpenRecipeEdit(recipe)
  }
  const setTransientResetHandler = (handler: (() => void) | null) => {
    resetTransientStateHandler = handler
  }
  watch(
    () => options.defaultRecipeDifficulty.value,
    nextValue => {
      if (!form.difficulty) {
        form.difficulty = nextValue
      }
    },
    { immediate: true }
  )
  return {
    form,
    ingredientDraft,
    stepDraft,
    coverImageDraft,
    formIngredients,
    formSteps,
    formTagsText,
    groupRecipeIngredientsByCategory,
    createPending,
    createError,
    createSuccess,
    editRecipeId,
    createRecipeId,
    editingIngredientIndex,
    hideFormIngredientSuggestions,
    ingredientActionInfo,
    editingRecipeModerationStatus,
    editingRecipeModerationNote,
    activeRecipeUploadId,
    isEditMode,
    hasRejectedModerationNote,
    formIngredientSuggestions,
    sanitizeIngredientAmountDraft,
    cancelIngredientEditing,
    selectFormIngredientSuggestion,
    onIngredientDraftNameInput,
    addFormIngredient,
    removeFormIngredient,
    startIngredientEditing,
    addFormStep,
    removeFormStep,
    updateFormStepText,
    moveFormStep,
    submitRecipe,
    resetRecipeForm,
    applyRecipeToForm,
    openRecipeEditing,
    setTransientResetHandler
  }
}
export type UseKitchenRecipeEditorResult = ReturnType<typeof useKitchenRecipeEditor>
