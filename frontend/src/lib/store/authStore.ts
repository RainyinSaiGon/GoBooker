import { create } from "zustand";

interface AuthState {
  token: string | null;
  setToken: (token: string | null) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => {
  // Safe SSR initial token
  const initialToken = typeof window !== "undefined" ? localStorage.getItem("token") : null;

  return {
    token: initialToken,
    setToken: (token) => {
      if (typeof window !== "undefined") {
        if (token) {
          localStorage.setItem("token", token);
        } else {
          localStorage.removeItem("token");
        }
      }
      set({ token });
    },
    logout: () => {
      if (typeof window !== "undefined") {
        localStorage.removeItem("token");
      }
      set({ token: null });
    },
  };
});
