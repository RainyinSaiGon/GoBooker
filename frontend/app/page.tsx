"use client";

const apiUrl = process.env.NEXT_PUBLIC_API_URL || "";

import { useState } from "react";

interface SignUpRequest {
  name: string;
  email: string;
  password: string;
}

interface FormErrors {
  name?: string;
  email?: string;
  password?: string;
}

export default function SignUpPage() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errors, setErrors] = useState<FormErrors>({});
  const [touched, setTouched] = useState<Record<string, boolean>>({});

  // Validation rules 
  const validate = (values: SignUpRequest): FormErrors => {
    const newErrors: FormErrors = {};

    if (!values.name.trim()) {
      newErrors.name = "Full name is required";
    }

    if (!values.email.trim()) {
      newErrors.email = "Email is required";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(values.email)) {
      newErrors.email = "Please enter a valid email address";
    }

    if (!values.password) {
      newErrors.password = "Password is required";
    } else if (values.password.length < 8) {
      newErrors.password = "Password must be at least 8 characters";
    }

    return newErrors;
  };

  const handleBlur = (field: keyof FormErrors) => {
    setTouched((prev) => ({ ...prev, [field]: true }));
    const newErrors = validate({ name, email, password });
    setErrors(newErrors);
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const newUser: SignUpRequest = { name, email, password };
    const newErrors = validate(newUser);
    setErrors(newErrors);
    setTouched({ name: true, email: true, password: true });

    // Stop here if there are validation errors — don't hit the API
    if (Object.keys(newErrors).length > 0) {
      return;
    }

    setIsSubmitting(true);
    try {
      const signUpResponse = await fetch(`${apiUrl}/users`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(newUser),
      });

      if (!signUpResponse.ok) {
        const errorData = await signUpResponse.json().catch(() => ({}));
        // Surface backend validation errors too (e.g. "email already exists")
        throw new Error(errorData.message || "Failed to create user");
      }

      setName("");
      setEmail("");
      setPassword("");
      setErrors({});
      setTouched({});
      alert("User created successfully!");
    } catch (error) {
      console.error("Error creating user:", error);
      setErrors((prev) => ({ ...prev, email: (error as Error).message }));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-zinc-50 font-sans dark:bg-zinc-950 p-6 sm:p-12 text-zinc-900 dark:text-zinc-50">
      <main className="w-full max-w-xl mx-auto bg-white dark:bg-zinc-900 shadow-sm rounded-xl border border-zinc-200 dark:border-zinc-800 p-6 sm:p-8 space-y-8">
        <div>
          <h1 className="text-xl text-center font-semibold tracking-tight">Sign Up</h1>
        </div>

        <form className="space-y-4" onSubmit={handleSubmit} noValidate>
          {/* Name field */}
          <div>
            <label className="block text-xs font-medium mb-1 text-zinc-600 dark:text-zinc-400">
              Full Name <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              onBlur={() => handleBlur("name")}
              placeholder="John Doe"
              className={`w-full px-3 py-2 text-sm bg-zinc-50 dark:bg-zinc-800 border rounded-lg focus:outline-none focus:ring-2 ${
                touched.name && errors.name
                  ? "border-red-500 focus:ring-red-500"
                  : "border-zinc-200 dark:border-zinc-700 focus:ring-blue-500"
              }`}
              disabled={isSubmitting}
            />
            {touched.name && errors.name && (
              <p className="mt-1 text-xs text-red-500">{errors.name}</p>
            )}
          </div>

          {/* Email field */}
          <div>
            <label className="block text-xs font-medium mb-1 text-zinc-600 dark:text-zinc-400">
              Email Address <span className="text-red-500">*</span>
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              onBlur={() => handleBlur("email")}
              placeholder="john@example.com"
              className={`w-full px-3 py-2 text-sm bg-zinc-50 dark:bg-zinc-800 border rounded-lg focus:outline-none focus:ring-2 ${
                touched.email && errors.email
                  ? "border-red-500 focus:ring-red-500"
                  : "border-zinc-200 dark:border-zinc-700 focus:ring-blue-500"
              }`}
              disabled={isSubmitting}
            />
            {touched.email && errors.email && (
              <p className="mt-1 text-xs text-red-500">{errors.email}</p>
            )}
          </div>

          {/* Password field */}
          <div>
            <label className="block text-xs font-medium mb-1 text-zinc-600 dark:text-zinc-400">
              Password <span className="text-red-500">*</span>
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              onBlur={() => handleBlur("password")}
              placeholder="••••••••"
              className={`w-full px-3 py-2 text-sm bg-zinc-50 dark:bg-zinc-800 border rounded-lg focus:outline-none focus:ring-2 ${
                touched.password && errors.password
                  ? "border-red-500 focus:ring-red-500"
                  : "border-zinc-200 dark:border-zinc-700 focus:ring-blue-500"
              }`}
              disabled={isSubmitting}
            />
            {touched.password && errors.password ? (
              <p className="mt-1 text-xs text-red-500">{errors.password}</p>
            ) : (
              <p className="mt-1 text-xs text-zinc-400">Must be at least 8 characters</p>
            )}
          </div>

          <div className="flex gap-2 pt-2">
            <button
              type="submit"
              disabled={isSubmitting}
              className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white text-sm font-medium py-2 px-4 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-zinc-900"
            >
              {isSubmitting ? "Adding..." : "Sign Up"}
            </button>
          </div>
        </form>
      </main>
    </div>
  );
}