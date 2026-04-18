const CYRILLIC_TO_LATIN: Record<string, string> = {
  а: "a",
  б: "b",
  в: "v",
  г: "g",
  д: "d",
  е: "e",
  ё: "e",
  ж: "zh",
  з: "z",
  и: "i",
  й: "y",
  к: "k",
  л: "l",
  м: "m",
  н: "n",
  о: "o",
  п: "p",
  р: "r",
  с: "s",
  т: "t",
  у: "u",
  ф: "f",
  х: "h",
  ц: "ts",
  ч: "ch",
  ш: "sh",
  щ: "sch",
  ъ: "",
  ы: "y",
  ь: "",
  э: "e",
  ю: "yu",
  я: "ya"
}
export const transliterateCyrillicToLatin = (value: string) =>
  Array.from(String(value || "").toLowerCase())
    .map(char => CYRILLIC_TO_LATIN[char] ?? char)
    .join("")
type SlugifyOptions = {
  fallback?: string
  maxLength?: number
}
export const slugifyLatin = (raw: string, options: SlugifyOptions = {}) => {
  const fallback = String(options.fallback || "item").trim() || "item"
  const maxLength = Number(options.maxLength || 0)
  const source = transliterateCyrillicToLatin(raw)
    .normalize("NFKD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-+|-+$/g, "")
  const clipped = maxLength > 0 ? source.slice(0, maxLength).replace(/-+$/g, "") : source
  return clipped || fallback
}
