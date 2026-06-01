import { Link } from "@tanstack/react-router";
import type { Post } from "../../types/Post"
import PostCard from "./PostCard";

interface TopPostProps{
    postList: Post[];
}
export default function TopPostList({postList }: TopPostProps) {
    
    return (
        <div>
            <div>
                {Array.isArray(postList) ? postList.map((post: Post) => {
                    return (
                        <PostCard
                            key={post.id} 
                            post={post}
                        />
                    )
                }) : null}
            </div>
        </div>
    )
}