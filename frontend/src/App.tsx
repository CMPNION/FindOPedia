import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useEffect } from 'react'
import { Navbar } from './components/Navbar'
import { ProtectedRoute } from './components/ProtectedRoute'
import { HomePage } from './pages/HomePage'
import { ArticlePage } from './pages/ArticlePage'
import { QuizPage } from './pages/QuizPage'
import { ResultPage } from './pages/ResultPage'
import { CollectionPage } from './pages/CollectionPage'
import { LoginPage } from './pages/LoginPage'
import { RegisterPage } from './pages/RegisterPage'
import { SettingsPage } from './pages/SettingsPage'
import { LeaderboardPage } from './pages/LeaderboardPage'
import { useAuthStore } from './store/auth'
import { authApi } from './api/auth'

function AuthLayout({ children }: { children: React.ReactNode }) {
  const { token, ready } = useAuthStore()
  if (!ready) return null
  if (token) return <Navigate to="/" replace />
  return <>{children}</>
}

function AppInner() {
  const { token, setAuth, logout, setReady } = useAuthStore()

  useEffect(() => {
    if (!token) {
      setReady(true)
      return
    }
    authApi.me()
      .then((user) => {
        setAuth(user, token)
        setReady(true)
      })
      .catch(() => {
        logout()
        setReady(true)
      })
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <>
      <Navbar />
      <Routes>
        <Route path="/login" element={<AuthLayout><LoginPage /></AuthLayout>} />
        <Route path="/register" element={<AuthLayout><RegisterPage /></AuthLayout>} />

        <Route path="/" element={<ProtectedRoute><HomePage /></ProtectedRoute>} />
        <Route path="/articles/:slug" element={<ProtectedRoute><ArticlePage /></ProtectedRoute>} />
        <Route path="/articles/:slug/quiz" element={<ProtectedRoute><QuizPage /></ProtectedRoute>} />
        <Route path="/result/:slug" element={<ProtectedRoute><ResultPage /></ProtectedRoute>} />
        <Route path="/users/:username" element={<ProtectedRoute><CollectionPage /></ProtectedRoute>} />
        <Route path="/settings" element={<ProtectedRoute><SettingsPage /></ProtectedRoute>} />
        <Route path="/leaderboard" element={<ProtectedRoute><LeaderboardPage /></ProtectedRoute>} />

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <AppInner />
    </BrowserRouter>
  )
}
