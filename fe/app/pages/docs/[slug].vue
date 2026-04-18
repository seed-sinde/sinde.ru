<script setup lang="ts">
definePageMeta({
  validate: route => isDocsSlug(route.params.slug)
})

const route = useRoute()
const { locale, key, load, t: docsT } = useI18nSection('docs')

await useAsyncData(key, load, { watch: [locale] })

const activeSlug = computed<DocsSlug>(() => normalizeDocsSlug(route.params.slug) ?? DOCS_DEFAULT_SLUG)

const resolveMessage = (messageKey: string) => docsT(messageKey)

const activeDoc = computed<DocContent>(() => {
  const slug = activeSlug.value
  const schema = DOCS_SCHEMA[slug]

  return {
    title: resolveMessage(`${slug}.title`),
    description: resolveMessage(`${slug}.description`) ,
    introTitle: resolveMessage(`${slug}.intro.title`) ,
    introItems: mapDocsItems(slug, 'intro', schema.introItems, resolveMessage),
    sectionsTitle: resolveMessage(`${slug}.sections.title`),
    sections: mapDocsItems(slug, 'sections', schema.sections, resolveMessage)
  }
})

const introItems = computed(() => activeDoc.value.introItems ?? [])

const tabLabels = computed<Record<DocsSlug, string>>(() =>
  Object.fromEntries(DOCS_SLUGS.map(slug => [slug, resolveMessage(`${slug}.title`)])) as Record<DocsSlug, string>
)

type DocsTabItem = {
  value: DocsSlug
  label: string
  badge?: string | number
}

const tabs = computed<DocsTabItem[]>(() =>
  getDocsTabItems(tabLabels.value).map(tab => ({
    value: tab.value as DocsSlug,
    label: tab.label,
    ...(tab.badge !== undefined && tab.badge !== null ? { badge: tab.badge } : {})
  }))
)
const docsTabRouteTargetMap = computed<TabRouteTargetMap>(() =>
  Object.fromEntries(DOCS_SLUGS.map(slug => [slug, getDocsHref(slug)]))
)
const activeTab = computed<DocsSlug>({
  get: () => activeSlug.value,
  set: () => {}
})
const breadcrumbItems = computed<BreadcrumbItem[]>(() => [
  { label: 'Документы', to: getDocsHref(DOCS_DEFAULT_SLUG) },
  { label: activeDoc.value.title, current: true, kind: 'tab' }
])
const isExternalHref = (href?: string) => Boolean(href && /^(https?:|mailto:|tel:)/.test(href))

usePageSeo({
  title: computed(() => activeDoc.value.title),
  description: computed(() => activeDoc.value.description)
})
</script>

<template>
  <div>
    <LabNavHeader :title="activeDoc.title" :breadcrumb-items="breadcrumbItems" />
    <LabNavTabs
      v-model="activeTab"
      :items="tabs"
      :render-panels="false"
      route-param-key="slug"
      :route-target-map="docsTabRouteTargetMap"
    />
    <section class="space-y-6 p-4">
      <div class="grid gap-6 md:grid-cols-[minmax(0,18rem)_minmax(0,1fr)]">
        <section v-if="introItems.length" class="space-y-3">
          <div v-if="activeDoc.introTitle" class="text-sm tracking-[0.18em] text-(--lab-text-muted) uppercase">
            {{ activeDoc.introTitle }}
          </div>
          <div class="divide-y divide-(--lab-border)">
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
              <p v-else class="mt-2 text-sm leading-6 sm:text-base">
                {{ item.value }}
              </p>
            </article>
          </div>
        </section>

        <section class="space-y-3" :class="introItems.length ? '' : 'xl:col-span-2'">
          <div v-if="activeDoc.sectionsTitle" class="text-sm tracking-[0.18em] text-(--lab-text-muted) uppercase">
            {{ activeDoc.sectionsTitle }}
          </div>
          <div class="divide-y divide-(--lab-border)">
            <article
              v-for="item in activeDoc.sections"
              :key="item.label"
              class="grid gap-2 py-3 xl:grid-cols-[12rem_minmax(0,1fr)] xl:gap-3"
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
