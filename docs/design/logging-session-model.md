# Logging Session Model

## Purpose

This document describes the target product model for Mugo's AI-assisted food logging flow. It translates the domain language in `CONTEXT.md` into implementation-facing boundaries for backend schema, API routes, mobile flow, and AI orchestration.

## Decided Foundations

### Ownership

An `Account` is the authenticated identity. One account has exactly one `Profile`, and the account owns all data associated with the individual, including logging sessions, meal logs, memories, and media artifacts.

Memories are assigned to profiles. Dietary preferences are memories, not a separate profile concept.

### Primary Workflow Object

`Logging Session` is the primary workflow object. It owns the capture flow, clarification prompts, correction context, and produced meal logs.

The target hierarchy is:

```text
Account
Profile
Memory
Logging Session
Capture Turn
Input Artifact
Clarification Prompt
Correction Chat Message
Meal Log
Food Item
Assumption
```

`Logging Session` belongs to exactly one account. Meal logs produced by a logging session belong to the same account.

### Creation Flow

A new logging session starts with structured multimodal capture, not a free-form chat UI.

Supported input modes are:

- Text
- Voice transcript
- Images
- Image-specific notes

Raw audio does not need to be retained after transcription for the MVP. Images are durable artifacts for the logging session.

The assistant should prefer intelligent guessing with visible assumptions over excessive clarification. Clarifying questions during creation should be targeted multiple-choice prompts with reasonable guesses and optional text input. Broad clarification questions are not allowed.

If the user skips or ignores a clarification, the assistant should proceed with the best reasonable guess and record it as an assumption unless the missing answer makes the target food or requested action impossible to infer.

Before required clarification is resolved or skipped, the logging session may remain unresolved without creating a meal log. Meal logs represent committed estimates.

### Correction Flow

After a meal log exists, corrections use a full chat interface. The correction chat is the user-visible context for modifying logs and understanding why logs changed.

Corrections modify meal logs in place. Mugo does not need meal log revision history for the MVP because change history can be inferred from the logging session's conversation.

### Nutrition Model

A `Meal Log` is a committed nutrition record for one eating occasion. A meal log contains one or more `Food Items`.

A `Food Item` is a nutritionally meaningful food or drink that can receive its own nutrition estimate. Food item boundaries are inferred from user intent and evidence across inputs, not from input count.

Meal log totals are derived from food item estimates. Corrections update affected food items, then recalculate meal log totals.

Food item estimates represent consumed amount, not necessarily visible or served amount. Shared food and partial consumption are context for consumed amount.

### Meal Type

Meal type is required on a meal log and falls back to `unknown` when it cannot be confidently assigned.

For a user on a given day, there is at most one meal log for each of:

- Breakfast
- Brunch
- Lunch
- Dinner

Brunch can coexist with breakfast and lunch. Snack can occur multiple times. Unknown remains ambiguous and should not force unrelated items to merge.

The meal log day is the user-intended local calendar day for the eating occasion. If intent is absent, use the user's current local date.

### Assumptions

An `Assumption` is a visible uncertain input used in a nutrition estimate. Assumptions primarily attach to food items, but may attach to meal logs when the uncertainty applies to the whole eating occasion.

Every assumption has a user-visible source or reason. Initial assumption source categories are:

- User input
- Memory
- Typical value
- Previous correction
- Food database

Source categories can expand over time.

### Memories

A `Memory` is durable user-specific food-logging knowledge assigned to a profile. Memories make future assumptions more accurate and can make something that would otherwise be a guess into known user context.

Memories are created from explicit user instruction or repeated confirmed behavior, not from a single unconfirmed assistant guess. Memories must be user-visible and editable.

## Source Of Truth

Mugo's own persisted logging session context is the canonical source of truth. Provider-side assistant session memory is execution infrastructure, not product state.

AI output becomes product data only after it is normalized into structured records such as logging sessions, meal logs, food items, assumptions, clarifications, and memories. Raw model output may be retained for debugging or audit, but user-facing behavior should not depend on raw model output.

## Current Code Mismatches

The current implementation does not yet match this target model:

- `meal_logs` currently represents one food item with one `food_name` and one macros object; target model requires one meal log with many food items.
- `conversations` currently acts as the closest session concept, but target language is `Logging Session` and it owns more than chat history.
- The mobile home screen currently assumes `createMealLog` returns a single `meal`, while the server response currently returns `meals`; both need to align with the target `Meal Log` shape.
- Profile data currently stores dietary preferences separately; target model treats dietary preferences as memories assigned to a profile.
- Authentication and authorization are pending, so current `users` code is not yet the target `Account` model.
- Images are not yet durable logging session artifacts.
- Creation-flow clarification prompts are not yet modeled.

## MVP Implementation Slice

`ready-for-agent`: Introduce the backend foundation for Logging Sessions and nested Meal Logs.

Scope:

- Add database schema for logging sessions as the primary workflow object.
- Add food items under meal logs so one meal log can contain multiple food items.
- Add assumptions with visible source/reason and attachment to either food item or meal log.
- Preserve existing AI execution path where possible, but normalize AI output into the new structured records.
- Keep raw assistant output as debug/audit data only.
- Keep correction behavior as in-place updates.

Out of scope for this slice:

- Full auth/account implementation.
- Durable image upload/storage.
- Memory promotion heuristics.
- Auto-discard timing for unresolved sessions.
- Complete mobile redesign.
- Full correction chat UI.

## Remaining High-Impact Decisions

- Exact database table shape and migration strategy from existing `meal_logs`.
- Whether to rename existing `users` to `accounts` now or defer until authentication is implemented.
- API route shape for logging session creation, clarification answers, and correction chat.
- Storage provider and lifecycle for durable images.
- How memories are represented in the database and surfaced in profile UI.
