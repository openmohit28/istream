import { Navigate, Route, Routes } from 'react-router-dom'
import { ProtectedRoute } from './components/ProtectedRoute'
import { DashboardPage } from './pages/DashboardPage'
import { LoginPage } from './pages/LoginPage'
import { RegisterPage } from './pages/RegisterPage'
import { JobSearchPage } from './pages/JobSearchPage'
import { ResultPage } from './pages/ResultPage'
import { ResumeBuilderPage } from './pages/ResumeBuilderPage'
import { ResumeListPage } from './pages/ResumeListPage'
import { ResumePreviewPage } from './pages/ResumePreviewPage'
import { TestPage } from './pages/TestPage'

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <DashboardPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/test"
        element={
          <ProtectedRoute>
            <TestPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/results/:id"
        element={
          <ProtectedRoute>
            <ResultPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/jobs"
        element={
          <ProtectedRoute>
            <JobSearchPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/resumes"
        element={
          <ProtectedRoute>
            <ResumeListPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/resumes/new"
        element={
          <ProtectedRoute>
            <ResumeBuilderPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/resumes/:id/edit"
        element={
          <ProtectedRoute>
            <ResumeBuilderPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/resumes/:id"
        element={
          <ProtectedRoute>
            <ResumePreviewPage />
          </ProtectedRoute>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App