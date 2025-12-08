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
import { postAgentsNutrition } from "@/lib/api/default/default";
import useGlobalStore from "@/lib/store";

export default function MealDetailScreen() {
  const [isSubmitting, setIsSubmitting] = useState(false);

  const { id } = useLocalSearchParams<{ id: string }>();
  const meal = useGlobalStore((state) => state.meals.find((m) => m.id === id));

  if (!meal) {
    return (
      <ThemedView className="flex-1 items-center justify-center">
        <ThemedText>Meal not found</ThemedText>
      </ThemedView>
    );
  }

  const assumptions = meal.nutrition.assumptions;

  const handleSubmitCorrection = async (text: string) => {
    if (isSubmitting) return;

    setIsSubmitting(true);
    try {
      // TODO: Call API endpoint for AI correction
      const response = await postAgentsNutrition({
        text,
        session_id: meal.sessionId,
        user_id: "user-1",
      });
      console.log("Submitting correction:", text);
      // await submitCorrection({ text, mealId: params.id });
    } catch (error) {
      console.error("Failed to submit correction:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
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
        {assumptions.length > 0 && (
          <View className="mb-6">
            <ThemedText type="subtitle" className="mb-3">
              AI Assumptions
            </ThemedText>
            <ThemedText className="text-sm text-stone-500 dark:text-stone-400 mb-3">
              The following values were estimated by AI based on your input:
            </ThemedText>
            {assumptions.map((assumption, index) => (
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
