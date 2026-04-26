export type ScrollableEdgesAxis = "x" | "y" | "both"
export type ScrollableEdgesState = {left: boolean; right: boolean; top: boolean; bottom: boolean}
interface Options {
  axis?: ScrollableEdgesAxis
  threshold?: number
}
type Target = Ref<HTMLElement | null>
type ScrollableEdgesEntry = {edges: ScrollableEdgesState; sync: () => void}
type ScrollableEdgesConfig = Options & {target: Target}
type ScrollableEdgesConfigs = Record<string, ScrollableEdgesConfig>
type ScrollableEdgesItems<T extends ScrollableEdgesConfigs> = {[K in keyof T]: ScrollableEdgesEntry}
const blank = () =>
  reactive<ScrollableEdgesState>({left: false, right: false, top: false, bottom: false})
const reset = (edges: ScrollableEdgesState) => {
  edges.left = false
  edges.right = false
  edges.top = false
  edges.bottom = false
}
const read = (el: HTMLElement | null, edges: ScrollableEdgesState, options: Options) => {
  if (!el) return reset(edges)
  const axis = options.axis ?? "both"
  const threshold = options.threshold ?? 2
  const x = axis === "x" || axis === "both"
  const y = axis === "y" || axis === "both"
  edges.left = x && el.scrollLeft > threshold
  edges.right = x && el.scrollLeft < Math.max(0, el.scrollWidth - el.clientWidth) - threshold
  edges.top = y && el.scrollTop > threshold
  edges.bottom = y && el.scrollTop < Math.max(0, el.scrollHeight - el.clientHeight) - threshold
}
export function useScrollableEdges(target: Target, options?: Options): ScrollableEdgesEntry
export function useScrollableEdges<const T extends ScrollableEdgesConfigs>(
  targets: T
): {items: ScrollableEdgesItems<T>; sync: () => void}
export function useScrollableEdges(target: Target | ScrollableEdgesConfigs, options: Options = {}) {
  const configs = isRef(target) ? {single: {target, ...options}} : target
  const items = Object.fromEntries(
    Object.keys(configs).map(key => [key, {edges: blank(), sync: () => {}}])
  ) as ScrollableEdgesItems<ScrollableEdgesConfigs>
  Object.entries(configs).forEach(([key, config]) => {
    items[key]!.sync = () => read(config.target.value, items[key]!.edges, config)
  })
  const sync = () =>
    Object.entries(configs).forEach(([key, config]) =>
      read(config.target.value, items[key]!.edges, config)
    )
  let resizeObserver: ResizeObserver | null = null
  let stops: Array<() => void> = []
  const detach = () => {
    stops.forEach(stop => stop())
    stops = []
    resizeObserver?.disconnect()
    resizeObserver = null
  }
  const attach = () => {
    detach()
    if (!import.meta.client) return sync()
    resizeObserver = new ResizeObserver(sync)
    for (const config of Object.values(configs)) {
      const el = config.target.value
      if (!el) continue
      el.addEventListener("scroll", sync, {passive: true})
      stops.push(() => el.removeEventListener("scroll", sync))
      resizeObserver.observe(el)
    }
    requestAnimationFrame(sync)
  }
  watch(() => Object.values(configs).map(config => config.target.value), attach, {flush: "post"})
  onMounted(attach)
  onBeforeUnmount(detach)
  return isRef(target) ? items.single! : {items, sync}
}
