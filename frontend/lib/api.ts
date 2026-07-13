/**
 * Shared API base URL used by every fetch call in the app.
 *
 * Override via NEXT_PUBLIC_API_URL in your .env.local for staging / prod:
 *   NEXT_PUBLIC_API_URL=https://api.mybooker.com
 *
 * FQ2: extracted from the duplicate const in page.tsx and dashboard/page.tsx.
 */
export const API_BASE: string =
  process.env.NEXT_PUBLIC_API_URL
    ? `${process.env.NEXT_PUBLIC_API_URL}/api/v1`
    : "http://localhost:8081/api/v1";
