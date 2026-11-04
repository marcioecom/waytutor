import { z } from 'zod'

const storedSchema = z.object({ accessToken: z.string().optional() })

export async function readToken(): Promise<string | null> {
  const stored = storedSchema.parse(await chrome.storage.local.get('accessToken'))
  return stored.accessToken ?? null
}

export async function saveToken(token: string): Promise<void> {
  await chrome.storage.local.set({ accessToken: token })
}

export async function clearToken(): Promise<void> {
  await chrome.storage.local.remove('accessToken')
}
