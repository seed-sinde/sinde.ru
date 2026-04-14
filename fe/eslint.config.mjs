import withNuxt from './.nuxt/eslint.config.mjs'
export default withNuxt({
  rules: {
    '@typescript-eslint/no-explicit-any': 'off',
    'vue/require-default-prop': 'off',
    'vue/no-required-prop-with-default': 'off',
    'vue/no-v-html': 'off',
    'no-console': ['warn', { allow: ['warn', 'error'] }]
  }
})
