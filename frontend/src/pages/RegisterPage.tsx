import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { authApi } from '../api/auth'
import { useAuthStore } from '../store/auth'

export function RegisterPage() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const setAuth = useAuthStore((s) => s.setAuth)
  const navigate = useNavigate()

  async function submit(e: FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const res = await authApi.register(username, password)
      setAuth(res.user, res.token)
      navigate('/')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="page-narrow" style={{ paddingTop: 80 }}>
      <div className="card">
        <h2 style={{ marginBottom: 6, fontSize: '1.5rem', fontWeight: 800 }}>Create account</h2>
        <p className="text-muted text-sm" style={{ marginBottom: 28 }}>Start your Wikipedia collection</p>

        <form onSubmit={submit} className="form-stack">
          <div className="form-group">
            <label className="form-label">Username</label>
            <input className="form-input" placeholder="your_username" value={username}
              onChange={(e) => setUsername(e.target.value)} required autoFocus />
          </div>
          <div className="form-group">
            <label className="form-label">Password</label>
            <input className="form-input" type="password" placeholder="••••••••" value={password}
              onChange={(e) => setPassword(e.target.value)} required />
          </div>
          {error && <div className="alert alert-danger">{error}</div>}
          <button className="btn btn-primary btn-full" type="submit" disabled={loading}>
            {loading ? 'Creating account...' : 'Register'}
          </button>
        </form>

        <hr className="divider" />
        <p className="text-sm text-muted" style={{ textAlign: 'center' }}>
          Have an account? <Link to="/login">Log in</Link>
        </p>
      </div>
    </div>
  )
}
