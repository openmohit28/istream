import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getResult, RIASEC_LABELS } from '../api/quiz'
import type { QuizResult, RiasecKey } from '../api/quiz'

const DEMAND_LABELS = {
  growing: 'Growing demand',
  stable: 'Stable demand',
  declining: 'Declining demand',
} as const

const AI_RISK_LABELS = {
  low: 'Low AI exposure',
  medium: 'Medium AI exposure',
  high: 'High AI exposure',
} as const

export function ResultPage() {
  const { id } = useParams<{ id: string }>()
  const [result, setResult] = useState<QuizResult | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!id) return
    getResult(id)
      .then(setResult)
      .catch((err) => setError(err instanceof Error ? err.message : 'Could not load result'))
  }, [id])

  if (error) {
    return (
      <main className="result-page">
        <p role="alert" className="error">{error}</p>
        <Link to="/">Back to dashboard</Link>
      </main>
    )
  }
  if (!result) {
    return <p className="status">Loading your results...</p>
  }

  const dimensions = Object.keys(RIASEC_LABELS) as RiasecKey[]

  return (
    <main className="result-page">
      <header>
        <h1>Your results</h1>
        <Link to="/">Back to dashboard</Link>
      </header>

      <section aria-label="interest profile">
        <h2>Your interest profile</h2>
        <ul className="profile-bars">
          {dimensions.map((d) => (
            <li key={d}>
              <span className="bar-label">{RIASEC_LABELS[d]}</span>
              <div className="bar-track">
                <div className="bar-fill" style={{ width: `${result.scores[d]}%` }} />
              </div>
              <span className="bar-value">{result.scores[d]}</span>
            </li>
          ))}
        </ul>
      </section>

      <section aria-label="job matches">
        <h2>Jobs that fit you</h2>
        <p className="hint">
          Ranked by fit with your profile, with a boost for fields growing through 2030.
        </p>
        <ol className="matches">
          {result.matches.map((m) => (
            <li key={m.title} className="match-card">
              <div className="match-head">
                <h3>{m.title}</h3>
                <span className="fit">{m.fit}% fit</span>
              </div>
              <p className="badges">
                <span className={`badge demand-${m.demand}`}>{DEMAND_LABELS[m.demand]}</span>
                <span className={`badge risk-${m.aiRisk}`}>{AI_RISK_LABELS[m.aiRisk]}</span>
                <span className="badge category">{m.category}</span>
              </p>
              <p>{m.blurb}</p>
              <Link className="search-link" to="/jobs" state={{ keywords: m.title }}>
                Search openings for this role
              </Link>
            </li>
          ))}
        </ol>
      </section>

      <p>
        <Link to="/test">Retake the test</Link>
      </p>
    </main>
  )
}