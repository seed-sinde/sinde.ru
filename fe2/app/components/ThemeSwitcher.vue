<script setup lang="ts">
import IcBaselineLightMode from "~icons/ic/baseline-light-mode"
import IcBaselineDarkMode from "~icons/ic/baseline-dark-mode"
import IcBaselineAutoMode from "~icons/ic/baseline-auto-mode"
const {mode, toggle} = useTheme()
const icon = computed(() =>
  mode.value === "light" ? IcBaselineLightMode : mode.value === "dark" ? IcBaselineDarkMode : IcBaselineAutoMode
)
const modeIndex = computed(() => (mode.value === "light" ? 0 : mode.value === "dark" ? 1 : 2))
</script>
<template>
  <Button
    :icon="icon"
    icon-class="text-lg"
    :icon-tooltip="['Theme: Light', 'Theme: Dark', 'Theme: System']"
    :tooltip-index="modeIndex"
    aria-label="Toggle theme"
    @click="toggle"
  />
</template>
