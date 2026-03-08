package db

import "strings"

// ToMigrateURL converts a standard postgres:// or postgresql:// URL to a pgx5:// URL
// required by the golang-migrate pgx/v5 driver.
func ToMigrateURL(databaseURL string) string {
	if strings.HasPrefix(databaseURL, "postgres://") {
		return strings.Replace(databaseURL, "postgres://", "pgx5://", 1)
	} else if strings.HasPrefix(databaseURL, "postgresql://") {
		return strings.Replace(databaseURL, "postgresql://", "pgx5://", 1)
	} else if !strings.HasPrefix(databaseURL, "pgx5://") {
		return "pgx5://" + databaseURL
	}
	return databaseURL
}
