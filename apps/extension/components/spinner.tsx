const sizes = {
  sm: 'size-3 border',
  md: 'size-4 border-2',
} as const

type Props = {
  size?: keyof typeof sizes
  className?: string
}

export function Spinner({ size = 'sm', className = '' }: Props) {
  return (
    <span
      aria-hidden="true"
      className={`inline-block shrink-0 animate-spin rounded-full border-current border-t-transparent motion-reduce:animate-none ${sizes[size]} ${className}`}
    />
  )
}
