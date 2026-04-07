<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      version?: 'status' | 'preview' | 'profile'
      user?: AuthUser | null
      showLabel?: boolean
      linkClass?: string
      labelClass?: string
      alt?: string
      clickable?: boolean
    }>(),
    {
      version: 'status',
      showLabel: true,
      linkClass: '',
      labelClass: 'text-xs',
      alt: '',
      clickable: false
    }
  )
  const emit = defineEmits<{
    (e: 'request-close'): void
    (e: 'click'): void
  }>()
  const { user, isAdmin, sharedAdminSummary, sharedUserSummary, loadSharedAdminSummary, ensureSummaryRealtime } =
    useAuth()
  const route = useRoute()
  const { t } = useInterfacePreferences()
  const resolvedUser = computed(() => (props.user === undefined ? user.value : props.user))
  const isStatusVersion = computed(() => props.version === 'status')
  const isProfileVersion = computed(() => props.version === 'profile')
  const isAuthenticated = computed(() => Boolean(resolvedUser.value))
  const authStatusTo = computed(() => (isAuthenticated.value ? '/auth/account' : buildLoginPath(route.fullPath)))
  const authStatusIcon = computed(() => (isAuthenticated.value ? 'ic:round-account-circle' : 'ic:round-login'))
  const avatarUrls = computed(() => getAuthAvatarUrls(resolvedUser.value))
  const avatarImageUrl = computed(() =>
    isProfileVersion.value ? avatarUrls.value.profileImageUrl : avatarUrls.value.iconImageUrl
  )
  const hasAvatarImage = computed(() => Boolean(avatarImageUrl.value))
  const fallbackIconName = computed(() => (isStatusVersion.value ? authStatusIcon.value : 'ic:round-account-circle'))
  const avatarAlt = computed(() => {
    if (props.alt) return props.alt
    return isProfileVersion.value ? 'Аватар профиля' : 'Аватар пользователя'
  })
  const previewRootClass = computed(() =>
    isProfileVersion.value
      ? 'bg-(--lab-bg-surface-subtle) flex h-36 w-36 items-center justify-center overflow-hidden border sm:h-40 sm:w-40'
      : 'relative inline-flex h-6 w-6 shrink-0 items-center justify-center overflow-hidden rounded-full'
  )
  const previewImageClass = computed(() =>
    isProfileVersion.value
      ? 'h-full w-full object-cover'
      : 'h-6 w-6 shrink-0 rounded-full border object-cover'
  )
  const fallbackIconClass = computed(() =>
    isProfileVersion.value ? 'text-(--lab-text-muted) h-20 w-20 text-6xl' : 'text-(--lab-text-muted) h-5 w-5 shrink-0'
  )
  const showPreviewButton = computed(() => props.clickable && hasAvatarImage.value)
  const isAccountPage = computed(() => route.path.startsWith('/auth/account'))
  const authStatusLabel = computed(() => {
    if (!isAuthenticated.value) return t('avatar.login')
    const displayName = String(resolvedUser.value?.display_name || '').trim()
    if (displayName) return displayName
    const email = String(resolvedUser.value?.email || '').trim()
    if (email) return email
    return t('avatar.account')
  })
  const labelClassList = computed(() => ['min-w-0 max-w-[11rem] truncate', props.labelClass])
  const userSummaryData = sharedUserSummary
  const adminSummaryData = sharedAdminSummary
  const hasUserNotifications = computed(() => {
    if (!isStatusVersion.value || !isAuthenticated.value || props.user !== undefined) return false
    const summary = userSummaryData.value
    if (!summary) return false
    if (typeof summary.has_unread === 'boolean') return summary.has_unread
    return (
      Number(summary.new_approved_recipes_since_last_login || 0) > 0 ||
      Number(summary.new_rejected_recipes_since_last_login || 0) > 0
    )
  })
  const hasAdminNotifications = computed(() => {
    if (!isStatusVersion.value || props.user !== undefined || !isAdmin.value) return false
    const summary = adminSummaryData.value
    if (!summary) return false
    if (typeof summary.has_unread === 'boolean') return summary.has_unread
    return (
      Number(summary.new_users_since_last_login || 0) > 0 ||
      Number(summary.new_pending_recipes_since_last_login || 0) > 0
    )
  })
  const hasAnyNotifications = computed(() => hasUserNotifications.value || hasAdminNotifications.value)
  const authStatusTitle = computed(() => {
    const base = authStatusLabel.value
    const parts: string[] = []
    if (hasUserNotifications.value) {
      const approvedNew = Number(userSummaryData.value?.new_approved_recipes_since_last_login || 0)
      const rejectedNew = Number(userSummaryData.value?.new_rejected_recipes_since_last_login || 0)
      parts.push(t('avatar.user_recipes', { approved: approvedNew, rejected: rejectedNew }))
    }
    if (hasAdminNotifications.value) {
      const usersNew = Number(adminSummaryData.value?.new_users_since_last_login || 0)
      const pendingNew = Number(adminSummaryData.value?.new_pending_recipes_since_last_login || 0)
      parts.push(t('avatar.admin_summary', { users: usersNew, pending: pendingNew }))
    }
    return parts.length ? `${base} · ${parts.join(' · ')}` : base
  })
  const isAdminPage = computed(() => route.path.startsWith('/auth/admin'))
  watch(isAdminPage, async value => {
    if (!value || !isAdmin.value || props.user !== undefined) return
    await loadSharedAdminSummary()
  })
  const onAvatarClick = (event: MouseEvent) => {
    if (!isAuthenticated.value || !isAccountPage.value) return
    event.preventDefault()
    emit('request-close')
  }
  const onPreviewClick = () => {
    if (!showPreviewButton.value) return
    emit('click')
  }
  onMounted(() => {
    if (!import.meta.client || props.user !== undefined || !isStatusVersion.value) return
    ensureSummaryRealtime()
  })
</script>
<template>
  <NuxtLink
    v-if="isStatusVersion"
    :to="authStatusTo"
    class="text-(--lab-text-muted) hover:text-(--lab-text-primary) relative inline-flex min-w-0 items-center gap-1.5 transition-colors"
    :class="linkClass"
    :title="authStatusTitle"
    :aria-label="authStatusTitle"
    @click="onAvatarClick">
    <span class="relative inline-flex shrink-0">
      <span
        v-if="hasAnyNotifications"
        class="ring-(--lab-bg-overlay) absolute -right-0.5 -top-0.5 h-2.5 w-2.5 rounded-full bg-orange-400 ring-2" />
      <img v-if="hasAvatarImage" :src="avatarImageUrl" alt="" :class="previewImageClass" />
      <Icon v-else :name="fallbackIconName" :class="fallbackIconClass" />
    </span>
    <span v-if="showLabel" :class="labelClassList">{{ authStatusLabel }}</span>
  </NuxtLink>

  <button
    v-else-if="showPreviewButton"
    type="button"
    :class="previewRootClass"
    :aria-label="avatarAlt"
    @click="onPreviewClick">
    <img :src="avatarImageUrl" :alt="avatarAlt" :class="previewImageClass" />
  </button>

  <div v-else :class="previewRootClass" :aria-label="avatarAlt">
    <img v-if="hasAvatarImage" :src="avatarImageUrl" :alt="avatarAlt" :class="previewImageClass" />
    <Icon v-else :name="fallbackIconName" :class="fallbackIconClass" />
  </div>
</template>
