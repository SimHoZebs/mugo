package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	dbgenerated "github.com/simhozebs/mugo/internal/db/dbgenerated"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	"github.com/simhozebs/mugo/internal/models"
)

type mealLogRepository struct {
	queries *dbgenerated.Queries
}

func NewMealLogRepository(queries *dbgenerated.Queries) MealLogRepository {
	return &mealLogRepository{queries: queries}
}

func (r *mealLogRepository) Create(ctx context.Context, userID, conversationID, foodName, mealType string, recordedAt time.Time, macros models.Macros, assumptions []models.Assumption, foodSource string, rawResponse interface{}) (*models.MealLog, error) {
	pgUUID, err := pgutil.ParseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}

	convUUID, err := pgutil.ParseUUIDPtr(conversationID)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation UUID: %w", err)
	}

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
		UserID:         pgUUID,
		ConversationID: convUUID,
		FoodName:       foodName,
		MealType:       mealType,
		RecordedAt:     pgutil.Timestamp(recordedAt),
		Macros:         macrosJSON,
		Assumptions:    assumptionsJSON,
		FoodSource:     foodSource,
		RawResponse:    rawResponseJSON,
	}
	result, err := r.queries.CreateMealLog(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create meal log: %w", err)
	}
	return mapToMealLog(result)
}

func (r *mealLogRepository) GetByID(ctx context.Context, id string) (*models.MealLog, error) {
	pgUUID, err := pgutil.ParseUUID(id)
	if err != nil {
		return nil, err
	}
	result, err := r.queries.GetMealLog(ctx, pgUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal log: %w", err)
	}
	return mapToMealLog(result)
}

func (r *mealLogRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]*models.MealLog, error) {
	pgUUID, err := pgutil.ParseUUID(userID)
	if err != nil {
		return nil, err
	}
	arg := dbgenerated.ListMealLogsByUserParams{
		UserID: pgUUID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	results, err := r.queries.ListMealLogsByUser(ctx, arg)
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

func (r *mealLogRepository) ListByUserAndDate(ctx context.Context, userID string, date time.Time) ([]*models.MealLog, error) {
	pgUUID, err := pgutil.ParseUUID(userID)
	if err != nil {
		return nil, err
	}
	arg := dbgenerated.ListMealLogsByUserAndDateParams{
		UserID:     pgUUID,
		RecordedAt: pgutil.Timestamp(date),
	}
	results, err := r.queries.ListMealLogsByUserAndDate(ctx, arg)
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

func (r *mealLogRepository) ListByUserAndDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*models.MealLog, error) {
	pgUUID, err := pgutil.ParseUUID(userID)
	if err != nil {
		return nil, err
	}
	arg := dbgenerated.ListMealLogsByUserAndDateRangeParams{
		UserID:    pgUUID,
		StartDate: pgutil.Timestamp(startDate),
		EndDate:   pgutil.Timestamp(endDate),
	}
	results, err := r.queries.ListMealLogsByUserAndDateRange(ctx, arg)
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

func (r *mealLogRepository) ListByConversation(ctx context.Context, conversationID string) ([]*models.MealLog, error) {
	pgUUID, err := pgutil.ParseUUID(conversationID)
	if err != nil {
		return nil, err
	}
	results, err := r.queries.ListMealLogsByConversation(ctx, pgUUID)
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

func (r *mealLogRepository) Update(ctx context.Context, id string, foodName, mealType string, macros models.Macros, assumptions []models.Assumption, rawResponse interface{}) (*models.MealLog, error) {
	pgUUID, err := pgutil.ParseUUID(id)
	if err != nil {
		return nil, err
	}

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
		ID:          pgUUID,
		FoodName:    foodName,
		MealType:    mealType,
		Macros:      macrosJSON,
		Assumptions: assumptionsJSON,
		RawResponse: rawResponseJSON,
	}
	result, err := r.queries.UpdateMealLog(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to update meal log: %w", err)
	}
	return mapToMealLog(result)
}

func (r *mealLogRepository) Delete(ctx context.Context, id string) error {
	pgUUID, err := pgutil.ParseUUID(id)
	if err != nil {
		return err
	}
	return r.queries.DeleteMealLog(ctx, pgUUID)
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
