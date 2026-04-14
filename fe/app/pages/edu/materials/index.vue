<script setup lang="ts">
const title = 'Инженерные материалы'
usePageSeo({
  title,
  description: 'Заготовка страницы для инженерных материалов.'
})
</script>
<template>
  <div>
    <LabNavHeader
      :title="title"
      :breadcrumb-items="[
        { label: 'Вики', to: '/edu' },
        { label: title, current: true }
      ]"
    />
    <section />
  </div>
</template>
