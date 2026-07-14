const TOKEN_KEY = 'istream_token'

// window.localStorage (not the bare global): Node 22+ defines its own
// localStorage global that shadows the browser one under Vitest/jsdom.
export function getToken(): string | null {
  return window.localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string) {
  window.localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken() {
  window.localStorage.removeItem(TOKEN_KEY)
}

export class ApiError extends Error {
  status: number

  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

export async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }
  const token = getToken()
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(path, { ...options, headers })
  const body = await res.json().catch(() => null)

  if (!res.ok) {
    throw new ApiError(res.status, body?.error ?? `Request failed (${res.status})`)
  }
  return body as T
}