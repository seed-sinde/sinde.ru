import type {
  DataType,
  KeyMeta,
  TraitMetaUi,
  TraitBooleanDisplay,
  TraitColorMode,
  TraitGeoType,
  TraitListValueType,
  TraitRangeValueType,
  TraitUnitCategory,
  TraitUnitOption
} from "../../shared/types/traits"
export const DATA_TYPES: DataType[] = [
  "string",
  "number",
  "list",
  "datetime",
  "range",
  "interval",
  "schedule",
  "geo",
  "color",
  "boolean",
  "validity",
  "surface"
]
export const LEGACY_DATA_TYPE_ALIASES: Partial<Record<DataType, DataType>> = {
  enum: "list",
  "datetime-range": "range",
  "geo-point": "geo"
}
export const DATE_DATA_TYPES: DataType[] = ["datetime", "range", "interval", "schedule"]
export const DATA_TYPE_SELECT_LABELS: Record<DataType, string> = {
  string: "Строка",
  number: "Число",
  list: "Список",
  datetime: "Момент времени",
  range: "Диапазон",
  interval: "Интервал времени",
  schedule: "Расписание",
  geo: "Геоданные",
  color: "Цвет",
  boolean: "Булево",
  validity: "Срок действия",
  surface: "Поверхность",
  enum: "Список",
  "datetime-range": "Диапазон дат/времени",
  "geo-point": "Геоточка"
}
export const DATA_TYPE_LABELS: Record<DataType, string> = {
  string: "строка",
  number: "число",
  list: "список",
  datetime: "момент времени",
  range: "диапазон",
  interval: "интервал времени",
  schedule: "расписание",
  geo: "гео",
  color: "цвет",
  boolean: "булево",
  validity: "валидность",
  surface: "поверхность",
  enum: "список",
  "datetime-range": "диапазон",
  "geo-point": "гео"
}
export const BOOLEAN_DISPLAY_OPTIONS: Array<{value: TraitBooleanDisplay; label: string; trueLabel: string; falseLabel: string}> = [
  {value: "truth", label: "Истина/Ложь", trueLabel: "Истина", falseLabel: "Ложь"},
  {value: "answer", label: "Да/Нет", trueLabel: "Да", falseLabel: "Нет"},
  {value: "state", label: "Вкл/Выкл", trueLabel: "Вкл", falseLabel: "Выкл"},
  {value: "presence", label: "Есть/Нет", trueLabel: "Есть", falseLabel: "Нет"}
]
export const LIST_VALUE_TYPE_OPTIONS: Array<{value: TraitListValueType; label: string}> = [
  {value: "string", label: "Строка"},
  {value: "number", label: "Число"},
  {value: "datetime", label: "Момент времени"},
  {value: "date", label: "Дата"},
  {value: "time", label: "Время"}
]
export const RANGE_VALUE_TYPE_OPTIONS: Array<{value: TraitRangeValueType; label: string}> = [
  {value: "number", label: "Число"},
  {value: "datetime", label: "Дата и время"},
  {value: "date", label: "Дата"},
  {value: "time", label: "Время"}
]
export const GEO_TYPE_OPTIONS: Array<{value: TraitGeoType; label: string}> = [
  {value: "point", label: "Точка"},
  {value: "polygon", label: "Полигон"},
  {value: "zone", label: "Зона"}
]
export const COLOR_MODE_OPTIONS: Array<{value: TraitColorMode; label: string}> = [
  {value: "hex", label: "HEX-код"},
  {value: "lab", label: "LAB"},
  {value: "spectrum", label: "спектр"}
]
export const COLOR_MODE_LABELS: Record<TraitColorMode, string> = {
  hex: "HEX-код",
  lab: "LAB",
  spectrum: "спектр"
}
export const UNIT_CATEGORIES: Array<{id: TraitUnitCategory; label: string}> = [
  {id: "unitless", label: "Безразмерная"},
  {id: "length", label: "Длина"},
  {id: "mass", label: "Масса"},
  {id: "temperature", label: "Температура"},
  {id: "time", label: "Время"},
  {id: "area", label: "Площадь"},
  {id: "volume", label: "Объем"},
  {id: "speed", label: "Скорость"},
  {id: "energy", label: "Энергия"},
  {id: "angle", label: "Угол"},
  {id: "pressure", label: "Давление"},
  {id: "density", label: "Плотность"},
  {id: "force", label: "Сила"},
  {id: "power", label: "Мощность"},
  {id: "frequency", label: "Частота"},
  {id: "electric", label: "Электрические величины"},
  {id: "error", label: "Погрешность ±"}
]
export const UNIT_CATEGORY_LABELS: Record<TraitUnitCategory, string> = {
  unitless: "безразмерная",
  length: "длина",
  mass: "масса",
  temperature: "температура",
  time: "время",
  area: "площадь",
  volume: "объем",
  speed: "скорость",
  energy: "энергия",
  angle: "угол",
  pressure: "давление",
  density: "плотность",
  force: "сила",
  power: "мощность",
  frequency: "частота",
  electric: "электрические величины",
  error: "погрешность ±"
}
export const UNIT_OPTIONS: Record<TraitUnitCategory, TraitUnitOption[]> = {
  unitless: [{value: "", label: "—"}],
  length: [
    {value: "m", label: "метр"},
    {value: "cm", label: "сантиметр"},
    {value: "mm", label: "миллиметр"},
    {value: "km", label: "километр"}
  ],
  mass: [
    {value: "kg", label: "килограмм"},
    {value: "g", label: "грамм"},
    {value: "mg", label: "миллиграмм"},
    {value: "t", label: "тонна"}
  ],
  temperature: [
    {value: "C", label: "цельсий"},
    {value: "K", label: "кельвин"}
  ],
  time: [
    {value: "s", label: "секунда"},
    {value: "ms", label: "миллисекунда"},
    {value: "min", label: "минута"},
    {value: "h", label: "час"}
  ],
  area: [
    {value: "m2", label: "квадратный метр"},
    {value: "cm2", label: "квадратный сантиметр"},
    {value: "km2", label: "квадратный километр"}
  ],
  volume: [
    {value: "m3", label: "кубический метр"},
    {value: "L", label: "литр"},
    {value: "mL", label: "миллилитр"}
  ],
  speed: [
    {value: "m/s", label: "метр в секунду"},
    {value: "km/h", label: "километр в час"}
  ],
  energy: [
    {value: "J", label: "джоуль"},
    {value: "kJ", label: "килоджоуль"},
    {value: "kWh", label: "киловатт-час"}
  ],
  angle: [
    {value: "deg", label: "градус"},
    {value: "rad", label: "радиан"}
  ],
  pressure: [
    {value: "Pa", label: "паскаль"},
    {value: "kPa", label: "килопаскаль"},
    {value: "bar", label: "бар"}
  ],
  density: [
    {value: "kg/m3", label: "кг/м³"},
    {value: "g/cm3", label: "г/см³"}
  ],
  force: [
    {value: "N", label: "ньютон"},
    {value: "kN", label: "килоньютон"}
  ],
  power: [
    {value: "W", label: "ватт"},
    {value: "kW", label: "киловатт"}
  ],
  frequency: [
    {value: "Hz", label: "герц"},
    {value: "kHz", label: "килогерц"},
    {value: "MHz", label: "мегагерц"}
  ],
  electric: [
    {value: "V", label: "вольт"},
    {value: "A", label: "ампер"},
    {value: "Ohm", label: "ом"}
  ],
  error: [
    {value: "%", label: "процент"},
    {value: "abs", label: "абсолютная"}
  ]
}
export const UNIT_LABELS: Record<string, string> = {
  m: "метр",
  cm: "сантиметр",
  mm: "миллиметр",
  km: "километр",
  kg: "килограмм",
  g: "грамм",
  mg: "миллиграмм",
  t: "тонна",
  C: "градус Цельсия",
  K: "кельвин",
  s: "секунда",
  ms: "миллисекунда",
  min: "минута",
  h: "час",
  m2: "квадратный метр",
  cm2: "квадратный сантиметр",
  km2: "квадратный километр",
  m3: "кубический метр",
  L: "литр",
  mL: "миллилитр",
  "m/s": "метр в секунду",
  "km/h": "километр в час",
  J: "джоуль",
  kJ: "килоджоуль",
  kWh: "киловатт-час",
  deg: "градус",
  rad: "радиан",
  Pa: "паскаль",
  kPa: "килопаскаль",
  bar: "бар",
  "kg/m3": "килограмм на кубометр",
  "g/cm3": "грамм на кубический сантиметр",
  N: "ньютон",
  kN: "килоньютон",
  W: "ватт",
  kW: "киловатт",
  Hz: "герц",
  kHz: "килогерц",
  MHz: "мегагерц",
  V: "вольт",
  A: "ампер",
  Ohm: "ом",
  "%": "процент",
  abs: "абсолютная"
}
/**
 * Checks whether the provided value is one of the supported trait data types.
 */
const isDataType = (value: unknown): value is DataType => {
  return typeof value === "string" && (DATA_TYPES.includes(value as DataType) || value in LEGACY_DATA_TYPE_ALIASES)
}
export const resolveDataType = (value?: string | null): DataType => {
  const raw = String(value || "").trim() as DataType
  if (LEGACY_DATA_TYPE_ALIASES[raw]) return LEGACY_DATA_TYPE_ALIASES[raw] as DataType
  return DATA_TYPES.includes(raw) ? raw : "string"
}
export const resolveBooleanDisplay = (value?: string | null): TraitBooleanDisplay => {
  const raw = String(value || "").trim() as TraitBooleanDisplay
  return BOOLEAN_DISPLAY_OPTIONS.some(option => option.value === raw) ? raw : "answer"
}
export const resolveListValueType = (value?: string | null): TraitListValueType => {
  const raw = String(value || "").trim() as TraitListValueType
  return LIST_VALUE_TYPE_OPTIONS.some(option => option.value === raw) ? raw : "string"
}
export const resolveRangeValueType = (value?: string | null): TraitRangeValueType => {
  const raw = String(value || "").trim() as TraitRangeValueType
  return RANGE_VALUE_TYPE_OPTIONS.some(option => option.value === raw) ? raw : "datetime"
}
export const resolveGeoType = (value?: string | null): TraitGeoType => {
  const raw = String(value || "").trim() as TraitGeoType
  return GEO_TYPE_OPTIONS.some(option => option.value === raw) ? raw : "point"
}
export const getDefaultBooleanLabels = (display?: string | null) => {
  const resolved = resolveBooleanDisplay(display)
  return BOOLEAN_DISPLAY_OPTIONS.find(option => option.value === resolved) || BOOLEAN_DISPLAY_OPTIONS[1]!
}
export const getDefaultUiInput = (dataType: DataType): TraitMetaUi["input"] => {
  switch (dataType) {
    case "number":
      return "number"
    case "list":
      return "array"
    case "datetime":
      return "datetime"
    case "range":
      return "date-range"
    case "interval":
      return "interval"
    case "schedule":
      return "schedule"
    case "geo":
      return "geo"
    case "color":
      return "color"
    case "boolean":
      return "boolean"
    case "validity":
      return "validity"
    default:
      return "text"
  }
}
/**
 * Normalizes a string list by trimming entries and removing duplicates.
 */
const normalizeStringList = (items: unknown[]): string[] => {
  const unique = new Set<string>()
  for (const item of items) {
    const text = String(item ?? "").trim()
    if (!text) continue
    unique.add(text)
  }
  return Array.from(unique)
}
/**
 * Parses a comma-separated enum options string into normalized values.
 */
export const parseEnumOptionsRaw = (raw: string): string[] => {
  return normalizeStringList(String(raw || "").split(","))
}
/**
 * Creates the default meta configuration for a trait data type.
 */
export const defaultKeyMeta = (dataType: DataType = "string"): KeyMeta => {
  const resolvedType = resolveDataType(dataType)
  switch (resolvedType) {
    case "list":
      return {
        dataType: "list",
        ui: {input: "array"},
        listValueType: "string",
        listOrdered: false,
        listUnique: false,
        minItems: 2
      }
    case "range":
      return {
        dataType: "range",
        ui: {input: "date-range"},
        rangeType: "datetime"
      }
    case "geo":
      return {
        dataType: "geo",
        ui: {input: "geo"},
        geoType: "point"
      }
    case "color":
      return {
        dataType: "color",
        ui: {input: "color"},
        mode: "hex"
      }
    case "number":
      return {
        dataType: "number",
        ui: {input: "number"},
        unitCategory: "unitless",
        unit: ""
      }
    case "boolean": {
      const labels = getDefaultBooleanLabels("answer")
      return {
        dataType: "boolean",
        ui: {input: "boolean"},
        booleanDisplay: labels.value,
        booleanTrueLabel: labels.trueLabel,
        booleanFalseLabel: labels.falseLabel
      }
    }
    case "interval":
      return {
        dataType: "interval",
        ui: {input: "interval"},
        durationUnit: "minutes",
        timezone: "UTC"
      }
    case "schedule":
      return {
        dataType: "schedule",
        ui: {input: "schedule"},
        timezone: "UTC"
      }
    case "datetime":
      return {
        dataType: "datetime",
        ui: {input: "datetime"},
        timezone: "UTC"
      }
    case "validity":
      return {
        dataType: "validity",
        ui: {input: "validity"},
        validityMode: "temporary",
        timezone: "UTC"
      }
    case "surface":
      return {
        dataType: "surface",
        ui: {input: "text"}
      }
    default:
      return {dataType: resolvedType, ui: {input: getDefaultUiInput(resolvedType)}}
  }
}
/**
 * Resolves the effective color mode from key meta.
 */
export const resolveColorMode = (meta?: Partial<KeyMeta> | Record<string, unknown> | null): TraitColorMode | "" => {
  const raw = String((meta as any)?.mode || (meta as any)?.format || "")
    .trim()
    .toLowerCase()
  if (raw === "hex" || raw === "lab" || raw === "spectrum") return raw
  return ""
}
/**
 * Resolves the effective unit category with a safe fallback.
 */
export const resolveUnitCategory = (category?: string | null): TraitUnitCategory => {
  const raw = String(category || "").trim()
  const found = UNIT_CATEGORIES.find(item => item.id === raw)
  return found?.id || "unitless"
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
  if (!meta) return defaultKeyMeta("string")
  const dataType = resolveDataType(isDataType(meta.dataType) ? meta.dataType : "string")
  const base = {...defaultKeyMeta(dataType), ...meta, dataType}
  if (dataType === "list") {
    const options =
      enumOptionsRaw !== undefined
        ? parseEnumOptionsRaw(enumOptionsRaw)
        : normalizeStringList(Array.isArray(base.options) ? base.options : [])
    const normalized: KeyMeta = {
      ...base,
      dataType: "list",
      ui: {input: "array"},
      listValueType: resolveListValueType(base.listValueType || base.optionType),
      listOrdered: Boolean(base.listOrdered),
      listUnique: Boolean(base.listUnique),
      minItems: typeof base.minItems === "number" ? Math.max(0, base.minItems) : 2
    }
    if (options.length) {
      normalized.options = options
    }
    return normalized
  }
  if (dataType === "range") {
    return {
      ...base,
      dataType: "range",
      ui: {input: "date-range"},
      rangeType: resolveRangeValueType(base.rangeType)
    }
  }
  if (dataType === "geo") {
    return {
      ...base,
      dataType: "geo",
      ui: {input: "geo"},
      geoType: resolveGeoType(base.geoType),
      heightUnit: String(base.heightUnit || "").trim(),
      radiusUnit: String(base.radiusUnit || "").trim()
    }
  }
  if (dataType === "datetime" || dataType === "interval" || dataType === "schedule" || dataType === "validity") {
    return {
      ...base,
      dataType,
      ui: {input: getDefaultUiInput(dataType)},
      timezone: String(base.timezone || "UTC").trim() || "UTC"
    }
  }
  if (dataType === "number") {
    const unitCategory = resolveUnitCategory(base.unitCategory)
    const units = getUnitOptionsByCategory(unitCategory)
    const unit = units.find(item => item.value === base.unit)?.value || units[0]?.value || ""
    return {
      ...base,
      ui: {input: "number"},
      unitCategory,
      unit
    }
  }
  if (dataType === "color") {
    const mode = resolveColorMode(base) || "hex"
    return {
      ...base,
      dataType: "color",
      ui: {input: "color"},
      mode,
      ...(typeof base.opacity === "number" ? {opacity: Math.max(0, Math.min(1, base.opacity))} : {})
    }
  }
  if (dataType === "boolean") {
    const labels = getDefaultBooleanLabels(base.booleanDisplay)
    return {
      ...base,
      dataType: "boolean",
      ui: {input: "boolean"},
      booleanDisplay: labels.value,
      booleanTrueLabel: String(base.booleanTrueLabel || labels.trueLabel).trim() || labels.trueLabel,
      booleanFalseLabel: String(base.booleanFalseLabel || labels.falseLabel).trim() || labels.falseLabel
    }
  }
  return {
    ...base,
    dataType,
    ui: {input: getDefaultUiInput(dataType)}
  }
}
/**
 * Produces a stable meta shape for structural equality checks.
 */
export const normalizeMetaForEquality = (meta?: KeyMeta | null): KeyMeta => {
  const normalized = normalizeKeyMeta(meta)
  const base = defaultKeyMeta(normalized.dataType || "string")
  const next: KeyMeta = {
    ...base,
    ...normalized
  }
  if (normalized.dataType === "number") {
    next.unitCategory = normalized.unitCategory || "unitless"
    next.unit = normalized.unit || ""
  }
  if (normalized.dataType === "list") {
    next.options = normalizeStringList(Array.isArray(normalized.options) ? normalized.options : [])
    next.listValueType = resolveListValueType(normalized.listValueType)
    next.listOrdered = Boolean(normalized.listOrdered)
    next.listUnique = Boolean(normalized.listUnique)
    next.minItems = typeof normalized.minItems === "number" ? normalized.minItems : 2
  }
  if (normalized.dataType === "range") {
    next.rangeType = resolveRangeValueType(normalized.rangeType)
  }
  if (normalized.dataType === "geo") {
    next.geoType = resolveGeoType(normalized.geoType)
  }
  if (normalized.dataType === "boolean") {
    const labels = getDefaultBooleanLabels(normalized.booleanDisplay)
    next.booleanDisplay = labels.value
    next.booleanTrueLabel = String(normalized.booleanTrueLabel || labels.trueLabel).trim() || labels.trueLabel
    next.booleanFalseLabel = String(normalized.booleanFalseLabel || labels.falseLabel).trim() || labels.falseLabel
  }
  if (normalized.dataType === "datetime" || normalized.dataType === "interval" || normalized.dataType === "schedule" || normalized.dataType === "validity") {
    next.timezone = String(normalized.timezone || "UTC").trim() || "UTC"
  }
  return next
}
/**
 * Compares two meta objects after canonical normalization.
 */
export const keyMetaEquals = (a?: KeyMeta | null, b?: KeyMeta | null): boolean => {
  return JSON.stringify(normalizeMetaForEquality(a)) === JSON.stringify(normalizeMetaForEquality(b))
}
const re = {
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
const isSafeStringValue = (value: string): boolean =>
  Array.from(value).every(char => {
    const code = char.charCodeAt(0)
    return (code >= 0x20 && code !== 0x7f) || code > 0x7f
  })
/**
 * Parses a canonical local datetime string into a UTC timestamp.
 */
const parseCanonicalDateTimeLocal = (value: string): number | null => {
  const m = String(value || "")
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
  const m = String(value || "")
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
  const parts = value.split(",")
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
  const v = typeof value === "string" ? value.trim() : value
  switch (resolveDataType(meta.dataType)) {
    case "string":
      return typeof v === "string" && v.length > 0 && isSafeStringValue(v)
    case "list": {
      if (Array.isArray(v)) return v.length >= Math.max(Number(meta.minItems || 0), 1)
      if (typeof v !== "string") return false
      const raw = v.trim()
      if (!raw) return false
      if (raw.startsWith("[")) {
        try {
          const parsed = JSON.parse(raw)
          return Array.isArray(parsed) && parsed.length >= Math.max(Number(meta.minItems || 0), 1)
        } catch {
          return false
        }
      }
      return raw.length > 0
    }
    case "number":
      return typeof v === "number" ? !Number.isNaN(v) : typeof v === "string" && re.number.test(v)
    case "datetime":
      return typeof v === "string" && isCanonicalDateTimeLocal(v)
    case "range": {
      if (typeof v !== "string") return false
      if (v.startsWith("{")) {
        try {
          const parsed = JSON.parse(v) as {start?: string; end?: string}
          const start = String(parsed?.start || "").trim()
          const end = String(parsed?.end || "").trim()
          const rangeType = resolveRangeValueType(meta.rangeType)
          if (rangeType === "number") return re.number.test(start) && re.number.test(end) && Number(start) <= Number(end)
          if (rangeType === "time") return isCanonicalTimeRangeValid(start, end)
          if (rangeType === "date") return /^\d{4}-\d{2}-\d{2}$/.test(start) && /^\d{4}-\d{2}-\d{2}$/.test(end) && start <= end
          return isCanonicalRangeValid(start, end)
        } catch {
          return false
        }
      }
      const m = v.match(re.rangeDash)
      if (!m) return false
      return isCanonicalRangeValid(String(m[1] || ""), String(m[2] || ""))
    }
    case "interval": {
      if (typeof v !== "string") return false
      const duration = v.match(re.intervalDuration)
      if (duration) {
        const start = String(duration[1] || "")
        const end = String(duration[2] || "")
        const amount = Number(duration[3])
        const unit = String(duration[4] || "")
        if (!isCanonicalRangeValid(start, end)) return false
        if (!Number.isFinite(amount) || amount < 0) return false
        return ["seconds", "minutes", "hours", "days", "years"].includes(unit)
      }
      // Backward compatibility for old interval formats.
      const weekly = v.match(re.intervalWeekly)
      if (weekly) {
        const fromDay = Number(weekly[1])
        const toDay = Number(weekly[2])
        const fromTime = String(weekly[3] || "")
        const toTime = String(weekly[4] || "")
        if (!Number.isInteger(fromDay) || !Number.isInteger(toDay)) return false
        if (fromDay < 1 || fromDay > 7 || toDay < 1 || toDay > 7) return false
        if (parseCanonicalTimeLocal(fromTime) === null || parseCanonicalTimeLocal(toTime) === null) return false
        if (fromDay === toDay) return isCanonicalTimeRangeValid(fromTime, toTime)
        return true
      }
      // Backward compatibility: previously interval was stored as datetime/datetime with "/".
      const legacy = v.match(re.rangeSlash)
      if (!legacy) return false
      return isCanonicalRangeValid(String(legacy[1] || ""), String(legacy[2] || ""))
    }
    case "schedule": {
      if (typeof v !== "string") return false
      const weekly = v.match(re.intervalWeekly)
      if (!weekly) return false
      const fromDay = Number(weekly[1])
      const toDay = Number(weekly[2])
      const fromTime = String(weekly[3] || "")
      const toTime = String(weekly[4] || "")
      if (!Number.isInteger(fromDay) || !Number.isInteger(toDay)) return false
      if (fromDay < 1 || fromDay > 7 || toDay < 1 || toDay > 7) return false
      if (parseCanonicalTimeLocal(fromTime) === null || parseCanonicalTimeLocal(toTime) === null) return false
      if (fromDay === toDay) return isCanonicalTimeRangeValid(fromTime, toTime)
      return true
    }
    case "geo": {
      if (typeof v !== "string") return false
      if (re.geoPoint.test(v)) {
        const parts = v.split(",")
        if (parts.length !== 2) return false
        const [latS, lonS] = parts
        const lat = Number(latS)
        const lon = Number(lonS)
        return lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
      }
      try {
        const parsed = JSON.parse(v) as {type?: string; lat?: string; lng?: string; points?: Array<{lat?: string; lng?: string}>; radius?: string}
        const geoType = resolveGeoType(parsed?.type || meta.geoType)
        if (geoType === "point") {
          const lat = Number(parsed?.lat)
          const lng = Number(parsed?.lng)
          return Number.isFinite(lat) && Number.isFinite(lng) && lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180
        }
        if (geoType === "zone") {
          const lat = Number(parsed?.lat)
          const lng = Number(parsed?.lng)
          const radius = Number(parsed?.radius)
          return (
            Number.isFinite(lat) &&
            Number.isFinite(lng) &&
            Number.isFinite(radius) &&
            lat >= -90 &&
            lat <= 90 &&
            lng >= -180 &&
            lng <= 180 &&
            radius >= 0
          )
        }
        const points = Array.isArray(parsed?.points) ? parsed.points : []
        return points.length >= 3 && points.every(point => Number.isFinite(Number(point?.lat)) && Number.isFinite(Number(point?.lng)))
      } catch {
        return false
      }
    }
    case "color": {
      if (typeof v !== "string" || v.length === 0) return false
      const colorMode = resolveColorMode(meta)
      if (colorMode === "hex") return re.hex.test(v)
      if (colorMode === "spectrum") return /^[a-z]+$/i.test(v)
      if (colorMode === "lab") return isLabTupleValid(v)
      if (re.hex.test(v)) return true
      if (isLabTupleValid(v)) return true
      try {
        const obj = JSON.parse(v)
        if (!obj || typeof obj !== "object") return false
        const mode = (obj as any).mode
        if (mode === "hex") return typeof (obj as any).hex === "string" && re.hex.test((obj as any).hex)
        if (mode === "spectrum") return typeof (obj as any).spectrum === "string" && (obj as any).spectrum.length > 0
        if (mode === "lab" && (obj as any).lab && typeof (obj as any).lab === "object") {
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
    case "boolean":
      return typeof v === "boolean" || (typeof v === "string" && (re.boolean01.test(v) || re.booleanTF.test(v)))
    case "validity": {
      if (typeof v !== "string" || !v) return false
      const permanent = v.match(re.validityPermanent)
      if (permanent) {
        const since = String(permanent[1] || "")
        return !since || isCanonicalDateTimeLocal(since)
      }
      const temporary = v.match(re.validityTemporary)
      if (!temporary) return false
      return isCanonicalDateTimeLocal(String(temporary[1] || ""))
    }
    case "surface": {
      if (typeof v !== "string" || !v.trim().startsWith("{")) return false
      try {
        const parsed = JSON.parse(v) as {glossCategory?: string; glossGU?: number | string; reliefType?: string; microReliefHeight?: number | string}
        const glossCategory = String(parsed?.glossCategory || "").trim()
        const glossGU = String(parsed?.glossGU ?? "").trim()
        const reliefType = String(parsed?.reliefType || "").trim()
        const microReliefHeight = String(parsed?.microReliefHeight ?? "").trim()
        return Boolean((glossCategory && (re.number.test(glossGU) || !glossGU)) || (reliefType && (re.number.test(microReliefHeight) || !microReliefHeight)))
      } catch {
        return false
      }
    }
    default:
      return false
  }
}
