package repository

import (
	"context"
	"database/sql"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/database/sqlcgen"
)

type RecordRepository struct {
	q *sqlcgen.Queries
}

func NewRecordRepository(db *sql.DB) *RecordRepository {
	return &RecordRepository{q: sqlcgen.New(db)}
}

func (r *RecordRepository) ListAll(ctx context.Context, userID int64) ([]sqlcgen.Record, error) {
	return r.q.ListRecords(ctx, userID)
}

func (r *RecordRepository) ListByOrigin(ctx context.Context, userID int64, origin string) ([]sqlcgen.Record, error) {
	return r.q.ListRecordsByOrigin(ctx, sqlcgen.ListRecordsByOriginParams{
		UserID: userID,
		Origin: sql.NullString{String: origin, Valid: true},
	})
}

func (r *RecordRepository) ListByRoastLevel(ctx context.Context, userID int64, roastLevel string) ([]sqlcgen.Record, error) {
	return r.q.ListRecordsByRoastLevel(ctx, sqlcgen.ListRecordsByRoastLevelParams{
		UserID:     userID,
		RoastLevel: sql.NullString{String: roastLevel, Valid: true},
	})
}

func (r *RecordRepository) ListByRatingMin(ctx context.Context, userID int64, ratingMin int8) ([]sqlcgen.Record, error) {
	return r.q.ListRecordsByRatingMin(ctx, sqlcgen.ListRecordsByRatingMinParams{
		UserID: userID,
		Rating: ratingMin,
	})
}

func (r *RecordRepository) ListByBrewMethod(ctx context.Context, userID int64, brewMethod string) ([]sqlcgen.Record, error) {
	return r.q.ListRecordsByBrewMethod(ctx, sqlcgen.ListRecordsByBrewMethodParams{
		UserID:     userID,
		BrewMethod: sql.NullString{String: brewMethod, Valid: true},
	})
}

func (r *RecordRepository) Get(ctx context.Context, id, userID int64) (sqlcgen.Record, error) {
	return r.q.GetRecord(ctx, sqlcgen.GetRecordParams{ID: id, UserID: userID})
}

func (r *RecordRepository) GetRelated(ctx context.Context, userID, id int64, name string) ([]sqlcgen.Record, error) {
	return r.q.GetRelatedRecords(ctx, sqlcgen.GetRelatedRecordsParams{
		UserID: userID,
		ID:     id,
		Name:   name,
	})
}

func (r *RecordRepository) Create(ctx context.Context, arg sqlcgen.CreateRecordParams) (int64, error) {
	return r.q.CreateRecord(ctx, arg)
}

func (r *RecordRepository) Update(ctx context.Context, arg sqlcgen.UpdateRecordParams) error {
	return r.q.UpdateRecord(ctx, arg)
}

func (r *RecordRepository) Delete(ctx context.Context, id, userID int64) error {
	return r.q.DeleteRecord(ctx, sqlcgen.DeleteRecordParams{ID: id, UserID: userID})
}

func (r *RecordRepository) Count(ctx context.Context, userID int64) (int64, error) {
	return r.q.CountRecords(ctx, userID)
}

func (r *RecordRepository) ListAllTastingNotes(ctx context.Context, userID int64) ([]sql.NullString, error) {
	return r.q.ListAllTastingNotes(ctx, userID)
}
