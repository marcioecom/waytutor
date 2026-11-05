import { Button } from '@/components/button'
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
      <div className="flex items-center justify-between">
        <h1 className="font-serif text-base text-ink">Resumo</h1>
        <Button variant="ghost" className="px-0 py-1" onClick={onClose}>
          Voltar
        </Button>
      </div>

      {summary.isPending ? <p className="text-sm text-ink-faint">Preparando o resumo…</p> : null}

      {summary.isError ? (
        <div className="flex flex-col gap-3">
          <p role="alert" className="text-sm text-danger">
            {messageFor(summary.error)}
          </p>
          <Button onClick={() => summary.refetch()}>Tentar de novo</Button>
        </div>
      ) : null}

      {text ? (
        <p className="min-h-0 flex-1 overflow-y-auto whitespace-pre-line text-sm leading-relaxed text-ink-muted">
          {text}
        </p>
      ) : null}

      {text ? (
        <Button
          onClick={() => (status === 'idle' ? onSpeak(text) : onToggle())}
          disabled={status === 'loading'}
        >
          {labels[status]}
        </Button>
      ) : null}
    </div>
  )
}
