import { NutritionPayload } from "./api/conversationAPI.schemas";

export type Meal = {
  id: string;
  sessionId: string;
  nutrition: NutritionPayload;
};

export type UnitSystem = "metric" | "imperial";

export type DietaryPreference = {
  id: string;
  text: string;
};

export type UserProfile = {
  name: string;
  dietaryPreferences: DietaryPreference[];
  weight: number;
  height: number;
  unitSystem: UnitSystem;
};
