import { createFileRoute, useNavigate } from '@tanstack/react-router'

import PostList from '../../components/posts/PostList'
import { postQueries } from '../../hooks/usePosts'
import { useQuery } from '@tanstack/react-query'

// url search parameters
type PostsURLSearch = {
	offset: number,
}

export const Route = createFileRoute('/posts')({
	validateSearch: (search: Record<string, unknown>): PostsURLSearch => {
		return {
			offset: Number(search?.offset ?? 0),
		}
	},
	loaderDeps: ({ search: { offset } }) => ({ offset }),
	loader: ({ context: { queryClient }, deps: { offset }}) => {
		return queryClient.ensureQueryData(postQueries.list(20, offset))
	},
  	component: PostsComponent,
})

function PostsComponent() {
	const { offset } = Route.useSearch()
	const navigate = useNavigate({ from: Route.id })
	const { data, isPending } = useQuery(postQueries.list(20, offset))
	const handleOffset = (newOffset: number) => {
		navigate({
			search: (prev) => ({...prev, offset: newOffset}),
			resetScroll: false,
		})
	}
  	return (
		<div>
			{data && !isPending && (
				<PostList
					postList={data}
					offset={offset}
					handleOffset={handleOffset}
				/>
			)}
			
		</div>
	) 
	
}
