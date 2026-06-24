package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	"github.com/simhozebs/mugo/internal/models"
	"github.com/simhozebs/mugo/internal/runner"
)

type CreateMealInput struct {
	UserUUID    pgtype.UUID
	UserID      string
	SessionID   string
	Description string
}

type CreateMealResult struct {
	SessionID string
	Meals     []*models.MealLog
}

func CreateMeal(ctx context.Context, input CreateMealInput, run runner.RunFunc, createSession runner.CreateSessionFunc, database db.DB) (*CreateMealResult, error) {
	today := time.Now().Format("2006-01-02")
	message := fmt.Sprintf("Today's date is %s. %s", today, input.Description)

	conv, err := database.Conversations().Create(ctx, input.UserUUID, input.SessionID, "New Meal Log")
	if err != nil {
		return nil, fmt.Errorf("failed to create new conversation: %w", err)
	}

	convUUID, err := pgutil.ParseUUID(conv.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conversation ID: %w", err)
	}

	if err := createSession(ctx, input.UserID, conv.SessionID); err != nil {
		return nil, fmt.Errorf("failed to create ADK session: %w", err)
	}

	result, err := run(ctx, input.UserID, conv.SessionID, message)
	if err != nil {
		return nil, fmt.Errorf("nutrition agent processing failed: %w", err)
	}

	var batch models.MealsBatchPayload
	if err := json.Unmarshal([]byte(result.FinalText), &batch); err != nil {
		return nil, fmt.Errorf("failed to parse nutrition response: %w", err)
	}

	var meals []*models.MealLog
	for _, payload := range batch.Meals {
		mealDate := parseMealDate(payload.Date)
		for i := range payload.Assumptions {
			payload.Assumptions[i].Unit = "g"
		}

		meal, err := database.Meals().Create(ctx,
			input.UserUUID,
			convUUID,
			payload.Name,
			string(payload.MealType),
			mealDate,
			payload.Macros,
			payload.Assumptions,
			"ai_estimated",
			payload,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create meal: %w", err)
		}
		meals = append(meals, meal)
	}

	return &CreateMealResult{
		SessionID: conv.SessionID,
		Meals:     meals,
	}, nil
}
