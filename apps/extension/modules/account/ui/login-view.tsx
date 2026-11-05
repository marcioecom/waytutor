import { messageFor } from '../messages'
import { useLogin } from '../hooks/use-account'
import { Button } from '@/components/button'
import { Spinner } from '@/components/spinner'

export function LoginView() {
  const login = useLogin()

  return (
    <div className="flex flex-col gap-4">
      <p className="text-sm leading-relaxed text-ink-muted">
        Entre com sua conta Google para ouvir o material do portal em voz alta e receber o resumo
        do que foi lido.
      </p>

      <Button
        className="inline-flex items-center justify-center gap-2"
        onClick={() => login.mutate()}
        disabled={login.isPending}
      >
        {login.isPending ? <Spinner /> : null}
        {login.isPending ? 'Abrindo o Google…' : 'Entrar com Google'}
      </Button>

      {login.isError ? (
        <p role="alert" className="text-sm text-danger">
          {messageFor(login.error)}
        </p>
      ) : null}
    </div>
  )
}
