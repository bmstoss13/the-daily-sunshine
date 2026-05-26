import { Link, Outlet, createRootRoute } from '@tanstack/react-router'

export const Route = createRootRoute({
    component: RootComponent,
})

function RootComponent() {
    return (
        <>
            <div>
                <Link
                    to="/"
                    activeProps={{
                        className: "font-bold",
                    }}
                    activeOptions={{ exact: true }}
                >
                    Home
                </Link>{' '}
                <Link
                    to="/about"
                    activeProps={{
                        className: "font-bold",
                    }}
                >
                    About
                </Link>{' '}
                <Link
                    to="/posts"
                    activeProps={{
                        className: "font-bold",
                    }}
                >
                    Posts 
                </Link>{' '}
                <Link
                    to="/profiles"
                    activeProps={{
                        className: "font-bold"
                    }}
                >
                    Profiles 
                </Link>
            </div>
            <Outlet />
        </>
    )
}
