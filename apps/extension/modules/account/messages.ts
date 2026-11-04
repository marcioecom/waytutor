import { ApiError } from './api'

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
}

export function messageFor(error: unknown): string {
  return error instanceof ApiError ? (messages[error.code] ?? fallback) : fallback
}
