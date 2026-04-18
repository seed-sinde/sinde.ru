<script setup lang="ts">
const { locale, key, load, t } = useI18nSection('ui')
await useAsyncData(key.value, load, { watch: [locale] })
const copy = computed(() => ({
  title: t('not_found.title'),
  description: t('not_found.description'),
  headline: t('not_found.headline'),
  text: t('not_found.text'),
  home: t('not_found.home')
}))
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
    <section class="space-y-4 p-4">
      <div class="flex items-center gap-4">
        <div class="text-4xl font-bold text-red-400">404</div>
        <div class="text-lg">{{ copy.headline }}</div>
      </div>
      <p>{{ copy.text }}</p>
      <LabBaseButton :label="copy.home" variant="info" @click="goHome" />
    </section>
  </div>
</template>
