
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient, Concert, ConcertInput } from "../api-client";
import { concertQueryKeys } from "./concert-query-keys";

export function useCreateConcert(options?: { onSuccess?: () => void; onError?: (err: any) => void }) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (concert: ConcertInput) => {
      const response = await apiClient.post<Concert>("/concerts", { data: concert });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: concertQueryKeys.all });
      options?.onSuccess?.();
    },
    onError: (err) => {
      options?.onError?.(err);
    }
  });
}

export function useUpdateConcert(options?: { onSuccess?: () => void; onError?: (err: any) => void }) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: ConcertInput }) => {
      const response = await apiClient.put<Concert>(`/concerts/${id}`, { data });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: concertQueryKeys.all });
      options?.onSuccess?.();
    },
    onError: (err) => {
      options?.onError?.(err);
    }
  });
}

export function useDeleteConcert(options?: { onSuccess?: () => void; onError?: (err: any) => void }) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      const response = await apiClient.delete(`/concerts/${id}`);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: concertQueryKeys.all });
      options?.onSuccess?.();
    },
    onError: (err) => {
      options?.onError?.(err);
    }
  });
}