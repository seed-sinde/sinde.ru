<script setup lang="ts">
  const route = useRoute()
  const slug = computed(() => {
    const param = route.params.slug
    return typeof param === 'string' ? param : ''
  })
  const { data: doc } = await useAsyncData(
    () => `edu-doc:${slug.value}`,
    async () => {
      const item = await edu.loadDoc(slug.value)
      if (!item) {
        throw createError({ statusCode: 404, statusMessage: 'Раздел не найден' })
      }
      return item
    },
    { watch: [slug] }
  )
  const title = computed(() => doc.value?.title || 'вики')
  usePageSeo({
    title,
    description: () => doc.value?.summary || 'Материал из раздела Вики.'
  })
</script>
<template>
  <div class="space-y-2">
    <LabNavHeader
      :title
      :breadcrumb-items="[
        { label: 'Вики', to: '/edu' },
        { label: title || 'Материал', current: true }
      ]" />
    <section class="bg-zinc-900/45 px-3 sm:px-5">
      <LabViewerMarkdown :source="doc?.content || ''" />
    </section>
  </div>
</template>
