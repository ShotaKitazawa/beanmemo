package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/api"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/database/sqlcgen"
)

// --- mocks ---

type stubStatsRepo struct {
	statsByOriginFn func(context.Context, int64) ([]sqlcgen.StatsByOriginRow, error)
	statsByRoastFn  func(context.Context, int64) ([]sqlcgen.StatsByRoastLevelRow, error)
	statsByBrewFn   func(context.Context, int64) ([]sqlcgen.StatsByBrewMethodRow, error)
	avgByOriginFn   func(context.Context, int64, string) (sql.NullFloat64, error)
	avgByNameFn     func(context.Context, int64, string) (sql.NullFloat64, error)
}

func (s *stubStatsRepo) StatsByOrigin(ctx context.Context, userID int64) ([]sqlcgen.StatsByOriginRow, error) {
	if s.statsByOriginFn != nil {
		return s.statsByOriginFn(ctx, userID)
	}
	return nil, nil
}

func (s *stubStatsRepo) StatsByRoastLevel(ctx context.Context, userID int64) ([]sqlcgen.StatsByRoastLevelRow, error) {
	if s.statsByRoastFn != nil {
		return s.statsByRoastFn(ctx, userID)
	}
	return nil, nil
}

func (s *stubStatsRepo) StatsByBrewMethod(ctx context.Context, userID int64) ([]sqlcgen.StatsByBrewMethodRow, error) {
	if s.statsByBrewFn != nil {
		return s.statsByBrewFn(ctx, userID)
	}
	return nil, nil
}

func (s *stubStatsRepo) AvgRatingByOrigin(ctx context.Context, userID int64, origin string) (sql.NullFloat64, error) {
	if s.avgByOriginFn != nil {
		return s.avgByOriginFn(ctx, userID, origin)
	}
	return sql.NullFloat64{}, nil
}

func (s *stubStatsRepo) AvgRatingByName(ctx context.Context, userID int64, name string) (sql.NullFloat64, error) {
	if s.avgByNameFn != nil {
		return s.avgByNameFn(ctx, userID, name)
	}
	return sql.NullFloat64{}, nil
}

type stubStatsRecordRepo struct {
	countFn               func(context.Context, int64) (int64, error)
	listAllTastingNotesFn func(context.Context, int64) ([]sql.NullString, error)
}

func (s *stubStatsRecordRepo) Count(ctx context.Context, userID int64) (int64, error) {
	if s.countFn != nil {
		return s.countFn(ctx, userID)
	}
	return 0, nil
}

func (s *stubStatsRecordRepo) ListAllTastingNotes(ctx context.Context, userID int64) ([]sql.NullString, error) {
	if s.listAllTastingNotesFn != nil {
		return s.listAllTastingNotesFn(ctx, userID)
	}
	return nil, nil
}

// --- tokenize ---

func TestTokenize_Basic(t *testing.T) {
	words := tokenize("fruity and bright")
	if len(words) != 3 {
		t.Fatalf("expected 3 words, got %d: %v", len(words), words)
	}
	if words[0] != "fruity" || words[1] != "and" || words[2] != "bright" {
		t.Errorf("unexpected words: %v", words)
	}
}

func TestTokenize_Punctuation(t *testing.T) {
	words := tokenize("citrus, floral. nutty")
	if len(words) != 3 {
		t.Fatalf("expected 3 words, got %d: %v", len(words), words)
	}
}

func TestTokenize_Japanese(t *testing.T) {
	words := tokenize("フルーティー・爽やか")
	if len(words) != 2 {
		t.Fatalf("expected 2 words, got %d: %v", len(words), words)
	}
}

func TestTokenize_Empty(t *testing.T) {
	words := tokenize("")
	if len(words) != 0 {
		t.Errorf("expected 0 words, got %d", len(words))
	}
}

func TestTokenize_LowercasesWords(t *testing.T) {
	words := tokenize("Fruity BRIGHT")
	for _, w := range words {
		for _, r := range w {
			if r >= 'A' && r <= 'Z' {
				t.Errorf("word %q should be lowercase", w)
			}
		}
	}
}

// --- toFloat32 ---

func TestToFloat32_ValidValue(t *testing.T) {
	got := toFloat32(sql.NullFloat64{Float64: 3.5, Valid: true})
	if got != 3.5 {
		t.Errorf("expected 3.5, got %v", got)
	}
}

func TestToFloat32_Invalid(t *testing.T) {
	got := toFloat32(sql.NullFloat64{Valid: false})
	if got != 0 {
		t.Errorf("expected 0 for invalid, got %v", got)
	}
}

// --- StatsUsecase.Summary ---

func TestStatsUsecaseSummary_Basic(t *testing.T) {
	statsRepo := &stubStatsRepo{
		statsByOriginFn: func(_ context.Context, _ int64) ([]sqlcgen.StatsByOriginRow, error) {
			return []sqlcgen.StatsByOriginRow{
				{Label: sql.NullString{String: "Ethiopia", Valid: true}, Count: 3, AvgRating: sql.NullFloat64{Float64: 4.0, Valid: true}},
			}, nil
		},
		statsByRoastFn: func(_ context.Context, _ int64) ([]sqlcgen.StatsByRoastLevelRow, error) {
			return []sqlcgen.StatsByRoastLevelRow{
				{Label: sql.NullString{String: "light", Valid: true}, Count: 2, AvgRating: sql.NullFloat64{Float64: 4.5, Valid: true}},
			}, nil
		},
		statsByBrewFn: func(_ context.Context, _ int64) ([]sqlcgen.StatsByBrewMethodRow, error) {
			return nil, nil
		},
	}
	recordRepo := &stubStatsRecordRepo{
		countFn: func(_ context.Context, _ int64) (int64, error) {
			return 5, nil
		},
	}
	uc := NewStatsUsecase(statsRepo, recordRepo)
	summary, err := uc.Summary(context.Background(), testUserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.TotalRecords != 5 {
		t.Errorf("expected TotalRecords 5, got %d", summary.TotalRecords)
	}
	if len(summary.ByOrigin) != 1 || summary.ByOrigin[0].Label != "Ethiopia" {
		t.Errorf("unexpected ByOrigin: %v", summary.ByOrigin)
	}
	if len(summary.ByRoastLevel) != 1 || summary.ByRoastLevel[0].Label != "light" {
		t.Errorf("unexpected ByRoastLevel: %v", summary.ByRoastLevel)
	}
	if len(summary.ByBrewMethod) != 0 {
		t.Errorf("expected empty ByBrewMethod, got %v", summary.ByBrewMethod)
	}
}

func TestStatsUsecaseSummary_NullLabelsSkipped(t *testing.T) {
	statsRepo := &stubStatsRepo{
		statsByOriginFn: func(_ context.Context, _ int64) ([]sqlcgen.StatsByOriginRow, error) {
			return []sqlcgen.StatsByOriginRow{
				{Label: sql.NullString{Valid: false}, Count: 2, AvgRating: sql.NullFloat64{Float64: 3.0, Valid: true}},
				{Label: sql.NullString{String: "Brazil", Valid: true}, Count: 1, AvgRating: sql.NullFloat64{Float64: 4.0, Valid: true}},
			}, nil
		},
		statsByRoastFn: func(_ context.Context, _ int64) ([]sqlcgen.StatsByRoastLevelRow, error) {
			return nil, nil
		},
		statsByBrewFn: func(_ context.Context, _ int64) ([]sqlcgen.StatsByBrewMethodRow, error) {
			return nil, nil
		},
	}
	recordRepo := &stubStatsRecordRepo{
		countFn: func(_ context.Context, _ int64) (int64, error) { return 3, nil },
	}
	uc := NewStatsUsecase(statsRepo, recordRepo)
	summary, err := uc.Summary(context.Background(), testUserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(summary.ByOrigin) != 1 || summary.ByOrigin[0].Label != "Brazil" {
		t.Errorf("null label should be skipped, got %v", summary.ByOrigin)
	}
}

// --- StatsUsecase.FlavorWords ---

func TestStatsUsecaseFlavorWords_Basic(t *testing.T) {
	recordRepo := &stubStatsRecordRepo{
		listAllTastingNotesFn: func(_ context.Context, _ int64) ([]sql.NullString, error) {
			return []sql.NullString{
				{String: "fruity bright fruity", Valid: true},
				{String: "bright citrus", Valid: true},
			}, nil
		},
	}
	uc := NewStatsUsecase(&stubStatsRepo{}, recordRepo)
	words, err := uc.FlavorWords(context.Background(), testUserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(words) == 0 {
		t.Fatal("expected flavor words, got none")
	}
	found := map[string]int{}
	for _, w := range words {
		found[w.Word] = w.Count
	}
	if found["fruity"] != 2 {
		t.Errorf("expected fruity count 2, got %d", found["fruity"])
	}
	if found["bright"] != 2 {
		t.Errorf("expected bright count 2, got %d", found["bright"])
	}
	if found["citrus"] != 1 {
		t.Errorf("expected citrus count 1, got %d", found["citrus"])
	}
}

func TestStatsUsecaseFlavorWords_ShortWordsExcluded(t *testing.T) {
	recordRepo := &stubStatsRecordRepo{
		listAllTastingNotesFn: func(_ context.Context, _ int64) ([]sql.NullString, error) {
			return []sql.NullString{
				{String: "a bb ccc", Valid: true},
			}, nil
		},
	}
	uc := NewStatsUsecase(&stubStatsRepo{}, recordRepo)
	words, err := uc.FlavorWords(context.Background(), testUserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, w := range words {
		if len([]rune(w.Word)) < 2 {
			t.Errorf("word %q should have been excluded (< 2 runes)", w.Word)
		}
	}
}

func TestStatsUsecaseFlavorWords_MaxTwenty(t *testing.T) {
	notes := make([]sql.NullString, 0)
	for i := 0; i < 30; i++ {
		word := string(rune('a'+i/26)) + string(rune('a'+i%26))
		notes = append(notes, sql.NullString{String: word, Valid: true})
	}
	recordRepo := &stubStatsRecordRepo{
		listAllTastingNotesFn: func(_ context.Context, _ int64) ([]sql.NullString, error) {
			return notes, nil
		},
	}
	uc := NewStatsUsecase(&stubStatsRepo{}, recordRepo)
	words, err := uc.FlavorWords(context.Background(), testUserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(words) > 20 {
		t.Errorf("expected at most 20 flavor words, got %d", len(words))
	}
}

// --- StatsUsecase.Recommend ---

func TestStatsUsecaseRecommend_Locked(t *testing.T) {
	recordRepo := &stubStatsRecordRepo{
		countFn: func(_ context.Context, _ int64) (int64, error) { return 3, nil },
	}
	uc := NewStatsUsecase(&stubStatsRepo{}, recordRepo)
	result, err := uc.Recommend(context.Background(), testUserID, api.GetRecommendParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Locked {
		t.Error("expected Locked=true when total < 5")
	}
	if result.TotalRecords != 3 {
		t.Errorf("expected TotalRecords 3, got %d", result.TotalRecords)
	}
	if !result.RecordsNeeded.Set || result.RecordsNeeded.Value != 2 {
		t.Errorf("expected RecordsNeeded 2, got %v", result.RecordsNeeded)
	}
}

func TestStatsUsecaseRecommend_UnlockedWithOrigin(t *testing.T) {
	statsRepo := &stubStatsRepo{
		avgByOriginFn: func(_ context.Context, _ int64, origin string) (sql.NullFloat64, error) {
			if origin != "Ethiopia" {
				return sql.NullFloat64{}, nil
			}
			return sql.NullFloat64{Float64: 4.2, Valid: true}, nil
		},
	}
	recordRepo := &stubStatsRecordRepo{
		countFn: func(_ context.Context, _ int64) (int64, error) { return 10, nil },
	}
	uc := NewStatsUsecase(statsRepo, recordRepo)
	params := api.GetRecommendParams{
		Origin: api.OptString{Set: true, Value: "Ethiopia"},
	}
	result, err := uc.Recommend(context.Background(), testUserID, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Locked {
		t.Error("expected Locked=false when total >= 5")
	}
	if !result.OriginAvg.Set {
		t.Error("expected OriginAvg to be set")
	}
	if result.OriginAvg.Value != float32(4.2) {
		t.Errorf("expected OriginAvg ~4.2, got %v", result.OriginAvg.Value)
	}
	if !result.Score.Set {
		t.Error("expected Score to be set")
	}
}

func TestStatsUsecaseRecommend_UnlockedNoParams(t *testing.T) {
	recordRepo := &stubStatsRecordRepo{
		countFn: func(_ context.Context, _ int64) (int64, error) { return 10, nil },
	}
	uc := NewStatsUsecase(&stubStatsRepo{}, recordRepo)
	result, err := uc.Recommend(context.Background(), testUserID, api.GetRecommendParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Locked {
		t.Error("expected Locked=false")
	}
	if !result.Score.Set || !result.Score.Null {
		t.Errorf("expected Score to be null when no params, got %v", result.Score)
	}
}
