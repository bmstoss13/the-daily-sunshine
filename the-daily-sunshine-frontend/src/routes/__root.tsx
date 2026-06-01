import type { QueryClient } from '@tanstack/react-query'
import { Link, Outlet, createRootRouteWithContext } from '@tanstack/react-router'
import { useAuth } from '../utils/supabase/authContext'
import { handleSignOut } from '../hooks/useAuth'
import { useState } from 'react'

export const Route = createRootRouteWithContext<{
    queryClient: QueryClient
}>()({
    component: RootComponent
})

function RootComponent() {
    const user = useAuth()
    const [loading, setLoading] = useState<boolean>(false);

    const signOut = async () => {
        if (!user) return
        setLoading(true)
        try{
            console.log("user ", user, " logging out")
            await handleSignOut()
        } catch (e) {
            console.error("Failed to sign user out: ", e);
        } finally {
            setLoading(false)
        }
    }
    return (
        <>
            <div className='flex gap-4'>
                <Link
                    to="/"
                    activeProps={{
                        className: "font-bold",
                    }}
                    activeOptions={{ exact: true }}
                >
                    Home
                </Link>
                <Link
                    to="/about"
                    activeProps={{
                        className: "font-bold",
                    }}
                >
                    About
                </Link>
                <Link
                    search={{offset: 0}}
                    to="/posts"
                    activeProps={{
                        className: "font-bold",
                    }}
                >
                    Posts 
                </Link>
                <Link
                    to="/profiles"
                    activeProps={{
                        className: "font-bold"
                    }}
                >
                    Profiles 
                </Link>
                {!user.user && (
                    <div className='flex'>                            
                        <Link 
                            to="/sign-in"
                            activeProps={{
                                className: "font-bold"
                            }}
                        >
                            Sign In
                        </Link>
                        <div className='ml-1 pr-1 border-l border-gray-300'/>
                        <Link
                            to="/sign-up"
                            activeProps={{
                                className: "font-bold"
                            }}
                        >
                            Sign Up
                        </Link>
                    </div>
                )}
                {user.user && (
                    <button
                        onClick={signOut}
                        disabled={loading}
                    >
                        Sign Out
                    </button>
                )}

            </div>
            <Outlet />
        </>
    )
}
