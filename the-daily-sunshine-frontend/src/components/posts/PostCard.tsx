import { Link } from "@tanstack/react-router";

import type { Post } from "../../types/Post";

interface PostCardProps{
    post: Post
}

export default function PostCard ({
    post,
}: PostCardProps){
    return(
        <div                                 
            className="flex flex-col max-w-140 p-2"
        >
            <Link
                to="/post/$slug"
                params={{ slug: post.slug }}
            >
                <img 
                    src={post.CoverImage?.image_url}
                />
                <h1 className="font-bold hover:underline">
                    {post.title}
                </h1>
                <p className="italic">
                    {post.subtitle}
                </p>                
            </Link>
        </div>        
    )
}