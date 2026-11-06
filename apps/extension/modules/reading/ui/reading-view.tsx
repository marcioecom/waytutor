import { useEffect, useState } from 'react'
import type { Workspace } from '@/modules/account/api'
import { useScan } from '@/modules/material/hooks/use-scan'
import { MaterialView } from '@/modules/material/ui/material-view'
import { usePreferences } from '@/modules/preferences/hooks/use-preferences'
import { PreferencesBar } from '@/modules/preferences/ui/preferences-bar'
import { SummaryView } from '@/modules/summary/ui/summary-view'
import { useReading } from '../hooks/use-reading'
import { PlayerBar } from './player-bar'

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

  useEffect(() => {
    const onKey = (event: KeyboardEvent) => {
      const target = event.target as HTMLElement | null
      if (target?.closest('select, input, textarea')) return

      if (event.code === 'Space') {
        event.preventDefault()
        reading.toggle()
        return
      }
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
  }, [reading.toggle, reading.next, reading.previous])

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

  return (
    <div className="flex min-h-0 flex-1 flex-col gap-4">
      <MaterialView
        currentIndex={reading.index}
        onSelect={reading.goTo}
        onSummary={() => setShowSummary(true)}
      />

      <PlayerBar
        status={reading.status}
        index={reading.index}
        total={material?.passages.length ?? 0}
        remainingMinutes={remainingMinutes}
        errorCode={reading.errorCode}
        onToggle={reading.toggle}
        onPrevious={reading.previous}
        onNext={reading.next}
      />

      <PreferencesBar voices={workspace.voices} preferences={preferences} onChange={update} />
    </div>
  )
}
