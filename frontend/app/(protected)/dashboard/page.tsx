'use client';

import { useAuth } from '@/hooks/useAuth';

export default function DashboardPage() {
  const { user, logout } = useAuth();

  return (
    <main className="min-h-screen p-8 bg-slate-50">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <div className="flex items-center gap-3">
          <span className="text-sm text-slate-600">
            {user?.name} ({user?.email})
          </span>
          <button
            onClick={logout}
            className="px-3 py-1 text-sm rounded bg-red-500 text-white"
          >
            Logout
          </button>
        </div>
      </div>
      <section className="grid gap-4 md:grid-cols-3">
        <div className="rounded-xl border bg-white p-4 shadow-sm">
          <h2 className="font-semibold mb-2 text-slate-800">Summary</h2>
          <p className="text-sm text-slate-600">
            ใส่ widget, metrics, charts ฯลฯ ของระบบตรงนี้
          </p>
        </div>
      </section>
    </main>
  );
}
