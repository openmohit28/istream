import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { createThread, deleteThread, listThreads } from '../api/pivot'
import type { PivotThread } from '../api/pivot'

export function PivotListPage() {
  const navigate = useNavigate()
  const [threads, setThreads] = useState<PivotThread[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    listThreads()
      .then((res) => setThreads(res?.threads ?? []))
      .catch(() => setThreads([]))
      .finally(() => setLoading(false))
  }, [])

  async function handleNew() {
    setError(null)
    try {
      const thread = await createThread()
      navigate(`/pivot/${thread.id}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not start exploration')
    }
  }

  async function handleDelete(id: string) {
    if (!window.confirm('Delete this exploration? This cannot be undone.')) return
    await deleteThread(id)
    setThreads((ts) => ts.filter((t) => t.id !== id))
  }

  return (
    <main className="pivot-page">
      <header>
        <h1>Pivot your career</h1>
        <Link to="/">Back to dashboard</Link>
      </header>
      <p className="hint">
        Answer a short series of questions to land on a concrete pivot plan - reduce hours, move
        within your field, switch out, or go consulting. Explore several directions side by side:
        fork any exploration at an earlier answer.
      </p>
      <p>
        <button type="button" className="cta-btn" onClick={handleNew}>
          Start a new exploration
        </button>
      </p>
      {error && <p role="alert" className="error">{error}</p>}

      {loading ? (
        <p className="status">Loading...</p>
      ) : threads.length === 0 ? (
        <p>No explorations yet. Start one - it takes about a minute.</p>
      ) : (
        <ul className="thread-list">
          {threads.map((t) => (
            <li key={t.id}>
              <Link to={`/pivot/${t.id}`}>
                <strong>{t.outcome ? t.outcome.title : 'In progress'}</strong>
              </Link>
              <span className="meta">
                {t.steps.length > 0 ? `"${t.steps[0].option}"` : 'not started'}
                {t.forkedFrom && ' · forked'}
                {' · '}
                {new Date(t.updatedAt).toLocaleDateString()}
              </span>
              <button type="button" className="remove" onClick={() => handleDelete(t.id)}>
                Delete
              </button>
            </li>
          ))}
        </ul>
      )}
    </main>
  )
}