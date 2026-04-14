<script setup lang="ts">
const title = 'Производство'
usePageSeo({
  title,
  description: 'Заготовка страницы для производственных технологий.'
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
