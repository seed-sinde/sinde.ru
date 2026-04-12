<script setup lang="ts">
  const { localeCode } = useInterfacePreferences()
  const page = computed(() => {
    const pages = {
      ru: {
        title: 'Компания',
        description: 'Контакты и реквизиты ИП Воротников Дмитрий Владимирович.',
        contactsTitle: 'Контакты',
        requisitesTitle: 'Реквизиты',
        contacts: [
          { label: 'Контактное лицо', value: 'Дмитрий В. Воротников' },
          { label: 'Email', value: 'bank@sinde.ru', href: 'mailto:bank@sinde.ru' },
          { label: 'Телефон', value: '+7 932 007-35-24', href: 'tel:+79320073524' }
        ],
        requisites: [
          { label: 'Название организации', value: 'ИНДИВИДУАЛЬНЫЙ ПРЕДПРИНИМАТЕЛЬ ВОРОТНИКОВ ДМИТРИЙ ВЛАДИМИРОВИЧ' },
          { label: 'Юридический адрес организации', value: '298517, РОССИЯ, РЕСП КРЫМ, Г АЛУШТА, УЛ ЯЛТИНСКАЯ, Д 3, КВ 31' },
          { label: 'ИНН', value: '910102272967' },
          { label: 'ОГРН/ОГРНИП', value: '326910000004858' },
          { label: 'Расчетный счет', value: '40802810700009527502' },
          { label: 'Банк', value: 'АО «ТБанк»' },
          { label: 'ИНН банка', value: '7710140679' },
          { label: 'БИК банка', value: '044525974' },
          { label: 'Корреспондентский счет банка', value: '30101810145250000974' },
          { label: 'Юридический адрес банка', value: '127287, г. Москва, ул. Хуторская 2-я, д. 38А, стр. 26' }
        ]
      },
      en: {
        title: 'Company',
        description: 'Contacts and legal details of Dmitry Vladimirovich Vorotnikov, Individual Entrepreneur.',
        contactsTitle: 'Contacts',
        requisitesTitle: 'Legal details',
        contacts: [
          { label: 'Contact person', value: 'Dmitry V. Vorotnikov' },
          { label: 'Email', value: 'bank@sinde.ru', href: 'mailto:bank@sinde.ru' },
          { label: 'Phone', value: '+7 932 007-35-24', href: 'tel:+79320073524' }
        ],
        requisites: [
          { label: 'Organization name', value: 'INDIVIDUAL ENTREPRENEUR VOROTNIKOV DMITRY VLADIMIROVICH' },
          { label: 'Registered address', value: '298517, RUSSIA, REPUBLIC OF CRIMEA, ALUSHTA, YALTINSKAYA ST., 3, APT. 31' },
          { label: 'TIN', value: '910102272967' },
          { label: 'Primary state registration number', value: '326910000004858' },
          { label: 'Settlement account', value: '40802810700009527502' },
          { label: 'Bank', value: 'TBank JSC' },
          { label: 'Bank TIN', value: '7710140679' },
          { label: 'Bank BIC', value: '044525974' },
          { label: 'Correspondent account', value: '30101810145250000974' },
          { label: 'Bank legal address', value: '127287, Moscow, 2nd Khutorskaya St., 38A, bld. 26' }
        ]
      },
      ch: {
        title: '公司',
        description: '个体 предприниматель Дмитрий Владимирович Воротников 的联系方式与法定信息。',
        contactsTitle: '联系方式',
        requisitesTitle: '公司信息',
        contacts: [
          { label: '联系人', value: 'Dmitry V. Vorotnikov' },
          { label: 'Email', value: 'bank@sinde.ru', href: 'mailto:bank@sinde.ru' },
          { label: '电话', value: '+7 932 007-35-24', href: 'tel:+79320073524' }
        ],
        requisites: [
          { label: '组织名称', value: '个体 предприниматель VOROTNIKOV DMITRY VLADIMIROVICH' },
          { label: '注册地址', value: '298517, 俄罗斯, 克里米亚共和国, 阿卢什塔, 雅尔京斯卡娅街 3 号, 31 室' },
          { label: '纳税号', value: '910102272967' },
          { label: '国家注册号', value: '326910000004858' },
          { label: '结算账户', value: '40802810700009527502' },
          { label: '银行', value: 'TBank JSC' },
          { label: '银行纳税号', value: '7710140679' },
          { label: '银行 BIC', value: '044525974' },
          { label: '对应账户', value: '30101810145250000974' },
          { label: '银行法定地址', value: '127287, 莫斯科, 第二胡托尔斯卡娅街 38A, 26 号楼' }
        ]
      },
      jp: {
        title: '会社',
        description: '個人事業主ドミトリー・ウラジーミロヴィチ・ヴォロトニコフの連絡先と法的情報です。',
        contactsTitle: '連絡先',
        requisitesTitle: '事業者情報',
        contacts: [
          { label: '担当者', value: 'Dmitry V. Vorotnikov' },
          { label: 'Email', value: 'bank@sinde.ru', href: 'mailto:bank@sinde.ru' },
          { label: '電話番号', value: '+7 932 007-35-24', href: 'tel:+79320073524' }
        ],
        requisites: [
          { label: '事業者名', value: '個人事業主 VOROTNIKOV DMITRY VLADIMIROVICH' },
          { label: '登録住所', value: '298517, ロシア, クリミア共和国, アルシタ, ヤルチンスカヤ通り 3, 31' },
          { label: '納税者番号', value: '910102272967' },
          { label: '登録番号', value: '326910000004858' },
          { label: '決済口座', value: '40802810700009527502' },
          { label: '銀行', value: 'TBank JSC' },
          { label: '銀行納税者番号', value: '7710140679' },
          { label: '銀行 BIC', value: '044525974' },
          { label: 'コルレス口座', value: '30101810145250000974' },
          { label: '銀行の法的住所', value: '127287, モスクワ, 2nd Khutorskaya St., 38A, bld. 26' }
        ]
      }
    } as const
    return pages[localeCode.value] || pages.ru
  })
  usePageSeo({
    title: () => page.value.title,
    description: () => page.value.description
  })
</script>
<template>
  <div>
    <LabNavHeader :title="page.title" />
    <section class="w-fit p-4">
      <div class="grid gap-6 xl:grid-cols-[minmax(0,18rem)_minmax(0,1fr)]">
        <section class="space-y-3">
          <div class="text-(--lab-text-muted) mb-3 text-sm uppercase tracking-[0.18em]">{{ page.contactsTitle }}</div>
          <div class="divide-y">
            <article v-for="item in page.contacts" :key="item.label" class="py-3">
              <div class="min-w-0">
                <p class="text-(--lab-text-muted) text-xs uppercase tracking-[0.16em]">
                  {{ item.label }}
                </p>
                <a
                  v-if="item.href"
                  :href="item.href"
                  class="lab-focus mt-2 inline-flex min-w-0 text-sm sm:text-base">
                  <span class="wrap-break-word">{{ item.value }}</span>
                </a>
                <p v-else class="mt-2 text-sm sm:text-base">{{ item.value }}</p>
              </div>
            </article>
          </div>
        </section>
        <section class="space-y-3">
          <div class="text-(--lab-text-muted) mb-3 text-sm uppercase tracking-[0.18em]">{{ page.requisitesTitle }}</div>
          <div class="divide-y">
            <article
              v-for="item in page.requisites"
              :key="item.label"
              class="grid gap-2 px-3 py-2 sm:grid-cols-[12rem_minmax(0,1fr)] sm:gap-3">
              <p class="text-(--lab-text-muted) text-xs uppercase tracking-[0.16em]">
                {{ item.label }}
              </p>
              <p class="wrap-break-word text-sm leading-6 sm:text-[0.95rem]">{{ item.value }}</p>
            </article>
          </div>
        </section>
      </div>
    </section>
  </div>
</template>
