const useApiJson = <T>(path: string, options?: NonNullable<Parameters<ReturnType<typeof useAPI>["json"]>[1]>) =>
  useAPI().json<T>(path, options)
type ApiErrorLike = {
  status?: number
  statusCode?: number
  response?: {status?: number}
}
/**
 * Creates a simple trait payload for set-building requests.
 */
const toTraitPairPayload = (trait: Pick<TraitInput, "t_key" | "t_value" | "meta">) => ({
  t_key: trait.t_key,
  t_value: String(trait.t_value ?? ""),
  ...(trait.meta ? {meta: trait.meta} : {})
})
/**
 * Returns true when the payload looks like a serialized trait record.
 */
const isTraitRecord = (value: unknown): value is Trait => {
  return Boolean(
    value &&
    typeof value === "object" &&
    typeof (value as Trait).t_uuid === "string" &&
    typeof (value as Trait).t_key === "string" &&
    typeof (value as Trait).t_value === "string"
  )
}
/**
 * Returns true when the payload looks like a compact trait tuple.
 */
const isTraitTuple = (value: unknown): value is [string, string] => {
  return Array.isArray(value) && value.length === 2 && typeof value[0] === "string" && typeof value[1] === "string"
}
/**
 * Builds a trait record from a compact tuple response.
 */
const createTraitFromTuple = (uuid: string, value: [string, string]): Trait => {
  return {t_uuid: uuid, t_key: value[0], t_value: value[1]}
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
const getApiErrorStatus = (error: unknown) => {
  const e = error as ApiErrorLike
  return Number(e?.statusCode ?? e?.status ?? e?.response?.status ?? 0)
}
/**
 * Executes a POST request against the keys API.
 */
const postTraitKeyRequest = async <T>(path: string, body: Record<string, unknown>) => {
  return await useApiJson<T>(path, {
    method: "POST",
    body
  })
}
/**
 * Creates a new standalone trait and returns its uuid.
 */
const createTrait = (t_key: string, t_value: string, meta?: KeyMeta) =>
  useApiJson<Trait>("/traits", {
    method: "POST",
    body: {t_key, t_value, ...(meta ? {meta} : {})}
  })
/**
 * Creates or resolves a set uuid by appending a trait to the current left-side uuid.
 */
const createOrGetSet = (leftUuid: string, right: Pick<TraitInput, "t_key" | "t_value" | "meta">) =>
  useApiJson<ApiResponseWithData<{s_uuid: string; t_uuid?: string; trait?: Trait}>>("/sets", {
    method: "POST",
    body: {
      s_childs: [leftUuid, toTraitPairPayload(right)]
    }
  }).then(response => response?.data)
/**
 * Finds the set uuid that contains the provided trait uuid collection.
 */
const findSetByTraitUuids = (t_uuids: string[]) =>
  useApiJson<ApiResponseWithData<{s_uuid: string}>>("/sets/find", {
    method: "POST",
    body: t_uuids
  }).then(response => response?.data)
/**
 * Fetches key meta by synonym.
 */
const getKeyMeta = async (syn: string, meta?: KeyMeta) => {
  return await postTraitKeyRequest<ApiResponseWithData<TraitKey>>("/keys/meta", {syn, ...(meta ? {meta} : {})})
}
/**
 * Fetches key meta for several ids in one request.
 */
const getKeysMetaBulk = async (ids: number[]) => {
  return await postTraitKeyRequest<ApiResponseWithData<{items: TraitKey[]}>>("/keys/meta/bulk", {ids})
}
/**
 * Fetches enum options for one key synonym.
 */
const getEnumOptions = async (syn: string) => {
  return await postTraitKeyRequest<ApiResponseWithData<{options: string[]}>>("/keys/enum/options", {syn})
}
/**
 * Updates the stored meta definition for one key id.
 */
const updateKeyMeta = async (id: number, meta: KeyMeta) => {
  return await postTraitKeyRequest<ApiResponseWithData<TraitKey>>("/keys/meta/update", {id, meta})
}
/**
 * Loads all traits behind one uuid, supporting both NDJSON sets and JSON single traits.
 */
const fetchTraitsByUuid = async (uuid: string): Promise<Trait[]> => {
  if (!uuid) return []
  const items: Trait[] = []
  try {
    await useAPI().stream(`/traits/resolve/${uuid}`, line => {
      const item = normalizeTraitPayload(uuid, line)
      if (item) items.push(item)
    })
  } catch (error) {
    if (getApiErrorStatus(error) === 404) return []
    throw error
  }
  return items
}
/**
 * Loads one trait by uuid and falls back to the resolve endpoint when needed.
 */
const fetchTraitByUuid = async (uuid: string): Promise<Trait | null> => {
  if (!uuid) return null
  try {
    const payload = await useApiJson<[string, string] | Trait>(`/traits/${uuid}`)
    return normalizeTraitPayload(uuid, payload)
  } catch (error: unknown) {
    if (getApiErrorStatus(error) === 404) return null
    throw error
  }
}
/**
 * Shortens a uuid for compact UI display.
 */
const shortUuid = (value?: string, keep = 8): string => {
  if (!value) return "—"
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
