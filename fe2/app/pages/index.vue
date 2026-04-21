<script setup lang="ts">
const details = ref(false)
const lorem = `Lorem, ipsum dolor sit amet consectetur adipisicing elit. Vel doloremque itaque numquam ad.`
const email = ref("")
const hasEmailError = ref(true)
const hasPassError = ref(true)
const email_error = ref("Неверный формат почты")
const password = ref("")
const password_error = ref("Пароль не верный")
</script>
<template>
  <Stack>
    <Spoiler v-model="details" label="Скрытый блок">
      <p class="my-1">Содержимое скрытого блока!</p>
    </Spoiler>
    <Button label="Кнопка с иконкой слева" icon="ic:round-more-vert" icon-class="text-sm" />
    <Button label="Кнопка с иконкой справа" icon="ic:round-more-vert" icon-class="text-sm" icon-position="right" />
    <Tooltip :text="lorem" side="right">
      <span class="border-b border-dotted">Показать всплывающую подсказку на тексте</span>
    </Tooltip>
    <div class="flex max-w-fit flex-col">
      <Field label="Почта" hint="Введите пароль не менее 12 символов" :error="email_error" required>
        <Input v-model="email" type="email" placeholder="Email" :invalid="hasEmailError" />
      </Field>
      <Field label="Пароль" hint="Введите пароль не менее 12 символов" :error="password_error" required>
        <Input type="password" v-model="password" placeholder="Пароль" :invalid="hasPassError" />
      </Field>
    </div>
  </Stack>
</template>
