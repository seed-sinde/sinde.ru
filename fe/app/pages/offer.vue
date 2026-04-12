<script setup lang="ts">
  const { localeCode } = useInterfacePreferences()
  const page = computed(() => {
    const pages = {
      ru: {
        title: 'Условия предоставления доступа',
        description:
          'Условия предоставления пользователю доступа к функционалу, интерфейсам и цифровому контенту платформы sinde.ru',
        summaryTitle: 'Основное',
        conditionsTitle: 'Условия',
        facts: [
          { label: 'Формат', value: 'Доступ к функционалу сервиса' },
          { label: 'Модель', value: 'Простая (неисключительная) лицензия' },
          { label: 'Срок доступа', value: 'Определяется выбранным тарифом' },
          { label: 'Возврат', value: 'В случаях, предусмотренных законодательством РФ' },
          { label: 'Персональные данные', value: 'Политика конфиденциальности', href: '/terms' }
        ],
        sections: [
          {
            label: '1. Общие положения',
            value:
              'Исполнитель: ИП Воротников Дмитрий Владимирович. Настоящие условия регулируют предоставление пользователю доступа к функционалу, интерфейсам и цифровому контенту платформы sinde.'
          },
          {
            label: '2. Предмет',
            value:
              'Пользователю за плату предоставляется ограниченный доступ к функционалу сервиса и простая (неисключительная) лицензия на использование программного обеспечения и цифрового контента платформы в пределах выбранного тарифа. Исключительные права, право собственности на программный код, интерфейсы, данные и иные результаты интеллектуальной деятельности пользователю не передаются.'
          },
          {
            label: '3. Порядок оплаты',
            value:
              'Актуальная стоимость, срок доступа и доступные способы оплаты указываются в интерфейсе оформления заказа на момент оплаты. Оплата считается совершённой после подтверждения успешного списания денежных средств платёжным провайдером.'
          },
          {
            label: '4. Предоставление доступа',
            value:
              'После подтверждения оплаты пользователю предоставляется доступ к функционалу сервиса в объёме, соответствующем выбранному тарифу. Обязательство исполнителя по предоставлению доступа считается исполненным с момента фактического открытия доступа к соответствующему функционалу.'
          },
          {
            label: '5. Срок доступа',
            value:
              'Доступ предоставляется на срок, указанный в выбранном тарифе или заказе. По окончании оплаченного периода доступ прекращается, если иное прямо не указано в условиях конкретного тарифа.'
          },
          {
            label: '6. Возврат денежных средств',
            value:
              'После предоставления доступа денежные средства не возвращаются, за исключением случаев, прямо предусмотренных законодательством Российской Федерации, а также случаев, когда доступ не был предоставлен или услуга оказана ненадлежащим образом.'
          },
          {
            label: '7. Условия использования',
            value:
              'Доступ предоставляется только пользователю, оформившему заказ, и не подлежит передаче третьим лицам без согласия исполнителя. Исполнитель вправе ограничить или прекратить доступ при нарушении пользователем условий сервиса, требований законодательства либо по требованию уполномоченных органов.'
          },
          {
            label: '8. Ответственность',
            value:
              'Сервис предоставляется по модели «как есть». Исполнитель не несёт ответственность за косвенные убытки, упущенную выгоду, а также за невозможность использования сервиса, вызванную действиями третьих лиц, внешними интеграциями, каналами связи, оборудованием или программным обеспечением пользователя.'
          },
          {
            label: '9. Персональные данные',
            value:
              'Обработка персональных данных осуществляется в объёме, необходимом для регистрации, оплаты, предоставления доступа и поддержки пользователей, в соответствии с Политикой конфиденциальности.'
          },
          {
            label: '10. Акцепт',
            value:
              'Оплата заказа и, если предусмотрено интерфейсом, нажатие кнопки подтверждения означают полное и безоговорочное принятие пользователем настоящих условий.'
          }
        ]
      },
      en: {
        title: 'Access Terms',
        description:
          'Terms governing user access to the functionality, interfaces, and digital content of the sinde.ru platform.',
        summaryTitle: 'Overview',
        conditionsTitle: 'Terms',
        facts: [
          { label: 'Format', value: 'Access to the service functionality' },
          { label: 'Model', value: 'Simple non-exclusive license' },
          { label: 'Access period', value: 'Defined by the selected plan' },
          { label: 'Refunds', value: 'In cases required by Russian law' },
          { label: 'Personal data', value: 'Privacy policy', href: '/terms' }
        ],
        sections: [
          {
            label: '1. General provisions',
            value:
              'Provider: Dmitry Vladimirovich Vorotnikov, Individual Entrepreneur. These terms govern the granting of user access to the functionality, interfaces, and digital content of the sinde platform.'
          },
          {
            label: '2. Subject matter',
            value:
              'For a fee, the user receives limited access to the service functionality and a simple non-exclusive license to use the platform software and digital content within the selected plan. Exclusive rights, ownership of source code, interfaces, data, and other intellectual property results are not transferred to the user.'
          },
          {
            label: '3. Payment procedure',
            value:
              'The current price, access period, and available payment methods are shown in the checkout interface at the time of payment. Payment is deemed completed after the payment provider confirms successful withdrawal of funds.'
          },
          {
            label: '4. Grant of access',
            value:
              'After payment confirmation, the user is granted access to the service functionality in the scope of the selected plan. The provider’s obligation to grant access is considered fulfilled once access to the relevant functionality is actually opened.'
          },
          {
            label: '5. Access period',
            value:
              'Access is granted for the term specified in the selected plan or order. After the paid period ends, access is terminated unless otherwise explicitly stated in the conditions of a specific plan.'
          },
          {
            label: '6. Refund policy',
            value:
              'After access has been granted, funds are not refunded except where expressly required by the laws of the Russian Federation or where access was not provided or the service was rendered improperly.'
          },
          {
            label: '7. Terms of use',
            value:
              'Access is granted only to the user who placed the order and may not be transferred to third parties without the provider’s consent. The provider may restrict or terminate access if the user violates service terms, legal requirements, or upon request of authorized bodies.'
          },
          {
            label: '8. Liability',
            value:
              'The service is provided on an “as is” basis. The provider is not liable for indirect losses, lost profits, or inability to use the service caused by third parties, external integrations, communication channels, equipment, or the user’s software.'
          },
          {
            label: '9. Personal data',
            value:
              'Personal data is processed to the extent necessary for registration, payment, access provision, and user support in accordance with the Privacy Policy.'
          },
          {
            label: '10. Acceptance',
            value:
              'Payment for an order and, where required by the interface, clicking the confirmation button constitute the user’s full and unconditional acceptance of these terms.'
          }
        ]
      },
      ch: {
        title: '访问条款',
        description: '规范用户访问 sinde.ru 平台功能、界面和数字内容的条款。',
        summaryTitle: '摘要',
        conditionsTitle: '条款',
        facts: [
          { label: '形式', value: '访问服务功能' },
          { label: '模式', value: '普通（非独占）许可' },
          { label: '访问期限', value: '由所选套餐决定' },
          { label: '退款', value: '按俄罗斯联邦法律规定的情形处理' },
          { label: '个人数据', value: '隐私政策', href: '/terms' }
        ],
        sections: [
          { label: '1. 一般规定', value: '服务提供者：个体 предприниматель Дмитрий Владимирович Воротников。本条款规范向用户提供 sinde 平台功能、界面和数字内容访问权限的条件。' },
          { label: '2. 标的', value: '用户付费后可在所选套餐范围内获得对服务功能的有限访问权，以及使用平台软件和数字内容的普通（非独占）许可。源代码、界面、数据和其他知识产权成果的专有权及所有权不转让给用户。' },
          { label: '3. 支付程序', value: '支付时有效的价格、访问期限和可用支付方式，以下单界面显示的信息为准。支付提供商确认扣款成功后，视为付款完成。' },
          { label: '4. 提供访问', value: '付款确认后，用户将获得与所选套餐相对应范围的服务功能访问权限。自相关功能实际开放之时起，服务提供者的访问提供义务视为履行完毕。' },
          { label: '5. 访问期限', value: '访问权限在所选套餐或订单规定的期限内有效。付费期结束后，除非特定套餐另有明确规定，访问将终止。' },
          { label: '6. 退款', value: '访问权限提供后，款项通常不予退还，但俄罗斯联邦法律明确规定的情形，或未提供访问、服务提供不当的情形除外。' },
          { label: '7. 使用条件', value: '访问权限仅授予下单用户本人，未经服务提供者同意不得转让给第三方。如用户违反服务条款、法律要求，或应主管机关要求，服务提供者有权限制或终止访问。' },
          { label: '8. 责任', value: '服务按“现状”提供。对于因第三方、外部集成、通信渠道、设备或用户软件导致的间接损失、利润损失或无法使用服务，服务提供者不承担责任。' },
          { label: '9. 个人数据', value: '个人数据仅在注册、支付、提供访问和支持用户所必需的范围内处理，并遵循隐私政策。' },
          { label: '10. 接受', value: '支付订单以及在界面要求时点击确认按钮，均表示用户对本条款的完全且无保留接受。' }
        ]
      },
      jp: {
        title: 'アクセス条件',
        description: 'sinde.ru プラットフォームの機能、インターフェース、デジタルコンテンツへのアクセス条件です。',
        summaryTitle: '概要',
        conditionsTitle: '条件',
        facts: [
          { label: '形式', value: 'サービス機能へのアクセス' },
          { label: 'モデル', value: '通常の非独占ライセンス' },
          { label: 'アクセス期間', value: '選択したプランにより決定' },
          { label: '返金', value: 'ロシア連邦法で定められた場合' },
          { label: '個人データ', value: 'プライバシーポリシー', href: '/terms' }
        ],
        sections: [
          { label: '1. 総則', value: '提供者: 個人事業主ドミトリー・ウラジーミロヴィチ・ヴォロトニコフ。本条件は、sinde プラットフォームの機能、インターフェース、デジタルコンテンツへのアクセス提供を規定します。' },
          { label: '2. 対象', value: '利用者は料金を支払うことで、選択したプランの範囲内でサービス機能への限定的アクセスと、プラットフォームのソフトウェアおよびデジタルコンテンツを利用するための通常の非独占ライセンスを受けます。ソースコード、インターフェース、データその他の知的財産の排他的権利や所有権は利用者に移転しません。' },
          { label: '3. 支払い手続き', value: '有効な価格、アクセス期間、支払い方法は、支払い時点の注文画面に表示されます。決済プロバイダーによる引き落とし成功の確認後、支払い完了とみなされます。' },
          { label: '4. アクセス提供', value: '支払い確認後、選択したプランに応じた範囲でサービス機能へのアクセスが提供されます。関連機能へのアクセスが実際に開かれた時点で、提供者の義務は履行済みとみなされます。' },
          { label: '5. アクセス期間', value: 'アクセスは、選択したプランまたは注文に記載された期間に提供されます。別途明示がない限り、支払い済み期間の終了後にアクセスは終了します。' },
          { label: '6. 返金', value: 'アクセス提供後の返金は、ロシア連邦法で明示的に定められている場合、またはアクセスが提供されなかった場合やサービスが不適切に提供された場合を除き、行われません。' },
          { label: '7. 利用条件', value: 'アクセスは注文した利用者本人のみに提供され、提供者の同意なく第三者へ譲渡することはできません。利用者がサービス条件や法令に違反した場合、または権限ある機関の要求があった場合、提供者はアクセスを制限または終了できます。' },
          { label: '8. 責任', value: 'サービスは「現状のまま」提供されます。第三者、外部連携、通信回線、機器、または利用者のソフトウェアに起因する間接損害、逸失利益、サービス利用不能について、提供者は責任を負いません。' },
          { label: '9. 個人データ', value: '個人データは、登録、支払い、アクセス提供、ユーザーサポートに必要な範囲で、プライバシーポリシーに従って処理されます。' },
          { label: '10. 承諾', value: '注文の支払い、およびインターフェースで求められる場合の確認ボタンの押下は、本条件への完全かつ無条件の承諾を意味します。' }
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
          <div class="text-(--lab-text-muted) mb-3 text-sm uppercase tracking-[0.18em]">{{ page.summaryTitle }}</div>
          <div class="divide-y">
            <article v-for="item in page.facts" :key="item.label" class="py-3">
              <div class="min-w-0">
                <p class="text-(--lab-text-muted) text-xs uppercase tracking-[0.16em]">
                  {{ item.label }}
                </p>
                <NuxtLink
                  v-if="item.href"
                  :to="item.href"
                  class="lab-focus mt-2 inline-flex min-w-0 text-sm sm:text-base">
                  <span class="wrap-break-word">{{ item.value }}</span>
                </NuxtLink>
                <p v-else class="mt-2 text-sm leading-6 sm:text-base">{{ item.value }}</p>
              </div>
            </article>
          </div>
        </section>

        <section class="space-y-3">
          <div class="text-(--lab-text-muted) mb-3 text-sm uppercase tracking-[0.18em]">{{ page.conditionsTitle }}</div>
          <div class="divide-y">
            <article
              v-for="item in page.sections"
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
