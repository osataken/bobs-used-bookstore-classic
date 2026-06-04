import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import type { ReactNode } from 'react';

export default function Layout({ children }: { children: ReactNode }) {
  const { auth, login, logout } = useAuth();

  return (
    <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <nav style={{ background: '#333', color: 'white', padding: '1rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
          <Link to="/" style={{ color: 'white', textDecoration: 'none', fontWeight: 'bold' }}>Bob's Used Bookstore</Link>
          <Link to="/search" style={{ color: 'white', textDecoration: 'none' }}>Browse</Link>
          <Link to="/cart" style={{ color: 'white', textDecoration: 'none' }}>Cart</Link>
          <Link to="/wishlist" style={{ color: 'white', textDecoration: 'none' }}>Wishlist</Link>
          {auth.authenticated && (
            <>
              <Link to="/orders" style={{ color: 'white', textDecoration: 'none' }}>Orders</Link>
              <Link to="/addresses" style={{ color: 'white', textDecoration: 'none' }}>Addresses</Link>
              <Link to="/resale" style={{ color: 'white', textDecoration: 'none' }}>Sell Books</Link>
            </>
          )}
          {auth.role === 'Administrators' && (
            <Link to="/admin" style={{ color: '#ffd700', textDecoration: 'none' }}>Admin</Link>
          )}
        </div>
        <div>
          {auth.authenticated ? (
            <span>
              <span style={{ marginRight: '1rem' }}>{auth.username}</span>
              <button onClick={logout} style={{ cursor: 'pointer' }}>Logout</button>
            </span>
          ) : (
            <button onClick={login} style={{ cursor: 'pointer' }}>Login</button>
          )}
        </div>
      </nav>
      <main style={{ flex: 1, padding: '2rem', maxWidth: '1200px', margin: '0 auto', width: '100%' }}>
        {children}
      </main>
    </div>
  );
}
