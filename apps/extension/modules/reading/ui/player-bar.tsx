import { Button } from '@/components/button'
import { Spinner } from '@/components/spinner'
import { messageForCode } from '@/modules/account/messages'
import type { ReadingStatus } from '../hooks/use-reading'

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

const labels: Record<ReadingStatus, string> = {
  idle: 'Ouvir',
  loading: 'Carregando…',
  playing: 'Pausar',
  paused: 'Continuar',
}

export function PlayerBar({
  status,
  index,
  total,
  remainingMinutes,
  errorCode,
  onToggle,
  onPrevious,
  onNext,
}: Props) {
  const percent = total === 0 ? 0 : Math.round(((index + 1) / total) * 100)

  return (
    <div className="flex flex-col gap-2 border-t border-line pt-3">
      <div className="h-0.5 w-full overflow-hidden rounded-full bg-night-4">
        <div
          className="h-full rounded-full bg-lamp transition-[width] duration-300"
          style={{ width: `${percent}%` }}
        />
      </div>

      <div className="flex items-center gap-2">
        <Button
          variant="ghost"
          className="px-2 py-1"
          onClick={onPrevious}
          disabled={index === 0}
          aria-label="Trecho anterior"
        >
          ◀
        </Button>

        <Button
          className="inline-flex flex-1 items-center justify-center gap-2"
          onClick={onToggle}
          disabled={total === 0 || status === 'loading'}
        >
          {status === 'loading' ? <Spinner /> : null}
          {labels[status]}
        </Button>

        <Button
          variant="ghost"
          className="px-2 py-1"
          onClick={onNext}
          disabled={index >= total - 1}
          aria-label="Próximo trecho"
        >
          ▶
        </Button>
      </div>

      <div className="flex items-baseline justify-between text-xs text-ink-faint">
        <span className="tabular-nums">
          {total === 0 ? '—' : `${index + 1} de ${total}`}
        </span>
        <span>{status === 'loading' ? '' : `~${remainingMinutes} min restantes`}</span>
      </div>

      {errorCode ? (
        <p role="alert" className="text-xs text-danger">
          {messageForCode(errorCode)}
        </p>
      ) : null}
    </div>
  )
}
