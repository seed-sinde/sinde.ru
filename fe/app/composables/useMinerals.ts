type MineralsListParams = {
  q?: string
  limit?: number
  offset?: number
  sort?: 'name_asc' | 'name_desc'
  imageFilter?: MineralImageFilter
  crystalSystems?: MineralCrystalSystem[]
  crystalSystemMode?: MineralCrystalSystemMode
  chemistryAll?: string[]
  chemistryAny?: string[]
  chemistryNone?: string[]
}
const useApiJson = <T>(path: string, options?: NonNullable<Parameters<ReturnType<typeof useAPI>['json']>[1]>) =>
  useAPI().json<T>(path, options)
const appendMineralsArrayParam = (query: URLSearchParams, key: string, values?: string[]) => {
  const normalized = (values || []).map(value => String(value || '').trim()).filter(Boolean)
  if (!normalized.length) return
  query.set(key, normalized.join(','))
}
export const getMineralsList = async (params: MineralsListParams = {}) => {
  const query = new URLSearchParams()
  if (String(params.q || '').trim()) query.set('q', String(params.q || '').trim())
  if (typeof params.limit === 'number' && Number.isFinite(params.limit)) query.set('limit', String(params.limit))
  if (typeof params.offset === 'number' && Number.isFinite(params.offset)) query.set('offset', String(params.offset))
  if (params.sort) query.set('sort', params.sort)
  if (params.imageFilter === 'with' || params.imageFilter === 'without') query.set('image', params.imageFilter)
  appendMineralsArrayParam(query, 'crystalSystem', params.crystalSystems)
  if (params.crystalSystemMode === 'all') query.set('crystalSystemMode', 'all')
  appendMineralsArrayParam(query, 'chemistryAll', params.chemistryAll)
  appendMineralsArrayParam(query, 'chemistryAny', params.chemistryAny)
  appendMineralsArrayParam(query, 'chemistryNone', params.chemistryNone)
  const suffix = query.toString()
  return await useApiJson<ApiResponseWithData<MineralsListPayload>>(`/minerals${suffix ? `?${suffix}` : ''}`, {
    method: 'GET'
  })
}
export const getMineralByDatabaseId = async (databaseId: string | number) => {
  return await useApiJson<ApiResponseWithData<Mineral>>(
    `/minerals/${encodeURIComponent(String(databaseId || '').trim())}`,
    {
      method: 'GET'
    }
  )
}
