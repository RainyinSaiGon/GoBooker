"use client";

/**
 * app/providers.tsx
 *
 * A single "use client" boundary that sets up all React context providers
 * needed by the app. Keeping this in one file means layout.tsx stays a pure
 * Server Component (no "use client" directive required there).
 *
 * QueryClient config:
 *   - staleTime: 60s  — data is considered fresh for 1 minute before a
 *                        background refetch is triggered. Most pages are
 *                        unlikely to need sub-second freshness.
 *   - gcTime:    5min — unused cache entries are garbage-collected after 5 min.
 */

import { useState } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

function makeQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60 * 1000,
        gcTime:    5 * 60 * 1000,
        retry:     1,
        refetchOnWindowFocus: true,
      },
    },
  });
}

let browserQueryClient: QueryClient | undefined;

/**
 * Returns a stable QueryClient.
 *
 * On the server we always create a new client (each request is independent).
 * In the browser we create one once and reuse it across re-renders so the
 * cache survives navigation.
 */
function getQueryClient() {
  if (typeof window === "undefined") {
    return makeQueryClient();
  }
  if (!browserQueryClient) {
    browserQueryClient = makeQueryClient();
  }
  return browserQueryClient;
}

export default function Providers({ children }: { children: React.ReactNode }) {
  // useState ensures the QueryClient is not re-created on every render.
  // See: https://tanstack.com/query/latest/docs/framework/react/guides/advanced-ssr
  const [queryClient] = useState(() => getQueryClient());

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}
