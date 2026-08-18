import { Outlet, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function Layout() {
  const { user, login, logout } = useAuth();

  return (
    <>
      <nav>
        <div className="container" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <Link to="/">Home</Link>
            <Link to="/search">Browse</Link>
            <Link to="/cart">Cart</Link>
            <Link to="/wishlist">Wishlist</Link>
            {user && <Link to="/orders">Orders</Link>}
            {user && <Link to="/addresses">Addresses</Link>}
            {user && <Link to="/resale">Sell Books</Link>}
            {user?.isAdmin && <Link to="/admin">Admin</Link>}
          </div>
          <div>
            {user ? (
              <span>
                {user.username} | <button onClick={logout} style={{ background: 'none', border: 'none', color: 'white', cursor: 'pointer' }}>Logout</button>
              </span>
            ) : (
              <button onClick={login} style={{ background: 'none', border: 'none', color: 'white', cursor: 'pointer' }}>Login</button>
            )}
          </div>
        </div>
      </nav>
      <main className="container">
        <Outlet />
      </main>
    </>
  );
}
