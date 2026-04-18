import type {InterfaceLocaleCode} from "./ui"

export type LocalizedRecord<T> = Record<InterfaceLocaleCode, T>
export type I18nNamespace = "auth" | "common" | "docs" | "minerals" | "nav" | "payments" | "traits" | "ui"
export type I18nMessages = Record<string, string>

export type LocalizedLinkItem = {
  label: string
  value: string
  href?: string
}

export type LocalizedTextItem = {
  label: string
  value: string
}
