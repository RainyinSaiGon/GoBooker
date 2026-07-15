import axios from "axios";
import { useAuthStore } from "@/lib/store/authStore";

const API_BASE_URL =  "http://localhost:3001/api/v1";

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

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
}

export interface PaginatedUsersResponse {
  users: User[];
  total: number;
  page: number;
  size: number;
  totalPages: number;
}


