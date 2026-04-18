import type {AdminModerationStatus} from "./admin"
export type ImageCropTarget = "cover" | "step"
export type KitchenStepImageTarget = "draft" | "item"
export type FavoriteToggleVisualState = "off" | "partial" | "on"
export type RecipeToolbarAction = "favorite" | "search" | "reset"
export type KitchenMainTab = "recipes" | "ingredients" | "my-recipes" | "create" | "edit"
export type KitchenFormStep = {text: string; image_key?: string}
export type KitchenRecipeForm = {
  title: string
  description: string
  kcal: number | ""
  prep_minutes: number | ""
  cook_minutes: number | ""
  servings: number | ""
  // Catalog code from kitchen_filter_options(kind='difficulty'), not a fixed ordinal level.
  difficulty: string
  meal_type: string
  cooking_method: string
  cuisine: string
  diet_type: string
  is_public: boolean
}
export type KitchenStepDragPreview = {
  step: KitchenFormStep
  stepNumber: number
  imageUrl: string
  width: number
  x: number
  y: number
}
export type KitchenCatalogItem = {
  id: string
  ingredient_id?: number
  name: string
  category: string
  description?: string
  protein_g?: number | null
  fat_g?: number | null
  carbs_g?: number | null
  kcal?: number | null
}
export type KitchenCatalogGroup = {category: string; items: KitchenCatalogItem[]}
export type KitchenIngredient = {
  name: string
  amount?: string
  unit?: string
  note?: string
}
export type KitchenStep = {
  order: number
  text: string
  image_key?: string
}
export type KitchenRecipe = {
  id: string
  owner_user_id?: string | null
  title: string
  description: string
  cover_image_key?: string
  kcal: number
  prep_minutes: number
  cook_minutes: number
  servings: number
  // Catalog code from kitchen_filter_options(kind='difficulty'), not a fixed ordinal level.
  difficulty: string
  meal_type: string
  cooking_method?: string
  cuisine: string
  diet_type?: string
  ingredients: KitchenIngredient[]
  steps: KitchenStep[]
  tags: string[]
  is_public: boolean
  moderation_status: AdminModerationStatus
  moderated_by_user_id?: string | null
  moderated_at?: string | null
  moderation_note?: string
  created_at: string
  updated_at: string
}
export type KitchenUserIngredient = {
  ingredient_id: string
  user_id: string
  name: string
  category: string
  normalized_name: string
  created_at: string
  updated_at: string
}
export type KitchenFavoriteIngredient = {
  ingredient_id: number
  user_id: string
  list_type: "include" | "exclude"
  name: string
  category: string
  normalized_name: string
  created_at: string
}
export type KitchenCategory = {
  category_id: number
  slug: string
  label: string
  description: string
  sort_order: number
}
export type KitchenCatalogIngredient = {
  ingredient_id: number
  name: string
  normalized_name: string
  category_id: number
  category_label: string
  description: string
  protein_g?: number | null
  fat_g?: number | null
  carbs_g?: number | null
  kcal?: number | null
  sort_order: number
}
export type KitchenFilterOption = {
  option_id: number
  kind: string
  code: string
  label: string
  description: string
  sort_order: number
}
export type KitchenCatalogPayload = {
  categories: KitchenCategory[]
  ingredients: KitchenCatalogIngredient[]
  filter_options: KitchenFilterOption[]
}
export type KitchenRecipeCreateInput = {
  recipe_id?: string
  title: string
  description: string
  cover_image_key?: string
  kcal?: number
  prep_minutes: number
  cook_minutes: number
  servings: number
  // Catalog code from kitchen_filter_options(kind='difficulty'), not a fixed ordinal level.
  difficulty: string
  meal_type: string
  cooking_method?: string
  cuisine: string
  diet_type?: string
  ingredients: KitchenIngredient[]
  steps: {text: string; image_key?: string}[]
  tags?: string[]
  is_public?: boolean
}
export type MediaUploadResult = {
  image_key: string
  image_url: string
  reused?: boolean
  file_hash?: string
  section?: string
  collection?: string
}
export type KitchenAccountIngredients = {
  custom: KitchenUserIngredient[]
  favorites: KitchenFavoriteIngredient[]
  favorites_include?: KitchenFavoriteIngredient[]
  favorites_exclude?: KitchenFavoriteIngredient[]
}
export type KitchenSearchInput = {
  query?: string
  ingredients?: string[]
  excludedIngredients?: string[]
  requireAllIngredients?: boolean
  // Filter by catalog code from kitchen_filter_options(kind='difficulty').
  difficulty?: string
  mealType?: string
  cookingMethod?: string
  cuisine?: string
  dietType?: string
  kcalMin?: number | null
  kcalMax?: number | null
  maxTotalMinutes?: number | null
  servingsMin?: number | null
  servingsMax?: number | null
  limit?: number
  offset?: number
}
export type IndexedIngredient = {index: number; item: KitchenIngredient}
