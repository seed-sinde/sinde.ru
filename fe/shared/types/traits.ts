export type DataType =
  | "string"
  | "number"
  | "list"
  | "datetime"
  | "range"
  | "interval"
  | "schedule"
  | "geo"
  | "color"
  | "boolean"
  | "validity"
  | "surface"
  | "enum"
  | "datetime-range"
  | "geo-point"
export type TraitColorMode = "hex" | "lab" | "spectrum"
export type TraitMetaUiInput =
  | "text"
  | "number"
  | "array"
  | "color"
  | "boolean"
  | "datetime"
  | "date-range"
  | "geo"
  | "interval"
  | "schedule"
  | "validity"
export type TraitMetaUi = {
  input: TraitMetaUiInput
}
export type TraitBooleanDisplay = "truth" | "answer" | "state" | "presence"
export type TraitListValueType = "string" | "number" | "datetime" | "date" | "time"
export type TraitRangeValueType = "number" | "datetime" | "date" | "time"
export type TraitGeoType = "point" | "polygon" | "zone"
export type TraitUnitCategory =
  | "unitless"
  | "length"
  | "mass"
  | "temperature"
  | "time"
  | "area"
  | "volume"
  | "speed"
  | "energy"
  | "angle"
  | "pressure"
  | "density"
  | "force"
  | "power"
  | "frequency"
  | "electric"
  | "error"
export type TraitUnitOption = {
  value: string
  label: string
}
export type TraitResolvePayload = [string, string] | Trait | Trait[]
export type KeyMeta = {
  dataType: DataType
  ui?: TraitMetaUi
  optionType?: DataType
  listValueType?: TraitListValueType
  rangeType?: TraitRangeValueType
  geoType?: TraitGeoType
  listOrdered?: boolean
  listUnique?: boolean
  minItems?: number
  maxItems?: number
  minLength?: number
  maxLength?: number
  validate?: string
  mode?: TraitColorMode
  format?: string
  options?: string[]
  min?: number
  max?: number
  precision?: number
  unitCategory?: TraitUnitCategory
  unit?: string
  durationUnit?: string
  heightUnit?: string
  radiusUnit?: string
  timezone?: string
  booleanDisplay?: TraitBooleanDisplay
  booleanTrueLabel?: string
  booleanFalseLabel?: string
  colorDomain?: string
  colorStandard?: string
  measurementConditions?: string
  surfaceTexture?: string
  opacity?: number
  namedColor?: string
  palette?: string
  validityMode?: "permanent" | "temporary"
  eventLabel?: string
  glossCategory?: string
  reliefType?: string
  microReliefHeight?: number
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
export type TraitResolveLockOwner = "loader" | "stream"
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
  mode: "permanent" | "temporary"
  since?: string
  until?: string
}
export type Color = {
  mode: TraitColorMode
  hex?: string
  lab?: {L: number; a: number; b: number}
  spectrum?: string
}
export type TraitDateRangeValue = {
  start: string
  end: string
}
export type TraitListModel = {
  items: string[]
}
export type TraitGeoPointModel = {
  lat: string
  lng: string
}
export type TraitGeoPolygonPoint = {
  lat: string
  lng: string
}
export type TraitGeoModel = {
  type: TraitGeoType
  lat?: string
  lng?: string
  height?: string
  heightUnit?: string
  radius?: string
  radiusUnit?: string
  points?: TraitGeoPolygonPoint[]
}
export type TraitRangeModel = {
  start: string
  end: string
}
export type TraitScheduleModel = {
  fromDay: string
  toDay: string
  fromTime: string
  toTime: string
}
export type TraitIntervalUnit = "seconds" | "minutes" | "hours" | "days" | "years"
export type TraitIntervalModel = TraitDateRangeValue & {
  unit: TraitIntervalUnit
}
export type TraitSurfaceModel = {
  glossCategory: string
  glossGU: string
  reliefType: string
  microReliefHeight: string
}
export type TraitValueModelMap = {
  string: string
  number: string | number
  list: TraitListModel
  datetime: string
  range: TraitRangeModel
  interval: TraitIntervalModel
  schedule: TraitScheduleModel
  geo: TraitGeoModel
  color: Color
  boolean: boolean
  validity: Validity
  surface: TraitSurfaceModel
  enum: string
  "datetime-range": TraitDateRangeValue
  "geo-point": TraitGeoPointModel
}
export type TraitValueModel = TraitValueModelMap[DataType]
export type TraitDynamicInputDataType =
  | "string"
  | "number"
  | "boolean"
  | "datetime"
  | "range"
  | "interval"
  | "schedule"
  | "geo"
  | "list"
  | "validity"
  | "color"
  | "surface"
export type TraitDynamicValueModel = TraitValueModelMap[TraitDynamicInputDataType]
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
