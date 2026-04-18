type ChemistryElementsResponse = ApiResponseWithData<ChemistryElementApi[]>
const useApiJson = <T>(path: string, options?: NonNullable<Parameters<ReturnType<typeof useAPI>["json"]>[1]>) =>
  useAPI().json<T>(path, options)

export const getChemistryElements = async () => {
  return await useApiJson<ChemistryElementsResponse>("/chemistry/elements", {
    method: "GET"
  })
}

export const useChemistryElements = () =>
  useAsyncData(
    "chemistry-elements",
    async () => {
      const response = await getChemistryElements()
      return sortPeriodicTableElements((response.data || []).map(mapChemistryElementApi))
    },
    {
      default: () => []
    }
  )
