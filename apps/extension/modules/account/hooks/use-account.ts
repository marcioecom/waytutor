import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ApiError } from '@/modules/api/errors'
import { apiUrl, endSession, exchangeLoginCode } from '../api'
import { loadAccountState, type AccountState } from '../state'
import { clearToken, readToken, saveToken } from '../storage'

const accountQueryKey = ['account'] as const

export function useAccount() {
  return useQuery({ queryKey: accountQueryKey, queryFn: loadAccountState })
}

export function useLogin() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async () => {
      const redirectUri = chrome.identity.getRedirectURL()
      const callback = await chrome.identity.launchWebAuthFlow({
        url: `${apiUrl('/auth/google')}?redirect_uri=${encodeURIComponent(redirectUri)}`,
        interactive: true,
      })

      const code = callback ? new URL(callback).searchParams.get('code') : null
      if (!code) throw new ApiError('login_cancelled')

      const { token } = await exchangeLoginCode(code)
      await saveToken(token)
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: accountQueryKey }),
  })
}

export function useLogout() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async () => {
      const token = await readToken()
      if (token) await endSession(token)
      await clearToken()
    },
    onSuccess: () =>
      queryClient.setQueryData(accountQueryKey, { kind: 'anonymous' } satisfies AccountState),
  })
}
