import {
  canonicalDateText,
  compareDateText,
  currentDateText,
  isValidDateText,
  parseDateText
} from '../composables/useDateText'
import { isHexColor, ensureHexHash } from './traitColor'
import { resolveColorMode } from './traitMeta'
import type {
  Color,
  GeoCoordinateParseResult,
  KeyMeta,
  TraitIntervalModel,
  TraitIntervalUnit,
  TraitScheduleModel
} from '../../shared/types/traits'
export const GEO_MIN_DECIMALS = 6
/**
 * Parses one geo coordinate string into a validated numeric form.
 */
export const parseGeoCoordinate = (value: string): GeoCoordinateParseResult | null => {
  const trimmed = String(value || '').trim()
  if (!trimmed) return null
  if (!/^-?\d+(\.\d+)?$/.test(trimmed)) return null
  const num = Number(trimmed)
  if (!Number.isFinite(num)) return null
  const dot = trimmed.indexOf('.')
  const decimals = dot >= 0 ? trimmed.length - dot - 1 : 0
  return { trimmed, num, decimals }
}
/**
 * Formats a parsed geo coordinate while preserving meaningful precision.
 */
const formatGeoCoordinate = (parsed: GeoCoordinateParseResult): string => {
  return parsed.decimals > GEO_MIN_DECIMALS ? parsed.num.toFixed(GEO_MIN_DECIMALS) : parsed.trimmed
}
/**
 * Normalizes a latitude/longitude pair and marks low-precision input.
 */
export const normalizeGeoPoint = (
  latRaw: string,
  lngRaw: string
): { lat: string; lng: string; lowPrecision: boolean } | null => {
  const lat = parseGeoCoordinate(latRaw)
  const lng = parseGeoCoordinate(lngRaw)
  if (!lat || !lng) return null
  return {
    lat: formatGeoCoordinate(lat),
    lng: formatGeoCoordinate(lng),
    lowPrecision: lat.decimals < GEO_MIN_DECIMALS || lng.decimals < GEO_MIN_DECIMALS
  }
}
/**
 * Validates that both dates are canonical and ordered.
 */
const isDatePairValid = (left: string, right: string, mode: 'datetime'): boolean => {
  if (!left || !right) return false
  if (!isValidDateText(left, mode) || !isValidDateText(right, mode)) return false
  const compare = compareDateText(left, right, mode)
  return compare !== null && compare <= 0
}
const INTERVAL_UNIT_MS: Record<TraitIntervalUnit, number> = {
  seconds: 1000,
  minutes: 60_000,
  hours: 3_600_000,
  days: 86_400_000,
  years: 31_536_000_000
}
/**
 * Returns true when the value is a valid ISO-like weekday index from 1 to 7.
 */
const isWeekDayValue = (value: string): boolean => {
  return /^[1-7]$/.test(value)
}
/**
 * Parses and validates the schedule form model.
 */
const parseIntervalSchedule = (value: unknown): TraitScheduleModel | null => {
  if (!(value && typeof value === 'object')) return null
  const fromDay = String((value as any).fromDay ?? '').trim()
  const toDay = String((value as any).toDay ?? '').trim()
  const fromTime = String((value as any).fromTime ?? '').trim()
  const toTime = String((value as any).toTime ?? '').trim()
  if (!isWeekDayValue(fromDay) || !isWeekDayValue(toDay)) return null
  if (!isValidDateText(fromTime, 'time') || !isValidDateText(toTime, 'time')) return null
  if (fromDay === toDay) {
    const cmp = compareDateText(fromTime, toTime, 'time')
    if (cmp === null || cmp > 0) return null
  }
  return { fromDay, toDay, fromTime, toTime }
}
/**
 * Parses and validates the interval form model.
 */
const parseIntervalDuration = (value: unknown): TraitIntervalModel | null => {
  if (!(value && typeof value === 'object')) return null
  const start = String((value as any).start ?? '').trim()
  const end = String((value as any).end ?? '').trim()
  const unit = String((value as any).unit ?? '').trim() as TraitIntervalUnit
  if (!isValidDateText(start, 'datetime') || !isValidDateText(end, 'datetime')) return null
  if (!Object.keys(INTERVAL_UNIT_MS).includes(unit)) return null
  const cmp = compareDateText(start, end, 'datetime')
  if (cmp === null || cmp > 0) return null
  return { start, end, unit }
}
/**
 * Formats a duration amount without trailing zero noise.
 */
const formatIntervalAmount = (value: number): string => {
  if (Number.isInteger(value)) return String(value)
  return value.toFixed(6).replace(/\.?0+$/, '')
}
/**
 * Returns true when the form value satisfies the required shape for its meta.
 */
export const isTraitFormValueFilled = (meta: KeyMeta, value: unknown): boolean => {
  const type = meta.dataType
  const dateMode = 'datetime'
  switch (type) {
    case 'string':
    case 'enum':
      return typeof value === 'string' && value.trim().length > 0
    case 'number':
      if (typeof value === 'number') return !Number.isNaN(value)
      if (typeof value === 'string') return value.trim().length > 0 && !Number.isNaN(Number(value))
      return false
    case 'boolean':
      return typeof value === 'boolean'
    case 'datetime':
      return typeof value === 'string' && isValidDateText(value.trim(), dateMode)
    case 'datetime-range': {
      if (!(value && typeof value === 'object')) return false
      const start = String((value as any).start ?? '').trim()
      const end = String((value as any).end ?? '').trim()
      return isDatePairValid(start, end, dateMode)
    }
    case 'interval': {
      return Boolean(parseIntervalDuration(value))
    }
    case 'schedule': {
      return Boolean(parseIntervalSchedule(value))
    }
    case 'geo-point': {
      if (!(value && typeof value === 'object')) return false
      const lat = String((value as any).lat ?? '').trim()
      const lng = String((value as any).lng ?? '').trim()
      return Boolean(normalizeGeoPoint(lat, lng))
    }
    case 'validity': {
      if (!(value && typeof value === 'object')) return false
      const mode = String((value as any).mode || '').trim()
      if (mode === 'temporary') {
        const until = String((value as any).until ?? '').trim()
        if (!isValidDateText(until, 'datetime')) return false
        const cmp = compareDateText(until, currentDateText('datetime'), 'datetime')
        return cmp !== null && cmp >= 0
      }
      if (mode === 'permanent') return isValidDateText(String((value as any).since ?? '').trim(), 'datetime')
      return false
    }
    case 'color': {
      const configuredMode = resolveColorMode(meta)
      if (typeof value === 'string') {
        const raw = value.trim()
        if (!raw) return false
        if (configuredMode === 'hex') return isHexColor(raw)
        if (configuredMode === 'lab') {
          const parts = raw.split(',').map((part) => part.trim())
          return parts.length === 3 && parts.every((part) => Number.isFinite(Number(part)))
        }
        return true
      }
      if (!(value && typeof value === 'object')) return false
      const mode =
        configuredMode ||
        String((value as any).mode || '')
          .trim()
          .toLowerCase()
      if (mode === 'hex') return isHexColor(String((value as any).hex ?? '').trim())
      if (mode === 'lab') {
        const lab = (value as any).lab
        return Boolean(lab && ['L', 'a', 'b'].every((key) => Number.isFinite(Number(lab[key]))))
      }
      if (mode === 'spectrum') return Boolean(String((value as any).spectrum ?? '').trim())
      return Boolean(String((value as any).hex ?? (value as any).spectrum ?? '').trim())
    }
    default:
      return false
  }
}
/**
 * Serializes color form data into the stored string representation.
 */
const serializeColorValue = (meta: KeyMeta, value: unknown): string => {
  const configuredMode = resolveColorMode(meta)
  if (typeof value === 'string') {
    const raw = value.trim()
    if (!raw) return ''
    if (configuredMode === 'hex') return ensureHexHash(raw)
    return raw
  }
  if (!(value && typeof value === 'object')) return ''
  const mode =
    configuredMode ||
    String((value as any).mode || '')
      .trim()
      .toLowerCase()
  if (mode === 'hex') {
    const hex = String((value as any).hex ?? '').trim()
    if (!hex) return ''
    return ensureHexHash(hex)
  }
  if (mode === 'lab' && (value as any).lab) {
    const lab = (value as any).lab
    if (!['L', 'a', 'b'].every((key) => Number.isFinite(Number(lab[key])))) return ''
    return `${Number(lab.L)},${Number(lab.a)},${Number(lab.b)}`
  }
  if (mode === 'spectrum') {
    return String((value as any).spectrum ?? '').trim()
  }
  return ''
}
/**
 * Serializes one trait form model into the backend string payload.
 */
export const serializeTraitFormValue = (meta: KeyMeta, value: unknown): string => {
  const type = meta.dataType
  const dateMode = 'datetime'
  switch (type) {
    case 'number':
      if (typeof value === 'string') return value.trim()
      return Number.isFinite(value as any) ? String(value) : ''
    case 'geo-point': {
      if (!(value && typeof value === 'object')) return ''
      const latRaw = String((value as any).lat ?? '').trim()
      const lngRaw = String((value as any).lng ?? '').trim()
      const normalized = normalizeGeoPoint(latRaw, lngRaw)
      if (!normalized) return ''
      return `${normalized.lat},${normalized.lng}`
    }
    case 'datetime': {
      if (typeof value !== 'string') return String(value ?? '').trim()
      const trimmed = value.trim()
      if (!trimmed) return ''
      return canonicalDateText(trimmed, dateMode)
    }
    case 'datetime-range': {
      if (!(value && typeof value === 'object' && 'start' in value && 'end' in value)) return ''
      const start = String((value as any).start ?? '').trim()
      const end = String((value as any).end ?? '').trim()
      if (!start || !end) return ''
      return `${canonicalDateText(start, dateMode)} — ${canonicalDateText(end, dateMode)}`
    }
    case 'interval': {
      const duration = parseIntervalDuration(value)
      if (!duration) return ''
      const startParsed = parseDateText(duration.start, 'datetime')
      const endParsed = parseDateText(duration.end, 'datetime')
      if (!startParsed || !endParsed) return ''
      const ms = endParsed.ts - startParsed.ts
      const unitMs = INTERVAL_UNIT_MS[duration.unit]
      if (!unitMs || ms < 0) return ''
      const amount = formatIntervalAmount(ms / unitMs)
      const startCanonical = canonicalDateText(duration.start, 'datetime')
      const endCanonical = canonicalDateText(duration.end, 'datetime')
      return `duration:${startCanonical}/${endCanonical}|${amount}|${duration.unit}`
    }
    case 'schedule': {
      const schedule = parseIntervalSchedule(value)
      if (!schedule) return ''
      const fromTime = canonicalDateText(schedule.fromTime, 'time')
      const toTime = canonicalDateText(schedule.toTime, 'time')
      return `weekly:${schedule.fromDay}-${schedule.toDay}/${fromTime}-${toTime}`
    }
    case 'validity': {
      if (!(value && typeof value === 'object')) return ''
      const mode = String((value as any).mode || 'permanent').trim()
      if (mode === 'permanent') {
        const since = String((value as any).since ?? '').trim()
        return since ? `permanent:${canonicalDateText(since, 'datetime')}` : 'permanent'
      }
      const until = String((value as any).until ?? '').trim()
      return until ? `temporary:${canonicalDateText(until, 'datetime')}` : ''
    }
    case 'color':
      return serializeColorValue(meta, value)
    default:
      return Array.isArray(value) || typeof value === 'object' ? JSON.stringify(value) : String(value ?? '')
  }
}
/**
 * Creates the default form model for the provided trait meta.
 */
export const defaultTraitFormValue = (meta: KeyMeta): string | number | boolean | Record<string, any> => {
  switch (meta.dataType) {
    case 'number':
      return ''
    case 'boolean':
      return false
    case 'datetime':
      return ''
    case 'datetime-range':
      return { start: '', end: '' }
    case 'interval':
      return { start: '', end: '', unit: 'minutes' }
    case 'schedule':
      return { fromDay: '1', toDay: '5', fromTime: '', toTime: '' }
    case 'geo-point':
      return { lat: '', lng: '' }
    case 'enum':
      return ''
    case 'validity':
      return { mode: 'permanent', since: '', until: '' }
    case 'color': {
      const mode = resolveColorMode(meta) || 'hex'
      if (mode === 'lab') return { mode: 'lab', lab: { L: 50, a: 0, b: 0 } }
      if (mode === 'spectrum') return { mode: 'spectrum', spectrum: '' }
      return { mode: 'hex', hex: '' }
    }
    default:
      return ''
  }
}
/**
 * Converts a canonical datetime string into the editable text format.
 */
const parseDateTimeCanonicalToText = (value: string): string => {
  const match = String(value || '')
    .trim()
    .match(/^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2})$/)
  if (!match) return String(value || '').trim()
  const text = `${match[3]}.${match[2]}.${match[1]} ${match[4]}:${match[5]}`
  return isValidDateText(text, 'datetime') ? text : String(value || '').trim()
}
/**
 * Parses a stored boolean string into a boolean value.
 */
const parseBooleanStoredValue = (raw: string): boolean => {
  const normalized = String(raw || '')
    .trim()
    .toLowerCase()
  return normalized === '1' || normalized === 'true'
}
/**
 * Parses a stored color string into the color form model.
 */
const parseColorStoredValue = (meta: KeyMeta, raw: string): Color => {
  const value = String(raw || '').trim()
  const configuredMode = resolveColorMode(meta)
  if (configuredMode === 'hex') {
    return { mode: 'hex', hex: value ? ensureHexHash(value) : '' }
  }
  const labMatch = value.match(/^\s*(-?\d+(?:\.\d+)?)\s*,\s*(-?\d+(?:\.\d+)?)\s*,\s*(-?\d+(?:\.\d+)?)\s*$/)
  const parsedLab = labMatch
    ? {
        L: Number(labMatch[1]),
        a: Number(labMatch[2]),
        b: Number(labMatch[3])
      }
    : null
  if (configuredMode === 'lab') {
    return parsedLab ? { mode: 'lab', lab: parsedLab } : { mode: 'lab', lab: { L: 50, a: 0, b: 0 } }
  }
  if (configuredMode === 'spectrum') {
    return { mode: 'spectrum', spectrum: value }
  }
  if (isHexColor(value)) return { mode: 'hex', hex: ensureHexHash(value) }
  if (parsedLab) return { mode: 'lab', lab: parsedLab }
  return { mode: 'spectrum', spectrum: value }
}
/**
 * Parses the stored backend value into the corresponding trait form model.
 */
export const parseTraitStoredValue = (
  meta: KeyMeta,
  rawValue: string
): string | number | boolean | Record<string, any> => {
  const raw = String(rawValue ?? '').trim()
  switch (meta.dataType) {
    case 'number':
    case 'string':
    case 'enum':
      return raw
    case 'boolean':
      return parseBooleanStoredValue(raw)
    case 'datetime':
      return parseDateTimeCanonicalToText(raw)
    case 'datetime-range': {
      const dashMatch = raw.match(/^(.+)\s—\s(.+)$/)
      if (dashMatch) {
        return {
          start: parseDateTimeCanonicalToText(String(dashMatch[1] || '').trim()),
          end: parseDateTimeCanonicalToText(String(dashMatch[2] || '').trim())
        }
      }
      const slashMatch = raw.match(/^(.+)\/(.+)$/)
      if (slashMatch) {
        return {
          start: parseDateTimeCanonicalToText(String(slashMatch[1] || '').trim()),
          end: parseDateTimeCanonicalToText(String(slashMatch[2] || '').trim())
        }
      }
      return defaultTraitFormValue(meta)
    }
    case 'interval': {
      const durationMatch = raw.match(
        /^duration:([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2})\/([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2})\|(-?\d+(?:\.\d+)?)\|(seconds|minutes|hours|days|years)$/
      )
      if (durationMatch) {
        return {
          start: parseDateTimeCanonicalToText(String(durationMatch[1] || '')),
          end: parseDateTimeCanonicalToText(String(durationMatch[2] || '')),
          unit: String(durationMatch[4] || 'minutes')
        }
      }
      const slashMatch = raw.match(/^(.+)\/(.+)$/)
      if (slashMatch) {
        return {
          start: parseDateTimeCanonicalToText(String(slashMatch[1] || '').trim()),
          end: parseDateTimeCanonicalToText(String(slashMatch[2] || '').trim()),
          unit: 'minutes'
        }
      }
      return defaultTraitFormValue(meta)
    }
    case 'schedule': {
      const weekly = raw.match(/^weekly:([1-7])-([1-7])\/([0-9]{2}:[0-9]{2})-([0-9]{2}:[0-9]{2})$/)
      if (!weekly) return defaultTraitFormValue(meta)
      return {
        fromDay: String(weekly[1] || '1'),
        toDay: String(weekly[2] || '5'),
        fromTime: String(weekly[3] || ''),
        toTime: String(weekly[4] || '')
      }
    }
    case 'geo-point': {
      const parts = raw.split(',').map((part) => part.trim())
      if (parts.length !== 2) return defaultTraitFormValue(meta)
      return { lat: parts[0], lng: parts[1] }
    }
    case 'validity': {
      const permanent = raw.match(/^permanent(?::(.+))?$/)
      if (permanent) {
        return {
          mode: 'permanent',
          since: parseDateTimeCanonicalToText(String(permanent[1] || '').trim()),
          until: ''
        }
      }
      const temporary = raw.match(/^temporary:(.+)$/)
      if (temporary) {
        return {
          mode: 'temporary',
          since: '',
          until: parseDateTimeCanonicalToText(String(temporary[1] || '').trim())
        }
      }
      return defaultTraitFormValue(meta)
    }
    case 'color':
      return parseColorStoredValue(meta, raw)
    default:
      return raw
  }
}
