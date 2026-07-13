import {create} from "zustand";


interface User { 
    name: string; 
    email: string; 
    role: string; 
}

interface AuthState { 
  user: User | null
  accessToken: string | null
  setUser: (user: User) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  accessToken: null,
  setUser: (user) => set({ user }),
  logout: () => set({ user: null, accessToken: null }),
}))