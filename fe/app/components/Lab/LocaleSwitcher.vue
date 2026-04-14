<script setup lang="ts">
const { localeCode, t } = useInterfacePreferences()
const { saveInterfacePreferences } = useInterfacePreferencesSync()
const localeButtons = computed(() =>
  INTERFACE_LOCALE_OPTIONS.map((option) => ({
    ...option,
    label: t(`locale.${option.code}`)
  }))
)
const pendingLocale = ref<InterfaceLocaleCode | null>(null)
const open = ref(false)
const selectedLocale = computed<InterfaceLocaleCode>(() => pendingLocale.value || localeCode.value)
const selectedLocaleBadge = computed(
  () => INTERFACE_LOCALE_OPTIONS.find((option) => option.code === selectedLocale.value)?.shortLabel || 'RU'
)
const getLocaleMessageKey = (code: InterfaceLocaleCode): InterfaceMessageKey => `locale.${code}` as InterfaceMessageKey
const triggerTitle = computed(() => `${t('locale.label')}: ${t(getLocaleMessageKey(selectedLocale.value))}`)
const selectLocale = async (code: InterfaceLocaleCode) => {
  if (pendingLocale.value || selectedLocale.value === code) return
  open.value = false
  pendingLocale.value = code
  try {
    await saveInterfacePreferences({ locale: code })
  } catch {
    // Rollback is handled by the shared sync composable.
  } finally {
    pendingLocale.value = null
  }
}
</script>
<template>
  <LabDropdown v-model="open" side="top" width-class="w-44" :match-trigger-width="false">
    <template #trigger="{ toggle }">
      <div class="relative inline-flex">
        <LabBaseButton
          :aria-label="triggerTitle"
          :disabled="Boolean(pendingLocale)"
          :title="triggerTitle"
          icon="ic:round-language"
          icon-only
          icon-size="sm"
          size="sm"
          variant="ghost"
          button-class="h-8 w-8 rounded-full p-0"
          @click="toggle"
        />
        <span
          class="h-2.8 pointer-events-none absolute top-1.5 -right-0.5 inline-flex min-w-1 items-center justify-center border border-zinc-600 bg-zinc-900 px-0.5 text-[8px] leading-none text-zinc-200"
        >
          {{ selectedLocaleBadge }}
        </span>
      </div>
    </template>
    <button
      v-for="option in localeButtons"
      :key="option.code"
      type="button"
      class="lab-dropdown-option"
      :class="selectedLocale === option.code ? 'lab-dropdown-option-active' : ''"
      :disabled="Boolean(pendingLocale)"
      @click="selectLocale(option.code)"
    >
      <span class="min-w-0 truncate">{{ option.label }}</span>
      <div class="flex items-center gap-2">
        <span class="lab-dropdown-option-meta text-xs">{{ option.shortLabel }}</span>
        <Icon
          v-if="selectedLocale === option.code"
          name="ic:round-check"
          class="lab-dropdown-option-meta shrink-0 text-sm"
        />
      </div>
    </button>
  </LabDropdown>
</template>
