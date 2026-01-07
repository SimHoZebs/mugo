package routes

import (
	"context"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/models"
)

type GetDailySummaryResponse struct {
	Body struct {
		Summary *models.DailyNutritionSummary `json:"summary"`
	}
}

type ListDailySummariesResponse struct {
	Body struct {
		Summaries []*models.DailyNutritionSummary `json:"summaries"`
	}
}

type GetWeeklySummaryResponse struct {
	Body struct {
		Summary *models.WeeklyNutritionSummary `json:"summary"`
	}
}

type ListWeeklySummariesResponse struct {
	Body struct {
		Summaries []*models.WeeklyNutritionSummary `json:"summaries"`
	}
}

// RegisterAnalyticsEndpoints registers nutrition analytics endpoints.
func RegisterAnalyticsEndpoints(humaAPI huma.API, prefix string, database *db.Database) {
	analyticsGroup := huma.NewGroup(humaAPI, prefix)

	huma.Get(analyticsGroup, "/daily/{user_id}", func(ctx context.Context, input *struct {
		UserID string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		Date   string `query:"date" example:"2025-01-07" doc:"Date (YYYY-MM-DD), defaults to today"`
	}) (*GetDailySummaryResponse, error) {
		date := parseDate(input.Date)
		if input.Date == "" {
			date = time.Now()
		}

		summary, err := database.NutritionRepository.GetDaily(ctx, input.UserID, date)
		if err != nil {
			return nil, fmt.Errorf("failed to get daily summary: %w", err)
		}

		resp := &GetDailySummaryResponse{}
		resp.Body.Summary = summary
		return resp, nil
	})

	huma.Get(analyticsGroup, "/daily/{user_id}/range", func(ctx context.Context, input *struct {
		UserID    string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		StartDate string `query:"start_date" example:"2025-01-01" doc:"Start date (YYYY-MM-DD)"`
		EndDate   string `query:"end_date" example:"2025-01-31" doc:"End date (YYYY-MM-DD)"`
		Limit     int    `query:"limit" default:"30" doc:"Maximum number of days to return"`
		Offset    int    `query:"offset" default:"0" doc:"Number of days to skip"`
	}) (*ListDailySummariesResponse, error) {
		summaries, err := database.NutritionRepository.ListDailyByDateRange(ctx, input.UserID, parseDate(input.StartDate), parseDate(input.EndDate))
		if err != nil {
			return nil, fmt.Errorf("failed to list daily summaries: %w", err)
		}

		resp := &ListDailySummariesResponse{}
		resp.Body.Summaries = summaries
		return resp, nil
	})

	huma.Get(analyticsGroup, "/weekly/{user_id}", func(ctx context.Context, input *struct {
		UserID        string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		WeekStartDate string `query:"week_start_date" example:"2025-01-06" doc:"Week start date (YYYY-MM-DD), defaults to current week"`
	}) (*GetWeeklySummaryResponse, error) {
		weekStart := parseDate(input.WeekStartDate)
		if input.WeekStartDate == "" {
			now := time.Now()
			weekStart = now.AddDate(0, 0, -int(now.Weekday()-1))
		}

		summary, err := database.NutritionRepository.GetWeekly(ctx, input.UserID, weekStart)
		if err != nil {
			return nil, fmt.Errorf("failed to get weekly summary: %w", err)
		}

		resp := &GetWeeklySummaryResponse{}
		resp.Body.Summary = summary
		return resp, nil
	})

	huma.Get(analyticsGroup, "/weekly/{user_id}/range", func(ctx context.Context, input *struct {
		UserID    string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		StartDate string `query:"start_date" example:"2025-01-01" doc:"Start date (YYYY-MM-DD)"`
		EndDate   string `query:"end_date" example:"2025-01-31" doc:"End date (YYYY-MM-DD)"`
		Limit     int    `query:"limit" default:"12" doc:"Maximum number of weeks to return"`
		Offset    int    `query:"offset" default:"0" doc:"Number of weeks to skip"`
	}) (*ListWeeklySummariesResponse, error) {
		summaries, err := database.NutritionRepository.ListWeeklyByDateRange(ctx, input.UserID, parseDate(input.StartDate), parseDate(input.EndDate))
		if err != nil {
			return nil, fmt.Errorf("failed to list weekly summaries: %w", err)
		}

		resp := &ListWeeklySummariesResponse{}
		resp.Body.Summaries = summaries
		return resp, nil
	})
}
