import { Button } from '@/components/button'
import type { Account } from '../api'
import { useLogout } from '../hooks/use-account'

export function SessionFooter({ account }: { account: Account }) {
  const logout = useLogout()

  return (
    <footer className="flex items-center justify-between gap-2 border-t border-line pt-3 text-xs text-ink-faint">
      <span className="truncate">{account.name || account.email}</span>
      <Button
        variant="ghost"
        className="shrink-0 px-0 py-1"
        onClick={() => logout.mutate()}
        disabled={logout.isPending}
      >
        Sair
      </Button>
    </footer>
  )
}
