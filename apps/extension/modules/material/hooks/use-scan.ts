import { useQuery } from '@tanstack/react-query'
import { displayTitle, scanTab } from '../scan'

export function useScan() {
  return useQuery({
    queryKey: ['material'],
    queryFn: async () => {
      const [tab] = await chrome.tabs.query({ active: true, currentWindow: true })
      if (!tab?.id) return null

      const result = await scanTab(tab.id)
      if (!result) return null

      return {
        ...result,
        material: { ...result.material, title: displayTitle(tab.title, result.material.title) },
      }
    },
  })
}
