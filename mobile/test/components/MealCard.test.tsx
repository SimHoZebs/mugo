import { render, screen, fireEvent } from "@testing-library/react-native";
import MealCard from "@/components/MealCard";

const mockPush = jest.fn();
jest.mock("expo-router", () => ({
  useRouter: () => ({
    push: mockPush,
    replace: jest.fn(),
    back: jest.fn(),
  }),
  useLocalSearchParams: () => ({}),
  useGlobalSearchParams: () => ({}),
  usePathname: () => "/",
  Link: ({ children }) => children,
}));

const mockMeal = {
  id: "meal-1",
  createdAt: "2024-01-01T00:00:00Z",
  nutrition: {
    name: "Grilled Chicken Salad",
    macros: {
      calories: 450,
      protein: 35,
      carbs: 20,
      fat: 25,
    },
    assumptions: [],
  },
};

const mockMealWithAssumptions = {
  ...mockMeal,
  nutrition: {
    ...mockMeal.nutrition,
    assumptions: [{ text: "Assumed olive oil" }, { text: "Assumed dressing" }],
  },
};

describe("MealCard", () => {
  beforeEach(() => {
    mockPush.mockClear();
  });

  it("should render meal name when meal is provided", () => {
    render(<MealCard meal={mockMeal as any} />);
    expect(screen.getByText("Grilled Chicken Salad")).toBeOnTheScreen();
  });

  it("should render macros correctly", () => {
    render(<MealCard meal={mockMeal as any} />);

    expect(screen.getByText("450 kcal")).toBeOnTheScreen();
    expect(screen.getByText("35g P")).toBeOnTheScreen();
    expect(screen.getByText("20g C")).toBeOnTheScreen();
    expect(screen.getByText("25g F")).toBeOnTheScreen();
  });

  it("should show assumption count when meal has assumptions", () => {
    render(<MealCard meal={mockMealWithAssumptions as any} />);
    expect(screen.getByText("2 assumptions")).toBeOnTheScreen();
  });

  it("should show singular 'assumption' for single assumption", () => {
    const mealWithOneAssumption = {
      ...mockMeal,
      nutrition: {
        ...mockMeal.nutrition,
        assumptions: [{ text: "Assumed olive oil" }],
      },
    };

    render(<MealCard meal={mealWithOneAssumption as any} />);
    expect(screen.getByText("1 assumption")).toBeOnTheScreen();
  });

  it("should navigate to meal detail on press", () => {
    render(<MealCard meal={mockMeal as any} />);

    fireEvent.press(screen.getByText("Grilled Chicken Salad"));

    expect(mockPush).toHaveBeenCalledWith({
      pathname: "/meal-detail",
      params: { id: "meal-1" },
    });
  });

  it("should not be pressable when loading", () => {
    render(<MealCard loading />);

    expect(mockPush).not.toHaveBeenCalled();
  });

  it("should render loading skeleton when loading is true", () => {
    const { queryByText } = render(<MealCard loading />);
    expect(queryByText("Grilled Chicken Salad")).toBeNull();
  });

  it("should render nothing when not loading and no meal", () => {
    const { queryByText } = render(<MealCard />);
    expect(queryByText("Grilled Chicken Salad")).toBeNull();
  });
});
