import { onScopeDispose, ref, type Ref } from 'vue'
export type UseKitchenRecipeStepsDnDOptions = {
  formSteps: Ref<KitchenFormStep[]>
  moveFormStep: (fromIndex: number, toIndex: number) => void
  buildStepImageUrl: (imageKey?: string) => string
}
export const useKitchenRecipeStepsDnD = (options: UseKitchenRecipeStepsDnDOptions) => {
  const stepItemRefs = ref<Array<HTMLElement | null>>([])
  const draggingStepIndex = ref<number | null>(null)
  const stepDragPreview = ref<KitchenStepDragPreview | null>(null)
  let activeStepPointerId: number | null = null
  const setStepItemRef = (index: number, el: Element | { $el?: Element } | null) => {
    const node =
      el instanceof HTMLElement
        ? el
        : el && typeof el === 'object' && '$el' in el && el.$el instanceof HTMLElement
          ? el.$el
          : null
    stepItemRefs.value[index] = node
  }
  const updateStepDragPreviewPosition = (clientX: number, clientY: number) => {
    if (!stepDragPreview.value) return
    const width = stepDragPreview.value.width
    const maxX = Math.max(12, window.innerWidth - width - 12)
    const nextX = Math.min(Math.max(12, clientX + 16), maxX)
    const nextY = Math.max(12, clientY - 28)
    stepDragPreview.value = {
      ...stepDragPreview.value,
      x: nextX,
      y: nextY
    }
  }
  const findStepIndexByPointerY = (clientY: number) => {
    let targetIndex = options.formSteps.value.length - 1
    for (let index = 0; index < options.formSteps.value.length; index += 1) {
      const element = stepItemRefs.value[index]
      if (!element) continue
      const rect = element.getBoundingClientRect()
      if (clientY < rect.top + rect.height / 2) {
        targetIndex = index
        break
      }
    }
    return targetIndex
  }
  const cleanupStepReorderDrag = () => {
    activeStepPointerId = null
    draggingStepIndex.value = null
    stepDragPreview.value = null
    document.body.style.userSelect = ''
    window.removeEventListener('pointermove', onStepReorderPointerMove)
    window.removeEventListener('pointerup', onStepReorderPointerUp)
    window.removeEventListener('pointercancel', onStepReorderPointerUp)
  }
  const onStepReorderPointerMove = (event: PointerEvent) => {
    if (draggingStepIndex.value === null) return
    if (activeStepPointerId !== null && event.pointerId !== activeStepPointerId) return
    updateStepDragPreviewPosition(event.clientX, event.clientY)
    const nextIndex = findStepIndexByPointerY(event.clientY)
    if (nextIndex === draggingStepIndex.value) return
    options.moveFormStep(draggingStepIndex.value, nextIndex)
    draggingStepIndex.value = nextIndex
  }
  const onStepReorderPointerUp = (event: PointerEvent) => {
    if (activeStepPointerId !== null && event.pointerId !== activeStepPointerId) return
    cleanupStepReorderDrag()
  }
  const startStepReorderDrag = (index: number, event: PointerEvent) => {
    if (event.pointerType === 'mouse' && event.button !== 0) return
    const step = options.formSteps.value[index]
    if (!step) return
    event.preventDefault()
    const sourceElement = stepItemRefs.value[index]
    const sourceRect = sourceElement?.getBoundingClientRect()
    const previewWidth = sourceRect
      ? Math.min(Math.max(sourceRect.width, 280), window.innerWidth - 24)
      : Math.min(420, window.innerWidth - 24)
    activeStepPointerId = event.pointerId
    draggingStepIndex.value = index
    stepDragPreview.value = {
      step: {
        text: step.text,
        image_key: step.image_key
      },
      stepNumber: index + 1,
      imageUrl: options.buildStepImageUrl(step.image_key),
      width: previewWidth,
      x: 12,
      y: 12
    }
    updateStepDragPreviewPosition(event.clientX, event.clientY)
    document.body.style.userSelect = 'none'
    window.addEventListener('pointermove', onStepReorderPointerMove)
    window.addEventListener('pointerup', onStepReorderPointerUp)
    window.addEventListener('pointercancel', onStepReorderPointerUp)
  }
  onScopeDispose(() => {
    if (!import.meta.client) return
    cleanupStepReorderDrag()
  })
  return {
    stepItemRefs,
    draggingStepIndex,
    stepDragPreview,
    setStepItemRef,
    startStepReorderDrag
  }
}
export type UseKitchenRecipeStepsDnDResult = ReturnType<typeof useKitchenRecipeStepsDnD>
