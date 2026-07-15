import { useQuery } from "@tanstack/react-query";
import { apiClient, User, PaginatedUsersResponse } from "../api-client";
import { userQueryKeys } from "./user-query-keys";
import { useAuthStore } from "@/lib/store/authStore";

const getUsersFn = async (): Promise<User[]> => {
  const response = await apiClient.get<PaginatedUsersResponse>("/users");
  return response.data.users;
};

export function useUsers() {
  const { token } = useAuthStore();
  return useQuery<User[]>({
    queryKey: userQueryKeys.all,
    queryFn: getUsersFn,
    enabled: !!token,
  });
}
