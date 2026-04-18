export const INTERFACE_LOCALE_OPTIONS = [
  {code: "ru", shortLabel: "RU", nativeLabel: "Русский"},
  {code: "en", shortLabel: "EN", nativeLabel: "English"},
  {code: "zh", shortLabel: "ZH", nativeLabel: "中文"},
  {code: "ja", shortLabel: "JA", nativeLabel: "日本語"}
] as const satisfies ReadonlyArray<{
  code: InterfaceLocaleCode
  shortLabel: string
  nativeLabel: string
}>

export const normalizeInterfaceLocale = (raw: unknown): InterfaceLocaleCode => {
  const normalized = String(raw || "")
    .trim()
    .toLowerCase()
  if (normalized === "ru" || normalized === "ru-ru") return "ru"
  if (normalized === "zh" || normalized === "zh-cn" || normalized === "zh-hans" || normalized === "ch") return "zh"
  if (normalized === "ja" || normalized === "ja-jp" || normalized === "jp") return "ja"
  return "en"
}

export const interfaceLocaleToTag = (locale: InterfaceLocaleCode) =>
  ({ru: "ru-RU", en: "en-US", zh: "zh-CN", ja: "ja-JP"})[normalizeInterfaceLocale(locale)]

export const INTERFACE_THEME_OPTIONS = ["system", "dark", "light"] as const satisfies ReadonlyArray<ThemePreference>
export const normalizeThemePreference = (raw: unknown): ThemePreference => {
  const normalized = String(raw || "")
    .trim()
    .toLowerCase()
  if (normalized === "light") return "light"
  if (normalized === "dark") return "dark"
  return "system"
}
