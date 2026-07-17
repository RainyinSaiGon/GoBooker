"use client";

import { useState, useEffect } from "react";
import { User, UserInput, useCreateUser, useUpdateUser } from "@/hooks";
import { UserPlus, Edit2 } from "lucide-react";

interface UserFormProps {
  editingUser: User | null;
  onFinished: () => void;
}

export default function UserForm({ editingUser, onFinished }: UserFormProps) {
  const [nameInput, setNameInput] = useState("");
  const [emailInput, setEmailInput] = useState("");
  const [passwordInput, setPasswordInput] = useState("");
  
  const [errorMsg, setErrorMsg] = useState("");
  const [successMsg, setSuccessMsg] = useState("");

  // Sync inputs when editingUser changes
  useEffect(() => {
    if (editingUser) {
      setNameInput(editingUser.name);
      setEmailInput(editingUser.email);
      setPasswordInput(""); // Blank password unless changing
      setErrorMsg("");
      setSuccessMsg("");
    } else {
      resetForm();
    }
  }, [editingUser]);

  const resetForm = () => {
    setNameInput("");
    setEmailInput("");
    setPasswordInput("");
  };

  const createMutation = useCreateUser({
    onSuccess: () => {
      setSuccessMsg("User created successfully!");
      resetForm();
      onFinished();
    },
    onError: (err: any) => {
      const msg = err.response?.data?.error || err.message || "Failed to create user";
      setErrorMsg(msg);
    },
  });

  const updateMutation = useUpdateUser({
    onSuccess: () => {
      setSuccessMsg("User updated successfully!");
      resetForm();
      onFinished();
    },
    onError: (err: any) => {
      const msg = err.response?.data?.error || err.message || "Failed to update user";
      setErrorMsg(msg);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg("");
    setSuccessMsg("");

    if (!nameInput || !emailInput) {
      setErrorMsg("Name and email are required");
      return;
    }

    const inputData: UserInput = {
      name: nameInput,
      email: emailInput
    };

    if (editingUser) {
      if (passwordInput) {
        inputData.password = passwordInput;
      }
      updateMutation.mutate({ id: editingUser.id, data: inputData });
    } else {
      if (!passwordInput) {
        setErrorMsg("Password is required for new users");
        return;
      }
      inputData.password = passwordInput;
      createMutation.mutate(inputData);
    }
  };

  return (
    <div className="card" style={{ height: "fit-content" }}>
      <h2 style={{ fontSize: "1.25rem", marginBottom: "1.5rem" }} className="flex-between">
        <span>{editingUser ? "Edit User" : "Create New User"}</span>
        {editingUser && (
          <button 
            type="button"
            onClick={onFinished} 
            className="btn btn-secondary" 
            style={{ padding: "0.25rem 0.5rem", fontSize: "0.8rem" }}
          >
            Cancel
          </button>
        )}
      </h2>

      {errorMsg && <div className="alert alert-danger" style={{ marginBottom: "1.25rem" }}>{errorMsg}</div>}
      {successMsg && <div className="alert alert-success" style={{ marginBottom: "1.25rem" }}>{successMsg}</div>}

      <form onSubmit={handleSubmit}>
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

          {/* Email */}
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

        {/* Password */}
        <div className="form-group">
          <label className="form-label" htmlFor="userPassword">
            {editingUser ? "Password (leave blank to keep current)" : "Password"}
          </label>
          <input
            id="userPassword"
            type="password"
            className="form-input"
            placeholder="••••••••"
            value={passwordInput}
            onChange={(e) => setPasswordInput(e.target.value)}
            required={!editingUser}
          />
        </div>

  

        <button 
          type="submit" 
          className="btn btn-primary" 
          style={{ width: "100%" }} 
          disabled={createMutation.isPending || updateMutation.isPending}
        >
          {editingUser ? <Edit2 size={16} /> : <UserPlus size={16} />}
          {createMutation.isPending || updateMutation.isPending 
            ? "Saving..." 
            : editingUser 
            ? "Update User" 
            : "Create User"}
        </button>
      </form>
    </div>
  );
}
