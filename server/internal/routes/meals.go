package routes

import (
	"context"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/models"
)

type ListMealsResponse struct {
	Body struct {
		Meals []*models.MealLog `json:"meals"`
	}
}

type GetMealResponse struct {
	Body struct {
		Meal *models.MealLog `json:"meal"`
	}
}

type ListMealsByDateRangeRequest struct {
	UserID    string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
	StartDate string `query:"start_date" example:"2025-01-01" doc:"Start date (YYYY-MM-DD)"`
	EndDate   string `query:"end_date" example:"2025-01-31" doc:"End date (YYYY-MM-DD)"`
}

// RegisterMealEndpoints registers meal log endpoints.
func RegisterMealEndpoints(humaAPI huma.API, prefix string, database *db.Database) {
	mealsGroup := huma.NewGroup(humaAPI, prefix)

	huma.Get(mealsGroup, "/{user_id}", func(ctx context.Context, input *struct {
		UserID string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		Limit  int    `query:"limit" default:"50" doc:"Maximum number of meals to return"`
		Offset int    `query:"offset" default:"0" doc:"Number of meals to skip"`
	}) (*ListMealsResponse, error) {
		meals, err := database.MealLogRepository.ListByUser(ctx, input.UserID, input.Limit, input.Offset)
		if err != nil {
			return nil, fmt.Errorf("failed to list meals: %w", err)
		}

		resp := &ListMealsResponse{}
		resp.Body.Meals = meals
		return resp, nil
	})

	huma.Get(mealsGroup, "/{user_id}/date/{date}", func(ctx context.Context, input *struct {
		UserID string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		Date   string `path:"date" example:"2025-01-07" doc:"Date (YYYY-MM-DD)"`
	}) (*ListMealsResponse, error) {
		meals, err := database.MealLogRepository.ListByUserAndDate(ctx, input.UserID, parseDate(input.Date))
		if err != nil {
			return nil, fmt.Errorf("failed to list meals by date: %w", err)
		}

		resp := &ListMealsResponse{}
		resp.Body.Meals = meals
		return resp, nil
	})

	huma.Get(mealsGroup, "/{user_id}/range", func(ctx context.Context, input *struct {
		UserID    string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		StartDate string `query:"start_date" example:"2025-01-01" doc:"Start date (YYYY-MM-DD)"`
		EndDate   string `query:"end_date" example:"2025-01-31" doc:"End date (YYYY-MM-DD)"`
	}) (*ListMealsResponse, error) {
		meals, err := database.MealLogRepository.ListByUserAndDateRange(ctx, input.UserID, parseDate(input.StartDate), parseDate(input.EndDate))
		if err != nil {
			return nil, fmt.Errorf("failed to list meals by date range: %w", err)
		}

		resp := &ListMealsResponse{}
		resp.Body.Meals = meals
		return resp, nil
	})

	huma.Get(mealsGroup, "/{user_id}/conversation/{conversation_id}", func(ctx context.Context, input *struct {
		UserID         string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		ConversationID string `path:"conversation_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Conversation ID"`
	}) (*ListMealsResponse, error) {
		meals, err := database.MealLogRepository.ListByConversation(ctx, input.ConversationID)
		if err != nil {
			return nil, fmt.Errorf("failed to list meals by conversation: %w", err)
		}

		resp := &ListMealsResponse{}
		resp.Body.Meals = meals
		return resp, nil
	})

	huma.Get(mealsGroup, "/meal/{meal_id}", func(ctx context.Context, input *struct {
		MealID string `path:"meal_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Meal ID"`
	}) (*GetMealResponse, error) {
		meal, err := database.MealLogRepository.GetByID(ctx, input.MealID)
		if err != nil {
			return nil, fmt.Errorf("failed to get meal: %w", err)
		}

		resp := &GetMealResponse{}
		resp.Body.Meal = meal
		return resp, nil
	})
}

func parseDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}
