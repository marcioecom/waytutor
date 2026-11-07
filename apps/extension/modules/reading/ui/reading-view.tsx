import { useEffect, useState } from 'react'
import { Button } from '@/components/button'
import { Capybara } from '@/components/capybara'
import { Spinner } from '@/components/spinner'
import type { Workspace } from '@/modules/account/api'
import { useScan } from '@/modules/material/hooks/use-scan'
import { PassageList } from '@/modules/material/ui/passage-list'
import { usePreferences } from '@/modules/preferences/hooks/use-preferences'
import { PreferencesBar } from '@/modules/preferences/ui/preferences-bar'
import { SummaryView } from '@/modules/summary/ui/summary-view'
import { useReading } from '../hooks/use-reading'
import { Stage } from './stage'
import { StatusBadge } from './status-badge'

// Roughly 900 characters per minute at 1x, the pace the voice lands on.
const charsPerMinute = 900

export function ReadingView({ workspace }: { workspace: Workspace }) {
  const scan = useScan()
  const { preferences, update } = usePreferences()
  const [showSummary, setShowSummary] = useState(false)

  const material = scan.data?.material ?? null

  const reading = useReading({
    tabId: scan.data?.tabId ?? -1,
    frameId: scan.data?.frameId ?? -1,
    materialId: scan.data?.id ?? null,
    material,
    preferences,
    defaultVoice: workspace.defaultVoice,
  })

  const remaining = material?.passages
    .slice(reading.index)
    .reduce((total, passage) => total + passage.text.length, 0)

  const remainingMinutes = Math.max(
    1,
    Math.round((remaining ?? 0) / (charsPerMinute * preferences.speed)),
  )

  const characters = material?.passages.reduce((total, passage) => total + passage.text.length, 0)
  const totalMinutes = Math.max(1, Math.round((characters ?? 0) / (charsPerMinute * preferences.speed)))

  useEffect(() => {
    const onKey = (event: KeyboardEvent) => {
      const target = event.target as HTMLElement | null
      if (target?.closest('select, input, textarea')) return

      if (event.code === 'Space') {
        event.preventDefault()
        reading.toggle()
        return
      }
      // Arrows belong to the Reading, not to the Summary that covers it.
      if (showSummary) return

      if (event.code === 'ArrowRight') {
        reading.next()
        return
      }
      if (event.code === 'ArrowLeft') {
        reading.previous()
      }
    }

    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [reading.toggle, reading.next, reading.previous, showSummary])

  const materialText = material?.passages.map((passage) => passage.text).join('\n\n') ?? ''

  if (showSummary) {
    return (
      <SummaryView
        materialId={scan.data?.id ?? null}
        materialText={materialText}
        status={reading.status}
        onSpeak={reading.speak}
        onToggle={reading.toggle}
        onClose={() => {
          reading.stop()
          setShowSummary(false)
        }}
      />
    )
  }

  if (scan.isPending) {
    return (
      <p className="flex flex-1 items-center justify-center gap-2 text-sm text-ink-faint">
        <Spinner />
        Lendo a página…
      </p>
    )
  }

  if (scan.isError) {
    return (
      <div className="flex flex-1 flex-col items-center justify-center gap-3 text-center">
        <p role="alert" className="max-w-[32ch] text-sm leading-relaxed text-danger">
          Não foi possível ler esta página.
        </p>
        <Button onClick={() => scan.refetch()}>Tentar de novo</Button>
      </div>
    )
  }

  if (!material) {
    return (
      <div className="flex flex-1 flex-col items-center justify-center gap-4 pb-10 text-center">
        <div className="relative">
          <div className="nook-halo absolute -inset-8 rounded-full opacity-45" />
          <Capybara className="relative h-24 w-[6.5rem]" />
        </div>
        <div className="flex max-w-[34ch] flex-col gap-1.5">
          <h1 className="text-xl font-semibold leading-tight tracking-[-0.02em] text-ink">
            Abra um tópico de texto
          </h1>
          <p className="text-[13px] leading-relaxed text-ink-muted">
            O painel procura o conteúdo dentro do visualizador do portal, ignora menus e lê em
            voz alta.
          </p>
        </div>
        <Button onClick={() => scan.refetch()} disabled={scan.isFetching}>
          {scan.isFetching ? 'Procurando…' : 'Procurar nesta página'}
        </Button>
      </div>
    )
  }

  return (
    <div className="flex min-h-0 flex-1 flex-col gap-3.5 overflow-y-auto">
      <div>
        <h1 className="line-clamp-3 text-[15px] font-semibold leading-[1.35] tracking-[-0.01em] text-ink">
          {material.title}
        </h1>
        <div className="mt-1 flex items-baseline justify-between gap-3">
          <p className="text-[10.5px] tabular-nums text-ink-faint">
            {material.passages.length} trechos <span className="text-line-2">·</span> ~{totalMinutes} min
          </p>
          <StatusBadge status={reading.status} errorCode={reading.errorCode} />
        </div>
      </div>

      <Stage
        status={reading.status}
        index={reading.index}
        total={material.passages.length}
        remainingMinutes={remainingMinutes}
        errorCode={reading.errorCode}
        onToggle={reading.toggle}
        onPrevious={reading.previous}
        onNext={reading.next}
      />

      <PassageList
        passages={material.passages}
        figureAfter={material.figureAfter}
        currentIndex={reading.index}
        onSelect={reading.goTo}
      />

      <div className="flex items-center gap-5">
        <Button
          variant="quiet"
          className="inline-flex min-h-9 items-center gap-1.5 text-xs"
          onClick={() => scan.refetch()}
          disabled={scan.isFetching}
        >
          {scan.isFetching ? <Spinner /> : null}
          {scan.isFetching ? 'Escaneando…' : 'Escanear de novo'}
        </Button>

        <Button
          variant="quiet"
          className="inline-flex min-h-9 items-center text-xs"
          onClick={() => setShowSummary(true)}
        >
          Resumo
        </Button>
      </div>

      <PreferencesBar voices={workspace.voices} preferences={preferences} onChange={update} />
    </div>
  )
}
