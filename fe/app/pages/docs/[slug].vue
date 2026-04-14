<script setup lang="ts">
definePageMeta({
  validate: (route) => isDocsSlug(route.params.slug)
})

const route = useRoute()
const { localeCode } = useInterfacePreferences()

const page = resolveLocalizedPage(docsPages, localeCode)

const activeSlug = computed<DocsSlug>(() => normalizeDocsSlug(route.params.slug) ?? DOCS_DEFAULT_SLUG)
const activeDoc = computed(() => page.value.documents[activeSlug.value])
const introItems = computed(() => activeDoc.value.introItems ?? [])
const tabItems = computed(() => getDocsTabItems(page.value.documents))
const tabRouteMap = getDocsTabRouteMap()

const isExternalHref = (href?: string): boolean => Boolean(href && /^(https?:|mailto:|tel:)/.test(href))

usePageSeo({
  title: () => activeDoc.value.title,
  description: () => activeDoc.value.description
})
</script>
<template>
  <div>
    <LabNavHeader :title="activeDoc.title" />
    <LabNavTabs :model-value="activeSlug" :items="tabItems" :route-to-map="tabRouteMap" :render-panels="false" />
    <section class="space-y-6 p-4">
      <div class="grid gap-6 xl:grid-cols-[minmax(0,18rem)_minmax(0,1fr)]">
        <section v-if="introItems.length" class="space-y-3">
          <div v-if="activeDoc.introTitle" class="text-sm tracking-[0.18em] text-(--lab-text-muted) uppercase">
            {{ activeDoc.introTitle }}
          </div>
          <div class="divide-y">
            <article v-for="item in introItems" :key="item.label" class="py-3">
              <p class="text-xs tracking-[0.16em] text-(--lab-text-muted) uppercase">
                {{ item.label }}
              </p>
              <a
                v-if="item.href && isExternalHref(item.href)"
                :href="item.href"
                class="lab-focus mt-2 inline-flex min-w-0 text-sm leading-6 sm:text-base"
              >
                <span class="wrap-break-word">{{ item.value }}</span>
              </a>
              <NuxtLink
                v-else-if="item.href"
                :to="item.href"
                class="lab-focus mt-2 inline-flex min-w-0 text-sm leading-6 sm:text-base"
              >
                <span class="wrap-break-word">{{ item.value }}</span>
              </NuxtLink>
              <p v-else class="mt-2 text-sm leading-6 sm:text-base">{{ item.value }}</p>
            </article>
          </div>
        </section>

        <section class="space-y-3" :class="introItems.length ? '' : 'xl:col-span-2'">
          <div v-if="activeDoc.sectionsTitle" class="text-sm tracking-[0.18em] text-(--lab-text-muted) uppercase">
            {{ activeDoc.sectionsTitle }}
          </div>
          <div class="divide-y">
            <article
              v-for="item in activeDoc.sections"
              :key="item.label"
              class="grid gap-2 py-3 sm:grid-cols-[12rem_minmax(0,1fr)] sm:gap-3"
            >
              <p class="text-xs tracking-[0.16em] text-(--lab-text-muted) uppercase">
                {{ item.label }}
              </p>
              <a
                v-if="item.href && isExternalHref(item.href)"
                :href="item.href"
                class="lab-focus text-sm leading-6 wrap-break-word sm:text-[0.95rem]"
              >
                {{ item.value }}
              </a>
              <NuxtLink
                v-else-if="item.href"
                :to="item.href"
                class="lab-focus text-sm leading-6 wrap-break-word sm:text-[0.95rem]"
              >
                {{ item.value }}
              </NuxtLink>
              <p v-else class="text-sm leading-6 wrap-break-word sm:text-[0.95rem]">
                {{ item.value }}
              </p>
            </article>
          </div>
        </section>
      </div>
    </section>
  </div>
</template>
