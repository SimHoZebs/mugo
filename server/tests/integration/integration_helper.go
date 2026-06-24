package integration

import (
	"context"
	"os"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/simhozebs/mugo/internal/runner"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/routes"
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

	// Initialize real database connection
	// We use LazyDatabase which handles migrations on first use
	provider := db.NewLazyDatabase(ctx)

	// Verify connection immediately for early failure if DB is down
	_, err := provider.GetDatabase()
	require.NoError(t, err, "Failed to connect to database for integration test")

	// Create test API
	_, api := humatest.New(t)

	// Register common endpoints
	routes.RegisterUserEndpoints(api, "/users", provider)
	routes.RegisterAnalyticsEndpoints(api, "/analytics", provider)
	routes.RegisterConversationEndpoints(api, "/conversations", provider)

	// Note: Meal endpoints require an AgentRunner which might need mocking
	// or specific setup depending on the test. For now, we skip it here
	// and let specific tests register it if needed, or we can add a mock runner.

	return &TestSuite{
		T:          t,
		Ctx:        ctx,
		DBProvider: provider,
		API:        api,
	}
}

// Teardown cleans up resources used by the test suite.
func (s *TestSuite) Teardown() {
	if s.DBProvider != nil {
		s.DBProvider.Close()
	}
}

// RegisterMeals registers the meal endpoints with specific run and createSession functions.
func (s *TestSuite) RegisterMeals(run runner.RunFunc, createSession runner.CreateSessionFunc) {
	routes.RegisterMealEndpoints(s.API, "/meals", run, createSession, s.DBProvider)
}
