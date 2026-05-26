import { useQuery } from "@tanstack/react-query";
import axios from "axios";

export function useTopPosts() {
    const { data, isPending } = useQuery({
        queryKey: ['posts', 'top'],
        queryFn: async () => {
            try{
                const response = await axios.get(`api/posts/top`)
                return response.data.data
            } catch (e) {
                if(axios.isAxiosError(e)){
                    console.error(e.response?.status)
                }
                console.error("Failed to fetch top posts: ", e)
            }
        },
    })

    return{
        data,
        isPending
    }
} 

export function useListPosts(limit: number, offset: number) {
    const { data, isPending } = useQuery({
        queryKey: ['posts', 'list', limit, offset],
        queryFn: async () => {
            try {
                const response = await axios.get(`api/posts/list`, {
                    params: {
                        limit: limit,
                        offset: offset,
                    }
                })
                return response.data.data
            } catch (e) {
                if (axios.isAxiosError(e)) {
                    console.error(e.response?.status)
                }
                console.error("Failed to fetch list of posts: ", e)
            }
        }
    })
    return {
        data,
        isPending
    }
}
