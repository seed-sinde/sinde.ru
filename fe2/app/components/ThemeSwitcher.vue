<script setup lang="ts">
import IcBaselineLightMode from "~icons/ic/baseline-light-mode"
import IcBaselineDarkMode from "~icons/ic/baseline-dark-mode"
import IcBaselineAutoMode from "~icons/ic/baseline-auto-mode"
const {mode, toggle} = useTheme()
const icons = {
  light: IcBaselineLightMode,
  dark: IcBaselineDarkMode,
  system: IcBaselineAutoMode
} as const
const icon = computed(() => icons[mode.value] ?? icons.system)
const tooltip = computed(
  () => `Theme: ${mode.value === "light" ? "Light" : mode.value === "dark" ? "Dark" : "System"}`
)
</script>
<template>
  <UiButton aria-label="Toggle theme" variant="ghost" @click="toggle">
    <template #left>
      <UiTooltip :text="tooltip">
        <component :is="icon" />
      </UiTooltip>
    </template>
  </UiButton>
</template>
