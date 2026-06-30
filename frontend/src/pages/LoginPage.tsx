import { useAuth } from '../hooks/useAuth'
import { login } from '../services/api'
import { useNavigate } from 'react-router-dom'
import { useEffect } from 'react'

export default function LoginPage() {
  const { user } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (user) navigate('/')
  }, [user, navigate])

  const handleLogin = async () => {
    try {
      const response = await login('/')
      if (response.data.redirect) {
        window.location.href = response.data.redirect
      }
    } catch {
      // Login failed
    }
  }

  return (
    <div style={{ textAlign: 'center', marginTop: '50px' }}>
      <h1>Login</h1>
      <p>Please log in to access your account.</p>
      <button onClick={handleLogin} style={{ padding: '10px 20px', fontSize: '16px' }}>Log In</button>
    </div>
  )
}
