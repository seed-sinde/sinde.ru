export type ScrollableEdgesAxis = 'x' | 'y' | 'both'

export const useScrollableEdges = (
  target: Ref<HTMLElement | null>,
  options?: {
    axis?: ScrollableEdgesAxis
    threshold?: number
  }
) => {
  const edges = reactive({
    left: false,
    right: false,
    top: false,
    bottom: false
  })

  const axis = options?.axis || 'both'
  const threshold = options?.threshold ?? 2
  let resizeObserver: ResizeObserver | null = null
  let removeScrollListener: (() => void) | null = null

  const sync = () => {
    const element = target.value
    if (!element) {
      edges.left = false
      edges.right = false
      edges.top = false
      edges.bottom = false
      return
    }
    if (axis === 'x' || axis === 'both') {
      const maxScrollLeft = Math.max(0, element.scrollWidth - element.clientWidth)
      edges.left = element.scrollLeft > threshold
      edges.right = element.scrollLeft < maxScrollLeft - threshold
    } else {
      edges.left = false
      edges.right = false
    }
    if (axis === 'y' || axis === 'both') {
      const maxScrollTop = Math.max(0, element.scrollHeight - element.clientHeight)
      edges.top = element.scrollTop > threshold
      edges.bottom = element.scrollTop < maxScrollTop - threshold
    } else {
      edges.top = false
      edges.bottom = false
    }
  }

  const detach = () => {
    removeScrollListener?.()
    removeScrollListener = null
    resizeObserver?.disconnect()
    resizeObserver = null
  }

  const attach = (element: HTMLElement | null) => {
    detach()
    if (!import.meta.client || !element) {
      sync()
      return
    }
    const onScroll = () => sync()
    element.addEventListener('scroll', onScroll, { passive: true })
    removeScrollListener = () => {
      element.removeEventListener('scroll', onScroll)
    }
    resizeObserver = new ResizeObserver(() => {
      sync()
    })
    resizeObserver.observe(element)
    requestAnimationFrame(sync)
  }

  watch(
    target,
    (element) => {
      attach(element)
    },
    { flush: 'post' }
  )

  onMounted(() => {
    attach(target.value)
  })

  onBeforeUnmount(() => {
    detach()
  })

  return {
    edges,
    sync
  }
}
