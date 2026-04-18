<script setup lang="ts">
type IpInfo = {
  ip: string
  success: boolean
  type: string
  continent: string
  continent_code: string
  country: string
  country_code: string
  region: string
  region_code: string
  city: string
  latitude: number
  longitude: number
  is_eu: boolean
  postal: string
  calling_code: string
  capital: string
  borders: string
  connection: {
    asn: number
    org: string
    isp: string
    domain: string
  }
  timezone: {
    id: string
    abbr: string
    is_dst: boolean
    offset: number
    utc: string
  }
}
type ApiRes<T> = {
  ok: boolean
  data: T
  message?: string
}
const useApiJson = <T,>(p: string, o?: NonNullable<Parameters<ReturnType<typeof useAPI>['json']>[1]>) =>
  useAPI().json<T>(p, o)
const { data, pending, error } = await useAsyncData(
  'admin-ip-info',
  () => useApiJson<ApiRes<IpInfo>>('/auth/admin/ip'),
  { server: true, lazy: false }
)
const d = computed(() => data.value?.data || null)
const em = computed(() => {
  const s = error.value as { data?: { message?: string } | string; statusMessage?: string; message?: string } | null
  return s?.data && typeof s.data === 'object' && 'message' in s.data
    ? String(s.data.message || 'failed to load ip info')
    : typeof s?.data === 'string'
      ? s.data
      : s?.statusMessage || s?.message || 'failed to load ip info'
})
</script>

<template>
  <div class="space-y-1 text-xs">
    <div v-if="pending">loading...</div>
    <div v-else-if="error">{{ em }}</div>
    <template v-else-if="d">
      <div v-if="d.ip || d.country_code || d.city">
        <span v-if="d.ip">{{ d.ip }}</span>
        <span v-if="d.country_code">· {{ d.country_code }}</span>
        <span v-if="d.city">· {{ d.city }}</span>
      </div>
      <div v-if="d.region || d.postal || d.timezone?.utc">
        <span v-if="d.region">{{ d.region }}</span>
        <span v-if="d.postal">· {{ d.postal }}</span>
        <span v-if="d.timezone?.utc">· {{ d.timezone.utc }}</span>
      </div>
      <div v-if="d.connection?.isp || d.connection?.asn">
        <span v-if="d.connection?.isp">{{ d.connection.isp }}</span>
        <span v-if="d.connection?.asn">· AS{{ d.connection.asn }}</span>
      </div>
    </template>
  </div>
</template>
