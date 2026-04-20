<script setup lang="ts">
const { localeTag } = useInterfacePreferences()
const { locale, key, load, t } = useI18nSection('payments')
await useAsyncData(key.value, load, { watch: [locale] })
const title = computed(() => t('success.seo_title'))
const description = computed(() => t('success.seo_description'))

usePageSeo({
  title,
  description
})

const route = useRoute()

const { lookupPublicOrder } = usePayments()
const { formatAbsoluteDateTime } = useLocalizedDateTime()

const loading = ref(true)
const syncing = ref(false)
const errorMessage = ref('')
const order = ref<PaymentOrderView | null>(null)

const sectionClass = 'bg-(--lab-bg-elevated) p-4 '

const orderId = computed(() => String(route.query.order_id || '').trim())
const token = computed(() => String(route.query.token || '').trim())
const nextPath = computed(() => normalizeInternalPath(String(route.query.next || '').trim()))

const statusTitle = computed(() => {
  switch (order.value?.status) {
    case 'success':
      return t('lookup.status.success')
    case 'pending':
      return t('lookup.status.pending')
    case 'failed':
      return t('lookup.status.failed')
    case 'canceled':
      return t('lookup.status.canceled')
    case 'refunded':
      return t('lookup.status.refunded')
    default:
      return t('lookup.status.unknown')
  }
})

const statusDescription = computed(() => {
  switch (order.value?.status) {
    case 'success':
      return order.value?.access_until
        ? t('lookup.success.active_until', { date: formatAbsoluteDateTime(order.value.access_until) })
        : t('lookup.success.default')
    case 'pending':
      return t('lookup.pending.description')
    case 'failed':
      return t('lookup.failed.description')
    case 'canceled':
      return t('lookup.canceled.description')
    case 'refunded':
      return t('lookup.refunded.description')
    default:
      return t('lookup.unknown_status')
  }
})

const amountText = computed(() => formatPaymentAmount(order.value?.amount || 0, localeTag.value))
const createdAtText = computed(() => {
  if (!order.value?.created_at) return t('status.unknown')
  return formatAbsoluteDateTime(order.value.created_at)
})
const paidAtText = computed(() => {
  if (!order.value?.paid_at) return t('status.unknown')
  return formatAbsoluteDateTime(order.value.paid_at)
})
const planLabel = computed(() =>
  order.value?.plan_code === 'donation' ? t('plan.donation') : t('plan.pro')
)
const orderStatusLabel = computed(() => {
  switch (order.value?.status) {
    case 'success':
      return t('status.success')
    case 'pending':
      return t('status.pending')
    case 'failed':
      return t('status.failed')
    case 'canceled':
      return t('status.canceled')
    case 'refunded':
      return t('status.refunded')
    default:
      return t('status.unknown')
  }
})

function normalizeInternalPath(value: string) {
  if (!value) return ''
  if (value.includes('://') || value.startsWith('//')) return ''
  return value.startsWith('/') ? value : `/${value}`
}

async function loadOrder(syncState = true) {
  if (!orderId.value || !token.value) {
    errorMessage.value = t('lookup.missing_data')
    loading.value = false
    return
  }

  errorMessage.value = ''

  if (loading.value) {
    loading.value = true
  } else {
    syncing.value = true
  }

  try {
    const res = await lookupPublicOrder({
      order_id: orderId.value,
      token: token.value,
      sync_state: syncState
    })

    order.value = res?.data?.order || null
  } catch (error: any) {
    errorMessage.value =
      String(error?.data?.message || error?.message || '').trim() || t('lookup.error_generic')
  } finally {
    loading.value = false
    syncing.value = false
  }
}

await loadOrder(true)
</script>

<template>
  <div class="space-y-6">
    <LabNavHeader :title="t('success.title')" />

    <div v-if="loading" variant="plain" :section-class="sectionClass">
      <p class="text-sm text-(--lab-text-secondary)">{{ t('lookup.retrying') }}</p>
    </div>

    <div
      v-else-if="errorMessage"
      :title="t('lookup.error_title')"
      variant="plain"
      section-class="bg-(--lab-danger)/10 p-4 "
      content-class="space-y-4"
    >
      <p class="text-sm leading-6 text-(--lab-text-primary)">
        {{ errorMessage }}
      </p>

      <div class="flex flex-wrap gap-3">
        <LabBaseButton variant="secondary" size="lg" :label="t('lookup.retry')" class="" @click="loadOrder(true)" />
        <LabBaseButton variant="plain" size="lg" :label="t('lookup.back')" class="" @click="navigateTo('/payments')" />
      </div>
    </div>

    <template v-else-if="order">
      <div variant="plain" :section-class="sectionClass" content-class="space-y-4">
        <div class="space-y-3">
          <p class="text-xs tracking-[0.22em] text-(--lab-success) uppercase">{{ t('success.title') }}</p>
          <h1 class="text-2xl font-semibold text-(--lab-text-primary)">{{ statusTitle }}</h1>
          <p class="max-w-2xl text-sm leading-6 text-(--lab-text-secondary)">{{ statusDescription }}</p>
        </div>

        <div class="grid gap-3 sm:grid-cols-[1fr_auto] sm:items-end">
          <div>
            <p class="text-xs tracking-[0.18em] text-(--lab-success) uppercase">{{ t('lookup.status') }}</p>
            <p class="mt-2 text-sm font-medium text-(--lab-text-primary)">{{ orderStatusLabel }}</p>
          </div>

          <div class="flex flex-wrap gap-3">
            <LabBaseButton
              v-if="order.status === 'success' && nextPath"
              variant="primary"
              size="lg"
              :label="t('lookup.continue')"
              class=""
              @click="navigateTo(nextPath)"
            />

            <LabBaseButton
              v-else-if="order.status === 'success'"
              variant="primary"
              size="lg"
              :label="t('lookup.account')"
              class=""
              @click="navigateTo('/account')"
            />

            <LabBaseButton
              v-if="order.status === 'pending'"
              variant="secondary"
              size="lg"
              :loading="syncing"
              :label="t('lookup.retry_short')"
              class=""
              @click="loadOrder(true)"
            />

            <LabBaseButton
              variant="plain"
              size="lg"
              :label="t('lookup.back')"
              class=""
              @click="navigateTo('/payments')"
            />
          </div>
        </div>
      </div>

      <div
        :title="t('index.latest_order_title')"
        variant="plain"
        :section-class="sectionClass"
        content-class="grid gap-x-8 gap-y-5 sm:grid-cols-2 xl:grid-cols-3"
      >
        <div>
          <p class="text-xs tracking-[0.18em] text-(--lab-success) uppercase">{{ t('lookup.order') }}</p>
          <p class="mt-2 text-sm wrap-break-word text-(--lab-text-primary)">{{ order.order_id }}</p>
        </div>
        <div>
          <p class="text-xs tracking-[0.18em] text-(--lab-success) uppercase">{{ t('lookup.status') }}</p>
          <p class="mt-2 text-sm text-(--lab-text-primary)">{{ orderStatusLabel }}</p>
        </div>
        <div>
          <p class="text-xs tracking-[0.18em] text-(--lab-success) uppercase">{{ t('lookup.amount') }}</p>
          <p class="mt-2 text-sm text-(--lab-text-primary)">{{ amountText }}</p>
        </div>
        <div>
          <p class="text-xs tracking-[0.18em] text-(--lab-success) uppercase">{{ t('lookup.plan') }}</p>
          <p class="mt-2 text-sm text-(--lab-text-primary)">{{ planLabel }}</p>
        </div>
        <div>
          <p class="text-xs tracking-[0.18em] text-(--lab-success) uppercase">{{ t('lookup.created') }}</p>
          <p class="mt-2 text-sm text-(--lab-text-primary)">{{ createdAtText }}</p>
        </div>
        <div>
          <p class="text-xs tracking-[0.18em] text-(--lab-success) uppercase">{{ t('lookup.paid') }}</p>
          <p class="mt-2 text-sm text-(--lab-text-primary)">{{ paidAtText }}</p>
        </div>
      </div>
    </template>
  </div>
</template>
