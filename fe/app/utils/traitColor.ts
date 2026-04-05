import type { DataType, KeyMeta, TraitColorPreviewInfo } from '../../shared/types/traits'
export const VALUE_COMPONENTS = {
  string: 'TraitsFormInputString',
  number: 'TraitsFormInputNumber',
  boolean: 'TraitsFormInputBoolean',
  datetime: 'TraitsFormInputDatetime',
  'datetime-range': 'TraitsFormInputDateRange',
  interval: 'TraitsFormInputInterval',
  schedule: 'TraitsFormInputSchedule',
  'geo-point': 'TraitsFormInputGeoPoint',
  enum: 'TraitsFormInputEnum',
  validity: 'TraitsFormInputValidity',
  color: 'TraitsFormInputColor'
} as const satisfies Record<DataType, string>
export const COLOR_SPECTRUM_OPTIONS = [
  'red',
  'orange',
  'yellow',
  'green',
  'cyan',
  'blue',
  'purple',
  'pink',
  'brown',
  'gray',
  'black',
  'white'
] as const
export const COLOR_SPECTRUM_LABELS: Record<(typeof COLOR_SPECTRUM_OPTIONS)[number], string> = {
  red: 'красный',
  orange: 'оранжевый',
  yellow: 'желтый',
  green: 'зеленый',
  cyan: 'голубой',
  blue: 'синий',
  purple: 'фиолетовый',
  pink: 'розовый',
  brown: 'коричневый',
  gray: 'серый',
  black: 'черный',
  white: 'белый'
}
export const COLOR_SPECTRUM_MAP: Record<string, string> = {
  red: '#ef4444',
  orange: '#f97316',
  yellow: '#eab308',
  green: '#22c55e',
  cyan: '#06b6d4',
  blue: '#3b82f6',
  // purple: '#8b5cf6',
  purple: '#663399',
  pink: '#ec4899',
  brown: '#92400e',
  gray: '#6b7280',
  black: '#0f172a',
  white: '#e5e7eb'
}
/**
 * Returns true when the string is a 6-digit hexadecimal color.
 */
export const isHexColor = (value: string): boolean => /^#?[0-9a-fA-F]{6}$/.test(String(value || '').trim())
/**
 * Ensures that a hexadecimal color value starts with `#`.
 */
export const ensureHexHash = (value: string): string => {
  const raw = String(value || '').trim()
  if (!raw) return ''
  return raw.startsWith('#') ? raw : `#${raw}`
}
/**
 * Converts a hexadecimal color string into an RGB tuple.
 */
export const hexToRgb = (hex: string): { r: number; g: number; b: number } | null => {
  const clean = String(hex || '')
    .trim()
    .replace('#', '')
  if (![3, 6].includes(clean.length)) return null
  const full =
    clean.length === 3
      ? clean
          .split('')
          .map(char => char + char)
          .join('')
      : clean
  const num = Number.parseInt(full, 16)
  if (Number.isNaN(num)) return null
  return {
    r: (num >> 16) & 255,
    g: (num >> 8) & 255,
    b: num & 255
  }
}
/**
 * Clamps a number to the inclusive `[0, 1]` range.
 */
const clamp01 = (value: number) => Math.min(Math.max(value, 0), 1)
/**
 * Converts LAB coordinates into an RGB tuple.
 */
export const labToRgb = (L: number, a: number, b: number) => {
  const pivot = (n: number) => {
    const n3 = n * n * n
    return n3 > 0.008856 ? n3 : (n - 16 / 116) / 7.787
  }
  let y = (L + 16) / 116
  let x = a / 500 + y
  let z = y - b / 200
  x = pivot(x)
  y = pivot(y)
  z = pivot(z)
  x *= 0.95047
  y *= 1
  z *= 1.08883
  let r = x * 3.2406 + y * -1.5372 + z * -0.4986
  let g = x * -0.9689 + y * 1.8758 + z * 0.0415
  let blue = x * 0.0557 + y * -0.204 + z * 1.057
  const toSrgb = (channel: number) => {
    const normalized = clamp01(channel)
    return normalized <= 0.0031308 ? 12.92 * normalized : 1.055 * Math.pow(normalized, 1 / 2.4) - 0.055
  }
  r = toSrgb(r)
  g = toSrgb(g)
  blue = toSrgb(blue)
  return {
    r: Math.round(clamp01(r) * 255),
    g: Math.round(clamp01(g) * 255),
    b: Math.round(clamp01(blue) * 255)
  }
}
/**
 * Converts an RGB tuple into a CSS `rgb(...)` string.
 */
export const rgbToCss = (rgb: { r: number; g: number; b: number }) => `rgb(${rgb.r}, ${rgb.g}, ${rgb.b})`
/**
 * Parses a `L,a,b` string into numeric LAB coordinates.
 */
export const parseLabCsv = (value: string): { L: number; a: number; b: number } | null => {
  const parts = String(value || '')
    .split(',')
    .map(part => Number(part.trim()))
  if (parts.length !== 3 || !parts.every(Number.isFinite)) return null
  return {
    L: parts[0] ?? 0,
    a: parts[1] ?? 0,
    b: parts[2] ?? 0
  }
}
/**
 * Builds preview data for one named spectrum color.
 */
const spectrumInfo = (name: string): TraitColorPreviewInfo => {
  return {
    css: COLOR_SPECTRUM_MAP[name] || '#6b7280',
    mode: 'spectrum',
    text: name
  }
}
/**
 * Builds preview data for one LAB color value.
 */
const labInfo = (L: number, a: number, b: number): TraitColorPreviewInfo => {
  const rgb = labToRgb(L, a, b)
  return {
    css: rgbToCss(rgb),
    mode: 'lab',
    text: `L${Math.round(L)} a${Math.round(a)} b${Math.round(b)}`
  }
}
/**
 * Builds preview data for one hexadecimal color.
 */
const hexInfo = (value: string): TraitColorPreviewInfo => {
  const hex = ensureHexHash(value)
  const rgb = hexToRgb(hex) || { r: 0, g: 0, b: 0 }
  return {
    css: hex,
    mode: 'hex',
    text: `${hex.toUpperCase()} (${rgb.r},${rgb.g},${rgb.b})`
  }
}
/**
 * Resolves the best preview representation for a serialized color value.
 */
export const resolveColorPreviewInfo = (
  rawValue: string,
  meta?: Partial<KeyMeta> | Record<string, unknown> | null
): TraitColorPreviewInfo | null => {
  const raw = String(rawValue || '').trim()
  if (!raw) return null
  const configuredMode = resolveColorMode(meta)
  const parsed = (() => {
    try {
      return JSON.parse(raw)
    } catch {
      return null
    }
  })()
  if (parsed && typeof parsed === 'object' && 'mode' in parsed) {
    const mode = String((parsed as any).mode || '')
      .trim()
      .toLowerCase()
    if (mode === 'hex' && typeof (parsed as any).hex === 'string') {
      return hexInfo((parsed as any).hex)
    }
    if (mode === 'lab' && (parsed as any).lab) {
      const lab = (parsed as any).lab
      return labInfo(Number(lab.L || 0), Number(lab.a || 0), Number(lab.b || 0))
    }
    if (mode === 'spectrum' && typeof (parsed as any).spectrum === 'string') {
      return spectrumInfo((parsed as any).spectrum)
    }
  }
  if (configuredMode === 'spectrum' && /^[a-z]+$/i.test(raw)) {
    return spectrumInfo(raw)
  }
  if (configuredMode === 'lab') {
    const lab = parseLabCsv(raw)
    if (lab) return labInfo(lab.L, lab.a, lab.b)
  }
  if (isHexColor(raw)) {
    return hexInfo(raw)
  }
  if (/^[a-z]+$/i.test(raw)) {
    return spectrumInfo(raw)
  }
  const lab = parseLabCsv(raw)
  if (lab) return labInfo(lab.L, lab.a, lab.b)
  return null
}
