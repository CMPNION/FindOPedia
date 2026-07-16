import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { articlesApi } from '../api/articles'
import { useAuthStore } from '../store/auth'

export function HomePage() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const user = useAuthStore((s) => s.user)

  async function discover() {
    setLoading(true)
    setError('')
    try {
      const article = await articlesApi.getRandom()
      navigate(`/articles/${article.slug}`)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to fetch article')
      setLoading(false)
    }
  }

  return (
    <div style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
      <div className="hero">
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'center' }}>
          <span className="chip">Welcome back, {user?.username ?? 'explorer'}</span>
        </div>
        <h1 className="hero-title">
          Explore. Learn.<br />
          <span style={{ color: 'var(--primary)' }}>Own.</span>
        </h1>
        <p className="hero-sub">
          Discover random Wikipedia articles. Pass the quiz with a perfect score to claim permanent ownership.
        </p>
        <button
          className="btn btn-primary btn-lg"
          onClick={discover}
          disabled={loading}
        >
          {loading ? 'Finding article...' : 'Discover Article'}
        </button>
        {error && <p style={{ color: 'var(--danger)', marginTop: 16, fontSize: '.9rem' }}>{error}</p>}
      </div>

      <div style={{ width: '100%', maxWidth: 640, padding: '0 20px 60px', display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))', gap: 12 }}>
        {[
          { label: 'Common', desc: '< 5 KB', cls: 'badge-common' },
          { label: 'Uncommon', desc: '5–20 KB', cls: 'badge-uncommon' },
          { label: 'Rare', desc: '20–50 KB', cls: 'badge-rare' },
          { label: 'Epic', desc: '50–100 KB', cls: 'badge-epic' },
          { label: 'Legendary', desc: '100+ KB', cls: 'badge-legendary' },
        ].map((r) => (
          <div key={r.label} className="card" style={{ textAlign: 'center', padding: 16 }}>
            <span className={`badge ${r.cls}`}>{r.label}</span>
            <p className="text-xs text-muted" style={{ marginTop: 8 }}>{r.desc}</p>
          </div>
        ))}
      </div>
    </div>
  )
}
