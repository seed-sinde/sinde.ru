import {slugifyLatin} from "~/utils/slug"
export const extractKitchenRecipeId = (value: string) => {
  const raw = String(value || "").trim()
  if (!raw) return ""
  const uuidMatch = raw.match(/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/i)
  return uuidMatch ? uuidMatch[0] : ""
}
export const extractKitchenEditSlug = (path: string) => {
  const rawPath = String(path || "").trim()
  const match = rawPath.match(/^\/kitchen\/edit\/(.+)$/)
  return match ? decodeURIComponent(match[1] || "").trim() : ""
}
export const kitchenStepWord = (count: number) => {
  const abs = Math.abs(count) % 100
  const last = abs % 10
  if (abs > 10 && abs < 20) return "шагов"
  if (last === 1) return "шаг"
  if (last >= 2 && last <= 4) return "шага"
  return "шагов"
}
export const kitchenStepsCountLabel = (count: number) => `${count} ${kitchenStepWord(count)}`
export const slugifyKitchenRecipeTitle = (title: string) => slugifyLatin(title, {fallback: "recipe"})
export const kitchenRecipeSlug = (recipe: {id: string; title: string}) =>
  `${slugifyKitchenRecipeTitle(recipe.title)}-${recipe.id}`
export const kitchenRecipeModerationLabel = (status?: string) => {
  const key = String(status || "").trim()
  if (key === "pending") return "На модерации"
  if (key === "approved") return "Публичный"
  if (key === "rejected") return "Отклонён"
  return "Черновик"
}
