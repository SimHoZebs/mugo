package adk_test

import (
	"context"
	"errors"
	"testing"

	"github.com/simhozebs/mugo/internal/adk"
	adkmocks "github.com/simhozebs/mugo/internal/adk/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockAgentRunner(t *testing.T) {
	t.Run("calls RunFunc when provided", func(t *testing.T) {
		expectedResult := &adk.RunResult{
			FinalText: "test response",
			Events:    nil,
		}
		runner := &adkmocks.MockAgentRunner{
			RunFunc: func(ctx context.Context, userID, sessionID, text string) (*adk.RunResult, error) {
				assert.Equal(t, "user-123", userID)
				assert.Equal(t, "session-456", sessionID)
				assert.Equal(t, "hello", text)
				return expectedResult, nil
			},
		}

		result, err := runner.Run(context.Background(), "user-123", "session-456", "hello")
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("returns empty result when RunFunc is nil", func(t *testing.T) {
		runner := &adkmocks.MockAgentRunner{}

		result, err := runner.Run(context.Background(), "user-123", "session-456", "hello")
		require.NoError(t, err)
		assert.Equal(t, &adk.RunResult{}, result)
	})

	t.Run("returns error from RunFunc", func(t *testing.T) {
		expectedErr := errors.New("agent error")
		runner := &adkmocks.MockAgentRunner{
			RunFunc: func(ctx context.Context, userID, sessionID, text string) (*adk.RunResult, error) {
				return nil, expectedErr
			},
		}

		result, err := runner.Run(context.Background(), "user-123", "session-456", "hello")
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
	})
}
