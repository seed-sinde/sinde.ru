<script setup lang="ts">
type AdminUserRow = LabDataTableRow & AdminUserView
type ModerationRow = LabDataTableRow & KitchenRecipe
type AdminTraitKeyRow = LabDataTableRow & AdminTraitKeySearchItem

definePageMeta({
  path: '/auth/admin/:tab(users|moderation|keys|analysis)?',
  middleware: ['admin-only']
})
const title = 'Админка'
usePageSeo({
  title,
  description: 'Управление пользователями, ролями и базовая аналитика.'
})
const {
  user,
  adminListUsers,
  sharedAdminSummary,
  loadSharedAdminSummary,
  adminMarkSummaryRead,
  adminSetUserRole,
  adminBlockUser,
  adminUnblockUser,
  adminForceLogoutUser,
  adminDeleteUser,
  adminSearchKeys,
  adminTraitsSetsAnalysis,
  adminListKitchenModerationRecipes,
  adminModerateKitchenRecipe,
  adminChangeKitchenRecipeOwner
} = useAuth()
const { adminSummary: loadPaymentsAdminSummary } = usePayments()
const { formatAbsoluteDateTime } = useLocalizedDateTime()
const { localeTag } = useInterfacePreferences()
const ADMIN_TABS = ['users', 'moderation', 'keys', 'analysis'] as const
const route = useRoute()
const router = useRouter()
const adminTabPathMap: Record<AdminTab, string> = {
  users: '/auth/admin/users',
  moderation: '/auth/admin/moderation',
  keys: '/auth/admin/keys',
  analysis: '/auth/admin/analysis'
}
const routeAdminTabParam = Array.isArray(route.params.tab) ? route.params.tab[0] : route.params.tab
const routeAdminTabQuery = Array.isArray(route.query.tab) ? route.query.tab[0] : route.query.tab
const routeAdminTab = normalizeAdminTab(routeAdminTabParam || routeAdminTabQuery)
if (routeAdminTabQuery || routeAdminTabParam !== routeAdminTab) {
  const { tab: _tab, ...query } = route.query
  await navigateTo(
    {
      path: adminTabPathMap[routeAdminTab],
      query,
      hash: route.hash
    },
    {
      redirectCode: 301,
      replace: true
    }
  )
}

const users = ref<AdminUserView[]>([])
const usersTotal = ref(0)
const usersLimit = ref(50)
const usersOffset = ref(0)
const usersInfo = ref('')
const usersError = ref('')
const userSearch = ref('')
const userStatusFilter = ref('')
const userRoleFilter = ref('')
const userStatusOptions: SelectOptionInput[] = [
  { value: '', label: 'все' },
  { value: 'active', label: 'Активен' },
  { value: 'pending_verification', label: 'Ожидает подтверждения' },
  { value: 'blocked', label: 'Заблокирован' }
]
const userRoleOptions: SelectOptionInput[] = [
  { value: '', label: 'все' },
  { value: 'admin', label: 'admin' },
  { value: 'user', label: 'user' }
]
const editableRoleOptions: SelectOptionInput[] = [
  { value: 'user', label: 'user' },
  { value: 'admin', label: 'admin' }
]
const roleDrafts = reactive<Record<string, 'admin' | 'user'>>({})
const roleSavingUserId = ref<string | null>(null)
const summaryInfo = ref('')
const summaryActionError = ref('')
const summary = computed(() => sharedAdminSummary.value)
const summaryReadPending = ref(false)
const adminTab = computed<AdminTab>({
  get: () => normalizeAdminTab(route.params.tab || route.query.tab),
  set: value => {
    void router.replace(adminTabPathMap[value])
  }
})
const adminTabRouteTargetMap = computed<TabRouteTargetMap>(() => ({
  users: adminTabPathMap.users,
  moderation: adminTabPathMap.moderation,
  keys: adminTabPathMap.keys,
  analysis: adminTabPathMap.analysis
}))
const keysQuery = ref('')
const keysQueryKey = computed(() => `admin-keys:${adminTab.value}:${keysQuery.value.trim()}`)
const {
  pending: summaryLoading,
  error: summaryAsyncError,
  refresh: refreshAdminSummary
} = await useAsyncData(
  () => `admin-summary:${adminTab.value}`,
  async () => {
    if (adminTab.value !== 'analysis') return null
    return await loadSharedAdminSummary()
  },
  {
    server: true,
    watch: [adminTab],
    default: () => null
  }
)
const summaryLoadError = computed(() =>
  String((summaryAsyncError.value as any)?.data?.message || (summaryAsyncError.value as any)?.message || '').trim()
)
const summaryError = computed(() => summaryActionError.value || summaryLoadError.value)
const {
  data: analysisResponse,
  pending: analysisLoading,
  error: analysisAsyncError
} = await useAsyncData(
  () => `admin-analysis:${adminTab.value}`,
  async () => {
    if (adminTab.value !== 'analysis') return null
    const res = await adminTraitsSetsAnalysis()
    return res
  },
  {
    server: true,
    watch: [adminTab],
    default: () => null
  }
)
const analysis = computed<AdminTraitsSetsAnalysis | null>(() => analysisResponse.value?.data || null)
const analysisError = computed(() =>
  String((analysisAsyncError.value as any)?.data?.message || (analysisAsyncError.value as any)?.message || '').trim()
)
const {
  data: paymentsSummaryResponse,
  pending: paymentsSummaryLoading,
  error: paymentsSummaryAsyncError
} = await useAsyncData(
  () => `admin-payments-summary:${adminTab.value}`,
  async () => {
    if (adminTab.value !== 'analysis') return null
    const res = await loadPaymentsAdminSummary()
    return res
  },
  {
    server: true,
    watch: [adminTab],
    default: () => null
  }
)
const paymentsSummary = computed<PaymentAdminOrdersSummary | null>(() => paymentsSummaryResponse.value?.data || null)
const paymentsSummaryError = computed(() =>
  String(
    (paymentsSummaryAsyncError.value as any)?.data?.message || (paymentsSummaryAsyncError.value as any)?.message || ''
  ).trim()
)
const {
  data: keysResponse,
  pending: keysLoading,
  error: keysAsyncError
} = await useAsyncData(
  () => keysQueryKey.value,
  async () => {
    if (adminTab.value !== 'keys') return null
    const res = await adminSearchKeys(keysQuery.value.trim(), 50)
    return res
  },
  {
    server: true,
    watch: [adminTab, keysQuery],
    default: () => null
  }
)
const keys = computed<AdminTraitKeySearchItem[]>(() => keysResponse.value?.data?.items || [])
const keysError = computed(() =>
  String((keysAsyncError.value as any)?.data?.message || (keysAsyncError.value as any)?.message || '').trim()
)
const usersQueryKey = computed(
  () =>
    `admin-users:${adminTab.value}:${userSearch.value}:${userStatusFilter.value}:${userRoleFilter.value}:${usersLimit.value}:${usersOffset.value}`
)

const {
  data: usersResponse,
  pending: usersLoading,
  error: usersAsyncError,
  refresh: refreshUsers
} = await useAsyncData(
  () => usersQueryKey.value,
  async () => {
    if (adminTab.value !== 'users') return null
    return await adminListUsers({
      q: userSearch.value.trim(),
      status: userStatusFilter.value.trim(),
      role: userRoleFilter.value.trim(),
      limit: usersLimit.value,
      offset: usersOffset.value
    })
  },
  {
    server: true,
    watch: [adminTab, userSearch, userStatusFilter, userRoleFilter, usersLimit, usersOffset],
    default: () => null
  }
)
const loadUsers = async () => {
  await refreshUsers()
}

watchEffect(() => {
  const res = usersResponse.value
  if (!res?.data) {
    users.value = []
    usersTotal.value = 0
    return
  }
  users.value = res.data.items || []
  usersTotal.value = Number(res.data.total || 0)
  syncRoleDrafts(users.value)
})

watchEffect(() => {
  usersError.value = (usersAsyncError.value as any)?.data?.message || (usersAsyncError.value as any)?.message || ''
})
const adminBreadcrumbItems = computed<BreadcrumbItem[]>(() => {
  const labels: Record<Exclude<AdminTab, 'users'>, string> = {
    moderation: 'Модерация рецептов',
    keys: 'Ключи',
    analysis: 'Аналитика'
  }
  return [
    { label: 'Админка', to: adminTabPathMap.users },
    {
      label: adminTab.value === 'users' ? 'Пользователи' : labels[adminTab.value as Exclude<AdminTab, 'users'>],
      current: true,
      kind: 'tab'
    }
  ]
})
const adminTabItems = computed<LabTabItem[]>(() => {
  const items: LabTabItem[] = [
    { value: 'users', label: 'Пользователи' },
    { value: 'moderation', label: 'Модерация рецептов' },
    { value: 'keys', label: 'Ключи' },
    { value: 'analysis', label: 'Аналитика' }
  ]
  if (usersTotal.value) items[0]!.badge = usersTotal.value
  if (moderationStatusTotals.value.pending) items[1]!.badge = moderationStatusTotals.value.pending
  if (keys.value.length) items[2]!.badge = keys.value.length
  if (summary.value?.has_unread) {
    items[3]!.badge =
      Number(summary.value.new_users_since_last_login || 0) +
      Number(summary.value.new_pending_recipes_since_last_login || 0)
  }
  return items
})
const adminSectionClass = 'space-y-4 py-1'
const adminSubsectionClass = 'space-y-3 py-1'
const adminConfirmButtonClass = 'h-8 min-w-28'
const moderationInfo = ref('')
const moderationError = ref('')
const moderationStatus = ref<AdminModerationStatus | 'all'>('pending')
const moderationItems = ref<KitchenRecipe[]>([])
const moderationTotal = ref(0)
const moderationLimit = ref(30)
const moderationOffset = ref(0)
const moderationStatusTotals = ref<Record<string, number>>({
  all: 0,
  pending: 0,
  approved: 0,
  rejected: 0,
  draft: 0
})
const moderationOwnerEditorRecipeId = ref<string | null>(null)
const moderationOwnerQuery = ref('')
const moderationOwnerCandidates = ref<AdminUserView[]>([])
const moderationOwnerCandidatesLoading = ref(false)
const moderationOwnerCandidatesError = ref('')
const moderationOwnerSavingRecipeId = ref<string | null>(null)
const moderationOwnerDirectory = reactive<Record<string, { email: string; display_name: string }>>({})
const moderationNoteDrafts = reactive<Record<string, string>>({})
let moderationOwnerSearchRequestId = 0
const moderationHasPrevPage = computed(() => moderationOffset.value > 0)
const moderationHasNextPage = computed(() => moderationOffset.value + moderationLimit.value < moderationTotal.value)
const analysisStats = computed(() => [
  { label: 'Traits', value: displayNumber(analysis.value?.total_traits) },
  { label: 'Уникальные пары (key+value)', value: displayNumber(analysis.value?.unique_trait_pairs) },
  { label: 'Уникальные ключи', value: displayNumber(analysis.value?.unique_trait_keys) },
  { label: 'Sets', value: displayNumber(analysis.value?.total_sets) },
  { label: 'Уникальность наборов', value: analysis.value ? formatPercent(analysis.value.set_uniqueness_rate) : '—' },
  {
    label: 'Исхожесть (доля производных наборов)',
    value: analysis.value ? formatPercent(analysis.value.derived_set_rate) : '—'
  },
  { label: 'Traits в наборах', value: displayNumber(analysis.value?.traits_referenced_in_sets) },
  {
    label: 'Покрытие traits наборами',
    value: analysis.value ? formatPercent(analysis.value.trait_coverage_in_sets_rate) : '—'
  },
  { label: 'Traits вне наборов', value: displayNumber(analysis.value?.orphan_traits) }
])
function normalizeAdminTab(value: unknown): AdminTab {
  const raw = typeof value === 'string' ? value : Array.isArray(value) ? value[0] : ''
  return ADMIN_TABS.includes(raw as AdminTab) ? (raw as AdminTab) : 'users'
}
const normalizeModerationStatusTotals = (raw?: Record<string, number> | null) => {
  const base = {
    all: 0,
    pending: 0,
    approved: 0,
    rejected: 0,
    draft: 0
  }
  if (!raw || typeof raw !== 'object') return base
  for (const key of Object.keys(base)) {
    const value = Number(raw[key] || 0)
    base[key as keyof typeof base] = Number.isFinite(value) && value > 0 ? value : 0
  }
  return base
}
const hasPrevPage = computed(() => usersOffset.value > 0)
const hasNextPage = computed(() => usersOffset.value + usersLimit.value < usersTotal.value)
const unreadAdminTotal = computed(
  () =>
    Number(summary?.value?.new_users_since_last_login || 0) +
    Number(summary?.value?.new_pending_recipes_since_last_login || 0)
)
const formatPercent = (value?: number) => `${Math.round(Math.max(0, Number(value || 0)) * 100)}%`
const displayNumber = (value?: number | null) => (value === null || value === undefined ? '—' : value)
const displayMoney = (value?: number | null) => formatPaymentAmount(Number(value || 0), localeTag.value)
const displayDateTime = (value?: string | null) => (value ? formatAbsoluteDateTime(value) : '—')
const paymentsStats = computed(() => [
  { key: 'orders_total', label: 'Всего заказов', value: displayNumber(paymentsSummary.value?.orders_total) },
  { key: 'orders_success', label: 'Успешных', value: displayNumber(paymentsSummary.value?.orders_success) },
  { key: 'orders_pending', label: 'В обработке', value: displayNumber(paymentsSummary.value?.orders_pending) },
  { key: 'orders_failed', label: 'Неуспешных', value: displayNumber(paymentsSummary.value?.orders_failed) },
  { key: 'orders_refunded', label: 'Возвратов', value: displayNumber(paymentsSummary.value?.orders_refunded) },
  {
    key: 'paid_users_total',
    label: 'Платящих пользователей',
    value: displayNumber(paymentsSummary.value?.paid_users_total)
  },
  {
    key: 'patron_users_total',
    label: 'Поддержавших проект',
    value: displayNumber(paymentsSummary.value?.patron_users_total)
  },
  {
    key: 'active_access_users',
    label: 'Активный доступ',
    value: displayNumber(paymentsSummary.value?.active_access_users)
  },
  { key: 'gross_revenue', label: 'Оборот', value: displayMoney(paymentsSummary.value?.gross_revenue) },
  { key: 'net_revenue', label: 'Чистая выручка', value: displayMoney(paymentsSummary.value?.net_revenue) },
  { key: 'tip_revenue', label: 'Донаты', value: displayMoney(paymentsSummary.value?.tip_revenue) },
  { key: 'mrr', label: 'MRR', value: displayMoney(paymentsSummary.value?.mrr) },
  { key: 'churn_rate', label: 'Churn', value: formatPercent(paymentsSummary.value?.churn_rate) },
  { key: 'patron_share', label: 'Доля поддержавших', value: formatPercent(paymentsSummary.value?.patron_share) }
])
function currentRole(item: AdminUserView): 'admin' | 'user' {
  return item.roles.includes('admin') ? 'admin' : 'user'
}
const isAdminAccount = (item: AdminUserView) => currentRole(item) === 'admin'
const userStatusLabel = (value?: string) => {
  const key = String(value || '').trim()
  if (key === 'active') return 'Активен'
  if (key === 'pending_verification') return 'Ожидает подтверждения'
  if (key === 'blocked') return 'Заблокирован'
  return key || '—'
}
function syncRoleDrafts(items: AdminUserView[]) {
  for (const item of items) {
    roleDrafts[item.user_id] = currentRole(item)
  }
}
const adminListUserToAuthUser = (item: AdminUserView): AuthUser => ({
  user_id: item.user_id,
  email: item.email,
  status: item.status,
  display_name: item.display_name,
  locale: item.locale,
  timezone: item.timezone,
  roles: item.roles,
  is_two_factor_enabled: item.is_two_factor_enabled,
  last_login_at: item.last_login_at || null,
  blocked_reason: item.blocked_reason || '',
  blocked_at: item.blocked_at || null,
  created_at: item.created_at,
  profile: item.profile || {},
  settings: item.settings || {}
})
const usersTableColumns: LabDataTableColumn[] = [
  { key: 'avatar', label: '', nowrap: true, widthClass: 'w-12 min-w-12' },
  { key: 'user', label: 'Пользователь' },
  { key: 'status', label: 'Статус' },
  { key: 'role', label: 'Роль' },
  { key: 'login', label: 'Вход', nowrap: true },
  { key: 'actions', label: 'Действия', nowrap: true, widthClass: 'min-w-44' }
]
const keysTableColumns: LabDataTableColumn[] = [
  { key: 'id', label: 'ID', nowrap: true },
  { key: 'key', label: 'Ключ' },
  { key: 'meta', label: 'Meta' },
  { key: 'traits', label: 'Traits', nowrap: true }
]
const moderationTableColumns: LabDataTableColumn[] = [
  { key: 'recipe', label: 'Рецепт', widthClass: 'w-full min-w-96' },
  { key: 'owner', label: 'Владелец', widthClass: 'min-w-32' },
  { key: 'status', label: 'Статус', nowrap: true, widthClass: 'min-w-32' },
  { key: 'created', label: 'Создан', nowrap: true, widthClass: 'min-w-36' },
  { key: 'actions', label: 'Действия', nowrap: true, widthClass: 'min-w-44' }
]
const setRole = async (item: AdminUserView) => {
  const role = roleDrafts[item.user_id]
  if (!role) return
  if (item.user_id === user.value?.user_id && role !== 'admin') {
    roleDrafts[item.user_id] = 'admin'
    usersError.value = 'Нельзя изменить собственную роль администратора.'
    return
  }
  if (role === currentRole(item)) return
  if (roleSavingUserId.value === item.user_id) return
  roleSavingUserId.value = item.user_id
  usersError.value = ''
  usersInfo.value = ''
  try {
    await adminSetUserRole(item.user_id, role)
    usersInfo.value = `Роль обновлена: ${item.email} -> ${role}.`
    await loadUsers()
  } catch (err: any) {
    usersError.value = err?.data?.message || err?.message || 'Не удалось обновить роль.'
  } finally {
    if (roleSavingUserId.value === item.user_id) {
      roleSavingUserId.value = null
    }
  }
}
const blockUser = async (item: AdminUserView) => {
  if (isAdminAccount(item)) {
    usersError.value = 'Нельзя блокировать пользователя с ролью admin.'
    return
  }
  usersError.value = ''
  usersInfo.value = ''
  try {
    await adminBlockUser(item.user_id, '')
    usersInfo.value = `Пользователь ${item.email} заблокирован.`
    await loadUsers()
  } catch (err: any) {
    usersError.value = err?.data?.message || err?.message || 'Не удалось заблокировать пользователя.'
  }
}
const unblockUser = async (item: AdminUserView) => {
  usersError.value = ''
  usersInfo.value = ''
  try {
    await adminUnblockUser(item.user_id)
    usersInfo.value = `Пользователь ${item.email} разблокирован.`
    await loadUsers()
  } catch (err: any) {
    usersError.value = err?.data?.message || err?.message || 'Не удалось разблокировать пользователя.'
  }
}
const forceLogoutUser = async (item: AdminUserView) => {
  if (isAdminAccount(item)) {
    usersError.value = 'Нельзя сбрасывать сессии пользователя с ролью admin.'
    return
  }
  usersError.value = ''
  usersInfo.value = ''
  try {
    await adminForceLogoutUser(item.user_id)
    usersInfo.value = `Все сессии пользователя ${item.email} завершены.`
  } catch (err: any) {
    usersError.value = err?.data?.message || err?.message || 'Не удалось принудительно завершить сессии.'
  }
}
const deleteUser = async (item: AdminUserView) => {
  if (isAdminAccount(item)) {
    usersError.value = 'Нельзя удалять пользователя с ролью admin.'
    return
  }
  usersError.value = ''
  usersInfo.value = ''
  try {
    await adminDeleteUser(item.user_id)
    usersInfo.value = `Аккаунт ${item.email} удалён.`
    await loadUsers()
  } catch (err: any) {
    usersError.value = err?.data?.message || err?.message || 'Не удалось удалить пользователя.'
  }
}
const keyMetaRecord = (meta?: Record<string, any>) => {
  if (!meta || typeof meta !== 'object' || Array.isArray(meta)) return {}
  return meta
}
const keyMetaSummaryItems = (meta?: Record<string, any>) => {
  const safe = keyMetaRecord(meta)
  const items: string[] = []
  const dataType = String(safe.dataType || '').trim()
  const unitCategory = String(safe.unitCategory || '').trim()
  const unit = String(safe.unit || '').trim()
  const mode = String(safe.mode || '').trim()
  const optionType = String(safe.optionType || '').trim()
  const optionsCount = Array.isArray(safe.options) ? safe.options.length : 0
  if (dataType) items.push(`тип: ${dataType}`)
  if (unitCategory) items.push(`категория ед.: ${unitCategory}`)
  if (unit) items.push(`ед. изм.: ${unit}`)
  if (mode) items.push(`режим: ${mode}`)
  if (optionType) items.push(`тип опций: ${optionType}`)
  if (optionsCount > 0) items.push(`опций: ${optionsCount}`)
  return items
}
const keyMetaJson = (meta?: Record<string, any>) => {
  const safe = keyMetaRecord(meta)
  if (!Object.keys(safe).length) return '—'
  return JSON.stringify(safe, null, 2)
}
const moderationStatusLabel = (value: string) => {
  if (value === 'all') return 'Все'
  if (value === 'pending') return 'На модерации'
  if (value === 'approved') return 'Одобрен'
  if (value === 'rejected') return 'Отклонён'
  if (value === 'draft') return 'Черновик'
  return value || '—'
}
const userStatusTone = (value?: string): 'success' | 'warning' | 'danger' | 'muted' => {
  if (value === 'active') return 'success'
  if (value === 'pending_verification') return 'warning'
  if (value === 'blocked') return 'danger'
  return 'muted'
}
const moderationStatusTone = (value?: string): 'success' | 'warning' | 'danger' | 'muted' => {
  if (value === 'approved') return 'success'
  if (value === 'pending') return 'warning'
  if (value === 'rejected') return 'danger'
  return 'muted'
}
const moderationStatusStatItems = computed(
  () =>
    [
      { key: 'all', label: 'Все' },
      { key: 'pending', label: 'На модерации' },
      { key: 'approved', label: 'Одобрено' },
      { key: 'rejected', label: 'Отклонено' },
      { key: 'draft', label: 'Черновики' }
    ] as const
)
const moderationCanApprove = (item: KitchenRecipe) => item.moderation_status !== 'approved'
const moderationCanReject = (item: KitchenRecipe) =>
  item.moderation_status === 'pending' || item.moderation_status === 'approved'
const moderationRecipeLink = (item: KitchenRecipe) => `/kitchen/${item.id}`
const moderationNoteDraft = (item: KitchenRecipe) => {
  const recipeId = String(item.id || '').trim()
  if (!recipeId) return ''
  if (!(recipeId in moderationNoteDrafts)) {
    moderationNoteDrafts[recipeId] = String(item.moderation_note || '').trim()
  }
  return moderationNoteDrafts[recipeId]
}
const cacheOwnerUsers = (items: AdminUserView[]) => {
  for (const item of items) {
    const id = String(item.user_id || '').trim()
    if (!id) continue
    moderationOwnerDirectory[id] = {
      email: String(item.email || '').trim(),
      display_name: String(item.display_name || '').trim()
    }
  }
}
const moderationOwnerDisplay = (ownerUserId?: string | null) => {
  const id = String(ownerUserId || '').trim()
  if (!id) return '—'
  const cached = moderationOwnerDirectory[id]
  if (!cached) return id
  if (cached.display_name) return `${cached.display_name} · ${cached.email}`
  if (cached.email) return cached.email
  return id
}
const moderationOwnerSubLabel = (ownerUserId?: string | null) => {
  const id = String(ownerUserId || '').trim()
  if (!id) return 'без владельца'
  const cached = moderationOwnerDirectory[id]
  if (!cached) return id
  return id
}
const isModerationOwnerEditorOpen = (recipeId?: string) =>
  String(recipeId || '').trim() !== '' && moderationOwnerEditorRecipeId.value === recipeId
const closeModerationOwnerEditor = () => {
  moderationOwnerEditorRecipeId.value = null
  moderationOwnerQuery.value = ''
  moderationOwnerCandidates.value = []
  moderationOwnerCandidatesError.value = ''
}
const searchModerationOwnerCandidates = async (queryRaw: string) => {
  if (!moderationOwnerEditorRecipeId.value) return
  const query = String(queryRaw || '').trim()
  if (query.length < 2) {
    moderationOwnerCandidates.value = []
    moderationOwnerCandidatesError.value = ''
    moderationOwnerCandidatesLoading.value = false
    return
  }
  const requestId = ++moderationOwnerSearchRequestId
  moderationOwnerCandidatesLoading.value = true
  moderationOwnerCandidatesError.value = ''
  try {
    const res = await adminListUsers({ q: query, limit: 12, offset: 0 })
    if (requestId !== moderationOwnerSearchRequestId) return
    const items = res.data.items || []
    moderationOwnerCandidates.value = items
    cacheOwnerUsers(items)
  } catch (err: any) {
    if (requestId !== moderationOwnerSearchRequestId) return
    moderationOwnerCandidates.value = []
    moderationOwnerCandidatesError.value = err?.data?.message || err?.message || 'Не удалось найти пользователей.'
  } finally {
    if (requestId === moderationOwnerSearchRequestId) {
      moderationOwnerCandidatesLoading.value = false
    }
  }
}
const debouncedSearchModerationOwnerCandidates = debounce((query: string) => {
  void searchModerationOwnerCandidates(query)
}, 250)
const openModerationOwnerEditor = (item: KitchenRecipe) => {
  const recipeId = String(item.id || '').trim()
  if (!recipeId) return
  if (moderationOwnerSavingRecipeId.value === recipeId) return
  if (moderationOwnerEditorRecipeId.value === recipeId) {
    closeModerationOwnerEditor()
    return
  }
  moderationOwnerEditorRecipeId.value = recipeId
  moderationOwnerCandidates.value = []
  moderationOwnerCandidatesError.value = ''
  const ownerId = String(item.owner_user_id || '').trim()
  const cached = ownerId ? moderationOwnerDirectory[ownerId] : null
  moderationOwnerQuery.value = String(cached?.email || cached?.display_name || '').trim()
  if (moderationOwnerQuery.value.length >= 2) {
    debouncedSearchModerationOwnerCandidates(moderationOwnerQuery.value)
  }
}
const changeModerationRecipeOwner = async (item: KitchenRecipe, nextOwner: AdminUserView) => {
  const recipeId = String(item.id || '').trim()
  const nextOwnerId = String(nextOwner.user_id || '').trim()
  if (!recipeId || !nextOwnerId) return
  if (moderationOwnerSavingRecipeId.value === recipeId) return
  if (String(item.owner_user_id || '').trim() === nextOwnerId) {
    closeModerationOwnerEditor()
    return
  }
  moderationOwnerSavingRecipeId.value = recipeId
  moderationError.value = ''
  moderationInfo.value = ''
  try {
    const res = await adminChangeKitchenRecipeOwner(recipeId, nextOwnerId)
    moderationItems.value = moderationItems.value.map(row => (row.id === recipeId ? res.data : row))
    cacheOwnerUsers([nextOwner])
    moderationInfo.value = `Владелец рецепта «${item.title}» изменён на ${nextOwner.email}.`
    closeModerationOwnerEditor()
  } catch (err: any) {
    moderationError.value = err?.data?.message || err?.message || 'Не удалось изменить владельца рецепта.'
  } finally {
    if (moderationOwnerSavingRecipeId.value === recipeId) {
      moderationOwnerSavingRecipeId.value = null
    }
  }
}
const {
  data: moderationResponse,
  pending: moderationLoading,
  refresh: refreshModeration
} = await useAsyncData(
  () =>
    `admin-moderation:${adminTab.value}:${moderationStatus.value}:${moderationLimit.value}:${moderationOffset.value}`,
  async () => {
    if (adminTab.value !== 'moderation') return null
    return await adminListKitchenModerationRecipes({
      status: moderationStatus.value,
      limit: moderationLimit.value,
      offset: moderationOffset.value
    })
  },
  {
    server: true,
    watch: [adminTab, moderationStatus, moderationLimit, moderationOffset],
    default: () => null
  }
)

watchEffect(() => {
  const res = moderationResponse.value
  if (!res?.data) {
    moderationItems.value = []
    moderationTotal.value = 0
    return
  }
  moderationItems.value = res.data.items || []
  moderationTotal.value = Number(res.data.total || 0)
  moderationStatusTotals.value = normalizeModerationStatusTotals(res.data.status_totals)
  for (const item of moderationItems.value) {
    const recipeId = String(item.id || '').trim()
    if (!recipeId) continue
    moderationNoteDrafts[recipeId] = String(item.moderation_note || '').trim()
  }
})
const moderationPrevPage = async () => {
  if (!moderationHasPrevPage.value) return
  moderationOffset.value = Math.max(0, moderationOffset.value - moderationLimit.value)
}
const moderationNextPage = async () => {
  if (!moderationHasNextPage.value) return
  moderationOffset.value = moderationOffset.value + moderationLimit.value
}
const moderateRecipe = async (item: KitchenRecipe, approve: boolean) => {
  moderationError.value = ''
  moderationInfo.value = ''
  const note = String(moderationNoteDraft(item) || '').trim()
  if (!approve && !note) {
    moderationError.value = 'Укажите причину отклонения рецепта.'
    return
  }
  try {
    await adminModerateKitchenRecipe(item.id, approve, note)
    moderationInfo.value = approve
      ? `Рецепт «${item.title}» одобрен и опубликован.`
      : `Рецепт «${item.title}» отклонён.`
    await refreshModeration()
  } catch (err: any) {
    moderationError.value = err?.data?.message || err?.message || 'Не удалось изменить статус модерации.'
  }
}
const markSummaryRead = async () => {
  if (summaryReadPending.value) return
  summaryReadPending.value = true
  summaryActionError.value = ''
  summaryInfo.value = ''
  try {
    await adminMarkSummaryRead()
    await refreshAdminSummary()
    summaryInfo.value = 'Уведомления отмечены как прочитанные.'
  } catch (err: any) {
    summaryActionError.value = err?.data?.message || err?.message || 'Не удалось отметить уведомления как прочитанные.'
  } finally {
    summaryReadPending.value = false
  }
}
const prevPage = async () => {
  if (!hasPrevPage.value) return
  usersOffset.value = Math.max(0, usersOffset.value - usersLimit.value)
}
const nextPage = async () => {
  if (!hasNextPage.value) return
  usersOffset.value = usersOffset.value + usersLimit.value
}
const applyUserFilters = async () => {
  usersOffset.value = 0
}
const debouncedApplyUserFilters = debounce(() => {
  void applyUserFilters()
}, 300)
watch(moderationStatus, async (next, prev) => {
  if (next === prev) return
  moderationOffset.value = 0
  closeModerationOwnerEditor()
  if (adminTab.value === 'moderation') {
    await refreshModeration()
  }
})
watch(moderationOwnerQuery, next => {
  if (!moderationOwnerEditorRecipeId.value) return
  debouncedSearchModerationOwnerCandidates(String(next || ''))
})
watch([userSearch, userStatusFilter, userRoleFilter], () => {
  if (adminTab.value !== 'users') return
  debouncedApplyUserFilters()
})
watch(users, items => {
  cacheOwnerUsers(items || [])
})
watch(
  sharedAdminSummary,
  nextSummary => {
    if (!nextSummary?.recipe_status_totals) return
    moderationStatusTotals.value = normalizeModerationStatusTotals(nextSummary.recipe_status_totals)
  },
  { immediate: true }
)
</script>
<template>
  <div>
    <LabNavHeader :title :breadcrumb-items="adminBreadcrumbItems">
      <template #actions>
        <LabBaseBadge v-if="summary?.has_unread" variant="warning" size="xs" :rounded="false">
          новое {{ unreadAdminTotal }}
        </LabBaseBadge>
      </template>
    </LabNavHeader>
    <LabNavTabs
      v-model="adminTab"
      :items="adminTabItems"
      :render-panels="false"
      route-param-key="tab"
      :route-target-map="adminTabRouteTargetMap"
    />
    <div class="mt-4 space-y-4 px-3 md:px-4">
      <section v-show="adminTab === 'users'" :class="adminSectionClass">
        <AdminSectionHeader
          title="Пользователи и роли"
          description="Поиск, блокировки, смена роли и принудительный выход из сессий."
        />
        <div class="grid gap-3 md:grid-cols-3 xl:grid-cols-4">
          <LabField label="Поиск" for-id="admin_user_search">
            <LabBaseInput
              id="admin_user_search"
              v-model="userSearch"
              name="admin_user_search"
              type="text"
              placeholder="email или имя"
            />
          </LabField>
          <LabField label="Статус" for-id="admin_user_status">
            <LabBaseSelect
              id="admin_user_status"
              v-model="userStatusFilter"
              name="admin_user_status"
              :options="userStatusOptions"
            />
          </LabField>
          <LabField label="Роль" for-id="admin_user_role_filter">
            <LabBaseSelect
              id="admin_user_role_filter"
              v-model="userRoleFilter"
              name="admin_user_role_filter"
              :options="userRoleOptions"
            />
          </LabField>
        </div>
        <AdminNotifyStack
          :items="[
            { text: usersError, tone: 'error' },
            { text: usersInfo, tone: 'success' }
          ]"
        />
        <LabDataTable
          :columns="usersTableColumns"
          :rows="users"
          :loading="usersLoading"
          empty-text="Пользователи не найдены."
          row-key="user_id"
        >
          <template #cell-avatar="{ row }">
            <NuxtLink :to="`/users/${(row as AdminUserRow).user_id}`" class="inline-flex">
              <LabAvatar version="preview" :user="adminListUserToAuthUser(row as AdminUserRow)" :show-label="false" />
            </NuxtLink>
          </template>
          <template #cell-user="{ row }">
            <div class="space-y-1">
              <NuxtLink
                :to="`/users/${(row as AdminUserRow).user_id}`"
                class="text-(--lab-text-primary) hover:text-(--lab-accent)"
              >
                {{ (row as AdminUserRow).email }}
              </NuxtLink>
              <p class="text-zinc-500">{{ (row as AdminUserRow).display_name || '—' }}</p>
            </div>
          </template>
          <template #cell-status="{ row }">
            <div class="space-y-1">
              <LabBaseBadge :variant="userStatusTone((row as AdminUserRow).status)" size="xs" :rounded="false">
                {{ userStatusLabel((row as AdminUserRow).status) }}
              </LabBaseBadge>
              <p v-if="(row as AdminUserRow).blocked_reason">причина: {{ (row as AdminUserRow).blocked_reason }}</p>
            </div>
          </template>
          <template #cell-role="{ row }">
            <div class="flex items-center gap-2">
              <LabBaseSelect
                v-model="roleDrafts[(row as AdminUserRow).user_id]"
                :name="`admin_user_role_${(row as AdminUserRow).user_id}`"
                :options="editableRoleOptions"
                class="min-w-28"
                :disabled="roleSavingUserId === (row as AdminUserRow).user_id || (row as AdminUserRow).user_id === user?.user_id"
                @change="setRole(row as AdminUserRow)"
              />
              <span v-if="roleSavingUserId === (row as AdminUserRow).user_id">сохранение...</span>
              <span v-else-if="(row as AdminUserRow).user_id === user?.user_id">текущий админ</span>
            </div>
          </template>
          <template #cell-login="{ row }">
            <LabRelativeTime :datetime="(row as AdminUserRow).last_login_at" compact />
          </template>
          <template #cell-actions="{ row }">
            <div class="flex min-w-max items-center gap-1 whitespace-nowrap">
              <span v-if="(row as AdminUserRow).status !== 'blocked' && !isAdminAccount(row as AdminUserRow)" class="max-sm:hidden">
                <LabConfirmActionButton
                  label="Блок"
                  confirm-label="Подтвердить"
                  tooltip="Подтвердить блокировку аккаунта?"
                  icon="ic:round-block"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border-[color-mix(in_srgb,var(--lab-danger)_42%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-danger)_12%,var(--lab-bg-surface))] text-(--lab-text-primary) hover:bg-[color-mix(in_srgb,var(--lab-danger)_18%,var(--lab-bg-surface-hover))]"
                  confirm-class="border-(--lab-danger) bg-(--lab-danger) text-white hover:bg-[color-mix(in_srgb,var(--lab-danger)_88%,black)]"
                  progress-class="bg-[color-mix(in_srgb,var(--lab-danger)_30%,transparent)]"
                  @confirm="blockUser(row as AdminUserRow)"
                />
              </span>
              <span v-if="(row as AdminUserRow).status !== 'blocked' && !isAdminAccount(row as AdminUserRow)" class="sm:hidden">
                <LabConfirmActionButton
                  icon-only
                  icon="ic:round-block"
                  aria-label="Заблокировать"
                  confirm-aria-label="Подтвердить блокировку"
                  confirm-label="Ок"
                  tooltip="Подтвердить блокировку аккаунта?"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border-[color-mix(in_srgb,var(--lab-danger)_42%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-danger)_12%,var(--lab-bg-surface))] text-(--lab-text-primary) hover:bg-[color-mix(in_srgb,var(--lab-danger)_18%,var(--lab-bg-surface-hover))]"
                  confirm-class="border-(--lab-danger) bg-(--lab-danger) text-white hover:bg-[color-mix(in_srgb,var(--lab-danger)_88%,black)]"
                  progress-class="bg-[color-mix(in_srgb,var(--lab-danger)_30%,transparent)]"
                  @confirm="blockUser(row as AdminUserRow)"
                />
              </span>
              <LabBaseButton
                v-if="(row as AdminUserRow).status === 'blocked'"
                class="max-sm:hidden"
                variant="secondary"
                size="xs"
                icon="ic:round-lock-open"
                label="Разблок"
                @click="unblockUser(row as AdminUserRow)"
              />
              <LabBaseButton
                v-if="(row as AdminUserRow).status === 'blocked'"
                class="sm:hidden"
                variant="secondary"
                size="xs"
                icon="ic:round-lock-open"
                icon-only
                aria-label="Разблокировать"
                @click="unblockUser(row as AdminUserRow)"
              />
              <span v-if="!isAdminAccount(row as AdminUserRow)" class="max-sm:hidden">
                <LabConfirmActionButton
                  label="Сброс сессий"
                  confirm-label="Подтвердить"
                  tooltip="Подтвердить принудительный сброс всех сессий?"
                  icon="ic:round-logout"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border-(--lab-border) bg-(--lab-bg-surface) text-(--lab-text-primary) hover:bg-(--lab-bg-surface-hover)"
                  confirm-class="border-(--lab-warning) bg-(--lab-warning) text-white hover:bg-[color-mix(in_srgb,var(--lab-warning)_88%,black)]"
                  progress-class="bg-[color-mix(in_srgb,var(--lab-warning)_30%,transparent)]"
                  tooltip-class="border-[color-mix(in_srgb,var(--lab-warning)_52%,var(--lab-border))] text-(--lab-text-primary)"
                  @confirm="forceLogoutUser(row as AdminUserRow)"
                />
              </span>
              <span v-if="!isAdminAccount(row as AdminUserRow)" class="sm:hidden">
                <LabConfirmActionButton
                  icon-only
                  icon="ic:round-logout"
                  aria-label="Сбросить сессии"
                  confirm-aria-label="Подтвердить сброс сессий"
                  confirm-label="Ок"
                  tooltip="Подтвердить принудительный сброс всех сессий?"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border-(--lab-border) bg-(--lab-bg-surface) text-(--lab-text-primary) hover:bg-(--lab-bg-surface-hover)"
                  confirm-class="border-(--lab-warning) bg-(--lab-warning) text-white hover:bg-[color-mix(in_srgb,var(--lab-warning)_88%,black)]"
                  progress-class="bg-[color-mix(in_srgb,var(--lab-warning)_30%,transparent)]"
                  tooltip-class="border-[color-mix(in_srgb,var(--lab-warning)_52%,var(--lab-border))] text-(--lab-text-primary)"
                  @confirm="forceLogoutUser(row as AdminUserRow)"
                />
              </span>
              <span v-if="!isAdminAccount(row as AdminUserRow)" class="max-sm:hidden">
                <LabConfirmActionButton
                  label="Удалить"
                  confirm-label="Подтвердить"
                  tooltip="Подтвердить удаление аккаунта?"
                  icon="ic:round-delete"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border-[color-mix(in_srgb,var(--lab-danger)_46%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-danger)_14%,var(--lab-bg-surface))] text-(--lab-text-primary) hover:bg-[color-mix(in_srgb,var(--lab-danger)_20%,var(--lab-bg-surface-hover))]"
                  confirm-class="border-(--lab-danger) bg-(--lab-danger) text-white hover:bg-[color-mix(in_srgb,var(--lab-danger)_88%,black)]"
                  progress-class="bg-[color-mix(in_srgb,var(--lab-danger)_30%,transparent)]"
                  @confirm="deleteUser(row as AdminUserRow)"
                />
              </span>
              <span v-if="!isAdminAccount(row as AdminUserRow)" class="sm:hidden">
                <LabConfirmActionButton
                  icon-only
                  icon="ic:round-delete"
                  aria-label="Удалить аккаунт"
                  confirm-aria-label="Подтвердить удаление аккаунта"
                  confirm-label="Ок"
                  tooltip="Подтвердить удаление аккаунта?"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border-[color-mix(in_srgb,var(--lab-danger)_46%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-danger)_14%,var(--lab-bg-surface))] text-(--lab-text-primary) hover:bg-[color-mix(in_srgb,var(--lab-danger)_20%,var(--lab-bg-surface-hover))]"
                  confirm-class="border-(--lab-danger) bg-(--lab-danger) text-white hover:bg-[color-mix(in_srgb,var(--lab-danger)_88%,black)]"
                  progress-class="bg-[color-mix(in_srgb,var(--lab-danger)_30%,transparent)]"
                  @confirm="deleteUser(row as AdminUserRow)"
                />
              </span>
            </div>
          </template>
        </LabDataTable>
        <AdminPager
          :prev-disabled="!hasPrevPage || usersLoading"
          :next-disabled="!hasNextPage || usersLoading"
          @prev="prevPage"
          @next="nextPage"
        />
      </section>
      <section v-show="adminTab === 'moderation'" :class="adminSectionClass">
        <div class="flex flex-wrap gap-2">
          <LabBaseButton
            v-for="item in moderationStatusStatItems"
            :key="`moderation-stat:${item.key}`"
            :variant="moderationStatus === item.key ? 'primary' : 'secondary'"
            size="xs"
            @click="moderationStatus = item.key"
          >
            {{ item.label }} {{ displayNumber(moderationStatusTotals[item.key]) }}
          </LabBaseButton>
        </div>
        <AdminNotifyStack
          :items="[
            { text: moderationError, tone: 'error' },
            { text: moderationInfo, tone: 'success' }
          ]"
        />
        <LabDataTable
          :columns="moderationTableColumns"
          :rows="moderationItems"
          :loading="moderationLoading"
          empty-text="Рецепты для выбранного статуса не найдены."
          row-key="id"
        >
          <template #cell-recipe="{ row }">
            <div class="min-w-0">
              <NuxtLink
                :to="moderationRecipeLink(row as ModerationRow)"
                class="inline-flex max-w-full items-center gap-1 text-zinc-200 hover:text-zinc-100 hover:underline"
              >
                <span class="truncate">{{ (row as ModerationRow).title }}</span>
                <Icon name="ic:round-link" class="h-3.5 w-3.5 shrink-0 text-zinc-500" />
              </NuxtLink>
              <p class="line-clamp-2 text-zinc-500">{{ (row as ModerationRow).description || '—' }}</p>
            </div>
          </template>
          <template #cell-owner="{ row }">
            <div class="min-w-0 space-y-1">
              <LabBaseButton
                variant="ghost"
                size="none"
                button-class="truncate inline-flex max-w-full items-start justify-start px-0 text-left text-zinc-300 hover:text-zinc-100"
                :disabled="moderationOwnerSavingRecipeId === (row as ModerationRow).id"
                :label="moderationOwnerDisplay((row as ModerationRow).owner_user_id)"
                @click="openModerationOwnerEditor(row as ModerationRow)"
              />
              <p class="truncate text-xs text-zinc-500">
                {{ moderationOwnerSubLabel((row as ModerationRow).owner_user_id) }}
              </p>
              <div v-if="isModerationOwnerEditorOpen((row as ModerationRow).id)" :class="adminSubsectionClass">
                <LabBaseInput
                  v-model="moderationOwnerQuery"
                  :name="`moderation_owner_query_${(row as ModerationRow).id}`"
                  type="text"
                  placeholder="Поиск: email или имя"
                />
                <LabNotify :text="moderationOwnerCandidatesError" tone="error" size="xs" />
                <p v-if="moderationOwnerCandidatesLoading">Поиск...</p>
                <div v-else-if="moderationOwnerCandidates.length" class="max-h-36 space-y-1 overflow-y-auto pr-1">
                  <LabBaseButton
                    v-for="candidate in moderationOwnerCandidates"
                    :key="`owner-candidate:${(row as ModerationRow).id}:${candidate.user_id}`"
                    variant="secondary"
                    size="none"
                    button-class="w-full justify-start px-2 py-1 text-left text-xs truncate"
                    :disabled="moderationOwnerSavingRecipeId === (row as ModerationRow).id"
                    :label="candidate.display_name || candidate.email"
                    @click="changeModerationRecipeOwner(row as ModerationRow, candidate)"
                  />
                </div>
                <p v-else-if="moderationOwnerQuery.trim().length >= 2">Пользователи не найдены.</p>
                <p v-else>Введите минимум 2 символа для поиска.</p>
                <LabBaseButton
                  variant="ghost"
                  size="xs"
                  button-class="px-0 text-zinc-400 hover:text-zinc-200"
                  @click="closeModerationOwnerEditor"
                >
                  Отмена
                </LabBaseButton>
              </div>
            </div>
          </template>
          <template #cell-status="{ row }">
            <LabBaseBadge :variant="moderationStatusTone((row as ModerationRow).moderation_status)" size="xs" :rounded="false">
              {{ moderationStatusLabel((row as ModerationRow).moderation_status) }}
            </LabBaseBadge>
          </template>
          <template #cell-created="{ row }">
            <LabRelativeTime :datetime="(row as ModerationRow).created_at" compact />
          </template>
          <template #cell-actions="{ row }">
            <div class="min-w-64 space-y-2">
              <LabBaseTextarea
                v-if="moderationCanReject(row as ModerationRow)"
                :model-value="moderationNoteDraft(row as ModerationRow)"
                :name="`moderation_note_${(row as ModerationRow).id}`"
                rows="2"
                placeholder="Причина отклонения для автора"
                class="w-full text-xs"
                @update:model-value="moderationNoteDrafts[(row as ModerationRow).id] = String($event || '')"
              />
              <p v-else-if="(row as ModerationRow).moderation_note">Последняя причина: {{ (row as ModerationRow).moderation_note }}</p>
              <div class="flex flex-wrap gap-1">
                <LabConfirmActionButton
                  v-if="moderationCanApprove(row as ModerationRow)"
                  label="Одобрить"
                  confirm-label="Подтвердить"
                  tooltip="Подтвердить публикацию рецепта?"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border border-emerald-500/40 bg-emerald-500/10 text-emerald-100 hover:bg-emerald-500/20"
                  confirm-class="border border-emerald-300/90 bg-emerald-600 text-white hover:bg-emerald-500"
                  progress-class="bg-emerald-300/45"
                  @confirm="moderateRecipe(row as ModerationRow, true)"
                />
                <LabConfirmActionButton
                  v-if="moderationCanReject(row as ModerationRow)"
                  label="Отклонить"
                  confirm-label="Подтвердить"
                  tooltip="Подтвердить отклонение рецепта?"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border border-rose-500/50 bg-rose-500/10 text-rose-200 hover:bg-rose-500/20"
                  confirm-class="border border-rose-300/90 bg-rose-600 text-white hover:bg-rose-500"
                  progress-class="bg-rose-300/45"
                  @confirm="moderateRecipe(row as ModerationRow, false)"
                />
              </div>
            </div>
          </template>
        </LabDataTable>
        <AdminPager
          :prev-disabled="!moderationHasPrevPage || moderationLoading"
          :next-disabled="!moderationHasNextPage || moderationLoading"
          @prev="moderationPrevPage"
          @next="moderationNextPage"
        />
      </section>
      <section v-show="adminTab === 'keys'" :class="adminSectionClass">
        <AdminSectionHeader
          title="Поиск по ключам особенностей"
          description="Проверка syn, meta и количества traits для ключа."
        />
        <LabField label="Запрос" for-id="admin_keys_query" class="max-w-xl">
          <LabBaseInput
            id="admin_keys_query"
            v-model="keysQuery"
            name="admin_keys_query"
            type="text"
            placeholder="часть названия ключа"
          />
        </LabField>
        <LabNotify :text="keysError" tone="error" size="xs" />
        <LabDataTable
          :columns="keysTableColumns"
          :rows="keys"
          :loading="keysLoading"
          empty-text="Ничего не найдено."
          row-key="key_id"
        >
          <template #cell-id="{ row }">
            <span class="text-zinc-400">{{ (row as AdminTraitKeyRow).key_id }}</span>
          </template>
          <template #cell-key="{ row }">
            <span class="text-zinc-200">{{ (row as AdminTraitKeyRow).syn }}</span>
          </template>
          <template #cell-meta="{ row }">
            <div class="space-y-1">
              <div class="flex flex-wrap gap-1">
                <LabBaseBadge
                  v-for="label in keyMetaSummaryItems((row as AdminTraitKeyRow).meta)"
                  :key="`${(row as AdminTraitKeyRow).key_id}:${label}`"
                  variant="muted"
                  size="xs"
                  :rounded="false"
                >
                  {{ label }}
                </LabBaseBadge>
                <LabBaseBadge
                  v-if="!keyMetaSummaryItems((row as AdminTraitKeyRow).meta).length"
                  variant="muted"
                  size="xs"
                  :rounded="false"
                >
                  —
                </LabBaseBadge>
              </div>
              <pre class="max-w-80 overflow-x-auto p-2 text-xs leading-5 text-zinc-400">{{
                keyMetaJson((row as AdminTraitKeyRow).meta)
              }}</pre>
            </div>
          </template>
          <template #cell-traits="{ row }">
            <span class="text-zinc-300">{{ (row as AdminTraitKeyRow).trait_count }}</span>
          </template>
        </LabDataTable>
      </section>
      <section v-show="adminTab === 'analysis'" :class="adminSectionClass">
        <AdminSectionHeader
          title="Аналитика"
          description="Сводка по уведомлениям, coverage traits и ключам с максимальной нагрузкой."
        />
        <AdminNotifyStack
          :items="[
            { text: analysisError, tone: 'error' },
            { text: summaryError, tone: 'error' },
            { text: paymentsSummaryError, tone: 'error' },
            { text: summaryInfo, tone: 'success' }
          ]"
        />
        <p v-if="analysisLoading || summaryLoading || paymentsSummaryLoading">Обновление данных...</p>
        <div :class="adminSubsectionClass">
          <div class="flex flex-wrap items-center justify-between gap-2">
            <p class="text-zinc-200">Сводка уведомлений</p>
            <LabBaseButton
              variant="secondary"
              size="xs"
              label="Прочитано"
              :disabled="summaryReadPending || !summary?.has_unread"
              @click="markSummaryRead"
            />
          </div>
          <div class="grid gap-2 sm:grid-cols-2">
            <div class="space-y-1 p-2.5">
              <p class="text-zinc-500">Новых регистраций</p>
              <p class="text-zinc-100">{{ Number(summary?.new_users_since_last_login || 0) }}</p>
            </div>
            <div class="space-y-1 p-2.5">
              <p class="text-zinc-500">Новых рецептов на модерации</p>
              <p class="text-zinc-100">{{ Number(summary?.new_pending_recipes_since_last_login || 0) }}</p>
            </div>
          </div>
        </div>
        <div :class="adminSubsectionClass">
          <div class="flex flex-wrap items-center justify-between gap-2">
            <div class="space-y-1">
              <p class="text-zinc-200">Payments</p>
              <p>Сводка заказов, выручки, активного доступа и доли поддержки проекта.</p>
            </div>
            <p>
              Последний успешный платёж:
              <span class="text-zinc-300">{{ displayDateTime(paymentsSummary?.last_successful_paid) }}</span>
            </p>
          </div>
        </div>
        <div class="grid gap-3 text-xs sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          <AdminStatCard v-for="item in paymentsStats" :key="item.key" :label="item.label" :value="item.value" />
        </div>
        <div class="grid gap-3 text-xs sm:grid-cols-2 lg:grid-cols-3">
          <AdminStatCard v-for="item in analysisStats" :key="item.label" :label="item.label" :value="item.value" />
        </div>
        <div :class="adminSubsectionClass">
          <p class="text-zinc-400">Топ ключей по количеству traits:</p>
          <ul v-if="analysis?.top_keys?.length" class="space-y-1">
            <li
              v-for="item in analysis?.top_keys || []"
              :key="item.key_id"
              class="flex items-center justify-between gap-2 px-2 py-1.5 text-zinc-300"
            >
              <span class="truncate">{{ item.syn }} (id {{ item.key_id }})</span>
              <span class="text-zinc-100">{{ item.trait_count }}</span>
            </li>
          </ul>
          <p v-else class="text-zinc-500">Данные пока не получены.</p>
        </div>
      </section>
    </div>
  </div>
</template>
