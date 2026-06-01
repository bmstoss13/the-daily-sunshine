import { Button } from "../../../components/ui/button"
import type { SignUpError } from "../../types/Errors";
import LoginError from "./LoginError";

interface SignUpProps{
    loading: boolean;
    error: SignUpError | null;
    setEmail: (email: string) => void;
    setPassword: (password: string) => void;
    handleSignIn: (e: React.FormEvent<HTMLFormElement>) => void;
}
export default function SignIn({
    loading, 
    error,
    setEmail,
    setPassword,
    handleSignIn,
}: SignUpProps){
    return(
        <div>
            <form
                onSubmit={(e) => handleSignIn(e)}
                className="max-w-md m-auto pt-24"
            >
                <h2 className="font-bold pb-2">
                    Welcome Back!
                </h2>
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

                    <Button
                        className=""
                        type="submit"
                        disabled={loading}
                    >
                        Sign in
                    </Button>
                </div>
            </form>
        </div>
    )
}