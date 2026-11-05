import { Button } from '@/components/button'
import { messageForCode } from '@/modules/account/messages'
import type { ReadingStatus } from '../hooks/use-reading'

type Props = {
  status: ReadingStatus
  index: number
  total: number
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

export function PlayerBar({ status, index, total, errorCode, onToggle, onPrevious, onNext }: Props) {
  return (
    <div className="flex flex-col gap-2 border-t border-line pt-3">
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
          className="flex-1"
          onClick={onToggle}
          disabled={total === 0 || status === 'loading'}
        >
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

        <span className="w-12 text-right text-xs tabular-nums text-ink-faint">
          {total === 0 ? '—' : `${index + 1}/${total}`}
        </span>
      </div>

      {errorCode ? (
        <p role="alert" className="text-xs text-danger">
          {messageForCode(errorCode)}
        </p>
      ) : null}
    </div>
  )
}
