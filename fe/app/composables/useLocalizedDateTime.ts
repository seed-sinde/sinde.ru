const normalizeLocalizedDate = (value: number | string | Date | null | undefined): Date | null => {
  if (!value && value !== 0) return null
  if (value instanceof Date) {
    return Number.isNaN(value.getTime()) ? null : value
  }
  if (typeof value === "number") {
    if (!Number.isFinite(value)) return null
    const ms = value > 1e11 ? value : value * 1000
    const result = new Date(ms)
    return Number.isNaN(result.getTime()) ? null : result
  }
  const normalized = String(value).trim()
  if (!normalized) return null
  if (/^\d+$/.test(normalized)) {
    const asNumber = Number(normalized)
    if (!Number.isFinite(asNumber)) return null
    const ms = asNumber > 1e11 ? asNumber : asNumber * 1000
    const result = new Date(ms)
    return Number.isNaN(result.getTime()) ? null : result
  }
  const result = new Date(normalized)
  return Number.isNaN(result.getTime()) ? null : result
}

export const useLocalizedDateTime = () => {
  const {localeTag} = useInterfacePreferences()
  const {load, t} = useI18nSection("ui")
  onServerPrefetch(load)
  if (import.meta.client) void load()
  const formatAbsoluteDateTime = (
    value: number | string | Date | null | undefined,
    options?: Intl.DateTimeFormatOptions,
    fallback = "—"
  ) => {
    const date = normalizeLocalizedDate(value)
    if (!date) return fallback
    return new Intl.DateTimeFormat(localeTag.value, options || {dateStyle: "medium", timeStyle: "short"}).format(date)
  }
  const formatRelativeTime = (
    value: number | string | Date | null | undefined,
    fallback = "—",
    numeric: Intl.RelativeTimeFormatNumeric = "always"
  ) => {
    const date = normalizeLocalizedDate(value)
    if (!date) return fallback
    const diffMs = date.getTime() - Date.now()
    const absMs = Math.abs(diffMs)
    if (absMs < 60_000) {
      return diffMs <= 0 ? t("time.less_than_minute_ago") : t("time.in_less_than_minute")
    }
    const rtf = new Intl.RelativeTimeFormat(localeTag.value, {numeric, style: "long"})
    const units: ReadonlyArray<[Intl.RelativeTimeFormatUnit, number]> = [
      ["year", 1000 * 60 * 60 * 24 * 365],
      ["month", 1000 * 60 * 60 * 24 * 30],
      ["week", 1000 * 60 * 60 * 24 * 7],
      ["day", 1000 * 60 * 60 * 24],
      ["hour", 1000 * 60 * 60],
      ["minute", 1000 * 60]
    ]
    for (const [unit, size] of units) {
      if (absMs >= size) {
        return rtf.format(Math.round(diffMs / size), unit)
      }
    }
    return fallback
  }
  return {
    localeTag,
    normalizeLocalizedDate,
    formatAbsoluteDateTime,
    formatRelativeTime
  }
}
