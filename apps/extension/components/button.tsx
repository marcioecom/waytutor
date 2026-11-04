import type { ButtonHTMLAttributes } from 'react'

const variants = {
  primary: 'bg-lamp text-night hover:bg-lamp-bright',
  ghost: 'text-ink-muted hover:text-ink',
} as const

type Props = ButtonHTMLAttributes<HTMLButtonElement> & { variant?: keyof typeof variants }

export function Button({ variant = 'primary', className = '', ...props }: Props) {
  return (
    <button
      {...props}
      className={`rounded-md px-4 py-2 text-sm font-medium transition disabled:opacity-50 ${variants[variant]} ${className}`}
    />
  )
}
