import { useCallback, useEffect, useState } from 'react'
import { readToken } from '@/modules/account/storage'
import type { Passage } from '@/modules/material/material'
import { highlightPassage } from '../highlight'
import { audioMessage } from '../messages'
import { ensureOffscreen, pauseAudio, playPassage, prefetchPassage, resumeAudio, stopAudio } from '../player'

const prefetchAhead = 2

export type ReadingStatus = 'idle' | 'loading' | 'playing' | 'paused'

type Props = {
  tabId: number
  frameId: number
  passages: Passage[]
}

export function useReading({ tabId, frameId, passages }: Props) {
  const [index, setIndex] = useState(0)
  const [status, setStatus] = useState<ReadingStatus>('idle')
  const [errorCode, setErrorCode] = useState<string | null>(null)

  const start = useCallback(
    async (from: number) => {
      const passage = passages[from]
      if (!passage) return

      const token = await readToken()
      if (!token) return

      setIndex(from)
      setErrorCode(null)
      setStatus('loading')

      await ensureOffscreen()
      await playPassage(from, passage.text, token)
      void highlightPassage(tabId, frameId, from)

      for (let ahead = 1; ahead <= prefetchAhead; ahead++) {
        const next = passages[from + ahead]
        if (next) void prefetchPassage(from + ahead, next.text, token)
      }
    },
    [frameId, passages, tabId],
  )

  useEffect(() => {
    const onMessage = (message: unknown) => {
      if (!message || typeof message !== 'object') return

      const event = message as { type?: string; index?: number; code?: string }
      if (event.index !== index) return

      switch (event.type) {
        case audioMessage.started:
          setStatus('playing')
          return
        case audioMessage.ended:
          if (passages[index + 1]) void start(index + 1)
          else setStatus('idle')
          return
        case audioMessage.error:
          setStatus('idle')
          setErrorCode(event.code ?? 'audio_failed')
          return
      }
    }

    chrome.runtime.onMessage.addListener(onMessage)
    return () => chrome.runtime.onMessage.removeListener(onMessage)
  }, [index, passages, start])

  // Closing the panel ends the Reading, so the panel takes the audio with it.
  useEffect(() => () => void stopAudio(), [])

  const toggle = useCallback(() => {
    if (status === 'playing') {
      void pauseAudio()
      setStatus('paused')
      return
    }
    if (status === 'paused') {
      void resumeAudio()
      setStatus('playing')
      return
    }
    void start(index)
  }, [index, start, status])

  const goTo = useCallback(
    (target: number) => {
      if (!passages[target]) return
      void start(target)
    },
    [passages, start],
  )

  return {
    index,
    status,
    errorCode,
    toggle,
    goTo,
    next: () => goTo(index + 1),
    previous: () => goTo(index - 1),
  }
}
