import { z } from 'zod'

export const passageSchema = z.object({
  text: z.string(),
  section: z.string(),
})

export const materialSchema = z.object({
  title: z.string(),
  url: z.string(),
  passages: z.array(passageSchema),
  figureAfter: z.array(z.number().int().nonnegative()),
})

export type Passage = z.infer<typeof passageSchema>
export type Material = z.infer<typeof materialSchema>
