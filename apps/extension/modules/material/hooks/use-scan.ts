import { useQuery } from '@tanstack/react-query'
import { displayTitle, scanTab } from '../scan'

export function useScan() {
  return useQuery({
    queryKey: ['material'],
    queryFn: async () => {
      const [tab] = await chrome.tabs.query({ active: true, currentWindow: true })
      if (!tab?.id) return null

      const material = await scanTab(tab.id)
      if (!material) return null

      return { ...material, title: displayTitle(tab.title, material.title) }
    },
  })
}
