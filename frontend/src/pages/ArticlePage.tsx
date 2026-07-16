import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { articlesApi, type Article } from '../api/articles'
import { RarityBadge } from '../components/RarityBadge'
import { OwnershipBanner } from '../components/OwnershipBanner'

function ArticleContent({ content }: { content: string }) {
  const paragraphs = content.split(/\n\n+/).filter(Boolean)
  return (
    <div className="article-content">
      {paragraphs.map((p, i) => (
        <p key={i} className="article-paragraph">{p.trim()}</p>
      ))}
    </div>
  )
}

export function ArticlePage() {
  const { slug } = useParams<{ slug: string }>()
  const [article, setArticle] = useState<Article | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const navigate = useNavigate()

  useEffect(() => {
    if (!slug) return
    setLoading(true)
    articlesApi.getBySlug(slug)
      .then(setArticle)
      .catch(() => setError('Article not found'))
      .finally(() => setLoading(false))
  }, [slug])

  if (loading) return (
    <div style={{ textAlign: 'center', marginTop: 80, color: 'var(--text-muted)' }}>Loading article...</div>
  )
  if (error || !article) return (
    <div className="page-narrow" style={{ textAlign: 'center', marginTop: 60 }}>
      <div className="alert alert-danger">{error || 'Article not found'}</div>
      <button className="btn btn-ghost" style={{ marginTop: 16 }} onClick={() => navigate('/')}>Go home</button>
    </div>
  )

  const wikiUrl = `https://en.wikipedia.org/wiki/${article.slug}`

  return (
    <div className="page">
      <button
        onClick={() => navigate('/')}
        className="btn btn-ghost"
        style={{ marginBottom: 24, padding: '6px 14px', fontSize: '.875rem' }}
      >
        ← Discover another
      </button>

      <div className="card">
        <div className="flex-center gap-12" style={{ marginBottom: 16, flexWrap: 'wrap' }}>
          <h1 style={{ fontSize: 'clamp(1.3rem, 4vw, 2rem)', fontWeight: 800, flex: 1, minWidth: 0 }}>{article.title}</h1>
          <RarityBadge tier={article.rarity_tier} />
        </div>

        <OwnershipBanner owner={article.owner} />

        <hr className="divider" />

        <ArticleContent content={article.content || article.summary} />

        <hr className="divider" />

        <div className="flex-center gap-12" style={{ flexWrap: 'wrap' }}>
          <a href={wikiUrl} target="_blank" rel="noopener noreferrer" className="btn btn-ghost">
            Open on Wikipedia ↗
          </a>
          <button
            className="btn btn-primary"
            onClick={() => navigate(`/articles/${slug}/quiz`)}
          >
            Take Quiz →
          </button>
        </div>

        <p className="text-xs text-muted" style={{ marginTop: 16 }}>
          {article.content_length.toLocaleString()} bytes · {article.rarity_tier} article
        </p>
      </div>
    </div>
  )
}
