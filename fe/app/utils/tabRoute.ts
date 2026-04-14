import type { RouteLocationNormalizedLoaded, RouteLocationRaw } from 'vue-router'
export type TabRouteValue = string | number
export type TabRouteTargetMap = Record<string, RouteLocationRaw | undefined>
export type BuildTabRouteOptions = {
  queryKey?: string
  defaultValue?: TabRouteValue
  path?: string
  persistDefault?: boolean
  targetMap?: TabRouteTargetMap
}
const readRouteValue = (raw: unknown): string => {
  if (typeof raw === 'string') return raw
  if (Array.isArray(raw)) return typeof raw[0] === 'string' ? raw[0] : ''
  return ''
}
export const stringifyTabValue = (value: TabRouteValue): string => String(value)
export const normalizeTabRouteValue = (
  raw: unknown,
  allowedValues: TabRouteValue[],
  fallback: TabRouteValue
): TabRouteValue => {
  const candidate = readRouteValue(raw)
  const match = allowedValues.find((value) => stringifyTabValue(value) === candidate)
  return match ?? fallback
}
export const buildTabRouteLocation = (
  route: RouteLocationNormalizedLoaded,
  value: TabRouteValue,
  options: BuildTabRouteOptions
): RouteLocationRaw | undefined => {
  const valueKey = stringifyTabValue(value)
  const mapped = options.targetMap?.[valueKey]
  if (mapped) return mapped
  if (!options.queryKey) return undefined
  const query = { ...route.query }
  const defaultKey = stringifyTabValue(options.defaultValue ?? value)
  if (!options.persistDefault && valueKey === defaultKey) {
    const { [options.queryKey]: _removed, ...restQuery } = query
    return {
      path: options.path || route.path,
      query: restQuery,
      hash: route.hash
    }
  }
  query[options.queryKey] = valueKey
  return {
    path: options.path || route.path,
    query,
    hash: route.hash
  }
}
