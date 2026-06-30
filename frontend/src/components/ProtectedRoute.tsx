import type { ReactNode } from 'react'
import { Navigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

interface ProtectedRouteProps {
  children: ReactNode
  admin?: boolean
}

export default function ProtectedRoute({ children, admin }: ProtectedRouteProps) {
  const { user, loading } = useAuth()

  if (loading) return <div>Loading...</div>

  if (!user) return <Navigate to="/login" replace />

  if (admin && !user.isAdmin) return <Navigate to="/" replace />

  return <>{children}</>
}
