import type {
  Assumption,
  Macros,
  MealLog as ApiMealLog,
} from "@/lib/api/mugoAPI.schemas";

type NutritionPayload = {
  name: string;
  macros: Macros;
  assumptions: Assumption[];
  meal_type: string;
};

export type MealLog = {
  id: string;
  sessionId: string;
  nutrition: NutritionPayload;
};

export function toLocalMealLog(meal: ApiMealLog, sessionId: string): MealLog {
  return {
    id: meal.id,
    sessionId,
    nutrition: {
      name: meal.food_name,
      macros: meal.macros,
      assumptions: meal.assumptions ?? [],
      meal_type: meal.meal_type,
    },
  };
}

export type UnitSystem = "metric" | "imperial";

export type DietaryPreference = {
  id: string;
  text: string;
};

export type UserProfile = {
  username: string;
  name: string;
  dietaryPreferences: DietaryPreference[];
  weight: number;
  height: number;
  unitSystem: UnitSystem;
};
