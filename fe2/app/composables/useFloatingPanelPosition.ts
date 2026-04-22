export type Side = "top" | "bottom" | "left" | "right"

interface Options {
  triggerRef: Ref<HTMLElement | null>
  panelRef: Ref<HTMLElement | null>
  side?: () => Side
  offset?: number
}

export function useFloatingPanelPosition({triggerRef, panelRef, side = () => "top", offset = 8}: Options) {
  const panelStyle = ref<Record<string, string>>({})

  const clamp = (v: number, max: number) => Math.max(8, Math.min(v, max - 8))

  const updatePosition = () => {
    const trigger = triggerRef.value
    const panel = panelRef.value
    if (!trigger || !panel) return

    const t = trigger.getBoundingClientRect()
    const p = panel.getBoundingClientRect()

    const s = side()

    const isVertical = s === "top" || s === "bottom"

    const top = isVertical
      ? s === "top"
        ? t.top - p.height - offset
        : t.bottom + offset
      : t.top + (t.height - p.height) / 2

    const left = !isVertical
      ? s === "left"
        ? t.left - p.width - offset
        : t.right + offset
      : t.left + (t.width - p.width) / 2

    panelStyle.value = {
      position: "fixed",
      top: `${clamp(top, window.innerHeight - p.height)}px`,
      left: `${clamp(left, window.innerWidth - p.width)}px`
    }
  }

  const handler = () => updatePosition()

  onMounted(() => {
    window.addEventListener("scroll", handler, true)
    window.addEventListener("resize", handler)
  })

  onBeforeUnmount(() => {
    window.removeEventListener("scroll", handler, true)
    window.removeEventListener("resize", handler)
  })

  return {panelStyle, updatePosition}
}
