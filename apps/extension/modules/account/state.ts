import { ApiError } from '@/modules/api/errors'
import { fetchWorkspace, type Workspace } from './api'
import { clearToken, readToken } from './storage'

export type AccountState =
  | { kind: 'anonymous' }
  | { kind: 'pending' }
  | { kind: 'paused' }
  | { kind: 'active'; workspace: Workspace }

export async function loadAccountState(): Promise<AccountState> {
  const token = await readToken()
  if (!token) return { kind: 'anonymous' }

  try {
    return { kind: 'active', workspace: await fetchWorkspace(token) }
  } catch (error) {
    if (!(error instanceof ApiError)) throw error

    if (error.code === 'unauthenticated') {
      await clearToken()
      return { kind: 'anonymous' }
    }
    if (error.code === 'approval_required') return { kind: 'pending' }
    if (error.code === 'account_paused') return { kind: 'paused' }

    throw error
  }
}
