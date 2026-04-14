export type MediaUploadTarget =
  | {
      section: 'users'
      collection: 'avatars'
    }
  | {
      section: 'kitchen'
      collection: 'recipes'
      recipeId?: string
    }
const useApiJson = <T>(path: string, options?: NonNullable<Parameters<ReturnType<typeof useAPI>['json']>[1]>) =>
  useAPI().json<T>(path, options)
export const uploadMediaFile = async (file: File, target: MediaUploadTarget) => {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('section', target.section)
  formData.append('collection', target.collection)
  const recipeId =
    target.section === 'kitchen' && target.collection === 'recipes' ? String(target.recipeId || '').trim() : ''
  if (recipeId) {
    formData.append('recipe_id', recipeId)
  }
  return await useApiJson<ApiResponseWithData<MediaUploadResult>>('/media/upload', {
    method: 'POST',
    body: formData
  })
}
