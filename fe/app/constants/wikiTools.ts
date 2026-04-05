export type WikiToolCategory = 'chemistry' | 'astronomy' | 'engineering'
export type WikiToolArt =
  | 'elements'
  | 'crystal'
  | 'orbit'
  | 'compounds'
  | 'isotopes'
  | 'materials'
  | 'products'
  | 'production'
export interface WikiTool {
  title: string
  text: string
  to: string
  icon: string
  category: WikiToolCategory
  categoryLabel: string
  accentColor: string
  accentSoftColor: string
  art: WikiToolArt
}
export const wikiTools: WikiTool[] = [
  {
    title: 'Периодическая система',
    text: 'Интерактивная таблица химических элементов с карточками и справочной информацией.',
    to: '/edu/chemistry/elements',
    icon: 'ic:round-science',
    category: 'chemistry',
    categoryLabel: 'Химия',
    accentColor: '#0ea5e9',
    accentSoftColor: '#67e8f9',
    art: 'compounds'
  },
  {
    title: 'Минералы',
    text: 'Каталог минералов с поиском и химическими фильтрами по составу и признакам.',
    to: '/edu/chemistry/minerals',
    icon: 'ic:round-diamond',
    category: 'chemistry',
    categoryLabel: 'Химия',
    accentColor: '#8b5cf6',
    accentSoftColor: '#22d3ee',
    art: 'crystal'
  },
  {
    title: 'Соединения',
    text: 'Заготовка для справочника химических соединений, их свойств и классификаций.',
    to: '/edu/chemistry/compounds',
    icon: 'ic:round-bubble-chart',
    category: 'chemistry',
    categoryLabel: 'Химия',
    accentColor: '#06b6d4',
    accentSoftColor: '#67e8f9',
    art: 'isotopes'
  },
  {
    title: 'Изотопы',
    text: 'Заготовка для каталога изотопов, стабильности, распадов и прикладных данных.',
    to: '/edu/chemistry/isotopes',
    icon: 'ic:round-scatter-plot',
    category: 'chemistry',
    categoryLabel: 'Химия',
    accentColor: '#3b82f6',
    accentSoftColor: '#93c5fd',
    art: 'elements'
  },
  {
    title: 'Астрономия',
    text: 'Обзор Солнечной системы, галактики и сравнительных масштабов космических объектов.',
    to: '/edu/astronomy',
    icon: 'ic:round-auto-awesome',
    category: 'astronomy',
    categoryLabel: 'Астрономия',
    accentColor: '#f97316',
    accentSoftColor: '#facc15',
    art: 'orbit'
  },
  {
    title: 'Инженерные материалы',
    text: 'Заготовка для материаловедения: металлы, полимеры, керамика, композиты и их свойства.',
    to: '/edu/materials',
    icon: 'ic:round-category',
    category: 'engineering',
    categoryLabel: 'Инженерия',
    accentColor: '#f97316',
    accentSoftColor: '#fdba74',
    art: 'materials'
  },
  {
    title: 'Изделия',
    text: 'Заготовка для инженерных изделий, узлов, стандартных форм и примеров применения.',
    to: '/edu/products',
    icon: 'ic:round-inventory-2',
    category: 'engineering',
    categoryLabel: 'Инженерия',
    accentColor: '#eab308',
    accentSoftColor: '#fde047',
    art: 'products'
  },
  {
    title: 'Производство',
    text: 'Заготовка для процессов изготовления, маршрутов обработки и производственных технологий.',
    to: '/edu/production',
    icon: 'ic:round-precision-manufacturing',
    category: 'engineering',
    categoryLabel: 'Инженерия',
    accentColor: '#ef4444',
    accentSoftColor: '#fca5a5',
    art: 'production'
  }
]
