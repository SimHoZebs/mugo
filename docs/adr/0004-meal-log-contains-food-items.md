# Model meal logs as containing food items

Mugo will model a `Meal Log` as the committed record for one eating occasion, containing one or more `Food Items`, instead of treating each food item as its own meal log. This better matches user intent for multimodal meal capture: a lunch can contain a burger, fries, and a drink while still being one meal log with totals derived from item estimates.

## Considered Options

- Keep one meal log per estimated food item.
- Store only one aggregate estimate per meal log.
- Model a meal log containing multiple food items with derived totals.

## Consequences

The current `meal_logs` shape is not the target model because it has one `food_name` and one macros object. The backend should introduce food items under meal logs, attach most assumptions to food items, and calculate meal log totals from food item estimates.
