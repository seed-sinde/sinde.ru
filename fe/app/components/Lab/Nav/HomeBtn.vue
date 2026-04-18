<template>
  <NuxtLink
    to="/"
    class="lab-focus m-1 inline-flex h-8 w-8 items-center justify-center rounded-full text-(--lab-text-primary) transition-colors hover:bg-(--lab-bg-surface-hover) focus-visible:bg-(--lab-bg-surface-hover)"
    :aria-label="t('home')"
  >
    <picture>
      <source :srcset="faviconDarkSrc" media="(prefers-color-scheme: dark)" />
      <img :src="faviconLightSrc" alt="" class="h-5 w-5 object-contain" />
    </picture>
  </NuxtLink>
</template>

<script setup lang="ts">
const { faviconLightSrc, faviconDarkSrc } = useInterfacePreferences()
const { locale, key, load, t } = useI18nSection('nav')
await useAsyncData(key.value, load, { watch: [locale] })
</script>
