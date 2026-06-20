# Mugo

Mugo is an AI-assisted food logging context for turning multimodal user input into nutrition records and summaries.

## Language

**Mugo**:
An AI-assisted food logging product that helps users record food through natural language, voice, images, or combinations of those inputs.
_Avoid_: LazyFood

**Account**:
The authenticated identity that owns a person's Mugo data. An account owns exactly one profile, its meal logs, and all other data associated with that individual.
_Avoid_: User

**Multimodal Input**:
A user-provided food description made from one or more input modes, including text, voice, images, and optional notes.
_Avoid_: Text input only, meal description

**Logging Session**:
An interactive context for recording food, usually for one user-determined eating occasion. It can contain multimodal inputs, AI clarification turns, and modification requests, and may produce multiple meal logs when the user explicitly describes multiple eating occasions.
_Avoid_: Meal Session, Conversation, Session

Each logging session belongs to exactly one account. Meal logs produced by a logging session belong to the same account.

A logging session is the primary context for capture inputs, clarification prompts, correction chat, and produced meal logs.

Mugo's own persisted logging session context is the source of truth for later reconstruction. Provider-side assistant session memory is execution infrastructure, not canonical product memory.

Assistant output becomes product data only after it is normalized into structured logging session, meal log, food item, assumption, clarification, or memory records. Raw assistant output is not product state.

A logging session can be unresolved when it has not produced a meal log. Unresolved logging sessions may be resumed or discarded, and should eventually be discarded automatically.

A logging session can be active, logged, unresolved, or discarded. Logged sessions can still accept corrections because logging preserves the relevant context for later modification.

The assistant should prefer intelligent guessing with visible assumptions over excessive clarification. It should ask a clarifying question only when missing information would make the estimate misleading, the target item or meal log ambiguous, or the user's requested action unsafe to infer.

A logging session starts with a structured capture flow for multimodal input. Clarifying questions during creation should be targeted multiple-choice prompts with reasonable guesses and optional text input.

During creation, the assistant should ask only about glaring ambiguity. Broad clarification questions are not allowed.

If the user skips or ignores a creation-flow clarification, proceed with the best reasonable guess and record it as an assumption unless the missing answer makes the target food or requested action impossible to infer.

After a meal log exists, corrections use a full chat interface. The correction chat is the user-visible context for modifying logs and understanding how meal logs changed.

A capture turn may contain text, voice transcript, images, and notes. Related input artifacts captured together can be grouped into one user turn.

Voice input contributes a transcript to the chat turn. Raw audio does not need to be retained after transcription for the MVP.

Images are durable chat artifacts for a logging session. Assistant-derived image observations can support estimates, but do not replace the original image.

An image note is user-provided text attached to a specific image and used to interpret that image.

**Food Item**:
A nutritionally meaningful food or drink identified within a logging session that can receive its own nutrition estimate. Incidental details may remain descriptive context instead of separate food items.
_Avoid_: Meal, dish

Food item boundaries are inferred from user intent and evidence across inputs, not from input count.

A food item is presented with one combined nutrition estimate. Components can support estimation and correction but are not required to be user-facing unless they explain an assumption.

Shared food provides context for the consumed portion of a food item. A meal log records what the account consumed, not everything visible or ordered.

Food item nutrition estimates represent consumed amount. When only served or visible amount is known, the assistant may assume consumed amount and expose that assumption.

Post-log corrections can change the consumed amount for a food item without creating a new food item or meal log.

**Assumption**:
A visible uncertain input used in a nutrition estimate, such as portion size, ingredient choice, preparation method, or interpretation of an image.
_Avoid_: Hidden guess

Assumptions primarily attach to food items. Assumptions can attach to a meal log when the uncertainty applies to the whole eating occasion.

Every assumption has a user-visible source or reason. A memory can be the source of an assumption.

Assumption sources can expand over time. Initial sources include user input, memory, typical values, previous corrections, and food databases.

**Memory**:
A durable user-specific food-logging fact or pattern that helps future assumptions become more accurate. A memory can make something that would otherwise be an assumption into a known fact for that user.
_Avoid_: Temporary guess, session context

A memory is created from explicit user instruction or repeated confirmed behavior, not from a single unconfirmed assistant guess. One-off corrections should not automatically become memory unless the user frames them as durable.

Memories are part of the user's profile context and must be inspectable and editable by the user.

Memories are assigned to profiles. Dietary preferences are memories, not a separate kind of profile data.

**Profile**:
A user-owned personal context for Mugo. Memories are assigned to profiles and personalize food interpretation for that profile.
_Avoid_: Account, User

One account has exactly one profile.

**Meal Log**:
A saved nutrition record for one eating occasion. A meal log contains one or more food items.
_Avoid_: Food item, individual estimate

A committed meal log has at least one consumed food item. If correction removes all consumed food items, the meal log should be deleted or discarded.

Meal log nutrition totals are the sum of its food items. Corrections update affected food items and then recalculate the meal log total.

Meal logs represent committed estimates. Before required clarification is resolved or skipped, the logging session may remain unresolved without creating a meal log.

Corrections modify the meal log in place. Change history is inferred from the logging session's conversation, not from meal log revisions.

The meal log day is the user-intended local calendar day for the eating occasion. If user intent is absent, use the user's current local date.

**Meal Type**:
A required classification for a meal log, such as breakfast, brunch, lunch, dinner, snack, or unknown. Explicit user-provided meal type takes precedence and can clarify that ambiguous food entries belong to the same meal log.
_Avoid_: Optional tag

Use unknown when nutrition can be estimated but meal type cannot be confidently assigned. Ask for clarification only when the missing meal type affects merging or the user's goal.

For a user on a given day, there is at most one meal log for breakfast, brunch, lunch, or dinner. Brunch can coexist with breakfast and lunch. Snack can occur multiple times, and unknown remains ambiguous.
