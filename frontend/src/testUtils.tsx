import { vi } from 'vitest'
import { render } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import App from './App'
import { AuthProvider } from './auth/AuthContext'

export interface RouteSpec {
  status?: number
  body: unknown
}

// mockApi stubs fetch with method-aware routes ("POST /api/x" or plain
// "/api/x" for GET). The most specific (longest) matching path wins, and a
// fresh Response is built per call so bodies are never double-consumed.
export function mockApi(routes: Record<string, RouteSpec>) {
  const fetchMock = vi.fn().mockImplementation((input: RequestInfo | URL, init?: RequestInit) => {
    const url = typeof input === 'string' ? input : input.toString()
    const method = init?.method ?? 'GET'
    const match = Object.entries(routes)
      .filter(([key]) => {
        const [routeMethod, routePath] = key.includes(' ') ? key.split(' ') : ['GET', key]
        return routeMethod === method && url.startsWith(routePath)
      })
      .sort(([a], [b]) => b.length - a.length)[0]
    const { status = 200, body } = match?.[1] ?? { status: 404, body: { error: 'not found' } }
    // 204 responses must have a null body or the Response constructor throws.
    const payload = status === 204 ? null : JSON.stringify(body)
    return Promise.resolve(
      new Response(payload, {
        status,
        headers: { 'Content-Type': 'application/json' },
      }),
    )
  })
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

export const testUser = {
  id: 'u1',
  email: 'mohit@example.com',
  name: 'Mohit',
  createdAt: '2026-07-14',
}

export function renderLoggedInAt(path: string) {
  window.localStorage.setItem('istream_token', 'test-token')
  return render(
    <MemoryRouter initialEntries={[path]}>
      <AuthProvider>
        <App />
      </AuthProvider>
    </MemoryRouter>,
  )
}