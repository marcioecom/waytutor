import { useEffect, useRef } from 'react'
import { Button } from '@/components/button'
import { Spinner } from '@/components/spinner'
import { useScan } from '../hooks/use-scan'

type Props = {
  currentIndex: number
  onSelect: (index: number) => void
  onSummary: () => void
}

export function MaterialView({ currentIndex, onSelect, onSummary }: Props) {
  const scan = useScan()
  const current = useRef<HTMLButtonElement | null>(null)

  // block: 'nearest' only scrolls when the Passage is out of view, so the list
  // does not fight the reader while they browse it.
  useEffect(() => {
    current.current?.scrollIntoView({ block: 'nearest', behavior: 'smooth' })
  }, [currentIndex])

  if (scan.isPending) {
    return (
      <p className="flex h-5 items-center gap-2 text-sm text-ink-faint">
        <Spinner />
        Lendo a página…
      </p>
    )
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

  const material = scan.data?.material
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

      <ol className="flex min-h-0 flex-1 animate-fade-in flex-col gap-2 overflow-y-auto pr-1 motion-reduce:animate-none">
        {material.passages.map((passage, index) => (
          <li key={index}>
            <button
              ref={index === currentIndex ? current : undefined}
              type="button"
              onClick={() => onSelect(index)}
              aria-current={index === currentIndex || undefined}
              className={`flex w-full gap-2 rounded-sm border-l-2 py-1 pl-2 text-left text-sm leading-relaxed transition ${
                index === currentIndex
                  ? 'animate-pulse-soft border-lamp text-ink motion-reduce:animate-none'
                  : 'border-transparent text-ink-muted hover:border-line-2 hover:text-ink'
              }`}
            >
              <span className="mt-0.5 shrink-0 text-xs tabular-nums text-ink-faint">
                {index + 1}
              </span>
              <span>
                {passage.text}
                {material.figureAfter.includes(index) ? (
                  <span className="mt-1 block text-xs text-lamp">▢ figura</span>
                ) : null}
              </span>
            </button>
          </li>
        ))}
      </ol>

      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          className="inline-flex items-center gap-2 px-0 py-1"
          onClick={() => scan.refetch()}
          disabled={scan.isFetching}
        >
          {scan.isFetching ? <Spinner /> : null}
          {scan.isFetching ? 'Escaneando…' : 'Escanear de novo'}
        </Button>

        <Button variant="ghost" className="px-0 py-1" onClick={onSummary}>
          Resumo
        </Button>
      </div>
    </div>
  )
}
