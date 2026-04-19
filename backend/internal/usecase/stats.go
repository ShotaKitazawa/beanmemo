package usecase

import (
	"context"
	"database/sql"
	"sort"
	"strings"
	"unicode"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/api"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/database/sqlcgen"
)

const recommendMinRecords = 5

type statsRepository interface {
	StatsByOrigin(ctx context.Context, userID int64) ([]sqlcgen.StatsByOriginRow, error)
	StatsByRoastLevel(ctx context.Context, userID int64) ([]sqlcgen.StatsByRoastLevelRow, error)
	StatsByBrewMethod(ctx context.Context, userID int64) ([]sqlcgen.StatsByBrewMethodRow, error)
	AvgRatingByOrigin(ctx context.Context, userID int64, origin string) (sql.NullFloat64, error)
	AvgRatingByName(ctx context.Context, userID int64, name string) (sql.NullFloat64, error)
}

type statsRecordRepository interface {
	Count(ctx context.Context, userID int64) (int64, error)
	ListAllTastingNotes(ctx context.Context, userID int64) ([]sql.NullString, error)
}

type StatsUsecase struct {
	statsRepo  statsRepository
	recordRepo statsRecordRepository
}

func NewStatsUsecase(statsRepo statsRepository, recordRepo statsRecordRepository) *StatsUsecase {
	return &StatsUsecase{statsRepo: statsRepo, recordRepo: recordRepo}
}

func (u *StatsUsecase) Summary(ctx context.Context, userID int64) (*api.StatsSummary, error) {
	total, err := u.recordRepo.Count(ctx, userID)
	if err != nil {
		return nil, err
	}

	origins, err := u.statsRepo.StatsByOrigin(ctx, userID)
	if err != nil {
		return nil, err
	}
	roasts, err := u.statsRepo.StatsByRoastLevel(ctx, userID)
	if err != nil {
		return nil, err
	}
	brews, err := u.statsRepo.StatsByBrewMethod(ctx, userID)
	if err != nil {
		return nil, err
	}

	summary := &api.StatsSummary{
		TotalRecords: int(total),
	}

	for _, row := range origins {
		if !row.Label.Valid {
			continue
		}
		summary.ByOrigin = append(summary.ByOrigin, api.GroupStat{
			Label:     row.Label.String,
			Count:     int(row.Count),
			AvgRating: toFloat32(row.AvgRating),
		})
	}
	for _, row := range roasts {
		if !row.Label.Valid {
			continue
		}
		summary.ByRoastLevel = append(summary.ByRoastLevel, api.GroupStat{
			Label:     row.Label.String,
			Count:     int(row.Count),
			AvgRating: toFloat32(row.AvgRating),
		})
	}
	for _, row := range brews {
		if !row.Label.Valid {
			continue
		}
		summary.ByBrewMethod = append(summary.ByBrewMethod, api.GroupStat{
			Label:     row.Label.String,
			Count:     int(row.Count),
			AvgRating: toFloat32(row.AvgRating),
		})
	}

	if summary.ByOrigin == nil {
		summary.ByOrigin = []api.GroupStat{}
	}
	if summary.ByRoastLevel == nil {
		summary.ByRoastLevel = []api.GroupStat{}
	}
	if summary.ByBrewMethod == nil {
		summary.ByBrewMethod = []api.GroupStat{}
	}

	return summary, nil
}

func (u *StatsUsecase) FlavorWords(ctx context.Context, userID int64) ([]api.FlavorWord, error) {
	notes, err := u.recordRepo.ListAllTastingNotes(ctx, userID)
	if err != nil {
		return nil, err
	}

	freq := make(map[string]int)
	for _, n := range notes {
		if !n.Valid || n.String == "" {
			continue
		}
		words := tokenize(n.String)
		for _, w := range words {
			if len([]rune(w)) >= 2 {
				freq[w]++
			}
		}
	}

	type kv struct {
		word  string
		count int
	}
	var pairs []kv
	for w, c := range freq {
		pairs = append(pairs, kv{w, c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].word < pairs[j].word
	})

	const maxWords = 20
	result := make([]api.FlavorWord, 0, maxWords)
	for i, p := range pairs {
		if i >= maxWords {
			break
		}
		result = append(result, api.FlavorWord{Word: p.word, Count: p.count})
	}
	return result, nil
}

func (u *StatsUsecase) Recommend(ctx context.Context, userID int64, params api.GetRecommendParams) (*api.RecommendResult, error) {
	total, err := u.recordRepo.Count(ctx, userID)
	if err != nil {
		return nil, err
	}

	if total < recommendMinRecords {
		needed := recommendMinRecords - int(total)
		return &api.RecommendResult{
			Locked:       true,
			TotalRecords: int(total),
			RecordsNeeded: api.OptNilInt{
				Set:   true,
				Value: needed,
			},
		}, nil
	}

	result := &api.RecommendResult{
		Locked:       false,
		TotalRecords: int(total),
	}

	var scores []float32
	if params.Origin.Set && params.Origin.Value != "" {
		avg, err := u.statsRepo.AvgRatingByOrigin(ctx, userID, params.Origin.Value)
		if err == nil && avg.Valid {
			v := toFloat32(avg)
			result.OriginAvg = api.OptNilFloat32{Set: true, Value: v}
			scores = append(scores, v)
		}
	}
	if params.Name.Set && params.Name.Value != "" {
		avg, err := u.statsRepo.AvgRatingByName(ctx, userID, params.Name.Value)
		if err == nil && avg.Valid {
			v := toFloat32(avg)
			result.NameMatchAvg = api.OptNilFloat32{Set: true, Value: v}
			scores = append(scores, v)
		}
	}

	if len(scores) > 0 {
		var sum float32
		for _, s := range scores {
			sum += s
		}
		result.Score = api.OptNilFloat32{Set: true, Value: sum / float32(len(scores))}
	} else {
		result.Score = api.OptNilFloat32{Set: true, Null: true}
	}

	return result, nil
}

// tokenize splits text by whitespace and punctuation
func tokenize(text string) []string {
	f := func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r) || r == ',' || r == '。' || r == '、' || r == '・'
	}
	raw := strings.FieldsFunc(text, f)
	var result []string
	for _, w := range raw {
		w = strings.ToLower(strings.TrimSpace(w))
		if w != "" {
			result = append(result, w)
		}
	}
	return result
}

func toFloat32(v sql.NullFloat64) float32 {
	if !v.Valid {
		return 0
	}
	return float32(v.Float64)
}
