<script setup lang="ts">
  const title = 'Мастерская'
  usePageSeo({
    title,
    description:
      'Раздел для изучения производства, цепочек, ролей, гипотез и улучшений вокруг любого изделия, вещества или технологического процесса.'
  })
  const coreSignals = [
    { label: 'Фокус', value: 'как производится объект, из чего он состоит и почему итог получается именно таким' },
    {
      label: 'Подход',
      value: 'единая модель от ресурса и процесса до детали, продукта, закупок, поставок и результата'
    },
    { label: 'Режим', value: 'настраиваемая глубина представления для разных ролей и задач' }
  ]
  const audienceCards = [
    {
      title: 'Для изучения',
      text: 'Подходит каждому, кто хочет понять, как производится инструмент, товар, вещество, устройство или любой технологический процесс.'
    },
    {
      title: 'Для проектной работы',
      text: 'Полезно тем, кто участвует в разработке, улучшении, анализе, производстве, закупках, поставках, эксплуатации или продаже.'
    },
    {
      title: 'Для коллективного улучшения',
      text: 'Команда может совместно вносить гипотезы, сравнивать изменения, фиксировать результаты и улучшать весь процесс без разрыва между слоями данных.'
    }
  ]
  const userValueCards = [
    {
      title: 'Понять объект целиком',
      points: [
        'видеть путь от сырья, материала и технологии до детали и конечного продукта',
        'понимать, какие этапы влияют на качество, стоимость, срок и результат',
        'быстро находить узкие места, зависимости и альтернативные маршруты'
      ]
    },
    {
      title: 'Настроить глубину просмотра',
      points: [
        'переключать уровень представления без смены самой модели',
        'смотреть объект как процесс, технологию, инженерную систему или экономическую цепочку',
        'оставлять интерфейс понятным даже при большом объёме данных'
      ]
    },
    {
      title: 'Работать по временной линии',
      points: [
        'прослеживать этапы производства и изменение состояния проекта во времени',
        'видеть, какие изменения в процессе привели к конкретному результату',
        'сравнивать разные версии решений на одной линии производства'
      ]
    }
  ]
  const representationLayers = [
    'Образный — контекст применения, сценарии, визуальные референсы',
    'Процессный — этапы, переходы, последовательности действий и зависимости',
    'Технологический — материалы, режимы, ограничения, качество, допуски',
    'Инженерный — расчёты, спецификации, чертежи, прототипы, тесты',
    'Индустриальный — оборудование, повторяемость, обслуживание, поставки, закупки',
    'Экономический — стоимость этапов, риски, влияние на итог',
    'Стратегический — зависимость от стран, индустрий, цепочек и сценариев масштабирования'
  ]
  const timelineStages = [
    'Ресурс',
    'Переработка',
    'Материал',
    'Технология',
    'Оборудование',
    'Деталь',
    'Сборка',
    'Продукт',
    'Эксплуатация',
    'Изменение результата'
  ]
  const roleViews = [
    {
      title: 'Разные роли',
      text: 'Инженер, технолог, оператор, аналитик, закупщик, исследователь или архитектор цепочки видят один и тот же объект, но в нужной им глубине.'
    },
    {
      title: 'Один источник модели',
      text: 'Не нужно разрывать проект между отдельными таблицами, документами и объяснениями. Всё связывается внутри одного пространства.'
    },
    {
      title: 'Контролируемая сложность',
      text: 'Сложные данные раскрываются постепенно: от общего описания до точных параметров, цепочек, причин и последствий.'
    }
  ]
  const collaborationCards = [
    {
      title: 'Совместные гипотезы',
      text: 'Можно предлагать улучшения процесса, конструкции, маршрута производства, закупки, логистики или состава и затем проверять их по фактическому результату.'
    },
    {
      title: 'Прозрачность изменений',
      text: 'Система помогает видеть, что именно изменили, на каком этапе, кем это было предложено и как это повлияло на итог.'
    },
    {
      title: 'Артефакты проекта',
      text: 'Внутри могут жить чертежи, схемы, рисунки прототипов, отчётность, данные по поставщикам, закупкам, качеству, производству и продажам.'
    }
  ]
  const usageFlow = [
    {
      title: '1. Выбрать объект',
      text: 'Открывается изделие, материал, вещество, технологический процесс или проектная линия.'
    },
    {
      title: '2. Настроить представление',
      text: 'Пользователь выбирает глубину, роль и нужный уровень детализации без перегрузки интерфейса.'
    },
    {
      title: '3. Просмотреть цепочку',
      text: 'На временной линии становятся видны этапы, связи, зависимости, решения и результаты по объекту.'
    },
    {
      title: '4. Сравнить и улучшить',
      text: 'Можно смотреть альтернативы, вносить гипотезы, обсуждать изменения и фиксировать эффект от них.'
    }
  ]
  const contentAreas = [
    'цепочки производства и этапы',
    'материалы, составы и свойства',
    'оборудование и производственные режимы',
    'детали, узлы и конечные продукты',
    'чертежи, схемы и прототипы',
    'закупки, поставщики и логистика',
    'качество, испытания и контроль',
    'стоимость, ограничения и итоговые эффекты'
  ]
</script>
<template>
  <div class="space-y-4">
    <LabNavHeader :title />
    <section class="space-y-4 p-4">
      <section class="overflow-hidden border border-zinc-800 bg-zinc-950/78 p-4 sm:p-5">
        <div class="grid gap-4 xl:grid-cols-5">
          <div class="space-y-4 xl:col-span-3">
            <div
              class="inline-flex items-center rounded-full border border-zinc-700 bg-zinc-900/80 px-3 py-1 text-[11px] font-medium tracking-wide text-zinc-300">
              sinde • мастерская процессов и изделий
            </div>
            <div class="space-y-3">
              <h1 class="text-3xl font-semibold leading-tight text-zinc-100 sm:text-4xl">
                Один раздел для изучения, анализа и улучшения любого производственного процесса.
              </h1>
              <p class="max-w-3xl text-sm leading-7 text-zinc-400 sm:text-base">
                Мастерская показывает путь объекта во времени: из чего он состоит, как производится, какие роли влияют
                на результат, где находятся узкие места и какие изменения реально улучшают итог.
              </p>
            </div>
            <div class="grid gap-3 sm:grid-cols-3">
              <article
                v-for="item in coreSignals"
                :key="item.label"
                class="rounded-2xl border border-zinc-800 bg-zinc-900/65 px-4 py-4">
                <div class="text-[11px] uppercase tracking-[0.08em] text-zinc-500">{{ item.label }}</div>
                <p class="mt-2 text-sm leading-6 text-zinc-200">{{ item.value }}</p>
              </article>
            </div>
          </div>
          <div class="space-y-3 xl:col-span-2">
            <article class="rounded-2xl border border-zinc-800 bg-zinc-900/65 p-4">
              <div
                class="inline-flex rounded-full border border-amber-500/30 bg-amber-500/10 px-3 py-1 text-[11px] font-medium text-amber-200">
                Что даёт раздел
              </div>
              <p class="mt-3 text-sm leading-6 text-zinc-300">
                Пользователь получает не просто описание объекта, а управляемое пространство для понимания причин,
                зависимостей, решений и последствий на всей производственной линии.
              </p>
            </article>
            <article class="rounded-2xl border border-zinc-800 bg-zinc-900/65 p-4">
              <div
                class="inline-flex rounded-full border border-sky-500/30 bg-sky-500/10 px-3 py-1 text-[11px] font-medium text-sky-200">
                Как использовать
              </div>
              <p class="mt-3 text-sm leading-6 text-zinc-400">
                Открыть объект, выбрать глубину представления, пройти по временной линии, найти узкие места, сравнить
                изменения и оценить их влияние на итоговый результат.
              </p>
            </article>
          </div>
        </div>
      </section>
      <section class="grid gap-4 xl:grid-cols-3">
        <article
          v-for="item in audienceCards"
          :key="item.title"
          class="rounded-3xl border border-zinc-800 bg-zinc-950/78 p-5">
          <h2 class="text-xl font-semibold text-zinc-100">{{ item.title }}</h2>
          <p class="mt-3 text-sm leading-7 text-zinc-400">{{ item.text }}</p>
        </article>
      </section>
      <section class="space-y-4">
        <div class="space-y-2">
          <div class="text-xs uppercase tracking-[0.2em] text-amber-300">Польза для пользователя</div>
          <h2 class="text-2xl font-semibold text-zinc-100 sm:text-3xl">
            Коллективная работа над оптимизацией производства
          </h2>
          <p class="max-w-2xl text-sm leading-6 text-zinc-400">
            Интерфейс раскрывает сложность по необходимости и остаётся понятным на любом уровне детализации.
          </p>
        </div>
        <div class="grid gap-4 xl:grid-cols-3">
          <article
            v-for="item in userValueCards"
            :key="item.title"
            class="rounded-3xl border border-zinc-800 bg-zinc-950/78 p-5">
            <h3 class="text-xl font-semibold text-zinc-100">{{ item.title }}</h3>
            <ul class="mt-4 space-y-2 text-sm text-zinc-300">
              <li v-for="point in item.points" :key="point" class="flex items-start gap-2">
                <span class="mt-2 h-1.5 w-1.5 shrink-0 rounded-full bg-emerald-400"></span>
                <span class="leading-6">{{ point }}</span>
              </li>
            </ul>
          </article>
        </div>
      </section>
      <section class="grid gap-4 xl:grid-cols-[minmax(0,1.2fr)_minmax(0,1fr)]">
        <article class="rounded-3xl border border-zinc-800 bg-zinc-950/78 p-5">
          <div class="space-y-2">
            <div class="text-xs uppercase tracking-[0.2em] text-sky-300">Уровни представления</div>
            <h2 class="text-2xl font-semibold text-zinc-100">Глубина настраивается, модель остаётся одной</h2>
          </div>
          <div class="mt-4 space-y-3">
            <div
              v-for="item in representationLayers"
              :key="item"
              class="rounded-2xl border border-zinc-800 bg-zinc-900/65 px-4 py-3 text-sm leading-6 text-zinc-300">
              {{ item }}
            </div>
          </div>
        </article>
        <article class="rounded-3xl border border-zinc-800 bg-zinc-950/78 p-5">
          <div class="space-y-2">
            <div class="text-xs uppercase tracking-[0.2em] text-emerald-300">Временная линия</div>
            <h2 class="text-2xl font-semibold text-zinc-100">Производство читается как последовательность этапов</h2>
          </div>
          <div class="mt-4 flex flex-wrap gap-x-3 gap-y-2 text-sm leading-6 text-zinc-300">
            <template v-for="(stage, index) in timelineStages" :key="stage">
              <span>{{ stage }}</span>
              <span v-if="index !== timelineStages.length - 1" class="text-zinc-600">→</span>
            </template>
          </div>
          <p class="mt-4 text-sm leading-6 text-zinc-400">
            На этой линии в моменте переключается представление: можно смотреть общий путь, инженерные детали,
            экономический контур или роль-зависимый слой без смены объекта.
          </p>
        </article>
      </section>
      <section class="grid gap-4 xl:grid-cols-3">
        <article
          v-for="item in roleViews"
          :key="item.title"
          class="rounded-3xl border border-zinc-800 bg-zinc-950/78 p-5">
          <h2 class="text-xl font-semibold text-zinc-100">{{ item.title }}</h2>
          <p class="mt-3 text-sm leading-7 text-zinc-400">{{ item.text }}</p>
        </article>
      </section>
      <section class="space-y-4">
        <div class="space-y-2">
          <div class="text-xs uppercase tracking-[0.2em] text-cyan-300">Совместная работа</div>
          <h2 class="text-2xl font-semibold text-zinc-100 sm:text-3xl">
            Изучение, изменения и проверка результата в одном месте
          </h2>
          <p class="max-w-2xl text-sm leading-6 text-zinc-400">
            Раздел ориентирован не только на просмотр, но и на коллективное улучшение процессов, изделий и решений.
          </p>
        </div>
        <div class="grid gap-4 xl:grid-cols-3">
          <article
            v-for="item in collaborationCards"
            :key="item.title"
            class="rounded-3xl border border-zinc-800 bg-zinc-950/78 p-5">
            <h3 class="text-xl font-semibold text-zinc-100">{{ item.title }}</h3>
            <p class="mt-3 text-sm leading-7 text-zinc-400">{{ item.text }}</p>
          </article>
        </div>
      </section>
      <section class="grid gap-4 xl:grid-cols-[minmax(0,1.1fr)_minmax(0,1fr)]">
        <article class="rounded-3xl border border-zinc-800 bg-zinc-950/78 p-5">
          <div class="space-y-2">
            <div class="text-xs uppercase tracking-[0.2em] text-lime-300">Как пользоваться</div>
            <h2 class="text-2xl font-semibold text-zinc-100">Рабочий сценарий пользователя</h2>
          </div>
          <div class="mt-4 space-y-3">
            <article
              v-for="item in usageFlow"
              :key="item.title"
              class="rounded-2xl border border-zinc-800 bg-zinc-900/65 p-4">
              <h3 class="text-base font-semibold text-zinc-100">{{ item.title }}</h3>
              <p class="mt-2 text-sm leading-6 text-zinc-400">{{ item.text }}</p>
            </article>
          </div>
        </article>
        <article class="rounded-3xl border border-zinc-800 bg-zinc-950/78 p-5">
          <div class="space-y-2">
            <div class="text-xs uppercase tracking-[0.2em] text-rose-300">Что может содержать раздел</div>
            <h2 class="text-2xl font-semibold text-zinc-100">Данные и артефакты проекта</h2>
          </div>
          <div class="mt-4 flex flex-wrap gap-x-3 gap-y-2 text-sm leading-6 text-zinc-300">
            <template v-for="(item, index) in contentAreas" :key="item">
              <span>{{ item }}</span>
              <span v-if="index !== contentAreas.length - 1" class="text-zinc-600">•</span>
            </template>
          </div>
          <p class="mt-4 text-sm leading-6 text-zinc-400">
            Всё это нужно не ради накопления данных, а ради улучшения продукта, уменьшения потерь, повышения качества и
            лучшего понимания причин результата.
          </p>
        </article>
      </section>
    </section>
  </div>
</template>
