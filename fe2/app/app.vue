<script setup lang="ts">
const {locale} = useI18nSection("common")
useHead({
  htmlAttrs: {
    lang: locale.value
  },
  link: [
    {
      key: "favicon",
      rel: "icon",
      type: "image/svg+xml",
      href: "/favicon.svg"
    },
    {
      key: "apple-touch-icon",
      rel: "apple-touch-icon",
      href: "/pwa-180.png"
    }
  ],
  meta: [
    {
      key: "referrer",
      name: "referrer",
      content: "strict-origin-when-cross-origin"
    },
    {
      key: "viewport",
      name: "viewport",
      content: "width=device-width, initial-scale=1, viewport-fit=cover"
    },
    {
      key: "theme-color",
      name: "theme-color",
      content: "#000"
    },
    {
      key: "mobile-web-app-capable",
      name: "mobile-web-app-capable",
      content: "yes"
    },
    {
      key: "apple-mobile-web-app-capable",
      name: "apple-mobile-web-app-capable",
      content: "yes"
    },
    {
      key: "apple-mobile-web-app-status-bar-style",
      name: "apple-mobile-web-app-status-bar-style",
      content: "black-translucent"
    },
    {
      key: "apple-mobile-web-app-title",
      name: "apple-mobile-web-app-title",
      content: "sinde.ru"
    },
    {key: "twitter:card", name: "twitter:card", content: "summary"},
    {key: "og:site_name", property: "og:site_name", content: "sinde.ru"},
    {key: "og:image", property: "og:image", content: "/favicon.svg"},
    {key: "twitter:image", name: "twitter:image", content: "/favicon.svg"}
  ]
})
const seoMeta = {
  title: "sinde.ru",
  description: "Коллективная разработка и проектирование"
}
useSeoMeta({
  title: seoMeta.title,
  description: seoMeta.description,
  ogTitle: seoMeta.title,
  ogDescription: seoMeta.description
})
</script>
<template>
  <VitePwaManifest />
  <NuxtLayout>
    <UiNav />
    <NuxtPage />
  </NuxtLayout>
</template>
