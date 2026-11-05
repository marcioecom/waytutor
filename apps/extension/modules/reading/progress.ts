import { z } from 'zod'

const progressSchema = z.object({
  index: z.number().int().min(0),
  updatedAt: z.number(),
})

export type ReadingProgress = z.infer<typeof progressSchema>

const storedSchema = z.record(z.string(), progressSchema)

async function readAll(): Promise<Record<string, ReadingProgress>> {
  const stored = await chrome.storage.local.get('readingProgress')
  const parsed = storedSchema.safeParse(stored.readingProgress)

  return parsed.success ? parsed.data : {}
}

export async function readProgress(materialId: string): Promise<ReadingProgress | null> {
  return (await readAll())[materialId] ?? null
}

export async function saveProgress(materialId: string, index: number): Promise<void> {
  const all = await readAll()
  all[materialId] = { index, updatedAt: Date.now() }

  await chrome.storage.local.set({ readingProgress: all })
}
