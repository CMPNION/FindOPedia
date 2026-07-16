import { useLocation, useNavigate, useParams } from 'react-router-dom'
import type { QuizResult } from '../api/quiz'

export function ResultPage() {
  const { slug } = useParams<{ slug: string }>()
  const location = useLocation()
  const navigate = useNavigate()
  const result = location.state as QuizResult | null

  if (!result) {
    navigate('/')
    return null
  }

  const passed = result.status === 'passed'

  return (
    <div className="page-narrow">
      <div className="card" style={{ textAlign: 'center', padding: 40 }}>
        <div style={{ fontSize: '4rem', marginBottom: 16 }}>{passed ? '🏆' : '❌'}</div>

        <h2 style={{ fontSize: '1.8rem', fontWeight: 800, marginBottom: 8 }}>
          {passed ? 'Perfect Score!' : 'Not Quite'}
        </h2>

        <p style={{ fontSize: '1.1rem', color: 'var(--text-muted)', marginBottom: 24 }}>
          {result.correct_count} / {result.total_count} correct
        </p>

        <div style={{
          display: 'flex', justifyContent: 'center', marginBottom: 28,
          gap: 24, flexWrap: 'wrap'
        }}>
          <div style={{ textAlign: 'center' }}>
            <div style={{ fontSize: '2.5rem', fontWeight: 900, color: passed ? 'var(--success)' : 'var(--danger)' }}>
              {result.score}%
            </div>
            <div className="text-xs text-muted">score</div>
          </div>
          <div style={{ textAlign: 'center' }}>
            <div style={{ fontSize: '2.5rem', fontWeight: 900 }}>{result.correct_count}</div>
            <div className="text-xs text-muted">correct</div>
          </div>
          <div style={{ textAlign: 'center' }}>
            <div style={{ fontSize: '2.5rem', fontWeight: 900 }}>{result.total_count - result.correct_count}</div>
            <div className="text-xs text-muted">wrong</div>
          </div>
        </div>

        {passed && result.ownership_claimed && (
          <div className="alert alert-success" style={{ marginBottom: 20, textAlign: 'left' }}>
            <strong>You now own this article!</strong> It's permanently in your collection.
          </div>
        )}

        {passed && !result.ownership_claimed && result.owner && (
          <div className="alert alert-warning" style={{ marginBottom: 20, textAlign: 'left' }}>
            You passed — but <strong>{result.owner.username}</strong> already owns this article.
          </div>
        )}

        {!passed && (
          <div className="alert alert-danger" style={{ marginBottom: 20, textAlign: 'left' }}>
            100% required. One attempt per article — no retries.
          </div>
        )}

        <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap', justifyContent: 'center' }}>
          {slug && (
            <button className="btn btn-ghost" onClick={() => navigate(`/articles/${slug}`)}>
              Back to Article
            </button>
          )}
          <button className="btn btn-primary" onClick={() => navigate('/')}>
            Discover Another →
          </button>
        </div>
      </div>
    </div>
  )
}
