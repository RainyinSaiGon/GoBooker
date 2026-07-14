import { useAuthStore } from "./store/authStore";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001/api/v1";

export interface User {
  id: string;
  name: string;
  email: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface UserInput {
  name: string;
  email: string;
  password?: string;
  role: string;
}

class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

async function apiRequest<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = useAuthStore.getState().token;
  
  const headers = new Headers(options.headers);
  if (!headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    let errorMessage = "An error occurred";
    try {
      const errorData = await response.json();
      errorMessage = errorData.error || errorMessage;
    } catch {
      // ignore
    }
    throw new ApiError(errorMessage, response.status);
  }

  if (response.status === 204) {
    return {} as T;
  }

  return response.json();
}

export const api = {
  login: async (email: string, password: string): Promise<{ token: string }> => {
    return apiRequest<{ token: string }>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  },

  logout: async (): Promise<{ status: string }> => {
    return apiRequest<{ status: string }>("/auth/logout", {
      method: "POST",
    });
  },

  getUsers: async (): Promise<User[]> => {
    return apiRequest<User[]>("/users");
  },

  createUser: async (user: UserInput): Promise<User> => {
    return apiRequest<User>("/users", {
      method: "POST",
      body: JSON.stringify(user),
    });
  },

  updateUser: async (id: string, user: UserInput): Promise<User> => {
    return apiRequest<User>(`/users/${id}`, {
      method: "PUT",
      body: JSON.stringify(user),
    });
  },

  deleteUser: async (id: string): Promise<void> => {
    return apiRequest<void>(`/users/${id}`, {
      method: "DELETE",
    });
  },
};
