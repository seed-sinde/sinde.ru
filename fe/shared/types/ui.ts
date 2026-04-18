import type {HTMLAttributes, StyleValue} from "vue"
export type TableColumn = {
  key: string
  label: string
  widthClass?: string
  headerClass?: string
  cellClass?: string
  nowrap?: boolean
}
export type TableRow = Record<string, unknown>
export type TableConfig = {
  columns: TableColumn[]
  rows: TableRow[]
  rowKey: string
}
export type UiSectionTab = "forms" | "actions" | "feedback" | "data" | "media"
export type DifficultyScaleOption = {
  value: string
  label?: string
  activeColor?: string
}
export type CropperAspectPreset = "1:1" | "4:3" | "16:9" | "custom" | "free"
export type NotifyTone = "error" | "success" | "info" | "warning"
export type NotifySize = "xs" | "sm" | "base"
export type NotifyTag = "p" | "div" | "span"
export type LabButtonVariant = "default" | "plain" | "primary" | "secondary" | "danger" | "success" | "info" | "ghost"
export type LabButtonSize = "none" | "xs" | "sm" | "md" | "lg" | "xl" | "2xl" | "3xl"
export type LabButtonType = "button" | "submit" | "reset"
export type LabButtonClass = HTMLAttributes["class"]
export type LabButtonStyle = StyleValue | Record<string, string | number | undefined>
export type ToggleTone = "cyan" | "emerald" | "rose" | "amber"
export type ToggleVisualState = "off" | "partial" | "on"
export type SelectOptionInput = {
  value: string | number | null | undefined
  label: string
  disabled?: boolean
  swatchColor?: string
}
export type NormalizedSelectOption = {
  key: string
  value: string
  label: string
  disabled: boolean
  swatchColor?: string
}
export type LabDataTableColumn = TableColumn
export type LabDataTableRow = TableRow
export type ViewerMode = "fit" | "original"
export type ImageViewerItem = {
  src: string
  title?: string | null
  author?: string | null
  attribution?: string | null
  sourceUrl?: string | null
  licenseUrl?: string | null
  license?: string | null
  alt?: string | null
  thumbnailSrc?: string | null
}
export type BreadcrumbBase = {
  label: string
  to?: string
  current?: boolean
  kind?: "page" | "tab"
  badge?: string
}
export type BreadcrumbItem = BreadcrumbBase
export type NormalizedBreadcrumbItem = BreadcrumbBase & Required<Pick<BreadcrumbBase, "current" | "kind">>
export type LabTabValue = string | number
export type LabTabItem = {
  value: LabTabValue
  label: string
  disabled?: boolean
  badge?: string | number
}
export type LanguageOption = {
  key: string
  label: string
  category: string
}
export type InterfaceLocaleCode = "ru" | "en" | "zh" | "ja"
export type ThemePreference = "system" | "dark" | "light"
