import { describe, expect, it, vi } from 'vitest'
import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import App from '../App'
import { AuthProvider } from '../auth/AuthContext'

const user = { id: 'u1', email: 'mohit@example.com', name: 'Mohit', createdAt: '2026-07-14' }

const savedResume = {
  id: 'cv-1',
  title: 'Backend Engineer',
  createdAt: '2026-07-14T10:00:00Z',
  updatedAt: '2026-07-14T10:00:00Z',
  data: {
    targetTitle: 'Backend Engineer',
    jobDescription: 'Go PostgreSQL Kubernetes',
    contact: {
      fullName: 'Mohit Rawat',
      email: 'mohit@example.com',
      phone: '',
      location: 'Bengaluru',
      linkedin: '',
      website: '',
    },
    summary: 'Backend engineer building APIs.',
    experience: [
      {
        company: 'Acme',
        title: 'Software Engineer',
        location: '',
        startDate: 'Jan 2023',
        endDate: '',
        current: true,
        bullets: ['Built REST APIs with Go'],
      },
    ],
    education: [],
    skills: ['Go', 'PostgreSQL'],
    certifications: [],
  },
}

const keywordReport = {
  score: 67,
  matched: ['go', 'postgresql'],
  missing: ['kubernetes'],
}

interface RouteSpec {
  status?: number
  body: unknown
  method?: string
}

function mockApi(routes: Record<string, RouteSpec>) {
  const fetchMock = vi.fn().mockImplementation((input: RequestInfo | URL, init?: RequestInit) => {
    const url = typeof input === 'string' ? input : input.toString()
    const method = init?.method ?? 'GET'
    // Pick the most specific (longest) matching route so e.g.
    // "POST /api/resumes/cv-1/keyword-check" beats "POST /api/resumes".
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

describe('job search', () => {
  it('builds a LinkedIn search link from filters', async () => {
    mockApi({
      'GET /api/auth/me': { body: user },
      'GET /api/jobs/search-url': {
        body: { provider: 'linkedin', url: 'https://www.linkedin.com/jobs/search/?keywords=nurse&f_WT=2' },
      },
    })
    renderLoggedInAt('/jobs')

    await userEvent.type(await screen.findByLabelText('Job title or keywords'), 'nurse')
    await userEvent.selectOptions(screen.getByLabelText('Workplace'), 'remote')
    await userEvent.click(screen.getByRole('button', { name: 'Build my search' }))

    const link = await screen.findByRole('link', { name: 'Open your LinkedIn search' })
    expect(link).toHaveAttribute('href', 'https://www.linkedin.com/jobs/search/?keywords=nurse&f_WT=2')
    expect(link).toHaveAttribute('target', '_blank')
  })

  it('surfaces API errors', async () => {
    mockApi({
      'GET /api/auth/me': { body: user },
      'GET /api/jobs/search-url': { status: 400, body: { error: 'keywords are required' } },
    })
    renderLoggedInAt('/jobs')

    await userEvent.type(await screen.findByLabelText('Job title or keywords'), 'x')
    await userEvent.click(screen.getByRole('button', { name: 'Build my search' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('keywords are required')
  })
})

describe('resume builder', () => {
  it('walks the wizard and saves a resume', async () => {
    const fetchMock = mockApi({
      'GET /api/auth/me': { body: user },
      'POST /api/resumes': { status: 201, body: savedResume },
      'GET /api/resumes/cv-1': { body: savedResume },
      'POST /api/resumes/cv-1/keyword-check': { body: keywordReport },
    })
    renderLoggedInAt('/resumes/new')

    // Step 1: target role (Next disabled until filled).
    const next = await screen.findByRole('button', { name: 'Next' })
    expect(next).toBeDisabled()
    await userEvent.type(screen.getByLabelText('Target job title'), 'Backend Engineer')
    expect(next).toBeEnabled()
    await userEvent.click(next)

    // Step 2: contact (Next disabled until name + email).
    await userEvent.type(screen.getByLabelText('Full name'), 'Mohit Rawat')
    expect(screen.getByRole('button', { name: 'Next' })).toBeDisabled()
    await userEvent.type(screen.getByLabelText('Email'), 'mohit@example.com')
    await userEvent.click(screen.getByRole('button', { name: 'Next' }))

    // Step 3: summary.
    await userEvent.type(screen.getByLabelText('Professional summary'), 'Backend engineer building APIs.')
    await userEvent.click(screen.getByRole('button', { name: 'Next' }))

    // Step 4: experience - add one position.
    await userEvent.click(screen.getByRole('button', { name: 'Add experience' }))
    await userEvent.type(screen.getByLabelText('Job title'), 'Software Engineer')
    await userEvent.type(screen.getByLabelText('Company'), 'Acme')
    await userEvent.click(screen.getByRole('button', { name: 'Next' }))

    // Step 5: education - skip.
    await userEvent.click(screen.getByRole('button', { name: 'Next' }))

    // Step 6: skills.
    await userEvent.type(
      screen.getByLabelText('Skills (comma-separated - mirror the job description wording)'),
      'Go, PostgreSQL',
    )
    await userEvent.click(screen.getByRole('button', { name: 'Next' }))

    // Step 7: review + save.
    expect(screen.getByText('Ready to save')).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Save resume' }))

    // Lands on the preview with the saved content.
    expect(await screen.findByText('Mohit Rawat')).toBeInTheDocument()

    // The POST body carried the wizard data.
    const post = fetchMock.mock.calls.find(([, init]) => (init as RequestInit)?.method === 'POST')
    expect(post).toBeDefined()
    const payload = JSON.parse((post![1] as RequestInit).body as string)
    expect(payload.targetTitle).toBe('Backend Engineer')
    expect(payload.experience[0].company).toBe('Acme')
    expect(payload.skills).toEqual(['Go', 'PostgreSQL'])
  })
})

describe('resume preview', () => {
  it('renders the ATS sheet and keyword report', async () => {
    mockApi({
      'GET /api/auth/me': { body: user },
      'GET /api/resumes/cv-1': { body: savedResume },
      'POST /api/resumes/cv-1/keyword-check': { body: keywordReport },
    })
    renderLoggedInAt('/resumes/cv-1')

    const sheet = within(await screen.findByRole('article', { name: 'resume' }))
    expect(sheet.getByText('Mohit Rawat')).toBeInTheDocument()
    expect(sheet.getByText('Backend Engineer')).toBeInTheDocument()
    expect(sheet.getByText('Built REST APIs with Go')).toBeInTheDocument()
    expect(sheet.getByText('Go, PostgreSQL')).toBeInTheDocument()

    // Keyword report auto-runs because the resume stored a job description.
    expect(await screen.findByText('67%')).toBeInTheDocument()
    expect(screen.getByText('kubernetes')).toBeInTheDocument()
  })
})

describe('resume list', () => {
  it('lists resumes and deletes on confirm', async () => {
    mockApi({
      'GET /api/auth/me': { body: user },
      'GET /api/resumes': {
        body: { resumes: [{ id: 'cv-1', title: 'Backend Engineer', createdAt: '2026-07-14', updatedAt: '2026-07-14' }] },
      },
      'DELETE /api/resumes/cv-1': { status: 204, body: null },
    })
    vi.spyOn(window, 'confirm').mockReturnValue(true)
    renderLoggedInAt('/resumes')

    expect(await screen.findByText('Backend Engineer')).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Delete' }))
    expect(await screen.findByText(/No resumes yet/)).toBeInTheDocument()
  })
})