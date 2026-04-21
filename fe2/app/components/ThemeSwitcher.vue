<script setup lang="ts">
const { mode, toggle } = useTheme();
</script>
<template>
  <Button :label="`theme: ${mode}`" @click="toggle" />
</template>
