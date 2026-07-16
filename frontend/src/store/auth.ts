import { create } from 'zustand'
import type { AuthUser } from '../api/auth'

interface AuthState {
  user: AuthUser | null
  token: string | null
  ready: boolean
  setAuth: (user: AuthUser, token: string) => void
  logout: () => void
  setReady: (ready: boolean) => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: localStorage.getItem('token'),
  ready: false,
  setAuth: (user, token) => {
    localStorage.setItem('token', token)
    set({ user, token })
  },
  logout: () => {
    localStorage.removeItem('token')
    set({ user: null, token: null })
  },
  setReady: (ready) => set({ ready }),
}))
