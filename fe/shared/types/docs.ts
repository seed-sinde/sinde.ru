export type DocsSlug = 'company' | 'offer' | 'terms' | 'privacy'

export type DocSectionItem = {
  label: string
  value: string
  href?: string
}

export type DocContent = {
  title: string
  description: string
  introTitle?: string
  introItems?: DocSectionItem[]
  sectionsTitle?: string
  sections: DocSectionItem[]
}

export type DocsPageContent = {
  tabsTitle: string
  documents: Record<DocsSlug, DocContent>
}
