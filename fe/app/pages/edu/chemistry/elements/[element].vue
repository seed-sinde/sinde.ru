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
  const color_hex_to_rgba = (value: string, alpha: number) => colorHexToRgba(value, alpha)
  const periodic_table_element_route = (value: Parameters<typeof getPeriodicTableElementRoute>[0]) =>
    getPeriodicTableElementRoute(value)
  const elementParam = computed(() => route.params.element)
  const periodicTableElements = computed(() => periodicTableElementsData.value || [])
  const element = computed(() => resolvePeriodicTableElement(periodicTableElements.value, elementParam.value))
  if (!element.value) {
    throw createError({ statusCode: 404, statusMessage: 'Элемент не найден' })
  }
  const currentElement = computed(() => element.value!)
  const imageViewerOpen = ref(false)
  const imageViewerIndex = ref(0)
  const selectedSampleImageIndex = ref(0)
  const formatText = (value: unknown, fallback = NO_DATA) => {
    const text = String(value ?? '').trim()
    return text || fallback
  }
  const formatNumberValue = (value: unknown, fallback = NO_DATA) => {
    return typeof value === 'number' && Number.isFinite(value) ? `${value}` : fallback
  }
  const formatNumberList = (values: number[], fallback = NO_DATA) => {
    return values.length ? values.join(', ') : fallback
  }
  const formatPhase = (value: string | null) => {
    const phase = String(value ?? '').trim()
    if (!phase) return NO_DATA
    return PHASE_LABELS[phase] || phase
  }
  const formatBlock = (value: string | null) => {
    const block = String(value ?? '').trim()
    return block ? `${block.toUpperCase()}-блок` : NO_DATA
  }
  const resolveAssetUrl = (value: string | null | undefined) => {
    const source = String(value ?? '').trim()
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
  const cpkHexLabel = computed(() =>
    currentElement.value.cpkHex ? `#${currentElement.value.cpkHex.toUpperCase()}` : ''
  )
  const cpkChipStyle = computed(() => {
    if (!cpkHexLabel.value) return null
    return {
      backgroundColor: cpkHexLabel.value,
      color: getContrastTextColor(cpkHexLabel.value)
    }
  })
  const detailPanelStyle = computed(() => ({
    borderColor: color_hex_to_rgba(currentElement.value.accentColor, 0.38),
    backgroundImage: `linear-gradient(160deg, ${color_hex_to_rgba(currentElement.value.accentColor, 0.14)} 0%, rgba(9, 9, 11, 0.95) 68%)`
  }))
  const heroStats = computed(() => {
    const item = currentElement.value
    return [
      { label: 'Атомный номер', value: `${item.number}` },
      { label: 'Атомная масса', value: formatNumberValue(item.atomicMass) },
      { label: 'Период', value: formatNumberValue(item.period) },
      { label: 'Группа', value: formatNumberValue(item.group) },
      { label: 'Блок', value: formatBlock(item.block) },
      { label: 'Фаза', value: formatPhase(item.phase) }
    ]
  })
  const headerDescription = computed(() => {
    return currentElement.value.summary || `Химический элемент №${currentElement.value.number}.`
  })
  const detailSections = computed<DetailSection[]>(() => {
    const item = currentElement.value
    return [
      {
        title: 'Основные данные',
        items: [
          { label: 'Внешний вид', value: formatText(item.appearance) },
          { label: 'Открыл', value: formatText(item.discoveredBy) },
          { label: 'Назвал', value: formatText(item.namedBy) },
          {
            label: 'Источник',
            value: formatText(item.source),
            href: item.source || null
          }
        ]
      },
      {
        title: 'Физические свойства',
        items: [
          { label: 'Плотность', value: formatNumberValue(item.density) },
          { label: 'Температура плавления, K', value: formatNumberValue(item.melt) },
          { label: 'Температура кипения, K', value: formatNumberValue(item.boil) },
          { label: 'Молярная теплоёмкость', value: formatNumberValue(item.molarHeat) }
        ]
      },
      {
        title: 'Электронная структура',
        items: [
          { label: 'Электронная конфигурация', value: formatText(item.electronConfiguration) },
          {
            label: 'Семантическая конфигурация',
            value: formatText(item.electronConfigurationSemantic)
          },
          { label: 'Оболочки', value: formatNumberList(item.shells) },
          {
            label: 'Электроотрицательность',
            value: formatNumberValue(item.electronegativityPauling)
          },
          {
            label: 'Сродство к электрону',
            value: formatNumberValue(item.electronAffinity)
          },
          {
            label: 'Энергии ионизации',
            value: formatNumberList(item.ionizationEnergies)
          }
        ]
      },
      {
        title: 'Положение в системе',
        items: [
          { label: 'Период', value: formatNumberValue(item.period) },
          { label: 'Группа', value: formatNumberValue(item.group) },
          { label: 'Блок', value: formatBlock(item.block) },
          { label: 'X-позиция', value: formatNumberValue(item.xpos) },
          { label: 'Y-позиция', value: formatNumberValue(item.ypos) },
          { label: 'WX-позиция', value: formatNumberValue(item.wxpos) },
          { label: 'WY-позиция', value: formatNumberValue(item.wypos) }
        ]
      }
    ]
  })
  const viewerItems = computed<ViewerImageItem[]>(() => {
    const items: ViewerImageItem[] = []
    currentElement.value.samples.forEach((sample, index) => {
      const src = resolveAssetUrl(sample.url) || sample.url
      items.push({
        key: `sample:${index}`,
        kind: 'sample',
        src,
        title: sample.title || `${currentElement.value.russianName} · образец ${index + 1}`,
        alt: sample.title || `${currentElement.value.russianName} sample ${index + 1}`,
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
        title: `Модель Бора ${currentElement.value.russianName}`,
        alt: `Модель Бора ${currentElement.value.russianName}`,
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
        title: `Спектр ${currentElement.value.russianName}`,
        alt: `Спектр ${currentElement.value.russianName}`,
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
  watch(
    () => currentElement.value.number,
    () => {
      selectedSampleImageIndex.value = 0
    },
    { immediate: true }
  )
  const title = computed(() => `(${formattedCurrentSymbol.value}) ${currentElement.value.russianName}`)
  usePageSeo({
    title,
    description: () => headerDescription.value
  })
</script>
<template>
  <div class="space-y-4 pb-4">
    <LabNavHeader
      :title
      :breadcrumb-items="[
        { label: 'Вики', to: '/edu' },
        { label: 'Химия', to: '/edu/chemistry' },
        {
          label: 'Элементы',
          to: '/edu/chemistry/elements'
        },
        { label: currentElement.russianName, current: true }
      ]" />
    <section
      class="mx-3 sm:mx-6 overflow-hidden border"
      :style="{
        borderColor: color_hex_to_rgba(currentElement.accentColor, 0.5),
        backgroundImage: `linear-gradient(155deg, ${color_hex_to_rgba(currentElement.accentColor, 0.24)} 0%, rgba(9, 9, 11, 0.96) 72%)`
      }">
      <div class="grid gap-0 xl:grid-cols-[minmax(0,1fr)_minmax(0,24rem)_auto_auto] xl:grid-rows-[auto_auto_1fr]">
        <div
          class="mx-4 mt-4 flex flex-wrap items-center gap-2 text-xs text-zinc-200 sm:mx-6 sm:mt-6 xl:col-start-1 xl:row-start-1">
          <span class="border border-white/15 bg-white/8 px-3 py-1">Элемент №{{ currentElement.number }}</span>
          <span class="border border-white/15 bg-white/8 px-3 py-1">
            {{ currentElement.categoryLabel }}
          </span>
          <span class="border border-white/15 bg-white/8 px-3 py-1">
            {{ formatPhase(currentElement.phase) }}
          </span>
          <span
            v-if="cpkChipStyle && cpkHexLabel"
            class="border border-black/10 px-3 py-1 font-medium"
            :style="cpkChipStyle">
            цвет CPK
            {{ cpkHexLabel }}
          </span>
        </div>
        <div class="mx-4 mt-4 flex flex-wrap items-end gap-4 sm:mx-6 xl:col-start-1 xl:row-start-2">
          <div
            class="grid h-20 w-20 place-items-center rounded-3xl border border-white/15 bg-black/20 text-5xl font-semibold text-white sm:h-24 sm:w-24 sm:text-6xl">
            {{ formattedCurrentSymbol }}
          </div>
          <div class="space-y-1">
            <h1 class="text-3xl font-semibold tracking-tight text-white sm:text-4xl">
              {{ currentElement.russianName }}
            </h1>
            <p class="text-base text-zinc-200 sm:text-lg">
              {{ currentElement.name }}
            </p>
          </div>
        </div>
        <p
          class="mx-4 mb-4 mt-4 max-w-4xl text-sm leading-7 text-zinc-100/90 sm:mx-6 sm:mb-6 sm:text-base xl:col-start-1 xl:row-start-3">
          {{ headerDescription }}
        </p>
        <div
          class="mx-4 my-4 grid gap-3 sm:mx-6 sm:my-6 sm:grid-cols-2 xl:col-start-2 xl:row-[1/4] xl:mx-0 xl:mr-6 xl:self-stretch xl:content-start">
          <div
            v-for="stat in heroStats"
            :key="stat.label"
            class="rounded-2xl border border-white/12 bg-black/20 px-4 py-3 backdrop-blur-sm">
            <div class="text-[11px] uppercase tracking-[0.07em] text-zinc-400">
              {{ stat.label }}
            </div>
            <div class="mt-2 text-lg font-semibold text-white sm:text-xl">
              {{ stat.value }}
            </div>
          </div>
        </div>
        <div
          v-if="resolvedBohrModel3d"
          class="hidden min-h-0 self-stretch overflow-hidden border-l border-white/10 xl:col-start-3 xl:row-[1/4] xl:block">
          <LabViewer3D
            compact
            :title="currentElement.russianName"
            :poster-src="resolvedBohrModelImage"
            :model-src="resolvedBohrModel3d"
            class="h-full" />
        </div>
        <div v-if="resolvedBohrModel3d" class="mt-2 border-t border-white/10 xl:hidden">
          <LabViewer3D
            compact
            :title="currentElement.russianName"
            :poster-src="resolvedBohrModelImage"
            :model-src="resolvedBohrModel3d"
            class="aspect-square sm:aspect-auto sm:min-h-100" />
        </div>
      </div>
    </section>
    <section
      class="mx-3 overflow-hidden border sm:mx-6"
      :style="detailPanelStyle">
      <div class="grid gap-0 xl:grid-cols-[minmax(0,1.05fr)_minmax(0,1fr)]">
        <article class="min-w-0 px-4 py-4 sm:px-6 sm:py-6 xl:border-r xl:border-white/10">
          <div class="space-y-6">
            <section v-for="section in detailSections" :key="section.title" class="space-y-3">
              <div class="flex items-center gap-2">
                <span class="h-px flex-1 bg-white/10"></span>
                <span class="whitespace-nowrap text-[11px] uppercase tracking-[0.07em] text-zinc-300">
                  {{ section.title }}
                </span>
                <span class="h-px flex-1 bg-white/10"></span>
              </div>
              <dl class="divide-y divide-white/8">
                <div
                  v-for="item in section.items"
                  :key="`${section.title}:${item.label}`"
                  class="grid gap-1 py-2.5 sm:grid-cols-[14rem_minmax(0,1fr)] sm:gap-4">
                  <dt class="text-xs uppercase tracking-[0.07em] text-zinc-400">
                    {{ item.label }}
                  </dt>
                  <dd class="min-w-0 text-sm leading-6 text-zinc-100 sm:text-right">
                    <NuxtLink
                      v-if="item.href && item.value !== NO_DATA"
                      :to="item.href"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="break-all text-amber-300 transition hover:text-amber-200">
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
        <article class="min-w-0 border-t border-white/10 px-4 py-4 sm:px-6 sm:py-6 xl:border-t-0">
          <div class="space-y-8">
            <section class="space-y-4">
              <div>
                <h2 class="text-base font-semibold text-zinc-100">Образцы вещества</h2>
                <p class="mt-1 text-sm text-zinc-300/80">Фотографии образцов и материалов, связанных с элементом.</p>
              </div>
              <div v-if="selectedSampleImage">
                <LabViewerPreviewButton
                  :src="selectedSampleImage.src"
                  :alt="selectedSampleImage.alt"
                  :label="selectedSampleImage.title || 'Открыть образец'"
                  free-height
                  button-class="w-full"
                  frame-class="max-h-80"
                  image-class="max-h-80 w-full object-contain"
                  @preview="openImageViewer(selectedSampleImage.viewerIndex)" />
                <div class="space-y-3 px-4 py-4">
                  <div>
                    <p class="text-sm font-medium text-zinc-100">
                      {{ selectedSampleImage.title }}
                    </p>
                    <LabViewerImageCredits
                      class="mt-1"
                      :author="selectedSampleImage.author"
                      :attribution="selectedSampleImage.attribution"
                      :source-url="selectedSampleImage.sourceUrl"
                      :license-url="selectedSampleImage.licenseUrl"
                      :license="selectedSampleImage.license" />
                  </div>
                  <LabViewerImageThumbnails
                    :items="sampleImages"
                    :active-index="selectedSampleImageIndex"
                    @select="selectedSampleImageIndex = $event" />
                </div>
              </div>
              <section
                v-else-if="currentElement.sampleFallback"
                class="border border-amber-500/30 bg-amber-500/8 p-4">
                <div class="text-xs uppercase tracking-[0.07em] text-amber-300">
                  {{ currentElement.sampleFallback.type }}
                </div>
                <div class="mt-2 text-sm font-semibold text-amber-100">
                  {{ currentElement.sampleFallback.label }}
                </div>
                <p class="mt-2 text-sm leading-6 text-amber-50/85">
                  {{ currentElement.sampleFallback.description }}
                </p>
              </section>
              <div
                v-else
                class="border border-white/10 bg-black/15 px-4 py-6 text-center text-sm leading-6 text-zinc-300/80">
                Для этого элемента изображения образцов пока не добавлено.
              </div>
            </section>
            <section v-if="resolvedSpectralImage" class="space-y-4 border-t border-white/10 pt-6">
              <div>
                <h2 class="text-base font-semibold text-zinc-100">Спектральное изображение</h2>
                <p class="mt-1 text-sm text-zinc-300/80">
                  Приблизительная шкала видимого диапазона 380–780 нм. Яркие линии показывают длины волн, на которых
                  элемент излучает наиболее заметно.
                </p>
              </div>
              <LabViewerPreviewButton
                :src="resolvedSpectralImage"
                :alt="`Спектр ${currentElement.russianName}`"
                label="Открыть спектр"
                button-class="w-full"
                frame-class="max-h-24"
                image-class="h-full w-full object-contain"
                @preview="openImageViewer(spectralViewerIndex)" />
              <div class="pointer-events-none">
                <div class="relative h-5 text-[11px] text-zinc-400">
                  <span class="absolute left-0 translate-x-0">380</span>
                  <span class="absolute left-[12.5%] -translate-x-1/2">430</span>
                  <span class="absolute left-[25%] -translate-x-1/2">480</span>
                  <span class="absolute left-[37.5%] -translate-x-1/2">530</span>
                  <span class="absolute left-1/2 -translate-x-1/2">580</span>
                  <span class="absolute left-[62.5%] -translate-x-1/2">630</span>
                  <span class="absolute left-[75%] -translate-x-1/2">680</span>
                  <span class="absolute left-[87.5%] -translate-x-1/2">730</span>
                  <span class="absolute right-0 translate-x-0">780 нм</span>
                </div>
                <div class="relative mt-1 h-2 overflow-hidden rounded-full bg-zinc-900">
                  <div
                    class="absolute inset-0 opacity-80"
                    style="
                      background: linear-gradient(
                        90deg,
                        #6a00ff 0%,
                        #2440ff 18%,
                        #00b7ff 32%,
                        #00ff87 50%,
                        #b7ff00 63%,
                        #ffb300 79%,
                        #ff3b00 100%
                      );
                    " />
                </div>
              </div>
            </section>
            <section class="space-y-4 border-t border-white/10 pt-6">
              <div>
                <h2 class="text-base font-semibold text-zinc-100">Связанные элементы</h2>
                <p class="mt-1 text-sm text-zinc-300/80">Элементы той же группы, периода или категории.</p>
              </div>
              <div class="flex flex-wrap gap-2">
                <NuxtLink
                  v-for="item in relatedElements"
                  :key="item.number"
                  :to="periodic_table_element_route(item)"
                  class="min-w-0 flex-1 basis-56 border border-white/10 bg-black/12 px-3 py-3 transition hover:border-white/20 hover:bg-black/20">
                  <div class="text-sm font-semibold text-zinc-100">{{ item.displaySymbol }} · {{ item.russianName }}</div>
                  <div class="mt-1 text-xs leading-5 text-zinc-300/70">№{{ item.number }} · {{ item.categoryLabel }}</div>
                </NuxtLink>
              </div>
            </section>
          </div>
        </article>
      </div>
    </section>
    <LabViewerImage
      v-model="imageViewerOpen"
      :items="viewerItems"
      :initial-index="imageViewerIndex"
      :title="activeViewerTitle"
      @active-index-change="imageViewerIndex = $event" />
  </div>
</template>
