import { create } from "zustand";
import { devtools, persist } from "zustand/middleware";
import type {} from "@redux-devtools/extension"; // required for devtools typing
import { Meal, UserProfile, DietaryPreference } from "./types";

interface GlobalState {
  meals: Meal[];
  setMeals: (meals: Meal[]) => void;
  updateMeal: (id: string, meal: Meal) => void;
  userProfile: UserProfile;
  setUserProfile: (profile: UserProfile) => void;
  updateUserProfile: (profile: Partial<UserProfile>) => void;
  addDietaryPreference: (preference: DietaryPreference) => void;
  updateDietaryPreference: (id: string, text: string) => void;
  removeDietaryPreference: (id: string) => void;
}

const useGlobalStore = create<GlobalState>()(
  devtools(
    persist(
      (set) => ({
        meals: [
          {
            id: "1",
            sessionId: "meal-1",
            nutrition: {
              name: "Oatmeal with Berries",
              assumptions: [
                {
                  id: "1",
                  field: "portion_size",
                  assumed_value: 1,
                  unit: "cup",
                  confidence: "high",
                  rationale: "Standard serving size for oatmeal",
                },
                {
                  id: "8",
                  field: "milk_type",
                  assumed_value: 1,
                  unit: "cup",
                  confidence: "medium",
                  rationale: "Assumed whole milk for cooking",
                },
              ],
              macros: {
                calories: 320,
                protein: 12,
                carbs: 45,
                fat: 8,
              },
            },
          },
          {
            id: "2",
            sessionId: "meal-2",
            nutrition: {
              name: "Grilled Chicken Salad",
              assumptions: [
                {
                  id: "2",
                  field: "portion_size",
                  assumed_value: 150,
                  unit: "grams",
                  confidence: "medium",
                  rationale: "Typical chicken breast portion",
                },
                {
                  id: "9",
                  field: "cooking_method",
                  assumed_value: 1,
                  unit: "grilled",
                  confidence: "high",
                  rationale: "Common preparation method",
                },
                {
                  id: "10",
                  field: "dressing_amount",
                  assumed_value: 2,
                  unit: "tbsp",
                  confidence: "low",
                  rationale: "Light dressing application",
                },
              ],
              macros: {
                calories: 450,
                protein: 35,
                carbs: 20,
                fat: 18,
              },
            },
          },
          {
            id: "3",
            sessionId: "meal-3",
            nutrition: {
              name: "Greek Yogurt Parfait",
              assumptions: [
                {
                  id: "3",
                  field: "portion_size",
                  assumed_value: 200,
                  unit: "grams",
                  confidence: "high",
                  rationale: "Standard yogurt container size",
                },
                {
                  id: "11",
                  field: "toppings",
                  assumed_value: 50,
                  unit: "grams granola",
                  confidence: "medium",
                  rationale: "Typical granola addition",
                },
              ],
              macros: {
                calories: 280,
                protein: 20,
                carbs: 30,
                fat: 10,
              },
            },
          },
          {
            id: "4",
            sessionId: "meal-4",
            nutrition: {
              name: "Turkey Sandwich",
              assumptions: [
                {
                  id: "4",
                  field: "portion_size",
                  assumed_value: 100,
                  unit: "grams",
                  confidence: "medium",
                  rationale: "Typical deli meat portion",
                },
                {
                  id: "12",
                  field: "bread_type",
                  assumed_value: 2,
                  unit: "slices whole wheat",
                  confidence: "medium",
                  rationale: "Common bread choice",
                },
                {
                  id: "13",
                  field: "condiments",
                  assumed_value: 1,
                  unit: "tbsp mayo",
                  confidence: "low",
                  rationale: "Light spread",
                },
              ],
              macros: {
                calories: 350,
                protein: 25,
                carbs: 35,
                fat: 12,
              },
            },
          },
          {
            id: "5",
            sessionId: "meal-5",
            nutrition: {
              name: "Vegetable Stir Fry",
              assumptions: [
                {
                  id: "5",
                  field: "portion_size",
                  assumed_value: 300,
                  unit: "grams",
                  confidence: "medium",
                  rationale: "Standard vegetable serving",
                },
                {
                  id: "14",
                  field: "cooking_oil",
                  assumed_value: 1,
                  unit: "tbsp",
                  confidence: "medium",
                  rationale: "Minimal oil for stir frying",
                },
              ],
              macros: {
                calories: 200,
                protein: 8,
                carbs: 25,
                fat: 10,
              },
            },
          },
          {
            id: "6",
            sessionId: "meal-6",
            nutrition: {
              name: "Salmon with Quinoa",
              assumptions: [
                {
                  id: "6",
                  field: "portion_size",
                  assumed_value: 150,
                  unit: "grams",
                  confidence: "medium",
                  rationale: "Typical fish fillet size",
                },
                {
                  id: "15",
                  field: "cooking_method",
                  assumed_value: 1,
                  unit: "baked",
                  confidence: "high",
                  rationale: "Healthy cooking method",
                },
                {
                  id: "16",
                  field: "seasoning",
                  assumed_value: 1,
                  unit: "tsp herbs",
                  confidence: "medium",
                  rationale: "Light seasoning",
                },
              ],
              macros: {
                calories: 400,
                protein: 30,
                carbs: 25,
                fat: 20,
              },
            },
          },
          {
            id: "7",
            sessionId: "meal-7",
            nutrition: {
              name: "Banana Smoothie",
              assumptions: [
                {
                  id: "7",
                  field: "portion_size",
                  assumed_value: 1,
                  unit: "medium banana",
                  confidence: "high",
                  rationale: "Average banana size",
                },
                {
                  id: "17",
                  field: "liquid_base",
                  assumed_value: 1,
                  unit: "cup milk",
                  confidence: "medium",
                  rationale: "Common smoothie base",
                },
              ],
              macros: {
                calories: 250,
                protein: 5,
                carbs: 40,
                fat: 3,
              },
            },
          },
        ],
        setMeals: (meals) => set(() => ({ meals })),
        updateMeal: (id, meal) =>
          set((state) => ({
            meals: state.meals.map((m) => (m.id === id ? meal : m)),
          })),
        userProfile: {
          name: "",
          dietaryPreferences: [],
          weight: 0,
          height: 0,
          unitSystem: "metric",
        },
        setUserProfile: (profile) => set(() => ({ userProfile: profile })),
        updateUserProfile: (profile) =>
          set((state) => ({
            userProfile: { ...state.userProfile, ...profile },
          })),
        addDietaryPreference: (preference: DietaryPreference) =>
          set((state) => ({
            userProfile: {
              ...state.userProfile,
              dietaryPreferences: [
                ...state.userProfile.dietaryPreferences,
                preference,
              ],
            },
          })),
        updateDietaryPreference: (id: string, text: string) =>
          set((state) => ({
            userProfile: {
              ...state.userProfile,
              dietaryPreferences: state.userProfile.dietaryPreferences.map((p) =>
                p.id === id ? { ...p, text } : p,
              ),
            },
          })),
        removeDietaryPreference: (id: string) =>
          set((state) => ({
            userProfile: {
              ...state.userProfile,
              dietaryPreferences: state.userProfile.dietaryPreferences.filter(
                (p) => p.id !== id,
              ),
            },
          })),
      }),
      {
        name: "global-storage",
      },
    ),
  ),
);

export default useGlobalStore;
