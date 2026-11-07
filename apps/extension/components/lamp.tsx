import type { Mood } from './capybara'

const bulbByMood: Record<Mood, string> = {
  reading: 'var(--color-bulb)',
  thinking: 'var(--color-bulb)',
  waiting: 'var(--color-bulb-dim)',
  error: 'var(--color-bulb-off)',
}

export function Lamp({ mood = 'waiting', className = '' }: { mood?: Mood; className?: string }) {
  const bulb = bulbByMood[mood]

  return (
    <svg className={className} viewBox="0 0 40 96" aria-hidden="true">
      <ellipse cx="16" cy="93" rx="11" ry="3" fill="var(--color-brass)" opacity=".42" />
      <path
        d="M16 91 L12 54 L22 36"
        stroke="var(--color-brass)"
        strokeWidth="2.6"
        fill="none"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <circle cx="22" cy="36" r="2.2" fill="var(--color-brass)" />
      <path d="M14 26 L32 30 L28 44 L10 40 Z" fill="var(--color-moss-deep)" />
      <path d="M14 26 L32 30" stroke="var(--color-moss)" strokeWidth="1.4" opacity=".5" />
      <path d="M28 43.5 L10 39.5" stroke={bulb} strokeWidth="2.4" strokeLinecap="round" />
      <circle cx="19" cy="42" r="3.2" fill={bulb} />
    </svg>
  )
}
