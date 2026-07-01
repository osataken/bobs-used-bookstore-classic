import { Link, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

function Layout() {
  const { user, login, logout } = useAuth();

  const isAdmin = user?.roles?.includes('Administrators');

  return (
    <div>
      <nav style={{ padding: '1rem', background: '#333', color: '#fff', display: 'flex', gap: '1rem', alignItems: 'center' }}>
        <Link to="/" style={{ color: '#fff', textDecoration: 'none', fontWeight: 'bold' }}>Bob's Used Bookstore</Link>
        <Link to="/search" style={{ color: '#fff', textDecoration: 'none' }}>Browse</Link>
        <Link to="/cart" style={{ color: '#fff', textDecoration: 'none' }}>Cart</Link>
        <Link to="/wishlist" style={{ color: '#fff', textDecoration: 'none' }}>Wishlist</Link>
        {user?.authenticated && (
          <>
            <Link to="/orders" style={{ color: '#fff', textDecoration: 'none' }}>Orders</Link>
            <Link to="/resale" style={{ color: '#fff', textDecoration: 'none' }}>Resale</Link>
            <Link to="/addresses" style={{ color: '#fff', textDecoration: 'none' }}>Addresses</Link>
          </>
        )}
        {isAdmin && (
          <Link to="/admin" style={{ color: '#ffd700', textDecoration: 'none' }}>Admin</Link>
        )}
        <div style={{ marginLeft: 'auto' }}>
          {user?.authenticated ? (
            <button onClick={logout} style={{ cursor: 'pointer' }}>Logout ({user.username})</button>
          ) : (
            <button onClick={login} style={{ cursor: 'pointer' }}>Login</button>
          )}
        </div>
      </nav>
      <main style={{ padding: '1rem', maxWidth: '1200px', margin: '0 auto' }}>
        <Outlet />
      </main>
    </div>
  );
}

export default Layout;
