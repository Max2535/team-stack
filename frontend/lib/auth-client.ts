'use client';

import { apiClient } from './api-client';
import type { ApiResponse, User } from './types';

const TOKEN_KEY = 'team_token';

export async function login(email: string, password: string) {
  const res = await apiClient.post<ApiResponse<{ token: string; user: User }>>(
    '/v1/auth/login',
    { email, password }
  );
  if (res.success && res.data) {
    localStorage.setItem(TOKEN_KEY, res.data.token);
  }
  return res;
}

export function logout() {
  if (typeof window === 'undefined') return;
  localStorage.removeItem(TOKEN_KEY);
}

export function getToken(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(TOKEN_KEY);
}
