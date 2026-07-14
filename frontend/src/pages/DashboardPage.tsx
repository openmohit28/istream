import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'
import { getResults } from '../api/quiz'
import type { QuizResultSummary } from '../api/quiz'

export function DashboardPage() {
  const { user, logout } = useAuth()
  const [results, setResults] = useState<QuizResultSummary[]>([])

  useEffect(() => {
    getResults()
      .then((res) => setResults(res?.results ?? []))
      .catch(() => setResults([]))
  }, [])

  return (
    <main className="dashboard">
      <header>
        <h1>istream</h1>
        <button type="button" onClick={logout}>
          Log out
        </button>
      </header>
      <h2>Welcome, {user?.name}</h2>
      <p>Pick where you are in your career journey:</p>
      <section className="options">
        <article>
          <h3>Discover your fit</h3>
          <p>Take the personality test to find the jobs that suit you best.</p>
          <Link className="cta" to="/test">
            Take the test
          </Link>
        </article>
        <article>
          <h3>Land the job</h3>
          <p>Search openings and build a resume customized for each role.</p>
          <p className="card-links">
            <Link className="cta" to="/jobs">
              Search jobs
            </Link>{' '}
            <Link className="cta secondary" to="/resumes">
              My resumes
            </Link>
          </p>
        </article>
        <article>
          <h3>Pivot your career</h3>
          <p>Switch fields, reduce hours, or move to consulting with a guided plan.</p>
          <p className="soon">Coming in Phase 4</p>
        </article>
      </section>

      {results.length > 0 && (
        <section className="history" aria-label="past results">
          <h3>Your past results</h3>
          <ul>
            {results.map((r) => (
              <li key={r.id}>
                <Link to={`/results/${r.id}`}>
                  {new Date(r.createdAt).toLocaleDateString()} - top match:{' '}
                  {r.topMatch ? `${r.topMatch.title} (${r.topMatch.fit}% fit)` : 'view details'}
                </Link>
              </li>
            ))}
          </ul>
        </section>
      )}
    </main>
  )
}