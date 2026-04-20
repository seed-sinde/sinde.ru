<script setup lang="ts">
import {
  deleteKitchenRecipe,
  favoriteKitchenRecipe,
  getKitchenManageRecipeById,
  getKitchenRecipeById,
  getKitchenRecipeFavoriteStatus,
  unfavoriteKitchenRecipe
} from '~/composables/useKitchen'
import { buildMediaFileUrl } from '~/utils/mediaUrl'
import { extractKitchenRecipeId, kitchenRecipeSlug, kitchenStepsCountLabel } from '~/utils/kitchenRecipe'
import { buildKitchenRecipeMetaItems } from '~/utils/kitchenRecipeMeta'
definePageMeta({
  title: 'рецепт',
  description: 'Полная версия рецепта.'
})
const route = useRoute()
const router = useRouter()
const slug = computed(() => String(route.params.slug || '').trim())
const recipeId = computed(() => extractKitchenRecipeId(slug.value))
const catalogStore = useKitchenCatalogStore()
await callOnce('kitchen-catalog', async () => {
  await catalogStore.ensureLoaded()
})
const { user, isAuthenticated, ensureLoaded } = useAuth()
const buildStepImageUrl = (imageKey?: string) => buildMediaFileUrl(imageKey)
const stepsCountLabel = (count: number) => kitchenStepsCountLabel(count)
const normalizeTag = (value: string) => value.trim().toLowerCase()
const ingredientCategoryByName = computed(() => {
  const map: Record<string, string> = {}
  for (const item of catalogStore.ingredientItems) {
    const key = normalizeTag(item.name)
    if (!map[key]) {
      map[key] = item.category
    }
  }
  return map
})
const difficultyLabel = (value: string) => catalogStore.labelFor('difficulty', value) || value
const mealTypeLabel = (value: string) => catalogStore.labelFor('meal_type', value) || value
const dietTypeLabel = (value?: string) => {
  const key = String(value || '').trim()
  return key ? catalogStore.labelFor('diet_type', key) || key : 'без ограничений'
}
const cookingMethodLabel = (value?: string) => {
  const key = String(value || '').trim()
  if (!key) return 'не указан'
  return catalogStore.labelFor('cooking_method', key) || key
}
const { data, pending, error } = await useAsyncData(
  () => `kitchen-recipe:${recipeId.value || 'invalid'}`,
  async () => {
    if (!recipeId.value) {
      throw createError({ statusCode: 404, statusMessage: 'Рецепт не найден' })
    }
    try {
      const res = await getKitchenRecipeById(recipeId.value)
      return res?.data || null
    } catch (publicErr: any) {
      const publicStatus = Number(publicErr?.statusCode || publicErr?.status || publicErr?.response?.status || 0)
      if (publicStatus && publicStatus !== 404) {
        throw publicErr
      }
      try {
        const managedRes = await getKitchenManageRecipeById(recipeId.value)
        return managedRes?.data || null
      } catch (manageErr: any) {
        const manageStatus = Number(manageErr?.statusCode || manageErr?.status || manageErr?.response?.status || 0)
        if (!publicStatus || publicStatus === 404) {
          if (!manageStatus || manageStatus === 401 || manageStatus === 403 || manageStatus === 404) {
            throw createError({ statusCode: 404, statusMessage: 'Рецепт не найден' })
          }
        }
        throw manageErr
      }
    }
  },
  { watch: [recipeId] }
)
const recipe = computed(() => data.value || null)
const canManageRecipe = computed(() => {
  const item = recipe.value
  const currentUserId = user.value?.user_id
  if (!item || !currentUserId) return false
  return Boolean((item.owner_user_id && item.owner_user_id === currentUserId) || user.value?.roles?.includes('admin'))
})
const recipeEditLink = computed(() => {
  const item = recipe.value
  if (!item) return '/kitchen/recipes'
  return `/kitchen/edit/${kitchenRecipeSlug(item)}`
})
const recipeMetaItems = computed(() => {
  const item = recipe.value
  if (!item) return [] as Array<{ label: string; value: string | number }>
  return buildKitchenRecipeMetaItems(item, {
    mealTypeLabel,
    dietTypeLabel,
    cookingMethodLabel,
    difficultyLabel,
    stepsCountLabel
  })
})
const groupIngredientsByCategory = ref(false)
const servingsEditing = ref(false)
const servingsInput = ref('')
const servingsOverride = ref<number | null>(null)
const SERVINGS_INPUT_ID = 'recipe-servings-inline'
const actionError = ref<string | null>(null)
const deletePending = ref(false)
const favoritePending = ref(false)
const isFavoriteRecipe = ref(false)
const baseServings = computed(() => {
  const n = Number(recipe.value?.servings || 0)
  return Number.isFinite(n) && n > 0 ? n : 1
})
const parseDecimal = (value: string): number | null => {
  const normalized = String(value || '')
    .trim()
    .replace(',', '.')
    .replace(/\s+/g, '')
  if (!normalized || !/^\d+(?:\.\d+)?$/.test(normalized)) return null
  const n = Number.parseFloat(normalized)
  if (!Number.isFinite(n) || n < 0) return null
  return n
}
const parsePositiveDecimal = (value: string): number | null => {
  const n = parseDecimal(value)
  if (n === null || n <= 0) return null
  return n
}
const formatDecimal = (value: number) => {
  const rounded = Math.round(value * 1000) / 1000
  return String(rounded)
    .replace(/(\.\d*?)0+$/, '$1')
    .replace(/\.$/, '')
}
const currentServings = computed(() => servingsOverride.value || baseServings.value)
const servingsDisplayText = computed(() => formatDecimal(currentServings.value))
const servingsScale = computed(() => {
  const base = baseServings.value
  if (!base || base <= 0) return 1
  return currentServings.value / base
})
const scaleIngredientAmount = (amount?: string, factor = 1) => {
  const raw = String(amount || '').trim()
  if (!raw || factor === 1) return raw
  const parsed = parseDecimal(raw)
  if (parsed === null) return raw
  return formatDecimal(parsed * factor)
}
const scaledIngredients = computed<KitchenIngredient[]>(() => {
  const item = recipe.value
  if (!item) return []
  const factor = servingsScale.value
  if (factor === 1) return item.ingredients
  return item.ingredients.map(ingredient => ({
    ...ingredient,
    amount: scaleIngredientAmount(ingredient.amount, factor)
  }))
})
const startServingsEditing = async () => {
  servingsEditing.value = true
  servingsInput.value = servingsDisplayText.value
  await nextTick()
  if (!import.meta.client) return
  requestAnimationFrame(() => {
    const input = document.getElementById(SERVINGS_INPUT_ID) as HTMLInputElement | null
    input?.focus()
    input?.select()
  })
}
const commitServingsEditing = () => {
  const parsed = parsePositiveDecimal(servingsInput.value)
  servingsOverride.value = parsed
  servingsInput.value = parsed ? formatDecimal(parsed) : formatDecimal(baseServings.value)
  servingsEditing.value = false
}
const cancelServingsEditing = () => {
  servingsInput.value = servingsDisplayText.value
  servingsEditing.value = false
}
const deleteRecipe = async () => {
  const item = recipe.value
  if (!item || !canManageRecipe.value || deletePending.value) return
  deletePending.value = true
  actionError.value = null
  try {
    await deleteKitchenRecipe(item.id)
    await router.push('/kitchen/recipes')
  } catch (err: any) {
    actionError.value = err?.data?.message || err?.message || 'Не удалось удалить рецепт.'
  } finally {
    deletePending.value = false
  }
}
const loadFavoriteStatus = async () => {
  if (!isAuthenticated.value || !recipe.value?.id) {
    isFavoriteRecipe.value = false
    return
  }
  try {
    const res = await getKitchenRecipeFavoriteStatus(recipe.value.id)
    isFavoriteRecipe.value = Boolean(res?.data?.favorited)
  } catch {
    isFavoriteRecipe.value = false
  }
}
const toggleRecipeFavorite = async () => {
  if (!isAuthenticated.value || !recipe.value?.id || favoritePending.value) return
  favoritePending.value = true
  actionError.value = null
  try {
    if (isFavoriteRecipe.value) {
      await unfavoriteKitchenRecipe(recipe.value.id)
      isFavoriteRecipe.value = false
    } else {
      await favoriteKitchenRecipe(recipe.value.id)
      isFavoriteRecipe.value = true
    }
  } catch (err: any) {
    actionError.value = err?.data?.message || err?.message || 'Не удалось обновить избранное.'
  } finally {
    favoritePending.value = false
  }
}
const imagePreviewDialog = reactive({
  open: false,
  url: '',
  title: ''
})
const openImagePreview = (url: string, title: string) => {
  const cleanUrl = String(url || '').trim()
  if (!cleanUrl) return
  imagePreviewDialog.open = true
  imagePreviewDialog.url = cleanUrl
  imagePreviewDialog.title = title
}
const onImagePreviewModelUpdate = (value: boolean) => {
  imagePreviewDialog.open = value
  if (value) return
  imagePreviewDialog.url = ''
  imagePreviewDialog.title = ''
}
watch(
  () => recipe.value?.id,
  async () => {
    await ensureLoaded()
    servingsEditing.value = false
    servingsOverride.value = null
    servingsInput.value = formatDecimal(baseServings.value)
    await loadFavoriteStatus()
  },
  { immediate: true }
)
useHead(() => {
  const item = recipe.value
  if (!item) {
    return {
      title: 'рецепт',
      meta: [{ name: 'description', content: 'Полная версия рецепта.' }]
    }
  }
  return {
    title: `Кухня · ${item.title}`,
    meta: [{ name: 'description', content: item.description || 'Полная версия рецепта.' }]
  }
})
</script>
<template>
  <div>
    <LabNavHeader
      :title="recipe?.title || 'Кухня'"
      :breadcrumb-items="[
        { label: 'Кухня', to: '/kitchen/recipes' },
        { label: recipe?.title || 'Рецепт', current: true }
      ]"
    />
    <section class="space-y-4 p-4">
      <p v-if="pending" class="text-sm text-zinc-400">Загрузка рецепта...</p>
      <LabNotify v-else-if="error" text="Не удалось загрузить рецепт." tone="error" />
      <article v-else-if="recipe" class="space-y-4 bg-transparent">
        <div class="flex flex-wrap items-center gap-2">
          <LabBaseButton
            v-if="isAuthenticated"
            :class="[
              isFavoriteRecipe
                && 'border-rose-500/50 bg-rose-500/15 text-rose-200 hover:bg-rose-500/25'
            ]"
            :disabled="favoritePending"
            :icon="isFavoriteRecipe ? 'ic:round-favorite' : 'ic:round-favorite-border'"
            @click="toggleRecipeFavorite"
          >
            {{ isFavoriteRecipe ? 'В избранном' : 'В избранное' }}
          </LabBaseButton>
          <div v-if="canManageRecipe" class="flex flex-wrap items-center gap-2">
            <NuxtLink
              :to="recipeEditLink"
              class="inline-flex h-9 items-center border border-(--lab-border) px-3 text-xs"
            >
              Редактировать
            </NuxtLink>
            <LabConfirmActionButton
              label="Удалить"
              confirm-label="Подтвердить"
              tooltip="Подтвердить удаление рецепта?"
              class="h-9 px-3 text-xs"
              idle-class="border border-rose-500/50 bg-rose-500/10 text-rose-300 hover:bg-rose-500/20"
              confirm-class="border border-rose-300/90 bg-rose-600 text-white hover:bg-rose-500"
              progress-class="bg-rose-300/45"
              :disabled="deletePending"
              @confirm="deleteRecipe"
            />
          </div>
        </div>
        <LabNotify :text="actionError" tone="error" size="xs" />
        <div class="grid items-start gap-4 lg:grid-cols-3">
          <section class="order-2 space-y-2">
            <div class="flex flex-wrap items-center gap-3 border-b border-(--lab-border) pb-1">
              <h2 class="text-base font-semibold text-amber-300">Ингредиенты</h2>
              <KitchenGroupByCategoryToggle v-model="groupIngredientsByCategory" />
            </div>
            <KitchenGroupedIngredients
              :ingredients="scaledIngredients"
              :category-by-name="ingredientCategoryByName"
              :group-by-category="groupIngredientsByCategory"
              :show-actions="false"
              empty-text="Ингредиенты не указаны."
            />
          </section>
          <section class="order-1 space-y-2">
            <h2 class="border-b border-(--lab-border) pb-1 text-base font-semibold text-amber-300">Фото блюда</h2>
            <LabViewerPreviewButton
              v-if="recipe.cover_image_key"
              :src="buildStepImageUrl(recipe.cover_image_key)"
              alt="Фото готового блюда"
              step-format
              class="block w-full max-w-sm"
              label="Посмотреть"
              @preview="openImagePreview(buildStepImageUrl(recipe.cover_image_key), 'Фото готового блюда')"
            />
            <p v-else class="text-xs text-zinc-500">Фото готового блюда не добавлено.</p>
          </section>
          <section class="order-3 space-y-3">
            <h2 class="border-b border-(--lab-border) pb-1 text-base font-semibold text-amber-300">Детали рецепта</h2>
            <div class="grid grid-cols-2 gap-2">
              <div
                v-for="item in recipeMetaItems"
                :key="`meta:${item.label}`"
                class="rounded-md border px-3 py-2"
                :class="item.label === 'Порций' ? 'cursor-text' : ''"
                @click="item.label === 'Порций' && !servingsEditing ? startServingsEditing() : undefined"
              >
                <div v-if="item.label === 'Порций'" class="flex items-center justify-between gap-2">
                  <p class="text-xs tracking-[0.06em] text-zinc-500 uppercase">{{ item.label }}</p>
                  <p
                    class="max-w-36 truncate text-right text-[10px]"
                    :class="servingsEditing ? 'text-zinc-500' : 'text-cyan-300/80'"
                  >
                    {{ servingsEditing ? 'Enter/Esc' : 'Изменить' }}
                  </p>
                </div>
                <p v-else class="text-xs tracking-[0.06em] text-zinc-500 uppercase">{{ item.label }}</p>
                <template v-if="item.label === 'Порций'">
                  <div class="mt-1">
                    <LabBaseInput
                      v-if="servingsEditing"
                      :id="SERVINGS_INPUT_ID"
                      v-model="servingsInput"
                      type="text"
                      class="h-8 min-h-0 rounded-md px-2 py-1 text-sm"
                      @blur="commitServingsEditing"
                      @keydown.enter.prevent="commitServingsEditing"
                      @keydown.esc.prevent="cancelServingsEditing"
                      @click.stop
                    />
                    <p v-else class="text-sm font-medium text-zinc-100">{{ servingsDisplayText }}</p>
                  </div>
                </template>
                <p v-else class="mt-1 text-sm font-medium text-zinc-100">{{ item.value }}</p>
              </div>
            </div>
            <div v-if="recipe.tags.length" class="flex flex-wrap gap-1">
              <span
                v-for="tag in recipe.tags"
                :key="`tag:${tag}`"
                class="rounded border border-(--lab-border) px-2 py-0.5 text-xs text-zinc-300"
              >
                {{ tag }}
              </span>
            </div>
          </section>
        </div>
        <div class="space-y-2">
          <h2 class="border-b border-(--lab-border) pb-1 text-base font-semibold text-amber-300">Шаги приготовления</h2>
          <ol class="space-y-4 text-sm text-zinc-300">
            <li
              v-for="step in recipe.steps"
              :key="`step:${step.order}`"
              class="border-t pt-3 first:border-t-0 first:pt-0"
            >
              <div class="mb-2 flex items-start gap-2">
                <div
                  class="mt-0.5 inline-flex h-6 w-6 shrink-0 items-center justify-center rounded-full border border-zinc-600 bg-zinc-800 text-xs font-semibold text-zinc-200"
                >
                  {{ step.order }}
                </div>
                <p class="text-sm leading-6 text-zinc-200">{{ step.text }}</p>
              </div>
              <LabViewerPreviewButton
                v-if="step.image_key"
                :src="buildStepImageUrl(step.image_key)"
                alt="Фото шага"
                compact
                step-format
                class="block w-40"
                label="Посмотреть"
                @preview="openImagePreview(buildStepImageUrl(step.image_key), `Фото шага ${step.order}`)"
              />
            </li>
          </ol>
        </div>
      </article>
    </section>
    <LabViewerImage
      v-model="imagePreviewDialog.open"
      :src="imagePreviewDialog.url"
      :title="imagePreviewDialog.title"
      @update:model-value="onImagePreviewModelUpdate"
    />
  </div>
</template>
