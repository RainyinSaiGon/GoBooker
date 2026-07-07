"use client";

const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001";
import { useState, useEffect } from "react";

interface User {
  id?: string | number;
  name: string;
  email: string;
}

export default function DashBoardPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await fetch(`${apiUrl}/users`);
        if (!response.ok) throw new Error("Failed to fetch users");
        const data = await response.json();
        setUsers(data || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : "An error occurred");
      } finally {
        setIsLoading(false);
      }
    };
    fetchUsers();
  }, []);

  return (
    /* Removed full screen wrappers here; layout takes care of placing this card */
    <main className="w-full max-w-xl bg-white dark:bg-zinc-900 shadow-sm rounded-xl border border-zinc-200 dark:border-zinc-800 p-6 sm:p-8 space-y-6">
      <div>
        <h1 className="text-xl text-center font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">Dashboard</h1>
        <p className="text-xs text-zinc-400 text-center mt-1">Registered System Users</p>
      </div>

      <hr className="border-zinc-200 dark:border-zinc-800" />

      {isLoading && (
        <div className="text-center py-4 text-sm text-zinc-500 animate-pulse">
          Loading users from database...
        </div>
      )}

      {error && (
        <div className="text-center py-3 px-4 text-sm bg-red-50 text-red-600 dark:bg-red-950/30 dark:text-red-400 rounded-lg border border-red-200 dark:border-red-900">
          {error}
        </div>
      )}

      {!isLoading && !error && (
        <div className="space-y-3">
          {users.length === 0 ? (
            <p className="text-sm text-zinc-500 text-center py-4">No users found.</p>
          ) : (
            <ul className="divide-y divide-zinc-100 dark:divide-zinc-800 border border-zinc-100 dark:border-zinc-800 rounded-xl overflow-hidden">
              {users.map((user, index) => (
                <li 
                  key={user.id || index} 
                  className="flex flex-col p-4 bg-zinc-50/50 dark:bg-zinc-900/50 hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-colors"
                >
                  <span className="text-sm font-medium text-zinc-800 dark:text-zinc-200">
                    {user.name}
                  </span>
                  <span className="text-xs text-zinc-500 dark:text-zinc-400">
                    {user.email}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </div>
      )}
    </main>
  );
}