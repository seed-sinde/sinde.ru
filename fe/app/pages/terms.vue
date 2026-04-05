<script setup lang="ts">
  const title = 'Кодекс участника'
  usePageSeo({
    title,
    description: 'Соглашение между участником и платформой sinde.'
  })
</script>
<template>
  <div>
    <LabNavHeader :title />
    <section class="p-4 space-y-4">
      <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        <!-- 1 -->
        <div class="rounded-2xl border border-zinc-800 bg-zinc-900/50 p-4 hover:bg-zinc-900/80 transition">
          <div class="mb-2 flex items-center gap-2 text-sm text-cyan-400">
            <Icon name="lucide:user-check" class="h-4 w-4" />
            <span>Контроль</span>
          </div>
          <p class="text-sm text-gray-300">
            Ты добавляешь только то, что считаешь уместным. Все действия и данные — твоя зона ответственности.
          </p>
        </div>
        <!-- 2 -->
        <div class="rounded-2xl border border-zinc-800 bg-zinc-900/50 p-4 hover:bg-zinc-900/80 transition">
          <div class="mb-2 flex items-center gap-2 text-sm text-emerald-400">
            <Icon name="lucide:database" class="h-4 w-4" />
            <span>Использование данных</span>
          </div>
          <p class="text-sm text-gray-300">
            Данные используются только внутри системы: хранение, отображение, связи и развитие.
          </p>
        </div>
        <!-- 3 -->
        <div class="rounded-2xl border border-zinc-800 bg-zinc-900/50 p-4 hover:bg-zinc-900/80 transition">
          <div class="mb-2 flex items-center gap-2 text-sm text-amber-400">
            <Icon name="lucide:filter" class="h-4 w-4" />
            <span>Минимизация</span>
          </div>
          <p class="text-sm text-gray-300">
            Ничего лишнего не собирается — только то, что ты явно вводишь или передаёшь.
          </p>
        </div>
        <!-- 4 -->
        <div class="rounded-2xl border border-zinc-800 bg-zinc-900/50 p-4 hover:bg-zinc-900/80 transition">
          <div class="mb-2 flex items-center gap-2 text-sm text-violet-400">
            <Icon name="lucide:refresh-ccw" class="h-4 w-4" />
            <span>Актуальность</span>
          </div>
          <p class="text-sm text-gray-300">Данные обновляются или удаляются тобой. Система не хранит лишнего.</p>
        </div>
        <!-- 5 -->
        <div class="rounded-2xl border border-zinc-800 bg-zinc-900/50 p-4 hover:bg-zinc-900/80 transition">
          <div class="mb-2 flex items-center gap-2 text-sm text-blue-400">
            <Icon name="lucide:lock" class="h-4 w-4" />
            <span>Доступ</span>
          </div>
          <p class="text-sm text-gray-300">Доступ к данным ограничен логикой системы и твоими действиями.</p>
        </div>
        <!-- 6 -->
        <div class="rounded-2xl border border-zinc-800 bg-zinc-900/50 p-4 hover:bg-zinc-900/80 transition">
          <div class="mb-2 flex items-center gap-2 text-sm text-green-400">
            <Icon name="lucide:shield" class="h-4 w-4" />
            <span>Безопасность</span>
          </div>
          <p class="text-sm text-gray-300">Защита встроена в архитектуру и не требует отдельного контроля.</p>
        </div>
        <!-- 7 -->
        <div class="rounded-2xl border border-zinc-800 bg-zinc-900/50 p-4 hover:bg-zinc-900/80 transition">
          <div class="mb-2 flex items-center gap-2 text-sm text-pink-400">
            <Icon name="lucide:eye" class="h-4 w-4" />
            <span>Прозрачность</span>
          </div>
          <p class="text-sm text-gray-300">Все изменения системы доступны сразу. Ничего скрытого.</p>
        </div>
        <!-- 8 -->
        <div class="rounded-2xl border border-zinc-800 bg-zinc-900/50 p-4 hover:bg-zinc-900/80 transition">
          <div class="mb-2 flex items-center gap-2 text-sm text-orange-400">
            <Icon name="lucide:check-circle" class="h-4 w-4" />
            <span>Согласие</span>
          </div>
          <p class="text-sm text-gray-300">
            Использование сервиса означает согласие с принципами: прозрачность и контроль у пользователя.
          </p>
        </div>
      </div>
    </section>
  </div>
</template>
