<!-- <script setup lang="ts">
const title = 'Мастерская'
usePageSeo({
  title,
  description:
    'Раздел для изучения производства, цепочек, ролей, гипотез и улучшений вокруг любого изделия, вещества или технологического процесса.'
})
</script>
<template>
  <div>
    <LabNavHeader :title />
    <section class="p-4"></section>
  </div>
</template> -->
<template>
  <section class="min-h-screen bg-(--lab-bg) text-(--lab-text-primary)">
    <div class="border-b border-(--lab-border) px-4 py-3">
      <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div class="min-w-0">
          <h1 class="truncate text-lg font-semibold">Мастерская</h1>
          <p class="text-sm text-(--lab-text-muted)">Список → структура → данные → запуск → история → анализ</p>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <button
            class="lab-focus inline-flex h-10 items-center justify-center border border-(--lab-border) px-3 text-sm hover:bg-white/2"
            type="button"
            @click="createProject"
          >
            Новый процесс
          </button>
          <button
            class="lab-focus inline-flex h-10 items-center justify-center border border-(--lab-border) px-3 text-sm hover:bg-white/2"
            type="button"
            @click="runProject"
          >
            {{ activeProject.status === 'running' ? 'Остановить запуск' : 'Запустить' }}
          </button>
        </div>
      </div>
    </div>
    <div class="grid min-h-[calc(100vh-73px)] grid-cols-1 xl:grid-cols-[320px_minmax(0,1fr)_320px]">
      <aside class="border-b border-(--lab-border) xl:border-r xl:border-b-0">
        <div class="border-b border-(--lab-border) p-4">
          <div class="flex flex-col gap-3">
            <input
              v-model="search"
              class="lab-control lab-focus min-h-10 border border-(--lab-border) bg-transparent px-3 text-sm"
              type="text"
              placeholder="Поиск процесса"
            />
            <div class="flex flex-wrap gap-2">
              <button
                v-for="item in statusFilters"
                :key="item"
                class="lab-focus inline-flex h-9 items-center justify-center border px-3 text-xs"
                :class="statusFilter === item ? 'border-(--lab-accent) text-(--lab-accent)' : 'border-(--lab-border) text-(--lab-text-muted) hover:bg-white/2'"
                type="button"
                @click="statusFilter = item"
              >
                {{ statusLabel(item) }}
              </button>
            </div>
          </div>
        </div>
        <div class="p-2">
          <button
            v-for="project in filteredProjects"
            :key="project.id"
            class="lab-focus flex w-full flex-col items-start gap-2 border px-3 py-3 text-left"
            :class="activeProjectId === project.id ? 'border-(--lab-accent) bg-white/2' : 'border-(--lab-border) hover:bg-white/2'"
            type="button"
            @click="activeProjectId = project.id"
          >
            <div class="flex w-full items-start justify-between gap-3">
              <div class="min-w-0">
                <div class="truncate text-sm font-medium">{{ project.name }}</div>
                <div class="truncate text-xs text-(--lab-text-muted)">{{ project.description }}</div>
              </div>
              <span
                class="shrink-0 border px-2 py-0.5 text-[11px]"
                :class="project.status === 'running' ? 'border-(--lab-accent) text-(--lab-accent)' : 'border-(--lab-border) text-(--lab-text-muted)'"
              >
                {{ statusLabel(project.status) }}
              </span>
            </div>
            <div class="flex flex-wrap gap-1">
              <span
                v-for="tag in project.tags"
                :key="tag"
                class="border border-(--lab-border) px-2 py-0.5 text-[11px] text-(--lab-text-muted)"
              >
                {{ tag }}
              </span>
            </div>
          </button>
          <div v-if="!filteredProjects.length" class="p-4 text-sm text-(--lab-text-muted)">Ничего не найдено</div>
        </div>
      </aside>
      <main class="min-w-0">
        <div class="border-b border-(--lab-border) px-4 py-3">
          <div class="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
            <div class="min-w-0">
              <div class="flex flex-wrap items-center gap-2">
                <input
                  v-model="activeProject.name"
                  class="lab-control lab-focus min-h-10 min-w-0 border border-(--lab-border) bg-transparent px-3 text-base font-semibold"
                  type="text"
                  aria-label="Название процесса"
                />
                <span class="border border-(--lab-border) px-2 py-1 text-xs text-(--lab-text-muted)">
                  {{ activeProject.steps.length }} шагов
                </span>
                <span class="border border-(--lab-border) px-2 py-1 text-xs text-(--lab-text-muted)">
                  {{ activeProject.events.length }} событий
                </span>
              </div>
              <textarea
                v-model="activeProject.description"
                class="lab-control lab-focus mt-2 min-h-20 w-full border border-(--lab-border) bg-transparent px-3 py-2 text-sm"
                aria-label="Описание процесса"
              />
            </div>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="tab in tabs"
                :key="tab"
                class="lab-focus inline-flex h-10 items-center justify-center border px-3 text-sm"
                :class="activeTab === tab ? 'border-(--lab-accent) text-(--lab-accent)' : 'border-(--lab-border) text-(--lab-text-muted) hover:bg-white/2'"
                type="button"
                @click="activeTab = tab"
              >
                {{ tabLabel(tab) }}
              </button>
            </div>
          </div>
        </div>
        <div class="p-4">
          <section v-if="activeTab === 'structure'" class="space-y-4">
            <div class="flex flex-wrap items-center justify-between gap-3">
              <div>
                <h2 class="text-sm font-semibold">Структура процесса</h2>
                <p class="text-xs text-(--lab-text-muted)">Inline-редактирование узлов и связей без переходов</p>
              </div>
              <button
                class="lab-focus inline-flex h-10 items-center justify-center border border-(--lab-border) px-3 text-sm hover:bg-white/2"
                type="button"
                @click="addStep"
              >
                Добавить шаг
              </button>
            </div>
            <div class="grid gap-3">
              <article
                v-for="(step, index) in activeProject.steps"
                :key="step.id"
                class="border border-(--lab-border) p-2"
              >
                <div class="flex flex-col gap-3 xl:flex-row xl:items-start">
                  <div
                    class="flex h-9 w-9 shrink-0 items-center justify-center border border-(--lab-border) text-xs text-(--lab-text-muted)"
                  >
                    {{ index + 1 }}
                  </div>
                  <div class="min-w-0 flex-1 space-y-3">
                    <div class="grid gap-3 lg:grid-cols-[minmax(0,1fr)_160px_160px_auto]">
                      <input
                        v-model="step.title"
                        class="lab-control lab-focus min-h-10 border border-(--lab-border) bg-transparent px-3 text-sm"
                        type="text"
                        aria-label="Название шага"
                      />
                      <select
                        v-model="step.kind"
                        class="lab-control lab-focus min-h-10 border border-(--lab-border) bg-transparent px-3 text-sm"
                        aria-label="Тип шага"
                      >
                        <option value="input">Ввод</option>
                        <option value="transform">Обработка</option>
                        <option value="review">Проверка</option>
                        <option value="output">Результат</option>
                      </select>
                      <select
                        v-model="step.status"
                        class="lab-control lab-focus min-h-10 border border-(--lab-border) bg-transparent px-3 text-sm"
                        aria-label="Статус шага"
                      >
                        <option value="idle">Ожидает</option>
                        <option value="running">В работе</option>
                        <option value="done">Готово</option>
                        <option value="error">Ошибка</option>
                      </select>
                      <button
                        class="lab-focus inline-flex h-10 items-center justify-center border border-(--lab-border) px-3 text-sm hover:bg-white/2"
                        type="button"
                        @click="removeStep(step.id)"
                      >
                        Удалить
                      </button>
                    </div>
                    <textarea
                      v-model="step.note"
                      class="lab-control lab-focus min-h-20 w-full border border-(--lab-border) bg-transparent px-3 py-2 text-sm"
                      aria-label="Описание шага"
                    />
                    <div class="grid gap-3 lg:grid-cols-2">
                      <div class="border border-(--lab-border) p-2">
                        <div class="mb-2 flex items-center justify-between gap-2">
                          <h3 class="text-xs font-medium">Входные данные</h3>
                          <button
                            class="lab-focus inline-flex h-8 items-center justify-center border border-(--lab-border) px-2 text-xs hover:bg-white/2"
                            type="button"
                            @click="addInput(step.id)"
                          >
                            Добавить
                          </button>
                        </div>
                        <div class="space-y-2">
                          <div
                            v-for="input in step.inputs"
                            :key="input.id"
                            class="grid gap-2 lg:grid-cols-[140px_minmax(0,1fr)_auto]"
                          >
                            <input
                              v-model="input.key"
                              class="lab-control lab-focus min-h-9 border border-(--lab-border) bg-transparent px-3 text-xs"
                              type="text"
                              placeholder="Ключ"
                            />
                            <input
                              v-model="input.value"
                              class="lab-control lab-focus min-h-9 border border-(--lab-border) bg-transparent px-3 text-xs"
                              type="text"
                              placeholder="Значение"
                            />
                            <button
                              class="lab-focus inline-flex h-9 items-center justify-center border border-(--lab-border) px-2 text-xs hover:bg-white/2"
                              type="button"
                              @click="removeInput(step.id, input.id)"
                            >
                              ×
                            </button>
                          </div>
                        </div>
                      </div>
                      <div class="border border-(--lab-border) p-2">
                        <div class="mb-2 flex items-center justify-between gap-2">
                          <h3 class="text-xs font-medium">Выходные данные</h3>
                          <button
                            class="lab-focus inline-flex h-8 items-center justify-center border border-(--lab-border) px-2 text-xs hover:bg-white/2"
                            type="button"
                            @click="addOutput(step.id)"
                          >
                            Добавить
                          </button>
                        </div>
                        <div class="space-y-2">
                          <div
                            v-for="output in step.outputs"
                            :key="output.id"
                            class="grid gap-2 lg:grid-cols-[140px_minmax(0,1fr)_auto]"
                          >
                            <input
                              v-model="output.key"
                              class="lab-control lab-focus min-h-9 border border-(--lab-border) bg-transparent px-3 text-xs"
                              type="text"
                              placeholder="Ключ"
                            />
                            <input
                              v-model="output.value"
                              class="lab-control lab-focus min-h-9 border border-(--lab-border) bg-transparent px-3 text-xs"
                              type="text"
                              placeholder="Значение"
                            />
                            <button
                              class="lab-focus inline-flex h-9 items-center justify-center border border-(--lab-border) px-2 text-xs hover:bg-white/2"
                              type="button"
                              @click="removeOutput(step.id, output.id)"
                            >
                              ×
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </article>
            </div>
          </section>
          <section v-else-if="activeTab === 'data'" class="space-y-4">
            <div>
              <h2 class="text-sm font-semibold">Сводка данных</h2>
              <p class="text-xs text-(--lab-text-muted)">Единое место для входов и результатов по всем шагам</p>
            </div>
            <div class="grid gap-4 lg:grid-cols-2">
              <div class="border border-(--lab-border) p-2">
                <h3 class="mb-3 text-xs font-medium">Все входные параметры</h3>
                <div class="space-y-2">
                  <div
                    v-for="row in allInputs"
                    :key="row.id"
                    class="grid gap-2 border border-(--lab-border) p-2 lg:grid-cols-[140px_minmax(0,1fr)_120px]"
                  >
                    <div class="truncate text-xs text-(--lab-text-muted)">{{ row.step }}</div>
                    <div class="truncate text-sm">{{ row.key }}: {{ row.value }}</div>
                    <div class="text-xs text-(--lab-text-muted)">{{ row.kind }}</div>
                  </div>
                </div>
              </div>
              <div class="border border-(--lab-border) p-2">
                <h3 class="mb-3 text-xs font-medium">Все выходные результаты</h3>
                <div class="space-y-2">
                  <div
                    v-for="row in allOutputs"
                    :key="row.id"
                    class="grid gap-2 border border-(--lab-border) p-2 lg:grid-cols-[140px_minmax(0,1fr)_120px]"
                  >
                    <div class="truncate text-xs text-(--lab-text-muted)">{{ row.step }}</div>
                    <div class="truncate text-sm">{{ row.key }}: {{ row.value }}</div>
                    <div class="text-xs text-(--lab-text-muted)">{{ row.kind }}</div>
                  </div>
                </div>
              </div>
            </div>
          </section>
          <section v-else-if="activeTab === 'history'" class="space-y-4">
            <div>
              <h2 class="text-sm font-semibold">История и состояния</h2>
              <p class="text-xs text-(--lab-text-muted)">Журнал изменений и запусков процесса</p>
            </div>
            <div class="space-y-2">
              <article
                v-for="event in activeProject.events"
                :key="event.id"
                class="grid gap-2 border border-(--lab-border) p-2 lg:grid-cols-[180px_120px_minmax(0,1fr)]"
              >
                <div class="text-xs text-(--lab-text-muted)">{{ event.time }}</div>
                <div class="text-xs">{{ event.type }}</div>
                <div class="text-sm">{{ event.text }}</div>
              </article>
            </div>
          </section>
          <section v-else class="space-y-4">
            <div>
              <h2 class="text-sm font-semibold">Анализ процесса</h2>
              <p class="text-xs text-(--lab-text-muted)">Сводные показатели и точки риска</p>
            </div>
            <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
              <article class="border border-(--lab-border) p-2">
                <div class="text-xs text-(--lab-text-muted)">Всего шагов</div>
                <div class="mt-2 text-2xl font-semibold">{{ activeProject.steps.length }}</div>
              </article>
              <article class="border border-(--lab-border) p-2">
                <div class="text-xs text-(--lab-text-muted)">Готово</div>
                <div class="mt-2 text-2xl font-semibold">{{ doneSteps }}</div>
              </article>
              <article class="border border-(--lab-border) p-2">
                <div class="text-xs text-(--lab-text-muted)">Ошибок</div>
                <div class="mt-2 text-2xl font-semibold">{{ errorSteps }}</div>
              </article>
              <article class="border border-(--lab-border) p-2">
                <div class="text-xs text-(--lab-text-muted)">Запусков</div>
                <div class="mt-2 text-2xl font-semibold">{{ runEvents }}</div>
              </article>
            </div>
            <div class="border border-(--lab-border) p-2">
              <h3 class="mb-3 text-xs font-medium">Узкие места</h3>
              <div class="space-y-2">
                <div
                  v-for="item in analysisItems"
                  :key="item.id"
                  class="grid gap-2 border border-(--lab-border) p-2 lg:grid-cols-[180px_minmax(0,1fr)]"
                >
                  <div class="text-xs text-(--lab-text-muted)">{{ item.label }}</div>
                  <div class="text-sm">{{ item.value }}</div>
                </div>
              </div>
            </div>
          </section>
        </div>
      </main>
      <aside class="border-t border-(--lab-border) xl:border-t-0 xl:border-l">
        <div class="border-b border-(--lab-border) p-4">
          <h2 class="text-sm font-semibold">Контекст</h2>
          <p class="mt-1 text-xs text-(--lab-text-muted)">Текущее состояние активного процесса</p>
        </div>
        <div class="space-y-4 p-4">
          <section class="border border-(--lab-border) p-2">
            <h3 class="mb-3 text-xs font-medium">Статус</h3>
            <div class="space-y-2 text-sm">
              <div class="flex items-center justify-between gap-3">
                <span class="text-(--lab-text-muted)">Состояние</span>
                <span>{{ statusLabel(activeProject.status) }}</span>
              </div>
              <div class="flex items-center justify-between gap-3">
                <span class="text-(--lab-text-muted)">Активная вкладка</span>
                <span>{{ tabLabel(activeTab) }}</span>
              </div>
              <div class="flex items-center justify-between gap-3">
                <span class="text-(--lab-text-muted)">Последнее изменение</span>
                <span>{{ lastEventTime }}</span>
              </div>
            </div>
          </section>
          <section class="border border-(--lab-border) p-2">
            <h3 class="mb-3 text-xs font-medium">Быстрые действия</h3>
            <div class="grid gap-2">
              <button
                class="lab-focus inline-flex h-10 items-center justify-center border border-(--lab-border) px-3 text-sm hover:bg-white/2"
                type="button"
                @click="addStep"
              >
                Добавить шаг
              </button>
              <button
                class="lab-focus inline-flex h-10 items-center justify-center border border-(--lab-border) px-3 text-sm hover:bg-white/2"
                type="button"
                @click="activeTab = 'history'"
              >
                Открыть историю
              </button>
              <button
                class="lab-focus inline-flex h-10 items-center justify-center border border-(--lab-border) px-3 text-sm hover:bg-white/2"
                type="button"
                @click="activeTab = 'analysis'"
              >
                Открыть анализ
              </button>
            </div>
          </section>
          <section class="border border-(--lab-border) p-2">
            <h3 class="mb-3 text-xs font-medium">Последние события</h3>
            <div class="space-y-2">
              <article
                v-for="event in activeProject.events.slice(0, 5)"
                :key="event.id"
                class="border border-(--lab-border) p-2"
              >
                <div class="text-[11px] text-(--lab-text-muted)">{{ event.time }}</div>
                <div class="mt-1 text-xs">{{ event.text }}</div>
              </article>
            </div>
          </section>
        </div>
      </aside>
    </div>
  </section>
</template>

<script setup lang="ts">
type ProjectStatus = 'idle' | 'running' | 'done' | 'error'
type StepKind = 'input' | 'transform' | 'review' | 'output'
type MainTab = 'structure' | 'data' | 'history' | 'analysis'
type IOField = {
  id: string
  key: string
  value: string
}
type StepItem = {
  id: string
  title: string
  kind: StepKind
  status: ProjectStatus
  note: string
  inputs: IOField[]
  outputs: IOField[]
}
type ProjectEvent = {
  id: string
  time: string
  type: 'create' | 'update' | 'run' | 'status'
  text: string
}
type ProjectItem = {
  id: string
  name: string
  description: string
  status: ProjectStatus
  tags: string[]
  steps: StepItem[]
  events: ProjectEvent[]
}
const tabs: MainTab[] = ['structure', 'data', 'history', 'analysis']
const statusFilters = ['all', 'idle', 'running', 'done', 'error'] as const
type StatusFilter = typeof statusFilters[number]
const createId = () => Math.random().toString(36).slice(2, 10)
const createField = (key = '', value = ''): IOField => ({ id: createId(), key, value })
const createStep = (title: string, kind: StepKind, status: ProjectStatus, note: string, inputs: IOField[], outputs: IOField[]): StepItem => ({
  id: createId(),
  title,
  kind,
  status,
  note,
  inputs,
  outputs
})
const createEvent = (type: ProjectEvent['type'], text: string): ProjectEvent => ({
  id: createId(),
  time: new Date().toLocaleString('ru-RU'),
  type,
  text
})
const projects = ref<ProjectItem[]>([
  {
    id: createId(),
    name: 'Подготовка публикации',
    description: 'Процесс от черновика до итогового выпуска.',
    status: 'running',
    tags: ['контент', 'редактура'],
    steps: [
      createStep('Сбор материалов', 'input', 'done', 'Получение исходных источников и тезисов.', [createField('sources', '12'), createField('brief', 'готов')], [createField('draft', 'v1')]),
      createStep('Редактура', 'review', 'running', 'Проверка структуры и тона.', [createField('draft', 'v1')], [createField('draft', 'v2')]),
      createStep('Публикация', 'output', 'idle', 'Выгрузка на площадку.', [createField('draft', 'v2')], [createField('url', '')])
    ],
    events: [
      createEvent('create', 'Процесс создан'),
      createEvent('run', 'Запуск процесса'),
      createEvent('status', 'Шаг «Редактура» переведен в работу')
    ]
  },
  {
    id: createId(),
    name: 'Исследование интерфейса',
    description: 'Флоу от гипотезы до прототипа.',
    status: 'idle',
    tags: ['ux', 'исследование'],
    steps: [
      createStep('Формулировка гипотезы', 'input', 'done', 'Определить цель и ограничения.', [createField('goal', 'повысить ясность')], [createField('hypothesis', 'готово')]),
      createStep('Прототипирование', 'transform', 'idle', 'Собрать черновой макет.', [createField('hypothesis', 'готово')], [createField('prototype', '')])
    ],
    events: [
      createEvent('create', 'Процесс создан')
    ]
  }
])
const search = ref('')
const statusFilter = ref<StatusFilter>('all')
const activeTab = ref<MainTab>('structure')
const activeProjectId = ref(projects.value[0]?.id || '')
const filteredProjects = computed(() =>
  projects.value.filter(project =>
    (statusFilter.value === 'all' ? true : project.status === statusFilter.value) &&
    (!search.value.trim() ? true : `${project.name} ${project.description} ${project.tags.join(' ')}`.toLowerCase().includes(search.value.trim().toLowerCase()))
  )
)
const emptyProject = (): ProjectItem => ({
  id: '',
  name: '',
  description: '',
  status: 'idle',
  tags: [],
  steps: [],
  events: []
})

const activeProject = computed<ProjectItem>(() =>
  projects.value.find(p => p.id === activeProjectId.value) ??
  projects.value[0] ??
  emptyProject()
)
const allInputs = computed(() =>
  activeProject.value.steps.flatMap(step =>
    step.inputs.map(input => ({ id: input.id, step: step.title, kind: step.kind, key: input.key, value: input.value }))
  )
)
const allOutputs = computed(() =>
  activeProject.value.steps.flatMap(step =>
    step.outputs.map(output => ({ id: output.id, step: step.title, kind: step.kind, key: output.key, value: output.value }))
  )
)
const doneSteps = computed(() => activeProject.value.steps.filter(step => step.status === 'done').length)
const errorSteps = computed(() => activeProject.value.steps.filter(step => step.status === 'error').length)
const runEvents = computed(() => activeProject.value.events.filter(event => event.type === 'run').length)
const lastEventTime = computed(() => activeProject.value.events[0]?.time || '—')
const analysisItems = computed(() => [
  { id: 'a1', label: 'Шаги без результатов', value: `${activeProject.value.steps.filter(step => !step.outputs.some(output => output.value.trim())).length}` },
  { id: 'a2', label: 'Шаги в ожидании', value: `${activeProject.value.steps.filter(step => step.status === 'idle').length}` },
  { id: 'a3', label: 'Шаги с риском', value: activeProject.value.steps.filter(step => step.status === 'error' || !step.note.trim()).map(step => step.title).join(', ') || 'Не обнаружены' }
])
const statusLabel = (status: StatusFilter | ProjectStatus) =>
  status === 'all' ? 'Все' :
  status === 'idle' ? 'Ожидает' :
  status === 'running' ? 'В работе' :
  status === 'done' ? 'Готово' :
  'Ошибка'
const tabLabel = (tab: MainTab) =>
  tab === 'structure' ? 'Структура' :
  tab === 'data' ? 'Данные' :
  tab === 'history' ? 'История' :
  'Анализ'
const appendEvent = (text: string, type: ProjectEvent['type'] = 'update') => activeProject.value.events.unshift(createEvent(type, text))
const createProject = () => {
  const project: ProjectItem = {
    id: createId(),
    name: 'Новый процесс',
    description: 'Краткое описание процесса.',
    status: 'idle',
    tags: ['новый'],
    steps: [createStep('Первый шаг', 'input', 'idle', 'Заполните описание шага.', [createField('input', '')], [createField('output', '')])],
    events: [createEvent('create', 'Процесс создан')]
  }
  projects.value.unshift(project)
  activeProjectId.value = project.id
  activeTab.value = 'structure'
}
const addStep = () => {
  activeProject.value.steps.push(createStep('Новый шаг', 'transform', 'idle', 'Опишите действие шага.', [createField('param', '')], [createField('result', '')]))
  appendEvent('Добавлен новый шаг')
}
const removeStep = (stepId: string) => {
  activeProject.value.steps = activeProject.value.steps.filter(step => step.id !== stepId)
  appendEvent('Шаг удален')
}
const addInput = (stepId: string) => {
  const step = activeProject.value.steps.find(item => item.id === stepId)
  step?.inputs.push(createField('input', ''))
  appendEvent('Добавлен входной параметр')
}
const removeInput = (stepId: string, inputId: string) => {
  const step = activeProject.value.steps.find(item => item.id === stepId)
  if (!step) return
  step.inputs = step.inputs.filter(input => input.id !== inputId)
  appendEvent('Удален входной параметр')
}
const addOutput = (stepId: string) => {
  const step = activeProject.value.steps.find(item => item.id === stepId)
  step?.outputs.push(createField('output', ''))
  appendEvent('Добавлен выходной результат')
}
const removeOutput = (stepId: string, outputId: string) => {
  const step = activeProject.value.steps.find(item => item.id === stepId)
  if (!step) return
  step.outputs = step.outputs.filter(output => output.id !== outputId)
  appendEvent('Удален выходной результат')
}
const runProject = () => {
  activeProject.value.status = activeProject.value.status === 'running' ? 'idle' : 'running'
  appendEvent(activeProject.value.status === 'running' ? 'Запуск процесса' : 'Запуск остановлен', 'run')
}
watch(
  () => activeProject.value.steps.map(step => `${step.title}|${step.kind}|${step.status}|${step.note}|${step.inputs.length}|${step.outputs.length}`).join('::'),
  () => {
    if (!activeProject.value?.id) return
    activeProject.value.events = activeProject.value.events.slice(0, 24)
  }
)
</script>
