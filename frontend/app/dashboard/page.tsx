import { API_BASE } from "@/lib/api";
import Link from "next/link";
import SearchBar from "./_components/SearchBar";
import Pagination from "./_components/Pagination";
import UserList from "./_components/UserList";

// ─── Types ───────────────────────────────────────────────────────────────────

export interface User { name: string; email: string; role: string; }

export interface PaginatedUsers {
  users: User[];
  total: number;
  page: number;
  size: number;
  totalPages: number;
}

// ─── Constants ───────────────────────────────────────────────────────────────

export const PAGE_SIZE_OPTIONS = [5, 10, 25] as const;
export type PageSize = typeof PAGE_SIZE_OPTIONS[number];

// ─── Data fetching ───────────────────────────────────────────────────────────

async function fetchUsers(
  query: string,
  page: number,
  size: number
): Promise<{ data: PaginatedUsers | null; error: string | null }> {
  try {
    const res = await fetch(
      `${API_BASE}/users?query=${encodeURIComponent(query)}&page=${page}&size=${size}`,
    );
    if (!res.ok) throw new Error(`Server responded with ${res.status}`);
    const json = await res.json();
    return { data: json, error: null };
  } catch (err) {
    return { data: null, error: err instanceof Error ? err.message : "Failed to load users" };
  }
}

// ─── Page (Server Component) ─────────────────────────────────────────────────

export default async function DashboardPage({
  searchParams,
}: {
  searchParams: Promise<{ query?: string; page?: string; size?: string }>;
}) {
  const params = await searchParams;

  const query = params.query ?? "";
  const page  = Math.max(1, Number(params.page) || 1);
  const size: PageSize = PAGE_SIZE_OPTIONS.includes(Number(params.size) as PageSize)
    ? (Number(params.size) as PageSize)
    : 10;

  const { data, error } = await fetchUsers(query, page, size);

  const users      = data?.users ?? [];
  const total      = data?.total ?? 0;
  const totalPages = data?.totalPages ?? 1;
  const safePage   = Math.min(page, totalPages);

  // Pagination bounds
  const from = total > 0 ? (safePage - 1) * size + 1 : 0;
  const to   = Math.min(safePage * size, total);

  return (
    <div className="animate-fade-in mx-auto max-w-3xl space-y-8">

      {/* ── Header ── */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-[var(--text-primary)]">Users</h1>
          <p className="mt-0.5 text-sm text-[var(--text-secondary)]">
            {query
              ? `Found ${total} matching user${total !== 1 ? "s" : ""}`
              : `${total} user${total !== 1 ? "s" : ""}`}
          </p>
        </div>
        <Link
          href="/"
          className="inline-flex items-center gap-2 self-start rounded-xl bg-[var(--brand-500)] px-4 py-2 text-sm font-semibold text-white shadow-md transition hover:bg-[var(--brand-600)] active:scale-[.98] focus:outline-none focus:ring-2 focus:ring-[var(--brand-400)] focus:ring-offset-2"
        >
          <svg width="14" height="14" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.4}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
          </svg>
          Add User
        </Link>
      </div>

      {/* ── Search + page-size + refresh (client) ── */}
      <SearchBar />

      {/* ── User cards ── */}
      <UserList users={users} error={error} search={query} />

      {/* ── Pagination (client) — only shown when there is something to page ── */}
      {!error && total > 0 && (
        <Pagination
          page={safePage}
          totalPages={totalPages}
          pageSize={size}
          from={from}
          to={to}
          total={total}
        />
      )}
    </div>
  );
}
