import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { quizApi, type Leaderboards } from '../api/quiz'

const RARITIES = ['common', 'uncommon', 'rare', 'epic', 'legendary']
const TAB_LABELS: Record<string, string> = {
  total: 'All articles',
  common: 'Common',
  uncommon: 'Uncommon',
  rare: 'Rare',
  epic: 'Epic',
  legendary: 'Legendary',
}

function rankClass(i: number) {
  if (i === 0) return 'top1'
  if (i === 1) return 'top2'
  if (i === 2) return 'top3'
  return ''
}

function rankLabel(i: number) {
  if (i === 0) return '🥇'
  if (i === 1) return '🥈'
  if (i === 2) return '🥉'
  return `#${i + 1}`
}

export function LeaderboardPage() {
  const [boards, setBoards] = useState<Leaderboards | null>(null)
  const [tab, setTab] = useState('total')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    quizApi.getLeaderboards()
      .then(setBoards)
      .finally(() => setLoading(false))
  }, [])

  const entries = boards
    ? tab === 'total'
      ? boards.Total ?? []
      : boards.ByRarity?.[tab] ?? []
    : []

  return (
    <div className="page-wide">
      <h2 style={{ fontSize: '1.4rem', fontWeight: 800, marginBottom: 6 }}>Leaderboard</h2>
      <p className="text-sm text-muted" style={{ marginBottom: 24 }}>Top collectors across all categories</p>

      <div className="lb-tabs">
        {['total', ...RARITIES].map((key) => (
          <button
            key={key}
            className={`lb-tab${tab === key ? ' active' : ''}`}
            onClick={() => setTab(key)}
          >
            {key !== 'total' && <span className={`badge badge-${key}`} style={{ marginRight: 4 }}>{key}</span>}
            {TAB_LABELS[key]}
          </button>
        ))}
      </div>

      <div className="card" style={{ padding: 0, overflow: 'hidden' }}>
        {loading ? (
          <div style={{ padding: 32, textAlign: 'center', color: 'var(--text-muted)' }}>Loading...</div>
        ) : entries.length === 0 ? (
          <div style={{ padding: 32, textAlign: 'center', color: 'var(--text-muted)' }}>
            No data yet. Be the first!
          </div>
        ) : (
          entries.map((e, i) => (
            <div key={e.Username} className="lb-row">
              <span className={`lb-rank ${rankClass(i)}`}>{rankLabel(i)}</span>
              <Link to={`/users/${e.Username}`} className="lb-username" style={{ color: 'inherit' }}>
                {e.Username}
              </Link>
              <span className="lb-count">{e.Count} {e.Count === 1 ? 'article' : 'articles'}</span>
            </div>
          ))
        )}
      </div>
    </div>
  )
}
