import type {DocSectionItem, DocsSlug} from "../types/docs"
import type {LabTabItem} from "../types/ui"

export const DOCS_SLUGS = ["company", "offer", "terms", "privacy"] as const satisfies ReadonlyArray<DocsSlug>

export const DOCS_DEFAULT_SLUG: DocsSlug = "company"

type DocsItemSchema = {
  key: string
  href?: string
}

type DocsSchema = {
  introItems: readonly DocsItemSchema[]
  sections: readonly DocsItemSchema[]
}

export const DOCS_SCHEMA: Record<DocsSlug, DocsSchema> = {
  company: {
    introItems: [
      {key: "provider"},
      {key: "email", href: "company.intro.items.email.href"},
      {key: "phone", href: "company.intro.items.phone.href"},
      {key: "response_time"}
    ],
    sections: [
      {key: "registered_name"},
      {key: "registered_address"},
      {key: "tin"},
      {key: "primary_registration_number"},
      {key: "settlement_account"},
      {key: "bank"},
      {key: "bic"},
      {key: "correspondent_account"}
    ]
  },
  offer: {
    introItems: [
      {key: "format"},
      {key: "access_period"},
      {key: "refunds"},
      {key: "personal_data", href: "/docs/privacy"}
    ],
    sections: [
      {key: "1_general"},
      {key: "2_subject_matter"},
      {key: "3_price_and_payment"},
      {key: "4_delivery_of_access"},
      {key: "5_access_term_and_termination"},
      {key: "6_refunds"},
      {key: "7_payment_documents"},
      {key: "8_related_documents"}
    ]
  },
  terms: {
    introItems: [{key: "purpose"}, {key: "related_documents"}],
    sections: [
      {key: "1_access_and_account"},
      {key: "2_acceptable_use"},
      {key: "3_user_data"},
      {key: "4_access_restrictions"},
      {key: "5_intellectual_property"},
      {key: "6_security"},
      {key: "7_changes_and_communications"}
    ]
  },
  privacy: {
    introItems: [
      {key: "operator"},
      {key: "contacts", href: "privacy.intro.items.contacts.href"},
      {key: "payment_document", href: "/docs/offer"}
    ],
    sections: [
      {key: "1_data_processed"},
      {key: "2_purposes"},
      {key: "3_legal_grounds"},
      {key: "4_payment_data"},
      {key: "5_storage_and_transfer"},
      {key: "6_user_rights"},
      {key: "7_security_and_requests"}
    ]
  }
}

export const normalizeDocsSlug = (value: unknown): DocsSlug | null => {
  if (typeof value === "string" && DOCS_SLUGS.includes(value as DocsSlug)) {
    return value as DocsSlug
  }
  if (Array.isArray(value) && typeof value[0] === "string" && DOCS_SLUGS.includes(value[0] as DocsSlug)) {
    return value[0] as DocsSlug
  }
  return null
}

export const isDocsSlug = (value: unknown): value is DocsSlug => normalizeDocsSlug(value) !== null

export const getDocsHref = (slug: DocsSlug): string => `/docs/${slug}`

export const getDocsTabItems = (titles: Record<DocsSlug, string>): LabTabItem[] =>
  DOCS_SLUGS.map(
    (slug): LabTabItem => ({
      value: slug,
      label: titles[slug]
    })
  )

export const getDocsTabRouteMap = (): Record<string, string> =>
  Object.fromEntries(DOCS_SLUGS.map(slug => [slug, getDocsHref(slug)]))

export const getDocsI18nKey = (
  slug: DocsSlug,
  section: "intro" | "sections",
  item: DocsItemSchema,
  field: "label" | "value"
) => `${slug}.${section}.items.${item.key}.${field}`

export const mapDocsItems = (
  slug: DocsSlug,
  section: "intro" | "sections",
  items: readonly DocsItemSchema[],
  resolve: (key: string) => string
): DocSectionItem[] =>
  items.map(item => {
    const href = !item.href ? "" : item.href.startsWith("/") ? item.href : resolve(item.href)
    const labelKey = getDocsI18nKey(slug, section, item, "label")
    const valueKey = getDocsI18nKey(slug, section, item, "value")
    return {
      label: resolve(labelKey),
      value: resolve(valueKey),
      ...(href ? {href} : {})
    }
  })
