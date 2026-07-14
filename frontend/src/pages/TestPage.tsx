import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { getQuestions, submitQuiz } from '../api/quiz'
import type { QuizQuestion, QuizScale } from '../api/quiz'

const SCALE_LABELS = ['Strongly disagree', 'Disagree', 'Neutral', 'Agree', 'Strongly agree']

export function TestPage() {
  const navigate = useNavigate()
  const [questions, setQuestions] = useState<QuizQuestion[]>([])
  const [scale, setScale] = useState<QuizScale | null>(null)
  const [index, setIndex] = useState(0)
  const [answers, setAnswers] = useState<Record<string, number>>({})
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    getQuestions()
      .then((res) => {
        setQuestions(res.questions)
        setScale(res.scale)
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Could not load questions'))
  }, [])

  if (error && questions.length === 0) {
    return (
      <main className="test-page">
        <p role="alert" className="error">{error}</p>
        <Link to="/">Back to dashboard</Link>
      </main>
    )
  }
  if (questions.length === 0 || !scale) {
    return <p className="status">Loading questions...</p>
  }

  const question = questions[index]
  const answered = Object.keys(answers).length
  const isLast = index === questions.length - 1
  const allAnswered = answered === questions.length

  function select(value: number) {
    setAnswers((prev) => ({ ...prev, [question.id]: value }))
    if (!isLast) {
      setIndex((i) => i + 1)
    }
  }

  async function handleSubmit() {
    setError(null)
    setSubmitting(true)
    try {
      const result = await submitQuiz(answers)
      navigate(`/results/${result.id}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not submit answers')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <main className="test-page">
      <header>
        <h1>Personality test</h1>
        <Link to="/">Exit</Link>
      </header>

      <div
        className="progress"
        role="progressbar"
        aria-valuenow={answered}
        aria-valuemin={0}
        aria-valuemax={questions.length}
        aria-label="questions answered"
      >
        <div className="progress-fill" style={{ width: `${(answered / questions.length) * 100}%` }} />
      </div>
      <p className="progress-text">
        Question {index + 1} of {questions.length}
      </p>

      <h2>{question.text}</h2>

      <div className="likert" role="group" aria-label="answer options">
        {SCALE_LABELS.map((label, i) => {
          const value = scale.min + i
          const selected = answers[question.id] === value
          return (
            <button
              key={value}
              type="button"
              className={selected ? 'likert-option selected' : 'likert-option'}
              aria-pressed={selected}
              onClick={() => select(value)}
            >
              {label}
            </button>
          )
        })}
      </div>

      {error && <p role="alert" className="error">{error}</p>}

      <footer className="test-nav">
        <button type="button" onClick={() => setIndex((i) => i - 1)} disabled={index === 0}>
          Back
        </button>
        {isLast && (
          <button type="button" className="primary" onClick={handleSubmit} disabled={!allAnswered || submitting}>
            {submitting ? 'Scoring...' : 'See my results'}
          </button>
        )}
      </footer>
    </main>
  )
}