import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'

export const Route = createFileRoute('/create')({
    component: CreatePostPage,
})

function CreatePostPage() {
    const [title, setTitle] = useState<string>("New Post")
    const [content, setContent] = useState<string>("")
    
    return <div>Create post</div>
}