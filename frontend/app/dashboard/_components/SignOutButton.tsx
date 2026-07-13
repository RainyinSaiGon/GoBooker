"use client";

/**
 * SignOutButton — Client Component
 *
 * Performs a clean client-side logout:
 *   1. Clear auth tokens from Zustand store (also clears localStorage via
 *      the persist middleware).
 *   2. Call queryClient.clear() so no user-scoped cached data leaks into
 *      the next session.
 *   3. Redirect to the sign-in page.
 *
 * NOTE: The backend currently has no /logout endpoint (no server-side token
 * revocation). This is tracked as a follow-up task. Client-side logout is
 * safe for now since access tokens are short-lived (~30 min).
 */

import { useRouter } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import { useAuthStore } from "@/lib/store/auth";
import { apiLogout } from "@/lib/api";

export default function SignOutButton() {
  const router      = useRouter();
  const queryClient = useQueryClient();
  const logout      = useAuthStore((s) => s.logout);

  const handleSignOut = async () => {
    try {
      await apiLogout();
    } catch (e) {
      console.error("Failed to sign out from backend:", e);
    }
    // 1. Clear Zustand auth state (+ localStorage via persist middleware).
    logout();
    // 2. Wipe all cached query data so the next user starts with a clean slate.
    queryClient.clear();
    // 3. Navigate back to the sign-in page.
    router.push("/signin");
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
