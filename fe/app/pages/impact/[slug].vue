<script setup lang="ts">
  const route = useRoute()
  const slug = computed(() => {
    const param = route.params.slug
    return typeof param === 'string' ? param : ''
  })
  const { data: doc } = await useAsyncData(
    () => `impact-doc:${slug.value}`,
    async () => {
      const item = await impact.loadDoc(slug.value)
      if (!item) {
        throw createError({ statusCode: 404, statusMessage: 'Раздел не найден' })
      }
      return item
    },
    { watch: [slug] }
  )
  const title = computed(() => doc.value?.title || 'влияние')
  usePageSeo({
    title,
    description: () => doc.value?.summary || 'Материал из раздела Влияние.'
  })
</script>
<template>
  <div class="space-y-4">
    <LabNavHeader
      :title
      :breadcrumb-items="[
        { label: 'Влияние', to: '/impact' },
        { label: title || 'Материал', current: true }
      ]" />
    <section class="p-3 sm:p-5">
      <LabViewerMarkdown :source="doc?.content || ''" />
    </section>
  </div>
</template>
