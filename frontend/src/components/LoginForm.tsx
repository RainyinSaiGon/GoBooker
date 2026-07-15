"use client";

import { useState } from "react";
import { useAuthStore } from "@/lib/store/authStore";
import { apiClient } from "@/api";

export default function LoginForm() {
  const { setToken } = useAuthStore();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loginError, setLoginError] = useState("");
  const [loginLoading, setLoginLoading] = useState(false);

  const handleLogin = async (e: React.SubmitEvent) => {
    e.preventDefault();
    setLoginError("");
    setLoginLoading(true);
    try {
      const response = await apiClient.post<{ token: string }>("/auth/login", { email, password });
      setToken(response.data.token);
    } catch (err: any) {
      const msg = err.response?.data?.error || err.message || "Invalid credentials";
      setLoginError(msg);
    } finally {
      setLoginLoading(false);
    }
  };

  return (
    <div className="flex-center animate-fade-in" style={{ minHeight: "100vh", padding: "1rem" }}>
      <div className="card" style={{ maxWidth: "450px", width: "100%" }}>
        <div style={{ textAlign: "center", marginBottom: "2rem" }}>
          <h1 style={{ fontSize: "2rem", marginBottom: "0.5rem" }}>GoBooker</h1>
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

            {loginLoading ? "Signing in..." : "Sign In"}
          </button>
        </form>
      </div>
    </div>
  );
}
