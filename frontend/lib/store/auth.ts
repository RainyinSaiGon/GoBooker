/**
 * lib/store/auth.ts
 *
 * Zustand auth store with localStorage persistence via the `persist`
 * middleware. This holds the JWT access token and refresh token returned
 * by POST /api/v1/auth/login.
 *
 * Rule of thumb (from the architectural decision):
 *   - Tokens + auth status  → this store (persisted)
 *   - Search/page/pageSize  → URL search params (shareable)
 *   - UI-only state         → local component state or a separate store
 */

import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AuthState {
  /** JWT access token (short-lived, ~30 min per backend). */
  token:        string | null;
  /** JWT refresh token (long-lived, ~7 days per backend). */
  refreshToken: string | null;
  /** Derived: true when a non-null token is stored. */
  isAuthenticated: boolean;

  /** Store both tokens after a successful login. */
  setTokens: (token: string, refreshToken: string) => void;

  /**
   * Clear all auth state.
   * Callers should also call queryClient.clear() to prevent stale data
   * from leaking into the next session.
   */
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token:           null,
      refreshToken:    null,
      isAuthenticated: false,

      setTokens: (token, refreshToken) =>
        set({ token, refreshToken, isAuthenticated: true }),

      logout: () =>
        set({ token: null, refreshToken: null, isAuthenticated: false }),
    }),
    {
      name:    "gobooker-auth",
      // Only persist what is needed — not the action functions.
      partialize: (state) => ({
        token:        state.token,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
