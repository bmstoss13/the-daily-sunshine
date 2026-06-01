import { Link } from "@tanstack/react-router";
import type { Post } from "../../types/Post"

interface PostListProps {
    postList: Post[],
    offset: number,
    handleOffset: (offset: number) => void;
}

export default function PostList({postList, offset, handleOffset}: PostListProps) {
    return (
        <div>
            <div>
                {Array.isArray(postList) ? postList.map((post: Post) => {
                    return (
                        <Link
                            key={post.id}
                            to="/post/$slug"
                            params={{ slug: post.slug }}
                        >
                            <div                                 
                                className="flex flex-col"
                            >
                                <h1 className="font-bold">
                                    {post.title}
                                </h1>
                                <p className="italic">
                                    {post.subtitle}
                                </p>
                            </div>
                        </Link>
                    )
                }) : null}
                <button onClick={() => handleOffset(offset + 20)}>
                    Load more
                </button>
            </div>

        </div>
    )
}