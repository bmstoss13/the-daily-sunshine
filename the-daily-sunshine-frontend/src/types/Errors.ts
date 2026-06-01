const SignUpErrors = {
    NoEmail: "An email is required",
    NoPassword: "A password is required",
    NoConfirmPassword: "Password confirmation is required",
    InvalidEmail: "Invalid email",
    PasswordMismatch: "Passwords must match",
} as const;

export type SignUpError = typeof SignUpErrors[keyof typeof SignUpErrors]