const kitchenFavoriteStatusInFlight = new Map<string, Promise<ApiResponseWithData<{ favorited: boolean }>>>()
const { json: useApiJson } = useAPI()
export const getKitchenIngredients = async () => {
  return await useApiJson<ApiResponseWithData<{ categories: string[] }>>('/kitchen/ingredients', {
    method: 'GET'
  })
}
export const getKitchenCatalog = async () => {
  return await useApiJson<ApiResponseWithData<KitchenCatalogPayload>>('/kitchen/catalog', {
    method: 'GET'
  })
}
export const createKitchenRecipe = async (payload: KitchenRecipeCreateInput) => {
  return await useApiJson<ApiResponseWithData<KitchenRecipe>>('/kitchen/recipes', {
    method: 'POST',
    body: payload,
  })
}
export const updateKitchenRecipe = async (id: string, payload: KitchenRecipeCreateInput) => {
  return await useApiJson<ApiResponseWithData<KitchenRecipe>>(`/kitchen/recipes/${id}`, {
    method: 'PATCH',
    body: payload,
  })
}
export const deleteKitchenRecipe = async (id: string) => {
  return await useApiJson<ApiResponseWithData<{ deleted: boolean }>>(`/kitchen/recipes/${id}`, {
    method: 'DELETE',
  })
}
export const getKitchenLatestRecipes = async (limit = 20) => {
  return await useApiJson<ApiResponseWithData<{ items: KitchenRecipe[], total_public?: number }>>(`/kitchen/recipes/latest?limit=${limit}`, {
    method: 'GET',
  })
}
export const searchKitchenRecipes = async (payload: KitchenSearchInput) => {
  return await useApiJson<ApiResponseWithData<{
    items: KitchenRecipe[]
    mode: 'any' | 'all'
    ingredients: string[]
  }>>('/kitchen/recipes/search', {
    method: 'POST',
    body: payload,
  })
}
export const getKitchenRecipeById = async (id: string) => {
  return await useApiJson<ApiResponseWithData<KitchenRecipe>>(`/kitchen/recipes/${id}`, {
    method: 'GET',
  })
}
export const getKitchenManageRecipeById = async (id: string) => {
  return await useApiJson<ApiResponseWithData<KitchenRecipe>>(`/kitchen/recipes/manage/${id}`, {
    method: 'GET',
  })
}
export const getKitchenRecipeFavoriteStatus = async (id: string) => {
  const normalizedId = String(id || '').trim()
  if (!normalizedId) {
    return {
      ok: true,
      data: { favorited: false },
    } as ApiResponseWithData<{ favorited: boolean }>
  }
  const existing = kitchenFavoriteStatusInFlight.get(normalizedId)
  if (existing) return await existing
  const request = useApiJson<ApiResponseWithData<{ favorited: boolean }>>(`/kitchen/recipes/${normalizedId}/favorite`, {
    method: 'GET',
  })
  kitchenFavoriteStatusInFlight.set(normalizedId, request)
  try {
    return await request
  } finally {
    kitchenFavoriteStatusInFlight.delete(normalizedId)
  }
}
export const favoriteKitchenRecipe = async (id: string) => {
  return await useApiJson<ApiResponseWithData<{ favorited: boolean }>>(`/kitchen/recipes/${id}/favorite`, {
    method: 'POST',
  })
}
export const unfavoriteKitchenRecipe = async (id: string) => {
  return await useApiJson<ApiResponseWithData<{ favorited: boolean }>>(`/kitchen/recipes/${id}/favorite`, {
    method: 'DELETE',
  })
}
export const getMyKitchenRecipes = async () => {
  return await useApiJson<ApiResponseWithData<{ items: KitchenRecipe[] }>>('/kitchen/recipes/mine', {
    method: 'GET',
  })
}
export const getKitchenAccountIngredients = async () => {
  return await useApiJson<ApiResponseWithData<KitchenAccountIngredients>>('/kitchen/ingredients/account', {
    method: 'GET',
  })
}
export const createKitchenCustomIngredient = async (payload: { name: string, category: string }) => {
  return await useApiJson<ApiResponseWithData<KitchenUserIngredient>>('/kitchen/ingredients/custom', {
    method: 'POST',
    body: payload,
  })
}
export const deleteKitchenCustomIngredient = async (id: string) => {
  return await useApiJson<ApiResponseWithData<{ deleted: boolean }>>(`/kitchen/ingredients/custom/${id}`, {
    method: 'DELETE',
  })
}
export const favoriteKitchenIngredient = async (payload: { ingredient_id: number, list_type?: 'include' | 'exclude' }) => {
  return await useApiJson<ApiResponseWithData<KitchenFavoriteIngredient>>('/kitchen/ingredients/favorites', {
    method: 'POST',
    body: payload,
  })
}
export const unfavoriteKitchenIngredient = async (ingredientId: number, listType: 'include' | 'exclude' = 'include') => {
  const query = new URLSearchParams()
  if (listType) query.set('list_type', listType)
  const suffix = query.toString()
  return await useApiJson<ApiResponseWithData<{ deleted: boolean }>>(`/kitchen/ingredients/favorites/${ingredientId}${suffix ? `?${suffix}` : ''}`, {
    method: 'DELETE',
  })
}
export const adminListKitchenModerationRecipes = async (params?: {
  status?: AdminModerationStatus | 'all' | ''
  limit?: number
  offset?: number
}) => {
  const query = new URLSearchParams()
  if (params?.status) query.set('status', String(params.status))
  if (typeof params?.limit === 'number') query.set('limit', String(params.limit))
  if (typeof params?.offset === 'number') query.set('offset', String(params.offset))
  const suffix = query.toString()
  return await useApiJson<ApiResponseWithData<{
    items: KitchenRecipe[]
    total: number
    limit: number
    offset: number
    status?: AdminModerationStatus | 'all' | ''
    status_totals?: Record<string, number>
  }>>(`/kitchen/admin/recipes/moderation${suffix ? `?${suffix}` : ''}`, {
    method: 'GET',
  })
}
export const adminModerateKitchenRecipe = async (id: string, approve: boolean, note = '') => {
  return await useApiJson<ApiResponseWithData<KitchenRecipe>>(`/kitchen/admin/recipes/${id}/moderate`, {
    method: 'POST',
    body: { approve, note },
  })
}
export const adminChangeKitchenRecipeOwner = async (id: string, ownerUserId: string) => {
  return await useApiJson<ApiResponseWithData<KitchenRecipe>>(`/kitchen/admin/recipes/${id}/owner`, {
    method: 'POST',
    body: { owner_user_id: ownerUserId },
  })
}
