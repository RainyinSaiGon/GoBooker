/**
 * loading.tsx — Next.js Instant Loading UI
 *
 * Next.js automatically shows this file while the dashboard Server Component
 * is fetching data (awaiting searchParams + the fetchUsers call).
 * It renders instantly from the server with no JavaScript needed.
 *
 * Docs: https://nextjs.org/docs/app/api-reference/file-conventions/loading
 */

function SkeletonRow() {
  return (
    <div className="flex items-center gap-4 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] px-4 py-3.5">
      <div className="skeleton h-10 w-10 rounded-full shrink-0" />
      <div className="flex-1 space-y-2">
        <div className="skeleton h-3.5 w-36 rounded" />
        <div className="skeleton h-3 w-52 rounded" />
      </div>
      <div className="skeleton h-5 w-16 rounded-full" />
    </div>
  );
}

export default function DashboardLoading() {
  return (
    <div className="mx-auto max-w-3xl space-y-8">
      {/* Header skeleton */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="space-y-2">
          <div className="skeleton h-7 w-24 rounded" />
          <div className="skeleton h-4 w-40 rounded" />
        </div>
        <div className="skeleton h-9 w-28 rounded-xl" />
      </div>

      {/* Search bar skeleton */}
      <div className="flex items-center gap-3">
        <div className="skeleton h-10 flex-1 rounded-xl" />
        <div className="skeleton h-10 w-24 rounded-xl" />
        <div className="skeleton h-10 w-10 rounded-xl" />
      </div>

      {/* Row skeletons */}
      <div className="space-y-2.5">
        {[...Array(5)].map((_, i) => (
          <SkeletonRow key={i} />
        ))}
      </div>

      {/* Pagination skeleton */}
      <div className="flex justify-between">
        <div className="skeleton h-4 w-32 rounded" />
        <div className="flex gap-1.5">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="skeleton h-8 w-8 rounded-lg" />
          ))}
        </div>
      </div>
    </div>
  );
}
