import type { ReadingStatus } from '../hooks/use-reading'

const labels: Record<ReadingStatus, string> = {
  idle: 'Pronto',
  loading: 'Preparando',
  playing: 'Lendo',
  paused: 'Pausado',
}

const tones = {
  live: 'text-ink',
  quiet: 'text-ink-muted',
  error: 'text-danger',
} as const

export function StatusBadge({
  status,
  errorCode,
}: {
  status: ReadingStatus
  errorCode: string | null
}) {
  const live = status === 'playing' || status === 'loading'
  const tone = errorCode ? tones.error : live ? tones.live : tones.quiet

  return (
    <span
      role="status"
      className={`inline-flex shrink-0 items-center gap-1.5 text-[10px] uppercase tracking-[0.14em] ${tone}`}
    >
      <span
        className={`size-1.5 rounded-full ${
          errorCode
            ? 'bg-danger'
            : live
              ? 'animate-status-pulse bg-signal motion-reduce:animate-none'
              : 'bg-line-2'
        }`}
      />
      {errorCode ? 'Erro' : labels[status]}
    </span>
  )
}
