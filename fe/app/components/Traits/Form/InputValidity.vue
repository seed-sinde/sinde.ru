<template>
  <div class="space-y-2 bg-transparent p-0">
    <div class="flex flex-wrap items-center gap-2">
      <label
        for="TraitValidityModePermanent"
        class="inline-flex cursor-pointer items-center gap-2 rounded-full border border-zinc-700/80 bg-zinc-950/75 px-3 py-1.5 text-zinc-200 transition hover:border-amber-400/60 hover:text-amber-200"
      >
        <input
          id="TraitValidityModePermanent"
          name="TraitValidity"
          type="radio"
          value="permanent"
          v-model="mode"
          class="accent-amber-400"
        />
        <span>Постоянно</span>
      </label>
      <label
        for="TraitValidityModeTemporary"
        class="inline-flex cursor-pointer items-center gap-2 rounded-full border border-zinc-700/80 bg-zinc-950/75 px-3 py-1.5 text-zinc-200 transition hover:border-amber-400/60 hover:text-amber-200"
      >
        <input
          id="TraitValidityModeTemporary"
          name="TraitValidity"
          type="radio"
          value="temporary"
          v-model="mode"
          class="accent-amber-400"
        />
        <span>Временно</span>
      </label>
    </div>
    <div v-if="mode === 'permanent'" class="space-y-1">
      <LabField label="С" for-id="TraitValiditySince">
        <TraitsFormDateTextInput
          id="TraitValiditySince"
          v-model="since"
          :placeholder="placeholder"
          :invalid="Boolean(since) && !sinceValid"
        />
      </LabField>
      <LabErrorMessage v-if="since && !sinceValid" :text="`Неверный формат: ${placeholder}`" />
    </div>
    <div v-else class="space-y-1">
      <LabField label="Действует до" for-id="TraitValidityUntil">
        <TraitsFormDateTextInput
          id="TraitValidityUntil"
          v-model="until"
          :placeholder="placeholder"
          :invalid="Boolean(until) && (!untilValid || !isUntilNotPast)"
        />
      </LabField>
      <LabErrorMessage v-if="until && !untilValid" :text="`Неверный формат: ${placeholder}`" />
      <LabErrorMessage v-else-if="until && !isUntilNotPast" text='Дата "Действует до" не может быть меньше текущего времени.' />
    </div>
  </div>
</template>
<script setup lang="ts">
import TraitsFormDateTextInput from './DateTextInput.vue'
defineProps<{ meta?: unknown }>()
const model = defineModel<Validity>({ required: true })
const mode = ref<Validity['mode']>(model.value.mode ?? 'permanent')
const since = ref(model.value.since ?? '')
const until = ref(model.value.until ?? '')
const placeholder = datePlaceholder('datetime')
const nowDateTime = ref(currentDateText('datetime'))
let nowTimer: ReturnType<typeof setInterval> | null = null
const sinceValid = computed(() => !since.value || isValidDateText(since.value, 'datetime'))
const untilValid = computed(() => !until.value || isValidDateText(until.value, 'datetime'))
const isUntilNotPast = computed(() => {
  if (!until.value || !untilValid.value) return true
  const cmp = compareDateText(until.value, nowDateTime.value, 'datetime')
  return cmp !== null && cmp >= 0
})
watchEffect(() => {
  model.value = mode.value === 'permanent'
    ? { mode: 'permanent', since: since.value }
    : { mode: 'temporary', until: until.value }
})
watch(model, (next) => {
  mode.value = next.mode
  since.value = next.since ?? ''
  until.value = next.until ?? ''
})
onMounted(() => {
  nowTimer = setInterval(() => {
    nowDateTime.value = currentDateText('datetime')
  }, 30_000)
})
onBeforeUnmount(() => {
  if (!nowTimer) return
  clearInterval(nowTimer)
  nowTimer = null
})
</script>
