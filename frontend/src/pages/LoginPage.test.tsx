import { describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import App from '../App'
import { AuthProvider } from '../auth/AuthContext'

function renderAt(path: string) {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <AuthProvider>
        <App />
      </AuthProvider>
    </MemoryRouter>,
  )
}

function mockFetch(status: number, body: unknown) {
  const fetchMock = vi.fn().mockResolvedValue(
    new Response(JSON.stringify(body), {
      status,
      headers: { 'Content-Type': 'application/json' },
    }),
  )
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

describe('LoginPage', () => {
  it('renders the login form', () => {
    renderAt('/login')
    expect(screen.getByRole('heading', { name: 'Log in' })).toBeInTheDocument()
    expect(screen.getByLabelText('Email')).toBeInTheDocument()
    expect(screen.getByLabelText('Password')).toBeInTheDocument()
  })

  it('shows the API error message when login fails', async () => {
    mockFetch(401, { error: 'invalid email or password' })
    renderAt('/login')

    await userEvent.type(screen.getByLabelText('Email'), 'mohit@example.com')
    await userEvent.type(screen.getByLabelText('Password'), 'wrongpassword')
    await userEvent.click(screen.getByRole('button', { name: 'Log in' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('invalid email or password')
  })

  it('stores the token and shows the dashboard after successful login', async () => {
    const fetchMock = mockFetch(200, {
      token: 'test-token',
      user: { id: '1', email: 'mohit@example.com', name: 'Mohit', createdAt: '2026-07-14' },
    })
    renderAt('/login')

    await userEvent.type(screen.getByLabelText('Email'), 'mohit@example.com')
    await userEvent.type(screen.getByLabelText('Password'), 'supersecret1')
    await userEvent.click(screen.getByRole('button', { name: 'Log in' }))

    expect(await screen.findByText('Welcome, Mohit')).toBeInTheDocument()
    expect(window.localStorage.getItem('istream_token')).toBe('test-token')
    expect(fetchMock).toHaveBeenCalledWith(
      '/api/auth/login',
      expect.objectContaining({ method: 'POST' }),
    )
  })

  it('redirects unauthenticated visitors from the dashboard to login', () => {
    renderAt('/')
    expect(screen.getByRole('heading', { name: 'Log in' })).toBeInTheDocument()
  })
})