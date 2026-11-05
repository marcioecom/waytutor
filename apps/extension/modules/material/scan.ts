import { materialSchema, type Material } from './material'
import { scanMessageType } from './messages'

const MIN_CONTENT_CHARS = 200
const FRAME_TIMEOUT_MS = 2000
const CONTENT_HOST = /conteudo\.catolica\.edu\.br|\/content\/enforced\//i

export type ScanResult = {
  material: Material
  tabId: number
  frameId: number
}

type Candidate = {
  material: Material
  frameId: number
}

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

async function scanFrame(tabId: number, frameId: number): Promise<Candidate | null> {
  const response = await chrome.tabs
    .sendMessage(tabId, { type: scanMessageType }, { frameId })
    .catch(() => null)

  const parsed = materialSchema.safeParse(response)
  return parsed.success ? { material: parsed.data, frameId } : null
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
export async function scanTab(tabId: number): Promise<ScanResult | null> {
  const frames = await chrome.webNavigation.getAllFrames({ tabId })
  if (!frames?.length) return null

  const scanned = await Promise.all(
    frames.map((frame) => withTimeout(scanFrame(tabId, frame.frameId))),
  )
  const candidates = scanned.filter((candidate): candidate is Candidate => candidate !== null)
  if (!candidates.length) return null

  const content = candidates.filter((candidate) => CONTENT_HOST.test(candidate.material.url))
  const pool = content.length ? content : candidates
  const best = pool.sort((a, b) => score(b.material) - score(a.material))[0]

  if (!best || textLength(best.material) < MIN_CONTENT_CHARS) return null
  return { material: best.material, tabId, frameId: best.frameId }
}

export function displayTitle(tabTitle: string | undefined, fallback: string): string {
  const raw = (tabTitle ?? '').trim()
  if (!raw) return fallback
  if (!/^leitura\s*:/i.test(raw) && !raw.includes(' - ')) return fallback

  return raw.replace(/^leitura\s*:\s*/i, '').trim() || fallback
}
