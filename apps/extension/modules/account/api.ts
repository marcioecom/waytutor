import { z } from 'zod'
import { apiUrl } from '@/modules/api/config'
import { ApiError, codeOf } from '@/modules/api/errors'

const accountSchema = z.object({
  id: z.string(),
  email: z.string(),
  name: z.string(),
  status: z.string(),
})

const workspaceSchema = z.object({
  account: accountSchema,
  voices: z.array(z.string()),
  defaultVoice: z.string(),
})

export type Account = z.infer<typeof accountSchema>
export type Workspace = z.infer<typeof workspaceSchema>

export async function exchangeLoginCode(code: string): Promise<{ token: string; account: Account }> {
  const response = await fetch(apiUrl('/v1/auth/exchange'), {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify({ code }),
  })
  if (!response.ok) throw new ApiError(await codeOf(response))

  return z.object({ token: z.string(), account: accountSchema }).parse(await response.json())
}

export async function fetchWorkspace(token: string): Promise<Workspace> {
  const response = await fetch(apiUrl('/v1/me'), {
    headers: { authorization: `Bearer ${token}` },
  })
  if (!response.ok) throw new ApiError(await codeOf(response))

  return workspaceSchema.parse(await response.json())
}

export async function endSession(token: string): Promise<void> {
  await fetch(apiUrl('/v1/auth/logout'), {
    method: 'POST',
    headers: { authorization: `Bearer ${token}` },
  })
}
