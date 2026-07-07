"use client";

import { useState, useEffect, useCallback } from "react";
import Link from "next/link";

const API = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001";

interface User { name: string; email: string; }

function SkeletonRow() {
  return (
    <div className="flex items-center gap-4 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] p-4">
      <div className="skeleton h-10 w-10 rounded-full shrink-0" />
      <div className="flex-1 space-y-2">
        <div className="skeleton h-3.5 w-32 rounded" />
        <div className="skeleton h-3 w-48 rounded" />
      </div>
    </div>
  );
}

function Avatar({ name }: { name: string }) {
  const initials = name.split(" ").map((n) => n[0]).join("").slice(0, 2).toUpperCase();
  const hue = [...name].reduce((acc, c) => acc + c.charCodeAt(0), 0) % 360;
  return (
    <div
      className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-sm font-semibold text-white"
      style={{ background: `hsl(${hue},60%,48%)` }}
      aria-hidden
    >
      {initials}
    </div>
  );
}

export default function DashboardPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState("");

  const fetchUsers = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await fetch(`${API}/users`);
      if (!res.ok) throw new Error("Failed to fetch users");
      setUsers((await res.json()) ?? []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong");
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => { fetchUsers(); }, [fetchUsers]);

  const filtered = users.filter(
    (u) =>
      u.name.toLowerCase().includes(search.toLowerCase()) ||
      u.email.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="animate-fade-in mx-auto max-w-3xl space-y-8">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-[var(--text-primary)]">Users</h1>
          <p className="mt-0.5 text-sm text-[var(--text-secondary)]">
            {isLoading ? "Loading…" : `${users.length} registered user${users.length !== 1 ? "s" : ""}`}
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

      {/* Search + refresh row */}
      <div className="flex items-center gap-3">
        <div className="relative flex-1">
          <svg
            className="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-[var(--text-muted)]"
            width="15" height="15" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}
          >
            <circle cx="11" cy="11" r="8"/><path strokeLinecap="round" d="m21 21-4.35-4.35"/>
          </svg>
          <input
            type="search" value={search} onChange={(e) => setSearch(e.target.value)}
            placeholder="Search by name or email…"
            className="w-full rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] py-2.5 pl-10 pr-4 text-sm text-[var(--text-primary)] placeholder:text-[var(--text-muted)] transition focus:border-[var(--brand-500)] focus:outline-none focus:ring-2 focus:ring-[var(--brand-200)]"
          />
        </div>
        <button
          onClick={fetchUsers} disabled={isLoading}
          title="Refresh"
          className="flex h-10 w-10 items-center justify-center rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] text-[var(--text-secondary)] transition hover:bg-[var(--bg-subtle)] hover:text-[var(--brand-500)] disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-[var(--brand-200)]"
        >
          <svg width="15" height="15" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2} className={isLoading ? "animate-spin" : ""}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M4 4v6h6M20 20v-6h-6M4.93 15A9 9 0 1019.07 9" />
          </svg>
        </button>
      </div>

      {/* States */}
      {isLoading && (
        <div className="space-y-3">
          {[...Array(5)].map((_, i) => <SkeletonRow key={i} />)}
        </div>
      )}

      {!isLoading && error && (
        <div className="flex flex-col items-center gap-4 rounded-2xl border border-red-200 bg-red-50 py-12 text-center dark:border-red-900 dark:bg-red-950/20">
          <div className="flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-950/50">
            <svg className="text-red-500" width="22" height="22" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.8}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v3m0 3h.01M5.07 19H19a2 2 0 001.73-3L13.73 4a2 2 0 00-3.46 0L4.27 16A2 2 0 005.07 19z"/>
            </svg>
          </div>
          <div>
            <p className="font-semibold text-red-600 dark:text-red-400">Could not load users</p>
            <p className="mt-0.5 text-sm text-red-500/80 dark:text-red-400/70">{error}</p>
          </div>
          <button
            onClick={fetchUsers}
            className="rounded-xl border border-red-300 px-4 py-1.5 text-sm text-red-600 transition hover:bg-red-100 dark:border-red-800 dark:text-red-400 dark:hover:bg-red-950/40"
          >
            Try again
          </button>
        </div>
      )}

      {!isLoading && !error && filtered.length === 0 && (
        <div className="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-[var(--border)] bg-[var(--bg-surface)] py-16 text-center">
          <div className="flex h-12 w-12 items-center justify-center rounded-full bg-[var(--bg-subtle)]">
            <svg className="text-[var(--text-muted)]" width="20" height="20" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.6}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2M9 11a4 4 0 100-8 4 4 0 000 8z"/>
            </svg>
          </div>
          <div>
            <p className="font-semibold text-[var(--text-primary)]">{search ? "No matches found" : "No users yet"}</p>
            <p className="mt-0.5 text-sm text-[var(--text-secondary)]">
              {search ? `No users match "${search}"` : "Add your first user to get started."}
            </p>
          </div>
          {!search && (
            <Link href="/" className="mt-1 rounded-xl bg-[var(--brand-500)] px-4 py-2 text-sm font-semibold text-white transition hover:bg-[var(--brand-600)]">
              Add First User
            </Link>
          )}
        </div>
      )}

      {!isLoading && !error && filtered.length > 0 && (
        <div className="space-y-2.5">
          {filtered.map((user, i) => (
            <div
              key={i}
              className="group flex items-center gap-4 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] px-4 py-3.5 shadow-[var(--shadow-sm)] transition hover:shadow-[var(--shadow-md)] hover:border-[var(--brand-200)]"
              style={{ animationDelay: `${i * 40}ms` }}
            >
              <Avatar name={user.name} />
              <div className="flex-1 min-w-0">
                <p className="truncate text-sm font-semibold text-[var(--text-primary)] group-hover:text-[var(--brand-600)] transition-colors">
                  {user.name}
                </p>
                <p className="truncate text-xs text-[var(--text-muted)]">{user.email}</p>
              </div>
              <div className="shrink-0 flex items-center gap-1.5">
                <span className="inline-flex items-center rounded-full bg-[var(--bg-subtle)] px-2.5 py-0.5 text-[10px] font-semibold uppercase tracking-wide text-[var(--brand-600)]">
                  customer
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
