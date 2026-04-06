import { mountSuspended } from '@nuxt/test-utils/runtime'
import { describe, expect, it } from 'vitest'
import LabBaseButton from '~/components/Lab/Base/Button.vue'

describe('LabBaseButton', () => {
  it('renders label and loading state', async () => {
    const wrapper = await mountSuspended(LabBaseButton, {
      props: {
        label: 'Сохранить',
        loading: true,
        loadingLabel: 'Сохраняем',
        icon: 'ic:round-save'
      }
    })

    expect(wrapper.text()).toContain('Сохраняем')
    expect(wrapper.get('button').attributes('aria-busy')).toBe('true')
    expect(wrapper.get('button').classes()).toContain('inline-flex')
  })
})
