"use client";

import { useAuth } from "@/contexts/AuthContext";
import { api, type Account } from "@/lib/api";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

export default function DashboardPage() {
  const { token, isReady, logout } = useAuth();
  const router = useRouter();
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [createName, setCreateName] = useState("");
  const [creating, setCreating] = useState(false);

  const load = useCallback(async () => {
    const { data, error: err, status } = await api.getAccounts();
    if (status === 401) {
      router.replace("/login");
      return;
    }
    if (err) setError(err);
    else setAccounts(data ?? []);
    setLoading(false);
  }, [router]);

  useEffect(() => {
    if (!isReady) return;
    if (!token) {
      router.replace("/login");
      return;
    }
    load();
  }, [isReady, token, router, load]);

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!createName.trim()) return;
    setCreating(true);
    setError("");
    const { data, error: err, status } = await api.createAccount(createName.trim());
    setCreating(false);
    if (status === 401) router.replace("/login");
    if (err) setError(err);
    else if (data) {
      setAccounts((prev) => [...prev, data]);
      setCreateName("");
    }
  }

  if (!isReady || !token) return <div className="flex min-h-screen items-center justify-center">Loading...</div>;

  return (
    <main className="min-h-screen bg-gray-50">
      <header className="border-b border-gray-200 bg-white">
        <div className="mx-auto flex max-w-4xl items-center justify-between px-4 py-3">
          <h1 className="text-lg font-semibold text-gray-900">Dashboard</h1>
          <div className="flex items-center gap-4">
            <Link href="/dashboard" className="text-sm text-gray-600 hover:text-gray-900">
              Accounts
            </Link>
            <button
              type="button"
              onClick={() => logout()}
              className="text-sm text-gray-600 hover:text-gray-900"
            >
              Sign out
            </button>
          </div>
        </div>
      </header>
      <div className="mx-auto max-w-4xl px-4 py-6">
        {error && (
          <p className="mb-4 rounded bg-red-50 px-3 py-2 text-sm text-red-700" role="alert">
            {error}
          </p>
        )}
        <section>
          <h2 className="text-base font-medium text-gray-900">Add account</h2>
          <form onSubmit={handleCreate} className="mt-2 flex gap-2">
            <input
              type="text"
              value={createName}
              onChange={(e) => setCreateName(e.target.value)}
              placeholder="Account name"
              className="rounded-md border border-gray-300 px-3 py-2 text-gray-900 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
            <button
              type="submit"
              disabled={creating || !createName.trim()}
              className="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            >
              {creating ? "Adding…" : "Add"}
            </button>
          </form>
        </section>
        <section className="mt-8">
          <h2 className="text-base font-medium text-gray-900">Your accounts</h2>
          {loading ? (
            <p className="mt-2 text-sm text-gray-500">Loading…</p>
          ) : accounts.length === 0 ? (
            <p className="mt-2 text-sm text-gray-500">No accounts yet. Create one above.</p>
          ) : (
            <ul className="mt-2 space-y-2">
              {accounts.map((acc) => (
                <li key={acc.id}>
                  <Link
                    href={`/accounts/${acc.id}`}
                    className="block rounded-lg border border-gray-200 bg-white px-4 py-3 shadow-sm hover:border-gray-300"
                  >
                    <span className="font-medium text-gray-900">{acc.name}</span>
                    {acc.balance != null && (
                      <span className="ml-2 text-sm text-gray-500">
                        Balance: {acc.balance} {acc.base_currence ?? ""}
                      </span>
                    )}
                  </Link>
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>
    </main>
  );
}