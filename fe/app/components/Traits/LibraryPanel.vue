<script setup lang="ts">
const { locale, key, load, t } = useI18nSection('traits')
await useAsyncData(key.value, load, { watch: [locale] })
const formatShortUuid = shortUuid
const props = withDefaults(
  defineProps<{
    currentTraitUuid?: string | null
  }>(),
  {
    currentTraitUuid: ''
  }
)
const {
  user,
  ensureLoaded,
  isAuthenticated,
  listTraitSets,
  saveTraitSet,
  updateTraitSet,
  deleteTraitSet,
  setPrimaryTraitUuid
} = useAuth()
const currentTraitUuid = computed(() => String(props.currentTraitUuid || '').trim())
const savedSets = ref<SavedTraitSetView[]>([])
const loading = ref(false)
const ready = ref(false)
const formPending = ref(false)
const primaryPending = ref(false)
const deletingId = ref('')
const infoText = ref('')
const errorText = ref('')
const formName = ref('')
const formDescription = ref('')
const primaryTraitUuid = computed(() => String(user.value?.primary_trait_uuid || '').trim())
const currentSavedSet = computed(() => savedSets.value.find(item => item.set_uuid === currentTraitUuid.value) || null)
const isCurrentPrimary = computed(
  () => Boolean(currentTraitUuid.value) && currentTraitUuid.value === primaryTraitUuid.value
)
const hasPrimaryTrait = computed(() => Boolean(primaryTraitUuid.value))
const formatCountText = (template: string, count: number) => template.replace('{count}', String(count))
const readErrorMessage = (error: unknown, fallback: string) => {
  if (!error || typeof error !== 'object') return fallback
  const apiError = error as {
    data?: { message?: string }
    statusMessage?: string
    message?: string
  }
  return apiError.data?.message || apiError.statusMessage || apiError.message || fallback
}
const syncUserPrimary = (next: string | null | undefined) => {
  if (!user.value) return
  user.value = {
    ...user.value,
    primary_trait_uuid: next || null
  }
}
const hydrateCurrentForm = () => {
  formName.value = currentSavedSet.value?.name || ''
  formDescription.value = currentSavedSet.value?.description || ''
}
const loadSavedSets = async () => {
  ready.value = false
  errorText.value = ''
  try {
    await ensureLoaded()
    if (!isAuthenticated.value) {
      savedSets.value = []
      hydrateCurrentForm()
      return
    }
    loading.value = true
    const res = await listTraitSets()
    savedSets.value = Array.isArray(res.data.items) ? res.data.items : []
    syncUserPrimary(res.data.primary_trait_uuid)
    hydrateCurrentForm()
  } catch (error: unknown) {
    errorText.value = readErrorMessage(error, t('library.load_failed'))
  } finally {
    loading.value = false
    ready.value = true
  }
}
const saveCurrentSet = async () => {
  if (!currentTraitUuid.value) return
  infoText.value = ''
  errorText.value = ''
  formPending.value = true
  try {
    if (currentSavedSet.value) {
      await updateTraitSet(currentSavedSet.value.saved_set_id, {
        name: formName.value,
        description: formDescription.value
      })
      infoText.value = t('library.updated_description')
    } else {
      await saveTraitSet({
        set_uuid: currentTraitUuid.value,
        name: formName.value,
        description: formDescription.value
      })
      infoText.value = t('library.saved_to_list')
    }
    await loadSavedSets()
  } catch (error: unknown) {
    errorText.value = readErrorMessage(error, t('library.save_failed'))
  } finally {
    formPending.value = false
  }
}
const setAsPrimary = async (setUuid: string) => {
  if (!setUuid) return
  infoText.value = ''
  errorText.value = ''
  primaryPending.value = true
  try {
    const res = await setPrimaryTraitUuid(setUuid)
    syncUserPrimary(res.data.changes.primary_trait_uuid)
    await loadSavedSets()
    infoText.value = t('library.primary_updated')
  } catch (error: unknown) {
    errorText.value = readErrorMessage(error, t('library.primary_failed'))
  } finally {
    primaryPending.value = false
  }
}
const removeSavedSet = async (savedSetId: string) => {
  if (!savedSetId) return
  infoText.value = ''
  errorText.value = ''
  deletingId.value = savedSetId
  try {
    await deleteTraitSet(savedSetId)
    infoText.value = t('library.removed')
    await loadSavedSets()
  } catch (error: unknown) {
    errorText.value = readErrorMessage(error, t('library.delete_failed'))
  } finally {
    deletingId.value = ''
  }
}
watch(
  () => [currentTraitUuid.value, currentSavedSet.value?.saved_set_id || ''].join(':'),
  () => {
    hydrateCurrentForm()
  },
  { immediate: true }
)
watch(
  () => user.value?.user_id || '',
  () => {
    if (!import.meta.client) return
    loadSavedSets()
  },
  { immediate: true }
)
</script>
<template>
  <div>
    <LabNotify :text="errorText" tone="error" />
    <LabNotify :text="infoText" tone="success" />
    <LabLoader
      v-if="!ready || loading"
      variant="cards"
      :count="2"
      :columns="1"
      :label="t('library.loading')"
      :show-media="false"
      :show-footer="false"
      class="max-w-lg"
    />
    <AuthFeatureGateNotice v-else-if="!isAuthenticated" :message="t('library.feature_gate')" />
    <div v-else>
      <article class="space-y-2">
        <div class="flex flex-wrap items-center gap-2">
          {{ t('library.primary_label') }}
          <NuxtLink
            v-if="hasPrimaryTrait"
            :to="`/traits/${primaryTraitUuid}`"
            class="lab-focus font-mono text-sm hover:text-(--lab-accent)"
          >
            {{ formatShortUuid(primaryTraitUuid, 6) }}
          </NuxtLink>
          <span v-else class="text-sm text-(--lab-text-muted)">{{ t('library.not_assigned') }}</span>
          <LabBaseBadge v-if="isCurrentPrimary" variant="info">{{ t('library.opened_badge') }}</LabBaseBadge>
        </div>
        <LabBaseButton
          v-if="currentTraitUuid && !isCurrentPrimary"
          :label="t('library.make_primary')"
          variant="secondary"
          size="sm"
          class="text-(--lab-text-primary)"
          :disabled="primaryPending"
          @click="setAsPrimary(currentTraitUuid)"
        />
        <div class="space-y-2">
          <div class="space-y-1">
            <h3 class="text-sm font-medium">
              {{ currentSavedSet ? t('library.edit_current.title') : t('library.save_current.title') }}
            </h3>
            <p class="text-xs text-(--lab-text-muted)">
              {{ currentTraitUuid ? formatShortUuid(currentTraitUuid, 6) : t('library.open_set.hint') }}
            </p>
          </div>
          <div v-if="currentTraitUuid" class="space-y-3">
            <div class="max-w-2xl space-y-3">
              <LabBaseField :label="t('library.name.label')" for-id="traits-library-name">
                <LabBaseInput
                  id="traits-library-name"
                  v-model="formName"
                  :placeholder="t('library.name.placeholder')"
                  maxlength="120"
                />
              </LabBaseField>
              <LabBaseField
                :label="t('library.description.label')"
                for-id="traits-library-description"
                :hint="t('library.description.hint')"
              >
                <LabBaseTextarea
                  id="traits-library-description"
                  v-model="formDescription"
                  rows="3"
                  maxlength="280"
                  :placeholder="t('library.description.placeholder')"
                />
              </LabBaseField>
            </div>
            <div class="flex flex-col gap-2 sm:flex-row sm:items-center">
              <LabBaseButton
                variant="primary"
                size="sm"
                :disabled="formPending || !formName.trim()"
                @click="saveCurrentSet"
              >
                {{ currentSavedSet ? t('library.update_entry') : t('library.save_set') }}
              </LabBaseButton>
              <span v-if="currentSavedSet" class="text-xs text-(--lab-text-muted)">
                {{ t('library.existing_record') }}
              </span>
            </div>
          </div>
          <p v-else class="text-sm text-(--lab-text-muted)">
            {{ t('library.waiting_form_hint') }}
          </p>
        </div>
      </article>
      <article class="space-y-2">
        <div class="flex items-center justify-between gap-3">
          <h3 class="text-sm font-medium">{{ t('library.title') }}</h3>
          <span class="text-xs text-(--lab-text-muted)">
            {{
            formatCountText(t('library.total'), savedSets.length)
            }}
          </span>
        </div>
        <div v-if="savedSets.length === 0" class="px-4 py-3 text-sm text-(--lab-text-muted)">
          {{ t('library.empty') }}
        </div>
        <div v-else class="space-y-2">
          <div v-for="item in savedSets" :key="item.saved_set_id" class="px-2 py-2">
            <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
              <div class="min-w-0 flex-1 space-y-1">
                <div class="flex flex-wrap items-center gap-2">
                  <NuxtLink :to="`/traits/${item.set_uuid}`" class="lab-focus min-w-0 truncate text-sm font-medium">
                    {{ item.name }}
                  </NuxtLink>
                  <span
                    v-if="item.set_uuid === primaryTraitUuid"
                    class="inline-flex items-center px-1.5 py-0.5 text-[11px] tracking-wide text-(--lab-accent) uppercase"
                  >
                    {{ t('library.primary.chip') }}
                  </span>
                  <span
                    v-if="item.set_uuid === currentTraitUuid"
                    class="inline-flex items-center px-1.5 py-0.5 text-[11px] tracking-wide uppercase"
                  >
                    {{ t('library.current.chip') }}
                  </span>
                </div>
                <p class="font-mono text-xs text-(--lab-text-muted)">{{ formatShortUuid(item.set_uuid, 6) }}</p>
                <p v-if="item.description" class="text-sm text-(--lab-text-secondary)">{{ item.description }}</p>
              </div>
              <div class="flex flex-wrap items-center gap-2">
                <LabBaseButton
                  v-if="item.set_uuid !== primaryTraitUuid"
                  variant="primary"
                  size="xs"
                  :label="t('library.make_mine')"
                  :disabled="primaryPending"
                  @click="setAsPrimary(item.set_uuid)"
                />
                <LabBaseButton
                  variant="danger"
                  size="xs"
                  :label="t('library.delete')"
                  :disabled="deletingId === item.saved_set_id"
                  @click="removeSavedSet(item.saved_set_id)"
                />
              </div>
            </div>
          </div>
        </div>
      </article>
    </div>
  </div>
</template>
