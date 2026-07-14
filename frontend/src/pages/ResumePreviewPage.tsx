import { useCallback, useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getResume, keywordCheck } from '../api/resumes'
import type { KeywordReport, Resume } from '../api/resumes'

export function ResumePreviewPage() {
  const { id } = useParams<{ id: string }>()
  const [resume, setResume] = useState<Resume | null>(null)
  const [report, setReport] = useState<KeywordReport | null>(null)
  const [jd, setJd] = useState('')
  const [error, setError] = useState<string | null>(null)

  const runCheck = useCallback(
    (jobDescription: string) => {
      if (!id || !jobDescription.trim()) return
      keywordCheck(id, jobDescription)
        .then(setReport)
        .catch(() => setReport(null))
    },
    [id],
  )

  useEffect(() => {
    if (!id) return
    getResume(id)
      .then((r) => {
        setResume(r)
        setJd(r.data.jobDescription)
        if (r.data.jobDescription.trim()) runCheck(r.data.jobDescription)
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Could not load resume'))
  }, [id, runCheck])

  if (error) {
    return (
      <main className="preview-page">
        <p role="alert" className="error">{error}</p>
        <Link to="/resumes">Back to resumes</Link>
      </main>
    )
  }
  if (!resume) return <p className="status">Loading resume...</p>

  const d = resume.data
  const contactLine = [d.contact.location, d.contact.email, d.contact.phone, d.contact.linkedin, d.contact.website]
    .filter(Boolean)
    .join(' | ')

  return (
    <main className="preview-page">
      <header className="no-print">
        <h1>Resume preview</h1>
        <nav>
          <Link to={`/resumes/${resume.id}/edit`}>Edit</Link>
          <Link to="/resumes">All resumes</Link>
          <button type="button" onClick={() => window.print()}>
            Print / Save as PDF
          </button>
        </nav>
      </header>
      <p className="hint no-print">
        Single-column, standard headings, no graphics - the format ATS parsers read reliably.
      </p>

      <article className="resume-sheet" aria-label="resume">
        <h2 className="resume-name">{d.contact.fullName}</h2>
        <p className="resume-target">{d.targetTitle}</p>
        {contactLine && <p className="resume-contact">{contactLine}</p>}

        {d.summary && (
          <>
            <h3>Summary</h3>
            <p>{d.summary}</p>
          </>
        )}

        {d.experience.length > 0 && (
          <>
            <h3>Experience</h3>
            {d.experience.map((e, i) => (
              <div key={i} className="resume-entry">
                <p className="entry-head">
                  <strong>{e.title}</strong>
                  {e.company && <> - {e.company}</>}
                  {e.location && <>, {e.location}</>}
                </p>
                {(e.startDate || e.endDate || e.current) && (
                  <p className="entry-dates">
                    {e.startDate} - {e.current ? 'Present' : e.endDate}
                  </p>
                )}
                {e.bullets.length > 0 && (
                  <ul>
                    {e.bullets.map((b, j) => (
                      <li key={j}>{b}</li>
                    ))}
                  </ul>
                )}
              </div>
            ))}
          </>
        )}

        {d.education.length > 0 && (
          <>
            <h3>Education</h3>
            {d.education.map((e, i) => (
              <p key={i}>
                <strong>{e.degree}</strong>
                {e.field && <> in {e.field}</>}
                {e.school && <> - {e.school}</>}
                {e.gradYear && <> ({e.gradYear})</>}
              </p>
            ))}
          </>
        )}

        {d.skills.length > 0 && (
          <>
            <h3>Skills</h3>
            <p>{d.skills.join(', ')}</p>
          </>
        )}

        {d.certifications.length > 0 && (
          <>
            <h3>Certifications</h3>
            <p>{d.certifications.join(', ')}</p>
          </>
        )}
      </article>

      <section className="keyword-panel no-print" aria-label="ATS keyword check">
        <h3>ATS keyword check</h3>
        <p className="hint">
          64% of ATS setups auto-reject poor keyword matches. Paste the job description to see your coverage.
        </p>
        <textarea
          rows={4}
          value={jd}
          onChange={(e) => setJd(e.target.value)}
          placeholder="Paste the job description here"
        />
        <button type="button" onClick={() => runCheck(jd)} disabled={!jd.trim()}>
          Check keywords
        </button>

        {report && (
          <div className="keyword-report">
            <p className="score">
              Keyword coverage: <strong>{report.score}%</strong>
            </p>
            {report.matched.length > 0 && (
              <p>
                <span className="report-label">Covered:</span>{' '}
                {report.matched.map((k) => (
                  <span key={k} className="badge matched">{k}</span>
                ))}
              </p>
            )}
            {report.missing.length > 0 && (
              <p>
                <span className="report-label">Missing - work these in:</span>{' '}
                {report.missing.map((k) => (
                  <span key={k} className="badge missing">{k}</span>
                ))}
              </p>
            )}
          </div>
        )}
      </section>
    </main>
  )
}