import type { ButtonHTMLAttributes } from 'react'

const variants = {
  primary:
    'bg-moss px-4 py-2 text-night hover:bg-moss-bright active:scale-[.99] disabled:bg-night-3 disabled:text-ink-faint',
  ghost: 'px-4 py-2 text-ink-muted hover:text-ink disabled:text-ink-faint',
  quiet: 'text-ink-muted hover:text-moss disabled:text-ink-faint',
} as const

type Props = ButtonHTMLAttributes<HTMLButtonElement> & { variant?: keyof typeof variants }

export function Button({ variant = 'primary', className = '', ...props }: Props) {
  return (
    <button
      {...props}
      className={`rounded-md text-sm font-medium transition duration-150 ease-soft ${variants[variant]} ${className}`}
    />
  )
}
