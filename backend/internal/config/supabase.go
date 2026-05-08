package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

var (
	supabaseURL     string
	supabaseAnonKey string
)

func InitEnv() error {
	_ = godotenv.Load()

	supabaseURL = os.Getenv("PUBLIC_SUPABASE_URL")
	supabaseAnonKey = os.Getenv("PUBLIC_SUPABASE_ANON_KEY")

	if supabaseURL == "" || supabaseAnonKey == "" {
		return fmt.Errorf("missing Supabase environment variables")
	}

	return nil
}

func NewUserClient(userJWT string) (*supabase.Client, error) {
	client, err := supabase.NewClient(supabaseURL, supabaseAnonKey, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize client: %w", err)
	}

	// override auth header with specific user's jwt
	// tells supabase exactly who is making the request so RLS works
	client.Auth.WithToken(userJWT)

	return client, nil
}

// NewAdminClient uses the Service Role key ONLY for specific background tasks
// that strictly require bypassing RLS (e.g., a cron job syncing data).
func NewAdminClient() (*supabase.Client, error) {
	serviceKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	return supabase.NewClient(supabaseURL, serviceKey, nil)
}
