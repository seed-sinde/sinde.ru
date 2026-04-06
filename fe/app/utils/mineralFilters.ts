export const parseQueryStringValue = (value: unknown) => {
  if (Array.isArray(value)) return String(value[0] || '').trim()
  return String(value || '').trim()
}

export const normalizeMineralElementSymbol = (value: unknown) => formatPeriodicTableSymbol(value)

export const uniqueElementSymbols = (values: string[]) =>
  Array.from(new Set(values.map(normalizeMineralElementSymbol).filter(Boolean)))

export const sortMineralElementSymbols = (values: string[], compareElementOrder: (left: string, right: string) => number) =>
  uniqueElementSymbols(values).sort(compareElementOrder)

export const parseQueryElementList = (value: unknown, compareElementOrder: (left: string, right: string) => number) =>
  sortMineralElementSymbols(
    parseQueryStringValue(value)
      .split(',')
      .map(item => normalizeMineralElementSymbol(item))
      .filter(Boolean),
    compareElementOrder
  )

export const parseQueryCrystalSystems = (value: unknown, allowedCrystalSystems: Set<MineralCrystalSystem>) =>
  parseQueryStringValue(value)
    .split(',')
    .map(item => item.trim() as MineralCrystalSystem)
    .filter(item => allowedCrystalSystems.has(item))

export const parseQueryCrystalSystemMode = (value: unknown): MineralCrystalSystemMode =>
  parseQueryStringValue(value) === 'all' ? 'all' : 'any'

export const parseQueryLimit = (value: unknown, defaultLimit: number) => {
  const parsed = Number.parseInt(parseQueryStringValue(value), 10)
  if (!Number.isFinite(parsed) || parsed <= 0) return defaultLimit
  return Math.min(parsed, 100)
}

export const parseQueryOffset = (value: unknown) => {
  const parsed = Number.parseInt(parseQueryStringValue(value), 10)
  if (!Number.isFinite(parsed) || parsed < 0) return 0
  return parsed
}

export const parseQuerySort = (value: unknown): MineralsRouteState['sort'] => {
  const parsed = parseQueryStringValue(value)
  return parsed === 'name_desc' ? 'name_desc' : 'name_asc'
}

export const parseQueryOnlyWithImages = (value: unknown) => {
  const parsed = parseQueryStringValue(value).toLowerCase()
  return parsed === '1' || parsed === 'true' || parsed === 'yes' || parsed === 'on'
}

export const hasMineralsRouteQueryValues = (query: Record<string, unknown>) =>
  Object.values(query).some(value => {
    if (Array.isArray(value)) return value.some(item => String(item || '').trim())
    return String(value || '').trim() !== ''
  })

export const readMineralsRouteState = (
  query: Record<string, unknown>,
  options: {
    defaultLimit: number
    allowedCrystalSystems: Set<MineralCrystalSystem>
    compareElementOrder: (left: string, right: string) => number
  }
): MineralsRouteState => ({
  q: parseQueryStringValue(query.q),
  sort: parseQuerySort(query.sort),
  limit: parseQueryLimit(query.limit, options.defaultLimit),
  offset: parseQueryOffset(query.offset),
  onlyWithImages: parseQueryOnlyWithImages(query.onlyWithImages),
  crystalSystems: parseQueryCrystalSystems(query.crystalSystem, options.allowedCrystalSystems),
  crystalSystemMode: parseQueryCrystalSystemMode(query.crystalSystemMode),
  chemistryAll: parseQueryElementList(query.chemistryAll, options.compareElementOrder),
  chemistryAny: parseQueryElementList(query.chemistryAny, options.compareElementOrder),
  chemistryNone: parseQueryElementList(query.chemistryNone, options.compareElementOrder)
})

export const buildMineralsRouteQuery = (
  state: MineralsRouteState,
  options: {
    defaultLimit: number
    sortElements: (values: string[]) => string[]
  }
) => {
  const query: Record<string, string> = {}
  const search = state.q.trim()
  if (search) query.q = search
  if (state.sort !== 'name_asc') query.sort = state.sort
  if (state.limit !== options.defaultLimit) query.limit = String(state.limit)
  if (state.offset > 0) query.offset = String(state.offset)
  if (state.onlyWithImages) query.onlyWithImages = '1'
  if (state.crystalSystems.length) query.crystalSystem = state.crystalSystems.join(',')
  if (state.crystalSystemMode !== 'any') query.crystalSystemMode = state.crystalSystemMode
  if (state.chemistryAll.length) query.chemistryAll = options.sortElements(state.chemistryAll).join(',')
  if (state.chemistryAny.length) query.chemistryAny = options.sortElements(state.chemistryAny).join(',')
  if (state.chemistryNone.length) query.chemistryNone = options.sortElements(state.chemistryNone).join(',')
  return query
}
