import { useEffect, useRef } from 'react'
import { Icon } from '@/components/icon'
import type { Passage } from '../material'

type Props = {
  passages: Passage[]
  figureAfter: number[]
  currentIndex: number
  onSelect: (index: number) => void
}

export function PassageList({ passages, figureAfter, currentIndex, onSelect }: Props) {
  const current = useRef<HTMLButtonElement | null>(null)

  // block: 'nearest' only scrolls when the Passage is out of view, so the list
  // does not fight the reader while they browse it.
  useEffect(() => {
    const reduce = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    current.current?.scrollIntoView({ block: 'nearest', behavior: reduce ? 'auto' : 'smooth' })
  }, [currentIndex])

  return (
    <ol className="flex min-h-24 flex-1 animate-fade-in flex-col gap-1 overflow-y-auto pr-1 motion-reduce:animate-none">
      {passages.map((passage, index) => (
        <li key={index}>
          <button
            ref={index === currentIndex ? current : undefined}
            type="button"
            onClick={() => onSelect(index)}
            aria-current={index === currentIndex || undefined}
            className={`flex w-full gap-2 rounded-sm border-l-2 py-1 pl-2 text-left text-[13px] leading-relaxed transition-colors duration-150 ${
              index === currentIndex
                ? 'animate-pulse-soft border-moss text-ink motion-reduce:animate-none'
                : 'border-transparent text-ink-muted hover:border-moss-deep hover:text-ink'
            }`}
          >
            <span className="mt-0.5 w-5 shrink-0 text-right text-[10px] tabular-nums text-ink-faint">
              {index + 1}
            </span>
            <span>
              {passage.text}
              {figureAfter.includes(index) ? (
                <span className="mt-1 flex items-center gap-1 text-[10.5px] text-bulb-dim">
                  <Icon name="image" className="size-3" />
                  Figura a seguir
                </span>
              ) : null}
            </span>
          </button>
        </li>
      ))}
    </ol>
  )
}
