<script setup lang="ts">
  const title = 'ТЕСТ: Мастерская'
  definePageMeta({
    middleware: 'admin-only'
  })
  usePageSeo({
    title,
    description: 'Управление и оптимизация производства.'
  })
  const objectTitle = ref('Твердосплавная концевая фреза Ø10')
  const selectedRole = ref<LabTestRoleKey>('engineer')
  const selectedDepth = ref<LabTestDepthKey>('process')
  const selectedStageId = ref('resource')
  const roles: LabTestOption<LabTestRoleKey>[] = [
    { key: 'engineer', label: 'Инженер' },
    { key: 'technologist', label: 'Технолог' },
    { key: 'operator', label: 'Оператор' },
    { key: 'analyst', label: 'Аналитик' },
    { key: 'procurement', label: 'Закупщик' },
    { key: 'researcher', label: 'Исследователь' },
    { key: 'architect', label: 'Архитектор цепочки' }
  ]
  const depthOptions: LabTestOption<LabTestDepthKey>[] = [
    { key: 'figurative', label: 'Образный' },
    { key: 'process', label: 'Процессный' },
    { key: 'technology', label: 'Технологический' },
    { key: 'engineering', label: 'Инженерный' },
    { key: 'industrial', label: 'Индустриальный' },
    { key: 'economic', label: 'Экономический' },
    { key: 'strategic', label: 'Стратегический' }
  ]
  const metrics = computed<LabTestMetricCard[]>(() => [
    { label: 'Себестоимость', value: '€18.40', delta: '-6.2%' },
    { label: 'Брак', value: '2.8%', delta: '-1.1 п.п.' },
    { label: 'Цикл производства', value: '4 д 6 ч', delta: '-9.4%' },
    { label: 'Стабильность качества', value: '93 / 100', delta: '+4' }
  ])
  const timeline = ref<LabTestTimelineStage[]>([
    {
      id: 'resource',
      title: 'Ресурс',
      short: 'WC, Co, добавки',
      owner: 'Закупки / аналитика сырья',
      status: 'warning',
      duration: '6 д',
      costImpact: 'Высокое',
      qualityImpact: 'Среднее',
      bottleneck: true,
      details: {
        figurative: 'Исходная база изделия: карбид вольфрама, кобальтовая связка и легирующие добавки.',
        process: 'Проверка происхождения сырья, стабильности поставок, чистоты партий и допустимых отклонений.',
        technology: 'Критичны гранулометрия, чистота порошка, влажность и стабильность состава по партиям.',
        engineering:
          'Влияние на износостойкость и прочность идёт через долю Co, размер зерна WC и равномерность смеси.',
        industrial: 'Узкое место — зависимость от ограниченного пула поставщиков по ключевому порошку.',
        economic: 'До 41% стоимости изделия формируется здесь; колебания цены сырья напрямую бьют по марже.',
        strategic: 'Риск привязки к странам и цепочкам поставки редких компонентов; нужно сценарное резервирование.'
      }
    },
    {
      id: 'processing',
      title: 'Переработка',
      short: 'Подготовка порошка',
      owner: 'Технолог',
      status: 'stable',
      duration: '10 ч',
      costImpact: 'Среднее',
      qualityImpact: 'Высокое',
      details: {
        figurative: 'Подготовка смеси до состояния, пригодного для стабильного прессования и спекания.',
        process: 'Смешивание, дозирование, диспергирование и контроль однородности порошковой массы.',
        technology: 'Режимы смешивания, связующее, влажность и время цикла определяют равномерность массы.',
        engineering: 'На этом шаге закладывается равномерность структуры будущей заготовки.',
        industrial: 'Точка с хорошей повторяемостью при наличии формализованных режимов и контроля.',
        economic: 'Потери здесь умеренные, но любая ошибка умножается на следующих стадиях.',
        strategic: 'Стандартизация режима снижает зависимость результата от смены персонала и площадки.'
      }
    },
    {
      id: 'material',
      title: 'Материал',
      short: 'Прессовка / спекание',
      owner: 'Технолог / оператор',
      status: 'critical',
      duration: '14 ч',
      costImpact: 'Высокое',
      qualityImpact: 'Высокое',
      bottleneck: true,
      details: {
        figurative: 'Материал приобретает финальную структуру и внутренние свойства.',
        process: 'Прессование, удаление связки, спекание, усадка, контроль геометрии и плотности.',
        technology: 'Температура, атмосфера, скорость нагрева и выдержка критичны для пористости и усадки.',
        engineering: 'Основной источник отклонений по микроструктуре, твёрдости и прочности.',
        industrial: 'Ограничение по печам и времени цикла создаёт очередь и держит throughput.',
        economic: 'Брак на этой стадии дорогой: в изделие уже вложено сырьё и предшествующая обработка.',
        strategic: 'Масштабирование упирается в парк печей, энергоёмкость и доступность сервисного обслуживания.'
      }
    },
    {
      id: 'technology',
      title: 'Технология',
      short: 'Шлифование / покрытие',
      owner: 'Инженер / оператор',
      status: 'improving',
      duration: '18 ч',
      costImpact: 'Среднее',
      qualityImpact: 'Высокое',
      details: {
        figurative: 'Формирование режущей геометрии и поверхностных свойств инструмента.',
        process: 'Шлифование канавок, заточка, подготовка поверхности, нанесение PVD-покрытия.',
        technology: 'Точность настройки, износ круга, режимы съёма и параметры покрытия критичны.',
        engineering: 'Именно здесь формируется реальный режущий профиль и рабочее поведение инструмента.',
        industrial: 'Есть потенциал ускорения без просадки качества за счёт маршрутизации и переналадки.',
        economic: 'Оптимизация цикла даёт быстрый эффект по сроку и OEE оборудования.',
        strategic: 'Собственная компетенция по покрытию снижает зависимость от внешних подрядчиков.'
      }
    },
    {
      id: 'equipment',
      title: 'Оборудование',
      short: 'Парк станков и печей',
      owner: 'Производство',
      status: 'warning',
      duration: 'Постоянно',
      costImpact: 'Высокое',
      qualityImpact: 'Среднее',
      details: {
        figurative: 'Материальная база производства и носитель повторяемости результата.',
        process: 'Загрузка линий, обслуживание, сменные ограничения, фактическая доступность оборудования.',
        technology: 'Паспортные режимы сами по себе ничего не гарантируют без обслуживания и стабильных настроек.',
        engineering: 'Важны точность осей, вибрации, температура, состояние шпинделей и оснастки.',
        industrial: 'Нужен слой MTBF, OEE и причин простоев по единицам парка.',
        economic: 'Простои и сервис сильно влияют на фактическую себестоимость единицы продукции.',
        strategic: 'Парк должен быть анализируем как ограничение роста и как карта зависимости от поставщиков сервиса.'
      }
    },
    {
      id: 'assembly',
      title: 'Деталь / сборка',
      short: 'Комплектность и партия',
      owner: 'ОТК / склад',
      status: 'stable',
      duration: '6 ч',
      costImpact: 'Низкое',
      qualityImpact: 'Среднее',
      details: {
        figurative: 'Формирование товарной единицы и проверка готовности к отгрузке.',
        process: 'Контроль, маркировка, упаковка, привязка к партии, документация.',
        technology: 'Нужны формализованные правила приемки и трассировки.',
        engineering: 'Тут важна связь изделия с параметрами партии и режимами производства.',
        industrial: 'Этап не должен быть «немой зоной» без связи с upstream-данными.',
        economic: 'Ошибки стоят дешевле, но создают репутационный и логистический ущерб.',
        strategic: 'Трассировка партий усиливает доверие и ускоряет анализ рекламаций.'
      }
    },
    {
      id: 'product',
      title: 'Продукт / эксплуатация',
      short: 'Работа у клиента',
      owner: 'Продажи / сервис / аналитика',
      status: 'warning',
      duration: '30–180 д',
      costImpact: 'Высокое',
      qualityImpact: 'Высокое',
      details: {
        figurative: 'Изделие проявляет себя в реальной среде применения.',
        process: 'Сбор обратной связи, ресурсные испытания, отказные случаи, сценарии применения.',
        technology: 'Нужна привязка эксплуатационного результата к геометрии, материалу и режимам производства.',
        engineering: 'Здесь подтверждаются или опровергаются инженерные гипотезы.',
        industrial: 'Эксплуатационные данные должны возвращаться в производство как источник улучшений.',
        economic: 'Именно здесь становится виден реальный итог: рекламации, повторные продажи, удержание клиента.',
        strategic: 'Рынок, отрасль применения и масштабируемость решения оцениваются на фактических данных.'
      }
    }
  ])
  const hypotheses = ref<LabTestHypothesis[]>([
    {
      id: 'h1',
      title: 'Сменить поставщика WC-порошка на более стабильную гранулометрию',
      stageId: 'resource',
      author: 'Аналитик цепочки',
      effect: 'Снижение разброса усадки и брака после спекания',
      risks: 'Рост закупочной цены на 4.8%',
      status: 'testing'
    },
    {
      id: 'h2',
      title: 'Перенести часть номенклатуры на ночное окно покрытия',
      stageId: 'technology',
      author: 'Технолог',
      effect: 'Ускорение цикла на 7–9%',
      risks: 'Нужна корректировка графика обслуживания',
      status: 'accepted'
    },
    {
      id: 'h3',
      title: 'Добавить трассировку по партиям до эксплуатации клиента',
      stageId: 'product',
      author: 'Инженер качества',
      effect: 'Быстрее находить причины рекламаций',
      risks: 'Нужно доработать модель данных и отчёты',
      status: 'draft'
    }
  ])
  const emptyStage: LabTestTimelineStage = {
    id: 'empty',
    title: 'Этап не найден',
    short: 'Нет данных',
    owner: '—',
    status: 'warning',
    duration: '—',
    costImpact: '—',
    qualityImpact: '—',
    details: {
      figurative: 'Нет данных для отображения.',
      process: 'Нет данных для отображения.',
      technology: 'Нет данных для отображения.',
      engineering: 'Нет данных для отображения.',
      industrial: 'Нет данных для отображения.',
      economic: 'Нет данных для отображения.',
      strategic: 'Нет данных для отображения.'
    }
  }
  const selectedStage = computed<LabTestTimelineStage>(() => {
    return timeline.value.find(stage => stage.id === selectedStageId.value) ?? timeline.value[0] ?? emptyStage
  })
  const stageHypotheses = computed(() => {
    return hypotheses.value.filter(item => item.stageId === selectedStage.value.id)
  })
  const rolePerspective = computed(() => {
    switch (selectedRole.value) {
      case 'engineer':
        return 'Фокус: параметры, геометрия, допуски, испытания и технические причины отклонений.'
      case 'technologist':
        return 'Фокус: режимы, маршруты, повторяемость, настройка процесса и узкие места по стадиям.'
      case 'operator':
        return 'Фокус: понятные действия на этапе, ограничения, состояние оборудования и фактический цикл.'
      case 'analyst':
        return 'Фокус: метрики процесса, причинно-следственные связи, динамика изменений и эффект решений.'
      case 'procurement':
        return 'Фокус: поставщики, риски сырья, срок поставки, стоимость и альтернативные цепочки.'
      case 'researcher':
        return 'Фокус: гипотезы, сравнение версий, новые материалы, проверка результатов и выводы.'
      case 'architect':
        return 'Фокус: целостная модель цепочки, зависимости между слоями и масштабирование системы.'
    }
  })
  const statusClasses: Record<LabTestStageStatus, string> = {
    stable: 'border-emerald-500/30 bg-emerald-500/10 text-emerald-300',
    warning: 'border-amber-500/30 bg-amber-500/10 text-amber-300',
    critical: 'border-rose-500/30 bg-rose-500/10 text-rose-300',
    improving: 'border-sky-500/30 bg-sky-500/10 text-sky-300'
  }
  const statusLabel: Record<LabTestStageStatus, string> = {
    stable: 'Стабильно',
    warning: 'Риск',
    critical: 'Узкое место',
    improving: 'Улучшается'
  }
</script>
<template>
  <div class="min-h-screen bg-zinc-950 text-zinc-100">
    <LabNavHeader :title />
    <div class="mx-auto flex max-w-400 flex-col gap-6 px-4 py-4 lg:px-6 xl:px-8">
      <header class="rounded-3xl border border-white/10 bg-white/5 p-5 backdrop-blur-sm lg:p-6">
        <div class="flex flex-col gap-5 xl:flex-row xl:items-start xl:justify-between">
          <div class="max-w-4xl space-y-4">
            <div
              class="inline-flex items-center rounded-full border border-cyan-400/20 bg-cyan-500/10 px-3 py-1 text-xs font-medium tracking-wide text-cyan-200">
              Мастерская · единая модель производственной линии
            </div>
            <div class="space-y-3">
              <h1 class="text-2xl font-semibold tracking-tight text-white sm:text-3xl lg:text-4xl">
                {{ objectTitle }}
              </h1>
              <p class="max-w-3xl text-sm leading-6 text-zinc-300 sm:text-base">
                Один раздел для изучения, анализа и улучшения любого производственного процесса: от ресурса и материала
                до детали, эксплуатации и изменения результата.
              </p>
            </div>
          </div>
          <div class="grid gap-3 sm:grid-cols-2 xl:w-110">
            <div
              v-for="metric in metrics"
              :key="metric.label"
              class="rounded-2xl border border-white/10 bg-black/20 px-4 py-3">
              <div class="text-xs uppercase tracking-wide text-zinc-500">{{ metric.label }}</div>
              <div class="mt-2 flex items-end justify-between gap-3">
                <div class="text-xl font-semibold text-white">{{ metric.value }}</div>
                <div class="text-sm text-emerald-300">{{ metric.delta }}</div>
              </div>
            </div>
          </div>
        </div>
      </header>
      <section class="grid gap-6 xl:grid-cols-[300px_minmax(0,1fr)_360px]">
        <aside class="space-y-4 rounded-3xl border border-white/10 bg-white/5 p-4">
          <div class="space-y-2">
            <div class="text-xs uppercase tracking-wide text-zinc-500">Роль</div>
            <div class="grid gap-2">
              <button
                v-for="role in roles"
                :key="role.key"
                type="button"
                class="rounded-2xl border px-3 py-2 text-left text-sm transition"
                :class="
                  selectedRole === role.key ?
                    'border-cyan-400/40 bg-cyan-500/10 text-white'
                  : 'border-white/10 bg-black/20 text-zinc-300 hover:border-white/20 hover:text-white'
                "
                @click="selectedRole = role.key">
                {{ role.label }}
              </button>
            </div>
          </div>
          <div class="space-y-2">
            <div class="text-xs uppercase tracking-wide text-zinc-500">Глубина</div>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="depth in depthOptions"
                :key="depth.key"
                type="button"
                class="rounded-full border px-3 py-1.5 text-sm transition"
                :class="
                  selectedDepth === depth.key ?
                    'border-indigo-400/40 bg-indigo-500/15 text-white'
                  : 'border-white/10 bg-black/20 text-zinc-300 hover:border-white/20 hover:text-white'
                "
                @click="selectedDepth = depth.key">
                {{ depth.label }}
              </button>
            </div>
          </div>
          <div class="rounded-2xl border border-white/10 bg-black/20 p-4">
            <div class="text-xs uppercase tracking-wide text-zinc-500">Активный взгляд</div>
            <p class="mt-3 text-sm leading-6 text-zinc-300">
              {{ rolePerspective }}
            </p>
          </div>
          <div class="rounded-2xl border border-white/10 bg-black/20 p-4">
            <div class="text-xs uppercase tracking-wide text-zinc-500">Что даёт раздел</div>
            <ul class="mt-3 space-y-2 text-sm leading-6 text-zinc-300">
              <li>Видеть объект как единую причинно-следственную цепочку.</li>
              <li>Переключать глубину без смены модели.</li>
              <li>Находить узкие места и сравнивать улучшения.</li>
              <li>Связывать производство, закупки, качество и эксплуатацию.</li>
            </ul>
          </div>
        </aside>
        <main class="space-y-6">
          <section class="rounded-3xl border border-white/10 bg-white/5 p-4 lg:p-5">
            <div class="mb-4 flex items-center justify-between gap-4">
              <div>
                <h2 class="text-lg font-semibold text-white">Временная линия</h2>
                <p class="mt-1 text-sm text-zinc-400">
                  Производство читается как последовательность этапов с переключением представления в моменте.
                </p>
              </div>
              <div
                class="hidden rounded-full border border-white/10 bg-black/20 px-3 py-1 text-xs text-zinc-400 lg:block">
                Один объект · одна модель · разные глубины
              </div>
            </div>
            <div class="overflow-x-auto pb-2">
              <div class="flex min-w-max items-stretch gap-3">
                <button
                  v-for="stage in timeline"
                  :key="stage.id"
                  type="button"
                  class="group relative w-55 rounded-3xl border p-4 text-left transition"
                  :class="
                    selectedStageId === stage.id ?
                      'border-cyan-400/40 bg-cyan-500/10'
                    : 'border-white/10 bg-black/20 hover:border-white/20 hover:bg-white/5'
                  "
                  @click="selectedStageId = stage.id">
                  <div class="flex items-start justify-between gap-3">
                    <div>
                      <div class="text-sm font-medium text-white">{{ stage.title }}</div>
                      <div class="mt-1 text-xs text-zinc-400">{{ stage.short }}</div>
                    </div>
                    <span class="rounded-full border px-2 py-1 text-xs" :class="statusClasses[stage.status]">
                      {{ statusLabel[stage.status] }}
                    </span>
                  </div>
                  <div class="mt-4 grid gap-2 text-xs text-zinc-400">
                    <div class="flex items-center justify-between gap-3">
                      <span>Время</span>
                      <span class="text-zinc-200">{{ stage.duration }}</span>
                    </div>
                    <div class="flex items-center justify-between gap-3">
                      <span>Владелец</span>
                      <span class="truncate text-zinc-200">{{ stage.owner }}</span>
                    </div>
                  </div>
                  <div
                    v-if="stage.bottleneck"
                    class="mt-4 rounded-2xl border border-amber-500/20 bg-amber-500/10 px-3 py-2 text-xs text-amber-200">
                    Узкое место линии
                  </div>
                </button>
              </div>
            </div>
          </section>
          <section class="grid gap-6 2xl:grid-cols-[minmax(0,1.15fr)_minmax(0,0.85fr)]">
            <article class="rounded-3xl border border-white/10 bg-white/5 p-5">
              <div class="flex flex-wrap items-start justify-between gap-4">
                <div>
                  <div class="text-xs uppercase tracking-wide text-zinc-500">Текущий этап</div>
                  <h3 class="mt-2 text-xl font-semibold text-white">{{ selectedStage.title }}</h3>
                  <p class="mt-1 text-sm text-zinc-400">{{ selectedStage.short }}</p>
                </div>
                <span class="rounded-full border px-3 py-1.5 text-xs" :class="statusClasses[selectedStage.status]">
                  {{ statusLabel[selectedStage.status] }}
                </span>
              </div>
              <div class="mt-5 grid gap-3 sm:grid-cols-3">
                <div class="rounded-2xl border border-white/10 bg-black/20 p-4">
                  <div class="text-xs uppercase tracking-wide text-zinc-500">Длительность</div>
                  <div class="mt-2 text-base font-medium text-white">{{ selectedStage.duration }}</div>
                </div>
                <div class="rounded-2xl border border-white/10 bg-black/20 p-4">
                  <div class="text-xs uppercase tracking-wide text-zinc-500">Влияние на стоимость</div>
                  <div class="mt-2 text-base font-medium text-white">{{ selectedStage.costImpact }}</div>
                </div>
                <div class="rounded-2xl border border-white/10 bg-black/20 p-4">
                  <div class="text-xs uppercase tracking-wide text-zinc-500">Влияние на качество</div>
                  <div class="mt-2 text-base font-medium text-white">{{ selectedStage.qualityImpact }}</div>
                </div>
              </div>
              <div class="mt-5 rounded-2xl border border-white/10 bg-black/20 p-4">
                <div class="text-xs uppercase tracking-wide text-zinc-500">Описание по выбранной глубине</div>
                <p class="mt-3 text-sm leading-7 text-zinc-200">
                  {{ selectedStage.details[selectedDepth] }}
                </p>
              </div>
            </article>
            <article class="rounded-3xl border border-white/10 bg-white/5 p-5">
              <div class="flex items-center justify-between gap-3">
                <div>
                  <div class="text-xs uppercase tracking-wide text-zinc-500">Контур влияния</div>
                  <h3 class="mt-2 text-lg font-semibold text-white">Причины и последствия</h3>
                </div>
              </div>
              <div class="mt-5 space-y-3">
                <div class="rounded-2xl border border-white/10 bg-black/20 p-4">
                  <div class="text-xs uppercase tracking-wide text-zinc-500">Что влияет сверху</div>
                  <p class="mt-2 text-sm leading-6 text-zinc-300">
                    Состав сырья, стабильность предыдущих стадий, фактическая загрузка оборудования, режимы и
                    доступность людей/оснастки.
                  </p>
                </div>
                <div class="rounded-2xl border border-white/10 bg-black/20 p-4">
                  <div class="text-xs uppercase tracking-wide text-zinc-500">Что меняет ниже</div>
                  <p class="mt-2 text-sm leading-6 text-zinc-300">
                    Геометрию, повторяемость, срок цикла, вероятность брака, качество эксплуатации и итоговую экономику.
                  </p>
                </div>
                <div class="rounded-2xl border border-white/10 bg-black/20 p-4">
                  <div class="text-xs uppercase tracking-wide text-zinc-500">Нужные артефакты</div>
                  <div class="mt-3 flex flex-wrap gap-2 text-sm text-zinc-200">
                    <span class="rounded-full border border-white/10 bg-white/5 px-3 py-1.5">Чертежи</span>
                    <span class="rounded-full border border-white/10 bg-white/5 px-3 py-1.5">Протоколы испытаний</span>
                    <span class="rounded-full border border-white/10 bg-white/5 px-3 py-1.5">Поставщики</span>
                    <span class="rounded-full border border-white/10 bg-white/5 px-3 py-1.5">Режимы обработки</span>
                    <span class="rounded-full border border-white/10 bg-white/5 px-3 py-1.5">Качество партии</span>
                  </div>
                </div>
              </div>
            </article>
          </section>
        </main>
        <aside class="space-y-4 rounded-3xl border border-white/10 bg-white/5 p-4">
          <div>
            <div class="text-xs uppercase tracking-wide text-zinc-500">Гипотезы и улучшения</div>
            <h2 class="mt-2 text-lg font-semibold text-white">Изменения по этапу</h2>
          </div>
          <div class="space-y-3">
            <article
              v-for="hypothesis in stageHypotheses"
              :key="hypothesis.id"
              class="rounded-2xl border border-white/10 bg-black/20 p-4">
              <div class="flex items-start justify-between gap-3">
                <h3 class="text-sm font-medium leading-6 text-white">{{ hypothesis.title }}</h3>
                <span
                  class="rounded-full border px-2 py-1 text-xs"
                  :class="
                    hypothesis.status === 'accepted' ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-300'
                    : hypothesis.status === 'testing' ? 'border-sky-500/30 bg-sky-500/10 text-sky-300'
                    : 'border-zinc-500/30 bg-zinc-500/10 text-zinc-300'
                  ">
                  {{
                    hypothesis.status === 'accepted' ? 'Принято'
                    : hypothesis.status === 'testing' ? 'Проверка'
                    : 'Черновик'
                  }}
                </span>
              </div>
              <div class="mt-3 space-y-2 text-sm leading-6 text-zinc-300">
                <p>
                  <span class="text-zinc-500">Автор:</span>
                  {{ hypothesis.author }}
                </p>
                <p>
                  <span class="text-zinc-500">Эффект:</span>
                  {{ hypothesis.effect }}
                </p>
                <p>
                  <span class="text-zinc-500">Риски:</span>
                  {{ hypothesis.risks }}
                </p>
              </div>
            </article>
            <div
              v-if="!stageHypotheses.length"
              class="rounded-2xl border border-dashed border-white/10 bg-black/20 p-4 text-sm leading-6 text-zinc-400">
              Для этого этапа гипотезы ещё не добавлены.
            </div>
          </div>
          <div class="rounded-2xl border border-white/10 bg-black/20 p-4">
            <div class="text-xs uppercase tracking-wide text-zinc-500">Как использовать</div>
            <ol class="mt-3 space-y-2 text-sm leading-6 text-zinc-300">
              <li>1. Открыть объект или производственную линию.</li>
              <li>2. Выбрать роль и глубину представления.</li>
              <li>3. Пройти по временной линии и найти узкие места.</li>
              <li>4. Сравнить гипотезы и фиксировать эффект изменений.</li>
            </ol>
          </div>
        </aside>
      </section>
    </div>
  </div>
</template>
