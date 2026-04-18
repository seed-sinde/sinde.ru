export type {CSSProperties} from "vue"
export type {RouteLocationRaw} from "vue-router"
export type ID = string | number
export type Nullable<T> = T | null
export type Optional<T> = T | undefined
export type Dictionary<T = unknown> = Record<string, T>
export type ClassValue = string | string[] | Record<string, boolean>
export type Timestamp = number
export type WithTimestamps = {
  createdAt?: Timestamp
  updatedAt?: Timestamp
}
export type CopyBlockVariant = "default" | "dark-cyan"
export type MenuItem = {
  to: string
  label: string
  icon: string
  iconColor?: string
  isDev?: boolean
  group?: "tools" | "knowledge"
}
export type ApiResponse = {
  ok: boolean
  message?: string
  details?: string
}
export type ApiResponseWithData<T> = ApiResponse & {
  data: T
}
