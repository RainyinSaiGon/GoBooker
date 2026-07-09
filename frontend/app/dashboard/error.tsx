"use client";

/**
 * Dashboard Error Boundary — app/dashboard/error.tsx
 *
 * Catches errors thrown within the /dashboard segment (layout + page).
 * Because it's segment-scoped, an error in the dashboard won't crash the
 * outer shell (sidebar, nav), which keeps rendering normally.
 *
 * Next.js docs: https://nextjs.org/docs/app/getting-started/error-handling
 */

import { useEffect } from "react";

export default function DashboardError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error("[DashboardError]", error);
  }, [error]);

  return (
    <div className="flex flex-col items-center justify-center gap-6 rounded-2xl border border-red-200 bg-red-50 py-20 text-center dark:border-red-900 dark:bg-red-950/20">
      {/* Icon */}
      <div className="flex h-14 w-14 items-center justify-center rounded-full bg-red-100 dark:bg-red-950/50">
        <svg
          className="h-7 w-7 text-red-500"
          fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.8}
        >
          <path strokeLinecap="round" strokeLinejoin="round"
            d="M12 9v3m0 3h.01M5.07 19H19a2 2 0 001.73-3L13.73 4a2 2 0 00-3.46 0L4.27 16A2 2 0 005.07 19z"
          />
        </svg>
      </div>

      {/* Copy */}
      <div className="space-y-1">
        <h2 className="text-lg font-bold text-red-700 dark:text-red-400">
          Dashboard failed to load
        </h2>
        <p className="text-sm text-red-600/80 dark:text-red-400/70">
          {error.message || "An unexpected error occurred."}
        </p>
        {error.digest && (
          <p className="font-mono text-[11px] text-red-400/60 dark:text-red-500/50">
            ID: {error.digest}
          </p>
        )}
      </div>

      {/* Retry — calls reset() which re-renders the segment without a full reload */}
      <button
        type="button"
        onClick={reset}
        className="rounded-xl border border-red-300 px-5 py-2 text-sm font-semibold text-red-600 transition hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-red-400 focus:ring-offset-2 dark:border-red-800 dark:text-red-400 dark:hover:bg-red-950/40"
      >
        Try again
      </button>
    </div>
  );
}
