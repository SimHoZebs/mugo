# Persist logging session context as source of truth

Mugo will persist its own structured logging session context and treat provider-side assistant session memory as execution infrastructure, not canonical product memory. This is harder than relying on ADK/provider session state, but it lets Mugo reconstruct corrections, explain assumptions, retain image and clarification context, and avoid losing product behavior when provider memory is unavailable or changes.

## Considered Options

- Depend on provider/session memory for correction context.
- Persist only final meal logs and raw model responses.
- Persist structured logging session context as Mugo's source of truth.

## Consequences

The backend must model logging sessions, capture artifacts, clarification prompts, correction messages, meal logs, food items, assumptions, and memory references explicitly. Raw model output can still be kept for debugging or audit, but user-facing behavior should read structured product records.
