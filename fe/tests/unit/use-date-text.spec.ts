import { describe, expect, it } from 'vitest'

describe('useDateText helpers', () => {
  it('masks raw digits into a datetime string', () => {
    expect(maskDateText('010220241530', 'datetime')).toBe('01.02.2024 15:30')
  })

  it('produces canonical values for valid russian-formatted inputs', () => {
    expect(canonicalDateText('01.02.2024', 'date')).toBe('2024-02-01')
    expect(canonicalDateText('15:30', 'time')).toBe('15:30')
    expect(canonicalDateText('01.02.2024 15:30', 'datetime')).toBe('2024-02-01T15:30')
  })

  it('compares valid dates and rejects invalid ones', () => {
    expect(compareDateText('01.02.2024', '02.02.2024', 'date')).toBeLessThan(0)
    expect(compareDateText('31.02.2024', '02.02.2024', 'date')).toBeNull()
  })
})
