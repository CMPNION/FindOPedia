import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { quizApi, type QuizQuestion } from '../api/quiz'

export function QuizPage() {
  const { slug } = useParams<{ slug: string }>()
  const navigate = useNavigate()
  const [questions, setQuestions] = useState<QuizQuestion[]>([])
  const [answers, setAnswers] = useState<Record<number, string>>({})
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [cooldownUntil, setCooldownUntil] = useState<Date | null>(null)

  useEffect(() => {
    if (!slug) return

    const provider = localStorage.getItem('ai_provider') ?? undefined
    const key = localStorage.getItem('ai_api_key') ?? undefined

    quizApi
      .getOrGenerateQuestions(slug, provider, key)
      .then((res) => {
        console.log('Question IDs:', res.questions.map((q) => q.id))
        setQuestions(res.questions)
      })
      .catch(async (err: unknown) => {
        const msg = err instanceof Error ? err.message : 'Failed'

        if (msg === 'already_attempted') {
          setError('You already attempted this article. One chance per article.')
        } else if (msg === 'cooldown_active') {
          try {
            const cd = await quizApi.getCooldown()
            if (cd.active) {
              setCooldownUntil(new Date(cd.next_claim))
            }
          } catch {
            setError('Cooldown active.')
          }
        } else if (!key) {
          setError('No AI key set. Go to Settings first.')
        } else {
          setError(msg)
        }
      })
      .finally(() => setLoading(false))
  }, [slug])

  const answered = Object.keys(answers).length
  const total = questions.length

  async function submit() {
    if (!slug) return

    if (answered < total) {
      setError(`Answer all ${total} questions first.`)
      return
    }

    setError('')
    setSubmitting(true)

    try {
      const result = await quizApi.submitAttempt(
        slug,
        questions.map((q, i) => ({
          question_id: q.id,
          chosen_answer: answers[i] ?? '',
        }))
      )

      navigate(`/result/${slug}`, { state: result })
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Submission failed')
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <div
        style={{
          textAlign: 'center',
          marginTop: 80,
          color: 'var(--text-muted)',
        }}
      >
        Generating questions...
      </div>
    )
  }

  if (cooldownUntil) {
    return (
      <div className="page-narrow" style={{ marginTop: 40 }}>
        <div className="cooldown-banner">
          <strong>Cooldown active</strong> — you can claim one article every 4 hours.
          <br />
          <span
            className="text-sm"
            style={{
              marginTop: 6,
              display: 'block',
            }}
          >
            Next claim available:{' '}
            <strong>{cooldownUntil.toLocaleTimeString()}</strong>
          </span>
        </div>

        <button
          className="btn btn-ghost"
          style={{ marginTop: 16 }}
          onClick={() => navigate(-1)}
        >
          ← Go back
        </button>
      </div>
    )
  }

  if (error && total === 0) {
    return (
      <div className="page-narrow" style={{ marginTop: 40 }}>
        <div className="alert alert-danger" style={{ marginBottom: 16 }}>
          {error}
        </div>

        <button
          className="btn btn-ghost"
          onClick={() => navigate(-1)}
        >
          ← Go back
        </button>
      </div>
    )
  }

  return (
    <div className="page">
      <div style={{ marginBottom: 24 }}>
        <h2 style={{ fontSize: '1.5rem', fontWeight: 800 }}>
          Quiz
        </h2>

        <p
          className="text-sm text-muted"
          style={{ marginTop: 4 }}
        >
          {answered}/{total} answered · 100% required to own this article
        </p>

        <div
          style={{
            marginTop: 8,
            height: 4,
            background: 'var(--border)',
            borderRadius: 999,
            overflow: 'hidden',
          }}
        >
          <div
            style={{
              height: '100%',
              background: 'var(--primary)',
              borderRadius: 999,
              width: `${total > 0 ? (answered / total) * 100 : 0}%`,
              transition: 'width .2s',
            }}
          />
        </div>
      </div>

      {questions.map((q, i) => (
        <div
          key={`${q.id}-${i}`}
          className="quiz-question"
        >
          <p className="quiz-question-text">
            <span
              style={{
                color: 'var(--text-muted)',
                marginRight: 6,
              }}
            >
              {i + 1}.
            </span>

            {q.question_text}
          </p>

          <div className="quiz-options">
            {q.options.map((opt) => (
              <label
                key={`${i}-${opt.value}`}
                className={`quiz-option${
                  answers[i] === opt.value ? ' selected' : ''
                }`}
              >
                <input
                  type="radio"
                  name={`question-${i}`}
                  value={opt.value}
                  checked={answers[i] === opt.value}
                  onChange={() =>
                    setAnswers((prev) => ({
                      ...prev,
                      [i]: opt.value,
                    }))
                  }
                />

                {q.type === 'multiple_choice'
                  ? `${opt.label}. ${opt.value}`
                  : opt.value}
              </label>
            ))}
          </div>
        </div>
      ))}

      {error && (
        <div
          className="alert alert-danger"
          style={{ marginBottom: 16 }}
        >
          {error}
        </div>
      )}

      <button
        className="btn btn-primary btn-full btn-lg"
        onClick={submit}
        disabled={submitting || answered < total}
        style={{ marginTop: 8 }}
      >
        {submitting
          ? 'Submitting...'
          : answered < total
            ? `Answer ${total - answered} more`
            : 'Submit Answers'}
      </button>
    </div>
  )
}