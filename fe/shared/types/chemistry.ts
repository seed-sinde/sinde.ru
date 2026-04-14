import type { InterfaceLocaleCode, ThemePreference } from './ui'
export type ElementCategory = {
  slug: string
  label: string
  count: number
}
export type MineralImage = {
  file: string
  rruff_id: string
  order: number
}
export type Mineral = {
  id: number
  database_id: number
  mineral_name: string
  mineral_name_plain: string
  valence_chemistry?: string | null
  ima_chemistry?: string | null
  chemistry_elements?: string[]
  valence_elements?: string[]
  ima_number?: string | null
  rruff_ids?: string[]
  ima_mineral_symbol?: string | null
  ima_status?: string | null
  status_notes?: string | null
  structural_groupname?: string | null
  crystal_systems?: string[]
  space_groups?: string[]
  country_of_type_locality?: string | null
  year_first_published?: number | null
  oldest_known_age_ma?: number | null
  paragenetic_modes?: string[]
  created_at: string
  updated_at: string
  images?: MineralImage[]
}
export type MineralCrystalSystem =
  | 'cubic'
  | 'hexagonal'
  | 'monoclinic'
  | 'orthorhombic'
  | 'tetragonal'
  | 'triclinic'
  | 'unknown'
export type MineralCrystalSystemMode = 'any' | 'all'
export type MineralListItem = Pick<Mineral, 'database_id' | 'mineral_name' | 'crystal_systems'>
export type MineralsListMeta = {
  limit: number
  offset: number
  total: number
}
export type MineralsListFacets = {
  chemistryElementsAvailable: Record<string, number>
}
export type MineralsListPayload = {
  items: MineralListItem[]
  meta: MineralsListMeta
  facets: MineralsListFacets
}
export type ChemistryElementApiSample = {
  file?: string | null
  title?: string | null
  author?: string | null
  attribution?: string | null
  source_url?: string | null
  license_url?: string | null
  license?: string | null
}
export type ChemistryElementApi = {
  name?: string | null
  russian_name?: string | null
  appearance?: string | null
  atomic_mass?: number | null
  boil?: number | null
  category?: string | null
  density?: number | null
  discovered_by?: string | null
  melt?: number | null
  molar_heat?: number | null
  named_by?: string | null
  number?: number | null
  period?: number | null
  group?: number | null
  phase?: string | null
  source?: string | null
  bohr_model_image?: string | null
  bohr_model_3d?: string | null
  spectral_img?: string | null
  summary?: string | null
  symbol?: string | null
  xpos?: number | null
  ypos?: number | null
  wxpos?: number | null
  wypos?: number | null
  shells?: number[] | null
  electron_configuration?: string | null
  electron_configuration_semantic?: string | null
  electron_affinity?: number | null
  electronegativity_pauling?: number | null
  ionization_energies?: number[] | null
  cpk_hex?: string | null
  block?: string | null
  samples?: ChemistryElementApiSample[] | null
  sample_fallback?: Partial<PeriodicTableSampleFallback> | null
}
export type PeriodicTableSample = {
  file: string
  url: string
  title: string | null
  author: string | null
  attribution: string | null
  sourceUrl: string | null
  licenseUrl: string | null
  license: string | null
}
type MineralSortOption = 'name_asc' | 'name_desc'
export type MineralImageFilter = 'without' | 'any' | 'with'
export type MineralsFiltersSnapshot = {
  q: string
  sort: MineralSortOption
  limit: number
  offset: number
  imageFilter: MineralImageFilter
  crystalSystems: MineralCrystalSystem[]
  crystalSystemMode: MineralCrystalSystemMode
  chemistryAll: string[]
  chemistryAny: string[]
  chemistryNone: string[]
}
export type UiPreferencesSnapshot = {
  interfaceLocale: InterfaceLocaleCode
  themePreference: ThemePreference
  mineralsFilters: MineralsFiltersSnapshot
}
export type PeriodicTableSampleFallbackType = 'synthetic' | 'unstable' | 'hypothetical' | 'not-isolated' | 'unknown'
export type PeriodicTableSampleFallback = {
  type: PeriodicTableSampleFallbackType
  label: string
  description: string
}
export type PeriodicTableElement = {
  slug: string
  aliases: string[]
  number: number
  name: string
  russianName: string
  symbol: string
  displaySymbol: string
  appearance: string | null
  atomicMass: number | null
  boil: number | null
  category: string
  categoryLabel: string
  categoryColor: string
  categoryLightColor: string
  density: number | null
  discoveredBy: string | null
  melt: number | null
  molarHeat: number | null
  namedBy: string | null
  period: number | null
  group: number | null
  phase: string | null
  source: string | null
  bohrModelImage: string | null
  bohrModel3d: string | null
  spectralImage: string | null
  summary: string
  xpos: number
  ypos: number
  wxpos: number | null
  wypos: number | null
  shells: number[]
  electronConfiguration: string | null
  electronConfigurationSemantic: string | null
  electronAffinity: number | null
  electronegativityPauling: number | null
  ionizationEnergies: number[]
  cpkHex: string | null
  block: string | null
  accentColor: string
  searchText: string
  samples: PeriodicTableSample[]
  primarySample: PeriodicTableSample | null
  sampleFallback: PeriodicTableSampleFallback | null
}
export type ViewerImageKind = 'spectral' | 'bohr' | 'sample'
export type ViewerImageItem = {
  key: string
  kind: ViewerImageKind
  src: string
  title: string
  alt: string
  author: string | null
  attribution: string | null
  sourceUrl: string | null
  licenseUrl: string | null
  license: string | null
  thumbnailSrc: string | null
}
export type SamplePreviewItem = ViewerImageItem & {
  viewerIndex: number
}
export type DetailItem = {
  label: string
  value: string
  href?: string | null
}
export type DetailSection = {
  title: string
  items: DetailItem[]
}
export type ChemistryBucket = 'all' | 'any' | 'none'
export type MineralsRouteState = {
  q: string
  sort: MineralSortOption
  limit: number
  offset: number
  imageFilter: MineralImageFilter
  crystalSystems: MineralCrystalSystem[]
  crystalSystemMode: MineralCrystalSystemMode
  chemistryAll: string[]
  chemistryAny: string[]
  chemistryNone: string[]
}
export type ChemistryBucketMeta = {
  title: string
  description: string
  accentClass: string
  borderClass: string
  activeClass: string
  activeTextClass: string
  dotClass: string
}
