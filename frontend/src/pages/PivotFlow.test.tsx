import { describe, expect, it, vi } from 'vitest'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { mockApi, renderLoggedInAt, testUser } from '../testUtils'

const driverNode = {
  id: 'driver',
  question: "What's really driving the urge to change?",
  options: [
    { label: 'Burnt out - I need time', next: 'hours-fix' },
    { label: 'The work no longer fits', next: 'what-broke' },
  ],
}

const hoursFixNode = {
  id: 'hours-fix',
  question: 'If you had 20% more free time, would your current job be fine?',
  options: [
    { label: 'Yes - just the hours', outcome: 'reduce-hours' },
    { label: 'No - deeper than that', next: 'what-broke' },
  ],
}

const outcomeThread = {
  id: 't1',
  steps: [
    { nodeId: 'driver', option: 'Burnt out - I need time' },
    { nodeId: 'hours-fix', option: 'Yes - just the hours' },
  ],
  outcome: {
    id: 'reduce-hours',
    path: 'reduce-hours',
    title: 'Negotiate reduced hours where you are',
    tagline: 'Keep the job, reclaim your time.',
    whyNow: 'Flexible arrangements are normal in 2026.',
    plan: ['Audit two weeks of workload', 'Draft the business case', 'Propose a 3-month trial'],
    resources: [
      { title: 'Mid-career pivot guide', url: 'https://example.com/guide' },
      { title: 'Build a resume for the target role', url: '/resumes/new' },
    ],
  },
  createdAt: '2026-07-14T10:00:00Z',
  updatedAt: '2026-07-14T10:00:00Z',
}

describe('pivot thread', () => {
  it('asks the current question and advances on answer', async () => {
    mockApi({
      'GET /api/auth/me': { body: testUser },
      'GET /api/pivot/threads/t1': {
        body: { id: 't1', steps: [], current: driverNode, createdAt: '', updatedAt: '' },
      },
      'PUT /api/pivot/threads/t1': {
        body: {
          id: 't1',
          steps: [{ nodeId: 'driver', option: 'Burnt out - I need time' }],
          current: hoursFixNode,
          createdAt: '',
          updatedAt: '',
        },
      },
    })
    renderLoggedInAt('/pivot/t1')

    expect(await screen.findByText("What's really driving the urge to change?")).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Burnt out - I need time' }))

    // Next question appears; the answer lands in the trail with a fork button.
    expect(await screen.findByText(/20% more free time/)).toBeInTheDocument()
    expect(screen.getByText('Burnt out - I need time')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Fork from here' })).toBeInTheDocument()
  })

  it('renders the outcome with plan and resources', async () => {
    mockApi({
      'GET /api/auth/me': { body: testUser },
      'GET /api/pivot/threads/t1': { body: outcomeThread },
    })
    renderLoggedInAt('/pivot/t1')

    expect(await screen.findByText('Negotiate reduced hours where you are')).toBeInTheDocument()
    expect(screen.getByText('Flexible arrangements are normal in 2026.')).toBeInTheDocument()
    expect(screen.getByText('Audit two weeks of workload')).toBeInTheDocument()
    // External resource opens in a new tab; internal one is a router link.
    expect(screen.getByRole('link', { name: 'Mid-career pivot guide' })).toHaveAttribute('target', '_blank')
    expect(screen.getByRole('link', { name: 'Build a resume for the target role' })).toHaveAttribute(
      'href',
      '/resumes/new',
    )
  })

  it('forks at an earlier answer into a new thread', async () => {
    const fetchMock = mockApi({
      'GET /api/auth/me': { body: testUser },
      'GET /api/pivot/threads/t1': { body: outcomeThread },
      'POST /api/pivot/threads/t1/fork': {
        status: 201,
        body: { id: 't2', steps: [], forkedFrom: 't1', current: driverNode, createdAt: '', updatedAt: '' },
      },
      'GET /api/pivot/threads/t2': {
        body: { id: 't2', steps: [], forkedFrom: 't1', current: driverNode, createdAt: '', updatedAt: '' },
      },
    })
    renderLoggedInAt('/pivot/t1')

    const forkButtons = await screen.findAllByRole('button', { name: 'Fork from here' })
    expect(forkButtons).toHaveLength(2)
    await userEvent.click(forkButtons[0])

    // Lands on the new thread, back at the re-asked question.
    expect(await screen.findByText("What's really driving the urge to change?")).toBeInTheDocument()
    expect(screen.getByText(/Forked from an earlier exploration/)).toBeInTheDocument()

    const forkCall = fetchMock.mock.calls.find(([url]) => String(url).includes('/fork'))
    expect(forkCall).toBeDefined()
    expect(JSON.parse((forkCall![1] as RequestInit).body as string)).toEqual({ atStep: 0 })
  })
})

describe('pivot list', () => {
  it('lists explorations and starts a new one', async () => {
    mockApi({
      'GET /api/auth/me': { body: testUser },
      'GET /api/pivot/threads': {
        body: {
          threads: [
            outcomeThread,
            {
              id: 't3',
              steps: [{ nodeId: 'driver', option: 'The work no longer fits' }],
              forkedFrom: 't1',
              current: hoursFixNode,
              createdAt: '2026-07-14T09:00:00Z',
              updatedAt: '2026-07-14T09:00:00Z',
            },
          ],
        },
      },
      'POST /api/pivot/threads': {
        status: 201,
        body: { id: 't9', steps: [], current: driverNode, createdAt: '', updatedAt: '' },
      },
      'GET /api/pivot/threads/t9': {
        body: { id: 't9', steps: [], current: driverNode, createdAt: '', updatedAt: '' },
      },
    })
    renderLoggedInAt('/pivot')

    expect(await screen.findByText('Negotiate reduced hours where you are')).toBeInTheDocument()
    expect(screen.getByText('In progress')).toBeInTheDocument()
    expect(screen.getByText(/forked/)).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Start a new exploration' }))
    expect(await screen.findByText("What's really driving the urge to change?")).toBeInTheDocument()
  })

  it('deletes an exploration after confirming', async () => {
    mockApi({
      'GET /api/auth/me': { body: testUser },
      'GET /api/pivot/threads': { body: { threads: [outcomeThread] } },
      'DELETE /api/pivot/threads/t1': { status: 204, body: null },
    })
    vi.spyOn(window, 'confirm').mockReturnValue(true)
    renderLoggedInAt('/pivot')

    await screen.findByText('Negotiate reduced hours where you are')
    await userEvent.click(screen.getByRole('button', { name: 'Delete' }))
    expect(await screen.findByText(/No explorations yet/)).toBeInTheDocument()
  })
})