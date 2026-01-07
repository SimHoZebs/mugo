package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	dbgenerated "github.com/simhozebs/mugo/internal/db/dbgenerated"
	"github.com/simhozebs/mugo/internal/models"
)

type MealLogRepository struct {
	queries *dbgenerated.Queries
}

func NewMealLogRepository(queries *dbgenerated.Queries) *MealLogRepository {
	return &MealLogRepository{queries: queries}
}

func (r *MealLogRepository) Create(ctx context.Context, userID, conversationID, foodName, mealType string, recordedAt time.Time, macros models.Macros, assumptions []models.Assumption, foodSource string, rawResponse interface{}) (*models.MealLog, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}

	var convUUID pgtype.UUID
	if conversationID != "" {
		parsedConvUUID, err := uuid.Parse(conversationID)
		if err != nil {
			return nil, fmt.Errorf("invalid conversation UUID: %w", err)
		}
		convUUID = pgtype.UUID{
			Bytes: [16]byte(parsedConvUUID),
			Valid: true,
		}
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
		RecordedAt:     pgtype.Timestamptz{Time: recordedAt, Valid: true},
		Macros:         macrosJSON,
		Assumptions:    assumptionsJSON,
		FoodSource:     foodSource,
		RawResponse:    rawResponseJSON,
	}
	result, err := r.queries.CreateMealLog(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create meal log: %w", err)
	}
	return mapToMealLog(result), nil
}

func (r *MealLogRepository) GetByID(ctx context.Context, id string) (*models.MealLog, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	result, err := r.queries.GetMealLog(ctx, pgUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal log: %w", err)
	}
	return mapToMealLog(result), nil
}

func (r *MealLogRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]*models.MealLog, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
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
		mealLogs[i] = mapToMealLog(m)
	}
	return mealLogs, nil
}

func (r *MealLogRepository) ListByUserAndDate(ctx context.Context, userID string, date time.Time) ([]*models.MealLog, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	arg := dbgenerated.ListMealLogsByUserAndDateParams{
		UserID:     pgUUID,
		RecordedAt: pgtype.Timestamptz{Time: date, Valid: true},
	}
	results, err := r.queries.ListMealLogsByUserAndDate(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list meal logs by date: %w", err)
	}
	mealLogs := make([]*models.MealLog, len(results))
	for i, m := range results {
		mealLogs[i] = mapToMealLog(m)
	}
	return mealLogs, nil
}

func (r *MealLogRepository) ListByUserAndDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*models.MealLog, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	arg := dbgenerated.ListMealLogsByUserAndDateRangeParams{
		UserID:     pgUUID,
		RecordedAt: pgtype.Timestamptz{Time: startDate, Valid: true},
	}
	results, err := r.queries.ListMealLogsByUserAndDateRange(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list meal logs by date range: %w", err)
	}
	mealLogs := make([]*models.MealLog, len(results))
	for i, m := range results {
		mealLogs[i] = mapToMealLog(m)
	}
	return mealLogs, nil
}

func (r *MealLogRepository) ListByConversation(ctx context.Context, conversationID string) ([]*models.MealLog, error) {
	parsedUUID, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	results, err := r.queries.ListMealLogsByConversation(ctx, pgUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to list meal logs by conversation: %w", err)
	}
	mealLogs := make([]*models.MealLog, len(results))
	for i, m := range results {
		mealLogs[i] = mapToMealLog(m)
	}
	return mealLogs, nil
}

func (r *MealLogRepository) Delete(ctx context.Context, id string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	return r.queries.DeleteMealLog(ctx, pgUUID)
}

func mapToMealLog(m dbgenerated.MealLog) *models.MealLog {
	var macros models.Macros
	if m.Macros != nil {
		json.Unmarshal(m.Macros, &macros)
	}

	var assumptions []models.Assumption
	if m.Assumptions != nil {
		json.Unmarshal(m.Assumptions, &assumptions)
	}

	var conversationID *string
	if m.ConversationID.Valid {
		s := m.ConversationID.String()
		conversationID = &s
	}

	var rawResponse interface{}
	if m.RawResponse != nil {
		json.Unmarshal(m.RawResponse, &rawResponse)
	}

	return &models.MealLog{
		ID:             m.ID.String(),
		UserID:         m.UserID.String(),
		ConversationID: conversationID,
		FoodName:       m.FoodName,
		MealType:       string(m.MealType.(string)),
		RecordedAt:     m.RecordedAt.Time.Format(time.RFC3339),
		Macros:         macros,
		Assumptions:    assumptions,
		FoodSource:     string(m.FoodSource.(string)),
		RawResponse:    rawResponse,
		CreatedAt:      m.CreatedAt.Time.Format(time.RFC3339),
	}
}
