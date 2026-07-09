"use client";

/**
 * Global Error Boundary — app/error.tsx
 *
 * Next.js docs: https://nextjs.org/docs/app/getting-started/error-handling
 *
 * Rules:
 *   • MUST be a Client Component ("use client") — it uses React's error boundary.
 *   • Receives `error` (the thrown Error) and `reset` (retries the segment).
 *   • The `error.digest` property is a hash for server-side log correlation.
 *   • This file catches errors in app/page.tsx (the sign-up page).
 *
 * A separate app/dashboard/error.tsx handles dashboard segment errors.
 */

import { useEffect } from "react";
import Link from "next/link";

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    // In production, pipe this to your error-reporting service (e.g. Sentry).
    console.error("[GlobalError]", error);
  }, [error]);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-6 p-8 text-center">
      {/* Icon */}
      <div className="flex h-16 w-16 items-center justify-center rounded-full bg-red-100 dark:bg-red-950/50">
        <svg
          className="h-8 w-8 text-red-500"
          fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.8}
        >
          <path strokeLinecap="round" strokeLinejoin="round"
            d="M12 9v3m0 3h.01M5.07 19H19a2 2 0 001.73-3L13.73 4a2 2 0 00-3.46 0L4.27 16A2 2 0 005.07 19z"
          />
        </svg>
      </div>

      {/* Copy */}
      <div className="space-y-1">
        <h1 className="text-xl font-bold text-[var(--text-primary)]">
          Something went wrong
        </h1>
        <p className="text-sm text-[var(--text-secondary)]">
          An unexpected error occurred. Our team has been notified.
        </p>
        {error.digest && (
          <p className="mt-1 font-mono text-[11px] text-[var(--text-muted)]">
            Error ID: {error.digest}
          </p>
        )}
      </div>

      {/* Actions */}
      <div className="flex gap-3">
        <button
          type="button"
          onClick={reset}
          className="rounded-xl bg-[var(--brand-500)] px-5 py-2 text-sm font-semibold text-white transition hover:bg-[var(--brand-600)] focus:outline-none focus:ring-2 focus:ring-[var(--brand-400)] focus:ring-offset-2"
        >
          Try again
        </button>
        <Link
          href="/"
          className="rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] px-5 py-2 text-sm font-medium text-[var(--text-secondary)] transition hover:bg-[var(--bg-subtle)]"
        >
          Go home
        </Link>
      </div>
    </div>
  );
}
