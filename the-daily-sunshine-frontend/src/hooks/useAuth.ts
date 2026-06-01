import type { Provider, AuthResponse, OAuthResponse } from "@supabase/supabase-js";
import { supabase } from "../utils/supabase/supabaseClient";

export async function handleSignUp(email: string, password: string): Promise<AuthResponse> {
    const response = await supabase.auth.signUp({ email, password })
    if (response.error) console.error('Sign up error:', response.error.message)
    return response
}

export async function handleSignIn(email: string, password: string): Promise<AuthResponse> {
    const response = await supabase.auth.signInWithPassword({ email, password })
    if (response.error) console.error('Sign in error:', response.error.message)
    return response
}

export async function handleOAuthSignIn(provider: Provider): Promise<OAuthResponse> {
    const response = await supabase.auth.signInWithOAuth({
        provider,
        options: {
            redirectTo: window.location.origin,
        },
    })
    if (response.error) console.error('OAuth error:', response.error.message)
    return response
}

export async function handleSignOut(): Promise<void> {
    const { error } = await supabase.auth.signOut()
    if (error) console.error('Sign out error:', error.message)
}