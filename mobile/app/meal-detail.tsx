import { useState } from "react";
import { View } from "react-native";
import { useLocalSearchParams } from "expo-router";
import {
  KeyboardAwareScrollView,
  KeyboardStickyView,
} from "react-native-keyboard-controller";
import InputBar from "@/components/InputBar";
import { MacroDisplay } from "@/components/MacroDisplay";
import { AssumptionCard } from "@/components/AssumptionCard";
import { Card } from "@/components/ui/Card";
import { SectionHeader } from "@/components/ui/SectionHeader";
import { ScreenLayout } from "@/components/ui/ScreenLayout";
import { Text } from "@/components/ui/Text";
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
    <ScreenLayout className="items-center justify-center">
      <Text variant="body">Meal not found</Text>
    </ScreenLayout>
  ) : (
    <View className="flex-1 bg-stone-50 dark:bg-stone-950">
      <KeyboardAwareScrollView className="flex-1 px-4">
        <View className="pt-4 pb-6">
          <Text variant="h1">
            {meal.nutrition.name || "Meal Details"}
          </Text>
        </View>

        {/* Macros Section */}
        <View className="mb-6">
          <SectionHeader title="Nutrition Facts" />
          <Card>
            <MacroDisplay
              variant="row"
              label="Calories"
              value={meal.nutrition.macros.calories}
              unit="kcal"
              colorClass="bg-amber-500"
            />
            <MacroDisplay
              variant="row"
              label="Protein"
              value={meal.nutrition.macros.protein}
              unit="g"
              colorClass="bg-emerald-500"
            />
            <MacroDisplay
              variant="row"
              label="Carbs"
              value={meal.nutrition.macros.carbs}
              unit="g"
              colorClass="bg-blue-500"
            />
            <MacroDisplay
              variant="row"
              label="Fat"
              value={meal.nutrition.macros.fat}
              unit="g"
              colorClass="bg-violet-500"
            />
          </Card>
        </View>

        {/* Assumptions Section */}
        {meal.nutrition.assumptions.length > 0 && (
          <View className="mb-6">
            <SectionHeader title="AI Assumptions" />
            <Text variant="caption" className="mb-3">
              The following values were estimated by AI based on your input:
            </Text>
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
    </View>
  );
}
