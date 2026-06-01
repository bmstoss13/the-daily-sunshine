import { createFileRoute, useParams } from '@tanstack/react-router'

import { postQueries } from '../../hooks/usePosts'
import { useQuery } from '@tanstack/react-query'

export const Route = createFileRoute('/post/$slug')({
	loader: ({ context: { queryClient }, params}) => {
		return queryClient.ensureQueryData(postQueries.slug(params.slug))
	},
  	component: PostComponent,
})

function PostComponent() {
	const { slug } = useParams({ from: '/post/$slug' })
	const { data, isPending } = useQuery(postQueries.slug(slug))
  	return (
		<div>
			{typeof data !== 'undefined' && data !== null && !isPending && (
				data.title
			)}
		</div>
	)
}
