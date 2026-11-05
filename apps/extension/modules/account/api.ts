import { z } from 'zod'
import { ApiError, codeOf } from '@/modules/api/errors'

const apiBaseUrl = import.meta.env.WXT_API_URL ?? 'http://localhost:8080'

const accountSchema = z.object({
  id: z.string(),
  email: z.string(),
  name: z.string(),
  status: z.string(),
})

export type Account = z.infer<typeof accountSchema>

export function apiUrl(path: string): string {
  return `${apiBaseUrl}${path}`
}

export async function exchangeLoginCode(code: string): Promise<{ token: string; account: Account }> {
  const response = await fetch(apiUrl('/v1/auth/exchange'), {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify({ code }),
  })
  if (!response.ok) throw new ApiError(await codeOf(response))

  return z.object({ token: z.string(), account: accountSchema }).parse(await response.json())
}

export async function fetchAccount(token: string): Promise<Account> {
  const response = await fetch(apiUrl('/v1/me'), {
    headers: { authorization: `Bearer ${token}` },
  })
  if (!response.ok) throw new ApiError(await codeOf(response))

  return accountSchema.parse((await response.json()).account)
}

export async function endSession(token: string): Promise<void> {
  await fetch(apiUrl('/v1/auth/logout'), {
    method: 'POST',
    headers: { authorization: `Bearer ${token}` },
  })
}
