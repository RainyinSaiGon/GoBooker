# GoBooker — Design Trade-offs & Decision Rationale

> This document accompanies the main Technical Report and explains the **alternatives considered** for each major design decision, and why the current approach was chosen.

---

## Table of Contents

1. [Backend: Pagination Strategy](#1-backend-pagination-strategy)
2. [Backend: ORM vs Raw SQL](#2-backend-orm-vs-raw-sql)
3. [Backend: HTTP Router Choice](#3-backend-http-router-choice)
4. [Backend: Password Hashing](#4-backend-password-hashing)
5. [Backend: Error Handling Architecture](#5-backend-error-handling-architecture)
6. [Backend: Sentinel Errors vs Wrapped Types](#6-backend-sentinel-errors-vs-wrapped-types)
7. [Frontend: Next.js Rendering Strategy](#7-frontend-nextjs-rendering-strategy)
8. [Frontend: URL State vs React useState](#8-frontend-url-state-vs-react-usestate)
9. [Frontend: Server vs Client Components](#9-frontend-server-vs-client-components)
10. [Frontend: Pagination UI Pattern](#10-frontend-pagination-ui-pattern)
11. [Frontend: Search Debounce Implementation](#11-frontend-search-debounce-implementation)
12. [Frontend: Error Boundary Scope](#12-frontend-error-boundary-scope)

---

## 1. Backend: Pagination Strategy

### The Problem
The database can have thousands of users. Returning all rows on every request is wasteful, slow, and a DoS risk.

### Alternatives Considered

| Approach | How It Works | Pros | Cons |
|----------|-------------|------|------|
| **In-memory (what we replaced)** | Fetch all rows, slice in Go/JS | Simple code | O(n) memory; slow on large tables; DB always under full load |
| **Cursor-based pagination** | `WHERE id > $last_seen_id LIMIT n` | Stable results; no skipped rows on inserts | Can't jump to page 5; no "total count"; complex UI |
| **Keyset pagination** | `WHERE (created_at, id) < ($last_ts, $last_id)` | Very fast on large indexed tables | Same as cursor — no random page access |
| **✅ LIMIT/OFFSET (chosen)** | `LIMIT $size OFFSET $(page-1)*size` | Easy to implement; supports random page jumps; compatible with any UI | Slight drift if rows are inserted between pages; slower on very high offsets (>100k) |

### Why LIMIT/OFFSET

For GoBooker's use case (user management with typically hundreds to low thousands of users), LIMIT/OFFSET is the **pragmatic choice**:
- Allows a user to jump directly to page 7 — impossible with cursor-based.
- The UI can display "Page 3 of 12" — requires a `COUNT(*)` which cursor pagination cannot easily provide.
- The "offset drift" problem (where inserting a row mid-query shifts results by one) is acceptable in an admin dashboard context — users are not reading a real-time feed.

**When to change**: If the user table grows to >100,000 rows and users frequently navigate to high page numbers (e.g., page 500), the `OFFSET 5000` scan becomes slow. At that point, cursor-based pagination is the correct migration.

---

## 2. Backend: ORM vs Raw SQL

### The Problem
Interacting with the database requires writing either raw SQL strings or using an abstraction layer.

### Alternatives Considered

| Approach | Examples | Pros | Cons |
|----------|---------|------|------|
| **Full ORM** | GORM, ent | Auto-migrations; less SQL knowledge needed; model-driven | Magic behaviour; hard to debug; N+1 query problems; poor control over complex queries |
| **Query builder** | sqlx, squirrel | Type-safe query building; composable | Extra dependency; still generates SQL dynamically |
| **✅ Raw SQL + database/sql (chosen)** | Standard library | Full control; readable queries; predictable performance; no magic | More boilerplate; SQL strings aren't type-checked at compile time |

### Why Raw SQL

GoBooker has a small, well-defined schema. Writing raw SQL means:

1. **The `ILIKE` pagination query is explicit** — an ORM would hide the `COUNT(*)` sub-query, making it hard to know if two round-trips are happening.
2. **The conditional `SET` for password updates is transparent** — ORMs often make "update only dirty fields" surprisingly complex.
3. **Performance is predictable** — no hidden JOINs, no lazy loading surprises.

The boilerplate cost (writing `rows.Scan(...)` manually) is low given the single `users` table.

**When to change**: With 10+ related tables and complex relational queries, `sqlx` (which adds struct scanning) or `ent` (which adds type-safe query generation) would reduce repetition without the downsides of a full ORM.

---

## 3. Backend: HTTP Router Choice

### The Problem
Go's standard library `http.ServeMux` is minimal. It doesn't support path parameters like `/users/{id}`.

### Alternatives Considered

| Approach | Examples | Pros | Cons |
|----------|---------|------|------|
| **stdlib ServeMux (Go 1.22+)** | `net/http` | Zero dependency; Go 1.22 added `{id}` patterns | Slightly verbose; no method-based routing before 1.22 |
| **Gin** | `github.com/gin-gonic/gin` | Very fast; middleware ecosystem; request binding | Large dependency; opinionated; harder to test handlers in isolation |
| **Chi** | `github.com/go-chi/chi` | Lightweight; idiomatic; stdlib-compatible | Less popular; smaller ecosystem |
| **✅ gorilla/mux (chosen)** | `github.com/gorilla/mux` | Battle-tested; `{id}` path vars via `mux.Vars(r)`; method-based routing | Gorilla is no longer actively maintained (though it's stable) |

### Why gorilla/mux

The project was already using gorilla/mux when we started. The key question was whether to migrate.

- **Gin** would have been faster (lower latency at high concurrency) but introduced opinionated response helpers that would require rewriting all handlers.
- **stdlib (Go 1.22+)** is a viable future migration since Go's project uses Go 1.26 — but the path variable extraction (`{id}`) syntax changed slightly between gorilla and stdlib patterns.
- gorilla/mux is **stable and widely deployed**. Its maintenance freeze means no new features, not instability.

**When to change**: If performance benchmarking shows HTTP routing as a bottleneck at high load, migrating to `net/http` (Go 1.22+ patterns) removes the external dependency entirely.

---

## 4. Backend: Password Hashing

### The Problem
User passwords must be stored securely — a compromised database should not expose plaintext or easily reversible hashes.

### Alternatives Considered

| Approach | Algorithm | Pros | Cons |
|----------|---------|------|------|
| **MD5 / SHA-256** | Fast hash | Extremely fast | Not suitable for passwords: no salt, GPU-crackable in hours |
| **PBKDF2** | Iterated hash | Tunable; FIPS-approved | Less common in Go; requires manual salt management |
| **Argon2** | Memory-hard | Strongest modern algorithm; winner of Password Hashing Competition | Not in Go stdlib; requires `golang.org/x/crypto/argon2` |
| **✅ bcrypt (chosen)** | Adaptive cost | Salt is embedded in output; widely supported; `golang.org/x/crypto/bcrypt` in stdlib family | Slower at very high registration throughput; capped at 72 bytes |

### Why bcrypt

- The **salt is embedded in the output hash** — no separate `password_salt` column needed.
- `bcrypt.DefaultCost` (10) is a good balance: ~100ms per hash, which is imperceptible to humans but expensive for an attacker trying millions of guesses.
- It is the **industry standard** for user passwords in most web applications built in the 2020s.
- The 72-byte input limit is not a real concern for typical user passwords.

**When to change**: Argon2id is the more modern recommendation (PHC winner, memory-hard). If security requirements are elevated, the service layer is the only place that needs to change — the repository and handler are unaffected because the hashing is fully encapsulated in `service.CreateUser`.

---

## 5. Backend: Error Handling Architecture

### The Problem
Errors occur at the database level but the HTTP handler must return the right status code. There are several ways to communicate error types across layers.

### Alternatives Considered

| Approach | How It Works | Pros | Cons |
|----------|-------------|------|------|
| **Return error strings** | `errors.New("duplicate email")`, check with `strings.Contains` | Simple | Fragile — any message change breaks callers; no compile-time checking |
| **HTTP status in repository** | Repository returns `(result, httpStatus, error)` | Single place | Violates separation of concerns — the DB layer should not know about HTTP |
| **Custom error types** | `type DuplicateEmailError struct{}; func (e *DuplicateEmailError) Error() string` | Full metadata | More boilerplate; `errors.As` instead of `errors.Is` |
| **✅ Sentinel errors (chosen)** | `var ErrDuplicateEmail = errors.New("email already in use")` | Simple; type-safe via `errors.Is`; idiomatic Go | Not suitable for carrying dynamic context (e.g. which field was duplicated) |

### Why Sentinel Errors

Sentinel errors (`var ErrXxx = errors.New(...)`) are the **idiomatic Go pattern** for well-known failure conditions:

```go
// Handler — clean comparison
if errors.Is(err, service.ErrDuplicateEmail) {
    writeError(w, http.StatusConflict, "email already taken")
}
```

This approach keeps each layer's concern separate:
- Repository: detects `SQLSTATE 23505` and returns `ErrDuplicateEmail`
- Service: wraps or re-returns it
- Handler: maps it to `409 Conflict`

**When to change**: If errors need to carry context (e.g., "field `email` already exists for user ID `abc-123`"), a custom error type with `errors.As` would be more appropriate. For an admin tool this is unnecessary.

---

## 6. Backend: Sentinel Errors vs Wrapped Types

### The `isDuplicateKey` Helper — Why Not a Typed Error from pgx?

pgx/v5 returns `*pgconn.PgError` which has a `Code` field (SQLSTATE). We could have written:

```go
// Alternative: type assertion on pgx error
var pgErr *pgconn.PgError
if errors.As(err, &pgErr) && pgErr.Code == "23505" { ... }
```

**Why we didn't**:
- Importing `pgconn` in the service or repository ties the error detection to the specific driver. If the driver is swapped, all detection breaks.
- The `isDuplicateKey(err)` helper checks `err.Error()` string for common patterns from both PostgreSQL and CockroachDB — making it **database-agnostic**.

**Trade-off**: String-matching on error messages is theoretically fragile if the DB vendor changes its error format. The pragmatic offset is that both PostgreSQL and CockroachDB have had stable `"duplicate key"` messages for years.

---

## 7. Frontend: Next.js Rendering Strategy

### The Problem
How should the dashboard page fetch and display users?

### Alternatives Considered

| Approach | How It Works | Pros | Cons |
|----------|-------------|------|------|
| **Client-side only (SPA)** | `useEffect → fetch → setState` | Simple to write; no server config | Flash of empty content; poor SEO; double render (mount + fetch) |
| **`getServerSideProps` (Pages Router)** | Run server code per request in Pages Router | Data on first HTML | Requires Pages Router; different API from App Router |
| **Static Generation (`getStaticProps`)** | Generate HTML at build time | Fastest; CDN cacheable | User list is dynamic — would be stale immediately |
| **✅ App Router Server Component (chosen)** | Async component with `await searchParams` | Zero client JS for data; first HTML contains data; built-in caching via `fetch` | Slightly harder mental model; requires understanding Server vs Client split |

### Why App Router Server Component

The Next.js 16 App Router is the **current official recommendation**. The key benefits for GoBooker:

1. **No loading flash**: The first HTML sent to the browser already contains the user list. With `useEffect`, the user sees a blank table for ~200ms.
2. **URL as state**: `searchParams` is a server-side prop — the page naturally re-renders whenever the URL changes, making search and pagination work without any client-side fetching logic.
3. **`cache: "no-store"`**: Each navigation fetches fresh data, equivalent to `getServerSideProps` behaviour.

**Trade-off accepted**: The server must be running and reachable for every page render. If the backend is down, the dashboard shows an error state (handled by `error.tsx`). With a fully static/client-side approach, the shell would still render.

---

## 8. Frontend: URL State vs React `useState`

### The Problem
The search query, current page, and page size need to be stored somewhere.

### Alternatives Considered

| Approach | Storage | Pros | Cons |
|----------|---------|------|------|
| **`useState`** | React memory | Simple; fast updates | Lost on refresh; can't be bookmarked; Back button doesn't work |
| **Global state (Zustand, Redux)** | Client JS store | Persists across component tree | More dependencies; still lost on refresh; overkill for URL params |
| **`sessionStorage` / `localStorage`** | Browser storage | Survives refresh | Not part of URL; can't share link; not SSR-compatible |
| **✅ URL search params (chosen)** | Browser URL | Bookmarkable; shareable; works with Back/Forward; SSR-compatible; free | Slightly more complex (`useSearchParams`, `useRouter.replace`) |

### Why URL Search Params

This is the **official Next.js recommendation** (Chapter 10: Adding Search and Pagination):

> "URL search params represent the URL's query string and are the single source of truth for the search and filter state."

Benefits:
- A user can bookmark `/dashboard?query=alice&page=2&size=25` and return to the exact same view.
- Sharing the URL with a colleague works out of the box.
- The browser's Back and Forward buttons naturally navigate between search states.
- The Server Component reads these directly from `searchParams` — no client-side state synchronisation needed.

**Trade-off accepted**: Every keystroke (after debounce) changes the URL, which triggers a server re-render. With `useState`, updates are instant (local memory). The 300ms debounce mitigates perceived latency, and `useTransition` keeps the UI responsive during the re-render.

---

## 9. Frontend: Server vs Client Components

### The Problem
Not all components need to run in the browser. Sending unnecessary JavaScript to the client wastes bandwidth and hurts performance (Time to Interactive).

### Alternatives Considered

| Approach | What it means | Pros | Cons |
|----------|-------------|------|------|
| **Everything `"use client"`** | Full SPA behavior; all JS in browser | Familiar React model | Larger bundle; no server rendering benefit; slower FCP |
| **Everything server-only** | No interactivity at all | Zero JS bundle | Search input, pagination clicks don't work |
| **✅ Server by default, Client at leaves (chosen)** | Server Components for data/layout; Client Components for interactivity | Minimum JS shipped; server renders data; best of both | Requires discipline to not accidentally mix |

### Why "Push Client to the Leaves"

The Next.js docs explicitly recommend this pattern:

> "Move Client Components to the leaves of your component tree."

In GoBooker, this means:
- `DashboardLayout` (the sidebar, nav) → **Server Component** — pure markup, no state
- `UserList` (the cards) → **Server Component** — receives data as props, no interactivity  
- `SearchBar` → **Client Component** — must use `useSearchParams`, `useRouter`
- `Pagination` → **Client Component** — must use `useSearchParams` to build URLs
- `SignOutButton` → **Client Component leaf** — isolated for `onClick` only

The critical bug we fixed: `DashboardLayout` originally had an inline `onClick` on the Sign Out button, which **forced the entire layout tree into Client rendering**. Extracting it to `SignOutButton.tsx` kept the layout as a Server Component.

**Trade-off accepted**: The component split adds more files. However, the benefit (smaller JS bundle, faster First Contentful Paint) outweighs the organisational cost.

---

## 10. Frontend: Pagination UI Pattern

### The Problem
How should the pagination component navigate between pages?

### Alternatives Considered

| Approach | Mechanism | Pros | Cons |
|----------|---------|------|------|
| **`onClick` + `useRouter.push()`** | Imperative JS navigation | Familiar React pattern | Adds JS event handlers; no free prefetch; pushes to history stack |
| **`onClick` + `useRouter.replace()`** | Replaces current history entry | Avoids history spam | Still imperative; no prefetch |
| **✅ `<Link>` components (chosen)** | Declarative Next.js Link | Free prefetch on hover; correct Back/Forward; no JS needed to render HTML | Requires building the full URL at render time |

### Why `<Link>` Components

The official Next.js pagination tutorial uses `<Link>` for a key reason:

> Next.js `<Link>` automatically **prefetches** the linked page's data when it enters the viewport. When the user clicks "Next Page", the data is already cached.

The `createPageURL` helper builds the complete URL at render time, preserving all existing params:
```tsx
const createPageURL = (pageNumber: number) => {
  const params = new URLSearchParams(searchParams.toString());
  params.set("page", pageNumber.toString());
  return `${pathname}?${params.toString()}`;
};
```

This ensures clicking page 3 while searching for "alice" produces:
`/dashboard?query=alice&size=10&page=3`
— not just `/dashboard?page=3` (losing the query).

**Trade-off accepted**: `<Link>` requires knowing the full URL at render time. If pagination logic were more dynamic (e.g., infinite scroll with a `fetchNextPage` function), `useRouter.push()` with an API call would be more appropriate. For standard page-based navigation, `<Link>` is superior.

---

## 11. Frontend: Search Debounce Implementation

### The Problem
Sending a backend request on every keystroke would hammer the database (e.g., typing "alice" = 5 requests in 200ms).

### Alternatives Considered

| Approach | How It Works | Pros | Cons |
|----------|-------------|------|------|
| **No debounce** | Fetch on every `onChange` | Simplest | Excessive API calls; poor UX; race conditions |
| **`use-debounce` library** | `npm install use-debounce` | Clean hook API | External dependency for a few lines of logic |
| **`setTimeout` + `useEffect` cleanup** | Effect with `clearTimeout` on re-render | No dependency | Cleanup is tricky; re-renders more than necessary |
| **✅ `useRef` + `setTimeout` (chosen)** | Store timer in ref; clear on next keystroke | No dependency; no extra re-renders; stable | Slightly more manual |

### Why `useRef` + `setTimeout`

```tsx
const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

const handleSearch = (value: string) => {
  if (debounceRef.current) clearTimeout(debounceRef.current);
  debounceRef.current = setTimeout(() => {
    updateParams({ query: value || null, page: "1" });
  }, 300);
};
```

- `useRef` holds the timer ID without triggering a re-render when it changes — unlike `useState`.
- Each keystroke cancels the previous timer and sets a new 300ms window.
- Only fires when the user **stops typing for 300ms**.

Additionally, `useTransition` wraps the `router.replace()` call:
```tsx
startTransition(() => replace(`${pathname}?${params}`));
```
This marks the navigation as non-urgent — React keeps the current page interactive (input is still typeable) while the Server Component re-renders in the background. `isPending` drives the spinner inside the search input.

**Trade-off accepted**: 300ms is a balance between responsiveness and server load. 100ms might feel more responsive but generates more requests. 500ms would reduce server load but feel sluggish. 300ms is the industry standard (used in Google Search, GitHub search).

---

## 12. Frontend: Error Boundary Scope

### The Problem
Where should errors be caught? A single global boundary vs granular per-segment boundaries.

### Alternatives Considered

| Approach | How It Works | Pros | Cons |
|----------|-------------|------|------|
| **No error boundaries** | Errors crash the whole app | No code | Terrible UX — a single failed fetch breaks everything |
| **Single global `app/error.tsx`** | One boundary for everything | Simple | A dashboard error also kills the sidebar/nav |
| **✅ Layered boundaries (chosen)** | `app/error.tsx` + `app/dashboard/error.tsx` | Dashboard errors are scoped; root shell stays rendered | Slightly more files |
| **Try/catch in Server Components** | Return error state from fetch | Fine for data errors | Can't catch unexpected thrown exceptions |

### Why Layered Error Boundaries

Next.js error boundaries are **segment-scoped**. `app/dashboard/error.tsx` only wraps the `<main>` content area, not the sidebar in `app/dashboard/layout.tsx`.

This means if the dashboard fails (e.g., backend is down):
- ✅ The sidebar navigation still renders — the user can navigate to other sections.
- ✅ The "Try again" button calls `reset()` which re-attempts only the dashboard segment — not a full page reload.
- ✅ The error shows the `error.digest` hash — a unique identifier that correlates with server logs for debugging.

For data-fetch errors (backend returning `503`), we chose **not** to throw an exception — instead returning `{ data: null, error: "..." }` and rendering an error card within `UserList`. This avoids triggering the error boundary for expected API failures, reserving the boundary for truly unexpected exceptions.

**Trade-off accepted**: Two `error.tsx` files instead of one means maintaining the UI in two places. Given they serve different scopes, this is a worthwhile separation.

---

## Summary: Decision Matrix

| Decision | Alternative | Chosen | Primary Reason |
|----------|-----------|--------|---------------|
| Pagination | In-memory slice | SQL LIMIT/OFFSET | DB does the work; only requested rows transferred |
| ORM | GORM / ent | Raw SQL | Full control; readable; no magic |
| HTTP Router | Gin / stdlib | gorilla/mux | Already in use; stable; path vars |
| Password hashing | SHA-256 | bcrypt | Adaptive cost; auto-salt; industry standard |
| Error propagation | String matching | Sentinel errors | Type-safe `errors.Is`; idiomatic Go |
| Rendering strategy | `useEffect` SPA | Server Component | No flash; data on first HTML; URL-driven |
| UI state | `useState` | URL search params | Bookmarkable; shareable; works with Back button |
| Component boundary | All client | Server at top, Client at leaves | Minimum JS bundle; faster First Contentful Paint |
| Pagination navigation | `onClick` + router | `<Link>` | Free prefetch; correct history behavior |
| Debounce | `use-debounce` package | `useRef` + `setTimeout` | Zero dependency; no extra re-renders |
| Error boundaries | Single global | Layered (global + dashboard) | Scoped recovery; sidebar stays rendered |
