<script setup lang="ts">
  const themes = ['dark', 'light'] as const
  const runtimeConfig = useRuntimeConfig()
  const requestUrl = useRequestURL()
  const { localeTag, effectiveTheme, themeColor, faviconSrc, faviconDarkSrc, faviconLightSrc, themePreference } =
    useInterfacePreferences()
  const siteUrl = computed(() => {
    const configured = String(runtimeConfig.public.baseURL || '').trim()
    if (configured) return configured.replace(/\/+$/, '')
    return `${requestUrl.protocol}//${requestUrl.host}`.replace(/\/+$/, '')
  })
  const defaultOgImage = computed(() => `${siteUrl.value}${faviconDarkSrc.value}`)
  const htmlThemeClass = computed(() => `theme-${effectiveTheme.value}`)
  const faviconLinks = computed(() => [
    ...(themePreference.value === 'system' ? [] : [{ rel: 'icon', href: faviconSrc.value }]),
    ...themes.map(t => ({
      rel: 'icon',
      href: t === 'dark' ? faviconDarkSrc.value : faviconLightSrc.value,
      media: `(prefers-color-scheme: ${t})`
    })),
    {
      rel: 'apple-touch-icon',
      href: '/pwa-180.png'
    }
  ])
  useHead({
    titleTemplate: t => (t ? `${t} · sinde` : 'sinde'),
    htmlAttrs: {
      lang: localeTag,
      class: htmlThemeClass
    },
    link: faviconLinks,
    meta: [
      { key: 'referrer', name: 'referrer', content: 'strict-origin-when-cross-origin' },
      {
        key: 'viewport',
        name: 'viewport',
        content: 'width=device-width, initial-scale=1, viewport-fit=cover'
      },
      { key: 'theme-color', name: 'theme-color', content: themeColor },
      { key: 'mobile-web-app-capable', name: 'mobile-web-app-capable', content: 'yes' },
      { key: 'apple-mobile-web-app-capable', name: 'apple-mobile-web-app-capable', content: 'yes' },
      {
        key: 'apple-mobile-web-app-status-bar-style',
        name: 'apple-mobile-web-app-status-bar-style',
        content: 'black-translucent'
      },
      { key: 'apple-mobile-web-app-title', name: 'apple-mobile-web-app-title', content: 'sinde' },
      { key: 'twitter:card', name: 'twitter:card', content: 'summary' },
      { key: 'og:site_name', property: 'og:site_name', content: 'sinde' },
      { key: 'og:image', property: 'og:image', content: defaultOgImage },
      { key: 'twitter:image', name: 'twitter:image', content: defaultOgImage }
    ]
  })
</script>
<template>
  <VitePwaManifest />
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
