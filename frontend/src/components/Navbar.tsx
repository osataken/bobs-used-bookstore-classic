import { Link } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import { login, logout } from '../services/api'
import toast from 'react-hot-toast'

export default function Navbar() {
  const { user } = useAuth()

  const handleLogin = async () => {
    try {
      const response = await login(window.location.pathname)
      if (response.data.redirect) {
        window.location.href = response.data.redirect
      }
    } catch {
      toast.error('Login failed')
    }
  }

  const handleLogout = async () => {
    try {
      const response = await logout()
      if (response.data.redirect) {
        window.location.href = response.data.redirect
      }
    } catch {
      toast.error('Logout failed')
    }
  }

  return (
    <nav style={{ background: '#333', padding: '10px 20px', display: 'flex', gap: '15px', alignItems: 'center' }}>
      <Link to="/" style={{ color: 'white', textDecoration: 'none', fontWeight: 'bold' }}>Bob's Used Bookstore</Link>
      <Link to="/search" style={{ color: 'white', textDecoration: 'none' }}>Search</Link>
      <Link to="/cart" style={{ color: 'white', textDecoration: 'none' }}>Cart</Link>
      <Link to="/wishlist" style={{ color: 'white', textDecoration: 'none' }}>Wishlist</Link>
      {user && (
        <>
          <Link to="/orders" style={{ color: 'white', textDecoration: 'none' }}>Orders</Link>
          <Link to="/addresses" style={{ color: 'white', textDecoration: 'none' }}>Addresses</Link>
          <Link to="/resale" style={{ color: 'white', textDecoration: 'none' }}>Resale</Link>
          {user.isAdmin && (
            <Link to="/admin/dashboard" style={{ color: 'white', textDecoration: 'none' }}>Admin</Link>
          )}
        </>
      )}
      <div style={{ marginLeft: 'auto' }}>
        {user ? (
          <span style={{ color: 'white' }}>
            {user.firstName} {user.lastName} |{' '}
            <button onClick={handleLogout} style={{ color: 'white', background: 'none', border: 'none', cursor: 'pointer' }}>Logout</button>
          </span>
        ) : (
          <button onClick={handleLogin} style={{ color: 'white', background: 'none', border: 'none', cursor: 'pointer' }}>Login</button>
        )}
      </div>
    </nav>
  )
}
