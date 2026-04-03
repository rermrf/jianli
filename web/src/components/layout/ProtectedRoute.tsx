import type { PropsWithChildren } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { isAuthenticated, setRedirectPath } from '../../lib/auth'

export function ProtectedRoute({ children }: PropsWithChildren) {
  const location = useLocation()

  if (!isAuthenticated()) {
    setRedirectPath(location.pathname)

    return <Navigate replace to="/login" />
  }

  return children
}
