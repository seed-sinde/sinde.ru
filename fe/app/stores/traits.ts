const dedupeTraits = (list: Trait[]) => {
  if (!Array.isArray(list) || list.length === 0) return []
  const seen = new Set<string>()
  const normalized: Trait[] = []
  for (const item of list) {
    const uuid = String(item?.t_uuid || "").trim()
    if (!uuid || seen.has(uuid)) continue
    seen.add(uuid)
    normalized.push(item)
  }
  return normalized
}
export const useTraitsStore = defineStore("traits", () => {
  const STORAGE_KEY = "traits.store.v1"
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
      if (typeof parsed.currentUuid === "string" || parsed.currentUuid === null) {
        currentUuid.value = parsed.currentUuid
      }
      if (parsed.keyMetaById && typeof parsed.keyMetaById === "object") {
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
    watch([traits, currentUuid, keyMetaById], persistState, {deep: true})
  }
  const applyTraitsSnapshot = (list: Trait[]) => {
    const next = dedupeTraits(list)
    const previousByUuid = new Map<string, Trait>()
    for (const item of traits.value) {
      const uuid = String(item?.t_uuid || "").trim()
      if (!uuid) continue
      previousByUuid.set(uuid, item)
    }
    const merged: Trait[] = []
    for (const item of next) {
      const uuid = String(item?.t_uuid || "").trim()
      if (!uuid) continue
      const previous = previousByUuid.get(uuid)
      if (previous) {
        Object.assign(previous, item)
        merged.push(previous)
        continue
      }
      merged.push({...item})
    }
    traits.value.splice(0, traits.value.length, ...merged)
  }
  const setTraits = (list: Trait[]) => {
    applyTraitsSnapshot(list)
  }
  const addTrait = (trait: Trait) => {
    const uuid = String(trait?.t_uuid || "").trim()
    if (!uuid) return
    const existingIndex = traits.value.findIndex(item => String(item?.t_uuid || "").trim() === uuid)
    if (existingIndex >= 0) {
      const existing = traits.value[existingIndex]
      if (!existing) return
      Object.assign(existing, trait)
      return
    }
    traits.value.push({...trait})
  }
  const addTraits = (list: Trait[]) => {
    if (!Array.isArray(list) || !list.length) return
    for (const item of list) {
      addTrait(item)
    }
  }
  const removeTrait = (t_uuid: string) => {
    const target = String(t_uuid || "").trim()
    if (!target) return
    const index = traits.value.findIndex(item => String(item?.t_uuid || "").trim() === target)
    if (index < 0) return
    traits.value.splice(index, 1)
  }
  const removeTraits = (uuids: string[]) => {
    const targets = new Set(uuids.map(value => String(value || "").trim()).filter(Boolean))
    if (!targets.size) return
    for (let index = traits.value.length - 1; index >= 0; index -= 1) {
      const uuid = String(traits.value[index]?.t_uuid || "").trim()
      if (!targets.has(uuid)) continue
      traits.value.splice(index, 1)
    }
  }
  const replaceTrait = (previousUuid: string, nextTrait: Trait) => {
    const target = String(previousUuid || "").trim()
    if (!target) return
    const index = traits.value.findIndex(item => String(item?.t_uuid || "").trim() === target)
    if (index < 0) {
      addTrait(nextTrait)
      return
    }
    traits.value.splice(index, 1, {...nextTrait})
  }
  const clear = () => {
    traits.value.splice(0, traits.value.length)
    currentUuid.value = null
    keyMetaById.value = {}
  }
  const setKeyMetaBulk = (items: TraitKey[]) => {
    if (!Array.isArray(items) || !items.length) return
    const nextMeta = {...keyMetaById.value}
    for (const item of items) {
      if (!item || typeof item.id !== "number") continue
      nextMeta[item.id] = item
    }
    keyMetaById.value = nextMeta
  }
  return {
    traits,
    currentUuid,
    keyMetaById,
    setTraits,
    addTrait,
    addTraits,
    removeTrait,
    removeTraits,
    replaceTrait,
    clear,
    setKeyMetaBulk
  }
})
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useTraitsStore, import.meta.hot))
}
