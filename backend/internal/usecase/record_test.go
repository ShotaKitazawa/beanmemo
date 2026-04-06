package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/api"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/database/sqlcgen"
)

// --- mock ---

type stubRecordRepo struct {
	listAllFn      func(context.Context, int64) ([]sqlcgen.Record, error)
	listByOriginFn func(context.Context, int64, string) ([]sqlcgen.Record, error)
	listByRoastFn  func(context.Context, int64, string) ([]sqlcgen.Record, error)
	listByRatingFn func(context.Context, int64, int8) ([]sqlcgen.Record, error)
	listByBrewFn   func(context.Context, int64, string) ([]sqlcgen.Record, error)
	getFn          func(context.Context, int64, int64) (sqlcgen.Record, error)
	getRelatedFn   func(context.Context, int64, int64, string) ([]sqlcgen.Record, error)
	createFn       func(context.Context, sqlcgen.CreateRecordParams) (int64, error)
	updateFn       func(context.Context, sqlcgen.UpdateRecordParams) error
	deleteFn       func(context.Context, int64, int64) error
}

func (s *stubRecordRepo) ListAll(ctx context.Context, userID int64) ([]sqlcgen.Record, error) {
	if s.listAllFn != nil {
		return s.listAllFn(ctx, userID)
	}
	return nil, nil
}

func (s *stubRecordRepo) ListByOrigin(ctx context.Context, userID int64, origin string) ([]sqlcgen.Record, error) {
	if s.listByOriginFn != nil {
		return s.listByOriginFn(ctx, userID, origin)
	}
	return nil, nil
}

func (s *stubRecordRepo) ListByRoastLevel(ctx context.Context, userID int64, roastLevel string) ([]sqlcgen.Record, error) {
	if s.listByRoastFn != nil {
		return s.listByRoastFn(ctx, userID, roastLevel)
	}
	return nil, nil
}

func (s *stubRecordRepo) ListByRatingMin(ctx context.Context, userID int64, ratingMin int8) ([]sqlcgen.Record, error) {
	if s.listByRatingFn != nil {
		return s.listByRatingFn(ctx, userID, ratingMin)
	}
	return nil, nil
}

func (s *stubRecordRepo) ListByBrewMethod(ctx context.Context, userID int64, brewMethod string) ([]sqlcgen.Record, error) {
	if s.listByBrewFn != nil {
		return s.listByBrewFn(ctx, userID, brewMethod)
	}
	return nil, nil
}

func (s *stubRecordRepo) Get(ctx context.Context, id, userID int64) (sqlcgen.Record, error) {
	if s.getFn != nil {
		return s.getFn(ctx, id, userID)
	}
	return sqlcgen.Record{}, nil
}

func (s *stubRecordRepo) GetRelated(ctx context.Context, userID, id int64, name string) ([]sqlcgen.Record, error) {
	if s.getRelatedFn != nil {
		return s.getRelatedFn(ctx, userID, id, name)
	}
	return nil, nil
}

func (s *stubRecordRepo) Create(ctx context.Context, arg sqlcgen.CreateRecordParams) (int64, error) {
	if s.createFn != nil {
		return s.createFn(ctx, arg)
	}
	return 0, nil
}

func (s *stubRecordRepo) Update(ctx context.Context, arg sqlcgen.UpdateRecordParams) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, arg)
	}
	return nil
}

func (s *stubRecordRepo) Delete(ctx context.Context, id, userID int64) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id, userID)
	}
	return nil
}

// --- helpers ---

func makeRecord(id int64, name string, rating int8) sqlcgen.Record {
	return sqlcgen.Record{
		ID:        id,
		UserID:    defaultUserID,
		Name:      name,
		Rating:    rating,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// --- filterRecords ---

func TestFilterRecords_NoFilter(t *testing.T) {
	records := []sqlcgen.Record{makeRecord(1, "A", 3), makeRecord(2, "B", 5)}
	got := filterRecords(records, api.ListRecordsParams{})
	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}
}

func TestFilterRecords_ByOrigin(t *testing.T) {
	r1 := makeRecord(1, "A", 3)
	r1.Origin = sql.NullString{String: "Ethiopia", Valid: true}
	r2 := makeRecord(2, "B", 4)
	r2.Origin = sql.NullString{String: "Brazil", Valid: true}
	r3 := makeRecord(3, "C", 5) // no origin

	params := api.ListRecordsParams{
		Origin: api.OptString{Set: true, Value: "Ethiopia"},
	}
	got := filterRecords([]sqlcgen.Record{r1, r2, r3}, params)
	if len(got) != 1 || got[0].ID != 1 {
		t.Fatalf("expected only record 1, got %v", got)
	}
}

func TestFilterRecords_ByRoastLevel(t *testing.T) {
	r1 := makeRecord(1, "A", 3)
	r1.RoastLevel = sql.NullString{String: "light", Valid: true}
	r2 := makeRecord(2, "B", 4)
	r2.RoastLevel = sql.NullString{String: "dark", Valid: true}

	params := api.ListRecordsParams{
		RoastLevel: api.OptListRecordsRoastLevel{Set: true, Value: "light"},
	}
	got := filterRecords([]sqlcgen.Record{r1, r2}, params)
	if len(got) != 1 || got[0].ID != 1 {
		t.Fatalf("expected only record 1, got %v", got)
	}
}

func TestFilterRecords_ByRatingMin(t *testing.T) {
	r1 := makeRecord(1, "A", 3)
	r2 := makeRecord(2, "B", 4)
	r3 := makeRecord(3, "C", 5)

	params := api.ListRecordsParams{
		RatingMin: api.OptInt{Set: true, Value: 4},
	}
	got := filterRecords([]sqlcgen.Record{r1, r2, r3}, params)
	if len(got) != 2 {
		t.Fatalf("expected 2 records with rating >= 4, got %d", len(got))
	}
}

func TestFilterRecords_ByBrewMethod(t *testing.T) {
	r1 := makeRecord(1, "A", 3)
	r1.BrewMethod = sql.NullString{String: "pour_over", Valid: true}
	r2 := makeRecord(2, "B", 4)
	r2.BrewMethod = sql.NullString{String: "espresso", Valid: true}

	params := api.ListRecordsParams{
		BrewMethod: api.OptListRecordsBrewMethod{Set: true, Value: "pour_over"},
	}
	got := filterRecords([]sqlcgen.Record{r1, r2}, params)
	if len(got) != 1 || got[0].ID != 1 {
		t.Fatalf("expected only record 1, got %v", got)
	}
}

func TestFilterRecords_MultipleFilters(t *testing.T) {
	r1 := makeRecord(1, "A", 5)
	r1.Origin = sql.NullString{String: "Ethiopia", Valid: true}
	r1.RoastLevel = sql.NullString{String: "light", Valid: true}

	r2 := makeRecord(2, "B", 5)
	r2.Origin = sql.NullString{String: "Ethiopia", Valid: true}
	r2.RoastLevel = sql.NullString{String: "dark", Valid: true}

	params := api.ListRecordsParams{
		Origin:     api.OptString{Set: true, Value: "Ethiopia"},
		RoastLevel: api.OptListRecordsRoastLevel{Set: true, Value: "light"},
	}
	got := filterRecords([]sqlcgen.Record{r1, r2}, params)
	if len(got) != 1 || got[0].ID != 1 {
		t.Fatalf("expected only record 1, got %v", got)
	}
}

// --- applyUpdate ---

func TestApplyUpdate_Name(t *testing.T) {
	existing := makeRecord(1, "Old Name", 3)
	req := &api.UpdateRecordRequest{
		Name: api.OptString{Set: true, Value: "New Name"},
	}
	result := applyUpdate(existing, req)
	if result.Name != "New Name" {
		t.Errorf("expected Name 'New Name', got %q", result.Name)
	}
	if result.Rating != 3 {
		t.Errorf("rating should be unchanged, got %d", result.Rating)
	}
}

func TestApplyUpdate_Rating(t *testing.T) {
	existing := makeRecord(1, "Coffee", 3)
	req := &api.UpdateRecordRequest{
		Rating: api.OptInt{Set: true, Value: 5},
	}
	result := applyUpdate(existing, req)
	if result.Rating != 5 {
		t.Errorf("expected Rating 5, got %d", result.Rating)
	}
}

func TestApplyUpdate_OriginSet(t *testing.T) {
	existing := makeRecord(1, "Coffee", 3)
	req := &api.UpdateRecordRequest{
		Origin: api.OptNilString{Set: true, Value: "Ethiopia"},
	}
	result := applyUpdate(existing, req)
	if !result.Origin.Valid || result.Origin.String != "Ethiopia" {
		t.Errorf("expected Origin 'Ethiopia', got %v", result.Origin)
	}
}

func TestApplyUpdate_OriginCleared(t *testing.T) {
	existing := makeRecord(1, "Coffee", 3)
	existing.Origin = sql.NullString{String: "Ethiopia", Valid: true}
	req := &api.UpdateRecordRequest{
		Origin: api.OptNilString{Set: true, Null: true},
	}
	result := applyUpdate(existing, req)
	if result.Origin.Valid {
		t.Errorf("expected Origin to be cleared (null), got %v", result.Origin)
	}
}

func TestApplyUpdate_TastingNote(t *testing.T) {
	existing := makeRecord(1, "Coffee", 3)
	req := &api.UpdateRecordRequest{
		TastingNote: api.OptNilString{Set: true, Value: "fruity and bright"},
	}
	result := applyUpdate(existing, req)
	if !result.TastingNote.Valid || result.TastingNote.String != "fruity and bright" {
		t.Errorf("expected TastingNote set, got %v", result.TastingNote)
	}
}

// --- RecordUsecase.List ---

func TestRecordUsecaseList_NoFilter(t *testing.T) {
	records := []sqlcgen.Record{makeRecord(1, "A", 3), makeRecord(2, "B", 4)}
	repo := &stubRecordRepo{
		listAllFn: func(_ context.Context, userID int64) ([]sqlcgen.Record, error) {
			if userID != defaultUserID {
				t.Errorf("unexpected userID %d", userID)
			}
			return records, nil
		},
	}
	uc := NewRecordUsecase(repo)
	got, err := uc.List(context.Background(), api.ListRecordsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 records, got %d", len(got))
	}
}

func TestRecordUsecaseList_ByOrigin(t *testing.T) {
	r := makeRecord(1, "A", 3)
	r.Origin = sql.NullString{String: "Ethiopia", Valid: true}
	repo := &stubRecordRepo{
		listByOriginFn: func(_ context.Context, _ int64, origin string) ([]sqlcgen.Record, error) {
			if origin != "Ethiopia" {
				t.Errorf("unexpected origin %q", origin)
			}
			return []sqlcgen.Record{r}, nil
		},
	}
	uc := NewRecordUsecase(repo)
	params := api.ListRecordsParams{Origin: api.OptString{Set: true, Value: "Ethiopia"}}
	got, err := uc.List(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 record, got %d", len(got))
	}
}

func TestRecordUsecaseList_ByRatingMin(t *testing.T) {
	r := makeRecord(1, "A", 4)
	repo := &stubRecordRepo{
		listByRatingFn: func(_ context.Context, _ int64, ratingMin int8) ([]sqlcgen.Record, error) {
			if ratingMin != 4 {
				t.Errorf("unexpected ratingMin %d", ratingMin)
			}
			return []sqlcgen.Record{r}, nil
		},
	}
	uc := NewRecordUsecase(repo)
	params := api.ListRecordsParams{RatingMin: api.OptInt{Set: true, Value: 4}}
	got, err := uc.List(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 record, got %d", len(got))
	}
}

// --- RecordUsecase.Get ---

func TestRecordUsecaseGet_Found(t *testing.T) {
	record := makeRecord(42, "Ethiopia Yirgacheffe", 5)
	related := []sqlcgen.Record{makeRecord(10, "Ethiopia Yirgacheffe", 4)}
	repo := &stubRecordRepo{
		getFn: func(_ context.Context, id, _ int64) (sqlcgen.Record, error) {
			if id != 42 {
				t.Errorf("unexpected id %d", id)
			}
			return record, nil
		},
		getRelatedFn: func(_ context.Context, _, _ int64, _ string) ([]sqlcgen.Record, error) {
			return related, nil
		},
	}
	uc := NewRecordUsecase(repo)
	got, gotRelated, err := uc.Get(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 42 {
		t.Errorf("expected ID 42, got %d", got.ID)
	}
	if len(gotRelated) != 1 {
		t.Errorf("expected 1 related record, got %d", len(gotRelated))
	}
}

func TestRecordUsecaseGet_NotFound(t *testing.T) {
	repo := &stubRecordRepo{
		getFn: func(_ context.Context, _, _ int64) (sqlcgen.Record, error) {
			return sqlcgen.Record{}, sql.ErrNoRows
		},
	}
	uc := NewRecordUsecase(repo)
	_, _, err := uc.Get(context.Background(), 99)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// --- RecordUsecase.Create ---

func TestRecordUsecaseCreate(t *testing.T) {
	created := makeRecord(1, "Test Coffee", 4)
	created.TastingNote = sql.NullString{String: "bright", Valid: true}
	created.IsNoteFilled = true

	repo := &stubRecordRepo{
		createFn: func(_ context.Context, p sqlcgen.CreateRecordParams) (int64, error) {
			if p.Name != "Test Coffee" {
				t.Errorf("expected Name 'Test Coffee', got %q", p.Name)
			}
			if p.Rating != 4 {
				t.Errorf("expected Rating 4, got %d", p.Rating)
			}
			if !p.IsNoteFilled {
				t.Error("expected IsNoteFilled true")
			}
			return 1, nil
		},
		getFn: func(_ context.Context, id, _ int64) (sqlcgen.Record, error) {
			if id != 1 {
				t.Errorf("expected Get with id 1, got %d", id)
			}
			return created, nil
		},
	}
	uc := NewRecordUsecase(repo)
	req := &api.CreateRecordRequest{
		Name:        "Test Coffee",
		Rating:      api.OptInt{Set: true, Value: 4},
		TastingNote: api.OptNilString{Set: true, Value: "bright"},
	}
	got, err := uc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("expected ID 1, got %d", got.ID)
	}
	if !got.IsNoteFilled {
		t.Error("expected IsNoteFilled true")
	}
}

// --- RecordUsecase.Update ---

func TestRecordUsecaseUpdate_Found(t *testing.T) {
	existing := makeRecord(1, "Old Name", 3)
	updated := makeRecord(1, "New Name", 3)

	repo := &stubRecordRepo{
		getFn: func(_ context.Context, _, _ int64) (sqlcgen.Record, error) {
			return existing, nil
		},
		updateFn: func(_ context.Context, p sqlcgen.UpdateRecordParams) error {
			if p.Name != "New Name" {
				t.Errorf("expected Name 'New Name', got %q", p.Name)
			}
			return nil
		},
	}
	// Second Get call (after update) returns the updated record
	callCount := 0
	repo.getFn = func(_ context.Context, _, _ int64) (sqlcgen.Record, error) {
		callCount++
		if callCount == 1 {
			return existing, nil
		}
		return updated, nil
	}

	uc := NewRecordUsecase(repo)
	req := &api.UpdateRecordRequest{
		Name: api.OptString{Set: true, Value: "New Name"},
	}
	got, err := uc.Update(context.Background(), 1, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "New Name" {
		t.Errorf("expected Name 'New Name', got %q", got.Name)
	}
}

func TestRecordUsecaseUpdate_NotFound(t *testing.T) {
	repo := &stubRecordRepo{
		getFn: func(_ context.Context, _, _ int64) (sqlcgen.Record, error) {
			return sqlcgen.Record{}, sql.ErrNoRows
		},
	}
	uc := NewRecordUsecase(repo)
	_, err := uc.Update(context.Background(), 99, &api.UpdateRecordRequest{})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// --- RecordUsecase.Delete ---

func TestRecordUsecaseDelete_Found(t *testing.T) {
	deleted := false
	repo := &stubRecordRepo{
		getFn: func(_ context.Context, _, _ int64) (sqlcgen.Record, error) {
			return makeRecord(1, "Coffee", 3), nil
		},
		deleteFn: func(_ context.Context, id, userID int64) error {
			if id != 1 || userID != defaultUserID {
				t.Errorf("unexpected delete args id=%d userID=%d", id, userID)
			}
			deleted = true
			return nil
		},
	}
	uc := NewRecordUsecase(repo)
	if err := uc.Delete(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Error("expected Delete to be called on repo")
	}
}

func TestRecordUsecaseDelete_NotFound(t *testing.T) {
	repo := &stubRecordRepo{
		getFn: func(_ context.Context, _, _ int64) (sqlcgen.Record, error) {
			return sqlcgen.Record{}, sql.ErrNoRows
		},
	}
	uc := NewRecordUsecase(repo)
	err := uc.Delete(context.Background(), 99)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
