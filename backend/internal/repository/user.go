package repository

import (
	"context"
	"database/sql"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/database/sqlcgen"
)

type UserRepository struct {
	q *sqlcgen.Queries
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{q: sqlcgen.New(db)}
}

// UpsertBySub inserts a new user with the given OIDC subject, or updates the name if already exists.
func (r *UserRepository) UpsertBySub(ctx context.Context, sub, name string) error {
	return r.q.UpsertUserBySub(ctx, sqlcgen.UpsertUserBySubParams{
		Sub:  sql.NullString{String: sub, Valid: true},
		Name: name,
	})
}

// GetBySub retrieves the internal user ID for the given OIDC subject.
func (r *UserRepository) GetBySub(ctx context.Context, sub string) (int64, error) {
	user, err := r.q.GetUserBySub(ctx, sql.NullString{String: sub, Valid: true})
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
