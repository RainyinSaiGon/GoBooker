/**
 * lib/queries/users.ts
 *
 * TanStack Query hooks and cache key factory for user-related data.
 * The `userKeys` factory ensures that every call site uses the same
 * cache key shape, so invalidation is predictable and type-safe.
 */

import {
  useQuery,
  useMutation,
  useQueryClient,
  queryOptions,
} from "@tanstack/react-query";
import {
  apiGetUsers,
  apiCreateUser,
  type GetUsersParams,
  type CreateUserPayload,
  type User,
} from "@/lib/api";

// ─── Query key factory ────────────────────────────────────────────────────────

export const userKeys = {
  /** Matches ALL user queries — use for broad invalidations. */
  all:    ["users"] as const,
  /** Matches paginated lists with specific filters. */
  list:   (params: GetUsersParams) => ["users", "list", params] as const,
};

// ─── queryOptions (reusable in server prefetch + client hooks) ────────────────

export function usersQueryOptions(params: GetUsersParams) {
  return queryOptions({
    queryKey: userKeys.list(params),
    queryFn:  () => apiGetUsers(params),
    staleTime: 30_000, // treat data as fresh for 30 s before background refetch
  });
}

// ─── Hooks ───────────────────────────────────────────────────────────────────

/**
 * useUsers — fetches a paginated, filtered user list.
 *
 * On first render after SSR, TanStack Query will find the dehydrated cache
 * from the server and return `data` immediately (no loading flash).
 * After `staleTime`, it will silently re-fetch in the background.
 */
export function useUsers(params: GetUsersParams) {
  return useQuery({
    ...usersQueryOptions(params),
    // Keep previous data visible while fetching the next page / new query —
    // this prevents the list from flickering to empty between navigations.
    placeholderData: (prev) => prev,
  });
}

/**
 * useCreateUser — mutation for the sign-up form.
 *
 * On success, invalidates ALL user queries so the dashboard list refreshes
 * automatically when the user navigates there after signing up.
 */
export function useCreateUser() {
  const queryClient = useQueryClient();
  return useMutation<User, Error, CreateUserPayload>({
    mutationFn: apiCreateUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.all });
    },
  });
}
