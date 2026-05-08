package repository

import (
	"context"
	"fmt"

	"github.com/github.com/bmstoss13/the-daily-sunshine/internal/config"
	"github.com/jackc/pgx/v5"
)

// WithRLS executes a database function securely within an RLS-scoped transaction.
func WithRLS(ctx context.Context, userID string, fn func(tx pgx.Tx) error) error {
	// 1. Begin a new transaction
	tx, err := config.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Automatically rolls back if something fails

	// 2. Set the Postgres role to 'authenticated' (Supabase's default role for logged-in users)
	_, err = tx.Exec(ctx, "SET LOCAL role = 'authenticated';")
	if err != nil {
		return fmt.Errorf("failed to set role: %w", err)
	}

	// 3. Set the auth.uid() by passing the user ID into the JWT claims setting.
	// We use Postgres parameterized queries ($1) to prevent SQL injection.
	_, err = tx.Exec(ctx, "SELECT set_config('request.jwt.claim.sub', $1, true)", userID)
	if err != nil {
		return fmt.Errorf("failed to set user id: %w", err)
	}

	// 4. Execute your actual database query using this scoped transaction
	if err := fn(tx); err != nil {
		return err // The defer will roll it back
	}

	// 5. Commit the transaction
	return tx.Commit(ctx)
}
