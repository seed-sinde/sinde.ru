<script setup lang="ts">
import { getAuthAvatarUrls } from '~/utils/authAvatar'
import { buildLoginPath } from '~/utils/authNavigation'
const props = withDefaults(defineProps<{
  showLabel?: boolean
  linkClass?: string
  labelClass?: string
}>(), {
  showLabel: true,
  linkClass: '',
  labelClass: 'text-xs',
})
const emit = defineEmits<{
  (e: 'request-close'): void
}>()
const {
  user,
  isAdmin,
  sharedAdminSummary,
  sharedUserSummary,
  loadSharedAdminSummary,
  ensureSummaryRealtime,
} = useAuth()
const route = useRoute()
const { t } = useInterfacePreferences()
const isAuthenticated = computed(() => Boolean(user.value))
const authStatusTo = computed(() => isAuthenticated.value ? '/auth/account' : buildLoginPath(route.fullPath))
const authStatusIcon = computed(() => isAuthenticated.value ? 'ic:round-account-circle' : 'ic:round-login')
const authAvatar = computed(() => getAuthAvatarUrls(user.value))
const isAccountPage = computed(() => route.path.startsWith('/auth/account'))
const authStatusLabel = computed(() => {
  if (!isAuthenticated.value) return t('avatar.login')
  const displayName = String(user.value?.display_name || '').trim()
  if (displayName) return displayName
  const email = String(user.value?.email || '').trim()
  if (email) return email
  return t('avatar.account')
})
const labelClassList = computed(() => ['min-w-0 max-w-[11rem] truncate', props.labelClass])
const userSummaryData = sharedUserSummary
const adminSummaryData = sharedAdminSummary
const hasUserNotifications = computed(() => {
  if (!isAuthenticated.value) return false
  const summary = userSummaryData.value
  if (!summary) return false
  if (typeof summary.has_unread === 'boolean') return summary.has_unread
  return Number(summary.new_approved_recipes_since_last_login || 0) > 0
    || Number(summary.new_rejected_recipes_since_last_login || 0) > 0
})
const hasAdminNotifications = computed(() => {
  if (!isAdmin.value) return false
  const summary = adminSummaryData.value
  if (!summary) return false
  if (typeof summary.has_unread === 'boolean') return summary.has_unread
  return Number(summary.new_users_since_last_login || 0) > 0
    || Number(summary.new_pending_recipes_since_last_login || 0) > 0
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
watch(isAdminPage, async (value) => {
  if (!value || !isAdmin.value) return
  await loadSharedAdminSummary()
})
const onAvatarClick = (event: MouseEvent) => {
  if (!isAuthenticated.value || !isAccountPage.value) return
  event.preventDefault()
  emit('request-close')
}
onMounted(() => {
  if (!import.meta.client) return
  ensureSummaryRealtime()
})
</script>
<template>
  <NuxtLink :to="authStatusTo"
    class="text-(--lab-text-muted) hover:text-(--lab-text-primary) relative inline-flex min-w-0 items-center gap-1.5 transition-colors"
    :class="linkClass" :title="authStatusTitle" :aria-label="authStatusTitle" @click="onAvatarClick">
    <span class="relative inline-flex shrink-0">
      <span v-if="hasAnyNotifications"
        class="ring-(--lab-bg-overlay) absolute -right-0.5 -top-0.5 h-2.5 w-2.5 rounded-full bg-orange-400 ring-2" />
      <img v-if="isAuthenticated && authAvatar.iconImageUrl" :src="authAvatar.iconImageUrl" alt=""
        class="border-(--lab-border) h-6 w-6 shrink-0 rounded-full border object-cover">
      <Icon v-else :name="authStatusIcon" class="h-5 w-5 shrink-0" />
    </span>
    <span v-if="showLabel" :class="labelClassList">{{ authStatusLabel }}</span>
  </NuxtLink>
</template>
