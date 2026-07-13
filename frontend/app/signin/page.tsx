"use client";

import { useState } from "react";
import {useRouter } from "next/navigation";
import { API_BASE } from "@/lib/api";
import router from "next/dist/shared/lib/router/router";

interface FormErrors { email?: string; password?: string; }

function FieldError({ msg }: { msg?: string }) {
  if (!msg) return null;
  return (
    <p className="mt-1.5 flex items-center gap-1 text-xs text-red-500">
      <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
        <circle cx="6" cy="6" r="5.5" stroke="currentColor"/>
        <path d="M6 4v3M6 8.5v.5" stroke="currentColor" strokeLinecap="round"/>
      </svg>
      {msg}
    </p>
  );
}

export default function SignInPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [errors, setErrors] = useState<FormErrors>({});
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const router = useRouter()

  const validate = (v: { email: string; password: string }): FormErrors => {
    const e: FormErrors = {};
    if (!v.email.trim()) e.email = "Email is required";
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v.email)) e.email = "Enter a valid email address";
    if (!v.password) e.password = "Password is required";
    else if (v.password.length < 8) e.password = "Must be at least 8 characters";
    return e;
  };

  const handleBlur = (field: keyof FormErrors) => {
    setTouched((p) => ({ ...p, [field]: true }));
    setErrors(validate({ email, password }));
  };

  const inputClass = (field: keyof FormErrors) =>
    `w-full px-4 py-2.5 text-sm rounded-xl border transition-all duration-150 bg-white dark:bg-[var(--bg-subtle)] text-[var(--text-primary)] placeholder:text-[var(--text-muted)] focus:outline-none focus:ring-2 ${
      touched[field] && errors[field]
        ? "border-red-400 focus:ring-red-300 dark:focus:ring-red-800"
        : "border-[var(--border)] focus:ring-[var(--brand-200)] focus:border-[var(--brand-500)]"
    }`;

  const handleSubmit = async (e: React.SubmitEvent) => {
    e.preventDefault();
    const newErrors = validate({ email, password });
    setErrors(newErrors);
    setTouched({ email: true, password: true });
    if (Object.keys(newErrors).length > 0) return;

    try {
      const res = await fetch(`${API_BASE}/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });
      if (!res.ok) {
        const d = await res.json().catch(() => ({}));
        if (res.status === 409) {
          throw new Error(d.error || "An account with that email already exists");
        }
        throw new Error(d.error || "Failed to create account");
      }
      router.push("/dashboard");

    
    } catch (err) {
        console.error(err);
    } 
  };

  return (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden p-4">
      {/* Gradient backdrop blobs */}
      <div aria-hidden className="pointer-events-none absolute inset-0 overflow-hidden">
        <div className="absolute -top-32 -left-32 h-125 w-125 rounded-full bg-(--brand-500)opacity-[.07] blur-[80px]" />
        <div className="absolute -bottom-20 right-0 h-100 w-100 rounded-full bg-purple-500 opacity-[.07] blur-[80px]" />
      </div>

      <main className="animate-fade-in relative z-10 w-full max-w-md">
        {/* Logo wordmark */}
        <div className="mb-8 text-center">
          <span className="gradient-text text-3xl font-bold tracking-tight">GoBooker</span>
          <p className="mt-1 text-sm text-(--text-secondary)">Book smarter, together.</p>
        </div>

        <div className="glass rounded-2xl p-8 shadow-(--shadow-lg)">
        

              <form onSubmit={handleSubmit} noValidate className="space-y-5">
                {/* Email */}
                <div>
                  <label htmlFor="email" className="mb-1.5 block text-xs font-semibold uppercase tracking-wider text-(--text-secondary)">
                    Email <span className="text-red-500 normal-case tracking-normal font-normal">*</span>
                  </label>
                  <input
                    id="email" type="email" autoComplete="email"
                    value={email} onChange={(e) => setEmail(e.target.value)} onBlur={() => handleBlur("email")}
                    placeholder="jane@example.com" 
                    className={inputClass("email")}
                    
                  />
                  <FieldError msg={touched.email ? errors.email : undefined} />
                </div>

                {/* Password */}
                <div>
                  <label htmlFor="password" className="mb-1.5 block text-xs font-semibold uppercase tracking-wider text-(--text-secondary)">
                    Password <span className="text-red-500 normal-case tracking-normal font-normal">*</span>
                  </label>
                  <div className="relative">
                    <input
                      id="password" type={showPassword ? "text" : "password"} autoComplete="new-password"
                      value={password} onChange={(e) => setPassword(e.target.value)} onBlur={() => handleBlur("password")}
                      placeholder="Min. 8 characters"
                      className={`${inputClass("password")} pr-10`}
                    />
                    <button
                      type="button" onClick={() => setShowPassword((p) => !p)}
                      className="absolute right-3 top-1/2 -translate-y-1/2 text-(--text-muted) hover:text-(--text-secondary) transition-colors"
                      aria-label={showPassword ? "Hide password" : "Show password"}
                    >
                      {showPassword ? (
                        <svg width="16" height="16" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.8}>
                          <path strokeLinecap="round" strokeLinejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-5 0-9-4-9-7s4-7 9-7c1.06 0 2.08.2 3.025.55M6.1 6.1A9.978 9.978 0 003 12c0 3 4 7 9 7a9.978 9.978 0 005.9-1.9M21 21L3 3"/>
                        </svg>
                      ) : (
                        <svg width="16" height="16" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.8}>
                          <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                          <path strokeLinecap="round" strokeLinejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                        </svg>
                      )}
                    </button>
                  </div>
                  {touched.password && errors.password
                    ? <FieldError msg={errors.password} />
                    : <p className="mt-1.5 text-xs text-(--text-muted)">Must be at least 8 characters</p>
                  }
                </div>

                <button
                  type="submit" id="signup-submit"
                  className="relative w-full overflow-hidden rounded-xl bg-(--brand-500) py-2.5 text-sm font-semibold text-white shadow-md transition-all duration-200 hover:bg-[var(--brand-600)] active:scale-[.98] disabled:cursor-not-allowed disabled:opacity-60 focus:outline-none focus:ring-2 focus:ring-[var(--brand-400)] focus:ring-offset-2"
                >
                Sign in
                </button>
              </form>
         </div>

        <p className="mt-6 text-center text-xs text-(--text-muted)">
          By signing up you agree to our{" "}
          <span className="cursor-pointer underline underline-offset-2 hover:text-(--brand-500)">Terms</span>
          {" & "}
          <span className="cursor-pointer underline underline-offset-2 hover:text-(--brand-500)">Privacy Policy</span>.
        </p>
      </main>
    </div>
  );
}
