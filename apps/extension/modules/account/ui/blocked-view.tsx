import { Button } from '@/components/button'
import { useLogout } from '../hooks/use-account'

type Props = {
  title: string
  message: string
  onRetry: () => void
}

export function BlockedView({ title, message, onRetry }: Props) {
  const logout = useLogout()

  return (
    <div className="flex flex-1 flex-col justify-center gap-4">
      <h1 className="text-xl font-semibold tracking-[-0.02em] text-ink">{title}</h1>
      <p className="max-w-[36ch] text-[13px] leading-relaxed text-ink-muted">{message}</p>

      <div className="mt-1 flex items-center gap-2">
        <Button onClick={onRetry}>Verificar de novo</Button>
        <Button variant="ghost" onClick={() => logout.mutate()} disabled={logout.isPending}>
          Sair
        </Button>
      </div>
    </div>
  )
}
