import { KeyboardStickyView } from "react-native-keyboard-controller";
import { useState } from "react";
import "react-native-get-random-values";
import { v7 as uuid7 } from "uuid";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import TotalMacroPanel from "@/components/TotalMacroPanel";
import MealCard from "@/components/MealCard";
import InputBar from "@/components/InputBar";
import { postAgentsNutrition } from "@/lib/api/default/default";
import { NutritionResponseBody } from "@/lib/api/conversationAPI.schemas";
import { ScrollView } from "react-native";

const defaultMeals: Array<NutritionResponseBody["analysis"]> = [
  {
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
    ],
    macros: {
      calories: 320,
      protein: 12,
      carbs: 45,
      fat: 8,
    },
  },
  {
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
    ],
    macros: {
      calories: 450,
      protein: 35,
      carbs: 20,
      fat: 18,
    },
  },
];

export default function HomeScreen() {
  const [sessionId] = useState(() => uuid7());
  const [meals, setMeals] =
    useState<Array<NutritionResponseBody["analysis"]>>(defaultMeals);

  const handleSubmitNutrition = async (text: string) => {
    try {
      const response = await postAgentsNutrition({
        text,
        session_id: sessionId,
        user_id: "user-1",
      });
      if (response.status !== 200) {
        throw new Error(`Server error: ${response.data}`);
      }
      const nutritionData = response.data.analysis;
      setMeals((prevMeals) => [nutritionData, ...prevMeals]);
    } catch (error) {
      console.error("Error submitting nutrition:", error);
    }
  };

  return (
    <ThemedView className="h-full w-full">
      <ThemedText className="px-4" type="title">
        Tuesday
      </ThemedText>

      <TotalMacroPanel />
      <ScrollView className="flex gap-4">
        <MealCard />
        {meals.map((meal, index) => (
          <MealCard key={index} mealData={meal} />
        ))}
      </ScrollView>

      <KeyboardStickyView offset={{ closed: 0, opened: 0 }}>
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
