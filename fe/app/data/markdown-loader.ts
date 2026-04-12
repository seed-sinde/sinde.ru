export type DocMeta = {
  slug: string
  title: string
  section: string
  summary: string
}
export type DocRecord = DocMeta & {
  file: string
}
export type ContentDoc = DocMeta & {
  content: string
}
const CONTENT_LOADERS = import.meta.glob('../content/**/*.md', {
  query: '?raw',
  import: 'default'
}) as Record<string, () => Promise<string>>
type ContentModuleOptions = {
  dir: string
  index: readonly DocRecord[]
}
const stripFile = (record: DocRecord): DocMeta => {
  const { file: _file, ...meta } = record
  return meta
}
const withContent = (record: DocRecord, content: string): ContentDoc => {
  const { file: _file, ...meta } = record
  return { ...meta, content }
}
export const createContentModule = ({ dir, index }: ContentModuleOptions) => {
  const contentCache = new Map<string, string>()
  const findDoc = (slug: string): DocMeta | null => {
    const item = index.find(record => record.slug === slug)
    return item ? stripFile(item) : null
  }
  const docs: DocMeta[] = index.map(stripFile)
  const loadDoc = async (slug: string): Promise<ContentDoc | null> => {
    const item = index.find(record => record.slug === slug)
    if (!item) return null
    const cached = contentCache.get(slug)
    if (cached !== undefined) {
      return withContent(item, cached)
    }
    const loaderKey = `../content/${dir}/${item.file}`
    const loader = CONTENT_LOADERS[loaderKey]
    if (!loader) return null
    const content = String(await loader())
    contentCache.set(slug, content)
    return withContent(item, content)
  }
  return {
    findDoc,
    docs,
    loadDoc
  }
}
