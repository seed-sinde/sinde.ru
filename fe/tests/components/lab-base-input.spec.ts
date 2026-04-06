import { mountSuspended } from '@nuxt/test-utils/runtime'
import { describe, expect, it } from 'vitest'
import LabBaseInput from '~/components/Lab/Base/Input.vue'

describe('LabBaseInput', () => {
  it('marks invalid email and emits model updates', async () => {
    const wrapper = await mountSuspended(LabBaseInput, {
      props: {
        modelValue: 'bad-email',
        type: 'email'
      }
    })

    expect(wrapper.get('input').attributes('aria-invalid')).toBe('true')
    await wrapper.get('input').setValue('user@example.com')
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['user@example.com'])
  })
})
