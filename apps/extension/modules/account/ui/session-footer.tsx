import { Button } from '@/components/button'
import type { Account } from '../api'
import { useLogout } from '../hooks/use-account'

export function SessionFooter({ account }: { account: Account }) {
  const logout = useLogout()

  return (
    <footer className="flex items-center justify-between gap-2 border-t border-line pt-3 text-[10px] tracking-[0.03em] text-ink-faint">
      <span className="truncate">{account.name || account.email}</span>
      <Button
        variant="quiet"
        className="inline-flex min-h-9 shrink-0 items-center text-[10px]"
        onClick={() => logout.mutate()}
        disabled={logout.isPending}
      >
        Sair
      </Button>
    </footer>
  )
}
