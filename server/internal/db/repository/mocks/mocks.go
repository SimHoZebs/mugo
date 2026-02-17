package mocks

import (
	"context"
	"time"

	"github.com/simhozebs/mugo/internal/models"
	"github.com/stretchr/testify/mock"
)

// UserRepositoryMock is a mock implementation of the UserRepository interface.
type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) Exists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *UserRepositoryMock) List(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *UserRepositoryMock) Update(ctx context.Context, id string, username string) (*models.User, error) {
	args := m.Called(ctx, id, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ConversationRepositoryMock is a mock implementation of the ConversationRepository interface.
type ConversationRepositoryMock struct {
	mock.Mock
}

func (m *ConversationRepositoryMock) Create(ctx context.Context, userID, sessionID, title string) (*models.Conversation, error) {
	args := m.Called(ctx, userID, sessionID, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Conversation), args.Error(1)
}

func (m *ConversationRepositoryMock) GetByID(ctx context.Context, id string) (*models.Conversation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Conversation), args.Error(1)
}

func (m *ConversationRepositoryMock) GetBySessionID(ctx context.Context, userID, sessionID string) (*models.Conversation, error) {
	args := m.Called(ctx, userID, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Conversation), args.Error(1)
}

func (m *ConversationRepositoryMock) ListByUser(ctx context.Context, userID string) ([]*models.Conversation, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Conversation), args.Error(1)
}

func (m *ConversationRepositoryMock) UpdateTitle(ctx context.Context, id, title string) (*models.Conversation, error) {
	args := m.Called(ctx, id, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Conversation), args.Error(1)
}

func (m *ConversationRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MealLogRepositoryMock is a mock implementation of the MealLogRepository interface.
type MealLogRepositoryMock struct {
	mock.Mock
}

func (m *MealLogRepositoryMock) Create(ctx context.Context, userID, conversationID, foodName, mealType string, recordedAt time.Time, macros models.Macros, assumptions []models.Assumption, foodSource string, rawResponse interface{}) (*models.MealLog, error) {
	args := m.Called(ctx, userID, conversationID, foodName, mealType, recordedAt, macros, assumptions, foodSource, rawResponse)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MealLog), args.Error(1)
}

func (m *MealLogRepositoryMock) GetByID(ctx context.Context, id string) (*models.MealLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MealLog), args.Error(1)
}

func (m *MealLogRepositoryMock) ListByUser(ctx context.Context, userID string, limit, offset int) ([]*models.MealLog, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.MealLog), args.Error(1)
}

func (m *MealLogRepositoryMock) ListByUserAndDate(ctx context.Context, userID string, date time.Time) ([]*models.MealLog, error) {
	args := m.Called(ctx, userID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.MealLog), args.Error(1)
}

func (m *MealLogRepositoryMock) ListByUserAndDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*models.MealLog, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.MealLog), args.Error(1)
}

func (m *MealLogRepositoryMock) ListByConversation(ctx context.Context, conversationID string) ([]*models.MealLog, error) {
	args := m.Called(ctx, conversationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.MealLog), args.Error(1)
}

func (m *MealLogRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// NutritionSummaryRepositoryMock is a mock implementation of the NutritionSummaryRepository interface.
type NutritionSummaryRepositoryMock struct {
	mock.Mock
}

func (m *NutritionSummaryRepositoryMock) UpsertDaily(ctx context.Context, userID string, date time.Time, totalCalories, totalProtein, totalCarbs, totalFat float64, mealCount int) (*models.DailyNutritionSummary, error) {
	args := m.Called(ctx, userID, date, totalCalories, totalProtein, totalCarbs, totalFat, mealCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DailyNutritionSummary), args.Error(1)
}

func (m *NutritionSummaryRepositoryMock) GetDaily(ctx context.Context, userID string, date time.Time) (*models.DailyNutritionSummary, error) {
	args := m.Called(ctx, userID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DailyNutritionSummary), args.Error(1)
}

func (m *NutritionSummaryRepositoryMock) ListDailyByUser(ctx context.Context, userID string, limit, offset int) ([]*models.DailyNutritionSummary, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DailyNutritionSummary), args.Error(1)
}

func (m *NutritionSummaryRepositoryMock) ListDailyByDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*models.DailyNutritionSummary, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DailyNutritionSummary), args.Error(1)
}

func (m *NutritionSummaryRepositoryMock) UpsertWeekly(ctx context.Context, userID string, weekStartDate time.Time, totalCalories, totalProtein, totalCarbs, totalFat, avgDailyCalories, avgDailyProtein, avgDailyCarbs, avgDailyFat float64, mealCount int) (*models.WeeklyNutritionSummary, error) {
	args := m.Called(ctx, userID, weekStartDate, totalCalories, totalProtein, totalCarbs, totalFat, avgDailyCalories, avgDailyProtein, avgDailyCarbs, avgDailyFat, mealCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WeeklyNutritionSummary), args.Error(1)
}

func (m *NutritionSummaryRepositoryMock) GetWeekly(ctx context.Context, userID string, weekStartDate time.Time) (*models.WeeklyNutritionSummary, error) {
	args := m.Called(ctx, userID, weekStartDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WeeklyNutritionSummary), args.Error(1)
}

func (m *NutritionSummaryRepositoryMock) ListWeeklyByDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*models.WeeklyNutritionSummary, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.WeeklyNutritionSummary), args.Error(1)
}
