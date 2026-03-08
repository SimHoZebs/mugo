# TODO-013: Float-to-Int Truncation in nutrition_repository.go

## Status: Cancelled (Not a bug - nutritional values are integers)
## Priority: High (Critical)

### Summary
In `internal/db/repository/nutrition_repository.go`, nutrition values (calories, protein, carbs, fat) are passed as `float64` but truncated to `int64` using `big.NewInt(int64(value))` before being stored. This results in the loss of all decimal data.

### Affected Code
**File:** `internal/db/repository/nutrition_repository.go`
**Lines:** 31-34 (UpsertDaily), 109-116 (UpsertWeekly)

```go
// Line 31
Calories: big.NewInt(int64(totalCalories)),
Protein: big.NewInt(int64(totalProtein)),
Carbs: big.NewInt(int64(totalCarbs)),
Fat: big.NewInt(int64(totalFat)),
```

### Risk
Loss of nutritional precision. While minor for individual meals, this error accumulates in summaries, leading to inaccurate totals for users over time.

### Proposed Fix
The database uses `NUMERIC` types, which sqlc represents as `*big.Int` (if not configured for floats). Either reconfigure `sqlc` to use `float64` for these numeric types, or use a proper helper to convert `float64` to `*big.Int` with scaling (e.g., fixed-point arithmetic if the precision is needed).
