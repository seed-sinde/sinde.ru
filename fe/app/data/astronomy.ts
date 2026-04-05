export type AstronomyPlanet = {
  id: string
  symbol: string
  name: string
  kind: 'planet'
  diameterKm: number
  orbitalRadiusAu: number
  orbitalPeriodDays: number
  color: string
  description: string
  modelSrc?: string | null
  modelSizeLabel?: string | null
  modelRotationPerSecond?: string | null
}
export type AstronomyMoon = {
  id: string
  planetId: string
  symbol: string
  name: string
  kind: 'moon'
  diameterKm: number
  orbitalRadiusKm: number
  orbitalPeriodDays: number
  color: string
  description: string
  modelSrc?: string | null
  modelSizeLabel?: string | null
  modelRotationPerSecond?: string | null
}
export type AstronomyBody = {
  id: string
  symbol: string
  name: string
  kind: 'star' | 'planet' | 'moon'
  diameterKm: number
  color: string
  description: string
  modelSrc?: string | null
  modelSizeLabel?: string | null
  modelRotationPerSecond?: string | null
}
export const ASTRONOMY_SUN: AstronomyBody = {
  id: 'sun',
  symbol: '☉',
  name: 'Солнце',
  kind: 'star',
  diameterKm: 1392700,
  color: '#fbbf24',
  description:
    'Солнце содержит почти всю массу Солнечной системы и определяет орбитальную динамику планет, астероидов и комет.',
  modelSrc: '/astronomy/models/sun.glb',
  modelSizeLabel: '2.1 МБ'
}
export const ASTRONOMY_PLANETS: AstronomyPlanet[] = [
  {
    id: 'mercury',
    symbol: '☿',
    name: 'Меркурий',
    kind: 'planet',
    diameterKm: 4879,
    orbitalRadiusAu: 0.387,
    orbitalPeriodDays: 87.97,
    color: '#a1a1aa',
    description: 'Ближайшая к Солнцу планета с коротким годом и выраженными перепадами температуры.',
    modelSrc: '/astronomy/models/mercury_planet.glb',
    modelSizeLabel: '3.6 МБ'
  },
  {
    id: 'venus',
    symbol: '♀',
    name: 'Венера',
    kind: 'planet',
    diameterKm: 12104,
    orbitalRadiusAu: 0.723,
    orbitalPeriodDays: 224.7,
    color: '#d4a373',
    description: 'Плотная атмосфера Венеры создаёт мощный парниковый эффект и экстремально горячую поверхность.',
    modelSrc: '/astronomy/models/venus.glb',
    modelSizeLabel: '1.4 МБ'
  },
  {
    id: 'earth',
    symbol: '♁',
    name: 'Земля',
    kind: 'planet',
    diameterKm: 12742,
    orbitalRadiusAu: 1,
    orbitalPeriodDays: 365.256,
    color: '#60a5fa',
    description: 'Единственная известная планета с устойчивой жидкой водой на поверхности и сложной биосферой.',
    modelSrc: '/astronomy/models/earth_new.glb',
    modelSizeLabel: '20.3 МБ',
    modelRotationPerSecond: '10deg'
  },
  {
    id: 'mars',
    symbol: '♂',
    name: 'Марс',
    kind: 'planet',
    diameterKm: 6779,
    orbitalRadiusAu: 1.524,
    orbitalPeriodDays: 686.98,
    color: '#f97316',
    description: 'Холодная каменистая планета с тонкой атмосферой, полярными шапками и древними следами воды.',
    modelSrc: '/astronomy/models/mars.glb',
    modelSizeLabel: '6.0 МБ'
  },
  {
    id: 'jupiter',
    symbol: '♃',
    name: 'Юпитер',
    kind: 'planet',
    diameterKm: 139820,
    orbitalRadiusAu: 5.203,
    orbitalPeriodDays: 4332.59,
    color: '#f59e0b',
    description: 'Крупнейшая планета системы, газовый гигант с мощной магнитосферой и многочисленными спутниками.',
    modelSrc: '/astronomy/models/jupiter.glb',
    modelSizeLabel: '2.4 МБ'
  },
  {
    id: 'saturn',
    symbol: '♄',
    name: 'Сатурн',
    kind: 'planet',
    diameterKm: 116460,
    orbitalRadiusAu: 9.537,
    orbitalPeriodDays: 10759.22,
    color: '#eab308',
    description: 'Газовый гигант с выраженной кольцевой системой и богатым семейством ледяных спутников.',
    modelSrc: '/astronomy/models/saturn.glb',
    modelSizeLabel: '3.3 МБ'
  },
  {
    id: 'uranus',
    symbol: '♅',
    name: 'Уран',
    kind: 'planet',
    diameterKm: 50724,
    orbitalRadiusAu: 19.191,
    orbitalPeriodDays: 30688.5,
    color: '#67e8f9',
    description: 'Ледяной гигант, ось которого заметно наклонена, поэтому сезонные изменения здесь необычайно резкие.',
    modelSrc: '/astronomy/models/uranus.glb',
    modelSizeLabel: '229 КБ'
  },
  {
    id: 'neptune',
    symbol: '♆',
    name: 'Нептун',
    kind: 'planet',
    diameterKm: 49244,
    orbitalRadiusAu: 30.07,
    orbitalPeriodDays: 60182,
    color: '#38bdf8',
    description: 'Самая удалённая планета с быстрыми ветрами, глубокими облачными слоями и крупным спутником Тритоном.',
    modelSrc: '/astronomy/models/neptune.glb',
    modelSizeLabel: '1.7 МБ'
  }
]
export const ASTRONOMY_MOONS: AstronomyMoon[] = [
  {
    id: 'moon',
    planetId: 'earth',
    symbol: 'Lu',
    name: 'Луна',
    kind: 'moon',
    diameterKm: 3474,
    orbitalRadiusKm: 384400,
    orbitalPeriodDays: 27.32,
    color: '#d4d4d8',
    description: 'Спутник Земли, стабилизирующий наклон оси и заметно влияющий на приливы.',
    modelSrc: '/astronomy/models/moon.glb',
    modelSizeLabel: '5.7 МБ'
  },
  {
    id: 'phobos',
    planetId: 'mars',
    symbol: 'Ph',
    name: 'Фобос',
    kind: 'moon',
    diameterKm: 23,
    orbitalRadiusKm: 9376,
    orbitalPeriodDays: 0.319,
    color: '#a8a29e',
    description: 'Крупнейший спутник Марса, движущийся очень близко к планете и постепенно снижающий орбиту.'
  },
  {
    id: 'deimos',
    planetId: 'mars',
    symbol: 'De',
    name: 'Деймос',
    kind: 'moon',
    diameterKm: 12,
    orbitalRadiusKm: 23463,
    orbitalPeriodDays: 1.263,
    color: '#c4b5fd',
    description: 'Небольшой внешний спутник Марса с медленным орбитальным движением.'
  },
  {
    id: 'io',
    planetId: 'jupiter',
    symbol: 'Io',
    name: 'Ио',
    kind: 'moon',
    diameterKm: 3643,
    orbitalRadiusKm: 421700,
    orbitalPeriodDays: 1.769,
    color: '#facc15',
    description: 'Один из самых вулканически активных объектов Солнечной системы.'
  },
  {
    id: 'europa',
    planetId: 'jupiter',
    symbol: 'Eu',
    name: 'Европа',
    kind: 'moon',
    diameterKm: 3122,
    orbitalRadiusKm: 671100,
    orbitalPeriodDays: 3.551,
    color: '#bfdbfe',
    description: 'Ледяной спутник Юпитера с вероятным глобальным океаном под поверхностью.'
  },
  {
    id: 'ganymede',
    planetId: 'jupiter',
    symbol: 'Ga',
    name: 'Ганимед',
    kind: 'moon',
    diameterKm: 5268,
    orbitalRadiusKm: 1070400,
    orbitalPeriodDays: 7.155,
    color: '#93c5fd',
    description: 'Крупнейший спутник Солнечной системы, превосходящий по размеру Меркурий.'
  },
  {
    id: 'callisto',
    planetId: 'jupiter',
    symbol: 'Ca',
    name: 'Каллисто',
    kind: 'moon',
    diameterKm: 4821,
    orbitalRadiusKm: 1882700,
    orbitalPeriodDays: 16.689,
    color: '#94a3b8',
    description: 'Далёкий крупный спутник с древней кратерированной поверхностью.'
  },
  {
    id: 'titan',
    planetId: 'saturn',
    symbol: 'Ti',
    name: 'Титан',
    kind: 'moon',
    diameterKm: 5150,
    orbitalRadiusKm: 1221870,
    orbitalPeriodDays: 15.945,
    color: '#fcd34d',
    description: 'Крупный спутник Сатурна с плотной атмосферой и углеводородными морями.'
  },
  {
    id: 'enceladus',
    planetId: 'saturn',
    symbol: 'En',
    name: 'Энцелад',
    kind: 'moon',
    diameterKm: 505,
    orbitalRadiusKm: 237948,
    orbitalPeriodDays: 1.37,
    color: '#e0f2fe',
    description: 'Небольшой ледяной спутник с активными гейзерами и подповерхностным океаном.'
  },
  {
    id: 'titania',
    planetId: 'uranus',
    symbol: 'Ta',
    name: 'Титания',
    kind: 'moon',
    diameterKm: 1578,
    orbitalRadiusKm: 436300,
    orbitalPeriodDays: 8.706,
    color: '#cbd5e1',
    description: 'Крупнейший спутник Урана с системой разломов и древних ударных бассейнов.'
  },
  {
    id: 'oberon',
    planetId: 'uranus',
    symbol: 'Ob',
    name: 'Оберон',
    kind: 'moon',
    diameterKm: 1523,
    orbitalRadiusKm: 583500,
    orbitalPeriodDays: 13.463,
    color: '#94a3b8',
    description: 'Внешний крупный спутник Урана с тёмной древней поверхностью.'
  },
  {
    id: 'triton',
    planetId: 'neptune',
    symbol: 'Tr',
    name: 'Тритон',
    kind: 'moon',
    diameterKm: 2707,
    orbitalRadiusKm: 354759,
    orbitalPeriodDays: 5.877,
    color: '#bae6fd',
    description: 'Крупнейший спутник Нептуна с ретроградной орбитой и следами криовулканизма.'
  }
]
export const ASTRONOMY_GALACTIC_CONTEXT = {
  galacticCenterName: 'Стрелец A*',
  galacticCenterMassSolar: 4.1e6,
  solarSystemRadiusLy: 26600,
  solarSystemOrbitalPeriodMillionYears: 230,
  solarSystemOrbitalSpeedKmS: 220
}
export const ASTRONOMY_SIZE_BODIES: AstronomyBody[] = [
  ASTRONOMY_SUN,
  ...ASTRONOMY_PLANETS.map(
    ({ id, symbol, name, kind, diameterKm, color, description, modelSrc, modelSizeLabel, modelRotationPerSecond }) => ({
    id,
    symbol,
    name,
    kind,
    diameterKm,
    color,
    description,
    modelSrc,
      modelSizeLabel,
      modelRotationPerSecond
    })
  ),
  ...ASTRONOMY_MOONS.map(
    ({ id, symbol, name, kind, diameterKm, color, description, modelSrc, modelSizeLabel, modelRotationPerSecond }) => ({
    id,
    symbol,
    name,
    kind,
    diameterKm,
    color,
    description,
      modelSrc,
      modelSizeLabel,
      modelRotationPerSecond
    })
  )
].sort((left, right) => right.diameterKm - left.diameterKm)
