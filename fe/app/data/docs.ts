export const docsPages: Record<InterfaceLocaleCode, DocsPageContent> = {
  ru: {
    tabsTitle: 'Документы',
    documents: {
      company: {
        title: 'Реквизиты',
        description: 'Реквизиты, контакты и порядок связи с продавцом и исполнителем платформы sinde.',
        introTitle: 'Контакты',
        introItems: [
          { label: 'Исполнитель', value: 'ИП Воротников Дмитрий Владимирович' },
          { label: 'Email', value: 'bank@sinde.ru', href: 'mailto:bank@sinde.ru' },
          { label: 'Телефон', value: '+7 932 007-35-24', href: 'tel:+79320073524' },
          { label: 'Ответ на обращения', value: 'Обычно отвечаем в течение суток.' }
        ],
        sectionsTitle: 'Реквизиты',
        sections: [
          {
            label: 'Название',
            value: 'ИНДИВИДУАЛЬНЫЙ ПРЕДПРИНИМАТЕЛЬ ВОРОТНИКОВ ДМИТРИЙ ВЛАДИМИРОВИЧ'
          },
          {
            label: 'Юридический адрес',
            value: '298517, Россия, Республика Крым, г. Алушта, ул. Ялтинская, д. 3, кв. 31'
          },
          { label: 'ИНН', value: '910102272967' },
          { label: 'ОГРНИП', value: '326910000004858' },
          { label: 'Расчётный счёт', value: '40802810700009527502' },
          { label: 'Банк', value: 'АО «ТБанк»' },
          { label: 'БИК', value: '044525974' },
          { label: 'Корреспондентский счёт', value: '30101810145250000974' }
        ]
      },
      offer: {
        title: 'Оферта',
        description: 'Условия покупки цифрового доступа, оплаты, акцепта, предоставления доступа и возврата средств.',
        introTitle: 'Основное',
        introItems: [
          { label: 'Формат', value: 'Покупка цифрового доступа к функционалу платформы' },
          { label: 'Срок доступа', value: 'Определяется выбранным тарифом или заказом' },
          { label: 'Возврат', value: 'По применимому законодательству и условиям документа' },
          { label: 'Персональные данные', value: 'Политика конфиденциальности', href: '/docs/privacy' }
        ],
        sectionsTitle: 'Условия',
        sections: [
          {
            label: '1. Общие положения',
            value:
              'Исполнитель публикует настоящую оферту для покупки цифрового доступа к платформе sinde. Оплата заказа и подтверждение в интерфейсе, если оно требуется, означают полный и безоговорочный акцепт.'
          },
          {
            label: '2. Предмет',
            value:
              'Пользователю предоставляется ограниченный доступ к функционалу, интерфейсам и цифровому контенту платформы в объёме выбранного тарифа. Исключительные права на программное обеспечение, данные и интерфейсы пользователю не передаются.'
          },
          {
            label: '3. Цена и оплата',
            value:
              'Актуальная стоимость, срок доступа и доступные способы оплаты указываются в интерфейсе оформления заказа на момент покупки. Обязательство по оплате считается исполненным после подтверждения успешного списания денежных средств платёжным провайдером.'
          },
          {
            label: '4. Предоставление доступа',
            value:
              'Доступ открывается после подтверждения оплаты в объёме выбранного тарифа. Обязательство исполнителя считается исполненным с момента фактического открытия доступа к соответствующему функционалу.'
          },
          {
            label: '5. Срок и прекращение доступа',
            value:
              'Доступ действует в течение оплаченного периода. По окончании срока доступ прекращается, если иное прямо не указано в описании конкретного тарифа или дополнительного предложения.'
          },
          {
            label: '6. Возврат денежных средств',
            value:
              'После предоставления доступа возврат денежных средств возможен только в случаях, предусмотренных законодательством, либо если доступ не был предоставлен или сервис был оказан ненадлежащим образом.'
          },
          {
            label: '7. Платёжные документы',
            value:
              'При оплате чек формируется в порядке, установленном для НПД, в том числе через сервис «Мой налог», если это требуется по применимым правилам.'
          },
          {
            label: '8. Связанные документы',
            value:
              'Правила использования платформы публикуются отдельно и применяются вместе с настоящей офертой. Порядок обработки персональных данных раскрыт в политике конфиденциальности.'
          }
        ]
      },
      terms: {
        title: 'Условия использования',
        description: 'Правила использования платформы sinde, аккаунта, данных, ограничений доступа и мер безопасности.',
        introTitle: 'О документе',
        introItems: [
          {
            label: 'Назначение',
            value: 'Документ описывает правила использования платформы и аккаунта, а не условия оплаты.'
          },
          {
            label: 'Связанные документы',
            value: 'Оферта о покупке доступа и политика конфиденциальности действуют отдельно.'
          }
        ],
        sectionsTitle: 'Правила',
        sections: [
          {
            label: '1. Доступ и аккаунт',
            value:
              'Пользователь обязан указывать достоверные данные при регистрации, не передавать доступ третьим лицам и самостоятельно контролировать сохранность логина, пароля и иных средств аутентификации.'
          },
          {
            label: '2. Допустимое использование',
            value:
              'Платформа используется только в законных целях и в пределах предоставленного функционала. Запрещены попытки обхода ограничений, вмешательство в работу сервиса, массовый парсинг, распространение вредоносного кода и использование аккаунта от имени другого лица без разрешения.'
          },
          {
            label: '3. Пользовательские данные',
            value:
              'Пользователь отвечает за законность, актуальность и содержание материалов и данных, которые он добавляет в систему. Сервис вправе хранить и обрабатывать эти данные в объёме, необходимом для работы платформы, поддержки, безопасности и исполнения обязательств.'
          },
          {
            label: '4. Ограничения доступа',
            value:
              'Исполнитель вправе временно ограничить или полностью прекратить доступ к платформе при нарушении настоящих правил, требований законодательства, условий оферты, а также по техническим причинам или по требованию уполномоченных органов.'
          },
          {
            label: '5. Интеллектуальные права',
            value:
              'Права на программный код, интерфейсы, дизайн, базу данных и иные результаты интеллектуальной деятельности принадлежат исполнителю или законным правообладателям. Пользователь не получает право на копирование, публикацию, перепродажу или передачу материалов платформы за пределами разрешённого использования.'
          },
          {
            label: '6. Безопасность',
            value:
              'Платформа применяет технические и организационные меры защиты, однако пользователь также обязан соблюдать базовые меры безопасности: использовать надёжный пароль, не раскрывать данные входа и своевременно сообщать о подозрительной активности.'
          },
          {
            label: '7. Изменения и коммуникации',
            value:
              'Исполнитель вправе обновлять функционал и редакции документов. Актуальные версии публикуются на сайте и применяются с момента размещения, если иной срок не указан отдельно.'
          }
        ]
      },
      privacy: {
        title: 'Политика конфиденциальности',
        description: 'Порядок обработки персональных данных пользователей платформы sinde.',
        introTitle: 'Основное',
        introItems: [
          { label: 'Оператор', value: 'ИП Воротников Дмитрий Владимирович' },
          { label: 'Контакты', value: 'bank@sinde.ru', href: 'mailto:bank@sinde.ru' },
          { label: 'Документ об оплате', value: 'Оферта', href: '/docs/offer' }
        ],
        sectionsTitle: 'Обработка данных',
        sections: [
          {
            label: '1. Какие данные обрабатываются',
            value:
              'Оператор может обрабатывать данные, которые пользователь сообщает при регистрации, оплате, обращениях в поддержку и использовании сервиса: имя, email, телефон, сведения аккаунта, технические журналы доступа, а также информацию, необходимую для подтверждения оплаты и предоставления доступа.'
          },
          {
            label: '2. Цели обработки',
            value:
              'Персональные данные используются для регистрации, аутентификации, предоставления цифрового доступа, сопровождения заказов, обратной связи, безопасности, поддержки пользователей и исполнения требований законодательства.'
          },
          {
            label: '3. Правовые основания',
            value:
              'Обработка осуществляется на основании согласия пользователя, необходимости исполнения договора или оферты, а также исполнения обязанностей, предусмотренных законодательством.'
          },
          {
            label: '4. Платёжные данные',
            value:
              'Оплата проходит через платёжные сервисы. Оператор получает только те сведения, которые необходимы для подтверждения успешной оплаты, предоставления доступа, учёта операций и исполнения законных обязанностей.'
          },
          {
            label: '5. Хранение и передача',
            value:
              'Данные хранятся в объёме и в течение срока, необходимых для целей обработки, безопасности сервиса и исполнения обязательств. Передача третьим лицам возможна только при наличии законного основания, для работы инфраструктуры сервиса либо по поручению пользователя.'
          },
          {
            label: '6. Права пользователя',
            value:
              'Пользователь вправе запрашивать уточнение, обновление, удаление и иные действия с персональными данными в пределах, предусмотренных применимым законодательством и техническими возможностями сервиса.'
          },
          {
            label: '7. Безопасность и обращения',
            value:
              'Оператор применяет разумные технические и организационные меры защиты данных. По вопросам обработки персональных данных и реализации прав пользователя можно писать на bank@sinde.ru.'
          }
        ]
      }
    }
  },
  en: {
    tabsTitle: 'Documents',
    documents: {
      company: {
        title: 'Company',
        description: 'Legal details, contacts, and communication channel for the seller and service provider of sinde.',
        introTitle: 'Contacts',
        introItems: [
          { label: 'Provider', value: 'Dmitry Vladimirovich Vorotnikov, Individual Entrepreneur' },
          { label: 'Email', value: 'bank@sinde.ru', href: 'mailto:bank@sinde.ru' },
          { label: 'Phone', value: '+7 932 007-35-24', href: 'tel:+79320073524' },
          { label: 'Response time', value: 'Requests are usually answered within one day.' }
        ],
        sectionsTitle: 'Legal details',
        sections: [
          {
            label: 'Registered name',
            value: 'INDIVIDUAL ENTREPRENEUR VOROTNIKOV DMITRY VLADIMIROVICH'
          },
          {
            label: 'Registered address',
            value: '298517, Russia, Republic of Crimea, Alushta, Yaltinskaya St., 3, Apt. 31'
          },
          { label: 'TIN', value: '910102272967' },
          { label: 'Primary registration number', value: '326910000004858' },
          { label: 'Settlement account', value: '40802810700009527502' },
          { label: 'Bank', value: 'TBank JSC' },
          { label: 'BIC', value: '044525974' },
          { label: 'Correspondent account', value: '30101810145250000974' }
        ]
      },
      offer: {
        title: 'Purchase Terms',
        description: 'Terms for purchasing digital access, payment, acceptance, delivery of access, and refunds.',
        introTitle: 'Overview',
        introItems: [
          { label: 'Format', value: 'Purchase of digital access to platform functionality' },
          { label: 'Access period', value: 'Defined by the selected plan or order' },
          { label: 'Refunds', value: 'According to applicable law and these terms' },
          { label: 'Personal data', value: 'Privacy policy', href: '/docs/privacy' }
        ],
        sectionsTitle: 'Terms',
        sections: [
          {
            label: '1. General',
            value:
              'The provider publishes this offer for the purchase of digital access to the sinde platform. Payment for an order and any confirmation required in the interface constitute full and unconditional acceptance.'
          },
          {
            label: '2. Subject matter',
            value:
              'The user receives limited access to the platform functionality, interfaces, and digital content within the selected plan. Exclusive rights to software, data, and interfaces are not transferred to the user.'
          },
          {
            label: '3. Price and payment',
            value:
              'The current price, access period, and available payment methods are shown in the checkout interface at the moment of purchase. Payment is deemed completed after the payment provider confirms successful withdrawal of funds.'
          },
          {
            label: '4. Delivery of access',
            value:
              'Access is opened after payment confirmation in the scope of the selected plan. The provider’s obligation is considered fulfilled when access to the relevant functionality is actually granted.'
          },
          {
            label: '5. Access term and termination',
            value:
              'Access remains active during the paid period. After that period ends, access is terminated unless a specific plan or additional offer explicitly states otherwise.'
          },
          {
            label: '6. Refunds',
            value:
              'After access has been granted, refunds are possible only where required by law or when access was not provided or the service was rendered improperly.'
          },
          {
            label: '7. Payment documents',
            value:
              'Upon payment, a receipt is issued in the manner applicable to the NPD tax regime, including through the My Tax service when required by the applicable rules.'
          },
          {
            label: '8. Related documents',
            value:
              'The platform usage rules are published separately and apply together with this offer. Personal data processing is described in the privacy policy.'
          }
        ]
      },
      terms: {
        title: 'Terms of Use',
        description:
          'Rules for using the sinde platform, the account, user data, access restrictions, and security measures.',
        introTitle: 'About',
        introItems: [
          {
            label: 'Purpose',
            value: 'This document describes platform and account usage rules rather than payment terms.'
          },
          {
            label: 'Related documents',
            value: 'The purchase offer and the privacy policy apply separately.'
          }
        ],
        sectionsTitle: 'Rules',
        sections: [
          {
            label: '1. Access and account',
            value:
              'The user must provide accurate information during registration, must not transfer access to third parties, and is responsible for protecting login credentials and other authentication tools.'
          },
          {
            label: '2. Acceptable use',
            value:
              'The platform may be used only for lawful purposes and within the provided functionality. Attempts to bypass restrictions, interfere with the service, scrape data in bulk, distribute malicious code, or use another person’s account without permission are prohibited.'
          },
          {
            label: '3. User data',
            value:
              'The user is responsible for the legality, accuracy, and content of materials and data added to the system. The service may store and process such data to the extent necessary for platform operation, support, security, and performance of obligations.'
          },
          {
            label: '4. Access restrictions',
            value:
              'The provider may temporarily restrict or fully terminate access to the platform in case of violations of these rules, the law, the offer terms, as well as for technical reasons or upon request of competent authorities.'
          },
          {
            label: '5. Intellectual property',
            value:
              'Rights to the source code, interfaces, design, database, and other intellectual property belong to the provider or lawful right holders. The user does not receive the right to copy, publish, resell, or transfer platform materials beyond the permitted use.'
          },
          {
            label: '6. Security',
            value:
              'The platform uses technical and organizational safeguards, but the user must also follow basic security measures: use a strong password, keep credentials confidential, and report suspicious activity without delay.'
          },
          {
            label: '7. Changes and communications',
            value:
              'The provider may update the functionality and the document versions. Current versions are published on the site and apply from the moment of publication unless a different effective date is specified.'
          }
        ]
      },
      privacy: {
        title: 'Privacy Policy',
        description: 'Rules for processing personal data of sinde platform users.',
        introTitle: 'Overview',
        introItems: [
          { label: 'Operator', value: 'Dmitry Vladimirovich Vorotnikov, Individual Entrepreneur' },
          { label: 'Contacts', value: 'bank@sinde.ru', href: 'mailto:bank@sinde.ru' },
          { label: 'Payment document', value: 'Offer', href: '/docs/offer' }
        ],
        sectionsTitle: 'Data processing',
        sections: [
          {
            label: '1. Data processed',
            value:
              'The operator may process data provided during registration, payment, support requests, and service use: name, email, phone number, account details, technical access logs, and information required to confirm payment and grant access.'
          },
          {
            label: '2. Purposes',
            value:
              'Personal data is used for registration, authentication, granting digital access, supporting orders, communication, service security, user support, and compliance with legal obligations.'
          },
          {
            label: '3. Legal grounds',
            value:
              'Processing is based on the user’s consent, the need to perform a contract or offer, and the fulfillment of obligations established by law.'
          },
          {
            label: '4. Payment data',
            value:
              'Payments are processed through payment services. The operator receives only the information necessary to confirm payment success, grant access, keep transaction records, and comply with legal duties.'
          },
          {
            label: '5. Storage and transfer',
            value:
              'Data is stored only to the extent and for the period necessary for the processing purposes, service security, and performance of obligations. Transfer to third parties is possible only when there is a legal basis, when required for service infrastructure, or on the user’s instructions.'
          },
          {
            label: '6. User rights',
            value:
              'The user may request clarification, updating, deletion, and other actions regarding personal data to the extent allowed by applicable law and the service’s technical capabilities.'
          },
          {
            label: '7. Security and requests',
            value:
              'The operator applies reasonable technical and organizational safeguards. Questions about personal data processing and user rights may be sent to bank@sinde.ru.'
          }
        ]
      }
    }
  },
  ch: {
    tabsTitle: '文档',
    documents: {
      company: {
        title: '经营者信息',
        description: 'sinde 平台卖方与服务提供者的法定信息、联系方式和沟通方式。',
        introTitle: '联系方式',
        introItems: [
          { label: '服务提供者', value: '个体工商户 Dmitry Vladimirovich Vorotnikov' },
          { label: 'Email', value: 'bank@sinde.ru', href: 'mailto:bank@sinde.ru' },
          { label: '电话', value: '+7 932 007-35-24', href: 'tel:+79320073524' },
          { label: '回复时间', value: '通常会在 24 小时内回复。' }
        ],
        sectionsTitle: '公司信息',
        sections: [
          {
            label: '名称',
            value: '个体工商户 VOROTNIKOV DMITRY VLADIMIROVICH'
          },
          {
            label: '注册地址',
            value: '298517, 俄罗斯, 克里米亚共和国, 阿卢什塔, 雅尔京斯卡娅街 3 号, 31 室'
          },
          { label: '纳税号', value: '910102272967' },
          { label: '注册号', value: '326910000004858' },
          { label: '结算账户', value: '40802810700009527502' },
          { label: '银行', value: 'TBank JSC' },
          { label: 'BIC', value: '044525974' },
          { label: '对应账户', value: '30101810145250000974' }
        ]
      },
      offer: {
        title: '购买条款',
        description: '购买数字访问权限时适用的支付、承诺、开通访问和退款条件。',
        introTitle: '摘要',
        introItems: [
          { label: '形式', value: '购买平台功能的数字访问权限' },
          { label: '访问期限', value: '由所选套餐或订单决定' },
          { label: '退款', value: '按适用法律和本文件处理' },
          { label: '个人数据', value: '隐私政策', href: '/docs/privacy' }
        ],
        sectionsTitle: '条件',
        sections: [
          {
            label: '1. 一般规定',
            value:
              '服务提供者发布本条款，用于购买 sinde 平台的数字访问权限。支付订单以及界面要求时的确认操作，均表示对本条款的完全接受。'
          },
          {
            label: '2. 标的',
            value:
              '用户获得在所选套餐范围内对平台功能、界面和数字内容的有限访问权限。软件、数据和界面的专有权不转让给用户。'
          },
          {
            label: '3. 价格与支付',
            value: '购买时有效的价格、访问期限和支付方式以下单界面显示内容为准。支付服务确认扣款成功后，视为付款完成。'
          },
          {
            label: '4. 提供访问',
            value: '付款确认后，将按所选套餐开通访问权限。相关功能实际开放之时，即视为服务提供者已履行开通义务。'
          },
          {
            label: '5. 期限与终止',
            value: '访问权限在已付款期间内有效。期限结束后，除非具体套餐或附加说明另有规定，访问将终止。'
          },
          {
            label: '6. 退款',
            value: '访问权限已开通后，仅在法律要求的情形，或未提供访问、服务提供不当时，才可能退款。'
          },
          {
            label: '7. 支付凭证',
            value: '付款时，如适用规则要求，将按照 NPD 相关程序开具凭证，包括通过「Мой налог」服务生成。'
          },
          {
            label: '8. 相关文件',
            value: '平台使用规则另行发布，并与本条款一并适用。个人数据处理规则见隐私政策。'
          }
        ]
      },
      terms: {
        title: '使用条款',
        description: '关于 sinde 平台、账户、用户数据、访问限制和安全措施的使用规则。',
        introTitle: '关于文件',
        introItems: [
          {
            label: '目的',
            value: '本文件说明平台和账户的使用规则，而不是支付条件。'
          },
          {
            label: '相关文件',
            value: '购买条款和隐私政策分别适用。'
          }
        ],
        sectionsTitle: '规则',
        sections: [
          {
            label: '1. 访问与账户',
            value: '用户在注册时应提供真实信息，不得将访问权限转交第三方，并应自行保护登录信息及其他认证工具。'
          },
          {
            label: '2. 合理使用',
            value:
              '平台只能在合法目的和既定功能范围内使用。禁止绕过限制、干扰服务运行、批量抓取数据、传播恶意代码，或未经许可使用他人账户。'
          },
          {
            label: '3. 用户数据',
            value:
              '用户应对其上传或输入系统的资料与数据的合法性、准确性和内容负责。服务可在平台运行、支持、安全和履约所必需范围内存储和处理这些数据。'
          },
          {
            label: '4. 访问限制',
            value: '如用户违反本规则、法律要求或购买条款，或因技术原因、主管机关要求，服务提供者可临时限制或终止访问。'
          },
          {
            label: '5. 知识产权',
            value:
              '源代码、界面、设计、数据库及其他知识产权归服务提供者或合法权利人所有。用户无权在许可范围之外复制、发布、转售或转让平台材料。'
          },
          {
            label: '6. 安全',
            value:
              '平台采取技术和组织措施保护数据，但用户也必须遵守基本安全要求：使用强密码、不泄露登录信息，并及时报告可疑活动。'
          },
          {
            label: '7. 变更与通知',
            value: '服务提供者可更新功能与文件版本。最新版本发布于网站，自发布之时起生效，除非另有说明。'
          }
        ]
      },
      privacy: {
        title: '隐私政策',
        description: '关于 sinde 平台用户个人数据处理规则的说明。',
        introTitle: '摘要',
        introItems: [
          { label: '运营者', value: '个体工商户 Dmitry Vladimirovich Vorotnikov' },
          { label: '联系方式', value: 'bank@sinde.ru', href: 'mailto:bank@sinde.ru' },
          { label: '支付文件', value: '条款', href: '/docs/offer' }
        ],
        sectionsTitle: '数据处理',
        sections: [
          {
            label: '1. 处理哪些数据',
            value:
              '运营者可处理用户在注册、支付、联系支持和使用服务时提供的数据，包括姓名、邮箱、电话、账户信息、技术访问日志，以及确认支付和开通访问所需的信息。'
          },
          {
            label: '2. 处理目的',
            value: '个人数据用于注册、身份验证、提供数字访问、支持订单、沟通、服务安全、用户支持以及履行法律义务。'
          },
          {
            label: '3. 法律依据',
            value: '数据处理基于用户同意、履行合同或购买条款的必要性，以及法律规定的义务。'
          },
          {
            label: '4. 支付数据',
            value: '支付通过支付服务完成。运营者仅接收确认付款成功、开通访问、记录交易和履行法定义务所必需的信息。'
          },
          {
            label: '5. 存储与传输',
            value:
              '数据仅在实现处理目的、保障服务安全和履约所需的范围与期限内保存。只有在具备法律依据、服务基础设施需要或用户指示时，才会向第三方传输。'
          },
          {
            label: '6. 用户权利',
            value: '在适用法律和服务技术能力允许的范围内，用户可以请求更正、更新、删除或采取其他与个人数据相关的操作。'
          },
          {
            label: '7. 安全与联系',
            value:
              '运营者采取合理的技术和组织措施保护数据。关于个人数据处理和用户权利的问题，可以发送至 bank@sinde.ru。'
          }
        ]
      }
    }
  },
  jp: {
    tabsTitle: '文書',
    documents: {
      company: {
        title: '事業者情報',
        description: 'sinde の販売者兼提供者に関する法的情報、連絡先、問い合わせ方法です。',
        introTitle: '連絡先',
        introItems: [
          { label: '提供者', value: '個人事業主 Dmitry Vladimirovich Vorotnikov' },
          { label: 'Email', value: 'bank@sinde.ru', href: 'mailto:bank@sinde.ru' },
          { label: '電話', value: '+7 932 007-35-24', href: 'tel:+79320073524' },
          { label: '回答時間', value: '通常は 24 時間以内に回答します。' }
        ],
        sectionsTitle: '事業者情報',
        sections: [
          {
            label: '名称',
            value: '個人事業主 VOROTNIKOV DMITRY VLADIMIROVICH'
          },
          {
            label: '登録住所',
            value: '298517, ロシア, クリミア共和国, アルシタ, ヤルチンスカヤ通り 3, 31'
          },
          { label: '納税者番号', value: '910102272967' },
          { label: '登録番号', value: '326910000004858' },
          { label: '決済口座', value: '40802810700009527502' },
          { label: '銀行', value: 'TBank JSC' },
          { label: 'BIC', value: '044525974' },
          { label: 'コルレス口座', value: '30101810145250000974' }
        ]
      },
      offer: {
        title: '購入条件',
        description: 'デジタルアクセス購入時の支払い、承諾、アクセス提供、返金条件です。',
        introTitle: '概要',
        introItems: [
          { label: '形式', value: 'プラットフォーム機能へのデジタルアクセス購入' },
          { label: 'アクセス期間', value: '選択したプランまたは注文により決定' },
          { label: '返金', value: '適用法令および本書の条件による' },
          { label: '個人データ', value: 'プライバシーポリシー', href: '/docs/privacy' }
        ],
        sectionsTitle: '条件',
        sections: [
          {
            label: '1. 総則',
            value:
              '提供者は sinde プラットフォームのデジタルアクセス購入のために本オファーを公開します。注文の支払い、および必要な場合の画面上の確認操作は、本条件への完全な承諾を意味します。'
          },
          {
            label: '2. 対象',
            value:
              '利用者は選択したプランの範囲で、プラットフォーム機能、インターフェース、デジタルコンテンツへの限定的アクセスを受けます。ソフトウェア、データ、インターフェースに対する排他的権利は移転しません。'
          },
          {
            label: '3. 価格と支払い',
            value:
              '有効な価格、アクセス期間、支払い方法は購入時の注文画面に表示されます。決済サービスが引き落とし成功を確認した時点で支払い完了とみなされます。'
          },
          {
            label: '4. アクセス提供',
            value:
              '支払い確認後、選択したプランに応じてアクセスが開通します。該当機能への実際のアクセスが提供された時点で、提供者の義務は履行済みとみなされます。'
          },
          {
            label: '5. 期間と終了',
            value:
              'アクセスは支払い済み期間中有効です。期間終了後は、特定プランや追加条件に別段の定めがない限りアクセスは終了します。'
          },
          {
            label: '6. 返金',
            value:
              'アクセス提供後の返金は、法令上必要な場合、またはアクセス未提供もしくはサービス不履行の場合に限り可能です。'
          },
          {
            label: '7. 支払書類',
            value:
              '支払い時には、適用ルールに応じて NPD 制度の手順に従い、必要な場合は My Tax サービスを通じて領収情報が作成されます。'
          },
          {
            label: '8. 関連文書',
            value:
              'プラットフォーム利用ルールは別文書として公開され、本オファーとあわせて適用されます。個人データ処理はプライバシーポリシーで定められます。'
          }
        ]
      },
      terms: {
        title: '利用規約',
        description: 'sinde プラットフォーム、アカウント、利用者データ、アクセス制限、安全対策に関する利用ルールです。',
        introTitle: '文書について',
        introItems: [
          {
            label: '目的',
            value: '本書は支払い条件ではなく、プラットフォームおよびアカウントの利用ルールを定めます。'
          },
          {
            label: '関連文書',
            value: '購入オファーとプライバシーポリシーは別途適用されます。'
          }
        ],
        sectionsTitle: 'ルール',
        sections: [
          {
            label: '1. アクセスとアカウント',
            value:
              '利用者は登録時に正確な情報を提供し、第三者へアクセスを渡してはならず、ログイン情報その他の認証手段を自ら保護する責任を負います。'
          },
          {
            label: '2. 適正利用',
            value:
              'プラットフォームは適法な目的かつ提供された機能の範囲内でのみ利用できます。制限の回避、サービス妨害、大量スクレイピング、悪意あるコードの配布、他人のアカウントの無断利用は禁止されます。'
          },
          {
            label: '3. 利用者データ',
            value:
              '利用者は自ら追加した資料やデータの適法性、正確性、内容について責任を負います。サービスは、運用、サポート、安全性、義務履行に必要な範囲でそれらを保存・処理できます。'
          },
          {
            label: '4. アクセス制限',
            value:
              '本規約、法令、オファー条件への違反がある場合、または技術的理由や権限機関の要請がある場合、提供者はアクセスを一時制限または終了できます。'
          },
          {
            label: '5. 知的財産',
            value:
              'ソースコード、インターフェース、デザイン、データベースその他の知的財産権は提供者または正当な権利者に帰属します。利用者は許可された範囲を超えて、複製、公開、再販売、譲渡を行えません。'
          },
          {
            label: '6. セキュリティ',
            value:
              'プラットフォームは技術的・組織的な保護措置を講じますが、利用者も強固なパスワードの使用、認証情報の秘匿、疑わしい活動の速やかな報告など、基本的な安全対策を守る必要があります。'
          },
          {
            label: '7. 変更と通知',
            value:
              '提供者は機能や文書の内容を更新できます。最新版はサイト上で公表され、別途定めがない限り掲載時から適用されます。'
          }
        ]
      },
      privacy: {
        title: 'プライバシーポリシー',
        description: 'sinde プラットフォーム利用者の個人データ処理に関するルールです。',
        introTitle: '概要',
        introItems: [
          { label: '運営者', value: '個人事業主 Dmitry Vladimirovich Vorotnikov' },
          { label: '連絡先', value: 'bank@sinde.ru', href: 'mailto:bank@sinde.ru' },
          { label: '支払い文書', value: 'オファー', href: '/docs/offer' }
        ],
        sectionsTitle: 'データ処理',
        sections: [
          {
            label: '1. 処理されるデータ',
            value:
              '運営者は、登録、支払い、サポートへの問い合わせ、サービス利用時に提供されるデータを処理できます。これには氏名、メールアドレス、電話番号、アカウント情報、技術ログ、支払い確認とアクセス提供に必要な情報が含まれます。'
          },
          {
            label: '2. 利用目的',
            value:
              '個人データは、登録、認証、デジタルアクセス提供、注文対応、連絡、サービスの安全確保、利用者サポート、法令遵守のために使用されます。'
          },
          {
            label: '3. 法的根拠',
            value: '処理は、利用者の同意、契約またはオファーの履行の必要性、ならびに法令上の義務に基づいて行われます。'
          },
          {
            label: '4. 支払いデータ',
            value:
              '支払いは決済サービスを通じて処理されます。運営者が受け取るのは、支払い確認、アクセス提供、取引記録、法的義務の履行に必要な情報のみです。'
          },
          {
            label: '5. 保存と第三者提供',
            value:
              'データは、処理目的、サービスの安全性、義務履行に必要な範囲と期間に限って保存されます。第三者への提供は、法的根拠がある場合、サービス基盤に必要な場合、または利用者の指示がある場合に限られます。'
          },
          {
            label: '6. 利用者の権利',
            value:
              '適用法令およびサービスの技術的範囲内で、利用者は個人データの訂正、更新、削除その他の対応を求めることができます。'
          },
          {
            label: '7. 安全性と問い合わせ',
            value:
              '運営者は合理的な技術的・組織的保護措置を講じます。個人データ処理や利用者の権利に関する問い合わせは bank@sinde.ru まで送信できます。'
          }
        ]
      }
    }
  }
}
