# GoBooker — Technical Implementation Report

> **Project**: GoBooker — A user management and booking platform  
> **Backend**: Go 1.26 · CockroachDB/PostgreSQL · gorilla/mux  
> **Frontend**: Next.js 16.2.10 (App Router) · React 19 · TypeScript 5

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Technology Stack](#2-technology-stack)
3. [Backend Implementation](#3-backend-implementation)
4. [Frontend Implementation](#4-frontend-implementation)
5. [Feature Sequence Flows](#5-feature-sequence-flows)
6. [Implementation Summary](#6-implementation-summary)

---

## 1. Architecture Overview

GoBooker is a **full-stack** web application built with a clear separation of concerns:

```
Browser (Next.js 16 Server + Client)
        │  HTTP (JSON)
        ▼
Go REST API (port 3001)
  ├── Middleware Layer  (CORS · Logger · Recovery)
  ├── Handler Layer     (HTTP ↔ Service interface)
  ├── Service Layer     (business rules · hashing · error mapping)
  └── Repository Layer  (SQL · LIMIT/OFFSET · ILIKE search)
        │  SQL (pgx/v5)
        ▼
CockroachDB / PostgreSQL
```

### Domain Model

```go
// model/user.go
type User struct {
    ID        string    // gen_random_uuid() — database-generated
    Name      string
    Email     string    // unique constraint in DB
    Password  string    // bcrypt hash — json:"-" never returned in API
    Role      string    // "customer" | "admin"
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### API Routes (all under `/api/v1`)

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/health` | inline | Health probe |
| GET | `/api/v1/users` | `GetAllUsers` | Paginated + searchable list |
| POST | `/api/v1/users` | `CreateUser` | Register a new user |
| GET | `/api/v1/users/{id}` | `GetUser` | Fetch single user |
| PUT | `/api/v1/users/{id}` | `UpdateUser` | Edit name/email/password |
| DELETE | `/api/v1/users/{id}` | `DeleteUser` | Remove user |

---

## 2. Technology Stack

| Layer | Technology | Why |
|-------|-----------|-----|
| Backend language | Go 1.26 | Strong typing, fast runtime, excellent HTTP stdlib |
| HTTP router | gorilla/mux | Path variable extraction (`{id}`) |
| Database driver | pgx/v5 (stdlib mode) | Native PostgreSQL wire protocol, prepared statements |
| Password hashing | bcrypt (golang.org/x/crypto) | Adaptive cost, industry standard |
| Frontend framework | Next.js 16 App Router | Server Components, `searchParams` as async props |
| Language | TypeScript 5 | Type safety across the full stack |
| CSS | Vanilla CSS + custom properties | Zero runtime, consistent design tokens |
| Font | `next/font/google` (Geist) | Self-hosted, zero layout shift |
| Package manager | bun | Fast installs, native TS support |

---

## 3. Backend Implementation

### 3.1 Entrypoint & Infrastructure (`main.go`)

```go
// Connection pool — prevents DB exhaustion under load
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(5 * time.Minute)

// Middleware chain (outermost → innermost)
chain := CORSMiddleware(cfg.AllowedOrigin)(
    Recovery(
        Logger(router)
    )
)
```

The middleware is applied in three layers:

1. **CORS** — sets `Access-Control-Allow-Origin`, handles `OPTIONS` preflight.
2. **Recovery** — catches any `panic` in a handler, logs the full goroutine stack trace (`debug.Stack()`), and returns `500`.
3. **Logger** — records method, path, status code, and latency per request using a wrapped `ResponseWriter` to prevent double-writing headers.

### 3.2 Three-Layer Architecture

#### Repository Layer (`repository/user_repository.go`)

Executes raw SQL. No business logic lives here. All queries use `$1` positional parameters (PostgreSQL/CockroachDB style).

Key improvement — **database-level LIMIT/OFFSET pagination**:
```sql
-- Count matching records (for total/totalPages)
SELECT COUNT(*) FROM users WHERE name ILIKE $1 OR email ILIKE $1

-- Fetch only the requested page
SELECT id, email, name, role, created_at, updated_at
FROM users
WHERE name ILIKE $1 OR email ILIKE $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
```

`ILIKE` provides case-insensitive full-text search across both `name` and `email`. The `% + query + %` wildcard pattern is built at the service layer.

The `UpdateUser` method conditionally includes `password` in the `SET` clause only if a new value is provided — preventing accidental password wipes:
```sql
-- Without new password
UPDATE users SET name=$1, email=$2, updated_at=$3 WHERE id=$4

-- With new password
UPDATE users SET name=$1, email=$2, password=$3, updated_at=$4 WHERE id=$5
```

#### Service Layer (`service/user_service.go`)

Business rules are enforced here:

- **Password hashing**: every `CreateUser` and `UpdateUser` (with new password) bcrypt-hashes the plaintext before the repo call.
- **Pagination normalization**: converts user-facing 1-based pages to 0-based DB offsets.
  ```go
  offset := (page - 1) * size
  ```
- **Duplicate email detection**: the `isDuplicateKey` helper inspects the raw error string for PostgreSQL SQLSTATE `23505` and returns the typed `ErrDuplicateEmail` sentinel.
- **Not-found mapping**: `sql.ErrNoRows` → `ErrNotFound` sentinel, so handlers don't need to import `database/sql`.

#### Handler Layer (`handler/user_handler.go`)

HTTP plumbing only — decode JSON, validate input, call service, write response.

Key implementations:
- **Email validation** (`helpers.go`): RFC 5321-compatible regex. Rejects `a@b` and `@foo.com`.
  ```go
  var emailRE = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
  ```
- **Request body limit**: `http.MaxBytesReader(w, r.Body, 1<<20)` — 1 MB cap on create/update to prevent DoS.
- **Paginated response envelope**:
  ```json
  {
    "users": [{ "name": "...", "email": "...", "role": "..." }],
    "total": 45,
    "page": 2,
    "size": 10,
    "totalPages": 5
  }
  ```
- **Correct HTTP status codes**:
  | Scenario | Status |
  |---------|--------|
  | Successful creation | `201 Created` |
  | Successful deletion | `204 No Content` |
  | Duplicate email | `409 Conflict` |
  | Not found | `404 Not Found` |
  | Bad input | `400 Bad Request` |
  | Server error | `500 Internal Server Error` |

---

## 4. Frontend Implementation

### 4.1 File Structure

```
app/
  layout.tsx          ← Root layout (Server Component) with full metadata
  page.tsx            ← Sign-up form (Client Component)
  error.tsx           ← Global error boundary (Client Component)
  not-found.tsx       ← Custom 404 page (Server Component)
  globals.css         ← CSS custom property design tokens
  dashboard/
    layout.tsx        ← Dashboard shell with sidebar (Server Component)
    page.tsx          ← Users list (Server Component, async searchParams)
    loading.tsx       ← Instant Loading UI skeleton
    error.tsx         ← Dashboard error boundary (Client Component)
    _components/
      SearchBar.tsx   ← Search + page size + refresh (Client Component)
      Pagination.tsx  ← Link-based page navigation (Client Component)
      UserList.tsx    ← User card renderer (display component)
      SignOutButton.tsx ← Isolated onClick leaf (Client Component)
lib/
  api.ts              ← Shared API_BASE constant
```

### 4.2 Design System (globals.css)

All colours, shadows, and radii are defined as CSS custom properties:
```css
--brand-500, --brand-600, --brand-400   /* primary action colours */
--bg-base, --bg-surface, --bg-subtle    /* background hierarchy */
--text-primary, --text-secondary, --text-muted
--border                                /* consistent border colour */
--shadow-sm, --shadow-md                /* box shadows */
```

These tokens are used throughout via `var(--token-name)` in JSX class names.

### 4.3 Server vs Client Component Split

| Component | Type | Reason |
|-----------|------|--------|
| `layout.tsx` | Server | No interactivity needed |
| `dashboard/page.tsx` | Server (async) | Fetches data using `await searchParams` |
| `dashboard/loading.tsx` | Server | Static skeleton markup |
| `dashboard/error.tsx` | **Client** | React error boundary API requires `reset()` |
| `not-found.tsx` | Server | Static 404 markup |
| `SearchBar.tsx` | **Client** | `useSearchParams`, `useRouter`, `useState` |
| `Pagination.tsx` | **Client** | `useSearchParams`, `usePathname` |
| `UserList.tsx` | Server | Pure display, receives props from page |
| `SignOutButton.tsx` | **Client** | `onClick` handler (leaf component) |

**Key rule applied**: Client Components are pushed to the leaves of the tree so the maximum amount of code renders on the server (better performance, no client JS bundle for static parts).

### 4.4 SEO & Metadata (Next.js docs best practices)

Root `layout.tsx` exports a full `Metadata` object:
```tsx
export const metadata: Metadata = {
  title: { default: "GoBooker – Book smarter", template: "%s | GoBooker" },
  description: "The modern booking platform for teams.",
  metadataBase: new URL(process.env.NEXT_PUBLIC_SITE_URL ?? "http://localhost:3000"),
  openGraph: { type: "website", locale: "en_US" },
  twitter: { card: "summary_large_image" },
  robots: { index: true, follow: true },
};
```

The `template` pattern means any child page that sets `title: "Dashboard"` automatically renders as `"Dashboard | GoBooker"` in the browser tab.

---

## 5. Feature Sequence Flows

### Feature 1: User Registration (Sign-Up)

```
User fills form on /
        │
        │ onSubmit
        ▼
[ app/page.tsx ]  (Client Component)
  1. Client-side validation
     • name not empty
     • email matches /^[^\s@]+@[^\s@]+\.[^\s@]+$/
     • password ≥ 8 chars
  2. If invalid → shows inline field errors (no server call)
  3. If valid  → POST /api/v1/users
                 Body: { name, email, password }
        │
        │ HTTP POST
        ▼
[ Go middleware chain ]
  CORS → Recovery → Logger
        │
        ▼
[ handler/user_handler.go : CreateUser ]
  1. MaxBytesReader(1 MB)
  2. json.Decode → UserRequest
  3. Check name/email/password present
  4. len(password) ≥ 8
  5. isValidEmail (RFC 5321 regex)
  6. h.svc.CreateUser(...)
        │
        ▼
[ service/user_service.go : CreateUser ]
  1. Set default role = "customer" if empty
  2. bcrypt.GenerateFromPassword(password, DefaultCost)
  3. repo.CreateUser(user with hashed password)
        │
        ▼
[ repository/user_repository.go : CreateUser ]
  SQL: INSERT INTO users (id,email,name,password,role,created_at,updated_at)
       VALUES (gen_random_uuid(),$1,$2,$3,$4,$5,$5)
       RETURNING id, created_at, updated_at
        │
        ▼ (error path: unique constraint violation)
[ service ] isDuplicateKey(err) → return ErrDuplicateEmail
[ handler ] errors.Is(ErrDuplicateEmail) → 409 Conflict
            { "error": "An account with that email already exists" }
        │
        ▼ (success path)
[ handler ] 201 Created
            { "name": "...", "email": "...", "role": "customer" }
        │
        ▼
[ app/page.tsx ]
  • Shows success state with confirmation UI
  • Link to /dashboard
```

---

### Feature 2: Browse Users with Search (Dashboard)

```
User navigates to /dashboard
(or types in search box / changes page size)
        │
URL: /dashboard?query=alice&page=2&size=10
        │
        ▼
[ Next.js Server Router ]
  Shows loading.tsx skeleton immediately
        │
        ▼
[ app/dashboard/page.tsx ] (async Server Component)
  const params = await searchParams;   // Next.js 15+ Promise API
  const query  = params.query ?? ""
  const page   = Number(params.page) || 1
  const size   = params.size in [5,10,25] ? params.size : 10

  fetchUsers(query, page, size)
        │
        │ HTTP GET (server → backend, never leaves server)
        ▼
[ GET /api/v1/users?query=alice&page=2&size=10 ]
        │
        ▼
[ handler/user_handler.go : GetAllUsers ]
  q.Get("query") → "alice"
  strconv.Atoi("page")  → 2
  strconv.Atoi("size")  → 10
  h.svc.GetAllUsers("alice", 2, 10)
        │
        ▼
[ service/user_service.go : GetAllUsers ]
  offset = (2-1) * 10 = 10
  repo.GetAllUsers("alice", 10, 10)
        │
        ▼
[ repository/user_repository.go : GetAllUsers ]
  queryParam = "%alice%"
  
  ① SELECT COUNT(*) FROM users
     WHERE name ILIKE '%alice%' OR email ILIKE '%alice%'
     → total = 15
  
  ② SELECT id,email,name,role,created_at,updated_at
     FROM users
     WHERE name ILIKE '%alice%' OR email ILIKE '%alice%'
     ORDER BY created_at DESC
     LIMIT 10 OFFSET 10
     → []User (page 2 of matching records)
        │
        ▼
[ handler ]  totalPages = ceil(15 / 10) = 2
  Response: {
    "users": [...],   ← only 10 rows, not all 15
    "total": 15,
    "page": 2,
    "size": 10,
    "totalPages": 2
  }
        │
        ▼
[ app/dashboard/page.tsx ]
  Renders: <UserList users={users} />
           <Pagination page=2 totalPages=2 total=15 from=11 to=15 />
```

---

### Feature 3: Search with Debounce (Client Interaction)

```
User types "ali" in search box
        │
        ▼
[ SearchBar.tsx ] (Client Component)
  onChange → handleSearch("ali")
  clearTimeout(debounceRef.current)     // cancel previous timer
  debounceRef.current = setTimeout(..., 300ms)
  
  ... user continues typing "alice" within 300ms ...
  
  clearTimeout → debounceRef.current = setTimeout("alice", 300ms)
  
  ... 300ms passes without typing ...
  
  startTransition(() => {
    replace("/dashboard?query=alice&page=1&size=10")
    //       ─── preserve current size ───┘
  })
        │
        ▼
[ Browser URL changes to /dashboard?query=alice&page=1&size=10 ]
  Next.js detects URL change
  → Re-renders app/dashboard/page.tsx (Server Component)
  → New fetch to backend with query=alice
  → input shows isPending=true (spinner while server re-renders)
        │
        ▼
[ SearchBar.tsx re-renders ]
  isPending → false
  input still shows "alice" (uncontrolled, not reset by URL)
```

---

### Feature 4: Page Navigation (Pagination)

```
User clicks page "3" button in Pagination component
        │
        ▼
[ Pagination.tsx ] (Client Component)
  createPageURL(3):
    params = new URLSearchParams("query=alice&page=2&size=10")
    params.set("page", "3")
    → "/dashboard?query=alice&page=3&size=10"
  
  <Link href="/dashboard?query=alice&page=3&size=10">3</Link>
  (Link is rendered — no onClick needed)
        │
        ▼
[ Next.js Link prefetch ]
  On hover: prefetches the route data automatically
  On click: navigates → browser URL changes
        │
        ▼
[ app/dashboard/page.tsx ] (Server Component re-renders)
  page=3, query="alice", size=10
  → fetchUsers("alice", 3, 10)
  → backend SQL: LIMIT 10 OFFSET 20
        │
        ▼
[ Backend responds ]
  { "users": [...], "total": 15, "page": 3, "totalPages": 2 }
  
  Note: page=3 > totalPages=2
  → safePage = Math.min(3, 2) = 2 (graceful clamp)
```

---

### Feature 5: Error Boundary Recovery

```
Backend unreachable (e.g. network error)
        │
        ▼
[ app/dashboard/page.tsx ] (Server Component)
  fetchUsers() catch block
  → returns { data: null, error: "Server responded with 503" }
  
  (Error is not thrown — it's returned as a typed value)
  → UserList rendered with error prop
  → Shows error card with message
  
  (If an unhandled exception IS thrown in the Server Component)
        │
        ▼
[ app/dashboard/error.tsx ] (React Error Boundary — Client Component)
  receives: error.message, error.digest (hash for log correlation)
  
  Renders:
    "Dashboard failed to load"
    [Try again] button → calls reset()
  
  reset() → Next.js retries the segment without a full page reload
```

---

## 6. Implementation Summary

### Backend Changes Made

| File | What Changed | Why |
|------|-------------|-----|
| `main.go` | Added connection pool settings (`SetMaxOpenConns(25)`, `SetMaxIdleConns(10)`, `SetConnMaxLifetime`), graceful startup logging | Prevent DB exhaustion, faster idle connection reuse |
| `middleware/middleware.go` | Added `responseWriter` wrapper to prevent double-`WriteHeader`; added `debug.Stack()` in `Recovery` | Prevents subtle HTTP header panics; makes panics actionable in logs |
| `handler/helpers.go` | Replaced naive `strings.Contains` email check with RFC 5321 regex | Rejects clearly invalid emails like `a@b` or `@foo.com` |
| `handler/router.go` | Added `/api/v1` versioned sub-router; added `/health` probe at root | API versioning allows future breaking changes without breaking old clients |
| `handler/user_handler.go` | `201` for create; `204` for delete; `409` for duplicate email; 1 MB body limit; paginated response struct; query/page/size param parsing | Correct REST semantics; DoS prevention; database-level pagination |
| `service/user_service.go` | `ErrDuplicateEmail`/`ErrNotFound` sentinels; `isDuplicateKey` helper; bcrypt hashing; pagination offset calculation | Clean error mapping; prevents password leakage; proper page → offset math |
| `repository/user_repository.go` | Conditional `SET` for password-less updates; `ILIKE` search; `COUNT(*) + LIMIT/OFFSET` pagination | Prevents password wipe bug; case-insensitive search; only fetches requested rows |

### Frontend Changes Made

| File | What Changed | Why |
|------|-------------|-----|
| `app/layout.tsx` | Full `Metadata` object with title template, OG tags, twitter card, robots | SEO and social sharing per Next.js docs |
| `app/page.tsx` | Form validation, `409 Conflict` handling, `htmlFor` on all labels, shared `API_BASE` | Accessibility; correct duplicate email UX |
| `app/error.tsx` | Created global error boundary | Catches unexpected runtime errors, provides retry |
| `app/not-found.tsx` | Created custom 404 page | Better UX for unknown routes |
| `app/dashboard/page.tsx` | Converted to async Server Component; `await searchParams`; database-level pagination fetch | Server-side rendering; URL as single source of truth; no in-memory slicing |
| `app/dashboard/loading.tsx` | Created skeleton UI | Instant perceived performance while server fetches |
| `app/dashboard/error.tsx` | Created segment error boundary | Scoped error recovery, sidebar stays rendered |
| `app/dashboard/layout.tsx` | Moved `onClick` to `SignOutButton` leaf component | Keeps layout as Server Component |
| `SearchBar.tsx` | Created: debounced input, `useTransition` spinner, page-size selector, refresh button | Optimal UX; prevents stale page on search change |
| `Pagination.tsx` | Created: `<Link>`-based navigation, `createPageURL` param merger, ellipsis, `aria-current` | Free prefetch; bookmark-friendly URLs; accessibility |
| `UserList.tsx` | Created: avatar with deterministic colour hash, role badge, empty/error states | Consistent visual identity |
| `SignOutButton.tsx` | Created as isolated Client Component leaf | Required for `onClick` without converting layout |
| `lib/api.ts` | Created: shared `API_BASE` constant from env var | Single source of truth for backend URL |
| `.env.local.example` | Created developer environment guide | Onboarding documentation |

### Key Architectural Principles Applied

1. **Database-level pagination** — SQL `LIMIT/OFFSET` + `COUNT(*)` instead of in-memory slicing. The frontend never receives more rows than are displayed.
2. **URL as single source of truth** — search query, page number, and page size all live in the URL. Pages are bookmarkable, shareable, and work with the browser's Back button.
3. **Server Components at the top, Client Components at the leaves** — the dashboard page fetches and renders on the server; only interactive UI (search input, pagination buttons, sign-out) runs in the browser.
4. **Typed error sentinels** — `ErrDuplicateEmail`, `ErrNotFound` propagate cleanly from repository → service → handler without string matching across layers.
5. **Defensive HTTP** — 1 MB body limits, idempotent response codes, panic recovery with full stack traces, CORS on all routes.
