package routes

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/db"
)

func GetDB(provider db.DBProvider) (db.DB, error) {
	database, err := provider.GetDatabase()
	if err != nil {
		return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
	}
	return database, nil
}

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func parseMealDate(s string) time.Time {
	if s == "" {
		return time.Now()
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Now()
	}
	return t
}
