export type DocMeta = {
  slug: string
  title: string
  section: string
  summary: string
}
export type DocRecord = DocMeta & {
  file: string
}
export type ContentDoc = DocMeta & {
  content: string
}
const CONTENT_LOADERS = import.meta.glob('../content/**/*.md', {
  query: '?raw',
  import: 'default'
}) as Record<string, () => Promise<string>>
type ContentModuleOptions = {
  dir: string
  index: readonly DocRecord[]
}
const stripFile = (record: DocRecord): DocMeta => {
  const { file: _file, ...meta } = record
  return meta
}
const withContent = (record: DocRecord, content: string): ContentDoc => {
  const { file: _file, ...meta } = record
  return { ...meta, content }
}
export const createContentModule = ({ dir, index }: ContentModuleOptions) => {
  const contentCache = new Map<string, string>()
  const findDoc = (slug: string): DocMeta | null => {
    const item = index.find(record => record.slug === slug)
    return item ? stripFile(item) : null
  }
  const docs: DocMeta[] = index.map(stripFile)
  const loadDoc = async (slug: string): Promise<ContentDoc | null> => {
    const item = index.find(record => record.slug === slug)
    if (!item) return null
    const cached = contentCache.get(slug)
    if (cached !== undefined) {
      return withContent(item, cached)
    }
    const loaderKey = `../content/${dir}/${item.file}`
    const loader = CONTENT_LOADERS[loaderKey]
    if (!loader) return null
    const content = String(await loader())
    contentCache.set(slug, content)
    return withContent(item, content)
  }
  return {
    findDoc,
    docs,
    loadDoc
  }
}
const EDU_DOC_INDEX: DocRecord[] = [
  {
    slug: 'algebra',
    title: 'Алгебра',
    section: 'Математика',
    summary: 'Тождества, степени, логарифмы, квадратные уравнения и неравенства.',
    file: 'algebra.md'
  },
  {
    slug: 'geometry-trig',
    title: 'Геометрия и тригонометрия',
    section: 'Математика',
    summary: 'Треугольники, окружности, основные тригонометрические формулы.',
    file: 'geometry_trig.md'
  },
  {
    slug: 'calculus-linear-prob',
    title: 'Анализ, линейная алгебра, вероятность',
    section: 'Математика',
    summary: 'Производные, интегралы, матрицы, статистика и базовые вероятности.',
    file: 'calculus_linear_prob.md'
  },
  {
    slug: 'physics-core',
    title: 'Физика: ядро формул',
    section: 'Физика',
    summary: 'Механика, термодинамика, электричество и оптика.',
    file: 'physics_core.md'
  },
  {
    slug: 'materials',
    title: 'Материаловедение',
    section: 'Материаловедение',
    summary: 'Свойства материалов, термообработка и выбор материала.',
    file: 'materials.md'
  },
  {
    slug: 'mechanics-strength',
    title: 'Теоретическая механика и сопромат',
    section: 'Механика и расчеты',
    summary: 'Равновесие, динамика, напряжения, деформации и усталость.',
    file: 'mechanics_strength.md'
  },
  {
    slug: 'manufacturing',
    title: 'Технология производства',
    section: 'Производство',
    summary: 'Резание, литье, сварка, аддитив и проектирование процесса.',
    file: 'manufacturing.md'
  },
  {
    slug: 'design-standards',
    title: 'Проектирование и стандарты',
    section: 'Проектирование',
    summary: 'Чертежи, допуски, посадки, CAD-практика.',
    file: 'design_standards.md'
  },
  {
    slug: 'electrical',
    title: 'Электротехника и автоматизация',
    section: 'Электротехника',
    summary: 'Законы цепей, AC, мощность, микроконтроллеры и PID.',
    file: 'electrical.md'
  },
  {
    slug: 'economics-quality',
    title: 'Экономика и качество',
    section: 'Организация производства',
    summary: 'Себестоимость, точка безубыточности, OEE, PDCA, SPC.',
    file: 'economics_quality.md'
  },
  {
    slug: 'safety-ergonomics',
    title: 'Безопасность и эргономика',
    section: 'Безопасность',
    summary: 'Риск-оценка, LOTO, иерархия мер и эргономика рабочего места.',
    file: 'safety_ergonomics.md'
  }
]
const DEVELOPMENT_DOC_INDEX: DocRecord[] = [
  {
    slug: 'algorithms-core',
    title: 'Алгоритмы: базовая картина',
    section: 'Алгоритмы',
    summary: 'Что такое алгоритмы, как оценивают их качество и где применяют в реальных системах.',
    file: 'algorithms_core.md'
  },
  {
    slug: 'sorting-searching',
    title: 'Сортировка и поиск',
    section: 'Алгоритмы',
    summary: 'Ключевые подходы к сортировке и поиску, практические сценарии и типовые ошибки.',
    file: 'sorting_searching.md'
  },
  {
    slug: 'graphs-dynamic',
    title: 'Графы и динамическое программирование',
    section: 'Алгоритмы',
    summary: 'Когда использовать BFS/DFS/Dijkstra и как мыслить через состояния в DP.',
    file: 'graphs_dynamic.md'
  },
  {
    slug: 'data-representation',
    title: 'Представление информации',
    section: 'Данные',
    summary: 'Форматы данных, структуры и схемы: как они влияют на качество и производительность.',
    file: 'data_representation.md'
  },
  {
    slug: 'system-design-basics',
    title: 'Базовый системный дизайн',
    section: 'Практика разработки',
    summary: 'Компоненты, архитектурные варианты и минимальный набор для production.',
    file: 'system_design_basics.md'
  },
  {
    slug: 'public-network-security',
    title: 'Безопасность в общедоступных сетях',
    section: 'Безопасность',
    summary: 'Риски публичных сетей Wi-Fi и практические меры защиты для пользователя и команды.',
    file: 'public_network_security.md'
  },
  {
    slug: 'privacy-confidentiality',
    title: 'Конфиденциальность и цифровая гигиена',
    section: 'Безопасность',
    summary: 'Базовые принципы защиты персональных данных и снижения цифрового следа.',
    file: 'privacy_confidentiality.md'
  },
  {
    slug: 'secure-development-checklist',
    title: 'Secure SDLC: минимальный чеклист',
    section: 'Безопасность',
    summary: 'Что делать до, во время и после разработки, чтобы снизить риск инцидентов.',
    file: 'secure_development_checklist.md'
  },
  {
    slug: 'dev-scaling-ml',
    title: 'Масштабирование хранения данных, обработки и ML в продакшене',
    section: 'Масштабирование',
    summary:
      'Статья объясняет, как масштабировать продакшен-системы с машинным обучением: хранение данных, обработку, инференс и обучение моделей. Разбираются базовые архитектурные подходы — шардинг, репликация, микросервисы, batch и stream processing, data lake, data warehouse и lambda architecture — и их компромиссы по скорости, консистентности и отказоустойчивости.',
    file: 'dev_scaling.md'
  },
  {
    slug: 'docker-quick-start',
    title: 'Docker: Быстрый старт',
    section: 'Масштабирование',
    summary:
      'Краткий справочник по Docker: установка на Ubuntu, очистка системы, базовые команды управления контейнерами и образами, работа с сетью и Docker Compose, а также минимальный шаблон Dockerfile.',
    file: 'docker.md'
  }
]
const IMPACT_DOC_INDEX: DocRecord[] = [
  {
    slug: 'minimalism-lifestyle',
    title: 'Минимализм и разумное потребление',
    section: 'Личные привычки',
    summary: 'Как снижать нагрузку на природу через осознанные покупки и отказ от избыточного потребления.',
    file: 'minimalism_lifestyle.md'
  },
  {
    slug: 'hygiene-daily-ecology',
    title: 'Экологическая гигиена в быту',
    section: 'Личные привычки',
    summary: 'Ежедневные действия для чистой домашней среды и уменьшения экологического следа.',
    file: 'hygiene_daily_ecology.md'
  },
  {
    slug: 'water-energy-efficiency',
    title: 'Бережное использование воды и энергии',
    section: 'Ресурсы',
    summary: 'Практики эффективного потребления ресурсов дома, в учебе и на работе.',
    file: 'water_energy_efficiency.md'
  },
  {
    slug: 'circular-reuse-repair',
    title: 'Повторное использование и ремонт',
    section: 'Ресурсы',
    summary: 'Циркулярный подход: ремонт, переиспользование и продление жизненного цикла вещей.',
    file: 'circular_reuse_repair.md'
  },
  {
    slug: 'biodiversity-support',
    title: 'Поддержка видового разнообразия',
    section: 'Природа и экосистемы',
    summary: 'Локальные действия для сохранения биоразнообразия в городе и пригороде.',
    file: 'biodiversity_support.md'
  },
  {
    slug: 'urban-green-mobility',
    title: 'Городская мобильность и зеленая инфраструктура',
    section: 'Природа и экосистемы',
    summary: 'Позитивное влияние пешеходной среды, транспорта и озеленения на качество жизни.',
    file: 'urban_green_mobility.md'
  },
  {
    slug: 'community-initiatives',
    title: 'Локальные инициативы и взаимопомощь',
    section: 'Общество',
    summary: 'Как общественные практики формируют экологичную и устойчивую среду.',
    file: 'community_initiatives.md'
  },
  {
    slug: 'education-values',
    title: 'Экологические ценности и образование',
    section: 'Общество',
    summary: 'Роль образования, семейных практик и просвещения в долгосрочных позитивных изменениях.',
    file: 'education_values.md'
  }
]
const UPDATES_DOC_INDEX: DocRecord[] = [
  {
    slug: '2026-03-04-prerelease-functional-overview',
    title: 'Новое обновление',
    section: 'Интерфейс',
    summary: 'Краткое описание',
    file: '2026-03-04-prerelease-functional-overview.md'
  },
  {
    slug: '2026-03-11-functional-site-capabilities',
    title: 'Основной функционал',
    section: 'API',
    summary: 'Краткое описание',
    file: '2026-03-11-functional-site-capabilities.md'
  },
  {
    slug: '2026-03-16-finance',
    title: 'Финансы поют романсы',
    section: 'DoNation',
    summary: 'финансовая поддержка проекта',
    file: '2026-03-16-finance.md'
  }
]
const contentModules = {
  edu: createContentModule({ dir: 'edu', index: EDU_DOC_INDEX }),
  dev: createContentModule({ dir: 'dev', index: DEVELOPMENT_DOC_INDEX }),
  impact: createContentModule({ dir: 'impact', index: IMPACT_DOC_INDEX }),
  updates: createContentModule({ dir: 'updates', index: UPDATES_DOC_INDEX })
} as const
export const { edu, dev, impact, updates } = contentModules
