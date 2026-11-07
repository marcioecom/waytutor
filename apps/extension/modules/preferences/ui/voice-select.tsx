import { useEffect, useId, useRef, useState } from 'react'
import { Icon } from '@/components/icon'

type Props = {
  voices: string[]
  value: string | null
  favorites: string[]
  onSelect: (voice: string | null) => void
  onToggleFavorite: (voice: string) => void
}

export function VoiceSelect({ voices, value, favorites, onSelect, onToggleFavorite }: Props) {
  const [open, setOpen] = useState(false)
  const id = useId()
  const labelId = `${id}-label`
  const valueId = `${id}-value`
  const panelId = `${id}-options`
  const root = useRef<HTMLDivElement | null>(null)
  const trigger = useRef<HTMLButtonElement | null>(null)

  const choose = (voice: string | null) => {
    onSelect(voice)
    setOpen(false)
    trigger.current?.focus()
  }

  // The row moves between groups, so React remounts the star. Put focus back
  // on the same voice after the move.
  const toggleFavorite = (voice: string) => {
    onToggleFavorite(voice)
    requestAnimationFrame(() => {
      root.current?.querySelector<HTMLButtonElement>(`[data-voice="${CSS.escape(voice)}"]`)?.focus()
    })
  }

  useEffect(() => {
    if (!open) return

    const onPointerDown = (event: PointerEvent) => {
      if (!root.current?.contains(event.target as Node)) setOpen(false)
    }
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key !== 'Escape') return
      setOpen(false)
      trigger.current?.focus()
    }

    document.addEventListener('pointerdown', onPointerDown)
    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.removeEventListener('pointerdown', onPointerDown)
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [open])

  const currentFavorite = value !== null && favorites.includes(value)
  const favoriteVoices = voices.filter((voice) => favorites.includes(voice))
  const otherVoices = voices.filter((voice) => !favorites.includes(voice))

  const row = (voice: string | null, label: string) => {
    const isSelected = value === voice
    const isFavorite = voice !== null && favorites.includes(voice)

    return (
      <li key={label} className="flex items-center gap-0.5">
        <button
          type="button"
          aria-pressed={isSelected}
          onClick={() => choose(voice)}
          className={`flex h-9 flex-1 items-center gap-2 rounded-sm px-2 text-left text-[12.5px] transition-colors duration-150 ${
            isSelected
              ? 'bg-night-3 font-medium text-ink'
              : 'text-ink-muted hover:bg-night-2 hover:text-ink'
          }`}
        >
          <span className="truncate">{label}</span>
          {isSelected ? <Icon name="check" className="ml-auto size-3.5 shrink-0 text-signal" /> : null}
        </button>

        {voice === null ? (
          <span className="size-9 shrink-0" aria-hidden="true" />
        ) : (
          <button
            type="button"
            data-voice={voice}
            aria-label={`${isFavorite ? 'Desfavoritar' : 'Favoritar'} ${voice}`}
            aria-pressed={isFavorite}
            onClick={() => toggleFavorite(voice)}
            className="grid size-9 shrink-0 place-items-center rounded-sm text-ink-faint transition-colors duration-150 hover:bg-night-2 hover:text-ink"
          >
            <Star filled={isFavorite} className={isFavorite ? 'text-moss' : ''} />
          </button>
        )}
      </li>
    )
  }

  return (
    <div ref={root} className="flex flex-col gap-1.5">
      <span id={labelId} className="text-[11px] font-medium text-ink-muted">
        Voz
      </span>

      <div className="relative">
        <button
          ref={trigger}
          type="button"
          aria-expanded={open}
          aria-controls={panelId}
          aria-labelledby={`${labelId} ${valueId}`}
          onClick={() => setOpen(!open)}
          className="flex h-9 w-full items-center justify-between gap-2 rounded-md border border-line-2 bg-night-3 px-2.5 text-left text-[12.5px] text-ink transition-colors duration-150 hover:border-ink-faint"
        >
          <span id={valueId} className="truncate">
            {value ?? 'Padrão'}
          </span>
          <span className="flex shrink-0 items-center gap-1.5 text-ink-faint">
            {currentFavorite ? <Star filled className="size-3.5 text-moss" /> : null}
            <Icon
              name="down"
              className={`size-3.5 transition-transform duration-150 ${open ? 'rotate-180' : ''}`}
            />
          </span>
        </button>

        {open ? (
          <div
            id={panelId}
            role="group"
            aria-labelledby={labelId}
            className="mt-1 rounded-md border border-line bg-night p-1"
          >
            <div className="max-h-56 overflow-y-auto">
              <ul className="flex flex-col">{row(null, 'Padrão')}</ul>
              <div className="my-1 border-t border-line" />

              {favoriteVoices.length > 0 ? (
                <>
                  <p className="px-2 pt-1 pb-1 text-[10px] uppercase tracking-[0.12em] text-ink-faint">
                    Favoritas
                  </p>
                  <ul className="flex flex-col">{favoriteVoices.map((voice) => row(voice, voice))}</ul>
                </>
              ) : null}

              {otherVoices.length > 0 ? (
                <>
                  <p className="px-2 pt-2 pb-1 text-[10px] uppercase tracking-[0.12em] text-ink-faint">
                    Todas as vozes
                  </p>
                  <ul className="flex flex-col">
                    {otherVoices.map((voice) => row(voice, voice))}
                  </ul>
                </>
              ) : null}
            </div>
          </div>
        ) : null}
      </div>

      <span className="text-[11px] text-ink-faint">A troca vale a partir do próximo trecho.</span>
    </div>
  )
}

function Star({ filled, className = '' }: { filled: boolean; className?: string }) {
  return (
    <svg
      viewBox="0 0 24 24"
      className={`size-3.5 ${className}`}
      aria-hidden="true"
      fill={filled ? 'currentColor' : 'none'}
      stroke="currentColor"
      strokeWidth="1.8"
      strokeLinejoin="round"
    >
      <path d="m12 3 2.7 5.6 6.1.9-4.4 4.3 1 6.1-5.4-2.9-5.4 2.9 1-6.1L3.2 9.5l6.1-.9z" />
    </svg>
  )
}
