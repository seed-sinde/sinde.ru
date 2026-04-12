import { buildMediaFileUrl } from '../app/utils/mediaUrl'
import { slugifyLatin } from '../app/utils/slug'

Object.assign(globalThis, {
  buildMediaFileUrl,
  slugifyLatin
})

const { formatPeriodicTableSymbol } = await import('../app/utils/chemistryElements')

Object.assign(globalThis, {
  formatPeriodicTableSymbol
})
