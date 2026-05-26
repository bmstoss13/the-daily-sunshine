import { useState } from "react"
import { useListPosts } from "../../hooks/usePosts"
import type { Post } from "../../types/Post"

export default function PostList() {
    const [offset, setOffset] = useState<number>(0)
    const { data, isPending } = useListPosts(20, offset) 
    return (
        <div>
            {data !== null && !isPending && (
                <div>
                    {Array.isArray(data) ? data.map((post: Post) => {
                        return (
                            <div 
                                key={post.id}
                                className="flex flex-col"
                            >
                                <h1 className="font-bold">
                                    {post.title}
                                </h1>
                                <p className="italic">
                                    {post.subtitle}
                                </p>
                            </div>
                        )
                    }) : null}
                    <button onClick={() => setOffset(offset + 20)}>
                        Load more
                    </button>
                </div>
            )}
        </div>
    )
}