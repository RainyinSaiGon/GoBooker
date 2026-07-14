"use client";

import { useState, useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useAuthStore } from "@/lib/store/authStore";
import { api, User, UserInput } from "@/lib/api";
import { LogIn, LogOut, Trash2, Edit2, Search, UserPlus, RefreshCw, X } from "lucide-react";

export default function Home() {
  const [mounted, setMounted] = useState(false);
  const { token, setToken, logout } = useAuthStore();
  const queryClient = useQueryClient();

  // Login Form state
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loginError, setLoginError] = useState("");
  const [loginLoading, setLoginLoading] = useState(false);

  // Search state
  const [searchQuery, setSearchQuery] = useState("");

  // CRUD state
  const [isEditing, setIsEditing] = useState<string | null>(null); // user ID being edited
  const [nameInput, setNameInput] = useState("");
  const [emailInput, setEmailInput] = useState("");
  const [passwordInput, setPasswordInput] = useState("");
  const [roleInput, setRoleInput] = useState("customer");
  const [crudError, setCrudError] = useState("");
  const [crudSuccess, setCrudSuccess] = useState("");

  useEffect(() => {
    setMounted(true);
  }, []);

  // TanStack Query for fetching users
  const { data: users = [], isLoading, isError, error, refetch } = useQuery<User[]>({
    queryKey: ["users"],
    queryFn: api.getUsers,
    enabled: !!token && mounted,
  });

  // TanStack Query mutations
  const createMutation = useMutation({
    mutationFn: api.createUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      setCrudSuccess("User created successfully!");
      resetForm();
    },
    onError: (err: any) => {
      setCrudError(err.message || "Failed to create user");
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UserInput }) => api.updateUser(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      setCrudSuccess("User updated successfully!");
      resetForm();
    },
    onError: (err: any) => {
      setCrudError(err.message || "Failed to update user");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: api.deleteUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      setCrudSuccess("User deleted successfully!");
    },
    onError: (err: any) => {
      setCrudError(err.message || "Failed to delete user");
    },
  });

  const resetForm = () => {
    setNameInput("");
    setEmailInput("");
    setPasswordInput("");
    setRoleInput("customer");
    setIsEditing(null);
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoginError("");
    setLoginLoading(true);
    try {
      const response = await api.login(email, password);
      setToken(response.token);
      setEmail("");
      setPassword("");
    } catch (err: any) {
      setLoginError(err.message || "Invalid credentials");
    } finally {
      setLoginLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      await api.logout();
    } catch {
      // ignore
    }
    logout();
    queryClient.clear();
  };

  const handleSubmitUser = (e: React.FormEvent) => {
    e.preventDefault();
    setCrudError("");
    setCrudSuccess("");

    if (!nameInput || !emailInput) {
      setCrudError("Name and email are required");
      return;
    }

    const inputData: UserInput = {
      name: nameInput,
      email: emailInput,
      role: roleInput,
    };

    if (isEditing) {
      if (passwordInput) {
        inputData.password = passwordInput;
      }
      updateMutation.mutate({ id: isEditing, data: inputData });
    } else {
      if (!passwordInput) {
        setCrudError("Password is required for new users");
        return;
      }
      inputData.password = passwordInput;
      createMutation.mutate(inputData);
    }
  };

  const handleEditClick = (user: User) => {
    setIsEditing(user.id);
    setNameInput(user.name);
    setEmailInput(user.email);
    setRoleInput(user.role);
    setPasswordInput(""); // blank password unless changing
  };

  const handleDeleteClick = (id: string) => {
    if (confirm("Are you sure you want to delete this user?")) {
      setCrudError("");
      setCrudSuccess("");
      deleteMutation.mutate(id);
    }
  };

  // Filter users
  const filteredUsers = users.filter(
    (u) =>
      u.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      u.email.toLowerCase().includes(searchQuery.toLowerCase())
  );

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
    return (
      <div className="flex-center animate-fade-in" style={{ minHeight: "100vh", padding: "1rem" }}>
        <div className="card" style={{ maxWidth: "450px", width: "100%" }}>
          <div style={{ textAlign: "center", marginBottom: "2rem" }}>
            <h1 style={{ fontSize: "2rem", marginBottom: "0.5rem" }}>GoBooker</h1>
            <p className="text-muted">Sign in to manage users</p>
          </div>

          {loginError && <div className="alert alert-danger">{loginError}</div>}

          <form onSubmit={handleLogin}>
            <div className="form-group">
              <label className="form-label" htmlFor="email">Email Address</label>
              <input
                id="email"
                type="email"
                className="form-input"
                placeholder="admin@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>

            <div className="form-group" style={{ marginBottom: "1.5rem" }}>
              <label className="form-label" htmlFor="password">Password</label>
              <input
                id="password"
                type="password"
                className="form-input"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>

            <button type="submit" className="btn btn-primary" style={{ width: "100%" }} disabled={loginLoading}>
              <LogIn size={18} />
              {loginLoading ? "Signing in..." : "Sign In"}
            </button>
          </form>
        </div>
      </div>
    );
  }

  // Dashboard View
  return (
    <div className="container animate-fade-in" style={{ display: "flex", flexDirection: "column", gap: "2rem" }}>
      {/* Header */}
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

      {/* Main Content Area */}
      <div className="dashboard-grid">

        {/* Left Side: Users List */}
        <div className="card" style={{ display: "flex", flexDirection: "column", gap: "1.5rem" }}>
          <div className="flex-between">
            <h2 style={{ fontSize: "1.25rem" }}>Registered Users</h2>
            <button 
              onClick={() => refetch()} 
              className="btn btn-secondary" 
              style={{ padding: "0.5rem", borderRadius: "6px" }}
              title="Refresh"
            >
              <RefreshCw size={16} />
            </button>
          </div>

          {/* Search Box */}
          <div style={{ position: "relative" }}>
            <Search size={18} style={{ position: "absolute", left: "1rem", top: "50%", transform: "translateY(-50%)", color: "var(--text-muted)" }} />
            <input
              type="text"
              className="form-input"
              style={{ paddingLeft: "2.75rem" }}
              placeholder="Search by name or email..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>

          {crudSuccess && (
            <div className="alert alert-success flex-between" style={{ margin: 0 }}>
              <span>{crudSuccess}</span>
              <X size={16} style={{ cursor: "pointer" }} onClick={() => setCrudSuccess("")} />
            </div>
          )}

          {isLoading ? (
            <div className="flex-center" style={{ height: "200px" }}>
              <div className="text-muted">Loading users...</div>
            </div>
          ) : isError ? (
            <div className="alert alert-danger" style={{ margin: 0 }}>
              Error fetching users: {(error as any)?.message || "Unknown error"}
            </div>
          ) : filteredUsers.length === 0 ? (
            <div className="flex-center" style={{ height: "200px", border: "1px dashed var(--card-border)", borderRadius: "8px" }}>
              <div className="text-muted">No users found</div>
            </div>
          ) : (
            <div className="table-container">
              <table className="table">
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>Email</th>
                    <th>Role</th>
                    <th style={{ textAlign: "right" }}>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredUsers.map((user) => (
                    <tr key={user.id}>
                      <td>
                        <div style={{ fontWeight: 600 }}>{user.name}</div>
                      </td>
                      <td className="text-muted">{user.email}</td>
                      <td>
                        <span className={`badge badge-${user.role}`}>
                          {user.role}
                        </span>
                      </td>
                      <td style={{ textAlign: "right" }}>
                        <div style={{ display: "inline-flex", gap: "0.5rem" }}>
                          <button
                            onClick={() => handleEditClick(user)}
                            className="btn btn-secondary"
                            style={{ padding: "0.4rem", borderRadius: "6px" }}
                            title="Edit User"
                          >
                            <Edit2 size={14} />
                          </button>
                          <button
                            onClick={() => handleDeleteClick(user.id)}
                            className="btn btn-danger"
                            style={{ padding: "0.4rem", borderRadius: "6px" }}
                            title="Delete User"
                          >
                            <Trash2 size={14} />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>

        {/* Right Side: Create/Edit Form */}
        <div className="card" style={{ height: "fit-content" }}>
          <h2 style={{ fontSize: "1.25rem", marginBottom: "1.5rem" }} className="flex-between">
            <span>{isEditing ? "Edit User" : "Create New User"}</span>
            {isEditing && (
              <button onClick={resetForm} className="btn btn-secondary" style={{ padding: "0.25rem 0.5rem", fontSize: "0.8rem" }}>
                Cancel
              </button>
            )}
          </h2>

          {crudError && <div className="alert alert-danger" style={{ marginBottom: "1.25rem" }}>{crudError}</div>}

          <form onSubmit={handleSubmitUser}>
            <div className="form-group">
              <label className="form-label" htmlFor="userName">Name</label>
              <input
                id="userName"
                type="text"
                className="form-input"
                placeholder="John Doe"
                value={nameInput}
                onChange={(e) => setNameInput(e.target.value)}
                required
              />
            </div>

            <div className="form-group">
              <label className="form-label" htmlFor="userEmail">Email Address</label>
              <input
                id="userEmail"
                type="email"
                className="form-input"
                placeholder="john@example.com"
                value={emailInput}
                onChange={(e) => setEmailInput(e.target.value)}
                required
              />
            </div>

            <div className="form-group">
              <label className="form-label" htmlFor="userPassword">
                {isEditing ? "Password (leave blank to keep current)" : "Password"}
              </label>
              <input
                id="userPassword"
                type="password"
                className="form-input"
                placeholder="••••••••"
                value={passwordInput}
                onChange={(e) => setPasswordInput(e.target.value)}
                required={!isEditing}
              />
            </div>

            <div className="form-group" style={{ marginBottom: "1.5rem" }}>
              <label className="form-label" htmlFor="userRole">Role</label>
              <select
                id="userRole"
                className="form-select"
                value={roleInput}
                onChange={(e) => setRoleInput(e.target.value)}
              >
                <option value="customer">Customer</option>
                <option value="admin">Admin</option>
              </select>
            </div>

            <button type="submit" className="btn btn-primary" style={{ width: "100%" }} disabled={createMutation.isPending || updateMutation.isPending}>
              {isEditing ? <Edit2 size={16} /> : <UserPlus size={16} />}
              {createMutation.isPending || updateMutation.isPending 
                ? "Saving..." 
                : isEditing 
                ? "Update User" 
                : "Create User"}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
