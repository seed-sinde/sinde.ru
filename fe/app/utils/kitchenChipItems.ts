export type KitchenChipItem<T = unknown> = {
  key: string
  label: string
  payload?: T
  disabled?: boolean
}
export const toKitchenChipItems = <T extends { name: string }>(items: T[], prefix = ''): KitchenChipItem<T>[] =>
  items.map((item, index) => ({
    key: `${prefix}${index}:${item.name}`,
    label: item.name,
    payload: item
  }))
export const toKitchenTagChipItems = (items: string[], prefix = ''): KitchenChipItem[] =>
  items.map((item, index) => ({
    key: `${prefix}${index}:${item}`,
    label: item
  }))
