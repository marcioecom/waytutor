import { z } from 'zod'
import { apiUrl } from '@/modules/api/config'
import { ApiError, codeOf } from '@/modules/api/errors'

const summarySchema = z.object({ summary: z.string() })

export async function requestSummary(text: string, token: string): Promise<string> {
  const response = await fetch(apiUrl('/v1/summary'), {
    method: 'POST',
    headers: {
      'content-type': 'application/json',
      authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ text }),
  })
  if (!response.ok) throw new ApiError(await codeOf(response))

  return summarySchema.parse(await response.json()).summary
}
