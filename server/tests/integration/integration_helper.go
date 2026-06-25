package integration

import (
	"context"
	"os"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/routes"
	"github.com/simhozebs/mugo/internal/runner"
	"github.com/stretchr/testify/require"
)

// TestSuite represents the context for integration tests, including database and API.
type TestSuite struct {
	T          *testing.T
	Ctx        context.Context
	DBProvider *db.LazyDatabase
	API        humatest.TestAPI
}

// SetupTestSuite initializes the dependencies for an integration test.
// It skips the test if DATABASE_URL is not set.
func SetupTestSuite(t *testing.T) *TestSuite {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("Skipping integration test: DATABASE_URL not set")
	}

	ctx := context.Background()

	provider := db.NewLazyDatabase(ctx)

	_, err := provider.GetDatabase()
	require.NoError(t, err, "Failed to connect to database for integration test")

	_, api := humatest.New(t)

	routes.RegisterUserEndpoints(api, "/users", provider)
	routes.RegisterAnalyticsEndpoints(api, "/analytics", provider)
	routes.RegisterLoggingSessionEndpoints(api, "/loggingsessions", provider)

	return &TestSuite{
		T:          t,
		Ctx:        ctx,
		DBProvider: provider,
		API:        api,
	}
}

func (s *TestSuite) Teardown() {
	if s.DBProvider != nil {
		s.DBProvider.Close()
	}
}

func (s *TestSuite) RegisterMeals(mealRunner *runner.AgentRunner) {
	routes.RegisterMealEndpoints(s.API, "/meals", mealRunner, s.DBProvider)
}
