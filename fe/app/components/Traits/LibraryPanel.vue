<script setup lang="ts">
  const formatShortUuid = shortUuid
  const props = withDefaults(
    defineProps<{
      currentTraitUuid?: string | null
    }>(),
    {
      currentTraitUuid: ''
    }
  )
  const router = useRouter()
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
    } catch (error: any) {
      errorText.value =
        error?.data?.message || error?.statusMessage || error?.message || 'Не удалось загрузить сохранённые наборы.'
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
        infoText.value = 'Описание сохранённого набора обновлено.'
      } else {
        await saveTraitSet({
          set_uuid: currentTraitUuid.value,
          name: formName.value,
          description: formDescription.value
        })
        infoText.value = 'Набор сохранён в личный список.'
      }
      await loadSavedSets()
    } catch (error: any) {
      errorText.value = error?.data?.message || error?.statusMessage || error?.message || 'Не удалось сохранить набор.'
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
      syncUserPrimary(res.data.user.primary_trait_uuid)
      await loadSavedSets()
      infoText.value = 'Текущий набор теперь считается вашим основным.'
    } catch (error: any) {
      errorText.value =
        error?.data?.message || error?.statusMessage || error?.message || 'Не удалось обновить основной набор.'
    } finally {
      primaryPending.value = false
    }
  }
  const openPrimary = async () => {
    if (!primaryTraitUuid.value) return
    await router.replace(`/traits/${primaryTraitUuid.value}`)
  }
  const removeSavedSet = async (savedSetId: string) => {
    if (!savedSetId) return
    infoText.value = ''
    errorText.value = ''
    deletingId.value = savedSetId
    try {
      await deleteTraitSet(savedSetId)
      infoText.value = 'Набор удалён из личного списка.'
      await loadSavedSets()
    } catch (error: any) {
      errorText.value = error?.data?.message || error?.statusMessage || error?.message || 'Не удалось удалить набор.'
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
  <div class="px-3 py-3 sm:px-4 sm:py-4 lg:px-5">
    <div class="mx-auto max-w-6xl space-y-4">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div class="space-y-1">
          <h2 class="lab-text-primary text-base font-semibold tracking-tight sm:text-lg">Сохранённые наборы</h2>
          <p class="lab-text-muted max-w-2xl text-sm">
            Основной адрес хранится отдельно, а сохранённые наборы остаются приватными.
          </p>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <LabBaseButton
            v-if="hasPrimaryTrait && currentTraitUuid !== primaryTraitUuid"
            variant="secondary"
            size="sm"
            label="Открыть основной набор"
            @click="openPrimary" />
        </div>
      </div>
      <div class="space-y-2 h-4">
        <LabNotify :text="errorText" tone="error" />
        <LabNotify :text="infoText" tone="success" />
      </div>
    </div>
  </div>
  <div v-if="!ready || loading" class="mx-auto mt-4 max-w-6xl px-3 text-sm text-zinc-500 sm:px-4 lg:px-5">
    <LabLoader
      variant="cards"
      :count="2"
      :columns="2"
      label="Загрузка наборов..."
      :show-media="false"
      :show-footer="false" />
  </div>
  <div v-else-if="!isAuthenticated" class="mx-auto mt-4 max-w-6xl px-3 sm:px-4 lg:px-5">
    <div class="max-w-3xl">
      <AuthFeatureGateNotice
        message="Войдите в аккаунт, чтобы закрепить основной UUID и вести приватный список сохранённых наборов." />
    </div>
  </div>
  <div v-else class="px-3 pb-4 sm:px-4 lg:px-5">
    <div class="traits-library-layout mx-auto max-w-6xl gap-4">
      <article class="traits-library-card space-y-4 p-4">
        <div class="flex flex-wrap items-center gap-2">
          <span class="lab-text-muted text-sm">Основной набор:</span>
          <NuxtLink
            v-if="hasPrimaryTrait"
            :to="`/traits/${primaryTraitUuid}`"
            class="traits-library-link font-mono text-sm">
            {{ formatShortUuid(primaryTraitUuid, 6) }}
          </NuxtLink>
          <span v-else class="lab-text-muted text-sm">ещё не назначен</span>
          <LabBaseBadge v-if="isCurrentPrimary" variant="info">Открыт</LabBaseBadge>
        </div>
        <div class="flex flex-wrap gap-2">
          <LabBaseButton
            v-if="currentTraitUuid && !isCurrentPrimary"
            variant="secondary"
            size="sm"
            button-class="rounded-xl border border-amber-400/35 bg-amber-400/10 text-amber-100 transition-colors hover:bg-amber-400/18"
            :disabled="primaryPending"
            @click="setAsPrimary(currentTraitUuid)">
            Сделать основным набором
          </LabBaseButton>
        </div>
        <div class="space-y-3 border-t border-zinc-800/80 pt-4">
          <div class="space-y-1">
            <h3 class="lab-text-primary text-sm font-medium">
              {{ currentSavedSet ? 'Редактировать запись текущего набора' : 'Сохранить текущий набор' }}
            </h3>
            <p class="lab-text-muted text-xs">
              {{
                currentTraitUuid
                  ? formatShortUuid(currentTraitUuid, 6)
                  : 'Сначала откройте набор или отдельную особенность по UUID.'
              }}
            </p>
          </div>
          <div v-if="currentTraitUuid" class="space-y-3">
            <div class="max-w-2xl space-y-3">
              <LabField label="Название" for-id="traits-library-name">
                <LabBaseInput
                  id="traits-library-name"
                  v-model="formName"
                  placeholder="Например, Основной набор или Набор A"
                  maxlength="120" />
              </LabField>
              <LabField label="Короткое описание" for-id="traits-library-description" hint="До 280 символов.">
                <LabBaseTextarea
                  id="traits-library-description"
                  v-model="formDescription"
                  rows="3"
                  maxlength="280"
                  placeholder="Короткая заметка о наборе" />
              </LabField>
            </div>
            <div class="flex flex-col gap-2 sm:flex-row sm:items-center">
              <LabBaseButton
                variant="primary"
                size="sm"
                button-class="rounded-xl"
                :disabled="formPending || !formName.trim()"
                @click="saveCurrentSet">
                {{ currentSavedSet ? 'Обновить запись' : 'Сохранить набор' }}
              </LabBaseButton>
              <span v-if="currentSavedSet" class="lab-text-muted text-xs">
                Запись уже существует в вашем приватном списке.
              </span>
            </div>
          </div>
          <p v-else class="lab-text-muted text-sm">
            Здесь появится форма сохранения, как только у страницы будет конкретный UUID набора или особенности.
          </p>
        </div>
      </article>
      <article class="traits-library-card space-y-3 p-4">
        <div class="flex items-center justify-between gap-3">
          <h3 class="lab-text-primary text-sm font-medium">Сохранённые наборы</h3>
          <span class="lab-text-muted text-xs">Всего: {{ savedSets.length }}</span>
        </div>
        <div v-if="savedSets.length === 0" class="traits-library-empty px-4 py-3 text-sm">
          Пока нет ни одного сохранённого набора.
        </div>
        <div v-else class="space-y-2">
          <div v-for="item in savedSets" :key="item.saved_set_id" class="traits-library-item px-3 py-3">
            <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
              <div class="min-w-0 flex-1 space-y-1">
                <div class="flex flex-wrap items-center gap-2">
                  <NuxtLink
                    :to="`/traits/${item.set_uuid}`"
                    class="traits-library-link min-w-0 truncate text-sm font-medium">
                    {{ item.name }}
                  </NuxtLink>
                  <span
                    v-if="item.set_uuid === primaryTraitUuid"
                    class="inline-flex items-center rounded-full border border-emerald-500/30 bg-emerald-500/10 px-2 py-0.5 text-xs uppercase tracking-wide text-emerald-200">
                    Основной
                  </span>
                  <span
                    v-if="item.set_uuid === currentTraitUuid"
                    class="inline-flex items-center rounded-full border border-amber-400/30 bg-amber-400/10 px-2 py-0.5 text-xs uppercase tracking-wide text-amber-100">
                    Открыт
                  </span>
                </div>
                <p class="lab-text-muted font-mono text-xs">{{ formatShortUuid(item.set_uuid, 6) }}</p>
                <p v-if="item.description" class="lab-text-secondary text-sm">{{ item.description }}</p>
              </div>
              <div class="flex flex-wrap items-center gap-2">
                <LabBaseButton
                  v-if="item.set_uuid !== primaryTraitUuid"
                  variant="primary"
                  size="xs"
                  label="Сделать моим"
                  :disabled="primaryPending"
                  @click="setAsPrimary(item.set_uuid)" />
                <LabBaseButton
                  variant="danger"
                  size="xs"
                  label="Удалить"
                  :disabled="deletingId === item.saved_set_id"
                  @click="removeSavedSet(item.saved_set_id)" />
              </div>
            </div>
          </div>
        </div>
      </article>
    </div>
  </div>
</template>
<style scoped>
  .traits-library-layout {
    display: grid;
  }
  @media (min-width: 1280px) {
    .traits-library-layout {
      grid-template-columns: minmax(0, 1.05fr) minmax(18rem, 0.95fr);
    }
  }
  .traits-library-card {
    border: 1px solid color-mix(in srgb, var(--lab-border) 82%, transparent);
    background: color-mix(in srgb, var(--lab-bg-surface) 88%, transparent);
  }
  .traits-library-item {
    border: 1px solid color-mix(in srgb, var(--lab-border) 78%, transparent);
    background: color-mix(in srgb, var(--lab-bg-surface-subtle) 68%, transparent);
  }
  .traits-library-empty {
    border: 1px dashed color-mix(in srgb, var(--lab-border) 82%, transparent);
    background: color-mix(in srgb, var(--lab-bg-surface-subtle) 52%, transparent);
    color: var(--lab-text-muted);
  }
  .traits-library-link {
    color: var(--lab-text-primary);
    transition: color 0.15s ease;
  }
  .traits-library-link:hover,
  .traits-library-link:focus-visible {
    color: var(--lab-accent);
  }
</style>
