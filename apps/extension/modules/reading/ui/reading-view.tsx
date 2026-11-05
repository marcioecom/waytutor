import type { Workspace } from '@/modules/account/api'
import { useScan } from '@/modules/material/hooks/use-scan'
import { MaterialView } from '@/modules/material/ui/material-view'
import { usePreferences } from '@/modules/preferences/hooks/use-preferences'
import { PreferencesBar } from '@/modules/preferences/ui/preferences-bar'
import { useReading } from '../hooks/use-reading'
import { PlayerBar } from './player-bar'

export function ReadingView({ workspace }: { workspace: Workspace }) {
  const scan = useScan()
  const { preferences, update } = usePreferences()

  const material = scan.data?.material ?? null

  const reading = useReading({
    tabId: scan.data?.tabId ?? -1,
    frameId: scan.data?.frameId ?? -1,
    materialId: scan.data?.id ?? null,
    material,
    preferences,
    defaultVoice: workspace.defaultVoice,
  })

  return (
    <div className="flex min-h-0 flex-1 flex-col gap-4">
      <MaterialView currentIndex={reading.index} onSelect={reading.goTo} />

      <PlayerBar
        status={reading.status}
        index={reading.index}
        total={material?.passages.length ?? 0}
        errorCode={reading.errorCode}
        onToggle={reading.toggle}
        onPrevious={reading.previous}
        onNext={reading.next}
      />

      <PreferencesBar voices={workspace.voices} preferences={preferences} onChange={update} />
    </div>
  )
}
