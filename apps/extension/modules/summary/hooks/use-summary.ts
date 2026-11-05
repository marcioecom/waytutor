import { useQuery } from '@tanstack/react-query'
import { readToken } from '@/modules/account/storage'
import { ApiError } from '@/modules/api/errors'
import { requestSummary } from '../api'

export function useSummary(materialId: string | null, text: string) {
  return useQuery({
    queryKey: ['summary', materialId],
    // A summary costs a model call, so opening the screen twice in a session
    // reuses it instead of paying again.
    staleTime: Infinity,
    enabled: Boolean(materialId && text),
    queryFn: async () => {
      const token = await readToken()
      if (!token) throw new ApiError('unauthenticated')

      return requestSummary(text, token)
    },
  })
}
