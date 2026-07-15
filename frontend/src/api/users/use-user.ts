import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient, User, UserInput } from "../api-client";
import { userQueryKeys } from "./user-query-keys";

export function useCreateUser(options?: { onSuccess?: () => void; onError?: (err: any) => void }) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (user: UserInput): Promise<User> => {
      const response = await apiClient.post<User>("/users", user);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userQueryKeys.all });
      if (options?.onSuccess) options.onSuccess();
    },
    onError: (err: any) => {
      if (options?.onError) options.onError(err);
    },
  });
}

export function useUpdateUser(options?: { onSuccess?: () => void; onError?: (err: any) => void }) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: UserInput }): Promise<User> => {
      const response = await apiClient.put<User>(`/users/${id}`, data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userQueryKeys.all });
      if (options?.onSuccess) options.onSuccess();
    },
    onError: (err: any) => {
      if (options?.onError) options.onError(err);
    },
  });
}

export function useDeleteUser(options?: { onSuccess?: () => void; onError?: (err: any) => void }) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string): Promise<void> => {
      await apiClient.delete(`/users/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userQueryKeys.all });
      if (options?.onSuccess) options.onSuccess();
    },
    onError: (err: any) => {
      if (options?.onError) options.onError(err);
    },
  });
}
