import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/posts/$slug')({
  	component: PostComponent,
})

function PostComponent() {
  	return <div>Hello "/posts/$slug"!</div>
}
