import { Outlet, Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

function Layout() {
  const { user, login, logout } = useAuth()

  return (
    <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <nav style={{ padding: '1rem', background: '#333', color: '#fff', display: 'flex', gap: '1rem', alignItems: 'center' }}>
        <Link to="/" style={{ color: '#fff', textDecoration: 'none', fontWeight: 'bold' }}>Bob's Bookstore</Link>
        <Link to="/search" style={{ color: '#fff', textDecoration: 'none' }}>Browse</Link>
        <Link to="/cart" style={{ color: '#fff', textDecoration: 'none' }}>Cart</Link>
        <Link to="/wishlist" style={{ color: '#fff', textDecoration: 'none' }}>Wishlist</Link>
        {user && (
          <>
            <Link to="/orders" style={{ color: '#fff', textDecoration: 'none' }}>Orders</Link>
            <Link to="/addresses" style={{ color: '#fff', textDecoration: 'none' }}>Addresses</Link>
            <Link to="/resale" style={{ color: '#fff', textDecoration: 'none' }}>Sell Books</Link>
            {user.role === 'Administrators' && (
              <Link to="/admin" style={{ color: '#fff', textDecoration: 'none' }}>Admin</Link>
            )}
          </>
        )}
        <div style={{ marginLeft: 'auto' }}>
          {user ? (
            <span>
              {user.firstName} {user.lastName}{' '}
              <button onClick={logout} style={{ marginLeft: '0.5rem' }}>Logout</button>
            </span>
          ) : (
            <button onClick={login}>Login</button>
          )}
        </div>
      </nav>
      <main style={{ flex: 1, padding: '2rem' }}>
        <Outlet />
      </main>
    </div>
  )
}

export default Layout
