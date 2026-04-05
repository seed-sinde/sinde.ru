const rawElementModules = import.meta.glob('../../public/chemistry/elements/*/meta.json', {
  eager: true,
  import: 'default'
}) as Record<string, RawPeriodicTableElement>
const rawSampleModules = import.meta.glob('../../public/chemistry/elements/*/samples/meta.json', {
  eager: true,
  import: 'default'
}) as Record<string, RawPeriodicTableSamplesMeta>
const CATEGORY_LABELS: Record<string, string> = {
  actinide: 'Актиноиды',
  'alkali metal': 'Щелочные металлы',
  'alkaline earth metal': 'Щёлочноземельные металлы',
  'diatomic nonmetal': 'Двухатомные неметаллы',
  lanthanide: 'Лантаноиды',
  metalloid: 'Металлоиды',
  'noble gas': 'Благородные газы',
  'polyatomic nonmetal': 'Многоатомные неметаллы',
  'post-transition metal': 'Постпереходные металлы',
  'transition metal': 'Переходные металлы',
  'unknown, but predicted to be an alkali metal': 'Предположительно щелочные металлы',
  'unknown, predicted to be noble gas': 'Предположительно благородные газы',
  'unknown, probably metalloid': 'Предположительно металлоиды',
  'unknown, probably post-transition metal': 'Предположительно постпереходные металлы',
  'unknown, probably transition metal': 'Предположительно переходные металлы'
}
const CATEGORY_COLORS: Record<string, string> = {
  actinide: '#ef4444',
  'alkali metal': '#f97316',
  'alkaline earth metal': '#f59e0b',
  'diatomic nonmetal': '#22c55e',
  lanthanide: '#ec4899',
  metalloid: '#14b8a6',
  'noble gas': '#38bdf8',
  'polyatomic nonmetal': '#10b981',
  'post-transition metal': '#a78bfa',
  'transition metal': '#6366f1',
  'unknown, but predicted to be an alkali metal': '#f97316',
  'unknown, predicted to be noble gas': '#38bdf8',
  'unknown, probably metalloid': '#14b8a6',
  'unknown, probably post-transition metal': '#a78bfa',
  'unknown, probably transition metal': '#6366f1'
}
const CATEGORY_LIGHT_COLORS: Record<string, string> = {
  actinide: '#f4cccc',
  'alkali metal': '#f9d5c2',
  'alkaline earth metal': '#f6e0bf',
  'diatomic nonmetal': '#d8ead3',
  lanthanide: '#f4d3e6',
  metalloid: '#d4ebe5',
  'noble gas': '#d8e9f8',
  'polyatomic nonmetal': '#dcefd8',
  'post-transition metal': '#e2daf6',
  'transition metal': '#d8def8',
  'unknown, but predicted to be an alkali metal': '#f8dfd0',
  'unknown, predicted to be noble gas': '#e1edf8',
  'unknown, probably metalloid': '#e0efea',
  'unknown, probably post-transition metal': '#ebe4f7',
  'unknown, probably transition metal': '#e2e7f8'
}
const CATEGORY_DESCRIPTIONS: Record<string, string> = {
  actinide:
    'Актиноиды образуют тяжёлый f-блок с выраженной радиоактивностью. Для них характерны сложные электронные конфигурации и высокая химическая вариативность.',
  'alkali metal':
    'Щелочные металлы находятся в первой группе, быстро реагируют с водой и образуют мягкие металлы с одним внешним электроном.',
  'alkaline earth metal':
    'Щёлочноземельные металлы составляют вторую группу. Обычно они твёрже щелочных металлов и проявляют устойчивую валентность +2.',
  'diatomic nonmetal':
    'Двухатомные неметаллы часто существуют в виде молекул из двух атомов и обладают высокой химической активностью в обычных условиях.',
  lanthanide:
    'Лантаноиды объединяют элементы f-блока с близкими химическими свойствами. Многие из них используются в магнитах, люминофорах и специальных сплавах.',
  metalloid:
    'Металлоиды занимают промежуточное положение между металлами и неметаллами. Их свойства особенно важны для полупроводников и материаловедения.',
  'noble gas':
    'Благородные газы отличаются заполненной внешней электронной оболочкой, поэтому химически инертны и часто используются в защитных и световых средах.',
  'polyatomic nonmetal':
    'Многоатомные неметаллы образуют разнообразные ковалентные структуры и лежат в основе органической химии, атмосферы и биологических систем.',
  'post-transition metal':
    'Постпереходные металлы мягче и химически проще классических переходных металлов. Они часто встречаются в низкоплавких сплавах и функциональных покрытиях.',
  'transition metal':
    'Переходные металлы занимают центральную часть таблицы. Для них типичны высокая прочность, переменные степени окисления и каталитическая активность.',
  'unknown, but predicted to be an alkali metal':
    'Эта категория описывает сверхтяжёлые элементы, которые по расчётам должны вести себя как щелочные металлы, но экспериментальных данных пока недостаточно.',
  'unknown, predicted to be noble gas':
    'Предполагаемые благородные газы находятся на границе подтверждённых данных: расчёты указывают на инертность, но свойства ещё требуют уточнения.',
  'unknown, probably metalloid':
    'Для этих элементов ожидаются свойства, близкие к металлоидам, однако экспериментальная база пока ограничена или отсутствует.',
  'unknown, probably post-transition metal':
    'Элементы этой категории прогнозируются как постпереходные металлы, но из-за малой стабильности их химия описана лишь частично.',
  'unknown, probably transition metal':
    'Предполагаемые переходные металлы сверхтяжёлой области таблицы пока известны главным образом по расчётам и короткоживущим экспериментам.'
}
const extractSymbolFromPath = (modulePath: string) => {
  const match = modulePath.match(/elements\/([^/]+)\//)
  return match?.[1] || ''
}
const getCategoryLabel = (category: string) => {
  return CATEGORY_LABELS[category] || category
}
const getCategoryColor = (category: string) => {
  return CATEGORY_COLORS[category] || '#52525b'
}
const getCategoryLightColor = (category: string) => {
  return CATEGORY_LIGHT_COLORS[category] || '#e4e4e7'
}
const normalizeText = (value: unknown) => {
  const text = String(value ?? '').trim()
  return text || null
}
export const formatPeriodicTableSymbol = (value: unknown) => {
  const symbol = String(value ?? '').trim()
  if (!symbol) return ''
  return `${symbol.slice(0, 1).toUpperCase()}${symbol.slice(1).toLowerCase()}`
}
const normalizeNumber = (value: unknown) => {
  return typeof value === 'number' && Number.isFinite(value) ? value : null
}
const normalizeNumberList = (value: unknown) => {
  if (!Array.isArray(value)) return [] as number[]
  return value.filter((item): item is number => typeof item === 'number' && Number.isFinite(item))
}
const normalizeHex = (value: unknown) => {
  const hex = String(value ?? '')
    .trim()
    .replace(/^#/, '')
    .toLowerCase()
  return /^[0-9a-f]{6}$/i.test(hex) ? hex : null
}
const normalizeLookupKey = (value: unknown) => {
  const source = String(value ?? '').trim()
  if (!source) return ''
  return slugifyLatin(source.replace(/[/_]+/g, ' '), { fallback: 'element' })
}
const buildAliases = (raw: RawPeriodicTableElement, slug: string, symbol: string, number: number) => {
  const aliases = new Set<string>()
  const candidates = [
    slug,
    raw.name,
    symbol,
    `${number}`,
    `${number}-${symbol}`,
    `${symbol}-${number}`,
    `${raw.name || ''}-${symbol}`,
    `${symbol}-${raw.name || ''}`
  ]
  for (const candidate of candidates) {
    const key = normalizeLookupKey(candidate)
    if (key) aliases.add(key)
  }
  return Array.from(aliases)
}
const buildModelPath = (number: number, symbol: string, name: string) => {
  const num = String(number).padStart(3, '0')
  const normalizedName = slugifyLatin(name, {
    fallback: `element-${num}`
  }).toLowerCase()
  const normalizedSymbol = symbol.toLowerCase()
  return `/chemistry/elements/${normalizedSymbol}/${num}-${normalizedSymbol}-${normalizedName}-model.glb`
}
const normalizeFallback = (value: RawPeriodicTableSamplesMeta['fallback']) => {
  const type = String(value?.type || '').trim() as PeriodicTableSampleFallbackType
  const label = String(value?.label || '').trim()
  const description = String(value?.description || '').trim()
  if (!type || !label || !description) return null
  return {
    type,
    label,
    description
  }
}
const samplesBySymbol = new Map<string, RawPeriodicTableSamplesMeta>(
  Object.entries(rawSampleModules).map(([modulePath, rawSamples]) => [extractSymbolFromPath(modulePath), rawSamples])
)
const rawElements = Object.values(rawElementModules).filter(raw => {
  const number = normalizeNumber(raw.number)
  return number !== null && number >= 1 && number <= 118
})
export const PERIODIC_TABLE_ELEMENTS: PeriodicTableElement[] = rawElements
  .map((raw, index) => {
    const name = normalizeText(raw.name) || `Element ${index + 1}`
    const symbolRaw = normalizeText(raw.symbol) || `E${index + 1}`
    const symbol = formatPeriodicTableSymbol(symbolRaw) // Al
    const symbolLower = symbol.toLowerCase() // al
    const number = normalizeNumber(raw.number) ?? index + 1
    const slug = normalizeLookupKey(name)
    const cpkHex = normalizeHex(raw['cpk-hex'])
    const category = normalizeText(raw.category) || 'unknown'
    const sampleMeta = samplesBySymbol.get(symbolLower)
    const samples: PeriodicTableSample[] = Array.isArray(sampleMeta?.samples)
      ? sampleMeta.samples
          .map(sample => {
            const file = String(sample?.file || '').trim()
            if (!file) return null
            return {
              file,
              url: `/chemistry/elements/${symbolLower}/samples/${file}`,
              title: normalizeText(sample?.title),
              author: normalizeText(sample?.author),
              attribution: normalizeText(sample?.attribution),
              sourceUrl: normalizeText(sample?.source_url),
              licenseUrl: normalizeText(sample?.license_url),
              license: normalizeText(sample?.license)
            }
          })
          .filter((sample): sample is PeriodicTableSample => Boolean(sample))
      : []
    const element: PeriodicTableElement = {
      slug: slug || normalizeLookupKey(symbol) || `${number}`,
      aliases: [],
      number,
      name,
      russianName: normalizeText(raw.russian_name) || name,
      displaySymbol: symbol,
      symbol: symbol, // для логики и API
      appearance: normalizeText(raw.appearance),
      atomicMass: normalizeNumber(raw.atomic_mass),
      boil: normalizeNumber(raw.boil),
      category,
      categoryLabel: getCategoryLabel(category),
      categoryColor: getCategoryColor(category),
      categoryLightColor: getCategoryLightColor(category),
      density: normalizeNumber(raw.density),
      discoveredBy: normalizeText(raw.discovered_by),
      melt: normalizeNumber(raw.melt),
      molarHeat: normalizeNumber(raw.molar_heat),
      namedBy: normalizeText(raw.named_by),
      period: normalizeNumber(raw.period),
      group: normalizeNumber(raw.group),
      phase: normalizeText(raw.phase),
      source: normalizeText(raw.source),
      bohrModelImage: normalizeText(raw.bohr_model_image),
      bohrModel3d: buildModelPath(number, symbol, name),
      spectralImage: normalizeText(raw.spectral_img),
      summary: normalizeText(raw.summary) || '',
      xpos: normalizeNumber(raw.xpos) ?? 1,
      ypos: normalizeNumber(raw.ypos) ?? 1,
      wxpos: normalizeNumber(raw.wxpos),
      wypos: normalizeNumber(raw.wypos),
      shells: normalizeNumberList(raw.shells),
      electronConfiguration: normalizeText(raw.electron_configuration),
      electronConfigurationSemantic: normalizeText(raw.electron_configuration_semantic),
      electronAffinity: normalizeNumber(raw.electron_affinity),
      electronegativityPauling: normalizeNumber(raw.electronegativity_pauling),
      ionizationEnergies: normalizeNumberList(raw.ionization_energies),
      cpkHex,
      block: normalizeText(raw.block),
      accentColor: cpkHex ? `#${cpkHex}` : '#52525b',
      searchText: '',
      samples,
      primarySample: samples[0] || null,
      sampleFallback: normalizeFallback(sampleMeta?.fallback)
    }
    element.aliases = buildAliases(raw, element.slug, symbol, number)
    element.searchText = [
      element.name,
      element.russianName,
      element.symbol,
      `${element.number}`,
      element.category,
      element.categoryLabel,
      element.phase,
      element.block,
      element.electronConfiguration,
      element.electronConfigurationSemantic
    ]
      .filter(Boolean)
      .join(' ')
      .toLowerCase()
    return element
  })
  .sort((left, right) => {
    if (left.ypos !== right.ypos) return left.ypos - right.ypos
    if (left.xpos !== right.xpos) return left.xpos - right.xpos
    return left.number - right.number
  })
const ELEMENT_LOOKUP = new Map<string, PeriodicTableElement>()
for (const element of PERIODIC_TABLE_ELEMENTS) {
  for (const alias of element.aliases) {
    ELEMENT_LOOKUP.set(alias, element)
  }
}
export const PERIODIC_TABLE_GROUPS = Math.max(...PERIODIC_TABLE_ELEMENTS.map(element => element.xpos), 18)
export const PERIODIC_TABLE_ROWS = Math.max(...PERIODIC_TABLE_ELEMENTS.map(element => element.ypos), 10)
export const PERIODIC_TABLE_CATEGORY_COUNTS = Array.from(
  PERIODIC_TABLE_ELEMENTS.reduce((map, element) => {
    map.set(element.category, (map.get(element.category) || 0) + 1)
    return map
  }, new Map<string, number>())
)
  .sort((left, right) => {
    if (right[1] !== left[1]) return right[1] - left[1]
    return left[0].localeCompare(right[0])
  })
  .map(([category, count]) => ({
    category,
    label: getCategoryLabel(category),
    color: getCategoryColor(category),
    description: CATEGORY_DESCRIPTIONS[category] || '',
    count
  }))
export const getPeriodicTableElementRoute = (element: Pick<PeriodicTableElement, 'slug'>) =>
  `/edu/chemistry/elements/${element.slug}`
export const resolvePeriodicTableElement = (value: string | string[] | null | undefined) => {
  const rawSegments = Array.isArray(value)
    ? value
    : String(value ?? '')
        .split('/')
        .filter(Boolean)
  const segments = rawSegments.map(segment => decodeURIComponent(String(segment || '').trim())).filter(Boolean)
  if (!segments.length) return null
  const variants = [segments.join('-'), segments.join(' '), segments[segments.length - 1] || '']
  for (const variant of variants) {
    const key = normalizeLookupKey(variant)
    if (!key) continue
    const element = ELEMENT_LOOKUP.get(key)
    if (element) return element
  }
  return null
}
