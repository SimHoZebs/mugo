package pgutil

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type txKey struct{}

func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

func ParseUUID(s string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid UUID: %w", err)
	}
	return pgtype.UUID{
		Bytes: [16]byte(parsed),
		Valid: true,
	}, nil
}

func ParseUUIDPtr(s string) (pgtype.UUID, error) {
	if s == "" {
		return pgtype.UUID{}, nil
	}
	return ParseUUID(s)
}

func Timestamp(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func Date(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func Text(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func TextPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func GenerateUUID() string {
	return uuid.New().String()
}

func FromText(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}
