import { request } from './client'

export interface Contact {
  fullName: string
  email: string
  phone: string
  location: string
  linkedin: string
  website: string
}

export interface Experience {
  company: string
  title: string
  location: string
  startDate: string
  endDate: string
  current: boolean
  bullets: string[]
}

export interface Education {
  school: string
  degree: string
  field: string
  gradYear: string
}

export interface ResumeDocument {
  targetTitle: string
  jobDescription: string
  contact: Contact
  summary: string
  experience: Experience[]
  education: Education[]
  skills: string[]
  certifications: string[]
}

export interface Resume {
  id: string
  title: string
  data: ResumeDocument
  createdAt: string
  updatedAt: string
}

export interface ResumeSummary {
  id: string
  title: string
  createdAt: string
  updatedAt: string
}

export interface KeywordReport {
  score: number
  matched: string[]
  missing: string[]
}

export function emptyDocument(): ResumeDocument {
  return {
    targetTitle: '',
    jobDescription: '',
    contact: { fullName: '', email: '', phone: '', location: '', linkedin: '', website: '' },
    summary: '',
    experience: [],
    education: [],
    skills: [],
    certifications: [],
  }
}

export function createResume(doc: ResumeDocument) {
  return request<Resume>('/api/resumes', { method: 'POST', body: JSON.stringify(doc) })
}

export function updateResume(id: string, doc: ResumeDocument) {
  return request<Resume>(`/api/resumes/${id}`, { method: 'PUT', body: JSON.stringify(doc) })
}

export function listResumes() {
  return request<{ resumes: ResumeSummary[] }>('/api/resumes')
}

export function getResume(id: string) {
  return request<Resume>(`/api/resumes/${id}`)
}

export function deleteResume(id: string) {
  return request<null>(`/api/resumes/${id}`, { method: 'DELETE' })
}

export function keywordCheck(id: string, jobDescription: string) {
  return request<KeywordReport>(`/api/resumes/${id}/keyword-check`, {
    method: 'POST',
    body: JSON.stringify({ jobDescription }),
  })
}