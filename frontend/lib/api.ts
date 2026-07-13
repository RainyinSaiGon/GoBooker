/**
 * lib/api.ts
 *
 * Typed API layer. All network calls go through the functions here so the
 * rest of the app never hard-codes fetch() calls or URL strings.
 */

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
  refreshToken: string;
}

// ─── Helper ───────────────────────────────────────────────────────────────────

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
  const res = await fetch(`${API_BASE}/users?${params.toString()}`, {
    // next: { revalidate: 0 } lets the Server Component always fetch fresh data
    // for SSR prefetch; the client-side cache is managed by TanStack Query.
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
  });
  return handleResponse<LoginResponse>(res);
}
