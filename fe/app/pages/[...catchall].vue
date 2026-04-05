<script setup lang="ts">
  definePageMeta({
    robots: 'noindex, nofollow'
  })
  const title = 'страница не найдена'
  usePageSeo({
    title,
    description: 'Запрошенная страница отсутствует.'
  })
  const goHome = () => navigateTo('/')
</script>
<template>
  <div>
    <LabNavHeader :title />
    <div class="p-4">
      <div class="flex items-center gap-3 mb-4">
        <div class="text-4xl font-bold text-red-400">404</div>
        <div class="text-lg text-gray-200">Хорошо что проверили, такой страницы нет.</div>
      </div>
      <p class="text-gray-400 mb-6">Страница, которую вы ищете, не существует. Проверьте адрес URL.</p>
      <LabBaseButton @click="goHome" label="На главную" variant="info" />
    </div>
  </div>
</template>
