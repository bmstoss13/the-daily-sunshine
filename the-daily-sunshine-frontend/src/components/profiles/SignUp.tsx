import { Link } from "@tanstack/react-router";

import { Button } from "../../../components/ui/button"
import type { SignUpError } from "../../types/Errors";
import LoginError from "./LoginError";

interface SignUpProps{
    loading: boolean;
    error: SignUpError | null;
    setEmail: (email: string) => void;
    setPassword: (password: string) => void;
    setConfirmPassword: (confirmPassword: string) => void;
    handleSignUp: (e: React.FormEvent<HTMLFormElement>) => void;
}
export default function SignUp({
    loading, 
    error,
    setEmail,
    setPassword,
    setConfirmPassword,
    handleSignUp,
}: SignUpProps){
    return(
        <div>
            <form
                onSubmit={(e) => handleSignUp(e)}
                className="max-w-md m-auto pt-24"
            >
                <h2 className="font-bold pb-2">
                    Sign Up Today!
                </h2>
                <p>
                    Already have an account? 
                    <Link 
                        to="/sign-in"
                        className="pl-1 underline text-amber-600"
                    >
                        Sign in!
                    </Link>
                </p>
                <div
                    className="flex flex-col gap-4"
                >
                    <input 
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}  
                        type="email"                        
                        placeholder="Email" 
                        className="bg-gray-100 p-2" 
                    />
                    {error === "An email is required" && <LoginError error={error}/>}
                    {error === "Invalid email" && <LoginError error={error}/>}

                    <input 
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
                        type="password" 
                        placeholder="Password"
                        className="bg-gray-100 p-2" 
                    />
                    {error === "A password is required" && <LoginError error={error}/>}

                    <input
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setConfirmPassword(e.target.value)}
                        type="password" 
                        placeholder="Confirm Password"
                        className="bg-gray-100 p-2" 
                    />
                    {error === "Password confirmation is required" && <LoginError error={error}/>}
                    {error === "Passwords must match" && <LoginError error={error}/>}

                    <Button
                        className=""
                        type="submit"
                        disabled={loading}
                    >
                        Sign up
                    </Button>
                </div>
            </form>
        </div>
    )
}