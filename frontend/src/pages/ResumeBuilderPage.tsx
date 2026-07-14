import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { createResume, emptyDocument, getResume, updateResume } from '../api/resumes'
import type { Education, Experience, ResumeDocument } from '../api/resumes'

const STEPS = ['Target role', 'Contact', 'Summary', 'Experience', 'Education', 'Skills', 'Review'] as const

// Wizard-local shapes keep list fields as editable text; they are converted
// to arrays on save.
interface ExperienceDraft extends Omit<Experience, 'bullets'> {
  bulletsText: string
}

function emptyExperience(): ExperienceDraft {
  return { company: '', title: '', location: '', startDate: '', endDate: '', current: false, bulletsText: '' }
}

function emptyEducation(): Education {
  return { school: '', degree: '', field: '', gradYear: '' }
}

const splitLines = (text: string) => text.split('\n').map((s) => s.trim()).filter(Boolean)
const splitCommas = (text: string) => text.split(',').map((s) => s.trim()).filter(Boolean)

export function ResumeBuilderPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [step, setStep] = useState(0)
  const [doc, setDoc] = useState<ResumeDocument>(emptyDocument())
  const [experience, setExperience] = useState<ExperienceDraft[]>([])
  const [education, setEducation] = useState<Education[]>([])
  const [skillsText, setSkillsText] = useState('')
  const [certsText, setCertsText] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)
  const [loading, setLoading] = useState(Boolean(id))

  useEffect(() => {
    if (!id) return
    getResume(id)
      .then((r) => {
        setDoc(r.data)
        setExperience(r.data.experience.map((e) => ({ ...e, bulletsText: e.bullets.join('\n') })))
        setEducation(r.data.education)
        setSkillsText(r.data.skills.join(', '))
        setCertsText(r.data.certifications.join(', '))
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Could not load resume'))
      .finally(() => setLoading(false))
  }, [id])

  if (loading) return <p className="status">Loading resume...</p>

  const patch = (partial: Partial<ResumeDocument>) => setDoc((d) => ({ ...d, ...partial }))
  const patchContact = (field: string, value: string) =>
    setDoc((d) => ({ ...d, contact: { ...d.contact, [field]: value } }))

  const canAdvance = () => {
    if (step === 0) return doc.targetTitle.trim() !== ''
    if (step === 1) return doc.contact.fullName.trim() !== '' && doc.contact.email.trim() !== ''
    return true
  }

  async function handleSave() {
    setError(null)
    setSaving(true)
    const payload: ResumeDocument = {
      ...doc,
      experience: experience.map(({ bulletsText, ...e }) => ({ ...e, bullets: splitLines(bulletsText) })),
      education,
      skills: splitCommas(skillsText),
      certifications: splitCommas(certsText),
    }
    try {
      const saved = id ? await updateResume(id, payload) : await createResume(payload)
      navigate(`/resumes/${saved.id}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not save resume')
    } finally {
      setSaving(false)
    }
  }

  return (
    <main className="builder-page">
      <header>
        <h1>{id ? 'Edit resume' : 'Build your resume'}</h1>
        <Link to="/resumes">Exit</Link>
      </header>

      <ol className="steps" aria-label="wizard steps">
        {STEPS.map((label, i) => (
          <li key={label} className={i === step ? 'active' : i < step ? 'done' : ''} aria-current={i === step ? 'step' : undefined}>
            {label}
          </li>
        ))}
      </ol>

      {step === 0 && (
        <section className="step-body">
          <h2>What job are you targeting?</h2>
          <label>
            Target job title
            <input value={doc.targetTitle} onChange={(e) => patch({ targetTitle: e.target.value })} required />
          </label>
          <label>
            Paste the job description (optional - unlocks the ATS keyword check)
            <textarea
              rows={6}
              value={doc.jobDescription}
              onChange={(e) => patch({ jobDescription: e.target.value })}
            />
          </label>
        </section>
      )}

      {step === 1 && (
        <section className="step-body">
          <h2>How can employers reach you?</h2>
          <label>
            Full name
            <input value={doc.contact.fullName} onChange={(e) => patchContact('fullName', e.target.value)} required />
          </label>
          <label>
            Email
            <input type="email" value={doc.contact.email} onChange={(e) => patchContact('email', e.target.value)} required />
          </label>
          <label>
            Phone
            <input value={doc.contact.phone} onChange={(e) => patchContact('phone', e.target.value)} />
          </label>
          <label>
            Location
            <input value={doc.contact.location} onChange={(e) => patchContact('location', e.target.value)} />
          </label>
          <label>
            LinkedIn URL
            <input value={doc.contact.linkedin} onChange={(e) => patchContact('linkedin', e.target.value)} />
          </label>
          <label>
            Website / portfolio
            <input value={doc.contact.website} onChange={(e) => patchContact('website', e.target.value)} />
          </label>
        </section>
      )}

      {step === 2 && (
        <section className="step-body">
          <h2>Summarize yourself in 2-3 sentences</h2>
          <p className="hint">Lead with your strongest match for the target role. Recruiters scan for 20-30 seconds.</p>
          <label>
            Professional summary
            <textarea rows={5} value={doc.summary} onChange={(e) => patch({ summary: e.target.value })} />
          </label>
        </section>
      )}

      {step === 3 && (
        <section className="step-body">
          <h2>Where have you worked?</h2>
          {experience.map((exp, i) => (
            <fieldset key={i}>
              <legend>Position {i + 1}</legend>
              <label>
                Job title
                <input
                  value={exp.title}
                  onChange={(e) => setExperience(experience.map((x, j) => (j === i ? { ...x, title: e.target.value } : x)))}
                />
              </label>
              <label>
                Company
                <input
                  value={exp.company}
                  onChange={(e) => setExperience(experience.map((x, j) => (j === i ? { ...x, company: e.target.value } : x)))}
                />
              </label>
              <div className="row">
                <label>
                  Start (e.g. Jan 2023)
                  <input
                    value={exp.startDate}
                    onChange={(e) => setExperience(experience.map((x, j) => (j === i ? { ...x, startDate: e.target.value } : x)))}
                  />
                </label>
                <label>
                  End
                  <input
                    value={exp.endDate}
                    disabled={exp.current}
                    onChange={(e) => setExperience(experience.map((x, j) => (j === i ? { ...x, endDate: e.target.value } : x)))}
                  />
                </label>
              </div>
              <label className="checkbox">
                <input
                  type="checkbox"
                  checked={exp.current}
                  onChange={(e) =>
                    setExperience(experience.map((x, j) => (j === i ? { ...x, current: e.target.checked, endDate: '' } : x)))
                  }
                />
                I currently work here
              </label>
              <label>
                Achievements (one per line - start with a verb, include numbers)
                <textarea
                  rows={4}
                  value={exp.bulletsText}
                  onChange={(e) => setExperience(experience.map((x, j) => (j === i ? { ...x, bulletsText: e.target.value } : x)))}
                />
              </label>
              <button type="button" className="remove" onClick={() => setExperience(experience.filter((_, j) => j !== i))}>
                Remove position
              </button>
            </fieldset>
          ))}
          <button type="button" onClick={() => setExperience([...experience, emptyExperience()])}>
            Add experience
          </button>
        </section>
      )}

      {step === 4 && (
        <section className="step-body">
          <h2>Education</h2>
          {education.map((edu, i) => (
            <fieldset key={i}>
              <legend>Education {i + 1}</legend>
              <label>
                School
                <input
                  value={edu.school}
                  onChange={(e) => setEducation(education.map((x, j) => (j === i ? { ...x, school: e.target.value } : x)))}
                />
              </label>
              <div className="row">
                <label>
                  Degree
                  <input
                    value={edu.degree}
                    onChange={(e) => setEducation(education.map((x, j) => (j === i ? { ...x, degree: e.target.value } : x)))}
                  />
                </label>
                <label>
                  Field of study
                  <input
                    value={edu.field}
                    onChange={(e) => setEducation(education.map((x, j) => (j === i ? { ...x, field: e.target.value } : x)))}
                  />
                </label>
              </div>
              <label>
                Graduation year
                <input
                  value={edu.gradYear}
                  onChange={(e) => setEducation(education.map((x, j) => (j === i ? { ...x, gradYear: e.target.value } : x)))}
                />
              </label>
              <button type="button" className="remove" onClick={() => setEducation(education.filter((_, j) => j !== i))}>
                Remove
              </button>
            </fieldset>
          ))}
          <button type="button" onClick={() => setEducation([...education, emptyEducation()])}>
            Add education
          </button>
        </section>
      )}

      {step === 5 && (
        <section className="step-body">
          <h2>Skills and certifications</h2>
          <label>
            Skills (comma-separated - mirror the job description wording)
            <textarea rows={3} value={skillsText} onChange={(e) => setSkillsText(e.target.value)} />
          </label>
          <label>
            Certifications (comma-separated, optional)
            <textarea rows={2} value={certsText} onChange={(e) => setCertsText(e.target.value)} />
          </label>
        </section>
      )}

      {step === 6 && (
        <section className="step-body">
          <h2>Ready to save</h2>
          <ul className="review">
            <li><strong>Target:</strong> {doc.targetTitle}</li>
            <li><strong>Name:</strong> {doc.contact.fullName} ({doc.contact.email})</li>
            <li><strong>Experience:</strong> {experience.length} position(s)</li>
            <li><strong>Education:</strong> {education.length} entr(ies)</li>
            <li><strong>Skills:</strong> {splitCommas(skillsText).length}</li>
            <li>
              <strong>ATS keyword check:</strong>{' '}
              {doc.jobDescription.trim() ? 'will run on the preview page' : 'add a job description in step 1 to enable'}
            </li>
          </ul>
        </section>
      )}

      {error && <p role="alert" className="error">{error}</p>}

      <footer className="test-nav">
        <button type="button" onClick={() => setStep((s) => s - 1)} disabled={step === 0}>
          Back
        </button>
        {step < STEPS.length - 1 ? (
          <button type="button" className="primary" onClick={() => setStep((s) => s + 1)} disabled={!canAdvance()}>
            Next
          </button>
        ) : (
          <button type="button" className="primary" onClick={handleSave} disabled={saving}>
            {saving ? 'Saving...' : 'Save resume'}
          </button>
        )}
      </footer>
    </main>
  )
}