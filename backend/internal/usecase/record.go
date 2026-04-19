package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/api"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/database/sqlcgen"
)

type recordRepository interface {
	ListAll(ctx context.Context, userID int64) ([]sqlcgen.Record, error)
	ListByOrigin(ctx context.Context, userID int64, origin string) ([]sqlcgen.Record, error)
	ListByRoastLevel(ctx context.Context, userID int64, roastLevel string) ([]sqlcgen.Record, error)
	ListByRatingMin(ctx context.Context, userID int64, ratingMin int64) ([]sqlcgen.Record, error)
	ListByBrewMethod(ctx context.Context, userID int64, brewMethod string) ([]sqlcgen.Record, error)
	Get(ctx context.Context, id, userID int64) (sqlcgen.Record, error)
	GetRelated(ctx context.Context, userID, id int64, name string) ([]sqlcgen.Record, error)
	Create(ctx context.Context, arg sqlcgen.CreateRecordParams) (int64, error)
	Update(ctx context.Context, arg sqlcgen.UpdateRecordParams) error
	Delete(ctx context.Context, id, userID int64) error
}

type RecordUsecase struct {
	repo recordRepository
}

func NewRecordUsecase(repo recordRepository) *RecordUsecase {
	return &RecordUsecase{repo: repo}
}

func (u *RecordUsecase) List(ctx context.Context, userID int64, params api.ListRecordsParams) ([]sqlcgen.Record, error) {
	if params.Origin.Set {
		records, err := u.repo.ListByOrigin(ctx, userID, params.Origin.Value)
		if err != nil {
			return nil, err
		}
		return filterRecords(records, params), nil
	}
	if params.RoastLevel.Set {
		records, err := u.repo.ListByRoastLevel(ctx, userID, string(params.RoastLevel.Value))
		if err != nil {
			return nil, err
		}
		return filterRecords(records, params), nil
	}
	if params.RatingMin.Set {
		records, err := u.repo.ListByRatingMin(ctx, userID, int64(params.RatingMin.Value))
		if err != nil {
			return nil, err
		}
		return filterRecords(records, params), nil
	}
	if params.BrewMethod.Set {
		records, err := u.repo.ListByBrewMethod(ctx, userID, string(params.BrewMethod.Value))
		if err != nil {
			return nil, err
		}
		return filterRecords(records, params), nil
	}
	return u.repo.ListAll(ctx, userID)
}

func filterRecords(records []sqlcgen.Record, params api.ListRecordsParams) []sqlcgen.Record {
	var out []sqlcgen.Record
	for _, r := range records {
		if params.Origin.Set && (!r.Origin.Valid || r.Origin.String != params.Origin.Value) {
			continue
		}
		if params.RoastLevel.Set && (!r.RoastLevel.Valid || r.RoastLevel.String != string(params.RoastLevel.Value)) {
			continue
		}
		if params.RatingMin.Set && int(r.Rating) < params.RatingMin.Value {
			continue
		}
		if params.BrewMethod.Set && (!r.BrewMethod.Valid || r.BrewMethod.String != string(params.BrewMethod.Value)) {
			continue
		}
		out = append(out, r)
	}
	return out
}

func (u *RecordUsecase) Get(ctx context.Context, userID, id int64) (sqlcgen.Record, []sqlcgen.Record, error) {
	record, err := u.repo.Get(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlcgen.Record{}, nil, ErrNotFound
		}
		return sqlcgen.Record{}, nil, err
	}
	related, err := u.repo.GetRelated(ctx, userID, id, record.Name)
	if err != nil {
		return sqlcgen.Record{}, nil, err
	}
	return record, related, nil
}

func (u *RecordUsecase) Create(ctx context.Context, userID int64, req *api.CreateRecordRequest) (sqlcgen.Record, error) {
	isNoteFilled := req.TastingNote.IsSet() && !req.TastingNote.Null && req.TastingNote.Value != ""

	params := sqlcgen.CreateRecordParams{
		UserID:       userID,
		Name:         req.Name,
		Rating:       int64(req.Rating.Or(0)),
		IsNoteFilled: isNoteFilled,
	}
	if req.Origin.Set && !req.Origin.Null {
		params.Origin = sql.NullString{String: req.Origin.Value, Valid: true}
	}
	if req.RoastLevel.Set && !req.RoastLevel.Null {
		params.RoastLevel = sql.NullString{String: string(req.RoastLevel.Value), Valid: true}
	}
	if req.Shop.Set && !req.Shop.Null {
		params.Shop = sql.NullString{String: req.Shop.Value, Valid: true}
	}
	if req.Price.Set && !req.Price.Null {
		params.Price = sql.NullInt64{Int64: int64(req.Price.Value), Valid: true}
	}
	if req.PurchasedAt.Set && !req.PurchasedAt.Null {
		params.PurchasedAt = sql.NullTime{Time: req.PurchasedAt.Value, Valid: true}
	}
	if req.TastingNote.Set && !req.TastingNote.Null {
		params.TastingNote = sql.NullString{String: req.TastingNote.Value, Valid: true}
	}
	if req.BrewMethod.Set && !req.BrewMethod.Null {
		params.BrewMethod = sql.NullString{String: string(req.BrewMethod.Value), Valid: true}
	}
	if req.Recipe.Set && !req.Recipe.Null {
		params.Recipe = sql.NullString{String: req.Recipe.Value, Valid: true}
	}

	id, err := u.repo.Create(ctx, params)
	if err != nil {
		return sqlcgen.Record{}, err
	}
	return u.repo.Get(ctx, id, userID)
}

func (u *RecordUsecase) Update(ctx context.Context, userID, id int64, req *api.UpdateRecordRequest) (sqlcgen.Record, error) {
	existing, err := u.repo.Get(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlcgen.Record{}, ErrNotFound
		}
		return sqlcgen.Record{}, err
	}

	updated := applyUpdate(existing, req)
	updated.IsNoteFilled = updated.TastingNote.Valid && updated.TastingNote.String != ""

	if err := u.repo.Update(ctx, sqlcgen.UpdateRecordParams{
		ID:           id,
		UserID:       userID,
		Name:         updated.Name,
		Rating:       updated.Rating,
		Origin:       updated.Origin,
		RoastLevel:   updated.RoastLevel,
		Shop:         updated.Shop,
		Price:        updated.Price,
		PurchasedAt:  updated.PurchasedAt,
		TastingNote:  updated.TastingNote,
		BrewMethod:   updated.BrewMethod,
		Recipe:       updated.Recipe,
		IsNoteFilled: updated.IsNoteFilled,
	}); err != nil {
		return sqlcgen.Record{}, err
	}
	return u.repo.Get(ctx, id, userID)
}

func applyUpdate(existing sqlcgen.Record, req *api.UpdateRecordRequest) sqlcgen.Record {
	r := existing
	if req.Name.Set {
		r.Name = req.Name.Value
	}
	if req.Rating.Set {
		r.Rating = int64(req.Rating.Value)
	}
	if req.Origin.Set {
		if req.Origin.Null {
			r.Origin = sql.NullString{}
		} else {
			r.Origin = sql.NullString{String: req.Origin.Value, Valid: true}
		}
	}
	if req.RoastLevel.Set {
		if req.RoastLevel.Null {
			r.RoastLevel = sql.NullString{}
		} else {
			r.RoastLevel = sql.NullString{String: string(req.RoastLevel.Value), Valid: true}
		}
	}
	if req.Shop.Set {
		if req.Shop.Null {
			r.Shop = sql.NullString{}
		} else {
			r.Shop = sql.NullString{String: req.Shop.Value, Valid: true}
		}
	}
	if req.Price.Set {
		if req.Price.Null {
			r.Price = sql.NullInt64{}
		} else {
			r.Price = sql.NullInt64{Int64: int64(req.Price.Value), Valid: true}
		}
	}
	if req.PurchasedAt.Set {
		if req.PurchasedAt.Null {
			r.PurchasedAt = sql.NullTime{}
		} else {
			r.PurchasedAt = sql.NullTime{Time: req.PurchasedAt.Value, Valid: true}
		}
	}
	if req.TastingNote.Set {
		if req.TastingNote.Null {
			r.TastingNote = sql.NullString{}
		} else {
			r.TastingNote = sql.NullString{String: req.TastingNote.Value, Valid: true}
		}
	}
	if req.BrewMethod.Set {
		if req.BrewMethod.Null {
			r.BrewMethod = sql.NullString{}
		} else {
			r.BrewMethod = sql.NullString{String: string(req.BrewMethod.Value), Valid: true}
		}
	}
	if req.Recipe.Set {
		if req.Recipe.Null {
			r.Recipe = sql.NullString{}
		} else {
			r.Recipe = sql.NullString{String: req.Recipe.Value, Valid: true}
		}
	}
	r.UpdatedAt = time.Now()
	return r
}

func (u *RecordUsecase) Delete(ctx context.Context, userID, id int64) error {
	_, err := u.repo.Get(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return u.repo.Delete(ctx, id, userID)
}
