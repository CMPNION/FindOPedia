import { useEffect, useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { quizApi, type Collection } from '../api/quiz'
import { RarityBadge } from '../components/RarityBadge'
import { useAuthStore } from '../store/auth'

export function CollectionPage() {
  const { username } = useParams<{ username: string }>()
  const currentUser = useAuthStore((s) => s.user)
  const target = username ?? currentUser?.username
  const isOwn = !username || username === currentUser?.username

  const [collection, setCollection] = useState<Collection | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const navigate = useNavigate()

  useEffect(() => {
    if (!target) return
    quizApi.getCollection(target)
      .then(setCollection)
      .catch(() => setError('User not found'))
      .finally(() => setLoading(false))
  }, [target])

  if (loading) return <div style={{ textAlign: 'center', marginTop: 80, color: 'var(--text-muted)' }}>Loading...</div>
  if (error) return (
    <div className="page-narrow" style={{ marginTop: 40 }}>
      <div className="alert alert-danger">{error}</div>
    </div>
  )
  if (!collection) return null

  return (
    <div className="page">
      <div style={{ marginBottom: 28 }}>
        <h2 style={{ fontSize: '1.5rem', fontWeight: 800 }}>
          {isOwn ? 'My Collection' : `${collection.username}'s Collection`}
        </h2>
        <p className="text-sm text-muted" style={{ marginTop: 4 }}>
          {collection.total} article{collection.total !== 1 ? 's' : ''} owned
        </p>
      </div>

      {collection.articles.length === 0 ? (
        <div className="card" style={{ textAlign: 'center', padding: 48 }}>
          <p style={{ fontSize: '2rem', marginBottom: 12 }}>📚</p>
          <p className="text-muted">No articles owned yet.</p>
          {isOwn && (
            <button className="btn btn-primary" style={{ marginTop: 20 }} onClick={() => navigate('/')}>
              Discover your first article →
            </button>
          )}
        </div>
      ) : (
        <div className="collection-grid">
          {collection.articles.map((item) => (
            <Link key={item.slug} to={`/articles/${item.slug}`} className="collection-card">
              <div style={{ minWidth: 0 }}>
                <div style={{ fontWeight: 600, fontSize: '.95rem', marginBottom: 4 }}>{item.title}</div>
                <div className="text-xs text-muted">
                  Claimed {new Date(item.claimed_at).toLocaleDateString()}
                </div>
              </div>
              <RarityBadge tier={item.rarity_tier} />
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}
