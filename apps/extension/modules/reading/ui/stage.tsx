import { Books } from '@/components/books'
import { Capybara, type Mood } from '@/components/capybara'
import { Icon } from '@/components/icon'
import { Lamp } from '@/components/lamp'
import { Spinner } from '@/components/spinner'
import { messageForCode } from '@/modules/account/messages'
import type { ReadingStatus } from '../hooks/use-reading'

const RING_RADIUS = 38
const RING_LENGTH = 2 * Math.PI * RING_RADIUS

const glowByMood: Record<Mood, string> = {
  reading: 'animate-lamp-breathe opacity-95 motion-reduce:animate-none',
  thinking: 'animate-lamp-warm opacity-80 motion-reduce:animate-none',
  waiting: 'opacity-45',
  error: 'opacity-30',
}

const ditherByMood: Record<Mood, string> = {
  reading: 'opacity-35',
  thinking: 'opacity-25',
  waiting: 'opacity-15',
  error: 'opacity-10',
}

const labels: Record<ReadingStatus, string> = {
  idle: 'Ouvir',
  loading: 'Preparando',
  playing: 'Pausar',
  paused: 'Continuar',
}

function moodFor(status: ReadingStatus, errorCode: string | null): Mood {
  if (errorCode) return 'error'
  if (status === 'playing') return 'reading'
  if (status === 'loading') return 'thinking'
  return 'waiting'
}

type Props = {
  status: ReadingStatus
  index: number
  total: number
  remainingMinutes: number
  errorCode: string | null
  onToggle: () => void
  onPrevious: () => void
  onNext: () => void
}

export function Stage({
  status,
  index,
  total,
  remainingMinutes,
  errorCode,
  onToggle,
  onPrevious,
  onNext,
}: Props) {
  const mood = moodFor(status, errorCode)
  const percent = total === 0 ? 0 : Math.round(((index + 1) / total) * 100)
  const label = errorCode ? 'Ouvir' : labels[status]

  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-center justify-center gap-2.5">
        <div className="relative flex h-[clamp(84px,24vw,96px)] w-[clamp(112px,34vw,132px)] shrink-0 items-end justify-end">
          <div className={`nook-halo absolute -inset-x-3 -top-4 -bottom-1.5 ${glowByMood[mood]}`} />
          <div
            className={`nook-dither absolute top-1/2 left-1/2 size-[clamp(130px,38vw,170px)] -translate-x-1/2 -translate-y-1/2 rounded-full ${ditherByMood[mood]}`}
          />
          <Lamp
            mood={mood}
            className="absolute bottom-0 left-0 h-[clamp(72px,22vw,90px)] w-[clamp(28px,9vw,36px)]"
          />
          <Books className="absolute bottom-0 left-[clamp(19px,6.5vw,26px)] h-[clamp(17px,5.2vw,21px)] w-[clamp(19px,6vw,24px)]" />
          <Capybara
            mood={mood}
            className="relative h-[clamp(80px,24vw,96px)] w-[clamp(84px,26vw,104px)]"
          />
          <div className="nook-line absolute inset-x-0 bottom-0 h-px" />
        </div>

        <div className="flex items-center gap-2">
          <button
            type="button"
            aria-label="Trecho anterior"
            title="Trecho anterior (←)"
            onClick={onPrevious}
            disabled={index === 0}
            className="grid size-[clamp(32px,9vw,36px)] shrink-0 place-items-center rounded-full border border-line-2 text-ink-muted transition duration-150 ease-soft hover:border-moss hover:bg-night-2 hover:text-ink disabled:opacity-40"
          >
            <Icon name="prev" className="size-4" />
          </button>

          <button
            type="button"
            aria-label={label}
            title={`${label} (espaço)`}
            onClick={onToggle}
            disabled={total === 0 || status === 'loading'}
            className="group relative grid size-[clamp(76px,20vw,82px)] shrink-0 place-items-center rounded-full disabled:opacity-50"
          >
            <span className="absolute -top-1.5 left-1/2 h-1.5 w-2.5 -translate-x-1/2 rounded-xs bg-moss transition-colors duration-150 group-hover:bg-moss-bright" />
            <svg className="absolute inset-0 size-full -rotate-90" viewBox="0 0 88 88" aria-hidden="true">
              <circle className="fill-none stroke-night-4" cx="44" cy="44" r={RING_RADIUS} strokeWidth="3.5" />
              <circle
                className="fill-none stroke-signal drop-shadow-[0_0_4px_#b8c54673] transition-[stroke-dashoffset] duration-500 ease-soft motion-reduce:transition-none"
                cx="44"
                cy="44"
                r={RING_RADIUS}
                strokeWidth="3.5"
                strokeLinecap="round"
                strokeDasharray={RING_LENGTH}
                strokeDashoffset={RING_LENGTH * (1 - percent / 100)}
              />
            </svg>
            <span className="relative grid size-[clamp(52px,14vw,56px)] place-items-center rounded-full border border-line-2 bg-night-3 text-moss transition duration-150 ease-soft group-hover:bg-moss group-hover:text-night group-hover:shadow-[0_0_26px_rgba(169,204,143,0.3)] group-active:scale-95">
              {status === 'loading' ? <Spinner size="md" /> : <Icon name={status === 'playing' ? 'pause' : 'play'} className="size-5" />}
            </span>
          </button>

          <button
            type="button"
            aria-label="Próximo trecho"
            title="Próximo trecho (→)"
            onClick={onNext}
            disabled={index >= total - 1}
            className="grid size-[clamp(32px,9vw,36px)] shrink-0 place-items-center rounded-full border border-line-2 text-ink-muted transition duration-150 ease-soft hover:border-moss hover:bg-night-2 hover:text-ink disabled:opacity-40"
          >
            <Icon name="next" className="size-4" />
          </button>
        </div>
      </div>

      <p className="flex items-baseline gap-1.5 text-[11px] tabular-nums text-ink-muted">
        <span>Trecho</span>
        <b className="font-semibold text-moss">{total === 0 ? '—' : index + 1}</b>
        <span>de {total || '—'}</span>
        <span className="ml-auto whitespace-nowrap">
          {total === 0
            ? '—'
            : index >= total - 1
              ? 'último trecho'
              : `faltam ~${remainingMinutes} min`}
        </span>
      </p>

      {errorCode ? (
        <p role="alert" className="text-xs text-danger">
          {messageForCode(errorCode)}
        </p>
      ) : null}
    </div>
  )
}
