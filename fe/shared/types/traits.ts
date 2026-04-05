export type DataType =
  | 'string'
  | 'enum'
  | 'number'
  | 'datetime'
  | 'datetime-range'
  | 'interval'
  | 'schedule'
  | 'geo-point'
  | 'color'
  | 'boolean'
  | 'validity'
export type TraitColorMode = 'hex' | 'lab' | 'spectrum'
export type TraitUnitCategory =
  | 'unitless'
  | 'length'
  | 'mass'
  | 'temperature'
  | 'time'
  | 'area'
  | 'volume'
  | 'speed'
  | 'energy'
export type TraitUnitOption = {
  value: string
  label: string
}
export type TraitResolvePayload = [string, string] | Trait | Trait[]
export type KeyMeta = {
  dataType: DataType
  optionType?: DataType
  mode?: TraitColorMode
  format?: string
  options?: string[]
  min?: number
  max?: number
  precision?: number
  unitCategory?: TraitUnitCategory
  unit?: string
}
export type TraitKey = {
  id: number
  syn: string
  syn_id?: number
  meta: KeyMeta | Record<string, unknown>
}
export type TraitInput = {
  t_key: string
  t_value: string
  meta?: KeyMeta
}
export type Trait = TraitInput & {
  t_uuid: string
  t_key_id?: number
}
export type TraitResolveLockOwner = 'loader' | 'stream'
export type TraitResolveLockState = {
  uuid: string | null
  owner: TraitResolveLockOwner | null
}
export type TraitAddFormFocusSnapshot = {
  id: string
  name: string
  wasValueField: boolean
  selectionStart: number | null
  selectionEnd: number | null
}
export type DateTextMaskSpec = {
  template: string
  slots: number[]
}
export type Validity = {
  mode: 'permanent' | 'temporary'
  since?: string
  until?: string
}
export type Color = {
  mode: TraitColorMode
  hex?: string
  lab?: { L: number; a: number; b: number }
  spectrum?: string
}
export type TraitDateRangeValue = {
  start: string
  end: string
}
export type TraitGeoPointModel = {
  lat: string
  lng: string
}
export type TraitScheduleModel = {
  fromDay: string
  toDay: string
  fromTime: string
  toTime: string
}
export type TraitIntervalUnit = 'seconds' | 'minutes' | 'hours' | 'days' | 'years'
export type TraitIntervalModel = TraitDateRangeValue & {
  unit: TraitIntervalUnit
}
export type TraitColorPreviewInfo = {
  css: string
  text: string
  mode: TraitColorMode
}
export type GeoCoordinateParseResult = {
  trimmed: string
  num: number
  decimals: number
}
