import { createFileRoute } from '@tanstack/react-router'
import TopPostList from '../../components/posts/TopPostList'

export const Route = createFileRoute('/')({
  	component: HomeComponent,
	notFoundComponent: () => <div>404 Not Found</div>
})

function HomeComponent() {
  	return (
		<div>
			<TopPostList/>
		</div>
	)
}
