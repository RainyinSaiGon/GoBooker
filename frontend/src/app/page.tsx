"use client";

import { useState, useEffect } from "react";
import { useAuthStore } from "@/lib/store/authStore";
import { User } from "@/api";
import LoginForm from "@/components/LoginForm";
import Header from "@/components/Header";
import UserList from "@/components/UserList";
import UserForm from "@/components/UserForm";

export default function Home() {
  const [mounted, setMounted] = useState(false);
  const { token } = useAuthStore();
  const [editingUser, setEditingUser] = useState<User | null>(null);

  useEffect(() => {
    setMounted(true);
  }, []);

  // Prevent SSR flickering / hydration mismatches
  if (!mounted) {
    return (
      <div className="flex-center animate-fade-in" style={{ minHeight: "100vh" }}>
        <div className="text-muted">Loading GoBooker...</div>
      </div>
    );
  }

  // Login View
  if (!token) {
    return <LoginForm />;
  }

  // Dashboard View
  return (
    <div className="container animate-fade-in" style={{ display: "flex", flexDirection: "column", gap: "2rem" }}>
      <Header />

      <div className="dashboard-grid">
        {/* Left Side: Users List */}
        <UserList 
          onEditUser={(user) => setEditingUser(user)} 
          editingUserId={editingUser ? editingUser.id : null}
        />

        {/* Right Side: Create/Edit Form */}
        <UserForm 
          editingUser={editingUser} 
          onFinished={() => setEditingUser(null)} 
        />
      </div>
    </div>
  );
}
