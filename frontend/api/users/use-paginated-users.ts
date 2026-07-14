import { useQuery } from "@tanstack/react-query";

type Props = {
    page: number;
    pageLimit: number;
}

export function usePaginatedUsers({ page, pageLimit }: Props) { 
    const getPaginatedUsersFn = async () => {
        const response = await fetch(`/api/v1/users?page=${page}&limit=${pageLimit}`);
        const data = await response.json();
        return data;
    }

    return useQuery({
        queryKey: ['users', page, pageLimit],
        queryFn: getPaginatedUsersFn,
        placeholderData: true,
    });
}


