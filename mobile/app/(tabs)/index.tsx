import {
  KeyboardAwareScrollView,
  KeyboardStickyView,
} from "react-native-keyboard-controller";
import "react-native-get-random-values";
import { v7 as uuid7 } from "uuid";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import TotalMacroPanel from "@/components/TotalMacroPanel";
import MealCard from "@/components/MealCard";
import InputBar from "@/components/InputBar";
import useGlobalStore from "@/lib/store";

// Temporary: Direct fetch until we regenerate API client with orval
const API_URL = process.env.EXPO_PUBLIC_API_URL || "http://localhost:8888";

export default function HomeScreen() {
  const meals = useGlobalStore((state) => state.meals);
  const setMeals = useGlobalStore((state) => state.setMeals);

  const handleSubmitNutrition = async (text: string) => {
    const newSessionId = uuid7();

    const payload = {
      description: text,
      session_id: newSessionId,
      user_id: "user-1",
    };

    try {
      const response = await fetch(`${API_URL}/meals`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        throw new Error(`Server error: ${response.status}`);
      }

      const data = await response.json();
      const newMeal = {
        id: data.meal.id,
        sessionId: newSessionId,
        nutrition: {
          name: data.meal.food_name,
          macros: data.meal.macros,
          assumptions: data.meal.assumptions,
          meal_type: data.meal.meal_type,
        },
      };
      setMeals([...meals, newMeal]);
    } catch (error) {
      console.error("Error submitting meal to:", `${API_URL}/meals`);
      console.error("Payload:", payload);
      console.error("Error:", error);
    }
  };

  return (
    <ThemedView className="h-full w-full gap-4 pt-8">
      <ThemedText className="px-4" type="title">
        Tuesday
      </ThemedText>

      <TotalMacroPanel meals={meals} />

      <KeyboardAwareScrollView>
        <ThemedView className="gap-4 px-4">
          {meals.map((meal) => (
            <MealCard key={meal.id} meal={meal} />
          ))}
        </ThemedView>
      </KeyboardAwareScrollView>

      <KeyboardStickyView offset={{ closed: 0, opened: 80 }}>
        <InputBar onSubmit={handleSubmitNutrition}>
          <InputBar.Action onPress={() => console.log("Mic pressed")}>
            <ThemedText className="text-lg">🎤</ThemedText>
          </InputBar.Action>
          <InputBar.Input placeholder="Describe your meals..." />
          <InputBar.Action onPress={() => console.log("Camera pressed")}>
            <ThemedText className="text-lg">📷</ThemedText>
          </InputBar.Action>
        </InputBar>
      </KeyboardStickyView>
    </ThemedView>
  );
}
