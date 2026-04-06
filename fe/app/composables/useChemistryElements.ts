type ChemistryElementsResponse = ApiResponseWithData<ChemistryElementApi[]>

const { json: useApiJson } = useAPI()

export const getChemistryElements = async () => {
  return await useApiJson<ChemistryElementsResponse>('/chemistry/elements', {
    method: 'GET'
  })
}

export const useChemistryElements = () =>
  useAsyncData(
    'chemistry-elements',
    async () => {
      const response = await getChemistryElements()
      return sortPeriodicTableElements((response.data || []).map(mapChemistryElementApi))
    },
    {
      default: () => []
    }
  )
