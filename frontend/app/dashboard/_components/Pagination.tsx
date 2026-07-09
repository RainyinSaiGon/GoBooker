"use client";

/**
 * Pagination — Client Component
 *
 * Implements the Next.js recommended pagination pattern
 *
 *   • Use useSearchParams() + usePathname() to build page URLs that preserve
 *     ALL existing search params (query, size, etc.) — only the `page` key
 *     is updated.
 *   • Use Next.js <Link> components instead of onClick handlers — this gives
 *     free prefetching on hover and correct browser Back/Forward behaviour.
 *   • The `aria-current="page"` attribute on the active page button satisfies
 *     accessibility requirements.
 */

import Link from "next/link";
import { useSearchParams, usePathname } from "next/navigation";

interface PaginationProps {
  page: number;
  totalPages: number;
  pageSize: number;
  from: number;
  to: number;
  total: number;
}

export default function Pagination({
  page,
  totalPages,
  from,
  to,
  total,
}: PaginationProps) {
  const searchParams = useSearchParams();
  const pathname     = usePathname();

  /** Build a URL that keeps ALL current params but changes only `page`. */
  const createPageURL = (pageNumber: number): string => {
    const params = new URLSearchParams(searchParams.toString());
    params.set("page", pageNumber.toString());
    return `${pathname}?${params.toString()}`;
  };

  if (totalPages <= 1) {
    return (
      <p className="text-center text-xs text-[var(--text-muted)]">
        Showing {from}–{to} of {total}
      </p>
    );
  }

  // ── Page-number range ─────────────────────────────────────────────────────
  // Show at most 5 page buttons, centred on the current page.
  const WINDOW = 2;
  const start  = Math.max(1, Math.min(page - WINDOW, totalPages - WINDOW * 2));
  const end    = Math.min(totalPages, start + WINDOW * 2);
  const pages  = Array.from({ length: end - start + 1 }, (_, i) => start + i);

  // ── Shared class helpers ──────────────────────────────────────────────────
  const base =
    "flex h-8 min-w-[2rem] items-center justify-center rounded-lg px-1 text-sm font-medium transition-all focus:outline-none focus-visible:ring-2 focus-visible:ring-[var(--brand-300)]";
  const activeClass =
    `${base} bg-[var(--brand-500)] text-white shadow-sm pointer-events-none`;
  const inactiveClass =
    `${base} border border-[var(--border)] bg-[var(--bg-surface)] text-[var(--text-secondary)] hover:bg-[var(--bg-subtle)] hover:text-[var(--brand-500)]`;
  const arrowClass =
    `${base} border border-[var(--border)] bg-[var(--bg-surface)] text-[var(--text-secondary)] hover:bg-[var(--bg-subtle)] hover:text-[var(--brand-500)] aria-disabled:pointer-events-none aria-disabled:opacity-40`;

  return (
    <div className="flex flex-col items-center gap-3 pt-1 sm:flex-row sm:justify-between">
      {/* Item counter */}
      <p className="text-xs text-[var(--text-muted)]">
        Showing <span className="font-medium text-[var(--text-secondary)]">{from}–{to}</span> of{" "}
        <span className="font-medium text-[var(--text-secondary)]">{total}</span>
      </p>

      {/* Page buttons */}
      <nav aria-label="Pagination" className="flex items-center gap-1">

        {/* ← Prev */}
        <Link
          href={createPageURL(page - 1)}
          aria-label="Go to previous page"
          aria-disabled={page === 1}
          className={arrowClass}
          tabIndex={page === 1 ? -1 : undefined}
        >
          <svg width="14" height="14" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.2} aria-hidden>
            <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
          </svg>
        </Link>

        {/* First page + ellipsis */}
        {start > 1 && (
          <>
            <Link href={createPageURL(1)} className={inactiveClass} aria-label="Go to page 1">1</Link>
            {start > 2 && (
              <span className="px-0.5 text-sm text-[var(--text-muted)]" aria-hidden>…</span>
            )}
          </>
        )}

        {/* Page window */}
        {pages.map((p) => (
          <Link
            key={p}
            href={createPageURL(p)}
            aria-label={`Go to page ${p}`}
            aria-current={p === page ? "page" : undefined}
            className={p === page ? activeClass : inactiveClass}
          >
            {p}
          </Link>
        ))}

        {/* Ellipsis + last page */}
        {end < totalPages && (
          <>
            {end < totalPages - 1 && (
              <span className="px-0.5 text-sm text-[var(--text-muted)]" aria-hidden>…</span>
            )}
            <Link href={createPageURL(totalPages)} className={inactiveClass} aria-label={`Go to page ${totalPages}`}>
              {totalPages}
            </Link>
          </>
        )}

        {/* Next → */}
        <Link
          href={createPageURL(page + 1)}
          aria-label="Go to next page"
          aria-disabled={page === totalPages}
          className={arrowClass}
          tabIndex={page === totalPages ? -1 : undefined}
        >
          <svg width="14" height="14" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.2} aria-hidden>
            <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
          </svg>
        </Link>
      </nav>
    </div>
  );
}
