import { audioMessage, offscreenTarget } from './messages'

export async function ensureOffscreen(): Promise<void> {
  if (await chrome.offscreen.hasDocument()) return

  await chrome.offscreen.createDocument({
    url: chrome.runtime.getURL('/offscreen.html'),
    reasons: ['AUDIO_PLAYBACK'],
    justification: 'Toca o áudio da leitura em voz alta.',
  })
}

function send(message: unknown): Promise<unknown> {
  return chrome.runtime.sendMessage(message).catch(() => undefined)
}

export function playPassage(index: number, text: string, token: string): Promise<unknown> {
  return send({ target: offscreenTarget, type: audioMessage.play, index, text, token })
}

export function prefetchPassage(index: number, text: string, token: string): Promise<unknown> {
  return send({ target: offscreenTarget, type: audioMessage.prefetch, index, text, token })
}

export function pauseAudio(): Promise<unknown> {
  return send({ target: offscreenTarget, type: audioMessage.pause })
}

export function resumeAudio(): Promise<unknown> {
  return send({ target: offscreenTarget, type: audioMessage.resume })
}

export function stopAudio(): Promise<unknown> {
  return send({ target: offscreenTarget, type: audioMessage.stop })
}
