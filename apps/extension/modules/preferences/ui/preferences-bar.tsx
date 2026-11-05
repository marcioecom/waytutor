import type { Preferences } from '../storage'

const speeds = [0.75, 1, 1.25, 1.5, 2]
const pauses = [0, 3, 5, 8, 10]

const selectClass =
  'rounded-sm border border-line bg-night-2 px-2 py-1 text-sm text-ink outline-none focus:border-lamp'

type Props = {
  voices: string[]
  preferences: Preferences
  onChange: (patch: Partial<Preferences>) => void
}

export function PreferencesBar({ voices, preferences, onChange }: Props) {
  return (
    <div className="grid grid-cols-2 gap-2 text-xs text-ink-faint">
      <label className="flex flex-col gap-1">
        Voz
        <select
          className={selectClass}
          value={preferences.voice ?? ''}
          onChange={(event) => onChange({ voice: event.target.value || null })}
        >
          <option value="">Padrão</option>
          {voices.map((voice) => (
            <option key={voice} value={voice}>
              {voice}
            </option>
          ))}
        </select>
      </label>

      <label className="flex flex-col gap-1">
        Velocidade
        <select
          className={selectClass}
          value={preferences.speed}
          onChange={(event) => onChange({ speed: Number(event.target.value) })}
        >
          {speeds.map((speed) => (
            <option key={speed} value={speed}>
              {speed.toFixed(2).replace(/0$/, '')}×
            </option>
          ))}
        </select>
      </label>

      <label className="col-span-2 flex flex-col gap-1">
        Pausa em figura
        <select
          className={selectClass}
          value={preferences.figurePauseSeconds}
          onChange={(event) => onChange({ figurePauseSeconds: Number(event.target.value) })}
        >
          {pauses.map((seconds) => (
            <option key={seconds} value={seconds}>
              {seconds === 0 ? 'sem pausa' : `${seconds} s`}
            </option>
          ))}
        </select>
      </label>
    </div>
  )
}
