import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { deleteResume, listResumes } from '../api/resumes'
import type { ResumeSummary } from '../api/resumes'

export function ResumeListPage() {
  const [resumes, setResumes] = useState<ResumeSummary[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    listResumes()
      .then((res) => setResumes(res?.resumes ?? []))
      .catch(() => setResumes([]))
      .finally(() => setLoading(false))
  }, [])

  async function handleDelete(id: string) {
    if (!window.confirm('Delete this resume? This cannot be undone.')) return
    await deleteResume(id)
    setResumes((rs) => rs.filter((r) => r.id !== id))
  }

  return (
    <main className="resumes-page">
      <header>
        <h1>Your resumes</h1>
        <Link to="/">Back to dashboard</Link>
      </header>
      <p>
        <Link className="cta" to="/resumes/new">
          Build a new resume
        </Link>
      </p>

      {loading ? (
        <p className="status">Loading...</p>
      ) : resumes.length === 0 ? (
        <p>No resumes yet. Build one customized for the job you want.</p>
      ) : (
        <ul className="resume-list">
          {resumes.map((r) => (
            <li key={r.id}>
              <Link to={`/resumes/${r.id}`}>
                <strong>{r.title}</strong>
              </Link>
              <span className="meta">updated {new Date(r.updatedAt).toLocaleDateString()}</span>
              <span className="actions">
                <Link to={`/resumes/${r.id}/edit`}>Edit</Link>
                <button type="button" className="remove" onClick={() => handleDelete(r.id)}>
                  Delete
                </button>
              </span>
            </li>
          ))}
        </ul>
      )}
    </main>
  )
}