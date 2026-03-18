"use client";

import React, { createContext, useCallback, useContext, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { api, setToken, clearToken } from "@/lib/api";

type AuthState = {
  token: string | null;
  isReady: boolean;
  login: (login: string, password: string) => Promise<{ error?: string }>;
  signUp: (login: string, password: string) => Promise<{ error?: string }>;
  logout: () => void;
};

const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setTokenState] = useState<string | null>(null);
  const [isReady, setIsReady] = useState(false);
  const router = useRouter();

  useEffect(() => {
    const t = typeof window !== "undefined" ? localStorage.getItem("token") : null;
    setTokenState(t);
    setIsReady(true);
  }, []);

  const login = useCallback(async (loginId: string, password: string) => {
    const { data, error, status } = await api.signIn(loginId, password);
    if (error) return { error: status === 401 ? "Invalid login or password" : error };
    if (data?.token) {
      setToken(data.token);
      setTokenState(data.token);
      router.push("/dashboard");
      return {};
    }
    return { error: "Login failed" };
  }, [router]);

  const signUp = useCallback(async (loginId: string, password: string) => {
    const { data, error, status } = await api.signUp(loginId, password);
    if (error) return { error: status === 409 ? "Login already taken" : error };
    if (data?.token) {
      setToken(data.token);
      setTokenState(data.token);
      router.push("/dashboard");
      return {};
    }
    return { error: "Sign up failed" };
  }, [router]);

  const logout = useCallback(() => {
    clearToken();
    setTokenState(null);
    router.push("/login");
  }, [router]);

  return (
    <AuthContext.Provider value={{ token, isReady, login, signUp, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
