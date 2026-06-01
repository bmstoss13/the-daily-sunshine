import { createFileRoute, useNavigate } from "@tanstack/react-router";
import SignIn from "../../components/profiles/SignIn";
import { handleSignIn } from "../../hooks/useAuth";
import { useState } from "react";
import type { SignUpError } from "../../types/Errors";

export const Route = createFileRoute('/sign-in')({
    component: SignInComponent,
})

function SignInComponent(){
    const navigate = useNavigate();
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<SignUpError | null>(null);

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        if (!email) {
            setError('An email is required')
            return
        } else if (!password) {
            setError('A password is required')
            return
        }
        
        setLoading(true)
        try{
            const response = await handleSignIn(email, password)
            if(!response.error) {
                navigate({ to: '/', search: {offset: 0} }) //replace with user profile once created
            }
        } catch(e) {
            console.error("Failed to sign up user: ", e)
        } finally {
            setLoading(false)
        }
    }

    return (
        <SignIn
            loading={loading} 
            error={error}
            setEmail={setEmail}
            setPassword={setPassword}
            handleSignIn={handleSubmit}
        />  
    )
}