import { Button } from '@/components/button'
import { useScan } from '../hooks/use-scan'

export function MaterialView() {
  const scan = useScan()

  if (scan.isPending) {
    return <p className="text-sm text-ink-faint">Lendo a página…</p>
  }

  if (scan.isError) {
    return (
      <div className="flex flex-col gap-3">
        <p role="alert" className="text-sm text-danger">
          Não foi possível ler esta página.
        </p>
        <Button onClick={() => scan.refetch()}>Tentar de novo</Button>
      </div>
    )
  }

  const material = scan.data
  if (!material) {
    return (
      <div className="flex flex-col gap-3">
        <p className="text-sm leading-relaxed text-ink-muted">
          Nenhum material encontrado nesta aba. Abra um tópico de leitura no portal.
        </p>
        <Button
          variant="ghost"
          className="self-start px-0"
          onClick={() => scan.refetch()}
          disabled={scan.isFetching}
        >
          Escanear de novo
        </Button>
      </div>
    )
  }

  const characters = material.passages.reduce((total, passage) => total + passage.text.length, 0)

  return (
    <div className="flex min-h-0 flex-1 flex-col gap-3">
      <div>
        <h1 className="font-serif text-base leading-snug text-ink">{material.title}</h1>
        <p className="text-xs text-ink-faint">
          {material.passages.length} trechos · {characters} caracteres
        </p>
      </div>

      <ol className="flex min-h-0 flex-1 flex-col gap-3 overflow-y-auto pr-1">
        {material.passages.map((passage, index) => (
          <li key={index} className="flex gap-2 text-sm leading-relaxed text-ink-muted">
            <span className="mt-0.5 shrink-0 text-xs tabular-nums text-ink-faint">{index + 1}</span>
            <span>
              {passage.text}
              {material.figureAfter.includes(index) ? (
                <span className="mt-1 block text-xs text-lamp">▢ figura</span>
              ) : null}
            </span>
          </li>
        ))}
      </ol>

      <Button
        variant="ghost"
        className="self-start px-0 py-1"
        onClick={() => scan.refetch()}
        disabled={scan.isFetching}
      >
        {scan.isFetching ? 'Escaneando…' : 'Escanear de novo'}
      </Button>
    </div>
  )
}
