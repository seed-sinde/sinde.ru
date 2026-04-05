<script setup lang="ts">
  definePageMeta({
    title: 'химия',
    description: 'Раздел химии с переходом к таблице элементов и каталогу минералов.'
  })
  const dataWikiTools = wikiTools.filter(tool => tool.category === 'chemistry')
</script>
<template>
  <div class="space-y-4">
    <LabNavHeader
      title="Химия"
      :breadcrumb-items="[
        { label: 'Вики', to: '/edu' },
        { label: 'Химия', current: true }
      ]" />
    <section class="grid gap-3 px-3 sm:px-6 lg:grid-cols-2">
      <EduWikiToolCard v-for="link in dataWikiTools" :key="link.to" :tool="link" />
    </section>
  </div>
</template>
