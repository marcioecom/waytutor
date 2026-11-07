export type IconName = 'prev' | 'next' | 'play' | 'pause' | 'chevron' | 'down' | 'check' | 'image'

export function Icon({ name, className }: { name: IconName; className?: string }) {
  if (name === 'image') {
    return (
      <svg
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
        aria-hidden="true"
        className={className}
      >
        <rect x="3" y="4" width="18" height="16" rx="2" />
        <circle cx="9" cy="10" r="1.5" />
        <path d="m21 15-5-5L5 20" />
      </svg>
    )
  }

  if (name === 'play') {
    return (
      <svg viewBox="0 0 24 24" aria-hidden="true" className={className}>
        <path d="M8.5 5.6v12.8L19 12z" fill="currentColor" />
      </svg>
    )
  }

  if (name === 'pause') {
    return (
      <svg viewBox="0 0 24 24" aria-hidden="true" className={className}>
        <rect x="7.5" y="5.5" width="3.2" height="13" rx="1.1" fill="currentColor" />
        <rect x="13.3" y="5.5" width="3.2" height="13" rx="1.1" fill="currentColor" />
      </svg>
    )
  }

  const path =
    name === 'down'
      ? 'M6 9l6 6 6-6'
      : name === 'chevron'
        ? 'M9 6l6 6-6 6'
        : name === 'check'
          ? 'M5 12.5l4.5 4.5L19 7.5'
          : name === 'prev'
            ? 'M15 6l-6 6 6 6'
            : 'M9 6l6 6-6 6'

  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
      className={className}
    >
      <path d={path} />
    </svg>
  )
}
