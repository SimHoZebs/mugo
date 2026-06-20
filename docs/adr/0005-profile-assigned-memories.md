# Use profile-assigned memories for food interpretation

Mugo will treat memories as profile-assigned, user-visible food-logging facts and patterns, including dietary preferences. This is more flexible than fixed dietary preference fields and safer than hidden assistant memory because users can inspect and edit the context that influences future assumptions.

## Considered Options

- Store dietary preferences and personalization as fixed profile fields only.
- Let provider-side assistant memory personalize future logs invisibly.
- Store user-visible memories assigned to the profile.

## Consequences

Profile UI should expose memories as editable personalization context. Future estimation should be able to cite memory as a visible assumption source. Existing mobile profile fields for dietary preferences are not the target model and should eventually move under memories.
