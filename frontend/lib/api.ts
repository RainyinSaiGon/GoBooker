/**
 * lib/api.ts
 *
 * Typed API layer. All network calls go through the functions here so the
 * rest of the app never hard-codes fetch() calls or URL strings.
 */

import { useAuthStore } from "./store/auth";

// ─── Base URL ─────────────────────────────────────────────────────────────────

export const API_BASE: string =
  process.env.NEXT_PUBLIC_API_URL
    ? `${process.env.NEXT_PUBLIC_API_URL}/api/v1`
    : "http://localhost:3001/api/v1";

// ─── Shared types ─────────────────────────────────────────────────────────────

export interface User {
  id?:    string;
  name:   string;
  email:  string;
  role:   string;
}

export interface PaginatedUsers {
  users:      User[];
  total:      number;
  page:       number;
  size:       number;
  totalPages: number;
}

export interface GetUsersParams {
  query:    string;
  page:     number;
  pageSize: number;
}

export interface CreateUserPayload {
  name:     string;
  email:    string;
  password: string;
}

export interface LoginPayload {
  email:    string;
  password: string;
}

export interface LoginResponse {
  token:        string;
  refreshToken?: string;
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

async function handleResponse<T>(res: Response): Promise<T> {
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(
      (body as { error?: string }).error ??
        `Request failed with status ${res.status}`
    );
  }
  return res.json() as Promise<T>;
}

async function authenticatedFetch(url: string, options: RequestInit = {}): Promise<Response> {
  const headers = (options.headers as Record<string, string>) || {};
  const token = useAuthStore.getState().token;
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }
  options.headers = headers;
  options.credentials = "include";

  let res = await fetch(url, options);

  // If 401 and we have a token (expired), try to refresh
  if (res.status === 401 && token) {
    try {
      const refreshRes = await fetch(`${API_BASE}/auth/refresh`, {
        method: "POST",
        credentials: "include",
      });
      if (refreshRes.ok) {
        const refreshData = await refreshRes.json();
        const newToken = refreshData.token;
        useAuthStore.getState().setTokens(newToken);
        
        // Retry with new token
        headers["Authorization"] = `Bearer ${newToken}`;
        options.headers = headers;
        res = await fetch(url, options);
      } else {
        useAuthStore.getState().logout();
      }
    } catch (e) {
      console.error("Token refresh failed:", e);
      useAuthStore.getState().logout();
    }
  }
  return res;
}

// ─── API functions ────────────────────────────────────────────────────────────

export async function apiGetUsers({
  query,
  page,
  pageSize,
}: GetUsersParams): Promise<PaginatedUsers> {
  const params = new URLSearchParams({
    query: query,
    page:  String(page),
    size:  String(pageSize),
  });
  
  const res = await authenticatedFetch(`${API_BASE}/users?${params.toString()}`, {
    cache: "no-store",
  });
  return handleResponse<PaginatedUsers>(res);
}

export async function apiCreateUser(payload: CreateUserPayload): Promise<User> {
  const res = await fetch(`${API_BASE}/users`, {
    method:  "POST",
    headers: { "Content-Type": "application/json" },
    body:    JSON.stringify(payload),
  });
  return handleResponse<User>(res);
}

export async function apiLogin(payload: LoginPayload): Promise<LoginResponse> {
  const res = await fetch(`${API_BASE}/auth/login`, {
    method:  "POST",
    headers: { "Content-Type": "application/json" },
    body:    JSON.stringify(payload),
    credentials: "include",
  });
  return handleResponse<LoginResponse>(res);
}

export async function apiRefreshToken(): Promise<LoginResponse> {
  const res = await fetch(`${API_BASE}/auth/refresh`, {
    method:  "POST",
    credentials: "include",
  });
  return handleResponse<LoginResponse>(res);
}

export async function apiLogout(): Promise<{ status: string }> {
  const res = await fetch(`${API_BASE}/auth/logout`, {
    method:  "POST",
    credentials: "include",
  });
  return handleResponse<{ status: string }>(res);
}

