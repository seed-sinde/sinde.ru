export type Side = "top" | "bottom" | "left" | "right"

interface Options {
  triggerRef: Ref<HTMLElement | null>
  panelRef: Ref<HTMLElement | null>
  side?: () => Side
  offset?: number
}

export function useFloatingPanelPosition({triggerRef, panelRef, side = () => "top", offset = 8}: Options) {
  const panelStyle = ref<Record<string, string>>({})

  const updatePosition = () => {
    const trigger = triggerRef.value
    const panel = panelRef.value
    if (!trigger || !panel) return

    const t = trigger.getBoundingClientRect()
    const p = panel.getBoundingClientRect()

    let top = 0
    let left = 0

    switch (side()) {
      case "top":
        top = t.top - p.height - offset
        left = t.left + (t.width - p.width) / 2
        break
      case "bottom":
        top = t.bottom + offset
        left = t.left + (t.width - p.width) / 2
        break
      case "left":
        top = t.top + (t.height - p.height) / 2
        left = t.left - p.width - offset
        break
      case "right":
        top = t.top + (t.height - p.height) / 2
        left = t.right + offset
        break
    }

    const vw = window.innerWidth
    left = Math.max(8, Math.min(left, vw - p.width - 8))

    const vh = window.innerHeight
    top = Math.max(8, Math.min(top, vh - p.height - 8))

    panelStyle.value = {
      position: "fixed",
      top: `${top}px`,
      left: `${left}px`
    }
  }

  const onScrollOrResize = () => updatePosition()

  onMounted(() => {
    window.addEventListener("scroll", onScrollOrResize, true)
    window.addEventListener("resize", onScrollOrResize)
  })

  onBeforeUnmount(() => {
    window.removeEventListener("scroll", onScrollOrResize, true)
    window.removeEventListener("resize", onScrollOrResize)
  })

  return {
    panelStyle,
    updatePosition
  }
}
