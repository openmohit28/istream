import { request } from './client'

export interface PivotStep {
  nodeId: string
  option: string
}

export interface PivotOption {
  label: string
  next?: string
  outcome?: string
}

export interface PivotNode {
  id: string
  question: string
  options: PivotOption[]
}

export interface PivotResource {
  title: string
  url: string
}

export interface PivotOutcome {
  id: string
  path: string
  title: string
  tagline: string
  whyNow: string
  plan: string[]
  resources: PivotResource[]
}

export interface PivotThread {
  id: string
  steps: PivotStep[]
  forkedFrom?: string
  current?: PivotNode
  outcome?: PivotOutcome
  createdAt: string
  updatedAt: string
}

export function createThread() {
  return request<PivotThread>('/api/pivot/threads', { method: 'POST' })
}

export function listThreads() {
  return request<{ threads: PivotThread[] }>('/api/pivot/threads')
}

export function getThread(id: string) {
  return request<PivotThread>(`/api/pivot/threads/${id}`)
}

export function answerThread(id: string, steps: PivotStep[]) {
  return request<PivotThread>(`/api/pivot/threads/${id}`, {
    method: 'PUT',
    body: JSON.stringify({ steps }),
  })
}

// forkThread keeps the first atStep answers in a NEW thread, so the
// original exploration stays intact.
export function forkThread(id: string, atStep: number) {
  return request<PivotThread>(`/api/pivot/threads/${id}/fork`, {
    method: 'POST',
    body: JSON.stringify({ atStep }),
  })
}

export function deleteThread(id: string) {
  return request<null>(`/api/pivot/threads/${id}`, { method: 'DELETE' })
}