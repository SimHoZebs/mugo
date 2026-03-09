package mocks

import (
	"context"

	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/db/repository"
	"github.com/stretchr/testify/mock"
)

// DBMock is a mock implementation of the DB interface.
type DBMock struct {
	mock.Mock
}

func (m *DBMock) Users() repository.UserRepository {
	args := m.Called()
	return args.Get(0).(repository.UserRepository)
}

func (m *DBMock) Conversations() repository.ConversationRepository {
	args := m.Called()
	return args.Get(0).(repository.ConversationRepository)
}

func (m *DBMock) Meals() repository.MealLogRepository {
	args := m.Called()
	return args.Get(0).(repository.MealLogRepository)
}

func (m *DBMock) Nutrition() repository.NutritionSummaryRepository {
	args := m.Called()
	return args.Get(0).(repository.NutritionSummaryRepository)
}

func (m *DBMock) WithTx(ctx context.Context, fn func(ctx context.Context, txDB db.DB) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// DBProviderMock is a mock implementation of the DBProvider interface.
type DBProviderMock struct {
	mock.Mock
}

func (m *DBProviderMock) GetDatabase() (db.DB, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(db.DB), args.Error(1)
}
