<script setup lang="ts">
import { toRefs } from 'vue'
import type { UseKitchenRecipeEditorResult } from '~/composables/useKitchenRecipeEditor'
import type { UseKitchenRecipeImagesResult } from '~/composables/useKitchenRecipeImages'
import type { UseKitchenRecipeStepsDnDResult } from '~/composables/useKitchenRecipeStepsDnD'
const props = defineProps<{
  editor: UseKitchenRecipeEditorResult
  images: UseKitchenRecipeImagesResult
  stepsDnD: UseKitchenRecipeStepsDnDResult
  isCurrentUserAdmin: boolean
  ingredientCategoryByName: Record<string, string>
  categoryTagStyle: (category: string) => Record<string, string>
  recipeDifficultyScaleOptions: DifficultyScaleOption[]
  mealTypeSuggestions: string[]
  dietTypeSuggestions: string[]
  cookingMethodSuggestions: string[]
  nationalitySuggestions: string[]
  unitSuggestions: string[]
  imageMaxMb: number
  imageFormatsLabel: string
  recommendedStepImageMinWidth: number
  recommendedStepImageMinHeight: number
}>()
const emit = defineEmits<{
  cancel: []
}>()
const {
  isCurrentUserAdmin,
  ingredientCategoryByName,
  recipeDifficultyScaleOptions,
  mealTypeSuggestions,
  dietTypeSuggestions,
  cookingMethodSuggestions,
  nationalitySuggestions,
  unitSuggestions,
  imageMaxMb,
  imageFormatsLabel,
  recommendedStepImageMinWidth,
  recommendedStepImageMinHeight
} = toRefs(props)
const editor = props.editor
const images = props.images
const stepsDnD = props.stepsDnD
const {
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
  editingIngredientIndex,
  ingredientActionInfo,
  editingRecipeModerationNote,
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
  submitRecipe,
  resetRecipeForm
} = editor
const {
  stepImageUploading,
  stepImageError,
  coverImageUploading,
  coverImageError,
  imagePreviewDialog,
  imageCropDialog,
  imageCropUploading,
  stepImageSrc,
  closeImageCropDialog,
  openImagePreview,
  onImagePreviewModelUpdate,
  onImageCropConfirm,
  onStepImageFileChange,
  onStepItemImageFileChange,
  onCoverImageFileChange,
  clearStepDraftImage,
  clearFormStepImage,
  clearCoverImage
} = images
const { draggingStepIndex, stepDragPreview, setStepItemRef, startStepReorderDrag } = stepsDnD
const onCancel = () => {
  emit('cancel')
}
</script>
<template>
  <div class="max-w-xl space-y-4">
    <div v-if="hasRejectedModerationNote" class="border border-amber-500/30 bg-amber-500/10 p-3 text-sm text-amber-100">
      <p class="font-medium">Причина отклонения</p>
      <p class="mt-1 text-amber-50/90">{{ editingRecipeModerationNote }}</p>
    </div>
    <LabField for-id="recipe-title" label="Название блюда">
      <LabBaseInput
        id="recipe-title"
        v-model="form.title"
        name="recipe_title"
        type="text"
        placeholder="например: Куриная грудка с томатами в духовке."
      />
    </LabField>
    <label class="block text-xs text-zinc-400">
      <span class="inline-flex items-center gap-1">
        <span>Описание</span>
        <span class="text-rose-400">*</span>
      </span>
      <LabBaseTextarea
        id="recipe-description"
        v-model="form.description"
        name="recipe_description"
        rows="3"
        placeholder="Кратко опишите вкус, текстуру и идею блюда: для кого подходит, с чем подать и чем оно отличается от других."
      />
    </label>
    <LabSpoiler label="Детали">
      <div class="flex flex-wrap items-end gap-2">
        <label class="w-full text-xs text-zinc-400 sm:w-[calc(50%-0.25rem)] xl:w-52">
          Тип блюда
          <LabBaseInput
            id="recipe-meal-type"
            v-model="form.meal_type"
            name="recipe_meal_type"
            type="text"
            list="kitchen-meal-type-list"
            placeholder="например, основное"
            class="mt-1 w-full"
          />
          <datalist id="kitchen-meal-type-list">
            <option v-for="item in mealTypeSuggestions" :key="`meal:${item}`" :value="item" />
          </datalist>
        </label>
        <label class="w-full text-xs text-zinc-400 sm:w-[calc(50%-0.25rem)] xl:w-52">
          Тип питания
          <LabBaseInput
            id="recipe-diet-type"
            v-model="form.diet_type"
            name="recipe_diet_type"
            type="text"
            list="kitchen-diet-type-list"
            placeholder="например, веганское"
            class="mt-1 w-full"
          />
          <datalist id="kitchen-diet-type-list">
            <option v-for="item in dietTypeSuggestions" :key="`diet:${item}`" :value="item" />
          </datalist>
        </label>
        <label class="w-full text-xs text-zinc-400 sm:w-[calc(50%-0.25rem)] xl:w-52">
          Способ приготовления
          <LabBaseInput
            id="recipe-cooking-method"
            v-model="form.cooking_method"
            name="recipe_cooking_method"
            type="text"
            list="kitchen-cooking-method-list"
            placeholder="например, запекание"
            class="mt-1 w-full"
          />
          <datalist id="kitchen-cooking-method-list">
            <option v-for="item in cookingMethodSuggestions" :key="`cook:${item}`" :value="item" />
          </datalist>
        </label>
        <label class="w-full text-xs text-zinc-400 sm:w-[calc(50%-0.25rem)] xl:w-52">
          Подготовка (мин)
          <LabBaseInput
            id="recipe-prep-minutes"
            v-model.number="form.prep_minutes"
            name="recipe_prep_minutes"
            type="number"
            min="0"
            class="mt-1 w-full"
          />
        </label>
        <label class="w-full text-xs text-zinc-400 sm:w-[calc(50%-0.25rem)] xl:w-52">
          Приготовление (мин)
          <LabBaseInput
            id="recipe-cook-minutes"
            v-model.number="form.cook_minutes"
            name="recipe_cook_minutes"
            type="number"
            min="0"
            class="mt-1 w-full"
          />
        </label>
        <label class="w-full text-xs text-zinc-400 sm:w-[calc(50%-0.25rem)] xl:w-52">
          Порции
          <LabBaseInput
            id="recipe-servings"
            v-model.number="form.servings"
            name="recipe_servings"
            type="number"
            min="1"
            class="mt-1 w-full"
          />
        </label>
        <label class="w-full text-xs text-zinc-400 sm:w-[calc(50%-0.25rem)] xl:w-52">
          Калорийность (ккал)
          <LabBaseInput
            id="recipe-kcal"
            v-model.number="form.kcal"
            name="recipe_kcal"
            type="number"
            min="0"
            class="mt-1 w-full"
          />
        </label>
        <label class="w-full text-xs text-zinc-400 sm:w-[calc(50%-0.25rem)] xl:w-52">
          Национальность
          <LabBaseInput
            id="recipe-cuisine"
            v-model="form.cuisine"
            name="recipe_cuisine"
            type="text"
            list="kitchen-nationality-list"
            placeholder="например, итальянская"
            class="mt-1 w-full"
          />
          <datalist id="kitchen-nationality-list">
            <option v-for="item in nationalitySuggestions" :key="`nationality:${item}`" :value="item" />
          </datalist>
        </label>
        <LabBaseScale
          id="recipe-difficulty"
          v-model="form.difficulty"
          name="recipe_difficulty"
          label="Сложность"
          :options="recipeDifficultyScaleOptions"
          clear-value=""
          field-class="w-full sm:w-[calc(50%-0.25rem)] xl:w-52"
        />
      </div>
    </LabSpoiler>
    <div class="space-y-2">
      <p class="text-xs text-zinc-400">Фото готового блюда</p>
      <div class="flex flex-wrap items-center gap-2">
        <LabBaseFile
          id="recipe-cover-file"
          name="recipe_cover_file"
          accept="image/jpeg,image/png,image/webp,image/gif,image/avif"
          :disabled="coverImageUploading"
          variant="primary"
          :label="coverImageUploading ? 'Загрузка...' : 'Загрузить фото блюда'"
          @change="onCoverImageFileChange"
        />
        <LabBaseButton v-if="coverImageDraft.image_key" size="sm" variant="danger" @click="clearCoverImage">
          Убрать фото
        </LabBaseButton>
      </div>
      <p class="text-xs text-zinc-500">Размер файла до {{ imageMaxMb }} МБ, форматы: {{ imageFormatsLabel }}.</p>
      <LabNotify :text="coverImageError" tone="error" size="xs" />
      <LabViewerPreviewButton
        v-if="coverImageDraft.image_url"
        :src="coverImageDraft.image_url"
        alt="Фото готового блюда"
        step-format
        button-class="block w-full max-w-sm"
        @preview="openImagePreview(coverImageDraft.image_url, 'Фото блюда')"
      />
    </div>
    <div class="grid gap-2 md:grid-cols-6">
      <label class="text-xs text-zinc-400 md:col-span-2">
        <span class="inline-flex items-center gap-1">
          <span>Название ингредиента</span>
          <span class="text-rose-400">*</span>
        </span>
        <div class="relative mt-1">
          <LabBaseInput
            id="recipe-ingredient-name"
            v-model="ingredientDraft.name"
            name="recipe_ingredient_name"
            type="text"
            placeholder="например, томаты черри"
            :class="['w-full']"
            :input-class="
              !editor.hideFormIngredientSuggestions && ingredientDraft.name.trim() && formIngredientSuggestions.length
                ? 'rounded-b-none border-b-transparent'
                : ''
            "
            @input="onIngredientDraftNameInput"
            @keydown.enter.prevent="addFormIngredient"
          />
          <div
            v-if="
              !editor.hideFormIngredientSuggestions && ingredientDraft.name.trim() && formIngredientSuggestions.length
            "
            class="absolute inset-x-0 top-full z-30 -mt-px rounded-b-md border border-t-0 border-zinc-700 bg-zinc-900/95"
          >
            <p class="px-2 py-1 text-xs text-zinc-500">Найдено: {{ formIngredientSuggestions.length }}</p>
            <div class="max-h-40 overflow-x-hidden overflow-y-auto px-2 pb-2">
              <div class="flex flex-wrap gap-2">
                <LabBaseButton
                  v-for="item in formIngredientSuggestions"
                  :key="`form-suggest:${item.name}`"
                  button-class="rounded border px-2 py-1 text-xs transition hover:brightness-110"
                  :style="categoryTagStyle(item.category)"
                  @click="selectFormIngredientSuggestion(item.name)"
                >
                  {{ item.name }}
                </LabBaseButton>
              </div>
            </div>
          </div>
        </div>
      </label>
      <label class="text-xs text-zinc-400">
        Количество
        <LabBaseInput
          id="recipe-ingredient-amount"
          v-model="ingredientDraft.amount"
          name="recipe_ingredient_amount"
          type="number"
          min="0"
          step="any"
          inputmode="decimal"
          placeholder="например, 2"
          class="mt-1 w-full"
          @input="sanitizeIngredientAmountDraft"
          @keydown.enter.prevent="addFormIngredient"
        />
      </label>
      <label class="text-xs text-zinc-400">
        Ед. изм.
        <LabBaseInput
          id="recipe-ingredient-unit"
          v-model="ingredientDraft.unit"
          name="recipe_ingredient_unit"
          type="text"
          list="kitchen-unit-list"
          placeholder="например, гр"
          class="mt-1 w-full"
          @keydown.enter.prevent="addFormIngredient"
        />
      </label>
      <label class="text-xs text-zinc-400">
        Примечание
        <LabBaseInput
          id="recipe-ingredient-note"
          v-model="ingredientDraft.note"
          name="recipe_ingredient_note"
          type="text"
          placeholder="например, по вкусу"
          class="mt-1 w-full"
          @keydown.enter.prevent="addFormIngredient"
        />
      </label>
      <div class="flex flex-col items-start md:self-start">
        <span class="sr-only">Действие</span>
        <div class="mt-6 flex flex-wrap items-center gap-2">
          <LabBaseButton
            :label="editingIngredientIndex !== null ? 'Сохранить' : 'Добавить'"
            variant="primary"
            @click="addFormIngredient"
          />
          <LabBaseButton
            v-if="editingIngredientIndex !== null"
            label="Отмена"
            variant="danger"
            @click="cancelIngredientEditing"
          />
          <LabNotify
            v-if="isEditMode"
            :text="ingredientActionInfo"
            tone="success"
            size="xs"
            as="span"
            class-name="inline-flex items-center whitespace-nowrap"
          />
        </div>
      </div>
      <datalist id="kitchen-unit-list">
        <option v-for="item in unitSuggestions" :key="`unit:${item}`" :value="item" />
      </datalist>
    </div>
    <div class="space-y-2">
      <div class="flex flex-wrap items-center gap-3">
        <h3 class="text-sm font-semibold text-zinc-100">Ингредиенты рецепта</h3>
        <KitchenGroupByCategoryToggle v-model="groupRecipeIngredientsByCategory" />
      </div>
      <KitchenGroupedIngredients
        :ingredients="formIngredients"
        :category-by-name="ingredientCategoryByName"
        :group-by-category="groupRecipeIngredientsByCategory"
        :show-actions="true"
        empty-text="Еще не добавлены."
        @edit="startIngredientEditing"
        @delete="removeFormIngredient"
      />
    </div>
    <div class="space-y-3">
      <h3 class="inline-flex items-center gap-1 text-sm font-semibold text-zinc-100">
        <span>Шаги приготовления</span>
        <span class="text-rose-400">*</span>
      </h3>
      <LabNotify :text="stepImageError" tone="error" size="xs" />
      <ol class="space-y-3 text-sm text-zinc-300">
        <li
          v-for="(step, idx) in formSteps"
          :key="`step:${idx}`"
          :ref="el => setStepItemRef(idx, el)"
          class="rounded-xl border border-zinc-800 bg-zinc-900/55 p-3"
          :class="
            draggingStepIndex === idx ? 'border-cyan-500/60 bg-cyan-500/5 opacity-35 ring-1 ring-cyan-500/50' : ''
          "
        >
          <div class="flex flex-col gap-3 md:flex-row md:items-start">
            <div class="flex items-center gap-2 md:w-28 md:shrink-0">
              <LabBaseButton
                button-class="inline-flex h-9 w-9 touch-none cursor-grab items-center justify-center rounded-lg border border-zinc-700 bg-zinc-950 text-zinc-300 hover:bg-zinc-800 active:cursor-grabbing"
                icon="ic:round-drag-indicator"
                icon-class="h-5 w-5"
                icon-only
                :title="`Перетащить шаг ${idx + 1}`"
                @pointerdown="startStepReorderDrag(idx, $event)"
              />
              <div
                class="inline-flex items-center gap-2 text-[10px] tracking-[0.06em] text-(--lab-text-soft) uppercase"
              >
                <span
                  class="inline-flex h-5 w-5 items-center justify-center rounded-full border border-[color-mix(in_srgb,var(--lab-border)_72%,transparent)] bg-transparent text-[10px] font-semibold text-(--lab-text-primary)"
                >
                  {{ idx + 1 }}
                </span>
                <span>шаг</span>
              </div>
            </div>
            <div class="min-w-0 flex-1 space-y-3">
              <p v-if="draggingStepIndex === idx" class="text-xs font-medium tracking-[0.08em] text-cyan-300 uppercase">
                Отпустите, чтобы оставить шаг здесь
              </p>
              <LabBaseTextarea
                :id="`recipe-step-text-${idx}`"
                :model-value="step.text"
                :name="`recipe_step_text_${idx}`"
                rows="2"
                placeholder="описание шага"
                class="w-full max-w-3xl"
                @update:model-value="updateFormStepText(idx, String($event || ''))"
              />
              <div class="flex flex-wrap items-center gap-2">
                <LabBaseFile
                  :id="`recipe-step-file-${idx}`"
                  variant="primary"
                  :name="`recipe_step_file_${idx}`"
                  accept="image/jpeg,image/png,image/webp,image/gif,image/avif"
                  :disabled="stepImageUploading"
                  :label="step.image_key ? 'Заменить картинку' : 'Загрузить картинку'"
                  @change="onStepItemImageFileChange(idx, $event)"
                />
                <LabConfirmActionButton
                  v-if="step.image_key"
                  variant="danger"
                  label="Убрать картинку"
                  confirm-label="Подтвердить"
                  tooltip="Подтвердить удаление картинки шага?"
                  progress-class="bg-rose-300/45"
                  @click="clearFormStepImage(idx)"
                />
                <LabConfirmActionButton
                  variant="danger"
                  label="Удалить шаг"
                  confirm-label="Подтвердить"
                  tooltip="Подтвердить удаление шага?"
                  progress-class="bg-rose-300/45"
                  @confirm="removeFormStep(idx)"
                />
              </div>
              <div>
                <LabViewerPreviewButton
                  v-if="step.image_key"
                  :src="stepImageSrc(step)"
                  compact
                  step-format
                  label="Полностью"
                  @preview="openImagePreview(stepImageSrc(step), `Изображение шага ${idx + 1}`)"
                />
              </div>
            </div>
          </div>
        </li>
      </ol>
      <Teleport to="body">
        <div
          v-if="stepDragPreview"
          class="pointer-events-none fixed top-0 left-0 z-90"
          :style="{
            transform: `translate3d(${stepDragPreview.x}px, ${stepDragPreview.y}px, 0)`,
            width: `${stepDragPreview.width}px`
          }"
        >
          <div class="overflow-hidden rounded-xl border border-cyan-400/60 bg-zinc-950/95 p-3 backdrop-blur">
            <div class="flex items-center gap-2 text-[10px] tracking-[0.08em] uppercase">
              <Icon name="ic:round-drag-indicator" />
              <span>Шаг {{ stepDragPreview.stepNumber }}</span>
              <span
                class="rounded-full border border-cyan-400/40 bg-cyan-400/10 px-1.5 py-0.5 text-[9px] text-cyan-100"
              >
                Перемещение
              </span>
            </div>
            <p class="mt-2 text-sm leading-5 whitespace-pre-wrap text-zinc-100">
              {{ stepDragPreview.step.text || 'Без описания шага' }}
            </p>
            <div
              v-if="stepDragPreview.imageUrl"
              class="mt-3 inline-flex max-w-44 items-center justify-center rounded-lg border border-zinc-700 bg-zinc-900/80 p-1.5"
            >
              <img
                :src="stepDragPreview.imageUrl"
                alt=""
                class="block h-auto max-h-28 w-auto max-w-full rounded object-contain"
              />
            </div>
          </div>
        </div>
      </Teleport>
      <div class="space-y-3 border border-dashed border-zinc-700 p-3">
        <div class="flex flex-wrap items-center gap-2">
          <p class="text-sm font-medium text-zinc-100">Новый шаг</p>
        </div>
        <label class="block text-xs text-zinc-400">
          <span class="inline-flex items-center gap-1">
            <span>Описание шага</span>
            <span class="text-rose-400">*</span>
          </span>
          <LabBaseTextarea
            id="recipe-step-text"
            v-model="stepDraft.text"
            name="recipe_step_text"
            rows="2"
            placeholder="описание шага"
            class="mt-1 w-full"
          />
        </label>
        <div class="flex flex-wrap items-center gap-2">
          <LabBaseFile
            id="recipe-step-file"
            variant="primary"
            name="recipe_step_file"
            accept="image/jpeg,image/png,image/webp,image/gif,image/avif"
            :disabled="stepImageUploading"
            :label="stepImageUploading ? 'Загрузка...' : 'Загрузить картинку шага'"
            @change="onStepImageFileChange"
          />
          <LabBaseButton
            v-if="stepDraft.image_key"
            variant="danger"
            size="sm"
            label="Убрать картинку"
            @click="clearStepDraftImage"
          />
          <LabBaseButton label="Добавить шаг" size="sm" variant="primary" @click="addFormStep" />
        </div>
        <p class="text-xs text-zinc-500">
          Размер файла до {{ imageMaxMb }} МБ, форматы: {{ imageFormatsLabel }}. Рекомендуемый размер для шага: не
          меньше {{ recommendedStepImageMinWidth }}×{{ recommendedStepImageMinHeight }} px.
        </p>
        <LabViewerPreviewButton
          v-if="stepDraft.image_url"
          :src="stepDraft.image_url"
          step-format
          @preview="openImagePreview(stepDraft.image_url, 'Изображение шага (черновик)')"
        />
      </div>
    </div>
    <label class="block text-xs text-zinc-400">
      Теги (через запятую)
      <LabBaseInput
        id="recipe-tags"
        v-model="formTagsText"
        name="recipe_tags"
        type="text"
        placeholder="быстро, постное, ужин"
        class="mt-1 w-full"
      />
    </label>
    <label v-if="isCurrentUserAdmin" class="inline-flex items-center gap-2 text-sm text-zinc-300">
      <LabBaseCheckbox
        id="recipe-is-public"
        v-model="form.is_public"
        name="recipe_is_public"
        bare
        class="h-4 w-4 rounded border-zinc-700 bg-zinc-900"
      />
      Опубликовать рецепт
    </label>
    <p v-else class="text-xs text-zinc-500">После сохранения рецепт будет отправлен на модерацию администратора.</p>
    <div class="flex flex-wrap gap-2">
      <LabBaseButton :disabled="createPending" variant="primary" @click="submitRecipe">
        {{ editRecipeId ? 'Сохранить изменения' : 'Сохранить рецепт' }}
      </LabBaseButton>
      <LabBaseButton :disabled="createPending" variant="danger" @click="isEditMode ? onCancel() : resetRecipeForm()">
        {{ isEditMode ? 'Отменить редактирование' : 'Очистить форму' }}
      </LabBaseButton>
      <LabNotify :text="createSuccess" tone="success" as="span" class-name="inline-flex items-center px-1" />
    </div>
    <LabNotify :text="createError" tone="error" />
    <LabViewerImage
      v-model="imagePreviewDialog.open"
      :src="imagePreviewDialog.url"
      :title="imagePreviewDialog.title"
      @update:model-value="onImagePreviewModelUpdate"
    />
    <LabCropperModal
      v-if="imageCropDialog.open && imageCropDialog.file"
      :file="imageCropDialog.file"
      :title="imageCropDialog.title"
      :loading="imageCropUploading"
      @cancel="closeImageCropDialog"
      @confirm="onImageCropConfirm"
    />
  </div>
</template>
