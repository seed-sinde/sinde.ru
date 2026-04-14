import { keyMetaEquals, normalizeMetaForEquality } from '../utils/traitMeta'
import type { Ref } from 'vue'
import type {
  Trait,
  TraitInput,
  KeyMeta,
  TraitKey,
  TraitResolveLockOwner,
  TraitResolveLockState
} from '../../shared/types/traits'
const SKIP_FETCH_UUID_STATE_KEY = 'skip-fetch-uuid'
const TRAIT_RESOLVE_LOCK_STATE_KEY = 'traits-resolve-lock'
const useApiJson = <T>(path: string, options?: NonNullable<Parameters<ReturnType<typeof useAPI>['json']>[1]>) =>
  useAPI().json<T>(path, options)
const useApiStream = (
  path: string,
  onLine: Parameters<ReturnType<typeof useAPI>['stream']>[1],
  options?: Parameters<ReturnType<typeof useAPI>['stream']>[2]
) => useAPI().stream(path, onLine, options)
export const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i
/**
 * Returns the uuid route param only for set pages inside the traits workspace.
 */
const resolveTraitsRouteUuid = (route: ReturnType<typeof useRoute>): string | undefined => {
  if (route.path.startsWith('/traits/trait/')) return undefined
  const rawParam = route.params.uuid
  const candidate = typeof rawParam === 'string' ? rawParam : Array.isArray(rawParam) ? rawParam[0] : undefined
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
    data?: { message?: string } | string
    statusMessage?: string
    message?: string
  }
  return source?.data && typeof source.data === 'object' && 'message' in source.data
    ? String(source.data.message || fallback)
    : typeof source?.data === 'string'
      ? source.data
      : source?.statusMessage || source?.message || fallback
}
/**
 * Extracts an HTTP-like status code from an unknown API error.
 */
const resolveTraitErrorStatus = (error: unknown): number => {
  const source = error as {
    statusCode?: number
    status?: number
    response?: { status?: number }
  }
  return Number(source?.statusCode ?? source?.status ?? source?.response?.status ?? 0)
}
/**
 * Deduplicates traits by uuid while preserving the original order.
 */
const dedupeTraitsByUuid = (traits: Trait[]): Trait[] => {
  const result: Trait[] = []
  const seen = new Set<string>()
  for (const trait of traits) {
    const uuid = String(trait?.t_uuid || '').trim()
    if (!uuid || seen.has(uuid)) continue
    seen.add(uuid)
    result.push(trait)
  }
  return result
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
    return candidate && isValidUuid(candidate) ? router.replace(`/traits/${candidate}`) : router.replace('/traits/')
  }
  return { uuid, skipFetchUuid, isValidUuid, updateAddress }
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
      state.value = { uuid, owner }
      return true
    }
    return current.uuid === uuid && current.owner === owner
  }
  /**
   * Releases the resolve lock only when it is still owned by the caller.
   */
  const release = (uuid: string, owner: TraitResolveLockOwner) => {
    const current = state.value
    if (current.uuid === uuid && current.owner === owner) {
      state.value = { uuid: null, owner: null }
    }
  }
  /**
   * Waits until another owner finishes resolving the same uuid.
   */
  const waitUntilFree = async (uuid: string, signal?: AbortSignal) => {
    while (true) {
      if (signal?.aborted) throw new DOMException('Aborted', 'AbortError')
      const current = state.value
      if (!current.owner || current.uuid !== uuid) return
      await new Promise<void>((resolve) => setTimeout(resolve, 25))
    }
  }
  return { state, acquire, release, waitUntilFree }
}
/**
 * Exposes normalized trait data types and meta helpers for forms.
 */
export const useTraitTypes = () => {
  const dataTypes = computed(() => DATA_TYPES)
  /**
   * Creates a normalized meta object for the requested data type.
   */
  const createMeta = (dataType: DataType = 'string'): KeyMeta => defaultKeyMeta(dataType)
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

/**
 * Loads a full trait list for the current uuid while reusing SSR or in-memory state when possible.
 */
export const useTraitLoader = async (
  uuid: Ref<string | undefined>,
  skipFetchUuid: Ref<string | null>
): Promise<TraitLoaderResult> => {
  const lock = useTraitResolveLock()
  const store = useTraitsStore()
  const key = computed(() => `traits:${uuid.value || ''}`)
  /**
   * Returns a cloned store snapshot for the current uuid.
   */
  const resolveFromStore = (value: string): Trait[] => cloneResolvedTraitsFromStore(store, value)
  const { data, pending, error, refresh } = await useAsyncData(
    () => key.value,
    async () => {
      if (!uuid.value) return []
      const value = uuid.value
      if (canReuseSkippedTraitFetch(store, skipFetchUuid, value)) {
        return resolveFromStore(value)
      }
      let lockAcquired = lock.acquire(value, 'loader')
      if (!lockAcquired) {
        await lock.waitUntilFree(value)
        if (store.currentUuid === value) {
          return resolveFromStore(value)
        }
        lockAcquired = lock.acquire(value, 'loader')
      }
      if (!lockAcquired) {
        return resolveFromStore(value)
      }
      try {
        return await fetchTraitsByUuid(value)
      } finally {
        lock.release(value, 'loader')
      }
    },
    {
      server: true,
      lazy: false,
      watch: [uuid]
    }
  )
  return { data, pending, error, refresh }
}
/**
 * Subscribes to the NDJSON resolve stream and keeps the local trait store in sync.
 */
export const useTraitStream = (uuid: Ref<string | undefined>, skipFetchUuid: Ref<string | null>) => {
  const store = useTraitsStore()
  const route = useRoute()
  const lock = useTraitResolveLock()
  const ctrl = shallowRef<AbortController | null>(null)
  const prevUuid = ref<string | undefined>(undefined)
  const prevTraitsLength = ref(store.traits.length)
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
    let gotAny = false
    let lockAcquired = false
    const seen = new Set<string>(store.traits.map((trait) => trait.t_uuid))
    const buffer: Trait[] = []
    let flushTimer: ReturnType<typeof setTimeout> | null = null
    /**
     * Flushes the buffered stream records into the store in one batch.
     */
    const flush = () => {
      if (!buffer.length) return
      store.addTraits(buffer.splice(0, buffer.length))
    }
    /**
     * Schedules the next buffer flush to avoid store churn on dense streams.
     */
    const scheduleFlush = () => {
      if (flushTimer) return
      flushTimer = setTimeout(() => {
        flushTimer = null
        flush()
      }, 50)
    }
    try {
      try {
        if (!lock.acquire(value, 'stream')) {
          await lock.waitUntilFree(value, controller.signal)
          if (ctrl.value !== controller) return
          if (store.traits.length > 0) return
        }
        if (!lock.acquire(value, 'stream')) return
        lockAcquired = true
        await useApiStream(
          `/traits/resolve/${value}`,
          (line: Trait) => {
            if (ctrl.value !== controller) return
            if (seen.has(line.t_uuid)) return
            if (store.traits.some((item) => item.t_uuid === line.t_uuid)) {
              seen.add(line.t_uuid)
              return
            }
            buffer.push(line)
            seen.add(line.t_uuid)
            gotAny = true
            scheduleFlush()
          },
          { signal: controller.signal }
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
        if (error?.name === 'AbortError') return
      }
      if (!gotAny && store.traits.length === 0) {
        try {
          const payload = await useApiJson<[string, string]>(`/traits/resolve/${value}`)
          if (Array.isArray(payload) && payload.length === 2 && ctrl.value === controller) {
            const [t_key, t_value] = payload
            store.setTraits([{ t_uuid: value, t_key, t_value }])
          }
        } catch {
          // ignore fallback errors
        }
      }
    } finally {
      if (lockAcquired) {
        lock.release(value, 'stream')
      }
    }
  }
  watch(
    uuid,
    (value) => {
      if (!import.meta.client) return
      if (!value) {
        stop()
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
          prevUuid.value = value
          return
        }
      }
      if (!prevUuid.value && store.currentUuid === value && store.traits.length > 0) {
        prevUuid.value = value
        return
      }
      if (store.currentUuid && store.currentUuid !== value) {
        store.setTraits([])
      }
      store.currentUuid = value
      prevUuid.value = value
      void start(value)
    },
    { immediate: true }
  )
  watch(
    () => store.traits.length,
    (nextLength) => {
      const lastLength = prevTraitsLength.value
      prevTraitsLength.value = nextLength
      if (!import.meta.client) return
      if (!uuid.value) return
      if (route.path.startsWith('/traits/trait/')) return
      if (lastLength <= 0 || nextLength !== 0) return
      if (store.currentUuid !== uuid.value) return
      void start(uuid.value)
    }
  )
  onBeforeUnmount(stop)
  return { stop }
}
/**
 * Applies add, remove and edit mutations to the traits workspace and keeps the route in sync.
 */
export const useTraitActions = () => {
  const store = useTraitsStore()
  const { uuid: currentUuid, updateAddress, skipFetchUuid } = useTraitNavigation()
  const { ensureLoaded, user, setPrimaryTraitUuid } = useAuth()
  const inProcess = ref(false)
  const error = ref<string | null>(null)
  /**
   * Finds cached key meta by synonym when the current trait record has no key id yet.
   */
  const findStoredKeyBySyn = (syn: string): TraitKey | null => {
    const targetSyn = String(syn || '')
      .trim()
      .toLowerCase()
    if (!targetSyn) return null
    const entries = Object.values(store.keyMetaById || {}) as TraitKey[]
    const found = entries.find(
      (entry) =>
        String(entry?.syn || '')
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
    const targets = new Set(traitUuids.map((value) => String(value || '').trim()).filter(Boolean))
    if (!targets.size) return
    store.setTraits(
      store.traits.map((trait) =>
        targets.has(String(trait.t_uuid || '').trim()) ? { ...trait, t_key_id: normalizedKeyId } : trait
      )
    )
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
      const primaryUuid = String(authUser.primary_trait_uuid || '').trim()
      const startedFromScratch = !primaryUuid && !previousUuid
      const editingPrimary = Boolean(primaryUuid && previousUuid && previousUuid === primaryUuid)
      if (!startedFromScratch && !editingPrimary) return
      const resolvedNextUuid = String(nextUuid || '').trim()
      if (!resolvedNextUuid) {
        if (editingPrimary && primaryUuid) {
          await setPrimaryTraitUuid(null)
        }
        return
      }
      if (resolvedNextUuid === primaryUuid) return
      await setPrimaryTraitUuid(resolvedNextUuid)
    } catch (syncError) {
      console.warn('primary trait sync failed', syncError)
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
   * Reads key meta while treating 404 as an empty result.
   */
  const fetchKeyMetaSafe = async (syn: string) => {
    try {
      return await getKeyMeta(syn)
    } catch (requestError) {
      if (resolveTraitErrorStatus(requestError) === 404) return null
      throw requestError
    }
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
  const addTrait = async ({ t_key, t_value, meta }: TraitInput) => {
    if (inProcess.value) return
    inProcess.value = true
    error.value = null
    const previousUuid = currentUuid.value ?? null
    try {
      const key = t_key.trim()
      const value = String(t_value ?? '').trim()
      if (!key || !value) throw new Error('Укажите ключ и значение')
      const existingKeyMeta = await fetchKeyMetaSafe(key)
      let keyId = existingKeyMeta?.data?.id
      const existingMeta = existingKeyMeta?.data?.meta as KeyMeta | undefined
      if (existingKeyMeta?.data) {
        store.setKeyMetaBulk([existingKeyMeta.data])
      }
      const keyAlreadyExists = Boolean(keyId)
      /**
       * Writes meta for brand-new keys after the trait itself is created.
       */
      const applyMetaIfNewKey = async (traitUuids: string[] = []) => {
        if (keyAlreadyExists || !meta) return
        const freshMeta = await fetchKeyMetaSafe(key)
        keyId = freshMeta?.data?.id ?? keyId
        if (!keyId) return
        const updated = await updateKeyMeta(keyId, meta)
        if (updated?.data) {
          keyId = updated.data.id ?? keyId
          store.setKeyMetaBulk([updated.data])
          attachLocalKeyId(traitUuids, keyId)
        }
      }
      const targetMeta = normalizeMetaForEquality(meta || existingMeta)
      const existingSameKey = store.traits.find(
        (trait) =>
          trait.t_key.trim().toLowerCase() === key.toLowerCase() &&
          keyMetaEquals(resolveExistingTraitMeta(trait, existingMeta), targetMeta)
      )
      if (existingSameKey && existingSameKey.t_value.trim().toLowerCase() === value.toLowerCase()) {
        await applyMetaIfNewKey()
        return
      }
      const duplicateValue = store.traits.some(
        (trait) =>
          trait.t_uuid !== existingSameKey?.t_uuid &&
          trait.t_key.trim().toLowerCase() === key.toLowerCase() &&
          keyMetaEquals(resolveExistingTraitMeta(trait, existingMeta), targetMeta) &&
          trait.t_value.trim().toLowerCase() === value.toLowerCase()
      )
      if (duplicateValue) throw new Error('Эта особенность уже есть в списке')
      if (existingSameKey) {
        const replaceIndex = store.traits.findIndex((trait) => trait.t_uuid === existingSameKey.t_uuid)
        const baseTraits = store.traits.filter((trait) => trait.t_uuid !== existingSameKey.t_uuid)
        const createdUuid = await createTrait(key, value)
        const createdTrait: Trait = { t_uuid: createdUuid, t_key: key, t_value: value }
        if (keyId) createdTrait.t_key_id = keyId
        const nextTraits = [...baseTraits]
        const insertAt = replaceIndex >= 0 ? replaceIndex : nextTraits.length
        nextTraits.splice(insertAt, 0, createdTrait)
        store.setTraits(nextTraits)
        let finalUuid: string | null = createdTrait.t_uuid
        if (nextTraits.length > 1) {
          const resolvedSet = await findSetByTraitUuids(nextTraits.map((trait) => trait.t_uuid))
          finalUuid = resolvedSet?.s_uuid || finalUuid
        }
        await applyMetaIfNewKey([createdTrait.t_uuid])
        await applyNextUuid(previousUuid, finalUuid)
        return
      }
      if (!currentUuid.value || store.traits.length === 0) {
        const createdUuid = await createTrait(key, value)
        const createdTrait: Trait = { t_uuid: createdUuid, t_key: key, t_value: value }
        if (keyId) createdTrait.t_key_id = keyId
        store.setTraits([createdTrait])
        await applyMetaIfNewKey([createdTrait.t_uuid])
        await applyNextUuid(previousUuid, createdUuid)
        return
      }
      const { s_uuid, t_uuid: createdUuid } = await createOrGetSet(currentUuid.value, { t_key: key, t_value: value })
      if (createdUuid) {
        const createdTrait: Trait = { t_uuid: createdUuid, t_key: key, t_value: value }
        if (keyId) createdTrait.t_key_id = keyId
        store.setTraits([...store.traits, createdTrait])
        await applyMetaIfNewKey([createdTrait.t_uuid])
      } else {
        await applyMetaIfNewKey()
      }
      await applyNextUuid(previousUuid, s_uuid)
    } catch (mutationError) {
      error.value = resolveTraitErrorMessage(mutationError, 'Не удалось добавить особенность')
    } finally {
      inProcess.value = false
    }
  }
  /**
   * Removes selected traits and resolves the next active uuid for the workspace.
   */
  const removeSelected = async (selectedIds: string[] | { value: string[] }) => {
    const ids = Array.isArray(selectedIds)
      ? selectedIds
      : Array.isArray((selectedIds as { value?: unknown })?.value)
        ? (selectedIds as { value: string[] }).value
        : []
    if (!ids.length) return
    const previousUuid = currentUuid.value ?? null
    ids.forEach((id) => store.removeTrait(id))
    const remainingTraits = store.traits
    if (remainingTraits.length === 0) {
      await applyNextUuid(previousUuid, null)
      return
    }
    if (remainingTraits.length === 1) {
      await applyNextUuid(previousUuid, remainingTraits[0]?.t_uuid || null)
      return
    }
    const resolvedSet = await findSetByTraitUuids(remainingTraits.map((trait) => trait.t_uuid))
    store.currentUuid = resolvedSet.s_uuid
    await applyNextUuid(previousUuid, resolvedSet.s_uuid)
  }
  /**
   * Replaces one trait value with a newly created canonical trait record.
   */
  const editTraitValue = async (payload: { traitUuid: string; t_value: string }) => {
    if (inProcess.value) return null
    inProcess.value = true
    error.value = null
    const previousUuid = currentUuid.value ?? null
    try {
      const traitUuid = String(payload?.traitUuid || '').trim()
      const nextValue = String(payload?.t_value || '').trim()
      if (!traitUuid) throw new Error('Особенность для редактирования не найдена')
      if (!nextValue) throw new Error('Укажите значение')
      const replaceIndex = store.traits.findIndex((trait) => trait.t_uuid === traitUuid)
      if (replaceIndex < 0) throw new Error('Выбранная особенность отсутствует в текущем наборе')
      const currentTrait = store.traits[replaceIndex]
      if (!currentTrait) throw new Error('Не удалось получить выбранную особенность')
      const currentValue = String(currentTrait.t_value || '').trim()
      if (currentValue === nextValue) return currentTrait.t_uuid
      const createdUuid = await createTrait(currentTrait.t_key, nextValue)
      const replacedTrait: Trait = {
        ...currentTrait,
        t_uuid: createdUuid,
        t_value: nextValue
      }
      const nextTraits = [...store.traits]
      nextTraits.splice(replaceIndex, 1, replacedTrait)
      const normalizedTraits = dedupeTraitsByUuid(nextTraits)
      store.setTraits(normalizedTraits)
      let finalUuid: string | null = replacedTrait.t_uuid
      if (normalizedTraits.length > 1) {
        const resolvedSet = await findSetByTraitUuids(normalizedTraits.map((trait) => trait.t_uuid))
        finalUuid = resolvedSet?.s_uuid || finalUuid
      }
      await applyNextUuid(previousUuid, finalUuid)
      return replacedTrait.t_uuid
    } catch (mutationError) {
      error.value = resolveTraitErrorMessage(mutationError, 'Не удалось изменить значение особенности')
      return null
    } finally {
      inProcess.value = false
    }
  }
  return { addTrait, removeSelected, editTraitValue, inProcess, error }
}
/**
 * Reads a traits uuid from the clipboard and navigates to the resolved route.
 */
export const usePasteUuid = () => {
  const router = useRouter()
  const { pasteInto } = useClipboard()
  const error = ref('')
  /**
   * Extracts a validated uuid from the current clipboard contents.
   */
  const readUuidFromClipboard = async () => {
    error.value = ''
    try {
      const next = ref('')
      await pasteInto(next)
      const value = String(next.value || '').trim()
      if (!UUID_RE.test(value)) {
        error.value = 'В буфере обмена нет корректного UUID.'
        return null
      }
      return value
    } catch {
      error.value = 'Не удалось прочитать буфер обмена. Разрешите доступ в браузере.'
      return null
    }
  }
  /**
   * Navigates to the route built from the clipboard uuid.
   */
  const pasteUuidAndNavigate = async (to: string | ((uuid: string) => string)) => {
    const uuid = await readUuidFromClipboard()
    if (!uuid) return null
    const target = typeof to === 'function' ? to(uuid) : to
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
export const useSelectionMap = <T extends { t_uuid: string }>(items: Ref<T[]>) => {
  const map = reactive<Record<string, boolean>>({})
  const hasSelected = computed(() => Object.values(map).some(Boolean))
  const isAllSelected = computed(() => items.value.length > 0 && items.value.every((item) => map[item.t_uuid]))
  /**
   * Toggles the whole current collection selection.
   */
  const toggleSelectAll = () => {
    const nextValue = !isAllSelected.value
    const next: Record<string, boolean> = {}
    if (nextValue) {
      items.value.forEach((item) => {
        next[item.t_uuid] = true
      })
    }
    Object.keys(map).forEach((key) => {
      map[key] = false
    })
    Object.assign(map, next)
  }
  /**
   * Clears all selected ids.
   */
  const clear = () => {
    Object.keys(map).forEach((key) => {
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
    return items.value.filter((item) => set.has(item.t_uuid))
  })
  return { selectedMap: map, hasSelected, isAllSelected, toggleSelectAll, clear, selectedIds, selectedItems }
}
