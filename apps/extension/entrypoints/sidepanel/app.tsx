import { Button } from '@/components/button'
import { Spinner } from '@/components/spinner'
import { useAccount } from '@/modules/account/hooks/use-account'
import { messageFor } from '@/modules/account/messages'
import { BlockedView } from '@/modules/account/ui/blocked-view'
import { LoginView } from '@/modules/account/ui/login-view'
import { SessionFooter } from '@/modules/account/ui/session-footer'
import { ReadingView } from '@/modules/reading/ui/reading-view'

export function App() {
  const account = useAccount()

  return (
    <main className="flex h-dvh flex-col gap-4 overflow-hidden px-4 py-4">
      <header className="flex items-baseline justify-between gap-3">
        <p className="text-[19px] font-semibold leading-none tracking-[-0.02em]">
          Way<span className="text-moss">Tutor</span>
        </p>
        <p className="text-[10px] tracking-[0.03em] text-ink-faint">professor de bolso</p>
      </header>

      <Screen account={account} />
    </main>
  )
}

function Screen({ account }: { account: ReturnType<typeof useAccount> }) {
  if (account.isPending) {
    return (
      <p className="flex h-5 items-center gap-2 text-sm text-ink-faint">
        <Spinner />
        Carregando…
      </p>
    )
  }

  if (account.isError) {
    return (
      <div className="flex flex-col gap-4">
        <p role="alert" className="text-sm text-danger">
          {messageFor(account.error)}
        </p>
        <Button onClick={() => account.refetch()}>Tentar de novo</Button>
      </div>
    )
  }

  switch (account.data.kind) {
    case 'anonymous':
      return <LoginView />
    case 'pending':
      return (
        <BlockedView
          title="Aguardando liberação"
          message="Sua conta foi criada e está na fila. Assim que for liberada, a leitura começa a funcionar."
          onRetry={() => account.refetch()}
        />
      )
    case 'paused':
      return (
        <BlockedView
          title="Conta pausada"
          message="Sua conta está pausada. Fale com quem administra o WayTutor."
          onRetry={() => account.refetch()}
        />
      )
    case 'active':
      return (
        <div className="flex min-h-0 flex-1 flex-col gap-4">
          <ReadingView workspace={account.data.workspace} />
          <SessionFooter account={account.data.workspace.account} />
        </div>
      )
  }
}
