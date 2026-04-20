import {nextTick, onBeforeUnmount, ref, toValue, type MaybeRefOrGetter, type Ref} from "vue"

export type FloatingPanelSide = "top" | "bottom" | "left" | "right"

export const useFloatingPanelPosition = (options: {
  triggerRef: Ref<HTMLElement | null>
  panelRef: Ref<HTMLElement | null>
  side?: MaybeRefOrGetter<FloatingPanelSide>
  offset?: MaybeRefOrGetter<number>
  viewportPadding?: MaybeRefOrGetter<number>
  matchTriggerWidth?: MaybeRefOrGetter<boolean>
}) => {
  const panelStyle = ref<Record<string, string>>({})
  const resolvedSide = ref<FloatingPanelSide>(toValue(options.side) || "bottom")
  const rafId = ref<number | null>(null)

  const clamp = (value: number, min: number, max: number) => Math.min(Math.max(value, min), max)

  const cancelScheduledUpdate = () => {
    if (!import.meta.client || rafId.value === null) return
    cancelAnimationFrame(rafId.value)
    rafId.value = null
  }

  const resetPosition = () => {
    panelStyle.value = {}
    resolvedSide.value = toValue(options.side) || "bottom"
    cancelScheduledUpdate()
  }

  const updatePosition = async () => {
    if (!import.meta.client) return

    await nextTick()

    const triggerEl = options.triggerRef.value
    const panelEl = options.panelRef.value
    if (!triggerEl || !panelEl) return

    const side = toValue(options.side) || "bottom"
    const offset = toValue(options.offset) ?? 8
    const viewportPadding = toValue(options.viewportPadding) ?? 12
    const matchTriggerWidth = Boolean(toValue(options.matchTriggerWidth))
    const triggerRect = triggerEl.getBoundingClientRect()

    panelEl.style.visibility = "hidden"
    panelEl.style.left = "0px"
    panelEl.style.top = "0px"
    panelEl.style.minWidth = matchTriggerWidth ? `${triggerRect.width}px` : ""
    panelEl.style.width = matchTriggerWidth ? `${triggerRect.width}px` : ""

    const panelRect = panelEl.getBoundingClientRect()
    const viewportWidth = window.innerWidth
    const viewportHeight = window.innerHeight

    const spaceTop = triggerRect.top - viewportPadding
    const spaceBottom = viewportHeight - triggerRect.bottom - viewportPadding
    const spaceLeft = triggerRect.left - viewportPadding
    const spaceRight = viewportWidth - triggerRect.right - viewportPadding

    const preferredAxis = side === "left" || side === "right" ? "x" : "y"

    const nextSide =
      preferredAxis === "y"
        ? side === "top"
          ? panelRect.height <= spaceTop || spaceTop > spaceBottom
            ? "top"
            : "bottom"
          : panelRect.height <= spaceBottom || spaceBottom > spaceTop
            ? "bottom"
            : "top"
        : side === "left"
          ? panelRect.width <= spaceLeft || spaceLeft > spaceRight
            ? "left"
            : "right"
          : panelRect.width <= spaceRight || spaceRight > spaceLeft
            ? "right"
            : "left"

    const left =
      nextSide === "top" || nextSide === "bottom"
        ? clamp(
            triggerRect.left + triggerRect.width / 2 - panelRect.width / 2,
            viewportPadding,
            Math.max(viewportPadding, viewportWidth - panelRect.width - viewportPadding)
          )
        : nextSide === "right"
          ? triggerRect.right + offset
          : triggerRect.left - panelRect.width - offset

    const top =
      nextSide === "left" || nextSide === "right"
        ? clamp(
            triggerRect.top + triggerRect.height / 2 - panelRect.height / 2,
            viewportPadding,
            Math.max(viewportPadding, viewportHeight - panelRect.height - viewportPadding)
          )
        : nextSide === "bottom"
          ? triggerRect.bottom + offset
          : triggerRect.top - panelRect.height - offset

    resolvedSide.value = nextSide
    panelStyle.value = {
      position: "fixed",
      left: `${Math.round(clamp(left, viewportPadding, Math.max(viewportPadding, viewportWidth - panelRect.width - viewportPadding)))}px`,
      top: `${Math.round(clamp(top, viewportPadding, Math.max(viewportPadding, viewportHeight - panelRect.height - viewportPadding)))}px`,
      minWidth: matchTriggerWidth ? `${Math.round(triggerRect.width)}px` : "",
      width: matchTriggerWidth ? `${Math.round(triggerRect.width)}px` : "",
      maxHeight: `calc(100vh - ${viewportPadding * 2}px)`
    }

    panelEl.style.visibility = ""
  }

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
    updatePosition,
    schedulePositionUpdate,
    resetPosition
  }
}
