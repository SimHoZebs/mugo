package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	dbgenerated "github.com/simhozebs/mugo/internal/db/dbgenerated"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	"github.com/simhozebs/mugo/internal/models"
)

type mealLogRepository struct {
	queries *dbgenerated.Queries
}

func (repo *mealLogRepository) q(ctx context.Context) *dbgenerated.Queries {
	if tx, ok := pgutil.TxFromContext(ctx); ok {
		return repo.queries.WithTx(tx)
	}
	return repo.queries
}

func NewMealLogRepository(queries *dbgenerated.Queries) MealLogRepository {
	return &mealLogRepository{queries: queries}
}

func (repo *mealLogRepository) Create(ctx context.Context, userID, conversationID pgtype.UUID, foodName, mealType string, recordedAt time.Time, macros models.Macros, assumptions []models.Assumption, foodSource string, rawResponse interface{}) (*models.MealLog, error) {
	macrosJSON, err := json.Marshal(macros)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal macros: %w", err)
	}
	assumptionsJSON, err := json.Marshal(assumptions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal assumptions: %w", err)
	}

	var rawResponseJSON []byte
	if rawResponse != nil {
		rawResponseJSON, err = json.Marshal(rawResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal raw response: %w", err)
		}
	}

	arg := dbgenerated.CreateMealLogParams{
		UserID:         userID,
		ConversationID: conversationID,
		FoodName:       foodName,
		MealType:       mealType,
		RecordedAt:     pgutil.Timestamp(recordedAt),
		Macros:         macrosJSON,
		Assumptions:    assumptionsJSON,
		FoodSource:     foodSource,
		RawResponse:    rawResponseJSON,
	}
	result, err := repo.q(ctx).CreateMealLog(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create meal log: %w", err)
	}
	return mapToMealLog(result)
}

func (repo *mealLogRepository) GetByID(ctx context.Context, id pgtype.UUID) (*models.MealLog, error) {
	result, err := repo.q(ctx).GetMealLog(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal log: %w", err)
	}
	return mapToMealLog(result)
}

func (repo *mealLogRepository) GetByIDWithSession(ctx context.Context, id pgtype.UUID) (*models.MealLog, string, error) {
	result, err := repo.q(ctx).GetMealLogWithSession(ctx, id)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get meal log with session: %w", err)
	}

	// Map GetMealLogWithSessionRow to dbgenerated.MealLog so we can reuse mapToMealLog
	m := dbgenerated.MealLog{
		ID:             result.ID,
		UserID:         result.UserID,
		ConversationID: result.ConversationID,
		FoodName:       result.FoodName,
		MealType:       result.MealType,
		RecordedAt:     result.RecordedAt,
		Macros:         result.Macros,
		Assumptions:    result.Assumptions,
		FoodSource:     result.FoodSource,
		RawResponse:    result.RawResponse,
		CreatedAt:      result.CreatedAt,
	}

	meal, err := mapToMealLog(m)
	if err != nil {
		return nil, "", err
	}

	return meal, result.AdkSessionID.String, nil
}

func (repo *mealLogRepository) ListByUser(ctx context.Context, userID pgtype.UUID, limit, offset int) ([]*models.MealLog, error) {
	arg := dbgenerated.ListMealLogsByUserParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	results, err := repo.q(ctx).ListMealLogsByUser(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list meal logs: %w", err)
	}
	mealLogs := make([]*models.MealLog, len(results))
	for i, m := range results {
		var err error
		mealLogs[i], err = mapToMealLog(m)
		if err != nil {
			return nil, err
		}
	}
	return mealLogs, nil
}

func (repo *mealLogRepository) ListByUserAndDate(ctx context.Context, userID pgtype.UUID, date time.Time) ([]*models.MealLog, error) {
	arg := dbgenerated.ListMealLogsByUserAndDateParams{
		UserID:     userID,
		RecordedAt: pgutil.Timestamp(date),
	}
	results, err := repo.q(ctx).ListMealLogsByUserAndDate(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list meal logs by date: %w", err)
	}
	mealLogs := make([]*models.MealLog, len(results))
	for i, m := range results {
		var err error
		mealLogs[i], err = mapToMealLog(m)
		if err != nil {
			return nil, err
		}
	}
	return mealLogs, nil
}

func (repo *mealLogRepository) ListByUserAndDateRange(ctx context.Context, userID pgtype.UUID, startDate, endDate time.Time) ([]*models.MealLog, error) {
	arg := dbgenerated.ListMealLogsByUserAndDateRangeParams{
		UserID:    userID,
		StartDate: pgutil.Timestamp(startDate),
		EndDate:   pgutil.Timestamp(endDate),
	}
	results, err := repo.q(ctx).ListMealLogsByUserAndDateRange(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list meal logs by date range: %w", err)
	}
	mealLogs := make([]*models.MealLog, len(results))
	for i, m := range results {
		var err error
		mealLogs[i], err = mapToMealLog(m)
		if err != nil {
			return nil, err
		}
	}
	return mealLogs, nil
}

func (repo *mealLogRepository) ListByConversation(ctx context.Context, conversationID pgtype.UUID) ([]*models.MealLog, error) {
	results, err := repo.q(ctx).ListMealLogsByConversation(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list meal logs by conversation: %w", err)
	}
	mealLogs := make([]*models.MealLog, len(results))
	for i, m := range results {
		var err error
		mealLogs[i], err = mapToMealLog(m)
		if err != nil {
			return nil, err
		}
	}
	return mealLogs, nil
}

func (repo *mealLogRepository) Update(ctx context.Context, id pgtype.UUID, foodName, mealType string, macros models.Macros, assumptions []models.Assumption, rawResponse interface{}) (*models.MealLog, error) {
	macrosJSON, err := json.Marshal(macros)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal macros: %w", err)
	}
	assumptionsJSON, err := json.Marshal(assumptions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal assumptions: %w", err)
	}

	var rawResponseJSON []byte
	if rawResponse != nil {
		rawResponseJSON, err = json.Marshal(rawResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal raw response: %w", err)
		}
	}

	arg := dbgenerated.UpdateMealLogParams{
		ID:          id,
		FoodName:    foodName,
		MealType:    mealType,
		Macros:      macrosJSON,
		Assumptions: assumptionsJSON,
		RawResponse: rawResponseJSON,
	}
	result, err := repo.q(ctx).UpdateMealLog(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to update meal log: %w", err)
	}
	return mapToMealLog(result)
}

func (repo *mealLogRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	return repo.q(ctx).DeleteMealLog(ctx, id)
}

func mapToMealLog(m dbgenerated.MealLog) (*models.MealLog, error) {
	var macros models.Macros
	if m.Macros != nil {
		if err := json.Unmarshal(m.Macros, &macros); err != nil {
			return nil, fmt.Errorf("failed to unmarshal macros: %w", err)
		}
	}

	var assumptions []models.Assumption
	if m.Assumptions != nil {
		if err := json.Unmarshal(m.Assumptions, &assumptions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal assumptions: %w", err)
		}
	}

	var rawResponse interface{}
	if m.RawResponse != nil {
		if err := json.Unmarshal(m.RawResponse, &rawResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal raw response: %w", err)
		}
	}

	var conversationID *string
	if m.ConversationID.Valid {
		s := m.ConversationID.String()
		conversationID = &s
	}

	mealType, _ := m.MealType.(string)
	foodSource, _ := m.FoodSource.(string)

	return &models.MealLog{
		ID:             m.ID.String(),
		UserID:         m.UserID.String(),
		ConversationID: conversationID,
		FoodName:       m.FoodName,
		MealType:       mealType,
		RecordedAt:     m.RecordedAt.Time.Format(time.RFC3339),
		Macros:         macros,
		Assumptions:    assumptions,
		FoodSource:     foodSource,
		RawResponse:    rawResponse,
		CreatedAt:      m.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}
