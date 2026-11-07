export function Books({ className = '' }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 26 22" aria-hidden="true">
      <rect x="1" y="17.6" width="24" height="4.2" rx="1" fill="#40566b" />
      <rect x="1" y="16.9" width="24" height="0.8" rx="0.4" fill="#c9b28c" opacity=".5" />
      <rect x="2.4" y="13" width="21" height="4.2" rx="1" fill="#d5c2a0" />
      <rect x="2.4" y="12.3" width="21" height="0.8" rx="0.4" fill="#e3d5bb" opacity=".7" />
      <rect x="2" y="8.4" width="19" height="4.2" rx="1" fill="#7d3b32" />
      <rect x="2" y="7.7" width="19" height="0.8" rx="0.4" fill="#c9b28c" opacity=".5" />
      <g transform="rotate(-3.5 13 5)">
        <rect x="4.5" y="3.6" width="15" height="4.2" rx="1" fill="#e9dcc4" />
        <rect x="4.5" y="3.6" width="15" height="1.1" rx="0.5" fill="#f4e8d2" />
        <rect x="4.5" y="6.8" width="15" height="1" rx="0.5" fill="#bda886" />
      </g>
    </svg>
  )
}
