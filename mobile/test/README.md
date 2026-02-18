# Testing Framework

This project uses **Jest** with `jest-expo` preset and `@testing-library/react-native` for testing.

## Quick Start

```bash
# Run tests
pnpm test

# Run tests in watch mode
pnpm test:watch

# Run tests with coverage
pnpm test:coverage
```

## Test Structure

```
test/
├── unit/                 # Unit tests (stores, hooks, utilities)
│   └── store.test.ts     # Zustand store tests
├── components/           # Component tests
│   └── MealCard.test.tsx # React Native component tests
└── integration/          # Integration tests (future)
```

## Writing Tests

### Unit Tests (Stores, Hooks)

```typescript
import useGlobalStore from "@/lib/store";

describe("useGlobalStore", () => {
  beforeEach(() => {
    useGlobalStore.setState({ meals: [], userProfile: null });
  });

  it("should login and set userProfile", () => {
    const { login } = useGlobalStore.getState();
    login("testuser");
    expect(useGlobalStore.getState().userProfile).not.toBeNull();
  });
});
```

### Component Tests

```typescript
import { render, screen, fireEvent } from "@testing-library/react-native";
import MyComponent from "@/components/MyComponent";

describe("MyComponent", () => {
  it("should render correctly", () => {
    render(<MyComponent />);
    expect(screen.getByText("Hello")).toBeOnTheScreen();
  });

  it("should handle press", () => {
    const onPress = jest.fn();
    render(<MyComponent onPress={onPress} />);
    fireEvent.press(screen.getByText("Click me"));
    expect(onPress).toHaveBeenCalled();
  });
});
```

## Matchers

`@testing-library/react-native` provides these key matchers:

- `toBeOnTheScreen()` - Element is rendered
- `toBeDisabled()` - Element is disabled
- `toHaveTextContent(text)` - Element contains text

## Mocks

The `jest-expo` preset automatically mocks:
- React Native core components
- Expo modules
- Navigation (expo-router)

For custom mocks, add them to `jest.config.js` or create `__mocks__/` directories.

## Coverage

Run `pnpm test:coverage` to generate a coverage report. Coverage includes:
- `lib/**/*` - Utilities and store
- `components/**/*` - UI components
- `hooks/**/*` - Custom hooks

Excludes: `lib/api/**/*` (auto-generated)

## Notes

- Uses **Jest 29** for compatibility with `jest-expo`
- Test environment: Node (configured by jest-expo)
- TypeScript support via `babel-preset-expo`

## E2E Testing

For end-to-end testing on actual devices, consider:
- **Maestro** - Recommended by Expo
- **Detox** - Gray-box E2E testing
