import { useCallback, useEffect, useRef, useState } from 'react'
import { readToken } from '@/modules/account/storage'
import { ApiError } from '@/modules/api/errors'
import type { Material } from '@/modules/material/material'
import type { Preferences } from '@/modules/preferences/storage'
import { AudioPlayer } from '../audio'
import { highlightPassage } from '../highlight'
import { readProgress, saveProgress } from '../progress'

const prefetchAhead = 2

export type ReadingStatus = 'idle' | 'loading' | 'playing' | 'paused'

type Props = {
  tabId: number
  frameId: number
  materialId: string | null
  material: Material | null
  preferences: Preferences
  defaultVoice: string
}

export function useReading({
  tabId,
  frameId,
  materialId,
  material,
  preferences,
  defaultVoice,
}: Props) {
  const [player] = useState(() => new AudioPlayer())
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const [index, setIndex] = useState(0)
  const [status, setStatus] = useState<ReadingStatus>('idle')
  const [errorCode, setErrorCode] = useState<string | null>(null)

  const voice = preferences.voice ?? defaultVoice

  const start = useCallback(
    async (from: number) => {
      const passage = material?.passages[from]
      if (!passage) return

      const token = await readToken()
      if (!token) return

      // A queued advance from a figure pause must not fire after this one.
      if (timerRef.current) clearTimeout(timerRef.current)

      setIndex(from)
      setErrorCode(null)
      setStatus('loading')

      try {
        await player.play(from, passage.text, token, { rate: preferences.speed, voice })
      } catch (error) {
        setStatus('idle')
        setErrorCode(error instanceof ApiError ? error.code : 'audio_failed')
        return
      }

      void highlightPassage(tabId, frameId, from)

      for (let ahead = 1; ahead <= prefetchAhead; ahead++) {
        const next = material?.passages[from + ahead]
        if (next) player.prefetch(from + ahead, next.text, token, voice)
      }
    },
    [frameId, material, player, preferences.speed, tabId, voice],
  )

  const advance = useCallback(
    (from: number) => {
      const next = from + 1
      if (!material?.passages[next]) {
        setStatus('idle')
        return
      }

      const pause = material.figureAfter.includes(from)
        ? preferences.figurePauseSeconds * 1000
        : 0
      if (pause <= 0) {
        void start(next)
        return
      }

      setStatus('loading')
      timerRef.current = setTimeout(() => void start(next), pause)
    },
    [material, preferences.figurePauseSeconds, start],
  )

  useEffect(() => {
    player.onStarted = () => setStatus('playing')
    player.onEnded = () => advance(index)
    player.onError = (code) => {
      setStatus('idle')
      setErrorCode(code)
    }
  }, [advance, index, player])

  useEffect(() => {
    player.setRate(preferences.speed)
  }, [player, preferences.speed])

  useEffect(() => {
    if (!materialId) return

    let cancelled = false
    void readProgress(materialId).then((saved) => {
      if (!cancelled && saved) setIndex(saved.index)
    })

    return () => {
      cancelled = true
    }
  }, [materialId])

  useEffect(() => {
    if (!materialId || status === 'idle') return
    void saveProgress(materialId, index)
  }, [index, materialId, status])

  // Closing the panel ends the Reading, so the panel takes the audio with it.
  useEffect(
    () => () => {
      if (timerRef.current) clearTimeout(timerRef.current)
      player.stop()
    },
    [player],
  )

  const toggle = useCallback(() => {
    if (status === 'playing') {
      player.pause()
      setStatus('paused')
      return
    }
    if (status === 'paused') {
      player.resume()
      setStatus('playing')
      return
    }
    void start(index)
  }, [index, player, start, status])

  const goTo = useCallback(
    (target: number) => {
      if (!material?.passages[target]) return
      void start(target)
    },
    [material, start],
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
