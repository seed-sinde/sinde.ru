export type KitchenRecipeMetaItem = {
  label: string
  value: string | number
}
export type KitchenRecipeMetaFormatters = {
  mealTypeLabel: (value: string) => string
  dietTypeLabel: (value?: string) => string
  cookingMethodLabel: (value?: string) => string
  difficultyLabel: (value: string) => string
  stepsCountLabel: (count: number) => string
}
export const kitchenRecipeTotalMinutes = (recipe: Pick<KitchenRecipe, 'prep_minutes' | 'cook_minutes'>) =>
  Number(recipe.prep_minutes || 0) + Number(recipe.cook_minutes || 0)
export const buildKitchenRecipeMetaItems = (
  recipe: KitchenRecipe,
  formatters: KitchenRecipeMetaFormatters,
  options: { includeSteps?: boolean } = {}
): KitchenRecipeMetaItem[] => {
  const includeSteps = options.includeSteps !== false
  const entries: Array<KitchenRecipeMetaItem | null> = [
    { label: 'Тип', value: formatters.mealTypeLabel(recipe.meal_type) },
    { label: 'Питание', value: formatters.dietTypeLabel(recipe.diet_type) },
    recipe.cooking_method ? { label: 'Способ', value: formatters.cookingMethodLabel(recipe.cooking_method) } : null,
    { label: 'Время', value: `${kitchenRecipeTotalMinutes(recipe)} мин` },
    { label: 'Порций', value: recipe.servings },
    { label: 'Калорийность', value: `${recipe.kcal} ккал` },
    recipe.cuisine ? { label: 'Кухня', value: recipe.cuisine } : null,
    { label: 'Сложность', value: formatters.difficultyLabel(recipe.difficulty) },
    includeSteps ? { label: 'Шаги', value: formatters.stepsCountLabel(recipe.steps.length) } : null
  ]
  return entries.filter((entry): entry is KitchenRecipeMetaItem => Boolean(entry && String(entry.value).trim()))
}
