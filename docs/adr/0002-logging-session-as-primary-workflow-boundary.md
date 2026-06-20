# Use logging session as the primary workflow boundary

Mugo will make `Logging Session` the primary workflow boundary for capture inputs, clarification prompts, correction chat, and produced meal logs. This is more explicit than centering the model on meal logs or provider conversations, but it matches the product flow where a user can provide multimodal input, answer targeted clarifications, produce zero or more committed logs, and later correct those logs from preserved context.

## Considered Options

- Center the model on `Meal Log` and attach interaction state around it.
- Center the model on provider-style `Conversation` records.
- Use `Logging Session` as the product workflow boundary.

## Consequences

Database schema, API routes, mobile navigation, and AI orchestration should be organized around logging sessions first. `Conversation` can remain an implementation detail, but should not be the domain boundary.
