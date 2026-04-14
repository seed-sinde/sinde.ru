import type { InterfaceLocaleCode } from './ui'
export type LocalizedRecord<T> = Record<InterfaceLocaleCode, T>

export type LocalizedLinkItem = {
  label: string
  value: string
  href?: string
}

export type LocalizedTextItem = {
  label: string
  value: string
}
