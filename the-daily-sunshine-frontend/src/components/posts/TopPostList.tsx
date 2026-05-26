import { useTopPosts } from "../../hooks/usePosts"
import type { Post } from "../../types/Post"

export default function TopPostList() {
    const { data, isPending } = useTopPosts() 
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
                </div>
            )}
        </div>
    )
}