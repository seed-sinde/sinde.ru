export function useFloatingArrow({
  triggerRef,
  panelRef,
  actualSide,
  arrowSize = 4,
  color = "var(--tooltip-bg)"
}: {
  triggerRef: Ref<HTMLElement | null>
  panelRef: Ref<HTMLElement | null>
  actualSide: Ref<Side>
  arrowSize?: number
  color?: string
}) {
  const arrowStyle = ref<Record<string, string>>({})
  const clamp = (v: number, min: number, max: number) => Math.max(min, Math.min(v, max))
  const sideMap = {
    top: "bottom-0 left-0",
    bottom: "top-0 left-0",
    left: "right-0 top-0",
    right: "left-0 top-0"
  } as const
  const arrowClass = computed(
    () => `absolute h-0 w-0 border-transparent ${sideMap[actualSide.value]}`
  )
  const arrowTransform = computed(() =>
    actualSide.value === "top"
      ? "translate(-50%, 100%)"
      : actualSide.value === "bottom"
        ? "translate(-50%, -100%)"
        : actualSide.value === "left"
          ? "translate(100%, -50%)"
          : "translate(-100%, -50%)"
  )
  const arrowColorStyle = computed(() => {
    return actualSide.value === "top"
      ? {borderTopColor: color}
      : actualSide.value === "bottom"
        ? {borderBottomColor: color}
        : actualSide.value === "left"
          ? {borderLeftColor: color}
          : {borderRightColor: color}
  })
  const updateArrow = () => {
    const t = triggerRef.value?.getBoundingClientRect()
    const p = panelRef.value?.getBoundingClientRect()
    if (!t || !p) return
    const isVertical = actualSide.value === "top" || actualSide.value === "bottom"
    const raw = isVertical ? t.left + t.width / 2 - p.left : t.top + t.height / 2 - p.top
    const value = clamp(
      raw,
      arrowSize * 2,
      isVertical ? p.width - arrowSize * 2 : p.height - arrowSize * 2
    )
    arrowStyle.value = isVertical ? {left: `${value}px`} : {top: `${value}px`}
  }
  const arrowFullStyle = computed(() => ({
    ...arrowStyle.value,
    transform: arrowTransform.value,
    borderWidth: `${arrowSize}px`
  }))
  return {
    arrowClass,
    arrowFullStyle,
    arrowColorStyle,
    updateArrow
  }
}
