<template>
  <section class="grid grid-cols-1 items-start gap-3 xl:grid-cols-[22rem_minmax(0,1fr)]">
    <aside class="min-w-0 space-y-2">
      <TraitsFormAdd @add="onAdd" />
      <LabNotify
        :text="error"
        tone="error"
        as="div"
        class-name="px-3 py-2 text-sm text-rose-300"
      />
    </aside>
    <section class="min-w-0 space-y-2">
      <div class="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center">
        <LabBaseButton
          :label="isAllSelected ? t('manager.clear_selection') : t('manager.select_all')"
          variant="secondary"
          size="sm"
          @click="toggleSelectAll"
        />
        <LabBaseButton
          v-show="hasSelected || inProcess"
          icon="ic:round-edit"
          :label="t('manager.edit_value')"
          variant="secondary"
          size="sm"
          :disabled="!canEditSelected"
          @click="startEditSelected"
        />
        <LabBaseButton
          icon="ic:round-delete"
          :label="t('manager.delete_selected')"
          variant="danger"
          size="sm"
          :disabled="!hasSelected || inProcess"
          @click="removeSelectedTraits(selectedIds)"
        />
        <div class="text-xs text-(--lab-text-muted) sm:ml-auto">
          {{ t('manager.total', { count: String(totalTraits) }) }}
        </div>
      </div>
      <LabLoader v-if="showInlineSyncLoader" variant="inline" :label="t('library.loading')" class="text-xs" />
      <div class="grid grid-cols-1 gap-1.5 md:grid-cols-2 2xl:grid-cols-3">
        <TraitsCard
          v-for="trait in resolvedTraits"
          :key="trait.t_uuid"
          :trait="trait"
          :selected="Boolean(selectedMap[trait.t_uuid])"
          :editing="editingTraitUuid === trait.t_uuid"
          :editing-meta="editingTraitMetaMap[trait.t_uuid] || null"
          :edit-pending="inProcess"
          :state-tone="traitVisualState[trait.t_uuid] || ''"
          @update:selected="selectedMap[trait.t_uuid] = $event"
          @request-edit="startEditForTraitUuid(trait.t_uuid)"
          @save-inline="onSaveEditedValue"
          @cancel-inline="closeEditForm"
        />
      </div>
      <div v-if="totalTraits === 0" class="p-4 text-sm text-(--lab-text-muted) italic">
        {{ t('manager.empty') }}
      </div>
    </section>
  </section>
</template>
<script setup lang="ts">
const { locale, key, load, t } = useI18nSection('traits')
await useAsyncData(key.value, load, { watch: [locale] })
const store = useTraitsStore()
const streamStatus = useTraitStreamStatus()
const { uuid, skipFetchUuid } = useTraitNavigation()
const { data: traitsData, pending: traitsPending } = await useTraitLoader(uuid, skipFetchUuid)
const { traits: storeTraits, currentUuid: storeCurrentUuid } = storeToRefs(store)
const resolvedTraits = computed<Trait[]>(() =>
  storeTraits.value.length > 0
    ? storeTraits.value
    : Array.isArray(traitsData.value) && traitsData.value.length > 0
      ? traitsData.value
      : []
)
const {
  selectedMap,
  hasSelected,
  isAllSelected,
  toggleSelectAll,
  clear: clearSelection,
  selectedIds
} = useSelectionMap(resolvedTraits)
const {
  addTrait: addTraitAction,
  removeSelected: removeSelectedTraitsAction,
  editTraitValue,
  inProcess,
  error
} = useTraitActions()
const editingTraitUuid = ref<string | null>(null)
const traitVisualState = reactive<Record<string, 'added' | 'updated' | 'removing' | ''>>({})
const editingTraitMetaMap = reactive<Record<string, KeyMeta>>({})
const removingIds = reactive(new Set<string>())
const visualStateTimers = new Map<string, ReturnType<typeof setTimeout>>()
const selectedSingleTraitUuid = computed(() => (selectedIds.value.length === 1 ? selectedIds.value[0] : ''))
const canEditSelected = computed(() => Boolean(selectedSingleTraitUuid.value) && !inProcess.value)
const isCurrentStreamPending = computed(
  () => Boolean(uuid.value) && streamStatus.value.uuid === uuid.value && streamStatus.value.pending
)
const totalTraits = computed(() => resolvedTraits.value.length)
const showInlineSyncLoader = computed(
  () => traitsPending.value || inProcess.value || removingIds.size > 0 || isCurrentStreamPending.value
)
const editingTrait = computed(() =>
  (u => (u ? resolvedTraits.value.find(t => t.t_uuid === u) || null : null))(
    String(editingTraitUuid.value || '').trim()
  )
)
const setTraitVisualState = (u: string, s: 'added' | 'updated' | 'removing' | '', ttl = 900) => {
  const id = String(u || '').trim()
  if (!id) return
  const t = visualStateTimers.get(id)
  if (t) clearTimeout(t)
  traitVisualState[id] = s
  if (!s) return
  const nt = setTimeout(() => {
    if (traitVisualState[id] === s) traitVisualState[id] = ''
    visualStateTimers.delete(id)
  }, ttl)
  visualStateTimers.set(id, nt)
}
const onAdd = async (p: TraitInput) => {
  const r = await addTraitAction(p)
  if (r?.traitUuid) setTraitVisualState(r.traitUuid, r.action === 'updated' ? 'updated' : 'added')
}
const resolveSelectedIds = (s: string[] | { value?: string[] }) =>
  Array.isArray(s) ? s : Array.isArray(s?.value) ? s.value : []
const removeSelectedTraits = async (s: string[] | { value?: string[] }) => {
  const ids = resolveSelectedIds(s)
    .map(v => String(v || '').trim())
    .filter(Boolean)
  if (!ids.length || inProcess.value) return
  ids.forEach(id => removingIds.add(id))
  await removeSelectedTraitsAction(ids)
  ids.forEach(id => {
    removingIds.delete(id)
    selectedMap[id] = false
  })
}
const editingTraitMeta = computed<KeyMeta>(() => {
  const t = editingTrait.value
  if (!t) return { dataType: 'string' }
  const m = t.t_key_id ? store.keyMetaById?.[t.t_key_id]?.meta : null
  return (m as KeyMeta) || { dataType: 'string' }
})
const startEditSelected = async () => {
  if (!selectedSingleTraitUuid.value) return
  await startEditForTraitUuid(selectedSingleTraitUuid.value)
}
const startEditForTraitUuid = async (u: string) => {
  const id = String(u || '').trim()
  if (!id) return
  editingTraitUuid.value = id
  const t = resolvedTraits.value.find(x => x.t_uuid === id) || null
  if (t) {
    editingTraitMetaMap[id] = editingTraitMeta.value
  }
}
const closeEditForm = () => (editingTraitUuid.value = null)
const onSaveEditedValue = async (p: { traitUuid: string; t_key: string; t_value: string }) => {
  const r = await editTraitValue({ traitUuid: p.traitUuid, t_value: p.t_value })
  if (!r?.nextTraitUuid) return
  editingTraitUuid.value = null
  clearSelection()
  selectedMap[r.nextTraitUuid] = true
  setTraitVisualState(r.nextTraitUuid, 'updated')
}
watch(resolvedTraits, n => {
  if (editingTraitUuid.value && !n.some(t => t.t_uuid === editingTraitUuid.value)) {
    editingTraitUuid.value = null
  }
})
watch(
  resolvedTraits,
  n => {
    const a = new Set(n.map(t => String(t.t_uuid || '').trim()))
    Object.keys(traitVisualState).forEach(u => !a.has(u) && (traitVisualState[u] = ''))
  },
  { deep: true }
)
const collectMetaTargets = (s: Trait[]) => {
  return {
    ids: [...new Set(s.map(t => t?.t_key_id).filter((v): v is number => typeof v === 'number' && v > 0))]
  }
}
const fetchKeyMeta = async (s: Trait[]): Promise<TraitKey[]> => {
  const ids = collectMetaTargets(s).ids.filter(id => !store.keyMetaById[id])
  if (!ids.length) return []
  try {
    const res = await getKeysMetaBulk(ids)
    return res?.data?.items || []
  } catch (e) {
    console.warn('meta preload failed', e)
    return []
  }
}
const metaTargets = computed(() => collectMetaTargets(resolvedTraits.value))
const metaCacheKey = computed(() => `traits-meta:${uuid.value || ''}:${metaTargets.value.ids.join(',')}`)
const canFetchMeta = computed(
  () =>
    !!uuid.value &&
    resolvedTraits.value.length > 0 &&
    storeCurrentUuid.value === uuid.value &&
    !isCurrentStreamPending.value
)
const metaFetchStateKey = computed(() => `${metaCacheKey.value}:${canFetchMeta.value ? 'ready' : 'pending'}`)
const { data: keyMetaData } = await useAsyncData<TraitKey[]>(
  () => metaFetchStateKey.value,
  () => (canFetchMeta.value ? fetchKeyMeta(resolvedTraits.value) : Promise.resolve([])),
  {
    server: true,
    lazy: false,
    default: () => [],
    watch: [metaCacheKey, canFetchMeta]
  }
)
const applyMeta = (i?: TraitKey[] | null) => {
  if (!Array.isArray(i) || i.length === 0) return
  store.setKeyMetaBulk(i)
}
applyMeta(keyMetaData.value)
watch(
  () => keyMetaData.value,
  next => applyMeta(next)
)
onBeforeUnmount(() => {
  for (const timer of visualStateTimers.values()) {
    clearTimeout(timer)
  }
  visualStateTimers.clear()
})
</script>
