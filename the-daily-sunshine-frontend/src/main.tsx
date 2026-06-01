import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { createRouter, RouterProvider } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen.ts'
import { AuthProvider } from './utils/supabase/authContext.tsx'

const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			staleTime: 1000 * 60 * 5,
			retry: 3,
		}
	}
});

const router = createRouter({ 
	routeTree,
	context: {
		queryClient,
	},
	defaultPreload: 'intent',
	defaultPreloadStaleTime: 0,
});

declare module '@tanstack/react-router' {
    interface Register {
        router: typeof router
    }
}

createRoot(document.getElementById('root')!).render(
	<StrictMode>
		<AuthProvider>
			<QueryClientProvider client={queryClient}>
				<RouterProvider router={router} />
				<ReactQueryDevtools initialIsOpen={false}/>
			</QueryClientProvider>
		</AuthProvider>
	</StrictMode>,
)
