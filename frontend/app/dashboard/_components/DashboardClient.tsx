/**
 * app/dashboard/_components/DashboardClient.tsx
 *
 * Client shell for the dashboard page. Reads search/pagination state from
 * URL params (so links are shareable) and calls useUsers() which gets an
 * instant cache hit from the dehydrated SSR prefetch on first paint.
 */
"use client";

import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { useUsers } from "@/lib/queries/users";
import { PAGE_SIZE_OPTIONS, type PageSize } from "@/lib/constants";
import SearchBar from "./SearchBar";
import Pagination from "./Pagination";
import UserList from "./UserList";

export default function DashboardClient() {
  const searchParams = useSearchParams();

  const query    = searchParams.get("query") ?? "";
  const page     = Math.max(1, Number(searchParams.get("page")) || 1);
  const rawSize  = Number(searchParams.get("size"));
  const pageSize: PageSize = PAGE_SIZE_OPTIONS.includes(rawSize as PageSize)
    ? (rawSize as PageSize)
    : 10;

  const { data, isLoading, isError, error } = useUsers({ query, page, pageSize });

  const users      = data?.users      ?? [];
  const total      = data?.total      ?? 0;
  const totalPages = data?.totalPages ?? 1;
  const safePage   = Math.min(page, Math.max(1, totalPages));

  const from = total > 0 ? (safePage - 1) * pageSize + 1 : 0;
  const to   = Math.min(safePage * pageSize, total);

  return (
    <div className="animate-fade-in mx-auto max-w-3xl space-y-8">

      {/* ── Header ── */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-[var(--text-primary)]">Users</h1>
          <p className="mt-0.5 text-sm text-[var(--text-secondary)]">
            {isLoading
              ? "Loading…"
              : query
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

      {/* ── Search + refresh ── */}
      <SearchBar />

      {/* ── User cards ── */}
      <UserList
        users={users}
        isLoading={isLoading}
        error={isError ? (error?.message ?? "Failed to load users") : null}
        search={query}
      />

      {/* ── Pagination ── */}
      {!isError && total > 0 && (
        <Pagination
          page={safePage}
          totalPages={totalPages}
          pageSize={pageSize}
          from={from}
          to={to}
          total={total}
        />
      )}
    </div>
  );
}
