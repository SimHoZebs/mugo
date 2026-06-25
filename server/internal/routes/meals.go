package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	"github.com/simhozebs/mugo/internal/models"
	"github.com/simhozebs/mugo/internal/runner"
	adkrunner "google.golang.org/adk/runner"
	"google.golang.org/adk/session"
)

type CreateMealRequest struct {
	Body struct {
		UserID      string `json:"user_id" example:"user-123" doc:"User ID"`
		SessionID   string `json:"session_id,omitempty" example:"session-456" doc:"Optional Session ID. If not provided, a new logging session will be created."`
		Description string `json:"description" example:"I ate a chicken sandwich" doc:"Description of the meal"`
	}
}

type CreateMealResponse struct {
	Body struct {
		SessionID string            `json:"session_id" doc:"The session ID used for this request"`
		Meals     []*models.MealLog `json:"meals"`
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

func RegisterMealEndpoints(humaAPI huma.API, prefix string, mealRunner *adkrunner.Runner, sessionService session.Service, appName string, provider db.DBProvider) {
	mealsGroup := huma.NewGroup(humaAPI, prefix)

	huma.Register(mealsGroup, huma.Operation{
		OperationID: "create-meal-log",
		Method:      "POST",
		Path:        "",
		Summary:     "Create a new meal-log",
		Tags:        []string{mealLogTags},
	}, func(ctx context.Context, input *CreateMealRequest) (*CreateMealResponse, error) {
		if mealRunner == nil {
			return nil, huma.Error503ServiceUnavailable("AI meal parsing is currently disabled (missing configuration)")
		}

		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		userUUID, err := pgutil.ParseUUID(input.Body.UserID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid user ID", err)
		}

		conv, err := database.LoggingSessions().Create(ctx, userUUID, input.Body.SessionID, "New Meal Log")
		if err != nil {
			return nil, fmt.Errorf("failed to create new logging session: %w", err)
		}

		convUUID, err := pgutil.ParseUUID(conv.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse logging session ID: %w", err)
		}

		if err := runner.CreateSession(ctx, sessionService, appName, input.Body.UserID, conv.SessionID); err != nil {
			return nil, fmt.Errorf("failed to create ADK session: %w", err)
		}

		today := time.Now().Format("2006-01-02")
		message := fmt.Sprintf("Today's date is %s. %s", today, input.Body.Description)

		runResult, err := runner.Run(ctx, mealRunner, input.Body.UserID, conv.SessionID, message)
		if err != nil {
			return nil, fmt.Errorf("nutrition agent processing failed: %w", err)
		}

		meals, err := CreateMeal(ctx, userUUID, convUUID, runResult, database)
		if err != nil {
			return nil, err
		}

		resp := &CreateMealResponse{}
		resp.Body.SessionID = conv.SessionID
		resp.Body.Meals = meals
		return resp, nil
	})

	huma.Register(mealsGroup, huma.Operation{
		OperationID: "update-meal",
		Method:      "PUT",
		Path:        "/{meal_id}",
		Summary:     "Update/correct an existing meal",
		Tags:        []string{mealLogTags},
	}, func(ctx context.Context, input *UpdateMealRequest) (*UpdateMealResponse, error) {
		db, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		mealUUID, err := pgutil.ParseUUID(input.MealID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid meal ID", err)
		}

		meal, adkSessionID, err := db.Meals().GetByIDWithSession(ctx, mealUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get meal: %w", err)
		}

		fmt.Printf("Updating meal %s with correction: %s (session: %s)\n",
			input.MealID, input.Body.Correction, adkSessionID)

		if mealRunner == nil {
			return nil, huma.Error503ServiceUnavailable("AI meal parsing is currently disabled (missing configuration)")
		}

		result, err := runner.Run(ctx, mealRunner, meal.UserID, adkSessionID, input.Body.Correction)
		if err != nil {
			return nil, fmt.Errorf("nutrition agent processing failed: %w", err)
		}

		var batch models.MealsBatchPayload
		if err := json.Unmarshal([]byte(result.FinalText), &batch); err != nil {
			return nil, huma.Error422UnprocessableEntity("failed to parse nutrition response", err)
		}
		if len(batch.Meals) == 0 {
			return nil, huma.Error422UnprocessableEntity("nutrition agent returned no meals")
		}
		payload := batch.Meals[0]

		// Apply default units since LLM now implies grams
		for i := range payload.Assumptions {
			payload.Assumptions[i].Unit = "g"
		}

		var updatedMeal *models.MealLog
		newMeal, err := db.Meals().Update(ctx,
			mealUUID,
			payload.Name,
			string(payload.MealType),
			payload.Macros,
			payload.Assumptions,
			payload,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update meal: %w", err)
		}
		updatedMeal = newMeal

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
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		userUUID, err := pgutil.ParseUUID(input.UserID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid user ID", err)
		}

		meals, err := database.Meals().ListByUser(ctx, userUUID, input.Limit, input.Offset)
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
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		userUUID, err := pgutil.ParseUUID(input.UserID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid user ID", err)
		}

		date, err := parseDate(input.Date)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid date format. Expected YYYY-MM-DD", err)
		}

		meals, err := database.Meals().ListByUserAndDate(ctx, userUUID, date)
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
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		userUUID, err := pgutil.ParseUUID(input.UserID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid user ID", err)
		}

		start, err := parseDate(input.StartDate)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid start_date format. Expected YYYY-MM-DD", err)
		}

		end, err := parseDate(input.EndDate)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid end_date format. Expected YYYY-MM-DD", err)
		}

		meals, err := database.Meals().ListByUserAndDateRange(ctx, userUUID, start, end)
		if err != nil {
			return nil, fmt.Errorf("failed to list meals by date range: %w", err)
		}

		resp := &ListMealsResponse{}
		resp.Body.Meals = meals
		return resp, nil
	})

	huma.Register(mealsGroup, huma.Operation{
		OperationID: "list-meals-by-logging-session",
		Method:      "GET",
		Path:        "/{user_id}/logging_session/{logging_session_id}",
		Summary:     "List meals for a user in a logging session",
		Tags:        []string{mealLogTags},
	}, func(ctx context.Context, input *struct {
		UserID           string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		LoggingSessionID string `path:"logging_session_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Logging session ID"`
	}) (*ListMealsResponse, error) {
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		sessionUUID, err := pgutil.ParseUUID(input.LoggingSessionID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid logging session ID", err)
		}

		meals, err := database.Meals().ListByLoggingSession(ctx, sessionUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to list meals by logging session: %w", err)
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
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		mealUUID, err := pgutil.ParseUUID(input.MealID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid meal ID", err)
		}

		meal, err := database.Meals().GetByID(ctx, mealUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get meal: %w", err)
		}

		resp := &GetMealResponse{}
		resp.Body.Meal = meal
		return resp, nil
	})
}
