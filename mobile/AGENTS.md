# Mugo Mobile Agent Guidelines

The mobile app is an Expo Router React Native app using React 19, NativeWind, Zustand, Jest, Testing Library, and an Orval-generated API client.

## Commands
- Start dev server: `make mobile`
- Start Android emulator: `make emulator`
- Regenerate API client: `make orval`
- Lint: `make lint-mobile`
- Tests: `make test-mobile`
- Direct package commands should use pnpm from `mobile/`, not npm.

## Layout
- Screens/routes: `app/`
- Shared UI primitives: `components/ui/`
- Domain components: `components/`
- Hooks: `hooks/`
- Zustand store: `lib/store.ts`
- Shared types: `lib/types.ts`
- Generated API client: `lib/api/`
- Tests: `test/`

## Generated API Client
- Do not hand-edit `lib/api/`.
- Orval reads the Huma OpenAPI spec from `API_SERVER_URL + "/openapi.json"`.
- After backend route schema changes, run the API, then `make orval`, then `make lint-mobile && make test-mobile`.

## Styling Rules
- Use NativeWind `className` for visual styles. Do not use `StyleSheet.create()` for styling new UI.
- Do not hardcode color hex values inline; use Tailwind token classes.
- Every `bg-*`, `text-*`, and `border-*` class should include a matching `dark:` variant.
- Use loading skeleton views with `stone-200 dark:stone-700`, not spinners.

## Component Rules
- Compose screens from primitives in `components/ui/` instead of re-implementing Text, Button, Card, Badge, Divider, or screen wrappers inline.
- New domain-specific components go in `components/`.
- New primitive base elements go in `components/ui/`.
- Use `<ScreenLayout>` as the root wrapper for screens.
- Use the `<Text>` primitive with the matching `variant` rather than repeating typography classes inline.
- Badges/pills should use the `<Badge>` primitive.

## Layout Rules
- Horizontal screen padding: `px-4`
- Section spacing: `mb-6`
- Card internal padding: `p-4`
- Card radius: `rounded-xl`
- Input/button radius: `rounded-2xl`
- Pill/icon radius: `rounded-full`

## Typography
- Page title: `text-2xl font-bold leading-8`
- Section title: `text-xl font-bold`
- Card title: `text-base font-semibold leading-6`
- Body: `text-base leading-6`
- Label/caption: `text-sm`
- Metadata/pill: `text-xs`

## Color Semantics
- Primary action/success: `emerald-500`
- Warning/assumptions: `amber-*`
- Calories: `amber-500`
- Protein: `emerald-500`
- Carbs: `blue-500`
- Fat: `rose-500`
- Page background: `stone-50 dark:stone-950`
- Card background: `white dark:stone-900`
- Subtle fill: `stone-100` or `stone-200` with `dark:stone-800`

## UX Patterns
- Pressable feedback should use `active:opacity-70` or a specific `active:bg-*` variant.
- Disabled state should use muted fill (`stone-300 dark:stone-800`) and muted text (`stone-500`).
- Inline confirmation uses an emerald checkmark icon, not a text label.
