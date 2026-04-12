export const SIDEBAR_ITEMS = [
  ['Мастерская', 'ic:twotone-precision-manufacturing', '/lab', 'text-amber-400', false],
  ['Особенности', 'ic:twotone-psychology-alt', '/traits', 'text-emerald-400', false],
  ['Кухня', 'ic:twotone-soup-kitchen', '/kitchen', 'text-yellow-300', false],
  ['Вики', 'ic:twotone-school', '/edu', 'text-violet-400', false],
  ['Восприятие', 'ic:twotone-self-improvement', '/perception', 'text-fuchsia-400', true],
  ['Опыт', 'ic:twotone-work-history', '/experience', 'text-sky-400', true],
  ['Карьера', 'ic:twotone-work', '/careers', 'text-orange-400', true],
  ['Семья', 'ic:twotone-family-restroom', '/family', 'text-pink-400', true]
] as const
export const sidebarItems = SIDEBAR_ITEMS.map(([label, icon, to, iconColor, isDev]) => ({
  label,
  icon,
  to,
  iconColor,
  ...(isDev ? { isDev } : {})
}))
