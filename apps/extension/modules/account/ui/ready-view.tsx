import { Button } from '@/components/button'
import type { Account } from '../api'
import { useLogout } from '../hooks/use-account'

export function ReadyView({ account }: { account: Account }) {
  const logout = useLogout()

  return (
    <div className="flex flex-col gap-4">
      <p className="text-sm text-ink-muted">
        Conectado como <span className="text-ink">{account.name || account.email}</span>
      </p>
      <p className="text-sm leading-relaxed text-ink-faint">
        Nenhum material aberto. Abra um tópico no portal para começar.
      </p>
      <Button variant="ghost" className="self-start px-0" onClick={() => logout.mutate()} disabled={logout.isPending}>
        Sair
      </Button>
    </div>
  )
}
