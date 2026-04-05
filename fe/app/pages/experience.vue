<script setup lang="ts">
  const title = 'Опыт'
  usePageSeo({
    title,
    description: 'Раздел о практическом опыте и накопленных подходах.'
  })
</script>
<template>
  <div>
    <LabNavHeader :title />
    <section class="p-4">
      <h2>{{ title }} - Раздел о практическом опыте и накопленных подходах.</h2>
    </section>
  </div>
</template>
