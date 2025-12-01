'use client';

import { createContext, useContext, useEffect, useState } from 'react';
import type { User } from '@/lib/types';
import { apiClient } from '@/lib/api-client';
import { getToken, logout as doLogout } from '@/lib/auth-client';

type AuthState = {
  user: User | null;
  loading: boolean;
};

type AuthContextType = AuthState & {
  refresh: () => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  async function refresh() {
    const token = getToken();
    if (!token) {
      setUser(null);
      setLoading(false);
      return;
    }
    try {
      const res = await apiClient.get<{ success: boolean; data?: User }>(
        '/v1/me',
        token
      );
      if (res.success && res.data) {
        setUser(res.data);
      } else {
        setUser(null);
      }
    } catch {
      setUser(null);
    }
    setLoading(false);
  }

  useEffect(() => {
    void refresh();
  }, []);

  function logout() {
    doLogout();
    setUser(null);
  }

  return (
    <AuthContext.Provider value={{ user, loading, refresh, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return ctx;
}
