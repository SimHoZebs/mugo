package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	dbgenerated "github.com/simhozebs/mugo/internal/db/dbgenerated"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	"github.com/simhozebs/mugo/internal/models"
)

type conversationRepository struct {
	queries *dbgenerated.Queries
}

func NewConversationRepository(queries *dbgenerated.Queries) ConversationRepository {
	return &conversationRepository{queries: queries}
}

func (r *conversationRepository) Create(ctx context.Context, userID pgtype.UUID, sessionID, title string) (*models.Conversation, error) {
	if sessionID == "" {
		sessionID = pgutil.GenerateUUID()
	}

	arg := dbgenerated.CreateConversationParams{
		UserID:    userID,
		SessionID: sessionID,
		Title:     pgutil.Text(title),
	}
	result, err := r.queries.CreateConversation(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}
	return mapToConversation(result), nil
}

func (r *conversationRepository) GetByID(ctx context.Context, id pgtype.UUID) (*models.Conversation, error) {
	result, err := r.queries.GetConversation(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	return mapToConversation(result), nil
}

func (r *conversationRepository) GetBySessionID(ctx context.Context, userID pgtype.UUID, sessionID string) (*models.Conversation, error) {
	arg := dbgenerated.GetConversationBySessionIDParams{
		UserID:    userID,
		SessionID: sessionID,
	}
	result, err := r.queries.GetConversationBySessionID(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation by session ID: %w", err)
	}
	return mapToConversation(result), nil
}

func (r *conversationRepository) ListByUser(ctx context.Context, userID pgtype.UUID) ([]*models.Conversation, error) {
	results, err := r.queries.ListConversationsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	conversations := make([]*models.Conversation, len(results))
	for i, c := range results {
		conversations[i] = mapToConversation(c)
	}
	return conversations, nil
}

func (r *conversationRepository) UpdateTitle(ctx context.Context, id pgtype.UUID, title string) (*models.Conversation, error) {
	arg := dbgenerated.UpdateConversationTitleParams{
		ID:    id,
		Title: pgutil.Text(title),
	}
	result, err := r.queries.UpdateConversationTitle(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to update conversation title: %w", err)
	}
	return mapToConversation(result), nil
}

func (r *conversationRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	return r.queries.DeleteConversation(ctx, id)
}

func mapToConversation(c dbgenerated.Conversation) *models.Conversation {
	return &models.Conversation{
		ID:        c.ID.String(),
		UserID:    c.UserID.String(),
		SessionID: c.SessionID,
		Title:     pgutil.FromText(c.Title),
		CreatedAt: c.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Time.Format(time.RFC3339),
	}
}
