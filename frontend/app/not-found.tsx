/**
 * Not Found page — app/not-found.tsx
 *
 * Next.js docs: https://nextjs.org/docs/app/api-reference/file-conventions/not-found
 *
 * This file is rendered when:
 *   1. A route segment calls `notFound()` from "next/navigation".
 *   2. A URL is requested that matches no route segment.
 *
 * It is a Server Component by default (no "use client" needed).
 */

import Link from "next/link";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Page Not Found",
  description: "The page you are looking for does not exist.",
  robots: { index: false, follow: false },
};

export default function NotFound() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-8 p-8 text-center">
      {/* Large 404 */}
      <div className="select-none">
        <p className="gradient-text text-[8rem] font-black leading-none tracking-tighter">
          404
        </p>
        <div
          aria-hidden
          className="pointer-events-none absolute -translate-y-32 opacity-10 blur-3xl"
          style={{
            background: "radial-gradient(circle, var(--brand-400) 0%, transparent 70%)",
            width: "400px",
            height: "400px",
            left: "50%",
            transform: "translateX(-50%) translateY(-200px)",
          }}
        />
      </div>

      {/* Copy */}
      <div className="max-w-md space-y-2">
        <h1 className="text-2xl font-bold text-[var(--text-primary)]">
          Page not found
        </h1>
        <p className="text-[var(--text-secondary)]">
          Sorry, the page you are looking for doesn&apos;t exist or has been
          moved.
        </p>
      </div>

      {/* Actions */}
      <div className="flex flex-wrap items-center justify-center gap-3">
        <Link
          href="/"
          className="rounded-xl bg-[var(--brand-500)] px-6 py-2.5 text-sm font-semibold text-white shadow-md transition hover:bg-[var(--brand-600)] active:scale-[.98] focus:outline-none focus:ring-2 focus:ring-[var(--brand-400)] focus:ring-offset-2"
        >
          Go home
        </Link>
        <Link
          href="/dashboard"
          className="rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] px-6 py-2.5 text-sm font-medium text-[var(--text-secondary)] transition hover:bg-[var(--bg-subtle)] focus:outline-none focus:ring-2 focus:ring-[var(--brand-200)] focus:ring-offset-2"
        >
          Dashboard
        </Link>
      </div>
    </div>
  );
}
