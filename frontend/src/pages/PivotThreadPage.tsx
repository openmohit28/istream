import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { answerThread, forkThread, getThread } from '../api/pivot'
import type { PivotThread } from '../api/pivot'

export function PivotThreadPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [thread, setThread] = useState<PivotThread | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!id) return
    setThread(null)
    getThread(id)
      .then(setThread)
      .catch((err) => setError(err instanceof Error ? err.message : 'Could not load exploration'))
  }, [id])

  if (error) {
    return (
      <main className="pivot-page">
        <p role="alert" className="error">{error}</p>
        <Link to="/pivot">Back to explorations</Link>
      </main>
    )
  }
  if (!thread) return <p className="status">Loading exploration...</p>

  async function answer(option: string) {
    if (!thread?.current) return
    setError(null)
    try {
      const updated = await answerThread(thread.id, [
        ...thread.steps,
        { nodeId: thread.current.id, option },
      ])
      setThread(updated)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not save answer')
    }
  }

  // Fork so the chosen step gets re-asked in a new thread; this one is kept.
  async function forkAt(stepIndex: number) {
    setError(null)
    try {
      const fork = await forkThread(thread!.id, stepIndex)
      navigate(`/pivot/${fork.id}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not fork thread')
    }
  }

  return (
    <main className="pivot-page">
      <header>
        <h1>Career pivot</h1>
        <Link to="/pivot">All explorations</Link>
      </header>
      {thread.forkedFrom && (
        <p className="fork-note">Forked from an earlier exploration - the original is untouched.</p>
      )}

      {thread.steps.length > 0 && (
        <section className="pivot-trail" aria-label="your answers so far">
          <h2>Your path</h2>
          <ol>
            {thread.steps.map((step, i) => (
              <li key={i}>
                <span>{step.option}</span>
                <button type="button" className="fork-btn" onClick={() => forkAt(i)}>
                  Fork from here
                </button>
              </li>
            ))}
          </ol>
        </section>
      )}

      {error && <p role="alert" className="error">{error}</p>}

      {thread.current && (
        <section className="pivot-question" aria-label="current question">
          <h2>{thread.current.question}</h2>
          <div className="pivot-options">
            {thread.current.options.map((opt) => (
              <button key={opt.label} type="button" className="likert-option" onClick={() => answer(opt.label)}>
                {opt.label}
              </button>
            ))}
          </div>
        </section>
      )}

      {thread.outcome && (
        <section className="pivot-outcome" aria-label="your pivot plan">
          <p className="badge outcome-path">{thread.outcome.path}</p>
          <h2>{thread.outcome.title}</h2>
          <p className="tagline">{thread.outcome.tagline}</p>
          <p className="why-now">{thread.outcome.whyNow}</p>

          <h3>Your action plan</h3>
          <ol className="plan">
            {thread.outcome.plan.map((step, i) => (
              <li key={i}>{step}</li>
            ))}
          </ol>

          <h3>Resources</h3>
          <ul className="resources">
            {thread.outcome.resources.map((r) =>
              r.url.startsWith('/') ? (
                <li key={r.url}>
                  <Link to={r.url}>{r.title}</Link>
                </li>
              ) : (
                <li key={r.url}>
                  <a href={r.url} target="_blank" rel="noopener noreferrer">
                    {r.title}
                  </a>
                </li>
              ),
            )}
          </ul>

          <p className="hint">
            Curious about a different path? Fork from any answer above - this plan stays saved.
          </p>
        </section>
      )}
    </main>
  )
}