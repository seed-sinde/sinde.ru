<script setup lang="ts">
  const title = 'Вики'
  const dataWikiTools = wikiTools
  usePageSeo({
    title,
    description: 'Сжатый справочник знаний и интерактивные инструменты.'
  })
</script>
<template>
  <div class="space-y-6">
    <LabNavHeader :title />
    <section class="space-y-3 px-3 sm:px-6">
      <div class="grid gap-3 lg:grid-cols-2 xl:grid-cols-3">
        <EduWikiToolCard v-for="tool in dataWikiTools" :key="tool.to" :tool />
      </div>
    </section>
  </div>
</template>
