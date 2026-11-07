import { Button } from '@/components/button'
import { Capybara } from '@/components/capybara'
import { Spinner } from '@/components/spinner'
import { useLogin } from '../hooks/use-account'
import { messageFor } from '../messages'

export function LoginView() {
  const login = useLogin()

  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-5 pb-8 text-center">
      <div className="relative">
        <div className="nook-halo absolute -inset-8 rounded-full opacity-45" />
        <Capybara className="relative h-24 w-[6.5rem]" />
      </div>

      <div className="flex max-w-[34ch] flex-col gap-1.5">
        <h1 className="text-xl font-semibold tracking-[-0.02em] text-ink">
          Entre para ouvir
        </h1>
        <p className="text-[13px] leading-relaxed text-ink-muted">
          Sua conta guarda a cota e o cache da leitura. Entre com o Google para ouvir o material
          do portal e receber o resumo do que foi lido.
        </p>
      </div>

      <Button
        className="inline-flex w-full items-center justify-center gap-2"
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

      <p className="max-w-[38ch] text-[10px] leading-relaxed text-ink-faint">
        O texto do material passa pelo servidor para gerar o áudio. Nada é publicado e não há
        telemetria.
      </p>
    </div>
  )
}
