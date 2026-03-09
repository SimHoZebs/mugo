package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/simhozebs/mugo/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id pgtype.UUID) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Exists(ctx context.Context, username string) (bool, error)
	List(ctx context.Context) ([]*models.User, error)
	Update(ctx context.Context, id pgtype.UUID, username string) (*models.User, error)
	Delete(ctx context.Context, id pgtype.UUID) error
}

type ConversationRepository interface {
	Create(ctx context.Context, userID pgtype.UUID, sessionID, title string) (*models.Conversation, error)
	GetByID(ctx context.Context, id pgtype.UUID) (*models.Conversation, error)
	GetBySessionID(ctx context.Context, userID pgtype.UUID, sessionID string) (*models.Conversation, error)
	ListByUser(ctx context.Context, userID pgtype.UUID) ([]*models.Conversation, error)
	UpdateTitle(ctx context.Context, id pgtype.UUID, title string) (*models.Conversation, error)
	Delete(ctx context.Context, id pgtype.UUID) error
}

type MealLogRepository interface {
	Create(ctx context.Context, userID, conversationID pgtype.UUID, foodName, mealType string, recordedAt time.Time, macros models.Macros, assumptions []models.Assumption, foodSource string, rawResponse interface{}) (*models.MealLog, error)
	GetByID(ctx context.Context, id pgtype.UUID) (*models.MealLog, error)
	GetByIDWithSession(ctx context.Context, id pgtype.UUID) (*models.MealLog, string, error)
	ListByUser(ctx context.Context, userID pgtype.UUID, limit, offset int) ([]*models.MealLog, error)
	ListByUserAndDate(ctx context.Context, userID pgtype.UUID, date time.Time) ([]*models.MealLog, error)
	ListByUserAndDateRange(ctx context.Context, userID pgtype.UUID, startDate, endDate time.Time) ([]*models.MealLog, error)
	ListByConversation(ctx context.Context, conversationID pgtype.UUID) ([]*models.MealLog, error)
	Update(ctx context.Context, id pgtype.UUID, foodName, mealType string, macros models.Macros, assumptions []models.Assumption, rawResponse interface{}) (*models.MealLog, error)
	Delete(ctx context.Context, id pgtype.UUID) error
}

type NutritionSummaryRepository interface {
	UpsertDaily(ctx context.Context, userID pgtype.UUID, date time.Time, totalCalories, totalProtein, totalCarbs, totalFat float64, mealCount int) (*models.DailyNutritionSummary, error)
	GetDaily(ctx context.Context, userID pgtype.UUID, date time.Time) (*models.DailyNutritionSummary, error)
	ListDailyByUser(ctx context.Context, userID pgtype.UUID, limit, offset int) ([]*models.DailyNutritionSummary, error)
	ListDailyByDateRange(ctx context.Context, userID pgtype.UUID, startDate, endDate time.Time) ([]*models.DailyNutritionSummary, error)
	UpsertWeekly(ctx context.Context, userID pgtype.UUID, weekStartDate time.Time, totalCalories, totalProtein, totalCarbs, totalFat, avgDailyCalories, avgDailyProtein, avgDailyCarbs, avgDailyFat float64, mealCount int) (*models.WeeklyNutritionSummary, error)
	GetWeekly(ctx context.Context, userID pgtype.UUID, weekStartDate time.Time) (*models.WeeklyNutritionSummary, error)
	ListWeeklyByDateRange(ctx context.Context, userID pgtype.UUID, startDate, endDate time.Time) ([]*models.WeeklyNutritionSummary, error)
}
