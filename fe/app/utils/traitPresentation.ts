import type { TraitUnitCategory, DataType } from '../../shared/types/traits'
import { COLOR_MODE_LABELS, DATA_TYPE_LABELS, UNIT_CATEGORY_LABELS, UNIT_LABELS, resolveColorMode } from './traitMeta'
export const META_FIELD_LABELS: Record<string, string> = {
  dataType: 'тип данных',
  mode: 'режим',
  unitCategory: 'категория единицы',
  unit: 'единица измерения',
  optionType: 'тип опции',
  options: 'опции',
  min: 'минимум',
  max: 'максимум',
  precision: 'точность'
}
/**
 * Formats a canonical local datetime string into a compact Russian label.
 */
export const formatDateTimeIsoToRu = (iso: string): string => {
  const match = String(iso || '')
    .trim()
    .match(/^(\d{4})-(\d{2})-(\d{2})[T\s](\d{2}):(\d{2})$/)
  if (!match) return iso
  return `${match[3]}.${match[2]}.${match[1]} ${match[4]}:${match[5]}`
}
const WEEK_DAY_LABELS: Record<string, string> = {
  '1': 'ПН',
  '2': 'ВТ',
  '3': 'СР',
  '4': 'ЧТ',
  '5': 'ПТ',
  '6': 'СБ',
  '7': 'ВС'
}
const INTERVAL_UNIT_LABELS: Record<string, string> = {
  seconds: 'сек',
  minutes: 'мин',
  hours: 'ч',
  days: 'дн',
  years: 'г'
}
/**
 * Returns a human-readable label for the provided trait data type.
 */
export const getDataTypeLabel = (rawType: string, fallback = 'тип'): string => {
  const key = String(rawType || '').trim() as DataType
  return DATA_TYPE_LABELS[key] || key || fallback
}
/**
 * Returns a localized unit category label from meta.
 */
export const getDataCategoryLabel = (meta?: Record<string, any> | null): string => {
  const source = meta || {}
  const raw = String(source.unitCategory || source.category || source.group || '').trim()
  if (!raw) return ''
  return UNIT_CATEGORY_LABELS[raw as TraitUnitCategory] || raw
}
/**
 * Returns a localized label for one numeric unit.
 */
export const getNumberUnitLabel = (unit: string): string => {
  const raw = String(unit || '').trim()
  return UNIT_LABELS[raw] || raw
}
/**
 * Returns a localized label for the current color mode.
 */
export const getColorModeLabel = (meta?: Record<string, any> | null): string => {
  const mode = resolveColorMode(meta)
  if (!mode) return ''
  return COLOR_MODE_LABELS[mode] || mode
}
/**
 * Converts a stored trait value into a user-facing string.
 */
export const formatTraitValueForDisplay = (rawValue: string, dataType: string): string => {
  const raw = String(rawValue || '').trim()
  if (!raw) return ''
  if (dataType === 'datetime') {
    return formatDateTimeIsoToRu(raw)
  }
  if (dataType === 'datetime-range') {
    const [start = '', end = ''] = raw.split(' — ')
    return `${formatDateTimeIsoToRu(start)} — ${formatDateTimeIsoToRu(end)}`
  }
  if (dataType === 'interval') {
    const duration = raw.match(
      /^duration:([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2})\/([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2})\|(-?\d+(?:\.\d+)?)\|(seconds|minutes|hours|days|years)$/
    )
    if (duration) {
      const start = formatDateTimeIsoToRu(String(duration[1] || ''))
      const end = formatDateTimeIsoToRu(String(duration[2] || ''))
      const amount = String(duration[3] || '')
      const unit = INTERVAL_UNIT_LABELS[String(duration[4] || '')] || String(duration[4] || '')
      return `Интервал: ${amount} ${unit} (с ${start} до ${end})`
    }
    // Backward compatibility for old weekly format.
    const weekly = raw.match(/^weekly:([1-7])-([1-7])\/([0-9]{2}:[0-9]{2})-([0-9]{2}:[0-9]{2})$/)
    if (weekly) {
      const fromDay = WEEK_DAY_LABELS[String(weekly[1] || '')] || String(weekly[1] || '')
      const toDay = WEEK_DAY_LABELS[String(weekly[2] || '')] || String(weekly[2] || '')
      const fromTime = String(weekly[3] || '')
      const toTime = String(weekly[4] || '')
      return `Каждую неделю: ${fromDay}-${toDay}, ${fromTime}-${toTime}`
    }
    // Backward compatibility for legacy format (absolute datetime/datetime).
    const [from = '', to = ''] = raw.split('/')
    return `${formatDateTimeIsoToRu(from)} / ${formatDateTimeIsoToRu(to)}`
  }
  if (dataType === 'schedule') {
    const weekly = raw.match(/^weekly:([1-7])-([1-7])\/([0-9]{2}:[0-9]{2})-([0-9]{2}:[0-9]{2})$/)
    if (weekly) {
      const fromDay = WEEK_DAY_LABELS[String(weekly[1] || '')] || String(weekly[1] || '')
      const toDay = WEEK_DAY_LABELS[String(weekly[2] || '')] || String(weekly[2] || '')
      const fromTime = String(weekly[3] || '')
      const toTime = String(weekly[4] || '')
      return `Каждую неделю: ${fromDay}-${toDay}, ${fromTime}-${toTime}`
    }
    return raw
  }
  if (dataType === 'validity') {
    if (raw.startsWith('permanent:')) {
      return `Постоянно с ${formatDateTimeIsoToRu(raw.slice('permanent:'.length))}`
    }
    if (raw.startsWith('temporary:')) {
      return `Временно до ${formatDateTimeIsoToRu(raw.slice('temporary:'.length))}`
    }
  }
  return raw
}
/**
 * Converts one meta field value into a compact display string.
 */
export const formatMetaValue = (key: string, value: unknown, typeLabel: string): string => {
  if (value === null || value === undefined) return ''
  if (key === 'dataType') return typeLabel
  if (key === 'unitCategory') {
    const raw = String(value).trim()
    return UNIT_CATEGORY_LABELS[raw as TraitUnitCategory] || raw
  }
  if (key === 'unit') {
    return getNumberUnitLabel(String(value))
  }
  if (key === 'mode') {
    const rawMode = String(value).trim().toLowerCase() as keyof typeof COLOR_MODE_LABELS
    return COLOR_MODE_LABELS[rawMode] || rawMode
  }
  if (Array.isArray(value)) {
    return value
      .map((item) => String(item).trim())
      .filter(Boolean)
      .join(', ')
  }
  if (typeof value === 'object') {
    try {
      return JSON.stringify(value)
    } catch {
      return String(value)
    }
  }
  return String(value).trim()
}
