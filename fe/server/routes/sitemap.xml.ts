import {getRequestURL, setHeader} from "h3"
type SitemapEntry = {
  path: string
  lastmod?: string
  changefreq?: "daily" | "weekly" | "monthly"
  priority?: string
}
const normalizePath = (rawPath: string) => {
  const path = String(rawPath || "/").replace(/\/{2,}/g, "/")
  if (path.length <= 1) return "/"
  return path.replace(/\/+$/, "")
}
const xmlEscape = (value: string) =>
  value
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&apos;")
const staticEntries: SitemapEntry[] = [
  {path: "/", changefreq: "weekly", priority: "1.0"},
  {path: "/updates", changefreq: "daily", priority: "0.9"},
  {path: "/edu", changefreq: "weekly", priority: "0.9"},
  {path: "/impact", changefreq: "weekly", priority: "0.9"},
  {path: "/kitchen", changefreq: "weekly", priority: "0.8"},
  {path: "/lab", changefreq: "weekly", priority: "0.8"},
  {path: "/terms", changefreq: "monthly", priority: "0.4"},
  {path: "/dev", changefreq: "monthly", priority: "0.3"},
  {path: "/experience", changefreq: "monthly", priority: "0.3"},
  {path: "/perception", changefreq: "monthly", priority: "0.3"},
  {path: "/careers", changefreq: "monthly", priority: "0.3"},
  {path: "/family", changefreq: "monthly", priority: "0.3"}
]
const eduSlugs = [
  "algebra",
  "geometry-trig",
  "calculus-linear-prob",
  "physics-core",
  "chemistry",
  "materials",
  "mechanics-strength",
  "machine-elements",
  "manufacturing",
  "design-standards",
  "electrical",
  "economics-quality",
  "safety-ergonomics"
] as const
const impactSlugs = [
  "minimalism-lifestyle",
  "hygiene-daily-ecology",
  "water-energy-efficiency",
  "circular-reuse-repair",
  "biodiversity-support",
  "urban-green-mobility",
  "community-initiatives",
  "education-values"
] as const
const updateEntries: Array<{slug: string; lastmod: string}> = [
  {slug: "2026-03-04-prerelease-functional-overview", lastmod: "2026-03-04"}
]
const dynamicEntries = (): SitemapEntry[] => {
  const updates = updateEntries.map(entry => ({
    path: `/updates/${entry.slug}`,
    lastmod: entry.lastmod,
    changefreq: "monthly" as const,
    priority: "0.7"
  }))
  const edu = eduSlugs.map(slug => ({
    path: `/edu/${slug}`,
    changefreq: "monthly" as const,
    priority: "0.7"
  }))
  const impact = impactSlugs.map(slug => ({
    path: `/impact/${slug}`,
    changefreq: "monthly" as const,
    priority: "0.7"
  }))
  return [...updates, ...edu, ...impact]
}
export default defineEventHandler(event => {
  const config = useRuntimeConfig(event)
  const requestUrl = getRequestURL(event)
  const configuredBase = String(config.public.baseURL || "").trim()
  const siteUrl = (configuredBase || `${requestUrl.protocol}//${requestUrl.host}`).replace(/\/+$/, "")
  const entries = [...staticEntries, ...dynamicEntries()]
  const seen = new Set<string>()
  const uniqueEntries = entries.filter(entry => {
    const path = normalizePath(entry.path)
    if (seen.has(path)) return false
    seen.add(path)
    return true
  })
  const body = [
    '<?xml version="1.0" encoding="UTF-8"?>',
    '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">',
    ...uniqueEntries.map(entry => {
      const path = normalizePath(entry.path)
      const loc = xmlEscape(`${siteUrl}${path}`)
      const lastmod = entry.lastmod ? `\n    <lastmod>${xmlEscape(entry.lastmod)}</lastmod>` : ""
      const changefreq = entry.changefreq ? `\n    <changefreq>${entry.changefreq}</changefreq>` : ""
      const priority = entry.priority ? `\n    <priority>${entry.priority}</priority>` : ""
      return `  <url>\n    <loc>${loc}</loc>${lastmod}${changefreq}${priority}\n  </url>`
    }),
    "</urlset>",
    ""
  ].join("\n")
  setHeader(event, "content-type", "application/xml; charset=utf-8")
  setHeader(event, "cache-control", "public, max-age=3600")
  return body
})
