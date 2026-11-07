import { Button } from '@/components/button'
import { Spinner } from '@/components/spinner'
import { messageFor } from '@/modules/account/messages'
import type { ReadingStatus } from '@/modules/reading/hooks/use-reading'
import { useSummary } from '../hooks/use-summary'

type Props = {
  materialId: string | null
  materialText: string
  status: ReadingStatus
  onSpeak: (text: string) => void
  onToggle: () => void
  onClose: () => void
}

const labels: Record<ReadingStatus, string> = {
  idle: 'Ouvir resumo',
  loading: 'Preparando…',
  playing: 'Pausar',
  paused: 'Continuar',
}

export function SummaryView({ materialId, materialText, status, onSpeak, onToggle, onClose }: Props) {
  const summary = useSummary(materialId, materialText)
  const text = summary.data ?? ''

  return (
    <div className="flex min-h-0 flex-1 flex-col gap-3">
      <div className="flex items-baseline justify-between gap-3">
        <h1 className="text-lg font-semibold tracking-[-0.02em] text-ink">Resumo</h1>
        <Button
          variant="quiet"
          className="inline-flex min-h-9 items-center text-xs"
          onClick={onClose}
        >
          Voltar
        </Button>
      </div>

      {summary.isPending ? (
        <p className="flex h-5 items-center gap-2 text-sm text-ink-faint">
          <Spinner />
          Preparando o resumo…
        </p>
      ) : null}

      {summary.isError ? (
        <div className="flex flex-col gap-3">
          <p role="alert" className="text-sm text-danger">
            {messageFor(summary.error)}
          </p>
          <Button onClick={() => summary.refetch()}>Tentar de novo</Button>
        </div>
      ) : null}

      {text ? (
        <p className="min-h-0 flex-1 overflow-y-auto whitespace-pre-line text-[13px] leading-relaxed text-ink-muted">
          {text}
        </p>
      ) : null}

      {text ? (
        <Button
          className="inline-flex items-center gap-2"
          onClick={() => (status === 'idle' ? onSpeak(text) : onToggle())}
          disabled={status === 'loading'}
        >
          {status === 'loading' ? <Spinner /> : null}
          {labels[status]}
        </Button>
      ) : null}
    </div>
  )
}
