import { Icon } from '@/components/icon'
import type { Preferences } from '../storage'
import { VoiceSelect } from './voice-select'

const speeds = [0.75, 1, 1.25, 1.5, 1.75, 2]
const pauses = [0, 3, 5, 8, 10]

const selectClass =
  'h-9 w-full appearance-none rounded-md border border-line-2 bg-night-3 pl-2.5 pr-8 text-[12.5px] text-ink transition-colors duration-150 hover:border-ink-faint focus-visible:border-signal'

type Props = {
  voices: string[]
  preferences: Preferences
  onChange: (patch: Partial<Preferences>) => void
}

export function PreferencesBar({ voices, preferences, onChange }: Props) {
  const toggleFavorite = (voice: string) => {
    const next = preferences.favoriteVoices.includes(voice)
      ? preferences.favoriteVoices.filter((item) => item !== voice)
      : [...preferences.favoriteVoices, voice]

    onChange({ favoriteVoices: next })
  }

  return (
    <div className="flex flex-col gap-2">
      <div className="flex flex-col gap-1.5">
        <span id="speed-label" className="text-[11px] font-medium text-ink-muted">
          Velocidade
        </span>
        <div
          role="group"
          aria-labelledby="speed-label"
          className="flex gap-0.5 rounded-lg border border-line bg-night-2 p-0.5"
        >
          {speeds.map((speed) => {
            const active = speed === preferences.speed

            return (
              <button
                key={speed}
                type="button"
                aria-pressed={active}
                onClick={() => onChange({ speed })}
                className={`h-7.5 flex-1 rounded-[5px] text-[11px] tabular-nums transition-colors duration-150 ${
                  active
                    ? 'bg-moss font-semibold text-night'
                    : 'text-ink-muted hover:bg-night-3 hover:text-ink'
                }`}
              >
                {speed}×
              </button>
            )
          })}
        </div>
      </div>

      <details className="group rounded-md border border-line bg-night-2">
        <summary className="flex h-10 cursor-pointer list-none items-center gap-2 rounded-md px-3 transition-colors duration-150 hover:bg-night-3 [&::-webkit-details-marker]:hidden">
          <span className="text-xs font-medium text-ink">Ajustes da voz</span>
          <span className="ml-auto max-w-[45%] truncate text-[11px] text-ink-muted">
            {preferences.voice ?? 'Padrão'}
          </span>
          <Icon
            name="chevron"
            className="size-3.5 shrink-0 text-ink-muted transition-transform duration-150 group-open:rotate-90"
          />
        </summary>

        <div className="flex flex-col gap-3.5 border-t border-line px-3 py-3.5">
          <VoiceSelect
            voices={voices}
            value={preferences.voice}
            favorites={preferences.favoriteVoices}
            onSelect={(voice) => onChange({ voice })}
            onToggleFavorite={toggleFavorite}
          />

          <div className="flex flex-col gap-1.5">
            <label htmlFor="figure-pause" className="text-[11px] font-medium text-ink-muted">
              Pausa em figura
            </label>
            <div className="relative">
              <select
                id="figure-pause"
                className={selectClass}
                value={preferences.figurePauseSeconds}
                onChange={(event) =>
                  onChange({ figurePauseSeconds: Number(event.target.value) })
                }
              >
                {pauses.map((seconds) => (
                  <option key={seconds} value={seconds}>
                    {seconds === 0 ? 'Sem pausa' : `${seconds} s`}
                  </option>
                ))}
              </select>
              <Icon
                name="down"
                className="pointer-events-none absolute top-1/2 right-2.5 size-3.5 -translate-y-1/2 text-ink-faint"
              />
            </div>
          </div>
        </div>
      </details>
    </div>
  )
}
