import { describe, expect, it } from 'vitest'

describe('mineralFilters route helpers', () => {
  const compareElementOrder = (left: string, right: string) => ['H', 'O', 'Na', 'Cl'].indexOf(left) - ['H', 'O', 'Na', 'Cl'].indexOf(right)
  const allowedCrystalSystems = new Set<MineralCrystalSystem>(['cubic', 'hexagonal', 'monoclinic', 'orthorhombic', 'tetragonal', 'triclinic', 'unknown'])

  it('reads ALL / ANY / NONE chemistry filters from route query', () => {
    const state = readMineralsRouteState(
      {
        q: ' quartz ',
        chemistryAll: 'o,h',
        chemistryAny: 'na',
        chemistryNone: 'cl',
        crystalSystem: 'hexagonal,unknown',
        crystalSystemMode: 'all',
        onlyWithImages: '1',
        limit: '60',
        offset: '30',
        sort: 'name_desc'
      },
      {
        defaultLimit: 30,
        allowedCrystalSystems,
        compareElementOrder
      }
    )

    expect(state).toEqual({
      q: 'quartz',
      sort: 'name_desc',
      limit: 60,
      offset: 30,
      onlyWithImages: true,
      crystalSystems: ['hexagonal', 'unknown'],
      crystalSystemMode: 'all',
      chemistryAll: ['H', 'O'],
      chemistryAny: ['Na'],
      chemistryNone: ['Cl']
    })
  })

  it('builds compact route query without default values', () => {
    const query = buildMineralsRouteQuery(
      {
        q: ' quartz ',
        sort: 'name_asc',
        limit: 30,
        offset: 0,
        onlyWithImages: false,
        crystalSystems: [],
        crystalSystemMode: 'any',
        chemistryAll: ['O', 'H'],
        chemistryAny: [],
        chemistryNone: ['Cl']
      },
      {
        defaultLimit: 30,
        sortElements: values => values.slice().sort(compareElementOrder)
      }
    )

    expect(query).toEqual({
      q: 'quartz',
      chemistryAll: 'H,O',
      chemistryNone: 'Cl'
    })
  })

  it('detects whether route query contains any active filters', () => {
    expect(hasMineralsRouteQueryValues({})).toBe(false)
    expect(hasMineralsRouteQueryValues({ q: '  ' })).toBe(false)
    expect(hasMineralsRouteQueryValues({ chemistryAny: 'Na' })).toBe(true)
  })
})
