import { buildMediaFileUrl } from '~/utils/mediaUrl'
import type { AuthAvatarGalleryItem, AuthAvatarPayload } from '../../shared/types/authAvatar'
export const AVATAR_ICON_SIZE = 64
export const AVATAR_ICON_DISPLAY_SIZE = 32
export type AuthAvatarKeys = {
  iconImageKey: string
  profileImageKey: string
  originalImageKey: string
}
const isRecord = (value: unknown): value is Record<string, any> =>
  Boolean(value) && typeof value === 'object' && !Array.isArray(value)
const toKey = (value: unknown) => String(value || '').trim()
/** Reads raw avatar data from user profile settings with profile priority. */
export const readAuthAvatarPayload = (user?: AuthUser | null): AuthAvatarPayload => {
  if (!user) return {}
  const profile = isRecord(user.profile) ? user.profile : {}
  const settings = isRecord(user.settings) ? user.settings : {}
  if (isRecord(profile.avatar)) return profile.avatar as AuthAvatarPayload
  if (isRecord(settings.avatar)) return settings.avatar as AuthAvatarPayload
  return {}
}
const normalizeAvatarGalleryItem = (value: unknown, fallbackId: string): AuthAvatarGalleryItem | null => {
  if (!isRecord(value)) return null
  const item: AuthAvatarGalleryItem = {
    id: toKey(value.id) || fallbackId,
    icon_image_key: toKey(value.icon_image_key),
    profile_image_key: toKey(value.profile_image_key),
    original_image_key: toKey(value.original_image_key),
    icon_size: Number.isFinite(Number(value.icon_size)) ? Number(value.icon_size) : null,
    created_at: toKey(value.created_at) || undefined,
    updated_at: toKey(value.updated_at) || undefined
  }
  if (!item.icon_image_key && !item.profile_image_key && !item.original_image_key) return null
  return item
}
/** Returns the normalized avatar gallery and the active primary item id. */
export const getAuthAvatarGallery = (user?: AuthUser | null) => {
  const avatar = readAuthAvatarPayload(user)
  const rawGallery = Array.isArray(avatar.gallery) ? avatar.gallery : []
  const items = rawGallery
    .map((item, index) => normalizeAvatarGalleryItem(item, `avatar-${index + 1}`))
    .filter((item): item is AuthAvatarGalleryItem => Boolean(item))
  if (items.length === 0) {
    const legacyItem = normalizeAvatarGalleryItem(
      {
        id: toKey(avatar.primary_id) || 'avatar-primary',
        icon_image_key: avatar.icon_image_key,
        profile_image_key: avatar.profile_image_key,
        original_image_key: avatar.original_image_key,
        icon_size: avatar.icon_size,
        updated_at: avatar.updated_at
      },
      'avatar-primary'
    )
    if (legacyItem) {
      items.push(legacyItem)
    }
  }
  const preferredId = toKey(avatar.primary_id)
  const primaryId = items.some(item => item.id === preferredId) ? preferredId : items[0]?.id || ''
  return {
    items,
    primaryId
  }
}
/** Builds the stored avatar payload while mirroring the active primary item for backward compatibility. */
export const buildAuthAvatarPayload = (
  items: AuthAvatarGalleryItem[],
  primaryId: string,
  iconSize = AVATAR_ICON_SIZE
): AuthAvatarPayload | null => {
  const normalizedItems = items
    .map((item, index) => normalizeAvatarGalleryItem(item, `avatar-${index + 1}`))
    .filter((item): item is AuthAvatarGalleryItem => Boolean(item))
  if (!normalizedItems.length) return null
  const primaryItem = normalizedItems.find(item => item.id === primaryId) || normalizedItems[0]
  if (!primaryItem) return null
  return {
    primary_id: primaryItem.id,
    gallery: normalizedItems,
    icon_image_key: primaryItem.icon_image_key,
    profile_image_key: primaryItem.profile_image_key,
    original_image_key: primaryItem.original_image_key,
    icon_size: iconSize,
    updated_at: new Date().toISOString()
  }
}
export const getAuthAvatarKeys = (user?: AuthUser | null): AuthAvatarKeys => {
  const gallery = getAuthAvatarGallery(user)
  const avatar = gallery.items.find(item => item.id === gallery.primaryId) || gallery.items[0]
  return {
    iconImageKey: toKey(avatar?.icon_image_key),
    profileImageKey: toKey(avatar?.profile_image_key),
    originalImageKey: toKey(avatar?.original_image_key),
  }
}
export const getAuthAvatarUrls = (user?: AuthUser | null) => {
  const keys = getAuthAvatarKeys(user)
  return {
    ...keys,
    iconImageUrl: buildMediaFileUrl(keys.iconImageKey),
    profileImageUrl: buildMediaFileUrl(keys.profileImageKey || keys.iconImageKey),
    originalImageUrl: buildMediaFileUrl(keys.originalImageKey || keys.profileImageKey || keys.iconImageKey),
  }
}
