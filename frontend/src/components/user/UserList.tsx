"use client";

import { useState } from "react";
import { User, useUsers, useDeleteUser } from "@/hooks";
import { Trash2, Edit2, Search, RefreshCw, X } from "lucide-react";

interface UserListProps {
  onEditUser: (user: User) => void;
  editingUserId: string | null;
}

export default function UserList({ onEditUser, editingUserId }: UserListProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [crudError, setCrudError] = useState("");
  const [crudSuccess, setCrudSuccess] = useState("");

  const { data: users = [], isLoading, isError, error, refetch } = useUsers();

  const deleteMutation = useDeleteUser({
    onSuccess: () => {
      setCrudSuccess("User deleted successfully!");
    },
    onError: (err: any) => {
      const msg = err.response?.data?.error || err.message || "Failed to delete user";
      setCrudError(msg);
    },
  });

  const handleDelete = (id: string) => {
    if (confirm("Are you sure you want to delete this user?")) {
      setCrudError("");
      setCrudSuccess("");
      deleteMutation.mutate(id);
    }
  };

  const filteredUsers = users.filter(
    (u) =>
      u.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      u.email.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
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

      {crudError && (
        <div className="alert alert-danger flex-between" style={{ margin: 0 }}>
          <span>{crudError}</span>
          <X size={16} style={{ cursor: "pointer" }} onClick={() => setCrudError("")} />
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
              {filteredUsers.map((user) => {
                const isUserEditing = editingUserId === user.id;
                return (
                  <tr key={user.id} style={isUserEditing ? { background: "rgba(99, 102, 241, 0.05)" } : {}}>
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
                          onClick={() => onEditUser(user)}
                          className="btn btn-secondary"
                          style={{ 
                            padding: "0.4rem", 
                            borderRadius: "6px",
                            borderColor: isUserEditing ? "var(--primary)" : "var(--card-border)"
                          }}
                          title="Edit User"
                        >
                          <Edit2 size={14} style={isUserEditing ? { color: "var(--primary)" } : {}} />
                        </button>
                        <button
                          onClick={() => handleDelete(user.id)}
                          className="btn btn-danger"
                          style={{ padding: "0.4rem", borderRadius: "6px" }}
                          title="Delete User"
                          disabled={deleteMutation.isPending}
                        >
                          <Trash2 size={14} />
                        </button>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
