<script setup lang="ts">
  const title = 'Астрономия'
  const sun = ASTRONOMY_SUN
  const galacticContext = ASTRONOMY_GALACTIC_CONTEXT
  const planets = ASTRONOMY_PLANETS
  const moons = ASTRONOMY_MOONS
  const sizeBodies = ASTRONOMY_SIZE_BODIES
  const astronomicalUnitKm = 149597870.7
  const defaultPlanet = planets[0] as AstronomyPlanet
  const earthPlanet = planets.find(planet => planet.id === 'earth') || defaultPlanet
  const defaultSizeBody = sizeBodies[0] as AstronomyBody
  const selectedPlanetId = ref('earth')
  const selectedMoonId = ref('moon')
  const selectedSizeBodyId = ref(sun.id)
  const animationTick = ref(0)
  const galaxyZoom = ref(1)
  const planetZoom = ref(1)
  const moonZoom = ref(1)
  let animationFrameId = 0
  const numberFormatter = new Intl.NumberFormat('ru-RU')
  const decimalFormatter = new Intl.NumberFormat('ru-RU', {
    maximumFractionDigits: 2
  })
  const introMetrics = [
    {
      label: 'Центр нашей галактики',
      value: galacticContext.galacticCenterName,
      note: 'сверхмассивная чёрная дыра'
    },
    {
      label: 'Масса центра',
      value: `${decimalFormatter.format(galacticContext.galacticCenterMassSolar / 1_000_000)} млн масс Солнца`,
      note: 'ориентир для орбит в центральной области'
    },
    {
      label: 'Радиус орбиты Солнечной системы',
      value: `${numberFormatter.format(galacticContext.solarSystemRadiusLy)} св. лет`,
      note: 'расстояние до центра Млечного Пути'
    },
    {
      label: 'Галактический год',
      value: `${numberFormatter.format(galacticContext.solarSystemOrbitalPeriodMillionYears)} млн лет`,
      note: 'один оборот Солнечной системы вокруг центра'
    }
  ]
  const compositionFacts = [
    'Вселенная включает всё наблюдаемое вещество, энергию, пространство и время. Галактики объединяются в группы и скопления, связанные гравитацией.',
    'Млечный Путь содержит звёзды, газ, пыль и тёмную материю. Солнечная система находится примерно в 26,6 тыс. световых лет от центра галактики.',
    'По массе Солнце и протопланетный материал состояли преимущественно из водорода, гелия и небольшой доли более тяжёлых элементов.'
  ]
  const solarComposition = [
    { label: 'Водород', symbol: 'H', value: 0.715, color: '#38bdf8' },
    { label: 'Гелий', symbol: 'He', value: 0.271, color: '#f59e0b' },
    { label: 'Металлы', symbol: 'Me', value: 0.014, color: '#a78bfa' }
  ]
  const sectionAnchors = [
    { id: 'astronomy-intro', label: 'Обзор' },
    { id: 'astronomy-galaxy', label: 'Галактика' },
    { id: 'astronomy-planets', label: 'Планеты' },
    { id: 'astronomy-moons', label: 'Спутники' },
    { id: 'astronomy-sizes', label: 'Диаметры' }
  ]
  const selectedPlanet = computed<AstronomyPlanet>(() => {
    return planets.find(planet => planet.id === selectedPlanetId.value) || defaultPlanet
  })
  const selectedPlanetMoons = computed(() => moons.filter(moon => moon.planetId === selectedPlanet.value.id))
  const planetsWithMoons = computed(() => planets.filter(planet => moons.some(moon => moon.planetId === planet.id)))
  const selectedPlanetHasMoons = computed(() => selectedPlanetMoons.value.length > 0)
  watch(
    selectedPlanetMoons,
    next => {
      const firstMoon = next[0] || null
      if (!firstMoon) {
        selectedMoonId.value = ''
        return
      }
      if (!next.some(moon => moon.id === selectedMoonId.value)) {
        selectedMoonId.value = firstMoon.id
      }
    },
    { immediate: true }
  )
  const selectedMoon = computed(
    () =>
      selectedPlanetMoons.value.find(moon => moon.id === selectedMoonId.value) || selectedPlanetMoons.value[0] || null
  )
  const selectedSizeBody = computed<AstronomyBody>(() => {
    return sizeBodies.find(body => body.id === selectedSizeBodyId.value) || defaultSizeBody
  })
  const planetCardItems = computed<AstronomyCardStripItem[]>(() =>
    planets.map(planet => ({
      id: planet.id,
      symbol: planet.symbol,
      name: planet.name,
      meta: formatDistanceAu(planet.orbitalRadiusAu),
      color: planet.color
    }))
  )
  const moonPlanetCardItems = computed<AstronomyCardStripItem[]>(() =>
    planetsWithMoons.value.map(planet => ({
      id: planet.id,
      symbol: planet.symbol,
      name: planet.name,
      meta: `${moons.filter(moon => moon.planetId === planet.id).length} спут.`,
      color: planet.color
    }))
  )
  const moonCardItems = computed<AstronomyCardStripItem[]>(() =>
    selectedPlanetMoons.value.map(moon => ({
      id: moon.id,
      symbol: moon.symbol,
      name: moon.name,
      meta: formatDistanceKm(moon.orbitalRadiusKm),
      color: moon.color
    }))
  )
  const solarSystemMarkerAngle = computed(() => {
    return ((animationTick.value % 32000) / 32000) * Math.PI * 2 - Math.PI / 5
  })
  const solarSystemMarker = computed(() => ({
    x: Math.cos(solarSystemMarkerAngle.value) * galacticContext.solarSystemRadiusLy,
    y: Math.sin(solarSystemMarkerAngle.value) * galacticContext.solarSystemRadiusLy
  }))
  const galaxyViewBox = computed(() => {
    const half = 32000 / galaxyZoom.value
    const center = galaxyZoom.value > 1 ? solarSystemMarker.value : { x: 0, y: 0 }
    return `${center.x - half} ${center.y - half} ${half * 2} ${half * 2}`
  })
  const planetOrbitRadiusMaxAu = Math.max(...planets.map(planet => planet.orbitalRadiusAu))
  const planetSceneRadius = planetOrbitRadiusMaxAu + 2
  const selectedPlanetViewRadiusAu = computed(() => {
    const base = selectedPlanet.value.orbitalRadiusAu * 1.16 + 0.45
    return Math.min(planetSceneRadius, Math.max(0.9, base))
  })
  const planetViewHalf = computed(() => selectedPlanetViewRadiusAu.value / planetZoom.value)
  const planetViewBox = computed(() => {
    const half = planetViewHalf.value
    return `${-half} ${-half} ${half * 2} ${half * 2}`
  })
  const orbitDurationMs = (periodDays: number) => {
    return (8 + Math.log10(periodDays + 1) * 7) * 1000
  }
  const getOrbitalPosition = (radius: number, periodDays: number, seed: number) => {
    const durationMs = orbitDurationMs(periodDays)
    const progress = (animationTick.value % durationMs) / durationMs
    const angle = progress * Math.PI * 2 + seed
    return {
      x: Math.cos(angle) * radius,
      y: Math.sin(angle) * radius
    }
  }
  const planetPositions = computed(() => {
    return new Map(
      planets.map((planet, index) => [
        planet.id,
        getOrbitalPosition(planet.orbitalRadiusAu, planet.orbitalPeriodDays, index * 0.9)
      ])
    )
  })
  const planetPosition = (planetId: string) => {
    return planetPositions.value.get(planetId) || { x: 0, y: 0 }
  }
  const moonOrbitRadiusMaxKm = computed(() =>
    Math.max(...selectedPlanetMoons.value.map(moon => moon.orbitalRadiusKm), selectedPlanet.value.diameterKm * 12)
  )
  const moonSceneRadius = computed(() => moonOrbitRadiusMaxKm.value * 1.3)
  const selectedMoonOrbitRadiusKm = computed(() => {
    return selectedMoon.value?.orbitalRadiusKm || moonOrbitRadiusMaxKm.value
  })
  const moonViewHalf = computed(() => {
    const base = selectedMoonOrbitRadiusKm.value * 1.16 + selectedPlanet.value.diameterKm * 0.9
    return Math.min(moonSceneRadius.value, Math.max(selectedPlanet.value.diameterKm * 2.6, base)) / moonZoom.value
  })
  const moonViewBox = computed(() => {
    const half = moonViewHalf.value
    const center = { x: 0, y: 0 }
    return `${center.x - half} ${center.y - half} ${half * 2} ${half * 2}`
  })
  const moonPositions = computed(() => {
    return new Map(
      selectedPlanetMoons.value.map((moon, index) => [
        moon.id,
        getOrbitalPosition(moon.orbitalRadiusKm, moon.orbitalPeriodDays, index * 1.2)
      ])
    )
  })
  const moonPosition = (moonId: string) => {
    return moonPositions.value.get(moonId) || { x: 0, y: 0 }
  }
  const sizeRatioPercent = (diameterKm: number) => {
    return (diameterKm / sun.diameterKm) * 100
  }
  const bodyRadiusAu = (diameterKm: number) => diameterKm / 2 / astronomicalUnitKm
  const orbitMarkerRadius = computed(() => {
    return Math.min(0.22, Math.max(0.03, planetViewHalf.value / 5.5))
  })
  const sunReticleRadius = computed(() => Math.max(0.08, Math.min(0.22, planetViewHalf.value / 5.5)))
  const selectedPlanetReticleRadius = computed(() => Math.max(0.12, Math.min(0.62, planetViewHalf.value / 4.2)))
  const moonReticleRadius = computed(() => Math.max(selectedPlanet.value.diameterKm * 0.7, moonViewHalf.value / 7))
  const planetLabelFontSize = computed(() => Math.max(0.12, Math.min(1.6, planetViewHalf.value * 0.07)))
  const planetLabelOffset = computed(() => Math.max(0.16, planetViewHalf.value * 0.11))
  const reticleLineLength = (radius: number) => radius * 0.52
  const reticleStrokeWidth = (radius: number) => Math.max(radius / 12, radius * 0.08)
  const colorHexToRgba = (hexColor: string, alpha: number) => {
    const normalized = String(hexColor || '').replace('#', '')
    if (!/^[0-9a-f]{6}$/i.test(normalized)) {
      return `rgba(82, 82, 91, ${alpha})`
    }
    const red = Number.parseInt(normalized.slice(0, 2), 16)
    const green = Number.parseInt(normalized.slice(2, 4), 16)
    const blue = Number.parseInt(normalized.slice(4, 6), 16)
    return `rgba(${red}, ${green}, ${blue}, ${alpha})`
  }
  const bodyCardStyle = (color: string) => ({
    borderColor: colorHexToRgba(color, 0.4),
    backgroundImage: `linear-gradient(160deg, ${colorHexToRgba(color, 0.22)} 0%, color-mix(in srgb, var(--lab-bg-surface) 96%, transparent) 76%)`
  })
  const formatDistanceAu = (value: number) => `${decimalFormatter.format(value)} а.е.`
  const formatDistanceKm = (value: number) => `${numberFormatter.format(Math.round(value))} км`
  const formatDiameter = (value: number) => `${numberFormatter.format(Math.round(value))} км`
  const formatPeriodDays = (value: number) => `${decimalFormatter.format(value)} суток`
  const formatLargeMass = (value: number) => `${decimalFormatter.format(value / 1_000_000)} млн`
  const sizeBodyRelationToSun = computed(() => selectedSizeBody.value.diameterKm / sun.diameterKm)
  const sizeBodyRelationToEarth = computed(() => selectedSizeBody.value.diameterKm / earthPlanet.diameterKm)
  const onPlanetFocus = (planet: AstronomyPlanet) => {
    selectedPlanetId.value = planet.id
    selectedSizeBodyId.value = planet.id
    planetZoom.value = 1
  }
  const onPlanetCardSelect = (planetId: string) => {
    const planet = planets.find(item => item.id === planetId)
    if (planet) onPlanetFocus(planet)
  }
  const onMoonPlanetFocus = (planet: AstronomyPlanet) => {
    selectedPlanetId.value = planet.id
    selectedSizeBodyId.value = planet.id
    moonZoom.value = 1
  }
  const onMoonPlanetCardSelect = (planetId: string) => {
    const planet = planets.find(item => item.id === planetId)
    if (planet) onMoonPlanetFocus(planet)
  }
  const onMoonFocus = (moon: AstronomyMoon) => {
    selectedMoonId.value = moon.id
    selectedSizeBodyId.value = moon.id
    moonZoom.value = 1
  }
  const onMoonCardSelect = (moonId: string) => {
    const moon = selectedPlanetMoons.value.find(item => item.id === moonId)
    if (moon) onMoonFocus(moon)
  }
  const onSizeBodyFocus = (body: AstronomyBody) => {
    selectedSizeBodyId.value = body.id
  }
  const clampZoom = (value: number) => Math.max(1, Number(value.toFixed(2)))
  const zoomOut = (target: 'galaxy' | 'planet' | 'moon') => {
    if (target === 'galaxy') galaxyZoom.value = clampZoom(galaxyZoom.value / 1.5)
    if (target === 'planet') planetZoom.value = clampZoom(planetZoom.value / 2)
    if (target === 'moon') moonZoom.value = clampZoom(moonZoom.value / 1.5)
  }
  const zoomIn = (target: 'galaxy' | 'planet' | 'moon') => {
    if (target === 'galaxy') galaxyZoom.value = Math.min(150, clampZoom(galaxyZoom.value * 1.5))
    if (target === 'planet') planetZoom.value = Math.min(20000, clampZoom(planetZoom.value * 2))
    if (target === 'moon') moonZoom.value = Math.min(50, clampZoom(moonZoom.value * 1.5))
  }
  const resetZoom = (target: 'galaxy' | 'planet' | 'moon') => {
    if (target === 'galaxy') galaxyZoom.value = 1
    if (target === 'planet') planetZoom.value = 1
    if (target === 'moon') moonZoom.value = 1
  }
  const scrollToSection = (sectionId: string) => {
    if (!import.meta.client) return
    document.getElementById(sectionId)?.scrollIntoView({
      behavior: 'smooth',
      block: 'start'
    })
  }
  onMounted(() => {
    if (!import.meta.client) return
    const animate = (time: number) => {
      animationTick.value = time
      animationFrameId = window.requestAnimationFrame(animate)
    }
    animationFrameId = window.requestAnimationFrame(animate)
  })
  onBeforeUnmount(() => {
    if (!import.meta.client || !animationFrameId) return
    window.cancelAnimationFrame(animationFrameId)
  })
  usePageSeo({
    title,
    description: 'Визуализация галактической орбиты Солнечной системы, траекторий планет и крупных спутников.'
  })
</script>
<template>
  <div class="astro-page">
    <LabNavHeader
      :title="title"
      :breadcrumb-items="[
        { label: 'Вики', to: '/edu' },
        { label: title, current: true }
      ]" />
    <div class="astro-sticky-nav">
      <div class="astro-sticky-nav__inner">
        <LabBaseButton
          v-for="section in sectionAnchors"
          :key="section.id"
          :label="section.label"
          size="sm"
          variant="ghost"
          @click="scrollToSection(section.id)" />
      </div>
    </div>
    <section id="astronomy-intro" class="astro-section">
      <div class="astro-section__body astro-grid-2">
        <div class="astro-stack-lg">
          <div class="astro-stack-sm">
            <p class="astro-paragraph">
              Страница связывает несколько масштабов сразу: орбиту Солнечной системы вокруг центра Млечного Пути,
              планетные траектории вокруг Солнца и локальные орбиты крупнейших спутников. Исходная идея остаётся той же:
              Вселенная включает всё, галактики собираются в структуры, а наша система движется внутри Млечного Пути как
              часть большого гравитационного узора.
            </p>
            <p class="astro-paragraph">
              В центре нашей галактики расположен {{ galacticContext.galacticCenterName }} массой около
              {{ formatLargeMass(galacticContext.galacticCenterMassSolar) }} масс Солнца. Солнечная система находится
              примерно в {{ numberFormatter.format(galacticContext.solarSystemRadiusLy) }} световых лет от этой точки и
              делает полный оборот вокруг центра примерно за
              {{ numberFormatter.format(galacticContext.solarSystemOrbitalPeriodMillionYears) }} миллионов лет.
            </p>
          </div>
          <div class="astro-metric-grid">
            <div
              v-for="item in solarComposition"
              :key="item.label"
              class="astro-color-card"
              :style="bodyCardStyle(item.color)">
              <div class="astro-color-card__head">
                <p class="astro-kicker astro-kicker-strong">{{ item.label }}</p>
                <span class="astro-color-card__symbol">{{ item.symbol }}</span>
              </div>
              <p class="astro-color-card__value">{{ decimalFormatter.format(item.value) }}</p>
            </div>
          </div>
          <div class="astro-stack-sm">
            <p v-for="paragraph in compositionFacts" :key="paragraph" class="astro-note">
              {{ paragraph }}
            </p>
          </div>
        </div>
        <div class="astro-facts-grid">
          <article v-for="item in introMetrics" :key="item.label" class="astro-card">
            <p class="astro-kicker">{{ item.label }}</p>
            <p class="astro-card__value">{{ item.value }}</p>
            <p class="astro-note">{{ item.note }}</p>
          </article>
        </div>
      </div>
    </section>
    <section id="astronomy-galaxy" class="astro-section">
      <div class="astro-section__body astro-grid-visual">
        <div class="astro-stack-md astro-main-column">
          <div class="astro-stack-sm">
            <h2 class="astro-heading">Орбита Солнечной системы вокруг центра Галактики</h2>
            <p class="astro-paragraph astro-paragraph-muted">
              Кольцо ниже показывает положение Солнечной системы на орбите вокруг центра Млечного Пути. Радиус орбиты
              сохранён пропорционально в световых годах, а сама система отмечена как навигационная точка на этой шкале.
            </p>
          </div>
          <div class="astro-toolbar">
            <LabBaseButton icon="ic:round-remove" icon-only size="sm" variant="ghost" @click="zoomOut('galaxy')" />
            <LabBaseButton icon="ic:round-add" icon-only size="sm" variant="ghost" @click="zoomIn('galaxy')" />
            <LabBaseButton label="Сбросить масштаб" size="sm" variant="ghost" @click="resetZoom('galaxy')" />
            <span class="astro-zoom-label">x{{ decimalFormatter.format(galaxyZoom) }}</span>
          </div>
          <div class="astro-panel astro-panel-visual astro-panel-galaxy">
            <svg
              class="astro-svg"
              :viewBox="galaxyViewBox"
              role="img"
              aria-label="Орбита Солнечной системы вокруг центра Млечного Пути">
              <defs>
                <radialGradient id="galaxy-core-glow" cx="50%" cy="50%" r="50%">
                  <stop offset="0%" stop-color="#fca5a5" stop-opacity="0.9" />
                  <stop offset="55%" stop-color="#f97316" stop-opacity="0.24" />
                  <stop offset="100%" stop-color="#09090b" stop-opacity="0" />
                </radialGradient>
              </defs>
              <circle cx="0" cy="0" r="28000" fill="none" stroke="var(--astro-line-muted)" stroke-width="700" />
              <circle
                cx="0"
                cy="0"
                :r="galacticContext.solarSystemRadiusLy"
                fill="none"
                stroke="var(--astro-line)"
                stroke-width="360"
                stroke-dasharray="1800 1400" />
              <circle cx="0" cy="0" r="7200" fill="url(#galaxy-core-glow)" />
              <circle cx="0" cy="0" r="420" fill="#fca5a5" />
              <circle :cx="solarSystemMarker.x" :cy="solarSystemMarker.y" r="500" fill="#fbbf24" />
              <line
                x1="0"
                y1="0"
                :x2="solarSystemMarker.x"
                :y2="solarSystemMarker.y"
                stroke="var(--astro-line-strong)"
                stroke-width="180"
                stroke-dasharray="900 900" />
            </svg>
          </div>
        </div>
        <aside class="astro-side-column">
          <article class="astro-card">
            <p class="astro-kicker">Положение Солнечной системы</p>
            <p class="astro-card__value">
              {{ numberFormatter.format(galacticContext.solarSystemRadiusLy) }} св. лет от центра
            </p>
            <p class="astro-note">
              Анимация показывает движение вдоль галактической орбиты. Реальный путь сложнее круга, но круговая схема
              помогает удержать масштаб расстояния и общий ритм обращения.
            </p>
          </article>
          <article class="astro-card">
            <p class="astro-kicker">Галактическая скорость</p>
            <p class="astro-card__value">
              {{ numberFormatter.format(galacticContext.solarSystemOrbitalSpeedKmS) }} км/с
            </p>
            <p class="astro-note">
              На этом масштабе движение планет вокруг Солнца исчезающе мало по сравнению с орбитой всей системы вокруг
              центра галактики.
            </p>
          </article>
          <AstronomyModelSpoiler
            :title="`3D модель: ${sun.name}`"
            :viewer-title="sun.name"
            :model-src="sun.modelSrc"
            :model-size-label="sun.modelSizeLabel"
            :rotation-per-second="sun.modelRotationPerSecond"
            hint="Модель Солнца вынесена в этот слайд, потому что здесь показано положение всей Солнечной системы на галактической орбите." />
        </aside>
      </div>
    </section>
    <section id="astronomy-planets" class="astro-section">
      <div class="astro-section__body astro-stack-md">
        <div class="astro-stack-sm">
          <h2 class="astro-heading">Траектории планет вокруг Солнца</h2>
          <p class="astro-paragraph astro-paragraph-muted">
            Здесь сохранён масштаб расстояний от Солнца в астрономических единицах. Орбитальные линии пропорциональны, а
            реальные диски планет на общем плане почти исчезают. Поэтому рядом с каждой планетой есть отдельный маркер
            положения, а размер самого диска внутри него остаётся в реальном масштабе.
          </p>
        </div>
        <AstronomyCardStrip
          :items="planetCardItems"
          :selected-id="selectedPlanet.id"
          large-symbol
          wide
          @select="onPlanetCardSelect" />
        <div class="astro-toolbar">
          <LabBaseButton icon="ic:round-remove" icon-only size="sm" variant="ghost" @click="zoomOut('planet')" />
          <LabBaseButton icon="ic:round-add" icon-only size="sm" variant="ghost" @click="zoomIn('planet')" />
          <LabBaseButton label="Сбросить масштаб" size="sm" variant="ghost" @click="resetZoom('planet')" />
          <span class="astro-zoom-label">x{{ decimalFormatter.format(planetZoom) }}</span>
        </div>
        <div class="astro-visual-grid">
          <div class="astro-main-column">
            <div class="astro-panel astro-panel-visual astro-panel-map">
              <svg class="astro-svg" :viewBox="planetViewBox" role="img" aria-label="Орбиты планет вокруг Солнца">
                <circle
                  v-for="planet in planets"
                  :key="`${planet.id}-orbit`"
                  cx="0"
                  cy="0"
                  :r="planet.orbitalRadiusAu"
                  fill="none"
                  stroke="var(--astro-line-muted)"
                  stroke-width="0.08" />
                <circle cx="0" cy="0" :r="bodyRadiusAu(sun.diameterKm)" fill="#fbbf24" />
                <g class="astro-reticle" aria-hidden="true">
                  <circle
                    cx="0"
                    cy="0"
                    :r="sunReticleRadius"
                    fill="none"
                    stroke="#fde68a"
                    :stroke-width="reticleStrokeWidth(sunReticleRadius)" />
                  <line
                    :x1="-sunReticleRadius - reticleLineLength(sunReticleRadius)"
                    y1="0"
                    :x2="-sunReticleRadius"
                    y2="0"
                    stroke="#fde68a"
                    :stroke-width="reticleStrokeWidth(sunReticleRadius)" />
                  <line
                    :x1="sunReticleRadius"
                    y1="0"
                    :x2="sunReticleRadius + reticleLineLength(sunReticleRadius)"
                    y2="0"
                    stroke="#fde68a"
                    :stroke-width="reticleStrokeWidth(sunReticleRadius)" />
                  <line
                    x1="0"
                    :y1="-sunReticleRadius - reticleLineLength(sunReticleRadius)"
                    x2="0"
                    :y2="-sunReticleRadius"
                    stroke="#fde68a"
                    :stroke-width="reticleStrokeWidth(sunReticleRadius)" />
                  <line
                    x1="0"
                    :y1="sunReticleRadius"
                    x2="0"
                    :y2="sunReticleRadius + reticleLineLength(sunReticleRadius)"
                    stroke="#fde68a"
                    :stroke-width="reticleStrokeWidth(sunReticleRadius)" />
                </g>
                <g
                  v-for="planet in planets"
                  :key="planet.id"
                  class="astro-hit-area"
                  @mouseenter="onPlanetFocus(planet)"
                  @click="onPlanetFocus(planet)">
                  <circle
                    :cx="planetPosition(planet.id).x"
                    :cy="planetPosition(planet.id).y"
                    :r="orbitMarkerRadius"
                    fill="transparent" />
                  <circle
                    :cx="planetPosition(planet.id).x"
                    :cy="planetPosition(planet.id).y"
                    :r="orbitMarkerRadius"
                    fill="none"
                    :stroke="
                      planet.id === selectedPlanet.id ? 'var(--lab-text-primary)' : colorHexToRgba(planet.color, 0.7)
                    "
                    :stroke-width="orbitMarkerRadius / 5" />
                  <circle
                    :cx="planetPosition(planet.id).x"
                    :cy="planetPosition(planet.id).y"
                    :r="bodyRadiusAu(planet.diameterKm)"
                    :fill="planet.color"
                    :stroke="planet.id === selectedPlanet.id ? 'var(--lab-text-primary)' : 'none'"
                    :stroke-width="planet.id === selectedPlanet.id ? bodyRadiusAu(planet.diameterKm) * 0.3 + 0.01 : 0">
                    <title>{{ planet.name }}</title>
                  </circle>
                  <text
                    :x="planetPosition(planet.id).x"
                    :y="planetPosition(planet.id).y + orbitMarkerRadius + planetLabelOffset"
                    text-anchor="middle"
                    class="astro-planet-label"
                    :style="{ fontSize: `${planetLabelFontSize}px` }">
                    {{ planet.name }}
                  </text>
                </g>
                <g class="astro-reticle" aria-hidden="true">
                  <circle
                    :cx="planetPosition(selectedPlanet.id).x"
                    :cy="planetPosition(selectedPlanet.id).y"
                    :r="selectedPlanetReticleRadius"
                    fill="none"
                    stroke="var(--lab-text-primary)"
                    :stroke-width="reticleStrokeWidth(selectedPlanetReticleRadius)" />
                  <line
                    :x1="
                      planetPosition(selectedPlanet.id).x -
                      selectedPlanetReticleRadius -
                      reticleLineLength(selectedPlanetReticleRadius)
                    "
                    :y1="planetPosition(selectedPlanet.id).y"
                    :x2="planetPosition(selectedPlanet.id).x - selectedPlanetReticleRadius"
                    :y2="planetPosition(selectedPlanet.id).y"
                    stroke="var(--lab-text-primary)"
                    :stroke-width="reticleStrokeWidth(selectedPlanetReticleRadius)" />
                  <line
                    :x1="planetPosition(selectedPlanet.id).x + selectedPlanetReticleRadius"
                    :y1="planetPosition(selectedPlanet.id).y"
                    :x2="
                      planetPosition(selectedPlanet.id).x +
                      selectedPlanetReticleRadius +
                      reticleLineLength(selectedPlanetReticleRadius)
                    "
                    :y2="planetPosition(selectedPlanet.id).y"
                    stroke="var(--lab-text-primary)"
                    :stroke-width="reticleStrokeWidth(selectedPlanetReticleRadius)" />
                  <line
                    :x1="planetPosition(selectedPlanet.id).x"
                    :y1="
                      planetPosition(selectedPlanet.id).y -
                      selectedPlanetReticleRadius -
                      reticleLineLength(selectedPlanetReticleRadius)
                    "
                    :x2="planetPosition(selectedPlanet.id).x"
                    :y2="planetPosition(selectedPlanet.id).y - selectedPlanetReticleRadius"
                    stroke="var(--lab-text-primary)"
                    :stroke-width="reticleStrokeWidth(selectedPlanetReticleRadius)" />
                  <line
                    :x1="planetPosition(selectedPlanet.id).x"
                    :y1="planetPosition(selectedPlanet.id).y + selectedPlanetReticleRadius"
                    :x2="planetPosition(selectedPlanet.id).x"
                    :y2="
                      planetPosition(selectedPlanet.id).y +
                      selectedPlanetReticleRadius +
                      reticleLineLength(selectedPlanetReticleRadius)
                    "
                    stroke="var(--lab-text-primary)"
                    :stroke-width="reticleStrokeWidth(selectedPlanetReticleRadius)" />
                </g>
              </svg>
            </div>
          </div>
          <aside class="astro-side-column">
            <article class="astro-card">
              <p class="astro-kicker">Выбранная планета</p>
              <p class="astro-card__value">{{ selectedPlanet.name }}</p>
              <p class="astro-note wrap-break-word">{{ selectedPlanet.description }}</p>
            </article>
            <AstronomyModelSpoiler
              :title="`3D модель: ${selectedPlanet.name}`"
              :viewer-title="selectedPlanet.name"
              :model-src="selectedPlanet.modelSrc"
              :model-size-label="selectedPlanet.modelSizeLabel"
              :rotation-per-second="selectedPlanet.modelRotationPerSecond"
              hint="Модель встраивается только после раскрытия этого блока, поэтому не нагружает страницу заранее." />
            <article class="astro-card">
              <div class="astro-stat-list">
                <div class="astro-stat-row">
                  <span class="astro-stat-label">Диаметр</span>
                  <span>{{ formatDiameter(selectedPlanet.diameterKm) }}</span>
                </div>
                <div class="astro-stat-row">
                  <span class="astro-stat-label">Средняя орбита</span>
                  <span>{{ formatDistanceAu(selectedPlanet.orbitalRadiusAu) }}</span>
                </div>
                <div class="astro-stat-row">
                  <span class="astro-stat-label">Период обращения</span>
                  <span>{{ formatPeriodDays(selectedPlanet.orbitalPeriodDays) }}</span>
                </div>
                <div class="astro-stat-row">
                  <span class="astro-stat-label">Крупных спутников в модели</span>
                  <span>{{ selectedPlanetMoons.length }}</span>
                </div>
              </div>
            </article>
            <article class="astro-card">
              <p class="astro-kicker">Почему планеты почти не видны</p>
              <p class="astro-note wrap-break-word">
                Да, это ожидаемо на реальном масштабе: диаметр планеты ничтожно мал по сравнению с расстоянием до
                Солнца. Светлое кольцо показывает положение планеты, а цветной диск внутри него остаётся
                пропорциональным.
              </p>
            </article>
          </aside>
        </div>
      </div>
    </section>
    <section id="astronomy-moons" class="astro-section">
      <div class="astro-section__body astro-stack-md">
        <div class="astro-stack-sm">
          <h2 class="astro-heading">Спутники выбранной планеты</h2>
          <p class="astro-paragraph astro-paragraph-muted">
            Для локальной системы ниже сохранены пропорции и орбитальных расстояний, и размеров тел. Чем ближе спутник к
            планете, тем плотнее видна траектория.
          </p>
        </div>
        <AstronomyCardStrip
          :items="moonPlanetCardItems"
          :selected-id="selectedPlanet.id"
          compact
          large-symbol
          wide
          @select="onMoonPlanetCardSelect" />
        <div class="astro-toolbar">
          <LabBaseButton icon="ic:round-remove" icon-only size="sm" variant="ghost" @click="zoomOut('moon')" />
          <LabBaseButton icon="ic:round-add" icon-only size="sm" variant="ghost" @click="zoomIn('moon')" />
          <LabBaseButton label="Сбросить масштаб" size="sm" variant="ghost" @click="resetZoom('moon')" />
          <span class="astro-zoom-label">x{{ decimalFormatter.format(moonZoom) }}</span>
        </div>
        <div class="astro-visual-grid">
          <div class="astro-main-column">
            <div class="astro-panel astro-panel-visual astro-panel-map">
              <svg
                v-if="selectedPlanetMoons.length"
                class="astro-svg"
                :viewBox="moonViewBox"
                role="img"
                :aria-label="`Орбиты спутников планеты ${selectedPlanet.name}`">
                <circle
                  v-for="moon in selectedPlanetMoons"
                  :key="`${moon.id}-orbit`"
                  cx="0"
                  cy="0"
                  :r="moon.orbitalRadiusKm"
                  fill="none"
                  stroke="var(--astro-line-muted)"
                  :stroke-width="moonSceneRadius / 420" />
                <circle cx="0" cy="0" :r="selectedPlanet.diameterKm / 2" :fill="selectedPlanet.color" />
                <g class="astro-reticle" aria-hidden="true">
                  <circle
                    cx="0"
                    cy="0"
                    :r="moonReticleRadius"
                    fill="none"
                    stroke="var(--lab-text-primary)"
                    :stroke-width="reticleStrokeWidth(moonReticleRadius)" />
                  <line
                    :x1="-moonReticleRadius - reticleLineLength(moonReticleRadius)"
                    y1="0"
                    :x2="-moonReticleRadius"
                    y2="0"
                    stroke="var(--lab-text-primary)"
                    :stroke-width="reticleStrokeWidth(moonReticleRadius)" />
                  <line
                    :x1="moonReticleRadius"
                    y1="0"
                    :x2="moonReticleRadius + reticleLineLength(moonReticleRadius)"
                    y2="0"
                    stroke="var(--lab-text-primary)"
                    :stroke-width="reticleStrokeWidth(moonReticleRadius)" />
                  <line
                    x1="0"
                    :y1="-moonReticleRadius - reticleLineLength(moonReticleRadius)"
                    x2="0"
                    :y2="-moonReticleRadius"
                    stroke="var(--lab-text-primary)"
                    :stroke-width="reticleStrokeWidth(moonReticleRadius)" />
                  <line
                    x1="0"
                    :y1="moonReticleRadius"
                    x2="0"
                    :y2="moonReticleRadius + reticleLineLength(moonReticleRadius)"
                    stroke="var(--lab-text-primary)"
                    :stroke-width="reticleStrokeWidth(moonReticleRadius)" />
                </g>
                <g
                  v-for="moon in selectedPlanetMoons"
                  :key="moon.id"
                  class="astro-hit-area"
                  @mouseenter="onMoonFocus(moon)"
                  @click="onMoonFocus(moon)">
                  <circle
                    :cx="moonPosition(moon.id).x"
                    :cy="moonPosition(moon.id).y"
                    :r="moon.diameterKm / 2"
                    :fill="moon.color"
                    :stroke="selectedMoon?.id === moon.id ? 'var(--lab-text-primary)' : 'none'"
                    :stroke-width="selectedMoon?.id === moon.id ? moonSceneRadius / 520 : 0">
                    <title>{{ moon.name }}</title>
                  </circle>
                </g>
                <g v-if="selectedMoon" class="astro-reticle" aria-hidden="true">
                  <circle
                    :cx="moonPosition(selectedMoon.id).x"
                    :cy="moonPosition(selectedMoon.id).y"
                    :r="moonReticleRadius"
                    fill="none"
                    stroke="#93c5fd"
                    :stroke-width="reticleStrokeWidth(moonReticleRadius)" />
                  <line
                    :x1="moonPosition(selectedMoon.id).x - moonReticleRadius - reticleLineLength(moonReticleRadius)"
                    :y1="moonPosition(selectedMoon.id).y"
                    :x2="moonPosition(selectedMoon.id).x - moonReticleRadius"
                    :y2="moonPosition(selectedMoon.id).y"
                    stroke="#93c5fd"
                    :stroke-width="reticleStrokeWidth(moonReticleRadius)" />
                  <line
                    :x1="moonPosition(selectedMoon.id).x + moonReticleRadius"
                    :y1="moonPosition(selectedMoon.id).y"
                    :x2="moonPosition(selectedMoon.id).x + moonReticleRadius + reticleLineLength(moonReticleRadius)"
                    :y2="moonPosition(selectedMoon.id).y"
                    stroke="#93c5fd"
                    :stroke-width="reticleStrokeWidth(moonReticleRadius)" />
                  <line
                    :x1="moonPosition(selectedMoon.id).x"
                    :y1="moonPosition(selectedMoon.id).y - moonReticleRadius - reticleLineLength(moonReticleRadius)"
                    :x2="moonPosition(selectedMoon.id).x"
                    :y2="moonPosition(selectedMoon.id).y - moonReticleRadius"
                    stroke="#93c5fd"
                    :stroke-width="reticleStrokeWidth(moonReticleRadius)" />
                  <line
                    :x1="moonPosition(selectedMoon.id).x"
                    :y1="moonPosition(selectedMoon.id).y + moonReticleRadius"
                    :x2="moonPosition(selectedMoon.id).x"
                    :y2="moonPosition(selectedMoon.id).y + moonReticleRadius + reticleLineLength(moonReticleRadius)"
                    stroke="#93c5fd"
                    :stroke-width="reticleStrokeWidth(moonReticleRadius)" />
                </g>
              </svg>
              <div v-else class="astro-empty-state">
                Для {{ selectedPlanet.name.toLowerCase() }} в этой подборке нет спутников. Выберите другую планету.
              </div>
            </div>
          </div>
          <aside class="astro-side-column">
            <AstronomyCardStrip
              v-if="selectedPlanetHasMoons"
              :items="moonCardItems"
              :selected-id="selectedMoon?.id || null"
              @select="onMoonCardSelect" />
            <article v-else class="astro-card">
              <p class="astro-kicker">Что выбрать дальше</p>
              <p class="astro-note wrap-break-word">
                В этой модели спутники добавлены не для всех планет. Ниже сразу доступны системы, где есть что
                посмотреть.
              </p>
              <div class="astro-chip-row">
                <LabBaseButton
                  v-for="planet in planetsWithMoons"
                  :key="`${planet.id}-suggested`"
                  :label="planet.name"
                  size="sm"
                  variant="ghost"
                  @click="onMoonPlanetFocus(planet)" />
              </div>
            </article>
            <article v-if="selectedMoon" class="astro-card">
              <p class="astro-kicker">Выбранный спутник</p>
              <p class="astro-card__value">{{ selectedMoon.name }}</p>
              <p class="astro-note wrap-break-word">{{ selectedMoon.description }}</p>
            </article>
            <article v-if="selectedMoon" class="astro-card">
              <div class="astro-stat-list">
                <div class="astro-stat-row">
                  <span class="astro-stat-label">Диаметр</span>
                  <span>{{ formatDiameter(selectedMoon.diameterKm) }}</span>
                </div>
                <div class="astro-stat-row">
                  <span class="astro-stat-label">Орбитальный радиус</span>
                  <span>{{ formatDistanceKm(selectedMoon.orbitalRadiusKm) }}</span>
                </div>
                <div class="astro-stat-row">
                  <span class="astro-stat-label">Период обращения</span>
                  <span>{{ formatPeriodDays(selectedMoon.orbitalPeriodDays) }}</span>
                </div>
                <div class="astro-stat-row">
                  <span class="astro-stat-label">Планета</span>
                  <span>{{ selectedPlanet.name }}</span>
                </div>
              </div>
            </article>
            <AstronomyModelSpoiler
              v-if="selectedMoon?.modelSrc"
              :title="`3D модель: ${selectedMoon.name}`"
              :viewer-title="selectedMoon.name"
              :model-src="selectedMoon.modelSrc"
              :model-size-label="selectedMoon.modelSizeLabel"
              :rotation-per-second="selectedMoon.modelRotationPerSecond"
              hint="Пока удалось найти готовую GLB-модель только для Луны, поэтому этот блок показывается не для всех спутников." />
          </aside>
        </div>
      </div>
    </section>
    <section id="astronomy-sizes" class="astro-section">
      <div class="astro-section__body astro-grid-visual">
        <div class="astro-stack-md astro-main-column">
          <div class="astro-stack-sm">
            <h2 class="astro-heading">Диаметры Солнца, планет и крупных спутников</h2>
            <p class="astro-paragraph astro-paragraph-muted">
              Полосы ниже показывают реальные диаметры тел, нормированные по диаметру Солнца. У маленьких спутников
              полосы почти исчезают на экране, и это ожидаемый результат реального масштаба.
            </p>
          </div>
          <div class="astro-size-list">
            <button
              v-for="body in sizeBodies"
              :key="body.id"
              type="button"
              class="astro-size-row"
              :class="body.id === selectedSizeBody.id ? 'astro-size-row-active' : ''"
              @mouseenter="onSizeBodyFocus(body)"
              @click="onSizeBodyFocus(body)">
              <span class="astro-size-label">{{ body.name }}</span>
              <span class="astro-size-track">
                <span
                  class="astro-size-bar"
                  :style="{
                    width: `${sizeRatioPercent(body.diameterKm)}%`,
                    backgroundColor: body.color
                  }"></span>
              </span>
              <span class="astro-size-value">{{ formatDiameter(body.diameterKm) }}</span>
            </button>
          </div>
          <div class="astro-mobile-slider">
            <article
              v-for="body in sizeBodies"
              :key="`${body.id}-mobile`"
              class="astro-mobile-slide"
              @click="onSizeBodyFocus(body)">
              <div class="astro-stack-sm">
                <p class="astro-kicker">{{ body.kind }}</p>
                <h3 class="astro-heading-sm">{{ body.name }}</h3>
              </div>
              <div class="astro-mobile-track">
                <div
                  class="astro-mobile-bar"
                  :style="{
                    width: `${sizeRatioPercent(body.diameterKm)}%`,
                    backgroundColor: body.color
                  }"></div>
              </div>
              <div class="astro-stat-list">
                <div class="astro-stat-row">
                  <span class="astro-stat-label">Диаметр</span>
                  <span>{{ formatDiameter(body.diameterKm) }}</span>
                </div>
                <div class="astro-stat-row">
                  <span class="astro-stat-label">Относительно Солнца</span>
                  <span>{{ decimalFormatter.format(body.diameterKm / sun.diameterKm) }}</span>
                </div>
                <div class="astro-stat-row">
                  <span class="astro-stat-label">Относительно Земли</span>
                  <span>{{ decimalFormatter.format(body.diameterKm / earthPlanet.diameterKm) }}</span>
                </div>
              </div>
              <p class="astro-note wrap-break-word">{{ body.description }}</p>
              <AstronomyModelSpoiler
                v-if="body.modelSrc"
                :title="`3D модель: ${body.name}`"
                :viewer-title="body.name"
                :model-src="body.modelSrc"
                :model-size-label="body.modelSizeLabel"
                :rotation-per-second="body.modelRotationPerSecond"
                hint="Просмотр модели открывается по требованию и не загружается вместе со слайдом." />
            </article>
          </div>
        </div>
        <aside class="astro-side-column astro-side-column-desktop-only">
          <article class="astro-card">
            <p class="astro-kicker">Выбранное тело</p>
            <p class="astro-card__value">{{ selectedSizeBody.name }}</p>
            <p class="astro-note wrap-break-word">{{ selectedSizeBody.description }}</p>
          </article>
          <AstronomyModelSpoiler
            v-if="selectedSizeBody.modelSrc"
            :title="`3D модель: ${selectedSizeBody.name}`"
            :viewer-title="selectedSizeBody.name"
            :model-src="selectedSizeBody.modelSrc"
            :model-size-label="selectedSizeBody.modelSizeLabel"
            :rotation-per-second="selectedSizeBody.modelRotationPerSecond"
            hint="Здесь доступны модели Солнца, планет и Луны. Остальные спутники добавлю, когда появятся подходящие GLB-файлы." />
          <article class="astro-card">
            <div class="astro-stat-list">
              <div class="astro-stat-row">
                <span class="astro-stat-label">Диаметр</span>
                <span>{{ formatDiameter(selectedSizeBody.diameterKm) }}</span>
              </div>
              <div class="astro-stat-row">
                <span class="astro-stat-label">Относительно Солнца</span>
                <span>{{ decimalFormatter.format(sizeBodyRelationToSun) }}</span>
              </div>
              <div class="astro-stat-row">
                <span class="astro-stat-label">Относительно Земли</span>
                <span>{{ decimalFormatter.format(sizeBodyRelationToEarth) }}</span>
              </div>
            </div>
          </article>
        </aside>
      </div>
    </section>
  </div>
</template>
<style>
  :root {
    --astro-page-max-width: 108rem;
    --astro-gap: 1rem;
    --astro-gap-lg: 1.25rem;
    --astro-radius: 0;
    --astro-panel-bg: radial-gradient(
      circle at top,
      color-mix(in srgb, var(--lab-bg-surface) 98%, transparent),
      color-mix(in srgb, var(--lab-bg-canvas) 96%, transparent)
    );
    --astro-line: color-mix(in srgb, var(--lab-border) 88%, transparent);
    --astro-line-muted: color-mix(in srgb, var(--lab-border) 48%, transparent);
    --astro-line-strong: color-mix(in srgb, var(--lab-border-strong) 88%, transparent);
    --astro-card-bg: color-mix(in srgb, var(--lab-bg-surface) 94%, transparent);
    --astro-card-bg-hover: color-mix(in srgb, var(--lab-bg-surface-hover) 92%, transparent);
  }
  .astro-page {
    width: 100%;
  }
  .astro-sticky-nav {
    position: sticky;
    top: 0;
    z-index: 20;
    margin-bottom: 1rem;
    border-bottom: 1px solid var(--lab-border);
    background: color-mix(in srgb, var(--lab-bg-overlay) 96%, transparent);
    backdrop-filter: blur(8px);
  }
  .astro-sticky-nav__inner {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    padding: 0.75rem 0.75rem;
  }
  .astro-section {
    width: min(100%, var(--astro-page-max-width));
    margin: 0 auto 1rem;
    border: 1px solid var(--lab-border);
    background: color-mix(in srgb, var(--lab-bg-surface-muted) 92%, transparent);
    scroll-margin-top: 4rem;
  }
  .astro-section__body {
    padding: 1rem;
  }
  .astro-grid-2,
  .astro-grid-visual,
  .astro-visual-grid {
    display: grid;
    gap: var(--astro-gap);
  }
  .astro-stack-sm {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  .astro-stack-md {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .astro-stack-lg {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }
  .astro-main-column {
    min-width: 0;
  }
  .astro-side-column {
    display: grid;
    gap: 0.75rem;
    align-content: start;
    min-width: 0;
  }
  .astro-heading,
  .astro-heading-sm {
    margin: 0;
    color: var(--lab-text-primary);
    font-weight: 600;
    line-height: 1.2;
  }
  .astro-heading {
    font-size: 1.125rem;
  }
  .astro-heading-sm {
    font-size: 1rem;
  }
  .astro-paragraph,
  .astro-note {
    margin: 0;
    line-height: 1.7;
  }
  .astro-paragraph {
    font-size: 0.9375rem;
    color: var(--lab-text-secondary);
  }
  .astro-paragraph-muted,
  .astro-note {
    color: var(--lab-text-muted);
  }
  .astro-kicker {
    margin: 0;
    color: var(--lab-text-soft);
    font-size: 0.75rem;
    line-height: 1.2;
    letter-spacing: 0.06em;
    text-transform: uppercase;
  }
  .astro-kicker-strong {
    color: color-mix(in srgb, var(--lab-text-primary) 82%, transparent);
  }
  .astro-metric-grid {
    display: grid;
    grid-template-columns: repeat(1, minmax(0, 1fr));
    gap: 0.75rem;
  }
  .astro-facts-grid {
    display: grid;
    gap: 0.75rem;
  }
  .astro-color-card,
  .astro-card {
    min-width: 0;
    border: 1px solid var(--lab-border);
    background: var(--astro-card-bg);
    padding: 0.875rem;
  }
  .astro-color-card__head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 0.75rem;
  }
  .astro-color-card__symbol {
    color: var(--lab-text-primary);
    font-size: 1.25rem;
    font-weight: 600;
    line-height: 1;
  }
  .astro-color-card__value,
  .astro-card__value {
    margin: 0.5rem 0 0;
    color: var(--lab-text-primary);
    font-size: 1rem;
    font-weight: 600;
    line-height: 1.3;
    word-wrap: anywhere;
  }
  .astro-toolbar {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 0.5rem;
  }
  .astro-zoom-label {
    color: var(--lab-text-soft);
    font-size: 0.75rem;
  }
  .astro-panel {
    min-width: 0;
    border: 1px solid var(--lab-border);
    background: var(--astro-panel-bg);
    padding: 0.75rem;
  }
  .astro-panel-visual {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 18rem;
  }
  .astro-panel-galaxy {
    min-height: clamp(18rem, 52vw, 34rem);
  }
  .astro-panel-map {
    min-height: clamp(18rem, 54vw, 36rem);
  }
  .astro-svg {
    display: block;
    width: 100%;
    height: 100%;
    max-height: 100%;
  }
  .astro-hit-area {
    cursor: pointer;
  }
  .astro-reticle {
    pointer-events: none;
  }
  .astro-planet-label {
    fill: var(--lab-text-primary);
    font-weight: 600;
    line-height: 1;
    paint-order: stroke;
    stroke: color-mix(in srgb, var(--lab-bg-canvas) 92%, transparent);
    stroke-width: 0.08em;
    pointer-events: none;
  }
  .astro-stat-list {
    display: grid;
    gap: 0.5rem;
    color: var(--lab-text-secondary);
    font-size: 0.875rem;
  }
  .astro-stat-row {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 0.75rem;
  }
  .astro-stat-label {
    color: var(--lab-text-soft);
  }
  .astro-empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 16rem;
    color: var(--lab-text-muted);
    font-size: 0.875rem;
    text-align: center;
  }
  .astro-chip-row {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-top: 0.75rem;
  }
  .astro-size-list {
    display: none;
  }
  .astro-size-row {
    display: grid;
    grid-template-columns: minmax(0, 6rem) minmax(0, 1fr) minmax(0, 5.25rem);
    align-items: center;
    gap: 0.75rem;
    width: 100%;
    border: 1px solid var(--lab-border);
    background: color-mix(in srgb, var(--lab-bg-control) 88%, transparent);
    padding: 0.75rem;
    text-align: left;
    color: var(--lab-text-primary);
    transition:
      background-color 160ms ease,
      border-color 160ms ease;
  }
  .astro-size-row:hover,
  .astro-size-row-active {
    border-color: var(--lab-border-strong);
    background: var(--astro-card-bg-hover);
  }
  .astro-size-label {
    min-width: 0;
    color: var(--lab-text-secondary);
    font-size: 0.875rem;
    word-wrap: anywhere;
  }
  .astro-size-track,
  .astro-mobile-track {
    display: flex;
    align-items: center;
    width: 100%;
    background: color-mix(in srgb, var(--lab-bg-surface-subtle) 92%, transparent);
  }
  .astro-size-track {
    height: 0.75rem;
  }
  .astro-mobile-track {
    height: 1rem;
  }
  .astro-size-bar,
  .astro-mobile-bar {
    display: block;
    max-width: 100%;
  }
  .astro-size-bar {
    height: 0.75rem;
  }
  .astro-mobile-bar {
    height: 1rem;
  }
  .astro-size-value {
    color: var(--lab-text-muted);
    font-size: 0.75rem;
    text-align: right;
  }
  .astro-mobile-slider {
    display: flex;
    gap: 0.75rem;
    overflow-x: auto;
    scroll-snap-type: x mandatory;
    padding-bottom: 0.25rem;
  }
  .astro-mobile-slide {
    display: flex;
    flex: 0 0 calc(100% - 0.5rem);
    flex-direction: column;
    gap: 1rem;
    scroll-snap-align: center;
    border: 1px solid var(--lab-border);
    background: var(--astro-card-bg);
    padding: 1rem;
  }
  .astro-side-column-desktop-only {
    display: none;
  }
  @media (min-width: 640px) {
    .astro-section__body {
      padding: 1.25rem;
    }
    .astro-sticky-nav__inner {
      padding-inline: 1rem;
    }
    .astro-metric-grid {
      grid-template-columns: repeat(3, minmax(0, 1fr));
    }
    .astro-facts-grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
  }
  @media (min-width: 1024px) {
    .astro-section {
      margin-bottom: 1.25rem;
    }
    .astro-grid-2 {
      grid-template-columns: minmax(0, 1.55fr) minmax(19rem, 0.95fr);
      align-items: start;
      gap: var(--astro-gap-lg);
    }
    .astro-grid-visual {
      grid-template-columns: minmax(0, 1.85fr) minmax(18rem, 0.95fr);
      align-items: start;
      gap: var(--astro-gap-lg);
    }
    .astro-visual-grid {
      grid-template-columns: minmax(0, 2fr) minmax(18rem, 1fr);
      align-items: start;
      gap: var(--astro-gap-lg);
    }
    .astro-size-list {
      display: grid;
      gap: 0.5rem;
    }
    .astro-mobile-slider {
      display: none;
    }
    .astro-side-column-desktop-only {
      display: grid;
    }
  }
</style>
