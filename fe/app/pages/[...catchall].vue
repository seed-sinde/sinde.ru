<script setup lang="ts">
  const { localeCode } = useInterfacePreferences()
  const copy = computed(() => {
    const messages = {
      ru: {
        title: 'Страница не найдена',
        description: 'Запрошенная страница отсутствует.',
        headline: 'Хорошо что проверили, такой страницы нет.',
        text: 'Страница, которую вы ищете, не существует. Проверьте адрес URL.',
        home: 'На главную'
      },
      en: {
        title: 'Page Not Found',
        description: 'The requested page does not exist.',
        headline: 'Good catch. This page does not exist.',
        text: 'The page you are looking for does not exist. Check the URL and try again.',
        home: 'Go Home'
      },
      ch: {
        title: '页面未找到',
        description: '请求的页面不存在。',
        headline: '检查得好，这个页面并不存在。',
        text: '你要查找的页面不存在。请检查 URL 地址。',
        home: '返回首页'
      },
      jp: {
        title: 'ページが見つかりません',
        description: '指定されたページは存在しません。',
        headline: '確認できてよかったです。このページは存在しません。',
        text: 'お探しのページは存在しません。URL を確認してください。',
        home: 'ホームへ'
      }
    } as const
    return messages[localeCode.value] || messages.ru
  })
  definePageMeta({
    robots: 'noindex, nofollow'
  })
  usePageSeo({
    title: () => copy.value.title,
    description: () => copy.value.description
  })
  const goHome = () => navigateTo('/')
</script>
<template>
  <div>
    <LabNavHeader :title="copy.title" />
    <section class="p-4 space-y-4">
      <div class="flex items-center gap-4">
        <div class="text-4xl font-bold text-red-400">404</div>
        <div class="text-lg">{{ copy.headline }}</div>
      </div>
      <p>{{ copy.text }}</p>
      <LabBaseButton @click="goHome" :label="copy.home" variant="info" />
    </section>
  </div>
</template>
