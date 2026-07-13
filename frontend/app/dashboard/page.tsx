/**
 * app/dashboard/page.tsx — Server Component
 *
 * Responsibility: prefetch the initial user list on the server so the client
 * gets a cache hit on first render (no loading spinner flash).
 *
 * Pattern (TanStack Query Advanced SSR):
 *   1. Create a server-side QueryClient.
 *   2. Call queryClient.prefetchQuery() with the same key + fn used by the
 *      client hook (via usersQueryOptions).
 *   3. Dehydrate the cache and pass it to HydrationBoundary.
 *   4. The client DashboardClient calls useUsers() → instant cache hit.
 */

import {
  HydrationBoundary,
  QueryClient,
  dehydrate,
} from "@tanstack/react-query";
import { usersQueryOptions } from "@/lib/queries/users";
import { PAGE_SIZE_OPTIONS, type PageSize } from "@/lib/constants";
import DashboardClient from "./_components/DashboardClient";

export const metadata = { title: "Users" };

export default async function DashboardPage({
  searchParams,
}: {
  searchParams: Promise<{ query?: string; page?: string; size?: string }>;
}) {
  const params = await searchParams;

  const query    = params.query ?? "";
  const page     = Math.max(1, Number(params.page) || 1);
  const rawSize  = Number(params.size);
  const pageSize: PageSize = PAGE_SIZE_OPTIONS.includes(rawSize as PageSize)
    ? (rawSize as PageSize)
    : 10;

  const queryClient = new QueryClient();
  await queryClient.prefetchQuery(usersQueryOptions({ query, page, pageSize }));

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <DashboardClient />
    </HydrationBoundary>
  );
}
