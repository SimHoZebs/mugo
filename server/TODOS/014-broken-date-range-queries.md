# TODO-014: Broken Date Range Queries (endDate never used)

## Status: Completed
## Priority: High (Critical)

### Summary
In both `nutrition_repository.go` and `meal_repository.go`, repository functions that accept a date range (`startDate`, `endDate`) only pass `startDate` to the underlying SQL query. The `endDate` parameter is accepted in the Go signature but never used, meaning the queries return all logs *from* the start date, regardless of the end date.

### Affected Code
**File:** `internal/db/repository/nutrition_repository.go`
**Lines:** 81, 142

```go
// Line 81 (ListDailyByDateRange)
func (r *nutritionRepository) ListDailyByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*models.DailyNutritionSummary, error) {
    // ...
    // Line 88
    Summaries, err := r.db.ListDailyNutritionSummariesByUserAndDateRange(ctx, db.ListDailyNutritionSummariesByUserAndDateRangeParams{
        UserID: userID,
        StartDate: pgtype.Timestamp{Time: startDate, Valid: true},
    })
    // ...
}
```

**File:** `internal/db/repository/meal_repository.go`
**Line:** 121

```go
// Line 121 (ListByUserAndDateRange)
func (r *mealLogRepository) ListByUserAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*models.MealLog, error) {
    // ...
    // Line 128
    MealLogs, err := r.db.ListMealLogsByUserAndDateRange(ctx, db.ListMealLogsByUserAndDateRangeParams{
        UserID: userID,
        StartDate: pgtype.Timestamp{Time: startDate, Valid: true},
    })
    // ...
}
```

### Risk
Users who filter their history for a specific week or month will receive **all data** from the start of that range until the present day, making range filtering unusable.

### Proposed Fix
Update the repository calls to pass the `endDate` parameter to the `sqlc` generated params struct.

```go
StartDate: pgtype.Timestamp{Time: startDate, Valid: true},
EndDate: pgtype.Timestamp{Time: endDate, Valid: true},
```
