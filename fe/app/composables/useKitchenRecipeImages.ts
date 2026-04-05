import { computed, reactive, ref, type ComputedRef, type Ref } from 'vue'
import { uploadMediaFile } from '~/composables/useMediaUpload'
import { buildMediaFileUrl } from '~/utils/mediaUrl'
import type { KitchenRecipeImageDraft, KitchenRecipeStepDraft } from './useKitchenRecipeEditor'
type MaybeReadonlyRef<T> = Ref<T> | ComputedRef<T>
export type UseKitchenRecipeImagesOptions = {
  activeRecipeUploadId: MaybeReadonlyRef<string>
  formSteps: Ref<KitchenFormStep[]>
  stepDraft: KitchenRecipeStepDraft
  coverImageDraft: KitchenRecipeImageDraft
}
export const useKitchenRecipeImages = (options: UseKitchenRecipeImagesOptions) => {
  const stepImageUploading = ref(false)
  const stepImageError = ref<string | null>(null)
  const coverImageUploading = ref(false)
  const coverImageError = ref<string | null>(null)
  const imagePreviewDialog = reactive({
    open: false,
    url: '',
    title: ''
  })
  const imageCropDialog = reactive<{
    open: boolean
    file: File | null
    target: ImageCropTarget
    title: string
  }>({
    open: false,
    file: null,
    target: 'cover',
    title: ''
  })
  const imageCropUploading = ref(false)
  const stepImageTargetMode = ref<KitchenStepImageTarget>('draft')
  const stepImageTargetIndex = ref<number | null>(null)
  const buildStepImageUrl = (imageKey?: string) => buildMediaFileUrl(imageKey)
  const stepImageSrc = (step: { image_key?: string }) => buildStepImageUrl(step.image_key)
  const resolveUploadedImageKey = (response: any) => String(response?.data?.image_key || response?.image_key || '').trim()
  const openImageCropDialog = (file: File, target: ImageCropTarget) => {
    if (!file) return
    imageCropDialog.target = target
    imageCropDialog.title = target === 'cover' ? 'Кадрирование фото блюда' : 'Кадрирование изображения шага'
    imageCropDialog.file = file
    imageCropDialog.open = true
  }
  const closeImageCropDialog = (force = false) => {
    if (imageCropUploading.value && !force) return
    imageCropDialog.open = false
    imageCropDialog.file = null
  }
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
  const applyUploadedKitchenImage = (target: ImageCropTarget, imageKey: string) => {
    if (target === 'cover') {
      options.coverImageDraft.image_key = imageKey
      options.coverImageDraft.image_url = buildStepImageUrl(imageKey)
      return
    }
    if (
      stepImageTargetMode.value === 'item' &&
      stepImageTargetIndex.value !== null &&
      options.formSteps.value[stepImageTargetIndex.value]
    ) {
      const next = [...options.formSteps.value]
      next[stepImageTargetIndex.value] = {
        text: next[stepImageTargetIndex.value]?.text ?? '',
        image_key: imageKey
      }
      options.formSteps.value = next
      return
    }
    options.stepDraft.image_key = imageKey
    options.stepDraft.image_url = buildStepImageUrl(imageKey)
  }
  const onImageCropConfirm = async (file: File) => {
    const target = imageCropDialog.target
    imageCropUploading.value = true
    if (target === 'cover') {
      coverImageError.value = null
      coverImageUploading.value = true
    } else {
      stepImageError.value = null
      stepImageUploading.value = true
    }
    try {
      const response = await uploadMediaFile(file, {
        section: 'kitchen',
        collection: 'recipes',
        recipeId: String(options.activeRecipeUploadId.value || '').trim() || undefined
      })
      const imageKey = resolveUploadedImageKey(response)
      if (!imageKey) {
        throw new Error('Сервер не вернул ключ изображения.')
      }
      applyUploadedKitchenImage(target, imageKey)
      closeImageCropDialog(true)
    } catch (err: any) {
      const message = err?.data?.message || err?.message || 'Не удалось загрузить изображение.'
      if (target === 'cover') {
        coverImageError.value = message
      } else {
        stepImageError.value = message
      }
    } finally {
      imageCropUploading.value = false
      coverImageUploading.value = false
      stepImageUploading.value = false
    }
  }
  const onStepImageFileChange = (event: Event) => {
    const input = event.target as HTMLInputElement | null
    const file = input?.files?.[0] || null
    if (input) input.value = ''
    if (!file) return
    stepImageError.value = null
    stepImageTargetMode.value = 'draft'
    stepImageTargetIndex.value = null
    openImageCropDialog(file, 'step')
  }
  const onStepItemImageFileChange = (index: number, event: Event) => {
    const input = event.target as HTMLInputElement | null
    const file = input?.files?.[0] || null
    if (input) input.value = ''
    if (!file) return
    stepImageError.value = null
    stepImageTargetMode.value = 'item'
    stepImageTargetIndex.value = index
    openImageCropDialog(file, 'step')
  }
  const onCoverImageFileChange = (event: Event) => {
    const input = event.target as HTMLInputElement | null
    const file = input?.files?.[0] || null
    if (input) input.value = ''
    if (!file) return
    coverImageError.value = null
    openImageCropDialog(file, 'cover')
  }
  const clearStepDraftImage = () => {
    options.stepDraft.image_key = ''
    options.stepDraft.image_url = ''
    stepImageError.value = null
  }
  const clearFormStepImage = (index: number) => {
    if (!options.formSteps.value[index]) return
    const next = [...options.formSteps.value]
    next[index] = {
      text: next[index]?.text ?? '',
      image_key: ''
    }
    options.formSteps.value = next
    stepImageError.value = null
  }
  const clearCoverImage = () => {
    const previousUrl = String(options.coverImageDraft.image_url || '')
    options.coverImageDraft.image_key = ''
    options.coverImageDraft.image_url = ''
    coverImageError.value = null
    if (imagePreviewDialog.url === previousUrl) {
      onImagePreviewModelUpdate(false)
    }
  }
  const resetRecipeImageUi = () => {
    stepImageError.value = null
    coverImageError.value = null
    stepImageUploading.value = false
    coverImageUploading.value = false
    imageCropUploading.value = false
    stepImageTargetMode.value = 'draft'
    stepImageTargetIndex.value = null
    closeImageCropDialog(true)
    onImagePreviewModelUpdate(false)
  }
  return {
    stepImageUploading,
    stepImageError,
    coverImageUploading,
    coverImageError,
    imagePreviewDialog,
    imageCropDialog,
    imageCropUploading,
    buildStepImageUrl,
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
    clearCoverImage,
    resetRecipeImageUi
  }
}
export type UseKitchenRecipeImagesResult = ReturnType<typeof useKitchenRecipeImages>
