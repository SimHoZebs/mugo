# LazyFood Go Server - Agent Guidelines

## Build & Development Commands
- **Run API**: `make server` (runs on localhost:8888)
- **Run ADK**: `make adk` (runs ADK with web, api, and webui)
- **Build**: `make build` (builds the API binary)
- **Database**: `make db` (starts database services)
- **Generate SQL**: `make sqlc` (generates Go code from SQL)
- **Dependencies**: `make tidy` (installs and cleans up dependencies)
- **API Docs**: Available at `http://localhost:8888/docs` when the server is running

## Code Style Guidelines
- **Imports**: Group standard library, then third-party, then local imports (server/*)
- **API**: Use Huma v2 with Chi router, follow existing endpoint patterns in routes/
- **Environment**: Load .env with godotenv, never commit secrets

## ADK (Agent Development Kit) Guidelines
- **ADK Documentation**: https://google.github.io/adk-docs/ (primary reference)
- **Go ADK API**: https://pkg.go.dev/google.golang.org/adk
- **Agent Creation**: Use llmagent.New() with gemini models (see agents/weather.go:31)
- **Runner Pattern**: Create agents in agents/, initialize runners in runners/ (see runners/echo.go:10)
- **Session Management**: Use session.InMemoryService() for conversation state
- **Tools**: Leverage geminitool.GoogleSearch{} and other built-in tools
- **Model Config**: Use gemini.NewModel() with GOOGLE_API_KEY from environment

---

## Mobile Frontend Guidelines (React Native / Expo)

### Build & Development Commands
- **Start dev server**: `cd mobile && npx expo start`
- **Run on iOS**: press `i` in the Expo dev server, or `npx expo run:ios`
- **Run on Android**: press `a` in the Expo dev server, or `npx expo run:android`
- **Lint**: `cd mobile && npx eslint .`
- **Tests**: `cd mobile && npx jest`

### Styling Rules
- Use NativeWind (Tailwind) `className` exclusively. Never use `StyleSheet.create()` for visual styles.
- Do **NOT** hardcode color hex values inline (e.g., `color="#10B981"`). Use Tailwind token classes.
- Always include dark mode variants: every `bg-*`, `text-*`, `border-*` must have a `dark:` pair.

### Component Rules
- Compose screens from primitives in `mobile/components/ui/`. Do not re-implement Text, Button, Card inline.
- New domain-specific components go in `mobile/components/`. New primitive base elements go in `mobile/components/ui/`.
- Loading states use skeleton shimmer (`stone-200 dark:stone-700` placeholder views), not spinners.

### Layout Rules
- Horizontal screen padding: `px-4`
- Section spacing: `mb-6`
- Card internal padding: `p-4`
- Use `rounded-xl` for cards, `rounded-2xl` for inputs/buttons, `rounded-full` for pills/icons.
- Use `<ScreenLayout>` from `mobile/components/ui/ScreenLayout.tsx` as the root wrapper for each screen.

### Typography Hierarchy
- Page title: `text-2xl font-bold leading-8`
- Section title: `text-xl font-bold`
- Card title: `text-base font-semibold leading-6`
- Body: `text-base leading-6`
- Label/caption: `text-sm`
- Metadata/pill: `text-xs`

Use the `<Text>` primitive from `mobile/components/ui/Text.tsx` with the matching `variant` prop instead of repeating these classes inline.

### Color Semantics
- Primary action / success: `emerald-500`
- Warning / assumptions: `amber-*`
- Calories: `amber-500` | Protein: `emerald-500` | Carbs: `blue-500` | Fat: `rose-500`
- Page background: `stone-50 dark:stone-950`
- Card background: `white dark:stone-900`
- Subtle fill (inputs, sub-cards): `stone-100/200 dark:stone-800`

### UX Patterns
- Pressable feedback: `active:opacity-70` or specific `active:bg-*` variant
- Disabled state: muted fill (`stone-300 dark:stone-800`) + muted text (`stone-500`)
- Inline confirmation uses a checkmark icon (emerald), not a text label
- Badges/pills use the `<Badge>` primitive from `mobile/components/ui/Badge.tsx`

### Primitive UI Component Library (`mobile/components/ui/`)

| Component | Purpose |
|---|---|
| `Text.tsx` | Typed variants: `h1`, `h2`, `h3`, `body`, `caption`, `micro` — pre-wired with theme tokens |
| `Card.tsx` | Standard rounded + bordered + padded surface |
| `Button.tsx` | Variants: `primary`, `ghost`, `destructive` with disabled states |
| `Badge.tsx` | Inline label (replaces ad-hoc amber/stone badge patterns) |
| `Divider.tsx` | Consistent `h-px` separator |
| `ScreenLayout.tsx` | Page wrapper with standard `bg + px-4` + safe area insets |
| `SectionHeader.tsx` | Labeled section with optional trailing action |

### Screen Layout Reference
```
## Screen: Home Tab
┌─────────────────────────┐
│ [ScreenLayout px-4]     │
│  SectionHeader "Today"  │
│  ─────────────────────  │
│  TotalMacroPanel        │
│  ─────────────────────  │
│  MealCard               │
│  MealCard (loading)     │
│ [InputBar sticky]       │
└─────────────────────────┘
```

---

## Notion Pages Relevant to Mugo (LazyFood / ai-nutrition-tracker)
Below are Notion pages and databases I found that appear directly related to this project. Each entry is the page title followed by the Notion URL.


- Mugo — https://www.notion.so/Mugo-29cee37bd028800ba00cd39ec50dd022
- LazyFood Tasks (database) — https://www.notion.so/6934f50dcef646e6a374a6edc0b44319
- Build single nutrition agent — https://www.notion.so/Build-single-nutrition-agent-fb282cf142ed48d58358871f8c09b129
- LazyFood: Technical Architecture — https://www.notion.so/LazyFood-Technical-Architecture-ce5f20d767774a04bb89a83c4e029bd2
- LazyFood: Market Analysis — https://www.notion.so/LazyFood-Market-Analysis-7843fc72e72441a1afd1b5696490e5d0
- LazyFood: Business Plan — https://www.notion.so/LazyFood-Business-Plan-2b99b481337f4501b3a7df083ad73574
- Add LazyFood project to website — https://www.notion.so/Add-LazyFood-project-to-website-285ee37bd02880c89984eca5cbbde3fd
- LazyFood: October 27, 2025 → November 2, 2025 (weekly note) — https://www.notion.so/LazyFood-October-27-2025-November-2-2025-6d5a994b19694b9f9387b48cea796fed
- Integrate Gemini API for AI cleanup — https://www.notion.so/Integrate-Gemini-API-for-AI-cleanup-5f4892c3f93042fb8ca938c0959fd226
- Integrate Gemini API for AI summaries — https://www.notion.so/Integrate-Gemini-API-for-AI-summaries-90c3efc56a574de1a169557aa68b8829

If you want more pages added (e.g., broader research links, weekly notes, or task pages), tell me the scope and I can append them.
