<script setup lang="ts">
const title = 'Химия'
definePageMeta({
  title,
  description: 'Раздел химии с переходом к таблице элементов и каталогу минералов.'
})
const dataWikiTools = wikiTools.filter((tool) => tool.category === 'chemistry')
</script>
<template>
  <div>
    <LabNavHeader
      :title
      :breadcrumb-items="[
        { label: 'Вики', to: '/edu' },
        { label: 'Химия', current: true }
      ]"
    />
    <section class="grid grid-cols-1 gap-3 p-4 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4">
      <EduWikiToolCard v-for="link in dataWikiTools" :key="link.to" :tool="link" />
    </section>
  </div>
</template>
