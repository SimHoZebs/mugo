package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/adk"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/models"
)

type CreateMealRequest struct {
	Body struct {
		UserID      string `json:"user_id" example:"user-123" doc:"User ID"`
		SessionID   string `json:"session_id" example:"session-456" doc:"Session ID for conversation tracking"`
		Description string `json:"description" example:"I ate a chicken sandwich" doc:"Description of the meal"`
	}
}

type CreateMealResponse struct {
	Body struct {
		Meal *models.MealLog `json:"meal"`
	}
}

type UpdateMealRequest struct {
	MealID string `path:"meal_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Meal ID to update"`
	Body   struct {
		Correction string `json:"correction" example:"actually it was whole wheat bread" doc:"Correction or refinement to the meal"`
	}
}

type UpdateMealResponse struct {
	Body struct {
		Meal *models.MealLog `json:"meal"`
	}
}

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

const mealLogTags = "Logs"

func RegisterMealEndpoints(humaAPI huma.API, prefix string, macroRunner adk.AgentRunner, provider db.DBProvider) {
	mealsGroup := huma.NewGroup(humaAPI, prefix)

	huma.Register(mealsGroup, huma.Operation{
		OperationID: "create-meal-log",
		Method:      "POST",
		Path:        "",
		Summary:     "Create a new meal-log",
		Tags:        []string{mealLogTags},
	}, func(ctx context.Context, input *CreateMealRequest) (*CreateMealResponse, error) {
		fmt.Printf("Creating meal: %s (user: %s, session: %s)\n",
			input.Body.Description, input.Body.UserID, input.Body.SessionID)

		result, err := macroRunner.Run(ctx, input.Body.UserID, input.Body.SessionID, input.Body.Description)
		if err != nil {
			return nil, fmt.Errorf("nutrition agent processing failed: %w", err)
		}

		var payload models.NutritionPayload
		if err := json.Unmarshal([]byte(result.FinalText), &payload); err != nil {
			return nil, fmt.Errorf("failed to parse nutrition response: %w", err)
		}

		database, err := provider.GetDatabase()
		if err != nil {
			return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
		}

		meal, err := database.Meals().Create(ctx,
			input.Body.UserID,
			input.Body.SessionID,
			payload.Name,
			string(payload.MealType),
			time.Now(),
			payload.Macros,
			payload.Assumptions,
			"ai_estimated",
			payload,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create meal: %w", err)
		}

		resp := &CreateMealResponse{}
		resp.Body.Meal = meal
		return resp, nil
	})

	huma.Register(mealsGroup, huma.Operation{
		OperationID: "update-meal",
		Method:      "PUT",
		Path:        "/{meal_id}",
		Summary:     "Update/correct an existing meal",
		Tags:        []string{mealLogTags},
	}, func(ctx context.Context, input *UpdateMealRequest) (*UpdateMealResponse, error) {
		database, err := provider.GetDatabase()
		if err != nil {
			return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
		}

		var updatedMeal *models.MealLog
		txErr := database.WithTx(ctx, func(ctx context.Context, txDB *db.TxDatabase) error {
			meal, err := txDB.MealLogRepository.GetByID(ctx, input.MealID)
			if err != nil {
				return fmt.Errorf("failed to get meal: %w", err)
			}

			sessionID := ""
			if meal.ConversationID != nil {
				sessionID = *meal.ConversationID
			}

			fmt.Printf("Updating meal %s with correction: %s (session: %s)\n",
				input.MealID, input.Body.Correction, sessionID)

			result, err := macroRunner.Run(ctx, meal.UserID, sessionID, input.Body.Correction)
			if err != nil {
				return fmt.Errorf("nutrition agent processing failed: %w", err)
			}

			var payload models.NutritionPayload
			if err := json.Unmarshal([]byte(result.FinalText), &payload); err != nil {
				return fmt.Errorf("failed to parse nutrition response: %w", err)
			}

			newMeal, err := txDB.MealLogRepository.Update(ctx,
				input.MealID,
				payload.Name,
				string(payload.MealType),
				payload.Macros,
				payload.Assumptions,
				payload,
			)
			if err != nil {
				return fmt.Errorf("failed to update meal: %w", err)
			}
			updatedMeal = newMeal
			return nil
		})

		if txErr != nil {
			return nil, txErr
		}

		resp := &UpdateMealResponse{}
		resp.Body.Meal = updatedMeal
		return resp, nil
	})

	huma.Register(mealsGroup, huma.Operation{
		OperationID: "list-meals-by-user",
		Method:      "GET",
		Path:        "/{user_id}",
		Summary:     "List meals for a user",
		Tags:        []string{mealLogTags},
	}, func(ctx context.Context, input *struct {
		UserID string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		Limit  int    `query:"limit" default:"50" doc:"Maximum number of meals to return"`
		Offset int    `query:"offset" default:"0" doc:"Number of meals to skip"`
	}) (*ListMealsResponse, error) {
		database, err := provider.GetDatabase()
		if err != nil {
			return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
		}

		meals, err := database.Meals().ListByUser(ctx, input.UserID, input.Limit, input.Offset)
		if err != nil {
			return nil, fmt.Errorf("failed to list meals: %w", err)
		}

		resp := &ListMealsResponse{}
		resp.Body.Meals = meals
		return resp, nil
	})

	huma.Register(mealsGroup, huma.Operation{
		OperationID: "list-meals-by-date",
		Method:      "GET",
		Path:        "/{user_id}/date/{date}",
		Summary:     "List meals for a user on a specific date",
		Tags:        []string{mealLogTags},
	}, func(ctx context.Context, input *struct {
		UserID string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		Date   string `path:"date" example:"2025-01-07" doc:"Date (YYYY-MM-DD)"`
	}) (*ListMealsResponse, error) {
		database, err := provider.GetDatabase()
		if err != nil {
			return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
		}

		date, err := parseDate(input.Date)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid date format. Expected YYYY-MM-DD", err)
		}

		meals, err := database.Meals().ListByUserAndDate(ctx, input.UserID, date)
		if err != nil {
			return nil, fmt.Errorf("failed to list meals by date: %w", err)
		}

		resp := &ListMealsResponse{}
		resp.Body.Meals = meals
		return resp, nil
	})

	huma.Register(mealsGroup, huma.Operation{
		OperationID: "list-meals-by-range",
		Method:      "GET",
		Path:        "/{user_id}/range",
		Summary:     "List meals for a user in a date range",
		Tags:        []string{mealLogTags},
	}, func(ctx context.Context, input *struct {
		UserID    string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		StartDate string `query:"start_date" example:"2025-01-01" doc:"Start date (YYYY-MM-DD)"`
		EndDate   string `query:"end_date" example:"2025-01-31" doc:"End date (YYYY-MM-DD)"`
	}) (*ListMealsResponse, error) {
		database, err := provider.GetDatabase()
		if err != nil {
			return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
		}

		start, err := parseDate(input.StartDate)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid start_date format. Expected YYYY-MM-DD", err)
		}

		end, err := parseDate(input.EndDate)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid end_date format. Expected YYYY-MM-DD", err)
		}

		meals, err := database.Meals().ListByUserAndDateRange(ctx, input.UserID, start, end)
		if err != nil {
			return nil, fmt.Errorf("failed to list meals by date range: %w", err)
		}

		resp := &ListMealsResponse{}
		resp.Body.Meals = meals
		return resp, nil
	})

	huma.Register(mealsGroup, huma.Operation{
		OperationID: "list-meals-by-conversation",
		Method:      "GET",
		Path:        "/{user_id}/conversation/{conversation_id}",
		Summary:     "List meals for a user in a conversation",
		Tags:        []string{mealLogTags},
	}, func(ctx context.Context, input *struct {
		UserID         string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		ConversationID string `path:"conversation_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Conversation ID"`
	}) (*ListMealsResponse, error) {
		database, err := provider.GetDatabase()
		if err != nil {
			return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
		}

		meals, err := database.Meals().ListByConversation(ctx, input.ConversationID)
		if err != nil {
			return nil, fmt.Errorf("failed to list meals by conversation: %w", err)
		}

		resp := &ListMealsResponse{}
		resp.Body.Meals = meals
		return resp, nil
	})

	huma.Register(mealsGroup, huma.Operation{
		OperationID: "get-meal-by-id",
		Method:      "GET",
		Path:        "/meal/{meal_id}",
		Summary:     "Get a meal by ID",
		Tags:        []string{mealLogTags},
	}, func(ctx context.Context, input *struct {
		MealID string `path:"meal_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Meal ID"`
	}) (*GetMealResponse, error) {
		database, err := provider.GetDatabase()
		if err != nil {
			return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
		}

		meal, err := database.Meals().GetByID(ctx, input.MealID)
		if err != nil {
			return nil, fmt.Errorf("failed to get meal: %w", err)
		}

		resp := &GetMealResponse{}
		resp.Body.Meal = meal
		return resp, nil
	})
}
