import { joinURL } from 'ufo'
const useApiJson = <T>(path: string, options?: NonNullable<Parameters<ReturnType<typeof useAPI>['json']>[1]>) =>
  useAPI().json<T>(path, options)
/**
 * Creates a simple trait payload for set-building requests.
 */
const toTraitPairPayload = (trait: Pick<TraitInput, 't_key' | 't_value'>) => ({
  t_key: trait.t_key,
  t_value: String(trait.t_value ?? '')
})
/**
 * Returns true when the payload looks like a serialized trait record.
 */
const isTraitRecord = (value: unknown): value is Trait => {
  return Boolean(
    value &&
    typeof value === 'object' &&
    typeof (value as Trait).t_uuid === 'string' &&
    typeof (value as Trait).t_key === 'string' &&
    typeof (value as Trait).t_value === 'string'
  )
}
/**
 * Returns true when the payload looks like a compact trait tuple.
 */
const isTraitTuple = (value: unknown): value is [string, string] => {
  return Array.isArray(value) && value.length === 2 && typeof value[0] === 'string' && typeof value[1] === 'string'
}
/**
 * Builds a trait record from a compact tuple response.
 */
const createTraitFromTuple = (uuid: string, value: [string, string]): Trait => {
  return { t_uuid: uuid, t_key: value[0], t_value: value[1] }
}
/**
 * Parses one NDJSON payload into a normalized trait list.
 */
const parseNdjsonTraits = (source: string): Trait[] => {
  const items: Trait[] = []
  for (const rawLine of source.split('\n')) {
    const line = rawLine.trim()
    if (!line) continue
    try {
      const parsed = JSON.parse(line)
      if (isTraitRecord(parsed)) {
        items.push(parsed)
      }
    } catch {
      // ignore malformed lines
    }
  }
  return items
}
/**
 * Wraps an HTTP status into a regular Error instance.
 */
const toErrorWithStatus = (message: string, statusCode: number) => {
  const error = new Error(message) as Error & { statusCode?: number; status?: number }
  error.statusCode = statusCode
  error.status = statusCode
  return error
}
/**
 * Normalizes one JSON traits resolve response into a list of records.
 */
const parseTraitResolveJsonPayload = (uuid: string, payload: unknown): Trait[] => {
  if (isTraitTuple(payload)) {
    return [createTraitFromTuple(uuid, payload)]
  }
  if (isTraitRecord(payload)) {
    return [payload]
  }
  if (Array.isArray(payload)) {
    return payload.filter(isTraitRecord)
  }
  return []
}
/**
 * Normalizes one single-trait payload into a concrete record.
 */
const normalizeTraitPayload = (uuid: string, payload: unknown): Trait | null => {
  if (isTraitTuple(payload)) {
    return createTraitFromTuple(uuid, payload)
  }
  return isTraitRecord(payload) ? payload : null
}
/**
 * Returns proxy request headers required for SSR trait resolve fetches.
 */
const getTraitResolveHeaders = (): Record<string, string> => {
  const headers: Record<string, string> = {
    Accept: 'application/x-ndjson, application/json'
  }
  if (import.meta.server) {
    const forwardedHeaderNames = ['cookie', 'user-agent', 'accept-language', 'x-forwarded-for', 'x-real-ip'] as const
    const requestHeaders = useRequestHeaders(forwardedHeaderNames as unknown as string[])
    for (const name of forwardedHeaderNames) {
      const value = requestHeaders[name]
      if (value) headers[name] = value
    }
  }
  return headers
}
/**
 * Resolves the absolute proxy base URL for trait requests during SSR and CSR.
 */
const getTraitResolveBaseUrl = (): string => {
  const config = useRuntimeConfig()
  const base = import.meta.client ? '' : useRequestURL().origin || config.public.baseURL
  if (!base) {
    throw new Error('SSR base URL is not defined')
  }
  return base
}
/**
 * Executes a POST request against the keys API.
 */
const postTraitKeyRequest = async <T>(path: string, body: Record<string, unknown>) => {
  return await useApiJson<T>(path, {
    method: 'POST',
    body
  })
}
/**
 * Creates a new standalone trait and returns its uuid.
 */
const createTrait = (t_key: string, t_value: string) =>
  useApiJson<string>('/traits', {
    method: 'POST',
    body: { t_key, t_value }
  })
/**
 * Creates or resolves a set uuid by appending a trait to the current left-side uuid.
 */
const createOrGetSet = (leftUuid: string, right: Pick<TraitInput, 't_key' | 't_value'>) =>
  useApiJson<ApiResponseWithData<{ s_uuid: string; t_uuid?: string }>>('/sets', {
    method: 'POST',
    body: {
      s_childs: [leftUuid, toTraitPairPayload(right)]
    }
  }).then(response => response?.data)
/**
 * Finds the set uuid that contains the provided trait uuid collection.
 */
const findSetByTraitUuids = (t_uuids: string[]) =>
  useApiJson<ApiResponseWithData<{ s_uuid: string }>>('/sets/find', {
    method: 'POST',
    body: t_uuids
  }).then(response => response?.data)
/**
 * Fetches key meta by synonym.
 */
const getKeyMeta = async (syn: string) => {
  return await postTraitKeyRequest<ApiResponseWithData<TraitKey>>('/keys/meta', { syn })
}
/**
 * Fetches key meta for several ids in one request.
 */
const getKeysMetaBulk = async (ids: number[]) => {
  return await postTraitKeyRequest<ApiResponseWithData<{ items: TraitKey[] }>>('/keys/meta/bulk', { ids })
}
/**
 * Fetches enum options for one key synonym.
 */
const getEnumOptions = async (syn: string) => {
  return await postTraitKeyRequest<ApiResponseWithData<{ options: string[] }>>('/keys/enum/options', { syn })
}
/**
 * Updates the stored meta definition for one key id.
 */
const updateKeyMeta = async (id: number, meta: KeyMeta) => {
  return await postTraitKeyRequest<ApiResponseWithData<TraitKey>>('/keys/meta/update', { id, meta })
}
/**
 * Loads all traits behind one uuid, supporting both NDJSON sets and JSON single traits.
 */
const fetchTraitsByUuid = async (uuid: string): Promise<Trait[]> => {
  if (!uuid) return []
  const url = joinURL(getTraitResolveBaseUrl(), '/api/proxy/traits/resolve', uuid)
  const response = await fetch(url, {
    method: 'GET',
    credentials: 'include',
    headers: getTraitResolveHeaders()
  })
  if (!response.ok) {
    if (response.status === 404) return []
    throw toErrorWithStatus(`Failed to resolve traits: HTTP ${response.status}`, response.status)
  }
  const contentType = String(response.headers.get('content-type') || '').toLowerCase()
  const body = await response.text()
  if (!body.trim()) return []
  if (contentType.includes('application/x-ndjson')) {
    return parseNdjsonTraits(body)
  }
  if (contentType.includes('application/json')) {
    try {
      return parseTraitResolveJsonPayload(uuid, JSON.parse(body) as TraitResolvePayload)
    } catch {
      return []
    }
  }
  return []
}
/**
 * Loads one trait by uuid and falls back to the resolve endpoint when needed.
 */
const fetchTraitByUuid = async (uuid: string): Promise<Trait | null> => {
  if (!uuid) return null
  try {
    const payload = await useApiJson<[string, string] | Trait>(`/traits/${uuid}`)
    return normalizeTraitPayload(uuid, payload)
  } catch (error: any) {
    const status = Number(error?.statusCode ?? error?.status ?? error?.response?.status ?? 0)
    if (status === 404) return null
  }
  const fallbackItems = await fetchTraitsByUuid(uuid)
  return fallbackItems[0] || null
}
/**
 * Shortens a uuid for compact UI display.
 */
const shortUuid = (value?: string, keep = 8): string => {
  if (!value) return '—'
  return value.length <= keep * 2 ? value : `${value.slice(0, keep)}…${value.slice(-keep)}`
}
export {
  createOrGetSet,
  createTrait,
  fetchTraitByUuid,
  fetchTraitsByUuid,
  findSetByTraitUuids,
  getEnumOptions,
  getKeyMeta,
  getKeysMetaBulk,
  shortUuid,
  updateKeyMeta
}
