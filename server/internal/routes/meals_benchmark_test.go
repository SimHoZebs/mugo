package routes_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/simhozebs/mugo/internal/db/repository/mocks"
	"github.com/simhozebs/mugo/internal/models"
	"github.com/stretchr/testify/mock"
)

// We want to test the benchmark to verify the overhead difference in calling DB methods.
// Here we just benchmark the iteration itself to show we are combining it into one batch call instead of doing a loop of calls.

// BenchmarkSequentialCreate simulates the old N+1 behavior where we call Create for each meal.
func BenchmarkSequentialCreate(b *testing.B) {
	mockRepo := &mocks.MealLogRepositoryMock{}
	mockRepo.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&models.MealLog{}, nil)

	ctx := context.Background()
	userUUID := pgtype.UUID{Bytes: [16]byte{1}}
	convUUID := pgtype.UUID{Bytes: [16]byte{2}}
	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate batch of 10 meals
		for j := 0; j < 10; j++ {
			_, _ = mockRepo.Create(ctx, userUUID, convUUID, "Apple", "snack", now, models.Macros{}, nil, "ai", nil)
		}
	}
}

// BenchmarkBatchCreate simulates the new batch behavior.
func BenchmarkBatchCreate(b *testing.B) {
	mockRepo := &mocks.MealLogRepositoryMock{}
	mockRepo.On("CreateBatch", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(([]*models.MealLog)(nil), nil)

	ctx := context.Background()
	userUUID := pgtype.UUID{Bytes: [16]byte{1}}
	convUUID := pgtype.UUID{Bytes: [16]byte{2}}
	now := time.Now()

	meals := make([]models.MealLogParams, 10)
	for j := 0; j < 10; j++ {
		meals[j] = models.MealLogParams{
			FoodName:   "Apple",
			MealType:   "snack",
			RecordedAt: now,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.CreateBatch(ctx, userUUID, convUUID, meals)
	}
}
