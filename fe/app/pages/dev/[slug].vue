<script setup lang="ts">
  const route = useRoute()
  const slug = computed(() => {
    const param = route.params.slug
    return typeof param === 'string' ? param : ''
  })
  const { data: doc } = await useAsyncData(
    () => `development-doc:${slug.value}`,
    async () => {
      const item = await dev.loadDoc(slug.value)
      if (!item) {
        throw createError({ statusCode: 404, statusMessage: 'Раздел не найден' })
      }
      return item
    },
    { watch: [slug] }
  )
  const title = computed(() => doc.value?.title || 'разработка')
  usePageSeo({
    title,
    description: () => doc.value?.summary || 'Материал из раздела Разработка.'
  })
</script>
<template>
  <div class="space-y-4">
    <div class="space-y-2">
      <LabNavHeader
        :title
        :breadcrumb-items="[
          { label: 'Разработка', to: '/dev' },
          { label: doc?.title || 'Материал', current: true }
        ]" />
    </div>
    <section class="p-3 sm:p-5">
      <LabViewerMarkdown :source="doc?.content || ''" />
    </section>
  </div>
</template>
