export const SIDEBAR_ITEMS = [
  ['Мастерская', 'ic:twotone-precision-manufacturing', '/lab', 'text-amber-400', false, 'tools'],
  ['Особенности', 'ic:twotone-psychology-alt', '/traits', 'text-emerald-400', false, 'tools'],
  ['Кухня', 'ic:twotone-soup-kitchen', '/kitchen', 'text-yellow-300', false, 'tools'],
  ['Вики', 'ic:twotone-school', '/edu', 'text-violet-400', false, 'knowledge'],
  ['Разработка', 'ic:twotone-auto-graph', '/dev', 'text-cyan-400', false, 'knowledge'],
  ['Влияние', 'ic:twotone-eco', '/impact', 'text-lime-400', false, 'knowledge'],
  ['Восприятие', 'ic:twotone-self-improvement', '/perception', 'text-fuchsia-400', true, undefined],
  ['Опыт', 'ic:twotone-work-history', '/experience', 'text-sky-400', true, undefined],
  ['Карьера', 'ic:twotone-work', '/careers', 'text-orange-400', true, undefined],
  ['Семья', 'ic:twotone-family-restroom', '/family', 'text-pink-400', true, undefined]
] as const
export const sidebarItems = SIDEBAR_ITEMS.map(([label, icon, to, iconColor, isDev, group]) => ({
  label,
  icon,
  to,
  iconColor,
  ...(isDev ? { isDev } : {}),
  ...(group ? { group } : {})
}))
