package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	dbgenerated "github.com/simhozebs/mugo/internal/db/dbgenerated"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	"github.com/simhozebs/mugo/internal/models"
)

type nutritionSummaryRepository struct {
	queries *dbgenerated.Queries
}

func NewNutritionSummaryRepository(queries *dbgenerated.Queries) NutritionSummaryRepository {
	return &nutritionSummaryRepository{queries: queries}
}

func (r *nutritionSummaryRepository) UpsertDaily(ctx context.Context, userID pgtype.UUID, date time.Time, totalCalories, totalProtein, totalCarbs, totalFat float64, mealCount int) (*models.DailyNutritionSummary, error) {
	arg := dbgenerated.UpsertDailyNutritionSummaryParams{
		UserID:        userID,
		Date:          pgutil.Date(date),
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

func (r *nutritionSummaryRepository) GetDaily(ctx context.Context, userID pgtype.UUID, date time.Time) (*models.DailyNutritionSummary, error) {
	arg := dbgenerated.GetDailyNutritionSummaryParams{
		UserID: userID,
		Date:   pgutil.Date(date),
	}
	result, err := r.queries.GetDailyNutritionSummary(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily nutrition summary: %w", err)
	}
	return mapToDailySummary(result), nil
}

func (r *nutritionSummaryRepository) ListDailyByUser(ctx context.Context, userID pgtype.UUID, limit, offset int) ([]*models.DailyNutritionSummary, error) {
	arg := dbgenerated.ListDailyNutritionSummariesByUserParams{
		UserID: userID,
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

func (r *nutritionSummaryRepository) ListDailyByDateRange(ctx context.Context, userID pgtype.UUID, startDate, endDate time.Time) ([]*models.DailyNutritionSummary, error) {
	arg := dbgenerated.ListDailyNutritionSummariesByUserAndDateRangeParams{
		UserID:    userID,
		StartDate: pgutil.Date(startDate),
		EndDate:   pgutil.Date(endDate),
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

func (r *nutritionSummaryRepository) UpsertWeekly(ctx context.Context, userID pgtype.UUID, weekStartDate time.Time, totalCalories, totalProtein, totalCarbs, totalFat, avgDailyCalories, avgDailyProtein, avgDailyCarbs, avgDailyFat float64, mealCount int) (*models.WeeklyNutritionSummary, error) {
	arg := dbgenerated.UpsertWeeklyNutritionSummaryParams{
		UserID:           userID,
		WeekStartDate:    pgutil.Date(weekStartDate),
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

func (r *nutritionSummaryRepository) GetWeekly(ctx context.Context, userID pgtype.UUID, weekStartDate time.Time) (*models.WeeklyNutritionSummary, error) {
	arg := dbgenerated.GetWeeklyNutritionSummaryParams{
		UserID:        userID,
		WeekStartDate: pgutil.Date(weekStartDate),
	}
	result, err := r.queries.GetWeeklyNutritionSummary(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly nutrition summary: %w", err)
	}
	return mapToWeeklySummary(result), nil
}

func (r *nutritionSummaryRepository) ListWeeklyByDateRange(ctx context.Context, userID pgtype.UUID, startDate, endDate time.Time) ([]*models.WeeklyNutritionSummary, error) {
	arg := dbgenerated.ListWeeklyNutritionSummariesByUserAndDateRangeParams{
		UserID:    userID,
		StartDate: pgutil.Date(startDate),
		EndDate:   pgutil.Date(endDate),
	}
	results, err := r.queries.ListWeeklyNutritionSummariesByUserAndDateRange(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list weekly nutrition summaries by date range: %w", err)
	}
	summaries := make([]*models.WeeklyNutritionSummary, len(results))
	for i, s := range results {
		summaries[i] = mapToWeeklySummary(s)
	}
	return summaries, nil
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
