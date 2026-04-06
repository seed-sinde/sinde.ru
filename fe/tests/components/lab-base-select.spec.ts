import { mountSuspended } from '@nuxt/test-utils/runtime'
import { describe, expect, it } from 'vitest'
import LabBaseSelect from '~/components/Lab/Base/Select.vue'

describe('LabBaseSelect', () => {
  it('opens dropdown and emits chosen option', async () => {
    const wrapper = await mountSuspended(LabBaseSelect, {
      props: {
        modelValue: '',
        options: [
          { value: 'light', label: 'Светлая' },
          { value: 'dark', label: 'Тёмная' }
        ]
      },
      global: {
        stubs: {
          Icon: { template: '<i />' },
          Teleport: true
        }
      }
    })

    await wrapper.get('button').trigger('click')
    const options = wrapper.findAll('[role=\"option\"]')
    expect(options).toHaveLength(2)
    await options[1]!.trigger('click')

    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['dark'])
  })
})
