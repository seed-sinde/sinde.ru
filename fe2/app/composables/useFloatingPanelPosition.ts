export type Side = "top" | "bottom" | "left" | "right"
interface Options {
  triggerRef: Ref<HTMLElement | null>
  panelRef: Ref<HTMLElement | null>
  side?: () => Side
  offset?: number
  viewportPadding?: number
}
export const useFloatingPanelPosition = ({
  triggerRef,
  panelRef,
  side = () => "top",
  offset = 14,
  viewportPadding = 14
}: Options) => {
  const actualSide = ref<Side>(side())
  const panelStyle = ref<Record<string, string>>({})
  const clamp = (v: number, min: number, max: number) => Math.max(min, Math.min(v, max))
  const updatePosition = () => {
    const t = triggerRef.value?.getBoundingClientRect()
    const p = panelRef.value?.getBoundingClientRect()
    if (!t || !p) return
    let s = side()
    s =
      s === "top" && t.top < p.height + offset
        ? "bottom"
        : s === "bottom" && window.innerHeight - t.bottom < p.height + offset
          ? "top"
          : s === "left" && t.left < p.width + offset
            ? "right"
            : s === "right" && window.innerWidth - t.right < p.width + offset
              ? "left"
              : s
    actualSide.value = s
    const top =
      s === "top"
        ? t.top - p.height - offset
        : s === "bottom"
          ? t.bottom + offset
          : t.top + t.height / 2 - p.height / 2
    const left =
      s === "left"
        ? t.left - p.width - offset
        : s === "right"
          ? t.right + offset
          : t.left + t.width / 2 - p.width / 2
    panelStyle.value = {
      position: "fixed",
      left: `${clamp(left, viewportPadding, window.innerWidth - p.width - viewportPadding)}px`,
      top: `${clamp(top, viewportPadding, window.innerHeight - p.height - viewportPadding)}px`
    }
  }
  const handler = () => updatePosition()
  onMounted(() => {
    window.addEventListener("scroll", handler, true)
    window.addEventListener("resize", handler)
    requestAnimationFrame(updatePosition)
  })
  useLayoutUpdate(updatePosition)
  onBeforeUnmount(() => {
    window.removeEventListener("scroll", handler, true)
    window.removeEventListener("resize", handler)
  })
  return {actualSide, panelStyle, updatePosition}
}
export const useEventListener = (
  target: EventTarget,
  event: string,
  handler: EventListenerOrEventListenerObject
) => {
  onMounted(() => target.addEventListener(event, handler))
  onBeforeUnmount(() => target.removeEventListener(event, handler))
}
export const useLayoutUpdate = (fn: () => void) => useEventListener(window, "layout:update", fn)
export const useLayoutUpdateTrigger = () => () =>
  requestAnimationFrame(() => window.dispatchEvent(new Event("layout:update")))
