import { useState } from "react";
import { View, ScrollView } from "react-native";
import { useLocalSearchParams } from "expo-router";
import {
  KeyboardAwareScrollView,
  KeyboardStickyView,
} from "react-native-keyboard-controller";
import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import { Assumption } from "@/lib/api/conversationAPI.schemas";
import InputBar from "@/components/InputBar";
import { MacroRow } from "@/components/MacroRow";
import { AssumptionCard } from "@/components/AssumptionCard";

export default function MealDetailScreen() {
  const [isSubmitting, setIsSubmitting] = useState(false);

  const params = useLocalSearchParams<{
    title: string;
    description: string;
    calories: string;
    protein: string;
    carbs: string;
    fat: string;
    assumptions: string;
  }>();

  const assumptions: Assumption[] = params.assumptions
    ? JSON.parse(params.assumptions)
    : [];

  const handleSubmitCorrection = async (text: string) => {
    if (isSubmitting) return;

    setIsSubmitting(true);
    try {
      // TODO: Call API endpoint for AI correction
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
          <ThemedText type="title">{params.title || "Meal Details"}</ThemedText>
          {params.description && (
            <ThemedText className="text-stone-500 dark:text-stone-400 mt-1">
              {params.description}
            </ThemedText>
          )}
        </View>

        {/* Macros Section */}
        <View className="mb-6">
          <ThemedText type="subtitle" className="mb-3">
            Nutrition Facts
          </ThemedText>

          <View className="bg-white dark:bg-stone-900 rounded-xl p-4 border border-stone-200 dark:border-stone-700">
            <MacroRow
              label="Calories"
              value={parseFloat(params.calories || "0")}
              unit="kcal"
              colorClass="bg-amber-500"
            />
            <MacroRow
              label="Protein"
              value={parseFloat(params.protein || "0")}
              unit="g"
              colorClass="bg-emerald-500"
            />
            <MacroRow
              label="Carbs"
              value={parseFloat(params.carbs || "0")}
              unit="g"
              colorClass="bg-blue-500"
            />
            <MacroRow
              label="Fat"
              value={parseFloat(params.fat || "0")}
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
