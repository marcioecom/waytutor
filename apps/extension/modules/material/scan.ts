import { materialSchema, type Material } from './material'
import { scanMessageType } from './messages'

const MIN_CONTENT_CHARS = 200
const FRAME_TIMEOUT_MS = 2000
const CONTENT_HOST = /conteudo\.catolica\.edu\.br|\/content\/enforced\//i

function textLength(material: Material): number {
  return material.passages.reduce((total, passage) => total + passage.text.length, 0)
}

function score(material: Material): number {
  let value = textLength(material) + material.passages.length * 25
  if (/conteudo\.catolica\.edu\.br/i.test(material.url)) value += 4000
  if (/\/content\/enforced\//i.test(material.url)) value += 2500
  if (/^leitura\s*:/i.test(material.title)) value += 500
  return value
}

async function scanFrame(tabId: number, frameId: number): Promise<Material | null> {
  const response = await chrome.tabs
    .sendMessage(tabId, { type: scanMessageType }, { frameId })
    .catch(() => null)

  const parsed = materialSchema.safeParse(response)
  return parsed.success ? parsed.data : null
}

async function withTimeout<T>(promise: Promise<T>): Promise<T | null> {
  let timer: ReturnType<typeof setTimeout> | undefined
  const timeout = new Promise<null>((resolve) => {
    timer = setTimeout(() => resolve(null), FRAME_TIMEOUT_MS)
  })

  try {
    return await Promise.race([promise, timeout])
  } finally {
    clearTimeout(timer)
  }
}

// The portal renders the Material inside a frame, surrounded by the shell. Scan
// every frame and keep the one that actually holds the text.
export async function scanTab(tabId: number): Promise<Material | null> {
  const frames = await chrome.webNavigation.getAllFrames({ tabId })
  if (!frames?.length) return null

  const scanned = await Promise.all(
    frames.map((frame) => withTimeout(scanFrame(tabId, frame.frameId))),
  )
  const candidates = scanned.filter((material): material is Material => material !== null)
  if (!candidates.length) return null

  const content = candidates.filter((material) => CONTENT_HOST.test(material.url))
  const pool = content.length ? content : candidates
  const best = pool.sort((a, b) => score(b) - score(a))[0]

  if (!best || textLength(best) < MIN_CONTENT_CHARS) return null
  return best
}

export function displayTitle(tabTitle: string | undefined, fallback: string): string {
  const raw = (tabTitle ?? '').trim()
  if (!raw) return fallback
  if (!/^leitura\s*:/i.test(raw) && !raw.includes(' - ')) return fallback

  return raw.replace(/^leitura\s*:\s*/i, '').trim() || fallback
}
