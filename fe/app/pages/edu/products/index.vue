<script setup lang="ts">
  const title = 'Изделия'
  usePageSeo({
    title,
    description: 'Заготовка страницы для инженерных изделий.'
  })
</script>
<template>
  <div>
    <LabNavHeader
      :title="title"
      :breadcrumb-items="[
        { label: 'Вики', to: '/edu' },
        { label: title, current: true }
      ]" />
    <section></section>
  </div>
</template>
