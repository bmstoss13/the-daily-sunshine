import { createFileRoute, useNavigate } from "@tanstack/react-router";

import SignUp from "../../components/profiles/SignUp";
import { useState } from "react";
import { handleSignUp } from "../../hooks/useAuth";
import type { SignUpError } from "../../types/Errors";

export const Route = createFileRoute('/sign-up')({
    component: SignUpComponent,
})

function SignUpComponent(){
    const navigate = useNavigate();
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [confirmPassword, setConfirmPassword] = useState<string>('');
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
        } else if (!confirmPassword) {
            setError('Password confirmation is required')
            return
        } else if (confirmPassword !== password) {
            setError('Passwords must match')
            return
        }
        
        setLoading(true)
        try{
            const response = await handleSignUp(email, password)
            if(!response.error) {
                navigate({ to: '/posts', search: {offset: 0} }) //replace with user profile once created
            }
        } catch(e) {
            console.error("Failed to sign up user: ", e)
        } finally {
            setLoading(false)
        }
    }

    return (
        <SignUp 
            loading={loading} 
            error={error}
            setEmail={setEmail}
            setPassword={setPassword}
            setConfirmPassword={setConfirmPassword}
            handleSignUp={handleSubmit}
        />  
    )
}