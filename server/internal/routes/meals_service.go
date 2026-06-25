package routes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/models"
	"github.com/simhozebs/mugo/internal/runner"
)

func CreateMeal(ctx context.Context, userUUID, convUUID pgtype.UUID, runResult *runner.RunResult, database db.DB) ([]*models.MealLog, error) {
	var batch models.MealsBatchPayload
	if err := json.Unmarshal([]byte(runResult.FinalText), &batch); err != nil {
		return nil, fmt.Errorf("failed to parse nutrition response: %w", err)
	}

	var meals []*models.MealLog
	for _, payload := range batch.Meals {
		mealDate := parseMealDate(payload.Date)
		for i := range payload.Assumptions {
			payload.Assumptions[i].Unit = "g"
		}

		meal, err := database.Meals().Create(ctx,
			userUUID,
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

	return meals, nil
}
