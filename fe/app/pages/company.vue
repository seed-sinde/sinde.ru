<script setup lang="ts">
  const title = 'Компания'
  usePageSeo({
    title,
    description: 'Контакты и реквизиты ИП Воротников Дмитрий Владимирович.'
  })
  const contacts: Array<{
    label: string
    value: string
    href?: string
  }> = [
    {
      label: 'Контактное лицо',
      value: 'Дмитрий В. Воротников'
    },
    {
      label: 'Email',
      value: 'bank@sinde.ru',
      href: 'mailto:bank@sinde.ru'
    },
    {
      label: 'Телефон',
      value: '+7 932 007-35-24',
      href: 'tel:+79320073524'
    }
  ]
  const requisites: Array<{
    label: string
    value: string
  }> = [
    {
      label: 'Название организации',
      value: 'ИНДИВИДУАЛЬНЫЙ ПРЕДПРИНИМАТЕЛЬ ВОРОТНИКОВ ДМИТРИЙ ВЛАДИМИРОВИЧ'
    },
    {
      label: 'Юридический адрес организации',
      value: '298517, РОССИЯ, РЕСП КРЫМ, Г АЛУШТА, УЛ ЯЛТИНСКАЯ, Д 3, КВ 31'
    },
    {
      label: 'ИНН',
      value: '910102272967'
    },
    {
      label: 'ОГРН/ОГРНИП',
      value: '326910000004858'
    },
    {
      label: 'Расчетный счет',
      value: '40802810700009527502'
    },
    {
      label: 'Банк',
      value: 'АО «ТБанк»'
    },
    {
      label: 'ИНН банка',
      value: '7710140679'
    },
    {
      label: 'БИК банка',
      value: '044525974'
    },
    {
      label: 'Корреспондентский счет банка',
      value: '30101810145250000974'
    },
    {
      label: 'Юридический адрес банка',
      value: '127287, г. Москва, ул. Хуторская 2-я, д. 38А, стр. 26'
    }
  ]
</script>
<template>
  <div>
    <LabNavHeader :title />
    <section class="w-fit p-4">
      <div class="grid gap-6 xl:grid-cols-[minmax(0,18rem)_minmax(0,1fr)]">
        <section class="space-y-3">
          <div class="text-(--lab-text-muted) mb-3 text-sm uppercase tracking-[0.18em]">Контакты</div>
          <div class="divide-y">
            <article v-for="item in contacts" :key="item.label" class="py-3">
              <div class="min-w-0">
                <p class="text-(--lab-text-muted) text-xs uppercase tracking-[0.16em]">
                  {{ item.label }}
                </p>
                <a
                  v-if="item.href"
                  :href="item.href"
                  class="text-(--lab-text-primary) mt-2 inline-flex min-w-0 text-sm ring-1 ring-transparent transition focus-visible:outline-none focus-visible:ring-2 sm:text-base">
                  <span class="wrap-break-word">{{ item.value }}</span>
                </a>
                <p v-else class="text-(--lab-text-primary) mt-2 text-sm sm:text-base">
                  {{ item.value }}
                </p>
              </div>
            </article>
          </div>
        </section>
        <section class="space-y-3">
          <div class="text-(--lab-text-muted) mb-3 text-sm uppercase tracking-[0.18em]">Реквизиты</div>
          <div class="divide-y">
            <article
              v-for="item in requisites"
              :key="item.label"
              class="grid gap-2 px-3 py-2 sm:grid-cols-[12rem_minmax(0,1fr)] sm:gap-3">
              <p class="text-(--lab-text-muted) text-xs uppercase tracking-[0.16em]">
                {{ item.label }}
              </p>
              <p class="text-(--lab-text-primary) wrap-break-word text-sm leading-6 sm:text-[0.95rem]">
                {{ item.value }}
              </p>
            </article>
          </div>
        </section>
      </div>
    </section>
  </div>
</template>
