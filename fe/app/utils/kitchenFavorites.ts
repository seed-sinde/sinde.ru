export type KitchenFavoriteListType = KitchenFavoriteIngredient['list_type']
export const normalizeKitchenFavoriteListType = (value: unknown): KitchenFavoriteListType =>
  String(value || '')
    .trim()
    .toLowerCase() === 'exclude'
    ? 'exclude'
    : 'include'
export const normalizeKitchenFavoriteIngredient = (
  item: KitchenFavoriteIngredient,
  fallback: KitchenFavoriteListType,
): KitchenFavoriteIngredient => ({
  ...item,
  list_type: normalizeKitchenFavoriteListType(item?.list_type || fallback),
})
export const upsertKitchenFavoriteIngredient = (
  items: KitchenFavoriteIngredient[],
  item: KitchenFavoriteIngredient,
  fallback: KitchenFavoriteListType,
) => {
  const nextItem = normalizeKitchenFavoriteIngredient(item, fallback)
  const ingredientId = Number(nextItem.ingredient_id || 0)
  const nextItems = items.filter(existing => Number(existing.ingredient_id || 0) !== ingredientId)
  return [nextItem, ...nextItems]
}
export const removeKitchenFavoriteIngredientById = (
  items: KitchenFavoriteIngredient[],
  ingredientId: number,
) => items.filter(item => Number(item.ingredient_id || 0) !== ingredientId)
