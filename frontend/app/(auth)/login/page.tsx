'use client';

import { useState } from 'react';
import { login } from '@/lib/auth-client';
import { useRouter } from 'next/navigation';

export default function LoginPage() {
  const [email, setEmail] = useState('user@example.com');
  const [password, setPassword] = useState('password');
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    try {
      const res = await login(email, password);
      if (!res.success) {
        setError(res.error?.message ?? 'Login failed');
        return;
      }
      router.push('/dashboard');
    } catch (err) {
      setError((err as Error).message);
    }
  }

  return (
    <main className="min-h-screen flex items-center justify-center">
      <form
        onSubmit={onSubmit}
        className="w-full max-w-sm space-y-4 border bg-white p-6 rounded-xl shadow"
      >
        <h1 className="text-xl font-bold text-center">Login</h1>
        {error && <p className="text-sm text-red-500">{error}</p>}
        <div className="space-y-1">
          <label className="text-sm font-medium">Email</label>
          <input
            className="w-full border rounded px-3 py-2 text-sm"
            value={email}
            onChange={e => setEmail(e.target.value)}
          />
        </div>
        <div className="space-y-1">
          <label className="text-sm font-medium">Password</label>
          <input
            type="password"
            className="w-full border rounded px-3 py-2 text-sm"
            value={password}
            onChange={e => setPassword(e.target.value)}
          />
        </div>
        <button
          type="submit"
          className="w-full rounded bg-blue-600 text-white py-2 text-sm font-medium"
        >
          Sign in
        </button>
      </form>
    </main>
  );
}
