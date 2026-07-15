import { useQuery } from "@tanstack/react-query";
import { apiClient, Concert, PaginatedConcertsResponse } from "../api-client";
import { concertQueryKeys } from "./concert-query-keys";
import { useAuthStore } from "@/lib/store/authStore";


// TODO: Review the pagination logic and adjust the query function accordingly if needed. For now, it fetches all concerts without pagination.
const getConcertsFn = async (): Promise<Concert[]> => { 
    const response = await apiClient.get<PaginatedConcertsResponse>("/concerts");
    return response.data.concerts;
}


export function useConcerts() {
    const { token } = useAuthStore();
    return useQuery<Concert[]>({
        queryKey: concertQueryKeys.all,
        queryFn: getConcertsFn,
        enabled: !!token,
    });
}