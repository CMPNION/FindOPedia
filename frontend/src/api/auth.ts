import { api } from './client'

export interface AuthUser {
  id: number
  username: string
}

interface AuthResponse {
  user: AuthUser
  token: string
}

export const authApi = {
  register: (username: string, password: string) =>
    api.post<AuthResponse>('/auth/register', { username, password }),
  login: (username: string, password: string) =>
    api.post<AuthResponse>('/auth/login', { username, password }),
  me: () => api.get<AuthUser & { created_at: string }>('/auth/me'),
}
