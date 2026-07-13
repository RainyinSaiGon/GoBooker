"use client";

/**
 * SearchBar — Client Component
 *
 * Search and page-size controls. State is stored in URL params (shareable,
 * bookmarkable) following the Next.js tutorial pattern. The key change from
 * the original is the refresh button: instead of `router.refresh()` (which
 * triggers a full Server Component re-render), we call
 * `queryClient.invalidateQueries(userKeys.all)` so TanStack Query does a
 * background refetch and updates only what changed.
 */

import { useSearchParams, usePathname, useRouter } from "next/navigation";
import { useRef, useCallback, useTransition, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { PAGE_SIZE_OPTIONS, type PageSize } from "@/lib/constants";
import { userKeys } from "@/lib/queries/users";

export default function SearchBar() {
  const searchParams          = useSearchParams();
  const pathname              = usePathname();
  const { replace }           = useRouter();
  const queryClient           = useQueryClient();
  const debounceRef           = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [isPending, startTransition] = useTransition();
  const [isRefreshing, setIsRefreshing] = useState(false);

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
          query: value.trim() || null,
          page:  "1",
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

  /**
   * Refresh button — invalidates the user query so TanStack Query does a
   * background refetch. This is cheaper than router.refresh() because it
   * only re-fetches the affected data, not the whole Server Component tree.
   */
  const handleRefresh = useCallback(() => {
    setIsRefreshing(true);
    queryClient.invalidateQueries({ queryKey: userKeys.all });
    setTimeout(() => setIsRefreshing(false), 700);
  }, [queryClient]);

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
          defaultValue={currentQuery}
          onChange={(e) => handleSearch(e.target.value)}
          placeholder="Search by name or email…"
          aria-label="Search users"
          className={`w-full rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] py-2.5 pl-10 pr-4 text-sm text-[var(--text-primary)] placeholder:text-[var(--text-muted)] transition focus:border-[var(--brand-500)] focus:outline-none focus:ring-2 focus:ring-[var(--brand-200)] ${isPending ? "opacity-60" : ""}`}
        />
        {isPending && (
          <span className="absolute right-3 top-1/2 -translate-y-1/2">
            <svg className="h-4 w-4 animate-spin text-[var(--brand-400)]" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4l3-3-3-3v4a8 8 0 00-8 8h4z"/>
            </svg>
          </span>
        )}
      </div>

      {/* ── Page-size selector ── */}
      <select
        id="page-size"
        value={currentSize}
        onChange={(e) => handlePageSize(e.target.value)}
        aria-label="Rows per page"
        className="h-10 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] px-3 text-sm text-[var(--text-secondary)] transition focus:border-[var(--brand-500)] focus:outline-none focus:ring-2 focus:ring-[var(--brand-200)]"
      >
        {PAGE_SIZE_OPTIONS.map((s) => (
          <option key={s} value={s}>{s} / page</option>
        ))}
      </select>

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
