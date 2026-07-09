"use client";

/**
 * SignOutButton — Client Component leaf
 *
 * Extracted to satisfy the Server Components rules:
 *   - DashboardLayout is a Server Component (no "use client").
 *   - Event handlers (onClick) are only allowed in Client Components.
 *   - By keeping DashboardLayout as a Server Component, we avoid converting the
 *     entire layout and its subtree to client-rendering.
 */
export default function SignOutButton() {
  const handleSignOut = () => {
    // TODO: wire to real auth / signOut call
    console.log("Signing out...");
  };

  return (
    <button
      type="button"
      onClick={handleSignOut}
      className="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-[var(--text-secondary)] transition-all hover:bg-red-50 hover:text-red-500 dark:hover:bg-red-950/30 dark:hover:text-red-400 group"
    >
      <svg
        className="text-[var(--text-muted)] group-hover:text-red-500 transition-colors"
        width="16"
        height="16"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={1.8}
        aria-hidden
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H5a3 3 0 01-3-3V7a3 3 0 013-3h5a3 3 0 013 3v1"
        />
      </svg>
      Sign Out
    </button>
  );
}
