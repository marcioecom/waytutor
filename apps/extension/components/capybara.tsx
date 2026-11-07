// The mascot is a state machine: the eyes and mouth mirror what the Reading is
// doing, never decoration. Ported from v1's capybara.
export type Mood = 'reading' | 'waiting' | 'thinking' | 'error'

const eye = 'origin-center [transform-box:fill-box] animate-capy-blink [animation-delay:-2.4s] motion-reduce:animate-none'

export function Capybara({ mood = 'waiting', className = '' }: { mood?: Mood; className?: string }) {
  const closed = mood === 'thinking'
  const eyeY = mood === 'reading' ? 59 : 56
  const eyeR = mood === 'error' ? 4.2 : 3

  return (
    <svg className={`overflow-visible ${className}`} viewBox="0 22 120 110" aria-hidden="true">
      <path
        d="M30 106h60c4 0 7 3 7 7v6c0 7-6 13-13 13H36c-7 0-13-6-13-13v-6c0-4 3-7 7-7z"
        fill="var(--color-leaf)"
      />
      <path d="M48 106h24l-12 16z" fill="var(--color-capy-mark)" opacity=".9" />

      <rect x="29" y="30" width="10" height="9" rx="4" fill="var(--color-capy-ear)" />
      <rect x="81" y="30" width="10" height="9" rx="4" fill="var(--color-capy-ear)" />

      <path
        d="M18 48c0-9 6-15 15-15h54c9 0 15 6 15 15v44c0 12-9 21-21 21H39c-12 0-21-9-21-21z"
        fill="var(--color-capy)"
      />
      <path
        d="M20 77c0-7 6-12 12-12h56c6 0 12 5 12 12v12c0 15-11 26-26 26H46c-15 0-26-11-26-26z"
        fill="var(--color-capy-2)"
      />

      <rect x="45" y="84" width="30" height="13" rx="6" fill="var(--color-capy-ink)" />

      {mood === 'error' ? (
        <path
          d="M52 104q4-4 8 0t8 0"
          fill="none"
          stroke="var(--color-capy-ink)"
          strokeWidth="2.2"
          strokeLinecap="round"
        />
      ) : mood === 'reading' ? (
        <path
          d="M54 104h12"
          fill="none"
          stroke="var(--color-capy-ink)"
          strokeWidth="2.2"
          strokeLinecap="round"
        />
      ) : (
        <path
          d="M53 103q7 4 14 0"
          fill="none"
          stroke="var(--color-capy-ink)"
          strokeWidth="2.2"
          strokeLinecap="round"
        />
      )}

      <circle
        cx="39"
        cy="56"
        r="9"
        fill="rgba(255,220,160,.05)"
        stroke="var(--color-capy-ink)"
        strokeWidth="2.4"
      />
      <circle
        cx="81"
        cy="56"
        r="9"
        fill="rgba(255,220,160,.05)"
        stroke="var(--color-capy-ink)"
        strokeWidth="2.4"
      />
      <path d="M48 56h24" stroke="var(--color-capy-ink)" strokeWidth="2.4" strokeLinecap="round" />

      {closed ? (
        <>
          <path
            d="M34 56.5q5 4.5 10 0"
            fill="none"
            stroke="var(--color-capy-ink)"
            strokeWidth="2.4"
            strokeLinecap="round"
          />
          <path
            d="M76 56.5q5 4.5 10 0"
            fill="none"
            stroke="var(--color-capy-ink)"
            strokeWidth="2.4"
            strokeLinecap="round"
          />
        </>
      ) : (
        <>
          <circle className={eye} cx="39" cy={eyeY} r={eyeR} fill="var(--color-capy-ink)" />
          <circle className={eye} cx="81" cy={eyeY} r={eyeR} fill="var(--color-capy-ink)" />
        </>
      )}
    </svg>
  )
}
