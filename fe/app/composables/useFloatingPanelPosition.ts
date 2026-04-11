import { nextTick, onBeforeUnmount, ref, toValue, type MaybeRefOrGetter, type Ref } from 'vue'
export type FloatingPanelAlign = 'left' | 'right'
export type FloatingPanelSide = 'top' | 'bottom' | 'left' | 'right'
/**
 * Positions a teleported floating panel inside the viewport and flips it when space is limited.
 */
export const useFloatingPanelPosition = (options: {
  triggerRef: Ref<HTMLElement | null>
  panelRef: Ref<HTMLElement | null>
  align?: MaybeRefOrGetter<FloatingPanelAlign>
  side?: MaybeRefOrGetter<FloatingPanelSide>
  offset?: MaybeRefOrGetter<number>
  crossAxisOffset?: MaybeRefOrGetter<number>
  viewportPadding?: MaybeRefOrGetter<number>
  matchTriggerWidth?: MaybeRefOrGetter<boolean>
}) => {
  const panelStyle = ref<Record<string, string>>({})
  const resolvedSide = ref<FloatingPanelSide>(toValue(options.side) || 'bottom')
  const resolvedAlign = ref<FloatingPanelAlign>(toValue(options.align) || 'left')
  const rafId = ref<number | null>(null)
  /** Restricts a value to a numeric range. */
  const clamp = (value: number, min: number, max: number) => Math.min(Math.max(value, min), max)
  /** Cancels a queued animation-frame position update. */
  const cancelScheduledUpdate = () => {
    if (!import.meta.client) return
    if (rafId.value === null) return
    cancelAnimationFrame(rafId.value)
    rafId.value = null
  }
  /** Restores initial placement state when the panel closes. */
  const resetPosition = () => {
    panelStyle.value = {}
    resolvedSide.value = toValue(options.side) || 'bottom'
    resolvedAlign.value = toValue(options.align) || 'left'
    cancelScheduledUpdate()
  }
  /** Measures the panel and computes the most visible on-screen placement. */
  const updatePosition = async () => {
    if (!import.meta.client) return
    await nextTick()
    const triggerEl = options.triggerRef.value
    const panelEl = options.panelRef.value
    if (!triggerEl || !panelEl) return
    const triggerRect = triggerEl.getBoundingClientRect()
    const side = toValue(options.side) || 'bottom'
    const align = toValue(options.align) || 'left'
    const offset = toValue(options.offset) ?? 8
    const crossAxisOffset = toValue(options.crossAxisOffset) ?? 10
    const viewportPadding = toValue(options.viewportPadding) ?? 12
    const matchTriggerWidth = Boolean(toValue(options.matchTriggerWidth))
    panelEl.style.visibility = 'hidden'
    panelEl.style.left = '0px'
    panelEl.style.top = '0px'
    panelEl.style.right = 'auto'
    panelEl.style.bottom = 'auto'
    panelEl.style.minWidth = matchTriggerWidth ? `${triggerRect.width}px` : ''
    panelEl.style.width = matchTriggerWidth ? `${triggerRect.width}px` : ''
    const panelRect = panelEl.getBoundingClientRect()
    const viewportWidth = window.innerWidth
    const viewportHeight = window.innerHeight
    const spaceBottom = viewportHeight - triggerRect.bottom - viewportPadding
    const spaceTop = triggerRect.top - viewportPadding
    const spaceRight = viewportWidth - triggerRect.right - viewportPadding
    const spaceLeft = triggerRect.left - viewportPadding
    let nextSide: FloatingPanelSide = side
    let nextAlign: FloatingPanelAlign = align
    if (nextSide === 'bottom' && panelRect.height > spaceBottom && spaceTop > spaceBottom) {
      nextSide = 'top'
    } else if (nextSide === 'top' && panelRect.height > spaceTop && spaceBottom > spaceTop) {
      nextSide = 'bottom'
    } else if (nextSide === 'right' && panelRect.width > spaceRight && spaceLeft > spaceRight) {
      nextSide = 'left'
    } else if (nextSide === 'left' && panelRect.width > spaceLeft && spaceRight > spaceLeft) {
      nextSide = 'right'
    }
    if (
      (nextSide === 'top' || nextSide === 'bottom') &&
      nextAlign === 'right' &&
      panelRect.width > triggerRect.right - viewportPadding &&
      spaceRight > spaceLeft
    ) {
      nextAlign = 'left'
    } else if (
      (nextSide === 'top' || nextSide === 'bottom') &&
      nextAlign === 'left' &&
      panelRect.width > viewportWidth - triggerRect.left - viewportPadding &&
      spaceLeft > spaceRight
    ) {
      nextAlign = 'right'
    }
    let left = 0
    let top = 0
    if (nextSide === 'top' || nextSide === 'bottom') {
      left =
        nextAlign === 'left' ? triggerRect.left + crossAxisOffset : triggerRect.right - panelRect.width - crossAxisOffset
      top = nextSide === 'bottom' ? triggerRect.bottom + offset : triggerRect.top - panelRect.height - offset
    } else {
      left = nextSide === 'right' ? triggerRect.right + offset : triggerRect.left - panelRect.width - offset
      top = triggerRect.top + triggerRect.height / 2 - panelRect.height / 2
    }
    left = clamp(left, viewportPadding, Math.max(viewportPadding, viewportWidth - panelRect.width - viewportPadding))
    top = clamp(top, viewportPadding, Math.max(viewportPadding, viewportHeight - panelRect.height - viewportPadding))
    resolvedSide.value = nextSide
    resolvedAlign.value = nextAlign
    panelStyle.value = {
      position: 'fixed',
      left: `${Math.round(left)}px`,
      top: `${Math.round(top)}px`,
      minWidth: matchTriggerWidth ? `${Math.round(triggerRect.width)}px` : '',
      width: matchTriggerWidth ? `${Math.round(triggerRect.width)}px` : '',
      maxWidth: `calc(100vw - ${viewportPadding * 2}px)`,
      maxHeight: `calc(100vh - ${viewportPadding * 2}px)`
    }
    panelEl.style.visibility = ''
  }
  /** Queues a position update for the next animation frame. */
  const schedulePositionUpdate = () => {
    if (!import.meta.client) return
    cancelScheduledUpdate()
    rafId.value = requestAnimationFrame(() => {
      void updatePosition()
      rafId.value = null
    })
  }
  onBeforeUnmount(cancelScheduledUpdate)
  return {
    panelStyle,
    resolvedSide,
    resolvedAlign,
    updatePosition,
    schedulePositionUpdate,
    cancelScheduledUpdate,
    resetPosition
  }
}
