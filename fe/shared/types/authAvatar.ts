export type AuthAvatarGalleryItem = {
  id: string
  icon_image_key: string
  profile_image_key: string
  original_image_key: string
  icon_size?: number | null
  created_at?: string
  updated_at?: string
}
export type AuthAvatarPayload = {
  primary_id?: string
  gallery?: AuthAvatarGalleryItem[]
  icon_image_key?: string
  profile_image_key?: string
  original_image_key?: string
  icon_size?: number | null
  updated_at?: string
}
