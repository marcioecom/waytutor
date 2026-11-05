const apiBaseUrl = import.meta.env.WXT_API_URL ?? 'http://localhost:8080'

export function apiUrl(path: string): string {
  return `${apiBaseUrl}${path}`
}
