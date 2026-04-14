import type { DocContent, DocsSlug } from '../types/docs'
import type { LabTabItem } from '../types/ui'

export const DOCS_SLUGS = ['company', 'offer', 'terms', 'privacy'] as const satisfies ReadonlyArray<DocsSlug>

export const DOCS_DEFAULT_SLUG: DocsSlug = 'company'

export const normalizeDocsSlug = (value: unknown): DocsSlug | null => {
  if (typeof value === 'string' && DOCS_SLUGS.includes(value as DocsSlug)) {
    return value as DocsSlug
  }
  if (Array.isArray(value) && typeof value[0] === 'string' && DOCS_SLUGS.includes(value[0] as DocsSlug)) {
    return value[0] as DocsSlug
  }
  return null
}

export const isDocsSlug = (value: unknown): value is DocsSlug => normalizeDocsSlug(value) !== null

export const getDocsHref = (slug: DocsSlug): string => `/docs/${slug}`

export const getDocsTabItems = (documents: Record<DocsSlug, Pick<DocContent, 'title'>>): LabTabItem[] =>
  DOCS_SLUGS.map(
    (slug): LabTabItem => ({
      value: slug,
      label: documents[slug as DocsSlug].title
    })
  )

export const getDocsTabRouteMap = (): Record<string, string> =>
  Object.fromEntries(DOCS_SLUGS.map((slug) => [slug, getDocsHref(slug)]))
