const mineralFormulaMarkerPattern = /[_^\\]/
const mineralFormulaChargeScriptPattern = /\^([+-]?\d+[+-]?|[+-]+)/g
export const isMineralFormulaLike = (value: unknown) => {
  const source = String(value || "").trim()
  return Boolean(source) && mineralFormulaMarkerPattern.test(source)
}
export const mineralFormulaToLatex = (value: unknown) => {
  let source = String(value || "").trim()
  if (!source) return ""
  source = source.replace(mineralFormulaChargeScriptPattern, (_, charge: string) => `^{${charge}}`)
  source = source.replace(/\^([^^]+)\^/g, "^{$1}")
  source = source.replace(/_([^_]+)_/g, "_{$1}")
  return source
}
