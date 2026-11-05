export class ApiError extends Error {
  constructor(readonly code: string) {
    super(code)
    this.name = 'ApiError'
  }
}

export async function codeOf(response: Response): Promise<string> {
  const body = (await response.json().catch(() => null)) as { error?: { code?: unknown } } | null
  const code = body?.error?.code

  return typeof code === 'string' ? code : 'unknown_error'
}
