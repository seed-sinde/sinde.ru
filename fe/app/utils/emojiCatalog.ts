type EmojiIconData = {
  body: string
  left?: number
  top?: number
  width?: number
  height?: number
}
type EmojiAliasData = {
  parent: string
}
export type EmojiCollectionData = {
  icons?: Record<string, EmojiIconData>
  aliases?: Record<string, EmojiAliasData>
  width?: number
  height?: number
}
// Static grep across fe/app currently finds no direct emojione:* UI icons outside markdown/editor flows.
export const UI_EMOJI_NAMES: string[] = []
export const EMOJI_PICKER_FALLBACK = [
  'grinning-face',
  'beaming-face-with-smiling-eyes',
  'face-with-tears-of-joy',
  'winking-face',
  'thinking-face',
  'face-with-rolling-eyes',
  'face-with-medical-mask',
  'smiling-face',
  'smiling-face-with-heart-eyes',
  'smiling-face-with-sunglasses',
  'hugging-face',
  'loudly-crying-face',
]
let emojiCollectionPromise: Promise<EmojiCollectionData> | null = null
let emojiPickerNamesPromise: Promise<string[]> | null = null
export const loadEmojiCollection = async (): Promise<EmojiCollectionData> => {
  if (!emojiCollectionPromise) {
    emojiCollectionPromise = import('~/data/emojione-face-subset.json')
      .then((mod: any) => (mod?.default || mod) as EmojiCollectionData)
  }
  return await emojiCollectionPromise
}
export const loadEmojiPickerNames = async () => {
  if (!emojiPickerNamesPromise) {
    emojiPickerNamesPromise = loadEmojiCollection()
      .then(collection => Object.keys(collection.icons || {}).sort())
  }
  return await emojiPickerNamesPromise
}
