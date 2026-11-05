import { useQuery, useQueryClient } from '@tanstack/react-query'
import { defaults, readPreferences, writePreferences, type Preferences } from '../storage'

const preferencesQueryKey = ['preferences'] as const

export function usePreferences(): {
  preferences: Preferences
  update: (patch: Partial<Preferences>) => void
} {
  const queryClient = useQueryClient()
  const query = useQuery({ queryKey: preferencesQueryKey, queryFn: readPreferences })
  const preferences = query.data ?? defaults

  const update = (patch: Partial<Preferences>) => {
    const next = { ...preferences, ...patch }
    queryClient.setQueryData(preferencesQueryKey, next)
    void writePreferences(next)
  }

  return { preferences, update }
}
