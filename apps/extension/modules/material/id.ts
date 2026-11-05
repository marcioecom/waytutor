import type { Passage } from './material'

// The portal changes the URL of the same topic, so the Material is identified by
// what it says.
export async function materialId(passages: Passage[]): Promise<string> {
  const text = passages.map((passage) => passage.text).join('\n')
  const digest = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(text))

  return Array.from(new Uint8Array(digest).slice(0, 8))
    .map((byte) => byte.toString(16).padStart(2, '0'))
    .join('')
}
