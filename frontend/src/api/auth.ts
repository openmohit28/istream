import { request } from './client'

export interface User {
  id: string
  email: string
  name: string
  createdAt: string
}

export interface AuthResponse {
  token: string
  user: User
}

export function register(email: string, name: string, password: string) {
  return request<AuthResponse>('/api/auth/register', {
    method: 'POST',
    body: JSON.stringify({ email, name, password }),
  })
}

export function login(email: string, password: string) {
  return request<AuthResponse>('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export function me() {
  return request<User>('/api/auth/me')
}