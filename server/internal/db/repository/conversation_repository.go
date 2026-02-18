package repository

import (
	"context"
	"fmt"
	"time"

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

func (r *conversationRepository) Create(ctx context.Context, userID, sessionID, title string) (*models.Conversation, error) {
	pgUUID, err := pgutil.ParseUUID(userID)
	if err != nil {
		return nil, err
	}
	arg := dbgenerated.CreateConversationParams{
		UserID:    pgUUID,
		SessionID: sessionID,
		Title:     pgutil.Text(title),
	}
	result, err := r.queries.CreateConversation(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}
	return mapToConversation(result), nil
}

func (r *conversationRepository) GetByID(ctx context.Context, id string) (*models.Conversation, error) {
	pgUUID, err := pgutil.ParseUUID(id)
	if err != nil {
		return nil, err
	}
	result, err := r.queries.GetConversation(ctx, pgUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	return mapToConversation(result), nil
}

func (r *conversationRepository) GetBySessionID(ctx context.Context, userID, sessionID string) (*models.Conversation, error) {
	pgUUID, err := pgutil.ParseUUID(userID)
	if err != nil {
		return nil, err
	}
	arg := dbgenerated.GetConversationBySessionIDParams{
		UserID:    pgUUID,
		SessionID: sessionID,
	}
	result, err := r.queries.GetConversationBySessionID(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation by session ID: %w", err)
	}
	return mapToConversation(result), nil
}

func (r *conversationRepository) ListByUser(ctx context.Context, userID string) ([]*models.Conversation, error) {
	pgUUID, err := pgutil.ParseUUID(userID)
	if err != nil {
		return nil, err
	}
	results, err := r.queries.ListConversationsByUser(ctx, pgUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	conversations := make([]*models.Conversation, len(results))
	for i, c := range results {
		conversations[i] = mapToConversation(c)
	}
	return conversations, nil
}

func (r *conversationRepository) UpdateTitle(ctx context.Context, id, title string) (*models.Conversation, error) {
	pgUUID, err := pgutil.ParseUUID(id)
	if err != nil {
		return nil, err
	}
	arg := dbgenerated.UpdateConversationTitleParams{
		ID:    pgUUID,
		Title: pgutil.Text(title),
	}
	result, err := r.queries.UpdateConversationTitle(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to update conversation title: %w", err)
	}
	return mapToConversation(result), nil
}

func (r *conversationRepository) Delete(ctx context.Context, id string) error {
	pgUUID, err := pgutil.ParseUUID(id)
	if err != nil {
		return err
	}
	return r.queries.DeleteConversation(ctx, pgUUID)
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
