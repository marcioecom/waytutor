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
    <div className="flex flex-col gap-4">
      <h1 className="text-base font-medium text-ink">{title}</h1>
      <p className="text-sm leading-relaxed text-ink-muted">{message}</p>

      <div className="flex items-center gap-2">
        <Button onClick={onRetry}>Verificar de novo</Button>
        <Button variant="ghost" onClick={() => logout.mutate()} disabled={logout.isPending}>
          Sair
        </Button>
      </div>
    </div>
  )
}
