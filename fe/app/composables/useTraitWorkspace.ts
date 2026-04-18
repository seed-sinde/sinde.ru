import {keyMetaEquals, normalizeMetaForEquality} from "../utils/traitMeta"
import type {Ref} from "vue"
import type {
  Trait,
  TraitInput,
  KeyMeta,
  TraitKey,
  TraitResolveLockOwner,
  TraitResolveLockState
} from "../../shared/types/traits"
const SKIP_FETCH_UUID_STATE_KEY = "skip-fetch-uuid"
const TRAIT_RESOLVE_LOCK_STATE_KEY = "traits-resolve-lock"
const TRAIT_STREAM_STATUS_STATE_KEY = "traits-stream-status"
const useApiStream = (
  path: string,
  onLine: Parameters<ReturnType<typeof useAPI>["stream"]>[1],
  options?: Parameters<ReturnType<typeof useAPI>["stream"]>[2]
) => useAPI().stream(path, onLine, options)
export const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i
/**
 * Returns the uuid route param only for set pages inside the traits workspace.
 */
const resolveTraitsRouteUuid = (route: ReturnType<typeof useRoute>): string | undefined => {
  if (route.path.startsWith("/traits/trait/")) return undefined
  const rawParam = route.params.uuid
  const candidate = typeof rawParam === "string" ? rawParam : Array.isArray(rawParam) ? rawParam[0] : undefined
  return candidate && UUID_RE.test(candidate) ? candidate : undefined
}
/**
 * Creates a detached trait list copy when the store already matches the requested uuid.
 */
const cloneResolvedTraitsFromStore = (store: ReturnType<typeof useTraitsStore>, uuid: string): Trait[] => {
  if (store.currentUuid !== uuid) return []
  return Array.isArray(store.traits) ? [...store.traits] : []
}
/**
 * Returns true when the pending route change may reuse the current store snapshot.
 */
const canReuseSkippedTraitFetch = (
  store: ReturnType<typeof useTraitsStore>,
  skipFetchUuid: Ref<string | null>,
  uuid: string
): boolean => {
  return skipFetchUuid.value === uuid && store.currentUuid === uuid && store.traits.length > 0
}
/**
 * Converts different backend error shapes into one readable message.
 */
const resolveTraitErrorMessage = (error: unknown, fallback: string): string => {
  const source = error as {
    data?: {message?: string} | string
    statusMessage?: string
    message?: string
  }
  return source?.data && typeof source.data === "object" && "message" in source.data
    ? String(source.data.message || fallback)
    : typeof source?.data === "string"
      ? source.data
      : source?.statusMessage || source?.message || fallback
}
/**
 * Normalizes one stream line into a trait record.
 */
const normalizeStreamTraitLine = (line: unknown, fallbackUuid: string): Trait | null => {
  if (Array.isArray(line) && line.length === 2 && typeof line[0] === "string" && typeof line[1] === "string") {
    return {
      t_uuid: fallbackUuid,
      t_key: line[0],
      t_value: line[1]
    }
  }
  if (!line || typeof line !== "object") return null
  const record = line as Partial<Trait>
  const t_key = String(record.t_key || "").trim()
  const t_value = String(record.t_value || "").trim()
  const t_uuid = String(record.t_uuid || fallbackUuid).trim()
  if (!t_uuid || !t_key) return null
  return {
    t_uuid,
    t_key,
    t_value,
    ...(typeof record.t_key_id === "number" ? {t_key_id: record.t_key_id} : {})
  }
}
/**
 * Provides routing helpers and a shared skip-fetch marker for the traits workspace.
 */
export const useTraitNavigation = () => {
  const route = useRoute()
  const router = useRouter()
  const skipFetchUuid = useState<string | null>(SKIP_FETCH_UUID_STATE_KEY, () => null)
  /**
   * Checks whether the provided string looks like a valid UUID.
   */
  const isValidUuid = (uuid: string): boolean => UUID_RE.test(uuid)
  const uuid = computed<string | undefined>(() => resolveTraitsRouteUuid(route))
  /**
   * Updates the address bar to the next traits workspace location.
   */
  const updateAddress = async (override?: string | null) => {
    const candidate = override === undefined ? uuid.value : override
    return candidate && isValidUuid(candidate) ? router.replace(`/traits/${candidate}`) : router.replace("/traits/")
  }
  return {uuid, skipFetchUuid, isValidUuid, updateAddress}
}
/**
 * Serializes loader and stream access to one active traits resolve request per uuid.
 */
export const useTraitResolveLock = () => {
  const state = useState<TraitResolveLockState>(TRAIT_RESOLVE_LOCK_STATE_KEY, () => ({
    uuid: null,
    owner: null
  }))
  /**
   * Attempts to claim the resolve lock for the given owner and uuid.
   */
  const acquire = (uuid: string, owner: TraitResolveLockOwner): boolean => {
    const current = state.value
    if (!current.owner) {
      state.value = {uuid, owner}
      return true
    }
    return current.uuid === uuid && current.owner === owner
  }
  /**
   * Releases the resolve lock only when it is still owned by the caller.
   */
  const release = (uuid: string, owner: TraitResolveLockOwner) => {
    const current = state.value
    if (current.uuid === uuid && current.owner === owner) state.value = {uuid: null, owner: null}
  }
  /**
   * Waits until another owner finishes resolving the same uuid.
   */
  const waitUntilFree = async (uuid: string, signal?: AbortSignal) => {
    while (true) {
      if (signal?.aborted) throw new DOMException("Aborted", "AbortError")
      const current = state.value
      if (!current.owner || current.uuid !== uuid) return
      await new Promise<void>(resolve => setTimeout(resolve, 25))
    }
  }
  return {state, acquire, release, waitUntilFree}
}
/**
 * Exposes normalized trait data types and meta helpers for forms.
 */
export const useTraitTypes = () => {
  const dataTypes = computed(() => DATA_TYPES)
  /**
   * Creates a normalized meta object for the requested data type.
   */
  const createMeta = (dataType: DataType = "string"): KeyMeta => defaultKeyMeta(dataType)
  /**
   * Validates a value against the provided trait meta definition.
   */
  const isValidValue = (value: unknown, meta: KeyMeta): boolean => validateValue(value, meta)
  return {
    dataTypes,
    createMeta,
    isValidValue
  }
}
type TraitLoaderResult = {
  data: Ref<Trait[] | undefined>
  pending: Ref<boolean>
  error: Ref<unknown>
  refresh: () => Promise<void>
}
type TraitStreamStatus = {
  uuid: string | null
  pending: boolean
  completedAt: number
}

/**
 * Shares resolve-stream activity state between workspace modules.
 */
export const useTraitStreamStatus = () =>
  useState<TraitStreamStatus>(TRAIT_STREAM_STATUS_STATE_KEY, () => ({
    uuid: null,
    pending: false,
    completedAt: 0
  }))

/**
 * Loads a full trait list for the current uuid while reusing SSR or in-memory state when possible.
 */
export const useTraitLoader = async (
  uuid: Ref<string | undefined>,
  skipFetchUuid: Ref<string | null>
): Promise<TraitLoaderResult> => {
  const lock = useTraitResolveLock()
  const store = useTraitsStore()
  const key = computed(() => `traits:${uuid.value || ""}`)
  /**
   * Returns a cloned store snapshot for the current uuid.
   */
  const resolveFromStore = (value: string): Trait[] => cloneResolvedTraitsFromStore(store, value)
  const {data, pending, error, refresh} = await useAsyncData(
    () => key.value,
    async () => {
      if (!uuid.value) return []
      const value = uuid.value
      if (canReuseSkippedTraitFetch(store, skipFetchUuid, value)) return resolveFromStore(value)
      let lockAcquired = lock.acquire(value, "loader")
      if (!lockAcquired) {
        await lock.waitUntilFree(value)
        if (store.currentUuid === value) return resolveFromStore(value)
        lockAcquired = lock.acquire(value, "loader")
      }
      if (!lockAcquired) return resolveFromStore(value)
      try {
        return await fetchTraitsByUuid(value)
      } finally {
        lock.release(value, "loader")
      }
    },
    {
      server: true,
      lazy: false,
      watch: [uuid]
    }
  )
  return {data, pending, error, refresh}
}
/**
 * Subscribes to the NDJSON resolve stream and keeps the local trait store in sync.
 */
export const useTraitStream = (uuid: Ref<string | undefined>, skipFetchUuid: Ref<string | null>) => {
  const store = useTraitsStore()
  const route = useRoute()
  const lock = useTraitResolveLock()
  const status = useTraitStreamStatus()
  const ctrl = shallowRef<AbortController | null>(null)
  const prevUuid = ref<string | undefined>(undefined)
  const prevTraitsLength = ref(store.traits.length)
  let streamRunId = 0
  /**
   * Marks resolve-stream state as completed for a uuid.
   */
  const markCompleted = (value: string | null) =>
    (status.value = {uuid: value, pending: false, completedAt: Date.now()})
  /**
   * Returns true only on the empty traits workspace route.
   */
  const isTraitsRootRoute = (): boolean => /^\/traits\/?$/.test(route.path)
  /**
   * Returns true when the store already contains resolved traits for the target uuid.
   */
  const hasCurrentStoreTraits = (value: string): boolean => store.currentUuid === value && store.traits.length > 0
  /**
   * Stops the active resolve stream.
   */
  const stop = () => {
    ctrl.value?.abort()
    ctrl.value = null
  }
  /**
   * Starts a fresh resolve stream for the provided uuid.
   */
  const start = async (value: string) => {
    stop()
    const controller = new AbortController()
    ctrl.value = controller
    const runId = ++streamRunId
    status.value = {
      uuid: value,
      pending: true,
      completedAt: status.value.completedAt
    }
    let gotAny = false
    let lockAcquired = false
    const seen = new Set<string>(store.traits.map(trait => trait.t_uuid))
    const buffer: Trait[] = []
    let flushTimer: ReturnType<typeof setTimeout> | null = null
    const flush = () => {
      if (!buffer.length) return
      store.addTraits(buffer.splice(0, buffer.length))
    }
    const scheduleFlush = () => {
      if (flushTimer) return
      flushTimer = setTimeout(() => {
        flushTimer = null
        flush()
      }, 50)
    }
    try {
      try {
        if (!lock.acquire(value, "stream")) {
          await lock.waitUntilFree(value, controller.signal)
          if (ctrl.value !== controller) return
          if (store.traits.length > 0) return
        }
        if (!lock.acquire(value, "stream")) return
        lockAcquired = true
        await useApiStream(
          `/traits/resolve/${value}`,
          (line: unknown) => {
            if (ctrl.value !== controller) return
            const normalizedLine = normalizeStreamTraitLine(line, value)
            if (!normalizedLine) return
            if (seen.has(normalizedLine.t_uuid)) return
            if (store.traits.some(item => item.t_uuid === normalizedLine.t_uuid)) {
              seen.add(normalizedLine.t_uuid)
              return
            }
            buffer.push(normalizedLine)
            seen.add(normalizedLine.t_uuid)
            gotAny = true
            scheduleFlush()
          },
          {signal: controller.signal}
        )
        if (flushTimer) {
          clearTimeout(flushTimer)
          flushTimer = null
        }
        flush()
      } catch (error: any) {
        if (flushTimer) {
          clearTimeout(flushTimer)
          flushTimer = null
        }
        flush()
        if (error?.name === "AbortError") return
      }
      if (!gotAny && store.traits.length === 0) {
        try {
          const fallbackItems = await fetchTraitsByUuid(value)
          if (ctrl.value === controller && fallbackItems.length > 0) store.setTraits(fallbackItems)
        } catch (error) {
          if (import.meta.client) {
            console.error(`fallback failed`, error)
          }
        }
      }
    } finally {
      if (lockAcquired) lock.release(value, "stream")
      if (runId === streamRunId) {
        if (ctrl.value === controller) ctrl.value = null
        markCompleted(value)
      }
    }
  }
  watch(
    uuid,
    value => {
      if (!import.meta.client) return
      if (!value) {
        stop()
        streamRunId += 1
        markCompleted(null)
        if (isTraitsRootRoute()) {
          store.clear()
        }
        prevUuid.value = undefined
        return
      }
      if (skipFetchUuid.value && skipFetchUuid.value === value) {
        skipFetchUuid.value = null
        if (hasCurrentStoreTraits(value)) {
          store.currentUuid = value
          markCompleted(value)
          prevUuid.value = value
          return
        }
      }
      if (!prevUuid.value && store.currentUuid === value && store.traits.length > 0) {
        markCompleted(value)
        prevUuid.value = value
        return
      }
      if (store.currentUuid && store.currentUuid !== value) store.setTraits([])
      store.currentUuid = value
      prevUuid.value = value
      void start(value)
    },
    {immediate: true}
  )
  watch(
    () => store.traits.length,
    nextLength => {
      const lastLength = prevTraitsLength.value
      prevTraitsLength.value = nextLength
      if (!import.meta.client) return
      if (!uuid.value) return
      if (route.path.startsWith("/traits/trait/")) return
      if (lastLength <= 0 || nextLength !== 0) return
      if (store.currentUuid !== uuid.value) return
      void start(uuid.value)
    }
  )
  onBeforeUnmount(stop)
  return {stop}
}
/**
 * Applies add, remove and edit mutations to the traits workspace and keeps the route in sync.
 */
export const useTraitActions = () => {
  const store = useTraitsStore()
  const {uuid: currentUuid, updateAddress, skipFetchUuid} = useTraitNavigation()
  const {ensureLoaded, user, setPrimaryTraitUuid} = useAuth()
  const inProcess = ref(false)
  const error = ref<string | null>(null)
  /**
   * Finds cached key meta by synonym when the current trait record has no key id yet.
   */
  const findStoredKeyBySyn = (syn: string): TraitKey | null => {
    const targetSyn = String(syn || "")
      .trim()
      .toLowerCase()
    if (!targetSyn) return null
    const entries = Object.values(store.keyMetaById || {}) as TraitKey[]
    const found = entries.find(
      entry =>
        String(entry?.syn || "")
          .trim()
          .toLowerCase() === targetSyn
    )
    return found || null
  }
  /**
   * Attaches a resolved key id to local traits that were just created in this session.
   */
  const attachLocalKeyId = (traitUuids: string[], keyId?: number | null) => {
    const normalizedKeyId = Number(keyId || 0)
    if (normalizedKeyId <= 0 || !traitUuids.length) return
    const targets = new Set(traitUuids.map(value => String(value || "").trim()).filter(Boolean))
    if (!targets.size) return
    for (const trait of store.traits) {
      const uuid = String(trait?.t_uuid || "").trim()
      if (!targets.has(uuid)) continue
      trait.t_key_id = normalizedKeyId
    }
  }
  /**
   * Synchronizes the user's primary trait uuid after a local mutation changes the active target.
   */
  const syncPrimaryAfterMutation = async (previousUuid: string | null, nextUuid: string | null) => {
    if (!import.meta.client) return
    try {
      await ensureLoaded()
      const authUser = user.value
      if (!authUser) return
      const primaryUuid = String(authUser.primary_trait_uuid || "").trim()
      const startedFromScratch = !primaryUuid && !previousUuid
      const editingPrimary = Boolean(primaryUuid && previousUuid && previousUuid === primaryUuid)
      if (!startedFromScratch && !editingPrimary) return
      const resolvedNextUuid = String(nextUuid || "").trim()
      if (!resolvedNextUuid) {
        if (editingPrimary && primaryUuid) {
          await setPrimaryTraitUuid(null)
        }
        return
      }
      if (resolvedNextUuid === primaryUuid) return
      await setPrimaryTraitUuid(resolvedNextUuid)
    } catch (syncError) {
      console.warn("primary trait sync failed", syncError)
    }
  }
  /**
   * Updates the workspace route and marks the next uuid as a reusable local snapshot.
   */
  const applyNextUuid = async (previousUuid: string | null, nextUuid: string | null) => {
    store.currentUuid = nextUuid
    skipFetchUuid.value = nextUuid
    await updateAddress(nextUuid)
    await syncPrimaryAfterMutation(previousUuid, nextUuid)
  }
  /**
   * Finds existing meta for a store trait, preferring resolved key ids.
   */
  const resolveExistingTraitMeta = (trait: Trait, fallbackMeta?: KeyMeta): KeyMeta | undefined => {
    if (trait.t_key_id) {
      const entry = (store.keyMetaById as Record<number, TraitKey | undefined>)?.[trait.t_key_id]
      if (entry?.meta) return entry.meta as KeyMeta
    }
    const bySyn = findStoredKeyBySyn(trait.t_key)
    if (bySyn?.meta) return bySyn.meta as KeyMeta
    return fallbackMeta
  }
  /**
   * Adds one trait to the current workspace or creates the first workspace record.
   */
  const addTrait = async ({t_key, t_value, meta}: TraitInput) => {
    if (inProcess.value) return
    inProcess.value = true
    error.value = null
    const previousUuid = currentUuid.value ?? null
    try {
      const key = t_key.trim()
      const value = String(t_value ?? "").trim()
      if (!key || !value) throw new Error("Укажите ключ и значение")
      const targetMeta = normalizeMetaForEquality(meta)
      const existingSameKey = store.traits.find(
        trait =>
          trait.t_key.trim().toLowerCase() === key.toLowerCase() &&
          keyMetaEquals(resolveExistingTraitMeta(trait), targetMeta)
      )
      if (existingSameKey && existingSameKey.t_value.trim().toLowerCase() === value.toLowerCase()) {
        return null
      }
      const duplicateValue = store.traits.some(
        trait =>
          trait.t_uuid !== existingSameKey?.t_uuid &&
          trait.t_key.trim().toLowerCase() === key.toLowerCase() &&
          keyMetaEquals(resolveExistingTraitMeta(trait), targetMeta) &&
          trait.t_value.trim().toLowerCase() === value.toLowerCase()
      )
      if (duplicateValue) throw new Error("Эта особенность уже есть в списке")
      if (existingSameKey) {
        const replaceIndex = store.traits.findIndex(trait => trait.t_uuid === existingSameKey.t_uuid)
        const baseTraits = store.traits.filter(trait => trait.t_uuid !== existingSameKey.t_uuid)
        const createdTrait = await createTrait(key, value, targetMeta)
        if (createdTrait.t_key_id) attachLocalKeyId([createdTrait.t_uuid], createdTrait.t_key_id)
        const nextTraits = [...baseTraits]
        const insertAt = replaceIndex >= 0 ? replaceIndex : nextTraits.length
        nextTraits.splice(insertAt, 0, createdTrait)
        store.setTraits(nextTraits)
        let finalUuid: string | null = createdTrait.t_uuid
        if (nextTraits.length > 1) {
          const resolvedSet = await findSetByTraitUuids(nextTraits.map(trait => trait.t_uuid))
          finalUuid = resolvedSet?.s_uuid || finalUuid
        }
        await applyNextUuid(previousUuid, finalUuid)
        return {action: "updated" as const, traitUuid: createdTrait.t_uuid}
      }
      if (!currentUuid.value || store.traits.length === 0) {
        const createdTrait = await createTrait(key, value, targetMeta)
        if (createdTrait.t_key_id) attachLocalKeyId([createdTrait.t_uuid], createdTrait.t_key_id)
        store.setTraits([createdTrait])
        await applyNextUuid(previousUuid, createdTrait.t_uuid)
        return {action: "added" as const, traitUuid: createdTrait.t_uuid}
      }
      const {s_uuid, t_uuid: createdUuid, trait: createdTrait} = await createOrGetSet(currentUuid.value, {
        t_key: key,
        t_value: value,
        meta: targetMeta
      })
      if (createdUuid && createdTrait) {
        store.addTrait(createdTrait)
        if (createdTrait.t_key_id) attachLocalKeyId([createdTrait.t_uuid], createdTrait.t_key_id)
        await applyNextUuid(previousUuid, s_uuid)
        return {action: "added" as const, traitUuid: createdTrait.t_uuid}
      }
      await applyNextUuid(previousUuid, s_uuid)
      return null
    } catch (mutationError) {
      error.value = resolveTraitErrorMessage(mutationError, "Не удалось добавить особенность")
      return null
    } finally {
      inProcess.value = false
    }
  }
  /**
   * Removes selected traits and resolves the next active uuid for the workspace.
   */
  const removeSelected = async (selectedIds: string[] | {value: string[]}) => {
    const ids = Array.isArray(selectedIds)
      ? selectedIds
      : Array.isArray((selectedIds as {value?: unknown})?.value)
        ? (selectedIds as {value: string[]}).value
        : []
    if (!ids.length) return
    if (inProcess.value) return
    inProcess.value = true
    error.value = null
    const previousUuid = currentUuid.value ?? null
    try {
      store.removeTraits(ids)
      const remainingTraits = store.traits
      if (remainingTraits.length === 0) {
        await applyNextUuid(previousUuid, null)
        return
      }
      if (remainingTraits.length === 1) {
        await applyNextUuid(previousUuid, remainingTraits[0]?.t_uuid || null)
        return
      }
      const resolvedSet = await findSetByTraitUuids(remainingTraits.map(trait => trait.t_uuid))
      store.currentUuid = resolvedSet.s_uuid
      await applyNextUuid(previousUuid, resolvedSet.s_uuid)
    } catch (mutationError) {
      error.value = resolveTraitErrorMessage(mutationError, "Не удалось удалить особенности")
    } finally {
      inProcess.value = false
    }
  }
  /**
   * Replaces one trait value with a newly created canonical trait record.
   */
  const editTraitValue = async (payload: {traitUuid: string; t_value: string}) => {
    if (inProcess.value) return null
    inProcess.value = true
    error.value = null
    const previousUuid = currentUuid.value ?? null
    try {
      const traitUuid = String(payload?.traitUuid || "").trim()
      const nextValue = String(payload?.t_value || "").trim()
      if (!traitUuid) throw new Error("Особенность для редактирования не найдена")
      if (!nextValue) throw new Error("Укажите значение")
      const replaceIndex = store.traits.findIndex(trait => trait.t_uuid === traitUuid)
      if (replaceIndex < 0) throw new Error("Выбранная особенность отсутствует в текущем наборе")
      const currentTrait = store.traits[replaceIndex]
      if (!currentTrait) throw new Error("Не удалось получить выбранную особенность")
      const currentValue = String(currentTrait.t_value || "").trim()
      if (currentValue === nextValue) return {nextTraitUuid: currentTrait.t_uuid, replacedTraitUuid: traitUuid}
      const createdTrait = await createTrait(
        currentTrait.t_key,
        nextValue,
        normalizeMetaForEquality(resolveExistingTraitMeta(currentTrait))
      )
      const replacedTrait: Trait = {
        ...currentTrait,
        ...createdTrait,
        t_value: nextValue
      }
      store.replaceTrait(traitUuid, replacedTrait)
      let finalUuid: string | null = replacedTrait.t_uuid
      if (store.traits.length > 1) {
        const resolvedSet = await findSetByTraitUuids(store.traits.map(trait => trait.t_uuid))
        finalUuid = resolvedSet?.s_uuid || finalUuid
      }
      await applyNextUuid(previousUuid, finalUuid)
      return {nextTraitUuid: replacedTrait.t_uuid, replacedTraitUuid: traitUuid}
    } catch (mutationError) {
      error.value = resolveTraitErrorMessage(mutationError, "Не удалось изменить значение особенности")
      return null
    } finally {
      inProcess.value = false
    }
  }
  return {addTrait, removeSelected, editTraitValue, inProcess, error}
}
/**
 * Reads a traits uuid from the clipboard and navigates to the resolved route.
 */
export const usePasteUuid = () => {
  const router = useRouter()
  const {pasteInto} = useClipboard()
  const error = useState<string>("traits-paste-uuid-error", () => "")
  error.value = ""
  /**
   * Extracts a validated uuid from the current clipboard contents.
   */
  const readUuidFromClipboard = async () => {
    error.value = ""
    try {
      const next = ref("")
      await pasteInto(next)
      const value = String(next.value || "").trim()
      if (!UUID_RE.test(value)) {
        error.value = "В буфере обмена нет корректного UUID."
        return null
      }
      return value
    } catch {
      error.value = "Не удалось прочитать буфер обмена. Разрешите доступ в браузере."
      return null
    }
  }
  /**
   * Navigates to the route built from the clipboard uuid.
   */
  const pasteUuidAndNavigate = async (to: string | ((uuid: string) => string)) => {
    const uuid = await readUuidFromClipboard()
    if (!uuid) return null
    const target = typeof to === "function" ? to(uuid) : to
    await router.replace(target)
    return uuid
  }
  return {
    error,
    readUuidFromClipboard,
    pasteUuidAndNavigate
  }
}
/**
 * Tracks checkbox selection state for trait collections.
 */
export const useSelectionMap = <T extends {t_uuid: string}>(items: Ref<T[]>) => {
  const map = reactive<Record<string, boolean>>({})
  const hasSelected = computed(() => Object.values(map).some(Boolean))
  const isAllSelected = computed(() => items.value.length > 0 && items.value.every(item => map[item.t_uuid]))
  /**
   * Toggles the whole current collection selection.
   */
  const toggleSelectAll = () => {
    const nextValue = !isAllSelected.value
    const next: Record<string, boolean> = {}
    if (nextValue) {
      items.value.forEach(item => {
        next[item.t_uuid] = true
      })
    }
    Object.keys(map).forEach(key => {
      map[key] = false
    })
    Object.assign(map, next)
  }
  /**
   * Clears all selected ids.
   */
  const clear = () => {
    Object.keys(map).forEach(key => {
      map[key] = false
    })
  }
  const selectedIds = computed<string[]>(() =>
    Object.entries(map)
      .filter(([, value]) => value)
      .map(([id]) => id)
  )
  const selectedItems = computed<T[]>(() => {
    const set = new Set(selectedIds.value)
    return items.value.filter(item => set.has(item.t_uuid))
  })
  return {selectedMap: map, hasSelected, isAllSelected, toggleSelectAll, clear, selectedIds, selectedItems}
}
