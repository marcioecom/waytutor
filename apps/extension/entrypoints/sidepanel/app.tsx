import type { ReactNode } from 'react'
import { Button } from '@/components/button'
import { useAccount } from '@/modules/account/hooks/use-account'
import { messageFor } from '@/modules/account/messages'
import { BlockedView } from '@/modules/account/ui/blocked-view'
import { LoginView } from '@/modules/account/ui/login-view'
import { ReadyView } from '@/modules/account/ui/ready-view'

export function App() {
  const account = useAccount()

  return (
    <main className="flex min-h-dvh flex-col gap-6 bg-night px-5 py-6">
      <header>
        <p className="font-serif text-xl text-lamp">WayTutor</p>
        <p className="text-sm text-ink-faint">professor de bolso</p>
      </header>

      <Screen account={account} />
    </main>
  )
}

function Screen({ account }: { account: ReturnType<typeof useAccount> }) {
  if (account.isPending) {
    return <p className="text-sm text-ink-faint">Carregando…</p>
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
      return <ReadyView account={account.data.account} />
  }
}
