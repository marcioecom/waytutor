import { ApiError, codeOf } from '@/modules/api/errors'
import { audioMessage, offscreenTarget } from '@/modules/reading/messages'

const apiBaseUrl = import.meta.env.WXT_API_URL ?? 'http://localhost:8080'
const maxBuffered = 3

const audio = new Audio()
const buffered = new Map<number, string>()
const loading = new Map<number, Promise<string>>()

let currentIndex = 0

audio.addEventListener('ended', () => emit(audioMessage.ended))
audio.addEventListener('error', () => emit(audioMessage.error, 'audio_failed'))

function emit(type: string, code?: string) {
  void chrome.runtime.sendMessage({ type, index: currentIndex, code }).catch(() => undefined)
}

async function fetchAudio(text: string, token: string): Promise<string> {
  const response = await fetch(`${apiBaseUrl}/v1/tts`, {
    method: 'POST',
    headers: {
      'content-type': 'application/json',
      authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ text }),
  })
  if (!response.ok) throw new ApiError(await codeOf(response))

  return URL.createObjectURL(await response.blob())
}

// One request per Passage: a prefetch that already started is reused by play.
function load(index: number, text: string, token: string): Promise<string> {
  const ready = buffered.get(index)
  if (ready) return Promise.resolve(ready)

  const started = loading.get(index)
  if (started) return started

  const request = fetchAudio(text, token)
    .then((url) => {
      buffered.set(index, url)
      evict(index)
      return url
    })
    .finally(() => loading.delete(index))

  loading.set(index, request)
  return request
}

function evict(keep: number) {
  for (const index of buffered.keys()) {
    if (buffered.size <= maxBuffered) return
    if (index === keep) continue

    URL.revokeObjectURL(buffered.get(index)!)
    buffered.delete(index)
  }
}

async function play(message: { index: number; text: string; token: string }) {
  currentIndex = message.index
  audio.src = await load(message.index, message.text, message.token)
  await audio.play()
  emit(audioMessage.started)
}

chrome.runtime.onMessage.addListener((message: unknown) => {
  if (!message || typeof message !== 'object') return
  if ((message as { target?: unknown }).target !== offscreenTarget) return

  const command = message as { type: string; index: number; text: string; token: string }

  switch (command.type) {
    case audioMessage.play:
      play(command).catch((error: unknown) => {
        emit(audioMessage.error, error instanceof ApiError ? error.code : 'audio_failed')
      })
      return
    case audioMessage.prefetch:
      load(command.index, command.text, command.token).catch(() => undefined)
      return
    case audioMessage.pause:
      audio.pause()
      return
    case audioMessage.resume:
      void audio.play()
      return
    case audioMessage.stop:
      audio.pause()
      audio.removeAttribute('src')
      return
  }
})
