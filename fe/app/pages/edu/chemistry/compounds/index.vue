<script setup lang="ts">
  const title = 'Соединения'
  usePageSeo({
    title,
    description: 'Заготовка страницы для каталога химических соединений.'
  })
</script>
<template>
  <div>
    <LabNavHeader
      :title="title"
      :breadcrumb-items="[
        { label: 'Вики', to: '/edu' },
        { label: 'Химия', to: '/edu/chemistry' },
        { label: title, current: true }
      ]" />
    <section></section>
  </div>
</template>
