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

func RegisterAnalyticsEndpoints(humaAPI huma.API, prefix string, provider db.DBProvider) {
	analyticsGroup := huma.NewGroup(humaAPI, prefix)

	huma.Register(analyticsGroup, huma.Operation{
		OperationID: "get-daily-summary",
		Method:      "GET",
		Path:        "/daily/{user_id}",
		Summary:     "Get daily nutrition summary",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		UserID string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		Date   string `query:"date" example:"2025-01-07" doc:"Date (YYYY-MM-DD), defaults to today"`
	}) (*GetDailySummaryResponse, error) {
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		date, err := parseDate(input.Date)
		if err != nil && input.Date != "" {
			return nil, huma.Error400BadRequest("Invalid date format. Expected YYYY-MM-DD", err)
		}
		if input.Date == "" {
			date = time.Now()
		}

		summary, err := database.Nutrition().GetDaily(ctx, input.UserID, date)
		if err != nil {
			return nil, fmt.Errorf("failed to get daily summary: %w", err)
		}

		resp := &GetDailySummaryResponse{}
		resp.Body.Summary = summary
		return resp, nil
	})

	huma.Register(analyticsGroup, huma.Operation{
		OperationID: "list-daily-summaries",
		Method:      "GET",
		Path:        "/daily/{user_id}/range",
		Summary:     "List daily nutrition summaries in a range",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		UserID    string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		StartDate string `query:"start_date" example:"2025-01-01" doc:"Start date (YYYY-MM-DD)"`
		EndDate   string `query:"end_date" example:"2025-01-31" doc:"End date (YYYY-MM-DD)"`
		Limit     int    `query:"limit" default:"30" doc:"Maximum number of days to return"`
		Offset    int    `query:"offset" default:"0" doc:"Number of days to skip"`
	}) (*ListDailySummariesResponse, error) {
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		start, err := parseDate(input.StartDate)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid start_date format. Expected YYYY-MM-DD", err)
		}

		end, err := parseDate(input.EndDate)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid end_date format. Expected YYYY-MM-DD", err)
		}

		summaries, err := database.Nutrition().ListDailyByDateRange(ctx, input.UserID, start, end)
		if err != nil {
			return nil, fmt.Errorf("failed to list daily summaries: %w", err)
		}

		resp := &ListDailySummariesResponse{}
		resp.Body.Summaries = summaries
		return resp, nil
	})

	huma.Register(analyticsGroup, huma.Operation{
		OperationID: "get-weekly-summary",
		Method:      "GET",
		Path:        "/weekly/{user_id}",
		Summary:     "Get weekly nutrition summary",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		UserID        string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		WeekStartDate string `query:"week_start_date" example:"2025-01-06" doc:"Week start date (YYYY-MM-DD), defaults to current week"`
	}) (*GetWeeklySummaryResponse, error) {
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		weekStart, err := parseDate(input.WeekStartDate)
		if err != nil && input.WeekStartDate != "" {
			return nil, huma.Error400BadRequest("Invalid week_start_date format. Expected YYYY-MM-DD", err)
		}
		if input.WeekStartDate == "" {
			now := time.Now()
			weekStart = now.AddDate(0, 0, -int(now.Weekday()-1))
		}

		summary, err := database.Nutrition().GetWeekly(ctx, input.UserID, weekStart)
		if err != nil {
			return nil, fmt.Errorf("failed to get weekly summary: %w", err)
		}

		resp := &GetWeeklySummaryResponse{}
		resp.Body.Summary = summary
		return resp, nil
	})

	huma.Register(analyticsGroup, huma.Operation{
		OperationID: "list-weekly-summaries",
		Method:      "GET",
		Path:        "/weekly/{user_id}/range",
		Summary:     "List weekly nutrition summaries in a range",
		Tags:        []string{"Analytics"},
	}, func(ctx context.Context, input *struct {
		UserID    string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		StartDate string `query:"start_date" example:"2025-01-01" doc:"Start date (YYYY-MM-DD)"`
		EndDate   string `query:"end_date" example:"2025-01-31" doc:"End date (YYYY-MM-DD)"`
		Limit     int    `query:"limit" default:"12" doc:"Maximum number of weeks to return"`
		Offset    int    `query:"offset" default:"0" doc:"Number of weeks to skip"`
	}) (*ListWeeklySummariesResponse, error) {
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		start, err := parseDate(input.StartDate)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid start_date format. Expected YYYY-MM-DD", err)
		}

		end, err := parseDate(input.EndDate)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid end_date format. Expected YYYY-MM-DD", err)
		}

		summaries, err := database.Nutrition().ListWeeklyByDateRange(ctx, input.UserID, start, end)
		if err != nil {
			return nil, fmt.Errorf("failed to list weekly summaries: %w", err)
		}

		resp := &ListWeeklySummariesResponse{}
		resp.Body.Summaries = summaries
		return resp, nil
	})
}
