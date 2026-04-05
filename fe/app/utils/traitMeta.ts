import type { DataType, KeyMeta, TraitColorMode, TraitUnitCategory, TraitUnitOption } from '../../shared/types/traits'
export const DATA_TYPES: DataType[] = [
  'string',
  'enum',
  'number',
  'datetime',
  'datetime-range',
  'interval',
  'schedule',
  'geo-point',
  'color',
  'boolean',
  'validity'
]
export const DATE_DATA_TYPES: DataType[] = ['datetime', 'datetime-range', 'interval', 'schedule']
export const DATA_TYPE_SELECT_LABELS: Record<DataType, string> = {
  string: 'Строка',
  enum: 'Список',
  number: 'Число',
  datetime: 'Момент времени',
  'datetime-range': 'Диапазон дат/времени',
  interval: 'Интервал времени',
  schedule: 'Расписание',
  'geo-point': 'Геоточка',
  color: 'Цвет',
  boolean: 'Булево',
  validity: 'Срок действия'
}
export const DATA_TYPE_LABELS: Record<DataType, string> = {
  string: 'строка',
  enum: 'список',
  number: 'число',
  datetime: 'момент времени',
  'datetime-range': 'диапазон',
  interval: 'интервал времени',
  schedule: 'расписание',
  'geo-point': 'гео',
  color: 'цвет',
  boolean: 'булево',
  validity: 'валидность'
}
export const COLOR_MODE_OPTIONS: Array<{ value: TraitColorMode; label: string }> = [
  { value: 'hex', label: 'HEX-код' },
  { value: 'lab', label: 'LAB' },
  { value: 'spectrum', label: 'спектр' }
]
export const COLOR_MODE_LABELS: Record<TraitColorMode, string> = {
  hex: 'HEX-код',
  lab: 'LAB',
  spectrum: 'спектр'
}
export const UNIT_CATEGORIES: Array<{ id: TraitUnitCategory; label: string }> = [
  { id: 'unitless', label: 'Безразмерная' },
  { id: 'length', label: 'Длина' },
  { id: 'mass', label: 'Масса' },
  { id: 'temperature', label: 'Температура' },
  { id: 'time', label: 'Время' },
  { id: 'area', label: 'Площадь' },
  { id: 'volume', label: 'Объем' },
  { id: 'speed', label: 'Скорость' },
  { id: 'energy', label: 'Энергия' }
]
export const UNIT_CATEGORY_LABELS: Record<TraitUnitCategory, string> = {
  unitless: 'безразмерная',
  length: 'длина',
  mass: 'масса',
  temperature: 'температура',
  time: 'время',
  area: 'площадь',
  volume: 'объем',
  speed: 'скорость',
  energy: 'энергия'
}
export const UNIT_OPTIONS: Record<TraitUnitCategory, TraitUnitOption[]> = {
  unitless: [{ value: '', label: '—' }],
  length: [
    { value: 'm', label: 'метр' },
    { value: 'cm', label: 'сантиметр' },
    { value: 'mm', label: 'миллиметр' },
    { value: 'km', label: 'километр' }
  ],
  mass: [
    { value: 'kg', label: 'килограмм' },
    { value: 'g', label: 'грамм' },
    { value: 'mg', label: 'миллиграмм' },
    { value: 't', label: 'тонна' }
  ],
  temperature: [
    { value: 'C', label: 'цельсий' },
    { value: 'K', label: 'кельвин' }
  ],
  time: [
    { value: 's', label: 'секунда' },
    { value: 'ms', label: 'миллисекунда' },
    { value: 'min', label: 'минута' },
    { value: 'h', label: 'час' }
  ],
  area: [
    { value: 'm2', label: 'квадратный метр' },
    { value: 'cm2', label: 'квадратный сантиметр' },
    { value: 'km2', label: 'квадратный километр' }
  ],
  volume: [
    { value: 'm3', label: 'кубический метр' },
    { value: 'L', label: 'литр' },
    { value: 'mL', label: 'миллилитр' }
  ],
  speed: [
    { value: 'm/s', label: 'метр в секунду' },
    { value: 'km/h', label: 'километр в час' }
  ],
  energy: [
    { value: 'J', label: 'джоуль' },
    { value: 'kJ', label: 'килоджоуль' },
    { value: 'kWh', label: 'киловатт-час' }
  ]
}
export const UNIT_LABELS: Record<string, string> = {
  m: 'метр',
  cm: 'сантиметр',
  mm: 'миллиметр',
  km: 'километр',
  kg: 'килограмм',
  g: 'грамм',
  mg: 'миллиграмм',
  t: 'тонна',
  C: 'градус Цельсия',
  K: 'кельвин',
  s: 'секунда',
  ms: 'миллисекунда',
  min: 'минута',
  h: 'час',
  m2: 'квадратный метр',
  cm2: 'квадратный сантиметр',
  km2: 'квадратный километр',
  m3: 'кубический метр',
  L: 'литр',
  mL: 'миллилитр',
  'm/s': 'метр в секунду',
  'km/h': 'километр в час',
  J: 'джоуль',
  kJ: 'килоджоуль',
  kWh: 'киловатт-час'
}
/**
 * Checks whether the provided value is one of the supported trait data types.
 */
const isDataType = (value: unknown): value is DataType => {
  return typeof value === 'string' && DATA_TYPES.includes(value as DataType)
}
/**
 * Normalizes a string list by trimming entries and removing duplicates.
 */
const normalizeStringList = (items: unknown[]): string[] => {
  const unique = new Set<string>()
  for (const item of items) {
    const text = String(item ?? '').trim()
    if (!text) continue
    unique.add(text)
  }
  return Array.from(unique)
}
/**
 * Parses a comma-separated enum options string into normalized values.
 */
export const parseEnumOptionsRaw = (raw: string): string[] => {
  return normalizeStringList(String(raw || '').split(','))
}
/**
 * Creates the default meta configuration for a trait data type.
 */
export const defaultKeyMeta = (dataType: DataType = 'string'): KeyMeta => ({
  dataType,
  optionType: dataType === 'enum' ? 'string' : undefined,
  mode: dataType === 'color' ? 'hex' : undefined,
  options: dataType === 'enum' ? [] : undefined,
  unitCategory: dataType === 'number' ? 'unitless' : undefined,
  unit: dataType === 'number' ? '' : undefined
})
/**
 * Resolves the effective color mode from key meta.
 */
export const resolveColorMode = (meta?: Partial<KeyMeta> | Record<string, unknown> | null): TraitColorMode | '' => {
  const raw = String((meta as any)?.mode || (meta as any)?.format || '')
    .trim()
    .toLowerCase()
  if (raw === 'hex' || raw === 'lab' || raw === 'spectrum') return raw
  return ''
}
/**
 * Resolves the effective unit category with a safe fallback.
 */
export const resolveUnitCategory = (category?: string | null): TraitUnitCategory => {
  const raw = String(category || '').trim()
  const found = UNIT_CATEGORIES.find(item => item.id === raw)
  return found?.id || 'unitless'
}
/**
 * Returns the supported unit options for one category.
 */
export const getUnitOptionsByCategory = (category?: string | null): TraitUnitOption[] => {
  const safeCategory = resolveUnitCategory(category)
  return UNIT_OPTIONS[safeCategory] || UNIT_OPTIONS.unitless
}
/**
 * Normalizes backend or form meta into the canonical front-end shape.
 */
export const normalizeKeyMeta = (meta?: KeyMeta | null, enumOptionsRaw?: string): KeyMeta => {
  if (!meta) return defaultKeyMeta('string')
  const dataType = isDataType(meta.dataType) ? meta.dataType : 'string'
  const base = { ...defaultKeyMeta(dataType), ...meta, dataType }
  if (dataType === 'enum') {
    const options =
      enumOptionsRaw !== undefined
        ? parseEnumOptionsRaw(enumOptionsRaw)
        : normalizeStringList(Array.isArray(base.options) ? base.options : [])
    return {
      dataType: 'enum',
      optionType: base.optionType || 'string',
      options: options.length ? options : undefined
    }
  }
  if (DATE_DATA_TYPES.includes(dataType)) {
    return { dataType }
  }
  if (dataType === 'number') {
    const unitCategory = resolveUnitCategory(base.unitCategory)
    const units = getUnitOptionsByCategory(unitCategory)
    const unit = units.find(item => item.value === base.unit)?.value || units[0]?.value || ''
    return {
      ...base,
      unitCategory,
      unit
    }
  }
  if (dataType === 'color') {
    const mode = resolveColorMode(base) || 'hex'
    return {
      dataType: 'color',
      mode
    }
  }
  return { dataType }
}
/**
 * Produces a stable meta shape for structural equality checks.
 */
export const normalizeMetaForEquality = (meta?: KeyMeta | null): KeyMeta => {
  const normalized = normalizeKeyMeta(meta)
  const base = defaultKeyMeta(normalized.dataType || 'string')
  return {
    ...base,
    ...normalized,
    unitCategory: normalized.dataType === 'number' ? normalized.unitCategory || 'unitless' : undefined,
    unit: normalized.dataType === 'number' ? normalized.unit || '' : undefined
  }
}
/**
 * Compares two meta objects after canonical normalization.
 */
export const keyMetaEquals = (a?: KeyMeta | null, b?: KeyMeta | null): boolean => {
  return JSON.stringify(normalizeMetaForEquality(a)) === JSON.stringify(normalizeMetaForEquality(b))
}
const re = {
  stringSafe: /^[^\x00-\x1F\x7F]+$/,
  number: /^-?\d+(\.\d+)?$/,
  datetimeLocal: /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/,
  timeLocal: /^\d{2}:\d{2}$/,
  geoPoint: /^-?\d{1,3}\.\d+,-?\d{1,3}\.\d+$/,
  cielab: /^-?\d{1,3},-?\d{1,3},-?\d{1,3}$/,
  boolean01: /^[01]$/,
  booleanTF: /^(true|false)$/i,
  rangeDash: /^([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2})\s—\s([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2})$/,
  rangeSlash: /^([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2})\/([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2})$/,
  intervalWeekly: /^weekly:([1-7])-([1-7])\/([0-9]{2}:[0-9]{2})-([0-9]{2}:[0-9]{2})$/,
  intervalDuration:
    /^duration:([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2})\/([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2})\|(-?\d+(?:\.\d+)?)\|(seconds|minutes|hours|days|years)$/,
  validityPermanent: /^permanent(?::([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}))?$/,
  validityTemporary: /^temporary:([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2})$/,
  hex: /^#?[0-9a-fA-F]{6}$/
}
/**
 * Parses a canonical local datetime string into a UTC timestamp.
 */
const parseCanonicalDateTimeLocal = (value: string): number | null => {
  const m = String(value || '')
    .trim()
    .match(/^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2})$/)
  if (!m) return null
  const year = Number(m[1])
  const month = Number(m[2])
  const day = Number(m[3])
  const hour = Number(m[4])
  const minute = Number(m[5])
  if (month < 1 || month > 12) return null
  if (day < 1 || day > 31) return null
  if (hour < 0 || hour > 23 || minute < 0 || minute > 59) return null
  const dt = new Date(Date.UTC(year, month - 1, day, hour, minute, 0, 0))
  if (dt.getUTCFullYear() !== year || dt.getUTCMonth() !== month - 1 || dt.getUTCDate() !== day) return null
  return dt.getTime()
}
/**
 * Returns true when the datetime string is canonical and calendar-valid.
 */
const isCanonicalDateTimeLocal = (value: string): boolean => {
  return parseCanonicalDateTimeLocal(value) !== null
}
/**
 * Validates a left/right canonical datetime range.
 */
const isCanonicalRangeValid = (left: string, right: string): boolean => {
  const leftTs = parseCanonicalDateTimeLocal(left)
  const rightTs = parseCanonicalDateTimeLocal(right)
  if (leftTs === null || rightTs === null) return false
  return leftTs <= rightTs
}
/**
 * Parses a canonical local time string into minutes since midnight.
 */
const parseCanonicalTimeLocal = (value: string): number | null => {
  const m = String(value || '')
    .trim()
    .match(/^(\d{2}):(\d{2})$/)
  if (!m) return null
  const hour = Number(m[1])
  const minute = Number(m[2])
  if (hour < 0 || hour > 23 || minute < 0 || minute > 59) return null
  return hour * 60 + minute
}
/**
 * Validates a left/right canonical local time range.
 */
const isCanonicalTimeRangeValid = (left: string, right: string): boolean => {
  const leftMins = parseCanonicalTimeLocal(left)
  const rightMins = parseCanonicalTimeLocal(right)
  if (leftMins === null || rightMins === null) return false
  return leftMins <= rightMins
}
/**
 * Validates a serialized LAB triplet string.
 */
const isLabTupleValid = (value: string): boolean => {
  if (!re.cielab.test(value)) return false
  const parts = value.split(',')
  if (parts.length !== 3) return false
  const [lStr, aStr, bStr] = parts
  const l = Number(lStr)
  const a = Number(aStr)
  const b = Number(bStr)
  return l >= 0 && l <= 100 && a >= -128 && a <= 127 && b >= -128 && b <= 127
}
/**
 * Validates a serialized trait value against its meta definition.
 */
export const validateValue = (value: unknown, meta: KeyMeta): boolean => {
  const v = typeof value === 'string' ? value.trim() : value
  switch (meta.dataType) {
    case 'string':
      return typeof v === 'string' && v.length > 0 && re.stringSafe.test(v)
    case 'enum':
      return typeof v === 'string' && v.length > 0
    case 'number':
      return typeof v === 'number' ? !Number.isNaN(v) : typeof v === 'string' && re.number.test(v)
    case 'datetime':
      return typeof v === 'string' && isCanonicalDateTimeLocal(v)
    case 'datetime-range': {
      if (typeof v !== 'string') return false
      const m = v.match(re.rangeDash)
      if (!m) return false
      return isCanonicalRangeValid(String(m[1] || ''), String(m[2] || ''))
    }
    case 'interval': {
      if (typeof v !== 'string') return false
      const duration = v.match(re.intervalDuration)
      if (duration) {
        const start = String(duration[1] || '')
        const end = String(duration[2] || '')
        const amount = Number(duration[3])
        const unit = String(duration[4] || '')
        if (!isCanonicalRangeValid(start, end)) return false
        if (!Number.isFinite(amount) || amount < 0) return false
        return ['seconds', 'minutes', 'hours', 'days', 'years'].includes(unit)
      }
      // Backward compatibility for old interval formats.
      const weekly = v.match(re.intervalWeekly)
      if (weekly) {
        const fromDay = Number(weekly[1])
        const toDay = Number(weekly[2])
        const fromTime = String(weekly[3] || '')
        const toTime = String(weekly[4] || '')
        if (!Number.isInteger(fromDay) || !Number.isInteger(toDay)) return false
        if (fromDay < 1 || fromDay > 7 || toDay < 1 || toDay > 7) return false
        if (parseCanonicalTimeLocal(fromTime) === null || parseCanonicalTimeLocal(toTime) === null) return false
        if (fromDay === toDay) return isCanonicalTimeRangeValid(fromTime, toTime)
        return true
      }
      // Backward compatibility: previously interval was stored as datetime/datetime with "/".
      const legacy = v.match(re.rangeSlash)
      if (!legacy) return false
      return isCanonicalRangeValid(String(legacy[1] || ''), String(legacy[2] || ''))
    }
    case 'schedule': {
      if (typeof v !== 'string') return false
      const weekly = v.match(re.intervalWeekly)
      if (!weekly) return false
      const fromDay = Number(weekly[1])
      const toDay = Number(weekly[2])
      const fromTime = String(weekly[3] || '')
      const toTime = String(weekly[4] || '')
      if (!Number.isInteger(fromDay) || !Number.isInteger(toDay)) return false
      if (fromDay < 1 || fromDay > 7 || toDay < 1 || toDay > 7) return false
      if (parseCanonicalTimeLocal(fromTime) === null || parseCanonicalTimeLocal(toTime) === null) return false
      if (fromDay === toDay) return isCanonicalTimeRangeValid(fromTime, toTime)
      return true
    }
    case 'geo-point': {
      if (typeof v !== 'string' || !re.geoPoint.test(v)) return false
      const parts = v.split(',')
      if (parts.length !== 2) return false
      const [latS, lonS] = parts
      const lat = Number(latS)
      const lon = Number(lonS)
      return lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
    }
    case 'color': {
      if (typeof v !== 'string' || v.length === 0) return false
      const colorMode = resolveColorMode(meta)
      if (colorMode === 'hex') return re.hex.test(v)
      if (colorMode === 'spectrum') return /^[a-z]+$/i.test(v)
      if (colorMode === 'lab') return isLabTupleValid(v)
      if (re.hex.test(v)) return true
      if (isLabTupleValid(v)) return true
      try {
        const obj = JSON.parse(v)
        if (!obj || typeof obj !== 'object') return false
        const mode = (obj as any).mode
        if (mode === 'hex') return typeof (obj as any).hex === 'string' && re.hex.test((obj as any).hex)
        if (mode === 'spectrum') return typeof (obj as any).spectrum === 'string' && (obj as any).spectrum.length > 0
        if (mode === 'lab' && (obj as any).lab && typeof (obj as any).lab === 'object') {
          const l = Number((obj as any).lab.L)
          const a = Number((obj as any).lab.a)
          const b = Number((obj as any).lab.b)
          if (![l, a, b].every(Number.isFinite)) return false
          return l >= 0 && l <= 100 && a >= -128 && a <= 127 && b >= -128 && b <= 127
        }
      } catch {
        // ignore parse errors
      }
      return false
    }
    case 'boolean':
      return typeof v === 'boolean' || (typeof v === 'string' && (re.boolean01.test(v) || re.booleanTF.test(v)))
    case 'validity': {
      if (typeof v !== 'string' || !v) return false
      const permanent = v.match(re.validityPermanent)
      if (permanent) {
        const since = String(permanent[1] || '')
        return !since || isCanonicalDateTimeLocal(since)
      }
      const temporary = v.match(re.validityTemporary)
      if (!temporary) return false
      return isCanonicalDateTimeLocal(String(temporary[1] || ''))
    }
    default:
      return false
  }
}
