import { useState } from "react";
import { View } from "react-native";
import { useLocalSearchParams } from "expo-router";
import {
  KeyboardAwareScrollView,
  KeyboardStickyView,
} from "react-native-keyboard-controller";
import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import InputBar from "@/components/InputBar";
import { MacroRow } from "@/components/MacroRow";
import { AssumptionCard } from "@/components/AssumptionCard";
import useGlobalStore from "@/lib/store";

// Temporary: Direct fetch until we regenerate API client with orval
const API_URL = process.env.EXPO_PUBLIC_API_URL || "http://localhost:8888";

export default function MealDetailScreen() {
  const [isSubmitting, setIsSubmitting] = useState(false);

  const { id } = useLocalSearchParams<{ id: string }>();
  const meal = useGlobalStore((state) => state.meals.find((m) => m.id === id));

  const handleSubmitCorrection = async (text: string) => {
    if (isSubmitting || !meal) return;

    setIsSubmitting(true);
    try {
      const response = await fetch(`${API_URL}/meals/${meal.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          correction: text,
        }),
      });

      if (!response.ok) {
        throw new Error(`Server error: ${response.status}`);
      }

      const data = await response.json();
      console.log("Meal updated:", data.meal);
      // TODO: Update the meal in global state
    } catch (error) {
      console.error("Failed to submit correction:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return !meal ? (
    <ThemedView className="flex-1 items-center justify-center">
      <ThemedText>Meal not found</ThemedText>
    </ThemedView>
  ) : (
    <ThemedView className="flex-1 px-3">
      <KeyboardAwareScrollView className="flex-1 px-4">
        <View className="pt-4 pb-6">
          <ThemedText type="title">
            {meal.nutrition.name || "Meal Details"}
          </ThemedText>
        </View>

        {/* Macros Section */}
        <View className="mb-6">
          <ThemedText type="subtitle" className="mb-3">
            Nutrition Facts
          </ThemedText>

          <View className="bg-white dark:bg-stone-900 rounded-xl p-4 border border-stone-200 dark:border-stone-700">
            <MacroRow
              label="Calories"
              value={meal.nutrition.macros.calories}
              unit="kcal"
              colorClass="bg-amber-500"
            />
            <MacroRow
              label="Protein"
              value={meal.nutrition.macros.protein}
              unit="g"
              colorClass="bg-emerald-500"
            />
            <MacroRow
              label="Carbs"
              value={meal.nutrition.macros.carbs}
              unit="g"
              colorClass="bg-blue-500"
            />
            <MacroRow
              label="Fat"
              value={meal.nutrition.macros.fat}
              unit="g"
              colorClass="bg-violet-500"
            />
          </View>
        </View>

        {/* Assumptions Section */}
        {meal.nutrition.assumptions.length > 0 && (
          <View className="mb-6">
            <ThemedText type="subtitle" className="mb-3">
              AI Assumptions
            </ThemedText>
            <ThemedText className="text-sm text-stone-500 dark:text-stone-400 mb-3">
              The following values were estimated by AI based on your input:
            </ThemedText>
            {meal.nutrition.assumptions.map((assumption, index) => (
              <AssumptionCard
                key={assumption.id || index}
                assumption={assumption}
              />
            ))}
          </View>
        )}
      </KeyboardAwareScrollView>

      <KeyboardStickyView offset={{ closed: 0, opened: 0 }}>
        <InputBar onSubmit={handleSubmitCorrection} isLoading={isSubmitting}>
          <InputBar.Input placeholder="Correct this meal..." />
        </InputBar>
      </KeyboardStickyView>
    </ThemedView>
  );
}
