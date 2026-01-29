import {
  KeyboardAwareScrollView,
  KeyboardStickyView,
} from "react-native-keyboard-controller";
import "react-native-get-random-values";
import { v7 as uuid7 } from "uuid";

import { useRef, useState } from "react";
import { ScrollView } from "react-native";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import TotalMacroPanel from "@/components/TotalMacroPanel";
import MealCard from "@/components/MealCard";
import InputBar from "@/components/InputBar";
import useGlobalStore from "@/lib/store";
import { getPostMealsUrl, postMeals } from "@/lib/api/default/default";

export default function HomeScreen() {
  const meals = useGlobalStore((state) => state.meals);
  const setMeals = useGlobalStore((state) => state.setMeals);
  const [pendingMealId, setPendingMealId] = useState<string | null>(null);
  const scrollViewRef = useRef<ScrollView | null>(null);

  const handleSubmitNutrition = async (text: string) => {
    const newSessionId = uuid7();
    const loadingMealId = `loading-${newSessionId}`;

    const payload = {
      description: text,
      session_id: newSessionId,
      user_id: "user-1",
    };

    try {
      setPendingMealId(loadingMealId);
      requestAnimationFrame(() =>
        scrollViewRef.current?.scrollToEnd({ animated: true }),
      );
      const response = await postMeals(payload);

      if (response.status !== 200) {
        throw response.data;
      }

      const data = response.data;
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
      setPendingMealId(null);
    } catch (error) {
      console.error("Error submitting request to:", getPostMealsUrl());
      console.error("Error:", error);
      setPendingMealId(null);
    }
  };

  return (
    <ThemedView className="h-full w-full gap-4 pt-8">
      <ThemedText className="px-4" type="title">
        Tuesday
      </ThemedText>

      <TotalMacroPanel meals={meals} />

      <KeyboardAwareScrollView ref={scrollViewRef}>
        <ThemedView className="gap-4 px-4">
          {meals.map((meal) => (
            <MealCard key={meal.id} meal={meal} />
          ))}
          {pendingMealId && <MealCard key={pendingMealId} loading />}
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
