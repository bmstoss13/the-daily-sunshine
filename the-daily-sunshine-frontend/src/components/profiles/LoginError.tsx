import type { SignUpError } from "../../types/Errors";

interface LoginErrorProps {
    error: SignUpError;
}
export default function LoginError({error}: LoginErrorProps){
    return (
        <p
            className="text-red-600"
        >
            {error}
        </p>
    )
}