import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { Customer } from '../types';
import { authService } from '../services';

interface AuthContextType {
  user: Customer | null;
  loading: boolean;
  login: () => void;
  logout: () => void;
  refreshUser: () => void;
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  loading: true,
  login: () => {},
  logout: () => {},
  refreshUser: () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<Customer | null>(null);
  const [loading, setLoading] = useState(true);

  const refreshUser = async () => {
    try {
      const res = await authService.me();
      setUser(res.data);
    } catch {
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    refreshUser();
  }, []);

  const login = async () => {
    try {
      const res = await authService.login();
      if (res.data.redirectUrl) {
        window.location.href = res.data.redirectUrl;
      } else {
        await refreshUser();
      }
    } catch (e) {
      console.error('Login failed', e);
    }
  };

  const logout = async () => {
    try {
      const res = await authService.logout();
      setUser(null);
      if (res.data.redirectUrl) {
        window.location.href = res.data.redirectUrl;
      }
    } catch (e) {
      console.error('Logout failed', e);
    }
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, refreshUser }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
