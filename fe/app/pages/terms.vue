<script setup lang="ts">
  const { localeCode } = useInterfacePreferences()
  const page = computed(() => {
    const pages = {
      ru: {
        title: 'Кодекс участника',
        description: 'Соглашение между участником и платформой sinde.',
        about: 'О документе',
        purposeLabel: 'Назначение',
        purposeValue:
          'Документ фиксирует базовые принципы использования сервиса, доступа к функционалу, оплаты и обращения с данными участника.',
        scopeLabel: 'Применение',
        scopeValue: 'Положения действуют для регистрации, оформления заказа, оплаты и дальнейшего использования платформы sinde.',
        principlesTitle: 'Принципы',
        principles: [
          {
            label: 'Контроль',
            value:
              'Ты сам решаешь, какие данные добавлять в систему, если иное прямо не требуется для регистрации, оплаты, безопасности или работы отдельных функций сервиса. Ответственность за содержание передаваемых тобой данных несёшь ты.'
          },
          {
            label: 'Использование данных',
            value:
              'Данные используются внутри сервиса для хранения, отображения, установления связей между сущностями, работы функций и развития системы. Они не используются вне этих целей без отдельного основания.'
          },
          {
            label: 'Минимизация',
            value:
              'Сервис стремится не собирать лишнего и обрабатывает только те данные, которые ты явно вводишь, передаёшь или которые необходимы для работы сервиса, оплаты, безопасности и поддержки.'
          },
          {
            label: 'Актуальность',
            value:
              'Ты можешь обновлять, изменять или удалять доступные тебе данные через предусмотренный функционал сервиса. Часть данных может сохраняться в объёме, необходимом для работы системы, исполнения обязательств и соблюдения требований законодательства.'
          },
          {
            label: 'Доступ',
            value:
              'Услуга сервиса заключается в предоставлении доступа к функционалу платформы. Объём доступа зависит от выбранного режима, тарифа, статуса аккаунта и твоих действий внутри системы.'
          },
          {
            label: 'Оказание услуги',
            value:
              'Услуга считается оказанной с момента предоставления доступа к соответствующему функционалу сервиса. Если иное не указано отдельно, доступ к предусмотренному объёму функций предоставляется сразу после выполнения условий подключения.'
          },
          {
            label: 'Возврат',
            value:
              'После предоставления доступа возврат денежных средств, как правило, не осуществляется, кроме случаев ненадлежащего качества услуги или иных случаев, предусмотренных законодательством и документами сервиса.'
          },
          {
            label: 'Безопасность',
            value:
              'В сервисе применяются технические и организационные меры защиты данных. При этом абсолютная защищённость не гарантируется, а пользователь также должен соблюдать базовые меры безопасности со своей стороны.'
          },
          {
            label: 'Прозрачность',
            value:
              'Основные условия использования сервиса, доступа, оплаты и обработки данных публикуются в документах сервиса. Существенные изменения условий подлежат отражению в актуальной версии таких документов.'
          },
          {
            label: 'Согласие',
            value:
              'Регистрация, оформление заказа, оплата и использование сервиса означают согласие пользователя с условиями доступа, оказания услуги, возврата средств и иными правилами, опубликованными в документах сервиса.'
          }
        ]
      },
      en: {
        title: 'Participant Code',
        description: 'Agreement between the participant and the sinde platform.',
        about: 'About',
        purposeLabel: 'Purpose',
        purposeValue:
          'This document defines the core principles of using the service, accessing its functionality, making payments, and handling participant data.',
        scopeLabel: 'Scope',
        scopeValue: 'These provisions apply to registration, checkout, payment, and the further use of the sinde platform.',
        principlesTitle: 'Principles',
        principles: [
          {
            label: 'Control',
            value:
              'You decide which data to add to the system unless specific data is explicitly required for registration, payment, security, or certain service features. You remain responsible for the content of the data you provide.'
          },
          {
            label: 'Use of data',
            value:
              'Data is used within the service to store, display, link entities, operate functions, and improve the system. It is not used outside these purposes without a separate legal basis.'
          },
          {
            label: 'Minimization',
            value:
              'The service aims not to collect unnecessary data and processes only the information you explicitly enter or provide, or that is required for the service, payment, security, and support.'
          },
          {
            label: 'Accuracy',
            value:
              'You may update, change, or delete the data available to you through the service features provided. Some data may still be retained to the extent necessary for system operation, contractual obligations, and legal compliance.'
          },
          {
            label: 'Access',
            value:
              'The service consists of granting access to the platform functionality. The scope of access depends on the selected mode, plan, account status, and your actions within the system.'
          },
          {
            label: 'Service delivery',
            value:
              'The service is considered delivered once access to the relevant functionality is provided. Unless stated otherwise, access to the included functions is granted immediately after the connection requirements are met.'
          },
          {
            label: 'Refunds',
            value:
              'After access has been granted, refunds are generally not provided except in cases of improper service quality or in other cases required by law or by the service documents.'
          },
          {
            label: 'Security',
            value:
              'The service uses technical and organizational safeguards to protect data. Absolute protection is not guaranteed, and the user must also follow basic security measures on their side.'
          },
          {
            label: 'Transparency',
            value:
              'The key terms of service use, access, payment, and data processing are published in the service documents. Material changes must be reflected in the current versions of those documents.'
          },
          {
            label: 'Consent',
            value:
              'Registration, checkout, payment, and use of the service mean the user accepts the terms of access, service delivery, refunds, and other rules published in the service documents.'
          }
        ]
      },
      ch: {
        title: '参与者守则',
        description: '参与者与 sinde 平台之间的协议。',
        about: '关于文件',
        purposeLabel: '目的',
        purposeValue: '本文件规定了服务使用、功能访问、支付以及参与者数据处理的基本原则。',
        scopeLabel: '适用范围',
        scopeValue: '这些条款适用于 sinde 平台的注册、下单、付款以及后续使用。',
        principlesTitle: '原则',
        principles: [
          { label: '控制', value: '除非注册、支付、安全或特定服务功能明确要求，否则由你自行决定向系统添加哪些数据。你对所提交数据的内容负责。' },
          { label: '数据使用', value: '数据仅在服务内部用于存储、展示、建立实体关联、支持功能运行以及系统发展。没有单独依据时，不会用于这些目的之外。' },
          { label: '最小化', value: '服务力求不收集多余信息，只处理你明确输入、提交或服务、支付、安全和支持所必需的数据。' },
          { label: '及时性', value: '你可以通过服务提供的功能更新、修改或删除你可访问的数据。部分数据可能因系统运行、履约和法律要求而保留必要范围。' },
          { label: '访问', value: '服务本质上是向你提供平台功能访问权限。访问范围取决于所选模式、套餐、账户状态以及你在系统中的操作。' },
          { label: '服务交付', value: '当相关服务功能的访问权限被提供时，即视为服务已交付。除非另有说明，满足接入条件后会立即开放相应功能。' },
          { label: '退款', value: '在已提供访问权限后，通常不予退款，除非存在服务质量问题或法律及服务文件另有规定。' },
          { label: '安全', value: '服务采用技术和组织措施保护数据，但不保证绝对安全，用户也应自行遵守基本安全要求。' },
          { label: '透明度', value: '服务使用、访问、支付和数据处理的主要条件均发布于服务文件中。重要变更应反映在这些文件的最新版本中。' },
          { label: '同意', value: '注册、下单、付款和使用服务即表示用户同意服务文件中发布的访问、服务交付、退款及其他规则。' }
        ]
      },
      jp: {
        title: '参加者コード',
        description: '参加者と sinde プラットフォームとの間の合意です。',
        about: '文書について',
        purposeLabel: '目的',
        purposeValue: 'この文書は、サービス利用、機能へのアクセス、支払い、および参加者データの取り扱いに関する基本原則を定めます。',
        scopeLabel: '適用範囲',
        scopeValue: '本規定は、sinde プラットフォームの登録、注文、支払い、およびその後の利用に適用されます。',
        principlesTitle: '原則',
        principles: [
          { label: '管理', value: '登録、支払い、安全性、または特定機能のために明示的に必要とされる場合を除き、どのデータをシステムに追加するかはあなた自身が決定します。送信するデータ内容の責任はあなたにあります。' },
          { label: 'データ利用', value: 'データは、保存、表示、エンティティ間の関連付け、機能動作、およびシステム改善のためにサービス内部で使用されます。別の法的根拠がない限り、これら以外の目的には使用されません。' },
          { label: '最小化', value: 'サービスは不要なデータを収集しないよう努め、あなたが明示的に入力・提供した情報、またはサービス、支払い、安全性、サポートに必要な情報のみを処理します。' },
          { label: '最新性', value: '提供された機能を通じて、あなたは自分が利用可能なデータを更新、変更、削除できます。ただし、システム運用、義務履行、法令遵守のために必要な範囲で一部データが保持される場合があります。' },
          { label: 'アクセス', value: 'サービスは、プラットフォーム機能へのアクセスを提供することにあります。アクセス範囲は、選択したモード、プラン、アカウント状態、およびシステム内でのあなたの操作に依存します。' },
          { label: '提供完了', value: '該当機能へのアクセスが提供された時点で、サービスは提供済みとみなされます。別途定めがない限り、接続条件が満たされ次第、対象機能へのアクセスが直ちに提供されます。' },
          { label: '返金', value: 'アクセス提供後の返金は、原則として行われません。ただし、サービス品質に問題がある場合や、法令またはサービス文書で別途定めがある場合を除きます。' },
          { label: '安全性', value: 'サービスではデータ保護のために技術的・組織的対策を講じていますが、絶対的な安全性は保証されません。利用者側でも基本的な安全対策を守る必要があります。' },
          { label: '透明性', value: 'サービス利用、アクセス、支払い、データ処理に関する主要条件はサービス文書に掲載されます。重要な変更は、その最新文書に反映される必要があります。' },
          { label: '同意', value: '登録、注文、支払い、およびサービス利用は、アクセス条件、サービス提供、返金、およびサービス文書に掲載されたその他の規則に同意したものとみなされます。' }
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
          <div class="text-(--lab-text-muted) mb-3 text-sm uppercase tracking-[0.18em]">{{ page.about }}</div>
          <div class="divide-y">
            <article class="py-3">
              <div class="min-w-0">
                <p class="text-(--lab-text-muted) text-xs uppercase tracking-[0.16em]">{{ page.purposeLabel }}</p>
                <p class="mt-2 text-sm leading-6 sm:text-base">{{ page.purposeValue }}</p>
              </div>
            </article>
            <article class="py-3">
              <div class="min-w-0">
                <p class="text-(--lab-text-muted) text-xs uppercase tracking-[0.16em]">{{ page.scopeLabel }}</p>
                <p class="mt-2 text-sm leading-6 sm:text-base">{{ page.scopeValue }}</p>
              </div>
            </article>
          </div>
        </section>

        <section class="space-y-3">
          <div class="text-(--lab-text-muted) mb-3 text-sm uppercase tracking-[0.18em]">{{ page.principlesTitle }}</div>
          <div class="divide-y">
            <article
              v-for="item in page.principles"
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
