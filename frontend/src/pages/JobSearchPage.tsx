import { useState, type FormEvent } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { getSearchURL } from '../api/jobs'

export function JobSearchPage() {
  // A match card on the results page can hand us a role to search for.
  const prefill = (useLocation().state as { keywords?: string } | null)?.keywords ?? ''
  const [keywords, setKeywords] = useState(prefill)
  const [location, setLocation] = useState('')
  const [workplace, setWorkplace] = useState('')
  const [experience, setExperience] = useState('')
  const [jobType, setJobType] = useState('')
  const [postedWithin, setPostedWithin] = useState('')
  const [url, setUrl] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)
    setUrl(null)
    try {
      const res = await getSearchURL({ keywords, location, workplace, experience, jobType, postedWithin })
      setUrl(res.url)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not build search link')
    }
  }

  return (
    <main className="jobs-page">
      <header>
        <h1>Search jobs</h1>
        <Link to="/">Back to dashboard</Link>
      </header>
      <p className="hint">
        Set your filters and we build a pre-filtered LinkedIn search - one click, no scrolling
        through irrelevant listings.
      </p>

      <form onSubmit={handleSubmit} aria-label="job search filters">
        <label>
          Job title or keywords
          <input value={keywords} onChange={(e) => setKeywords(e.target.value)} required />
        </label>
        <label>
          Location
          <input
            value={location}
            onChange={(e) => setLocation(e.target.value)}
            placeholder="City, country, or leave empty"
          />
        </label>
        <label>
          Workplace
          <select value={workplace} onChange={(e) => setWorkplace(e.target.value)}>
            <option value="">Any</option>
            <option value="remote">Remote</option>
            <option value="hybrid">Hybrid</option>
            <option value="onsite">On-site</option>
          </select>
        </label>
        <label>
          Experience level
          <select value={experience} onChange={(e) => setExperience(e.target.value)}>
            <option value="">Any</option>
            <option value="internship">Internship</option>
            <option value="entry">Entry level</option>
            <option value="associate">Associate</option>
            <option value="mid-senior">Mid-senior</option>
            <option value="director">Director</option>
            <option value="executive">Executive</option>
          </select>
        </label>
        <label>
          Job type
          <select value={jobType} onChange={(e) => setJobType(e.target.value)}>
            <option value="">Any</option>
            <option value="fulltime">Full-time</option>
            <option value="parttime">Part-time</option>
            <option value="contract">Contract</option>
            <option value="temporary">Temporary</option>
            <option value="internship">Internship</option>
          </select>
        </label>
        <label>
          Posted within
          <select value={postedWithin} onChange={(e) => setPostedWithin(e.target.value)}>
            <option value="">Any time</option>
            <option value="day">Past 24 hours</option>
            <option value="week">Past week</option>
            <option value="month">Past month</option>
          </select>
        </label>
        <button type="submit">Build my search</button>
      </form>

      {error && <p role="alert" className="error">{error}</p>}
      {url && (
        <p className="search-result">
          <a href={url} target="_blank" rel="noopener noreferrer" className="cta">
            Open your LinkedIn search
          </a>
        </p>
      )}
    </main>
  )
}