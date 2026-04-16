package handler

import (
	"context"
	"errors"
	"time"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/api"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/auth"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/database/sqlcgen"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/usecase"
)

// UserinfoProvider fetches user info from the OIDC userinfo endpoint (or returns
// dummy data when OIDC is disabled).
type UserinfoProvider interface {
	FetchUserinfo(ctx context.Context, accessToken string) (auth.UserinfoResult, error)
}

type Handler struct {
	recordUsecase    *usecase.RecordUsecase
	statsUsecase     *usecase.StatsUsecase
	userinfoProvider UserinfoProvider
}

func New(recordUsecase *usecase.RecordUsecase, statsUsecase *usecase.StatsUsecase, userinfoProvider UserinfoProvider) *Handler {
	return &Handler{
		recordUsecase:    recordUsecase,
		statsUsecase:     statsUsecase,
		userinfoProvider: userinfoProvider,
	}
}

func (h *Handler) GetUserinfo(ctx context.Context) (api.GetUserinfoRes, error) {
	_, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return &api.GetUserinfoUnauthorized{Message: "unauthorized"}, nil
	}
	token, _ := auth.TokenFromContext(ctx)
	result, err := h.userinfoProvider.FetchUserinfo(ctx, token)
	if err != nil {
		return &api.GetUserinfoInternalServerError{Message: err.Error()}, nil
	}
	resp := api.UserinfoResponse{Sub: result.Sub}
	if result.Name != "" {
		resp.Name = api.OptNilString{Set: true, Value: result.Name}
	} else {
		resp.Name = api.OptNilString{Set: true, Null: true}
	}
	if result.Email != "" {
		resp.Email = api.OptNilString{Set: true, Value: result.Email}
	} else {
		resp.Email = api.OptNilString{Set: true, Null: true}
	}
	if result.Picture != "" {
		resp.Picture = api.OptNilString{Set: true, Value: result.Picture}
	} else {
		resp.Picture = api.OptNilString{Set: true, Null: true}
	}
	return &resp, nil
}

func (h *Handler) ListRecords(ctx context.Context, params api.ListRecordsParams) (api.ListRecordsRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return &api.ListRecordsUnauthorized{Message: "unauthorized"}, nil
	}
	records, err := h.recordUsecase.List(ctx, userID, params)
	if err != nil {
		return &api.ListRecordsInternalServerError{Message: err.Error()}, nil
	}
	result := make(api.ListRecordsOKApplicationJSON, 0, len(records))
	for _, r := range records {
		result = append(result, toOgenRecord(r))
	}
	return &result, nil
}

func (h *Handler) CreateRecord(ctx context.Context, req *api.CreateRecordRequest) (api.CreateRecordRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return &api.CreateRecordUnauthorized{Message: "unauthorized"}, nil
	}
	record, err := h.recordUsecase.Create(ctx, userID, req)
	if err != nil {
		return &api.CreateRecordInternalServerError{Message: err.Error()}, nil
	}
	r := toOgenRecord(record)
	return &r, nil
}

func (h *Handler) GetRecord(ctx context.Context, params api.GetRecordParams) (api.GetRecordRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return &api.GetRecordUnauthorized{Message: "unauthorized"}, nil
	}
	record, related, err := h.recordUsecase.Get(ctx, userID, params.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return &api.GetRecordNotFound{Message: "record not found"}, nil
		}
		return &api.GetRecordInternalServerError{Message: err.Error()}, nil
	}
	detail := toOgenRecordDetail(record, related)
	return &detail, nil
}

func (h *Handler) UpdateRecord(ctx context.Context, req *api.UpdateRecordRequest, params api.UpdateRecordParams) (api.UpdateRecordRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return &api.UpdateRecordUnauthorized{Message: "unauthorized"}, nil
	}
	record, err := h.recordUsecase.Update(ctx, userID, params.ID, req)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return &api.UpdateRecordNotFound{Message: "record not found"}, nil
		}
		return &api.UpdateRecordInternalServerError{Message: err.Error()}, nil
	}
	r := toOgenRecord(record)
	return &r, nil
}

func (h *Handler) DeleteRecord(ctx context.Context, params api.DeleteRecordParams) (api.DeleteRecordRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return &api.DeleteRecordUnauthorized{Message: "unauthorized"}, nil
	}
	err := h.recordUsecase.Delete(ctx, userID, params.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return &api.DeleteRecordNotFound{Message: "record not found"}, nil
		}
		return &api.DeleteRecordInternalServerError{Message: err.Error()}, nil
	}
	return &api.DeleteRecordNoContent{}, nil
}

func (h *Handler) GetStatsSummary(ctx context.Context) (api.GetStatsSummaryRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return &api.GetStatsSummaryUnauthorized{Message: "unauthorized"}, nil
	}
	summary, err := h.statsUsecase.Summary(ctx, userID)
	if err != nil {
		return &api.GetStatsSummaryInternalServerError{Message: err.Error()}, nil
	}
	return summary, nil
}

func (h *Handler) GetStatsFlavorWords(ctx context.Context) (api.GetStatsFlavorWordsRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return &api.GetStatsFlavorWordsUnauthorized{Message: "unauthorized"}, nil
	}
	words, err := h.statsUsecase.FlavorWords(ctx, userID)
	if err != nil {
		return &api.GetStatsFlavorWordsInternalServerError{Message: err.Error()}, nil
	}
	result := api.GetStatsFlavorWordsOKApplicationJSON(words)
	return &result, nil
}

func (h *Handler) GetRecommend(ctx context.Context, params api.GetRecommendParams) (api.GetRecommendRes, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return &api.GetRecommendUnauthorized{Message: "unauthorized"}, nil
	}
	result, err := h.statsUsecase.Recommend(ctx, userID, params)
	if err != nil {
		return &api.GetRecommendInternalServerError{Message: err.Error()}, nil
	}
	return result, nil
}

// --- conversion helpers ---

func toOgenRecord(r sqlcgen.Record) api.Record {
	rec := api.Record{
		ID:           r.ID,
		Name:         r.Name,
		Rating:       int(r.Rating),
		IsNoteFilled: r.IsNoteFilled,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
	if r.Origin.Valid {
		rec.Origin = api.OptNilString{Set: true, Value: r.Origin.String}
	} else {
		rec.Origin = api.OptNilString{Set: true, Null: true}
	}
	if r.RoastLevel.Valid {
		rec.RoastLevel = api.OptNilRecordRoastLevel{
			Set:   true,
			Value: api.RecordRoastLevel(r.RoastLevel.String),
		}
	} else {
		rec.RoastLevel = api.OptNilRecordRoastLevel{Set: true, Null: true}
	}
	if r.Shop.Valid {
		rec.Shop = api.OptNilString{Set: true, Value: r.Shop.String}
	} else {
		rec.Shop = api.OptNilString{Set: true, Null: true}
	}
	if r.Price.Valid {
		rec.Price = api.OptNilInt{Set: true, Value: int(r.Price.Int32)}
	} else {
		rec.Price = api.OptNilInt{Set: true, Null: true}
	}
	if r.PurchasedAt.Valid {
		rec.PurchasedAt = api.OptNilDate{Set: true, Value: r.PurchasedAt.Time.UTC().Truncate(24 * time.Hour)}
	} else {
		rec.PurchasedAt = api.OptNilDate{Set: true, Null: true}
	}
	if r.TastingNote.Valid {
		rec.TastingNote = api.OptNilString{Set: true, Value: r.TastingNote.String}
	} else {
		rec.TastingNote = api.OptNilString{Set: true, Null: true}
	}
	if r.BrewMethod.Valid {
		rec.BrewMethod = api.OptNilRecordBrewMethod{
			Set:   true,
			Value: api.RecordBrewMethod(r.BrewMethod.String),
		}
	} else {
		rec.BrewMethod = api.OptNilRecordBrewMethod{Set: true, Null: true}
	}
	if r.Recipe.Valid {
		rec.Recipe = api.OptNilString{Set: true, Value: r.Recipe.String}
	} else {
		rec.Recipe = api.OptNilString{Set: true, Null: true}
	}
	return rec
}

func toOgenRecordDetail(r sqlcgen.Record, related []sqlcgen.Record) api.RecordDetail {
	detail := api.RecordDetail{
		ID:           r.ID,
		Name:         r.Name,
		Rating:       int(r.Rating),
		IsNoteFilled: r.IsNoteFilled,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
	if r.Origin.Valid {
		detail.Origin = api.OptNilString{Set: true, Value: r.Origin.String}
	} else {
		detail.Origin = api.OptNilString{Set: true, Null: true}
	}
	if r.RoastLevel.Valid {
		detail.RoastLevel = api.OptNilRecordDetailRoastLevel{
			Set:   true,
			Value: api.RecordDetailRoastLevel(r.RoastLevel.String),
		}
	} else {
		detail.RoastLevel = api.OptNilRecordDetailRoastLevel{Set: true, Null: true}
	}
	if r.Shop.Valid {
		detail.Shop = api.OptNilString{Set: true, Value: r.Shop.String}
	} else {
		detail.Shop = api.OptNilString{Set: true, Null: true}
	}
	if r.Price.Valid {
		detail.Price = api.OptNilInt{Set: true, Value: int(r.Price.Int32)}
	} else {
		detail.Price = api.OptNilInt{Set: true, Null: true}
	}
	if r.PurchasedAt.Valid {
		detail.PurchasedAt = api.OptNilDate{Set: true, Value: r.PurchasedAt.Time.UTC().Truncate(24 * time.Hour)}
	} else {
		detail.PurchasedAt = api.OptNilDate{Set: true, Null: true}
	}
	if r.TastingNote.Valid {
		detail.TastingNote = api.OptNilString{Set: true, Value: r.TastingNote.String}
	} else {
		detail.TastingNote = api.OptNilString{Set: true, Null: true}
	}
	if r.BrewMethod.Valid {
		detail.BrewMethod = api.OptNilRecordDetailBrewMethod{
			Set:   true,
			Value: api.RecordDetailBrewMethod(r.BrewMethod.String),
		}
	} else {
		detail.BrewMethod = api.OptNilRecordDetailBrewMethod{Set: true, Null: true}
	}
	if r.Recipe.Valid {
		detail.Recipe = api.OptNilString{Set: true, Value: r.Recipe.String}
	} else {
		detail.Recipe = api.OptNilString{Set: true, Null: true}
	}

	detail.RelatedRecords = make([]api.Record, 0, len(related))
	for _, rel := range related {
		detail.RelatedRecords = append(detail.RelatedRecords, toOgenRecord(rel))
	}
	return detail
}
