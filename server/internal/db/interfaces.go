package db

import "github.com/simhozebs/mugo/internal/db/repository"

// DB interface defines the accessors for all repositories.
type DB interface {
	Users() repository.UserRepository
	Conversations() repository.ConversationRepository
	Meals() repository.MealLogRepository
	Nutrition() repository.NutritionSummaryRepository
}

// DBProvider interface defines how to obtain a DB connection.
type DBProvider interface {
	GetDatabase() (DB, error)
}
