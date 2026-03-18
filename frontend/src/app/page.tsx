"use client";

import { useAuth } from "@/contexts/AuthContext";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function Home() {
  const { token, isReady } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isReady) return;
    if (token) router.replace("/dashboard");
  }, [isReady, token, router]);

  if (!isReady) return <div className="flex min-h-screen items-center justify-center">Loading...</div>;
  if (token) return null;

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8">
      <h1 className="text-3xl font-bold text-gray-900">Web Application</h1>
      <p className="mt-2 text-gray-600">
        Sign in or create an account to manage your accounts and transactions.
      </p>
      <div className="mt-6 flex gap-4">
        <Link
          href="/login"
          className="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          Sign in
        </Link>
        <Link
          href="/signup"
          className="rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
        >
          Sign up
        </Link>
      </div>
      <p className="mt-8 text-sm text-gray-500">
        Backend API: <code className="rounded bg-gray-200 px-1">/api/*</code> (proxied to Go server).
        <br />
        Run backend with Air: <code className="rounded bg-gray-200 px-1">cd server && make watch</code>
      </p>
    </main>
  );
}
