export type PaymentOrderView = {
  order_id: string
  user_id: string
  user_email: string
  user_display_name: string
  plan_code: 'pro' | 'donation' | string
  base_amount: number
  amount: number
  tip_amount: number
  currency: string
  subscription_type: 'one_time' | 'recurring' | string
  status: 'pending' | 'success' | 'failed' | 'refunded' | 'canceled' | string
  provider: string
  provider_status: string
  provider_payment_id: string
  provider_error_code: string
  provider_message: string
  provider_fee_percent: number
  fee_amount: number
  net_amount: number
  access_from?: string | null
  access_until?: string | null
  return_to: string
  payment_url: string
  last_checked_at?: string | null
  notified_at?: string | null
  paid_at?: string | null
  failed_at?: string | null
  refunded_at?: string | null
  created_at: string
  updated_at: string
}
export type PaymentAccessSummary = {
  has_active_access: boolean
  plan_code: string
  amount: number
  tip_amount: number
  currency: string
  access_from?: string | null
  access_until?: string | null
  order_id?: string | null
  latest_order?: PaymentOrderView | null
}
export type PaymentCreateOrderResult = {
  order: PaymentOrderView
  payment_url: string
}
export type PaymentPublicLookupResult = {
  order: PaymentOrderView
}
export type PaymentAdminOrdersSummary = {
  orders_total: number
  orders_success: number
  orders_pending: number
  orders_failed: number
  orders_refunded: number
  paid_users_total: number
  patron_users_total: number
  active_access_users: number
  gross_revenue: number
  net_revenue: number
  tip_revenue: number
  mrr: number
  churn_rate: number
  patron_share: number
  last_successful_paid?: string | null
}
