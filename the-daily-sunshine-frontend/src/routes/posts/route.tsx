import { createFileRoute } from '@tanstack/react-router'
import PostList from '../../components/posts/PostList'

export const Route = createFileRoute('/posts')({
  	component: PostsComponent,
})

function PostsComponent() {
  	return (
		<div>
			<PostList/>
		</div>
	) 
	
}
