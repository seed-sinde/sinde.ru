<script setup lang="ts">
const NO_DATA = 'Нет данных'
const PHASE_LABELS: Record<string, string> = {
  Solid: 'Твёрдое',
  Liquid: 'Жидкое',
  Gas: 'Газообразное',
  Plasma: 'Плазма'
}
const route = useRoute()
const runtimeConfig = useRuntimeConfig()
const { data: periodicTableElementsData, error: chemistryElementsError } = await useChemistryElements()
if (chemistryElementsError.value) {
  throw createError({ statusCode: 500, statusMessage: 'Не удалось загрузить элементы' })
}
const periodic_table_element_route = (value: Parameters<typeof getPeriodicTableElementRoute>[0]) =>
  getPeriodicTableElementRoute(value)
const elementParam = computed(() => route.params.element)
const periodicTableElements = computed(() => periodicTableElementsData.value || [])
const element = computed(() => resolvePeriodicTableElement(periodicTableElements.value, elementParam.value))
if (!element.value) {
  throw createError({ statusCode: 404, statusMessage: 'Элемент не найден' })
}
const currentElement = computed(() => element.value!)
const orderedElements = computed(() =>
  periodicTableElements.value.slice().sort((left, right) => left.number - right.number)
)
const currentElementIndex = computed(() =>
  orderedElements.value.findIndex(item => item.number === currentElement.value.number)
)
const previousElement = computed(() => {
  const index = currentElementIndex.value
  return index > 0 ? orderedElements.value[index - 1] || null : null
})
const nextElement = computed(() => {
  const index = currentElementIndex.value
  return index >= 0 && index < orderedElements.value.length - 1 ? orderedElements.value[index + 1] || null : null
})
const imageViewerOpen = ref(false)
const imageViewerIndex = ref(0)
const selectedSampleImageIndex = ref(0)
const trimText = (value: unknown) => String(value ?? '').trim()
const formatText = (value: unknown, fallback = NO_DATA) => {
  const text = trimText(value)
  return text || fallback
}
const formatNumber = (value: unknown, fallback = NO_DATA) => {
  return typeof value === 'number' && Number.isFinite(value) ? String(value) : fallback
}
const formatNumberList = (values: number[] | null | undefined, fallback = NO_DATA) => {
  return Array.isArray(values) && values.length ? values.join(', ') : fallback
}
const formatPhase = (value: string | null | undefined) => {
  const phase = trimText(value)
  return phase ? PHASE_LABELS[phase] || phase : NO_DATA
}
const formatBlock = (value: string | null | undefined) => {
  const block = trimText(value)
  return block ? `${block.toUpperCase()}-блок` : NO_DATA
}
const formatKelvin = (value: unknown) => {
  return typeof value === 'number' && Number.isFinite(value) ? `${value} K` : NO_DATA
}
const formatMass = (value: unknown) => {
  return typeof value === 'number' && Number.isFinite(value) ? String(value) : NO_DATA
}
const formatDensity = (value: unknown) => {
  return typeof value === 'number' && Number.isFinite(value) ? String(value) : NO_DATA
}
const formatMolarHeat = (value: unknown) => {
  return typeof value === 'number' && Number.isFinite(value) ? String(value) : NO_DATA
}
const toLink = (value: string | null | undefined) => {
  const text = trimText(value)
  return /^(https?:)?\/\//i.test(text) ? text : null
}
const resolveAssetUrl = (value: string | null | undefined) => {
  const source = trimText(value)
  if (!source) return null
  if (/^(?:https?:)?\/\//i.test(source)) return source
  const base = String(runtimeConfig.app.baseURL || '/').replace(/\/+$/, '')
  const path = source.startsWith('/') ? source : `/${source}`
  return `${base}${path}` || path
}
const getContrastTextColor = (hexColor: string) => {
  const normalized = String(hexColor || '').replace('#', '')
  if (!/^[0-9a-f]{6}$/i.test(normalized)) {
    return '#ffffff'
  }
  const red = Number.parseInt(normalized.slice(0, 2), 16)
  const green = Number.parseInt(normalized.slice(2, 4), 16)
  const blue = Number.parseInt(normalized.slice(4, 6), 16)
  const luminance = (0.2126 * red + 0.7152 * green + 0.0722 * blue) / 255
  return luminance > 0.62 ? '#111827' : '#ffffff'
}
const formattedCurrentSymbol = computed(() => formatPeriodicTableSymbol(currentElement.value.symbol))
const resolvedSpectralImage = computed(() => resolveAssetUrl(currentElement.value.spectralImage))
const resolvedBohrModelImage = computed(() => resolveAssetUrl(currentElement.value.bohrModelImage))
const resolvedBohrModel3d = computed(() => resolveAssetUrl(currentElement.value.bohrModel3d))
const cpkHexLabel = computed(() => {
  const value = trimText(currentElement.value.cpkHex)
  return value ? `#${value.toUpperCase()}` : ''
})
const cpkChipStyle = computed(() => {
  if (!cpkHexLabel.value) return null
  return {
    backgroundColor: cpkHexLabel.value,
    color: getContrastTextColor(cpkHexLabel.value)
  }
})
const elementSymbolStyle = computed(() => ({
  borderColor: colorHexToRgba(currentElement.value.accentColor, 0.4),
  background: `color-mix(in srgb, ${currentElement.value.accentColor} 20%, var(--lab-bg-surface))`,
  color: getContrastTextColor(currentElement.value.accentColor)
}))
const headerDescription = computed(() => {
  return trimText(currentElement.value.summary) || `Химический элемент №${currentElement.value.number}.`
})
const heroStats = computed(() => {
  const item = currentElement.value
  return [
    { label: 'Атомная масса', value: formatMass(item.atomicMass) },
    { label: 'Конфигурация', value: formatText(item.electronConfiguration) },
    { label: 'Электроотрицательность', value: formatNumber(item.electronegativityPauling) },
    { label: 'Плотность', value: formatDensity(item.density) },
    { label: 'Плавление', value: formatKelvin(item.melt) },
    { label: 'Кипение', value: formatKelvin(item.boil) }
  ].filter(stat => stat.value !== NO_DATA)
})
const detailSections = computed<DetailSection[]>(() => {
  const item = currentElement.value
  return [
    {
      title: 'Основные данные',
      items: [
        { label: 'Период', value: formatNumber(item.period) },
        { label: 'Группа', value: formatNumber(item.group) },
        { label: 'Блок', value: formatBlock(item.block) },
        { label: 'Фаза', value: formatPhase(item.phase) },
        { label: 'Внешний вид', value: formatText(item.appearance) }
      ].filter(entry => entry.value !== NO_DATA)
    },
    {
      title: 'Физические свойства',
      items: [
        { label: 'Плотность', value: formatDensity(item.density) },
        { label: 'Плавление', value: formatKelvin(item.melt) },
        { label: 'Кипение', value: formatKelvin(item.boil) },
        { label: 'Молярная теплоёмкость', value: formatMolarHeat(item.molarHeat) }
      ].filter(entry => entry.value !== NO_DATA)
    },
    {
      title: 'Электронная структура',
      items: [
        { label: 'Электронная конфигурация', value: formatText(item.electronConfiguration) },
        { label: 'Семантическая конфигурация', value: formatText(item.electronConfigurationSemantic) },
        { label: 'Оболочки', value: formatNumberList(item.shells) },
        { label: 'Электроотрицательность', value: formatNumber(item.electronegativityPauling) },
        { label: 'Сродство к электрону', value: formatNumber(item.electronAffinity) },
        { label: 'Энергии ионизации', value: formatNumberList(item.ionizationEnergies) }
      ].filter(entry => entry.value !== NO_DATA)
    },
    {
      title: 'История и источники',
      items: [
        { label: 'Открыл', value: formatText(item.discoveredBy) },
        { label: 'Назвал', value: formatText(item.namedBy) },
        {
          label: 'Источник',
          value: formatText(item.source),
          href: toLink(item.source)
        }
      ].filter(entry => entry.value !== NO_DATA)
    },
    {
      title: 'Положение',
      items: [
        { label: 'X', value: formatNumber(item.xpos) },
        { label: 'Y', value: formatNumber(item.ypos) },
        { label: 'WX', value: formatNumber(item.wxpos) },
        { label: 'WY', value: formatNumber(item.wypos) }
      ].filter(entry => entry.value !== NO_DATA)
    }
  ].filter(section => section.items.length > 0)
})
const viewerItems = computed<ViewerImageItem[]>(() => {
  const item = currentElement.value
  const items: ViewerImageItem[] = []
  item.samples.forEach((sample, index) => {
    const src = resolveAssetUrl(sample.url) || sample.url
    items.push({
      key: `sample:${index}`,
      kind: 'sample',
      src,
      title: sample.title || `${item.russianName} · образец ${index + 1}`,
      alt: sample.title || `${item.russianName} sample ${index + 1}`,
      author: sample.author,
      attribution: sample.attribution,
      sourceUrl: sample.sourceUrl,
      licenseUrl: sample.licenseUrl,
      license: sample.license,
      thumbnailSrc: src
    })
  })
  if (resolvedBohrModelImage.value) {
    items.push({
      key: 'bohr',
      kind: 'bohr',
      src: resolvedBohrModelImage.value,
      title: `Модель Бора ${item.russianName}`,
      alt: `Модель Бора ${item.russianName}`,
      author: null,
      attribution: null,
      sourceUrl: null,
      licenseUrl: null,
      license: null,
      thumbnailSrc: null
    })
  }
  if (resolvedSpectralImage.value) {
    items.push({
      key: 'spectral',
      kind: 'spectral',
      src: resolvedSpectralImage.value,
      title: `Спектр ${item.russianName}`,
      alt: `Спектр ${item.russianName}`,
      author: null,
      attribution: null,
      sourceUrl: null,
      licenseUrl: null,
      license: null,
      thumbnailSrc: null
    })
  }
  return items
})
const spectralViewerIndex = computed(() => viewerItems.value.findIndex(item => item.kind === 'spectral'))
const sampleImages = computed<SamplePreviewItem[]>(() => {
  return viewerItems.value
    .map((item, index) => ({ ...item, viewerIndex: index }))
    .filter((item): item is SamplePreviewItem => item.kind === 'sample')
})
const selectedSampleImage = computed(() => {
  if (!sampleImages.value.length) return null
  return sampleImages.value[selectedSampleImageIndex.value] || sampleImages.value[0] || null
})
const relatedElements = computed(() => {
  const item = currentElement.value
  return periodicTableElements.value
    .filter(element => element.number !== item.number)
    .filter(
      element => element.group === item.group || element.period === item.period || element.category === item.category
    )
    .slice(0, 8)
})
const activeViewerTitle = computed(() => {
  return viewerItems.value[imageViewerIndex.value]?.title || currentElement.value.russianName
})
const openImageViewer = (index: number) => {
  if (index < 0 || !viewerItems.value[index]) return
  imageViewerIndex.value = index
  imageViewerOpen.value = true
}
const navigateToPreviousElement = () => {
  if (!previousElement.value) return
  void navigateTo(periodic_table_element_route(previousElement.value))
}
const navigateToNextElement = () => {
  if (!nextElement.value) return
  void navigateTo(periodic_table_element_route(nextElement.value))
}
const title = computed(() => `${currentElement.value.russianName} (${formattedCurrentSymbol.value})`)
usePageSeo({
  title,
  description: () => headerDescription.value
})
watch(
  () => currentElement.value.number,
  () => {
    selectedSampleImageIndex.value = 0
  },
  { immediate: true }
)
watch(
  sampleImages,
  items => {
    if (!items.length) {
      selectedSampleImageIndex.value = 0
      return
    }
    if (selectedSampleImageIndex.value > items.length - 1) {
      selectedSampleImageIndex.value = 0
    }
  },
  { immediate: true }
)
</script>
<template>
  <div>
    <LabNavHeader
      :title
      :breadcrumb-items="[
        { label: 'Вики', to: '/edu' },
        { label: 'Химия', to: '/edu/chemistry' },
        { label: 'Элементы', to: '/edu/chemistry/elements' },
        { label: currentElement.russianName, current: true }
      ]"
    />

    <section class="grid gap-4 p-4 xl:grid-cols-[22rem_minmax(0,1fr)]">
      <div class="min-w-0 space-y-4">
        <div v-if="resolvedBohrModel3d" class="aspect-square min-h-72">
          <LabViewer3D :src="resolvedBohrModel3d" :title="currentElement.russianName" class="h-full w-full" />
        </div>
        <section class="space-y-3">
          <div>
            <h2 class="text-sm font-semibold tracking-[0.07em] uppercase">Образцы вещества</h2>
            <p class="mt-1 text-sm text-(--lab-text-muted)">Фото и материалы рядом с основной информацией.</p>
          </div>
          <div v-if="selectedSampleImage" class="space-y-3">
            <LabViewerPreviewButton
              :src="selectedSampleImage.src"
              :alt="selectedSampleImage.alt"
              :label="selectedSampleImage.title || 'Открыть образец'"
              free-height
              class="w-full"
              frame-class="aspect-[16/9] max-h-none"
              image-class="h-full w-full object-cover"
              @preview="openImageViewer(selectedSampleImage.viewerIndex)"
            />
            <LabViewerImageThumbnails
              :items="sampleImages"
              :active-index="selectedSampleImageIndex"
              @select="selectedSampleImageIndex = $event"
            />
          </div>
          <section v-else-if="currentElement.sampleFallback" class="py-2">
            <div class="text-xs tracking-[0.07em] uppercase">
              {{ currentElement.sampleFallback.type }}
            </div>
            <div class="mt-2 text-sm font-semibold">
              {{ currentElement.sampleFallback.label }}
            </div>
            <p class="mt-2 text-sm leading-6">
              {{ currentElement.sampleFallback.description }}
            </p>
          </section>
          <div v-else class="py-2 text-sm text-(--lab-text-muted)">
            Для этого элемента изображения образцов пока не добавлено.
          </div>
        </section>
      </div>

      <div class="min-w-0 space-y-4">
        <div class="flex items-start justify-between gap-3">
          <div class="flex min-w-0 items-start gap-4">
            <div
              class="grid shrink-0 place-items-center text-4xl font-semibold sm:h-20 sm:w-20 sm:text-5xl"
              :style="elementSymbolStyle"
            >
              {{ formattedCurrentSymbol }}
            </div>
            <div class="min-w-0 pt-1">
              <h1 class="text-2xl font-semibold tracking-tight sm:text-3xl">
                {{ currentElement.russianName }}
              </h1>
              <p class="mt-1 text-sm text-(--lab-text-muted) sm:text-base">
                {{ currentElement.name }}
              </p>
            </div>
          </div>
          <div class="flex shrink-0 items-center gap-1">
            <LabBaseButton
              v-if="previousElement"
              icon="ic:round-navigate-before"
              :title="`Предыдущий элемент: ${previousElement.russianName}`"
              :aria-label="`Предыдущий элемент: ${previousElement.russianName}`"
              variant="ghost"
              size="sm"
              icon-only
              @click="navigateToPreviousElement"
            />
            <LabBaseButton
              v-if="nextElement"
              icon="ic:round-navigate-next"
              :title="`Следующий элемент: ${nextElement.russianName}`"
              :aria-label="`Следующий элемент: ${nextElement.russianName}`"
              variant="ghost"
              size="sm"
              icon-only
              @click="navigateToNextElement"
            />
          </div>
        </div>
        <p class="max-w-4xl text-sm leading-6 text-(--lab-text-secondary) sm:text-base">
          {{ headerDescription }}
        </p>
        <div class="flex flex-wrap gap-2 text-xs">
          <span class="px-2 py-1">
            {{ currentElement.categoryLabel }}
          </span>
          <span class="px-2 py-1">
            {{ formatPhase(currentElement.phase) }}
          </span>
          <span v-if="cpkChipStyle && cpkHexLabel" class="px-2 py-1 font-medium" :style="cpkChipStyle">
            CPK {{ cpkHexLabel }}
          </span>
        </div>
        <dl class="flex flex-wrap gap-2">
          <div v-for="stat in heroStats" :key="stat.label" class="max-w-full min-w-32 px-2 py-1.5">
            <dt class="text-[11px] tracking-[0.07em] text-(--lab-text-muted) uppercase">
              {{ stat.label }}
            </dt>
            <dd class="mt-1 text-sm leading-5 font-medium">
              {{ stat.value }}
            </dd>
          </div>
        </dl>
        <section v-if="resolvedSpectralImage" class="space-y-3">
          <div>
            <h2 class="text-sm font-semibold tracking-[0.07em] uppercase">Спектр</h2>
            <p class="mt-1 text-sm text-(--lab-text-muted)">Видимый диапазон 380–780 нм и основные линии излучения.</p>
          </div>
          <LabViewerPreviewButton
            :src="resolvedSpectralImage"
            :alt="`Спектр ${currentElement.russianName}`"
            label="Открыть спектр"
            class="w-full max-w-xl"
            frame-class="max-h-24"
            image-class="h-full w-full object-contain"
            @preview="openImageViewer(spectralViewerIndex)"
          />
        </section>
      </div>
    </section>

    <section>
      <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_20rem]">
        <article class="min-w-0 px-4 pb-4 sm:px-6 sm:pb-6">
          <div class="grid gap-x-6 gap-y-6 md:grid-cols-2 2xl:grid-cols-3">
            <section v-for="section in detailSections" :key="section.title" class="min-w-0">
              <h2 class="mb-3 text-sm font-semibold tracking-[0.07em] text-(--lab-text-muted) uppercase">
                {{ section.title }}
              </h2>

              <dl class="space-y-2">
                <div
                  v-for="item in section.items"
                  :key="`${section.title}:${item.label}`"
                  class="grid grid-cols-[max-content_minmax(0,1fr)] items-start gap-x-3"
                >
                  <dt class="text-xs leading-6 text-(--lab-text-muted)">
                    {{ item.label }}
                  </dt>

                  <dd class="min-w-0 text-sm leading-6">
                    <NuxtLink
                      v-if="item.href && item.value !== NO_DATA"
                      :to="item.href"
                      external
                      class="wrap-break-word underline underline-offset-2 transition hover:opacity-75"
                    >
                      {{ item.value }}
                    </NuxtLink>

                    <span v-else class="wrap-break-word">
                      {{ item.value }}
                    </span>
                  </dd>
                </div>
              </dl>
            </section>
          </div>
        </article>

        <div class="space-y-6 px-4 pb-4 sm:px-6 sm:pb-6 xl:pl-0">
          <article class="min-w-0">
            <section class="space-y-4">
              <div>
                <h2 class="text-base font-semibold">Связанные элементы</h2>
                <p class="mt-1 text-sm text-(--lab-text-muted)">Та же группа, период или категория.</p>
              </div>

              <div class="flex flex-wrap gap-2">
                <NuxtLink
                  v-for="item in relatedElements"
                  :key="item.number"
                  :to="periodic_table_element_route(item)"
                  class="lab-focus px-3 py-2"
                >
                  <div class="text-sm font-semibold">{{ item.displaySymbol }} · {{ item.russianName }}</div>
                  <div class="mt-1 text-xs text-(--lab-text-muted)">№{{ item.number }} · {{ item.categoryLabel }}</div>
                </NuxtLink>
              </div>
            </section>
          </article>
        </div>
      </div>
    </section>

    <LabViewerImage
      v-model="imageViewerOpen"
      :items="viewerItems"
      :initial-index="imageViewerIndex"
      :title="activeViewerTitle"
      @active-index-change="imageViewerIndex = $event"
    />
  </div>
</template>
