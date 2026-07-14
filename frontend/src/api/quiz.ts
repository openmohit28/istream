import { request } from './client'

export interface QuizQuestion {
  id: string
  text: string
}

export interface QuizScale {
  min: number
  max: number
  minLabel: string
  maxLabel: string
}

export interface QuestionsResponse {
  questions: QuizQuestion[]
  scale: QuizScale
}

export type RiasecKey = 'R' | 'I' | 'A' | 'S' | 'E' | 'C'

export const RIASEC_LABELS: Record<RiasecKey, string> = {
  R: 'Realistic (hands-on)',
  I: 'Investigative (analytical)',
  A: 'Artistic (creative)',
  S: 'Social (helping)',
  E: 'Enterprising (leading)',
  C: 'Conventional (organizing)',
}

export interface JobMatch {
  title: string
  hollandCode: string
  category: string
  demand: 'growing' | 'stable' | 'declining'
  aiRisk: 'low' | 'medium' | 'high'
  blurb: string
  fit: number
}

export interface QuizResult {
  id: string
  createdAt: string
  scores: Record<RiasecKey, number>
  matches: JobMatch[]
}

export interface QuizResultSummary {
  id: string
  createdAt: string
  scores: Record<RiasecKey, number>
  topMatch?: { title: string; fit: number }
}

export function getQuestions() {
  return request<QuestionsResponse>('/api/quiz/questions')
}

export function submitQuiz(answers: Record<string, number>) {
  return request<QuizResult>('/api/quiz/submit', {
    method: 'POST',
    body: JSON.stringify({ answers }),
  })
}

export function getResults() {
  return request<{ results: QuizResultSummary[] }>('/api/quiz/results')
}

export function getResult(id: string) {
  return request<QuizResult>(`/api/quiz/results/${id}`)
}