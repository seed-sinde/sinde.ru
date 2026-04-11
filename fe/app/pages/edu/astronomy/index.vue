<script setup lang="ts">
  const title = 'Астрономия'
  usePageSeo({
    title,
    description: 'Визуализация галактической орбиты Солнечной системы, траекторий планет и крупных спутников.'
  })
</script>

<template>
  <div>
    <LabNavHeader
      :title
      :breadcrumb-items="[
        { label: 'Вики', to: '/edu' },
        { label: title, current: true }
      ]" />
    <section>Astronomy</section>
  </div>
</template>
