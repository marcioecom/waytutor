import { z } from 'zod'

export const defaults = {
  voice: null,
  speed: 1.25,
  figurePauseSeconds: 5,
  favoriteVoices: [] as string[],
} as const

export const preferencesSchema = z.object({
  voice: z.string().min(1).nullable().catch(defaults.voice),
  speed: z.number().min(0.5).max(2).catch(defaults.speed),
  figurePauseSeconds: z.number().min(0).max(15).catch(defaults.figurePauseSeconds),
  favoriteVoices: z.array(z.string()).catch(defaults.favoriteVoices),
})

export type Preferences = z.infer<typeof preferencesSchema>

export async function readPreferences(): Promise<Preferences> {
  const stored = await chrome.storage.local.get('preferences')
  const parsed = preferencesSchema.safeParse(stored.preferences)

  return parsed.success ? parsed.data : { ...defaults }
}

export async function writePreferences(preferences: Preferences): Promise<void> {
  await chrome.storage.local.set({ preferences })
}
