import { ApiError, codeOf } from '@/modules/api/errors'
import { apiUrl } from '@/modules/api/config'

const maxBuffered = 3

type PlayOptions = {
  rate: number
  voice: string
}

export class AudioPlayer {
  private audio = new Audio()
  private buffered = new Map<string, string>()
  private loading = new Map<string, Promise<string>>()

  onStarted?: () => void
  onEnded?: () => void
  onError?: (code: string) => void

  constructor() {
    this.audio.addEventListener('ended', () => this.onEnded?.())
    this.audio.addEventListener('error', () => {
      if (this.audio.src) this.onError?.('audio_failed')
    })
  }

  async play(index: number, text: string, token: string, options: PlayOptions): Promise<void> {
    const key = this.key(index, options.voice)

    this.audio.src = await this.load(key, text, token, options.voice)
    this.audio.playbackRate = options.rate
    await this.audio.play()
    this.onStarted?.()
  }

  prefetch(index: number, text: string, token: string, voice: string): void {
    void this.load(this.key(index, voice), text, token, voice).catch(() => undefined)
  }

  pause(): void {
    this.audio.pause()
  }

  resume(): void {
    void this.audio.play()
  }

  setRate(rate: number): void {
    this.audio.playbackRate = rate
  }

  stop(): void {
    this.audio.pause()
    this.audio.removeAttribute('src')
  }

  private key(index: number, voice: string): string {
    return `${index}:${voice}`
  }

  // A prefetch that already started is reused by play, so one Passage costs one
  // request and one quota reservation.
  private load(key: string, text: string, token: string, voice: string): Promise<string> {
    const ready = this.buffered.get(key)
    if (ready) return Promise.resolve(ready)

    const started = this.loading.get(key)
    if (started) return started

    const request = this.fetchAudio(text, token, voice)
      .then((url) => {
        this.buffered.set(key, url)
        this.evict(key)
        return url
      })
      .finally(() => this.loading.delete(key))

    this.loading.set(key, request)
    return request
  }

  private async fetchAudio(text: string, token: string, voice: string): Promise<string> {
    const response = await fetch(apiUrl('/v1/tts'), {
      method: 'POST',
      headers: {
        'content-type': 'application/json',
        authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ text, voice }),
    })
    if (!response.ok) throw new ApiError(await codeOf(response))

    return URL.createObjectURL(await response.blob())
  }

  private evict(keep: string): void {
    for (const [key, url] of this.buffered) {
      if (this.buffered.size <= maxBuffered) return
      if (key === keep) continue

      URL.revokeObjectURL(url)
      this.buffered.delete(key)
    }
  }
}
