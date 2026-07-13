/**
 * lib/constants.ts
 *
 * Shared constants used by both server and client components.
 * Keeping them here avoids importing from page.tsx (a route module)
 * which can cause circular dependency warnings in some bundler configurations.
 */

export const PAGE_SIZE_OPTIONS = [5, 10, 25] as const;
export type PageSize = (typeof PAGE_SIZE_OPTIONS)[number];
