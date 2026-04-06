package handler

import (
	"context"
	"errors"
	"time"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/api"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/database/sqlcgen"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/usecase"
)

type Handler struct {
	recordUsecase *usecase.RecordUsecase
	statsUsecase  *usecase.StatsUsecase
}

func New(recordUsecase *usecase.RecordUsecase, statsUsecase *usecase.StatsUsecase) *Handler {
	return &Handler{recordUsecase: recordUsecase, statsUsecase: statsUsecase}
}

func (h *Handler) ListRecords(ctx context.Context, params api.ListRecordsParams) (api.ListRecordsRes, error) {
	records, err := h.recordUsecase.List(ctx, params)
	if err != nil {
		msg := err.Error()
		return &api.Error{Message: msg}, nil
	}
	result := make(api.ListRecordsOKApplicationJSON, 0, len(records))
	for _, r := range records {
		result = append(result, toOgenRecord(r))
	}
	return &result, nil
}

func (h *Handler) CreateRecord(ctx context.Context, req *api.CreateRecordRequest) (api.CreateRecordRes, error) {
	record, err := h.recordUsecase.Create(ctx, req)
	if err != nil {
		msg := err.Error()
		return &api.CreateRecordInternalServerError{Message: msg}, nil
	}
	r := toOgenRecord(record)
	return &r, nil
}

func (h *Handler) GetRecord(ctx context.Context, params api.GetRecordParams) (api.GetRecordRes, error) {
	record, related, err := h.recordUsecase.Get(ctx, params.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return &api.GetRecordNotFound{Message: "record not found"}, nil
		}
		msg := err.Error()
		return &api.GetRecordInternalServerError{Message: msg}, nil
	}

	detail := toOgenRecordDetail(record, related)
	return &detail, nil
}

func (h *Handler) UpdateRecord(ctx context.Context, req *api.UpdateRecordRequest, params api.UpdateRecordParams) (api.UpdateRecordRes, error) {
	record, err := h.recordUsecase.Update(ctx, params.ID, req)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return &api.UpdateRecordNotFound{Message: "record not found"}, nil
		}
		msg := err.Error()
		return &api.UpdateRecordInternalServerError{Message: msg}, nil
	}
	r := toOgenRecord(record)
	return &r, nil
}

func (h *Handler) DeleteRecord(ctx context.Context, params api.DeleteRecordParams) (api.DeleteRecordRes, error) {
	err := h.recordUsecase.Delete(ctx, params.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return &api.DeleteRecordNotFound{Message: "record not found"}, nil
		}
		msg := err.Error()
		return &api.DeleteRecordInternalServerError{Message: msg}, nil
	}
	return &api.DeleteRecordNoContent{}, nil
}

func (h *Handler) GetStatsSummary(ctx context.Context) (api.GetStatsSummaryRes, error) {
	summary, err := h.statsUsecase.Summary(ctx)
	if err != nil {
		msg := err.Error()
		return &api.Error{Message: msg}, nil
	}
	return summary, nil
}

func (h *Handler) GetStatsFlavorWords(ctx context.Context) (api.GetStatsFlavorWordsRes, error) {
	words, err := h.statsUsecase.FlavorWords(ctx)
	if err != nil {
		msg := err.Error()
		return &api.Error{Message: msg}, nil
	}
	result := api.GetStatsFlavorWordsOKApplicationJSON(words)
	return &result, nil
}

func (h *Handler) GetRecommend(ctx context.Context, params api.GetRecommendParams) (api.GetRecommendRes, error) {
	result, err := h.statsUsecase.Recommend(ctx, params)
	if err != nil {
		msg := err.Error()
		return &api.Error{Message: msg}, nil
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
