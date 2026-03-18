const API_BASE = "/api";

function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("token");
}

export function setToken(token: string): void {
  if (typeof window === "undefined") return;
  localStorage.setItem("token", token);
}

export function clearToken(): void {
  if (typeof window === "undefined") return;
  localStorage.removeItem("token");
}

export type SignInResponse = {
  token: string;
  expires_at: number;
  user_id: number;
  login: string;
};

export type Account = {
  id: number;
  name: string;
  base_currence?: string;
  balance?: number;
  created_at?: string;
  updated_at?: string;
};

export type Transaction = {
  id: number;
  amount: number;
  base_currency: string;
  type: string;
  short_description?: string;
  account_id: number;
  user_id: number;
  created_at?: string;
};

async function apiFetch<T>(
  path: string,
  options: RequestInit = {}
): Promise<{ data?: T; error?: string; status: number }> {
  const token = getToken();
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };
  if (token) (headers as Record<string, string>)["Authorization"] = `Bearer ${token}`;
  const res = await fetch(`${API_BASE}${path}`, { ...options, headers });
  let error: string | undefined;
  if (!res.ok) {
    const text = await res.text();
    try {
      const j = JSON.parse(text);
      error = (j.error ?? j.message ?? text) || res.statusText;
    } catch {
      error = text || res.statusText;
    }
  }
  let data: T | undefined;
  if (res.ok && res.status !== 204) {
    const text = await res.text();
    if (text) {
      try {
        data = JSON.parse(text) as T;
      } catch {
        data = text as unknown as T;
      }
    }
  }
  return { data, error, status: res.status };
}

export const api = {
  async signIn(login: string, password: string) {
    return apiFetch<SignInResponse>("/signin", {
      method: "POST",
      body: JSON.stringify({ login, password }),
    });
  },
  async signUp(login: string, password: string) {
    return apiFetch<SignInResponse>("/signup", {
      method: "POST",
      body: JSON.stringify({ login, password }),
    });
  },
  async getAccounts() {
    return apiFetch<Account[]>("/accounts");
  },
  async getAccount(id: number) {
    return apiFetch<Account>(`/accounts/${id}`);
  },
  async createAccount(name: string) {
    return apiFetch<Account>("/accounts", {
      method: "POST",
      body: JSON.stringify({ name }),
    });
  },
  async updateAccount(id: number, body: { name?: string }) {
    return apiFetch<Account>(`/accounts/${id}`, {
      method: "PUT",
      body: JSON.stringify(body),
    });
  },
  async deleteAccount(id: number) {
    return apiFetch<unknown>(`/accounts/${id}`, { method: "DELETE" });
  },
  async getAccountTransactions(accountId: number) {
    return apiFetch<Transaction[]>(`/accounts/${accountId}/transactions`);
  },
  async createTransaction(
    accountId: number,
    body: { amount: number; base_currency: string; type: string; short_description?: string }
  ) {
    return apiFetch<Transaction>(`/accounts/${accountId}/transactions`, {
      method: "POST",
      body: JSON.stringify(body),
    });
  },
};
