import { ApiError } from '@/modules/api/errors'

const fallback = 'Algo deu errado. Tente de novo.'

const messages: Record<string, string> = {
  approval_required: 'Sua conta ainda não foi liberada.',
  account_paused: 'Sua conta está pausada.',
  unauthenticated: 'Sua sessão expirou. Entre de novo.',
  invalid_code: 'O login não foi concluído. Tente de novo.',
  invalid_state: 'O login não foi concluído. Tente de novo.',
  invalid_redirect_uri: 'O login não foi concluído. Tente de novo.',
  oauth_failed: 'O Google não confirmou o login. Tente de novo.',
  login_cancelled: 'Login cancelado.',
  quota_exceeded: 'Sua cota do mês acabou.',
  rate_limited: 'Muitos pedidos seguidos. Espere um minuto.',
  provider_quota: 'O serviço de voz está indisponível agora.',
  provider_error: 'Não foi possível gerar o áudio. Tente de novo.',
  provider_auth: 'Não foi possível gerar o áudio. Tente de novo.',
  empty_text: 'O trecho está vazio.',
  text_too_long: 'O trecho é longo demais para uma leitura.',
}

export function messageForCode(code: string): string {
  return messages[code] ?? fallback
}

export function messageFor(error: unknown): string {
  return error instanceof ApiError ? messageForCode(error.code) : fallback
}
