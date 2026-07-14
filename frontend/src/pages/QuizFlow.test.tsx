import { describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import App from '../App'
import { AuthProvider } from '../auth/AuthContext'

const user = { id: 'u1', email: 'mohit@example.com', name: 'Mohit', createdAt: '2026-07-14' }

const questionsBody = {
  questions: [
    { id: 'q1', text: 'First statement?' },
    { id: 'q2', text: 'Second statement?' },
    { id: 'q3', text: 'Third statement?' },
  ],
  scale: { min: 1, max: 5, minLabel: 'Strongly disagree', maxLabel: 'Strongly agree' },
}

const resultBody = {
  id: 'res-1',
  createdAt: '2026-07-14T10:00:00Z',
  scores: { R: 10, I: 100, A: 20, S: 30, E: 40, C: 50 },
  matches: [
    {
      title: 'AI / Machine Learning Engineer',
      hollandCode: 'IRE',
      category: 'Technology',
      demand: 'growing',
      aiRisk: 'low',
      blurb: 'Build and deploy AI systems.',
      fit: 93,
    },
    {
      title: 'Graphic Designer',
      hollandCode: 'AER',
      category: 'Creative',
      demand: 'declining',
      aiRisk: 'high',
      blurb: 'Production design is automating.',
      fit: 41,
    },
  ],
}

// mockApi routes fetch calls by URL prefix, returning a fresh Response each
// time so bodies are never double-consumed.
function mockApi(routes: Record<string, { status?: number; body: unknown }>) {
  const fetchMock = vi.fn().mockImplementation((input: RequestInfo | URL) => {
    const url = typeof input === 'string' ? input : input.toString()
    const match = Object.entries(routes).find(([path]) => url.startsWith(path))
    const { status = 200, body } = match?.[1] ?? { status: 404, body: { error: 'not found' } }
    return Promise.resolve(
      new Response(JSON.stringify(body), {
        status,
        headers: { 'Content-Type': 'application/json' },
      }),
    )
  })
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

function renderLoggedInAt(path: string) {
  window.localStorage.setItem('istream_token', 'test-token')
  return render(
    <MemoryRouter initialEntries={[path]}>
      <AuthProvider>
        <App />
      </AuthProvider>
    </MemoryRouter>,
  )
}

describe('quiz flow', () => {
  it('walks through every question and lands on the results page', async () => {
    mockApi({
      '/api/auth/me': { body: user },
      '/api/quiz/questions': { body: questionsBody },
      '/api/quiz/submit': { status: 201, body: resultBody },
      '/api/quiz/results/res-1': { body: resultBody },
    })
    renderLoggedInAt('/test')

    // Answer all three questions; selection auto-advances.
    await userEvent.click(await screen.findByRole('button', { name: 'Agree' }))
    expect(screen.getByText('Second statement?')).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Strongly agree' }))
    expect(screen.getByText('Third statement?')).toBeInTheDocument()

    // Submit is disabled until the last question is answered.
    const submit = screen.getByRole('button', { name: 'See my results' })
    expect(submit).toBeDisabled()
    await userEvent.click(screen.getByRole('button', { name: 'Neutral' }))
    expect(submit).toBeEnabled()

    await userEvent.click(submit)

    // Results page renders profile and matches.
    expect(await screen.findByText('Your interest profile')).toBeInTheDocument()
    expect(screen.getByText('AI / Machine Learning Engineer')).toBeInTheDocument()
    expect(screen.getByText('93% fit')).toBeInTheDocument()
    expect(screen.getByText('Growing demand')).toBeInTheDocument()
    expect(screen.getByText('Declining demand')).toBeInTheDocument()
    expect(screen.getByText('High AI exposure')).toBeInTheDocument()
  })

  it('tracks progress and supports going back', async () => {
    mockApi({
      '/api/auth/me': { body: user },
      '/api/quiz/questions': { body: questionsBody },
    })
    renderLoggedInAt('/test')

    const progress = await screen.findByRole('progressbar')
    expect(progress).toHaveAttribute('aria-valuenow', '0')

    await userEvent.click(screen.getByRole('button', { name: 'Agree' }))
    expect(progress).toHaveAttribute('aria-valuenow', '1')
    expect(screen.getByText('Second statement?')).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Back' }))
    expect(screen.getByText('First statement?')).toBeInTheDocument()
    // The earlier answer is remembered.
    expect(screen.getByRole('button', { name: 'Agree' })).toHaveAttribute('aria-pressed', 'true')
  })

  it('shows an error when a result cannot be loaded', async () => {
    mockApi({
      '/api/auth/me': { body: user },
      '/api/quiz/results/missing': { status: 404, body: { error: 'result not found' } },
    })
    renderLoggedInAt('/results/missing')

    expect(await screen.findByRole('alert')).toHaveTextContent('result not found')
  })

  it('lists past results on the dashboard', async () => {
    mockApi({
      '/api/auth/me': { body: user },
      '/api/quiz/results': {
        body: {
          results: [
            {
              id: 'res-1',
              createdAt: '2026-07-14T10:00:00Z',
              scores: resultBody.scores,
              topMatch: { title: 'AI / Machine Learning Engineer', fit: 93 },
            },
          ],
        },
      },
    })
    renderLoggedInAt('/')

    expect(await screen.findByText('Your past results')).toBeInTheDocument()
    expect(
      screen.getByText(/AI \/ Machine Learning Engineer \(93% fit\)/),
    ).toBeInTheDocument()
  })
})