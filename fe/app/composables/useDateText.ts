export type DateTextMode = 'date' | 'time' | 'datetime'
type Parsed = { canonical: string; ts: number }
const pad = (n: number, len = 2) => String(n).padStart(len, '0')
const digitsLimitByMode: Record<DateTextMode, number> = {
  date: 8,
  time: 4,
  datetime: 12
}
const textLengthLimitByMode: Record<DateTextMode, number> = {
  date: 10,
  time: 5,
  datetime: 16
}
const parseDateRu = (value: string): { year: number; month: number; day: number; ts: number } | null => {
  const m = String(value || '')
    .trim()
    .match(/^(\d{2})\.(\d{2})\.(\d{4})$/)
  if (!m) return null
  const day = Number(m[1])
  const month = Number(m[2])
  const year = Number(m[3])
  if (year < 1 || month < 1 || month > 12 || day < 1 || day > 31) return null
  // Date.UTC maps years 0..99 to 1900..1999, so set full year explicitly.
  // Build date on a leap year baseline to keep day/month intent before validation.
  const dt = new Date(Date.UTC(2000, month - 1, day))
  dt.setUTCFullYear(year)
  if (dt.getUTCFullYear() !== year || dt.getUTCMonth() !== month - 1 || dt.getUTCDate() !== day) return null
  return { year, month, day, ts: dt.getTime() }
}
const parseTime = (value: string): { hour: number; minute: number; ts: number } | null => {
  const m = String(value || '')
    .trim()
    .match(/^(\d{2}):(\d{2})$/)
  if (!m) return null
  const hour = Number(m[1])
  const minute = Number(m[2])
  if (hour < 0 || hour > 23 || minute < 0 || minute > 59) return null
  return { hour, minute, ts: hour * 60 + minute }
}
const parseDateTimeRu = (value: string): Parsed | null => {
  const m = String(value || '')
    .trim()
    .match(/^(\d{2}\.\d{2}\.\d{4})\s+(\d{2}:\d{2})$/)
  if (!m) return null
  const datePart = m[1] || ''
  const timePart = m[2] || ''
  const d = parseDateRu(datePart)
  const t = parseTime(timePart)
  if (!d || !t) return null
  const ts = Date.UTC(d.year, d.month - 1, d.day, t.hour, t.minute, 0, 0)
  return { canonical: `${pad(d.year, 4)}-${pad(d.month)}-${pad(d.day)}T${pad(t.hour)}:${pad(t.minute)}`, ts }
}
export const maskDateText = (raw: string, mode: DateTextMode): string => {
  const digits = String(raw || '')
    .replace(/\D/g, '')
    .slice(0, digitsLimitByMode[mode])
  if (mode === 'time') {
    const d = digits
    const hh = d.slice(0, 2)
    const mm = d.slice(2, 4)
    if (d.length <= 2) return hh
    return `${hh}:${mm}`
  }
  if (mode === 'date') {
    const d = digits
    const dd = d.slice(0, 2)
    const mm = d.slice(2, 4)
    const yyyy = d.slice(4, 8)
    if (d.length <= 2) return dd
    if (d.length <= 4) return `${dd}.${mm}`
    return `${dd}.${mm}.${yyyy}`
  }
  const d = digits
  const dd = d.slice(0, 2)
  const mm = d.slice(2, 4)
  const yyyy = d.slice(4, 8)
  const hh = d.slice(8, 10)
  const mi = d.slice(10, 12)
  if (d.length <= 2) return dd
  if (d.length <= 4) return `${dd}.${mm}`
  if (d.length <= 8) return `${dd}.${mm}.${yyyy}`
  if (d.length <= 10) return `${dd}.${mm}.${yyyy} ${hh}`
  return `${dd}.${mm}.${yyyy} ${hh}:${mi}`
}
export const parseDateText = (value: string, mode: DateTextMode): Parsed | null => {
  if (mode === 'time') {
    const t = parseTime(value)
    if (!t) return null
    return { canonical: `${pad(t.hour)}:${pad(t.minute)}`, ts: t.ts }
  }
  if (mode === 'date') {
    const d = parseDateRu(value)
    if (!d) return null
    return { canonical: `${pad(d.year, 4)}-${pad(d.month)}-${pad(d.day)}`, ts: d.ts }
  }
  return parseDateTimeRu(value)
}
export const isValidDateText = (value: string, mode: DateTextMode): boolean => {
  return Boolean(parseDateText(value, mode))
}
export const compareDateText = (a: string, b: string, mode: DateTextMode): number | null => {
  const pa = parseDateText(a, mode)
  const pb = parseDateText(b, mode)
  if (!pa || !pb) return null
  return pa.ts - pb.ts
}
export const canonicalDateText = (value: string, mode: DateTextMode): string => {
  return parseDateText(value, mode)?.canonical || String(value || '').trim()
}
export const datePlaceholder = (mode: DateTextMode): string => {
  if (mode === 'time') return 'чч:мм'
  if (mode === 'date') return 'дд.мм.гггг'
  return 'дд.мм.гггг чч:мм'
}
export const dateTextMaxLength = (mode: DateTextMode): number => {
  return textLengthLimitByMode[mode]
}
export const currentDateText = (mode: DateTextMode): string => {
  const now = new Date()
  const year = now.getFullYear()
  const month = now.getMonth() + 1
  const day = now.getDate()
  const hour = now.getHours()
  const minute = now.getMinutes()
  if (mode === 'time') return `${pad(hour)}:${pad(minute)}`
  if (mode === 'date') return `${pad(day)}.${pad(month)}.${pad(year, 4)}`
  return `${pad(day)}.${pad(month)}.${pad(year, 4)} ${pad(hour)}:${pad(minute)}`
}
export const detectDateMode = (_format?: string): DateTextMode => {
  return 'datetime'
}
