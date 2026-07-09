import Link from "next/link";
import type { User } from "../page";

// ─── Avatar ───────────────────────────────────────────────────────────────────

function Avatar({ name }: { name: string }) {
  const safe     = name?.trim() || "?";
  const initials = safe.split(" ").map((n) => n[0] ?? "").join("").slice(0, 2).toUpperCase();
  const hue      = [...safe].reduce((acc, c) => acc + c.charCodeAt(0), 0) % 360;
  return (
    <div
      className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-sm font-semibold text-white"
      style={{ background: `hsl(${hue},60%,48%)` }}
      aria-hidden
    >
      {initials}
    </div>
  );
}

// ─── Role badge ───────────────────────────────────────────────────────────────

function RoleBadge({ role }: { role: string }) {
  const styles: Record<string, string> = {
    admin:    "bg-purple-100 text-purple-700 dark:bg-purple-950/50 dark:text-purple-300",
    customer: "bg-[var(--bg-subtle)] text-[var(--brand-600)]",
  };
  const cls = styles[role] ?? "bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-300";
  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-[10px] font-semibold uppercase tracking-wide ${cls}`}>
      {role}
    </span>
  );
}

// ─── UserList ─────────────────────────────────────────────────────────────────

interface UserListProps {
  users: User[];
  error: string | null;
  search: string;
}

export default function UserList({ users, error, search }: UserListProps) {
  // ── Error state ──
  if (error) {
    return (
      <div className="flex flex-col items-center gap-4 rounded-2xl border border-red-200 bg-red-50 py-12 text-center dark:border-red-900 dark:bg-red-950/20">
        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-950/50">
          <svg className="text-red-500" width="22" height="22" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.8}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v3m0 3h.01M5.07 19H19a2 2 0 001.73-3L13.73 4a2 2 0 00-3.46 0L4.27 16A2 2 0 005.07 19z"/>
          </svg>
        </div>
        <div>
          <p className="font-semibold text-red-600 dark:text-red-400">Could not load users</p>
          <p className="mt-0.5 text-sm text-red-500/80 dark:text-red-400/70">{error}</p>
        </div>
      </div>
    );
  }

  // ── Empty state ──
  if (users.length === 0) {
    return (
      <div className="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-[var(--border)] bg-[var(--bg-surface)] py-16 text-center">
        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-[var(--bg-subtle)]">
          <svg className="text-[var(--text-muted)]" width="20" height="20" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.6}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2M9 11a4 4 0 100-8 4 4 0 000 8z"/>
          </svg>
        </div>
        <div>
          <p className="font-semibold text-[var(--text-primary)]">
            {search ? "No matches found" : "No users yet"}
          </p>
          <p className="mt-0.5 text-sm text-[var(--text-secondary)]">
            {search
              ? `No users match "${search}"`
              : "Add your first user to get started."}
          </p>
        </div>
        {!search && (
          <Link
            href="/"
            className="mt-1 rounded-xl bg-[var(--brand-500)] px-4 py-2 text-sm font-semibold text-white transition hover:bg-[var(--brand-600)]"
          >
            Add First User
          </Link>
        )}
      </div>
    );
  }

  // ── User cards ──
  return (
    <div className="space-y-2.5">
      {users.map((user) => (
        <div
          key={user.email}
          className="group flex items-center gap-4 rounded-xl border border-[var(--border)] bg-[var(--bg-surface)] px-4 py-3.5 shadow-[var(--shadow-sm)] transition hover:border-[var(--brand-200)] hover:shadow-[var(--shadow-md)]"
        >
          <Avatar name={user.name} />
          <div className="flex-1 min-w-0">
            <p className="truncate text-sm font-semibold text-[var(--text-primary)] group-hover:text-[var(--brand-600)] transition-colors">
              {user.name}
            </p>
            <p className="truncate text-xs text-[var(--text-muted)]">{user.email}</p>
          </div>
          <div className="shrink-0">
            <RoleBadge role={user.role} />
          </div>
        </div>
      ))}
    </div>
  );
}
