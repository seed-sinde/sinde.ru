export type MathUnitCategoryKey =
  | 'length'
  | 'mass'
  | 'temperature'
  | 'area'
  | 'volume'
  | 'speed'
  | 'time'
  | 'pressure'
  | 'energy'
  | 'power'
export type MathUnitBase = {
  key: string
  label: string
  symbol: string
  aliases?: string[]
}
export type MathLinearUnitDef = MathUnitBase & {
  toBase: number
}
export type MathTemperatureUnitDef = MathUnitBase & {
  toBase: (value: number) => number
  fromBase: (value: number) => number
}
export type MathUnitDef = MathLinearUnitDef | MathTemperatureUnitDef
export type MathUnitCategoryBase = {
  key: MathUnitCategoryKey
  title: string
  description: string
  baseUnitKey: string
}
export type MathLinearCategoryDef = MathUnitCategoryBase & {
  kind: 'linear'
  units: MathLinearUnitDef[]
}
export type MathTemperatureCategoryDef = MathUnitCategoryBase & {
  kind: 'temperature'
  units: MathTemperatureUnitDef[]
}
export type MathUnitCategoryDef = MathLinearCategoryDef | MathTemperatureCategoryDef
