import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import type { Profile } from "../types/Profile";

export function useProfileByID(profileID: string) {
    const {data, isPending} = useQuery({
        queryKey: ['profiles', 'id', profileID],
        queryFn: async () => {
            if (profileID === null || profileID === "") {
                console.error("Profile ID required.")
                return
            }
            try {
                const { data } = await axios.get<Profile>(`api/profiles/id/${profileID}`)
                return data
            } catch (e) {
                if(axios.isAxiosError(e)) {
                    console.error(e.response?.status)
                } else {
                    console.error(`An error occurred while fetching profile by ID ${profileID}: ${e}`)
                }
            }
        }
    })

    return({
        data,
        isPending
    })
}

export function useProfileByUsername(username: string) {
    const {data, isPending} = useQuery({
        queryKey: ['profiles', 'username', username],
        queryFn: async () => {
            if (username === null || username === "") {
                console.error("Username required.")
                return
            }
            try {
                const { data } = await axios.get<Profile>(`api/username/${username}`)
                return data
            } catch (e) {
                if(axios.isAxiosError(e)) {
                    console.error(e.response?.status)
                } else {
                    console.error(`An error occurred while fetching profile by username ${username}: ${e}`)
                }
            }
        }
    })

    return {
        data,
        isPending
    }
}

export function useListProfiles(limit: number, offset: number) {
    const {data, isPending} = useQuery({
        queryKey: ['profiles', 'list', limit, offset],
        queryFn: async () => {
            try {
                const response = await axios.get(`api/profiles/list`, {
                    params: {
                        limit: limit,
                        offset: offset,
                    }
                })
                return response.data.data
            } catch (e) {
                if(axios.isAxiosError(e)) {
                    console.error(e.response?.status)
                } else {
                    console.error("An error occurred while fetching profiles: ", e)
                }                
            }
        }
    })

    return {
        data,
        isPending
    }
}