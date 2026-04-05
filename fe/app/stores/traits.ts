const dedupeTraits = (list: Trait[]) => {
  if (!Array.isArray(list) || list.length === 0) return []
  const seen = new Set<string>()
  const normalized: Trait[] = []
  for (const item of list) {
    const uuid = String(item?.t_uuid || '').trim()
    if (!uuid || seen.has(uuid)) continue
    seen.add(uuid)
    normalized.push(item)
  }
  return normalized
}
export const useTraitsStore = defineStore('traits', () => {
  const STORAGE_KEY = 'traits.store.v1'
  const traits = ref<Trait[]>([])
  const currentUuid = ref<string | null>(null)
  const keyMetaById = ref<Record<number, TraitKey>>({})
  const isStoreEmpty = () => {
    return traits.value.length === 0 && currentUuid.value === null && Object.keys(keyMetaById.value).length === 0
  }
  const restorePersisted = () => {
    if (!import.meta.client) return
    if (!isStoreEmpty()) return
    try {
      const raw = window.localStorage.getItem(STORAGE_KEY)
      if (!raw) return
      const parsed = JSON.parse(raw) as Partial<{
        traits: Trait[]
        currentUuid: string | null
        keyMetaById: Record<number, TraitKey>
      }>
      if (Array.isArray(parsed.traits)) traits.value = parsed.traits
      if (typeof parsed.currentUuid === 'string' || parsed.currentUuid === null) {
        currentUuid.value = parsed.currentUuid
      }
      if (parsed.keyMetaById && typeof parsed.keyMetaById === 'object') {
        keyMetaById.value = parsed.keyMetaById
      }
    } catch {
      window.localStorage.removeItem(STORAGE_KEY)
    }
  }
  const persistState = () => {
    if (!import.meta.client) return
    try {
      window.localStorage.setItem(
        STORAGE_KEY,
        JSON.stringify({
          traits: traits.value,
          currentUuid: currentUuid.value,
          keyMetaById: keyMetaById.value
        })
      )
    } catch {
      // Ignore storage quota/serialization issues.
    }
  }
  if (import.meta.client) {
    onNuxtReady(() => {
      restorePersisted()
    })
    watch([traits, currentUuid, keyMetaById], persistState, { deep: true })
  }
  const setTraits = (list: Trait[]) => {
    traits.value = dedupeTraits(list)
  }
  const addTrait = (trait: Trait) => {
    traits.value = dedupeTraits([...traits.value, trait])
  }
  const addTraits = (list: Trait[]) => {
    if (!Array.isArray(list) || !list.length) return
    traits.value = dedupeTraits([...traits.value, ...list])
  }
  const removeTrait = (t_uuid: string) => {
    traits.value = traits.value.filter(t => t.t_uuid !== t_uuid)
  }
  const clear = () => {
    traits.value = []
    currentUuid.value = null
    keyMetaById.value = {}
  }
  const setKeyMetaBulk = (items: TraitKey[]) => {
    if (!Array.isArray(items)) return
    const nextMeta = { ...keyMetaById.value }
    const metaBySyn: Record<string, TraitKey> = {}
    items.forEach(item => {
      if (!item || typeof item.id !== 'number') return
      nextMeta[item.id] = item
      if (item.syn) metaBySyn[item.syn.toLowerCase()] = item
    })
    keyMetaById.value = nextMeta
    if (Object.keys(metaBySyn).length === 0) return
    // const updated = traits.value.map(trait => {
    //   if (trait.t_key_id || !trait.t_key) return trait
    //   const match = metaBySyn[trait.t_key.toLowerCase()]
    //   return match ? { ...trait, t_key_id: match.id } : trait
    // })
    // traits.value = updated
  }
  return {
    traits,
    currentUuid,
    keyMetaById,
    setTraits,
    addTrait,
    addTraits,
    removeTrait,
    clear,
    setKeyMetaBulk
  }
})
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useTraitsStore, import.meta.hot))
}
