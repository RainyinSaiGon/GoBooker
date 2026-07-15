"use client";

import { useAuthStore } from "@/lib/store/authStore";
import { useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/api";
import { LogOut } from "lucide-react";

export default function Header() {
  const { logout } = useAuthStore();
  const queryClient = useQueryClient();

  const handleLogout = async () => {
    try {
      await apiClient.post("/auth/logout");
    } catch {
      // ignore
    }
    logout();
    queryClient.clear();
  };

  return (
    <header className="card flex-between" style={{ padding: "1rem 2rem" }}>
      <div>
        <h1 style={{ fontSize: "1.5rem" }}>GoBooker</h1>
        <p className="text-muted" style={{ fontSize: "0.85rem" }}>User Management Portal</p>
      </div>
      <button onClick={handleLogout} className="btn btn-secondary">
        <LogOut size={16} />
        Sign Out
      </button>
    </header>
  );
}
