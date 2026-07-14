import { request } from './client'

export interface JobSearchFilters {
  keywords: string
  location?: string
  workplace?: string
  experience?: string
  jobType?: string
  postedWithin?: string
}

export interface JobSearchURL {
  provider: string
  url: string
}

export function getSearchURL(filters: JobSearchFilters) {
  const params = new URLSearchParams()
  for (const [key, value] of Object.entries(filters)) {
    if (value) params.set(key, value)
  }
  return request<JobSearchURL>(`/api/jobs/search-url?${params.toString()}`)
}