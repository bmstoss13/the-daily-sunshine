import { mutationOptions, queryOptions } from "@tanstack/react-query";
import axios from "axios";
import type { NewPostRequest, Post, PostStatus } from "../types/Post";
import { useAuth } from "../utils/supabase/authContext";

export const postQueries = {
    all: () => ['posts'],
    top: () => queryOptions({
        queryKey: [...postQueries.all(), 'top'],
        queryFn: fetchTopPosts,
    }),
    list: (limit: number, offset: number) => queryOptions({
        queryKey: [...postQueries.all(), 'list', limit, offset],
        queryFn: () => fetchListPosts(limit, offset),
    }),
    id: (id: string) => queryOptions({
        queryKey: [...postQueries.all(), 'id', id],
        queryFn: () => fetchPostByID(id),
    }),
    slug: (slug: string) => queryOptions({
        queryKey: [...postQueries.all(), 'slug', slug],
        queryFn: () => fetchPostBySlug(slug)
    })
}

// export const postMutations = {
//     create: () => mutationOptions({
//         mutationFn: (newPost: NewPostRequest)
//     })
// }

export async function fetchTopPosts() {
    const response = await axios.get(`api/posts/top`)
    return response.data.data as Post[] 
}

export async function fetchListPosts(limit: number, offset: number) {
    const response = await axios.get(`api/posts/list`, {
        params: { limit, offset }
    })
    return response.data.data as Post[]
}

export async function fetchPostByID(id: string) {
    const response = await axios.get(`api/posts/id/${id}`)
    return response.data.data as Post
}

export async function fetchPostBySlug(slug: string) {
    const response = await axios.get(`api/posts/slug/${slug}`)
    return response.data.data as Post
}

export async function createPost(newPost: NewPostRequest) {
    const { user } = useAuth()
    const response = await axios.post(`api/posts`)
}

