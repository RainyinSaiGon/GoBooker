"use client";

import { useState } from "react";
import Link from "next/link";

const API = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001";

interface FormErrors { name?: string; email?: string; password?: string; }

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

export default function SignUpPage() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errors, setErrors] = useState<FormErrors>({});
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [success, setSuccess] = useState(false);
  const [serverError, setServerError] = useState("");

  const validate = (v: { name: string; email: string; password: string }): FormErrors => {
    const e: FormErrors = {};
    if (!v.name.trim()) e.name = "Full name is required";
    if (!v.email.trim()) e.email = "Email is required";
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v.email)) e.email = "Enter a valid email address";
    if (!v.password) e.password = "Password is required";
    else if (v.password.length < 8) e.password = "Must be at least 8 characters";
    return e;
  };

  const handleBlur = (field: keyof FormErrors) => {
    setTouched((p) => ({ ...p, [field]: true }));
    setErrors(validate({ name, email, password }));
  };

  const inputClass = (field: keyof FormErrors) =>
    `w-full px-4 py-2.5 text-sm rounded-xl border transition-all duration-150 bg-white dark:bg-[var(--bg-subtle)] text-[var(--text-primary)] placeholder:text-[var(--text-muted)] focus:outline-none focus:ring-2 ${
      touched[field] && errors[field]
        ? "border-red-400 focus:ring-red-300 dark:focus:ring-red-800"
        : "border-[var(--border)] focus:ring-[var(--brand-200)] focus:border-[var(--brand-500)]"
    }`;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setServerError("");
    const newErrors = validate({ name, email, password });
    setErrors(newErrors);
    setTouched({ name: true, email: true, password: true });
    if (Object.keys(newErrors).length > 0) return;

    setIsSubmitting(true);
    try {
      const res = await fetch(`${API}/users`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, email, password }),
      });
      if (!res.ok) {
        const d = await res.json().catch(() => ({}));
        throw new Error(d.error || "Failed to create account");
      }
      setSuccess(true);
      setName(""); setEmail(""); setPassword(""); setErrors({}); setTouched({});
    } catch (err) {
      setServerError((err as Error).message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden p-4">
      {/* Gradient backdrop blobs */}
      <div aria-hidden className="pointer-events-none absolute inset-0 overflow-hidden">
        <div className="absolute -top-32 -left-32 h-[500px] w-[500px] rounded-full bg-[var(--brand-500)] opacity-[.07] blur-[80px]" />
        <div className="absolute -bottom-20 right-0 h-[400px] w-[400px] rounded-full bg-purple-500 opacity-[.07] blur-[80px]" />
      </div>

      <main className="animate-fade-in relative z-10 w-full max-w-md">
        {/* Logo wordmark */}
        <div className="mb-8 text-center">
          <span className="gradient-text text-3xl font-bold tracking-tight">GoBooker</span>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">Book smarter, together.</p>
        </div>

        <div className="glass rounded-2xl p-8 shadow-[var(--shadow-lg)]">
          {success ? (
            <div className="flex flex-col items-center gap-4 py-6 text-center animate-fade-in">
              <div className="flex h-16 w-16 items-center justify-center rounded-full bg-green-100 dark:bg-green-950/50">
                <svg className="h-8 w-8 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <div>
                <h2 className="text-lg font-semibold text-[var(--text-primary)]">Account created!</h2>
                <p className="mt-1 text-sm text-[var(--text-secondary)]">Welcome aboard. Head to the dashboard.</p>
              </div>
              <Link
                href="/dashboard"
                className="mt-2 inline-flex items-center gap-2 rounded-xl bg-[var(--brand-500)] px-6 py-2.5 text-sm font-medium text-white shadow-md transition hover:bg-[var(--brand-600)] focus:outline-none focus:ring-2 focus:ring-[var(--brand-400)] focus:ring-offset-2"
              >
                Go to Dashboard
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M5 12h14M12 5l7 7-7 7" />
                </svg>
              </Link>
            </div>
          ) : (
            <>
              <div className="mb-6">
                <h1 className="text-xl font-bold text-[var(--text-primary)]">Create your account</h1>
                <p className="mt-1 text-sm text-[var(--text-secondary)]">Already have one?{" "}
                  <Link href="/dashboard" className="font-medium text-[var(--brand-500)] hover:underline">Sign in</Link>
                </p>
              </div>

              {serverError && (
                <div className="mb-4 flex items-start gap-2.5 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-600 dark:border-red-900 dark:bg-red-950/40 dark:text-red-400">
                  <svg className="mt-0.5 shrink-0" width="14" height="14" viewBox="0 0 12 12" fill="none">
                    <circle cx="6" cy="6" r="5.5" stroke="currentColor"/>
                    <path d="M6 4v3M6 8.5v.5" stroke="currentColor" strokeLinecap="round"/>
                  </svg>
                  {serverError}
                </div>
              )}

              <form onSubmit={handleSubmit} noValidate className="space-y-5">
                {/* Full name */}
                <div>
                  <label className="mb-1.5 block text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)]">
                    Full Name <span className="text-red-500 normal-case tracking-normal font-normal">*</span>
                  </label>
                  <input
                    id="name" type="text" autoComplete="name"
                    value={name} onChange={(e) => setName(e.target.value)} onBlur={() => handleBlur("name")}
                    placeholder="Jane Doe" disabled={isSubmitting}
                    className={inputClass("name")}
                  />
                  <FieldError msg={touched.name ? errors.name : undefined} />
                </div>

                {/* Email */}
                <div>
                  <label className="mb-1.5 block text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)]">
                    Email <span className="text-red-500 normal-case tracking-normal font-normal">*</span>
                  </label>
                  <input
                    id="email" type="email" autoComplete="email"
                    value={email} onChange={(e) => setEmail(e.target.value)} onBlur={() => handleBlur("email")}
                    placeholder="jane@example.com" disabled={isSubmitting}
                    className={inputClass("email")}
                  />
                  <FieldError msg={touched.email ? errors.email : undefined} />
                </div>

                {/* Password */}
                <div>
                  <label className="mb-1.5 block text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)]">
                    Password <span className="text-red-500 normal-case tracking-normal font-normal">*</span>
                  </label>
                  <div className="relative">
                    <input
                      id="password" type={showPassword ? "text" : "password"} autoComplete="new-password"
                      value={password} onChange={(e) => setPassword(e.target.value)} onBlur={() => handleBlur("password")}
                      placeholder="Min. 8 characters" disabled={isSubmitting}
                      className={`${inputClass("password")} pr-10`}
                    />
                    <button
                      type="button" onClick={() => setShowPassword((p) => !p)}
                      className="absolute right-3 top-1/2 -translate-y-1/2 text-[var(--text-muted)] hover:text-[var(--text-secondary)] transition-colors"
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
                    : <p className="mt-1.5 text-xs text-[var(--text-muted)]">Must be at least 8 characters</p>
                  }
                </div>

                <button
                  type="submit" id="signup-submit" disabled={isSubmitting}
                  className="relative w-full overflow-hidden rounded-xl bg-[var(--brand-500)] py-2.5 text-sm font-semibold text-white shadow-md transition-all duration-200 hover:bg-[var(--brand-600)] active:scale-[.98] disabled:cursor-not-allowed disabled:opacity-60 focus:outline-none focus:ring-2 focus:ring-[var(--brand-400)] focus:ring-offset-2"
                >
                  {isSubmitting ? (
                    <span className="flex items-center justify-center gap-2">
                      <svg className="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4l3-3-3-3v4a8 8 0 00-8 8h4l-3 3 3 3H4z"/>
                      </svg>
                      Creating account…
                    </span>
                  ) : (
                    "Create Account"
                  )}
                </button>
              </form>
            </>
          )}
        </div>

        <p className="mt-6 text-center text-xs text-[var(--text-muted)]">
          By signing up you agree to our{" "}
          <span className="cursor-pointer underline underline-offset-2 hover:text-[var(--brand-500)]">Terms</span>
          {" & "}
          <span className="cursor-pointer underline underline-offset-2 hover:text-[var(--brand-500)]">Privacy Policy</span>.
        </p>
      </main>
    </div>
  );
}
