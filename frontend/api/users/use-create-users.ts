
import { QueryClient, useMutation } from "@tanstack/react-query";
import {User} from "../../app/types";
import {apiClient} from "../api-client"

const createUserFn = async (user: User): Promise<User> => {
    const response = await apiClient.post('/users', user);
    return response.data;
}

export function useCreateUser() {
    const queryClient = new QueryClient();
    return useMutation({
        mutationFn: createUserFn,
        onMutate: async () => {
            await queryClient.cancelQueries({ queryKey: ['users'] });
        }
    }  
    )
}