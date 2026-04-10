<script setup lang="ts">
  const { t, localeTag } = useInterfacePreferences()
  const title = t('payments.success.seo_title')

  usePageSeo({
    title,
    description: t('payments.success.seo_description')
  })

  const route = useRoute()

  const { lookupPublicOrder } = usePayments()
  const { formatAbsoluteDateTime } = useLocalizedDateTime()

  const loading = ref(true)
  const syncing = ref(false)
  const errorMessage = ref('')
  const order = ref<PaymentOrderView | null>(null)

  const sectionClass = 'bg-(--lab-bg-elevated) p-5 sm:p-6'

  const orderId = computed(() => String(route.query.order_id || '').trim())
  const token = computed(() => String(route.query.token || '').trim())
  const nextPath = computed(() => normalizeInternalPath(String(route.query.next || '').trim()))

  const statusTitle = computed(() => {
    switch (order.value?.status) {
      case 'success':
        return t('payments.lookup.status.success')
      case 'pending':
        return t('payments.lookup.status.pending')
      case 'failed':
        return t('payments.lookup.status.failed')
      case 'canceled':
        return t('payments.lookup.status.canceled')
      case 'refunded':
        return t('payments.lookup.status.refunded')
      default:
        return t('payments.lookup.status.unknown')
    }
  })

  const statusDescription = computed(() => {
    switch (order.value?.status) {
      case 'success':
        return order.value?.access_until
          ? t('payments.lookup.success.active_until', { date: formatAbsoluteDateTime(order.value.access_until) })
          : t('payments.lookup.success.default')
      case 'pending':
        return t('payments.lookup.pending.description')
      case 'failed':
        return t('payments.lookup.failed.description')
      case 'canceled':
        return t('payments.lookup.canceled.description')
      case 'refunded':
        return t('payments.lookup.refunded.description')
      default:
        return t('payments.lookup.unknown_status')
    }
  })

  const amountText = computed(() => formatPaymentAmount(order.value?.amount || 0, localeTag.value))
  const createdAtText = computed(() => {
    if (!order.value?.created_at) return t('payments.status.unknown')
    return formatAbsoluteDateTime(order.value.created_at)
  })
  const paidAtText = computed(() => {
    if (!order.value?.paid_at) return t('payments.status.unknown')
    return formatAbsoluteDateTime(order.value.paid_at)
  })
  const planLabel = computed(() => (order.value?.plan_code === 'donation' ? t('payments.plan.donation') : t('payments.plan.pro')))
  const orderStatusLabel = computed(() => {
    switch (order.value?.status) {
      case 'success':
        return t('payments.status.success')
      case 'pending':
        return t('payments.status.pending')
      case 'failed':
        return t('payments.status.failed')
      case 'canceled':
        return t('payments.status.canceled')
      case 'refunded':
        return t('payments.status.refunded')
      default:
        return t('payments.status.unknown')
    }
  })

  function normalizeInternalPath(value: string) {
    if (!value) return ''
    if (value.includes('://') || value.startsWith('//')) return ''
    return value.startsWith('/') ? value : `/${value}`
  }

  async function loadOrder(syncState = true) {
    if (!orderId.value || !token.value) {
      errorMessage.value = t('payments.lookup.missing_data')
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
      errorMessage.value = String(error?.data?.message || error?.message || '').trim() || t('payments.lookup.error_generic')
    } finally {
      loading.value = false
      syncing.value = false
    }
  }

  await loadOrder(true)
</script>

<template>
  <div class="space-y-6">
    <LabNavHeader :title="t('payments.success.title')" />

    <LabBaseSection v-if="loading" variant="plain" :section-class="sectionClass">
      <p class="text-sm text-(--lab-text-secondary)">{{ t('payments.lookup.retrying') }}</p>
    </LabBaseSection>

    <LabBaseSection
      v-else-if="errorMessage"
      :title="t('payments.lookup.error_title')"
      variant="plain"
      section-class="bg-(--lab-danger)/10 p-5 sm:p-6"
      content-class="space-y-4">
      <p class="text-sm leading-6 text-(--lab-text-primary)">
        {{ errorMessage }}
      </p>

      <div class="flex flex-wrap gap-3">
        <LabBaseButton
          variant="secondary"
          size="lg"
          :label="t('payments.lookup.retry')"
          :button-style="{ borderWidth: '0px' }"
          button-class="focus-visible:ring-2 focus-visible:ring-(--lab-accent)"
          @click="loadOrder(true)" />
        <LabBaseButton
          variant="plain"
          size="lg"
          :label="t('payments.lookup.back')"
          :button-style="{ borderWidth: '0px' }"
          button-class="focus-visible:ring-2 focus-visible:ring-(--lab-accent)"
          @click="navigateTo('/payments')" />
      </div>
    </LabBaseSection>

    <template v-else-if="order">
      <LabBaseSection variant="plain" :section-class="sectionClass" content-class="space-y-5">
        <div class="space-y-3">
          <p class="text-xs uppercase tracking-[0.22em] text-(--lab-success)">{{ t('payments.success.title') }}</p>
          <h1 class="text-2xl font-semibold text-(--lab-text-primary)">{{ statusTitle }}</h1>
          <p class="max-w-2xl text-sm leading-6 text-(--lab-text-secondary)">{{ statusDescription }}</p>
        </div>

        <div class="grid gap-3 sm:grid-cols-[1fr_auto] sm:items-end">
          <div>
            <p class="text-xs uppercase tracking-[0.18em] text-(--lab-success)">{{ t('payments.lookup.status') }}</p>
            <p class="mt-2 text-sm font-medium text-(--lab-text-primary)">{{ orderStatusLabel }}</p>
          </div>

          <div class="flex flex-wrap gap-3">
            <LabBaseButton
              v-if="order.status === 'success' && nextPath"
              variant="primary"
              size="lg"
              :label="t('payments.lookup.continue')"
              :button-style="{ borderWidth: '0px' }"
              button-class="focus-visible:ring-2 focus-visible:ring-(--lab-accent)"
              @click="navigateTo(nextPath)" />

            <LabBaseButton
              v-else-if="order.status === 'success'"
              variant="primary"
              size="lg"
              :label="t('payments.lookup.account')"
              :button-style="{ borderWidth: '0px' }"
              button-class="focus-visible:ring-2 focus-visible:ring-(--lab-accent)"
              @click="navigateTo('/account')" />

            <LabBaseButton
              v-if="order.status === 'pending'"
              variant="secondary"
              size="lg"
              :loading="syncing"
              :loading-label="t('payments.lookup.retrying')"
              :label="t('payments.lookup.retry_short')"
              :button-style="{ borderWidth: '0px' }"
              button-class="focus-visible:ring-2 focus-visible:ring-(--lab-accent)"
              @click="loadOrder(true)" />

            <LabBaseButton
              variant="plain"
              size="lg"
              :label="t('payments.lookup.back')"
              :button-style="{ borderWidth: '0px' }"
              button-class="focus-visible:ring-2 focus-visible:ring-(--lab-accent)"
              @click="navigateTo('/payments')" />
          </div>
        </div>
      </LabBaseSection>

      <LabBaseSection
        :title="t('payments.index.latest_order_title')"
        variant="plain"
        :section-class="sectionClass"
        content-class="grid gap-x-8 gap-y-5 sm:grid-cols-2 xl:grid-cols-3">
        <div>
          <p class="text-xs uppercase tracking-[0.18em] text-(--lab-success)">{{ t('payments.lookup.order') }}</p>
          <p class="mt-2 wrap-break-word text-sm text-(--lab-text-primary)">{{ order.order_id }}</p>
        </div>

        <div>
          <p class="text-xs uppercase tracking-[0.18em] text-(--lab-success)">{{ t('payments.lookup.status') }}</p>
          <p class="mt-2 text-sm text-(--lab-text-primary)">{{ orderStatusLabel }}</p>
        </div>

        <div>
          <p class="text-xs uppercase tracking-[0.18em] text-(--lab-success)">{{ t('payments.lookup.amount') }}</p>
          <p class="mt-2 text-sm text-(--lab-text-primary)">{{ amountText }}</p>
        </div>

        <div>
          <p class="text-xs uppercase tracking-[0.18em] text-(--lab-success)">{{ t('payments.lookup.plan') }}</p>
          <p class="mt-2 text-sm text-(--lab-text-primary)">{{ planLabel }}</p>
        </div>

        <div>
          <p class="text-xs uppercase tracking-[0.18em] text-(--lab-success)">{{ t('payments.lookup.created') }}</p>
          <p class="mt-2 text-sm text-(--lab-text-primary)">{{ createdAtText }}</p>
        </div>

        <div>
          <p class="text-xs uppercase tracking-[0.18em] text-(--lab-success)">{{ t('payments.lookup.paid') }}</p>
          <p class="mt-2 text-sm text-(--lab-text-primary)">{{ paidAtText }}</p>
        </div>
      </LabBaseSection>
    </template>
  </div>
</template>
