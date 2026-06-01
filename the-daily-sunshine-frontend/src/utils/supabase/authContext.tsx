import React, { createContext, useContext, useEffect, useState, type ReactNode } from "react";
import { supabase } from "./supabaseClient";
import type { User } from "@supabase/supabase-js";

interface AuthContextType {
    user: User | null
    loading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider= ({ children }: {children: ReactNode}): React.JSX.Element => {
    const [user, setUser] = useState<User | null>(null)
    const [loading, setLoading] = useState<boolean>(true)

    useEffect(() => {
        supabase.auth.getSession().then(({ data: { session } }) => {
            setUser(session?.user ?? null)
            setLoading(false)
        })

        const { data: { subscription } } = supabase.auth.onAuthStateChange((_event, session) => {
            setUser(session?.user ?? null)
            setLoading(false)
        })

        return () => subscription.unsubscribe()
    }, [])

    return (
        <AuthContext.Provider value={{ user, loading }}>
            {children}  
        </AuthContext.Provider>
    )
}

export const useAuth = (): AuthContextType => {
    const context = useContext(AuthContext)
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider')
    }
    return context
}