import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { getAuthStatus, login as apiLogin, logout as apiLogout } from '../services/api';
import type { AuthStatus } from '../types';

interface AuthContextType {
  auth: AuthStatus;
  loading: boolean;
  login: () => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [auth, setAuth] = useState<AuthStatus>({ authenticated: false });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getAuthStatus().then(setAuth).finally(() => setLoading(false));
  }, []);

  const login = async () => {
    await apiLogin();
    const status = await getAuthStatus();
    setAuth(status);
  };

  const logout = async () => {
    await apiLogout();
    setAuth({ authenticated: false });
  };

  return (
    <AuthContext.Provider value={{ auth, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
}
