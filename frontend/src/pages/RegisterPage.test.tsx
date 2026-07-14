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

async function fillForm(name: string, email: string, password: string) {
  await userEvent.type(screen.getByLabelText('Name'), name)
  await userEvent.type(screen.getByLabelText('Email'), email)
  await userEvent.type(screen.getByLabelText('Password'), password)
  await userEvent.click(screen.getByRole('button', { name: 'Create account' }))
}

describe('RegisterPage', () => {
  it('renders the registration form', () => {
    renderAt('/register')
    expect(screen.getByRole('heading', { name: 'Create account' })).toBeInTheDocument()
    expect(screen.getByLabelText('Name')).toBeInTheDocument()
    expect(screen.getByLabelText('Email')).toBeInTheDocument()
    expect(screen.getByLabelText('Password')).toBeInTheDocument()
  })

  it('rejects short passwords without calling the API', async () => {
    const fetchMock = mockFetch(201, {})
    renderAt('/register')

    await fillForm('Mohit', 'mohit@example.com', 'short')

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Password must be at least 8 characters',
    )
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('shows the API error when the email is already registered', async () => {
    mockFetch(409, { error: 'email already registered' })
    renderAt('/register')

    await fillForm('Mohit', 'mohit@example.com', 'supersecret1')

    expect(await screen.findByRole('alert')).toHaveTextContent('email already registered')
  })

  it('stores the token and shows the dashboard after successful registration', async () => {
    mockFetch(201, {
      token: 'new-token',
      user: { id: '1', email: 'mohit@example.com', name: 'Mohit', createdAt: '2026-07-14' },
    })
    renderAt('/register')

    await fillForm('Mohit', 'mohit@example.com', 'supersecret1')

    expect(await screen.findByText('Welcome, Mohit')).toBeInTheDocument()
    expect(window.localStorage.getItem('istream_token')).toBe('new-token')
  })
})