const env = ((globalThis as {process?: {env?: Record<string, string | undefined>}}).process?.env ||
  {}) as Record<string, string | undefined>
const isDev = env.NODE_ENV !== "production"
const runtimeConfig = {
  apiInternalUrl: env.NUXT_API_INTERNAL_URL || env.API_INTERNAL_URL || "",
  apiVersion: env.API_VERSION || "v1",
  public: {
    baseURL: env.NUXT_PUBLIC_BASE_URL || "",
    apiUrl: env.NUXT_PUBLIC_API_URL || "",
    authPasswordMinLength: Number(env.AUTH_PASSWORD_MIN_LENGTH || 12),
    mediaImageMaxBytes: Number(
      env.NUXT_PUBLIC_MEDIA_IMAGE_MAX_BYTES || env.MEDIA_IMAGE_MAX_BYTES || 8388608
    )
  }
}
export {env, isDev, runtimeConfig}
