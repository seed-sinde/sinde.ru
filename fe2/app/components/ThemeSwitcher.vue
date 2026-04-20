<script setup lang="ts">
const { mode, set } = useTheme();

const triggerIcon = computed(() =>
  mode.value === "light"
    ? "ic:round-light-mode"
    : mode.value === "dark"
      ? "ic:round-dark-mode"
      : "ic:round-auto-mode",
);

const toggle = () =>
  set(
    mode.value === "light"
      ? "dark"
      : mode.value === "dark"
        ? "system"
        : "light",
  );
</script>

<template>
  <Button :icon="triggerIcon" @click="toggle" />
</template>
