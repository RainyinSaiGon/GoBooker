"use client";

/**
 * SearchBar — Client Component
 *
 * Implements the Next.js recommended search pattern from the official tutorial
 * (nextjs.org/learn — Chapter 10: Adding Search and Pagination):
 *
 *   1. Read current params with useSearchParams().
 *   2. On input change, debounce (300 ms) then call router.replace() so the
 *      Server Component page re-renders with the new query in the URL.
 *      Using `replace` (not `push`) keeps the browser history clean — the user
 *      won't have to hit Back dozens of times after typing a search query.
 *   3. Reset `page` to "1" whenever the query or page size changes so the user
 *      never lands on an empty page.
 *   4. The refresh button calls router.refresh(), which re-fetches the server
 *      data without a full navigation.
 */

import { useSearchParams, usePathname, useRouter } from "next/navigation";
import { useRef, useCallback, useTransition, useState } from "react";
import { PAGE_SIZE_OPTIONS, type PageSize } from "../page";

export default function SearchBar() {
  const searchParams  = useSearchParams();
  const pathname      = usePathname();
  const { replace, refresh } = useRouter();
  const debounceRef   = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [isPending, startTransition] = useTransition();
  const [isRefreshing, setIsRefreshing] = useState(false);

  // Current values from URL (used as controlled values for the select,
  // and as defaultValue for the uncontrolled search input).
  const currentQuery = searchParams.get("query") ?? "";
  const currentSize  = searchParams.get("size") ?? "10";

  /** Update URL search params, preserving unrelated keys. */
  const updateParams = useCallback(
    (updates: Record<string, string | null>) => {
      const params = new URLSearchParams(searchParams.toString());
      for (const [key, value] of Object.entries(updates)) {
        if (value === null) params.delete(key);
        else params.set(key, value);
      }
      startTransition(() => {
        replace(`${pathname}?${params.toString()}`);
      });
    },
    [searchParams, pathname, replace],
  );

  /** Debounced search — fires 300 ms after the user stops typing. */
  const handleSearch = useCallback(
    (value: string) => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
      debounceRef.current = setTimeout(() => {
        updateParams({
          query: value.trim() || null, // remove param when empty
          page: "1",                   // reset to page 1 on new query
        });
      }, 300);
    },
    [updateParams],
  );

  /** Page-size selector — takes effect immediately. */
  const handlePageSize = useCallback(
    (size: string) => {
      updateParams({ size, page: "1" });
    },
    [updateParams],
  );

  /** Refresh button — re-fetches server data without changing the URL. */
  const handleRefresh = useCallback(() => {
    setIsRefreshing(true);
    refresh();
    // The spinning animation runs for 700 ms regardless so it looks intentional.
    setTimeout(() => setIsRefreshing(false), 700);
  }, [refresh]);

  return (
    <div className="flex items-center gap-3">
      {/* ── Search input ── */}
      <div className="relative flex-1">
        <svg
          className="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-[var(--text-muted)]"
          width="15" height="15" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}
          aria-hidden
        >
          <circle cx="11" cy="11" r="8"/>
          <path strokeLinecap="round" d="m21 21-4.35-4.35"/>
        </svg>
        <input
          id="user-search"
          type="search"
          // defaultValue keeps the input uncontrolled so typing is not
          // interrupted by re-renders from URL changes (Next.js tutorial pattern).
          defaultValue={currentQuery}
          onChange={(e) => handleSearch(e.target.value)}
          placeholder="Search by name or email…"
          aria-label="Search users"
          className={`w-full rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] py-2.5 pl-10 pr-4 text-sm text-[var(--text-primary)] placeholder:text-[var(--text-muted)] transition focus:border-[var(--brand-500)] focus:outline-none focus:ring-2 focus:ring-[var(--brand-200)] ${isPending ? "opacity-60" : ""}`}
        />
        {/* Pending indicator — appears while the Server Component re-renders */}
        {isPending && (
          <span className="absolute right-3 top-1/2 -translate-y-1/2">
            <svg className="h-4 w-4 animate-spin text-[var(--brand-400)]" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4l3-3-3-3v4a8 8 0 00-8 8h4z"/>
            </svg>
          </span>
        )}
      </div>

   

      {/* ── Refresh button ── */}
      <button
        type="button"
        onClick={handleRefresh}
        title="Refresh user list"
        aria-label="Refresh user list"
        className="flex h-10 w-10 items-center justify-center rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] text-[var(--text-secondary)] transition hover:bg-[var(--bg-subtle)] hover:text-[var(--brand-500)] focus:outline-none focus:ring-2 focus:ring-[var(--brand-200)]"
      >
        <svg
          width="15" height="15" fill="none" viewBox="0 0 24 24"
          stroke="currentColor" strokeWidth={2}
          className={isRefreshing ? "animate-spin" : ""}
          aria-hidden
        >
          <path strokeLinecap="round" strokeLinejoin="round" d="M4 4v6h6M20 20v-6h-6M4.93 15A9 9 0 1019.07 9" />
        </svg>
      </button>
    </div>
  );
}
