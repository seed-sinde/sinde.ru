<script setup lang="ts">
const title = 'Изотопы'
usePageSeo({
  title,
  description: 'Заготовка страницы для каталога изотопов.'
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
      ]"
    />
    <section />
  </div>
</template>
