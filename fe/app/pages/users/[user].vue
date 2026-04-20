<script setup lang="ts">
definePageMeta({
  path: '/users/:user/:tab(profile|payments|activity)?/:activity(sessions|attempts|events)?',
  validate: route => {
    const tab = Array.isArray(route.params.tab) ? route.params.tab[0] : route.params.tab
    const activity = Array.isArray(route.params.activity) ? route.params.activity[0] : route.params.activity
    return !activity || tab === 'activity'
  }
})

type UserTab = 'profile' | 'payments' | 'activity'
type UserPageData = {
  publicProfile: PublicUserProfileView | null
  detail: AdminUserDetailView | null
  access: PaymentAccessSummary | null
}

const USER_TAB_VALUES: UserTab[] = ['profile', 'payments', 'activity']
const nuxtApp = useNuxtApp()
const route = useRoute()
const router = useRouter()
const readRouteValue = (value: unknown) =>
  typeof value === 'string' ? value : Array.isArray(value) ? String(value[0] || '') : ''
const targetUserId = computed(() => readRouteValue(route.params.user).trim())
const routeUserTab = computed(
  () => normalizeTabRouteValue(readRouteValue(route.params.tab) || readRouteValue(route.query.tab), USER_TAB_VALUES, 'profile') as UserTab
)
const buildUserTabPath = (tab: UserTab) => `/users/${targetUserId.value}/${tab}`
const {
  ensureLoaded,
  isAdmin,
  publicUserProfile,
  adminUserDetail,
  adminBlockUser,
  adminUnblockUser,
  adminForceLogoutUser,
  adminDeleteUser
} = useAuth()
const { adminUserAccess, adminUserOrders } = usePayments()
const { formatAbsoluteDateTime } = useLocalizedDateTime()
const { locale, key, load, t } = useI18nSection('auth')

await nuxtApp.runWithContext(() => useAsyncData(key, load, { watch: [locale] }))

await nuxtApp.runWithContext(() => ensureLoaded())

const adminTab = computed<UserTab>(() => (isAdmin.value ? routeUserTab.value : 'profile'))
const normalizeRoute = async () => {
  if (
    !(
      readRouteValue(route.query.tab) ||
      readRouteValue(route.query.activity) ||
      readRouteValue(route.params.tab) !== adminTab.value ||
      Boolean(readRouteValue(route.params.activity))
    )
  ) {
    return
  }
  const { tab: _tab, activity: _activity, ...query } = route.query
  await nuxtApp.runWithContext(() =>
    navigateTo(
      { path: buildUserTabPath(adminTab.value), query, hash: route.hash },
      { redirectCode: 301, replace: true }
    )
  )
}

await normalizeRoute()
watch([() => route.fullPath, () => isAdmin.value], () => void normalizeRoute())

const pageActionError = ref('')
const pageInfo = ref('')
const emptyUserPage = (): UserPageData => ({
  publicProfile: null,
  detail: null,
  access: null
})
const {
  data: pageData,
  pending: loading,
  error: pageLoadError,
  refresh: refreshPage
} = await nuxtApp.runWithContext(() => useAsyncData(
  computed(() => `user-page:${targetUserId.value}:${isAdmin.value ? 'admin' : 'public'}`),
  async () => {
    if (!targetUserId.value) {
      throw createError({ statusCode: 404, statusMessage: 'Пользователь не найден' })
    }
    const publicRes = await publicUserProfile(targetUserId.value)
    if (!isAdmin.value) {
      return {
        publicProfile: publicRes.data,
        detail: null,
        access: null
      }
    }
    const [detailRes, accessRes] = await Promise.all([adminUserDetail(targetUserId.value), adminUserAccess(targetUserId.value)])
    return {
      publicProfile: publicRes.data,
      detail: detailRes.data,
      access: accessRes.data || null
    }
  },
  {
    default: emptyUserPage
  }
))
const {
  data: ordersData,
  pending: ordersLoading,
  error: ordersLoadError,
  refresh: refreshOrders,
  execute: executeOrders
} = await nuxtApp.runWithContext(() => useAsyncData(
  computed(() => `user-orders:${targetUserId.value}`),
  async () => (!isAdmin.value || adminTab.value !== 'payments' || !targetUserId.value ? [] : (await adminUserOrders(targetUserId.value)).data.items || []),
  {
    default: () => [],
    immediate: false
  }
))

watch(
  [() => isAdmin.value, adminTab, targetUserId],
  ([admin, tab, userId], prev) => {
    if (!admin || tab !== 'payments' || !userId) return
    if (!prev || admin !== prev[0] || tab !== prev[1] || userId !== prev[2] || !ordersData.value.length) {
      void executeOrders().catch(() => {})
    }
  },
  { immediate: true }
)

watch(targetUserId, () => {
  pageActionError.value = ''
  pageInfo.value = ''
})

const publicProfile = computed(() => pageData.value.publicProfile)
const detail = computed(() => pageData.value.detail)
const access = computed(() => pageData.value.access)
const orders = computed(() => ordersData.value)
const pageError = computed(() =>
  pageActionError.value ||
  extractApiErrorMessage(pageLoadError.value, '') ||
  (adminTab.value === 'payments' ? extractApiErrorMessage(ordersLoadError.value, '') : '')
)

const title = computed(() => {
  const displayName =
    String(detail.value?.user.display_name || '').trim() || String(publicProfile.value?.display_name || '').trim()
  return displayName || 'Пользователь'
})
const breadcrumbItems = computed<BreadcrumbItem[]>(() => {
  const root = { label: title.value, to: buildUserTabPath('profile') }
  if (adminTab.value === 'activity') {
    return [root, { label: t('account.activity.title'), current: true, kind: 'tab' }]
  }
  return [
    root,
    {
      label: adminTabItems.value.find(item => item.value === adminTab.value)?.label || t('account.profile'),
      current: true,
      kind: 'tab'
    }
  ]
})
nuxtApp.runWithContext(() => {
  usePageSeo({
    title,
    description: 'Профиль пользователя'
  })
})
const adminTabItems = computed<LabTabItem[]>(() => [
  { value: 'profile', label: t('account.profile') },
  { value: 'payments', label: t('account.payments') },
  { value: 'activity', label: t('account.activity.title') }
])

const adminTabRouteTargetMap = computed<TabRouteTargetMap>(() => ({
  profile: buildUserTabPath('profile'),
  payments: buildUserTabPath('payments'),
  activity: buildUserTabPath('activity')
}))

const displayUser = computed<AuthUser | null>(() => {
  if (detail.value?.user) return detail.value.user
  if (!publicProfile.value) return null
  return {
    user_id: publicProfile.value.user_id,
    email: '',
    status: 'active',
    display_name: publicProfile.value.display_name,
    locale: 'ru-RU',
    timezone: 'Europe/Moscow',
    roles: [],
    is_two_factor_enabled: false,
    primary_trait_uuid: publicProfile.value.primary_trait_uuid ?? null,
    created_at: new Date(0).toISOString(),
    profile: publicProfile.value.profile || {},
    settings: {}
  }
})

const paymentAccessUntilText = computed(() => {
  if (!access.value?.access_until) return ''
  return formatAbsoluteDateTime(access.value.access_until)
})
const subscriptionTooltipText = computed(() => {
  if (!access.value?.has_active_access || !paymentAccessUntilText.value) return ''
  return `Подписка активна до ${paymentAccessUntilText.value}`
})
const adminDetailUser = computed(() => detail.value?.user || null)
const targetUserIsAdmin = computed(() => adminDetailUser.value?.roles.includes('admin') === true)
const showAdminView = computed(() => isAdmin.value && Boolean(detail.value))
const adminActivitySessions = computed(() => detail.value?.sessions || [])
const adminActivityLoginAttempts = computed(() => detail.value?.login_attempts || [])
const adminActivitySecurityEvents = computed(() => detail.value?.security_events || [])
const paymentHistoryColumns = computed<LabDataTableColumn[]>(() => [
  { key: 'createdAt', label: t('account.activity.created_at'), nowrap: true },
  { key: 'plan', label: t('account.activity.plan') },
  { key: 'amount', label: t('account.activity.amount'), nowrap: true },
  { key: 'status', label: t('account.activity.status'), nowrap: true },
  { key: 'access', label: t('account.activity.access') }
])
const paymentHistoryRows = computed<{
  id: string
  createdAt: string
  plan: string
  amount: string
  status: string
  access: string
}[]>(() =>
  orders.value.map(item => ({
    id: item.order_id,
    createdAt: formatDateTime(item.created_at),
    plan: paymentPlanLabel(item.plan_code),
    amount: formatPaymentWholeRubles(item.amount),
    status: paymentStatusLabel(item.status),
    access: item.access_until
      ? `До ${formatDateTime(item.access_until)}`
      : item.access_from
        ? `С ${formatDateTime(item.access_from)}`
        : '—'
  }))
)
const formatDateTime = (value?: string | null) =>
  formatAbsoluteDateTime(value, { dateStyle: 'medium', timeStyle: 'short' })
const paymentPlanLabel = (planCode?: string | null) => (String(planCode || '').trim() === 'donation' ? 'Донат' : 'Pro')
const paymentStatusLabel = (status?: string | null) => {
  switch (String(status || '').trim()) {
    case 'success':
      return 'Успешно'
    case 'pending':
      return 'В обработке'
    case 'failed':
      return 'Ошибка'
    case 'canceled':
      return 'Отменён'
    case 'refunded':
      return 'Возврат'
    default:
      return '—'
  }
}
const refreshAll = async () => {
  await refreshPage()
  if (isAdmin.value && adminTab.value === 'payments') await refreshOrders()
}
const blockUser = async () => {
  if (!adminDetailUser.value) return
  if (targetUserIsAdmin.value) {
    pageActionError.value = 'Нельзя блокировать пользователя с ролью admin.'
    return
  }
  pageActionError.value = ''
  pageInfo.value = ''
  try {
    await adminBlockUser(adminDetailUser.value.user_id, '')
    pageInfo.value = `Пользователь ${adminDetailUser.value.email} заблокирован.`
    await refreshAll()
  } catch (err: unknown) {
    pageActionError.value = extractApiErrorMessage(err, 'Не удалось заблокировать пользователя.')
  }
}
const unblockUser = async () => {
  if (!adminDetailUser.value) return
  pageActionError.value = ''
  pageInfo.value = ''
  try {
    await adminUnblockUser(adminDetailUser.value.user_id)
    pageInfo.value = `Пользователь ${adminDetailUser.value.email} разблокирован.`
    await refreshAll()
  } catch (err: unknown) {
    pageActionError.value = extractApiErrorMessage(err, 'Не удалось разблокировать пользователя.')
  }
}
const forceLogout = async () => {
  if (!adminDetailUser.value) return
  if (targetUserIsAdmin.value) {
    pageActionError.value = 'Нельзя сбрасывать сессии пользователя с ролью admin.'
    return
  }
  pageActionError.value = ''
  pageInfo.value = ''
  try {
    await adminForceLogoutUser(adminDetailUser.value.user_id)
    pageInfo.value = `Все сессии пользователя ${adminDetailUser.value.email} завершены.`
    await refreshAll()
  } catch (err: unknown) {
    pageActionError.value = extractApiErrorMessage(err, 'Не удалось завершить сессии пользователя.')
  }
}
const deleteUser = async () => {
  if (!adminDetailUser.value) return
  if (targetUserIsAdmin.value) {
    pageActionError.value = 'Нельзя удалять пользователя с ролью admin.'
    return
  }
  pageActionError.value = ''
  pageInfo.value = ''
  try {
    await adminDeleteUser(adminDetailUser.value.user_id)
    await router.push('/auth/admin/users')
  } catch (err: unknown) {
    pageActionError.value = extractApiErrorMessage(err, 'Не удалось удалить пользователя.')
  }
}
</script>

<template>
  <div>
    <LabNavHeader :title :breadcrumb-items="breadcrumbItems" />
    <LabLoader v-if="loading" variant="inline" label="Загрузка пользователя…" />
    <template v-else-if="displayUser">
      <LabNotify v-if="pageError" :text="pageError" tone="error" size="xs" />
      <LabNotify v-if="pageInfo" :text="pageInfo" tone="success" size="xs" />
      <section v-if="showAdminView && adminDetailUser">
        <LabNavTabs
          :model-value="adminTab"
          :items="adminTabItems"
          route-param-key="tab"
          :route-target-map="adminTabRouteTargetMap"
        >
          <template #panel-profile>
            <section class="p-4">
              <div class="flex flex-wrap items-start gap-4">
                <LabAvatar version="profile" :user="adminDetailUser" />
                <div class="flex min-w-0 flex-1 flex-col gap-3">
                  <div class="flex flex-wrap items-start gap-2">
                    <h2 class="text-xl font-semibold">{{ adminDetailUser.display_name || 'Пользователь' }}</h2>
                    <LabBaseTooltip
                      v-if="subscriptionTooltipText"
                      :text="subscriptionTooltipText"
                      side="right"
                      align="left"
                      :offset="10"
                      :cross-axis-offset="0"
                    >
                      <template #trigger>
                        <LabBaseButton
                          icon="ic:round-auto-awesome"
                          icon-only
                          variant="ghost"
                          size="sm"
                          class="h-8 w-8 rounded-full border-transparent text-orange-300 hover:bg-(--lab-bg-surface-hover) focus:bg-(--lab-bg-surface-hover) focus-visible:bg-(--lab-bg-surface-hover)"
                          aria-label="Статус подписки"
                        />
                      </template>
                    </LabBaseTooltip>
                  </div>
                  <div class="flex flex-wrap items-baseline gap-x-6 gap-y-2 text-sm">
                    <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                      <span class="shrink-0 text-xs tracking-wide text-(--lab-text-muted) uppercase">Email</span>
                      <NuxtLink :to="`mailto:${adminDetailUser.email}`" external>{{ adminDetailUser.email }}</NuxtLink>
                    </div>
                    <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                      <span class="shrink-0 text-xs tracking-wide text-(--lab-text-muted) uppercase">
                        Последний вход
                      </span>
                      <span>{{ formatDateTime(adminDetailUser.last_login_at) }}</span>
                    </div>
                    <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                      <span class="shrink-0 text-xs tracking-wide text-(--lab-text-muted) uppercase">Регистрация</span>
                      <span>{{ formatDateTime(adminDetailUser.created_at) }}</span>
                    </div>
                  </div>
                  <div class="flex flex-nowrap items-center gap-1 overflow-x-auto whitespace-nowrap">
                    <span v-if="adminDetailUser.status !== 'blocked' && !targetUserIsAdmin" class="max-sm:hidden">
                      <LabConfirmActionButton
                        label="Блок"
                        confirm-label="Подтвердить"
                        tooltip="Подтвердить блокировку аккаунта?"
                        icon="ic:round-block"
                        @confirm="blockUser"
                      />
                    </span>
                    <span v-if="adminDetailUser.status !== 'blocked' && !targetUserIsAdmin" class="sm:hidden">
                      <LabConfirmActionButton
                        icon-only
                        icon="ic:round-block"
                        aria-label="Заблокировать"
                        confirm-aria-label="Подтвердить блокировку"
                        confirm-label="Ок"
                        tooltip="Подтвердить блокировку аккаунта?"
                        @confirm="blockUser"
                      />
                    </span>
                    <LabBaseButton
                      v-if="adminDetailUser.status === 'blocked'"
                      class="max-sm:hidden"
                      variant="secondary"
                      size="xs"
                      icon="ic:round-lock-open"
                      label="Разблок"
                      @click="unblockUser"
                    />
                    <LabBaseButton
                      v-if="adminDetailUser.status === 'blocked'"
                      class="sm:hidden"
                      variant="secondary"
                      size="xs"
                      icon="ic:round-lock-open"
                      icon-only
                      aria-label="Разблокировать"
                      @click="unblockUser"
                    />
                    <span v-if="!targetUserIsAdmin" class="max-sm:hidden">
                      <LabConfirmActionButton
                        label="Сброс сессий"
                        confirm-label="Подтвердить"
                        tooltip="Подтвердить принудительный сброс всех сессий?"
                        icon="ic:round-logout"
                        @confirm="forceLogout"
                      />
                    </span>
                    <span v-if="!targetUserIsAdmin" class="sm:hidden">
                      <LabConfirmActionButton
                        icon-only
                        icon="ic:round-logout"
                        aria-label="Сбросить сессии"
                        confirm-aria-label="Подтвердить сброс сессий"
                        confirm-label="Ок"
                        tooltip="Подтвердить принудительный сброс всех сессий?"
                        @confirm="forceLogout"
                      />
                    </span>
                    <span v-if="!targetUserIsAdmin" class="max-sm:hidden">
                      <LabConfirmActionButton
                        label="Удалить"
                        confirm-label="Подтвердить"
                        tooltip="Подтвердить удаление аккаунта?"
                        icon="ic:round-delete"
                        @confirm="deleteUser"
                      />
                    </span>
                    <span v-if="!targetUserIsAdmin" class="sm:hidden">
                      <LabConfirmActionButton
                        icon-only
                        icon="ic:round-delete"
                        aria-label="Удалить аккаунт"
                        confirm-aria-label="Подтвердить удаление аккаунта"
                        confirm-label="Ок"
                        tooltip="Подтвердить удаление аккаунта?"
                        @confirm="deleteUser"
                      />
                    </span>
                  </div>
                </div>
              </div>
            </section>
          </template>
          <template #panel-payments>
            <LabDataTable
              :columns="paymentHistoryColumns"
              :rows="paymentHistoryRows"
              :loading="ordersLoading"
              empty-text="История транзакций пуста."
            />
          </template>
          <template #panel-activity>
            <div class="space-y-4 p-4">
              <AuthAccountActivityTables
                :sessions="adminActivitySessions"
                :login-attempts="adminActivityLoginAttempts"
                :security-events="adminActivitySecurityEvents"
              />
            </div>
          </template>
        </LabNavTabs>
      </section>
      <section v-else class="space-y-4 p-4">
        <div class="flex flex-wrap items-start gap-4">
          <LabAvatar version="profile" :user="displayUser" />
          <div class="space-y-2">
            <h1 class="text-2xl font-semibold">{{ displayUser.display_name || 'Пользователь' }}</h1>
          </div>
        </div>
      </section>
    </template>
    <LabNotify v-else :text="pageError || 'Пользователь не найден.'" tone="error" size="xs" />
  </div>
</template>
