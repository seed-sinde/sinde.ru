<script setup lang="ts">
const { themePreference } = useInterfacePreferences()
const { locale, key, load, t } = useI18nSection('ui')
await useAsyncData(key.value, load, { watch: [locale] })
const { saveInterfacePreferences } = useInterfacePreferencesSync()
const hydrated = ref(false)
onMounted(() => {
  hydrated.value = true
})
const themeButtons = computed(() =>
  INTERFACE_THEME_OPTIONS.map(value => ({
    value,
    label: t(`theme.${value}`)
  }))
)
const pendingTheme = ref<ThemePreference | null>(null)
const open = ref(false)
const displayedTheme = computed<ThemePreference>(() => {
  if (pendingTheme.value) return pendingTheme.value
  if (!hydrated.value) return 'system'
  return themePreference.value
})
const triggerTitle = computed(() => `${t('theme.label')}: ${t(`theme.${displayedTheme.value}`)}`)
const triggerIcon = computed(() => {
  switch (displayedTheme.value) {
    case 'light':
      return 'ic:round-light-mode'
    case 'dark':
      return 'ic:round-dark-mode'
    default:
      return 'ic:round-auto-mode'
  }
})
const selectTheme = async (value: ThemePreference) => {
  if (pendingTheme.value || themePreference.value === value) return
  open.value = false
  pendingTheme.value = value
  await saveInterfacePreferences({ theme: value }).finally(() => (pendingTheme.value = null))
}
</script>
<template>
  <LabDropdown v-model="open" side="top" width-class="w-40" :match-trigger-width="false">
    <template #trigger="{ toggle }">
      <LabBaseButton
        :key="triggerIcon"
        :aria-label="triggerTitle"
        :disabled="Boolean(pendingTheme)"
        :icon="triggerIcon"
        :title="triggerTitle"
        icon-only
        icon-size="sm"
        size="sm"
        variant="ghost"
        button-class="h-8 w-8 rounded-full p-0"
        @click="toggle"
      />
    </template>
    <button
      v-for="option in themeButtons"
      :key="option.value"
      type="button"
      class="lab-dropdown-option"
      :class="displayedTheme === option.value ? 'lab-dropdown-option-active' : ''"
      :disabled="Boolean(pendingTheme)"
      @click="selectTheme(option.value)"
    >
      <span class="min-w-0 truncate">{{ option.label }}</span>
      <Icon
        v-if="displayedTheme === option.value"
        name="ic:round-check"
        class="lab-dropdown-option-meta shrink-0 text-sm"
      />
    </button>
  </LabDropdown>
</template>
