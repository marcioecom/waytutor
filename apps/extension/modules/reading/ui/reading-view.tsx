import { MaterialView } from '@/modules/material/ui/material-view'
import type { Passage } from '@/modules/material/material'
import { useScan } from '@/modules/material/hooks/use-scan'
import { useReading } from '../hooks/use-reading'
import { PlayerBar } from './player-bar'

const noPassages: Passage[] = []

export function ReadingView() {
  const scan = useScan()
  const passages = scan.data?.material.passages ?? noPassages

  const reading = useReading({
    tabId: scan.data?.tabId ?? -1,
    frameId: scan.data?.frameId ?? -1,
    passages,
  })

  return (
    <div className="flex min-h-0 flex-1 flex-col gap-4">
      <MaterialView currentIndex={reading.index} onSelect={reading.goTo} />

      <PlayerBar
        status={reading.status}
        index={reading.index}
        total={passages.length}
        errorCode={reading.errorCode}
        onToggle={reading.toggle}
        onPrevious={reading.previous}
        onNext={reading.next}
      />
    </div>
  )
}
