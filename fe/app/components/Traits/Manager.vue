<template>
  <section class="px-3 py-3 sm:px-4 sm:py-4 lg:px-5 lg:py-5">
    <div class="grid gap-4 xl:grid-cols-[minmax(18rem,24rem)_minmax(0,1fr)] xl:items-start">
      <aside class="space-y-3 xl:sticky xl:top-4">
        <TraitsFormAdd @add="onAdd" />
        <LabNotify
          :text="error"
          tone="error"
          as="div"
          class-name="border border-rose-500/30 bg-rose-500/10 px-3 py-2 text-sm text-rose-300" />
      </aside>
      <section class="space-y-3">
        <div class="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center">
          <LabBaseButton
            :label="isAllSelected ? 'Снять выбор' : 'Выбрать все'"
            variant="secondary"
            size="sm"
            @click="toggleSelectAll" />
          <LabBaseButton
            v-show="hasSelected || inProcess"
            icon="ic:round-edit"
            label="Изменить значение"
            variant="secondary"
            size="sm"
            :disabled="!canEditSelected"
            @click="startEditSelected" />
          <LabBaseButton
            icon="ic:round-delete"
            label="Удалить выбранное"
            variant="danger"
            size="sm"
            :disabled="!hasSelected || inProcess"
            @click="removeSelectedTraits(selectedIds)" />
          <div class="lab-text-muted text-xs sm:ml-auto">Всего: {{ traits.length }}</div>
        </div>
        <TraitsEditValueForm
          v-if="editingTrait"
          :trait="editingTrait"
          :meta="editingTraitMeta"
          :pending="inProcess"
          @save="onSaveEditedValue"
          @cancel="closeEditForm" />
        <div class="grid grid-cols-1 gap-2.5 md:grid-cols-2 xl:grid-cols-3">
          <TraitsCard
            v-for="trait in traits"
            :key="trait.t_uuid"
            :trait="trait"
            :selected="Boolean(selectedUuids[trait.t_uuid])"
            @update:selected="selectedUuids[trait.t_uuid] = $event" />
        </div>
        <div
          v-if="traits.length === 0"
          class="border border-dashed border-(--lab-border) px-4 py-4 text-sm italic text-(--lab-text-muted)">
          Список особенностей пустой.
        </div>
      </section>
    </div>
  </section>
</template>
<script setup lang="ts">
  const store = useTraitsStore()
  const route = useRoute()
  const { uuid, skipFetchUuid } = useTraitNavigation()
  watchEffect(() => {
    if (!uuid.value && /^\/traits\/?$/.test(route.path)) {
      store.clear()
    }
  })
  const traits = computed(() => store.traits)
  const {
    selectedMap: selectedUuids,
    hasSelected,
    isAllSelected,
    toggleSelectAll,
    clear: clearSelection,
    selectedIds
  } = useSelectionMap(traits)
  const { addTrait: onAdd, removeSelected: removeSelectedTraits, editTraitValue, inProcess, error } = useTraitActions()
  const editingTraitUuid = ref<string | null>(null)
  const { data: traitsData } = await useTraitLoader(uuid, skipFetchUuid)
  const hasResolvedStoreTraits = () => {
    return Boolean(uuid.value) && store.currentUuid === uuid.value && store.traits.length > 0
  }
  const applyResolvedTraits = (items?: Trait[] | null) => {
    if (!uuid.value) return
    if ((!Array.isArray(items) || items.length === 0) && hasResolvedStoreTraits()) {
      return
    }
    store.currentUuid = uuid.value
    store.setTraits(Array.isArray(items) ? items : [])
  }
  applyResolvedTraits(traitsData.value)
  watch(
    () => traitsData.value,
    next => {
      applyResolvedTraits(next)
    }
  )
  const selectedSingleTraitUuid = computed(() => (selectedIds.value.length === 1 ? selectedIds.value[0] : ''))
  const canEditSelected = computed(() => Boolean(selectedSingleTraitUuid.value) && !inProcess.value)
  const editingTrait = computed(() => {
    const targetUuid = String(editingTraitUuid.value || '').trim()
    if (!targetUuid) return null
    return traits.value.find(trait => trait.t_uuid === targetUuid) || null
  })
  const findMetaBySyn = (syn: string): TraitKey | null => {
    const targetSyn = String(syn || '')
      .trim()
      .toLowerCase()
    if (!targetSyn) return null
    const entries = Object.values(store.keyMetaById || {})
    const found = entries.find(
      entry =>
        String(entry?.syn || '')
          .trim()
          .toLowerCase() === targetSyn
    )
    return found || null
  }
  const ensureMetaForTrait = async (trait: Trait | null) => {
    if (!trait) return
    if (trait.t_key_id && store.keyMetaById?.[trait.t_key_id]) return
    if (findMetaBySyn(trait.t_key)) return
    try {
      const res = await getKeyMeta(trait.t_key)
      if (res?.data) store.setKeyMetaBulk([res.data])
    } catch {
      // optional meta preload: ignore errors and fallback to string input
    }
  }
  const editingTraitMeta = computed<KeyMeta>(() => {
    const trait = editingTrait.value
    if (!trait) return { dataType: 'string' }
    if (trait.t_key_id) {
      const byId = store.keyMetaById?.[trait.t_key_id]
      if (byId?.meta) return byId.meta as KeyMeta
    }
    const bySyn = findMetaBySyn(trait.t_key)
    if (bySyn?.meta) return bySyn.meta as KeyMeta
    return { dataType: 'string' }
  })
  const startEditSelected = async () => {
    if (!selectedSingleTraitUuid.value) return
    editingTraitUuid.value = selectedSingleTraitUuid.value
    await ensureMetaForTrait(editingTrait.value)
  }
  const closeEditForm = () => {
    editingTraitUuid.value = null
  }
  const onSaveEditedValue = async (payload: { traitUuid: string; t_key: string; t_value: string }) => {
    const nextTraitUuid = await editTraitValue({
      traitUuid: payload.traitUuid,
      t_value: payload.t_value
    })
    if (!nextTraitUuid) return
    editingTraitUuid.value = null
    clearSelection()
    selectedUuids[nextTraitUuid] = true
  }
  watch(traits, nextTraits => {
    if (!editingTraitUuid.value) return
    if (!nextTraits.some(trait => trait.t_uuid === editingTraitUuid.value)) {
      editingTraitUuid.value = null
    }
  })
  const collectMetaTargets = (sourceTraits: Trait[]) => {
    const ids = Array.from(
      new Set(sourceTraits.map(t => t?.t_key_id).filter((v): v is number => typeof v === 'number' && v > 0))
    )
    const keysWithoutId = Array.from(
      new Set(
        sourceTraits
          .filter(t => !t?.t_key_id && typeof t?.t_key === 'string')
          .map(t => String(t.t_key).trim())
          .filter(Boolean)
      )
    )
    return { ids, keysWithoutId }
  }
  const hasMetaForTrait = (trait: Trait) => {
    const id = Number(trait.t_key_id || 0)
    if (id > 0 && store.keyMetaById[id]) return true
    const syn = String(trait.t_key || '')
      .trim()
      .toLowerCase()
    if (!syn) return true
    return Object.values(store.keyMetaById).some(
      item =>
        String(item?.syn || '')
          .trim()
          .toLowerCase() === syn
    )
  }
  const fetchKeyMeta = async (sourceTraits: Trait[]): Promise<TraitKey[]> => {
    const unresolvedTraits = sourceTraits.filter(trait => !hasMetaForTrait(trait))
    if (!unresolvedTraits.length) return []
    const { ids, keysWithoutId } = collectMetaTargets(sourceTraits)
    if (!ids.length && !keysWithoutId.length) return []
    try {
      const [bulkRes, keyMetaResults] = await Promise.all([
        ids.length
          ? getKeysMetaBulk(ids)
          : Promise.resolve({ data: { items: [] as TraitKey[] } } as ApiResponseWithData<{ items: TraitKey[] }>),
        Promise.all(
          keysWithoutId.map(async syn => {
            try {
              const res = await getKeyMeta(syn)
              return res?.data || null
            } catch {
              return null
            }
          })
        )
      ])
      const byId = new Map<number, TraitKey>()
      for (const item of bulkRes?.data?.items || []) {
        if (item?.id) byId.set(item.id, item)
      }
      for (const item of keyMetaResults) {
        if (item?.id) byId.set(item.id, item)
      }
      return Array.from(byId.values())
    } catch (e) {
      console.warn('meta preload failed', e)
      return []
    }
  }
  const metaTargets = computed(() => collectMetaTargets(traits.value))
  const metaCacheKey = computed(() => {
    const idPart = metaTargets.value.ids.join(',')
    const keyPart = metaTargets.value.keysWithoutId.join(',')
    return `traits-meta:${uuid.value || ''}:${idPart}:${keyPart}`
  })
  const { data: keyMetaData } = await useAsyncData<TraitKey[]>(
    () => metaCacheKey.value,
    async () => await fetchKeyMeta(traits.value),
    {
      server: true,
      lazy: false,
      default: () => [],
      watch: [metaCacheKey]
    }
  )
  const applyMeta = (items?: TraitKey[] | null) => {
    if (!Array.isArray(items) || items.length === 0) return
    store.setKeyMetaBulk(items)
  }
  applyMeta(keyMetaData.value)
  watch(
    () => keyMetaData.value,
    next => {
      applyMeta(next)
    }
  )
</script>
