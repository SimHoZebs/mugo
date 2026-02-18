package routes

import "time"

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// parseMealDate parses a YYYY-MM-DD date string from the agent response.
// Falls back to the current time if the string is empty or unparseable.
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
