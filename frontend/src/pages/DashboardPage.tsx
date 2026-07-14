import { useAuth } from '../auth/AuthContext'

export function DashboardPage() {
  const { user, logout } = useAuth()

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
          <p className="soon">Coming in Phase 2</p>
        </article>
        <article>
          <h3>Land the job</h3>
          <p>Search openings and build a resume customized for each role.</p>
          <p className="soon">Coming in Phase 3</p>
        </article>
        <article>
          <h3>Pivot your career</h3>
          <p>Switch fields, reduce hours, or move to consulting with a guided plan.</p>
          <p className="soon">Coming in Phase 4</p>
        </article>
      </section>
    </main>
  )
}