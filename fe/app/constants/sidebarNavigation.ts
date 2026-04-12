export const SIDEBAR_ITEMS = [
  ['Мастерская', 'ic:twotone-precision-manufacturing', '/lab', 'text-amber-400'],
  ['Особенности', 'ic:twotone-psychology-alt', '/traits', 'text-emerald-400'],
  ['Кухня', 'ic:twotone-soup-kitchen', '/kitchen', 'text-yellow-300'],
  ['Вики', 'ic:twotone-school', '/edu', 'text-violet-400']
] as const
export const sidebarItems = SIDEBAR_ITEMS.map(([label, icon, to, iconColor]) => ({
  label,
  icon,
  to,
  iconColor
}))
