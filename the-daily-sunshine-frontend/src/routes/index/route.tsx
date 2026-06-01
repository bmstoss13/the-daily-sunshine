import { createFileRoute } from '@tanstack/react-router'

import TopPostList from '../../components/posts/TopPostList'
import { postQueries } from '../../hooks/usePosts'
import { useQuery } from '@tanstack/react-query'

export const Route = createFileRoute('/')({
	loader: ({ context: { queryClient }}) => {
		return queryClient.ensureQueryData(postQueries.top())
	},
  	component: HomeComponent,
})

function HomeComponent() {
	const { data, isPending } = useQuery(postQueries.top()) 

  	return (
		<div>
			{typeof data !== 'undefined' && data !== null && !isPending ? (
				<TopPostList
					postList={data}
				/>
			) : (
				<div>
					No posts
				</div>
			)}
		</div>
	)
}
