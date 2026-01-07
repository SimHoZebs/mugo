package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	dbgenerated "github.com/simhozebs/mugo/internal/db/dbgenerated"
	"github.com/simhozebs/mugo/internal/models"
)

type ConversationRepository struct {
	queries *dbgenerated.Queries
}

func NewConversationRepository(queries *dbgenerated.Queries) *ConversationRepository {
	return &ConversationRepository{queries: queries}
}

func (r *ConversationRepository) Create(ctx context.Context, userID, sessionID, title string) (*models.Conversation, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	titleText := pgtype.Text{String: title, Valid: title != ""}
	arg := dbgenerated.CreateConversationParams{
		UserID:    pgUUID,
		SessionID: sessionID,
		Title:     titleText,
	}
	result, err := r.queries.CreateConversation(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}
	return mapToConversation(result), nil
}

func (r *ConversationRepository) GetByID(ctx context.Context, id string) (*models.Conversation, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	result, err := r.queries.GetConversation(ctx, pgUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	return mapToConversation(result), nil
}

func (r *ConversationRepository) GetBySessionID(ctx context.Context, userID, sessionID string) (*models.Conversation, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
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

func (r *ConversationRepository) ListByUser(ctx context.Context, userID string) ([]*models.Conversation, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
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

func (r *ConversationRepository) UpdateTitle(ctx context.Context, id, title string) (*models.Conversation, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	titleText := pgtype.Text{String: title, Valid: true}
	arg := dbgenerated.UpdateConversationTitleParams{
		ID:    pgUUID,
		Title: titleText,
	}
	result, err := r.queries.UpdateConversationTitle(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to update conversation title: %w", err)
	}
	return mapToConversation(result), nil
}

func (r *ConversationRepository) Delete(ctx context.Context, id string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	return r.queries.DeleteConversation(ctx, pgUUID)
}

func mapToConversation(c dbgenerated.Conversation) *models.Conversation {
	return &models.Conversation{
		ID:        c.ID.String(),
		UserID:    c.UserID.String(),
		SessionID: c.SessionID,
		Title:     mapTextToPtr(c.Title),
		CreatedAt: c.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Time.Format(time.RFC3339),
	}
}

func mapTextToPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}
