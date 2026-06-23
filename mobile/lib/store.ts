import { create } from "zustand";
import { devtools } from "zustand/middleware";
import { MealLog, UserProfile, DietaryPreference } from "./types";

interface GlobalState {
  meals: MealLog[];
  setMeals: (meals: MealLog[]) => void;
  updateMeal: (id: string, meal: MealLog) => void;
  userProfile: UserProfile | null;
  setUserProfile: (profile: UserProfile | null) => void;
  updateUserProfile: (profile: Partial<UserProfile>) => void;
  addDietaryPreference: (preference: DietaryPreference) => void;
  updateDietaryPreference: (id: string, text: string) => void;
  removeDietaryPreference: (id: string) => void;
  login: (username: string) => void;
  logout: () => void;
}

const useGlobalStore = create<GlobalState>()(
  devtools(
    (set) => ({
      meals: [],
      setMeals: (meals) => set(() => ({ meals })),
      updateMeal: (id, meal) =>
        set((state) => ({
          meals: state.meals.map((m) => (m.id === id ? meal : m)),
        })),
      userProfile: null,
      setUserProfile: (profile) => set(() => ({ userProfile: profile })),
      login: (username: string) =>
        set({
          userProfile: {
            username,
            name: "",
            dietaryPreferences: [],
            weight: 0,
            height: 0,
            unitSystem: "metric",
          },
        }),
      logout: () => set({ userProfile: null }),
      updateUserProfile: (profile) =>
        set((state) => ({
          userProfile: state.userProfile
            ? { ...state.userProfile, ...profile }
            : null,
        })),
      addDietaryPreference: (preference: DietaryPreference) =>
        set((state) => ({
          userProfile: state.userProfile
            ? {
                ...state.userProfile,
                dietaryPreferences: [
                  ...state.userProfile.dietaryPreferences,
                  preference,
                ],
              }
            : null,
        })),
      updateDietaryPreference: (id: string, text: string) =>
        set((state) => ({
          userProfile: state.userProfile
            ? {
                ...state.userProfile,
                dietaryPreferences: state.userProfile.dietaryPreferences.map(
                  (p) => (p.id === id ? { ...p, text } : p),
                ),
              }
            : null,
        })),
      removeDietaryPreference: (id: string) =>
        set((state) => ({
          userProfile: state.userProfile
            ? {
                ...state.userProfile,
                dietaryPreferences: state.userProfile.dietaryPreferences.filter(
                  (p) => p.id !== id,
                ),
              }
            : null,
        })),
    }),
    {
      name: "global-storage",
    },
  ),
);

export default useGlobalStore;
