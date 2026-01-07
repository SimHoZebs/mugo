package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	dbgenerated "github.com/simhozebs/mugo/internal/db/dbgenerated"
	"github.com/simhozebs/mugo/internal/models"
)

type NutritionSummaryRepository struct {
	queries *dbgenerated.Queries
}

func NewNutritionSummaryRepository(queries *dbgenerated.Queries) *NutritionSummaryRepository {
	return &NutritionSummaryRepository{queries: queries}
}

func (r *NutritionSummaryRepository) UpsertDaily(ctx context.Context, userID string, date time.Time, totalCalories, totalProtein, totalCarbs, totalFat float64, mealCount int) (*models.DailyNutritionSummary, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	arg := dbgenerated.UpsertDailyNutritionSummaryParams{
		UserID:        pgUUID,
		Date:          pgtype.Date{Time: date, Valid: true},
		TotalCalories: pgtype.Numeric{Int: big.NewInt(int64(totalCalories)), Valid: true},
		TotalProtein:  pgtype.Numeric{Int: big.NewInt(int64(totalProtein)), Valid: true},
		TotalCarbs:    pgtype.Numeric{Int: big.NewInt(int64(totalCarbs)), Valid: true},
		TotalFat:      pgtype.Numeric{Int: big.NewInt(int64(totalFat)), Valid: true},
		MealCount:     int32(mealCount),
	}
	result, err := r.queries.UpsertDailyNutritionSummary(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert daily nutrition summary: %w", err)
	}
	return mapToDailySummary(result), nil
}

func (r *NutritionSummaryRepository) GetDaily(ctx context.Context, userID string, date time.Time) (*models.DailyNutritionSummary, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	arg := dbgenerated.GetDailyNutritionSummaryParams{
		UserID: pgUUID,
		Date:   pgtype.Date{Time: date, Valid: true},
	}
	result, err := r.queries.GetDailyNutritionSummary(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily nutrition summary: %w", err)
	}
	return mapToDailySummary(result), nil
}

func (r *NutritionSummaryRepository) ListDailyByUser(ctx context.Context, userID string, limit, offset int) ([]*models.DailyNutritionSummary, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	arg := dbgenerated.ListDailyNutritionSummariesByUserParams{
		UserID: pgUUID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	results, err := r.queries.ListDailyNutritionSummariesByUser(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list daily nutrition summaries: %w", err)
	}
	summaries := make([]*models.DailyNutritionSummary, len(results))
	for i, s := range results {
		summaries[i] = mapToDailySummary(s)
	}
	return summaries, nil
}

func (r *NutritionSummaryRepository) ListDailyByDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*models.DailyNutritionSummary, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	arg := dbgenerated.ListDailyNutritionSummariesByUserAndDateRangeParams{
		UserID: pgUUID,
		Date:   pgtype.Date{Time: startDate, Valid: true},
	}
	results, err := r.queries.ListDailyNutritionSummariesByUserAndDateRange(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list daily nutrition summaries by date range: %w", err)
	}
	summaries := make([]*models.DailyNutritionSummary, len(results))
	for i, s := range results {
		summaries[i] = mapToDailySummary(s)
	}
	return summaries, nil
}

func (r *NutritionSummaryRepository) UpsertWeekly(ctx context.Context, userID string, weekStartDate time.Time, totalCalories, totalProtein, totalCarbs, totalFat, avgDailyCalories, avgDailyProtein, avgDailyCarbs, avgDailyFat float64, mealCount int) (*models.WeeklyNutritionSummary, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	arg := dbgenerated.UpsertWeeklyNutritionSummaryParams{
		UserID:           pgUUID,
		WeekStartDate:    pgtype.Date{Time: weekStartDate, Valid: true},
		TotalCalories:    pgtype.Numeric{Int: big.NewInt(int64(totalCalories)), Valid: true},
		TotalProtein:     pgtype.Numeric{Int: big.NewInt(int64(totalProtein)), Valid: true},
		TotalCarbs:       pgtype.Numeric{Int: big.NewInt(int64(totalCarbs)), Valid: true},
		TotalFat:         pgtype.Numeric{Int: big.NewInt(int64(totalFat)), Valid: true},
		AvgDailyCalories: pgtype.Numeric{Int: big.NewInt(int64(avgDailyCalories)), Valid: true},
		AvgDailyProtein:  pgtype.Numeric{Int: big.NewInt(int64(avgDailyProtein)), Valid: true},
		AvgDailyCarbs:    pgtype.Numeric{Int: big.NewInt(int64(avgDailyCarbs)), Valid: true},
		AvgDailyFat:      pgtype.Numeric{Int: big.NewInt(int64(avgDailyFat)), Valid: true},
		MealCount:        int32(mealCount),
	}
	result, err := r.queries.UpsertWeeklyNutritionSummary(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert weekly nutrition summary: %w", err)
	}
	return mapToWeeklySummary(result), nil
}

func (r *NutritionSummaryRepository) GetWeekly(ctx context.Context, userID string, weekStartDate time.Time) (*models.WeeklyNutritionSummary, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	arg := dbgenerated.GetWeeklyNutritionSummaryParams{
		UserID:        pgUUID,
		WeekStartDate: pgtype.Date{Time: weekStartDate, Valid: true},
	}
	result, err := r.queries.GetWeeklyNutritionSummary(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly nutrition summary: %w", err)
	}
	return mapToWeeklySummary(result), nil
}

func mapToDailySummary(s dbgenerated.DailyNutritionSummary) *models.DailyNutritionSummary {
	return &models.DailyNutritionSummary{
		ID:            s.ID.String(),
		UserID:        s.UserID.String(),
		Date:          s.Date.Time.Format("2006-01-02"),
		TotalCalories: parseNumeric(s.TotalCalories),
		TotalProtein:  parseNumeric(s.TotalProtein),
		TotalCarbs:    parseNumeric(s.TotalCarbs),
		TotalFat:      parseNumeric(s.TotalFat),
		MealCount:     int(s.MealCount),
		CreatedAt:     s.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:     s.UpdatedAt.Time.Format(time.RFC3339),
	}
}

func mapToWeeklySummary(s dbgenerated.WeeklyNutritionSummary) *models.WeeklyNutritionSummary {
	return &models.WeeklyNutritionSummary{
		ID:               s.ID.String(),
		UserID:           s.UserID.String(),
		WeekStartDate:    s.WeekStartDate.Time.Format("2006-01-02"),
		TotalCalories:    parseNumeric(s.TotalCalories),
		TotalProtein:     parseNumeric(s.TotalProtein),
		TotalCarbs:       parseNumeric(s.TotalCarbs),
		TotalFat:         parseNumeric(s.TotalFat),
		AvgDailyCalories: parseNumeric(s.AvgDailyCalories),
		AvgDailyProtein:  parseNumeric(s.AvgDailyProtein),
		AvgDailyCarbs:    parseNumeric(s.AvgDailyCarbs),
		AvgDailyFat:      parseNumeric(s.AvgDailyFat),
		MealCount:        int(s.MealCount),
		CreatedAt:        s.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:        s.UpdatedAt.Time.Format(time.RFC3339),
	}
}

func parseNumeric(n pgtype.Numeric) float64 {
	if !n.Valid || n.Int == nil {
		return 0
	}
	return float64(n.Int.Int64())
}
