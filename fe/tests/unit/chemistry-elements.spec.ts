import { describe, expect, it } from 'bun:test'
import { buildPeriodicTableCategoryCounts, mapChemistryElementApi, resolvePeriodicTableElement } from '../../app/utils/chemistryElements'
import type { ChemistryElementApi } from '../../shared/types/chemistry'

describe('chemistryElements utils', () => {
  const apiElement: ChemistryElementApi = {
    number: 8,
    symbol: 'o',
    name: 'Oxygen',
    russian_name: 'Кислород',
    category: 'diatomic nonmetal',
    phase: 'Gas',
    xpos: 16,
    ypos: 2,
    bohr_model_image: 'chemistry/elements/o/008-o-oxygen-bohr.webp',
    bohr_model_3d: 'chemistry/elements/o/008-o-oxygen-model.glb',
    spectral_img: 'chemistry/elements/o/008-o-oxygen-spectral.webp',
    cpk_hex: 'ff0d0d',
    samples: [
      {
        file: 'sample-1.webp',
        title: 'Жидкий кислород',
        attribution: 'NASA'
      }
    ]
  }

  it('maps backend element payload into frontend view model', () => {
    const item = mapChemistryElementApi(apiElement, 0)

    expect(item.slug).toBe('oxygen')
    expect(item.displaySymbol).toBe('O')
    expect(item.bohrModelImage).toBe('/api/proxy/media/files/chemistry/elements/o/008-o-oxygen-bohr.webp')
    expect(item.bohrModel3d).toBe('/api/proxy/media/files/chemistry/elements/o/008-o-oxygen-model.glb')
    expect(item.samples[0]?.url).toBe('/api/proxy/media/files/chemistry/elements/o/samples/sample-1.webp')
    expect(item.categoryLabel).toBe('Двухатомные неметаллы')
    expect(item.accentColor).toBe('#ff0d0d')
  })

  it('resolves element by symbol, number and slug aliases', () => {
    const item = mapChemistryElementApi(apiElement, 0)
    const items = [item]

    expect(resolvePeriodicTableElement(items, 'o')?.number).toBe(8)
    expect(resolvePeriodicTableElement(items, '8')?.number).toBe(8)
    expect(resolvePeriodicTableElement(items, 'oxygen')?.number).toBe(8)
    expect(resolvePeriodicTableElement(items, ['edu', 'chemistry', 'elements', 'oxygen'])?.number).toBe(8)
  })

  it('builds category counts from hydrated element list', () => {
    const items = [
      mapChemistryElementApi(apiElement, 0),
      mapChemistryElementApi({ ...apiElement, number: 16, symbol: 'S', name: 'Sulfur', russian_name: 'Сера' }, 1)
    ]

    expect(buildPeriodicTableCategoryCounts(items)).toEqual([
      expect.objectContaining({
        category: 'diatomic nonmetal',
        count: 2
      })
    ])
  })
})
