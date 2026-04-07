<script setup lang="ts">
  const title = 'Оплата не завершена'

  usePageSeo({
    title,
    description: 'Проверка статуса неуспешного или отменённого платежа.'
  })

  const route = useRoute()

  const { lookupPublicOrder } = usePayments()
  const { formatAbsoluteDateTime } = useLocalizedDateTime()

  const loading = ref(true)
  const syncing = ref(false)
  const errorMessage = ref('')
  const order = ref<PaymentOrderView | null>(null)

  const orderId = computed(() => String(route.query.order_id || '').trim())
  const token = computed(() => String(route.query.token || '').trim())
  const nextPath = computed(() => normalizeInternalPath(String(route.query.next || '').trim()))

  const statusTitle = computed(() => {
    switch (order.value?.status) {
      case 'pending':
        return 'Платёж ещё обрабатывается'
      case 'failed':
        return 'Платёж завершился ошибкой'
      case 'canceled':
        return 'Платёж отменён'
      case 'refunded':
        return 'Платёж возвращён'
      case 'success':
        return 'Платёж всё же подтверждён'
      default:
        return 'Оплата не завершена'
    }
  })

  const statusDescription = computed(() => {
    switch (order.value?.status) {
      case 'pending':
        return 'Финальный статус ещё не получен. Можно повторить проверку.'
      case 'failed':
        return 'Банк не подтвердил платёж.'
      case 'canceled':
        return 'Пользователь или платёжный провайдер отменил операцию.'
      case 'refunded':
        return 'По заказу зарегистрирован возврат.'
      case 'success':
        return order.value?.access_until
          ? `Оплата подтверждена, доступ активен до ${formatAbsoluteDateTime(order.value.access_until)}.`
          : 'Оплата всё же завершилась успешно.'
      default:
        return 'Не удалось определить состояние заказа.'
    }
  })

  const amountText = computed(() => formatPrice(order.value?.amount || 0))

  function normalizeInternalPath(value: string) {
    if (!value) return ''
    if (value.includes('://') || value.startsWith('//')) return ''
    return value.startsWith('/') ? value : `/${value}`
  }

  function formatPrice(amount: number) {
    const rub = Math.floor(Number(amount || 0) / 100)
    return new Intl.NumberFormat('ru-RU').format(rub) + ' ₽'
  }

  async function loadOrder(syncState = true) {
    if (!orderId.value || !token.value) {
      errorMessage.value = 'Недостаточно данных для проверки платежа'
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
        String(error?.data?.message || error?.message || '').trim() || 'Не удалось получить состояние заказа'
    } finally {
      loading.value = false
      syncing.value = false
    }
  }

  await loadOrder(true)
</script>

<template>
  <div class="space-y-8">
    <LabNavHeader :title />

    <section v-if="loading" class="border bg-(--lab-bg-elevated) p-5 text-sm text-(--lab-text-muted)">
      Проверка статуса платежа…
    </section>

    <section v-else-if="errorMessage" class="space-y-4 border border-(--lab-danger)/30 bg-(--lab-danger)/8 p-5">
      <div class="text-base font-medium text-(--lab-text-primary)">Ошибка проверки платежа</div>
      <p class="text-sm text-(--lab-text-muted)">
        {{ errorMessage }}
      </p>

      <div class="flex flex-wrap gap-3">
        <button
          type="button"
          class="inline-flex min-h-11 items-center justify-center border border-(--lab-border-strong) px-4 text-sm font-medium text-(--lab-text-primary)"
          @click="loadOrder(true)">
          Повторить проверку
        </button>

        <NuxtLink
          to="/payments"
          class="inline-flex min-h-11 items-center justify-center border px-4 text-sm text-(--lab-text-primary)">
          Назад к оплате
        </NuxtLink>
      </div>
    </section>

    <template v-else-if="order">
      <section class="space-y-3 border bg-(--lab-bg-elevated) p-5">
        <h1 class="text-lg font-medium text-(--lab-text-primary)">
          {{ statusTitle }}
        </h1>
        <p class="text-sm text-(--lab-text-muted)">
          {{ statusDescription }}
        </p>

        <div class="flex flex-wrap gap-3 pt-2">
          <button
            v-if="order.status === 'pending'"
            type="button"
            class="inline-flex min-h-11 items-center justify-center border border-(--lab-border-strong) px-4 text-sm font-medium text-(--lab-text-primary) disabled:opacity-60"
            :disabled="syncing"
            @click="loadOrder(true)">
            {{ syncing ? 'Проверка…' : 'Проверить снова' }}
          </button>

          <NuxtLink
            v-if="order.status === 'success' && nextPath"
            :to="nextPath"
            class="inline-flex min-h-11 items-center justify-center border border-(--lab-border-strong) px-4 text-sm font-medium text-(--lab-text-primary)">
            Продолжить
          </NuxtLink>

          <NuxtLink
            to="/payments"
            class="inline-flex min-h-11 items-center justify-center border px-4 text-sm text-(--lab-text-primary)">
            Попробовать снова
          </NuxtLink>
        </div>
      </section>

      <section class="border bg-(--lab-bg-elevated) p-5">
        <dl class="grid gap-3 sm:grid-cols-2">
          <div class="space-y-1">
            <dt class="text-xs uppercase tracking-wide text-(--lab-text-muted)">Заказ</dt>
            <dd class="wrap-break-word text-sm text-(--lab-text-primary)">
              {{ order.order_id }}
            </dd>
          </div>

          <div class="space-y-1">
            <dt class="text-xs uppercase tracking-wide text-(--lab-text-muted)">Статус</dt>
            <dd class="text-sm text-(--lab-text-primary)">
              {{ order.status }}
            </dd>
          </div>

          <div class="space-y-1">
            <dt class="text-xs uppercase tracking-wide text-(--lab-text-muted)">Сумма</dt>
            <dd class="text-sm text-(--lab-text-primary)">
              {{ amountText }}
            </dd>
          </div>

          <div class="space-y-1">
            <dt class="text-xs uppercase tracking-wide text-(--lab-text-muted)">Создан</dt>
            <dd class="text-sm text-(--lab-text-primary)">
              {{ order.created_at ? formatAbsoluteDateTime(order.created_at) : '—' }}
            </dd>
          </div>
        </dl>
      </section>
    </template>
  </div>
</template>
