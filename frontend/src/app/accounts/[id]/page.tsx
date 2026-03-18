"use client";

import { useAuth } from "@/contexts/AuthContext";
import { api, type Account, type Transaction } from "@/lib/api";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

export default function AccountDetailPage() {
  const { token, isReady } = useAuth();
  const params = useParams();
  const router = useRouter();
  const id = Number(params.id);
  const [account, setAccount] = useState<Account | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [form, setForm] = useState({
    amount: 0,
    base_currency: "USD",
    type: "income",
    short_description: "",
  });
  const [submitting, setSubmitting] = useState(false);

  const load = useCallback(async () => {
    if (!id || isNaN(id)) return;
    const [accRes, txRes] = await Promise.all([
      api.getAccount(id),
      api.getAccountTransactions(id),
    ]);
    if (accRes.status === 401 || txRes.status === 401) {
      router.replace("/login");
      return;
    }
    if (accRes.error) setError(accRes.error);
    else setAccount(accRes.data ?? null);
    if (txRes.data) setTransactions(txRes.data);
    setLoading(false);
  }, [id, router]);

  useEffect(() => {
    if (!isReady) return;
    if (!token) {
      router.replace("/login");
      return;
    }
    load();
  }, [isReady, token, router, load]);

  async function handleAddTransaction(e: React.FormEvent) {
    e.preventDefault();
    if (!id || form.amount <= 0) return;
    setSubmitting(true);
    setError("");
    const { data, error: err, status } = await api.createTransaction(id, {
      amount: form.amount,
      base_currency: form.base_currency,
      type: form.type,
      short_description: form.short_description || undefined,
    });
    setSubmitting(false);
    if (status === 401) router.replace("/login");
    if (err) setError(err);
    else if (data) {
      setTransactions((prev) => [data, ...prev]);
      setForm((f) => ({ ...f, amount: 0, short_description: "" }));
    }
  }

  if (!isReady || !token) return <div className="flex min-h-screen items-center justify-center">Loading...</div>;
  if (!id || isNaN(id)) return <div className="p-6">Invalid account</div>;

  return (
    <main className="min-h-screen bg-gray-50">
      <header className="border-b border-gray-200 bg-white">
        <div className="mx-auto flex max-w-4xl items-center justify-between px-4 py-3">
          <Link href="/dashboard" className="text-sm text-gray-600 hover:text-gray-900">
            ← Dashboard
          </Link>
          <h1 className="text-lg font-semibold text-gray-900">
            {account ? account.name : "Account"}
          </h1>
          <span />
        </div>
      </header>
      <div className="mx-auto max-w-4xl px-4 py-6">
        {error && (
          <p className="mb-4 rounded bg-red-50 px-3 py-2 text-sm text-red-700" role="alert">
            {error}
          </p>
        )}
        {loading ? (
          <p className="text-sm text-gray-500">Loading…</p>
        ) : (
          <>
            {account && (
              <div className="mb-6 rounded-lg border border-gray-200 bg-white p-4">
                <p className="text-sm text-gray-500">Account name</p>
                <p className="font-medium text-gray-900">{account.name}</p>
                {account.balance != null && (
                  <p className="mt-1 text-sm text-gray-600">
                    Balance: {account.balance} {account.base_currence ?? ""}
                  </p>
                )}
              </div>
            )}
            <section>
              <h2 className="text-base font-medium text-gray-900">Add transaction</h2>
              <form onSubmit={handleAddTransaction} className="mt-2 grid gap-3 sm:grid-cols-2 lg:grid-cols-5">
                <input
                  type="number"
                  min={1}
                  value={form.amount || ""}
                  onChange={(e) => setForm((f) => ({ ...f, amount: Number(e.target.value) || 0 }))}
                  placeholder="Amount"
                  required
                  className="rounded-md border border-gray-300 px-3 py-2 text-gray-900"
                />
                <input
                  type="text"
                  value={form.base_currency}
                  onChange={(e) => setForm((f) => ({ ...f, base_currency: e.target.value }))}
                  placeholder="Currency"
                  className="rounded-md border border-gray-300 px-3 py-2 text-gray-900"
                />
                <select
                  value={form.type}
                  onChange={(e) => setForm((f) => ({ ...f, type: e.target.value }))}
                  className="rounded-md border border-gray-300 px-3 py-2 text-gray-900"
                >
                  <option value="income">Income</option>
                  <option value="expense">Expense</option>
                  <option value="transfer">Transfer</option>
                </select>
                <input
                  type="text"
                  value={form.short_description}
                  onChange={(e) => setForm((f) => ({ ...f, short_description: e.target.value }))}
                  placeholder="Description"
                  className="rounded-md border border-gray-300 px-3 py-2 text-gray-900"
                />
                <button
                  type="submit"
                  disabled={submitting || form.amount <= 0}
                  className="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
                >
                  {submitting ? "Adding…" : "Add"}
                </button>
              </form>
            </section>
            <section className="mt-8">
              <h2 className="text-base font-medium text-gray-900">Transactions</h2>
              {transactions.length === 0 ? (
                <p className="mt-2 text-sm text-gray-500">No transactions yet.</p>
              ) : (
                <ul className="mt-2 space-y-2">
                  {transactions.map((tx) => (
                    <li
                      key={tx.id}
                      className="rounded-lg border border-gray-200 bg-white px-4 py-3 text-sm"
                    >
                      <span className="font-medium text-gray-900">{tx.amount}</span>
                      <span className="ml-2 text-gray-600">{tx.base_currency}</span>
                      <span className="ml-2 text-gray-500">{tx.type}</span>
                      {tx.short_description && (
                        <span className="ml-2 text-gray-500">— {tx.short_description}</span>
                      )}
                    </li>
                  ))}
                </ul>
              )}
            </section>
          </>
        )}
      </div>
    </main>
  );
}