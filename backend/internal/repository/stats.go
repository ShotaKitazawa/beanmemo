package repository

import (
	"context"
	"database/sql"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/database/sqlcgen"
)

type StatsRepository struct {
	q *sqlcgen.Queries
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{q: sqlcgen.New(db)}
}

func (r *StatsRepository) StatsByOrigin(ctx context.Context, userID int64) ([]sqlcgen.StatsByOriginRow, error) {
	return r.q.StatsByOrigin(ctx, userID)
}

func (r *StatsRepository) StatsByRoastLevel(ctx context.Context, userID int64) ([]sqlcgen.StatsByRoastLevelRow, error) {
	return r.q.StatsByRoastLevel(ctx, userID)
}

func (r *StatsRepository) StatsByBrewMethod(ctx context.Context, userID int64) ([]sqlcgen.StatsByBrewMethodRow, error) {
	return r.q.StatsByBrewMethod(ctx, userID)
}

func (r *StatsRepository) AvgRatingByOrigin(ctx context.Context, userID int64, origin string) (sql.NullFloat64, error) {
	return r.q.AvgRatingByOrigin(ctx, sqlcgen.AvgRatingByOriginParams{
		UserID: userID,
		Origin: sql.NullString{String: origin, Valid: true},
	})
}

func (r *StatsRepository) AvgRatingByName(ctx context.Context, userID int64, name string) (sql.NullFloat64, error) {
	return r.q.AvgRatingByName(ctx, sqlcgen.AvgRatingByNameParams{
		UserID: userID,
		Name:   name,
	})
}
