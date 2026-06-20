# Use structured capture for creation and chat for correction

Mugo will use a structured multimodal capture flow for creating meal logs, with only targeted multiple-choice clarifications and optional free text, then use a full chat interface after logs exist for corrections. This avoids turning initial logging into a slow chat interview while still preserving conversational context where it is most useful: modifying and explaining existing meal logs.

## Considered Options

- Use full chat for both creation and correction.
- Use structured forms for both creation and correction.
- Use structured capture for creation and chat for correction.

## Consequences

Creation APIs and UI should model capture artifacts and clarification prompts, not arbitrary chat messages. Correction APIs and UI should support chat-style interaction against an existing logging session and its produced meal logs.
